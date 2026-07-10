---
kind: plugin
last_verified: 2026-07-03
composes_with: [plan-mode, worktree-isolation, autonomous-loops]
install: /plugin install code-modernization@claude-plugins-official
facts: "Structured migration of legacy codebases. Hands-on (start verified): the preflight phase engages coherently."
session_signal: "code-modernization is installed (its skills/commands are visible in the session) or its commands run in this conversation"
source: https://github.com/anthropics/claude-plugins-official
pitfalls:
  - "Needs a generous turn budget, and its multi-command probes fragment under default permissions — allowlist Bash for real runs."
  - "Biggest component surface in the catalog — expect meaningful always-on context."
---
