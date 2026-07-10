# Built-in Review Skills
*Last verified: 2026-07-06*

## What It Is

Built-in Review Skills are ready-to-use commands in Claude Code that analyze your code changes for bugs, security issues, and simplification opportunities. You run a slash command — `/code-review`, `/security-review`, or `/simplify` — and Claude systematically reviews your diff using structured analysis.

Two companion skills close the loop from the behavior side: `/verify` exercises a change end-to-end in the running application to confirm it actually does what it's supposed to (driving the affected flow, not just the tests), and `/run` launches the project's app so you can see a change working for real.

## Why It Works

A reviewer following a defined methodology catches more than one who "just reads through the diff" — and separating correctness, security, and simplification into distinct passes lets each go deep on its category.

## When to Use It

- Self-reviewing your own changes before pushing — catch bugs your eyes skipped after hours of writing the code
- Reviewing a teammate's PR when you want a structured first pass before doing your own read-through
- Post-refactoring cleanup — run `/simplify` after a large restructuring to catch missed deduplication
- Confirming a nontrivial change actually works before committing — `/verify` drives the affected flow end-to-end instead of trusting typecheck and unit tests alone
- CI integration — run reviews automatically on every PR via GitHub Actions

## When NOT to Use It

- Reviewing changes you have not made yet — these skills analyze the current diff, not hypothetical code
- As a substitute for human review on critical paths — use them as a first pass, not the only pass
- Running `/verify` on diffs with no runtime surface (docs-only or test-only changes) — there is no behavior to observe

## How It Works

### Basic (Beginner)

1. Make your code changes and stage them (or leave them unstaged — both work)
2. Run `/code-review` in your Claude Code session — optionally pass an effort level and a target such as a path or PR number (`/code-review high src/api/`). Effort levels:
   - Low/Medium: fewer findings, higher confidence — good for quick sanity checks
   - High: broader coverage, may surface uncertain findings — good for thorough review
   - xhigh/max: the deepest local levels — exhaustive coverage for critical changes (available levels depend on the model)
3. Claude analyzes the diff and reports findings grouped by severity
4. Optionally, add `--fix` to have Claude auto-apply its findings: `/code-review --fix`
5. Or add `--comment` to post findings as inline PR comments: `/code-review --comment`

### Composing with Other Approaches (Intermediate)

- **Review then verify**: `/code-review` reads the code; `/verify` observes the behavior. Static review catches bugs that would never reproduce in a quick manual test, while verification catches integration failures no diff reader can see. Together they cover failure modes neither catches alone.
- **Subagent changes then review**: After spawning subagents to make parallel changes, run `/code-review` on each worktree branch to verify the agents' work before merging.
- **Plan Mode then review**: Use Plan Mode to design and execute a change, then immediately run `/security-review` to catch security implications the plan did not consider.

### Advanced Patterns

- **CI pipeline integration**: Run reviews on every PR without a manual step: use the `claude-code-action` GitHub Action with a review skill as the prompt, so findings post as inline comments alongside the diff. From any CI script, `claude ultrareview` runs the deep cloud review non-interactively.
- **Targeted review with context**: Before running the review, tell Claude about specific concerns: "This change modifies our rate limiter. Run /security-review with extra attention to bypass vectors." The skill uses your context to focus its analysis.
- **Ultrareview for pre-merge confidence** (research preview, v2.1.86+): `/code-review ultra` launches a fleet of reviewer agents in a cloud sandbox — every finding is independently reproduced and verified, so results skew toward real bugs rather than style notes. Reviews your branch diff or a PR (`/code-review ultra 1234`), takes ~5-10 minutes in the background, and bills to usage credits — see the ultrareview docs for current pricing. Requires claude.ai auth; unavailable on Bedrock, Google Cloud's Agent Platform, and Microsoft Foundry, and to Zero Data Retention orgs. From CI, `claude ultrareview` runs the same review non-interactively.

## Common Pitfalls

- **Running only `/simplify` without `/code-review` first**: Simplification can mask bugs. If you have a logic error in duplicated code, `/simplify` might deduplicate it — now the bug is in one place but still exists. Review for correctness first.
- **Treating AI review as authoritative**: These skills catch real bugs, but they also produce false positives. A finding that says "possible null pointer" might be guarded by an upstream check the AI did not trace. Always evaluate findings with your domain knowledge.
- **Ignoring effort levels**: Running high-effort review on every trivial change wastes time. Use low/medium for routine work, save high effort for complex or risky changes.
- **Not using `--comment` in team workflows**: If you review locally but do not post comments, your teammates do not benefit. Use `--comment` to make findings visible on the PR.

## Sources

- [Claude Code Skills](https://code.claude.com/docs/en/skills) — Official docs for skills including built-in /code-review and /security-review
- [Commands](https://code.claude.com/docs/en/commands) — Reference for /code-review effort levels and --fix/--comment flags, plus /simplify, /security-review, /verify, and /run
- [Find bugs with ultrareview](https://code.claude.com/docs/en/ultrareview) — Official docs for /code-review ultra: cloud fleet review, pricing, and the CI subcommand

## Signals

- Setup: Review commands wired into CI workflows
- Session: `/code-review`, `/security-review`, `/simplify`, `/verify`, or `/run` in the transcript
