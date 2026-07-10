package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func approachMD() string {
	var b strings.Builder
	b.WriteString("# Solution\n*Last verified: 2026-07-03*\n\n")
	for _, s := range approachSections[:len(approachSections)-2] {
		b.WriteString(s + "\n\nfiller\n\n")
	}
	for range 30 {
		b.WriteString("filler\n")
	}
	b.WriteString("## Sources\n\n- [Doc](https://example.com/doc)\n\n")
	b.WriteString("## Signals\n\n- Setup: —\n- Session: some signal\n")
	return b.String()
}

// recordMD has no goals/best_when — every record derives both from its
// ranked rows (the generator enforces that; the audit only needs kind).
func recordMD(kind string) string {
	return "# record\n*Last verified: 2026-07-03*\n\nkind: " + kind + "\nsession_signal: seen\n"
}

func validTree() map[string]string {
	return map[string]string{
		"skills/mentor/problems/test-goal.md": `# test-goal
*Last verified: 2026-07-03*

**Hidden gem:** Alpha — because.

| # | Solution | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Alpha](../solutions/alpha.md) | x | y |
| 2 | [Beta](../solutions/beta.md) | x | y |
| 3 | [Gamma](../solutions/gamma.md) | x | y |
| 4 | [shiny-plugin](../solutions/shiny-plugin.md) | x | y |
| 5 | [some-integration](../solutions/some-integration.md) | x | y |
`,
		"skills/mentor/solutions/alpha.md":            approachMD(),
		"skills/mentor/solutions/beta.md":             approachMD(),
		"skills/mentor/solutions/gamma.md":            approachMD(),
		"skills/mentor/solutions/some-integration.md": recordMD("integration"),
		"skills/mentor/solutions/shiny-plugin.md":     recordMD("plugin"),
		"skills/mentor/processed-changelogs.md": `# Ledger
*Updated: 2026-07-03*

| Week | Processed | Outcome |
|------|-----------|---------|
| [2026-w26](https://example.com) | 2026-07-01 | processed |
`,
		"skills/mentor/marketplace.md": `# Catalog
*Last synced: 2026-07-03*

| ` + "`alpha-tool`" + ` | does a thing | ` + "`test-goal`" + ` | ☑️ desk-checked |
`,
		"skills/mentor/profile-schema.md": "# Profile schema\n",
		"skills/mentor/problem-mode.md": `# Problem mode

Full schema: ` + "`profile-schema.md`" + `

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
	a := &auditor{root: repo, skill: filepath.Join(repo, skillDir)}
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
	routing := "skills/mentor/problems/test-goal.md"
	cases := []struct {
		name   string
		mutate func(f map[string]string)
		expect string
	}{
		{"non-sequential rows", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "| 3 |", "| 5 |", 1)
		}, "row numbering not sequential"},
		{"fewer than 3 rows", func(f map[string]string) {
			i := strings.Index(f[routing], "| 3 |")
			f[routing] = f[routing][:i]
			delete(f, "skills/mentor/solutions/gamma.md")
		}, "only 2 rows"},
		{"missing marketplace directory", func(f map[string]string) {
			delete(f, "skills/mentor/marketplace.md")
		}, "missing marketplace directory"},
		{"missing hidden gem", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "**Hidden gem:** Alpha — because.\n", "", 1)
		}, "missing Hidden gem line"},
		{"gem names unranked solution", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "**Hidden gem:** Alpha", "**Hidden gem:** Omega", 1)
		}, "does not match any ranked row"},
		{"broken reference", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "solutions/beta.md", "solutions/missing.md", 1)
		}, "broken reference solutions/missing.md"},
		{"orphan technique", func(f map[string]string) {
			f["skills/mentor/solutions/orphan.md"] = approachMD()
		}, "orphan: not ranked by any problems file"},
		{"bad routing date line", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "*Last verified: 2026-07-03*", "verified recently", 1)
		}, "line 2 must be"},
		{"bad record date line", func(f map[string]string) {
			f["skills/mentor/solutions/some-integration.md"] = strings.Replace(
				f["skills/mentor/solutions/some-integration.md"], "*Last verified: 2026-07-03*", "recently", 1)
		}, "line 2 must be"},
		{"missing approach section", func(f map[string]string) {
			f["skills/mentor/solutions/alpha.md"] = strings.Replace(
				f["skills/mentor/solutions/alpha.md"], "## Common Pitfalls", "## Pitfalls", 1)
		}, "missing section '## Common Pitfalls'"},
		{"sections out of order", func(f map[string]string) {
			a := f["skills/mentor/solutions/alpha.md"]
			a = strings.Replace(a, "## Why It Works\n\nfiller\n\n", "", 1)
			f["skills/mentor/solutions/alpha.md"] = a + "## Why It Works\n\nfiller\n"
		}, "out of order"},
		{"approach too short", func(f map[string]string) {
			f["skills/mentor/solutions/alpha.md"] = strings.ReplaceAll(
				f["skills/mentor/solutions/alpha.md"], "filler\n", "")
		}, "(expected at least 40)"},
		{"optional example section out of order", func(f map[string]string) {
			f["skills/mentor/solutions/alpha.md"] = strings.Replace(
				f["skills/mentor/solutions/alpha.md"], "## Common Pitfalls",
				"## Real-World Example\n\nfiller\n\n## Common Pitfalls", 1)
		}, "section '## Real-World Example' out of order"},
		{"no sources", func(f map[string]string) {
			f["skills/mentor/solutions/alpha.md"] = strings.Replace(
				f["skills/mentor/solutions/alpha.md"], "- [Doc](https://example.com/doc)\n", "", 1)
		}, "0 Sources entries"},
		{"approach missing signals section", func(f map[string]string) {
			f["skills/mentor/solutions/alpha.md"] = strings.Replace(
				f["skills/mentor/solutions/alpha.md"], "## Signals", "## Adoption evidence", 1)
		}, "missing section '## Signals'"},
		{"bad week slug", func(f map[string]string) {
			f["skills/mentor/processed-changelogs.md"] = strings.Replace(
				f["skills/mentor/processed-changelogs.md"], "[2026-w26]", "[week-26]", 1)
		}, "not a week slug"},
		{"duplicate ledger row", func(f map[string]string) {
			f["skills/mentor/processed-changelogs.md"] +=
				"| [2026-w26](https://example.com) | 2026-07-02 | again |\n"
		}, "duplicate ledger row"},
		{"invalid ledger date", func(f map[string]string) {
			f["skills/mentor/processed-changelogs.md"] = strings.Replace(
				f["skills/mentor/processed-changelogs.md"], "2026-07-01", "yesterday", 1)
		}, "invalid processed date"},
		{"impossible ledger date", func(f map[string]string) {
			f["skills/mentor/processed-changelogs.md"] = strings.Replace(
				f["skills/mentor/processed-changelogs.md"], "2026-07-01", "2026-99-99", 1)
		}, "invalid processed date"},
		{"malformed short ledger row", func(f map[string]string) {
			f["skills/mentor/processed-changelogs.md"] +=
				"| [2026-w27](https://example.com) | 2026-07-02 |\n"
		}, "is malformed"},
		{"impossible routing date", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "*Last verified: 2026-07-03*", "*Last verified: 2026-99-99*", 1)
		}, "line 2 must be"},
		{"broken doc reference in a mode file", func(f map[string]string) {
			delete(f, "skills/mentor/profile-schema.md")
		}, "broken reference profile-schema.md"},
		{"empty ledger outcome", func(f map[string]string) {
			f["skills/mentor/processed-changelogs.md"] = strings.Replace(
				f["skills/mentor/processed-changelogs.md"], "| processed |", "|  |", 1)
		}, "empty outcome"},
		{"goal missing from classification table", func(f map[string]string) {
			f["skills/mentor/problem-mode.md"] = strings.Replace(
				f["skills/mentor/problem-mode.md"], "`test-goal`", "`other-goal`", 1)
		}, "missing from the problem-mode classification table"},
		{"stale classification row", func(f map[string]string) {
			f["skills/mentor/problem-mode.md"] = strings.Replace(
				f["skills/mentor/problem-mode.md"], "`test-goal`", "`other-goal`", 1)
		}, "has no matching routing section"},
		{"wrong goal count in prose", func(f map[string]string) {
			f["skills/mentor/problem-mode.md"] = strings.Replace(
				f["skills/mentor/problem-mode.md"], "1 goal categories", "24 goal categories", 1)
		}, "prose says '24 goal categories' but there are 1"},
		{"missing problems directory", func(f map[string]string) {
			delete(f, routing)
		}, "missing problems directory"},
		{"built-ins line reintroduced", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "| # |",
				"**Built-ins:** `/testcmd` — does the thing.\n\n| # |", 1)
		}, "capability line found"},
		{"plugins line reintroduced", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "| # |",
				"**Plugins:** `alpha-tool` ☑️ something useful.\n\n| # |", 1)
		}, "capability line found"},
		{"orphan integration record", func(f map[string]string) {
			f["skills/mentor/solutions/unrouted.md"] = recordMD("integration")
		}, "orphan: not ranked by any problems file"},
		{"orphan plugin record", func(f map[string]string) {
			f["skills/mentor/solutions/lonely-plugin.md"] = recordMD("plugin")
		}, "orphan: not ranked by any problems file"},
		{"promoted plugin still in the directory", func(f map[string]string) {
			f["skills/mentor/marketplace.md"] += "| `shiny-plugin` | dup row | `test-goal` | ☑️ desk-checked |\n"
		}, "promoted plugin still has a marketplace.md row"},
		{"unknown record kind", func(f map[string]string) {
			f["skills/mentor/solutions/some-integration.md"] = strings.Replace(
				f["skills/mentor/solutions/some-integration.md"], "kind: integration", "kind: gadget", 1)
		}, "unknown kind 'gadget'"},
		{"missing ledger", func(f map[string]string) {
			delete(f, "skills/mentor/processed-changelogs.md")
		}, "missing processed-changelog ledger"},
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

func TestPluginNamesIgnoresGoalSlugs(t *testing.T) {
	catalog := strings.Join([]string{
		"### Dev workflow",
		"| `real-plugin` | desc | `debugging` | ✅ |",
		"### Language servers (LSPs)",
		"Family: `gopls-lsp` (Go).",
	}, "\n")
	got := strings.Join(pluginNames(catalog), ",")
	if got != "real-plugin,gopls-lsp" {
		t.Errorf("pluginNames = %q, want \"real-plugin,gopls-lsp\" (goal slug 'debugging' must be excluded)", got)
	}
}

func TestEmptySolutionsDirIsFatal(t *testing.T) {
	repo := t.TempDir()
	a := &auditor{root: repo, skill: filepath.Join(repo, skillDir)}
	if err := a.run(); err == nil {
		t.Error("empty solutions directory should be fatal")
	}
}
