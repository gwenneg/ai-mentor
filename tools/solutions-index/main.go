// Generates skills/mentor/solutions/index.md — the compiled solutions index —
// from the authored sources of truth:
//
//   - problems/*.md    goal membership, rank, and best-when triggers for
//     technique solutions (the ranked rows)
//   - solutions/*.md   every solution, one file each. A file with a `kind:`
//     line is a flat record (built-in command, integration, or doc — the
//     filename is its id); any other file is a technique deep-dive with a
//     "## Signals" section.
//
// index.md is a build artifact: never hand-edit it. After editing any source
// above, regenerate with `go -C tools/solutions-index run .`. In CI, `-check`
// regenerates in memory and exits 1 if the on-disk file is stale.
//
// Deterministic by construction: rows sorted by id, no timestamps. Exits 1 on
// any inconsistency in the sources (a technique without a ranked row, a
// missing Signals section, a ranked row pointing at a record), 2 on a fatal
// setup problem. No network, no LLM. Stdlib only.
//
// Usage: go -C tools/solutions-index run . [-check] [repo-root]
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
	reRow     = regexp.MustCompile(`^\| (\d+) \| \[([^\]]+)\]\(\.\./solutions/([a-z0-9-]+)\.md\)`)
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

// solutions parses every solutions/*.md file into a row. A `kind:` line makes
// the file a flat record (filename = id); otherwise it is a technique
// deep-dive whose goals and best_when come from the problems tables.
func (g *gen) solutions() map[string]*row {
	rows := map[string]*row{}
	isRecord := map[string]bool{}

	files, _ := filepath.Glob(filepath.Join(g.skill, "solutions", "*.md"))
	for _, f := range files {
		id := strings.TrimSuffix(filepath.Base(f), ".md")
		if id == "index" {
			continue // the build artifact itself
		}
		r := &row{id: id}
		rows[id] = r
		inSignals := false
		for _, l := range lines(f) {
			switch {
			case reRegKind.MatchString(l):
				r.kind = reRegKind.FindStringSubmatch(l)[1]
				isRecord[id] = true
			case reRegGoal.MatchString(l):
				for _, gl := range strings.Split(reRegGoal.FindStringSubmatch(l)[1], ",") {
					r.goals = append(r.goals, strings.TrimSpace(gl))
				}
			case reRegBest.MatchString(l):
				r.bestWhen = reRegBest.FindStringSubmatch(l)[1]
			case reRegSig.MatchString(l):
				r.sessionSig = reRegSig.FindStringSubmatch(l)[1]
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
		if isRecord[id] {
			if len(r.goals) == 0 || r.bestWhen == "" || r.sessionSig == "" {
				g.errf(f, "record '%s' is missing one of goals/best_when/session_signal", id)
			}
			if r.setupSig == "" {
				r.setupSig = "—"
			}
		} else {
			r.kind = "technique"
			if r.setupSig == "" || r.sessionSig == "" {
				g.errf(f, "missing or incomplete '## Signals' section (need '- Setup:' and '- Session:' lines)")
			}
		}
	}

	// technique goals and best_when come from the problems tables
	bestRank := map[string]int{}
	problems, _ := filepath.Glob(filepath.Join(g.skill, "problems", "*.md"))
	slices.Sort(problems)
	for _, f := range problems {
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
				g.errf(f, "ranked row references solutions/%s.md, which does not exist", slug)
				continue
			}
			if isRecord[slug] {
				g.errf(f, "ranked row ranks '%s', a record — only techniques are ranked; records ride the Built-ins/Integrations lines", slug)
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
		if !isRecord[id] && len(r.goals) == 0 {
			g.errf(filepath.Join(g.skill, "solutions", id+".md"), "technique has no ranked row in any problems file — cannot index it")
		}
	}
	return rows
}

func (g *gen) build() string {
	rows := g.solutions()

	ids := make([]string, 0, len(rows))
	for id := range rows {
		ids = append(ids, id)
	}
	slices.Sort(ids)

	var b strings.Builder
	b.WriteString("# Solutions Index\n")
	b.WriteString("*Generated — do not edit. Sources: the problems tables and each solutions/ file. Regenerate: `go -C tools/solutions-index run .`*\n\n")
	b.WriteString("One row per first-party solution (marketplace plugins live in `plugins.md`). Setup signals are re-checkable disk evidence; session signals are conversation evidence accumulated in the profile. `—` means no signal of that tier exists.\n\n")
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

	target := filepath.Join(g.skill, "solutions", "index.md")
	if *check {
		disk, err := os.ReadFile(target)
		if err != nil || string(disk) != out {
			fmt.Println("solutions/index.md is stale — regenerate with `go -C tools/solutions-index run .`")
			os.Exit(1)
		}
		fmt.Println("Solutions index: fresh")
		return
	}
	if err := os.WriteFile(target, []byte(out), 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("Wrote %s (%d solutions)\n", target, strings.Count(out, "\n|")-2)
}
