# Background Agents
*Last verified: 2026-07-06*

## What It Is

Background Agents are full Claude Code sessions that keep running without a terminal attached. You dispatch a task, close the view, and keep working — the agent works on, edits files in its own isolated worktree, and surfaces when it finishes or needs your input. Agent view (`claude agents`, currently a research preview) is the control room: every background session across all your projects on one screen, grouped by what's running, what's done, and what's waiting on you.

## Why It Works

Most engineering tasks don't need you watching them. A test-suite fix, a long audit, a mechanical refactor — the work is real but your attention adds nothing between kickoff and review. Background agents convert attention-bound work into queue-bound work: you spend thirty seconds specifying the task and ten minutes reviewing the result, and the hour in between belongs to something else. Because each session automatically moves into its own git worktree before editing, parallel agents can't trample each other or your working copy — the isolation that makes "fire and forget" actually safe.

## When to Use It

- A task is well-specified but long: fix the failing suite, update all call sites, write the missing tests
- You want to run several independent tasks in parallel without opening four terminals
- You're about to leave your desk and want work continuing while you're gone
- The current interactive session turned into a grind — push it to the background and reclaim your terminal

## When NOT to Use It

- Tasks needing frequent steering — ambiguous requirements mean the agent will either guess or sit blocked on "needs input"
- Quick edits where dispatch-and-review overhead exceeds just doing it interactively
- Parallel tasks that must edit the same files toward one coherent change — isolation defers the merge conflict, it doesn't remove it

## How It Works

### Basic (Beginner)

1. Run `claude agents` to open agent view, type a task, and press Enter to dispatch it as a background session. (`Shift+Enter` dispatches and attaches immediately.)
2. The session runs on its own. Before editing files, it moves into an isolated worktree under `.claude/worktrees/`, so your working copy stays untouched.
3. Watch the row update: sessions group under "Needs input", "Working", "Ready for review" (a PR is open), and "Completed". Peek at a row to read progress and reply without attaching.
4. Press `Enter` or `→` on a row to attach — the full conversation takes over your terminal, opening with a recap of what happened while you were away. Press `←` on an empty prompt to detach and go back to the table.
5. Review the result like any change: the diff lives on the session's worktree branch, ready to merge or discard.

Alternate entry points: `claude --bg "fix the flaky SettingsChangeDetector test"` dispatches straight from your shell (add `--name` to label it), and `/background` (alias `/bg`) pushes your *current* interactive session into the background mid-task.

### Composing with Other Approaches (Intermediate)

- **Background agents plus autonomous loops**: dispatch a session whose task is a measurable condition ("make all tests in packages/billing pass") — goal-driven iteration with nobody watching, surfacing only when green or stuck.
- **Background agents plus worktree isolation**: the worktree move is automatic, but you can also pre-create a worktree and dispatch the agent inside it when you want to control the branch and base yourself.
- **Background agents plus session & context management**: when an interactive task turns mechanical halfway through, `/bg run the remaining fixes and the full suite` hands it off with one instruction — your terminal is free and the context carries over.

### Advanced Patterns

- **Fleet dispatch**: dispatch several sessions from agent view in a row — one per module, one per bug — and use `claude agents --cwd <path>` to filter the board per project. `@<repo>` in a dispatch prompt targets a child repository from a parent directory, and `! <command>` runs a plain shell command as a monitored background job on the same board — no model in the loop.
- **Finish-line automation**: background agents that complete code work in a worktree commit, push, and open a draft PR when they finish (v2.1.198+), so "review the result" means reviewing a PR, not hunting for a branch.
- **Notification wiring**: while agent view is open, blocked or finished sessions fire the `Notification` hook with `agent_needs_input` / `agent_completed` events — wire it to a desktop notification or Slack webhook, park the board in a spare tab, and stop watching it.

## Common Pitfalls

- **Underspecified dispatches**: a background agent can't ask you cheap clarifying questions mid-flow the way an interactive session can — vague tasks come back wrong or blocked. Spend the extra minute specifying the finish line.
- **Forgetting the work lands in a worktree**: your working directory won't show the changes. Check the session's worktree (peek shows the path) or wait for the draft PR rather than concluding the agent did nothing.
- **Deleting sessions carelessly**: deleting a session from agent view also removes the worktree Claude created for it, *including uncommitted changes*. Merge or push what you want to keep first.
- **Dispatching conflicting tasks**: two agents editing the same module produce two divergent branches you must reconcile. Split parallel work along module boundaries, not within them.

## Real-World Example

Monday morning triage leaves you with three independent chores: a flaky date test in `packages/reports`, lint debt in `services/gateway`, and a missing integration test for last week's webhook handler. Instead of doing them serially, you open `claude agents` and dispatch all three, each with a concrete finish line:

```text
fix the flaky date test in packages/reports — done when pnpm test --filter reports passes 20 consecutive runs
```

You spend the morning on design work. The tab title reads `1 awaiting input · claude agents` around 11 — the lint agent found a rule violation that needs a human call. You peek, reply "disable that rule for generated files only", and go back to your doc.

By lunch, all three rows sit under Ready for review, each with a draft PR. You review three focused diffs in twenty minutes. Total attention spent on three chores: the dispatch minute, one peek, and the reviews — the rest of the work happened while you weren't looking.

## Sources

- [Manage multiple agents with agent view](https://code.claude.com/docs/en/agent-view) — Official docs for claude agents, /background, dispatching, attaching, and worktree isolation
