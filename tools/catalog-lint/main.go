// Deterministic structural lint for the ai-mentor catalog.
//
// Checks the playbooks tables, approach files (technique deep-dives and flat
// records), cross-references, the changelog ledger, and SKILL.md consistency.
// Exits 1 if any issue is found, 2 on a fatal setup problem. No network, no
// LLM — safe as a PR gate. Stdlib only.
//
// The compiled index (approaches/index.md) is NOT audited here:
// tools/approaches-index generates it from the same sources and its -check
// mode is the freshness gate in CI.
//
// Usage: go -C tools/catalog-lint run . [repo-root]
package main

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	refPat           = `approaches/(techniques|tools)/[a-z0-9-]+\.md`
	minGoalRows      = 3
	minApproachLines = 40
)

var (
	reRow         = regexp.MustCompile(`^\| (\d+) \| \[([^\]]+)\]`)
	reGem         = regexp.MustCompile(`^\*\*Hidden gem:\*\* ([^—\n]+)`)
	rePlugTok     = regexp.MustCompile("`([a-z0-9.-]+)`")
	reRowName     = regexp.MustCompile("^\\| `([a-z0-9.-]+)`")
	reDocRef      = regexp.MustCompile(`playbooks/[a-z0-9-]+\.md|` + refPat + `|approaches/index\.md|\b(marketplace|profile-schema|processed-changelogs)\.md`)
	reRef         = regexp.MustCompile(refPat)
	reSource      = regexp.MustCompile(`^- \[[^\]]+\]\(https?://`)
	reLedger      = regexp.MustCompile(`^\| *\[([^\]]+)\]`)
	reWeek        = regexp.MustCompile(`^\d{4}-w\d{2}$`)
	reRegKind     = regexp.MustCompile(`^kind: ([a-z-]+)$`)
	reClassifyRow = regexp.MustCompile("^\\| `([a-z0-9-]+)` \\|")
	reCount       = regexp.MustCompile(`(\d+) goal categories`)
)

