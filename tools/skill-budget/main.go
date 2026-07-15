// skill-budget prices the mentor skill definition in tokens.
//
// The always-loaded skill files are paid for on every invocation, and eval
// fixes tend to grow them a sentence at a time — growth must be a reviewed
// number, not a feeling. This tool counts each file with the API's free
// /v1/messages/count_tokens endpoint (exact, model-specific — never a
// tiktoken-style estimate) and prints a markdown budget table to stdout.
// With -base <git-ref> it also counts the same files at that ref and
// reports deltas; the CI workflow posts the result as a PR comment so every
// change to the skill definition carries a visible price tag.
//
// Auth: CLAUDE_CODE_OAUTH_TOKEN (Bearer + the oauth beta header) wins when
// set, else ANTHROPIC_API_KEY (x-api-key) — the same precedence evals.yml
// enforces. One credential, one code path: a rejected credential fails
// loudly rather than silently switching.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const endpoint = "https://api.anthropic.com/v1/messages/count_tokens"

// countModel is fixed on purpose: budget deltas are only meaningful when
// every measurement uses the same tokenizer. When this model deprecates,
// change the constant — the first run after the change re-baselines, which
// is visible in that PR's comment.
const countModel = "claude-opus-4-8"

// The always-loaded skill surface. SKILL.md loads on every invocation; the
// two mode files are the fork — exactly one of them is read per run, which
// is why the report prices the two invocation paths separately.
var budgetFiles = []string{
	"skills/mentor/SKILL.md",
	"skills/mentor/problem-mode.md",
	"skills/mentor/growth-mode.md",
}

type row struct {
	name string
	base int // -1 when no base ref was given
	head int
}

func main() {
	baseRef := flag.String("base", "", "git ref to diff against; omit for absolute counts only")
	flag.Parse()

	root, err := findRoot()
	if err != nil {
		fatal(err)
	}
	if *baseRef != "" {
		// A typo'd ref must fail loudly here — gitShow treats per-file errors
		// as "file absent at base", which would silently zero every baseline.
		if err := verifyRef(root, *baseRef); err != nil {
			fatal(err)
		}
	}
	c, err := newCounter()
	if err != nil {
		fatal(err)
	}

	var rows []row
	for _, f := range budgetFiles {
		content, err := os.ReadFile(filepath.Join(root, f))
		if err != nil {
			fatal(fmt.Errorf("reading %s: %w", f, err))
		}
		head, err := c.count(string(content))
		if err != nil {
			fatal(fmt.Errorf("counting %s: %w", f, err))
		}
		base := -1
		if *baseRef != "" {
			base = 0 // ref verified above, so absence means the file is new
			if old := gitShow(root, *baseRef, f); old != "" {
				if base, err = c.count(old); err != nil {
					fatal(fmt.Errorf("counting %s at %s: %w", f, *baseRef, err))
				}
			}
		}
		rows = append(rows, row{name: f, base: base, head: head})
	}

	fmt.Print(render(rows, countModel, *baseRef))
}

// render produces the markdown budget report. With no base ref the delta
// column is omitted. The two invocation paths are priced because a real run
// loads SKILL.md plus exactly one mode file, never all three.
func render(rows []row, model, baseRef string) string {
	withBase := baseRef != ""
	var b strings.Builder
	fmt.Fprintf(&b, "### Skill-def token budget (`%s` tokenizer)\n\n", model)
	if withBase {
		fmt.Fprintf(&b, "| File | %s | this PR | Δ |\n|---|--:|--:|--:|\n", baseRef)
	} else {
		b.WriteString("| File | tokens |\n|---|--:|\n")
	}
	var baseTotal, headTotal int
	byName := map[string]row{}
	for _, r := range rows {
		byName[r.name] = r
		baseTotal += r.base
		headTotal += r.head
		if withBase {
			fmt.Fprintf(&b, "| %s | %d | %d | %s |\n", r.name, r.base, r.head, delta(r.head-r.base))
		} else {
			fmt.Fprintf(&b, "| %s | %d |\n", r.name, r.head)
		}
	}
	if withBase {
		fmt.Fprintf(&b, "| **Total** | **%d** | **%d** | **%s** |\n", baseTotal, headTotal, delta(headTotal-baseTotal))
	} else {
		fmt.Fprintf(&b, "| **Total** | **%d** |\n", headTotal)
	}

	skill, problem, growth := byName[budgetFiles[0]], byName[budgetFiles[1]], byName[budgetFiles[2]]
	b.WriteString("\n")
	b.WriteString(pathLine("problem mode (SKILL.md + problem-mode.md)", skill, problem, withBase))
	b.WriteString(pathLine("growth mode (SKILL.md + growth-mode.md)", skill, growth, withBase))
	b.WriteString("\n<sub>Counted via `/v1/messages/count_tokens` (free, exact); per-file counts include a few tokens of message framing.</sub>\n")
	return b.String()
}

