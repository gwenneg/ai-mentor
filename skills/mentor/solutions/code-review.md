# /code-review
*Last verified: 2026-07-06*

kind: builtin-command
goals: code-review, tech-debt
best_when: a diff or PR needs a structured correctness pass before humans read it
composes_with: verify, built-in-review-skills, subagent-delegation
exemplar: /code-review high src/api/
session_signal: user ran /code-review or /review in this conversation
pitfalls:
- Effort levels trade coverage for confidence: low/medium give fewer, higher-confidence findings; high through max widen coverage and may include uncertain ones. Pick by stakes, not habit.
- `--fix` applies findings to the working tree and `--comment` posts them as inline PR comments — say which one the situation wants.
source: https://code.claude.com/docs/en/commands