var (
	skillDir = filepath.Join("skills", "mentor")

	// "## Real-World Example" is deliberately absent: it is optional — kept
	// only where the example embeds exact syntax (see templates/technique.md).
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
	root                     string // repo root; issue paths print relative to it
	skill                    string // skills/mentor directory
	issues                   []string
	goals, approaches, weeks int
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

// slug returns an approach file's id: its base name without .md.
func slug(path string) string {
	return strings.TrimSuffix(filepath.Base(path), ".md")
}

// cells splits a Markdown table row on '|' and trims each cell.
func cells(l string) []string {
	cs := strings.Split(l, "|")
	for i, c := range cs {
		cs[i] = strings.TrimSpace(c)
	}
	return cs
}

// validDate reports whether s is exactly a real YYYY-MM-DD calendar date —
// time.Parse enforces zero-padding and rejects any surrounding text.
func validDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

// pluginNames extracts the plugin ids the catalog declares: the first
// backticked token of each table row plus the backticked tokens in the prose
// sections (Language servers, Specialty). Nothing else — so goal slugs and
// backticked command names are never mistaken for plugins. Keep in sync with
// the copies in tools/catalog-drift/main.go and tools/eval-runner/main.go.
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

// checkRefs stats every catalog path re finds in text, reporting broken ones
// against at.
func (a *auditor) checkRefs(at, text string, re *regexp.Regexp) {
	for _, ref := range dedup(re.FindAllString(text, -1)) {
		if _, err := os.Stat(filepath.Join(a.skill, ref)); err != nil {
			a.issue(at, "broken reference %s", ref)
		}
	}
}

// checkDocRefs verifies that catalog paths named in SKILL.md and the mode
// files resolve — the load-bearing references (profile-schema, registries,
// routing, approaches) that nothing else in the audit reads.
func (a *auditor) checkDocRefs(files ...string) {
	for _, f := range files {
		b, _ := os.ReadFile(f)
		a.checkRefs(f, string(b), reDocRef)
	}
}

// dateLine checks that line 2 is e.g. "*Last verified: 2026-07-03*" with a
// real date.
func (a *auditor) dateLine(path, label string, ls []string) {
	ok := false
	if len(ls) >= 2 {
		if rest, found := strings.CutPrefix(ls[1], "*"+label+": "); found {
			if date, _, closed := strings.Cut(rest, "*"); closed {
				ok = validDate(date)
			}
		}
	}
	if !ok {
		a.issue(path, "line 2 must be '*%s: YYYY-MM-DD*' with a real date", label)
	}
}

// checkRouting audits every per-goal file under skills/mentor/playbooks/:
// date line, hidden gem, at least minGoalRows sequentially numbered rows,
// cross-references, and orphans. ranked maps every approach id — technique
// or record — to its file path; each must appear in at least one ranked row:
// the ranking is the only routing surface.
func (a *auditor) checkRouting(dir string, ranked map[string]string) []string {
	files, _ := filepath.Glob(filepath.Join(dir, "*.md"))
	if len(files) == 0 {
		a.issue(dir, "missing playbooks directory (per-goal files)")
		return nil
	}
	var goals []string
	var all strings.Builder
	for _, f := range files {
		b, _ := os.ReadFile(f)
		ls := strings.Split(string(b), "\n")
		goals = append(goals, slug(f))
		a.dateLine(f, "Last verified", ls)
		rows, gem := []string{}, ""
		for _, l := range ls {
			if m := reGem.FindStringSubmatch(l); m != nil {
				gem = strings.TrimSpace(m[1])
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
		if gem == "" {
			a.issue(f, "missing Hidden gem line")
		} else {
			g := strings.ToLower(gem)
			ok := slices.ContainsFunc(rows, func(r string) bool {
				rl := strings.ToLower(r)
				return strings.Contains(g, rl) || strings.Contains(rl, g)
			})
			if !ok {
				a.issue(f, "Hidden gem '%s' does not match any ranked row", gem)
			}
		}
		all.Write(b)
		all.WriteByte('\n')
	}

	text := all.String()
	a.checkRefs(dir, text, reRef)
	for _, name := range slices.Sorted(maps.Keys(ranked)) {
		if !strings.Contains(text, "/"+name+".md") {
			a.issue(ranked[name], "orphan: not ranked by any playbooks file")
		}
	}
	return goals
}

func (a *auditor) checkApproach(path string) {
	ls := lines(path)
	a.dateLine(path, "Last verified", ls)

	find := func(s string) int {
		return slices.IndexFunc(ls, func(l string) bool { return strings.Contains(l, s) })
	}
	pos := 0
	for _, s := range approachSections {
		ln := find(s) + 1
		switch {
		case ln == 0:
			a.issue(path, "missing section '%s'", s)
		case ln < pos:
			a.issue(path, "section '%s' out of order", s)
		default:
			pos = ln
		}
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

// checkLedger returns the number of ledger rows for the summary line.
func (a *auditor) checkLedger(path string) (weeks int) {
	ls := lines(path)
	if ls == nil {
		a.issue(path, "missing processed-changelog ledger")
		return 0
	}
	a.dateLine(path, "Updated", ls)

	seen := map[string]bool{}
	for _, l := range ls {
		m := reLedger.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		weeks++
		week := m[1]
		if !reWeek.MatchString(week) {
			a.issue(path, "row '%s' is not a week slug like 2026-w26", week)
		}
		if seen[week] {
			a.issue(path, "duplicate ledger row for '%s'", week)
		}
		seen[week] = true
		cs := cells(l)
		if len(cs) < 5 {
			a.issue(path, "row '%s' is malformed (need | Week | Processed | Outcome |)", week)
			continue
		}
		if !validDate(cs[2]) {
			a.issue(path, "row '%s' has invalid processed date '%s'", week, cs[2])
		}
		if cs[3] == "" {
			a.issue(path, "row '%s' has an empty outcome", week)
		}
	}
	return weeks
}

func (a *auditor) checkProblemMode(path string, goals []string) {
	ls := lines(path)
	var table []string
	for _, l := range ls {
		if m := reClassifyRow.FindStringSubmatch(l); m != nil {
			table = append(table, m[1])
		}
	}
	missing, stale := subtract(goals, table), subtract(dedup(table), goals)
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

// fileKind returns the value of an approach file's kind: line, or "" for a
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
	techDir := filepath.Join(a.skill, "approaches", "techniques")
	recDir := filepath.Join(a.skill, "approaches", "tools")
	techFiles, _ := filepath.Glob(filepath.Join(techDir, "*.md"))
	if len(techFiles) == 0 {
		return fmt.Errorf("approaches/techniques directory empty/missing")
	}
	recFiles, _ := filepath.Glob(filepath.Join(recDir, "*.md"))

	catalog := map[string]bool{}
	catPath := filepath.Join(a.skill, "marketplace.md")
	catText, catErr := os.ReadFile(catPath)
	if catErr != nil {
		a.issue(catPath, "missing marketplace directory")
	}
	for _, n := range pluginNames(string(catText)) {
		catalog[n] = true
	}

	// ranked maps every approach id to its file path — the orphan check
	// reports against the real subfolder location. Every approach — technique
	// or record — must be a ranked row; kind is a semantic label, not a
	// routing tier. Promoted plugins additionally must not retain a
	// marketplace.md directory row.
	ranked := map[string]string{}
	for _, f := range techFiles {
		ranked[slug(f)] = f
	}
	for _, f := range recFiles {
		id := slug(f)
		ranked[id] = f
		switch kind := fileKind(f); kind { // "" (missing kind) is the generator's error
		case "plugin":
			if catalog[id] {
				a.issue(f, "promoted plugin still has a marketplace.md row — remove the directory row")
			}
		case "integration", "doc", "":
		default:
			a.issue(f, "unknown kind '%s'", kind)
		}
	}
	a.approaches = len(ranked)

	goals := a.checkRouting(filepath.Join(a.skill, "playbooks"), ranked)
	a.goals = len(goals)
	for _, f := range recFiles {
		a.checkRecord(f)
	}
	for _, f := range techFiles {
		a.checkApproach(f)
	}
	a.weeks = a.checkLedger(filepath.Join(a.skill, "processed-changelogs.md"))
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

// subtract returns the elements of xs not in ys, preserving order.
func subtract(xs, ys []string) []string {
	var out []string
	for _, x := range xs {
		if !slices.Contains(ys, x) {
			out = append(out, x)
		}
	}
	return out
}

// checkRecord audits one flat record under approaches/tools/: it must be a
// pure ----delimited YAML frontmatter file with a valid last_verified date.
// Content completeness (kind, session_signal, no inline goals/best_when) is
// the generator's job.
func (a *auditor) checkRecord(path string) {
	ls := lines(path)
	if len(ls) < 3 || ls[0] != "---" {
		a.issue(path, "record must be a pure YAML frontmatter file starting with ---")
		return
	}
	closed, dated := false, false
	for _, l := range ls[1:] {
		if l == "---" {
			closed = true
			break
		}
		if v, found := strings.CutPrefix(l, "last_verified: "); found && validDate(v) {
			dated = true
		}
	}
	if !closed {
		a.issue(path, "frontmatter never closed with ---")
	}
	if !dated {
		a.issue(path, "missing 'last_verified: YYYY-MM-DD' with a real date")
	}
}

// findRoot walks upward from dir to the first directory containing
// skills/mentor, so the audit works from anywhere in the repo — including
// tools/catalog-lint itself, where `go -C tools/catalog-lint run .` lands.
// Keep in sync with the copies in tools/catalog-drift, tools/approaches-index,
// and tools/eval-runner.
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
	fmt.Printf("Audited %d playbooks, %d approaches, %d processed changelogs.\n", a.goals, a.approaches, a.weeks)
	if len(a.issues) > 0 {
		fmt.Printf("\n%d issue(s) found (listed above).\n", len(a.issues))
		os.Exit(1)
	}
	fmt.Println("Catalog lint: PASS")
}
