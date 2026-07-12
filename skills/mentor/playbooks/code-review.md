# code-review
*Last verified: 2026-07-12*

**Hidden gem:** MCP Context — reviewing against the ticket's acceptance criteria instead of just the diff catches the most expensive bugs: code that works but solves the wrong problem.

**Exemplar move:** Run `/code-review high` on the branch diff — focus on src/services/billing/ discount stacking edge cases and Stripe webhook failure-mode error handling.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Built-In Review Skills](../approaches/techniques/built-in-review-skills.md) | Quick review of a focused diff or your own pre-PR code | Most review value is systematic checking — codified reviewer instincts applied consistently to every line |
| 2 | [Subagent Delegation](../approaches/techniques/subagent-delegation.md) | Large PR touching security, perf, and correctness | Attention is finite — specialized agents keep focus per concern while parallelism keeps wall-clock time low |
| 3 | [Fan-Out Workflows](../approaches/techniques/fan-out-workflows.md) | Critical change needing adversarial verification | Automated review's biggest problem is false positives; adversarial verification mimics human pushback for higher-signal results |
| 4 | [MCP Context](../approaches/techniques/mcp-context.md) | PR implements a design doc or addresses an issue | The most expensive bugs are specification bugs — grounding review in requirements catches code solving the wrong problem |
| 5 | [Cloud Sessions](../approaches/techniques/cloud-sessions.md) | Deep review off your machine, or PRs that fix themselves | Most post-review churn is mechanical; a watching cloud agent handles it so human reviewers focus on design |
| 6 | [pr-review-toolkit](../approaches/tools/pr-review-toolkit.md) | A substantial PR deserves specialized review angles beyond one general pass | Six specialized reviewers each read one concern deeply — found a planted off-by-one at the exact line in evaluation |
| 7 | [claude-code-action](../approaches/tools/claude-code-action.md) | Claude should review every PR automatically, findings as inline comments | Review that runs on every PR catches what reviewers skip on busy days; the Action makes it one workflow file |
