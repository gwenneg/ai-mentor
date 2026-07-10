# frontend-design
*Last verified: 2026-07-03*

kind: plugin
goals: greenfield
best_when: UI work should come out production-grade and visually bold by default
composes_with: visual-artifacts, browser-integration
install: /plugin install frontend-design@claude-plugins-official
facts: Auto-invoked skill for production-grade UI design. Hands-on: auto-engaged (invocation observed in the transcript) and produced a branded page in 4 turns.
session_signal: frontend-design is installed (its skills/commands are visible in the session) or its commands run in this conversation
pitfalls:
- Its "self-contained" output included a Google Fonts link — check outputs when offline-safe artifacts matter.
source: https://github.com/anthropics/claude-plugins-official
