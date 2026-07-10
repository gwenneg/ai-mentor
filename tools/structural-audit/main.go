// Deterministic structural audit for the ai-mentor catalog.
//
// Checks the problems tables, solution files (technique deep-dives and flat
// records), cross-references, the changelog ledger, and SKILL.md consistency.
// Exits 1 if any issue is found, 2 on a fatal setup problem. No network, no
// LLM — safe as a PR gate. Stdlib only.
//
// The compiled index (solutions/index.md) is NOT audited here:
// tools/solutions-index generates it from the same sources and its -check
// mode is the freshness gate in CI.
//
// Usage: go -C tools/structural-audit run . [repo-root]
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	datePat          = `\d{4}-\d{2}-\d{2}`
	minGoalRows      = 3
	minApproachLines = 40
)

var (
	reRow         = regexp.MustCompile(`^\| (\d+) \| \[([^\]]+)\]`)
	reGem         = regexp.MustCompile(`^\*\*Hidden gem:\*\* ([^—\n]+)`)
	rePlugLine    = regexp.MustCompile(`^\*\*Plugins:\*\* `)
	rePlugTok     = regexp.MustCompile("`([a-z0-9.-]+)`")
	reRowName     = regexp.MustCompile("^\\| `([a-z0-9.-]+)`")
	reDocRef      = regexp.MustCompile(`(problems|solutions)/[a-z0-9-]+\.md|\b(plugins|profile-schema|processed-changelogs)\.md`)
	reRef         = regexp.MustCompile(`solutions/[a-z0-9-]+\.md`)
	reSource      = regexp.MustCompile(`^- \[[^\]]+\]\(https?://`)
	reLedger      = regexp.MustCompile(`^\| *\[([^\]]+)\]`)
	reWeek        = regexp.MustCompile(`^\d{4}-w\d{2}$`)
	reDateTail    = regexp.MustCompile(`^` + datePat + `\*`)
	reBuiltin     = regexp.MustCompile("`/([a-z0-9-]+)`")
	reBuiltinL    = regexp.MustCompile(`^\*\*Built-ins:\*\* `)
	reIntegL      = regexp.MustCompile(`^\*\*Integrations:\*\* `)
	reIntegTok    = regexp.MustCompile("`([a-z0-9-]+)`")
	reRegKind     = regexp.MustCompile(`^kind: ([a-z-]+)$`)
	reRegGoals    = regexp.MustCompile(`^goals: (.+)$`)
	reClassifyRow = regexp.MustCompile("^\\| `([a-z0-9-]+)` \\|")
	reCount       = regexp.MustCompile(`(\d+) goal categories`)
)

var (
	skillDir = filepath.Join("skills", "mentor")

	// "## Real-World Example" is deliberately absent: it is optional — kept
	// only where the example embeds exact syntax (see templates/approach.md).
	// When present it must sit between Common Pitfalls and Sources
	// (checkApproach enforces that).
	approachSections = []string{
		"## What It Is", "## Why It Works", "## When to Use It", "## When NOT to Use It",
		"## How It Works", "### Basic (Beginner)",
		"### Composing with Other Approaches (Intermediate)", "### Advanced Patterns",
		"## Common Pitfalls", "## Sources", "## Signals",
	}
)

type auditor struct {
	root                    string // repo root; issue paths print relative to it
	skill                   string // skills/mentor directory
	issues                  []string
	goals, solutions, weeks int
}

func (a *auditor) issue(path, format string, args ...any) {
	rel := path
	if r, err := filepath.Rel(a.root, path); err == nil {
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

// cells splits a Markdown table row on '|' and trims each cell.
func cells(l string) []string {
	cs := strings.Split(l, "|")
	for i, c := range cs {
		cs[i] = strings.TrimSpace(c)
	}
	return cs
}

// validDate reports whether s is a real calendar date (rejects 2026-99-99,
// which the format regex alone would accept).
func validDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

// pluginNames extracts the plugin ids the catalog declares: the first
// backticked token of each table row plus the backticked tokens in the prose
// sections (Language servers, Specialty). Nothing else — so goal slugs and
// backticked command names are never mistaken for plugins. Keep in sync with
// the copy in tools/catalog-drift/main.go.
func pluginNames(text string) []string {
	var names []string
	proseList := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "#") {
			h := strings.ToLower(line)
			proseList = strings.Contains(h, "language server") || strings.Contains(h, "specialty")
			continue
		}
		if m := reRowName.FindStringSubmatch(line); m != nil {
			names = append(names, m[1])
			continue
		}
		if proseList {
			for _, m := range rePlugTok.FindAllStringSubmatch(line, -1) {
				names = append(names, m[1])
			}
		}
	}
	return names
}

