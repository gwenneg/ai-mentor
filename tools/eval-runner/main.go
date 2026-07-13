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
// -smoke runs the curated per-change tier (see smokeCases); -epochs N runs
// every selected case N times, passing on a strict majority and flagging
// mixed results FLAKY.
//
// Usage: go -C tools/eval-runner run . -repo ../.. -groups A,B,C -gate
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"
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
	// Plugin-name extraction for the judge's ground-truth block. Keep in sync
	// with the copies in tools/catalog-lint and tools/catalog-drift.
	reRowName = regexp.MustCompile("^\\| `([a-z0-9.-]+)`")
	reTok     = regexp.MustCompile("`([a-z0-9.-]+)`")
)

// smokeCases is the curated smoke tier (-smoke): one case per behavior class,
// for cheap per-change runs — the full suite stays the release gate.
// A01 canonical classified shape · A08 misroute trap · A13 graceful decline ·
// A19 promoted-plugin displacement pin · A20 stack-match tier label + portable
// prompt · B01 first meeting + profile creation · B06 saturated ignorance
// map · C01 adopted-not-retaught · C02 same-HOME double run.
// selectCases fails loudly if any of these IDs drops out of cases.md.
var smokeCases = []string{"A01", "A08", "A13", "A19", "A20", "B01", "B06", "C01", "C02"}

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
// that parsed to zero cases is fatal, and so is a requested ID that matches
// nothing — format drift (or a stale smoke list) must be loud, never a
// silently smaller run.
func selectCases(all map[string][]evalCase, groups, ids []string) ([]evalCase, error) {
	var out []evalCase
	matched := map[string]bool{}
	for _, g := range groups {
		gc := all[g]
		if len(gc) == 0 {
			return nil, fmt.Errorf("group %s parsed to zero cases — cases.md format drift?", g)
		}
		for _, c := range gc {
			if len(ids) == 0 || slices.Contains(ids, c.ID) {
				out = append(out, c)
				matched[c.ID] = true
			}
		}
	}
	for _, id := range ids {
		if !matched[id] {
			return nil, fmt.Errorf("case %s not found in the requested groups — typo or cases.md drift?", id)
		}
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

// approachNames enumerates every teachable unit for the B06 all-adopted
// profile: one approaches/<id>.md file per capability (index.md excluded) —
// B06's "honest empty answer" only holds when the WHOLE ignorance map is
// saturated.
func approachNames(repo string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(repo, "skills", "mentor", "approaches", "*", "*.md"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no approach files under %s/skills/mentor/approaches", repo)
	}
	var names []string
	for _, f := range files {
		if n := strings.TrimSuffix(filepath.Base(f), ".md"); n != "index" {
			names = append(names, n)
		}
	}
	return names, nil
}

// groundTruth is the set of real capabilities and fixture paths inlined into
// the judge prompt so fabrication and grounding are checked against the repo,
// not the judge's own memory.
type groundTruth struct {
	fixture      []string
	plugins      []string // directory plugins ∪ promoted (the fabrication whitelist)
	promoted     []string // promoted plugin ids — first-class approaches, no tier label
	techniques   []string
	integrations []string
}

// buildGroundTruth reads the catalog and fixture once so every judge call can
// check recommendations against them. Any failure here must be fatal to the
// caller: the judge prompt frames the plugin list as COMPLETE and fails
// anything absent as a fabrication, so a silently empty or truncated list
// would mass-fail every plugin recommendation.
func buildGroundTruth(repo, fixture string) (groundTruth, error) {
	skill := filepath.Join(repo, "skills", "mentor")
	gt := groundTruth{fixture: fixtureFiles(fixture)}
	b, err := os.ReadFile(filepath.Join(skill, "marketplace.md"))
	if err != nil {
		return gt, fmt.Errorf("judge ground truth: %w", err)
	}
	gt.plugins = pluginNames(string(b))
	if len(gt.plugins) == 0 {
		return gt, fmt.Errorf("judge ground truth: zero plugin names parsed from marketplace.md — format drift?")
	}
	files, err := filepath.Glob(filepath.Join(skill, "approaches", "*", "*.md"))
	if err != nil || len(files) == 0 {
		return gt, fmt.Errorf("judge ground truth: no approach files under %s: %v", filepath.Join(skill, "approaches"), err)
	}
	for _, f := range files {
		id := strings.TrimSuffix(filepath.Base(f), ".md")
		if id == "index" {
			continue
		}
		switch approachKind(f) {
		case "integration", "doc":
			gt.integrations = append(gt.integrations, id)
		case "plugin":
			gt.plugins = append(gt.plugins, id)
			gt.promoted = append(gt.promoted, id)
		default:
			// techniques — built-in commands live inside their covering
			// technique files now, so the judge gets no separate command list.
			gt.techniques = append(gt.techniques, id)
		}
	}
	return gt, nil
}

// fixtureFiles lists the fixture repo's files as repo-relative paths.
func fixtureFiles(dir string) []string {
	var out []string
	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if rel, e := filepath.Rel(dir, p); e == nil {
			out = append(out, rel)
		}
		return nil
	})
	slices.Sort(out)
	return out
}