func pathLine(label string, a, c row, withBase bool) string {
	if !withBase {
		return fmt.Sprintf("Invocation cost — %s: **%d**\n", label, a.head+c.head)
	}
	d := (a.head + c.head) - (a.base + c.base)
	return fmt.Sprintf("Invocation cost — %s: **%d** (%s)\n", label, a.head+c.head, delta(d))
}

func delta(d int) string {
	switch {
	case d > 0:
		return fmt.Sprintf("+%d", d)
	case d < 0:
		return fmt.Sprintf("%d", d)
	default:
		return "±0"
	}
}

// counter holds the one credential chosen at startup (OAuth wins, matching
// evals.yml) and calls count_tokens with it.
type counter struct {
	client   *http.Client
	cred     string
	useOAuth bool
}

func newCounter() (*counter, error) {
	c := &counter{client: &http.Client{Timeout: 60 * time.Second}}
	if oauth := os.Getenv("CLAUDE_CODE_OAUTH_TOKEN"); oauth != "" {
		c.cred, c.useOAuth = oauth, true
	} else if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		c.cred = key
	} else {
		return nil, errors.New("set CLAUDE_CODE_OAUTH_TOKEN or ANTHROPIC_API_KEY")
	}
	return c, nil
}

func (c *counter) count(text string) (int, error) {
	payload, err := json.Marshal(map[string]any{
		"model":    countModel,
		"messages": []map[string]string{{"role": "user", "content": text}},
	})
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return 0, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")
	if c.useOAuth {
		req.Header.Set("Authorization", "Bearer "+c.cred)
		req.Header.Set("anthropic-beta", "oauth-2025-04-20")
	} else {
		req.Header.Set("x-api-key", c.cred)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("count_tokens: HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return parseCount(body)
}

func parseCount(body []byte) (int, error) {
	var r struct {
		InputTokens int `json:"input_tokens"`
	}
	if err := json.Unmarshal(body, &r); err != nil {
		return 0, fmt.Errorf("parsing count_tokens response: %w", err)
	}
	if r.InputTokens <= 0 {
		return 0, fmt.Errorf("count_tokens returned no input_tokens: %s", strings.TrimSpace(string(body)))
	}
	return r.InputTokens, nil
}

// verifyRef fails when the ref doesn't resolve to a commit, so a typo'd
// -base can never silently zero the baseline.
func verifyRef(root, ref string) error {
	if err := exec.Command("git", "-C", root, "rev-parse", "--verify", "--quiet", ref+"^{commit}").Run(); err != nil {
		return fmt.Errorf("base ref %q does not resolve to a commit", ref)
	}
	return nil
}

// gitShow returns the file's content at ref, or "" when it doesn't exist
// there — safe to conflate with errors only because main verified the ref.
func gitShow(root, ref, path string) string {
	out, err := exec.Command("git", "-C", root, "show", ref+":"+path).Output()
	if err != nil {
		return ""
	}
	return string(out)
}

// findRoot walks up from the working directory to the repo root. It is
// intentionally duplicated across tools/ — keep in sync with the copies in
// catalog-lint, catalog-drift, approaches-index, and eval-runner.
func findRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no .git found walking up from the working directory")
		}
		dir = parent
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "skill-budget:", err)
	os.Exit(1)
}
