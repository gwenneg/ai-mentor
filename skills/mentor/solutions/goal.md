# /goal
*Last verified: 2026-07-06*

kind: builtin-command
goals: debugging, testing, migration
best_when: the finish line is machine-verifiable (tests green, zero type errors) and human supervision of each iteration adds nothing
composes_with: autonomous-loops, safe-autonomy, worktree-isolation
exemplar: /goal all tests in tests/payment/ pass
session_signal: user ran /goal in this conversation
pitfalls:
- The evaluator optimizes for the stated condition — a "tests pass" goal can be satisfied by deleting the test. Review the diff (`git diff --stat` first).
- Vague conditions ("improve performance") give the evaluator nothing checkable; check progress with `/goal`, stop with `/goal clear`.
source: https://code.claude.com/docs/en/goal
