# incident-response
*Last verified: 2026-07-03*

**Hidden gem:** Deep Research — two minutes checking whether it's a known upstream outage can save an hour of debugging code that was never the problem.

**Exemplar move:** Connect to Grafana and Datadog MCP servers. payment-service returns 503s on ~12% of requests — pull error timeline, latest deploy, top service:payment-service log errors; correlate spike with deploy timestamp.

**Plugins:** `sentry` · `datadog` · `grafana-mcp` · `pagerduty` · `rootly` — all ☑️, all need the vendor account already in use at your org.

**Built-ins:** `/loop` — poll the deploy or error rate while you investigate; `/deep-research` — check whether the root cause is a known upstream issue. Facts and pitfalls per command: `registry/builtin-commands.md`.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [MCP Context](../approaches/mcp-context.md) | Intermediate | Error spike correlating with a recent deploy | Incidents are data problems — MCP makes Claude a unified query layer across your observability stack, collapsing data-gathering |
| 2 | [Deep Research](../approaches/deep-research.md) | Beginner | Outage symptoms match something you've seen in a dependency | Many incidents originate outside your codebase; checking known issues first saves hours debugging code that wasn't the problem |
| 3 | [Plan Mode](../approaches/plan-mode.md) | Beginner | Service down and you don't know where to start | In incidents a wrong action costs more than a delayed one — structured triage prevents fixing the wrong component |
| 4 | [Channels](../approaches/channels.md) | Advanced | Alerts should land in the session with context already loaded | Triage speed is dominated by context assembly; channels deliver the event where that assembly already happened |
