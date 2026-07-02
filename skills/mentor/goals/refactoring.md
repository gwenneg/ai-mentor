# Refactoring & Large-Scale Changes
*Last reviewed: 2026-07-02*

## When You're Here

You need to restructure code without changing behavior — maybe extracting a module that's grown too large, migrating from one API to another across dozens of files, or paying down tech debt that's been on the backlog for three sprints. Refactoring is high-reward but high-risk: every file you touch is a chance to break something. AI shines here because it can hold more context, work more files in parallel, and verify continuously.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Complex refactor that could go sideways | Plan mode | Design the strategy before touching code |
| Changes spanning many files independently | Subagent delegation | Parallelize across files safely |
| Refactor with comprehensive test coverage | Autonomous loops | Iterate until all tests pass again |
| Multiple agents modifying overlapping areas | Worktree isolation | Prevents file conflicts between agents |
| Post-refactor cleanup and polish | Built-in review skills | /simplify catches what you missed |
| Risky refactor you might need to undo | Checkpoints & rewind | Safety net for bold changes |

**Hidden gem:** Checkpoints & Rewind — knowing any restructuring is instantly reversible changes which refactors you dare to attempt.

## Approaches (Ranked)

### 1. Plan Mode — Architect the change before writing a line
**Level:** Beginner

Large refactors fail when developers start changing code before understanding the full dependency graph. Plan mode forces you to map out what moves where, what interfaces change, what breaks, and in what order to make changes so the codebase stays buildable at each step. This is the difference between a clean refactor and a three-day merge conflict nightmare.

**Try it now:**
> Enter plan mode. I need to extract our authentication logic from the monolithic `UserController` into a dedicated `AuthService`. The controller is in `src/controllers/user_controller.rb` (~800 lines) and handles login, signup, password reset, OAuth callbacks, session management, and user profile CRUD. Plan the extraction: which methods move, what new interfaces are needed, which existing callers need updating, and what order to make changes so tests pass at each step.

**Why this works:** Refactoring is fundamentally a dependency management problem. Planning reveals the dependency graph before you start cutting, so you know exactly what breaks and in what order to fix it. Without this, you discover dependencies one crash at a time.

**Pros:**
- Prevents the "I didn't realize X depended on Y" surprises
- Creates a step-by-step roadmap you can follow or delegate
- The plan itself becomes documentation of the change

**Cons:**
- Feels like overhead when you "just want to start coding"
- Plan may need revision as you discover hidden dependencies

**Deeper:** See `approaches/plan-mode.md`

---

### 2. Subagent Delegation — Parallelize across files
**Level:** Advanced

When your refactoring plan identifies independent work streams (e.g., "update all API controllers to use the new service layer" while "migrate database queries to the new ORM syntax"), subagents can execute these in parallel. Each agent gets a clear scope, works independently, and reports back. This turns a day-long refactor into an hour-long one.

**Try it now:**
> I'm migrating from `moment.js` to `dayjs` across the frontend. Spawn three parallel agents: (1) update all date formatting calls in `src/components/` — replace `moment().format()` with `dayjs().format()` and adjust format strings, (2) update all duration calculations in `src/utils/` — replace `moment.duration()` with dayjs duration plugin equivalents, (3) update all relative time displays in `src/views/` — replace `moment().fromNow()` with dayjs relativeTime plugin. Each agent should run `npm test` in its scope after making changes.

**Why this works:** Most large refactors decompose into independent subtasks that touch different files. Parallelizing these is safe because there are no conflicts, and the speedup is roughly linear with the number of agents. This is where AI-assisted refactoring dramatically outperforms manual work.

**Pros:**
- Linear speedup for independent file changes
- Each agent has focused context, reducing mistakes
- Can process hundreds of files in minutes

**Cons:**
- Requires clear task boundaries — overlapping scopes cause conflicts
- Needs coordination if changes create cross-file dependencies

**Deeper:** See `approaches/subagent-delegation.md`

---

### 3. Autonomous Loops — Refactor until green
**Level:** Intermediate

After making structural changes, there's often a cascade of small fixes needed — updated imports, renamed references, adjusted type signatures. Autonomous loops handle this grind: set "all tests pass" as the goal and let Claude chase down every broken reference until the suite is green again.

**Try it now:**
> /goal: I've moved `src/utils/helpers.ts` to `src/shared/helpers/index.ts` and updated the main exports. Now there are broken imports across the project. Find every file that imports from the old path, update it to the new path, and verify by running `npm run typecheck && npm test`. Keep going until there are zero errors.

**Why this works:** Refactoring creates a known-good end state (tests pass, types check) with an unknown number of steps to get there. Autonomous loops excel at this pattern because they can iterate without human intervention, handling the tedious mechanical fixes that are easy individually but numerous collectively.

**Pros:**
- Handles cascading fixes without manual intervention
- Self-verifying — runs tests after each change
- Ideal for the "mechanical" phase of refactoring

