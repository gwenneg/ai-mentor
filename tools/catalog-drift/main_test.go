package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestFetchLiveNamesParsesManifest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"name":"claude-plugins-official","plugins":[
			{"name":"alpha-one","source":"./plugins/alpha-one"},
			{"name":"partner-tool","source":{"source":"github","repo":"partner/tool"}},
			{"name":""}]}`))
	}))
	defer srv.Close()
	names, paths, err := fetchLiveNames(&http.Client{Timeout: 5 * time.Second}, srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "alpha-one,partner-tool"
	if strings.Join(names, ",") != want {
		t.Errorf("fetchLiveNames = %v, want %s (external-source plugins must be included, empty names dropped)", names, want)
	}
	if paths["alpha-one"] != "plugins/alpha-one" {
		t.Errorf("in-repo source must map to its repo-relative path, got %q", paths["alpha-one"])
	}
	if _, ok := paths["partner-tool"]; ok {
		t.Error("an external (non-string) source must have no path entry — its content is not in the marketplace repo")
	}
}

func TestFetchLiveNamesEmptyManifestIsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"plugins":[]}`))
	}))
	defer srv.Close()
	if _, _, err := fetchLiveNames(&http.Client{Timeout: 5 * time.Second}, srv.URL); err == nil {
		t.Error("empty plugin list should be an error, not an in-sync verdict")
	}
}

func TestPluginNames(t *testing.T) {
	// Only names the catalog structurally declares: the first backticked token
	// of each table row (a "(Vendor)" suffix is tolerated), and backticked
	// tokens in the two prose-list sections. Goal slugs in later cells and
	// tokens in intro prose must NOT be counted.
	text := strings.Join([]string{
		"Intro prose mentioning `ui5-modernization` as an example — not a declaration.",
		"### Dev workflow",
		"| Plugin | What it does | Relevant goal | Verdict |",
		"|--------|-------------|--------------|---------|",
		"| `security-guidance` | hooks | `security` | ✅ |",
		"| `context7` (Upstash) | docs | `code-understanding` | ✅ |",
		"### Language servers (LSPs)",
		"Drop-in: `gopls-lsp` (Go), `pyright-lsp` (Python).",
		"### Specialty",
		"Listed for completeness: `math-olympiad` and `cwc-makers`.",
		"### Database",
		"| `alloydb` | db | Google | `devops` | ☑️ |",
	}, "\n")
	got := pluginNames(text)
	want := []string{"security-guidance", "context7", "gopls-lsp", "pyright-lsp", "math-olympiad", "cwc-makers", "alloydb"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("pluginNames = %v, want %v", got, want)
	}
	for _, slug := range []string{"security", "code-understanding", "devops", "ui5-modernization"} {
		if slices.Contains(got, slug) {
			t.Errorf("goal slug / intro token %q must not be counted as a plugin: %v", slug, got)
		}
	}
}

func TestReportInSync(t *testing.T) {
	var out strings.Builder
	drift := report(&out, []string{"alpha-one", "beta-two"}, []string{"beta-two", "alpha-one"})
	if drift {
		t.Error("identical live and documented sets must not be drift")
	}
	s := out.String()
	if !strings.Contains(s, "Marketplace plugins: 2; documented plugins: 2") {
		t.Errorf("wrong summary:\n%s", s)
	}
	if !strings.Contains(s, "Plugin catalog: in sync.") {
		t.Errorf("missing in-sync line:\n%s", s)
	}
}

func TestReportDrift(t *testing.T) {
	var out strings.Builder
	drift := report(&out, []string{"alpha-one", "brand-new"}, []string{"alpha-one", "gone-away"})
	if !drift {
		t.Error("a missing or removed plugin must be drift")
	}
	s := out.String()
	if !strings.Contains(s, "NEW plugins not yet in the catalog:\n  + brand-new") {
		t.Errorf("missing NEW listing:\n%s", s)
	}
	if !strings.Contains(s, "no longer in the marketplace (renamed or removed):\n  - gone-away") {
		t.Errorf("missing removed listing:\n%s", s)
	}
	if !strings.Contains(s, "Drift detected") {
		t.Errorf("missing drift verdict:\n%s", s)
	}
}

