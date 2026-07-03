# CI/CD & Automation
*Last verified: 2026-06-27*

## When You're Here

You want AI to work without you in the loop — reviewing PRs automatically, running checks in a pipeline, generating reports on a schedule, or orchestrating multi-step workflows that trigger on events. This is where AI goes from "tool I use at my desk" to "teammate that works while I sleep." The challenge is reliability: automated AI needs to handle edge cases gracefully, produce structured output, and fail loudly rather than silently.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Running Claude in GitHub Actions or GitLab CI | Headless mode | No terminal, structured JSON output, non-interactive |
| Multi-step pipeline with verification between stages | Fan-out workflows | Orchestrate steps with gates between them |
| Automated PR review on every push | Subagent delegation | Parallel reviewers (security, style, logic) per event |
| Want to enforce rules automatically on every edit | Hooks | Pre/PostToolUse hooks for formatting, linting, test runs |
| Pipeline needs data from GitHub, monitoring, or wikis | MCP context | Connect Claude to external systems at runtime |

**Hidden gem:** Custom Agents — project-specific reviewer agents catch the issues only someone who knows your codebase would notice, automatically on every PR.

## Approaches (Ranked)

### 1. Headless Mode — Run Claude in CI without a terminal
**Level:** Intermediate

Headless mode is how you run Claude in environments with no human and no terminal — GitHub Actions, GitLab CI, cron jobs, and webhook handlers. You pass a prompt via stdin or command-line argument, get structured JSON output on stdout, and use exit codes for pass/fail signaling. This is the bridge between "AI tool" and "CI step."

**Try it now:**
> Write a GitHub Actions workflow step that runs Claude in headless mode on every PR to `main`. The prompt should: read all changed files in the PR, check for TODO/FIXME/HACK comments that were added (not existing ones), and output a JSON array of `{"file": "...", "line": N, "comment": "...", "suggestion": "..."}`. If the array is non-empty, post it as a PR comment using `gh pr comment`. Use `claude -p` with `--output-format json` to get structured output.

**Why this works:** CI systems are fundamentally stdin/stdout/exit-code machines. Headless mode adapts Claude to this interface, letting you slot AI into any pipeline step that accepts a command-line tool. Structured JSON output means downstream steps can parse results programmatically.

**Pros:**
- Drops into any CI system that can run shell commands
- Structured output integrates with existing tooling and dashboards
- Non-interactive by design, so no risk of hanging on prompts

**Cons:**
- No iterative feedback loop — the prompt must be right on the first try
- Debugging failed runs requires reading CI logs, not interactive conversation

**Deeper:** See `approaches/headless-mode.md`

---

### 2. Fan-Out Workflows — Orchestrate multi-step CI pipelines with verification
**Level:** Advanced

Real CI pipelines are not linear — they fan out into parallel tracks (lint, test, build), converge for integration checks, and branch again for deployment. Fan-out workflows let you orchestrate these stages with AI-driven verification between them. Each stage gets its own focused prompt, and the orchestrator decides whether to proceed or abort based on results.

**Try it now:**
> Design a fan-out workflow for our PR validation pipeline. Stage 1 (parallel): run lint, unit tests, and type checking. Stage 2 (gate): only proceed if all three pass. Stage 3 (parallel): run integration tests and generate a coverage diff report. Stage 4 (gate): flag if coverage dropped more than 2%. Stage 5: post a summary comment on the PR with results from all stages, including timing for each. Output the workflow as a reusable script that takes a PR number as input.

**Why this works:** Fan-out workflows decompose complex pipelines into independently verifiable stages. Each stage has a clear contract (inputs, outputs, pass/fail criteria), making the overall pipeline both faster (parallelism) and more reliable (explicit gates prevent bad state from propagating).

**Pros:**
- Maximizes pipeline speed through parallelism where dependencies allow
- Explicit gates between stages make failures easy to diagnose
- Each stage can be tested and debugged independently

**Cons:**
- More complex to set up than a linear script
- Requires careful dependency mapping to know what can parallelize

**Deeper:** See `approaches/fan-out-workflows.md`

---

