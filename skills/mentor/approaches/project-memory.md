# Project Memory & Context Docs
*Last reviewed: 2026-07-02*

## What It Is

Project Memory gives your AI assistant persistent knowledge that survives across sessions. Two mechanisms work together: `CLAUDE.md` files hold instructions you write — build commands, conventions, architecture facts — and auto memory holds notes Claude writes for itself as it learns your project. Both load at the start of every conversation, so each session begins already knowing what previous sessions had to be told.

## Why It Works

Every session starts with a fresh context window; without persistent memory, you re-explain the same conventions forever and the AI re-makes the same mistakes. Writing knowledge down once converts a per-session cost into a one-time cost — the same reason teams write onboarding docs instead of re-explaining the codebase to every new hire. The quality of AI output is proportional to the quality of its context, and project memory is the highest-leverage context there is: a 50-line CLAUDE.md improves every response in every future session.

## When to Use It

- Claude makes the same mistake a second time — that's the signal to write the correction down
- You type the same clarification into chat that you typed last session
- A code review catches something Claude should have known about this codebase
- Starting work in a new repository — generate the initial file before doing anything else
- A new teammate would need the same context to be productive

## When NOT to Use It

- Multi-step procedures — those belong in a skill, which loads on demand, not in CLAUDE.md, which costs context every session
- Rules that must be enforced rather than followed — CLAUDE.md is context, not configuration; use a PreToolUse hook to actually block an action
- Instructions relevant to only one part of the codebase — use a path-scoped rule so it loads only when matching files are touched

## How It Works

### Basic (Beginner)

1. Run `/init` in your project. Claude analyzes the codebase and generates a starter `CLAUDE.md` with build commands, test instructions, and discovered conventions. If one already exists, it suggests improvements instead of overwriting.
2. Add the facts you'd otherwise re-explain, written to be verifiable: "Run `npm test` before committing", "API handlers live in `src/api/handlers/`" — not "test your changes" or "keep files organized".
3. Keep it under 200 lines with markdown headers — longer files consume more context and reduce adherence.
4. Auto memory works on its own (enabled by default, v2.1.59+): Claude saves build quirks, debugging insights, and preferences it discovers to `~/.claude/projects/<project>/memory/` and recalls them next session.
5. Run `/memory` anytime to see every loaded instruction file, open one in your editor, or toggle auto memory.

### Composing with Other Approaches (Intermediate)

- **Memory plus Plan Mode**: a plan built on documented architecture facts is grounded from the first sentence — Claude doesn't waste the analysis phase rediscovering what CLAUDE.md already states.
- **Memory plus hooks**: split guidance from enforcement. "Prefer small commits" belongs in CLAUDE.md; "never push to main" belongs in a hook that blocks it regardless of what the model decides.
- **Memory plus subagents**: subagents can maintain their own auto memory, so a recurring reviewer agent accumulates project knowledge across invocations just like the main session does.

### Advanced Patterns

- **Path-scoped rules**: split instructions into `.claude/rules/*.md` files with `paths:` frontmatter (glob patterns like `src/api/**/*.ts`). The rule loads only when Claude works on matching files — modular for teams, cheaper on context.
- **Imports**: `@path/to/file` syntax pulls other files into CLAUDE.md at launch (up to four hops). Use `@AGENTS.md` to share one instruction file with other coding agents, plus Claude-specific lines below the import.
- **Scope layering**: `~/.claude/CLAUDE.md` for personal preferences everywhere, `./CLAUDE.md` for team-shared project facts, gitignored `./CLAUDE.local.md` for personal project-specific notes, and `claudeMdExcludes` to skip other teams' files in a monorepo.

## Common Pitfalls

- **Treating CLAUDE.md as enforcement**: instructions shape behavior but nothing guarantees compliance. If a rule must always hold, encode it as a hook or a permission rule; keep CLAUDE.md for guidance.
- **Letting it bloat**: every line loads every session. Past ~200 lines adherence drops and you pay context for instructions that rarely matter — move niche content to path-scoped rules or skills.
- **Contradicting yourself**: conflicting instructions across CLAUDE.md, nested files, and rules make Claude pick one arbitrarily. Review periodically and delete stale entries.
- **Never auditing auto memory**: Claude's own notes are plain markdown under `~/.claude/projects/<project>/memory/` — browse them via `/memory` occasionally and delete anything wrong, or it stays wrong in every future session.

## Real-World Example

You inherit a Go service where the previous owner left no docs. First session, you run `/init` — Claude generates a CLAUDE.md capturing the `make test-integration` command, the `internal/` package layout, and the fact that migrations run through a custom `./scripts/migrate.sh` wrapper.

Over the next week you append the corrections you catch yourself repeating: "gRPC handlers must call `auth.Verify()` before touching the store", "the `orders` table is append-only — never generate UPDATE migrations against it". Meanwhile auto memory quietly records that the integration tests need Docker running and that you prefer table-driven tests.

Two weeks later a teammate opens the repo with Claude Code for the first time. Their first session already builds correctly, tests correctly, respects the append-only table, and writes table-driven tests — the entire ramp-up you went through, inherited for free from `git pull`.

## Sources

- [How Claude remembers your project](https://code.claude.com/docs/en/memory) — Official docs for CLAUDE.md files, .claude/rules/, imports, and auto memory
