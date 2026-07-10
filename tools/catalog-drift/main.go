// Check the documented plugins (marketplace.md directory rows plus promoted
// kind: plugin records in approaches/records/) against the live official marketplace.
//
// Pure manifest diffing — no LLM. Reads the marketplace.json manifest, the
// authoritative plugin list: it includes externally-hosted plugins that have
// no directory in the marketplace repo, so directory listings undercount.
// Exits 1 on drift so a scheduled workflow can open an issue or feed the
// diff to the maintenance run, 2 on fetch or setup errors. Stdlib only.
//
// Usage: go -C tools/catalog-drift run .
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

const manifestURL = "https://raw.githubusercontent.com/anthropics/claude-plugins-official/main/.claude-plugin/marketplace.json"

var (
	skillDir = filepath.Join("skills", "mentor")

	// A plugin id is the first backticked token of a catalog table row (the
	// name cell, which may carry a "(Vendor)" suffix) or a backticked token in
	// the prose-list sections. Dots are valid in names (e.g. wordpress.com).
	reRowName = regexp.MustCompile("^\\| `([a-z0-9.-]+)`")
	reTok     = regexp.MustCompile("`([a-z0-9.-]+)`")
)

// fetchLiveNames downloads the marketplace manifest and returns every listed
// plugin name — in-repo and externally-sourced alike.
func fetchLiveNames(client *http.Client, url string) ([]string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET marketplace manifest: %s", resp.Status)
	}
	var manifest struct {
		Plugins []struct {
			Name string `json:"name"`
		} `json:"plugins"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, err
	}
	if len(manifest.Plugins) == 0 {
		return nil, fmt.Errorf("manifest has no plugins — format change?")
	}
	var names []string
	for _, p := range manifest.Plugins {
		if p.Name != "" {
			names = append(names, p.Name)
		}
	}
	return names, nil
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

// promotedPlugins returns the ids of `kind: plugin` records under
// approaches/ — marketplace plugins promoted out of the directory. They are
// documented plugins too: the drift check covers directory ∪ promoted.
func promotedPlugins(repo string) []string {
	var ids []string
	files, _ := filepath.Glob(filepath.Join(repo, skillDir, "approaches", "records", "*.md"))
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		if strings.Contains(string(b), "\nkind: plugin\n") {
			ids = append(ids, strings.TrimSuffix(filepath.Base(f), ".md"))
		}
	}
	return ids
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
	documented := append(pluginNames(string(catalog)), promotedPlugins(repo)...)

	client := &http.Client{Timeout: 30 * time.Second}
	live, err := fetchLiveNames(client, manifestURL)
	if err != nil {
		fatal(err)
	}

	if report(os.Stdout, live, documented) {
		os.Exit(1)
	}
}