### 3. Subagent Delegation — Parallel reviewers triggered on PR events
**Level:** Advanced

Instead of one monolithic PR review, delegate specialized reviewers that run in parallel: one checks for security issues, another evaluates test coverage, a third reviews API contract changes, and a fourth scans for performance regressions. Each subagent has focused expertise and a narrow scope, producing higher-quality feedback than a single pass that tries to catch everything.

**Try it now:**
> Set up a PR review workflow with three parallel subagents. Agent 1: Security reviewer — check for SQL injection, XSS, hardcoded secrets, and insecure deserialization in the changed files. Agent 2: API contract reviewer — compare changed API endpoints against `openapi.yaml` and flag any undocumented breaking changes. Agent 3: Test coverage reviewer — identify changed lines that lack test coverage and suggest specific test cases. Each agent should output structured JSON with `severity`, `file`, `line`, and `message` fields. Combine all results into a single PR comment grouped by reviewer.

**Why this works:** Specialization beats generalization for review quality. A security-focused prompt with security-specific instructions catches more vulnerabilities than a general "review this PR" prompt, because the model's attention is concentrated on one concern at a time.

**Pros:**
- Higher-quality feedback through focused, specialized prompts
- Parallel execution means review time equals the slowest reviewer, not the sum
- Easy to add or remove reviewers as your needs change

**Cons:**
- Multiple agents mean higher token costs per PR
- Requires orchestration logic to aggregate and deduplicate findings

**Deeper:** See `approaches/subagent-delegation.md`

---

### 4. Hooks — Enforce rules automatically on every edit
**Level:** Intermediate

Hooks let you wire automatic actions to Claude Code's lifecycle events. A PostToolUse hook can run your linter after every file edit. A PreToolUse hook can block edits to sensitive files. A Stop hook can evaluate whether a goal was met. This turns manual discipline ("remember to run tests") into automated enforcement.

**Try it now:**
> Set up a PostToolUse hook that runs `npm test -- --related` after every file edit in `src/`, so I get instant feedback on whether my changes broke anything. Also add a PreToolUse hook that blocks any edit to files in `config/production/` unless I explicitly confirm.

**Why this works:** Automation eliminates the gap between intention and execution. Developers know they should run tests after every change, but under pressure they skip it. Hooks remove the choice — the right thing happens automatically, every time.

**Pros:**
- Zero-effort enforcement of best practices
- Catches issues at edit time, not at PR review time
- Composable — stack multiple hooks for layered protection

**Cons:**
- Slow hooks (full test suites) can break your flow
- Overly aggressive blocking hooks create friction

**Deeper:** See `approaches/hooks-as-workflow.md`

---

### 5. MCP Context — Connect to GitHub, GitLab, monitoring tools
**Level:** Intermediate

Automated workflows become significantly more powerful when Claude can reach beyond the local filesystem. MCP connections let your CI workflows pull context from GitHub (PR metadata, issue links, review comments), monitoring tools (error rates, latency dashboards), artifact stores (build outputs, test reports), and documentation systems. This turns Claude from "a tool that reads files" into "a tool that understands your entire delivery pipeline."

**Try it now:**
> In our nightly CI pipeline, I want Claude to pull the last 24 hours of error logs from our monitoring system, cross-reference them with PRs merged today via GitHub, and generate a morning report that maps new errors to the PRs that likely introduced them. The report should include error count, first occurrence timestamp, affected service, and a link to the suspected PR. Output as markdown posted to our team's Slack channel via webhook.

**Why this works:** The value of automated AI scales with the breadth of context it can access. A code review that also checks whether the change correlates with a spike in error rates catches problems that a code-only review never would.

**Pros:**
- Connects AI workflows to the full software delivery lifecycle
- Enables cross-system correlation that would take a human hours
- Context can be refreshed dynamically, keeping automated workflows current

**Cons:**
- Each MCP connection adds a dependency that can break your pipeline
- Sensitive data flowing through MCP requires careful access control

**Deeper:** See `approaches/mcp-context.md`

---

### 6. Custom Skills — Reusable skills for CI-triggered workflows
**Level:** Advanced

