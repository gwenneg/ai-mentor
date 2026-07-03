# Session & Context Management
*Last verified: 2026-07-02*

## What It Is

Session & Context Management is the discipline of keeping your AI session's working memory healthy: seeing where the context window is going, compacting it when it fills, keeping side questions from polluting the main thread, and branching or rewinding the conversation instead of fighting a degraded one. Claude Code ships a full toolkit for this — `/context`, `/compact`, `/btw`, `/clear`, `/branch`, `/rewind`, `/resume` — that most developers never touch.

## Why It Works

Context windows are finite, and response quality degrades as they fill with stale tool output, abandoned tangents, and half-relevant history. The model pays attention to everything in context — including the noise. Managing context is the same principle as managing working memory on a team: summarize what's settled, archive what's done, and start focused threads for new topics. Developers who treat the conversation as an infinite scrollback get mysteriously worse results over long sessions; developers who curate context keep the model operating at full quality for hours.

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
4. Starting an unrelated task? `/clear` gives you an empty context; the previous conversation stays available in `/resume`. Pass a name (`/clear payments-work`) to label the old one for easy retrieval.
5. Use `/rewind` to step the conversation and/or code back to an earlier point when a direction didn't pan out.

### Composing with Other Approaches (Intermediate)

- **Context management plus project memory**: `/compact` and `/clear` are safe when durable facts live in CLAUDE.md — the project-root CLAUDE.md is re-injected after compaction, so instructions survive even when conversation details don't.
- **Branching plus checkpoints**: `/branch` forks the conversation at the current point to try a different direction while keeping the original intact — the conversation-level twin of checkpoint-based code rewind.
- **Context management plus subagents**: delegate bulk reading to subagents so their file dumps never enter your main context — the orchestrator receives summaries, and `/context` stays green through a whole afternoon.

### Advanced Patterns

- **Named session workflows**: `/rename` sessions as you go and resume any of them by name with `/resume <name>` — turning sessions into durable, addressable workstreams rather than one anonymous scrollback.
- **Recovering from a premature `/clear`**: `/rewind` can restore the conversation from before a `/clear` (v2.1.191+), so clearing is no longer irreversible.
- **Compaction-aware task ordering**: settle decisions early ("we're using approach B") and let details accumulate after — compaction preserves conclusions better than meandering deliberation, so a session that decides-then-executes compacts cleanly while a session that deliberates forever compacts into mush.
- **Extended 1M-token context**: for sessions that legitimately need huge context — sprawling monorepos, massive migration diffs — append `[1m]` to the model: `/model opus[1m]` or `/model sonnet[1m]` (Sonnet 5 runs the 1M window natively on the Anthropic API; Opus availability varies by plan). A bigger window is not a substitute for curation — noise degrades quality at any size — but it raises the ceiling when the working set is genuinely large.

## Common Pitfalls

- **Compacting too late**: by the time responses visibly degrade, the noise has already been influencing answers for a while. Check `/context` at natural task boundaries instead of waiting for symptoms.
- **Using the main thread for everything**: every "quick question" answered in the main conversation permanently occupies context. `/btw` exists precisely for this — use it.
- **Preserving conversations that should die**: nursing a confused 3-hour session is usually worse than `/clear` plus a two-sentence restatement of where you are. The restatement forces clarity the old context lacked.
- **Confusing conversation memory with project memory**: anything you'd be sad to lose in a `/clear` belongs in CLAUDE.md or auto memory. If losing the conversation would lose knowledge, the knowledge is in the wrong place.

## Real-World Example

You're four hours into a gnarly migration in a monorepo. Claude starts re-asking about the package layout it knew earlier — the classic saturation sign. `/context` confirms it: the window is dominated by old test output from the morning's debugging.

You run `/compact keep the migration checklist, the three remaining failing packages, and the decision to use the compat shim`. The session shrinks to a tight summary plus your project's CLAUDE.md, which re-injects automatically. Claude's next answer is sharp again — it names the three failing packages without being reminded.

Mid-afternoon, a teammate pings you about an unrelated production question. Instead of contaminating the migration session, you ask it via `/btw` — answered, gone, zero context cost. At the end of the day you `/rename` the session `pkg-migration` and close the laptop. Tomorrow, `/resume pkg-migration` puts you exactly where you left off — checklist, decisions, and all.

## Sources

- [Claude Code commands reference](https://code.claude.com/docs/en/commands) — Official reference for /context, /compact, /btw, /clear, /branch, /rewind, and /resume
- [How Claude remembers your project](https://code.claude.com/docs/en/memory) — What survives compaction and how CLAUDE.md re-injection works
- [Model configuration](https://code.claude.com/docs/en/model-config) — Official docs for the [1m] extended context window and per-plan availability
