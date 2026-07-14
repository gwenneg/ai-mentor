package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// casesSnippet is copied from the real evals/cases.md — one representative
// row per group, plus the output-shape block and the interactive-only
// Group D table (which must never parse into runnable cases).
const casesSnippet = "# Eval cases\n\n" +
	"Cases for the discovery-first skill, in four groups.\n\n" +
	"## Group A — Classification (problem mode)\n\n" +
	"Run as `/ai-mentor:mentor <statement>` in the fixture repo.\n\n" +
	"| ID | Problem statement | Expected goal | Notes |\n" +
	"|----|-------------------|--------------|-------|\n" +
	"| A01 | `debug a flaky test that only fails in CI` | debugging | The README's canonical example |\n" +
	"| A03 | `refactor authentication across 30 files` | refactoring | Cross-file scale signal |\n\n" +
	"### Group A output-shape expectations (every classified case)\n\n" +
	"- Opens with the one-sentence load-state announcement, then a diagnosis naming observed evidence — never a questionnaire\n" +
	"- Exactly **one** primary move, with a fenced prompt using at least one real path or command from the fixture repo\n\n" +
	"## Group B — Growth mode (bare invocation)\n\n" +
	"Run as `/ai-mentor:mentor` with a controlled `~/.ai-mentor/profile.md` fixture.\n\n" +
	"| ID | Profile fixture | Expected behavior |\n" +
	"|----|----------------|-------------------|\n" +
	"| B01 | No profile file | First-meeting announcement (names the profile path once); teaches ONE capability from the ignorance map; creates the profile with correct schema |\n\n" +
	"## Group C — Never-repeat under problem mode\n\n" +
	"| ID | Setup | Expected behavior |\n" +
	"|----|-------|-------------------|\n" +
	"| C01 | Profile marks the matched goal's #1 approach `adopted`; run a Group A case for that goal | The move builds on the adopted approach or picks the next-best; it is NOT re-taught from scratch |\n\n" +
	"## Group D — Trigger calibration (interactive only)\n\n" +
	"| ID | Prompt (typed, no slash command) | Expected |\n" +
	"|----|----------------------------------|----------|\n" +
	"| D01 | `what's the best way to use AI to add tests to this codebase?` | Skill fires |\n"

func TestParseCases(t *testing.T) {
	all, shape, err := parseCases(casesSnippet)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a01 := all["A"][0]
	if len(all["A"]) != 2 || a01.ID != "A01" {
		t.Fatalf("group A parsed wrong: %+v", all["A"])
	}
	if a01.Statement != "debug a flaky test that only fails in CI" {
		t.Errorf("A01 statement should have backticks stripped, got %q", a01.Statement)
	}
	if a01.Expected != "debugging" || !strings.Contains(a01.Notes, "canonical") {
		t.Errorf("A01 expected/notes parsed wrong: %+v", a01)
	}
	if len(all["B"]) != 1 || !strings.Contains(all["B"][0].Expected, "First-meeting announcement") {
		t.Errorf("group B parsed wrong: %+v", all["B"])
	}
	if len(all["C"]) != 1 || !strings.Contains(all["C"][0].Statement, "#1 approach") {
		t.Errorf("group C parsed wrong: %+v", all["C"])
	}
	if len(all["D"]) != 0 {
		t.Errorf("group D is interactive-only and must not parse into cases: %+v", all["D"])
	}
	if !strings.Contains(shape, "Exactly **one** primary move") ||
		!strings.Contains(shape, "never a questionnaire") {
		t.Errorf("output-shape block not captured verbatim, got:\n%s", shape)
	}
	if strings.Contains(shape, "| B01 |") {
		t.Errorf("shape block leaked into the next section:\n%s", shape)
	}
}

// A missing group heading must be a loud fatal error for any run that
// requests that group — never a silently green, empty run.
func TestMissingGroupHeadingIsFatal(t *testing.T) {
	corrupted := strings.Replace(casesSnippet,
		"## Group B — Growth mode (bare invocation)", "Growth mode (bare invocation)", 1)
	all, _, err := parseCases(corrupted)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if _, err := selectCases(all, []string{"A", "B", "C"}, nil); err == nil {
		t.Error("requesting group B with its heading gone should be fatal")
	} else if !strings.Contains(err.Error(), "group B") {
		t.Errorf("error should name the empty group, got: %v", err)
	}
	if _, err := selectCases(all, []string{"A", "C"}, nil); err != nil {
		t.Errorf("intact groups should still be selectable: %v", err)
	}
}

