---
kind: plugin
last_verified: 2026-07-12
composes_with:
  - hooks-as-workflow
install: /plugin install hookify@claude-plugins-official
facts: "Creates hooks from conversation patterns or explicit rules. Hands-on: generated a working PostToolUse hook and the firing was verified."
session_signal: "hookify is installed (its skills/commands are visible in the session) or its commands run in this conversation"
source: https://github.com/anthropics/claude-plugins-official
pitfalls:
  - "Headless caveat: it cannot write settings files non-interactively — hook creation needs an interactive session."
---
