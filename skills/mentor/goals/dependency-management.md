# Dependency Management
*Last reviewed: 2026-07-02*

## When You're Here

You're responsible for the health of your dependency tree. Maybe you need to evaluate a new library before the team adopts it, audit existing dependencies for vulnerabilities or maintenance red flags, replace something that's been deprecated, or just make sense of the supply chain risk hiding in your `package.json` or `requirements.txt`. Dependency decisions compound — a careless adoption today becomes a painful migration next year.

This is different from migration (which is about upgrading a specific framework version) and research (which is general investigation). Dependency management is the ongoing discipline of knowing what you depend on, whether it's healthy, and what to do when it isn't.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Evaluating a library you haven't used before | Deep research | Surfaces maintenance health, CVEs, and community sentiment before you commit |
| Need to understand what depends on a package before removing it | Plan mode | Maps the blast radius so you know what breaks |
| Your org has approved dependency lists or security policies | MCP context | Check internal standards before researching externally |
| Bumping a dependency version with cascading breakage | Autonomous loops | Set "tests pass" and let Claude grind through the fixes |
| Want to test a major upgrade without risking your branch | Worktree isolation | Try it in a disposable copy and evaluate the damage |

**Hidden gem:** Worktree Isolation — trying the upgrade in a throwaway copy gives you a real damage report before you commit to anything.

## Approaches (Ranked)

### 1. Deep Research — Evaluate before you adopt
**Level:** Beginner

The most common dependency mistake is adopting a library based on one blog post or a compelling README. Deep research fans out across npm/PyPI download stats, GitHub activity (commit frequency, open issue count, response times), license terms, CVE databases, and community sentiment on Reddit and Hacker News. It gives you the full picture — not just whether a library does what you want, but whether it will still be maintained in two years.

**Try it now:**
> /deep-research We're considering adopting `zod` to replace `joi` for runtime validation in our Node.js API (`src/validators/`). Compare them on: TypeScript integration, bundle size, runtime performance, maintenance activity, breaking change history, and community adoption trends. Also check if there are any known CVEs for either library. We need something that will be actively maintained for at least 3 years.

**Why this works:** Dependency adoption is a long-term bet. Deep research automates the due diligence that most developers skip — checking issue response times, release cadence, and whether the maintainer is a single person or a funded team. This is the difference between a dependency that serves you for years and one that becomes abandonware.

**Pros:**
- Surfaces maintenance health signals that README stars don't show
- Cross-references multiple sources to catch bias from a single blog post
- Checks CVE databases and license compatibility you'd otherwise miss

**Cons:**
- Quality depends on publicly available data — internal or niche libraries get thin coverage
- Can overwhelm with information if you don't constrain the comparison criteria

**Deeper:** See `approaches/deep-research.md`

---

### 2. Plan Mode — Map your dependency graph before changing it
**Level:** Beginner

Before adding, removing, or replacing a dependency, you need to understand the blast radius. Plan mode walks your import graph to find every file that uses a package, identifies transitive dependents, and maps out what breaks if you remove or swap it. This prevents the all-too-common experience of yanking a library and discovering that three other packages depended on it transitively.

**Try it now:**
> Enter plan mode. I want to remove `moment.js` from our project and replace it with `date-fns`. First, find every import of `moment` across `src/` and `lib/`. Then identify which of those usages are for formatting, which are for date math, and which are for timezone handling. Map the `date-fns` equivalent for each usage pattern and flag any that don't have a direct replacement. Give me a replacement plan ordered by complexity.

**Why this works:** Dependency changes are graph operations, not file operations. Removing a package affects everything downstream. Plan mode reveals the full dependency graph before you start cutting edges, so you can sequence the work and catch transitive breakage before it happens.

**Pros:**
- Reveals transitive dependencies and hidden usages you'd miss with grep
- Creates a concrete replacement plan you can execute incrementally
- Prevents the "I'll just swap it" approach that leaves the codebase half-migrated

**Cons:**
- Takes time before any code changes happen — but prevents costly rework

**Deeper:** See `approaches/plan-mode.md`

---

### 3. MCP Context — Check internal standards and approved lists
**Level:** Intermediate

Many organizations maintain approved dependency lists, security policies, or architecture decision records about specific libraries. Before spending time researching a library externally, check whether your team has already evaluated it, banned it, or mandated an alternative. MCP context servers can pull from Confluence, internal wikis, Slack discussions, and artifact repositories to surface this institutional knowledge.

**Try it now:**
> Check our internal Confluence space "Engineering Standards" and the `approved-dependencies.yaml` file in our `platform-config` repo for any policies on HTTP client libraries. We're considering switching from `requests` to `httpx` in our Python services. Has any team already evaluated this? Are there security or compliance requirements that constrain our choice?

