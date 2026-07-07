# code-review
*Last verified: 2026-07-03*

**Hidden gem:** MCP Context — reviewing against the ticket's acceptance criteria instead of just the diff catches the most expensive bugs: code that works but solves the wrong problem.

**Exemplar move:** Run `/code-review high` on the branch diff — focus on src/services/billing/ discount stacking edge cases and Stripe webhook failure-mode error handling.

**Plugins:** `pr-review-toolkit` ✅ 6-agent deep review · `github` ☑️ PR management · `sonarqube` ☑️ and `qodo-skills` ☑️ quality/security scanning.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Built-In Review Skills](../approaches/built-in-review-skills.md) | Beginner | Quick review of a focused diff or your own pre-PR code | Most review value is systematic checking — codified reviewer instincts applied consistently to every line |
| 2 | [Subagent Delegation](../approaches/subagent-delegation.md) | Advanced | Large PR touching security, perf, and correctness | Attention is finite — specialized agents keep focus per concern while parallelism keeps wall-clock time low |
| 3 | [Fan-Out Workflows](../approaches/fan-out-workflows.md) | Advanced | Critical change needing adversarial verification | Automated review's biggest problem is false positives; adversarial verification mimics human pushback for higher-signal results |
| 4 | [MCP Context](../approaches/mcp-context.md) | Intermediate | PR implements a design doc or addresses an issue | The most expensive bugs are specification bugs — grounding review in requirements catches code solving the wrong problem |
| 5 | [Cloud Sessions](../approaches/cloud-sessions.md) | Intermediate | Deep review off your machine, or PRs that fix themselves | Most post-review churn is mechanical; a watching cloud agent handles it so human reviewers focus on design |
