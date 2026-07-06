# Permissions & Safe Autonomy
*Last verified: 2026-07-06*

## What It Is

Permissions & Safe Autonomy is the practice of tuning what your AI agent may do without asking — so it can work autonomously where that's safe and must stop where it isn't. Claude Code's permission system layers allow/ask/deny rules, permission modes (from read-only plan mode to fully autonomous), and OS-level sandboxing. Tuned well, the agent stops interrupting you for harmless commands and physically cannot touch the things that matter.

## Why It Works

Autonomy is a function of trust boundaries, not model capability. An agent that prompts for every `npm test` trains you to approve without reading — which is *less* safe, not more. An agent with crisp boundaries (these commands always fine, these paths never touchable, everything else asks) can run long autonomous stretches while the dangerous surface stays provably closed. This is least-privilege from security engineering applied to AI: precision about what's denied is what makes broad allowance safe — and crucially, permission rules are enforced by Claude Code itself, so unlike CLAUDE.md instructions they hold regardless of what the model decides.

## When to Use It

- Permission prompts interrupt you many times per session for commands you always approve
- You want autonomous loops or long unattended runs without granting blanket trust
- Certain paths must be untouchable: production configs, `.env` files, lockfiles, migration history
- You're setting up team-wide guardrails that individual sessions shouldn't be able to loosen

## When NOT to Use It

- Brand-new to Claude Code — run with defaults for a week first, then codify what you actually kept approving
- As a substitute for isolation on truly risky work — for that, combine permissions with sandboxing, a worktree, or a cloud session rather than relying on rules alone

## How It Works

### Basic (Beginner)

1. Run `/permissions` to see every allow/ask/deny rule and which settings file each comes from.
2. Allow the safe, frequent commands — add rules like `Bash(npm run test *)` and `Bash(git commit *)` to `.claude/settings.json`, or run `/fewer-permission-prompts` to have your transcripts scanned for read-only commands worth allowlisting.
3. Deny the untouchables: `Read(./.env)`, `Edit(/config/production/**)`, `Bash(git push *)`. Deny always wins — rules evaluate deny, then ask, then allow, and no other settings layer can override a deny.
4. Pick the mode for the work: `plan` to explore without edits, `default` for normal prompting, `acceptEdits` when you're fine with file changes landing without per-edit approval.
5. Mind the syntax details: `Bash(ls *)` needs the space before `*` (word boundary), and compound commands like `a && b` must match rules for *each* subcommand — Claude Code parses shell operators, so `Bash(safe-cmd *)` doesn't smuggle in `safe-cmd && rm -rf`.

### Composing with Other Approaches (Intermediate)

- **Permissions plus autonomous loops**: allowlist the loop's inner cycle (test runner, linter, file edits via `acceptEdits`) and deny the escape hatches (pushes, deletions of test files) — the loop runs unattended and can't cheat destructively.
- **Permissions plus hooks**: rules match patterns; hooks evaluate logic. A PreToolUse hook can block conditionally (e.g. edits to files with pending migrations) — and deny rules still apply regardless of what a hook returns.
- **Permissions plus sandboxing**: complementary layers (enable it with `/sandbox`) — rules govern what the agent may attempt; the OS-level sandbox restricts what Bash and its child processes can reach even if a prompt injection gets past the model's judgment. With the sandbox on, sandboxed commands run without prompts by default while deny rules and content-scoped ask rules like `Bash(git push *)` still hold (only a bare `Bash` ask rule is skipped for sandboxed commands).

### Advanced Patterns

- **Parameter-scoped rules**: deny/ask rules can match tool parameters — `Agent(isolation:worktree)`, `Bash(run_in_background:true)` — for policies on *how* tools are used, not just which.
- **Mode ceilings for real autonomy**: `bypassPermissions` skips prompts (keep it for containers/VMs — root/home `rm -rf` still trips a circuit breaker); `auto` mode auto-approves with background safety checks (research preview). Organizations can disable both via `disableBypassPermissionsMode` / `disableAutoMode` in managed settings.
- **Team guardrails via settings precedence**: managed settings > CLI flags > local project > shared project > user — a deny at any level cannot be re-allowed at any other level, so checked-in project denies become team-wide invariants.

## Common Pitfalls

- **Approval fatigue as policy**: leaving everything on "ask" feels safe but trains reflexive yes-clicking. Decide once in rules; save prompts for genuinely ambiguous actions.
- **Argument-constrained Bash patterns**: `Bash(curl http://github.com/ *)` is trivially bypassed by flags, redirects, or variables. Deny the network tools in Bash and grant `WebFetch(domain:github.com)` instead, or enforce with a PreToolUse hook.
- **Expecting Read/Edit denies to bind subprocesses**: they cover Claude's file tools and recognized file commands, not a Python script opening files itself. For OS-level enforcement, enable the sandbox.
- **`bypassPermissions` on your real machine**: it skips protections including writes to `.git` and `.claude`. Isolated environments only.

## Real-World Example

Your team's autonomous test-fixing loops keep stalling on permission prompts, so developers run them in `bypassPermissions` — on their laptops, against the repo with the production Terraform in it. You replace bravado with boundaries in `.claude/settings.json`:

```json
{
  "permissions": {
    "allow": ["Bash(npm run test *)", "Bash(npm run lint *)", "Bash(git add *)", "Bash(git commit *)"],
    "deny": ["Bash(git push *)", "Edit(/infra/**)", "Read(./.env)", "Edit(**/*.snap)"]
  }
}
```

Committed to the repo, this applies to everyone: loops run test/lint/commit cycles without a single prompt, nobody's session can push, touch infra, read secrets, or "fix" a failing test by editing its snapshot — and because a deny from any settings level can't be re-allowed by any other, a `--allowedTools` flag can't loosen it. The loops now run unattended to completion, and the two prompts that remain (`git push`, infra edits) are the two that deserve a human.

## Sources

- [Configure permissions](https://code.claude.com/docs/en/permissions) — Official reference for rules, modes, precedence, and hook interaction
- [Sandboxing](https://code.claude.com/docs/en/sandboxing) — OS-level filesystem and network isolation for Bash commands
