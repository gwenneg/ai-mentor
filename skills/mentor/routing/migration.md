# migration
*Last verified: 2026-07-03*

**Hidden gem:** Worktree Isolation — running the upgrade in a disposable copy first turns "should we migrate?" from changelog speculation into a concrete damage report.

**Exemplar move:** Enter plan mode. Migrate React Router v5 to v6: 34 routes in src/routes/, useHistory in ~20 components, guards in src/auth/ProtectedRoute.tsx — map changes, categorize, order to stay functional.

**Plugins:** `code-modernization` ✅ legacy-codebase migration (COBOL, old Java, monoliths) · `ui5-modernization`/`ui5-typescript-conversion` ☑️ SAPUI5 · `aws-transform`/`migration-to-aws` ☑️ moves to AWS — grep the catalog when the user names a stack.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | Beginner | Large migration touching dozens of files | Migrations are dependency graphs, not task lists — planning reveals the structure so you traverse it in the right order |
| 2 | [Subagent Delegation](../approaches/subagent-delegation.md) | Advanced | Migration across independent modules or services | Migrations are embarrassingly parallel when modules don't share state — subagents do in minutes what takes sequential hours |
| 3 | [Autonomous Loops](../approaches/autonomous-loops.md) | Intermediate | Upgrade where "it compiles and tests pass" is the goal | Mechanical migrations have a clear convergence criterion — AI doesn't get bored or lose focus on iteration 47 |
| 4 | [Worktree Isolation](../approaches/worktree-isolation.md) | Intermediate | Risky upgrade you want to test without polluting main | A throwaway environment changes how you approach risk — bolder strategies, faster discovery, nothing at stake |