func TestSelectCasesFiltersByID(t *testing.T) {
	all, _, err := parseCases(casesSnippet)
	if err != nil {
		t.Fatal(err)
	}
	sel, err := selectCases(all, []string{"A", "B"}, []string{"A03"})
	if err != nil || len(sel) != 1 || sel[0].ID != "A03" {
		t.Errorf("want exactly A03, got %+v (err %v)", sel, err)
	}
	if _, err := selectCases(all, []string{"A"}, []string{"A99"}); err == nil {
		t.Error("a filter matching nothing should be fatal")
	}
	// A partially matching filter must be fatal too, naming the missing ID —
	// this is what keeps the smoke list loud when cases.md drifts.
	if _, err := selectCases(all, []string{"A"}, []string{"A01", "A99"}); err == nil {
		t.Error("an unmatched requested ID should be fatal even when others match")
	} else if !strings.Contains(err.Error(), "A99") {
		t.Errorf("error should name the missing ID, got: %v", err)
	}
}

// The smoke tier is a hand-curated ID list; every entry must exist in the
// real cases.md, or a rename silently shrinks the per-change signal.
func TestSmokeCasesExistInSuite(t *testing.T) {
	root, err := findRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	text, err := os.ReadFile(filepath.Join(root, "evals", "cases.md"))
	if err != nil {
		t.Fatal(err)
	}
	all, _, err := parseCases(string(text))
	if err != nil {
		t.Fatal(err)
	}
	sel, err := selectCases(all, []string{"A", "B", "C"}, smokeCases)
	if err != nil {
		t.Fatalf("smoke tier drifted from cases.md: %v", err)
	}
	if len(sel) != len(smokeCases) {
		t.Errorf("want %d smoke cases, selected %d", len(smokeCases), len(sel))
	}
}

func TestParseVerdictLenient(t *testing.T) {
	v, err := parseVerdict("Here is my verdict:\n```json\n" +
		`{"pass": false, "checks": [{"name": "classification", "pass": false, "reason": "routed to testing"}]}` +
		"\n```\nHope that helps!")
	if err != nil {
		t.Fatalf("prose-wrapped JSON should parse: %v", err)
	}
	if v.Pass || len(v.Checks) != 1 || v.Checks[0].Reason != "routed to testing" {
		t.Errorf("verdict parsed wrong: %+v", v)
	}
	if _, err := parseVerdict("no json here at all"); err == nil {
		t.Error("prose without JSON must be an error")
	}
	if _, err := parseVerdict(`{"checks": []}`); err == nil {
		t.Error(`JSON without a "pass" field must be an error, not a silent fail`)
	}
	if _, err := parseVerdict(`{"pass": maybe}`); err == nil {
		t.Error("invalid JSON must be an error")
	}
}

func TestRenderReport(t *testing.T) {
	longResponse := strings.TrimSuffix(strings.Repeat("line\n", 70), "\n")
	results := []result{
		{c: evalCase{Group: "A", ID: "A01"}, verdict: vPass},
		{c: evalCase{Group: "A", ID: "A03"}, verdict: vFail,
			reason: "classification: routed to testing\nsecond | line", response: longResponse},
		{c: evalCase{Group: "B", ID: "B01"}, verdict: vError, reason: "judge reply not parseable"},
	}
	report := renderReport(results)
	for _, want := range []string{
		"Group A: 1 pass / 1 fail / 0 error",
		"Group B: 0 pass / 0 fail / 1 error",
		"| A01 | PASS |  |",
		"| A03 | FAIL | classification: routed to testing second / line |",
		"| B01 | ERROR | judge reply not parseable |",
		"## Raw failures",
		"### A03 (FAIL)",
		"... (10 more lines truncated)",
	} {
		if !strings.Contains(report, want) {
			t.Errorf("report missing %q; got:\n%s", want, report)
		}
	}
	if strings.Contains(report, "### B01") {
		t.Error("an ERROR case with no captured response should not appear in Raw failures")
	}
	if got := strings.Count(report, "line\n"); got > maxRawLines+5 {
		t.Errorf("raw failure not truncated: %d response lines in report", got)
	}
}

// newTestRunner returns a runner wired to temp dirs and the parsed snippet.
func newTestRunner(t *testing.T) *runner {
	t.Helper()
	all, shape, err := parseCases(casesSnippet)
	if err != nil {
		t.Fatal(err)
	}
	return &runner{
		repo: t.TempDir(), fixture: t.TempDir(), judge: "judge-model",
		subjectModel: "subject-model",
		timeout:      time.Minute,
		shape:        shape,
		statements:   statementsByID(all["A"]),
		approaches:   []string{"plan-mode", "hooks-as-workflow"},
		ground: groundTruth{
			fixture:      []string{"package.json", "src/orders.js"},
			plugins:      []string{"security-guidance", "code-modernization"},
			techniques:   []string{"plan-mode"},
			integrations: []string{"github-actions"},
		},
		today: "2026-07-07",
	}
}