When your CI pipeline triggers the same Claude Code workflows repeatedly — running a PR checklist, generating release notes, verifying API contract compatibility — a custom skill encodes each workflow as a reusable command. `/pr-check` runs your team's full PR validation sequence. `/release-notes` reads commits since the last tag and produces a categorized changelog. Skills make CI workflows reproducible and easy to invoke both interactively and in headless mode.

**Try it now:**
> Create a custom skill at `.claude/skills/pr-check.md`. When invoked with `/pr-check`, it should: (1) run `npm run lint` and report any violations, (2) run `npm run typecheck` and report errors, (3) run `npm test -- --coverage` and flag if coverage dropped, (4) check for `TODO` or `FIXME` comments added in the current diff, and (5) output a structured summary with pass/fail for each check. This should work both interactively and via `claude -p "/pr-check"` in CI.

**Why this works:** CI workflows that depend on Claude need the same reproducibility as any other CI step. Custom skills provide a versioned, reviewable definition of what the workflow does, so the behavior is consistent whether triggered by a human or a pipeline.

**Pros:**
- Same workflow works interactively and in headless CI mode
- Versioned in the repo alongside the code it validates
- Easy to iterate on — edit the skill definition and re-run

**Cons:**
- Skill definitions need maintenance as CI requirements evolve
- Complex multi-step workflows may outgrow a single skill file

**Deeper:** See `approaches/custom-skills.md`

---

### 7. Custom Agents — Specialized CI review agents
**Level:** Advanced

For CI pipelines that need domain-specific review — security scanning that knows your auth patterns, API contract validation against your OpenAPI spec, database migration safety checks — custom agents in `.claude/agents/` provide focused, reusable reviewers. Each agent carries project-specific knowledge and can be invoked as a CI step, producing structured findings that integrate with PR comments.

**Try it now:**
> Create two custom agents: (1) `.claude/agents/security-reviewer.md` that checks PRs for SQL injection (we use raw queries in `src/legacy/`), auth bypass (all `/api/v1/admin/*` routes must use `requireAdminAuth`), and secrets exposure. (2) `.claude/agents/api-contract-reviewer.md` that diffs changed API endpoints against `docs/openapi.yaml` and flags undocumented breaking changes. Both should output structured JSON so CI can post findings as PR comments.

**Why this works:** Generic CI review catches generic problems. Custom agents encode the specific risks and patterns of your codebase, catching the issues that only someone who knows the project would notice — but doing it automatically on every PR.

**Pros:**
- Project-specific review knowledge applied automatically on every PR
- Structured output integrates cleanly with CI comment workflows
- New agents can be added as new review concerns emerge

**Cons:**
- Agent definitions require domain expertise to write well
- Over-specified agents can produce false positives that erode developer trust

**Deeper:** See `approaches/custom-agents.md`

---

### 8. Channels — Push CI events into your live session
**Level:** Advanced

Where headless mode runs Claude *inside* the pipeline, channels run the pipeline's events *into your session*: a webhook-receiver channel forwards CI failures, deploy statuses, or error-tracker alerts to the Claude Code session you already have open on the relevant repo. The failure gets triaged with your working context loaded — and if you're away, a two-way chat channel (Telegram, Discord, iMessage) lets you read the diagnosis and approve the fix from your phone.

**Try it now:**
> Our CI posts webhook notifications on failed builds. Build a small channel server following the channels reference (capability declaration, notification event, reply tool) that accepts those webhooks and forwards them into my session, so a red build on `main` shows up here with the job name and failing step while I still have the branch context loaded.

**Why this works:** Reacting to CI needs the context that produced the change; pushing the event to the session that made the change beats cold-starting an agent that has to rediscover everything.

**Pros:**
- CI failures triaged with the authoring session's context intact
- Complements headless mode instead of replacing it — push events in, or run Claude in the pipe
- Custom channel servers cover any system that can send a webhook

**Cons:**
- Research preview with per-session opt-in (`--channels`) — not fire-and-forget infrastructure yet
- Building a custom webhook channel is a small development project of its own

**Deeper:** See `approaches/channels.md`
