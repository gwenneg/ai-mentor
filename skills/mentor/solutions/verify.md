# /verify
*Last verified: 2026-07-06*

kind: builtin-command
goals: testing, code-review
best_when: tests and typecheck pass but nobody has watched the change actually work end-to-end
composes_with: code-review, run, built-in-review-skills
exemplar: /verify
session_signal: user ran /verify in this conversation
pitfalls:
- Don't run it on diffs with no runtime surface (docs-only, test-only) — there is no behavior to observe.
source: https://code.claude.com/docs/en/commands
