# Tech Debt Assessment
*Last reviewed: 2026-07-02*

## When You're Here

You suspect the codebase has accumulated significant tech debt, but you need more than a gut feeling. Maybe you are preparing a proposal for leadership to justify a cleanup sprint, triaging debt before the next quarter's planning, or inheriting a codebase and trying to understand where the landmines are buried. The core challenge is not fixing the debt — it is finding it, measuring it, and deciding what to fix first.

Tech debt assessment is distinct from refactoring. Refactoring assumes you already know what to change. Assessment is the step before that: auditing the codebase across multiple dimensions (duplication, complexity, coverage gaps, deprecated patterns), quantifying the cost of inaction, and producing a prioritized list so your limited cleanup budget goes where it matters most.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Multi-dimensional audit across the whole codebase | Subagent delegation | Parallel auditors cover more ground faster |
| Quick sense of the worst offenders, no setup | Built-in review skills | One command surfaces obvious debt immediately |
| You have findings and need to prioritize them | Plan mode | Structures debt by impact and urgency |
| Very large codebase with many modules | Fan-out workflows | Scales auditing across module boundaries |
| Recurring patterns specific to your project | Custom agents | Codifies your team's known debt patterns |

**Hidden gem:** Custom Agents — encoding your team's own deprecated patterns as a detector turns tribal knowledge into a tracked migration metric.

## Approaches (Ranked)

### 1. Subagent Delegation — Parallel auditors by concern
**Level:** Advanced

Tech debt is multi-dimensional: duplication, deprecated APIs, missing tests, excessive complexity, inconsistent error handling. A single pass trying to track all of these at once produces shallow results. Subagent delegation spawns focused auditors — one per concern — that work in parallel and report back. The consolidated report gives you a multi-dimensional debt picture in a fraction of the time a sequential audit would take.

**Try it now:**
> Spawn four parallel agents to audit this codebase for tech debt. Agent 1: scan `src/` for code duplication — find functions or blocks that are near-identical and report clusters with file paths and line numbers. Agent 2: search for deprecated API usage — look for `TODO`, `FIXME`, `DEPRECATED`, `@deprecated` annotations and calls to known deprecated standard library methods. Agent 3: measure test coverage gaps — identify modules in `src/` that have no corresponding test files in `tests/` and functions over 20 lines with zero test references. Agent 4: find overly complex functions — report any function with cyclomatic complexity above 10 or more than 50 lines. Consolidate all findings into a single prioritized report.

**Why this works:** Each debt dimension requires a different scanning strategy. Duplication detection is pattern matching; deprecated API scanning is keyword search; coverage analysis is cross-referencing source and test directories. Specialized agents apply the right strategy per concern without conflating them.

**Pros:**
- Covers multiple debt dimensions simultaneously
- Each auditor stays focused, producing higher-quality findings
- Consolidated report reveals cross-cutting patterns (e.g., the most-duplicated module also has the worst coverage)

**Cons:**
- Requires clear scope definitions to avoid overlapping or redundant findings
- Consolidation step needs careful prompt design to avoid losing nuance

**Deeper:** See `approaches/subagent-delegation.md`

---

### 2. Built-In Review Skills — Quick health check with zero setup
**Level:** Beginner

Before building a sophisticated audit pipeline, start with what is already available. `/code-review` surfaces dead code, redundant logic, missing error handling, and overly complex functions. `/simplify` identifies code that could be consolidated or replaced with simpler constructs. Running these across the codebase gives you a fast baseline of obvious debt without any setup.

**Try it now:**
> Run `/code-review --effort high` across the project. I am assessing tech debt and want to identify the most problematic areas: dead code, functions that are too complex, missing error handling, and duplicated logic. Focus especially on `src/services/` and `src/utils/` — those directories have not been touched in over a year.

**Why this works:** The most impactful tech debt is often the most visible — functions that are 300 lines long, try/catch blocks that swallow errors silently, utility functions duplicated across three modules. Built-in skills catch these systematically without requiring you to know what to look for.

**Pros:**
- Zero setup, immediate results
- Good starting point before investing in more sophisticated approaches
- Catches the low-hanging fruit that often represents the highest-impact debt

**Cons:**
- Limited to patterns the built-in skills recognize
- Cannot detect project-specific debt like deprecated internal APIs or legacy patterns

**Deeper:** See `approaches/built-in-review-skills.md`

---

### 3. Plan Mode — Prioritize what to fix first
**Level:** Beginner

Identifying debt is only half the problem. The harder question is what to fix first. Plan mode takes your audit findings and helps you prioritize by impact: which debt causes the most bugs, which slows the team down most, which blocks upcoming features. The output is a prioritized backlog you can take to sprint planning or include in a tech debt proposal for leadership.

