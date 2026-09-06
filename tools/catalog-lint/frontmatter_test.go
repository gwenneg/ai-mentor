package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The mentor's profile grant lives in three frontmatter hook entries, not in
// allowed-tools: an allowed-tools grant clears on the user's next message,
// hooks stay registered for the session. These checks keep that shape from
// drifting — the three `if` rules stay exactly the profile paths, every
// hook command is a constant "allow" (no script ships with the plugin), and
// no `~/.ai-mentor` rule reappears under allowed-tools, so the grant has one
// source of truth.

func skillFrontmatter(t *testing.T) string {
	t.Helper()
	root, err := findRoot(".")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(root, skillDir, "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.SplitN(string(raw), "\n---\n", 2)
	if len(parts) != 2 || !strings.HasPrefix(parts[0], "---\n") {
		t.Fatal("SKILL.md has no frontmatter block")
	}
	return parts[0]
}

// section returns the lines of one top-level frontmatter key (indented
// continuation lines and comments included), or "" when the key is absent.
func section(fm, key string) string {
	var b strings.Builder
	in := false
	for _, line := range strings.Split(fm, "\n") {
		switch {
		case strings.HasPrefix(line, key+":"):
			in = true
		case in && (line == "" || (line[0] != ' ' && line[0] != '#')):
			return b.String()
		}
		if in {
			b.WriteString(line + "\n")
		}
	}
	return b.String()
}

var reIf = regexp.MustCompile(`(?m)^\s+if: "([^"]+)"`)
var reEcho = regexp.MustCompile(`(?m)^\s+echo '(\{.*\})'\s*$`)

func TestProfileGrantIsDeclarativeHooks(t *testing.T) {
	fm := skillFrontmatter(t)
	hooks := section(fm, "hooks")
	if hooks == "" {
		t.Fatal("SKILL.md frontmatter has no hooks section")
	}
	if !strings.Contains(hooks, "PreToolUse:") {
		t.Error("hooks section is not a PreToolUse registration")
	}

	var rules []string
	for _, m := range reIf.FindAllStringSubmatch(hooks, -1) {
		rules = append(rules, m[1])
	}
	// Inside an `if` the tool name is literal (unlike permission rules, where
	// Edit covers Write), so a first-meeting profile creation needs the Write rule.
	want := []string{"Read(~/.ai-mentor/**)", "Edit(~/.ai-mentor/profile.md)", "Write(~/.ai-mentor/profile.md)"}
	if strings.Join(rules, "|") != strings.Join(want, "|") {
		t.Errorf("hook if-rules = %q, want exactly %q", rules, want)
	}

	// Every hook command is an echo of a constant allow decision: a plugin
	// that promised "no code" must not grow a script here.
	commands := regexp.MustCompile(`(?m)^\s+command: `).FindAllStringIndex(hooks, -1)
	echos := reEcho.FindAllStringSubmatch(hooks, -1)
	if len(commands) != 3 || len(echos) != 3 {
		t.Fatalf("want 3 hook commands, each a one-line echo; got %d commands, %d echoes", len(commands), len(echos))
	}
	for _, m := range echos {
		var out struct {
			H struct {
				Event    string `json:"hookEventName"`
				Decision string `json:"permissionDecision"`
				Reason   string `json:"permissionDecisionReason"`
			} `json:"hookSpecificOutput"`
		}
		if err := json.Unmarshal([]byte(m[1]), &out); err != nil {
			t.Errorf("hook echo is not valid JSON: %v\n%s", err, m[1])
			continue
		}
		if out.H.Event != "PreToolUse" || out.H.Decision != "allow" || !strings.HasPrefix(out.H.Reason, "ai-mentor: ") {
			t.Errorf("hook echo must be a PreToolUse allow with an ai-mentor reason, got %+v", out.H)
		}
	}
	if strings.Contains(hooks, "deny") {
		t.Error("the mentor's hooks never deny; a user's own rules do that")
	}

	// The profile grant has one source of truth.
	if allowed := section(fm, "allowed-tools"); strings.Contains(allowed, ".ai-mentor") {
		t.Errorf("allowed-tools must not carry a ~/.ai-mentor rule (the hooks own the profile grant):\n%s", allowed)
	}
}
