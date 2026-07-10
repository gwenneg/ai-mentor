# /run
*Last verified: 2026-07-06*

kind: builtin-command
goals: testing, debugging
best_when: you want to see the change working in the real app, not infer it from green tests
composes_with: verify, browser-integration
exemplar: /run
session_signal: user ran /run in this conversation
pitfalls:
- It launches and drives the project's app; for a library with no runnable surface it has nothing to drive.
source: https://code.claude.com/docs/en/commands
