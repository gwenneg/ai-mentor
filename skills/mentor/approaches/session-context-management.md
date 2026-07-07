# Session & Context Management
*Last verified: 2026-07-06*

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
4. Starting an unrelated task? `/clear` gives you an empty context; the previous conversation stays available in `/resume` — and, on v2.1.191+, as a `previous session` entry at the top of the `/rewind` menu until you exit Claude Code. Pass a name (`/clear payments-work`) to label the old one for easy retrieval.
5. Use `/rewind` to step the conversation and/or code back to an earlier point when a direction didn't pan out.

### Composing with Other Approaches (Intermediate)

- **Context management plus project memory**: `/compact` and `/clear` are safe when durable facts live in CLAUDE.md — the project-root CLAUDE.md is re-injected after compaction, so instructions survive even when conversation details don't.
- **Branching plus checkpoints**: `/branch` forks the conversation at the current point to try a different direction while keeping the original intact — the conversation-level twin of checkpoint-based code rewind.
- **Context management plus subagents**: delegate bulk reading to subagents so their file dumps never enter your main context — the orchestrator receives summaries, and `/context` stays green through a whole afternoon.

### Advanced Patterns

- **Named session workflows**: `/rename` sessions as you go and resume any of them by name with `/resume <name>` — turning sessions into durable, addressable workstreams rather than one anonymous scrollback.
- **Extended 1M-token context**: for sessions that legitimately need huge context — sprawling monorepos, massive migration diffs — append `[1m]` to the model: `/model opus[1m]` or `/model sonnet[1m]` (Sonnet 5 runs the 1M window natively on the Anthropic API; Opus availability varies by plan). A bigger window is not a substitute for curation — noise degrades quality at any size — but it raises the ceiling when the working set is genuinely large.
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
