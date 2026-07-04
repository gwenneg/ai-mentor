// Check references/official-plugins.md against the live official marketplace.
//
// Pure GitHub API diffing — no LLM. Exits 1 on drift so a scheduled workflow
// can open an issue or feed the diff to the maintenance run, 2 on fetch or
// setup errors. Uses GITHUB_TOKEN / GH_TOKEN if set (else unauthenticated,
// subject to rate limits). Stdlib only.
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
	"sort"
	"strings"
	"time"
)

const api = "https://api.github.com/repos/anthropics/claude-plugins-official/contents"

var (
	reToken     = regexp.MustCompile("`([a-z0-9-]+)`")
	reMultiWord = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)+$`)
)

// fetchNames returns the directory names under one marketplace directory.
func fetchNames(client *http.Client, dir string) ([]string, error) {
	req, err := http.NewRequest("GET", api+"/"+dir, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GH_TOKEN")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", dir, resp.Status)
	}
	var entries []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.Type == "dir" {
			names = append(names, e.Name)
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
	seen := map[string]bool{}
	var out []string
	for _, x := range xs {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	sort.Strings(out)
	return out
}

// report prints the diff between the live and documented sets and returns
// whether real drift (live plugins missing from the catalog) was found.
func report(w io.Writer, live, documented []string) bool {
	live, documented = sortedUnique(live), sortedUnique(documented)
	doc := map[string]bool{}
	for _, d := range documented {
		doc[d] = true
	}
	inLive := map[string]bool{}
	var missing []string
	inBoth := 0
	for _, l := range live {
		inLive[l] = true
		if doc[l] {
			inBoth++
		} else {
			missing = append(missing, l)
		}
	}
	// Documented multi-word kebab tokens not in the marketplace (may be prose tokens)
	var removed []string
	for _, d := range documented {
		if !inLive[d] && reMultiWord.MatchString(d) {
			removed = append(removed, d)
		}
	}

	fmt.Fprintf(w, "Marketplace plugins: %d; documented names found: %d\n", len(live), inBoth)
	if len(missing) > 0 {
		fmt.Fprintf(w, "\nNEW plugins not yet in the catalog:\n")
		for _, m := range missing {
			fmt.Fprintf(w, "  + %s\n", m)
		}
	}
	if len(removed) > 0 {
		fmt.Fprintf(w, "\nDocumented names not in the marketplace (verify manually — may be prose tokens):\n")
		for _, r := range removed {
			fmt.Fprintf(w, "  ? %s\n", r)
		}
	}
	if len(missing) > 0 {
		fmt.Fprintf(w, "\nDrift detected: run the maintenance skill's catalog sync (step 5).\n")
		return true
	}
	fmt.Fprintf(w, "\nPlugin catalog: in sync.\n")
	return false
}

// findRoot walks upward from dir to the first directory containing
// skills/mentor, so the check works from anywhere in the repo — including
// tools/catalog-drift itself, where `go -C tools/catalog-drift run .` lands.
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

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
	os.Exit(2)
}

func main() {
	repo, err := findRoot(".")
	if err != nil {
		fatal(err)
	}
	catalog, err := os.ReadFile(filepath.Join(repo, "skills", "mentor", "references", "official-plugins.md"))
	if err != nil {
		fatal(err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	var live []string
	for _, dir := range []string{"plugins", "external_plugins"} {
		names, err := fetchNames(client, dir)
		if err != nil {
			fatal(err)
		}
		live = append(live, names...)
	}

	var out strings.Builder
	drift := report(&out, live, documentedNames(string(catalog)))
	fmt.Print(out.String())
	if drift {
		os.Exit(1)
	}
}
