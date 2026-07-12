---
kind: plugin
last_verified: 2026-07-12
composes_with:
  - headless-mode
install: /plugin install commit-commands@claude-plugins-official
facts: "/commit, /commit-push-pr, and /clean_gone git workflow commands. Hands-on: flawless first try."
session_signal: "commit-commands is installed (its skills/commands are visible in the session) or its commands run in this conversation"
source: https://github.com/anthropics/claude-plugins-official
pitfalls:
  - "Mostly duplicates native committing — the real value is team commit conventions and /clean_gone; skip it if neither matters to you."
---