// A live plugin whose name collides with a goal slug (e.g. a marketplace plugin
// literally named "debugging") must still be reported as drift — the false
// negative the precise parser closes.
func TestGoalSlugCollisionIsDrift(t *testing.T) {
	catalog := "### Dev workflow\n| `security-guidance` | hooks | `debugging` | ✅ |\n"
	documented := pluginNames(catalog)
	if slices.Contains(documented, "debugging") {
		t.Fatalf("goal slug 'debugging' must not be counted as a plugin: %v", documented)
	}
	var out strings.Builder
	if !report(&out, []string{"debugging", "security-guidance"}, documented) {
		t.Error("live plugin 'debugging' colliding with a goal slug must be drift")
	}
	if !strings.Contains(out.String(), "+ debugging") {
		t.Errorf("expected 'debugging' listed as NEW:\n%s", out.String())
	}
}

func TestReportDedupsAndSorts(t *testing.T) {
	var out strings.Builder
	report(&out, []string{"b-plugin", "a-plugin", "b-plugin"}, []string{})
	s := out.String()
	if !strings.Contains(s, "Marketplace plugins: 2;") {
		t.Errorf("live set should be deduplicated:\n%s", s)
	}
	if strings.Index(s, "+ a-plugin") > strings.Index(s, "+ b-plugin") {
		t.Errorf("missing plugins should be sorted:\n%s", s)
	}
}

// commitsServer serves the GitHub commits API shape, returning the given
// committer date for every path, and records the requests it saw.
func commitsServer(t *testing.T, date string, sawAuth *string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sawAuth != nil {
			*sawAuth = r.Header.Get("Authorization")
		}
		w.Write([]byte(`[{"commit":{"committer":{"date":"` + date + `"}}}]`))
	}))
}

