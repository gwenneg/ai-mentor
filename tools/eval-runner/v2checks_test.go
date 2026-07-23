package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const v2SpecSnippet = "## Machine expectations (V2 — deterministic layer)\n\n" +
	"| ID | goal | move | surprise | fence | judge |\n" +
	"|----|------|------|----------|-------|-------|\n" +
	"| A01 | debugging | - | required | grounded | - |\n" +
	"| A06 | code-understanding\\|onboarding | - | required | grounded | - |\n" +
	"| A19 | migration | code-modernization | required | portable | - |\n" +
	"| A30 | - | - | - | - | judge |\n" +
	"| C05 | debugging | !plan-mode | required | grounded | - |\n"

func TestParseV2Specs(t *testing.T) {
	specs, err := parseV2Specs(v2SpecSnippet)
	if err != nil {
		t.Fatal(err)
	}
	if got := specs["A06"].Goals; len(got) != 2 || got[0] != "code-understanding" || got[1] != "onboarding" {
		t.Errorf("escaped-pipe alternatives parsed wrong: %v", got)
	}
	if specs["A19"].Move != "code-modernization" || specs["C05"].Move != "!plan-mode" {
		t.Errorf("move column parsed wrong: %+v %+v", specs["A19"], specs["C05"])
	}
	if !specs["A30"].Judge || specs["A01"].Judge {
		t.Error("judge column parsed wrong")
	}
	if specs["A30"].Fence != "" || len(specs["A30"].Goals) != 0 {
		t.Errorf("'-' cells must mean unconstrained: %+v", specs["A30"])
	}
	if _, err := parseV2Specs("no table here"); err == nil {
		t.Error("missing section must be a loud error")
	}
}

// The real cases.md must satisfy the two-way drift guard the runner enforces.
func TestV2SpecsCoverRealSuite(t *testing.T) {
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
	specs, err := parseV2Specs(string(text))
	if err != nil {
		t.Fatal(err)
	}
	for g, cs := range all {
		if g == "D" {
			continue
		}
		for _, c := range cs {
			if _, ok := specs[c.ID]; !ok {
				t.Errorf("case %s has no machine-expectations row", c.ID)
			}
		}
	}
	known := map[string]bool{}
	for _, cs := range all {
		for _, c := range cs {
			known[c.ID] = true
		}
	}
	for id := range specs {
		if !known[id] {
			t.Errorf("machine-expectations row %s names no case", id)
		}
	}
}

// v2Checks2 aliases v2Checks for tests written against the gating return.
func v2Checks2(c evalCase, spec v2Spec, responses []string, plugins, promoted []string) (string, string) {
	return v2Checks(c, spec, responses, plugins, promoted)
}

// adv returns only the advisory (fence-discipline) tier.
func adv(c evalCase, spec v2Spec, responses []string) string {
	_, a := v2Checks(c, spec, responses, nil, nil)
	return a
}

const v2Closing = "\n\nMore options for this — say \"more\". (Calibrated for intermediate.)\n"

func v2resp(trailer, body string) string {
	return body + v2Closing + "\n<!-- mentor " + trailer + " -->"
}

func TestV2ChecksFailPaths(t *testing.T) {
	a := evalCase{Group: "A", ID: "A01"}
	spec := v2Spec{Goals: []string{"debugging"}, Surprise: "required", Fence: "grounded"}
	good := v2resp("mode=problem goal=debugging move=plan-mode surprise=hooks-as-workflow",
		"Diagnosis.\n**One thing you might not know exists:** hooks.\n\n```\n/plan debug orders.go — go test ./...\n```")

	if got, _ := v2Checks2(a, spec, []string{good}, nil, nil); got != "" {
		t.Fatalf("valid response must pass, got %q", got)
	}
	cases := []struct{ name, resp, wantSub string }{
		{"missing trailer", good[:strings.Index(good, "<!--")], "trailer"},
		{"wrong mode", strings.Replace(good, "mode=problem", "mode=growth", 1), "mode"},
		{"wrong goal", strings.Replace(good, "goal=debugging", "goal=testing", 1), "classification"},
		{"no surprise marker", strings.Replace(good, "One thing you might not know exists", "note", 1), "surprise"},
		{"trailer omits surprise", strings.Replace(good, "surprise=hooks-as-workflow", "surprise=omitted", 1), "trailer declares none"},
		{"no closing line", strings.Replace(good, "say \"more\"", "bye", 1), "closing line"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got, _ := v2Checks2(a, spec, []string{tc.resp}, nil, nil); !strings.Contains(got, tc.wantSub) {
				t.Errorf("want failure containing %q, got %q", tc.wantSub, got)
			}
		})
	}
	// Fence discipline is the ADVISORY tier in Phase 2 (true rates too high
	// to gate per-run; Phase 3 rate-gates them) — observed, never gating.
	advCases := []struct{ name, resp, wantSub string }{
		{"no fence", strings.Replace(good, "```", "", 2), "fence"},
		{"ungrounded fence", strings.Replace(good, "orders.go — go test ./...", "your test file", 1), "grounding"},
	}
	for _, tc := range advCases {
		t.Run(tc.name, func(t *testing.T) {
			gating, advisory := v2Checks(a, spec, []string{tc.resp}, nil, nil)
			if gating != "" {
				t.Errorf("discipline finding must not gate, got gating %q", gating)
			}
			if !strings.Contains(advisory, tc.wantSub) {
				t.Errorf("want advisory containing %q, got %q", tc.wantSub, advisory)
			}
		})
	}
}

