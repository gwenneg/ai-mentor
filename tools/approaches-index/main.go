// Generates skills/mentor/approaches/index.md — the compiled solutions index —
// from the authored sources of truth:
//
//   - playbooks/*.md               goal membership, rank, and best-when
//     triggers for every ranked approach (the ranked rows)
//   - approaches/techniques/*.md   prose deep-dives with a "## Signals" section
//   - approaches/tools/*.md      pure YAML-frontmatter fact sheets (plugin,
//     integration, or doc — the filename is the id; kind is a semantic label,
//     not a routing tier)
//
// Everything is ranked: goals and best_when always derive from the rows.
//
// index.md is a build artifact: never hand-edit it. After editing any source
// above, regenerate with `go -C tools/approaches-index run .`. In CI, `-check`
// regenerates in memory and exits 1 if the on-disk file is stale.
//
// Deterministic by construction: rows sorted by id, no timestamps. Exits 1 on
// any inconsistency in the sources (a technique without a ranked row, a
// missing Signals section, a ranked row pointing at a record), 2 on a fatal
// setup problem. No network, no LLM. Stdlib only.
//
// Usage: go -C tools/approaches-index run . [-check] [repo-root]
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
	reRow     = regexp.MustCompile(`^\| (\d+) \| \[([^\]]+)\]\(\.\./approaches/(?:techniques|tools)/([a-z0-9-]+)\.md\)`)
	reRegKind = regexp.MustCompile(`^kind: ([a-z-]+)$`)
	reRegGoal = regexp.MustCompile(`^goals: (.+)$`)
	reRegBest = regexp.MustCompile(`^best_when: (.+)$`)
	reRegSig  = regexp.MustCompile(`^session_signal: (.+)$`)
	reSetupL  = regexp.MustCompile(`^- Setup: (.+)$`)
	reSessL   = regexp.MustCompile(`^- Session: (.+)$`)
)

