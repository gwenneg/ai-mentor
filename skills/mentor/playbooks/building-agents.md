# building-agents
*Last verified: 2026-07-03*

**Hidden gem:** Custom Agents — prototyping your agent as a ten-line `.claude/agents/` file answers most of the design questions (tools, model, instructions) before you write a single line of SDK code.

**Exemplar move:** Enter plan mode. Design a support-inbox triage agent — classify severity, draft responses, escalate billing/data-loss to humans: what tools, forbidden actions, escalation boundaries, minimal first version worth shipping?

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | Unsure what the agent should even do | An agent is a policy wrapped around a model — explicit designs fail predictably, emergent ones fail creatively |
| 2 | [Custom Agents](../approaches/custom-agents.md) | Want a working prototype today | Editing markdown and re-running is the cheapest iteration loop — converge on instructions and tool surface in an afternoon |
| 3 | [Official Plugins](../approaches/official-plugins.md) | Ready to build a standalone product | Agent infrastructure is undifferentiated heavy lifting; a production-tested engine puts your effort into the agent's judgment |
| 4 | [Permissions & Safe Autonomy](../approaches/safe-autonomy.md) | Need to constrain what the agent can do | Capability is easy and trust hard to win back — agents earn adoption through provable boundaries, not demos |
| 5 | [agent-sdk-dev](../approaches/agent-sdk-dev.md) | Starting an Agent SDK project from a validated scaffold | A strict-TS scaffold with streaming already wired answers the setup questions so effort goes into the agent's judgment |
| 6 | [agent-sdk](../approaches/agent-sdk.md) | The agent is a product, not a Claude Code workflow | Programmatic sessions and custom tool definitions are a different altitude than .claude/agents/ files — the SDK is the supported path for standalone agents |
