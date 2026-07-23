// V2 deterministic layer ("grade forms, not essays"): structural checks
// against the machine-expectations table in cases.md and the response
// trailer. These verdicts are code, not judgment — a failure here is exact,
// reproducible, and skips the judge, exactly like the older det checks.
package main

import (
	"fmt"
	"regexp"
	"strings"
)

// v2Spec is one row of cases.md's "## Machine expectations" table.
type v2Spec struct {
	Goals    []string // accepted trailer goal values; empty = unconstrained
	Move     string   // required move id, "!id" = forbidden id, "" = unconstrained
	Surprise string   // "required", "omitted-ok", or ""
	Fence    string   // grounded | portable | either | setup | scan | none | ""
	Judge    bool     // substance stays LLM-judged (A30/B05/B06 class)
}

// escaped-pipe placeholder: goal cells hold alternatives like a\|b.
const pipePlaceholder = "\x00PIPE\x00"

// parseV2Specs reads the machine-expectations table. Loud on drift: the
// caller must verify one row per headless case and no unknown IDs.
func parseV2Specs(text string) (map[string]v2Spec, error) {
	idx := strings.Index(text, "## Machine expectations")
	if idx < 0 {
		return nil, fmt.Errorf("cases.md: '## Machine expectations' section missing — the V2 deterministic layer has no contract")
	}
	specs := map[string]v2Spec{}
	for _, line := range strings.Split(text[idx:], "\n") {
		line = strings.TrimSpace(strings.ReplaceAll(line, `\|`, pipePlaceholder))
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cs := strings.Split(strings.Trim(line, "|"), "|")
		if len(cs) < 6 {
			continue
		}
		id := strings.TrimSpace(cs[0])
		if !reCaseID.MatchString(id) {
			continue // header / separator rows
		}
		cell := func(i int) string {
			v := strings.TrimSpace(strings.ReplaceAll(cs[i], pipePlaceholder, "|"))
			if v == "-" {
				return ""
			}
			return v
		}
		spec := v2Spec{Move: cell(2), Surprise: cell(3), Fence: cell(4), Judge: cell(5) == "judge"}
		if g := cell(1); g != "" {
			spec.Goals = strings.Split(g, "|")
		}
		specs[id] = spec
	}
	if len(specs) == 0 {
		return nil, fmt.Errorf("cases.md: machine-expectations table parsed to zero rows — format drift?")
	}
	return specs, nil
}

var (
	reCaseID       = regexp.MustCompile(`^[ABC]\d{2}$`)
	reFence        = regexp.MustCompile("(?s)```[a-zA-Z]*\n(.*?)```")
	reSurpriseMark = regexp.MustCompile(`(?i)one (?:more )?thing you might not know`)
	rePluginLine   = regexp.MustCompile(`/plugin install ([a-zA-Z0-9._-]+)@`)
	reSetupCmd     = regexp.MustCompile(`(?i)claude mcp add|mcp connect|/mcp\b`)
	reSurface      = regexp.MustCompile(`(?i)server\.go|orders|handler`)
	reClosing      = regexp.MustCompile(`(?i)say "more"`)
	// flagged-for-replacement references are legal in portable fences
	// (maintainer ruling 2026-07-21).
	reFlaggedRef = regexp.MustCompile(`(?i)(adjust|swap|replace|change)[^\n]*|this repo'?s`)
)

// fixture-specific tokens: presence grounds a fence; presence in a portable
// fence (outside a flagged-reference line) is an import violation. Generic
// artifacts every repo has (go.mod, CLAUDE.md) ground but never violate.
var fixtureOnlyTokens = []string{"orders.go", "orders_test.go", "server.go", "orders-service"}
var groundingTokens = append([]string{"go.mod", "CLAUDE.md", "go test"}, fixtureOnlyTokens...)

func trailerFields(trailer string) map[string]string {
	f := map[string]string{}
	for _, kv := range strings.Fields(trailer) {
		if k, v, ok := strings.Cut(kv, "="); ok {
			f[k] = v
		}
	}
	return f
}

