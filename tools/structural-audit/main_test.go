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

func validTree() map[string]string {
	return map[string]string{
		"skills/mentor/routing/test-goal.md": `# test-goal
*Last verified: 2026-07-03*

**Hidden gem:** Alpha — because.

**Plugins:** ` + "`alpha-tool`" + ` ☑️ something useful.

**Built-ins:** ` + "`/testcmd`" + ` — does the thing.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Alpha](../approaches/alpha.md) | Beginner | x | y |
| 2 | [Beta](../approaches/beta.md) | Intermediate | x | y |
| 3 | [Gamma](../approaches/gamma.md) | Advanced | x | y |
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
		"skills/mentor/registry/integrations.md": `# Integrations
*Last verified: 2026-07-03*

## some-integration

id: some-integration
kind: integration
goals: test-goal
`,
		"skills/mentor/registry/builtin-commands.md": `# Registry
*Last verified: 2026-07-03*

## testcmd

id: testcmd
kind: builtin-command
goals: test-goal
best_when: testing the audit
`,
		"skills/mentor/references/official-plugins.md": `# Catalog
*Last synced: 2026-07-03*

| ` + "`alpha-tool`" + ` | does a thing | ` + "`test-goal`" + ` | ☑️ desk-checked |
`,
		"skills/mentor/references/profile-schema.md": "# Profile schema\n",
		"skills/mentor/problem-mode.md": `# Problem mode

Full schema: ` + "`references/profile-schema.md`" + `

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
	routing := "skills/mentor/routing/test-goal.md"
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
		}, "only 2 rows"},
		{"missing plugins line", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "**Plugins:** `alpha-tool` ☑️ something useful.\n\n", "", 1)
		}, "missing Plugins line"},
		{"phantom plugin token", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "`alpha-tool` ☑️ something useful", "`ghost-tool` ☑️ something useful", 1)
		}, "Plugins line names 'ghost-tool', not found"},
		{"plugins line names a goal slug (not a plugin)", func(f map[string]string) {
			// `test-goal` is backticked in the catalog's goal column but is not a
			// plugin — the old all-tokens parser let this pass.
			f[routing] = strings.Replace(f[routing], "`alpha-tool` ☑️ something useful", "`test-goal` ☑️ something useful", 1)
		}, "Plugins line names 'test-goal', not found"},
		{"missing plugin catalog", func(f map[string]string) {
			delete(f, "skills/mentor/references/official-plugins.md")
		}, "missing official-plugins catalog"},
		{"missing hidden gem", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "**Hidden gem:** Alpha — because.\n", "", 1)
		}, "missing Hidden gem line"},
		{"gem names unranked approach", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "**Hidden gem:** Alpha", "**Hidden gem:** Omega", 1)
		}, "does not match any ranked row"},
		{"broken reference", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "approaches/beta.md", "approaches/missing.md", 1)
		}, "broken reference approaches/missing.md"},
		{"orphan approach", func(f map[string]string) {
			f["skills/mentor/approaches/orphan.md"] = approachMD()
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
			f["skills/mentor/approaches/alpha.md"] = strings.ReplaceAll(
				f["skills/mentor/approaches/alpha.md"], "filler\n", "")
		}, "(expected at least 40)"},
		{"optional example section out of order", func(f map[string]string) {
			f["skills/mentor/approaches/alpha.md"] = strings.Replace(
				f["skills/mentor/approaches/alpha.md"], "## Common Pitfalls",
				"## Real-World Example\n\nfiller\n\n## Common Pitfalls", 1)
		}, "section '## Real-World Example' out of order"},
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
		{"impossible ledger date", func(f map[string]string) {
			f["skills/mentor/references/processed-changelogs.md"] = strings.Replace(
				f["skills/mentor/references/processed-changelogs.md"], "2026-07-01", "2026-99-99", 1)
		}, "invalid processed date"},
		{"malformed short ledger row", func(f map[string]string) {
			f["skills/mentor/references/processed-changelogs.md"] +=
				"| [2026-w27](https://example.com) | 2026-07-02 |\n"
		}, "is malformed"},
		{"impossible routing date", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "*Last verified: 2026-07-03*", "*Last verified: 2026-99-99*", 1)
		}, "line 2 must be"},
		{"broken doc reference in a mode file", func(f map[string]string) {
			delete(f, "skills/mentor/references/profile-schema.md")
		}, "broken reference references/profile-schema.md"},
		{"empty ledger outcome", func(f map[string]string) {
			f["skills/mentor/references/processed-changelogs.md"] = strings.Replace(
				f["skills/mentor/references/processed-changelogs.md"], "| processed |", "|  |", 1)
		}, "empty outcome"},
		{"approach missing signals section", func(f map[string]string) {
			f["skills/mentor/approaches/alpha.md"] = strings.Replace(
				f["skills/mentor/approaches/alpha.md"], "## Signals", "## Adoption evidence", 1)
		}, "missing section '## Signals'"},
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
		{"missing routing directory", func(f map[string]string) {
			delete(f, routing)
		}, "missing routing directory"},
		{"phantom builtin token", func(f map[string]string) {
			f[routing] = strings.Replace(f[routing], "`/testcmd`", "`/ghostcmd`", 1)
		}, "Built-ins line names '/ghostcmd', not found"},
		{"orphan registry id", func(f map[string]string) {
			f["skills/mentor/registry/builtin-commands.md"] += "\n## other\n\nid: othercmd\nkind: builtin-command\ngoals: test-goal\n"
		}, "registry id 'othercmd' not referenced"},
		{"duplicate registry id", func(f map[string]string) {
			f["skills/mentor/registry/builtin-commands.md"] += "\nid: testcmd\n"
		}, "duplicate registry id 'testcmd'"},
		{"registry goals name unknown goal", func(f map[string]string) {
			f["skills/mentor/registry/builtin-commands.md"] = strings.Replace(
				f["skills/mentor/registry/builtin-commands.md"], "goals: test-goal", "goals: no-such-goal", 1)
		}, "registry goals name 'no-such-goal'"},
		{"integration id collides with builtin", func(f map[string]string) {
			f["skills/mentor/registry/integrations.md"] = strings.Replace(
				f["skills/mentor/registry/integrations.md"], "id: some-integration", "id: testcmd", 1)
		}, "duplicate registry id 'testcmd'"},
		{"missing integrations registry", func(f map[string]string) {
			delete(f, "skills/mentor/registry/integrations.md")
		}, "missing integrations registry"},
		{"missing registry", func(f map[string]string) {
			delete(f, "skills/mentor/registry/builtin-commands.md")
			f[routing] = strings.Replace(f[routing], "\n**Built-ins:** `/testcmd` — does the thing.\n", "", 1)
		}, "missing builtin-commands registry"},
		{"missing ledger", func(f map[string]string) {
			delete(f, "skills/mentor/references/processed-changelogs.md")
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

func TestEmptyApproachDirIsFatal(t *testing.T) {
	repo := t.TempDir()
	a := &auditor{root: repo, skill: filepath.Join(repo, skillDir)}
	if err := a.run(); err == nil {
		t.Error("empty approach directory should be fatal")
	}
}
