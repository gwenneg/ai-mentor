# Hooks as Workflow
*Last verified: 2026-06-27*

## What It Is

Hooks are automatic actions that fire when Claude Code does something — edits a file, runs a command, finishes a session. Instead of remembering to run your linter or tests manually, you wire them to happen automatically at the right moment. Think of hooks as guardrails that enforce your workflow without requiring discipline.

## Why It Works

The biggest productivity losses in development come from delayed feedback. A bug caught at edit time costs seconds to fix. The same bug caught in CI costs minutes. Caught in production, hours. Hooks shift feedback left to the earliest possible moment — the instant Claude makes a change — by automating the checks that developers intend to do but often skip under pressure.

## When to Use It

- You want tests to run automatically after every edit in a specific directory
- You need to block edits to sensitive files (production configs, lock files, credentials)
- You want auto-formatting applied to every file Claude touches without asking
- You need a notification when Claude finishes a long-running task
- You want to enforce coding standards or security policies automatically

## When NOT to Use It

- One-off tasks where the hook setup takes longer than the task itself
- Hooks that run slow operations (full test suite, heavy compilation) — they break your flow
- When you need flexibility to bypass the rule frequently — hooks are for rules you always want enforced

## How It Works

### Basic (Beginner)

1. Identify the trigger event: `PreToolUse` (before an action), `PostToolUse` (after an action), `Stop` (session end), `SessionStart`, or `Notification`
2. Write the hook action: a shell command, an HTTP request, or a prompt
3. Add it to `.claude/settings.json` under the `hooks` key
4. The hook fires automatically — no manual invocation needed

Example — auto-format after every file edit:
```json
{
  "hooks": {
    "PostToolUse": [{
      "matcher": "Edit|Write",
      "hooks": [{
        "type": "command",
        "command": "prettier --write \"$(cat /dev/stdin | jq -r '.tool_input.file_path')\""
      }]
    }]
  }
}
```

### Composing with Other Approaches (Intermediate)

- **Hooks plus autonomous loops**: Add a PostToolUse hook that runs tests after every edit. Combined with `/goal "all tests pass"`, Claude gets instant test feedback on each iteration — tightening the autonomous loop.
- **Hooks plus worktree isolation**: Add a PreToolUse hook that blocks edits outside the worktree directory, preventing agents from accidentally modifying your main working tree.
- **Hooks plus subagent delegation**: Use a Stop hook to aggregate results from all subagents and post a summary when the last one finishes.

### Advanced Patterns

- **Layered protection**: Stack PreToolUse hooks — one blocks `.env` edits, another requires confirmation for database migration files, a third rejects any `rm -rf` command. Each layer catches what the others miss.
- **Async notification hooks**: Fire Slack webhooks or desktop notifications when Claude finishes a task, without blocking the session. Useful for long-running autonomous work.
- **Goal evaluation via Stop hooks**: When a session ends, a Stop hook can run a validation script and report whether the work met acceptance criteria — turning informal "is it done?" into automated verification.

## Common Pitfalls

- **Slow hooks kill productivity**: A hook that runs `npm test` on every edit adds 10+ seconds of latency per change. Use `--related` flags or scope tests to the changed file. If the full suite is needed, run it asynchronously.
- **Over-blocking with PreToolUse**: Too many "are you sure?" prompts train developers to click yes without reading. Reserve blocking hooks for genuinely dangerous operations.
- **Forgetting hook scope**: Hooks in `.claude/settings.json` apply to the project. Hooks in `~/.claude/settings.json` apply globally. A formatting hook for a JavaScript project should not fire in a Go project.

## Real-World Example

You are debugging a flaky integration test in `tests/integration/billing_test.py`. The test passes individually but fails when run with the full suite. You suspect shared database state.

```
> Add a PostToolUse hook that runs `pytest tests/integration/billing_test.py -x -v`
  after every edit to files in `src/billing/`. Also add a PreToolUse hook that
  blocks any edit to `tests/conftest.py` — I don't want to accidentally change
  the shared fixtures while debugging.
```

Claude adds both hooks to `.claude/settings.json`. Now every time you edit billing code, the failing test runs instantly. You discover the issue in `src/billing/invoice.py` — a class-level cache that persists between test runs. The PostToolUse hook confirms the fix immediately: the test passes. You remove the hooks and run the full suite to verify.

## Sources

- [Claude Code Hooks](https://code.claude.com/docs/en/hooks) — Hooks reference covering all hook events, matchers, and exit codes
- [Claude Code Hooks Guide](https://code.claude.com/docs/en/hooks-guide) — Practical guide for automating workflows with hooks
