# code-review
*Last verified: 2026-07-03*

**Hidden gem:** MCP Context — reviewing against the ticket's acceptance criteria instead of just the diff catches the most expensive bugs: code that works but solves the wrong problem.

**Exemplar move:** Run `/code-review high` on the branch diff — focus on src/services/billing/ discount stacking edge cases and Stripe webhook failure-mode error handling.

**Plugins:** `pr-review-toolkit` ✅ 6-agent deep review · `github` ☑️ PR management · `sonarqube` ☑️ and `qodo-skills` ☑️ quality/security scanning.

**Integrations:** `claude-code-action` — Claude reviews on every PR via the maintained GitHub Action, findings posted as inline comments. Facts and pitfalls: its `solutions/<id>.md` record.

**Built-ins:** `/code-review` — structured correctness pass on the diff; `/security-review` — security lens on pending changes; `/verify` — watch the change actually work before merging. Facts and pitfalls per command: its `solutions/<id>.md` record.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Built-In Review Skills](../solutions/built-in-review-skills.md) | Quick review of a focused diff or your own pre-PR code | Most review value is systematic checking — codified reviewer instincts applied consistently to every line |
| 2 | [Subagent Delegation](../solutions/subagent-delegation.md) | Large PR touching security, perf, and correctness | Attention is finite — specialized agents keep focus per concern while parallelism keeps wall-clock time low |
| 3 | [Fan-Out Workflows](../solutions/fan-out-workflows.md) | Critical change needing adversarial verification | Automated review's biggest problem is false positives; adversarial verification mimics human pushback for higher-signal results |
| 4 | [MCP Context](../solutions/mcp-context.md) | PR implements a design doc or addresses an issue | The most expensive bugs are specification bugs — grounding review in requirements catches code solving the wrong problem |
| 5 | [Cloud Sessions](../solutions/cloud-sessions.md) | Deep review off your machine, or PRs that fix themselves | Most post-review churn is mechanical; a watching cloud agent handles it so human reviewers focus on design |
