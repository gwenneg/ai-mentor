---
kind: plugin
last_verified: 2026-07-03
composes_with: [visual-artifacts]
install: /plugin install project-artifact@claude-plugins-official
facts: "Publishes a living project status page with workstreams and decisions. Hands-on: produced a project-specific tabbed status page with honest unverified-state markings."
session_signal: "project-artifact is installed (its skills/commands are visible in the session) or its commands run in this conversation"
source: https://github.com/anthropics/claude-plugins-official
pitfalls:
  - "Publishing needs an interactive claude.ai session — headless falls back to a local HTML file plus refresh config."
---
