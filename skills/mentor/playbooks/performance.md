# performance
*Last verified: 2026-07-03*

**Hidden gem:** Hooks — benchmarking after every single edit turns optimization from guesswork into a measured experiment per change.

**Exemplar move:** Enter plan mode. /api/v2/orders averages 2.4s — trace src/controllers/orders_controller.py through OrderService.list_with_details() and OrderRepository, rank bottlenecks (N+1 queries, missing indexes), no code changes yet.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | Not sure where the bottleneck actually is | Performance budgets are finite and intuition unreliable — analysis sends effort where the data says it should go |
| 2 | [Autonomous Loops](../approaches/autonomous-loops.md) | Clear target like "under 200ms" or "below 500KB" | Tuning is inherently iterative, and a measurable target prevents both premature stopping and over-optimization |
| 3 | [Deep Research](../approaches/deep-research.md) | Need framework-specific optimization techniques | Frameworks evolve faster than developers track — research applies current best practices, not two-year-old blog advice |
| 4 | [Hooks](../approaches/hooks-as-workflow.md) | Want to see performance impact of every edit | Optimization without measurement is guesswork; auto-benchmarks make every change a data point Claude adjusts to |