// v2Checks runs every structural check for one case. The first return is
// GATING (identity/promise tier — measured ~99%+ compliant, a failure is a
// real defect); the second is ADVISORY (the fence-discipline tier, whose
// newly-measured true rates of 40-70% would red every gate if enforced
// per-run — they record into results and rate-gate in Phase 3).
// responses carries each subject run's raw text (C02 has two).
func v2Checks(c evalCase, spec v2Spec, responses []string, plugins, promoted []string) (string, string) {
	if c.Group == "D" {
		return "", ""
	}
	last := responses[len(responses)-1]
	fields := trailerFields(parseTrailer(last))
	if len(fields) == 0 {
		return "trailer: missing or unparseable — the response's self-report is the deterministic contract", ""
	}

	// mode
	wantMode := "problem"
	if c.Group == "B" {
		wantMode = "growth"
	}
	if got := fields["mode"]; got != wantMode {
		return fmt.Sprintf("trailer mode: got %q, want %q", got, wantMode), ""
	}

	// goal
	if len(spec.Goals) > 0 {
		ok := false
		for _, g := range spec.Goals {
			if fields["goal"] == g {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Sprintf("classification: trailer goal %q not in expected {%s}", fields["goal"], strings.Join(spec.Goals, ", ")), ""
		}
	}

	// move identity
	if spec.Move != "" {
		if forbidden, isNeg := strings.CutPrefix(spec.Move, "!"); isNeg {
			if fields["move"] == forbidden {
				return fmt.Sprintf("move identity: %q is excluded for this case", forbidden), ""
			}
		} else if fields["move"] != spec.Move {
			return fmt.Sprintf("move identity: got %q, want %q", fields["move"], spec.Move), ""
		}
	}

	// growth-mode opener/taught expectations
	if c.Group == "B" {
		if reason := b06StyleChecks(c.ID, fields); reason != "" {
			return reason, ""
		}
	}

	// surprise count via the prose marker + trailer agreement
	marks := len(reSurpriseMark.FindAllString(stripTrailer(last), -1))
	switch spec.Surprise {
	case "required":
		if marks != 1 {
			return fmt.Sprintf("surprise: want exactly one labeled pick, found %d markers", marks), ""
		}
		if fields["surprise"] == "omitted" || fields["surprise"] == "" {
			return "surprise: prose carries a pick but the trailer declares none", ""
		}
	case "omitted-ok":
		if marks > 1 {
			return fmt.Sprintf("surprise: at most one labeled pick allowed, found %d", marks), ""
		}
	}

	// C02: the two runs' surprises must differ (never-repeat across runs)
	if c.ID == "C02" && len(responses) == 2 {
		s1 := trailerFields(parseTrailer(responses[0]))["surprise"]
		s2 := fields["surprise"]
		if s1 != "" && s1 == s2 {
			return fmt.Sprintf("never-repeat: both runs picked the same surprise (%s)", s1), ""
		}
	}

	// fence discipline: ADVISORY tier — recorded, rate-gated in Phase 3.
	advisory := fenceChecks(spec.Fence, stripTrailer(last))

	// closing line last (trailer exempt via stripTrailer); not for fence=none cases
	if spec.Fence != "none" && spec.Fence != "" {
		body := strings.TrimSpace(stripTrailer(last))
		tail := body
		if len(body) > 300 {
			tail = body[len(body)-300:]
		}
		if !reClosing.MatchString(tail) {
			return `closing line: the response must end with the single closing line (say "more" + calibration)`, ""
		}
	}

	// plugin fabrication + directory-plugin labels, on every response
	for _, resp := range responses {
		if reason := pluginChecks(resp, plugins, promoted); reason != "" {
			return reason, ""
		}
	}
	return "", advisory
}

func b06StyleChecks(id string, fields map[string]string) string {
	expect := map[string][]string{ // opener expectations; B05 is judge-owned
		"B01": {"lesson"}, "B02": {"followup"}, "B03": {"lesson"}, "B04": {"lesson"}, "B06": {"transfer", "empty"},
	}
	want, ok := expect[id]
	if !ok {
		return ""
	}
	got := fields["opener"]
	legal := false
	for _, w := range want {
		if got == w {
			legal = true
			break
		}
	}
	if !legal {
		return fmt.Sprintf("growth opener: got %q, want one of %v", got, want)
	}
	if id == "B04" && fields["taught"] != "hooks-as-workflow" {
		return fmt.Sprintf("growth lesson: got taught=%q, want hooks-as-workflow (the repo's own configured signal)", fields["taught"])
	}
	// A transfer opener legitimately carries the transferred (already
	// adopted) capability in taught= — only the empty-map answer must
	// teach nothing. (The trailer proved more honest than this check's
	// first version, 2026-07-24.)
	if id == "B06" && fields["opener"] == "empty" && fields["taught"] != "none" && fields["taught"] != "" {
		return fmt.Sprintf("growth lesson: the empty-map answer must teach nothing new, got taught=%q", fields["taught"])
	}
	return ""
}

// reSetupOnly matches a fence that is purely a setup line — the product
// tolerates these as separate blocks in practice; the one-fence rule's
// intent is one MOVE prompt, not zero setup visibility.
var reSetupOnly = regexp.MustCompile(`(?s)^\s*(/plugin install \S+|claude mcp add[^\n]*|claude --\S+)\s*$`)

func fenceChecks(mode, body string) string {
	all := reFence.FindAllStringSubmatch(body, -1)
	switch mode {
	case "", "none":
		return ""
	}
	var candidates []string
	if mode == "setup" {
		// In a setup-shaped move the setup command IS the substantive
		// fence — no filtering, exactly one block expected.
		for _, m := range all {
			candidates = append(candidates, m[1])
		}
	} else {
		for _, m := range all {
			if !reSetupOnly.MatchString(m[1]) {
				candidates = append(candidates, m[1])
			}
		}
	}
	if len(candidates) != 1 {
		return fmt.Sprintf("fence: want exactly one substantive fenced block (setup-only blocks aside), found %d of %d total", len(candidates), len(all))
	}
	fence := candidates[0]
	hasFixtureToken := func(s string) (string, bool) {
		for _, t := range fixtureOnlyTokens {
			if strings.Contains(s, t) {
				return t, true
			}
		}
		return "", false
	}
	switch mode {
	case "grounded", "scan":
		grounded := false
		for _, t := range groundingTokens {
			if strings.Contains(fence, t) {
				grounded = true
				break
			}
		}
		if !grounded {
			return "fence grounding: no verified fixture path or command inside the fenced block"
		}
		if mode == "scan" && !strings.Contains(fence, "server.go") {
			return "scan canary: the grounded fence must name server.go (only a real repo scan surfaces it)"
		}
	case "portable":
		for _, line := range strings.Split(fence, "\n") {
			if t, found := hasFixtureToken(line); found && !reFlaggedRef.MatchString(line) {
				return fmt.Sprintf("portability: fixture-specific %q imported into a portable fence without a replacement flag", t)
			}
		}
	case "setup":
		if !reSetupCmd.MatchString(fence) {
			return "setup move: the fence must carry the concrete setup command (e.g. claude mcp add ...)"
		}
		if !reSurface.MatchString(fence) {
			return "setup move: the fence must name the service/suspect surface alongside the setup command"
		}
	}
	return ""
}

func pluginChecks(resp string, plugins, promoted []string) string {
	known := map[string]bool{}
	for _, p := range plugins {
		known[p] = true
	}
	isPromoted := map[string]bool{}
	for _, p := range promoted {
		isPromoted[p] = true
	}
	for _, m := range rePluginLine.FindAllStringSubmatch(resp, -1) {
		name := m[1]
		if !known[name] {
			return fmt.Sprintf("plugin fabrication: %q is not in the marketplace whitelist", name)
		}
		if isPromoted[name] {
			continue
		}
		// directory plugin recommended by install line: the user-facing label
		// must appear near a mention (marker alone is taxonomy shorthand).
		idx := strings.Index(resp, name)
		lo, hi := max(0, idx-250), min(len(resp), idx+len(name)+250)
		if !strings.Contains(strings.ToLower(resp[lo:hi]), "not hands-on evaluated") {
			return fmt.Sprintf("tier label: directory plugin %q recommended without the literal \"not hands-on evaluated\" label nearby", name)
		}
	}
	return ""
}
