// skill-budget prices the mentor skill definition in tokens.
//
// The always-loaded skill files are paid for on every invocation, and eval
// fixes tend to grow them a sentence at a time — growth must be a reviewed
// number, not a feeling. This tool counts each file with the API's free
// /v1/messages/count_tokens endpoint (exact, model-specific — never a
// tiktoken-style estimate) and renders a markdown budget table. With
// -base <git-ref> it also counts the same files at that ref and reports
// deltas; the CI workflow posts the result as a PR comment so every change
// to the skill definition carries a visible price tag.
//
// Auth: ANTHROPIC_API_KEY (x-api-key) or CLAUDE_CODE_OAUTH_TOKEN (Bearer +
// the oauth beta header). When both are set the API key is tried first —
// count_tokens is free, so even an out-of-credit key may serve it — and any
// non-200 falls back to the OAuth token once.
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
	model := flag.String("model", "claude-opus-4-8", "model whose tokenizer counts (token counts are model-specific)")
	baseRef := flag.String("base", "", "git ref to diff against; omit for absolute counts only")
	out := flag.String("out", "", "write the markdown here instead of stdout")
	flag.Parse()

	root, err := findRoot()
	if err != nil {
		fatal(err)
	}
	c, err := newCounter(*model)
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
			base = 0 // file absent at base counts as new
			if old := gitShow(root, *baseRef, f); old != "" {
				if base, err = c.count(old); err != nil {
					fatal(fmt.Errorf("counting %s at %s: %w", f, *baseRef, err))
				}
			}
		}
		rows = append(rows, row{name: f, base: base, head: head})
	}

	md := render(rows, *model, *baseRef)
	if *out == "" {
		fmt.Print(md)
		return
	}
	if err := os.WriteFile(*out, []byte(md), 0o644); err != nil {
		fatal(err)
	}
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

// counter calls count_tokens with whichever credential works. The API key is
// tried first (the endpoint is free, so a billing-dead key may still serve
// it); any non-200 switches to the OAuth token once, permanently.
type counter struct {
	model    string
	client   *http.Client
	apiKey   string
	oauth    string
	useOAuth bool
}

func newCounter(model string) (*counter, error) {
	c := &counter{
		model:  model,
		client: &http.Client{Timeout: 60 * time.Second},
		apiKey: os.Getenv("ANTHROPIC_API_KEY"),
		oauth:  os.Getenv("CLAUDE_CODE_OAUTH_TOKEN"),
	}
	if c.apiKey == "" && c.oauth == "" {
		return nil, errors.New("set ANTHROPIC_API_KEY or CLAUDE_CODE_OAUTH_TOKEN")
	}
	c.useOAuth = c.apiKey == ""
	return c, nil
}

func (c *counter) count(text string) (int, error) {
	payload, err := json.Marshal(map[string]any{
		"model":    c.model,
		"messages": []map[string]string{{"role": "user", "content": text}},
	})
	if err != nil {
		return 0, err
	}
	for {
		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
		if err != nil {
			return 0, err
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("anthropic-version", "2023-06-01")
		if c.useOAuth {
			req.Header.Set("Authorization", "Bearer "+c.oauth)
			req.Header.Set("anthropic-beta", "oauth-2025-04-20")
		} else {
			req.Header.Set("x-api-key", c.apiKey)
		}
		resp, err := c.client.Do(req)
		if err != nil {
			return 0, err
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return parseCount(body)
		}
		if !c.useOAuth && c.oauth != "" {
			c.useOAuth = true // API key rejected; fall back to the subscription token
			continue
		}
		return 0, fmt.Errorf("count_tokens: HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
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

// gitShow returns the file's content at ref, or "" when it doesn't exist
// there (a new file has no base to diff against).
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
