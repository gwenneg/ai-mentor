# Built-in Review Skills
*Last reviewed: 2026-07-01*

## What It Is

Built-in Review Skills are ready-to-use commands in Claude Code that analyze your code changes for bugs, security issues, and simplification opportunities. You run a slash command — `/code-review`, `/security-review`, or `/simplify` — and Claude systematically reviews your diff using structured analysis. No prompt engineering needed: the skills encode expert review strategies so you get consistent, thorough reviews every time.

## Why It Works

Code review quality depends heavily on structure. A human reviewer who "just reads through the diff" catches fewer bugs than one who follows a checklist: error handling, edge cases, concurrency, input validation. Built-in review skills apply this principle mechanically — each skill runs a defined review methodology against your changes, ensuring that common bug categories are never skipped. By separating review concerns (correctness vs. security vs. simplification), each pass can go deep on its category rather than doing a shallow scan of everything.

## When to Use It

- Self-reviewing your own changes before pushing — catch bugs your eyes skipped after hours of writing the code
- Reviewing a teammate's PR when you want a structured first pass before doing your own read-through
- Post-refactoring cleanup — run `/simplify` after a large restructuring to catch missed deduplication
- CI integration — run reviews automatically on every PR via GitHub Actions

## When NOT to Use It

- Reviewing changes you have not made yet — these skills analyze the current diff, not hypothetical code
- As a substitute for human review on critical paths — use them as a first pass, not the only pass
- On massive diffs (500+ files) without narrowing scope — the review quality degrades at extreme scale

## How It Works

### Basic (Beginner)

1. Make your code changes and stage them (or leave them unstaged — both work)
2. Run `/code-review` in your Claude Code session. Choose an effort level:
   - Low/Medium: fewer findings, higher confidence — good for quick sanity checks
   - High: broader coverage, may surface uncertain findings — good for thorough review
   - xhigh/max: deepest analysis available on Opus 4.7+ models — exhaustive coverage for critical changes
3. Claude analyzes the diff and reports findings grouped by severity
4. Optionally, add `--fix` to have Claude auto-apply its findings: `/code-review --fix`
5. Or add `--comment` to post findings as inline PR comments: `/code-review --comment`

### Composing with Other Approaches (Intermediate)

- **Review then simplify**: Run `/code-review` first to catch bugs, then `/simplify` to clean up the code. This ensures correctness before optimization — fixing a bug in already-simplified code is harder.
- **Subagent changes then review**: After spawning subagents to make parallel changes, run `/code-review` on each worktree branch to verify the agents' work before merging.
- **Plan Mode then review**: Use Plan Mode to design and execute a change, then immediately run `/security-review` to catch security implications the plan did not consider.

### Advanced Patterns

- **CI pipeline integration**: Add Claude Code to your GitHub Actions workflow in headless mode. On every PR, it runs `/code-review --comment` and posts inline findings directly on the PR. Reviewers see AI findings alongside the diff without any manual step.
- **Layered review strategy**: Run all three skills in sequence for critical changes: `/code-review` for correctness, `/security-review` for security (uses Opus-tier reasoning for deeper analysis), then `/simplify` to clean up. Each pass focuses on its specialty.
- **Targeted review with context**: Before running the review, tell Claude about specific concerns: "This change modifies our rate limiter. Run /security-review with extra attention to bypass vectors." The skill uses your context to focus its analysis.

## Common Pitfalls

- **Running only `/simplify` without `/code-review` first**: Simplification can mask bugs. If you have a logic error in duplicated code, `/simplify` might deduplicate it — now the bug is in one place but still exists. Review for correctness first.
- **Treating AI review as authoritative**: These skills catch real bugs, but they also produce false positives. A finding that says "possible null pointer" might be guarded by an upstream check the AI did not trace. Always evaluate findings with your domain knowledge.
- **Ignoring effort levels**: Running high-effort review on every trivial change wastes time. Use low/medium for routine work, save high effort for complex or risky changes.
- **Not using `--comment` in team workflows**: If you review locally but do not post comments, your teammates do not benefit. Use `--comment` to make findings visible on the PR.

## Real-World Example

You have just finished adding request validation to your API's `/users` endpoint. The change touches `api/handlers/users.go`, `api/middleware/validation.go`, and `api/handlers/users_test.go`.

```
> /code-review
```

Claude analyzes the diff and reports:

> **High**: In `api/handlers/users.go:47`, the email regex validation accepts
> strings up to 10,000 characters before checking format. A malicious request
> with a 10KB email field will cause catastrophic regex backtracking. Move the
> length check before the regex match.
>
> **Medium**: In `api/handlers/users_test.go:89`, `TestCreateUser_InvalidEmail`
> asserts status 400 but does not check the error message body. If the handler
> returns 400 for a different reason, this test passes silently.

You fix the regex ordering, then run `/simplify`:

> The validation logic in `api/middleware/validation.go:23-45` duplicates the
> email format check already present in `api/handlers/users.go:44`. Extract to
> `internal/validate/email.go` and call from both locations.

Claude applies the extraction automatically. You run `/security-review` as a final pass — it confirms the length-before-regex fix resolves the ReDoS vector and finds no additional issues. Three focused review passes, each catching something the others would not have prioritized.

## Sources

- [Claude Code Slash Commands](https://docs.anthropic.com/en/docs/claude-code/slash-commands) — Official docs for skills including built-in /code-review and /security-review
