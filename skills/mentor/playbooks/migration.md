# migration
*Last verified: 2026-07-03*

**Hidden gem:** Worktree Isolation — running the upgrade in a disposable copy first turns "should we migrate?" from changelog speculation into a concrete damage report.

**Exemplar move:** Enter plan mode. Migrate React Router v5 to v6: 34 routes in src/routes/, useHistory in ~20 components, guards in src/auth/ProtectedRoute.tsx — map changes, categorize, order to stay functional.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Plan Mode](../approaches/techniques/plan-mode.md) | Large migration touching dozens of files | Migrations are dependency graphs, not task lists — planning reveals the structure so you traverse it in the right order |
| 2 | [Subagent Delegation](../approaches/techniques/subagent-delegation.md) | Migration across independent modules or services | Migrations are embarrassingly parallel when modules don't share state — subagents do in minutes what takes sequential hours |
| 3 | [Autonomous Loops](../approaches/techniques/autonomous-loops.md) | Upgrade where "it compiles and tests pass" is the goal | Mechanical migrations have a clear convergence criterion — AI doesn't get bored or lose focus on iteration 47 |
| 4 | [Worktree Isolation](../approaches/techniques/worktree-isolation.md) | Risky upgrade you want to test without polluting main | A throwaway environment changes how you approach risk — bolder strategies, faster discovery, nothing at stake |
| 5 | [code-modernization](../approaches/records/code-modernization.md) | A legacy codebase needs a structured migration, not one-shot edits | Preflight analysis before transformation keeps a monolith migration from becoming a pile of broken edits |
