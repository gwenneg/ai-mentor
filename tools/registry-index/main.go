// Generates skills/mentor/registry/index.md — the compiled capability index —
// from the authored sources of truth:
//
//   - routing/*.md          goal membership, rank, setup badge, best-when triggers
//   - approaches/*.md       the "## Signals" section (setup + session signals)
//   - registry/builtin-commands.md and registry/integrations.md
//     (id, kind, goals, best_when, setup, session_signal per record)
//
// index.md is a build artifact: never hand-edit it. After editing any source
// above, regenerate with `go -C tools/registry-index run .`. In CI, `-check`
// regenerates in memory and exits 1 if the on-disk file is stale.
//
// Merge rule: a capability that exists both as a technique and as a registry
// record under the same id (e.g. deep-research) is ONE row — kinds joined,
// goals unioned, and the record's best_when/setup/session_signal win over the
// routing-derived values (records are authored per-capability; badges are
// derived through routing).
//
// Deterministic by construction: rows sorted by id, no timestamps. Exits 1 on
// any inconsistency in the sources (conflicting setup badges, missing Signals
// section), 2 on a fatal setup problem. No network, no LLM. Stdlib only.
//
// Usage: go -C tools/registry-index run . [-check] [repo-root]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var skillDir = filepath.Join("skills", "mentor")

var (
	reRow     = regexp.MustCompile(`^\| (\d+) \| \[([^\]]+)\]\(\.\./approaches/([a-z0-9-]+)\.md\)`)
	reRegID   = regexp.MustCompile(`^id: ([a-z0-9-]+)$`)
	reRegKind = regexp.MustCompile(`^kind: ([a-z-]+)$`)
	reRegGoal = regexp.MustCompile(`^goals: (.+)$`)
	reRegBest = regexp.MustCompile(`^best_when: (.+)$`)
	reRegSet  = regexp.MustCompile(`^setup: ([a-z]+)$`)
	reRegSig  = regexp.MustCompile(`^session_signal: (.+)$`)
	reSetupL  = regexp.MustCompile(`^- Setup: (.+)$`)
	reSessL   = regexp.MustCompile(`^- Session: (.+)$`)
)

// validSetup is the one setup vocabulary, shared by routing cells and records.
var validSetup = map[string]bool{"none": true, "some": true, "involved": true}

type row struct {
	id, kind, bestWhen, setup, setupSig, sessionSig string
	goals                                           []string
}

type gen struct {
	root, skill string
	errs        []string
}

