package main

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
			commands:     []string{"/verify", "/loop"},
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
	if mentor[0] != "dir="+r.fixture {
		t.Errorf("mentor must run in the fixture dir, got %s", mentor[0])
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
	read := func(t *testing.T, id string) string {
		t.Helper()
		home := t.TempDir()
		if err := r.setupProfile(evalCase{ID: id}, home); err != nil {
			t.Fatal(err)
		}
		return readFile(filepath.Join(home, ".ai-mentor", "profile.md"))
	}
	if p := read(t, "B01"); p != "" {
		t.Errorf("B01 must start with no profile, got:\n%s", p)
	}
	if p := read(t, "B02"); !strings.Contains(p, "| autonomous-loops | shown |") {
		t.Errorf("B02 needs a past shown row:\n%s", p)
	}
	if p := read(t, "B05"); !strings.Contains(p, "Last new-capability check: 2026-w20") {
		t.Errorf("B05 needs the stale anchor:\n%s", p)
	}
	p := read(t, "B06")
	for _, a := range r.approaches {
		if !strings.Contains(p, "| "+a+" | adopted |") {
			t.Errorf("B06 must mark every approach adopted, missing %s:\n%s", a, p)
		}
	}
	if p := read(t, "C01"); !strings.Contains(p, "| plan-mode | adopted |") {
		t.Errorf("C01 needs plan-mode adopted:\n%s", p)
	}
	if p := read(t, "C04"); !strings.Contains(p, "| background-agents | declined |") || !strings.Contains(p, "| plan-mode | shown |") {
		t.Errorf("C04 needs a declined and a shown seeded row:\n%s", p)
	}
	if p := read(t, "C05"); !strings.Contains(p, "| plan-mode | declined |") {
		t.Errorf("C05 needs plan-mode declined:\n%s", p)
	}
}

func TestJudgePromptGroundTruth(t *testing.T) {
	r := newTestRunner(t)
	jp := r.judgePrompt(
		evalCase{Group: "A", ID: "A01", Statement: "x", Expected: "debugging"},
		[]string{"resp"}, "")
	for _, want := range []string{
		"src/orders.js",     // fixture manifest inlined (grounding — #8)
		"security-guidance", // authoritative plugin list inlined (fabrication — #6)
		"COMPLETE list",     // plugin list framed as exhaustive
		"is a fabrication",  // fabrication instruction present
		"/verify",           // known-real commands listed
	} {
		if !strings.Contains(jp, want) {
			t.Errorf("judge prompt missing %q", want)
		}
	}
}

func TestHookedFixtureForB04(t *testing.T) {
	r := newTestRunner(t)
	if err := os.WriteFile(filepath.Join(r.fixture, "package.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dir, err := r.hookedFixture()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	if dir == r.fixture {
		t.Fatal("B04 must run in a copy, not the shared fixture")
	}
	if readFile(filepath.Join(dir, "package.json")) == "" {
		t.Error("fixture contents not copied")
	}
	if !strings.Contains(readFile(filepath.Join(dir, ".claude", "settings.json")), `"hooks"`) {
		t.Error("hooks settings not written into the copy")
	}
	if _, err := os.Stat(filepath.Join(r.fixture, ".claude")); err == nil {
		t.Error("the shared fixture must stay untouched")
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