**Why this works:** The fastest dependency decision is one someone already made. Checking internal sources first prevents you from duplicating evaluation work another team already completed — and avoids the embarrassment of recommending a library that security already rejected.

**Pros:**
- Prevents duplicated evaluation effort across teams
- Surfaces compliance and security constraints early
- Incorporates lessons learned from past adoption decisions

**Cons:**
- Requires MCP server setup for internal tools (Confluence, Jira, etc.)
- Internal documentation may be outdated or incomplete

**Deeper:** See `approaches/mcp-context.md`

---

### 4. Autonomous Loops — Upgrade and fix until green
**Level:** Intermediate

For mechanical dependency updates — bump the version in your lockfile, fix the breaking API changes, update the import paths, get the tests passing — autonomous loops excel. You set the success condition ("build passes and all tests green") and Claude iterates through the cascade of changes: fix one type error, run the build, find the next error, fix it, repeat. This is especially effective for minor version bumps with a handful of breaking changes scattered across many files.

**Try it now:**
> Upgrade `typescript` from 5.3 to 5.5 in `package.json`. Run `npm install`, then `npm run build`. For each type error, fix it according to the TypeScript 5.5 migration notes. After each round of fixes, run `npm run build && npm test` again. Keep iterating until both the build and all tests pass with zero errors. Don't change any test expectations unless the old behavior was a TypeScript bug.

**Why this works:** Dependency upgrades are convergent problems — each fix brings you closer to green, and the test suite tells you exactly how far you have left. AI excels at these tight edit-compile-fix loops because it doesn't lose focus on iteration 30.

**Pros:**
- Handles tedious cascading fixes without supervision
- Self-verifies by running the build and tests after each round
- Catches transitive type errors that ripple across the codebase

**Cons:**
- Can make "creative" fixes to pass the build that aren't idiomatic — review the diff
- Not suitable when passing tests doesn't fully validate correctness

**Deeper:** See `approaches/autonomous-loops.md`

---

### 5. Worktree Isolation — Test an upgrade without risking main
**Level:** Intermediate

Before committing to a dependency upgrade, try it in a throwaway environment. Worktree isolation creates a disposable copy of your repo where you can bump the version, run the full test suite, and evaluate the damage — all without touching your working branch. If the upgrade is worse than expected, delete the worktree and lose nothing. If it's clean, you have a working branch ready to merge.

**Try it now:**
> Create a worktree from main. In it, upgrade `react` and `react-dom` from 18.2 to 19.0 in `package.json`. Run `npm install` and then `npm run build && npm test`. Capture every deprecation warning, type error, and test failure. Summarize the total damage: how many files need changes, what categories of breakage exist, and give an honest estimate of whether this upgrade is a half-day or a half-sprint.

**Why this works:** The psychological safety of a disposable environment changes how you evaluate risk. You'll try the upgrade you've been putting off when failure costs nothing, and you'll get a realistic damage assessment instead of guessing from the changelog.

**Pros:**
- Zero risk to your current work — the worktree is fully disposable
- Gives you a concrete damage report instead of changelog speculation
- Great for answering "should we upgrade now or wait?"

**Cons:**
- Adds disk space for the worktree copy
- Only evaluates build-time and test-time breakage, not runtime behavior

**Deeper:** See `approaches/worktree-isolation.md`

---

### 6. Scheduled & Recurring Agents — Nightly dependency triage that actually happens
**Level:** Intermediate

Dependency hygiene fails because it's important-but-never-urgent. A scheduled routine runs the triage nightly on cloud infrastructure: read the open dependency-update PRs, check changelogs, run the tests, label the safe ones, and summarize the breaking ones. You review pre-analyzed PRs over coffee instead of remembering to do the sweep.

**Try it now:**
> /schedule every weekday at 6am: check open Renovate PRs in this repo. For each, read the changelog diff and run the test suite. Label passing patch/minor bumps `auto-verified`; for major bumps or failures, comment with a summary of breaking changes and failing tests. Never merge anything.

**Why this works:** Recurring maintenance survives only when it stops depending on human initiative — a schedule converts "someone should" into "it happened, review it."

**Pros:**
- Runs on Anthropic cloud infrastructure — your machine stays closed
- Triggers can combine: schedule + GitHub events + API calls from your tooling
- Each run is a reviewable session with a full transcript

**Cons:**
- Runs autonomously with no prompts — the routine prompt must be fully self-contained
- Actions carry your identity (commits, comments appear as you)

**Deeper:** See `approaches/scheduled-agents.md`
