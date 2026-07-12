# building-mcp-integrations
*Last verified: 2026-07-12*

**Hidden gem:** Headless Mode — a shell script of `claude -p` calls against your server is an MCP integration test suite: repeatable, CI-friendly, and it tests the thing that actually matters (can a model use your tools correctly?).

**Exemplar move:** Enter plan mode. Design MCP tool surface for incident management (search, read detail, comment, change status, page on-call): tools vs. left out, descriptions, read-only ones, damage potential.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Plan Mode](../approaches/techniques/plan-mode.md) | Not sure what the server should expose | Models select tools by description — a well-designed five-tool surface outperforms a twenty-tool REST dump |
| 2 | [mcp-server-dev](../approaches/tools/mcp-server-dev.md) | Building or exposing an API as an MCP server | Hands-on validated guided flow with current SDK idioms — idioms move faster than training data, and design-plus-implementation in one pipeline is exactly the "expose our API to Claude" job |
| 3 | [Official Plugins](../approaches/techniques/official-plugins.md) | Ready to implement | MCP's protocol surface is undifferentiated work; a guided workflow encodes current conventions so effort goes into tool semantics |
| 4 | [MCP Context](../approaches/techniques/mcp-context.md) | Haven't used MCP as a consumer yet | Interface intuition comes from the consumer side — a week consuming MCP beats a month producing it blind |
| 5 | [Headless Mode](../approaches/techniques/headless-mode.md) | Server built, needs regression testing | Whether a model can use your tools correctly is the real acceptance criterion — headless runs make it executable |
