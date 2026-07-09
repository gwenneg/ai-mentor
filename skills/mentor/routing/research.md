# research
*Last verified: 2026-07-03*

**Hidden gem:** Plan Mode — defining evaluation criteria before gathering evidence is the only reliable defense against confirmation bias, and nobody thinks of plan mode for research.

**Exemplar move:** /deep-research Compare RabbitMQ, Kafka, NATS JetStream for ~50K events/sec, exactly-once payment delivery, Go and Python services — throughput, delivery guarantees, operational complexity, client maturity, AWS managed hosting.

**Plugins:** `firecrawl` ☑️ crawling and structured extraction · `microsoft-docs` ☑️ Azure/.NET references · `zyte-web-data` ☑️ scraping — ~19 more domain-specific research tools in the catalog.

**Built-ins:** `/deep-research` — multi-source, adversarially verified report. Facts and pitfalls per command: `registry/builtin-commands.md`.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Deep Research](../approaches/deep-research.md) | none | Comparing libraries, frameworks, or SaaS tools | Automates the 30-browser-tab workflow into a cited, cross-referenced report you can present to your team |
| 2 | [Plan Mode](../approaches/plan-mode.md) | none | Need a structured evaluation framework first | Research without a framework is browsing — criteria first, evidence second, decision last is how staff engineers decide |
| 3 | [Browser Integration](../approaches/browser-integration.md) | involved | Evaluating tools with interactive demos or UIs | Tools are chosen on experience as much as capability — browser testing evaluates the experience, not the feature list |
| 4 | [Model & Effort Selection](../approaches/model-effort-selection.md) | involved | Different phases need different reasoning depth | Reasoning depth is a budget — allocating it per phase gets better conclusions and a cheaper session at once |