type row struct {
	id, kind, bestWhen, setupSig, sessionSig string
	goals                                    []string
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

// unquote strips a surrounding double-quote pair from a YAML scalar and
// unescapes \" and \\ — records quote free-text values because they contain
// ": ", which YAML plain scalars forbid.
func unquote(v string) string {
	if len(v) >= 2 && strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`) {
		v = v[1 : len(v)-1]
		v = strings.ReplaceAll(v, `\"`, `"`)
		v = strings.ReplaceAll(v, `\\`, `\`)
	}
	return v
}

// approaches parses the two approach subfolders into rows: prose technique
// deep-dives under approaches/techniques/ and pure YAML-frontmatter records
// under approaches/tools/ (filename = id; the subfolder decides the format,
// kind: stays a semantic label). Goals and best_when always come from the
// playbooks tables.
func (g *gen) approaches() map[string]*row {
	rows := map[string]*row{}
	isRecord := map[string]bool{}

	techFiles, _ := filepath.Glob(filepath.Join(g.skill, "approaches", "techniques", "*.md"))
	for _, f := range techFiles {
		id := strings.TrimSuffix(filepath.Base(f), ".md")
		r := &row{id: id, kind: "technique"}
		rows[id] = r
		inSignals := false
		for _, l := range lines(f) {
			switch {
			case reRegKind.MatchString(l):
				g.errf(f, "technique file carries a kind: line — records live in approaches/tools/")
			case strings.HasPrefix(l, "## "):
				inSignals = l == "## Signals"
			case inSignals:
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

	recFiles, _ := filepath.Glob(filepath.Join(g.skill, "approaches", "tools", "*.md"))
	for _, f := range recFiles {
		id := strings.TrimSuffix(filepath.Base(f), ".md")
		if rows[id] != nil {
			g.errf(f, "id '%s' exists in both approaches/techniques/ and approaches/tools/ — one capability, one file", id)
			continue
		}
		r := &row{id: id, setupSig: "—"}
		rows[id] = r
		isRecord[id] = true
		for _, l := range lines(f) {
			switch {
			case reRegKind.MatchString(l):
				r.kind = reRegKind.FindStringSubmatch(l)[1]
			case reRegGoal.MatchString(l), reRegBest.MatchString(l):
				// Every record is ranked like a technique: goals and best_when
				// derive from the playbooks rows — inline copies would be a
				// second home for the fact.
				g.errf(f, "record carries inline goals:/best_when: — these derive from its ranked rows; remove them")
			case reRegSig.MatchString(l):
				r.sessionSig = unquote(reRegSig.FindStringSubmatch(l)[1])
			}
		}
		if r.kind == "" {
			g.errf(f, "record is missing kind:")
		}
		if r.sessionSig == "" {
			g.errf(f, "record is missing session_signal")
		}
	}

	// goals and best_when come from the playbooks tables — for every ranked approach
	bestRank := map[string]int{}
	playbooks, _ := filepath.Glob(filepath.Join(g.skill, "playbooks", "*.md"))
	slices.Sort(playbooks)
	for _, f := range playbooks {
		goal := strings.TrimSuffix(filepath.Base(f), ".md")
		for _, l := range lines(f) {
			m := reRow.FindStringSubmatch(l)
			if m == nil {
				continue
			}
			rank, _ := strconv.Atoi(m[1])
			slug := m[3]
			r := rows[slug]
			if r == nil {
				g.errf(f, "ranked row references a missing approach file for '%s'", slug)
				continue
			}
			cs := cells(l) // | # | Solution | Best when | Why it fits |
			if len(cs) < 4 {
				g.errf(f, "row for '%s' has too few columns", slug)
				continue
			}
			if bestRank[slug] == 0 {
				bestRank[slug] = 1 << 30
			}
			r.goals = append(r.goals, goal)
			if rank < bestRank[slug] {
				bestRank[slug] = rank
				r.bestWhen = strings.ToLower(cs[3][:1]) + cs[3][1:]
			}
		}
	}
	for id, r := range rows {
		if len(r.goals) == 0 {
			g.errf(filepath.Join(g.skill, "approaches"), "'"+id+"' (%s) has no ranked row in any playbooks file — cannot index it", r.kind)
		}
	}
	return rows
}

func (g *gen) build() string {
	rows := g.approaches()

	ids := make([]string, 0, len(rows))
	for id := range rows {
		ids = append(ids, id)
	}
	slices.Sort(ids)

	var b strings.Builder
	b.WriteString("# Approaches Index\n")
	b.WriteString("*Generated — do not edit. Sources: the playbooks tables and each approaches/ file. Regenerate: `go -C tools/approaches-index run .`*\n\n")
	b.WriteString("One row per first-party approach (unpromoted marketplace plugins live in `marketplace.md`). Setup signals are re-checkable disk evidence; session signals are conversation evidence accumulated in the profile. `—` means no signal of that tier exists.\n\n")
	b.WriteString("| Id | Kind | Goals | Best when | Setup signal | Session signal |\n")
	b.WriteString("|----|------|-------|-----------|--------------|----------------|\n")
	for _, id := range ids {
		r := rows[id]
		goals := dedup(r.goals)
		slices.Sort(goals)
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n",
			r.id, r.kind, strings.Join(goals, ", "), r.bestWhen, r.setupSig, r.sessionSig)
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

	target := filepath.Join(g.skill, "approaches", "index.md")
	if *check {
		disk, err := os.ReadFile(target)
		if err != nil || string(disk) != out {
			fmt.Println("approaches/index.md is stale — regenerate with `go -C tools/approaches-index run .`")
			os.Exit(1)
		}
		fmt.Println("Approaches index: fresh")
		return
	}
	if err := os.WriteFile(target, []byte(out), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("Wrote %s (%d approaches)\n", target, strings.Count(out, "\n|")-2)
}
