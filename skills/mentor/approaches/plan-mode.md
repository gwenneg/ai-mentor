# Plan Mode
*Last verified: 2026-07-06*

## What It Is

Plan Mode tells Claude Code to read your code and propose a plan before making any changes. Instead of jumping straight to editing files, the AI analyzes the problem, outlines its approach, and waits for your approval. You stay in control of the strategy while the AI handles the analysis and execution.

## Why It Works

Separating analysis from execution catches flawed assumptions at the outline stage, where redirection is cheap — before they become flawed code spread across ten files.

## When to Use It

- Debugging a complex issue where you want the AI to gather evidence and form a hypothesis before touching code
- Refactoring that spans multiple files — review the rename/move plan before execution
- Greenfield architecture where you want to discuss structure, not just generate files
- Understanding an unfamiliar codebase — let the AI read and explain before proposing changes
- Migration planning where a wrong first step can cascade into hours of cleanup

## When NOT to Use It

- Quick, well-scoped changes like fixing a typo or updating a config value — planning adds friction with no benefit
- When you already have a precise plan and just need execution — skip straight to editing
- Rapid prototyping where you want to iterate fast and inspect results, not review proposals

## How It Works

### Basic (Beginner)

1. Start Claude Code with plan permissions: `claude --permission-mode plan`
2. Or switch mid-session: press `Shift+Tab` to cycle into Plan Mode, or prefix a single prompt with `/plan`
3. Describe your task: "The user login flow is broken when SSO tokens expire. Find the root cause and propose a fix."
4. Claude reads relevant files, traces the logic, and presents a structured plan
5. Review the plan. Approve it — picking whether Claude executes in auto, accept-edits, or manual mode — or refine: "Good analysis, but let's handle the token refresh in the middleware instead."

### Composing with Other Approaches (Intermediate)

- **Plan Mode plus Subagent Delegation**: Use Plan Mode to design a refactoring strategy, then spawn subagents to execute each piece in parallel. The plan becomes the task list.
- **Plan Mode plus Built-in Review Skills**: Have Claude plan the fix, execute it, then run `/code-review` on its own changes to catch issues the plan missed.
- **Plan Mode plus Worktree Isolation**: Enter a worktree first, then use Plan Mode for risky changes. If the plan goes sideways after approval, you can discard the entire worktree.

### Advanced Patterns

- **Iterative narrowing**: Start with a broad question ("Why are our API response times degrading?"), review Claude's initial analysis, then ask it to plan a fix for only the specific bottleneck it identified. Each cycle narrows scope.
- **Plan as documentation**: Ask Claude to save its plan as a design doc or ADR before executing. The plan becomes a record of the reasoning behind the change.
- **Multi-option planning**: Ask Claude to propose two or three alternative approaches with trade-offs, then pick one. This surfaces options you might not have considered.

## Common Pitfalls

- **Over-planning simple tasks**: If you use Plan Mode for a one-line fix, you waste time reviewing a plan that is more complex than the change itself. Reserve it for tasks where a wrong approach has real cost.
- **Rubber-stamping plans**: The value comes from actually reading and questioning the plan. If you approve without reviewing, you get the latency cost of planning with none of the safety benefit.
- **Stale context after replanning**: If you reject a plan and redirect multiple times, the conversation can get long and confused. Consider starting a fresh session with your refined requirements.

## Sources

- [Permission modes](https://code.claude.com/docs/en/permission-modes) — Official docs on plan mode: how to enter it, what it can and cannot do, and the plan approval flow
- [Claude Code Best Practices](https://code.claude.com/docs/en/best-practices) — Official best-practices guide with the explore, plan, then code workflow
