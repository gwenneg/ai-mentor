# Subagent Delegation
*Last verified: 2026-07-12*

## What It Is

Subagent Delegation lets you spawn separate AI worker agents from your main session, each with its own context window and task. The main agent stays focused on the big picture while workers handle specific subtasks. When a worker finishes, it sends a summary back — not its entire conversation — so your main context stays clean and focused.

## Why It Works

Context windows are finite and degrade as they fill; delegating subtasks to fresh agents with clean context lets each operate at peak effectiveness while the orchestrator keeps only the summaries.

## When to Use It

- Multi-file code review where you want each agent to deeply analyze one module rather than skimming all of them
- Parallel refactoring across independent services or packages in a monorepo
- Research tasks where you need to search for patterns across a large codebase while keeping the main session focused on decision-making
- CI/CD orchestration where different agents handle linting, testing, and deployment checks concurrently

## When NOT to Use It

- Small, focused tasks where the overhead of spawning and summarizing outweighs the benefit
- Tightly coupled changes where Agent B needs to see exactly what Agent A wrote — the summary may lose critical details (a *fork* subagent, which inherits the whole conversation, covers some of these cases — see Advanced Patterns)
- When you need the agent to interact with you iteratively — subagents run autonomously and return results, they do not ask follow-up questions

## How It Works

### Basic (Beginner)

1. Describe a task with separable pieces, or ask for delegation outright: "Use a subagent to find every caller of the legacy API". Claude spawns subagents through its `Agent` tool — automatically when a side task would flood your context, or whenever you ask
2. Built-in agent types cover the common cases: `Explore` (fast, read-only search — Claude picks a thoroughness level from quick to very thorough), `Plan` (research during plan mode), and `general-purpose` (reads and writes, multi-step tasks)
3. The subagent works in its own context window and returns only a summary — the search results, logs, and file dumps it processed never enter your main conversation. Subagents run in the background by default (v2.1.198+): Claude keeps working while they run and picks up results when they finish, dropping to the foreground only when it needs the result before continuing, and background subagents surface their permission prompts in your main session
4. Name a custom subagent to invoke it directly: "Use the code-improver agent to suggest improvements in this project". Definitions live in `.claude/agents/` (see Custom Agent Definitions); pin an agent's foreground/background behavior with the `background` frontmatter field
5. The main session synthesizes the summaries and makes the decisions

For example:

```
> Review this PR for both correctness bugs and security issues.
```

Claude can spawn one subagent for correctness review and another for security review, then combine the findings.

### Composing with Other Approaches (Intermediate)

- **Subagents with worktree isolation**: Give each subagent `isolation: "worktree"` so they can all edit code in parallel without conflicts. Each produces a branch you can review independently; see Worktree Isolation for the `worktree.baseRef` and uncommitted-changes details.
- **Plan Mode then subagent execution**: Use Plan Mode to design the overall approach, then spawn subagents to execute each step of the plan. The plan becomes the task breakdown.
- **Subagent results into review skills**: After subagents make parallel changes, run `/code-review` on each worktree branch to verify the changes before merging.

### Advanced Patterns

- **Agent Teams** (experimental — requires `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` to be enabled): Multiple Claude instances share a task list and message each other directly via `SendMessage`. Unlike simple subagents (which report to a parent), team members are peers: tasks can declare dependencies, and teammates self-claim the next unblocked task when they finish one. Best for work requiring discussion — competing debugging hypotheses, multi-lens reviews — at a significantly higher token cost than subagents.
- **Scaling beyond one session**: Subagents live inside your session and do not appear in agent view (`claude agents`) — that board tracks background sessions. When you want many parallel workers you can monitor and steer from one screen, dispatch separate background sessions instead (see Background Agents).
- **Resume, fork, nest**: Subagents are resumable — ask Claude to resume one and it retains its full conversation history, picking up where it stopped (the built-in Explore and Plan types are one-shot and can't be). A *fork* inherits the entire parent conversation instead of starting fresh — the tool for tasks that need everything, not a summary. And subagents can spawn their own subagents (v2.1.172+, depth capped at five levels), so a reviewer can dispatch a verifier per finding without the intermediate output ever reaching your session.
- **Typed subagents**: Choose agent types based on the task — `Explore` for read-only code search (fast, cannot edit) or general-purpose for tasks that require both reading and writing. The built-in `Plan` type is mainly Claude's own: it spawns it for read-only codebase research while in Plan Mode.

## Common Pitfalls

- **Over-delegation**: Spawning a subagent for a task that takes 10 seconds in the main session wastes time on setup and summarization overhead. Use subagents when the task is large enough to justify isolation.
- **Lost nuance in summaries**: Subagents return summaries, not full transcripts. If a subagent found a subtle race condition, the summary might flatten it to "found a concurrency issue." For critical findings, ask the main agent to follow up with specific questions.
- **Context fragmentation**: If you spawn too many subagents with overlapping scopes, you get redundant work and potentially contradictory findings. Give each agent a clear, non-overlapping scope.
- **Ignoring agent types**: Using a general-purpose agent for a pure search task is slower and uses more resources than an Explore agent. Match the agent type to the task.

## Sources

- [Claude Code Sub-Agents](https://code.claude.com/docs/en/sub-agents) — Official docs for creating and using custom subagents
- [Agent Teams](https://code.claude.com/docs/en/agent-teams) — Official docs for coordinating peer Claude instances with shared tasks and inter-agent messaging
- [Multi-Agent Research System](https://www.anthropic.com/engineering/multi-agent-research-system) — Anthropic engineering blog on orchestrator-worker multi-agent architecture

## Signals

- Setup: `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` in settings env (teams)
- Session: Spawns subagents; asks for parallel investigation
