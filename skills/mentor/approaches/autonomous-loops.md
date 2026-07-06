# Autonomous Loops
*Last verified: 2026-07-06*

## What It Is

Autonomous Loops let you give your AI coding tool a measurable goal — like "all tests pass" or "test coverage above 80%" — and then step away while it works toward that goal on its own. The AI writes code, runs tests, reads errors, fixes problems, and repeats until the condition you set is satisfied. You define the finish line; the AI figures out how to cross it.

## Why It Works

Feedback loops are how all iterative processes converge on a solution. Autonomous Loops automate the change-test-observe cycle by letting the AI act as both developer and evaluator: after each turn, a fast model checks whether the goal condition is met, and if not, the AI keeps working. This removes the bottleneck of a human reading each intermediate result and typing "try again" — low-value labor when the success criterion is already well-defined. The key insight is that many development tasks have a clear, machine-verifiable definition of "done" — and for those tasks, human supervision of each iteration adds latency without adding judgment.

## When to Use It

- Fixing a failing test suite where the fix is mechanical — the AI can iterate until green without your guidance
- Generating tests until a specific coverage threshold is met (e.g., "line coverage above 80% for `src/auth/`")
- Refactoring code to satisfy a linter or type checker — the goal is "zero errors from `eslint` or `mypy`"
- Debugging a build failure in a dependency upgrade where the errors are cascading and require multiple incremental fixes
- Resolving type errors after a large-scale rename or interface change — the goal is "zero TypeScript errors" and each fix may reveal the next one

## When NOT to Use It

- When the goal is subjective or aesthetic — "make the code cleaner" has no verifiable finish line, and the AI will either loop forever or declare victory prematurely
- When you do not trust the AI to make unsupervised architectural decisions — autonomous loops optimize for the stated condition, not for code quality you have not specified
- When the problem requires domain knowledge or context that cannot be expressed in a pass/fail condition — for example, "migrate to the new pricing model" has business rules the AI does not know

## How It Works

### Basic (Beginner)

1. Let turns run unattended (optional — `/goal` works in any permission mode, but you would otherwise approve each tool call): press Shift+Tab to cycle to "auto" mode if your account supports it, or start Claude with `--permission-mode auto`
2. Start a goal: `/goal all tests in tests/payment/ pass`
3. Claude reads the first failing test output, edits code, runs `pytest tests/payment/`, reads the new failures, and edits again
4. After each turn, a fast model evaluates the test output against your condition
5. When all tests pass, the goal clears automatically and Claude stops — check progress anytime with `/goal`, or stop early with `/goal clear`

### Composing with Other Approaches (Intermediate)

- **Plan Mode then Autonomous Loop**: Use Plan Mode to design an approach for a complex migration, review and approve the plan, then switch to an autonomous loop to execute it: `/goal no TypeScript errors in src/`. The plan ensures strategic soundness; the loop handles mechanical iteration.
- **Worktree plus Autonomous Loop**: Enter a worktree first, then run the autonomous loop inside the isolated copy. If the AI produces changes that pass tests but look wrong, you can discard the entire worktree and try a different goal condition.
- **Autonomous Loop then Code Review**: Let the loop run to green, then run `/code-review` on the diff. The loop optimizes for your stated condition; the review catches everything the condition did not cover.

### Advanced Patterns

- **Headless goal loops**: Run `claude -p "/goal coverage above 85% for src/billing/"` in CI or a background terminal. The AI works unattended and exits when done. Combine with `--output-format json` to capture the final result programmatically.
- **Compound conditions**: Goal conditions can be up to 4,000 characters. Use this to set multi-part goals: `/goal all tests pass AND no eslint errors AND no TypeScript errors`. The evaluator checks all parts.
- **Know when you want `/loop` instead**: `/goal` starts the next turn as soon as the previous one finishes and stops when the evaluator confirms the condition; `/loop` re-runs a prompt on a time interval and stops when you stop it. Use `/loop` for time-triggered work like polling a deploy, and `/goal` for condition-driven convergence.

## Common Pitfalls

- **Vague conditions**: A goal like "improve performance" gives the evaluator nothing measurable to check. The AI will either stop after one change or loop indefinitely — a goal runs until the evaluator judges the condition met or you run `/goal clear`. Always use conditions that produce a clear pass/fail signal.
- **Superficial fixes**: The AI optimizes for the stated condition. If your goal is "all tests pass," it might delete a failing test, add `@pytest.mark.skip`, or catch-and-swallow an exception. Review the diff carefully — start with `git diff --stat` for the blast radius — passing tests are necessary but not sufficient for correct code.
- **Over-long loops**: If the AI has not converged after 15-20 turns, the goal may be underspecified or the problem may require a different approach entirely. There is no built-in turn limit, so bound the run in the condition itself, e.g. `/goal tests pass, or stop after 20 turns`.
- **Cost blindness**: Each iteration consumes tokens. A 20-turn loop that reads large files and runs long test suites can cost significantly more than a focused interactive session. Monitor iteration count and set the goal scope appropriately — target a specific directory or test file rather than the entire project.

## Real-World Example

**Scenario**: Upgrading a dependency with cascading test failures.

You are upgrading `axios` from v0.x to v1.x in a Node.js project. After updating `package.json` and running `npm install`, you run the test suite and 23 tests fail across `tests/api/` because the error response structure changed (`error.response.data` vs. `error.response`).

```
claude
> /goal all tests in tests/api/ pass
```

Claude reads the first failing test in `tests/api/orders.test.ts`, sees `TypeError: Cannot read properties of undefined (reading 'message')` at line 47, and traces it to the error handler in `src/api/client.ts`. It updates the error destructuring from `error.response.data.message` to `error.response?.data?.message` with a fallback. It runs the test suite again — 18 still fail.

It reads the next failure, finds that `src/api/interceptors.ts` references `error.request.headers` which moved to `error.config.headers` in axios v1. It fixes that — 9 failures remain. Over four more edit-test cycles, Claude updates the request config shape in `src/api/upload.ts` and fixes a changed default for `transformResponse` in `src/api/client.ts`. All 23 tests pass.

Total time: about three minutes, seven iterations. You run `git diff` and review 4 files changed, 31 insertions, 19 deletions — all mechanical changes to match the new axios API. You scan the diff to make sure nothing was skipped or hacked around, confirm the fixes are correct, and commit.

## Sources

- [Keep Claude Working Toward a Goal](https://code.claude.com/docs/en/goal) — Official docs for /goal: setting, checking, and clearing goals, and how the evaluator works
- [Choose a Permission Mode](https://code.claude.com/docs/en/permission-modes) — Official docs on auto mode and the other permission modes that let loop turns run unattended
