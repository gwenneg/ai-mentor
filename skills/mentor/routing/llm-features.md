# llm-features
*Last verified: 2026-07-03*

**Hidden gem:** Autonomous Loops — point a loop at your eval suite and prompt engineering becomes test-driven development: `/goal 90% of eval cases pass` turns the squishiest part of LLM work into the most mechanical.

**Exemplar move:** Enter plan mode. Design AI-generated ticket summaries: prompt strategy, UI handling of bad/refused summaries, eval plan, model tier and cost per 1K tickets, latency budget. No code.

**Plugins:** `langfuse` ☑️ tracing and evals · `huggingface-skills` ☑️ models and datasets · `pinecone` ☑️ vector search — 6 more in the catalog.

**Integrations:** `agent-sdk` — programmatic sessions and custom tools for AI features built into your own product. Facts and pitfalls per record: `registry/integrations.md`.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | none | Designing the feature before building | The hard part isn't calling the API — measurable quality and graceful failure are design decisions code inherits |
| 2 | [Deep Research](../approaches/deep-research.md) | none | Choosing models, patterns, or architecture | LLM platform facts have a shelf life of weeks — current documentation beats stale training-data recall |
| 3 | [Autonomous Loops](../approaches/autonomous-loops.md) | some | Prompt exists, quality is inconsistent | Prompts regress invisibly because quality is a distribution; evals plus a loop pin it to a number and push it up |
