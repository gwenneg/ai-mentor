---
kind: plugin
last_verified: 2026-07-03
composes_with: [project-memory, hooks-as-workflow]
install: /plugin install claude-code-setup@claude-plugins-official
facts: "Analyzes a codebase and recommends tailored Claude Code automations. Hands-on: recommendations were concretely repo-tailored — each hook justified from real project facts, unjustified MCP servers declined."
session_signal: "claude-code-setup is installed (its skills/commands are visible in the session) or its commands run in this conversation"
source: https://github.com/anthropics/claude-plugins-official
pitfalls:
  - "Conceptually overlaps this plugin's own growth mode — prefer /mentor when the question is about the person, this plugin when it is about the repo."
---
