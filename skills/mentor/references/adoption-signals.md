# Adoption Signals
*Last reviewed: 2026-07-03*

Observable evidence that an engineer already knows and uses each approach. The mentor computes the ignorance map by subtraction: the approaches below, minus those with positive signals, minus those recorded in the mentor profile (see `profile-schema.md`), ranked by leverage for the work observed in the session and repo.

Two signal tiers:

- **Setup & repo signals** — cheap to check at any moment (a glob or a grep against the project and `~/.claude/`). Re-verify these each invocation rather than trusting the profile; configuration is ground truth.
- **Session signals** — only visible when they happen (a command in the transcript, a tool use). The profile is the accumulator for these: record them when observed, because they cannot be re-checked later.

Absence is weak evidence, graded by tier: no hooks in settings after months of use genuinely suggests hooks are unknown; no `/rewind` in one transcript means nothing. Treat setup-signal absence as "likely unknown" and session-signal absence as "no information".

**Sourcing rule:** session signals come from the *current* conversation only — the skill runs inside the live session, so the transcript so far is already in context; nothing needs to be fetched. Never parse stored transcript files (`~/.claude/projects/<project>/*.jsonl`): their format is documented as internal and version-unstable. Past sessions are represented exclusively by what the profile accumulated when they happened; sessions before the mentor's first use are simply unobserved, and the cold start leans on setup signals plus the current conversation.

| Approach | Setup & repo signals | Session signals |
|----------|---------------------|-----------------|
| autonomous-loops | — | `/loop` or `/goal` in the transcript; goal-conditioned prompts |
| background-agents | — | `claude agents`, `claude --bg`, or `/bg` usage; talks about dispatching tasks |
| browser-integration | — | `claude --chrome` or `/chrome` usage; asks Claude to drive a browser |
| built-in-review-skills | Review commands wired into CI workflows | `/code-review`, `/security-review`, `/simplify`, `/verify`, or `/run` in the transcript |
| channels | A channel plugin installed (telegram, discord, imessage, fakechat) | `--channels` flag mentions; talks about pushing webhooks/chat into a session |
| checkpoints-rewind | — | `/rewind` usage; "undo that" / restore-checkpoint interactions |
| cloud-sessions | — | `claude --cloud` (or the deprecated `--remote`), `/teleport`, claude.ai/code or mobile-app mentions |
| custom-agents | `.claude/agents/*.md` or `~/.claude/agents/*.md` exists | References their own named agents |
| custom-plugins | `.claude-plugin/` in a repo they own; a marketplace file they maintain | Talks about publishing or packaging a plugin |
| custom-skills | `.claude/skills/*/SKILL.md` or `~/.claude/skills/*/SKILL.md` exists (their own, not plugin-installed) | Invokes their own slash commands |
| deep-research | — | `/deep-research` in the transcript |
| fan-out-workflows | `.claude/workflows/` exists | Workflow/ultracode usage; asks for parallel multi-agent runs |
| headless-mode | `claude -p` in `.github/workflows/`, `.gitlab-ci.yml`, or scripts | Discusses non-interactive/CI invocations |
| hooks-as-workflow | `hooks` configured in `.claude/settings.json`, `.claude/settings.local.json`, or `~/.claude/settings.json` | Asks about automating actions on edits/commits |
| lsp-self-correction | An LSP plugin installed for the project language | Mentions go-to-definition / diagnostics-driven fixes |
| mcp-context | `.mcp.json` exists, or MCP servers configured in settings | Uses MCP-backed tools; mentions connecting external systems |
| model-effort-selection | — | `/model`, `/effort`, `/fast`, or `/usage` in the transcript |
| official-plugins | Plugins installed from `claude-plugins-official` (visible via `/plugin list`) | `/plugin` usage; references marketplace plugins |
| plan-mode | — | Enters plan mode (Shift+Tab); asks for a plan before edits |
| project-memory | `CLAUDE.md` (root or `.claude/`), `.claude/rules/`, or `CLAUDE.local.md` exists; substantive auto-memory `MEMORY.md` | Asks Claude to remember things; `/memory` usage |
| safe-autonomy | `permissions` rules or `sandbox` configuration in settings | Discusses permission modes; runs with elevated autonomy deliberately |
| scheduled-agents | — | `/schedule` usage; mentions routines or recurring runs |
| session-context-management | — | `/context`, `/compact`, `/btw`, `/clear`, `/resume`, or `/branch` in the transcript |
| subagent-delegation | `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` in settings env (teams) | Spawns subagents; asks for parallel investigation |
| visual-artifacts | — | Publishes artifacts; asks for rendered/shareable pages |
| worktree-isolation | `.claude/worktrees/` exists | `claude --worktree` usage; mentions isolated parallel checkouts |