// stubClaude replaces runClaude for the test and records every call,
// including the profile content in the call's HOME at call time (the temp
// HOME is gone by the time runCase returns). Mentor calls return mentorOut;
// judge calls (spotted by --model) judgeOut.
func stubClaude(t *testing.T, mentorOut, judgeOut string) *[][]string {
	t.Helper()
	orig := runClaude
	t.Cleanup(func() { runClaude = orig })
	calls := &[][]string{}
	runClaude = func(ctx context.Context, dir string, env []string, args ...string) (string, error) {
		call := append([]string{"dir=" + dir, "env=" + strings.Join(env, "\x00")}, args...)
		if home := envValue(call, "HOME"); home != "" {
			call = append(call, "profile="+readFile(filepath.Join(home, ".ai-mentor", "profile.md")))
		}
		*calls = append(*calls, call)
		// The mentor runs with --plugin-dir; the judge runs in an empty dir
		// without it. (Both now pass --model, so --model no longer discriminates.)
		if slices.Contains(args, "--plugin-dir") {
			return mentorOut, nil
		}
		return judgeOut, nil
	}
	return calls
}

// callField returns the value of a "key=value" element recorded by stubClaude.
func callField(call []string, key string) string {
	for _, e := range call {
		if v, ok := strings.CutPrefix(e, key+"="); ok {
			return v
		}
	}
	return ""
}

// envValue returns one variable from the recorded child environment.
func envValue(call []string, key string) string {
	for kv := range strings.SplitSeq(callField(call, "env"), "\x00") {
		if v, ok := strings.CutPrefix(kv, key+"="); ok {
			return v
		}
	}
	return ""
}

func argAfter(call []string, flagName string) string {
	if i := slices.Index(call, flagName); i >= 0 && i+1 < len(call) {
		return call[i+1]
	}
	return ""
}

func TestRunCaseStubbed(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key") // skip the credentials copy
	r := newTestRunner(t)
	calls := stubClaude(t, `{"result": "growth-mode lesson"}`, `{"pass": true, "checks": []}`)

	b03 := evalCase{Group: "B", ID: "B03",
		Statement: "A declined row", Expected: "The declined capability is never offered"}
	res := r.runCase(b03)
	if res.verdict != vPass {
		t.Fatalf("stubbed B03 should pass, got %s (%s)", res.verdict, res.reason)
	}
	if len(*calls) != 2 {
		t.Fatalf("want 1 mentor call + 1 judge call, got %d", len(*calls))
	}
	mentor, judge := (*calls)[0], (*calls)[1]

	if got := argAfter(mentor, "-p"); got != mentorCmd {
		t.Errorf("Group B must be the bare invocation, got prompt %q", got)
	}
	if argAfter(mentor, "--plugin-dir") != r.repo {
		t.Errorf("mentor call missing --plugin-dir %s: %v", r.repo, mentor)
	}
	dir := callField(mentor, "dir")
	if dir == r.fixture || !strings.Contains(dir, "eval-fixture-") {
		t.Errorf("mentor must run in a per-case fixture copy, got %s", dir)
	}
	home := envValue(mentor, "HOME")
	if home == "" || home == os.Getenv("HOME") {
		t.Errorf("mentor HOME not isolated: %q", home)
	}
	profile := callField(mentor, "profile")
	if profile == "" {
		t.Fatal("B03 profile fixture not written into the temp HOME before the run")
	}
	if !strings.Contains(profile, `| fan-out-workflows | declined |`) ||
		!strings.Contains(profile, "too token-heavy") {
		t.Errorf("B03 fixture wrong:\n%s", profile)
	}

	jp := argAfter(judge, "-p")
	if argAfter(judge, "--model") != "judge-model" || argAfter(judge, "--max-turns") != "5" {
		t.Errorf("judge call flags wrong: %v", judge)
	}
	for _, want := range []string{"growth-mode lesson", "never offered", "STRICT JSON", "fan-out-workflows"} {
		if !strings.Contains(jp, want) {
			t.Errorf("judge prompt missing %q", want)
		}
	}
}

