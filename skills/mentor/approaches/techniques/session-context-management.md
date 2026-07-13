# Session & Context Management
*Last verified: 2026-07-13*

## What It Is

Session & Context Management is the discipline of keeping your AI session's working memory healthy: seeing where the context window is going, compacting it when it fills, keeping side questions from polluting the main thread, and branching or rewinding the conversation instead of fighting a degraded one. Claude Code ships a full toolkit for this — `/context`, `/compact`, `/btw`, `/clear`, `/branch`, `/rewind`, `/resume` — that most developers never touch.

## Why It Works

The model pays attention to everything in context, including the noise — curating what stays keeps it operating at full quality for hours.

## When to Use It

- A long session starts giving vaguer, slower, or more forgetful answers — the classic sign of a saturated context window
- You want to ask a quick side question ("what does this flag do?") without burying the main task under a tangent
- You finished one task and are starting an unrelated one in the same terminal
- You want to try a different direction without losing the conversation as it stands
- You need to pick up yesterday's session exactly where it left off

## When NOT to Use It

- Short, single-task sessions — the defaults handle context fine; management overhead buys nothing
- As a substitute for persistent knowledge — facts worth keeping across sessions belong in CLAUDE.md or auto memory, not in a carefully preserved conversation

## How It Works

### Basic (Beginner)

1. Run `/context` to see where the window is going — it renders usage as a colored grid and flags context-heavy tools, memory bloat, and capacity issues.
2. When the conversation gets long, run `/compact` to summarize it down. Pass focus instructions to control what survives: `/compact keep the migration plan and open bugs`.
3. Ask side questions with `/btw <question>` — you get an answer without the exchange being added to the conversation history.
4. Feed real command output into the conversation with shell mode: `! npm test` runs the command directly, and Claude responds once the output lands in the transcript (v2.1.186+; the response costs the same as a normal prompt) — no copy-pasting failures into the prompt. Set `respondToBashCommands` to `false` in settings.json to add the output to context without a response.
5. Starting an unrelated task? `/clear` gives you an empty context; the previous conversation stays available in `/resume` — and, on v2.1.191+, as a `previous session` entry at the top of the `/rewind` menu until you exit Claude Code or resume a different session. Pass a name (`/clear payments-work`) to label the old one for easy retrieval.
6. Use `/rewind` to step the conversation and/or code back to an earlier point when a direction didn't pan out.

### Composing with Other Approaches (Intermediate)

- **Context management plus project memory**: `/compact` and `/clear` are safe when durable facts live in CLAUDE.md — the project-root CLAUDE.md is re-injected after compaction, so instructions survive even when conversation details don't.
- **Branching plus checkpoints**: `/branch` forks the conversation at the current point to try a different direction while keeping the original intact — the conversation-level twin of checkpoint-based code rewind.
- **Context management plus subagents**: delegate bulk reading to subagents so their file dumps never enter your main context — the orchestrator receives summaries, and `/context` stays green through a whole afternoon.

### Advanced Patterns

- **Named session workflows**: `/rename` sessions as you go and resume any of them by name with `/resume <name>` — turning sessions into durable, addressable workstreams rather than one anonymous scrollback. Sessions are also linked to the pull requests they create via `gh pr create`: paste a PR URL into the `/resume` picker and the list filters to the session that created it (v2.1.122+; GitHub, GitHub Enterprise, GitLab, and Bitbucket URLs all work), or skip the picker entirely with `claude --from-pr 1234`.
- **Mid-session directory moves**: `/cd ../other-project` (v2.1.169+) relocates the session to a different working directory without rebuilding the prompt cache — the new directory's CLAUDE.md is appended as a message instead of replacing the system prompt. The session moves to the new directory's project storage, so `--resume` and `--continue` find it there, and Claude prompts you to trust a directory you haven't worked in before.
- **Extended 1M-token context**: for sessions that legitimately need huge context — sprawling monorepos, massive migration diffs — append `[1m]` to the model: `/model opus[1m]` or `/model sonnet[1m]` (Sonnet 5 runs the 1M window natively on the Anthropic API; Opus availability varies by plan). A bigger window is not a substitute for curation — noise degrades quality at any size — but it raises the ceiling when the working set is genuinely large.
- **Periodic setup checkups**: `/doctor` (alias `/checkup`, v2.1.205+) audits your whole setup and offers to fix what it finds — flagging unused skills, MCP servers, and plugins against their context cost, deduplicating local CLAUDE.md files against checked-in ones, migrating always-loaded guidance into skills and nested CLAUDE.md files that load on demand, and calling out slow hooks. It reports findings first and asks before changing anything — run it when `/context` shows overhead you can't account for.
- **Bisecting a broken setup**: when the problem is misbehavior rather than overhead, `claude --safe-mode` (v2.1.169+, or the `CLAUDE_CODE_SAFE_MODE` env var) launches with every customization disabled — CLAUDE.md, skills, plugins, hooks, MCP servers, and custom commands and agents don't load, while authentication, model selection, built-in tools, and permissions still work. If the problem disappears in safe mode, one of those surfaces is the cause.
- **Cache-aware session habits**: the prompt cache is keyed on your model, effort level, and the conversation prefix — so `/model` and `/effort` switches make the *next* turn reprocess the whole history uncached (one slow, expensive turn), while `/compact` deliberately rebuilds the conversation layer. Pick model and effort at the top of a session, save `/compact` for natural breaks between tasks, and when abandoning a bad path prefer `/rewind` over compaction: rewinding truncates back to a prefix that's already cached.

## Common Pitfalls

- **Compacting too late**: by the time responses visibly degrade, the noise has already been influencing answers for a while. Check `/context` at natural task boundaries instead of waiting for symptoms — compacting at a break you choose beats auto-compaction firing mid-task.
- **Using the main thread for everything**: every "quick question" answered in the main conversation permanently occupies context. `/btw` exists precisely for this — use it.
- **Preserving conversations that should die**: nursing a confused 3-hour session is usually worse than `/clear` plus a two-sentence restatement of where you are. The restatement forces clarity the old context lacked.
- **Confusing conversation memory with project memory**: anything you'd be sad to lose in a `/clear` belongs in CLAUDE.md or auto memory. If losing the conversation would lose knowledge, the knowledge is in the wrong place.

## Sources

- [Claude Code commands reference](https://code.claude.com/docs/en/commands) — Official reference for /context, /compact, /btw, /clear, /branch, /rewind, and /resume
- [How Claude remembers your project](https://code.claude.com/docs/en/memory) — What survives compaction and how CLAUDE.md re-injection works
- [Model configuration](https://code.claude.com/docs/en/model-config) — Official docs for the [1m] extended context window and per-plan availability
- [How Claude Code uses prompt caching](https://code.claude.com/docs/en/prompt-caching) — What invalidates the cache and what /compact and /rewind cost
- [Interactive mode](https://code.claude.com/docs/en/interactive-mode) — Shell mode with the ! prefix and how Claude responds to command output
- [Manage sessions](https://code.claude.com/docs/en/sessions) — The session picker, including resuming by PR URL
- [Debug your configuration](https://code.claude.com/docs/en/debug-your-config) — Testing against a clean configuration with --safe-mode

## Signals

- Setup: —
- Session: `/context`, `/compact`, `/btw`, `/clear`, `/cd`, `/resume`, `/branch`, or `/doctor` in the transcript
