# /simplify
*Last verified: 2026-07-06*

kind: builtin-command
goals: refactoring, tech-debt
best_when: working code needs reuse, simplification, and efficiency cleanup — quality only, not bug hunting
composes_with: code-review, built-in-review-skills
exemplar: /simplify
session_signal: user ran /simplify in this conversation
pitfalls:
- Run /code-review first: simplification can consolidate duplicated code that still contains a logic error, leaving the bug in one tidier place.
source: https://code.claude.com/docs/en/commands
