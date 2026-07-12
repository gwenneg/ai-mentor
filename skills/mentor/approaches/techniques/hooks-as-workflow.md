# Hooks as Workflow
*Last verified: 2026-07-12*

## What It Is

Hooks are automatic actions that fire when Claude Code does something — edits a file, runs a command, finishes a session. Instead of remembering to run your linter or tests manually, you wire them to happen automatically at the right moment. Think of hooks as guardrails that enforce your workflow without requiring discipline.

## Why It Works

Hooks shift feedback left to the earliest possible moment — the instant Claude makes a change — automating the checks developers intend to do but skip under pressure.

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

1. Identify the trigger event: `PreToolUse` (before an action), `PostToolUse` (after an action), `Stop` (when Claude finishes responding), `SessionStart`, or `Notification` — among others; see the hooks reference for the full event list
2. Write the hook action: a shell `command`, a `prompt` fed to the model, an `agent` (experimental), an `http` POST to a URL, or an `mcp_tool` call on a connected MCP server
3. Add it to `.claude/settings.json` under the `hooks` key
4. The hook fires automatically — no manual invocation needed

Don't want to write the JSON by hand? The official `hookify` plugin turns a plain-language rule — or a repeated pattern it spots in your conversation history — into a working hook for you.

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
- **Hooks plus subagent delegation**: Use `SubagentStop` hooks to log each subagent's result as it finishes, and a Stop hook to post a summary once the main agent completes the turn. Since v2.1.195, a hyphenated matcher like `code-reviewer` exact-matches instead of substring-matching — on earlier versions anchor it as `^code-reviewer$` so it doesn't also fire for `senior-code-reviewer`.

### Advanced Patterns

- **Layered protection**: Stack PreToolUse hooks — one blocks `.env` edits, another requires confirmation for database migration files, a third rejects any `rm -rf` command. Each layer catches what the others miss.
- **Async notification hooks**: Fire Slack webhooks or desktop notifications when Claude finishes a task, without blocking the session. Useful for long-running autonomous work.
- **Goal evaluation via Stop hooks**: When Claude finishes responding, a Stop hook can run a validation script — and exit with code 2 to block the stop, feeding the failure back so Claude keeps working until acceptance criteria are met. This turns informal "is it done?" into automated verification.

## Common Pitfalls

- **Slow hooks kill productivity**: A hook that runs `npm test` on every edit adds 10+ seconds of latency per change. Scope tests to the changed file (e.g. `jest --findRelatedTests`). If the full suite is needed, run it asynchronously.
- **Over-blocking with PreToolUse**: Too many "are you sure?" prompts train developers to click yes without reading. Reserve blocking hooks for genuinely dangerous operations.
- **Forgetting hook scope**: Hooks in `.claude/settings.json` apply to the project. Hooks in `~/.claude/settings.json` apply globally. A formatting hook for a JavaScript project should not fire in a Go project.

## Sources

- [Claude Code Hooks](https://code.claude.com/docs/en/hooks) — Hooks reference covering all hook events, matchers, and exit codes
- [Claude Code Hooks Guide](https://code.claude.com/docs/en/hooks-guide) — Practical guide for automating workflows with hooks

## Signals

- Setup: `hooks` configured in `.claude/settings.json`, `.claude/settings.local.json`, or `~/.claude/settings.json`
- Session: Asks about automating actions on edits/commits
