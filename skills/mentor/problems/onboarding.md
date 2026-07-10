# onboarding
*Last verified: 2026-07-03*

**Hidden gem:** Custom Skills — a `/setup-dev` skill is executable documentation: it can't silently go stale the way a wiki page does.

**Exemplar move:** Pull the Confluence "Engineering Onboarding" guide and #platform-team pinned messages, cross-reference with README.md and docker-compose.yml — flag outdated steps, produce one consolidated setup guide.

**Plugins:** `claude-code-setup` ✅ repo-tailored automation recommendations · `learn-with-coursera` ☑️ learning paths.

**Built-ins:** `/init` — generate the repo's starter CLAUDE.md. Facts and pitfalls per command: its `solutions/<id>.md` record.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [MCP Context](../solutions/mcp-context.md) | Team knowledge scattered across Confluence, Slack, Notion, and READMEs | The onboarding bottleneck is finding information, not understanding it — MCP brings every source into one queryable context |
| 2 | [Plan Mode](../solutions/plan-mode.md) | Need a systematic understanding of the codebase and architecture | Replicates how experts onboard — learn the architecture's shape first, fill in module details as tickets demand |
| 3 | [Deep Research](../solutions/deep-research.md) | Codebase uses frameworks or patterns you've never worked with | Patterns encode decisions into structure; learning the stack's vocabulary lets you read the codebase fluently from day one |
| 4 | [Custom Skills](../solutions/custom-skills.md) | Local dev setup takes a half-day of manual steps | A skill is executable documentation — when a step breaks you fix the skill and every future onboarder benefits |
