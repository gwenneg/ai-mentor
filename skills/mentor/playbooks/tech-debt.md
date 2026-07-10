# tech-debt
*Last verified: 2026-07-03*

**Hidden gem:** Custom Agents — encoding your team's own deprecated patterns as a detector turns tribal knowledge into a tracked migration metric.

**Exemplar move:** Spawn four parallel agents auditing src/: duplication clusters, deprecated APIs (TODO/FIXME/@deprecated), test-coverage gaps vs tests/, functions over cyclomatic complexity 10 — consolidate into one prioritized report.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Subagent Delegation](../solutions/subagent-delegation.md) | Multi-dimensional audit across the whole codebase | Each debt dimension needs a different scanning strategy — specialized agents apply the right one without conflating them |
| 2 | [Built-In Review Skills](../solutions/built-in-review-skills.md) | Quick sense of the worst offenders, no setup | The most impactful debt is often the most visible — built-in skills catch the low-hanging fruit systematically |
| 3 | [Plan Mode](../solutions/plan-mode.md) | You have findings and need to prioritize them | Prioritization weighs bug frequency, developer drag, and roadmap fit — plan mode forces explicit trade-off reasoning |
| 4 | [Custom Agents](../solutions/custom-agents.md) | Recurring patterns specific to your project | The most impactful debt is patterns your team decided to leave; custom agents enforce those decisions and track progress |
| 5 | [Background Agents](../solutions/background-agents.md) | Audit takes an hour but your attention shouldn't | Audits are attention-light but time-heavy — backgrounding turns them from a thing you do into a thing you review |