**Cons:**
- May make minimal fixes that pass tests but aren't idiomatic
- Can loop excessively on stubborn type errors — set iteration limits

**Deeper:** See `approaches/autonomous-loops.md`

---

### 4. Worktree Isolation — Each agent gets its own sandbox
**Level:** Intermediate

When multiple agents are refactoring simultaneously, file conflicts become a real risk. Worktree isolation gives each agent its own copy of the repo. They can make changes freely without stepping on each other, and you merge the results when they're done. This is essential for large-scale refactors where subagents touch adjacent areas.

**Try it now:**
> Create two isolated worktrees for a parallel refactor. In worktree A: refactor the `PaymentGateway` class to use the Strategy pattern, extracting `StripeGateway`, `PayPalGateway`, and `SquareGateway` from `src/payments/gateway.py`. In worktree B: refactor the `NotificationService` to use the Observer pattern, extracting channel-specific handlers from `src/notifications/service.py`. Both should run their module's tests before finishing.

**Why this works:** Isolation is a prerequisite for safe parallelism. Without it, two agents editing the same file create merge conflicts or silently overwrite each other's changes. Worktrees provide filesystem-level isolation with zero overhead compared to full repo clones.

**Pros:**
- Eliminates file conflicts between parallel agents
- Each worktree has a clean, testable state
- Merging is straightforward with git's tooling

**Cons:**
- Disk space multiplied by number of worktrees
- Merge conflicts still possible at merge time if scopes overlap

**Deeper:** See `approaches/worktree-isolation.md`

---

### 5. Built-In Review Skills — Polish after the heavy lifting
**Level:** Beginner

After a large refactor, there's always cleanup: dead imports, redundant helper functions that are now one-liners, duplicated code that can be consolidated. `/simplify` catches these systematically, turning your "it works" refactor into a "it works and it's clean" refactor.

**Try it now:**
> /simplify

**Why this works:** Refactoring shifts your attention to structure, which means you stop noticing local inefficiencies. A dedicated cleanup pass with fresh eyes (or fresh AI context) catches the small improvements that compound into significantly cleaner code.

**Pros:**
- Catches post-refactor dead code and redundancy
- Low effort, high polish
- Builds confidence before opening the PR

**Cons:**
- Only handles surface-level cleanup, not structural issues

**Deeper:** See `approaches/built-in-review-skills.md`

---

### 6. Checkpoints & Rewind — Bold refactoring with a safety net
**Level:** Beginner

Sometimes the right refactoring approach is "try it and see." Checkpoints let you commit to a bold structural change knowing you can instantly revert to a known-good state if it doesn't work out. This is especially valuable early in a refactor when you're still discovering the right decomposition.

**Try it now:**
> Create a checkpoint before I start. I'm going to try inlining the `AbstractBaseRepository` into each concrete repository — I think the abstraction layer is more complexity than it's worth, but I'm not sure until I see the result. If the concrete repos end up cleaner, we keep it. If they get messier, rewind.

**Why this works:** Fear of irreversibility makes developers conservative. Checkpoints remove that fear, enabling bolder experiments. The best refactoring strategies often emerge from trying an approach, evaluating the result, and either committing or reverting — which is only practical with cheap, reliable undo.

**Pros:**
- Enables experimental refactoring without risk
- Instant revert to any previous state
- Encourages bolder, potentially better designs

**Cons:**
- Can become a crutch that delays committing to a direction
- Only tracks Claude's changes, not manual edits between checkpoints

**Deeper:** See `approaches/checkpoints-rewind.md`

---

### 7. Hooks — Auto-format and auto-lint after every refactoring edit
**Level:** Intermediate

During large refactors, maintaining consistent formatting and lint compliance across dozens of files is tedious. A PostToolUse hook runs your formatter (Prettier, Black, gofmt) and linter after every file edit, so Claude's output is always clean. This prevents the "refactor is done but now I need to fix 200 lint errors" problem.

**Try it now:**
> Set up a PostToolUse hook that runs `npx prettier --write` on any edited file and then `npx eslint --fix` on it. I'm extracting the validation logic from `src/controllers/order_controller.ts` into a new `src/services/order-validator.ts`. Every file you touch should be formatted and lint-clean before moving to the next one.

**Why this works:** Refactoring generates many small edits across many files. Without automated formatting, style drift accumulates silently and creates a noisy diff that obscures the structural changes. Hooks keep every file clean as you go, so the final diff shows only the refactor — not formatting noise.

**Pros:**
- Eliminates the post-refactor formatting cleanup pass
- Keeps diffs clean and reviewable throughout the process
- Catches lint violations at edit time, not at PR review time

**Cons:**
- Formatter and linter runs add latency to each edit cycle
- Conflicting formatter opinions between tools can cause thrashing

**Deeper:** See `approaches/hooks-as-workflow.md`
