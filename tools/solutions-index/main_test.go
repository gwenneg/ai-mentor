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

func recordMD(kind, goals, bestWhen, sessionSig string) string {
	return "# record\n*Last verified: 2026-07-03*\n\n" +
		"kind: " + kind + "\ngoals: " + goals + "\nbest_when: " + bestWhen + "\nsession_signal: " + sessionSig + "\n"
}

func validTree() map[string]string {
	return map[string]string{
		"skills/mentor/problems/test-goal.md": `# test-goal
*Last verified: 2026-07-03*

| # | Solution | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Alpha](../solutions/alpha.md) | Alpha shines | y |
| 2 | [Beta](../solutions/beta.md) | Beta fits | y |
`,
		"skills/mentor/problems/other-goal.md": `# other-goal
*Last verified: 2026-07-03*

| # | Solution | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Beta](../solutions/beta.md) | Beta wins here | y |
`,
		"skills/mentor/solutions/alpha.md":            techniqueMD("`x` exists", "uses alpha"),
		"skills/mentor/solutions/beta.md":             techniqueMD("—", "uses beta"),
		"skills/mentor/solutions/solo.md":             recordMD("builtin-command", "test-goal", "solo fits", "ran /solo"),
		"skills/mentor/solutions/some-integration.md": recordMD("integration", "test-goal, other-goal", "integrating", "repo uses it"),
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
		// technique rows: goals from problems membership, best_when from the best-ranked row
		"| alpha | technique | test-goal | alpha shines | `x` exists | uses alpha |",
		"| beta | technique | other-goal, test-goal | beta wins here | — | uses beta |",
		// record rows: everything from the file's own fields, id = filename
		"| solo | builtin-command | test-goal | solo fits | — | ran /solo |",
		"| some-integration | integration | other-goal, test-goal | integrating | — | repo uses it |",
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
			f["skills/mentor/solutions/orphan.md"] = techniqueMD("—", "sig")
		}, "no ranked row"},
		{"missing signals section", func(f map[string]string) {
			f["skills/mentor/solutions/alpha.md"] = strings.Replace(
				f["skills/mentor/solutions/alpha.md"], "## Signals", "## Whatever", 1)
		}, "missing or incomplete '## Signals'"},
		{"incomplete signals section", func(f map[string]string) {
			f["skills/mentor/solutions/alpha.md"] = strings.Replace(
				f["skills/mentor/solutions/alpha.md"], "- Session: uses alpha\n", "", 1)
		}, "missing or incomplete '## Signals'"},
		{"ranked row to missing solution", func(f map[string]string) {
			delete(f, "skills/mentor/solutions/alpha.md")
		}, "solutions/alpha.md, which does not exist"},
		{"ranked row ranks a record", func(f map[string]string) {
			f["skills/mentor/problems/test-goal.md"] = strings.Replace(
				f["skills/mentor/problems/test-goal.md"],
				"| 2 | [Beta](../solutions/beta.md) | Beta fits | y |",
				"| 2 | [Solo](../solutions/solo.md) | Solo here | y |", 1)
			delete(f, "skills/mentor/solutions/beta.md") // beta would otherwise be unrouted
			delete(f, "skills/mentor/problems/other-goal.md")
		}, "ranks 'solo', a record"},
		{"record missing a field", func(f map[string]string) {
			f["skills/mentor/solutions/some-integration.md"] = strings.Replace(
				f["skills/mentor/solutions/some-integration.md"], "session_signal: repo uses it\n", "", 1)
		}, "missing one of goals/best_when/session_signal"},
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
