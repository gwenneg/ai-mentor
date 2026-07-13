---
kind: plugin
last_verified: 2026-07-13
composes_with:
  - hooks-as-workflow
  - built-in-review-skills
install: /plugin install security-guidance@claude-plugins-official
facts: "Three review layers, each deeper than the last: a per-edit pattern match (no model call, ~0 always-on tokens), an end-of-turn background LLM review of the turn's git diff, and an agentic review on each commit/push Claude makes (capped at 20 per rolling hour). All layers are advisory — none blocks writes or commits; findings re-prompt Claude to fix in-session. Hands-on: an injection attempt produced hardened parameterized code; invisible when quiet. Complements the on-demand /security-review built-in."
session_signal: "security-guidance is installed (its skills/commands are visible in the session) or its commands run in this conversation"
source: https://code.claude.com/docs/en/security-guidance
pitfalls:
  - "It reviews as you edit; it is not an audit of existing code — pair with /security-review for the branch-level pass."
  - "No dependency/supply-chain scanning or policy enforcement — the docs defer those to your CI scanners as a separate stage."
  - "Commit reviews cover only commits Claude makes through its Bash tool — commits from your own shell (including the `!` escape) are not reviewed."
---
