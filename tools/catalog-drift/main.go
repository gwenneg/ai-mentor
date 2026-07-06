// Check references/official-plugins.md against the live official marketplace.
//
// Pure manifest diffing — no LLM. Reads the marketplace.json manifest, which
// lists every plugin including externally-hosted ones; listing the repo's
// plugin directories instead (the pre-2026-07-07 approach) silently missed
// all external-source entries. Exits 1 on drift so a scheduled workflow can
// open an issue or feed the diff to the maintenance run, 2 on fetch or setup
// errors. Stdlib only.
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
	"time"
)

const manifestURL = "https://raw.githubusercontent.com/anthropics/claude-plugins-official/main/.claude-plugin/marketplace.json"

var (
	skillDir = filepath.Join("skills", "mentor")

	reToken     = regexp.MustCompile("`([a-z0-9-]+)`")
	reMultiWord = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)+$`)
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

// documentedNames extracts every backticked kebab/word token from the catalog.
func documentedNames(text string) []string {
	var names []string
	for _, m := range reToken.FindAllStringSubmatch(text, -1) {
		names = append(names, m[1])
	}
	return names
}

func sortedUnique(xs []string) []string {
	return slices.Compact(slices.Sorted(slices.Values(xs)))
}

// report prints the diff between the live and documented sets and returns
// whether real drift (live plugins missing from the catalog) was found.
func report(w io.Writer, live, documented []string) bool {
	live, documented = sortedUnique(live), sortedUnique(documented)

	var missing []string
	for _, l := range live {
		if !slices.Contains(documented, l) {
			missing = append(missing, l)
		}
	}
	// Documented multi-word kebab tokens not in the marketplace (may be prose tokens)
	var removed []string
	for _, d := range documented {
		if !slices.Contains(live, d) && reMultiWord.MatchString(d) {
			removed = append(removed, d)
		}
	}

	fmt.Fprintf(w, "Marketplace plugins: %d; documented names found: %d\n", len(live), len(live)-len(missing))
	if len(missing) > 0 {
		fmt.Fprint(w, "\nNEW plugins not yet in the catalog:\n")
		for _, m := range missing {
			fmt.Fprintf(w, "  + %s\n", m)
		}
	}
	if len(removed) > 0 {
		fmt.Fprint(w, "\nDocumented names not in the marketplace (verify manually — may be prose tokens):\n")
		for _, r := range removed {
			fmt.Fprintf(w, "  ? %s\n", r)
		}
	}
	if len(missing) > 0 {
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

func main() {
	repo, err := findRoot(".")
	if err != nil {
		fatal(err)
	}
	catalog, err := os.ReadFile(filepath.Join(repo, skillDir, "references", "official-plugins.md"))
	if err != nil {
		fatal(err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	live, err := fetchLiveNames(client, manifestURL)
	if err != nil {
		fatal(err)
	}

	if report(os.Stdout, live, documentedNames(string(catalog))) {
		os.Exit(1)
	}
}
