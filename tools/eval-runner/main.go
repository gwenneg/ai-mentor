// Headless runner for the ai-mentor eval suite (Groups A-C in evals/cases.md).
//
// Parses the case tables, runs each case through the claude CLI in an
// isolated HOME with a controlled ~/.ai-mentor/profile.md fixture, scores
// every response with an LLM judge, and writes a markdown report. Group D
// is interactive-only (see evals/README.md) and never runs here.
//
// Exits 0 on success, 1 when -gate is set and any case fails or errors,
// 2 on a fatal setup problem. Stdlib only.
//
// Usage: go -C tools/eval-runner run . -repo ../.. -groups A,B,C -gate
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"
)

const (
	mentorCmd   = "/ai-mentor:mentor"
	maxRawLines = 60

	vPass  = "PASS"
	vFail  = "FAIL"
	vError = "ERROR"
)

var (
	reCaseRow  = regexp.MustCompile(`^\| ([ABC]\d{2}) \|`)
	profileRel = filepath.Join(".ai-mentor", "profile.md")
)

// evalCase is one row from a cases.md group table. Statement holds the
// problem statement for Group A and the fixture/setup description otherwise.
type evalCase struct {
	Group, ID, Statement, Expected, Notes string
}

// check is one named judge check; verdict is the judge's full reply.
type check struct {
	Name   string `json:"name"`
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
}

type verdict struct {
	Pass   bool
	Checks []check
}

// result is the final outcome for one case.
type result struct {
	c        evalCase
	verdict  string // vPass, vFail, or vError
	reason   string
	response string
}

func errResult(c evalCase, err error) result {
	return result{c: c, verdict: vError, reason: err.Error()}
}

// cells splits a Markdown table row on '|' and trims each cell.
func cells(l string) []string {
	cs := strings.Split(l, "|")
	for i, c := range cs {
		cs[i] = strings.TrimSpace(c)
	}
	return cs
}

// parseCases extracts the Group A, B, and C case tables from cases.md, plus
// the Group A output-shape expectations block verbatim.
func parseCases(text string) (map[string][]evalCase, string, error) {
	cases := map[string][]evalCase{}
	var shape []string
	group, inShape := "", false
	for _, l := range strings.Split(text, "\n") {
		switch {
		case strings.HasPrefix(l, "### Group A output-shape expectations"):
			group, inShape = "", true
			continue
		case strings.HasPrefix(l, "## Group ") && len(l) > len("## Group "):
			group, inShape = "", false
			if g := string(l[len("## Group ")]); strings.Contains("ABC", g) {
				group = g
			}
			continue
		case strings.HasPrefix(l, "#"):
			group, inShape = "", false
			continue
		}
		if inShape {
			shape = append(shape, l)
			continue
		}
		m := reCaseRow.FindStringSubmatch(l)
		if group == "" || m == nil {
			continue
		}
		cs := cells(l)
		if len(cs) < 4 {
			return nil, "", fmt.Errorf("case row %s has too few columns", m[1])
		}
		c := evalCase{Group: group, ID: m[1], Expected: cs[3]}
		if group == "A" {
			c.Statement = strings.Trim(cs[2], "`")
			if len(cs) > 4 {
				c.Notes = cs[4]
			}
		} else {
			c.Statement = cs[2]
		}
		cases[group] = append(cases[group], c)
	}
	return cases, strings.TrimSpace(strings.Join(shape, "\n")), nil
}

