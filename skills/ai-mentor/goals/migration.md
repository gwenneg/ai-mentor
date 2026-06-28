# Migration & Upgrades
*Last reviewed: 2026-06-27*

## When You're Here

You need to move from one version of something to another — a framework upgrade, an API deprecation, a dependency bump with breaking changes. The codebase compiles today, and your job is to make it compile tomorrow on the new version without breaking what already works. Migrations are stressful because the blast radius is often unclear until you're halfway through.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Large migration touching dozens of files | Plan mode | You need the full scope mapped before touching anything |
| Migration across independent modules or services | Subagent delegation | Parallelize work that doesn't share state |
| Upgrade where "it compiles and tests pass" is the goal | Autonomous loops | Set the success criteria and let Claude grind |
| Risky upgrade you want to test without polluting main | Worktree isolation | Fail safely in a throwaway branch |
| Unfamiliar framework with unclear migration path | Deep research | Find the official guide and known pitfalls first |
| Multi-step migration where any step might break things | Checkpoints & rewind | Save progress so you can roll back individual steps |

## Approaches (Ranked)

### 1. Plan Mode — Map the full migration scope before starting
**Level:** Beginner | **Tools:** Any

Migrations fail when you discover halfway through that a change in module A cascades into modules B, C, and D. Plan mode forces you to inventory every affected file, every deprecated API call, and every breaking change before writing a single line. This upfront investment pays for itself by eliminating surprise rework.

**Try it now:**
> Enter plan mode. I need to migrate our React app from React Router v5 to v6. The app has 34 route definitions across `src/routes/`, uses `useHistory` in about 20 components, and has custom route guards in `src/auth/ProtectedRoute.tsx`. Map every file that needs changes, categorize them by type of change (API rename, pattern replacement, logic rewrite), and give me a migration order that keeps the app functional at each step.

**Why this works:** Migrations are dependency graphs, not task lists. Plan mode reveals the graph structure so you can traverse it in the right order instead of discovering edges the hard way.

**Pros:**
- Prevents the "one more file" surprise that turns a 2-hour task into 2 days
- Creates a checklist you can track progress against
- Identifies the riskiest changes upfront so you can plan testing around them

**Cons:**
- Takes 10-15 minutes before you write any code — feels slow under deadline pressure

**Deeper:** See `approaches/plan-mode.md`

---

### 2. Subagent Delegation — Parallelize migration across modules
**Level:** Advanced | **Tools:** Claude Code

When your migration plan reveals independent workstreams — say, updating API calls in the payments module has no overlap with updating them in the notifications module — subagent delegation lets you run those migrations simultaneously. Each subagent gets a focused scope, applies the changes, and reports back.

**Try it now:**
> I'm migrating from `axios` 0.x to `axios` 1.x across our monorepo. The breaking changes are: `AxiosRequestConfig` renamed to `InternalAxiosRequestConfig` in interceptors, `data` in error responses moved to `error.response.data`, and default timeout behavior changed. Delegate three subagents: one for `packages/api-client/`, one for `packages/dashboard/`, and one for `packages/admin-tools/`. Each should apply all three breaking changes and run the package's unit tests to verify.

**Why this works:** Migrations are embarrassingly parallel when modules don't share mutable state. Subagents exploit this by doing in 5 minutes what would take 15 sequentially, with each agent maintaining full context for its scope.

**Pros:**
- Dramatically faster for monorepos and multi-module projects
- Each subagent has focused context, reducing mistakes
- Failures in one module don't block others

**Cons:**
- Only works when modules are genuinely independent
- Coordinating shared code (utils, types) still needs a single pass first

**Deeper:** See `approaches/subagent-delegation.md`

---

### 3. Autonomous Loops — Set "build passes after migration", Claude iterates
**Level:** Intermediate | **Tools:** Claude Code

Some migrations are straightforward but tedious: rename 150 API calls, update import paths, fix type mismatches. You know exactly what "done" looks like — the build passes and tests are green. Autonomous loops let you define that success condition and walk away while Claude grinds through the mechanical work.

**Try it now:**
> /goal: Migrate all usages of our deprecated `@company/logger` v2 to `@company/logger` v3 in `services/order-service/`. In v3, `logger.info(message, meta)` becomes `logger.info({ msg: message, ...meta })`, `logger.child()` becomes `logger.createChild()`, and the `Logger` type is now imported from `@company/logger/types`. Update all files, then run `npm run build && npm test` to verify. Keep iterating until both pass with zero errors.

**Why this works:** Mechanical migrations have a clear convergence criterion — the build is either green or it isn't. AI excels at tight edit-compile-fix loops because it doesn't get bored or lose focus on iteration 47.

**Pros:**
- Handles tedious bulk changes without supervision
- Self-verifies by running the build after each round of changes
- Catches cascading type errors that manual search-and-replace misses

**Cons:**
- Can make "creative" fixes to pass the build that aren't what you intended — review the diff
- Not suitable when "tests pass" doesn't fully validate the migration

**Deeper:** See `approaches/autonomous-loops.md`

---