func (g *gen) errf(path, format string, args ...any) {
	rel := path
	if r, err := filepath.Rel(g.root, path); err == nil {
		rel = r
	}
	g.errs = append(g.errs, rel+": "+fmt.Sprintf(format, args...))
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

// techniques derives one row per approach file: goals and setup from the
// routing tables, best_when from the approach's best-ranked routing row
// (lowest rank; ties broken by goal-file order, which is alphabetical),
// signals from the approach file's "## Signals" section.
func (g *gen) techniques() map[string]*row {
	rows := map[string]*row{}
	bestRank := map[string]int{}

	routing, _ := filepath.Glob(filepath.Join(g.skill, "routing", "*.md"))
	slices.Sort(routing)
	for _, f := range routing {
		goal := strings.TrimSuffix(filepath.Base(f), ".md")
		for _, l := range lines(f) {
			m := reRow.FindStringSubmatch(l)
			if m == nil {
				continue
			}
			rank, _ := strconv.Atoi(m[1])
			slug := m[3]
			cs := cells(l) // | # | Approach | Setup | Best when | Why it fits |
			if len(cs) < 5 {
				g.errf(f, "row for '%s' has too few columns", slug)
				continue
			}
			r := rows[slug]
			if r == nil {
				r = &row{id: slug, kind: "technique"}
				rows[slug] = r
				bestRank[slug] = 1 << 30
			}
			r.goals = append(r.goals, goal)
			setup := cs[3]
			if !validSetup[setup] {
				g.errf(f, "row for '%s' has invalid setup '%s'", slug, setup)
			} else if r.setup == "" {
				r.setup = setup
			} else if r.setup != setup {
				g.errf(f, "setup for '%s' conflicts with another routing file ('%s' vs '%s')", slug, r.setup, setup)
			}
			if rank < bestRank[slug] {
				bestRank[slug] = rank
				r.bestWhen = strings.ToLower(cs[4][:1]) + cs[4][1:]
			}
		}
	}

	approaches, _ := filepath.Glob(filepath.Join(g.skill, "approaches", "*.md"))
	for _, f := range approaches {
		slug := strings.TrimSuffix(filepath.Base(f), ".md")
		r := rows[slug]
		if r == nil {
			g.errf(f, "approach has no routing row — cannot index it")
			continue
		}
		inSection := false
		for _, l := range lines(f) {
			switch {
			case strings.HasPrefix(l, "## "):
				inSection = l == "## Signals"
			case inSection:
				if m := reSetupL.FindStringSubmatch(l); m != nil {
					r.setupSig = m[1]
				}
				if m := reSessL.FindStringSubmatch(l); m != nil {
					r.sessionSig = m[1]
				}
			}
		}
		if r.setupSig == "" || r.sessionSig == "" {
			g.errf(f, "missing or incomplete '## Signals' section (need '- Setup:' and '- Session:' lines)")
		}
	}
	for slug := range rows {
		if _, err := os.Stat(filepath.Join(g.skill, "approaches", slug+".md")); err != nil {
			g.errf(filepath.Join(g.skill, "routing"), "routing rows reference approaches/%s.md, which does not exist", slug)
			delete(rows, slug)
		}
	}
	return rows
}

// records parses registry/builtin-commands.md or registry/integrations.md
// into rows (kind comes from each record's own kind: line).
func (g *gen) records(path string) []*row {
	var out []*row
	var cur *row
	for _, l := range lines(path) {
		if m := reRegID.FindStringSubmatch(l); m != nil {
			cur = &row{id: m[1]}
			out = append(out, cur)
			continue
		}
		if cur == nil {
			continue
		}
		switch {
		case reRegKind.MatchString(l):
			cur.kind = reRegKind.FindStringSubmatch(l)[1]
		case reRegGoal.MatchString(l):
			for _, gl := range strings.Split(reRegGoal.FindStringSubmatch(l)[1], ",") {
				cur.goals = append(cur.goals, strings.TrimSpace(gl))
			}
		case reRegBest.MatchString(l):
			cur.bestWhen = reRegBest.FindStringSubmatch(l)[1]
		case reRegSet.MatchString(l):
			cur.setup = reRegSet.FindStringSubmatch(l)[1]
		case reRegSig.MatchString(l):
			cur.sessionSig = reRegSig.FindStringSubmatch(l)[1]
		}
	}
	for _, r := range out {
		if r.kind == "" || len(r.goals) == 0 || r.bestWhen == "" || r.setup == "" || r.sessionSig == "" {
			g.errf(path, "record '%s' is missing one of kind/goals/best_when/setup/session_signal", r.id)
		}
		if r.setupSig == "" {
			r.setupSig = "—"
		}
	}
	return out
}

func (g *gen) build() string {
	rows := g.techniques()
	for _, path := range []string{
		filepath.Join(g.skill, "registry", "builtin-commands.md"),
		filepath.Join(g.skill, "registry", "integrations.md"),
	} {
		if lines(path) == nil {
			g.errf(path, "missing registry file")
			continue
		}
		for _, rec := range g.records(path) {
			if t := rows[rec.id]; t != nil {
				// One capability, two kinds: merge. Record fields win; goals union.
				t.kind = t.kind + " + " + rec.kind
				t.bestWhen, t.setup, t.sessionSig = rec.bestWhen, rec.setup, rec.sessionSig
				t.goals = append(t.goals, rec.goals...)
			} else {
				rows[rec.id] = rec
			}
		}
	}

	ids := make([]string, 0, len(rows))
	for id := range rows {
		ids = append(ids, id)
	}
	slices.Sort(ids)

	var b strings.Builder
	b.WriteString("# Capability Index\n")
	b.WriteString("*Generated — do not edit. Sources: routing tables, approach Signals sections, registry records. Regenerate: `go -C tools/registry-index run .`*\n\n")
	b.WriteString("One row per first-party capability (marketplace plugins live in `references/official-plugins.md`). Setup signals are re-checkable disk evidence; session signals are conversation evidence accumulated in the profile. `—` means no signal of that tier exists.\n\n")
	b.WriteString("| Id | Kind | Goals | Best when | Setup | Setup signal | Session signal |\n")
	b.WriteString("|----|------|-------|-----------|-------|--------------|----------------|\n")
	for _, id := range ids {
		r := rows[id]
		goals := dedup(r.goals)
		slices.Sort(goals)
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s |\n",
			r.id, r.kind, strings.Join(goals, ", "), r.bestWhen, r.setup, r.setupSig, r.sessionSig)
	}
	return b.String()
}

func dedup(xs []string) []string {
	var out []string
	for _, x := range xs {
		if !slices.Contains(out, x) {
			out = append(out, x)
		}
	}
	return out
}

// findRoot walks upward from dir to the first directory containing
// skills/mentor. Keep in sync with the copy in tools/structural-audit/main.go.
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
	check := flag.Bool("check", false, "verify index.md is up to date instead of writing it")
	flag.Parse()

	var repo string
	var err error
	if flag.NArg() > 0 {
		repo = flag.Arg(0)
	} else if repo, err = findRoot("."); err != nil {
		fatal(err)
	}
	g := &gen{root: repo, skill: filepath.Join(repo, skillDir)}
	out := g.build()
	if len(g.errs) > 0 {
		for _, e := range g.errs {
			fmt.Printf("  - %s\n", e)
		}
		fmt.Printf("\n%d source issue(s) — index not %s.\n", len(g.errs), map[bool]string{true: "verified", false: "written"}[*check])
		os.Exit(1)
	}

	target := filepath.Join(g.skill, "registry", "index.md")
	if *check {
		disk, err := os.ReadFile(target)
		if err != nil || string(disk) != out {
			fmt.Println("registry/index.md is stale — regenerate with `go -C tools/registry-index run .`")
			os.Exit(1)
		}
		fmt.Println("Capability index: fresh")
		return
	}
	if err := os.WriteFile(target, []byte(out), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("Wrote %s (%d capabilities)\n", target, strings.Count(out, "\n|")-2)
}
