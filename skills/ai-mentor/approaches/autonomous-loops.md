# Autonomous Loops
*Last reviewed: 2026-06-27*

## What It Is

Autonomous Loops let you give your AI coding tool a measurable goal — like "all tests pass" or "test coverage above 80%" — and then step away while it works toward that goal on its own. The AI writes code, runs tests, reads errors, fixes problems, and repeats until the condition you set is satisfied. You define the finish line; the AI figures out how to cross it.

## Why It Works

Feedback loops are how all iterative processes converge on a solution. A developer debugging manually follows the same pattern: change code, run tests, read output, change code again. Autonomous Loops automate this cycle by letting the AI act as both the developer and the evaluator. After each turn, a fast model checks whether the goal condition is met. If not, the AI continues. This removes the bottleneck of a human reading each intermediate result and typing "try again" — which is low-value labor when the success criterion is already well-defined. The key insight is that many development tasks have a clear, machine-verifiable definition of "done" — and for those tasks, human supervision of each iteration adds latency without adding judgment.

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

1. Enable auto mode: press Shift+Tab to cycle to "auto" mode, or start Claude with `--permission-mode auto`
2. Start a goal: `/goal all tests in tests/payment/ pass`
3. Claude reads the first failing test output to understand what is broken
4. It edits code, runs `pytest tests/payment/`, reads the new failures, and edits again
5. After each turn, a fast model evaluates the test output against your condition
6. When all tests pass, Claude stops and reports what it changed

### Composing with Other Approaches (Intermediate)

- **Plan Mode then Autonomous Loop**: Use Plan Mode to design an approach for a complex migration, review and approve the plan, then switch to an autonomous loop to execute it: `/goal no TypeScript errors in src/`. The plan ensures strategic soundness; the loop handles mechanical iteration.
- **Worktree plus Autonomous Loop**: Enter a worktree first, then run the autonomous loop inside the isolated copy. If the AI produces changes that pass tests but look wrong, you can discard the entire worktree and try a different goal condition.
- **Autonomous Loop then Code Review**: Let the loop run to green, then run `/code-review` on the diff. The loop optimizes for your stated condition; the review catches everything the condition did not cover.

### Advanced Patterns

- **Headless goal loops**: Run `claude -p "/goal coverage above 85% for src/billing/"` in CI or a background terminal. The AI works unattended and exits when done. Combine with `--output-format json` to capture the final result programmatically.
- **Compound conditions**: Goal conditions can be up to 4,000 characters. Use this to set multi-part goals: `/goal all tests pass AND no eslint errors AND no TypeScript errors`. The evaluator checks all parts.
- **Progressive tightening**: Start with a loose goal (`/goal tests pass`), review the result, then set a tighter goal (`/goal tests pass with no skipped tests and no console warnings`). Each round raises the bar.
- **Scoped loops for large projects**: In a monorepo, avoid `/goal all tests pass` across the entire project. Instead, target a specific package or directory: `/goal all tests in packages/billing/ pass`. Narrower scope means faster iterations and lower cost per loop.

## Tool Support

| Tool | Support | Notes |
|------|---------|-------|
| Claude Code | Native | `/goal <condition>`, evaluated by fast model after each turn, works with auto mode |
| OpenCode | None | No built-in goal loop; manual iteration required |
| Cursor | None | No autonomous loop feature; agent mode iterates but without goal-condition evaluation |
| aider | Partial | Can iterate on lint/test errors automatically, but no user-defined goal condition |

## Common Pitfalls

- **Vague conditions**: A goal like "improve performance" gives the evaluator nothing measurable to check. The AI will either stop after one change or loop until the turn limit. Always use conditions that produce a clear pass/fail signal.
- **Superficial fixes**: The AI optimizes for the stated condition. If your goal is "all tests pass," it might delete a failing test, add `@pytest.mark.skip`, or catch-and-swallow an exception. Review the diff carefully — passing tests are necessary but not sufficient for correct code.
- **Ignoring the diff**: After the loop completes, read what changed. Autonomous loops can touch many files across many iterations. A quick `git diff --stat` tells you the blast radius before you commit.
- **Over-long loops**: If the AI has not converged after 15-20 turns, the goal may be underspecified or the problem may require a different approach entirely. Set reasonable expectations for what can be achieved autonomously.
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

- [Claude Code Interactive Mode](https://docs.anthropic.com/en/docs/claude-code/interactive-mode) — Official docs covering /goal command and autonomous behavior
- [Enabling Claude Code to Work More Autonomously](https://www.anthropic.com/news/enabling-claude-code-to-work-more-autonomously) — Anthropic blog on goal-oriented autonomous execution