// selectCases returns the requested cases in table order. A requested group
// that parsed to zero cases is fatal — format drift must be loud, never a
// silently green run.
func selectCases(all map[string][]evalCase, groups, ids []string) ([]evalCase, error) {
	var out []evalCase
	for _, g := range groups {
		gc := all[g]
		if len(gc) == 0 {
			return nil, fmt.Errorf("group %s parsed to zero cases — cases.md format drift?", g)
		}
		for _, c := range gc {
			if len(ids) == 0 || slices.Contains(ids, c.ID) {
				out = append(out, c)
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no cases match -cases %s in the requested groups", strings.Join(ids, ","))
	}
	return out, nil
}

// statementsByID maps Group A case IDs to their problem statements, which
// the Group C cases reuse.
func statementsByID(as []evalCase) map[string]string {
	m := map[string]string{}
	for _, c := range as {
		m[c.ID] = c.Statement
	}
	return m
}

// approachNames enumerates approach basenames for the B06 all-adopted profile.
func approachNames(repo string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(repo, "skills", "mentor", "approaches", "*.md"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no approach files under %s/skills/mentor/approaches", repo)
	}
	names := make([]string, len(files))
	for i, f := range files {
		names[i] = strings.TrimSuffix(filepath.Base(f), ".md")
	}
	return names, nil
}

// runClaude is the seam between the runner and the claude CLI; tests stub it.
var runClaude = func(ctx context.Context, dir string, env []string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "claude", args...)
	cmd.Dir = dir
	cmd.Env = env
	var out, errOut bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errOut
	if err := cmd.Run(); err != nil {
		// stderr is often empty on failure — the CLI reports errors as a
		// JSON envelope on stdout — so surface both, capped.
		detail := strings.TrimSpace(errOut.String())
		if o := strings.TrimSpace(out.String()); o != "" {
			detail += " stdout: " + o
		}
		if len(detail) > 1000 {
			detail = detail[:1000] + " ..."
		}
		return "", fmt.Errorf("claude: %w: %s", err, detail)
	}
	return out.String(), nil
}

// runner holds everything shared across cases.
type runner struct {
	repo, fixture string // absolute paths
	judge         string
	timeout       time.Duration
	shape         string            // Group A output-shape expectations, verbatim
	statements    map[string]string // Group A ID -> problem statement
	approaches    []string          // approach basenames, for the B06 fixture
	today         string            // YYYY-MM-DD
}

// runCase executes one case in a fresh temp HOME and returns its judged result.
func (r *runner) runCase(c evalCase) result {
	home, err := os.MkdirTemp("", "eval-home-")
	if err != nil {
		return errResult(c, err)
	}
	defer os.RemoveAll(home)
	// macOS temp dirs live behind a /var -> /private/var symlink; resolve so
	// HOME matches the paths the CLI's file tools report and match rules on.
	if resolved, rerr := filepath.EvalSymlinks(home); rerr == nil {
		home = resolved
	}
	env, err := caseEnv(home)
	if err != nil {
		return errResult(c, err)
	}
	workdir := r.fixture
	if c.ID == "B04" {
		if workdir, err = r.hookedFixture(); err != nil {
			return errResult(c, err)
		}
		defer os.RemoveAll(workdir)
	}
	if err := r.setupProfile(c, home); err != nil {
		return errResult(c, err)
	}
	prompts, err := r.prompts(c)
	if err != nil {
		return errResult(c, err)
	}
	var responses []string
	for _, p := range prompts {
		resp, err := r.invoke(p, workdir, env)
		if err != nil {
			return errResult(c, err)
		}
		responses = append(responses, resp)
	}
	profile := readFile(filepath.Join(home, profileRel))
	return r.judgeCase(c, responses, profile)
}

// caseEnv builds the child environment: the parent env with HOME pointed at
// the isolated temp dir. When no ANTHROPIC_API_KEY is present (local runs),
// the developer's credential is copied in so auth still works; in CI the
// API key env var passing through is the whole auth story.
func caseEnv(home string) ([]string, error) {
	env := slices.DeleteFunc(os.Environ(), func(kv string) bool {
		return strings.HasPrefix(kv, "HOME=")
	})
	env = append(env, "HOME="+home)
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		return env, nil
	}
	creds, err := localCredentials()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	return env, os.WriteFile(filepath.Join(dir, ".credentials.json"), creds, 0o600)
}

// localCredentials finds the developer's Claude Code credential for local
// runs: ~/.claude/.credentials.json where the CLI stores it as a file
// (Linux), or the login keychain on macOS, where no file exists.
func localCredentials() ([]byte, error) {
	realHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	creds, err := os.ReadFile(filepath.Join(realHome, ".claude", ".credentials.json"))
	if err == nil {
		return creds, nil
	}
	if runtime.GOOS == "darwin" {
		out, kerr := exec.Command("security", "find-generic-password",
			"-s", "Claude Code-credentials", "-w").Output()
		if kerr == nil {
			return bytes.TrimSpace(out), nil
		}
	}
	return nil, fmt.Errorf("no ANTHROPIC_API_KEY, no credentials file, no macOS keychain entry: %w", err)
}

// hookedFixture copies the fixture project to a temp dir and adds a
// .claude/settings.json with hooks, so B04 can observe hooks-as-workflow as
// a setup signal without touching the shared fixture.
func (r *runner) hookedFixture() (string, error) {
	dir, err := os.MkdirTemp("", "eval-fixture-")
	if err != nil {
		return "", err
	}
	fail := func(err error) (string, error) {
		os.RemoveAll(dir)
		return "", err
	}
	if err := os.CopyFS(dir, os.DirFS(r.fixture)); err != nil {
		return fail(err)
	}
	settings := filepath.Join(dir, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(settings), 0o755); err != nil {
		return fail(err)
	}
	hooks := `{"hooks":{"PostToolUse":[{"matcher":"Edit|Write","hooks":[{"type":"command","command":"npm test"}]}]}}` + "\n"
	if err := os.WriteFile(settings, []byte(hooks), 0o644); err != nil {
		return fail(err)
	}
	return dir, nil
}

