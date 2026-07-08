package main

import (
	"net/http"
	"net/http/httptest"
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
	names, err := fetchLiveNames(&http.Client{Timeout: 5 * time.Second}, srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "alpha-one,partner-tool"
	if strings.Join(names, ",") != want {
		t.Errorf("fetchLiveNames = %v, want %s (external-source plugins must be included, empty names dropped)", names, want)
	}
}

func TestFetchLiveNamesEmptyManifestIsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"plugins":[]}`))
	}))
	defer srv.Close()
	if _, err := fetchLiveNames(&http.Client{Timeout: 5 * time.Second}, srv.URL); err == nil {
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