func TestRunCaseGroupAIncludesShape(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")
	r := newTestRunner(t)
	calls := stubClaude(t, `{"result": "a move"}`, `{"pass": true, "checks": []}`)

	a01 := evalCase{Group: "A", ID: "A01",
		Statement: "debug a flaky test that only fails in CI", Expected: "debugging"}
	if res := r.runCase(a01); res.verdict != vPass {
		t.Fatalf("stubbed A01 should pass, got %s (%s)", res.verdict, res.reason)
	}
	mentor, judge := (*calls)[0], (*calls)[1]
	if got := argAfter(mentor, "-p"); got != mentorCmd+" debug a flaky test that only fails in CI" {
		t.Errorf("Group A prompt wrong: %q", got)
	}
	jp := argAfter(judge, "-p")
	if !strings.Contains(jp, "Exactly **one** primary move") {
		t.Error("Group A judge prompt must include the output-shape expectations verbatim")
	}
}

func TestRunCaseC02RunsTwiceSameHome(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")
	r := newTestRunner(t)
	calls := stubClaude(t, `{"result": "a move"}`, `{"pass": true, "checks": []}`)

	c02 := evalCase{Group: "C", ID: "C02", Statement: "same case twice", Expected: "surprises differ"}
	if res := r.runCase(c02); res.verdict != vPass {
		t.Fatalf("stubbed C02 should pass, got %s (%s)", res.verdict, res.reason)
	}
	if len(*calls) != 3 { // two mentor runs, one judge
		t.Fatalf("C02 must run twice then judge once, got %d calls", len(*calls))
	}
	first, second := (*calls)[0], (*calls)[1]
	want := mentorCmd + " refactor authentication across 30 files"
	if argAfter(first, "-p") != want || argAfter(second, "-p") != want {
		t.Errorf("C02 must reuse A03's statement on both runs: %v / %v",
			argAfter(first, "-p"), argAfter(second, "-p"))
	}
	if envValue(first, "HOME") != envValue(second, "HOME") {
		t.Error("C02 runs must share the same HOME so the profile persists")
	}
}

// Every judge failure path must be able to fail a case, not silently pass.
func TestJudgeFailurePaths(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")
	cases := []struct {
		name, judgeOut, verdict, reason string
	}{
		{"garbage judge reply is ERROR", "sorry, I cannot help", vError, "not parseable"},
		{"pass=false is FAIL with the check reason",
			`{"pass": false, "checks": [{"name": "never-repeat", "pass": false, "reason": "re-taught hooks"}]}`,
			vFail, "never-repeat: re-taught hooks"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newTestRunner(t)
			stubClaude(t, `{"result": "a lesson"}`, tc.judgeOut)
			res := r.runCase(evalCase{Group: "B", ID: "B01", Expected: "first meeting"})
			if res.verdict != tc.verdict || !strings.Contains(res.reason, tc.reason) {
				t.Errorf("want %s with reason containing %q, got %s (%s)",
					tc.verdict, tc.reason, res.verdict, res.reason)
			}
			if res.response == "" {
				t.Error("non-PASS verdicts after a run must keep the response for the report")
			}
		})
	}
}

func TestMentorOutputMustBeJSONWithResult(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")
	for _, out := range []string{"plain text, not json", `{"no_result": true}`} {
		r := newTestRunner(t)
		stubClaude(t, out, `{"pass": true, "checks": []}`)
		if res := r.runCase(evalCase{Group: "B", ID: "B01"}); res.verdict != vError {
			t.Errorf("mentor output %q should yield ERROR, got %s", out, res.verdict)
		}
	}
}

func TestSetupProfileFixtures(t *testing.T) {
	r := newTestRunner(t)
	read := func(t *testing.T, id string) (string, []string) {
		t.Helper()
		home := t.TempDir()
		seeded, err := r.setupProfile(evalCase{ID: id}, home)
		if err != nil {
			t.Fatal(err)
		}
		return readFile(filepath.Join(home, ".ai-mentor", "profile.md")), seeded
	}
	if p, seeded := read(t, "B01"); p != "" || seeded != nil {
		t.Errorf("B01 must start with no profile and no seeded ids, got %v:\n%s", seeded, p)
	}
	if p, seeded := read(t, "B02"); !strings.Contains(p, "| autonomous-loops | shown |") || !slices.Equal(seeded, []string{"autonomous-loops"}) {
		t.Errorf("B02 needs a past shown row and its seeded id (%v):\n%s", seeded, p)
	}
	if p, _ := read(t, "B05"); !strings.Contains(p, "Last new-capability check: 2026-w20") {
		t.Errorf("B05 needs the stale anchor:\n%s", p)
	}
	p, seeded := read(t, "B06")
	for _, a := range r.approaches {
		if !strings.Contains(p, "| "+a+" | adopted |") {
			t.Errorf("B06 must mark every approach adopted, missing %s:\n%s", a, p)
		}
	}
	if !slices.Equal(seeded, r.approaches) {
		t.Errorf("B06 must report every approach as seeded, got %v", seeded)
	}
	if p, _ := read(t, "C01"); !strings.Contains(p, "| plan-mode | adopted |") {
		t.Errorf("C01 needs plan-mode adopted:\n%s", p)
	}
	if p, seeded := read(t, "C04"); !strings.Contains(p, "| background-agents | declined |") || !strings.Contains(p, "| plan-mode | shown |") || len(seeded) != 2 {
		t.Errorf("C04 needs a declined and a shown seeded row (%v):\n%s", seeded, p)
	}
	if p, _ := read(t, "C05"); !strings.Contains(p, "| plan-mode | declined |") {
		t.Errorf("C05 needs plan-mode declined:\n%s", p)
	}
}

