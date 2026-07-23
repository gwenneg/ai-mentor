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
// every selected case N times, passing on a strict majority ([strict]-marked
// cases need every epoch) and flagging majority-pass mixes FLAKY.
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
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	mentorCmd   = "/ai-mentor:mentor"
	maxRawLines = 60
	// runSeparator joins a case's multiple mentor runs everywhere a single
	// response string is needed (det checks, records, the report) — one
	// literal, so the texts can never diverge.
	runSeparator = "\n\n--- second run ---\n\n"

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
// prompt · A30 fabrication trap (the honesty behavior class needs a cheap
// per-change signal, not only the release gate) · B01 first meeting + profile
// creation · B06 saturated ignorance map · C01 adopted-not-retaught ·
// C02 same-HOME double run.
// selectCases fails loudly if any of these IDs drops out of cases.md.
var smokeCases = []string{"A01", "A08", "A13", "A19", "A20", "A30", "B01", "B06", "C01", "C02"}

// Strictness — hard invariants gated at pass^k (every epoch must pass,
// never a majority: a promise kept two epochs out of three is a broken
// promise) — is declared WHERE CASES ARE DEFINED: a cases.md row carrying
// the exact strictMarker parses into evalCase.Strict, and the marker is
// stripped before any text reaches the judge. parseCases fatals on
// near-miss spellings, so a typo cannot silently downgrade an invariant.
const strictMarker = "[strict]"

var reStrictNearMiss = regexp.MustCompile(`(?i)\[\s*strict\s*\]`)

// Parenthesized Expected values on Group A cases are markers, not goals.
// unclassifiedSentinels opt a case OUT of the classified-case shape (judged
// on its notes alone); classifiedMarkers keep full shape enforcement while
// having no single goal (A18's fallback path is still held to grounding and
// one-move-one-surprise). Any other parenthesized value is fatal at startup
// — this convention must never be extended by accident of punctuation.
var (
	unclassifiedSentinels = map[string]bool{
		"(catalog browse)":   true,
		"(out of scope)":     true,
		"(fabrication trap)": true,
	}
	classifiedMarkers = map[string]bool{"(no dedicated goal)": true}
)

// groupLabels make report headings self-explanatory; the letters stay the
// keys everywhere else (case IDs, -groups, coverage.md, PR markers).
var groupLabels = map[string]string{"A": "problem mode", "B": "growth mode", "C": "never-repeat"}

func groupTitle(g string) string {
	if l := groupLabels[g]; l != "" {
		return g + " — " + l
	}
	return g
}

// Seeded capability ids shared by setupProfile's fixtures and the
// deterministic checks — one identity const per capability, used at every
// seed site and every check, so a rename cannot silently desynchronize a
// check from its fixture (a desynced forbid check would pass everything
// forever, with no failure to notice).
const (
	capFanOut     = "fan-out-workflows"
	capBackground = "background-agents"
	capPlanMode   = "plan-mode"
)

// A deterministic pre-check decides mechanically checkable verdicts before
// the judge runs: cheaper, unbiased, and a FAIL skips the judge call
// entirely. Only verdicts plain code can decide belong here (file existence,
// row survival, forbidden strings); judgment stays with the judge, which
// still runs on every deterministic pass.
type detCheck struct {
	name string
	fn   func(responses, profile string) string // "" = pass, otherwise the failure reason
}

var detChecks = map[string][]detCheck{
	// First meeting must CREATE the profile with the lesson recorded — the
	// exact regression class of announced-file-empty-table.
	"B01": {{"profile-created-with-rows", profileHasRows}},
	// Declined capabilities are invisible: the technique's name or id
	// anywhere in the response is a violation by definition. The pattern is
	// derived from the seeded id (hyphen-or-space, word-bounded, case-
	// insensitive); paraphrases and singular/plural drift stay judge
	// territory.
	"B03": {{"declined-" + capFanOut + "-invisible", forbidCapability(capFanOut)}},
	"C05": {{"declined-" + capPlanMode + "-invisible", forbidCapability(capPlanMode)}},
	// Seeded rows survive and the declined status never regresses.
	"C04": {{"seeded-rows-survive", c04RowsSurvive}},
}

func runDetChecks(id, responses, profile string) string {
	for _, ch := range detChecks[id] {
		if reason := ch.fn(responses, profile); reason != "" {
			return ch.name + ": " + reason
		}
	}
	return ""
}

// profileHasRows fails when the profile is missing or its capability table
// carries zero data rows — a first meeting that teaches must record. It
// delegates to profileIDs, the runner's one profile-table parser, so note
// text ("---", the word "status") can never confuse the verdict and the
// check can't drift from how taughtIDs reads the same file.
func profileHasRows(_, profile string) string {
	if strings.TrimSpace(profile) == "" {
		return "no profile file was written"
	}
	if len(profileIDs(profile)) == 0 {
		return "profile exists but its capability table has zero data rows"
	}
	return ""
}