**Try it now:**
> Enter plan mode. Here is our tech debt audit summary: (1) `src/services/payment_gateway.py` has three near-identical implementations of retry logic, (2) `src/utils/date_helpers.ts` uses the deprecated `moment.js` library, (3) `src/controllers/` has 14 controller files with no corresponding test files, (4) `src/models/user.rb` is 1200 lines with a cyclomatic complexity of 47. Prioritize these by impact: which causes the most bugs, which slows developers down most, which blocks the planned API v2 migration. Produce a prioritized backlog with estimated effort for each item.

**Why this works:** Debt prioritization requires weighing multiple factors — frequency of bugs caused, developer time wasted, strategic alignment with roadmap — that are hard to hold in your head simultaneously. Plan mode structures this analysis and forces explicit trade-off reasoning.

**Pros:**
- Turns a wall of findings into an actionable, ordered backlog
- Forces explicit reasoning about impact versus effort
- Output is ready to share with stakeholders

**Cons:**
- Quality depends on the accuracy of your audit inputs
- Cannot measure actual bug frequency or developer time — relies on estimation

**Deeper:** See `approaches/plan-mode.md`

---

### 4. Fan-Out Workflows — Audit the entire codebase at scale
**Level:** Advanced

For large codebases with many modules, a single audit pass — even with subagents — may not cover enough ground. Fan-out workflows distribute auditors across modules in parallel: each agent scans one module and reports standardized debt metrics. The orchestrator consolidates findings, identifies cross-cutting patterns, and highlights modules that are disproportionately problematic.

**Try it now:**
> Fan out a tech debt audit across our microservices. For each directory in `services/` (`auth/`, `billing/`, `inventory/`, `notifications/`, `search/`), spawn an agent that reports: (1) number of `TODO`/`FIXME` annotations, (2) files with no test coverage, (3) functions over 40 lines, (4) any usage of the deprecated `v1` internal API client. Consolidate into a table ranking modules by overall debt load, and flag any cross-cutting patterns that appear in three or more modules.

**Why this works:** Large codebases have uneven debt distribution — typically a few modules carry most of the burden. Fan-out auditing reveals this distribution quantitatively, so you can focus cleanup effort on the modules that matter most instead of spreading it thinly across the entire codebase.

**Pros:**
- Scales to codebases with dozens of modules or services
- Standardized metrics make cross-module comparison meaningful
- Orchestrator catches patterns that per-module auditors miss individually

**Cons:**
- High token cost for very large codebases
- Standardized metrics may miss module-specific debt that does not fit the template

**Deeper:** See `approaches/fan-out-workflows.md`

---

### 5. Custom Agents — Project-specific debt detectors
**Level:** Advanced

Every codebase has its own unique debt patterns: the deprecated internal API your team is migrating away from, the legacy ORM pattern you have agreed to replace, the module that should have been split two years ago. Custom agents encode these patterns into reusable detectors that run consistently, catching debt that generic tools would never flag.

**Try it now:**
> Create a custom agent at `.claude/agents/debt-detector.md` for our project. It should detect: (1) any usage of `LegacyHttpClient` — we are migrating to `HttpService` and need to track remaining call sites, (2) models in `src/models/` that directly import `DatabaseConnection` instead of going through the repository layer, (3) any controller in `src/controllers/` that contains SQL strings (all queries should be in `src/repositories/`), (4) React class components in `src/components/` that should be converted to function components. Then run it against the codebase and report findings grouped by category.

**Why this works:** Generic debt detection tools flag universal code smells, but the most impactful debt in a mature codebase is project-specific: patterns your team has explicitly decided to move away from. Custom agents encode those decisions, ensuring they are enforced consistently and that progress is tracked over time.

**Pros:**
- Catches debt that only your team would recognize
- Runs consistently — no reliance on tribal knowledge or manual audits
- Findings serve as a migration progress tracker over time

**Cons:**
- Requires upfront investment to define your project's specific debt patterns
- Agent definitions need maintenance as migration targets evolve

**Deeper:** See `approaches/custom-agents.md`

---

### 6. Background Agents — Run the audit while you do something else
**Level:** Intermediate

A thorough debt audit takes an agent an hour and takes you thirty seconds to specify — the definition of work that shouldn't hold your terminal hostage. Dispatch the audit as a background session from `claude agents`; it works in its own isolated worktree, and you review a finished report (or draft PR of mechanical cleanups) instead of babysitting the scan.

**Try it now:**
> Run `claude agents` and dispatch: "Audit src/ for tech debt: duplicated logic, functions over 50 lines, TODO/FIXME annotations, and modules without test files. Produce a prioritized report in docs/debt-audit.md with file paths and line counts per finding."

**Why this works:** Audits are attention-light but time-heavy — backgrounding converts them from "a thing you do" into "a thing you review," which is the only way recurring debt assessment survives contact with a busy sprint.

**Pros:**
- Your terminal and your attention stay free during the scan
- Automatic worktree isolation keeps audit artifacts off your working copy
- Notification fires when it finishes or needs a judgment call

**Cons:**
- Underspecified audit criteria produce noisy reports — invest in the dispatch prompt
- Results land on a worktree branch, not your working directory

**Deeper:** See `approaches/background-agents.md`
