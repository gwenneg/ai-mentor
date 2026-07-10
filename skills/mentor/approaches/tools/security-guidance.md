---
kind: plugin
last_verified: 2026-07-03
composes_with:
  - hooks-as-workflow
  - built-in-review-skills
install: /plugin install security-guidance@claude-plugins-official
facts: "Per-edit security hooks plus a Stop-time LLM diff review — 12 hooks, ~0 always-on tokens. Hands-on: an injection attempt produced hardened parameterized code; invisible when quiet. Complements the on-demand /security-review built-in."
session_signal: "security-guidance is installed (its skills/commands are visible in the session) or its commands run in this conversation"
source: https://github.com/anthropics/claude-plugins-official
pitfalls:
  - "It reviews as you edit; it is not an audit of existing code — pair with /security-review for the branch-level pass."
---