// forbidCapability fails when the response names the capability in any of
// its natural spellings: hyphens in the id match hyphen or space, matching
// is case-insensitive, and word boundaries prevent false hits inside longer
// words ("plan mode" must not match "plan modeled").
func forbidCapability(id string) func(string, string) string {
	re := regexp.MustCompile(`(?i)\b` + strings.ReplaceAll(regexp.QuoteMeta(id), "-", "[- ]") + `\b`)
	return func(responses, _ string) string {
		if m := re.FindString(responses); m != "" {
			return fmt.Sprintf("response names the declined capability (%q)", m)
		}
		return ""
	}
}

// c04RowsSurvive enforces C04's mechanical half over profileRows (the one
// profile-table iterator): both seeded rows exist exactly once each — the
// case mandates one row per capability, so a duplicate appended row is a
// violation wherever it sits — and the declined status never regresses.
// Note refreshes and shown-row nuances stay with the judge.
func c04RowsSurvive(_, profile string) string {
	status := map[string]string{}
	for _, cs := range profileRows(profile) {
		id := cs[1]
		if id != capBackground && id != capPlanMode {
			continue
		}
		if _, dup := status[id]; dup {
			return "duplicate rows for " + id + " — one row per capability"
		}
		status[id] = cs[2]
	}
	switch {
	case status[capBackground] == "":
		return "seeded declined row (" + capBackground + ") missing from profile"
	case !strings.EqualFold(status[capBackground], "declined"):
		return capBackground + " row's status regressed to " + status[capBackground]
	case status[capPlanMode] == "":
		return "seeded shown row (" + capPlanMode + ") missing from profile"
	}
	return ""
}

// evalCase is one row from a cases.md group table. Statement holds the
// problem statement for Group A and the fixture/setup description otherwise.
type evalCase struct {
	Group, ID, Statement, Expected, Notes string
	Strict                                bool // parsed from cases.md's [strict] marker; gates at pass^k
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
	profile  string // after-run profile content, for records
	judgeRaw string // judge's raw reply, for records / calibration review
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
		// Strictness is syntax, not prose: the exact marker sets the field
		// and is STRIPPED, so gating metadata never reaches the judge. Any
		// other spelling ([Strict], [ strict]) is fatal — a typo'd marker
		// must never silently gate a hard invariant at plain majority.
		for _, f := range []*string{&c.Expected, &c.Notes} {
			if strings.Contains(*f, strictMarker) {
				c.Strict = true
				*f = strings.TrimSpace(strings.ReplaceAll(*f, strictMarker, ""))
			}
		}
		if m := reStrictNearMiss.FindString(c.Expected + " " + c.Notes); m != "" {
			return nil, "", fmt.Errorf("case %s: malformed strict marker %q — use exactly %q", c.ID, m, strictMarker)
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
	fixtureText  string   // every fixture file inlined verbatim, for command/stack claims
	plugins      []string // directory plugins ∪ promoted (the fabrication whitelist)
	promoted     []string // promoted plugin ids — first-class approaches, no tier label
	techniques   []string
	integrations []string
}

// b04Hooks is B04's per-case .claude/settings.json — the one project-level
// setup signal that case is about. Single-sourced so the judge's ground
// truth always matches what caseFixture actually writes: the file lands in
// the copy AFTER buildGroundTruth runs, and a judge told the fixture list
// is "the only real paths" would otherwise brand B04's own expected demo
// path a fabrication.
// The hook command is deliberately instant (gofmt, no compile): the hook's
// eval purpose is to EXIST as a project-config signal, and it fires on the
// subject's own profile writes too — a `go test ./...` hook injected a cold
// compile plus test output into the middle of the Record flow, and B04's
// empty-profile flips appeared with it (2 of 8 epochs post-fixture-v2).
const b04Hooks = `{"hooks":{"PostToolUse":[{"matcher":"Edit|Write","hooks":[{"type":"command","command":"gofmt -l ."}]}]}}` + "\n"

// buildGroundTruth reads the catalog and fixture once so every judge call can
// check recommendations against them. Any failure here must be fatal to the
// caller: the judge prompt frames the plugin list as COMPLETE and fails
// anything absent as a fabrication, so a silently empty or truncated list
// would mass-fail every plugin recommendation.
func buildGroundTruth(repo, fixture string) (groundTruth, error) {
	skill := filepath.Join(repo, "skills", "mentor")
	gt := groundTruth{}
	fixFiles, err := fixtureFiles(fixture)
	if err != nil {
		return gt, err
	}
	gt.fixture = fixFiles
	text, err := fixtureContents(fixture, gt.fixture)
	if err != nil {
		return gt, err
	}
	gt.fixtureText = text
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

// frameFile renders one file block in the judge's fixture-contents format —
// the single owner of the "--- name ---" delimiter, shared with the B04
// settings append so the two can never drift apart.
func frameFile(name, content string) string {
	return fmt.Sprintf("--- %s ---\n%s\n", name, strings.TrimRight(content, "\n"))
}

// fixtureContents inlines every fixture file verbatim so the judge can
// check fenced commands and stack claims against real content, not just
// path membership — a fenced `npm run <invented>` used to sail through a
// paths-only whitelist. Ground-truth doctrine applies: the judge is told
// these contents are COMPLETE, so a read failure or an over-cap fixture
// must be fatal, never a silent omission that mass-fails grounded fences.
func fixtureContents(dir string, files []string) (string, error) {
	var b strings.Builder
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			return "", fmt.Errorf("fixture ground truth: %w", err)
		}
		b.WriteString(frameFile(f, string(data)))
	}
	const max = 8 * 1024
	if b.Len() > max {
		return "", fmt.Errorf("fixture ground truth: contents are %d bytes (cap %d) — the judge prompt inlines them as COMPLETE; trim the fixture or raise the cap deliberately", b.Len(), max)
	}
	return b.String(), nil
}

