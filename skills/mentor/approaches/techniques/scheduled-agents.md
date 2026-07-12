# Scheduled & Recurring Agents
*Last verified: 2026-07-12*

## What It Is

Scheduled & Recurring Agents run Claude Code without you starting it. A routine packages a prompt, one or more repositories, and a set of connectors, then runs on Anthropic-managed cloud infrastructure whenever a trigger fires: on a schedule (hourly, nightly, weekly, or a one-off timestamp), on an HTTP call from your own systems, or on GitHub events like a PR opening. Your laptop can be closed; the work happens anyway.

## Why It Works

Work that runs automatically happens every time; work that requires initiative happens until the third busy week — the same shift that made CI the default for testing.

## When to Use It

- Recurring maintenance with a clear outcome: nightly issue triage, weekly docs-drift review, scheduled dependency audits
- Event-driven work: review every PR when it opens, port merged changes to a sibling repository, verify each production deploy
- Wiring Claude into existing systems: your alerting tool POSTs to a routine's API endpoint and on-call reviews a draft-fix PR instead of a blank terminal
- Future one-offs: "in two weeks, open a cleanup PR removing this feature flag"

## When NOT to Use It

- Tasks needing your judgment mid-run — routines run autonomously with no approval prompts; everything must be decided in the prompt
- Sub-hourly cadences — the minimum schedule interval is one hour
- Work that needs your local machine's state — routines clone the repo fresh from GitHub each run (for local scheduled work, use Desktop scheduled tasks or `/loop` in an open session)

## How It Works

### Basic (Beginner)

1. From any session, run `/schedule` with a natural-language description: `/schedule daily PR review at 9am` or a one-off like `/schedule in 2 weeks, open a cleanup PR that removes the feature flag`. Claude walks through the setup and saves the routine to your account. (Research preview; requires a claude.ai subscription login. Recurring runs count against a daily per-account allowance; one-off runs are exempt and draw down regular subscription usage instead. Also manageable at claude.ai/code/routines.)
2. Write the prompt as if briefing someone who can't ask questions: what to do, where, and what success looks like. The routine runs with no permission prompts, so the prompt is the whole specification.
3. Each run clones the selected repositories from the default branch and pushes changes only to `claude/`-prefixed branches unless you explicitly allow unrestricted pushes.
4. Review runs at claude.ai/code/routines — each run is a full session you can open, read, and continue. Manage from the CLI with `/schedule list`, `/schedule update`, and `/schedule run`.
5. Trust but verify: a green run status means the session completed without infrastructure errors, not that the task succeeded. Read the transcript.

### Composing with Other Approaches (Intermediate)

- **Routines plus MCP connectors**: a triage routine reads new Slack reports and files Linear issues — your claude.ai connectors are available during runs (all included by default; remove the ones a routine doesn't need).
- **Routines plus headless CI**: they overlap but split cleanly — GitHub Actions for checks tied to your pipeline and secrets, routines for account-level work spanning repos and external services with no CI config at all.
- **API trigger plus Cloud Sessions**: wire your monitoring tool to the routine's `/fire` endpoint with the alert body as `text` (the endpoint requires the `anthropic-beta: experimental-cc-routine-2026-04-01` header, and the generated token is shown only once) — the run correlates the stack trace with recent commits and opens a draft PR as a full cloud session on-call can open and continue, instead of starting from a blank terminal.

### Advanced Patterns

- **Multi-trigger routines**: one PR-review routine can run nightly, react to every `pull_request.opened`, and be fired from a deploy script — same prompt, three entry points. API and GitHub triggers are added from the web UI; the CLI's `/schedule` creates schedule triggers only.
- **Filtered GitHub triggers**: scope event triggers with filters (author, title, body, base/head branch, labels, draft or merged state) — e.g. run the security-focused reviewer only when the head branch contains `auth` or a maintainer applies a `needs-backport` label.
- **Custom cron cadence**: pick the closest preset in the UI, then `/schedule update` to set an exact cron expression (minimum interval one hour).

## Common Pitfalls

- **Vague prompts**: an interactive session recovers from ambiguity by asking; a routine just guesses, every night, unattended. Spell out scope, success criteria, and what to leave alone.
- **Over-connected routines**: all your connectors are included by default, and the routine can use any of their tools — including writes — without asking. Strip a routine down to the connectors it actually needs.
- **Forgetting runs act as you**: commits, PRs, Slack messages, and ticket updates from a routine carry your identity. Anything embarrassing it does, it does under your name.
- **Set-and-forget rot**: a routine written for last quarter's repo layout degrades silently. Skim recent run transcripts periodically — the run list won't flag task-level failure on its own.

## Sources

- [Automate work with routines](https://code.claude.com/docs/en/routines) — Official docs for scheduled, API, and GitHub triggers, /schedule, connectors, and limits

## Signals

- Setup: —
- Session: `/schedule` usage; mentions routines or recurring runs
