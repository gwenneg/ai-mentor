package main

import (
	"strings"
	"testing"
)

func TestDelta(t *testing.T) {
	for _, tc := range []struct {
		in   int
		want string
	}{{32, "+32"}, {-4, "-4"}, {0, "±0"}} {
		if got := delta(tc.in); got != tc.want {
			t.Errorf("delta(%d) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseCount(t *testing.T) {
	n, err := parseCount([]byte(`{"input_tokens": 2650}`))
	if err != nil || n != 2650 {
		t.Errorf("parseCount = %d, %v; want 2650, nil", n, err)
	}
	if _, err := parseCount([]byte(`{"type":"error"}`)); err == nil {
		t.Error("parseCount accepted a response with no input_tokens")
	}
	if _, err := parseCount([]byte(`not json`)); err == nil {
		t.Error("parseCount accepted invalid JSON")
	}
}

func TestRenderWithBase(t *testing.T) {
	rows := []row{
		{name: budgetFiles[0], base: 2650, head: 2682},
		{name: budgetFiles[1], base: 4100, head: 4096},
		{name: budgetFiles[2], base: 700, head: 700},
	}
	md := render(rows, "claude-opus-4-8", "origin/main")
	for _, want := range []string{
		"| skills/mentor/SKILL.md | 2650 | 2682 | +32 |",
		"| skills/mentor/problem-mode.md | 4100 | 4096 | -4 |",
		"| skills/mentor/growth-mode.md | 700 | 700 | ±0 |",
		"| **Total** | **7450** | **7478** | **+28** |",
		"problem mode (SKILL.md + problem-mode.md): **6778** (+28)",
		"growth mode (SKILL.md + growth-mode.md): **3382** (+32)",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("render output missing %q\n---\n%s", want, md)
		}
	}
}

func TestRenderWithoutBase(t *testing.T) {
	rows := []row{
		{name: budgetFiles[0], base: -1, head: 2682},
		{name: budgetFiles[1], base: -1, head: 4096},
		{name: budgetFiles[2], base: -1, head: 700},
	}
	md := render(rows, "claude-opus-4-8", "")
	for _, want := range []string{
		"| skills/mentor/SKILL.md | 2682 |",
		"| **Total** | **7478** |",
		"problem mode (SKILL.md + problem-mode.md): **6778**",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("render output missing %q\n---\n%s", want, md)
		}
	}
	if strings.Contains(md, "Δ") {
		t.Error("render without base must not show a delta column")
	}
}
