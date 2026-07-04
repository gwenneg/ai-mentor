package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func approachMD() string {
	var b strings.Builder
	b.WriteString("# Approach\n*Last verified: 2026-07-03*\n\n")
	for _, s := range approachSections[:len(approachSections)-1] {
		b.WriteString(s + "\n\nfiller\n\n")
	}
	for range 30 {
		b.WriteString("filler\n")
	}
	b.WriteString("## Sources\n\n- [Doc](https://example.com/doc)\n")
	return b.String()
}

func validTree() map[string]string {
	return map[string]string{
		"skills/mentor/routing.md": `# Routing
*Last verified: 2026-07-03*

## test-goal

**Hidden gem:** Alpha — because.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Alpha](approaches/alpha.md) | Beginner | x | y |
| 2 | [Beta](approaches/beta.md) | Intermediate | x | y |
| 3 | [Gamma](approaches/gamma.md) | Advanced | x | y |
`,
		"skills/mentor/approaches/alpha.md": approachMD(),
		"skills/mentor/approaches/beta.md":  approachMD(),
		"skills/mentor/approaches/gamma.md": approachMD(),
		"skills/mentor/references/processed-changelogs.md": `# Ledger
*Updated: 2026-07-03*

| Week | Processed | Outcome |
|------|-----------|---------|
| [2026-w26](https://example.com) | 2026-07-01 | processed |
`,
		"skills/mentor/references/adoption-signals.md": `# Signals
*Last reviewed: 2026-07-03*

| alpha | some signal |
| beta | some signal |
| gamma | some signal |
`,
		"skills/mentor/SKILL.md": `# Skill

| ` + "`test-goal`" + ` | keywords |

There is 1 goal categories here.
`,
	}
}

func runOn(t *testing.T, files map[string]string) []string {
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
	a := &auditor{skill: filepath.Join(repo, "skills", "mentor")}
	if err := a.run(); err != nil {
		t.Fatalf("unexpected fatal: %v", err)
	}
	return a.issues
}

func TestValidTreePasses(t *testing.T) {
	if issues := runOn(t, validTree()); len(issues) != 0 {
		t.Errorf("valid tree should pass, got issues:\n%s", strings.Join(issues, "\n"))
	}
}