// approachKind returns the value of an approach file's `kind:` line, or ""
// for a technique deep-dive (which has no kind: line).
func approachKind(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, l := range strings.Split(string(b), "\n") {
		if k, ok := strings.CutPrefix(l, "kind: "); ok {
			return strings.TrimSpace(k)
		}
	}
	return ""
}

// pluginNames extracts the plugin ids the catalog declares: the first
// backticked token of each table row plus backticked tokens in the prose
// sections (Language servers, Specialty). Keep in sync with the copies in
// tools/catalog-lint/main.go and tools/catalog-drift/main.go.
func pluginNames(text string) []string {
	var names []string
	proseList := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "#") {
			h := strings.ToLower(line)
			proseList = strings.Contains(h, "language server") || strings.Contains(h, "specialty")
			continue
		}
		if m := reRowName.FindStringSubmatch(line); m != nil {
			names = append(names, m[1])
			continue
		}
		if proseList {
			for _, m := range reTok.FindAllStringSubmatch(line, -1) {
				names = append(names, m[1])
			}
		}
	}
	return names
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
	subjectModel  string // model the mentor under test runs on (pinned for gate stability)
	timeout       time.Duration
	shape         string            // Group A output-shape expectations, verbatim
	statements    map[string]string // Group A ID -> problem statement
	approaches    []string          // approach basenames, for the B06 fixture
	ground        groundTruth       // capability + fixture facts inlined into judge prompts
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
	workdir, err := r.caseFixture(c.ID == "B04")
	if err != nil {
		return errResult(c, err)
	}
	defer os.RemoveAll(workdir)
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
// the isolated temp dir. When neither ANTHROPIC_API_KEY nor
// CLAUDE_CODE_OAUTH_TOKEN is present (local runs), the developer's credential
// is copied in so auth still works; in CI either env var passing through is
// the whole auth story — the CLI honors both.
func caseEnv(home string) ([]string, error) {
	env := slices.DeleteFunc(os.Environ(), func(kv string) bool {
		return strings.HasPrefix(kv, "HOME=")
	})
	env = append(env, "HOME="+home)
	if os.Getenv("ANTHROPIC_API_KEY") != "" || os.Getenv("CLAUDE_CODE_OAUTH_TOKEN") != "" {
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
	return nil, fmt.Errorf("no ANTHROPIC_API_KEY or CLAUDE_CODE_OAUTH_TOKEN, no credentials file, no macOS keychain entry: %w", err)
}

// caseFixture copies the fixture project to a temp dir, so concurrent cases
// can never observe each other's edits. B04's copy additionally gets a
// .claude/settings.json with hooks, so hooks-as-workflow is observable as a
// setup signal without touching the shared fixture.
func (r *runner) caseFixture(withHooks bool) (string, error) {
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
	if !withHooks {
		return dir, nil
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
	case "C04":
		content = r.profileMD(week,
			profileRow("background-agents", "declined", past, `"prefer local runs"`),
			profileRow("plan-mode", "shown", past, "tried it once"))
	case "C05":
		content = r.profileMD(week, profileRow("plan-mode", "declined", past, `"too slow for me"`))
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
	case c.ID == "C01", c.ID == "C03", c.ID == "C04", c.ID == "C05":
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
		"-p", prompt, "--model", r.subjectModel, "--plugin-dir", r.repo,
		"--output-format", "stream-json", "--verbose", "--max-turns", "30")
	if err != nil {
		return "", err
	}
	return assistantText(out)
}

// assistantText extracts what the user would have seen: the concatenated
// text blocks of every assistant message in a stream-json transcript.
// Judging only the final message hides everything the model said before a
// trailing profile write ("Recorded." would be the whole response). A plain
// json envelope's "result" field is kept as a fallback for older output.
func assistantText(out string) (string, error) {
	var texts []string
	envelope := ""
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			continue
		}
		if msg, ok := m["message"].(map[string]any); ok && m["type"] == "assistant" {
			content, _ := msg["content"].([]any)
			for _, c := range content {
				cm, _ := c.(map[string]any)
				if cm["type"] == "text" {
					if s, _ := cm["text"].(string); strings.TrimSpace(s) != "" {
						texts = append(texts, s)
					}
				}
			}
		}
		if s, ok := m["result"].(string); ok {
			envelope = s
		}
	}
	if len(texts) > 0 {
		return strings.Join(texts, "\n\n"), nil
	}
	if envelope != "" {
		return envelope, nil
	}
	return "", fmt.Errorf("no assistant text found in claude output")
}