func TestUpstreamStaleRecordIsDrift(t *testing.T) {
	srv := commitsServer(t, "2026-07-10T08:00:00Z", nil)
	defer srv.Close()
	records := []promotedRecord{{id: "alpha-one", lastVerified: "2026-07-03"}}
	paths := map[string]string{"alpha-one": "plugins/alpha-one"}
	findings, notes, err := checkUpstream(&http.Client{Timeout: 5 * time.Second}, srv.URL, "", records, paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notes) != 0 {
		t.Errorf("unexpected notes: %v", notes)
	}
	if len(findings) != 1 || findings[0].upstream != "2026-07-10" {
		t.Fatalf("upstream commit after last_verified must be a finding, got %v", findings)
	}
	var out strings.Builder
	if !reportUpstream(&out, 1, findings, notes) {
		t.Error("an upstream finding must be drift (gate property: printed issues fail the exit)")
	}
	if !strings.Contains(out.String(), "! alpha-one — verified 2026-07-03, upstream changed 2026-07-10") {
		t.Errorf("missing finding line:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "step 3") {
		t.Errorf("upstream drift must point at re-verification (step 3), not catalog sync:\n%s", out.String())
	}
}

func TestUpstreamFreshRecordNotFlagged(t *testing.T) {
	srv := commitsServer(t, "2026-06-19T00:46:00Z", nil)
	defer srv.Close()
	records := []promotedRecord{{id: "alpha-one", lastVerified: "2026-07-03"}}
	paths := map[string]string{"alpha-one": "plugins/alpha-one"}
	findings, _, err := checkUpstream(&http.Client{Timeout: 5 * time.Second}, srv.URL, "", records, paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("upstream commit before last_verified must not be flagged: %v", findings)
	}
	var out strings.Builder
	if reportUpstream(&out, 1, findings, nil) {
		t.Error("no findings must not be drift")
	}
	if !strings.Contains(out.String(), "hands-on evidence current") {
		t.Errorf("missing all-clear line:\n%s", out.String())
	}
}

// A commit on the verification day itself is not flagged — last_verified
// carries no time of day, so same-day order is unknowable. Documented blind
// spot: the next upstream commit flags the record.
func TestUpstreamSameDayNotFlagged(t *testing.T) {
	srv := commitsServer(t, "2026-07-03T23:59:59Z", nil)
	defer srv.Close()
	records := []promotedRecord{{id: "alpha-one", lastVerified: "2026-07-03"}}
	paths := map[string]string{"alpha-one": "plugins/alpha-one"}
	findings, _, err := checkUpstream(&http.Client{Timeout: 5 * time.Second}, srv.URL, "", records, paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("same-day upstream commit must not be flagged: %v", findings)
	}
}

func TestUpstreamExternalSourceIsNoteNotDrift(t *testing.T) {
	// No server: an external record must not trigger any API call.
	records := []promotedRecord{{id: "partner-tool", lastVerified: "2026-07-03"}}
	findings, notes, err := checkUpstream(&http.Client{Timeout: time.Second}, "http://127.0.0.1:0", "", records, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("uncheckable source must not be a finding: %v", findings)
	}
	if len(notes) != 1 || !strings.Contains(notes[0], "partner-tool") {
		t.Fatalf("uncheckable source must be a note: %v", notes)
	}
	var out strings.Builder
	if reportUpstream(&out, 1, findings, notes) {
		t.Error("notes alone must not be drift — a permanently external source would otherwise fail every run")
	}
	if !strings.Contains(out.String(), "? partner-tool") {
		t.Errorf("note must still be printed:\n%s", out.String())
	}
}

func TestCommitsAPIErrorPropagates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "rate limited", http.StatusForbidden)
	}))
	defer srv.Close()
	records := []promotedRecord{{id: "alpha-one", lastVerified: "2026-07-03"}}
	paths := map[string]string{"alpha-one": "plugins/alpha-one"}
	if _, _, err := checkUpstream(&http.Client{Timeout: 5 * time.Second}, srv.URL, "", records, paths); err == nil {
		t.Error("a commits API failure must be an error (exit 2), never a silent all-clear")
	}
}

func TestLatestCommitDateNoCommitsIsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()
	if _, err := latestCommitDate(&http.Client{Timeout: 5 * time.Second}, srv.URL, "", "plugins/gone"); err == nil {
		t.Error("an empty commit list means the path is wrong — must be an error, not a stale-free verdict")
	}
}

func TestLatestCommitDateSendsToken(t *testing.T) {
	var sawAuth string
	srv := commitsServer(t, "2026-07-01T00:00:00Z", &sawAuth)
	defer srv.Close()
	date, err := latestCommitDate(&http.Client{Timeout: 5 * time.Second}, srv.URL, "tok-123", "plugins/alpha-one")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if date != "2026-07-01" {
		t.Errorf("date = %q, want 2026-07-01 (UTC date only)", date)
	}
	if sawAuth != "Bearer tok-123" {
		t.Errorf("Authorization = %q, want Bearer tok-123", sawAuth)
	}
}

func TestPromotedRecordsParsesLastVerified(t *testing.T) {
	repo := t.TempDir()
	tools := filepath.Join(repo, "skills", "mentor", "approaches", "tools")
	if err := os.MkdirAll(tools, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(tools, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("alpha-one.md", "---\nkind: plugin\nlast_verified: 2026-07-03\n---\n")
	write("no-date.md", "---\nkind: plugin\n---\n")
	write("some-doc.md", "---\nkind: doc\nlast_verified: 2026-07-03\n---\n")

	records := promotedRecords(repo)
	if len(records) != 2 {
		t.Fatalf("want 2 plugin records (doc records excluded), got %v", records)
	}
	if records[0].id != "alpha-one" || records[0].lastVerified != "2026-07-03" {
		t.Errorf("alpha-one parsed wrong: %+v", records[0])
	}
	if records[1].id != "no-date" || records[1].lastVerified != "" {
		t.Errorf("missing last_verified must parse as empty: %+v", records[1])
	}
}