// fixtureFiles lists the fixture repo's files as repo-relative paths. The
// enumerator matches fixtureContents' strictness: a walk error is fatal (a
// silently dropped entry would make the judge's "only real paths" list
// incomplete — the same mass-fail mode the reader treats as fatal), and
// only well-known editor/OS junk is filtered — never dotfiles wholesale,
// because future fixture variants may legitimately carry .mcp.json or
// .claude/ files as ground truth.
func fixtureFiles(dir string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if name == ".git" || name == ".idea" {
				return fs.SkipDir
			}
			return nil
		}
		if name == ".DS_Store" || name == "Thumbs.db" || strings.HasSuffix(name, "~") ||
			strings.HasSuffix(name, ".swp") || strings.HasSuffix(name, ".swo") || strings.HasPrefix(name, ".#") {
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		out = append(out, rel)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fixture ground truth: %w", err)
	}
	slices.Sort(out)
	return out, nil
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
	specs         map[string]v2Spec // machine expectations — the V2 deterministic contract
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
	seeded, err := r.setupProfile(c, home)
	if err != nil {
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
	joined := strings.Join(responses, runSeparator)
	if reason := runDetChecks(c.ID, joined, profile); reason != "" {
		return result{c: c, verdict: vFail, reason: "deterministic: " + reason, response: joined, profile: profile}
	}
	gating, advisory := v2Checks(c, r.specs[c.ID], responses, r.ground.plugins, r.ground.promoted)
	if gating != "" {
		return result{c: c, verdict: vFail, reason: "structural: " + gating, response: joined, profile: profile}
	}
	sources := r.catalogSources(taughtIDs(profile, seeded))
	res := r.judgeCase(c, responses, joined, profile, sources)
	if advisory != "" && res.verdict == vPass {
		// Discipline-tier observation: visible in reports and records (the
		// Phase 3 rate baseline), never verdict-affecting in Phase 2.
		res.reason = strings.TrimSpace("advisory-structural: " + advisory + " | " + res.reason)
	}
	return res
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
	if err := os.WriteFile(settings, []byte(b04Hooks), 0o644); err != nil {
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

// setupProfile writes the per-case ~/.ai-mentor/profile.md fixture and
// returns the capability ids it seeded, so the after-run profile can be
// diffed for what the run actually taught. Cases without an entry here
// (B01, all of Group A, C02, C03) start profile-less.
func (r *runner) setupProfile(c evalCase, home string) ([]string, error) {
	past := time.Now().AddDate(0, 0, -21).Format("2006-01-02")
	week := currentWeek()
	var seeded []string
	row := func(name, status, note string) string {
		seeded = append(seeded, name)
		return profileRow(name, status, past, note)
	}
	var content string
	switch c.ID {
	case "B02":
		content = r.profileMD(week, row("autonomous-loops", "shown", "Demoed /loop on a flaky retry test"))
	case "B03":
		content = r.profileMD(week, row(capFanOut, "declined", `"too token-heavy"`))
	case "B04":
		content = r.profileMD(week) // empty profile; the hooks live in the fixture copy
	case "B05":
		content = r.profileMD("2026-w20")
	case "B06":
		rows := make([]string, len(r.approaches))
		for i, a := range r.approaches {
			rows[i] = row(a, "adopted", "eval fixture")
		}
		content = r.profileMD(week, rows...)
	case "C01":
		content = r.profileMD(week, row(capPlanMode, "adopted", "uses plan mode daily"))
	case "C04":
		content = r.profileMD(week,
			row(capBackground, "declined", `"prefer local runs"`),
			row(capPlanMode, "shown", "tried it once"))
	case "C05":
		content = r.profileMD(week, row(capPlanMode, "declined", `"too slow for me"`))
	default:
		return nil, nil
	}
	path := filepath.Join(home, profileRel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return seeded, os.WriteFile(path, []byte(content), 0o644)
}

// profileRows yields the data-row cell slices of a profile's capability
// table, skipping the header and separators (including |:-:| alignment
// forms). It is THE row iterator for profile files — every consumer
// (profileIDs/taughtIDs, the B01 and C04 deterministic checks) must go
// through it, so a format quirk is fixed once, never per-parser.
func profileRows(profile string) [][]string {
	var rows [][]string
	for _, l := range strings.Split(profile, "\n") {
		if !strings.HasPrefix(l, "|") {
			continue
		}
		cs := cells(l)
		if len(cs) < 3 || cs[1] == "" || cs[1] == "Capability" ||
			strings.HasPrefix(cs[1], "-") || strings.HasPrefix(cs[1], ":") {
			continue
		}
		rows = append(rows, cs)
	}
	return rows
}

// profileIDs extracts the capability ids (first column) from a profile's
// table rows.
func profileIDs(profile string) []string {
	var ids []string
	for _, cs := range profileRows(profile) {
		ids = append(ids, cs[1])
	}
	return ids
}

// capSource is one catalog file inlined into the judge prompt as ground
// truth for a capability the run taught.
type capSource struct{ id, content string }

// Bounds on inlined catalog sources: a run records a move plus a surprise
// (plus C04's refreshed row), so 4 files is headroom, and the byte cap keeps
// a pathological profile from exploding the judge prompt.
const (
	maxSources     = 4
	maxSourceBytes = 40_000
)

// catalogSources loads the approach files behind the taught capability ids,
// so the judge checks taught mechanisms against the doc-verified catalog
// instead of its own (older) training data. Ids with no approach file
// (marketplace directory plugins) are skipped: the plugin whitelist already
// covers those.
func (r *runner) catalogSources(taught []string) []capSource {
	var out []capSource
	total := 0
	for _, id := range taught {
		if len(out) == maxSources || total > maxSourceBytes {
			break
		}
		matches, err := filepath.Glob(filepath.Join(r.repo, "skills", "mentor", "approaches", "*", id+".md"))
		if err != nil || len(matches) == 0 {
			continue
		}
		if content := readFile(matches[0]); content != "" {
			out = append(out, capSource{id: id, content: content})
			total += len(content)
		}
	}
	return out
}

// taughtIDs returns the profile ids the run added beyond the seeded fixture.
func taughtIDs(profile string, seeded []string) []string {
	var out []string
	for _, id := range profileIDs(profile) {
		if !slices.Contains(seeded, id) {
			out = append(out, id)
		}
	}
	return out
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
// judgeCase scores the case; joined is the runSeparator-joined response text
// (computed once in runCase so det checks, records, and the report all see
// the identical string by construction).
// stripTrailer removes the mentor trailer from text handed to the judge:
// Phase 1 trailers are observational — nothing may gate on them, including
// a judge side-reading a trailer/prose mismatch (it did, A07, 2026-07-22).
// Det checks and records keep the raw response; accuracy analysis is
// offline until Phase 2 makes the trailer a first-class checked artifact.
func stripTrailer(text string) string {
	return strings.TrimSpace(reTrailer.ReplaceAllString(text, ""))
}

func (r *runner) judgeCase(c evalCase, responses []string, joined, profile string, sources []capSource) result {
	res := result{c: c, response: joined, profile: profile}
	stripped := make([]string, len(responses))
	for i, resp := range responses {
		stripped[i] = stripTrailer(resp)
	}
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
		"-p", r.judgePrompt(c, stripped, profile, sources), "--model", r.judge, "--max-turns", "5")
	if err != nil {
		res.verdict, res.reason = vError, err.Error()
		return res
	}
	res.judgeRaw = out
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

// consistencyPrompt is the V2 short form: one gating question.
func (r *runner) consistencyPrompt(c evalCase, responses []string, profile string) string {
	var b strings.Builder
	b.WriteString("You are grading one response from the 'mentor' skill. The deterministic layer has ALREADY verified classification, structure, grounding, labels, and profile state in code — none of that is your job, and none of it may affect your verdict.\n\n")
	fmt.Fprintf(&b, "Case %s (Group %s). Statement / setup: %s\nExpected behavior (context only, verified elsewhere): %s\n", c.ID, c.Group, c.Statement, c.Expected)
	trailer := parseTrailer(strings.Join(responses, "\n"))
	fmt.Fprintf(&b, "\nThe response declared about itself: %s\n", trailer)
	for i, resp := range responses {
		fmt.Fprintf(&b, "\nResponse %d (metadata removed): <<<\n%s\n>>>\n", i+1, stripTrailer(resp))
	}
	if profile != "" {
		fmt.Fprintf(&b, "\nAfter-run profile: <<<\n%s\n>>>\n", profile)
	}
	b.WriteString("\nYour ONE gating question: does the response's substance contradict its own declaration? (Declared move never actually taught, declared goal plainly not what the content addresses, declared surprise not the labeled pick.) Minor wording differences are not contradictions.\n")
	b.WriteString("Reply with STRICT JSON only: {\"pass\": bool, \"checks\": [{\"name\": string, \"pass\": bool, \"reason\": string}]}. pass=false ONLY for a real contradiction, in a check named 'consistency'. Anything else you notice — quality, tone, rule deviations — goes in checks entries with pass=true and the observation in reason (advisory; it never affects pass).\n")
	return b.String()
}

// judgePrompt builds the scoring prompt. V2 two-path design: cases whose
// substance stays LLM-judged (the machine-expectations table's `judge`
// column — fabrication/ledger/invented-lesson semantics) keep the full
// ground-truth prompt below; every other case gets the SHORT form, whose
// only gating question is trailer/prose consistency — classification,
// structure, grounding, labels, and profile state are all verified in
// code before the judge is ever called.
func (r *runner) judgePrompt(c evalCase, responses []string, profile string, sources []capSource) string {
	if !r.specs[c.ID].Judge {
		return r.consistencyPrompt(c, responses, profile)
	}
	var b strings.Builder
	b.WriteString("You are a strict evaluator for the ai-mentor Claude Code skill. Do not use any tools — judge only the material in this prompt and reply immediately. ")
	b.WriteString("You cannot see the run's tool calls: never infer from the text whether checks actually ran, and a missing profile file does not mean reads were skipped. ")
	b.WriteString("The response was produced INSIDE the fixture repo described by the problem context — its file paths are the fixture's, not this machine's; do not judge them against any other repository. ")
	b.WriteString("Approaches/techniques (plan mode, hooks, browser integration, ...) and built-in slash commands are NOT catalog plugins — plugin rules (tier labels, ⚠️ alternatives) apply only to marketplace plugins installed via /plugin install. ")
	b.WriteString("Judge whether the response(s) below meet the case expectation. ")
	b.WriteString("Reply with STRICT JSON only — no prose, no markdown fences: ")
	b.WriteString(`{"pass": bool, "checks": [{"name": string, "pass": bool, "reason": string}]}`)
	// A run started before midnight can be judged after it: the model stamps
	// profile rows with the real date at write time, so when the dates differ
	// the judge must accept either day as "today".
	if now := time.Now().Format("2006-01-02"); now != r.today {
		fmt.Fprintf(&b, "\n\nToday's date: %s or %s — this run crossed midnight, so a profile row dated either day counts as dated today.\nCase %s (Group %s).\n", r.today, now, c.ID, c.Group)
	} else {
		fmt.Fprintf(&b, "\n\nToday's date: %s\nCase %s (Group %s).\n", r.today, c.ID, c.Group)
	}
	if c.Group == "A" {
		// Only ALLOWLISTED sentinels opt out of the classified-case shape.
		// A18's "(no dedicated goal)" is deliberately NOT one: it has no
		// goal to name, but grounding, one-move-one-surprise, and the
		// closing line still bind it — punctuation alone never decides.
		if unclassifiedSentinels[c.Expected] {
			fmt.Fprintf(&b, "Problem statement: %s\nThis is NOT a classified case — its Expected column carries the marker %s instead of a goal. Judge it on its case notes alone; the classified-case shape expectations below bind only where the notes don't override them.\n", c.Statement, c.Expected)
		} else {
			fmt.Fprintf(&b, "Problem statement: %s\nExpected goal classification: %s\n", c.Statement, c.Expected)
		}
		if c.Notes != "" {
			fmt.Fprintf(&b, "Case notes: %s\n", c.Notes)
		}
		b.WriteString("\nGroup A output-shape expectations (verbatim from cases.md; every classified case must satisfy all of them):\n")
		b.WriteString(r.shape + "\n")
		b.WriteString("\nThe case notes take precedence over the shape expectations when they conflict.\n")
		b.WriteString("The move's \"Do it now\" offer and the surprise's \"want me to show it\" offer are mandated output shape — never fail a response for their presence.\n")
		b.WriteString("For a problem about the fixture repo itself, a fenced prompt that cites a code path not in the fixture file list below is fabricated grounding — fail it (unless the fence is portable: the problem targets a different repo or names code the fixture does not contain, or the case notes mark it portable).\n")
		b.WriteString("The only grounding exceptions are those written here or in the case notes — never invent one.\n")
	} else {
		fmt.Fprintf(&b, "Setup / profile fixture: %s\nExpected behavior: %s\n", c.Statement, c.Expected)
	}

	b.WriteString("\n--- Ground truth: judge the response's recommendations against THIS, not your own memory ---\n")
	if len(r.ground.fixture) > 0 {
		files := r.ground.fixture
		text := r.ground.fixtureText
		// Bulleted like the promoted list: judges misscan comma runs.
		b.WriteString("Fixture repo files (the only real paths in the fixture repo):\n")
		for _, f := range files {
			b.WriteString("- " + f + "\n")
		}
		if text != "" {
			b.WriteString("Complete fixture file contents — ground truth for any claim about the repo's stack, commands, or code. Judge each fence by the mode it chose: for a fence grounded in this repo, a command that could not work against these files is fabricated grounding; a portable fence — one for a problem that targets a different repo or names code the fixture does not contain — correctly carries foreign-stack commands and placeholders instead, but a fence mixing fixture paths with commands these files cannot support gets no portable excuse:\n" + text)
		}
	}
	b.WriteString("Real marketplace plugins (COMPLETE list; installed as `<name>@claude-plugins-official`). A recommended plugin whose name is NOT in this list is a fabrication — fail the case and name it:\n")
	b.WriteString(strings.Join(r.ground.plugins, ", ") + "\n")
	// Bulleted, not comma-joined: a judge scanning a 15-name comma run has
	// missed entries and failed cases on plugins that ARE in this list.
	b.WriteString("Of the plugins above, these are PROMOTED first-class approaches (hands-on validated, ranked in the playbooks). They carry NO tier label — expecting a label on them is an error; their record facts count as hands-on validated. Tier-label rules apply only to the remaining (directory) plugins. Check this list name by name before calling a plugin unpromoted:\n")
	for _, p := range r.ground.promoted {
		b.WriteString("- " + p + "\n")
	}
	b.WriteString("For directory (unpromoted) plugins, the tier marker and the \"not hands-on evaluated\" disclaimer label are DISTINCT requirements — a marker alone does not satisfy a case or rule that requires the label.\n")
	b.WriteString("Known-real techniques: " + strings.Join(r.ground.techniques, ", ") + ". Known-real integrations: " + strings.Join(r.ground.integrations, ", ") + ".\n")
	b.WriteString("These technique/integration lists are NOT exhaustive of Claude Code, and built-in slash commands are not listed at all (e.g. /code-review, /verify, /goal, /loop, /schedule, /init, /plan, /model, /effort, --worktree, Shift+Tab are all real). The plugin list above IS complete: judge plugin recommendations strictly against it.\n")
	// The judge's knowledge ends at its training cutoff; the catalog is
	// doc-verified and tracks features shipped after it. A hermetic judge
	// failed B01 for teaching .claude/rules paths-frontmatter (real,
	// documented, min-version 2.1.198) as "a fabricated feature" — and a
	// names-only trust rule did not stop it: the judge ruled the mechanism
	// "not part of project-memory" from that same stale memory. Only
	// inlining the catalog files themselves closes that hole, so the judge
	// checks taught mechanisms against verified text, not recollection.
	b.WriteString("Your knowledge of Claude Code features ends at your training cutoff; this repo's catalog is verified against current official docs and tracks features shipped after your cutoff — it is NEWER than your training data. A mechanism described in a catalog source below is real, whatever your memory says. A taught mechanism or command you find in NEITHER the catalog sources NOR the lists above must be recorded as a check named 'unverifiable' with pass=true and the detail in its reason — never fail the case from your own memory of what Claude Code supports. Fabrication FAILs are reserved for ground you actually hold: a recommended plugin absent from the complete plugin list, or a fenced path not in the fixture file list.\n")
	// What's-new behavior (B05) is judged against the ACTUAL ledger, not the
	// expectation's description of it: maintenance appends rows weekly, so any
	// prose snapshot of the ledger's state goes stale. Inline the file whenever
	// the case expectation defers to it.
	if strings.Contains(c.Expected, "ledger") || strings.Contains(c.Notes, "ledger") {
		if ledger := readFile(filepath.Join(r.repo, "skills", "mentor", "processed-changelogs.md")); ledger != "" {
			b.WriteString("\n--- Ledger ground truth (skills/mentor/processed-changelogs.md, verbatim). Judge what's-new claims against THESE rows only: a row whose action column records real catalog updates is real news; 'Initial bootstrap' and no-action rows are not. Never assume the ledger's state from the expectation text. ---\n")
			b.WriteString(ledger + "\n")
		}
	}
	if len(sources) > 0 {
		b.WriteString("\n--- Catalog sources for the capabilities this run recorded in the profile (the doc-verified files the lesson was drawn from) ---\n")
		for _, s := range sources {
			fmt.Fprintf(&b, "Source for %s:\n<<<\n%s\n>>>\n", s.id, s.content)
		}
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
// verdict: PASS on a strict majority of passing epochs (every epoch for
// [strict] cases), ERROR when every epoch errored (or a strict shortfall is
// error-only harness noise), FAIL otherwise. Majority-pass mixes are
// flagged FLAKY so they stay visible; gate-blocking verdicts carry a plain
// tally (or the STRICT label) instead — FLAKY never marks a red.
// Relies on expandEpochs's adjacency invariant.
func aggregateEpochs(results []result, n int) []result {
	if n <= 1 {
		// Single-epoch runs (smoke, default dispatch) keep their verdicts,
		// but a failed strict invariant still announces itself in the report.
		for i := range results {
			if results[i].c.Strict && results[i].verdict == vFail {
				results[i].reason = "STRICT invariant: " + results[i].reason
			}
		}
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
	pass, fail, firstFail, firstBad := 0, 0, -1, -1
	for i, r := range chunk {
		switch r.verdict {
		case vPass:
			pass++
		case vFail:
			fail++
			if firstFail < 0 {
				firstFail = i
			}
		}
		if r.verdict != vPass && firstBad < 0 {
			firstBad = i
		}
	}
	n := len(chunk)
	strict := agg.c.Strict
	switch {
	case pass == 0 && fail == 0:
		agg.verdict = vError
	case strict && pass < n:
		// pass^k: one failed epoch breaks the invariant. But an ERRORED
		// epoch is harness noise, not a broken promise — with zero actual
		// FAILs the verdict is ERROR, so triage doesn't blame the model.
		if fail == 0 {
			agg.verdict = vError
		} else {
			agg.verdict = vFail
		}
	case pass*2 > n: // a strict case reaching here has pass == n
		agg.verdict = vPass
	default:
		agg.verdict = vFail
	}
	// Cite a FAIL epoch when one exists: an ERROR epoch's transport message
	// misdirects triage when a real failure is available to quote.
	bad := firstFail
	if bad < 0 {
		bad = firstBad
	}
	if bad >= 0 {
		agg.reason = fmt.Sprintf("epoch %d: %s", bad+1, chunk[bad].reason)
		agg.response = chunk[bad].response
	}
	// FLAKY marks exactly one thing: mixed epochs that still PASSED on
	// majority. Gate-blocking verdicts never carry it — a strict red is
	// labeled STRICT whether it failed or error-downgraded.
	switch {
	case strict && agg.verdict != vPass:
		agg.reason = fmt.Sprintf("STRICT invariant (pass^%d required): %d/%d epochs passed — %s", n, pass, n, agg.reason)
	case agg.verdict == vPass && pass < n:
		agg.reason = fmt.Sprintf("FLAKY %d/%d epochs passed — %s", pass, n, agg.reason)
	case agg.verdict != vPass:
		agg.reason = fmt.Sprintf("%d/%d epochs passed — %s", pass, n, agg.reason)
	}
	return agg
}

// checkCoverage cross-checks evals/coverage.md against the parsed cases so
// the matrix can never silently go stale (the failure mode that rotted B05's
// case spec): every A/B/C ID the matrix references must exist in cases.md,
// and every runnable case must appear in the matrix — a new case ships with
// its coverage row, or the run fails loudly. Ranges like "A01–A29" expand.
func checkCoverage(repo string, known map[string]bool) error {
	full := readFile(filepath.Join(repo, "evals", "coverage.md"))
	if full == "" {
		return fmt.Errorf("evals/coverage.md missing or empty — the rule-coverage matrix is required")
	}
	// Only TABLE ROWS are coverage claims: prose (the gap queue, legend,
	// notes outside tables) is free text, so planning language may name
	// future cases without bricking every run, and a prose mention can
	// never satisfy the has-a-row requirement.
	var tableLines []string
	for _, l := range strings.Split(full, "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), "|") {
			tableLines = append(tableLines, l)
		}
	}
	text := strings.Join(tableLines, "\n")
	seen := map[string]bool{}
	// Ranges expand with every interior ID validated — an interior case
	// deleted from cases.md must fail here, not hide inside "A01–A29".
	for _, m := range regexp.MustCompile(`([ABC])(\d{2})\s*[–-]\s*([ABC])(\d{2})`).FindAllStringSubmatch(text, -1) {
		if m[1] != m[3] {
			return fmt.Errorf("coverage.md range %q mixes groups", m[0])
		}
		lo, _ := strconv.Atoi(m[2])
		hi, _ := strconv.Atoi(m[4])
		if hi < lo {
			return fmt.Errorf("coverage.md range %q is reversed or empty", m[0])
		}
		for i := lo; i <= hi; i++ {
			id := fmt.Sprintf("%s%02d", m[1], i)
			if !known[id] {
				return fmt.Errorf("coverage.md range %q includes %s, which is not in cases.md — the matrix went stale", m[0], id)
			}
			seen[id] = true
		}
	}
	for _, id := range regexp.MustCompile(`\b[ABC]\d{2}\b`).FindAllString(text, -1) {
		seen[id] = true
		if !known[id] {
			return fmt.Errorf("coverage.md references %s, which is not in cases.md — the matrix went stale", id)
		}
	}
	for id := range known {
		if !seen[id] {
			return fmt.Errorf("case %s has no row in coverage.md — map it (or record it as a deliberate gap) before it can gate", id)
		}
	}
	return nil
}

// record is one per-attempt verdict written to the -records JSONL:
// everything a human needs to independently re-judge the case (calibration)
// and everything flake analysis needs. Epoch and Attempt make the history
// self-describing: attempt 1 is the original epoch run, attempt 2 its
// bounded ERROR retry — without them, a retried record is indistinguishable
// from an extra epoch.
type record struct {
	Case     string `json:"case"`
	Group    string `json:"group"`
	Epoch    int    `json:"epoch"`
	Attempt  int    `json:"attempt"`
	Verdict  string `json:"verdict"`
	Trailer  string `json:"trailer,omitempty"` // raw mentor trailer fields (V2 Phase 1: observed, never gated)
	Reason   string `json:"reason,omitempty"`
	Judge    string `json:"judge,omitempty"`
	Response string `json:"response,omitempty"`
	Profile  string `json:"profile,omitempty"`
}

// reTrailer extracts the mentor's machine-readable trailer (V2 Phase 1:
// recorded for observation; no check reads it yet).
var reTrailer = regexp.MustCompile(`<!--\s*mentor\s+([^>]*?)-->`)

func parseTrailer(response string) string {
	if m := reTrailer.FindStringSubmatch(response); m != nil {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func toRecord(r result, epoch, attempt int) record {
	return record{
		Case: r.c.ID, Group: r.c.Group, Epoch: epoch, Attempt: attempt,
		Verdict: r.verdict, Trailer: parseTrailer(r.response), Reason: r.reason,
		Judge: r.judgeRaw, Response: r.response, Profile: r.profile,
	}
}

func writeRecords(path string, recs []record) error {
	var b strings.Builder
	for _, rec := range recs {
		line, err := json.Marshal(rec)
		if err != nil {
			return err
		}
		b.Write(line)
		b.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
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
		parts = append(parts, fmt.Sprintf("Group %s: %d pass / %d fail / %d error", groupTitle(g), pass, fail, errs))
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
	// V2 Phase 1 observability: emission fidelity for the mentor trailer.
	withTrailer := 0
	for _, r := range results {
		if parseTrailer(r.response) != "" {
			withTrailer++
		}
	}
	fmt.Fprintf(&b, "\nTrailers: %d/%d responses carried a parseable mentor trailer (observed only; nothing gates on this yet).\n", withTrailer, len(results))
	for _, g := range groupsIn(results) {
		fmt.Fprintf(&b, "\n## Group %s\n\n| Case | Verdict | Reason |\n|------|---------|--------|\n", groupTitle(g))
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
	epochs := flag.Int("epochs", 1, "independent runs per case; with N>1 a case passes on a strict majority of epochs ([strict]-marked cases need every epoch) and majority-pass mixes are flagged FLAKY")
	passk := flag.Bool("passk", false, "treat every selected case as a strict pass^k invariant: each must pass all -epochs runs — the bar for a PR that claims to fix a case")
	jobs := flag.Int("j", 3, "cases to run concurrently (keep modest: every case is a subject run plus a judge run against the same account)")
	judge := flag.String("model-judge", "claude-sonnet-5", "judge model for scoring")
	modelSubject := flag.String("model-subject", "claude-sonnet-5", "model the mentor under test runs on (pinned so a gate red is a regression, not CLI-default drift)")
	timeout := flag.Int("timeout", 300, "per-case timeout in seconds")
	records := flag.String("records", "", "write per-epoch JSONL verdict records (case, verdict, reason, judge reply, response, profile) — the raw material for judge calibration and run history")
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
	specs, err := parseV2Specs(string(text))
	if err != nil {
		fatal(err)
	}
	// Two-way drift guard: every headless case has a machine-expectations
	// row and every row names a real case — the V2 contract can never
	// silently drop a case's scrutiny or grade a ghost.
	for g, cs := range all {
		if g == "D" {
			continue
		}
		for _, c := range cs {
			if _, ok := specs[c.ID]; !ok {
				fatal(fmt.Errorf("case %s has no machine-expectations row — add one before it can run", c.ID))
			}
		}
	}
	for id := range specs {
		found := false
		for _, cs := range all {
			for _, c := range cs {
				if c.ID == id {
					found = true
				}
			}
		}
		if !found {
			fatal(fmt.Errorf("machine-expectations row %s names no case in the suite", id))
		}
	}
	selected, err := selectCases(all, splitList(*groups), idList)
	if err != nil {
		fatal(err)
	}
	if *passk {
		for i := range selected {
			selected[i].Strict = true
		}
	}
	// The deterministic-check registry is keyed by case ID in Go; a renamed
	// or deleted case must fail the run loudly, not silently drop its extra
	// scrutiny. Parenthesized Expected values must be known markers — the
	// sentinel convention is never extended by accident of punctuation.
	known := map[string]bool{}
	for _, cs := range all {
		for _, c := range cs {
			known[c.ID] = true
			if c.Group == "A" && strings.HasPrefix(c.Expected, "(") && !unclassifiedSentinels[c.Expected] && !classifiedMarkers[c.Expected] {
				fatal(fmt.Errorf("case %s: unknown parenthesized Expected %q — add it to the sentinel allowlists or use a goal name", c.ID, c.Expected))
			}
		}
	}
	for id := range detChecks {
		if !known[id] {
			fatal(fmt.Errorf("detChecks id %s is not in cases.md — update the registry", id))
		}
	}
	if err := checkCoverage(repoAbs, known); err != nil {
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
		specs:        specs,
		today:        time.Now().Format("2006-01-02"),
	}
	if err := preflight(r); err != nil {
		fatal(fmt.Errorf("auth pre-flight failed — expired login or missing credentials? %w", err))
	}
	results := r.runAll(expandEpochs(selected, *epochs), *jobs)
	// First attempts are captured for the records BEFORE the retry splice:
	// the flake history must contain transient ERRORs, not their retries'
	// disguises. Each record carries its epoch and attempt number so the
	// history is reconstructable from the file alone.
	recs := make([]record, 0, len(results))
	for i, r := range results {
		recs = append(recs, toRecord(r, i%*epochs+1, 1))
	}
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
			recs = append(recs, toRecord(rerun[k], i%*epochs+1, 2))
		}
	}
	// Records are optional raw material (calibration, flake history):
	// written first and warn-only, so neither their failure nor a later
	// report failure can discard the other artifact of a fully-paid run.
	if *records != "" {
		if err := writeRecords(*records, recs); err != nil {
			fmt.Fprintln(os.Stderr, "warning: writing records:", err)
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
