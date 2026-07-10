# greenfield
*Last verified: 2026-07-03*

**Hidden gem:** Custom Skills — encoding your conventions as a scaffold command before writing feature #2 pays off for every feature after it.

**Exemplar move:** Enter plan mode. Design a notification service (email, Slack, in-app; consumes Kafka events) — module structure, data models, API surface, delivery retries, per-channel rate limiting. Architecture only, no code.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | Starting a new service or module from scratch | In greenfield work code is cheap but architecture is expensive; planning makes structural decisions explicit, not accidental |
| 2 | [Deep Research](../approaches/deep-research.md) | Need to choose between libraries or frameworks | Technology choices compound — twenty minutes of research avoids the wrong-library realization after 10,000 lines of integration |
| 3 | [Autonomous Loops](../approaches/autonomous-loops.md) | Feature is well-defined and you want speed | Loops excel at mechanical build-test-fix iteration — clear, testable criteria tell Claude when it's done |
| 4 | [Custom Skills](../approaches/custom-skills.md) | Your team scaffolds the same structure for every new module | The first 20% of every module is patterned boilerplate — skills encode it once and skip to the interesting work |
| 5 | [feature-dev](../approaches/feature-dev.md) | A feature deserves a guided explore-design-review pipeline | Phased delivery packages the plan-build-review discipline so nobody has to remember to apply it |
| 6 | [frontend-design](../approaches/frontend-design.md) | UI work should come out production-grade by default | Auto-invoked design guidance turns the default output from 'works' into 'shippable' with zero prompting effort |
