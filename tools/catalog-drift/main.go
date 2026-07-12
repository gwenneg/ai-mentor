// Check the documented plugins (marketplace.md directory rows plus promoted
// kind: plugin records in approaches/tools/) against the live official marketplace.
//
// Two checks, no LLM. Membership: the marketplace.json manifest is the
// authoritative plugin list (it includes externally-hosted plugins that have
// no directory in the marketplace repo, so directory listings undercount);
// any name difference against the documented union is drift. Upstream drift:
// each promoted record carries hands-on claims dated by last_verified; if the
// plugin's path in the marketplace repo has a commit after that date, the
// evidence predates upstream changes and the record needs re-verification
// (maintenance step 3). Exits 1 on either kind of drift so a scheduled
// workflow can open an issue or feed the diff to the maintenance run, 2 on
// fetch or setup errors. Sends GITHUB_TOKEN as auth when set (the commits
// API is rate-limited unauthenticated). Stdlib only.
//
// Usage: go -C tools/catalog-drift run .
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

const (
	manifestURL     = "https://raw.githubusercontent.com/anthropics/claude-plugins-official/main/.claude-plugin/marketplace.json"
	githubAPIBase   = "https://api.github.com"
	marketplaceRepo = "anthropics/claude-plugins-official"
)

var (
	skillDir = filepath.Join("skills", "mentor")

	// A plugin id is the first backticked token of a catalog table row (the
	// name cell, which may carry a "(Vendor)" suffix) or a backticked token in
	// the prose-list sections. Dots are valid in names (e.g. wordpress.com).
	reRowName = regexp.MustCompile("^\\| `([a-z0-9.-]+)`")
	reTok     = regexp.MustCompile("`([a-z0-9.-]+)`")
)

// fetchLiveNames downloads the marketplace manifest and returns every listed
// plugin name — in-repo and externally-sourced alike — plus a map from name
// to its in-repo source path ("./plugins/x" → "plugins/x"). Plugins whose
// source is not a repo-relative string (external partner sources) have no
// entry in the map: their content lives outside the marketplace repo, so the
// upstream-drift check cannot see it.
func fetchLiveNames(client *http.Client, manifestURL string) ([]string, map[string]string, error) {
	resp, err := client.Get(manifestURL)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("GET marketplace manifest: %s", resp.Status)
	}
	var manifest struct {
		Plugins []struct {
			Name   string `json:"name"`
			Source any    `json:"source"`
		} `json:"plugins"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, nil, err
	}
	if len(manifest.Plugins) == 0 {
		return nil, nil, fmt.Errorf("manifest has no plugins — format change?")
	}
	var names []string
	paths := make(map[string]string)
	for _, p := range manifest.Plugins {
		if p.Name == "" {
			continue
		}
		names = append(names, p.Name)
		if s, ok := p.Source.(string); ok && strings.HasPrefix(s, "./") {
			paths[p.Name] = strings.TrimPrefix(s, "./")
		}
	}
	return names, paths, nil
}

// pluginNames extracts the plugin ids the catalog declares — nothing else, so
// goal slugs and command names backticked in prose or description cells are
// never mistaken for plugins. Two sources: the first backticked token of each
// table row, and every backticked token in the prose-list sections (Language
// servers, Specialty). Keep in sync with the copy in tools/structural-audit/main.go.
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
			for _, m := range reTok.FindAllStringSubmatch(line, -1) {
				names = append(names, m[1])
			}
		}
	}
	return names
}

func sortedUnique(xs []string) []string {
	return slices.Compact(slices.Sorted(slices.Values(xs)))
}

// report prints the diff between the live and documented plugin sets and
// returns whether drift was found: a live plugin missing from the catalog, or
// a documented plugin gone from the marketplace (renamed or removed).
func report(w io.Writer, live, documented []string) bool {
	live, documented = sortedUnique(live), sortedUnique(documented)

	var missing, removed []string
	for _, l := range live {
		if !slices.Contains(documented, l) {
			missing = append(missing, l)
		}
	}
	for _, d := range documented {
		if !slices.Contains(live, d) {
			removed = append(removed, d)
		}
	}

	fmt.Fprintf(w, "Marketplace plugins: %d; documented plugins: %d\n", len(live), len(documented))
	if len(missing) > 0 {
		fmt.Fprint(w, "\nNEW plugins not yet in the catalog:\n")
		for _, m := range missing {
			fmt.Fprintf(w, "  + %s\n", m)
		}
	}
	if len(removed) > 0 {
		fmt.Fprint(w, "\nDocumented plugins no longer in the marketplace (renamed or removed):\n")
		for _, r := range removed {
			fmt.Fprintf(w, "  - %s\n", r)
		}
	}
	if len(missing) > 0 || len(removed) > 0 {
		fmt.Fprint(w, "\nDrift detected: run the maintenance skill's catalog sync (step 5).\n")
		return true
	}
	fmt.Fprint(w, "\nPlugin catalog: in sync.\n")
	return false
}

// findRoot walks upward from dir to the first directory containing
// skills/mentor, so the check works from anywhere in the repo — including
// tools/catalog-drift itself, where `go -C tools/catalog-drift run .` lands.
// Keep in sync with the copy in tools/structural-audit/main.go.
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

var reLastVerified = regexp.MustCompile(`(?m)^last_verified: (\d{4}-\d{2}-\d{2})$`)

// promotedRecord is a `kind: plugin` file under approaches/tools/ — a
// marketplace plugin promoted out of the directory, whose hands-on claims
// are dated by its last_verified field.
type promotedRecord struct {
	id           string
	lastVerified string // YYYY-MM-DD; empty when the field is missing (the structural audit owns that invariant)
}

// promotedRecords returns the promoted plugins with their verification dates.
// They are documented plugins too: the membership check covers directory ∪ promoted.
func promotedRecords(repo string) []promotedRecord {
	var records []promotedRecord
	files, _ := filepath.Glob(filepath.Join(repo, skillDir, "approaches", "tools", "*.md"))
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		if !strings.Contains(string(b), "\nkind: plugin\n") {
			continue
		}
		rec := promotedRecord{id: strings.TrimSuffix(filepath.Base(f), ".md")}
		if m := reLastVerified.FindStringSubmatch(string(b)); m != nil {
			rec.lastVerified = m[1]
		}
		records = append(records, rec)
	}
	return records
}

// latestCommitDate returns the UTC date (YYYY-MM-DD) of the most recent
// commit touching path in the marketplace repo. token, when non-empty, is
// sent as a bearer token — unauthenticated calls share a small rate limit.
func latestCommitDate(client *http.Client, apiBase, token, path string) (string, error) {
	u := fmt.Sprintf("%s/repos/%s/commits?path=%s&per_page=1", apiBase, marketplaceRepo, url.QueryEscape(path))
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET commits for %s: %s", path, resp.Status)
	}
	var commits []struct {
		Commit struct {
			Committer struct {
				Date time.Time `json:"date"`
			} `json:"committer"`
		} `json:"commit"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return "", err
	}
	if len(commits) == 0 {
		return "", fmt.Errorf("no commits found for %s — path moved or renamed?", path)
	}
	return commits[0].Commit.Committer.Date.UTC().Format("2006-01-02"), nil
}

