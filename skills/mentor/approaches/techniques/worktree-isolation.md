# Worktree Isolation
*Last verified: 2026-07-12*

## What It Is

Worktree Isolation gives each AI coding session its own independent copy of your repository. Every session works on a separate git branch in a separate directory, so multiple sessions can edit the same files simultaneously without conflicts. Think of it as giving each AI worker their own desk with their own copy of the blueprints — they can scribble, rearrange, and experiment without affecting anyone else's workspace.

## Why It Works

A physically separate working directory on its own branch lets aggressive, speculative changes run in parallel without coordination overhead or risk to your working branch.

## When to Use It

- Debugging a flaky test where you want to add logging and instrumentation without polluting your working branch
- Running multiple refactoring tasks in parallel, each touching overlapping files
- Testing a risky migration (database schema, framework upgrade) where you want a clean rollback path
- Exploring two different implementation approaches simultaneously to compare results
- Letting an AI agent experiment freely without risk to your current working state

## When NOT to Use It

- Simple, sequential changes where you are only running one session at a time — worktrees add setup overhead for no benefit
- When the task requires modifying gitignored state (databases, build caches, local server state) that cannot be replicated via `.worktreeinclude`
- Very large monorepos where the disk cost of a full working copy is prohibitive

## How It Works

### Basic (Beginner)

1. Start Claude Code with a worktree: `claude --worktree feature-name` (omit the name and Claude generates one)
2. Claude creates the working directory `.claude/worktrees/<name>/` on a new branch named `worktree-<name>`
3. Work normally — all edits happen in the isolated copy, your original branch is untouched
4. When done, review changes with `git diff` in the worktree; merge the branch if they look good, or discard the worktree if they don't
5. Exit the session: if the worktree has changes, Claude prompts you to keep or remove it; a clean worktree is removed automatically (unless the session has a name — then Claude prompts even when clean, so you can keep it for later)

### Composing with Other Approaches (Intermediate)

- **Worktree plus Plan Mode**: Enter a worktree, then use Plan Mode for a risky refactoring. If the plan looks wrong after execution, discard the entire worktree — no `git reset` gymnastics needed.
- **Worktree plus subagents**: Spawn subagents with `isolation: "worktree"` so each agent edits in its own worktree. Three agents can refactor three different services in parallel without merge conflicts.
- **Worktree plus Checkpoints**: Inside a worktree you get two levels of undo — rewind a checkpoint when a single edit went wrong, discard the whole worktree when the entire approach did.

### Advanced Patterns

- **`.worktreeinclude` for environment parity**: Create a `.worktreeinclude` file listing gitignored files your project needs (`.env`, `.env.local`, vendored binaries). These get copied into each new worktree automatically, so the isolated environment actually works.
- **Subagent fan-out with worktree isolation**: In a fan-out workflow, each subagent gets `isolation: "worktree"`. This lets you run 5-10 parallel agents that all edit code without any coordination — each produces a branch you can review and merge independently.
- **Long-lived worktrees for feature branches**: Keep a worktree alive across sessions for a multi-day feature. Ask Claude to switch back into it — the `EnterWorktree` tool enters an existing worktree under `.claude/worktrees/` by path (entering one outside that directory asks for confirmation, v2.1.206+) — and pick up the file state where you left off. This is particularly useful for multi-day migrations where you want to checkpoint progress without merging incomplete work.
- **PR-based worktrees**: `claude --worktree "#1234"` (or a full PR URL) fetches `pull/1234/head` from `origin` and creates the worktree at `.claude/worktrees/pr-1234` — a teammate's PR checked out for review or fixes without touching your branch.

## Common Pitfalls

- **Forgetting to set `worktree.baseRef`**: By default, worktrees branch from `origin/<default-branch>`. If you want to branch from your current HEAD instead, set `worktree.baseRef` to `head` in your Claude settings. Uncommitted changes never carry over into a new worktree — commit first if the work must build on them.
- **Forgetting `.worktreeinclude`**: Your worktree will not have `.env`, local configs, or other gitignored files unless you list them. If your app fails to start in the worktree, this is almost always why.
- **Worktree accumulation**: Each worktree is a full working copy of your project files on disk (git history is shared). If you forget to clean them up, they consume disk space. Use `git worktree list` periodically and remove stale ones.
- **Merging divergent worktrees**: If two worktrees edit overlapping files, you still get merge conflicts when combining them — isolation defers the conflict, it does not eliminate it. Plan your parallel work to minimize overlap.

## Sources

- [Run parallel sessions with worktrees](https://code.claude.com/docs/en/worktrees) — Official docs for the `--worktree` flag, `.worktreeinclude`, `worktree.baseRef`, and worktree cleanup
- [Claude Code Sub-Agents](https://code.claude.com/docs/en/sub-agents) — Official docs covering worktree isolation mode for subagents
- [Claude Code IDE Integrations](https://code.claude.com/docs/en/ide-integrations) — VS Code integration docs covering the --worktree flag

## Signals

- Setup: `.claude/worktrees/` exists
- Session: `claude --worktree` usage; mentions isolated parallel checkouts
