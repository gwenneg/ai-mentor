# Code Review
*Last reviewed: 2026-06-27*

## When You're Here

You have code that needs a second pair of eyes — maybe your own PR before requesting human review, maybe a teammate's PR that touches unfamiliar subsystems, or maybe a critical change where you need confidence that nothing slipped through. AI doesn't replace human judgment on design decisions, but it excels at the systematic, exhaustive checks that humans rush through on the third PR of the day.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Quick review of a focused diff | Built-in review skills | One command, immediate feedback |
| Reviewing my own code before opening a PR | Built-in review skills | One command to catch issues before requesting human review |
| Large PR touching security, perf, and correctness | Subagent delegation | Parallel specialized reviewers catch more than one pass |
| Critical change needing adversarial verification | Fan-out workflows | Cross-checks findings against each other |
| PR implements a design doc or addresses an issue | MCP context | Review against the spec, not just the code |

## Approaches (Ranked)

### 1. Built-In Review Skills — One command, immediate results
**Level:** Beginner | **Tools:** Claude Code

For everyday reviews, built-in skills provide immediate value with zero setup. `/code-review` scans for correctness bugs and cleanup opportunities. `/security-review` focuses on vulnerabilities. `/simplify` finds reuse and efficiency improvements. These are your daily drivers — fast, reliable, and tuned for common issues.

**Try it now:**
> Run `/code-review --effort high` on the current branch diff. Focus on the payment processing changes in `src/services/billing/` — I'm most concerned about edge cases in the discount stacking logic and whether the error handling covers all the Stripe webhook failure modes.

**Why this works:** Most code review value comes from systematic checking, not creative insight. Built-in skills codify the checks that experienced reviewers do instinctively — null checks, error handling, boundary conditions — and apply them consistently to every line.

**Pros:**
- Zero setup, works immediately
- Effort levels let you match depth to risk
- Consistent — doesn't get tired or skip files

**Cons:**
- Less customizable than subagent approaches
- May miss domain-specific business logic issues

**Deeper:** See `approaches/built-in-review-skills.md`

---

### 2. Subagent Delegation — Specialist reviewers working in parallel
**Level:** Advanced | **Tools:** Claude Code

Instead of one reviewer trying to hold security, performance, correctness, and style in their head simultaneously, subagent delegation spawns focused reviewers. One agent checks for SQL injection and auth bypasses. Another profiles algorithmic complexity. A third verifies business logic. They work in parallel and report back, giving you multi-dimensional coverage in the time it takes for a single pass.

**Try it now:**
> Review PR #247 using three parallel subagents: (1) a security reviewer checking for injection, auth bypass, and secrets exposure in `src/api/`, (2) a performance reviewer looking for N+1 queries and unnecessary allocations in `src/services/`, and (3) a correctness reviewer verifying that the new discount logic in `OrderService.applyPromoCode()` handles edge cases like expired codes, stacking, and negative totals. Consolidate findings by severity.

**Why this works:** Attention is a finite resource. A single reviewer doing multiple types of analysis will miss things because they're context-switching between concerns. Specialized agents maintain focus, and parallelism means the wall-clock time stays low.

**Pros:**
- Catches issues across multiple dimensions simultaneously
- Each reviewer stays focused on its specialty
- Scales to large PRs without reviewer fatigue

**Cons:**
- Requires prompt engineering to define each reviewer's scope
- Can produce redundant or conflicting findings that need reconciliation

**Deeper:** See `approaches/subagent-delegation.md`

---

### 3. Fan-Out Workflows — Review with built-in adversarial checks
**Level:** Advanced | **Tools:** Claude Code

Fan-out workflows go beyond parallel review by adding adversarial verification. After initial reviewers flag issues, a second wave checks whether those findings are real or false positives. This is critical for code review, where false alarms erode trust and waste developer time.

**Try it now:**
> Review the diff on this branch using a fan-out workflow. First wave: have three reviewers independently analyze the changes for bugs, security issues, and performance regressions. Second wave: have an adversarial verifier check each finding — confirm it's a real issue by tracing the code path, and discard false positives. Final output: only verified findings, ranked by severity.