// Every corruption must produce at least one issue mentioning the expected
// text — i.e. every check must be able to fail the gate, not just print.
func TestCorruptionsAreCaught(t *testing.T) {
	routing := "skills/mentor/routing.md"
	cases := []struct {
		name   string
		mutate func(f map[string]string)
		expect string
	}{
		{"invalid level", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "| Advanced |", "| Expert |", 1)
		}, "invalid level Expert"},
		{"non-sequential rows", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "| 3 |", "| 5 |", 1)
		}, "row numbering not sequential"},
		{"fewer than 3 rows", func(f map[string]string) {
			i := strings.Index(f[routing], "| 3 |")
			f[routing] = f[routing][:i]
			delete(f, "skills/mentor/approaches/gamma.md")
			f["skills/mentor/references/adoption-signals.md"] = strings.Replace(
				f["skills/mentor/references/adoption-signals.md"], "| gamma | some signal |\n", "", 1)
		}, "only 2 rows"},
		{"missing hidden gem", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "**Hidden gem:** Alpha — because.\n", "", 1)
		}, "missing Hidden gem line"},
		{"gem names unranked approach", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "**Hidden gem:** Alpha", "**Hidden gem:** Omega", 1)
		}, "does not match any ranked row"},
		{"broken reference", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "approaches/beta.md", "approaches/missing.md", 1)
			f["skills/mentor/references/adoption-signals.md"] = strings.Replace(
				f["skills/mentor/references/adoption-signals.md"], "beta", "missing", 1)
		}, "broken reference approaches/missing.md"},
		{"orphan approach", func(f map[string]string) {
			f["skills/mentor/approaches/orphan.md"] = approachMD()
			f["skills/mentor/references/adoption-signals.md"] += "| orphan | some signal |\n"
		}, "orphan: not referenced"},
		{"bad routing date line", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "*Last verified: 2026-07-03*", "verified recently", 1)
		}, "line 2 must be"},
		{"missing approach section", func(f map[string]string) {
			f["skills/mentor/approaches/alpha.md"] = strings.Replace(
				f["skills/mentor/approaches/alpha.md"], "## Common Pitfalls", "## Pitfalls", 1)
		}, "missing section '## Common Pitfalls'"},
		{"sections out of order", func(f map[string]string) {
			a := f["skills/mentor/approaches/alpha.md"]
			a = strings.Replace(a, "## Why It Works\n\nfiller\n\n", "", 1)
			f["skills/mentor/approaches/alpha.md"] = a + "## Why It Works\n\nfiller\n"
		}, "out of order"},
		{"approach too short", func(f map[string]string) {
			f["skills/mentor/approaches/alpha.md"] = strings.Replace(
				f["skills/mentor/approaches/alpha.md"], strings.Repeat("filler\n", 30), "", 1)
		}, "(expected at least 60)"},
		{"no sources", func(f map[string]string) {
			f["skills/mentor/approaches/alpha.md"] = strings.Replace(
				f["skills/mentor/approaches/alpha.md"], "- [Doc](https://example.com/doc)\n", "", 1)
		}, "0 Sources entries"},
		{"bad week slug", func(f map[string]string) {
			f["skills/mentor/references/processed-changelogs.md"] = strings.Replace(
				f["skills/mentor/references/processed-changelogs.md"], "[2026-w26]", "[week-26]", 1)
		}, "not a week slug"},
		{"duplicate ledger row", func(f map[string]string) {
			f["skills/mentor/references/processed-changelogs.md"] +=
				"| [2026-w26](https://example.com) | 2026-07-02 | again |\n"
		}, "duplicate ledger row"},
		{"invalid ledger date", func(f map[string]string) {
			f["skills/mentor/references/processed-changelogs.md"] = strings.Replace(
				f["skills/mentor/references/processed-changelogs.md"], "2026-07-01", "yesterday", 1)
		}, "invalid processed date"},
		{"empty ledger outcome", func(f map[string]string) {
			f["skills/mentor/references/processed-changelogs.md"] = strings.Replace(
				f["skills/mentor/references/processed-changelogs.md"], "| processed |", "|  |", 1)
		}, "empty outcome"},
		{"approach missing signals row", func(f map[string]string) {
			f["skills/mentor/references/adoption-signals.md"] = strings.Replace(
				f["skills/mentor/references/adoption-signals.md"], "| beta | some signal |\n", "", 1)
		}, "'beta' has no adoption-signals row"},
		{"stale signals row", func(f map[string]string) {
			f["skills/mentor/references/adoption-signals.md"] += "| deleted | some signal |\n"
		}, "'deleted' has no matching approach file"},
		{"duplicate signals row", func(f map[string]string) {
			f["skills/mentor/references/adoption-signals.md"] += "| alpha | some signal |\n"
		}, "duplicate signals row"},
		{"goal missing from SKILL.md", func(f map[string]string) {
			f["skills/mentor/SKILL.md"] = strings.Replace(
				f["skills/mentor/SKILL.md"], "`test-goal`", "`other-goal`", 1)
		}, "missing from the Phase 1 classification table"},
		{"stale SKILL.md row", func(f map[string]string) {
			f["skills/mentor/SKILL.md"] = strings.Replace(
				f["skills/mentor/SKILL.md"], "`test-goal`", "`other-goal`", 1)
		}, "has no matching routing section"},
		{"wrong goal count in prose", func(f map[string]string) {
			f["skills/mentor/SKILL.md"] = strings.Replace(
				f["skills/mentor/SKILL.md"], "1 goal categories", "24 goal categories", 1)
		}, "prose says '24 goal categories' but there are 1"},
		{"missing routing table", func(f map[string]string) {
			delete(f, routing)
		}, "missing routing table"},
		{"missing ledger", func(f map[string]string) {
			delete(f, "skills/mentor/references/processed-changelogs.md")
		}, "missing processed-changelog ledger"},
		{"missing signals", func(f map[string]string) {
			delete(f, "skills/mentor/references/adoption-signals.md")
		}, "missing adoption-signals table"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := validTree()
			tc.mutate(files)
			issues := runOn(t, files)
			if len(issues) == 0 {
				t.Fatalf("corruption %q produced no issues — the gate would pass", tc.name)
			}
			for _, is := range issues {
				if strings.Contains(is, tc.expect) {
					return
				}
			}
			t.Errorf("no issue mentions %q; got:\n%s", tc.expect, strings.Join(issues, "\n"))
		})
	}
}

func TestEmptyApproachDirIsFatal(t *testing.T) {
	a := &auditor{skill: filepath.Join(t.TempDir(), "skills", "mentor")}
	if err := a.run(); err == nil {
		t.Error("empty approach directory should be fatal")
	}
}
