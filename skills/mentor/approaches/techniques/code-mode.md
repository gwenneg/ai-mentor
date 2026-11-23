# Code Mode
*Last verified: 2026-11-23*

## What It Is

Code Mode is an experimental, opt-in way of working in Claude Code where Claude writes a piece of JavaScript that makes several tool calls in one step, instead of asking for Bash, Read, or Edit one call at a time and reading each result before deciding what's next. Claude Code runs that code inside your session — it can run commands in parallel, feed one result into the next call, and return only the part of the output Claude actually needs to read.

## Why It Works

Writing code to orchestrate a batch of related tool calls collapses many call-and-read round trips into one, so filtering and chaining happen inside the code instead of burning a conversation turn per step.

## When to Use It

- A change that repeats across many files — one script edits every match instead of one Edit call per file
- A search whose results need filtering before Claude reads them — grep broadly, then keep only what matters in code
- Several independent commands that can run at the same time, such as linting, testing, and a build in parallel
- You've already hit the ten-minute foreground limit or context cost of many sequential ordinary tool calls on a repetitive task

## When NOT to Use It

- Routine or production-sensitive sessions where you want stable, well-exercised behavior — it's experimental and off by default, and Anthropic says it "may change or be removed in any release"
- Sessions that need the VS Code panel, the Desktop app, a cloud session, Remote Control, Claude in Chrome, or computer use — none of these work with Code Mode
- Sessions on Amazon Bedrock, Claude Platform on AWS, Google Cloud's Agent Platform, Microsoft Foundry, a gateway, or a custom `ANTHROPIC_BASE_URL` — Code Mode needs a direct Claude subscription or Console account, and is unavailable under HIPAA or zero-data-retention org configurations

## How It Works

### Basic (Beginner)

1. Check your version: `claude --version` must be v2.1.323 or later; run `claude update` if not.
2. Set `CLAUDE_CODE_EXPERIMENTAL_CODE_MODE=1` — in your shell to try it once, or in the `env` block of `~/.claude/settings.json` to turn it on for every new session. Claude Code ignores this variable in a project's own settings files, so a cloned repo can't silently turn it on for you.
3. Confirm **Yes, turn on Code Mode** at the one-time prompt. The session restarts and shows `Code Mode (experimental) is on · this session can be resumed only with Code Mode on` near the top.
4. Give Claude a task with many related tool calls, e.g. "Read every Markdown file in this project and list the links that point to local files that don't exist."
5. Press `Ctrl+O` and expand a collapsed **Code** row to read the JavaScript behind the tool calls — the ordinary conversation only shows the resulting tool calls (an edit's diff, a command's output), not the code itself.

### Composing with Other Approaches (Intermediate)

- **Code Mode plus Fan-Out Workflows**: for a batch of related operations that don't need separate context windows, one piece of code can replace several parallel subagents — reserve actual subagent fan-out for work that genuinely benefits from isolated context.
- **Code Mode plus Permissions & Safe Autonomy**: add an `ask` rule on `REPL` (the tool that runs the code) to see and approve each piece of code before it starts. In auto mode, the classifier reviews the code once before it runs — not each individual tool call — and can't see values computed at runtime, such as a path pulled from search output, so keep explicit `deny`/`ask` rules for calls that must always get individual scrutiny.
- **Code Mode plus Claude Mods**: a mod's `PostToolUse` hook matching `REPL` can inspect or redact what the code returns before Claude reads it — the same hook-on-a-tool pattern Claude Mods uses elsewhere, applied to a whole piece of code's output instead of one call's.

### Advanced Patterns

- **Long-running background code**: in auto or `bypassPermissions` mode, code running a long build or test can keep going in the background for up to two hours while you keep talking to Claude — check or stop it from `/tasks`. Outside those modes, Claude waits for the code to finish and it's capped at ten minutes.
- **OpenTelemetry correlation**: tool calls made from code emit the same `tool_decision`/`tool_result` events as ordinary calls, under the tool's usual name (e.g. `Bash`), with `tool_use_id` values prefixed `repl_` or `toolu_cc_` so you can tell them apart from direct calls in your backend.
- **Organization-wide block**: set `CLAUDE_CODE_EXPERIMENTAL_CODE_MODE` to `"0"` in managed settings to keep Code Mode off for everyone in your organization; once blocked, in-progress Code Mode sessions can't be resumed until the block is lifted.

## Common Pitfalls

- **Expecting it to work everywhere**: Code Mode only runs in an interactive terminal session, on a direct Anthropic subscription or Console sign-in, with Opus 5.5, Sonnet 5.5, Haiku 5.5, or Fable 5.2. It refuses to start on third-party providers, gateways, the VS Code extension, the Desktop app, and cloud sessions.
- **Losing visibility into what actually ran**: the main conversation shows only the resulting tool calls, not the JavaScript. If you need to audit what executed, press `Ctrl+O` and expand each **Code** row, or add an `ask` rule on `REPL` so you see the code before it runs.
- **Assuming the auto-mode classifier reviews every call like normal**: it reviews the code once, before execution, and has no visibility into values the code computes at runtime — a destructive call built from search output gets less scrutiny than the same call made directly. Keep explicit rules for calls you can't afford to miss.
- **Forgetting the resume flag**: a session started with Code Mode on can only be resumed with `CLAUDE_CODE_EXPERIMENTAL_CODE_MODE` still set to `1` — dropping it refuses the resume. Sessions that never had Code Mode on are unaffected.

## Sources

- [Run tool calls from code with Code Mode](https://code.claude.com/docs/en/code-mode) — Official docs: how Code Mode works, requirements, permission checks, long-running code, and troubleshooting
- [Claude Code changelog](https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md) — v2.1.323: "Added Code Mode, an experimental opt-in in which Claude writes JavaScript to call its tools"

## Signals

- Setup: `CLAUDE_CODE_EXPERIMENTAL_CODE_MODE` set in `~/.claude/settings.json` or the shell environment
- Session: `Code Mode (experimental) is on` banner near the top of the session; a collapsed **Code** row in the transcript viewer (`Ctrl+O`)
