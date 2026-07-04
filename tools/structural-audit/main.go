// Deterministic structural audit for the ai-mentor catalog.
//
// Checks goal routing, approach files, cross-references, the changelog
// ledger, adoption signals, and SKILL.md consistency. Exits 1 if any issue
// is found, 2 on a fatal setup problem. No network, no LLM — safe as a PR
// gate. Stdlib only.
//
// Usage: go run . <repo-root>
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	reRow      = regexp.MustCompile(`^\| (\d+) \| \[([^\]]+)\]`)
	reGem      = regexp.MustCompile(`^\*\*Hidden gem:\*\* ([^—\n]+)`)
	reRef      = regexp.MustCompile(`approaches/[a-z0-9-]+\.md`)
	reSource   = regexp.MustCompile(`^- \[[^\]]+\]\(https?://`)
	reLedger   = regexp.MustCompile(`^\| *\[([^\]]+)\]`)
	reWeek     = regexp.MustCompile(`^\d{4}-w\d{2}$`)
	reDate     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	reSignal   = regexp.MustCompile("^\\| ([a-z0-9-]+) \\|")
	reSkillRow = regexp.MustCompile("^\\| `([a-z0-9-]+)` \\|")
	reCount    = regexp.MustCompile(`(\d+) goal categories`)
)

var approachSections = []string{
	"## What It Is", "## Why It Works", "## When to Use It", "## When NOT to Use It",
	"## How It Works", "### Basic (Beginner)",
	"### Composing with Other Approaches (Intermediate)", "### Advanced Patterns",
	"## Common Pitfalls", "## Real-World Example", "## Sources",
}

type auditor struct {
	skill                    string // skills/mentor directory
	issues                   []string
	goals, approaches, weeks int
}

func (a *auditor) issue(path, format string, args ...any) {
	rel := path
	if r, err := filepath.Rel(filepath.Dir(filepath.Dir(a.skill)), path); err == nil {
		rel = r
	}
	a.issues = append(a.issues, rel+": "+fmt.Sprintf(format, args...))
}

func lines(path string) []string {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return strings.Split(string(b), "\n")
}

// dateLine checks that line 2 is e.g. "*Last verified: 2026-07-03*".
func (a *auditor) dateLine(path, label string, ls []string) {
	re := regexp.MustCompile(`^\*` + label + `: \d{4}-\d{2}-\d{2}\*`)
	if len(ls) < 2 || !re.MatchString(ls[1]) {
		a.issue(path, "line 2 must be '*%s: YYYY-MM-DD*'", label)
	}
}

func (a *auditor) checkRouting(path string, approachNames []string) []string {
	ls := lines(path)
	if ls == nil {
		a.issue(path, "missing routing table")
		return nil
	}
	a.dateLine(path, "Last verified", ls)

	var goals []string
	section, rows, gem := "", []string{}, ""
	flush := func() {
		if section == "" || section == "extraction-notes" {
			return
		}
		goals = append(goals, section)
		if len(rows) < 3 {
			a.issue(path, "section %s: only %d rows (expected at least 3)", section, len(rows))
		}
		if gem == "" {
			a.issue(path, "section %s: missing Hidden gem line", section)
		} else {
			g := strings.ToLower(strings.TrimSpace(gem))
			ok := false
			for _, r := range rows {
				rl := strings.ToLower(r)
				if strings.Contains(g, rl) || strings.Contains(rl, g) {
					ok = true
				}
			}
			if !ok {
				a.issue(path, "section %s: Hidden gem '%s' does not match any ranked row", section, strings.TrimSpace(gem))
			}
		}
	}
	for _, l := range ls {
		if strings.HasPrefix(l, "## ") {
			flush()
			section, rows, gem = l[3:], nil, ""
			continue
		}
		if m := reGem.FindStringSubmatch(l); m != nil {
			gem = m[1]
		}
		if m := reRow.FindStringSubmatch(l); m != nil {
			if n, _ := strconv.Atoi(m[1]); n != len(rows)+1 {
				a.issue(path, "section %s: row numbering not sequential at row %d", section, len(rows)+1)
			}
			rows = append(rows, m[2])
			cells := strings.Split(l, "|")
			if len(cells) > 3 {
				switch lvl := strings.TrimSpace(cells[3]); lvl {
				case "Beginner", "Intermediate", "Advanced":
				default:
					a.issue(path, "invalid level %s", lvl)
				}
			}
		}
	}
	flush()

	text := strings.Join(ls, "\n")
	for _, ref := range dedup(reRef.FindAllString(text, -1)) {
		if _, err := os.Stat(filepath.Join(a.skill, ref)); err != nil {
			a.issue(path, "broken reference %s", ref)
		}
	}
	for _, name := range approachNames {
		if !strings.Contains(text, "approaches/"+name+".md") {
			a.issue(filepath.Join(a.skill, "approaches", name+".md"), "orphan: not referenced by the routing table")
		}
	}
	return goals
}

func (a *auditor) checkApproach(path string) {
	ls := lines(path)
	a.dateLine(path, "Last verified", ls)

	pos := 0
	for _, s := range approachSections {
		ln := 0
		for i, l := range ls {
			if strings.Contains(l, s) {
				ln = i + 1
				break
			}
		}
		switch {
		case ln == 0:
			a.issue(path, "missing section '%s'", s)
		case ln < pos:
			a.issue(path, "section '%s' out of order", s)
		default:
			pos = ln
		}
	}

	if n := len(ls) - 1; n < 60 { // trailing newline yields one empty element
		a.issue(path, "%d lines (expected at least 60)", n)
	}
	srcs, in := 0, false
	for _, l := range ls {
		if l == "## Sources" {
			in = true
			continue
		}
		if in && reSource.MatchString(l) {
			srcs++
		}
	}
	if srcs < 1 {
		a.issue(path, "%d Sources entries (expected at least 1)", srcs)
	}
}