func TestV2MoveAndPortability(t *testing.T) {
	c := evalCase{Group: "A", ID: "A19"}
	spec := v2Spec{Goals: []string{"migration"}, Move: "code-modernization", Surprise: "required", Fence: "portable"}
	good := v2resp("mode=problem goal=migration move=code-modernization surprise=worktree-isolation",
		"Diagnosis.\n**One thing you might not know exists:** worktrees.\n\n```\nMigrate <your COBOL sources> to Java module by module.\n```")
	if got, _ := v2Checks2(c, spec, []string{good}, nil, nil); got != "" {
		t.Fatalf("valid portable response must pass, got %q", got)
	}
	if got, _ := v2Checks2(c, spec, []string{strings.Replace(good, "move=code-modernization", "move=plan-mode", 1)}, nil, nil); !strings.Contains(got, "move identity") {
		t.Errorf("wrong move must fail, got %q", got)
	}
	leaked := strings.Replace(good, "<your COBOL sources>", "server.go", 1)
	if got := adv(c, spec, []string{leaked}); !strings.Contains(got, "portability") {
		t.Errorf("fixture import into portable fence must be an advisory finding, got %q", got)
	}
	flagged := strings.Replace(good, "<your COBOL sources>", "your sources (adjust: this repo's is server.go)", 1)
	if got := adv(c, spec, []string{flagged}); got != "" {
		t.Errorf("flagged-for-replacement reference must pass (maintainer ruling), got %q", got)
	}
	// forbidden move (C05-style)
	c05 := evalCase{Group: "C", ID: "C05"}
	neg := v2Spec{Move: "!plan-mode", Fence: ""}
	bad := v2resp("mode=problem goal=debugging move=plan-mode surprise=hooks-as-workflow", "x")
	if got, _ := v2Checks2(c05, neg, []string{bad}, nil, nil); !strings.Contains(got, "excluded") {
		t.Errorf("forbidden move must fail, got %q", got)
	}
}

func TestV2SetupScanAndC02(t *testing.T) {
	a09 := evalCase{Group: "A", ID: "A09"}
	setup := v2Spec{Fence: "setup", Surprise: "omitted-ok"}
	good := v2resp("mode=problem goal=incident-response move=mcp-context surprise=omitted",
		"Triage.\n\n```\nclaude mcp add --transport http obs <url> — then correlate the orders handlers.\n```")
	if got, _ := v2Checks2(a09, setup, []string{good}, nil, nil); got != "" {
		t.Fatalf("valid setup fence must pass, got %q", got)
	}
	noCmd := strings.Replace(good, "claude mcp add --transport http obs <url>", "connect your telemetry", 1)
	if got := adv(a09, setup, []string{noCmd}); !strings.Contains(got, "setup command") {
		t.Errorf("setup fence without the command must be advisory, got %q", got)
	}
	noSurface := strings.Replace(good, " — then correlate the orders handlers", "", 1)
	if got := adv(a09, setup, []string{noSurface}); !strings.Contains(got, "surface") {
		t.Errorf("setup fence without the surface must be advisory, got %q", got)
	}

	a22 := evalCase{Group: "A", ID: "A22"}
	scan := v2Spec{Fence: "scan", Surprise: "required"}
	scanGood := v2resp("mode=problem goal=documentation move=custom-skills surprise=hooks-as-workflow",
		"Diagnosis.\n**One thing you might not know exists:** hooks.\n\n```\nDocument GET /orders/total from server.go.\n```")
	if got, _ := v2Checks2(a22, scan, []string{scanGood}, nil, nil); got != "" {
		t.Fatalf("scan fence naming server.go must pass, got %q", got)
	}
	if got := adv(a22, scan, []string{strings.ReplaceAll(scanGood, "server.go", "orders.go")}); !strings.Contains(got, "canary") {
		t.Errorf("scan fence without server.go must be an advisory finding, got %q", got)
	}

	c02 := evalCase{Group: "C", ID: "C02"}
	pspec := v2Spec{Surprise: "required", Fence: "portable"}
	r1 := v2resp("mode=problem goal=refactoring move=subagent-delegation surprise=checkpoints-rewind",
		"D.\n**One thing you might not know exists:** checkpoints.\n\n```\nRefactor <your files>.\n```")
	r2same := v2resp("mode=problem goal=refactoring move=plan-mode surprise=checkpoints-rewind",
		"D.\n**One thing you might not know exists:** checkpoints.\n\n```\nRefactor <your files>.\n```")
	if got, _ := v2Checks2(c02, pspec, []string{r1, r2same}, nil, nil); !strings.Contains(got, "never-repeat") {
		t.Errorf("same surprise across C02 runs must fail, got %q", got)
	}
	r2diff := strings.Replace(r2same, "surprise=checkpoints-rewind", "surprise=autonomous-loops", 1)
	if got, _ := v2Checks2(c02, pspec, []string{r1, r2diff}, nil, nil); got != "" {
		t.Errorf("differing surprises must pass, got %q", got)
	}
}