// The taught-capability diff feeds the judge's catalog sources: ids parse
// from the after-run profile, seeded fixture rows are excluded, and only ids
// with a real approach file produce sources.
func TestTaughtIDsAndCatalogSources(t *testing.T) {
	r := newTestRunner(t)
	profile := r.profileMD("2026-w28",
		profileRow("plan-mode", "shown", "2026-07-01", "seeded"),
		profileRow("project-memory", "shown", "2026-07-14", "taught today"),
		profileRow("security-guidance", "shown", "2026-07-14", "a plugin id with no approach file"))
	if ids := profileIDs(profile); !slices.Equal(ids, []string{"plan-mode", "project-memory", "security-guidance"}) {
		t.Fatalf("profileIDs wrong: %v", ids)
	}
	taught := taughtIDs(profile, []string{"plan-mode"})
	if !slices.Equal(taught, []string{"project-memory", "security-guidance"}) {
		t.Fatalf("taughtIDs must exclude seeded rows, got %v", taught)
	}
	dir := filepath.Join(r.repo, "skills", "mentor", "approaches", "techniques")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "project-memory.md"), []byte("path-scoped rules in .claude/rules/*.md\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	sources := r.catalogSources(taught)
	if len(sources) != 1 || sources[0].id != "project-memory" {
		t.Fatalf("want one source for project-memory (no file for the plugin id), got %+v", sources)
	}
	if !strings.Contains(sources[0].content, ".claude/rules") {
		t.Errorf("source content not loaded: %+v", sources[0])
	}
	if got := taughtIDs("", nil); got != nil {
		t.Errorf("a missing profile must teach nothing, got %v", got)
	}
}

func TestJudgePromptGroundTruth(t *testing.T) {
	r := newTestRunner(t)
	jp := r.judgePrompt(
		evalCase{Group: "A", ID: "A01", Statement: "x", Expected: "debugging"},
		[]string{"resp"}, "", nil)
	for _, want := range []string{
		"src/orders.js",     // fixture manifest inlined (grounding — #8)
		"security-guidance", // authoritative plugin list inlined (fabrication — #6)
		"COMPLETE list",     // plugin list framed as exhaustive
		"is a fabrication",  // fabrication instruction present
		"/verify",           // commands named in the not-exhaustive caveat
		"NEWER than your training data", // catalog outranks judge memory on technique detail
		"'unverifiable'",    // unknown-mechanism escape hatch, never a memory-based FAIL
	} {
		if !strings.Contains(jp, want) {
			t.Errorf("judge prompt missing %q", want)
		}
	}
	if strings.Contains(jp, "Catalog sources") {
		t.Error("no sources given, so the catalog-sources block must be absent")
	}

	jp = r.judgePrompt(
		evalCase{Group: "B", ID: "B01", Statement: "no profile", Expected: "teaches one capability"},
		[]string{"resp"}, "profile",
		[]capSource{{id: "project-memory", content: "path-scoped rules in .claude/rules/*.md"}})
	for _, want := range []string{
		"Catalog sources",
		"Source for project-memory",
		".claude/rules/*.md", // the verified text the judge must defer to
	} {
		if !strings.Contains(jp, want) {
			t.Errorf("judge prompt with sources missing %q", want)
		}
	}
}