func (a *auditor) checkLedger(path string) {
	ls := lines(path)
	if ls == nil {
		a.issue(path, "missing processed-changelog ledger")
		return
	}
	a.dateLine(path, "Updated", ls)

	seen := map[string]bool{}
	for _, l := range ls {
		m := reLedger.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		a.weeks++
		slug := m[1]
		if !reWeek.MatchString(slug) {
			a.issue(path, "row '%s' is not a week slug like 2026-w26", slug)
		}
		if seen[slug] {
			a.issue(path, "duplicate ledger row for '%s'", slug)
		}
		seen[slug] = true
		cells := strings.Split(l, "|")
		if len(cells) > 3 {
			if !reDate.MatchString(strings.TrimSpace(cells[2])) {
				a.issue(path, "row '%s' has invalid processed date '%s'", slug, strings.TrimSpace(cells[2]))
			}
			if strings.TrimSpace(cells[3]) == "" {
				a.issue(path, "row '%s' has an empty outcome", slug)
			}
		}
	}
}

func (a *auditor) checkSignals(path string, approachNames []string) {
	ls := lines(path)
	if ls == nil {
		a.issue(path, "missing adoption-signals table")
		return
	}
	a.dateLine(path, "Last reviewed", ls)

	var names []string
	for _, l := range ls {
		if m := reSignal.FindStringSubmatch(l); m != nil {
			names = append(names, m[1])
		}
	}
	missing, stale := diff(approachNames, names)
	for _, x := range missing {
		a.issue(path, "approach '%s' has no adoption-signals row", x)
	}
	for _, x := range stale {
		a.issue(path, "row '%s' has no matching approach file", x)
	}
	for _, d := range dups(names) {
		a.issue(path, "duplicate signals row for '%s'", d)
	}
}

func (a *auditor) checkSkillMD(path string, goals []string) {
	ls := lines(path)
	var table []string
	for _, l := range ls {
		if m := reSkillRow.FindStringSubmatch(l); m != nil {
			table = append(table, m[1])
		}
	}
	missing, stale := diff(goals, table)
	for _, x := range missing {
		a.issue(path, "routing goal %s missing from the Phase 1 classification table", x)
	}
	for _, x := range stale {
		a.issue(path, "Phase 1 table row %s has no matching routing section", x)
	}
	for _, l := range ls {
		for _, m := range reCount.FindAllStringSubmatch(l, -1) {
			if n, _ := strconv.Atoi(m[1]); n != len(goals) {
				a.issue(path, "prose says '%d goal categories' but there are %d", n, len(goals))
			}
		}
	}
}

func (a *auditor) run() error {
	approachDir := filepath.Join(a.skill, "approaches")
	files, _ := filepath.Glob(filepath.Join(approachDir, "*.md"))
	if len(files) == 0 {
		return fmt.Errorf("approach directory empty/missing")
	}
	a.approaches = len(files)
	var names []string
	for _, f := range files {
		names = append(names, strings.TrimSuffix(filepath.Base(f), ".md"))
	}

	goals := a.checkRouting(filepath.Join(a.skill, "routing.md"), names)
	a.goals = len(goals)
	for _, f := range files {
		a.checkApproach(f)
	}
	a.checkLedger(filepath.Join(a.skill, "references", "processed-changelogs.md"))
	a.checkSignals(filepath.Join(a.skill, "references", "adoption-signals.md"), names)
	a.checkSkillMD(filepath.Join(a.skill, "SKILL.md"), goals)
	return nil
}

func dedup(xs []string) []string {
	seen, out := map[string]bool{}, []string{}
	for _, x := range xs {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}

func dups(xs []string) []string {
	seen, out := map[string]int{}, []string{}
	for _, x := range xs {
		seen[x]++
		if seen[x] == 2 {
			out = append(out, x)
		}
	}
	return out
}

// diff returns (in a but not b, in b but not a), preserving order.
func diff(a, b []string) (onlyA, onlyB []string) {
	inA, inB := map[string]bool{}, map[string]bool{}
	for _, x := range a {
		inA[x] = true
	}
	for _, x := range b {
		inB[x] = true
	}
	for _, x := range a {
		if !inB[x] {
			onlyA = append(onlyA, x)
		}
	}
	for _, x := range dedup(b) {
		if !inA[x] {
			onlyB = append(onlyB, x)
		}
	}
	return
}

// findRoot walks upward from dir to the first directory containing
// skills/mentor, so the audit works from anywhere in the repo — including
// tools/structural-audit itself, where `go -C tools/structural-audit run .` lands.
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

func main() {
	var repo string
	var err error
	if len(os.Args) > 1 {
		repo = os.Args[1]
	} else if repo, err = findRoot("."); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
		os.Exit(2)
	}
	a := &auditor{skill: filepath.Join(repo, "skills", "mentor")}
	if err := a.run(); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
		os.Exit(2)
	}
	for _, is := range a.issues {
		fmt.Printf("  - %s\n", is)
	}
	fmt.Printf("Audited %d routing goals, %d approaches, %d processed changelogs.\n", a.goals, a.approaches, a.weeks)
	if len(a.issues) > 0 {
		fmt.Printf("\n%d issue(s) found (listed above).\n", len(a.issues))
		os.Exit(1)
	}
	fmt.Println("Structural audit: PASS")
}
