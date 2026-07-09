# building-agents
*Last verified: 2026-07-03*

**Hidden gem:** Custom Agents — prototyping your agent as a ten-line `.claude/agents/` file answers most of the design questions (tools, model, instructions) before you write a single line of SDK code.

**Exemplar move:** Enter plan mode. Design a support-inbox triage agent — classify severity, draft responses, escalate billing/data-loss to humans: what tools, forbidden actions, escalation boundaries, minimal first version worth shipping?

**Plugins:** `agent-sdk-dev` ✅ Agent SDK scaffolding · `pydantic-ai` ☑️ and `atomic-agents` ☑️ framework-specific patterns · `aws-agents` ☑️ Bedrock AgentCore — 2 more in the catalog.

**Integrations:** `agent-sdk` — the supported path when the agent is a product, not a Claude Code workflow. Facts and pitfalls per record: `registry/integrations.md`.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | Unsure what the agent should even do | An agent is a policy wrapped around a model — explicit designs fail predictably, emergent ones fail creatively |
| 2 | [Custom Agents](../approaches/custom-agents.md) | Want a working prototype today | Editing markdown and re-running is the cheapest iteration loop — converge on instructions and tool surface in an afternoon |
| 3 | [Official Plugins](../approaches/official-plugins.md) | Ready to build a standalone product | Agent infrastructure is undifferentiated heavy lifting; a production-tested engine puts your effort into the agent's judgment |
| 4 | [Permissions & Safe Autonomy](../approaches/safe-autonomy.md) | Need to constrain what the agent can do | Capability is easy and trust hard to win back — agents earn adoption through provable boundaries, not demos |