### 4. Worktree Isolation — Test migration in a sandbox without risking main
**Level:** Intermediate | **Tools:** Claude Code

Migrations are inherently risky. Worktree isolation gives you a disposable copy of your repo where you can attempt the upgrade, run the full test suite, and evaluate the results — all without touching your working branch. If the migration goes sideways, you delete the worktree and lose nothing.

**Try it now:**
> Create a worktree from main. In it, upgrade `Django` from 4.2 to 5.1 in `requirements.txt` and run `pip install -r requirements.txt`. Then run `python manage.py test --parallel` and capture every deprecation warning and failure. Summarize what breaks and categorize fixes as "simple rename", "logic change", or "needs investigation".

**Why this works:** The psychological safety of a throwaway environment changes how you approach risk. You'll try bolder migration strategies when failure costs nothing, and you'll discover issues faster because you're not tiptoeing around your working tree.

**Pros:**
- Zero risk to your current work — the worktree is fully disposable
- Can run the full test suite against the upgraded version immediately
- Great for evaluating "how bad is this upgrade?" before committing to it

**Cons:**
- Adds disk space for the worktree copy

**Deeper:** See `approaches/worktree-isolation.md`

---

### 5. Deep Research — Research migration guides, breaking changes, upgrade paths
**Level:** Beginner | **Tools:** Claude Code

Before writing any migration code, you need to know what changed. Deep research fans out across release notes, migration guides, GitHub issues, and community blog posts to build a comprehensive picture of breaking changes, recommended upgrade paths, and known pitfalls that others have already encountered.

**Try it now:**
> /deep-research I need to upgrade our Spring Boot app from 2.7 to 3.2. We use Spring Security, Spring Data JPA with Hibernate, and Flyway for migrations. What are all the breaking changes across these dependencies? Are there any known gotchas with the Jakarta EE namespace migration? What's the recommended upgrade order?

**Why this works:** Every major framework upgrade has a migration guide written by its maintainers, plus dozens of blog posts from developers who hit the undocumented edge cases. Deep research synthesizes all of these into a single briefing, saving hours of manual research.

**Pros:**
- Surfaces breaking changes you didn't know about before they bite you
- Finds community workarounds for known pain points
- Gives you confidence that your migration plan covers everything

**Cons:**
- Only as good as the documentation and community discussion available

**Deeper:** See `approaches/deep-research.md`

---

### 6. Checkpoints & Rewind — Revert if a migration step breaks things
**Level:** Beginner | **Tools:** Claude Code

Multi-step migrations often have a "point of no return" feel — you're three steps in, something breaks, and you're not sure which step caused it. Checkpoints let you save progress after each successful step and rewind to the last known-good state when something goes wrong, turning a linear gamble into a series of safe experiments.

**Try it now:**
> I'm migrating our database from MySQL 5.7 to 8.0 compatibility mode. After each of these steps, create a checkpoint: (1) update `GROUP BY` queries that rely on implicit sorting, (2) replace deprecated `PASSWORD()` function calls, (3) update character set declarations from `utf8` to `utf8mb4`. Run `npm run test:db` after each step. If any step fails, rewind to the last checkpoint and show me what went wrong.

**Why this works:** Checkpoints turn an irreversible process into a reversible one. When you can always go back, you move forward faster because the cost of a wrong step is seconds, not hours of manual untangling.

**Pros:**
- Eliminates the fear of making things worse during a complex migration
- Makes it easy to isolate which specific change broke things
- Preserves working states you can compare against

**Cons:**
- Adds overhead to create and manage checkpoints
- Won't help if the migration is all-or-nothing (e.g., a schema applied atomically)

**Deeper:** See `approaches/checkpoints-rewind.md`

---

### 7. Custom Agents — Migration helper with your ORM rules baked in
**Level:** Advanced | **Tools:** Claude Code

For teams that run migrations regularly — ORM upgrades, framework version bumps, API deprecation cycles — a custom agent in `.claude/agents/` encodes your project-specific migration rules. It knows your ORM conventions, your testing requirements for schema changes, your rollback patterns, and the gotchas specific to your stack. Invoke it by name and it applies all of that knowledge automatically instead of you repeating it in every prompt.

**Try it now:**
> Create a custom agent at `.claude/agents/migration-helper.md`. It should know that: we use SQLAlchemy with Alembic for migrations, every schema change needs a reversible migration with `upgrade()` and `downgrade()`, all migrations must be tested with `pytest tests/migrations/`, and we never use `op.execute()` for data migrations — those go in a separate data migration script. Then use this agent to migrate our `User` model from storing `full_name` as a single field to separate `first_name` and `last_name` fields.

**Why this works:** Migration rules are project-specific tribal knowledge — the kind of thing a senior engineer mentions in code review. A custom agent captures that knowledge in a reusable definition, so every migration follows the same standards regardless of who runs it.

**Pros:**
- Encodes project-specific migration standards and ORM conventions
- Produces consistent, reviewable migrations every time
- New team members get senior-level migration guidance automatically

**Cons:**
- Requires upfront effort to document your migration rules
- Agent needs updating when ORM or framework conventions change

**Deeper:** See `approaches/custom-agents.md`
