package main

import (
	"strings"
	"testing"
)

func TestDocumentedNames(t *testing.T) {
	text := "The `code-review` and `security-guardian` plugins, plus prose like `pr` and `multi-word-token`."
	got := documentedNames(text)
	want := []string{"code-review", "security-guardian", "pr", "multi-word-token"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("documentedNames = %v, want %v", got, want)
	}
}

func TestReportInSync(t *testing.T) {
	var out strings.Builder
	drift := report(&out, []string{"alpha-one", "beta-two"}, []string{"alpha-one", "beta-two", "prose", "stale-token"})
	if drift {
		t.Error("no missing plugins should not be drift")
	}
	s := out.String()
	if !strings.Contains(s, "Marketplace plugins: 2; documented names found: 2") {
		t.Errorf("wrong summary:\n%s", s)
	}
	if !strings.Contains(s, "  ? stale-token") {
		t.Errorf("multi-word stale token should be listed as advisory:\n%s", s)
	}
	if strings.Contains(s, "? prose") {
		t.Errorf("single-word token must not be listed as advisory:\n%s", s)
	}
	if !strings.Contains(s, "Plugin catalog: in sync.") {
		t.Errorf("missing in-sync line:\n%s", s)
	}
}

func TestReportDrift(t *testing.T) {
	var out strings.Builder
	drift := report(&out, []string{"alpha-one", "brand-new"}, []string{"alpha-one"})
	if !drift {
		t.Error("a live plugin missing from the catalog must be drift")
	}
	s := out.String()
	if !strings.Contains(s, "NEW plugins not yet in the catalog:\n  + brand-new") {
		t.Errorf("missing NEW listing:\n%s", s)
	}
	if !strings.Contains(s, "Drift detected") {
		t.Errorf("missing drift verdict:\n%s", s)
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
