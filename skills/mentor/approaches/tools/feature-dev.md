---
kind: plugin
last_verified: 2026-07-12
composes_with:
  - plan-mode
  - custom-agents
install: /plugin install feature-dev@claude-plugins-official
facts: "7-phase guided feature development with explorer/architect/reviewer agents. Hands-on (start verified): the phased flow engages correctly and scales down sensibly on small repos."
session_signal: "feature-dev is installed (its skills/commands are visible in the session) or its commands run in this conversation"
source: https://github.com/anthropics/claude-plugins-official
pitfalls:
  - "Overlaps plan mode — its value is the packaged pipeline; if the team already plans rigorously, plan mode alone may be enough."
---