// judgeCase scores the responses with the judge model. An unparseable judge
// reply is an ERROR verdict, reported distinctly from FAIL. The judge runs
// hermetically — isolated HOME and an empty working directory — so no
// CLAUDE.md, auto memory, or repo context can leak into verdicts. (--bare
// would be simpler but silently breaks macOS keychain auth.)
func (r *runner) judgeCase(c evalCase, responses []string, profile string) result {
	res := result{c: c, response: strings.Join(responses, "\n\n--- second run ---\n\n")}
	home, err := os.MkdirTemp("", "judge-home-")
	if err != nil {
		res.verdict, res.reason = vError, err.Error()
		return res
	}
	defer os.RemoveAll(home)
	if resolved, rerr := filepath.EvalSymlinks(home); rerr == nil {
		home = resolved
	}
	env, err := caseEnv(home)
	if err != nil {
		res.verdict, res.reason = vError, err.Error()
		return res
	}
	workdir := filepath.Join(home, "empty")
	if err := os.MkdirAll(workdir, 0o755); err != nil {
		res.verdict, res.reason = vError, err.Error()
		return res
	}
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	out, err := runClaude(ctx, workdir, env,
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
	b.WriteString("You are a strict evaluator for the ai-mentor Claude Code skill. Do not use any tools — judge only the material in this prompt and reply immediately. ")
	b.WriteString("You cannot see the run's tool calls: never infer from the text whether checks actually ran, and a missing profile file does not mean reads were skipped. ")
	b.WriteString("The response was produced INSIDE the fixture repo described by the problem context — its file paths are the fixture's, not this machine's; do not judge them against any other repository. ")
	b.WriteString("Approaches/techniques (plan mode, hooks, browser integration, ...) and built-in slash commands are NOT catalog plugins — plugin rules (tier labels, ⚠️ alternatives) apply only to marketplace plugins installed via /plugin install. ")
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
		b.WriteString("\nThe case notes take precedence over the shape expectations when they conflict: a case whose notes say it is not classified (catalog browse, graceful decline) is judged on its notes, not on the classified-case shape.\n")
		b.WriteString("For a problem about the fixture repo itself, a fenced prompt that cites a code path not in the fixture file list below is fabricated grounding — fail it (unless the case notes mark the prompt portable to a different repo).\n")
	} else {
		fmt.Fprintf(&b, "Setup / profile fixture: %s\nExpected behavior: %s\n", c.Statement, c.Expected)
	}

	b.WriteString("\n--- Ground truth: judge the response's recommendations against THIS, not your own memory ---\n")
	if len(r.ground.fixture) > 0 {
		b.WriteString("Fixture repo files (the only real paths in the fixture repo): " + strings.Join(r.ground.fixture, ", ") + "\n")
	}
	b.WriteString("Real marketplace plugins (COMPLETE list; installed as `<name>@claude-plugins-official`). A recommended plugin whose name is NOT in this list is a fabrication — fail the case and name it:\n")
	b.WriteString(strings.Join(r.ground.plugins, ", ") + "\n")
	// Bulleted, not comma-joined: a judge scanning a 15-name comma run has
	// missed entries and failed cases on plugins that ARE in this list.
	b.WriteString("Of the plugins above, these are PROMOTED first-class approaches (hands-on validated, ranked in the playbooks). They carry NO tier label — expecting a label on them is an error; their record facts count as hands-on validated. Tier-label rules apply only to the remaining (directory) plugins. Check this list name by name before calling a plugin unpromoted:\n")
	for _, p := range r.ground.promoted {
		b.WriteString("- " + p + "\n")
	}
	b.WriteString("Known-real techniques: " + strings.Join(r.ground.techniques, ", ") + ". Known-real integrations: " + strings.Join(r.ground.integrations, ", ") + ".\n")
	b.WriteString("These technique/integration lists are NOT exhaustive of Claude Code, and built-in slash commands are not listed at all (e.g. /code-review, /verify, /goal, /loop, /schedule, /init, /plan, /model, /effort, --worktree, Shift+Tab are all real) — judge those against your knowledge of current Claude Code, flagging only commands or flags you are confident do not exist. The plugin list above IS complete: judge plugin recommendations strictly against it.\n")

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

// expandEpochs repeats each case n times as adjacent copies, so epoch
// results come back as consecutive chunks that aggregateEpochs can fold
// per case. Copies are fully independent runs: each gets its own HOME and
// fixture copy in runCase.
func expandEpochs(cases []evalCase, n int) []evalCase {
	if n <= 1 {
		return cases
	}
	out := make([]evalCase, 0, len(cases)*n)
	for _, c := range cases {
		for range n {
			out = append(out, c)
		}
	}
	return out
}

// aggregateEpochs folds each case's n consecutive epoch results into one
// verdict: PASS on a strict majority of passing epochs, ERROR when every
// epoch errored, FAIL otherwise. Mixed results are flagged FLAKY in the
// reason so they stay visible in the report even when the majority passes.
// Relies on expandEpochs's adjacency invariant.
func aggregateEpochs(results []result, n int) []result {
	if n <= 1 {
		return results
	}
	out := make([]result, 0, len(results)/n)
	for i := 0; i+n <= len(results); i += n {
		out = append(out, foldEpochs(results[i:i+n]))
	}
	return out
}

func foldEpochs(chunk []result) result {
	agg := result{c: chunk[0].c}
	pass, fail, bad := 0, 0, -1
	for i, r := range chunk {
		switch r.verdict {
		case vPass:
			pass++
		case vFail:
			fail++
		}
		if r.verdict != vPass && bad < 0 {
			bad = i
		}
	}
	n := len(chunk)
	switch {
	case pass*2 > n:
		agg.verdict = vPass
	case pass == 0 && fail == 0:
		agg.verdict = vError
	default:
		agg.verdict = vFail
	}
	if bad >= 0 {
		agg.reason = fmt.Sprintf("epoch %d: %s", bad+1, chunk[bad].reason)
		agg.response = chunk[bad].response
	}
	switch {
	case pass > 0 && pass < n:
		agg.reason = fmt.Sprintf("FLAKY %d/%d epochs passed — %s", pass, n, agg.reason)
	case agg.verdict != vPass:
		agg.reason = fmt.Sprintf("%d/%d epochs passed — %s", pass, n, agg.reason)
	}
	return agg
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
	smoke := flag.Bool("smoke", false, "run the curated smoke tier (one case per behavior class) — the cheap per-change signal; the full suite stays the release gate")
	epochs := flag.Int("epochs", 1, "independent runs per case; with N>1 a case passes on a strict majority of epochs and mixed results are flagged FLAKY")
	jobs := flag.Int("j", 3, "cases to run concurrently (keep modest: every case is a subject run plus a judge run against the same account)")
	judge := flag.String("model-judge", "claude-sonnet-5", "judge model for scoring")
	modelSubject := flag.String("model-subject", "claude-sonnet-5", "model the mentor under test runs on (pinned so a gate red is a regression, not CLI-default drift)")
	timeout := flag.Int("timeout", 300, "per-case timeout in seconds")
	flag.Parse()

	if *epochs < 1 {
		fatal(fmt.Errorf("-epochs must be >= 1"))
	}
	idList := splitList(*ids)
	if *smoke {
		if len(idList) > 0 {
			fatal(fmt.Errorf("-smoke and -cases are mutually exclusive"))
		}
		idList = smokeCases
	}

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
	selected, err := selectCases(all, splitList(*groups), idList)
	if err != nil {
		fatal(err)
	}
	approaches, err := approachNames(repoAbs)
	if err != nil {
		fatal(err)
	}

	ground, err := buildGroundTruth(repoAbs, fix)
	if err != nil {
		fatal(err)
	}

	r := &runner{
		repo: repoAbs, fixture: fix, judge: *judge,
		subjectModel: *modelSubject,
		timeout:      time.Duration(*timeout) * time.Second,
		shape:        shape,
		statements:   statementsByID(all["A"]),
		approaches:   approaches,
		ground:       ground,
		today:        time.Now().Format("2006-01-02"),
	}
	if err := preflight(r); err != nil {
		fatal(fmt.Errorf("auth pre-flight failed — expired login or missing credentials? %w", err))
	}
	results := r.runAll(expandEpochs(selected, *epochs), *jobs)
	// One bounded retry for ERROR verdicts: transient API failures
	// (connection drops, judge hiccups) must not fail a gating run.
	var errored []int
	for i, res := range results {
		if res.verdict == vError {
			errored = append(errored, i)
		}
	}
	if len(errored) > 0 {
		retry := make([]evalCase, len(errored))
		for k, i := range errored {
			retry[k] = results[i].c
		}
		fmt.Printf("retrying %d errored case(s) ...\n", len(retry))
		rerun := r.runAll(retry, *jobs)
		for k, i := range errored {
			results[i] = rerun[k]
		}
	}
	results = aggregateEpochs(results, *epochs)

	if err := os.WriteFile(*out, []byte(renderReport(results)), 0o644); err != nil {
		fatal(err)
	}
	fmt.Println(summary(results))
	if *gate && slices.ContainsFunc(results, func(r result) bool { return r.verdict != vPass }) {
		os.Exit(1)
	}
}

// runAll runs cases through a bounded worker pool, printing each verdict as
// it lands; results keep table order regardless of completion order. Cases
// never share state — each gets its own HOME and fixture copy — so the only
// concurrency limit is the account's rate limit.
func (r *runner) runAll(cases []evalCase, jobs int) []result {
	if jobs < 1 {
		jobs = 1
	}
	results := make([]result, len(cases))
	sem := make(chan struct{}, jobs)
	var wg sync.WaitGroup
	var mu sync.Mutex // keeps a verdict line and its reason line together
	for i, c := range cases {
		wg.Add(1)
		go func(i int, c evalCase) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			res := r.runCase(c)
			mu.Lock()
			fmt.Printf("%s: %s\n", c.ID, res.verdict)
			if res.verdict != vPass {
				fmt.Printf("  reason: %s\n", res.reason)
			}
			mu.Unlock()
			results[i] = res
		}(i, c)
	}
	wg.Wait()
	return results
}

// preflight runs one trivial isolated-HOME prompt so an expired login
// fails the run in seconds with a clear message, not as N per-case errors.
func preflight(r *runner) error {
	home, err := os.MkdirTemp("", "preflight-home-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(home)
	env, err := caseEnv(home)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err = runClaude(ctx, home, env, "-p", "reply with: ok", "--max-turns", "1")
	return err
}

// findRoot walks upward from dir to the first directory containing
// skills/mentor, so the runner works from anywhere in the repo — including
// tools/eval-runner itself, where `go -C tools/eval-runner run .` lands.
// Keep in sync with the copies in tools/catalog-lint/main.go and
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