**Why this works:** The biggest problem with automated code review isn't missing issues — it's false positives. Adversarial verification mimics the natural pushback that happens in human review conversations ("Are you sure that's a bug? I think that's handled by..."), producing higher-signal results.

**Pros:**
- Dramatically reduces false positives
- Findings come pre-verified with reasoning
- Mimics the back-and-forth of good human review

**Cons:**
- More expensive in tokens and time than a single pass
- Overkill for small, low-risk changes

**Deeper:** See `approaches/fan-out-workflows.md`

---

### 4. MCP Context — Review against the spec, not just the code
**Level:** Intermediate | **Tools:** Claude Code / OpenCode (with MCP)

Code that is internally correct can still be wrong if it doesn't match the requirements. MCP context connects your review to issue trackers, design docs, and architecture decision records, so Claude can verify that the implementation actually addresses the ticket, follows the agreed-upon design, and doesn't contradict documented constraints.

**Try it now:**
> Connect to our Linear workspace and pull the requirements from issue ENG-1842. Then review the changes on this branch against those requirements. Flag any acceptance criteria that aren't covered by the implementation, and any implementation details that contradict the design notes in the issue comments.

**Why this works:** The most expensive bugs are specification bugs — code that works perfectly but solves the wrong problem. By grounding review in the original requirements, you catch these before they reach production.

**Pros:**
- Catches gaps between spec and implementation
- Brings context that code alone can't provide
- Reduces "I thought the requirement was..." conversations

**Cons:**
- Requires MCP server setup for your tools
- Only as good as the documentation it has access to

**Deeper:** See `approaches/mcp-context.md`

---

### 5. Custom Agents — Reusable specialist reviewers in your repo
**Level:** Advanced | **Tools:** Claude Code

Instead of writing detailed reviewer prompts every time, define reusable reviewer agents in `.claude/agents/`. A `security-reviewer.md` agent knows your auth patterns, sensitive data fields, and common vulnerability classes. A `perf-reviewer.md` agent knows your N+1 hotspots and caching strategy. Invoke them by name and they bring their full expertise to every review, consistently.

**Try it now:**
> Create a custom agent at `.claude/agents/security-reviewer.md` for our project. It should check for: SQL injection (we use Prisma, so focus on raw queries), auth bypass (all routes in `src/api/` must use `requireAuth` middleware), secrets exposure (no hardcoded API keys or connection strings), and insecure deserialization of user input. Then run it against the current branch diff.

**Why this works:** Review quality depends on the reviewer knowing the codebase's specific patterns and risks. Custom agents encode that knowledge once and apply it consistently — no more hoping the reviewer remembers to check for raw SQL queries in a codebase that normally uses an ORM.

**Pros:**
- Consistent, project-aware reviews on every PR
- New reviewers can be added for new concerns (accessibility, API contracts)
- Team shares the same review standards without tribal knowledge

**Cons:**
- Agent definitions need maintenance as codebase patterns evolve
- Over-specified agents can produce false positives on intentional exceptions

**Deeper:** See `approaches/custom-agents.md`

---

### 6. Plugins — Ready-made review workflows out of the box
**Level:** Intermediate | **Tools:** Claude Code

If you want structured code review without defining your own agents and prompts, plugins like `code-review` and `pr-review-toolkit` provide pre-built review workflows. They handle the orchestration — parallel reviewers, finding consolidation, severity ranking — so you can focus on acting on the results rather than building the pipeline.

**Try it now:**
> Browse available Claude Code plugins that provide code review capabilities. Install the most relevant one for our project and run it against the latest changes on this branch. I want structured review output with severity levels that I can share with the team.

**Why this works:** Review plugins encode best practices from the community — structured output, severity ranking, false-positive filtering — that take significant effort to build from scratch. For most teams, a good plugin covers 80% of review needs immediately.

**Pros:**
- Structured review workflow with zero setup effort
- Community-maintained with evolving best practices
- Standardized output format that integrates with PR comments

**Cons:**
- Less customizable than hand-crafted agent-based review pipelines
- May not cover domain-specific review concerns unique to your project

**Deeper:** See `approaches/plugins.md`
