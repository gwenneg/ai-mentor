# llm-features
*Last verified: 2026-07-03*

**Hidden gem:** Autonomous Loops — point a loop at your eval suite and prompt engineering becomes test-driven development: `/goal 90% of eval cases pass` turns the squishiest part of LLM work into the most mechanical.

**Exemplar move:** Enter plan mode. Design AI-generated ticket summaries: prompt strategy, UI handling of bad/refused summaries, eval plan, model tier and cost per 1K tickets, latency budget. No code.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Plan Mode](../approaches/techniques/plan-mode.md) | Designing the feature before building | The hard part isn't calling the API — measurable quality and graceful failure are design decisions code inherits |
| 2 | [LLM Evals](../approaches/techniques/llm-evals.md) | Prompt changes keep shipping regressions users find first | A change's effect is invisible in any single output — a graded case suite turns "feels better" into a number you can gate on |
| 3 | [Deep Research](../approaches/techniques/deep-research.md) | Choosing models, patterns, or architecture | LLM platform facts have a shelf life of weeks — current documentation beats stale training-data recall |
| 4 | [Autonomous Loops](../approaches/techniques/autonomous-loops.md) | Eval suite exists, quality is inconsistent | Prompts regress invisibly because quality is a distribution; evals plus a loop pin it to a number and push it up |
| 5 | [agent-sdk](../approaches/tools/agent-sdk.md) | Building AI features into your own product | Programmatic sessions and custom tools give product code the agent loop without reinventing it |
