# Claude Mods
*Last verified: 2026-10-12*

## What It Is

A Claude Mod is a plugin made of JavaScript or TypeScript event handlers (called hooks) that run inside the Claude Code process itself, instead of as an external script. Claude Code calls a mod's handler when something happens — a tool call, a submitted prompt, a part of the interface being drawn — and the handler can observe it, rewrite it, or answer it outright instead of letting the usual behavior run. Where a settings hook runs an outside script and a skill only gives Claude instructions, a mod runs *inside* Claude Code, so it can draw its own pane or band in the interface, redraw parts Claude Code draws itself (a tool call's row, the spinner, a dialog), and hold or intercept a tool call before it executes.

## Why It Works

Giving a plugin the same process-level access Claude Code's own UI and tool-execution loop have — rather than only a script's-eye view from outside — makes it possible to build the exact interface affordance a workflow needs (a live status pane, an inline confirmation) instead of approximating it with text in the transcript.

## When to Use It

- You want a persistent pane or a band above the prompt showing live state (context usage, a running cost tally, the current git branch) that a skill or settings hook can't draw
- You need to intercept a specific tool call and ask a custom question, show a custom confirmation, or substitute a different model for just that call
- You want a custom `/command` that runs instantly with no Claude turn involved, even while Claude is mid-task
- You're enforcing an organization policy on tool calls that a settings hook's allow/ask/deny rules can't express (e.g., a nuanced, stateful check across several calls)
- A settings hook, skill, or MCP server already covers the need — check the comparison in Sources before reaching for a mod

## When NOT to Use It

- The task is "give Claude more instructions" — that's a skill, with none of a mod's code-execution risk
- The task is "block or log this tool call based on a pattern" — a `PreToolUse` settings hook does this without writing JavaScript
- You don't fully trust the mod's source: mods are not sandboxed, run with your full permissions (file access, network, secrets in env vars and settings files), and can approve tool calls that an `ask` or even a `deny` rule would otherwise stop

## How It Works

### Basic (Beginner)

1. Ask for one in plain language in a Claude Code session (v2.1.287+): "make a mod that shows the current git branch above the prompt." Claude writes it from the built-in `plugin-authoring` skill into `~/.claude/dev-mods/<session-id>/<name>/` and asks whether to enable hot reloading for the session.
2. Approve each file write (mod directories under `~/.claude` are a protected path, so `default`/`acceptEdits` modes ask) and choose **Enable for this session** so the mod loads at the end of the current turn.
3. Run `/plugin`, go to the **Installed** tab, and confirm the mod is listed — or check the dim `N mods active · <names>` line under the tabs.
4. Use what you asked for, then iterate: tell Claude what to change, and the mod reloads at the end of the turn that edits its files.
5. To keep it beyond the session (which Claude Code deletes once older than `cleanupPeriodDays`), copy its directory elsewhere and load it with `claude --plugin-dir ~/mods/<name>` or add it to a marketplace.

A minimal hand-written mod — `hooks/register.js` — that counts tool calls and shows the count beside the spinner:
```javascript
let calls = 0

export function register(on) {
  on('tool.call', async ($, e, next) => {
    calls += 1
    $.ui.invalidate('ui.render')
    return next(e)
  })

  on('ui.render', { component: 'Spinner' }, async ($, e, next) => {
    return next({ ...e, props: { ...e.props, suffix: ' · tool calls: ' + calls + '…' } })
  })
}
```

### Composing with Other Approaches (Intermediate)

- **Mods plus Custom Plugins**: a mod is just one more thing a plugin can bundle alongside skills, hooks, agents, and an MCP server — the same `plugin.json` and marketplace distribution path applies; see Custom Plugins for packaging and versioning.
- **Mods plus Permissions & Safe Autonomy**: a mod can approve a tool call before the permission system asks, including one an `ask` rule would otherwise prompt for — so a policy mod can tighten what settings-file rules alone can express, but review one before installing with the same scrutiny as a deny-rule override.
- **Mods plus Hooks as Workflow**: both react to events, but a settings hook is a script you configure declaratively while a mod's handler is code that runs in-process and can draw UI — reach for a mod only once a settings hook's capabilities (block, allow, log) stop being enough.

### Advanced Patterns

- **Validate before installing**: `claude plugin validate <dir>` lists a mod's `hooks:` (which events, with their matchers) and `calls:` (which mods-API methods it uses) without running it — the way to audit what a mod does before trusting it with your permissions.
- **Turning mods off at different scopes**: disable one mod's plugin from `/plugin`; start a single session with `--safe-mode` or `--bare` to disable every mod for that run; or set `disableAllHooks: true` in `~/.claude/settings.json` to stop every mod (and your settings hooks and status line) in every session. Organizations gate this further with managed settings (e.g. `allowManagedModsOnly`).
- **Where mods do and don't draw**: hooks run in the terminal, the Desktop app's Code tab, the VS Code extension's chat panel, `claude -p`/the Agent SDK, Remote Control, and (for a plugin that reaches it) a cloud session — but drawing a pane or band only renders in the terminal and the Desktop app; everywhere else, a well-built mod falls back to transcript text or a command's reply.
- **Built-in mods as examples**: run `/plugin` → **Installed** → **Show disabled** to see Claude Code's own features implemented as mods (`cc-plugin-diff` for `/diff`, `cc-plugin-agents-md` for `AGENTS.md` loading, and the opt-in `cc-plugin-you-should-know`, v2.1.287 — a side agent that watches longer tasks and surfaces a note above the prompt on things you or Claude might miss). Several ship their source in the `claude-code` repo's `mods/` directory as working examples to read or fork.

## Common Pitfalls

- **Installing one without reading what it does**: a mod runs unsandboxed with your full permissions — it can read secrets, start processes, and approve tool calls silently. Run `claude plugin validate` and read the `hooks:`/`calls:` lines before trusting an unfamiliar mod, especially one from outside your own organization.
- **Expecting UI everywhere**: a mod's panes and bands only render in the terminal and the Desktop Code tab. If you build one for a team that works mostly through the VS Code extension or `claude -p`, its drawn interface won't show up for them — only its hooks will still run.
- **Confusing a mod with a settings hook**: both are called "hooks" in the docs, but a settings hook is a script/HTTP/prompt call configured in `settings.json`; a mod's hook is a function running inside Claude Code. If you just need to block or log a tool call with a script you already have, a settings hook is simpler and has none of a mod's code-execution surface.
- **Forgetting the version gate**: mods need Claude Code v2.1.287+ in the terminal and v2.1.286+ in the Desktop app. An older `CLAUDE_CODE_ENABLE_FUNCTION_HOOKS` env var from early access no longer does anything on current versions.

## Sources

- [Mods overview](https://code.claude.com/docs/en/plugins/mods/overview) — Official docs: what a mod can do, trust and installation, where mods run, and the mods-vs-settings-hooks-vs-skills-vs-MCP comparison
- [Create a mod](https://code.claude.com/docs/en/plugins/mods/create) — Official tutorial for asking Claude to write a mod or writing one by hand, plus validate/test loop
- [Claude Code changelog](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md) — v2.1.287: "Added Claude Mods: plugins may now modify deeper behavior" and the `cc-plugin-you-should-know` built-in mod

## Signals

- Setup: a plugin directory whose `hooks/hooks.json` has a `modules` key, or `~/.claude/dev-mods/` contains session-written mods
- Session: `/plugin` referencing an "Installed" mod; asks for a pane, a band above the prompt, or to intercept a tool call
