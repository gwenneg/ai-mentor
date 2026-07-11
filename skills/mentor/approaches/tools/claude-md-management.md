---
kind: plugin
last_verified: 2026-07-03
composes_with:
  - project-memory
install: /plugin install claude-md-management@claude-plugins-official
facts: "Audits and maintains CLAUDE.md files. Hands-on: scored audit with a rubric and real gaps found, cross-checked against the codebase."
session_signal: "claude-md-management is installed (its skills/commands are visible in the session) or its commands run in this conversation"
source: https://github.com/anthropics/claude-plugins-official
pitfalls:
  - "The skill is invoked as `claude-md-improver`, not the plugin name."
---
