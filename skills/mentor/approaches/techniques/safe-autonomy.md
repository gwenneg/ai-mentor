# Permissions & Safe Autonomy
*Last verified: 2026-08-15*

## What It Is

Permissions & Safe Autonomy is the practice of tuning what your AI agent may do without asking — so it can work autonomously where that's safe and must stop where it isn't. Claude Code's permission system layers allow/ask/deny rules, permission modes (from read-only plan mode to fully autonomous), and OS-level sandboxing. Tuned well, the agent stops interrupting you for harmless commands and physically cannot touch the things that matter.

## Why It Works

An agent that prompts for every `npm test` trains you to approve without reading; crisp boundaries — enforced by Claude Code itself, unlike CLAUDE.md instructions — let long autonomous stretches run while the dangerous surface stays provably closed.

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
4. Pick the mode for the work: `plan` to explore without edits, `default` (labeled Manual in the UIs) for per-action review, `acceptEdits` when you're fine with file changes landing without per-edit approval. Note the baseline has moved: new sessions start in `auto` mode by default when no permission mode is otherwise configured — a default you set yourself stays in place unless you accept the one-time switch prompt, and an organization-managed default doesn't change. This started on Pro, Max, and Team plans and, as of v2.1.283–v2.1.285, now applies on every plan and provider — interactive terminal and VS Code sessions, `claude -p`, and the Python Agent SDK alike, including third-party providers and telemetry-off setups — so the question is generally which mode to step *down* to for sensitive work, not whether auto is on. `permissions.defaultMode`/`--permission-mode` still overrides it. When auto mode is available, a Bash permission prompt in Manual or `acceptEdits` mode offers a **Yes, and switch to auto mode** option (v2.1.240+) that approves the command and flips the session to auto mode in one step.
5. Mind the syntax details: `Bash(ls *)` needs the space before `*` (word boundary), and compound commands like `a && b` must match rules for *each* subcommand — Claude Code parses shell operators, so `Bash(safe-cmd *)` doesn't smuggle in `safe-cmd && rm -rf`.

### Composing with Other Approaches (Intermediate)

- **Permissions plus autonomous loops**: allowlist the loop's inner cycle (test runner, linter, file edits via `acceptEdits`) and deny the escape hatches (pushes, deletions of test files) — the loop runs unattended and can't cheat destructively.
- **Permissions plus hooks**: rules match patterns; hooks evaluate logic. A PreToolUse hook can block conditionally (e.g. edits to files with pending migrations) — and deny rules still apply regardless of what a hook returns.
- **Permissions plus sandboxing**: complementary layers (enable it with `/sandbox`) — rules govern what the agent may attempt; the OS-level sandbox restricts what Bash and its child processes can reach even if a prompt injection gets past the model's judgment. With the sandbox on, sandboxed commands run without prompts by default while deny rules and content-scoped ask rules like `Bash(git push *)` still hold (only a bare `Bash` ask rule is skipped for sandboxed commands) — except where auto mode's classifier review runs server-side (see below), in which case read-only and sandboxed shell commands also wait for that review and can be blocked by it (v2.1.281). Sandbox credential entries can also *mask* instead of deny (`mode: "mask"`): commands see a sentinel value while the sandbox proxy injects the real credential on egress to hosts you allow — with `extract` patterns, JWT-aware `decode`, and AWS SigV4 re-signing for structured credentials (v2.1.224+) — honored only from user or managed settings, never a repository's checked-in settings. A per-command `allowed_domains` option on Bash, PowerShell, and Monitor (v2.1.271) narrows this further: the hosts a given command needs are reviewed and opened for that command alone, with every other host refused.

### Advanced Patterns

- **Parameter-scoped rules**: deny/ask rules can match tool parameters — `Agent(isolation:worktree)`, `Bash(run_in_background:true)` — for policies on *how* tools are used, not just which.
- **Mode ceilings for real autonomy**: `bypassPermissions` skips prompts except those forced by explicit `ask` rules (keep it for containers/VMs — root/home `rm -rf` still trips a circuit breaker); `auto` mode auto-approves with a classifier reviewing each action — all plans, available by default on Team and Enterprise (an Owner no longer has to enable it), and on Pro, Max, and Team the classifier's own calls don't count toward usage limits; `dontAsk` mode inverts the default — auto-denying anything not pre-approved by allow rules, for CI and restricted environments. Organizations can disable the first two via `disableBypassPermissionsMode` / `disableAutoMode` in managed settings.
- **Auto mode drops broad allow rules on entry**: blanket `Bash(*)`/`PowerShell(*)`, wildcarded interpreters like `Bash(python*)`, package-manager run commands, and `Agent` allow rules are suspended while auto mode is active (restored when you leave it); narrow rules like `Bash(npm test)` carry over — so tune allowlists to specific commands if you want them to keep working under the new default.
- **Team guardrails via settings precedence**: managed settings > CLI flags > local project > shared project > user — a deny at any level cannot be re-allowed at any other level, so checked-in project denies become team-wide invariants.
- **Third-party providers need no opt-in for auto mode** (v2.1.207+): auto mode now runs on Amazon Bedrock, Google Cloud's Agent Platform, and Microsoft Foundry without the `CLAUDE_CODE_ENABLE_AUTO_MODE` variable those platforms previously required; administrators can still turn it off with `disableAutoMode`.
- **"Always allow" persists across worktrees** (v2.1.211+): allow rules granted interactively now save to `.claude/settings.local.json` at the repository root, so an approval given in one git worktree or subdirectory holds in your other worktrees and future sessions instead of re-prompting.
- **The classifier's default blocks keep expanding**: as of v2.1.251+, auto mode also blocks by default actions such as requesting credentials from the cloud instance-metadata endpoint or connecting to sibling containers Claude didn't start, on top of the data-exfiltration and destructive-git protections already in place — and it now asks you before Claude's first read of a file outside your working directories, with an option to block such reads from then on. When the classifier blocks something, the denial reason given to Claude usually names the specific rule that matched (e.g. `[Data Exfiltration]`, v2.1.263+), which speeds up fixing a false-positive denial with an allow rule.
- **`defaultMode: "bypassPermissions"` no longer works from project settings** (v2.1.251+): setting it in a project's `.claude/settings.json` or `.claude/settings.local.json` is ignored and the session starts in Manual mode instead — set `bypassPermissions` in user or managed settings, or pass `--permission-mode`, to actually start a session in that mode.
- **Classifier locality and billing**: auto mode's safety classifier can run server-side (free of classifier overhead) or fall back to a local classifier call that counts toward usage. As of v2.1.278, this defaults to the server-side classifier for Claude API and Enterprise users, and on Amazon Bedrock, Google Vertex AI, Microsoft Foundry, and gateways (`CLAUDE_CODE_AUTO_MODE_SERVER=0` opts out on those platforms, with a warning if a session falls back to the billed local classifier); a direct Anthropic API connection with telemetry off defaults to server-side too (v2.1.282; `CLAUDE_CODE_AUTO_MODE_SERVER=0` opts out there as well). Where the classifier runs server-side, read-only and sandboxed shell commands now wait for that review and can be blocked by it too (v2.1.281). See [classifier billing](https://code.claude.com/docs/en/auto-mode-classifier-billing) if a session's usage looks higher than expected under auto mode.

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

## Signals

- Setup: `permissions` rules or `sandbox` configuration in settings
- Session: Discusses permission modes; runs with elevated autonomy deliberately