// checkDocRefs verifies that catalog paths named in SKILL.md and the mode
// files resolve — the load-bearing references (profile-schema, registries,
// routing, approaches) that nothing else in the audit reads.
func (a *auditor) checkDocRefs(files ...string) {
	for _, f := range files {
		text := strings.Join(lines(f), "\n")
		for _, ref := range dedup(reDocRef.FindAllString(text, -1)) {
			if _, err := os.Stat(filepath.Join(a.skill, ref)); err != nil {
				a.issue(f, "broken reference %s", ref)
			}
		}
	}
}

// dateLine checks that line 2 is e.g. "*Last verified: 2026-07-03*" with a
// real date.
func (a *auditor) dateLine(path, label string, ls []string) {
	ok := false
	if len(ls) >= 2 {
		if rest, found := strings.CutPrefix(ls[1], "*"+label+": "); found {
			if m := reDateTail.FindString(rest); m != "" {
				ok = validDate(strings.TrimSuffix(m, "*"))
			}
		}
	}
	if !ok {
		a.issue(path, "line 2 must be '*%s: YYYY-MM-DD*' with a real date", label)
	}
}

// checkRouting audits every per-goal file under skills/mentor/problems/:
// date line, hidden gem, Plugins line with catalog-known tokens, at least
// minGoalRows sequentially numbered rows, cross-references, and orphans.
// Built-ins and Integrations tokens must resolve to a solutions/<id>.md file
// of any kind (a command can be documented as a technique deep-dive).
func (a *auditor) checkRouting(dir string, techniqueNames []string, catalog, solutions, promoted map[string]bool, usedBuiltins, usedIntegrations, usedPlugins map[string]bool) []string {
	files, _ := filepath.Glob(filepath.Join(dir, "*.md"))
	if len(files) == 0 {
		a.issue(dir, "missing problems directory (per-goal files)")
		return nil
	}
	var goals []string
	var all strings.Builder
	for _, f := range files {
		ls := lines(f)
		goals = append(goals, strings.TrimSuffix(filepath.Base(f), ".md"))
		a.dateLine(f, "Last verified", ls)
		rows, gem, plugs := []string{}, "", false
		for _, l := range ls {
			if m := reGem.FindStringSubmatch(l); m != nil {
				gem = m[1]
			}
			if rePlugLine.MatchString(l) {
				plugs = true
				for _, m := range rePlugTok.FindAllStringSubmatch(l, -1) {
					if promoted[m[1]] {
						usedPlugins[m[1]] = true
						continue
					}
					if !catalog[m[1]] {
						a.issue(f, "Plugins line names '%s', which is neither a promoted plugin record nor a marketplace.md row", m[1])
					}
				}
			}
			if reBuiltinL.MatchString(l) {
				for _, m := range reBuiltin.FindAllStringSubmatch(l, -1) {
					if !solutions[m[1]] {
						a.issue(f, "Built-ins line names '/%s', which has no solutions/%s.md", m[1], m[1])
					}
					usedBuiltins[m[1]] = true
				}
			}
			if reIntegL.MatchString(l) {
				// reIntegTok's charset (no dots, no slashes) only matches solution
				// ids — backticked file names and placeholders never capture.
				for _, m := range reIntegTok.FindAllStringSubmatch(l, -1) {
					if !solutions[m[1]] {
						a.issue(f, "Integrations line names '%s', which has no solutions/%s.md", m[1], m[1])
					}
					usedIntegrations[m[1]] = true
				}
			}
			if m := reRow.FindStringSubmatch(l); m != nil {
				if n, _ := strconv.Atoi(m[1]); n != len(rows)+1 {
					a.issue(f, "row numbering not sequential at row %d", len(rows)+1)
				}
				rows = append(rows, m[2])
			}
		}
		if len(rows) < minGoalRows {
			a.issue(f, "only %d rows (expected at least %d)", len(rows), minGoalRows)
		}
		if !plugs {
			a.issue(f, "missing Plugins line")
		}
		if gem == "" {
			a.issue(f, "missing Hidden gem line")
		} else {
			g := strings.ToLower(strings.TrimSpace(gem))
			ok := slices.ContainsFunc(rows, func(r string) bool {
				rl := strings.ToLower(r)
				return strings.Contains(g, rl) || strings.Contains(rl, g)
			})
			if !ok {
				a.issue(f, "Hidden gem '%s' does not match any ranked row", strings.TrimSpace(gem))
			}
		}
		all.WriteString(strings.Join(ls, "\n"))
		all.WriteString("\n")
	}

	text := all.String()
	for _, ref := range dedup(reRef.FindAllString(text, -1)) {
		if _, err := os.Stat(filepath.Join(a.skill, ref)); err != nil {
			a.issue(dir, "broken reference %s", ref)
		}
	}
	for _, name := range techniqueNames {
		if !strings.Contains(text, "solutions/"+name+".md") {
			a.issue(filepath.Join(a.skill, "solutions", name+".md"), "orphan: not ranked by any problems file")
		}
	}
	return goals
}

