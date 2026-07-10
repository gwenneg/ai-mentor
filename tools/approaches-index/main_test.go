package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func techniqueMD(setupSig, sessionSig string) string {
	return "# Solution\n*Last verified: 2026-07-03*\n\nfiller\n\n" +
		"## Signals\n\n- Setup: " + setupSig + "\n- Session: " + sessionSig + "\n"
}

// recordMD has no goals/best_when — every record derives both from its
// ranked rows; inline copies are an error.
func recordMD(kind, sessionSig string) string {
	return "---\nkind: " + kind + "\nlast_verified: 2026-07-03\n" +
		"session_signal: \"" + sessionSig + "\"\n---\n"
}

func validTree() map[string]string {
	return map[string]string{
		"skills/mentor/playbooks/test-goal.md": `# test-goal
*Last verified: 2026-07-03*

| # | Solution | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Alpha](../approaches/techniques/alpha.md) | Alpha shines | y |
| 2 | [Beta](../approaches/techniques/beta.md) | Beta fits | y |
| 3 | [neat-plugin](../approaches/records/neat-plugin.md) | Plugin shines | y |
`,
		"skills/mentor/playbooks/other-goal.md": `# other-goal
*Last verified: 2026-07-03*

| # | Solution | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Beta](../approaches/techniques/beta.md) | Beta wins here | y |
| 2 | [some-integration](../approaches/records/some-integration.md) | Integrating | y |
`,
		"skills/mentor/approaches/techniques/alpha.md":         techniqueMD("`x` exists", "uses alpha"),
		"skills/mentor/approaches/techniques/beta.md":          techniqueMD("—", "uses beta"),
		"skills/mentor/approaches/records/neat-plugin.md":      recordMD("plugin", "neat-plugin installed"),
		"skills/mentor/approaches/records/some-integration.md": recordMD("integration", "repo uses it"),
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
		// every ranked solution: goals from membership, best_when from the best-ranked row
		"| alpha | technique | test-goal | alpha shines | `x` exists | uses alpha |",
		"| beta | technique | other-goal, test-goal | beta wins here | — | uses beta |",
		"| neat-plugin | plugin | test-goal | plugin shines | — | neat-plugin installed |",
		"| some-integration | integration | other-goal | integrating | — | repo uses it |",
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
		{"technique without ranked row", func(f map[string]string) {
			f["skills/mentor/approaches/techniques/orphan.md"] = techniqueMD("—", "sig")
		}, "'orphan' (technique) has no ranked row"},
		{"missing signals section", func(f map[string]string) {
			f["skills/mentor/approaches/techniques/alpha.md"] = strings.Replace(
				f["skills/mentor/approaches/techniques/alpha.md"], "## Signals", "## Whatever", 1)
		}, "missing or incomplete '## Signals'"},
		{"incomplete signals section", func(f map[string]string) {
			f["skills/mentor/approaches/techniques/alpha.md"] = strings.Replace(
				f["skills/mentor/approaches/techniques/alpha.md"], "- Session: uses alpha\n", "", 1)
		}, "missing or incomplete '## Signals'"},
		{"ranked row to missing solution", func(f map[string]string) {
			delete(f, "skills/mentor/approaches/techniques/alpha.md")
		}, "missing approach file for 'alpha'"},
		{"record with inline goals", func(f map[string]string) {
			f["skills/mentor/approaches/records/neat-plugin.md"] = strings.Replace(
				f["skills/mentor/approaches/records/neat-plugin.md"],
				"kind: plugin\n", "kind: plugin\ngoals: test-goal\n", 1)
		}, "carries inline goals:/best_when:"},
		{"record with inline best_when", func(f map[string]string) {
			f["skills/mentor/approaches/records/some-integration.md"] = strings.Replace(
				f["skills/mentor/approaches/records/some-integration.md"],
				"kind: integration\n", "kind: integration\nbest_when: something\n", 1)
		}, "carries inline goals:/best_when:"},
		{"plugin record not ranked", func(f map[string]string) {
			f["skills/mentor/playbooks/test-goal.md"] = strings.Replace(
				f["skills/mentor/playbooks/test-goal.md"],
				"| 3 | [neat-plugin](../approaches/records/neat-plugin.md) | Plugin shines | y |\n", "", 1)
		}, "'neat-plugin' (plugin) has no ranked row"},
		{"record missing session_signal", func(f map[string]string) {
			f["skills/mentor/approaches/records/some-integration.md"] = strings.Replace(
				f["skills/mentor/approaches/records/some-integration.md"], "session_signal: \"repo uses it\"\n", "", 1)
		}, "record is missing session_signal"},
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