// upstreamFinding is a promoted record whose upstream plugin changed after
// the record's hands-on claims were last verified.
type upstreamFinding struct {
	id                     string
	lastVerified, upstream string
}

// checkUpstream compares each promoted record's last_verified date against
// the latest upstream commit to its plugin's path. Records whose manifest
// source is not a marketplace-repo path are returned as notes — their content
// is hosted elsewhere and must be checked by hand. A record whose upstream
// changed on the verification day itself is NOT flagged (dates carry no time
// of day); the next upstream commit will flag it.
func checkUpstream(client *http.Client, apiBase, token string, records []promotedRecord, paths map[string]string) ([]upstreamFinding, []string, error) {
	var findings []upstreamFinding
	var notes []string
	for _, rec := range records {
		path, ok := paths[rec.id]
		if !ok {
			notes = append(notes, fmt.Sprintf("%s — source is not a marketplace-repo path; upstream must be checked by hand", rec.id))
			continue
		}
		if rec.lastVerified == "" {
			notes = append(notes, fmt.Sprintf("%s — no last_verified field; the structural audit should be failing", rec.id))
			continue
		}
		upstream, err := latestCommitDate(client, apiBase, token, path)
		if err != nil {
			return nil, nil, err
		}
		if upstream > rec.lastVerified {
			findings = append(findings, upstreamFinding{id: rec.id, lastVerified: rec.lastVerified, upstream: upstream})
		}
	}
	return findings, notes, nil
}

// reportUpstream prints the upstream-drift verdict and returns whether any
// promoted record needs re-verification.
func reportUpstream(w io.Writer, checked int, findings []upstreamFinding, notes []string) bool {
	fmt.Fprintf(w, "\nPromoted records checked against upstream: %d\n", checked)
	for _, n := range notes {
		fmt.Fprintf(w, "  ? %s\n", n)
	}
	if len(findings) == 0 {
		fmt.Fprint(w, "Promoted records: hands-on evidence current.\n")
		return false
	}
	fmt.Fprint(w, "\nUpstream changed AFTER last verification (re-verify — maintenance step 3):\n")
	for _, f := range findings {
		fmt.Fprintf(w, "  ! %s — verified %s, upstream changed %s\n", f.id, f.lastVerified, f.upstream)
	}
	return true
}

func main() {
	repo, err := findRoot(".")
	if err != nil {
		fatal(err)
	}
	catalog, err := os.ReadFile(filepath.Join(repo, skillDir, "marketplace.md"))
	if err != nil {
		fatal(err)
	}
	promoted := promotedRecords(repo)
	documented := pluginNames(string(catalog))
	for _, rec := range promoted {
		documented = append(documented, rec.id)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	live, sourcePaths, err := fetchLiveNames(client, manifestURL)
	if err != nil {
		fatal(err)
	}
	membershipDrift := report(os.Stdout, live, documented)

	findings, notes, err := checkUpstream(client, githubAPIBase, os.Getenv("GITHUB_TOKEN"), promoted, sourcePaths)
	if err != nil {
		fatal(err)
	}
	upstreamDrift := reportUpstream(os.Stdout, len(promoted), findings, notes)

	if membershipDrift || upstreamDrift {
		os.Exit(1)
	}
}