func (a *auditor) checkApproach(path string) {
	ls := lines(path)
	a.dateLine(path, "Last verified", ls)

	pos := 0
	for _, s := range approachSections {
		ln := slices.IndexFunc(ls, func(l string) bool { return strings.Contains(l, s) }) + 1
		switch {
		case ln == 0:
			a.issue(path, "missing section '%s'", s)
		case ln < pos:
			a.issue(path, "section '%s' out of order", s)
		default:
			pos = ln
		}
	}

	find := func(s string) int {
		return slices.IndexFunc(ls, func(l string) bool { return strings.Contains(l, s) })
	}
	if ex := find("## Real-World Example"); ex >= 0 {
		if src := find("## Sources"); ex < find("## Common Pitfalls") || (src >= 0 && ex > src) {
			a.issue(path, "section '## Real-World Example' out of order")
		}
	}

	if n := len(ls) - 1; n < minApproachLines { // trailing newline yields one empty element
		a.issue(path, "%d lines (expected at least %d)", n, minApproachLines)
	}
	srcs := 0
	if i := slices.Index(ls, "## Sources"); i >= 0 {
		for _, l := range ls[i+1:] {
			if reSource.MatchString(l) {
				srcs++
			}
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
		cs := cells(l)
		if len(cs) < 5 {
			a.issue(path, "row '%s' is malformed (need | Week | Processed | Outcome |)", slug)
			continue
		}
		if !validDate(cs[2]) {
			a.issue(path, "row '%s' has invalid processed date '%s'", slug, cs[2])
		}
		if cs[3] == "" {
			a.issue(path, "row '%s' has an empty outcome", slug)
		}
	}
}

func (a *auditor) checkProblemMode(path string, goals []string) {
	ls := lines(path)
	var table []string
	for _, l := range ls {
		if m := reClassifyRow.FindStringSubmatch(l); m != nil {
			table = append(table, m[1])
		}
	}
	missing, stale := diff(goals, table)
	for _, x := range missing {
		a.issue(path, "routing goal %s missing from the problem-mode classification table", x)
	}
	for _, x := range stale {
		a.issue(path, "classification table row %s has no matching routing section", x)
	}
	for _, l := range ls {
		for _, m := range reCount.FindAllStringSubmatch(l, -1) {
			if n, _ := strconv.Atoi(m[1]); n != len(goals) {
				a.issue(path, "prose says '%d goal categories' but there are %d", n, len(goals))
			}
		}
	}
}

// fileKind returns the value of a solution file's kind: line, or "" for a
// technique deep-dive (which has no kind: line).
func fileKind(path string) string {
	for _, l := range lines(path) {
		if m := reRegKind.FindStringSubmatch(l); m != nil {
			return m[1]
		}
	}
	return ""
}

func (a *auditor) run() error {
	solDir := filepath.Join(a.skill, "solutions")
	files, _ := filepath.Glob(filepath.Join(solDir, "*.md"))
	if len(files) == 0 {
		return fmt.Errorf("solutions directory empty/missing")
	}
	var techNames []string
	var techFiles, recFiles []string
	recordKind := map[string]string{} // id -> kind
	solutions := map[string]bool{}
	for _, f := range files {
		id := strings.TrimSuffix(filepath.Base(f), ".md")
		if id == "index" {
			continue // the compiled index; freshness is tools/solutions-index -check
		}
		solutions[id] = true
		if k := fileKind(f); k != "" {
			recordKind[id] = k
			recFiles = append(recFiles, f)
		} else {
			techNames = append(techNames, id)
			techFiles = append(techFiles, f)
		}
	}
	a.solutions = len(solutions)

	catalog := map[string]bool{}
	catPath := filepath.Join(a.skill, "marketplace.md")
	catText, catErr := os.ReadFile(catPath)
	if catErr != nil {
		a.issue(catPath, "missing marketplace directory")
	}
	for _, n := range pluginNames(string(catText)) {
		catalog[n] = true
	}
	promoted := map[string]bool{}
	for id, kind := range recordKind {
		if kind == "plugin" {
			if catalog[id] {
				a.issue(filepath.Join(solDir, id+".md"), "promoted plugin still has a marketplace.md row — remove the directory row")
			}
			promoted[id] = true
		}
	}

	usedBuiltins, usedIntegrations, usedPlugins := map[string]bool{}, map[string]bool{}, map[string]bool{}
	goals := a.checkRouting(filepath.Join(a.skill, "problems"), techNames, catalog, solutions, promoted, usedBuiltins, usedIntegrations, usedPlugins)
	a.goals = len(goals)
	for id, kind := range recordKind {
		recPath := filepath.Join(solDir, id+".md")
		switch kind {
		case "builtin-command":
			if !usedBuiltins[id] {
				a.issue(recPath, "command record not referenced by any Built-ins line")
			}
		case "integration", "doc":
			if !usedIntegrations[id] {
				a.issue(recPath, "integration record not referenced by any Integrations line")
			}
		case "plugin":
			if !usedPlugins[id] {
				a.issue(recPath, "plugin record not referenced by any Plugins line")
			}
		default:
			a.issue(recPath, "unknown kind '%s'", kind)
		}
	}
	for _, f := range recFiles {
		a.checkRecord(f, goals)
	}
	for _, f := range techFiles {
		a.checkApproach(f)
	}
	a.checkLedger(filepath.Join(a.skill, "processed-changelogs.md"))
	a.checkProblemMode(filepath.Join(a.skill, "problem-mode.md"), goals)
	a.checkDocRefs(
		filepath.Join(a.skill, "SKILL.md"),
		filepath.Join(a.skill, "problem-mode.md"),
		filepath.Join(a.skill, "growth-mode.md"),
	)
	return nil
}

// dedup is order-preserving on purpose: its order feeds issue order, which
// is part of the frozen output. Do not replace with a sorted form.
func dedup(xs []string) []string {
	var out []string
	for _, x := range xs {
		if !slices.Contains(out, x) {
			out = append(out, x)
		}
	}
	return out
}

// diff returns (in a but not b, in b but not a), preserving order.
func diff(a, b []string) (onlyA, onlyB []string) {
	for _, x := range a {
		if !slices.Contains(b, x) {
			onlyA = append(onlyA, x)
		}
	}
	for _, x := range dedup(b) {
		if !slices.Contains(a, x) {
			onlyB = append(onlyB, x)
		}
	}
	return
}

// checkRecord audits one flat record file under solutions/: date line, and a
// goals line naming only real goal slugs (files under problems/). Content
// completeness (best_when, session_signal) is the generator's job.
func (a *auditor) checkRecord(path string, goals []string) {
	ls := lines(path)
	a.dateLine(path, "Last verified", ls)
	for _, l := range ls {
		m := reRegGoals.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		for _, g := range strings.Split(m[1], ",") {
			g = strings.TrimSpace(g)
			if !slices.Contains(goals, g) {
				a.issue(path, "goals name '%s', which has no problems/%s.md", g, g)
			}
		}
	}
}

// findRoot walks upward from dir to the first directory containing
// skills/mentor, so the audit works from anywhere in the repo — including
// tools/structural-audit itself, where `go -C tools/structural-audit run .` lands.
// Keep in sync with the copy in tools/catalog-drift/main.go.
func findRoot(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, skillDir)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no %s directory found here or above", skillDir)
		}
		dir = parent
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
	os.Exit(2)
}

func main() {
	var repo string
	var err error
	if len(os.Args) > 1 {
		repo = os.Args[1]
	} else if repo, err = findRoot("."); err != nil {
		fatal(err)
	}
	a := &auditor{root: repo, skill: filepath.Join(repo, skillDir)}
	if err := a.run(); err != nil {
		fatal(err)
	}
	for _, is := range a.issues {
		fmt.Printf("  - %s\n", is)
	}
	fmt.Printf("Audited %d problems, %d solutions, %d processed changelogs.\n", a.goals, a.solutions, a.weeks)
	if len(a.issues) > 0 {
		fmt.Printf("\n%d issue(s) found (listed above).\n", len(a.issues))
		os.Exit(1)
	}
	fmt.Println("Structural audit: PASS")
}
