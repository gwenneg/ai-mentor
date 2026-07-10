---
kind: plugin
last_verified: 2026-07-03
composes_with: [model-effort-selection]
install: /plugin install session-report@claude-plugins-official
facts: "Generates an HTML report of session token usage and cache efficiency. Hands-on: self-contained HTML with real usage numbers; cheapest always-on cost of the evaluated set (~70 tokens)."
session_signal: "session-report is installed (its skills/commands are visible in the session) or its commands run in this conversation"
source: https://github.com/anthropics/claude-plugins-official
pitfalls:
  - "Needs >12 turns of history to say anything useful, and default permissions block its bundled analyzer."
  - "Reports a 7-day window, not strictly the current session."
---
