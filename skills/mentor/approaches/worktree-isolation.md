# Worktree Isolation
*Last reviewed: 2026-07-01*

## What It Is

Worktree Isolation gives each AI coding session its own independent copy of your repository. Every session works on a separate git branch in a separate directory, so multiple sessions can edit the same files simultaneously without conflicts. Think of it as giving each AI worker their own desk with their own copy of the blueprints — they can scribble, rearrange, and experiment without affecting anyone else's workspace.

## Why It Works

Isolation is one of the oldest reliability principles in software engineering — from process isolation in operating systems to container isolation in deployment. When AI agents edit code, they face the same concurrency problems that humans do: two agents editing the same file at the same time will produce conflicts. Worktrees solve this at the git level by giving each session a physically separate working directory on its own branch. This means you can run aggressive, speculative changes without risking your working branch, and you can run multiple agents in parallel without coordination overhead.

## When to Use It

- Debugging a flaky test where you want to add logging and instrumentation without polluting your working branch
- Running multiple refactoring tasks in parallel, each touching overlapping files
- Testing a risky migration (database schema, framework upgrade) where you want a clean rollback path
- Exploring two different implementation approaches simultaneously to compare results
- Letting an AI agent experiment freely without risk to your current working state

## When NOT to Use It

- Simple, sequential changes where you are only running one session at a time — worktrees add setup overhead for no benefit
- When the task requires modifying gitignored state (databases, build caches, local server state) that cannot be replicated via `.worktreeinclude`
- Quick questions or code explanations that do not involve any file edits
- Very large monorepos where the disk cost of a full working copy is prohibitive

## How It Works

### Basic (Beginner)

1. Start Claude Code with a worktree: `claude --worktree`
2. Claude creates a new branch and working directory under `.claude/worktrees/<name>/`
3. Work normally — all edits happen in the isolated copy, your original branch is untouched
4. When done, review changes with `git diff` or `git log` in the worktree
5. If the changes look good, merge the worktree branch into your main branch
6. If the changes are bad, discard the entire worktree — no cleanup needed
7. Exit the worktree: Claude prompts you to keep or remove it

### Composing with Other Approaches (Intermediate)

- **Worktree plus Plan Mode**: Enter a worktree, then use Plan Mode for a risky refactoring. If the plan looks wrong after execution, discard the entire worktree — no `git reset` gymnastics needed.
- **Worktree plus subagents**: Spawn subagents with `isolation: "worktree"` so each agent edits in its own worktree. Three agents can refactor three different services in parallel without merge conflicts.
- **Worktree for A/B approaches**: Create two worktrees, give each a different implementation strategy for the same problem, then compare the results side by side.

### Advanced Patterns

- **`.worktreeinclude` for environment parity**: Create a `.worktreeinclude` file listing gitignored files your project needs (`.env`, `.env.local`, vendored binaries). These get copied into each new worktree automatically, so the isolated environment actually works.
- **Subagent fan-out with worktree isolation**: In a fan-out workflow, each subagent gets `isolation: "worktree"`. This lets you run 5-10 parallel agents that all edit code without any coordination — each produces a branch you can review and merge independently.
- **Long-lived worktrees for feature branches**: Keep a worktree alive across sessions for a multi-day feature. Re-enter it with `EnterWorktree` using the path to resume exactly where you left off. This is particularly useful for multi-day migrations where you want to checkpoint progress without merging incomplete work.

## Common Pitfalls

- **Forgetting to set `worktree.baseRef`**: By default, worktrees branch from `origin/<default-branch>`. If you want to branch from your current HEAD instead, set `worktree.baseRef` to `head` in your Claude settings.
- **Forgetting `.worktreeinclude`**: Your worktree will not have `.env`, local configs, or other gitignored files unless you list them. If your app fails to start in the worktree, this is almost always why.
- **Worktree accumulation**: Each worktree is a full copy of your repo on disk. If you forget to clean them up, they consume disk space. Use `git worktree list` periodically and remove stale ones.
- **Merging divergent worktrees**: If two worktrees edit overlapping files, you still get merge conflicts when combining them — isolation defers the conflict, it does not eliminate it. Plan your parallel work to minimize overlap.
- **Assuming worktree state persists across sessions**: If you exit Claude Code without keeping the worktree, it may be cleaned up. Explicitly choose "keep" when prompted if you plan to return to the work later.

## Real-World Example

You need to upgrade your project from React Router v5 to v6, but you are not sure it will go smoothly. You also have a teammate waiting on a bug fix in the same codebase.

```
claude --worktree
> Upgrade all React Router v5 APIs to v6 in src/routes/ and src/components/.
  Update useHistory to useNavigate, Switch to Routes, and component prop
  to element prop.
```

Claude creates `.claude/worktrees/router-upgrade/` on branch `router-upgrade`, then methodically updates 14 files. Meanwhile, you open a second terminal and work on the bug fix in your normal working directory — no conflicts, no stashing, no branch switching.

After Claude finishes, you run `npm test` in the worktree directory. Three tests fail in `src/routes/__tests__/ProtectedRoute.test.tsx` because the test setup still uses `MemoryRouter` from v5. You ask Claude to fix the tests, it does, and all 47 tests pass. You merge the branch into `main` and remove the worktree with a clean history. Your teammate's bug fix, developed concurrently on the main working directory, never experienced a single conflict.

## Sources

- [Claude Code Sub-Agents](https://docs.anthropic.com/en/docs/claude-code/sub-agents) — Official docs covering worktree isolation mode for subagents
- [Claude Code IDE Integrations](https://docs.anthropic.com/en/docs/claude-code/ide-integrations) — VS Code integration docs covering the --worktree flag
