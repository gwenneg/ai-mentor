# Background Agents
*Last verified: 2026-07-12*

## What It Is

Background Agents are full Claude Code sessions that keep running without a terminal attached. You dispatch a task, close the view, and keep working — the agent works on, edits files in its own isolated worktree, and surfaces when it finishes or needs your input. Agent view (`claude agents`, currently a research preview) is the control room: every background session across all your projects on one screen, grouped by what's running, what's done, and what's waiting on you.

## Why It Works

Background agents convert attention-bound work into queue-bound work: thirty seconds specifying the task, ten minutes reviewing the result, and the hour in between belongs to something else.

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

Alternate entry points: `claude --bg "fix the flaky SettingsChangeDetector test"` dispatches straight from your shell (add `--name` to label it); `/background` (alias `/bg`) pushes your *current* interactive session into the background mid-task; and `/fork` (v2.1.207+) copies the conversation so far into a *new* background session with its own row in agent view, leaving your current session running untouched (the in-session subagent `/fork` used to launch is now `/subtask` — see Subagent Delegation).

### Composing with Other Approaches (Intermediate)

- **Background agents plus autonomous loops**: dispatch a session whose task is a measurable condition ("make all tests in packages/billing pass") — goal-driven iteration with nobody watching, surfacing only when green or stuck.
- **Background agents plus worktree isolation**: the worktree move is automatic, but you can also pre-create a worktree and dispatch the agent inside it when you want to control the branch and base yourself.
- **Background agents plus session & context management**: when an interactive task turns mechanical halfway through, `/bg run the remaining fixes and the full suite` hands it off with one instruction — your terminal is free and the context carries over.

### Advanced Patterns

- **Fleet dispatch**: dispatch several sessions from agent view in a row — one per module, one per bug — and use `claude agents --cwd <path>` to filter the board per project. `@<repo>` in a dispatch prompt targets a child repository from a parent directory, and `! <command>` runs a plain shell command as a monitored background job on the same board — no model in the loop.
- **Shell-first management**: skip the board entirely — `claude attach`, `claude logs`, `claude stop`, `claude respawn`, and `claude rm` manage background sessions straight from your shell, and `claude daemon status` checks on the machinery running them.
- **Finish-line automation**: background agents that complete code work in a worktree commit, push, and open a draft PR when they finish (v2.1.198+), so "review the result" means reviewing a PR, not hunting for a branch.
- **Notification wiring**: while agent view is open, blocked or finished sessions fire the `Notification` hook with `agent_needs_input` / `agent_completed` events — wire it to a desktop notification or Slack webhook, park the board in a spare tab, and stop watching it.

## Common Pitfalls

- **Underspecified dispatches**: a background agent can't ask you cheap clarifying questions mid-flow the way an interactive session can — vague tasks come back wrong or blocked. Spend the extra minute specifying the finish line.
- **Forgetting the work lands in a worktree**: your working directory won't show the changes. Check the session's worktree (peek shows the path) or wait for the draft PR rather than concluding the agent did nothing.
- **Deleting sessions carelessly**: deleting a session from agent view also removes the worktree Claude created for it, *including uncommitted changes*. Merge or push what you want to keep first.
- **Dispatching conflicting tasks**: two agents editing the same module produce two divergent branches you must reconcile. Split parallel work along module boundaries, not within them.

## Sources

- [Manage multiple agents with agent view](https://code.claude.com/docs/en/agent-view) — Official docs for claude agents, /background, dispatching, attaching, and worktree isolation

## Signals

- Setup: —
- Session: `claude agents`, `claude --bg`, or `/background` usage; talks about dispatching tasks
