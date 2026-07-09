# greenfield
*Last verified: 2026-07-03*

**Hidden gem:** Custom Skills — encoding your conventions as a scaffold command before writing feature #2 pays off for every feature after it.

**Exemplar move:** Enter plan mode. Design a notification service (email, Slack, in-app; consumes Kafka events) — module structure, data models, API surface, delivery retries, per-channel rate limiting. Architecture only, no code.

**Plugins:** `feature-dev` ✅ guided feature workflow · `frontend-design` ✅ production-grade UI · `figma` ☑️ design-to-code · `laravel-boost` ☑️ — ~35 more stack starters in the catalog; grep by framework.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | Beginner | Starting a new service or module from scratch | In greenfield work code is cheap but architecture is expensive; planning makes structural decisions explicit, not accidental |
| 2 | [Deep Research](../approaches/deep-research.md) | Beginner | Need to choose between libraries or frameworks | Technology choices compound — twenty minutes of research avoids the wrong-library realization after 10,000 lines of integration |
| 3 | [Autonomous Loops](../approaches/autonomous-loops.md) | Intermediate | Feature is well-defined and you want speed | Loops excel at mechanical build-test-fix iteration — clear, testable criteria tell Claude when it's done |
| 4 | [Custom Skills](../approaches/custom-skills.md) | Beginner | Your team scaffolds the same structure for every new module | The first 20% of every module is patterned boilerplate — skills encode it once and skip to the interesting work |