func TestCaseFixtureCopies(t *testing.T) {
	r := newTestRunner(t)
	if err := os.WriteFile(filepath.Join(r.fixture, "package.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	plain, err := r.caseFixture(false)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(plain)
	if plain == r.fixture {
		t.Fatal("every case must run in a copy, not the shared fixture")
	}
	if readFile(filepath.Join(plain, "package.json")) == "" {
		t.Error("fixture contents not copied")
	}
	if _, err := os.Stat(filepath.Join(plain, ".claude")); err == nil {
		t.Error("a plain copy must not carry hooks")
	}

	hooked, err := r.caseFixture(true)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(hooked)
	if !strings.Contains(readFile(filepath.Join(hooked, ".claude", "settings.json")), `"hooks"`) {
		t.Error("hooks settings not written into the B04 copy")
	}
	if _, err := os.Stat(filepath.Join(r.fixture, ".claude")); err == nil {
		t.Error("the shared fixture must stay untouched")
	}
}

// runAll must keep results in table order however the goroutines finish, and
// the pool must actually bound concurrency.
func TestRunAllOrderAndBound(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")
	r := newTestRunner(t)

	var mu sync.Mutex
	inflight, peak := 0, 0
	orig := runClaude
	t.Cleanup(func() { runClaude = orig })
	runClaude = func(ctx context.Context, dir string, env []string, args ...string) (string, error) {
		mu.Lock()
		inflight++
		if inflight > peak {
			peak = inflight
		}
		mu.Unlock()
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		inflight--
		mu.Unlock()
		if slices.Contains(args, "--plugin-dir") {
			return `{"result": "a lesson"}`, nil
		}
		return `{"pass": true, "checks": []}`, nil
	}

	cases := make([]evalCase, 8)
	for i := range cases {
		cases[i] = evalCase{Group: "B", ID: fmt.Sprintf("B%02d", i+1), Expected: "x"}
	}
	results := r.runAll(cases, 2)
	for i, res := range results {
		if res.c.ID != cases[i].ID {
			t.Fatalf("result %d is %s, want %s — order not preserved", i, res.c.ID, cases[i].ID)
		}
		if res.verdict != vPass {
			t.Errorf("%s: want PASS, got %s (%s)", res.c.ID, res.verdict, res.reason)
		}
	}
	// Each case runs a mentor call then a judge call, so with -j 2 at most
	// 2 claude invocations are ever in flight.
	if peak > 2 {
		t.Errorf("concurrency bound violated: %d claude calls in flight with jobs=2", peak)
	}
	if peak < 2 {
		t.Errorf("cases did not actually run concurrently (peak %d)", peak)
	}
}

// Either CI credential must short-circuit the local-credentials copy: with a
// token (or API key) in the env, caseEnv passes the env through untouched and
// never writes a .credentials.json into the temp HOME.
func TestCaseEnvHonorsOAuthToken(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("CLAUDE_CODE_OAUTH_TOKEN", "test-token")
	home := t.TempDir()
	env, err := caseEnv(home)
	if err != nil {
		t.Fatalf("a CLAUDE_CODE_OAUTH_TOKEN alone must satisfy caseEnv: %v", err)
	}
	if !slices.Contains(env, "HOME="+home) {
		t.Error("HOME not pointed at the isolated temp dir")
	}
	if !slices.Contains(env, "CLAUDE_CODE_OAUTH_TOKEN=test-token") {
		t.Error("the token must pass through to the child env")
	}
	if _, err := os.Stat(filepath.Join(home, ".claude", ".credentials.json")); err == nil {
		t.Error("no credentials file should be written when a CI credential is present")
	}
}

func TestExpandEpochs(t *testing.T) {
	cases := []evalCase{{ID: "A01"}, {ID: "B01"}}
	if got := expandEpochs(cases, 1); len(got) != 2 {
		t.Fatalf("epochs=1 must be a no-op, got %d cases", len(got))
	}
	got := expandEpochs(cases, 3)
	var ids []string
	for _, c := range got {
		ids = append(ids, c.ID)
	}
	want := []string{"A01", "A01", "A01", "B01", "B01", "B01"}
	if !slices.Equal(ids, want) {
		t.Errorf("copies must be adjacent per case: got %v", ids)
	}
}

func TestAggregateEpochs(t *testing.T) {
	mk := func(id, verdict, reason string) result {
		return result{c: evalCase{Group: "A", ID: id}, verdict: verdict, reason: reason, response: "resp-" + id}
	}
	results := []result{
		mk("A01", vPass, ""), mk("A01", vPass, ""), mk("A01", vPass, ""),
		mk("A02", vPass, ""), mk("A02", vFail, "bad shape"), mk("A02", vPass, ""),
		mk("A03", vFail, "wrong goal"), mk("A03", vPass, ""), mk("A03", vFail, "wrong goal"),
		mk("A04", vError, "boom"), mk("A04", vError, "boom"), mk("A04", vError, "boom"),
	}
	if got := aggregateEpochs(results, 1); len(got) != len(results) {
		t.Fatalf("epochs=1 must be a no-op, got %d results", len(got))
	}
	agg := aggregateEpochs(results, 3)
	if len(agg) != 4 {
		t.Fatalf("want 4 aggregated results, got %d", len(agg))
	}
	if agg[0].verdict != vPass || agg[0].reason != "" {
		t.Errorf("clean 3/3 must be a PASS with no reason, got %s (%q)", agg[0].verdict, agg[0].reason)
	}
	if agg[1].verdict != vPass {
		t.Errorf("2/3 pass is a strict majority — want PASS, got %s", agg[1].verdict)
	}
	if !strings.Contains(agg[1].reason, "FLAKY 2/3") || !strings.Contains(agg[1].reason, "bad shape") {
		t.Errorf("a flaky pass must stay visible with the failing epoch's reason, got %q", agg[1].reason)
	}
	if agg[1].response != "resp-A02" {
		t.Error("the failing epoch's response must be kept for the report")
	}
	if agg[2].verdict != vFail || !strings.Contains(agg[2].reason, "FLAKY 1/3") {
		t.Errorf("1/3 pass must be a flagged FAIL, got %s (%q)", agg[2].verdict, agg[2].reason)
	}
	if agg[3].verdict != vError || !strings.Contains(agg[3].reason, "0/3 epochs passed") {
		t.Errorf("all-ERROR epochs must aggregate to ERROR, got %s (%q)", agg[3].verdict, agg[3].reason)
	}
}

// A broken or missing catalog read must be fatal, never a silently empty
// fabrication whitelist (the judge fails any plugin absent from it).
func TestBuildGroundTruth(t *testing.T) {
	repo, fixture := t.TempDir(), t.TempDir()
	skill := filepath.Join(repo, "skills", "mentor")
	if _, err := buildGroundTruth(repo, fixture); err == nil {
		t.Error("missing marketplace.md must be an error")
	}
	if err := os.MkdirAll(skill, 0o755); err != nil {
		t.Fatal(err)
	}
	market := filepath.Join(skill, "marketplace.md")
	if err := os.WriteFile(market, []byte("# Plugins\n\nno table rows here\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := buildGroundTruth(repo, fixture); err == nil {
		t.Error("a marketplace.md yielding zero plugin names must be an error")
	}
	if err := os.WriteFile(market, []byte("# Plugins\n\n| `security-guidance` | desc |\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := buildGroundTruth(repo, fixture); err == nil {
		t.Error("zero approach files must be an error")
	}
	writeApproach := func(rel, content string) {
		p := filepath.Join(skill, "approaches", rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeApproach("quality/plan-mode.md", "# Plan mode\n")
	writeApproach("tools/code-modernization.md", "kind: plugin\n")
	writeApproach("integrations/github-actions.md", "kind: integration\n")
	gt, err := buildGroundTruth(repo, fixture)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(gt.plugins, "security-guidance") || !slices.Contains(gt.plugins, "code-modernization") {
		t.Errorf("plugins wrong: %v", gt.plugins)
	}
	if !slices.Equal(gt.promoted, []string{"code-modernization"}) {
		t.Errorf("promoted wrong: %v", gt.promoted)
	}
	if !slices.Equal(gt.techniques, []string{"plan-mode"}) || !slices.Equal(gt.integrations, []string{"github-actions"}) {
		t.Errorf("techniques/integrations wrong: %v / %v", gt.techniques, gt.integrations)
	}
}

// TestLiveJudgeAnchors scores three frozen B01 transcripts with the REAL
// judge (opt-in: LIVE_JUDGE=1) — the seed of a judge-drift anchor set. They
// pin the three load-bearing judge behaviors after the catalog-source fix:
// a catalog-sourced new-feature lesson must PASS (judge memory must not
// overrule the inlined source), a fabricated marketplace plugin must FAIL
// (whitelist ground truth keeps its teeth), and a fabricated built-in
// mechanism must NOT fail the case from judge memory (it has no ground
// truth either way — 'unverifiable', visible but not a false red).
func TestLiveJudgeAnchors(t *testing.T) {
	if os.Getenv("LIVE_JUDGE") == "" {
		t.Skip("set LIVE_JUDGE=1 to score the anchor transcripts with the real judge")
	}
	root, err := findRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	text, err := os.ReadFile(filepath.Join(root, "evals", "cases.md"))
	if err != nil {
		t.Fatal(err)
	}
	all, shape, err := parseCases(string(text))
	if err != nil {
		t.Fatal(err)
	}
	var b01 evalCase
	for _, c := range all["B"] {
		if c.ID == "B01" {
			b01 = c
		}
	}
	if b01.ID == "" {
		t.Fatal("B01 not found in cases.md")
	}
	fixture := filepath.Join(root, "evals", "fixture")
	ground, err := buildGroundTruth(root, fixture)
	if err != nil {
		t.Fatal(err)
	}
	r := &runner{
		repo: root, fixture: fixture, judge: "claude-sonnet-5",
		subjectModel: "claude-sonnet-5", timeout: 5 * time.Minute,
		shape: shape, ground: ground, today: time.Now().Format("2006-01-02"),
	}

	opener := "First time we meet: I keep a lightweight profile of what you already use at ~/.ai-mentor/profile.md — created just now, yours to edit or delete. Quick scan: CLAUDE.md loads (npm test, npm run lint, uppercase discount codes), no other .claude/ config, no MCP servers.\n\n"
	anchors := []struct {
		name, taughtID, response string
		wantPass                 bool
		wantInReason             string
	}{
		{
			name:    "catalog-sourced new feature must PASS",
			taughtID: "project-memory",
			response: opener + "**Today's capability: path-scoped rules.** Your CLAUDE.md loads whole into every session. You can split niche instructions into `.claude/rules/*.md` files with `paths:` frontmatter so a rule loads only when I work on matching files:\n\n```\nmkdir -p .claude/rules\ncat > .claude/rules/orders.md <<'EOF'\n---\npaths:\n  - \"src/orders.js\"\n  - \"tests/orders.test.js\"\n---\nDiscount codes are uppercase; run npm test after edits here.\nEOF\n```\n\nThat keeps the always-on CLAUDE.md lean while orders-specific rules appear exactly when they matter. Want me to create it now?",
			wantPass: true,
		},
		{
			name:    "fabricated marketplace plugin must FAIL",
			taughtID: "test-coverage-pro",
			response: opener + "**Today's capability: automated coverage gating.** The official marketplace covers this: install with `/plugin install test-coverage-pro@claude-plugins-official`, then `/coverage-gate 80` blocks any change that drops line coverage below 80% — it wires itself into npm test automatically. Want me to install it now?",
			wantPass: false, wantInReason: "test-coverage-pro",
		},
		{
			name:    "fabricated built-in must not FAIL from judge memory",
			taughtID: "autopilot-mode",
			response: opener + "**Today's capability: autopilot mode.** Claude Code can commit for you after every green test run: create `.claude/autopilot.yaml` with `autopilot: true` and `on: green-tests`, and each time npm test passes I commit the working tree with a generated message. Want me to enable it?",
			wantPass: true,
		},
	}
	for _, a := range anchors {
		t.Run(a.name, func(t *testing.T) {
			profile := r.profileMD(currentWeek(),
				profileRow(a.taughtID, "shown", r.today, "eval anchor"))
			sources := r.catalogSources([]string{a.taughtID})
			if a.taughtID == "project-memory" && len(sources) == 0 {
				t.Fatal("project-memory must resolve to a catalog source")
			}
			res := r.judgeCase(b01, []string{a.response}, profile, sources)
			t.Logf("verdict=%s reason=%s", res.verdict, res.reason)
			if a.wantPass && res.verdict != vPass {
				t.Errorf("want PASS, got %s: %s", res.verdict, res.reason)
			}
			if !a.wantPass && res.verdict != vFail {
				t.Errorf("want FAIL, got %s: %s", res.verdict, res.reason)
			}
			if a.wantInReason != "" && !strings.Contains(res.reason, a.wantInReason) {
				t.Errorf("reason should name %q, got: %s", a.wantInReason, res.reason)
			}
		})
	}
}

func TestTruncateLines(t *testing.T) {
	if got := truncateLines("a\nb\nc", 5); got != "a\nb\nc" {
		t.Errorf("short input must pass through, got %q", got)
	}
	got := truncateLines(strings.Repeat("x\n", 99)+"x", 60)
	if lines := strings.Count(got, "\n") + 1; lines != 61 {
		t.Errorf("want 60 lines + marker, got %d lines", lines)
	}
	if !strings.HasSuffix(got, "(40 more lines truncated)") {
		t.Errorf("missing truncation marker: %q", got[len(got)-40:])
	}
}

func TestPromptRules(t *testing.T) {
	r := newTestRunner(t)
	for _, id := range []string{"C01", "C04", "C05"} {
		if _, err := r.prompts(evalCase{Group: "C", ID: id}); err != nil {
			t.Errorf("%s should resolve A01's statement: %v", id, err)
		}
	}
	r.statements = map[string]string{}
	if _, err := r.prompts(evalCase{Group: "C", ID: "C01"}); err == nil {
		t.Error("a C case whose A dependency did not parse must be an error")
	}
	if _, err := r.prompts(evalCase{Group: "C", ID: "C99"}); err == nil {
		t.Error("an unknown C case must be an error, not a silent skip")
	}
}
