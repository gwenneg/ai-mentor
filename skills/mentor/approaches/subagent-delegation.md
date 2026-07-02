# Subagent Delegation
*Last reviewed: 2026-07-01*

## What It Is

Subagent Delegation lets you spawn separate AI worker agents from your main session, each with its own context window and task. The main agent stays focused on the big picture while workers handle specific subtasks. When a worker finishes, it sends a summary back — not its entire conversation — so your main context stays clean and focused.

## Why It Works

This mirrors the way effective engineering teams work: a tech lead breaks a problem into pieces and assigns each piece to a specialist, rather than trying to hold every detail in their own head. AI context windows are finite and degrade in quality as they fill up. By delegating subtasks to fresh agents with clean context, each agent operates at peak effectiveness on a narrow problem. The main agent acts as an orchestrator — it synthesizes results and makes architectural decisions without being buried in implementation details.

## When to Use It

- Multi-file code review where you want each agent to deeply analyze one module rather than skimming all of them
- Parallel refactoring across independent services or packages in a monorepo
- Research tasks where you need to search for patterns across a large codebase while keeping the main session focused on decision-making
- CI/CD orchestration where different agents handle linting, testing, and deployment checks concurrently

## When NOT to Use It

- Small, focused tasks where the overhead of spawning and summarizing outweighs the benefit
- Tightly coupled changes where Agent B needs to see exactly what Agent A wrote — the summary may lose critical details
- When you need the agent to interact with you iteratively — subagents run autonomously and return results, they do not ask follow-up questions

## How It Works

### Basic (Beginner)

1. In your Claude Code session, describe a task that has separable pieces
2. Claude spawns a subagent: it creates a new agent with a specific prompt and task scope
3. The subagent reads files, searches code, or makes edits independently
4. When the subagent finishes, it returns a summary to the main session
5. The main session uses the summary to inform its next steps or synthesize a final answer

For example: "Review this PR for both correctness bugs and security issues." Claude can spawn one subagent for correctness review and another for security review, then combine the findings.

### Composing with Other Approaches (Intermediate)

- **Subagents with worktree isolation**: Give each subagent `isolation: "worktree"` so they can all edit code in parallel without conflicts. Each produces a branch you can review independently.
- **Plan Mode then subagent execution**: Use Plan Mode to design the overall approach, then spawn subagents to execute each step of the plan. The plan becomes the task breakdown.
- **Subagent results into review skills**: After subagents make parallel changes, run `/code-review` on each worktree branch to verify the changes before merging.

### Advanced Patterns

- **Agent Teams** (experimental — requires `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` to be enabled): Multiple Claude instances share a task list and can send peer messages to coordinate. Unlike simple subagents (which report to a parent), team members are peers that pick up tasks and collaborate. This is useful for complex multi-step workflows where agents need to react to each other's findings.
- **Agent View monitoring**: When running 5-30 parallel subagents, use Agent View to monitor all sessions from one screen. You can see each agent's progress, intervene if one gets stuck, and track overall completion.
- **Typed subagents**: Choose agent types based on the task — `Explore` for read-only code search (faster, cannot edit), `Plan` for architecture design (cannot edit), or general-purpose for tasks that require both reading and writing.

## Common Pitfalls

- **Over-delegation**: Spawning a subagent for a task that takes 10 seconds in the main session wastes time on setup and summarization overhead. Use subagents when the task is large enough to justify isolation.
- **Lost nuance in summaries**: Subagents return summaries, not full transcripts. If a subagent found a subtle race condition, the summary might flatten it to "found a concurrency issue." For critical findings, ask the main agent to follow up with specific questions.
- **Context fragmentation**: If you spawn too many subagents with overlapping scopes, you get redundant work and potentially contradictory findings. Give each agent a clear, non-overlapping scope.
- **Ignoring agent types**: Using a general-purpose agent for a pure search task is slower and uses more resources than an Explore agent. Match the agent type to the task.

## Real-World Example

You are reviewing a large PR (47 files changed) that adds a new payment provider integration. Instead of asking one agent to review everything, you delegate:

```
> Review PR #342 which adds Stripe Connect support. I want three
  separate reviews: (1) correctness of the payment flow logic in
  services/payments/, (2) security review of the webhook handler
  in api/webhooks/stripe.go, and (3) test coverage analysis of the
  new test files in tests/payments/.
```

Claude spawns three subagents:
- Agent 1 reads `services/payments/stripe_connect.go` and `services/payments/payout.go`, finds that the idempotency key is generated after the API call instead of before, which could cause duplicate charges on retry.
- Agent 2 reads `api/webhooks/stripe.go`, confirms the webhook signature verification is correct but flags that the endpoint does not validate the event type before processing.
- Agent 3 reads the test files, reports that `TestPayoutCreation` does not cover the error path when the Stripe API returns a `402 Payment Required`.

Claude synthesizes all three reports into a unified review with prioritized findings, and you now have a thorough, multi-dimensional review that would have taken a single agent (or a single human) much longer.

## Sources

- [Claude Code Sub-Agents](https://code.claude.com/docs/en/sub-agents) — Official docs for creating and using custom subagents
- [Multi-Agent Research System](https://www.anthropic.com/engineering/multi-agent-research-system) — Anthropic engineering blog on orchestrator-worker multi-agent architecture