// profileMD renders a profile fixture per references/profile-schema.md.
// anchor may be "" to omit the what's-new line.
func (r *runner) profileMD(anchor string, rows ...string) string {
	var b strings.Builder
	b.WriteString("# Mentor Profile\n*Updated: " + r.today + "*\n\n")
	b.WriteString("Level: comfortable — eval fixture\n")
	if anchor != "" {
		b.WriteString("Last new-capability check: " + anchor + "\n")
	}
	b.WriteString("\n| Capability | Status | Date | Note |\n|------------|--------|------|------|\n")
	for _, row := range rows {
		b.WriteString(row + "\n")
	}
	return b.String()
}

func profileRow(name, status, date, note string) string {
	return "| " + name + " | " + status + " | " + date + " | " + note + " |"
}

// setupProfile writes the per-case ~/.ai-mentor/profile.md fixture. Cases
// without an entry here (B01, all of Group A, C02, C03) start profile-less.
func (r *runner) setupProfile(c evalCase, home string) error {
	past := time.Now().AddDate(0, 0, -21).Format("2006-01-02")
	week := currentWeek()
	var content string
	switch c.ID {
	case "B02":
		content = r.profileMD(week, profileRow("autonomous-loops", "shown", past, "Demoed /loop on a flaky retry test"))
	case "B03":
		content = r.profileMD(week, profileRow("fan-out-workflows", "declined", past, `"too token-heavy"`))
	case "B04":
		content = r.profileMD(week) // empty profile; the hooks live in the fixture copy
	case "B05":
		content = r.profileMD("2026-w20")
	case "B06":
		rows := make([]string, len(r.approaches))
		for i, a := range r.approaches {
			rows[i] = profileRow(a, "adopted", past, "eval fixture")
		}
		content = r.profileMD(week, rows...)
	case "C01":
		content = r.profileMD(week, profileRow("plan-mode", "adopted", past, "uses plan mode daily"))
	default:
		return nil
	}
	path := filepath.Join(home, profileRel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// prompts returns the mentor invocations for a case, in order. Group B is
// always the bare growth-mode invocation; Group C reuses Group A statements
// (C02 runs its statement twice against the same HOME).
func (r *runner) prompts(c evalCase) ([]string, error) {
	stmt := func(id string) (string, error) {
		s, ok := r.statements[id]
		if !ok {
			return "", fmt.Errorf("case %s needs Group A case %s, which did not parse", c.ID, id)
		}
		return mentorCmd + " " + s, nil
	}
	switch {
	case c.Group == "A":
		return []string{mentorCmd + " " + c.Statement}, nil
	case c.Group == "B":
		return []string{mentorCmd}, nil
	case c.ID == "C01", c.ID == "C03":
		p, err := stmt("A01")
		return []string{p}, err
	case c.ID == "C02":
		p, err := stmt("A03")
		return []string{p, p}, err
	}
	return nil, fmt.Errorf("case %s: no prompt rule", c.ID)
}

// invoke runs one mentor invocation and extracts the JSON "result" field.
func (r *runner) invoke(prompt, dir string, env []string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	out, err := runClaude(ctx, dir, env,
		"-p", prompt, "--plugin-dir", r.repo, "--output-format", "json", "--max-turns", "30")
	if err != nil {
		return "", err
	}
	return resultField(out)
}

// resultField loosely decodes claude's --output-format json envelope and
// returns its "result" string.
func resultField(out string) (string, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		return "", fmt.Errorf("claude output is not JSON: %w", err)
	}
	s, ok := m["result"].(string)
	if !ok {
		return "", fmt.Errorf(`claude output has no string "result" field`)
	}
	return s, nil
}

// judgeCase scores the responses with the judge model. An unparseable judge
// reply is an ERROR verdict, reported distinctly from FAIL.
func (r *runner) judgeCase(c evalCase, responses []string, profile string) result {
	res := result{c: c, response: strings.Join(responses, "\n\n--- second run ---\n\n")}
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	out, err := runClaude(ctx, ".", os.Environ(),
		"-p", r.judgePrompt(c, responses, profile), "--model", r.judge, "--max-turns", "5")
	if err != nil {
		res.verdict, res.reason = vError, err.Error()
		return res
	}
	v, err := parseVerdict(out)
	if err != nil {
		res.verdict, res.reason = vError, "judge reply not parseable: "+err.Error()
		return res
	}
	if v.Pass {
		res.verdict = vPass
	} else {
		res.verdict, res.reason = vFail, failReason(v)
	}
	return res
}

// judgePrompt builds the strict scoring prompt: the case expectation, the
// response(s), the after-run profile, and (for Group A) the output-shape
// expectations from cases.md verbatim.
func (r *runner) judgePrompt(c evalCase, responses []string, profile string) string {
	var b strings.Builder
	b.WriteString("You are a strict evaluator for the ai-mentor Claude Code skill. ")
	b.WriteString("Judge whether the response(s) below meet the case expectation. ")
	b.WriteString("Reply with STRICT JSON only — no prose, no markdown fences: ")
	b.WriteString(`{"pass": bool, "checks": [{"name": string, "pass": bool, "reason": string}]}`)
	fmt.Fprintf(&b, "\n\nToday's date: %s\nCase %s (Group %s).\n", r.today, c.ID, c.Group)
	if c.Group == "A" {
		fmt.Fprintf(&b, "Problem statement: %s\nExpected goal classification: %s\n", c.Statement, c.Expected)
		if c.Notes != "" {
			fmt.Fprintf(&b, "Case notes: %s\n", c.Notes)
		}
		b.WriteString("\nGroup A output-shape expectations (verbatim from cases.md; every classified case must satisfy all of them):\n")
		b.WriteString(r.shape + "\n")
	} else {
		fmt.Fprintf(&b, "Setup / profile fixture: %s\nExpected behavior: %s\n", c.Statement, c.Expected)
	}
	for i, resp := range responses {
		label := "Response"
		if len(responses) > 1 {
			label = fmt.Sprintf("Response from run %d", i+1)
		}
		fmt.Fprintf(&b, "\n%s:\n<<<\n%s\n>>>\n", label, resp)
	}
	if profile != "" {
		fmt.Fprintf(&b, "\nProfile after the run (~/.ai-mentor/profile.md):\n<<<\n%s\n>>>\n", profile)
	} else {
		b.WriteString("\nNo profile file existed after the run.\n")
	}
	return b.String()
}

// parseVerdict extracts the judge's JSON leniently: everything from the
// first '{' to the last '}' must decode to an object with a "pass" bool.
func parseVerdict(s string) (verdict, error) {
	i, j := strings.Index(s, "{"), strings.LastIndex(s, "}")
	if i < 0 || j <= i {
		return verdict{}, fmt.Errorf("no JSON object found")
	}
	var raw struct {
		Pass   *bool   `json:"pass"`
		Checks []check `json:"checks"`
	}
	if err := json.Unmarshal([]byte(s[i:j+1]), &raw); err != nil {
		return verdict{}, err
	}
	if raw.Pass == nil {
		return verdict{}, fmt.Errorf(`judge JSON has no "pass" field`)
	}
	return verdict{Pass: *raw.Pass, Checks: raw.Checks}, nil
}

// failReason picks the first failing check's reason as the one-line summary.
func failReason(v verdict) string {
	for _, ch := range v.Checks {
		if !ch.Pass {
			return ch.Name + ": " + ch.Reason
		}
	}
	return "judge returned pass=false without a failing check"
}

// groupsIn returns the group letters in first-seen order — order-preserving
// on purpose, so the report follows the case order.
func groupsIn(results []result) []string {
	var out []string
	for _, r := range results {
		if !slices.Contains(out, r.c.Group) {
			out = append(out, r.c.Group)
		}
	}
	return out
}

// summary renders the one-line per-group tallies used in the report and on
// stdout.
func summary(results []result) string {
	var parts []string
	for _, g := range groupsIn(results) {
		var pass, fail, errs int
		for _, r := range results {
			switch {
			case r.c.Group != g:
			case r.verdict == vPass:
				pass++
			case r.verdict == vFail:
				fail++
			default:
				errs++
			}
		}
		parts = append(parts, fmt.Sprintf("Group %s: %d pass / %d fail / %d error", g, pass, fail, errs))
	}
	return strings.Join(parts, " | ")
}

// oneLine flattens a reason so it fits a markdown table cell.
func oneLine(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.ReplaceAll(s, "|", "/")
}

// truncateLines caps s at n lines, appending a marker when cut.
func truncateLines(s string, n int) string {
	ls := strings.Split(s, "\n")
	if len(ls) <= n {
		return s
	}
	return strings.Join(ls[:n], "\n") + fmt.Sprintf("\n... (%d more lines truncated)", len(ls)-n)
}

// renderReport writes the full markdown report: summary, one table per
// group, then the raw responses of every non-passing case.
func renderReport(results []result) string {
	var b strings.Builder
	b.WriteString("# ai-mentor eval report\n\n")
	b.WriteString(summary(results) + "\n")
	for _, g := range groupsIn(results) {
		fmt.Fprintf(&b, "\n## Group %s\n\n| Case | Verdict | Reason |\n|------|---------|--------|\n", g)
		for _, r := range results {
			if r.c.Group == g {
				fmt.Fprintf(&b, "| %s | %s | %s |\n", r.c.ID, r.verdict, oneLine(r.reason))
			}
		}
	}
	wroteHeader := false
	for _, r := range results {
		if r.verdict == vPass || r.response == "" {
			continue
		}
		if !wroteHeader {
			b.WriteString("\n## Raw failures\n")
			wroteHeader = true
		}
		fmt.Fprintf(&b, "\n### %s (%s)\n\n```\n%s\n```\n", r.c.ID, r.verdict, truncateLines(r.response, maxRawLines))
	}
	return b.String()
}

// currentWeek returns the ISO week slug (e.g. 2026-w28) used as the
// what's-new anchor in profile fixtures.
func currentWeek() string {
	y, w := time.Now().ISOWeek()
	return fmt.Sprintf("%d-w%02d", y, w)
}

// splitList splits a comma list, trimming blanks; empty input yields nil.
func splitList(s string) []string {
	var out []string
	for _, x := range strings.Split(s, ",") {
		if x = strings.TrimSpace(x); x != "" {
			out = append(out, x)
		}
	}
	return out
}

// readFile returns the file's content, or "" when it doesn't exist.
func readFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
	os.Exit(2)
}

