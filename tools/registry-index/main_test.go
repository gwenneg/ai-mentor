package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func approachMD(setupSig, sessionSig string) string {
	return "# Approach\n*Last verified: 2026-07-03*\n\nfiller\n\n" +
		"## Signals\n\n- Setup: " + setupSig + "\n- Session: " + sessionSig + "\n"
}

func validTree() map[string]string {
	return map[string]string{
		"skills/mentor/routing/test-goal.md": `# test-goal
*Last verified: 2026-07-03*

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Alpha](../approaches/alpha.md) | none | Alpha shines | y |
| 2 | [Beta](../approaches/beta.md) | some | Beta fits | y |
`,
		"skills/mentor/routing/other-goal.md": `# other-goal
*Last verified: 2026-07-03*

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Beta](../approaches/beta.md) | some | Beta wins here | y |
`,
		"skills/mentor/approaches/alpha.md": approachMD("`x` exists", "uses alpha"),
		"skills/mentor/approaches/beta.md":  approachMD("—", "uses beta"),
		"skills/mentor/registry/builtin-commands.md": `# Registry
*Last verified: 2026-07-03*

## beta

id: beta
kind: builtin-command
goals: test-goal, extra-goal
best_when: record best-when wins
setup: involved
session_signal: ran /beta

## solo

id: solo
kind: builtin-command
goals: test-goal
best_when: solo fits
setup: none
session_signal: ran /solo
`,
		"skills/mentor/registry/integrations.md": `# Integrations
*Last verified: 2026-07-03*

## some-integration

id: some-integration
kind: integration
goals: test-goal
best_when: integrating
setup: some
session_signal: repo uses it
`,
	}
}

func runOn(t *testing.T, files map[string]string) (string, []string) {
	t.Helper()
	repo := t.TempDir()
	for path, content := range files {
		full := filepath.Join(repo, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	g := &gen{root: repo, skill: filepath.Join(repo, skillDir)}
	return g.build(), g.errs
}

func TestValidTreeGenerates(t *testing.T) {
	out, errs := runOn(t, validTree())
	if len(errs) != 0 {
		t.Fatalf("valid tree should generate, got:\n%s", strings.Join(errs, "\n"))
	}
	for _, want := range []string{
		// technique row: goals from routing, best_when from its #1 row, signals from the approach file
		"| alpha | technique | test-goal | alpha shines | none | `x` exists | uses alpha |",
		// merged row: kinds joined, goals unioned+sorted, record fields win, approach setup signal kept
		"| beta | technique + builtin-command | extra-goal, other-goal, test-goal | record best-when wins | involved | — | ran /beta |",
		// plain record rows pass through
		"| solo | builtin-command | test-goal | solo fits | none | — | ran /solo |",
		"| some-integration | integration | test-goal | integrating | some | — | repo uses it |",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing row:\n%s\ngot:\n%s", want, out)
		}
	}
	// deterministic: rows sorted by id
	if strings.Index(out, "| alpha |") > strings.Index(out, "| beta |") {
		t.Error("rows not sorted by id")
	}
}

// Every source inconsistency must fail generation — the -check gate depends on it.
func TestSourceIssuesAreCaught(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(f map[string]string)
		expect string
	}{
		{"conflicting setup values", func(f map[string]string) {
			f["skills/mentor/routing/other-goal.md"] = strings.Replace(
				f["skills/mentor/routing/other-goal.md"], "| some |", "| involved |", 1)
		}, "setup for 'beta' conflicts"},
		{"invalid setup", func(f map[string]string) {
			f["skills/mentor/routing/test-goal.md"] = strings.Replace(
				f["skills/mentor/routing/test-goal.md"], "| none |", "| Beginner |", 1)
		}, "invalid setup 'Beginner'"},
		{"approach without routing row", func(f map[string]string) {
			f["skills/mentor/approaches/orphan.md"] = approachMD("—", "sig")
		}, "no routing row"},
		{"missing signals section", func(f map[string]string) {
			f["skills/mentor/approaches/alpha.md"] = strings.Replace(
				f["skills/mentor/approaches/alpha.md"], "## Signals", "## Whatever", 1)
		}, "missing or incomplete '## Signals'"},
		{"incomplete signals section", func(f map[string]string) {
			f["skills/mentor/approaches/alpha.md"] = strings.Replace(
				f["skills/mentor/approaches/alpha.md"], "- Session: uses alpha\n", "", 1)
		}, "missing or incomplete '## Signals'"},
		{"routing row to missing approach", func(f map[string]string) {
			delete(f, "skills/mentor/approaches/alpha.md")
		}, "approaches/alpha.md, which does not exist"},
		{"record missing a field", func(f map[string]string) {
			f["skills/mentor/registry/integrations.md"] = strings.Replace(
				f["skills/mentor/registry/integrations.md"], "session_signal: repo uses it\n", "", 1)
		}, "missing one of kind/goals/best_when/setup/session_signal"},
		{"missing registry file", func(f map[string]string) {
			delete(f, "skills/mentor/registry/integrations.md")
		}, "missing registry file"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := validTree()
			tc.mutate(files)
			_, errs := runOn(t, files)
			if len(errs) == 0 {
				t.Fatalf("corruption %q produced no errors — the gate would pass", tc.name)
			}
			for _, e := range errs {
				if strings.Contains(e, tc.expect) {
					return
				}
			}
			t.Errorf("no error mentions %q; got:\n%s", tc.expect, strings.Join(errs, "\n"))
		})
	}
}
