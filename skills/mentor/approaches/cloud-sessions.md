# Cloud Sessions & Remote Work
*Last verified: 2026-07-06*

## What It Is

Cloud Sessions run Claude Code on Anthropic-managed infrastructure instead of your machine. You hand off a task from your terminal with `claude --cloud`, from claude.ai/code in a browser, or from the Claude mobile app; the session clones your repo into a sandboxed cloud environment and works there. Sessions persist when you close your laptop, you can monitor them from your phone, and `--teleport` pulls any cloud session back into your terminal to continue locally.

## Why It Works

Cloud sessions decouple the work from the workstation: you keep the judgment-heavy parts where interaction is cheap and ship the execution-heavy parts to an environment where autonomy is safe by construction.

## When to Use It

- Well-scoped tasks you want off your machine: "fix this bug, run the tests, open a PR" while you keep working locally
- Several independent tasks at once — each `claude --cloud` creates its own parallel cloud session
- Keeping a PR healthy without babysitting it: auto-fix responds to CI failures and review comments as they arrive
- Working away from your desk — kick off or steer sessions from the mobile app

## When NOT to Use It

- Work depending on local state: uncommitted changes (the cloud clones from GitHub, not your disk), local databases, devices, or credentials that don't belong in a cloud environment
- Tight interactive loops where you steer every step — the value is autonomy; constant intervention negates it
- Repositories your organization hasn't cleared for cloud execution (Owners can disable web sessions org-wide), or setups authenticated by API key or Bedrock/Foundry — cloud sessions and `--teleport` require claude.ai subscription sign-in

## How It Works

### Basic (Beginner)

1. Push your branch (the cloud clones from GitHub, so local-only commits won't be there), then hand off: `claude --cloud "Fix the authentication bug in src/auth/login.ts"` (the older `--remote` spelling still works as a deprecated alias).
2. The session provisions a cloud environment, clones your current repo at your current branch, and starts working. The CLI shows a live setup checklist (v2.1.195+); messages you type while it provisions are queued.
3. Check progress with `/tasks` in your terminal, or open the session at claude.ai/code or in the mobile app to steer it like any conversation.
4. When it finishes, review the changes and create the PR from the web session.
5. Want it back on your machine? Run `/teleport` (or `claude --teleport`) and pick the session — it continues in your terminal with context intact. Handoff from CLI is one-way: you can pull cloud sessions down, and `--cloud` starts *new* cloud sessions; the Desktop app can push an existing local session to the web.

### Composing with Other Approaches (Intermediate)

- **Plan locally, execute remotely**: collaborate on the plan in local plan mode, commit the plan file, push, then `claude --cloud "Execute the migration plan in docs/migration-plan.md"` — your judgment on strategy, cloud autonomy on execution.
- **Cloud sessions plus code review**: `/autofix-pr` on a PR branch (requires the Claude GitHub App on the repo) spawns a web session that watches the PR — pushing fixes for clear CI failures and review comments, asking you first when a request is ambiguous or architecturally significant.
- **Cloud sessions plus plan mode, via ultraplan**: draft and review a plan in a web session while you keep working; comment on sections in the browser, then execute remotely or send the plan back to your terminal.

### Advanced Patterns

- **Parallel fleets**: fire several `--cloud` tasks back to back — independent sandboxes, no shared state, no coordination cost beyond reviewing the PRs.
- **Tuned environments**: configure cloud environments with setup scripts (cached between runs), environment variables, and network access levels — the default "Trusted" allowlist covers package registries and common dev domains; use Custom to add your own hosts or Full for unrestricted access.
- **Mobile-first supervision**: dispatch from the terminal before a commute, answer the one clarifying question from your phone, review the finished PR when you arrive.
- **Remote Control — the local-machine counterpart**: when the task needs your local environment (filesystem, MCP servers, local credentials), run `/remote-control` in a running session or `claude remote-control` in server mode, then continue that session from claude.ai/code or the Claude mobile app. The session keeps executing on your machine — the phone is a window into it, and nothing moves to the cloud. Enable mobile push notifications in `/config` ("Push when Claude decides" / "Push when actions required") to get pinged when it finishes or needs a decision.
- **Slack dispatch**: with the Claude app installed in your workspace and a repo connected at claude.ai/code, mentioning `@Claude` with a coding task in a channel creates a cloud session that gathers context from the thread, posts progress updates, and offers "View Session" and "Create PR" buttons when done — delegation without leaving the conversation where the bug was reported.

## Common Pitfalls

- **Forgetting to push first**: the most common failure — the cloud VM clones from GitHub, so a handoff referencing local-only commits starts from stale code.
- **Auto-fix in automation-heavy repos**: Claude replies to review threads under your GitHub account (labeled as Claude Code). If PR comments trigger automation like Atlantis or Terraform Cloud, a reply can run privileged workflows — review your repo's comment triggers before enabling auto-fix.
- **Expecting conflict resolution**: GitHub emits no webhook when the base branch advances into conflict, so auto-fix can't react to merge conflicts — open the session and ask for a rebase.
- **Under-provisioned environments**: a session without your setup script or required env vars fails in ways that look like model failure. If the task needs dependencies installed, configure the environment before dispatching.

## Sources

- [Use Claude Code on the web](https://code.claude.com/docs/en/claude-code-on-the-web) — Official docs for cloud sessions, environments, --cloud, /teleport, and auto-fix PRs
- [Remote Control](https://code.claude.com/docs/en/remote-control) — Official docs for continuing local sessions from a phone or browser, plus mobile push notifications
- [Claude Code in Slack](https://code.claude.com/docs/en/slack) — Official docs for delegating coding tasks from Slack channels