// Calibration fixes from the first full-suite validation (2026-07-24):
// setup-only fences beside the move fence are tolerated style, and a
// B06 transfer legitimately names the transferred capability.
func TestV2CalibrationFixes(t *testing.T) {
	a19 := evalCase{Group: "A", ID: "A19"}
	spec := v2Spec{Fence: "portable", Surprise: "required"}
	twoFences := v2resp("mode=problem goal=migration move=code-modernization surprise=worktree-isolation",
		"D.\n**One thing you might not know exists:** worktrees.\n\n```\n/plugin install code-modernization@claude-plugins-official\n```\n\n```\nMigrate <your COBOL> module by module.\n```")
	if got, _ := v2Checks2(a19, spec, []string{twoFences}, []string{"code-modernization"}, []string{"code-modernization"}); got != "" {
		t.Errorf("a setup-only fence beside the move fence must pass, got %q", got)
	}
	twoMoves := v2resp("mode=problem goal=migration move=code-modernization surprise=worktree-isolation",
		"D.\n**One thing you might not know exists:** worktrees.\n\n```\nMigrate <your COBOL> plan A.\n```\n\n```\nMigrate <your COBOL> plan B.\n```")
	if got := adv(a19, spec, []string{twoMoves}); !strings.Contains(got, "fence") {
		t.Errorf("two substantive fences must still be flagged (advisory), got %q", got)
	}

	b06 := evalCase{Group: "B", ID: "B06"}
	transfer := "offer\n\n<!-- mentor mode=growth opener=transfer taught=hooks-as-workflow -->"
	if got, _ := v2Checks2(b06, v2Spec{}, []string{transfer}, nil, nil); got != "" {
		t.Errorf("B06 transfer naming the adopted capability must pass, got %q", got)
	}
	emptyTeach := "offer\n\n<!-- mentor mode=growth opener=empty taught=hooks-as-workflow -->"
	if got, _ := v2Checks2(b06, v2Spec{}, []string{emptyTeach}, nil, nil); !strings.Contains(got, "empty-map") {
		t.Errorf("B06 empty-map answer teaching something must fail, got %q", got)
	}
}

func TestV2PluginChecks(t *testing.T) {
	plugins := []string{"prisma", "code-modernization"}
	promoted := []string{"code-modernization"}
	if got := pluginChecks("install with /plugin install ghost-plugin@claude-plugins-official", plugins, promoted); !strings.Contains(got, "fabrication") {
		t.Errorf("unknown plugin must fail, got %q", got)
	}
	unlabeled := "use /plugin install prisma@claude-plugins-official for Postgres."
	if got := pluginChecks(unlabeled, plugins, promoted); !strings.Contains(got, "tier label") {
		t.Errorf("unlabeled directory plugin must fail, got %q", got)
	}
	labeled := "prisma (☑️ desk-checked, not hands-on evaluated): /plugin install prisma@claude-plugins-official"
	if got := pluginChecks(labeled, plugins, promoted); got != "" {
		t.Errorf("labeled directory plugin must pass, got %q", got)
	}
	if got := pluginChecks("/plugin install code-modernization@claude-plugins-official now", plugins, promoted); got != "" {
		t.Errorf("promoted plugin needs no label, got %q", got)
	}
}

func TestV2GrowthOpeners(t *testing.T) {
	b04 := evalCase{Group: "B", ID: "B04"}
	good := "lesson\n\n<!-- mentor mode=growth opener=lesson taught=hooks-as-workflow -->"
	if got, _ := v2Checks2(b04, v2Spec{}, []string{good}, nil, nil); got != "" {
		t.Fatalf("valid B04 must pass, got %q", got)
	}
	wrongTaught := strings.Replace(good, "taught=hooks-as-workflow", "taught=project-memory", 1)
	if got, _ := v2Checks2(b04, v2Spec{}, []string{wrongTaught}, nil, nil); !strings.Contains(got, "hooks-as-workflow") {
		t.Errorf("B04 must teach the configured hook signal, got %q", got)
	}
	b06 := evalCase{Group: "B", ID: "B06"}
	invented := "lesson\n\n<!-- mentor mode=growth opener=lesson taught=hooks-as-workflow -->"
	if got, _ := v2Checks2(b06, v2Spec{}, []string{invented}, nil, nil); !strings.Contains(got, "opener") {
		t.Errorf("B06 must not invent a lesson, got %q", got)
	}
}