func main() {
	groups := flag.String("groups", "A,B,C", "comma-separated groups to run")
	ids := flag.String("cases", "", "comma-separated case IDs to run (default: all in the groups)")
	repo := flag.String("repo", "", "plugin repo root (default: walk up to the directory containing skills/mentor)")
	fixture := flag.String("fixture", "", "fixture project dir (default <repo>/evals/fixture)")
	out := flag.String("out", "eval-report.md", "markdown report path")
	gate := flag.Bool("gate", false, "exit 1 when any case fails or errors")
	judge := flag.String("model-judge", "claude-sonnet-5", "judge model for scoring")
	timeout := flag.Int("timeout", 300, "per-case timeout in seconds")
	flag.Parse()

	repoAbs, err := filepath.Abs(*repo)
	if *repo == "" {
		repoAbs, err = findRoot(".")
	}
	if err != nil {
		fatal(err)
	}
	fix := *fixture
	if fix == "" {
		fix = filepath.Join(repoAbs, "evals", "fixture")
	}
	if fix, err = filepath.Abs(fix); err != nil {
		fatal(err)
	}
	if _, err := os.Stat(fix); err != nil {
		fatal(fmt.Errorf("fixture dir: %w", err))
	}

	text, err := os.ReadFile(filepath.Join(repoAbs, "evals", "cases.md"))
	if err != nil {
		fatal(err)
	}
	all, shape, err := parseCases(string(text))
	if err != nil {
		fatal(err)
	}
	selected, err := selectCases(all, splitList(*groups), splitList(*ids))
	if err != nil {
		fatal(err)
	}
	approaches, err := approachNames(repoAbs)
	if err != nil {
		fatal(err)
	}

	r := &runner{
		repo: repoAbs, fixture: fix, judge: *judge,
		timeout:    time.Duration(*timeout) * time.Second,
		shape:      shape,
		statements: statementsByID(all["A"]),
		approaches: approaches,
		today:      time.Now().Format("2006-01-02"),
	}
	var results []result
	for _, c := range selected {
		fmt.Printf("running %s ...", c.ID)
		res := r.runCase(c)
		fmt.Printf(" %s\n", res.verdict)
		if res.verdict != vPass {
			fmt.Printf("  reason: %s\n", res.reason)
		}
		results = append(results, res)
	}

	if err := os.WriteFile(*out, []byte(renderReport(results)), 0o644); err != nil {
		fatal(err)
	}
	fmt.Println(summary(results))
	if *gate && slices.ContainsFunc(results, func(r result) bool { return r.verdict != vPass }) {
		os.Exit(1)
	}
}

// findRoot walks upward from dir to the first directory containing
// skills/mentor, so the runner works from anywhere in the repo — including
// tools/eval-runner itself, where `go -C tools/eval-runner run .` lands.
// Keep in sync with the copies in tools/structural-audit/main.go and
// tools/catalog-drift/main.go.
func findRoot(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "skills", "mentor")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no skills/mentor directory found here or above")
		}
		dir = parent
	}
}
