# building-mcp-integrations
*Last verified: 2026-07-03*

**Hidden gem:** Headless Mode — a shell script of `claude -p` calls against your server is an MCP integration test suite: repeatable, CI-friendly, and it tests the thing that actually matters (can a model use your tools correctly?).

**Exemplar move:** Enter plan mode. Design MCP tool surface for incident management (search, read detail, comment, change status, page on-call): tools vs. left out, descriptions, read-only ones, damage potential.

**Plugins:** `mcp-server-dev` ✅ guided server design · `mcp-apps` ☑️ MCP Apps SDK · `mcp-tunnels` ☑️ private-server tunnels.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | none | Not sure what the server should expose | Models select tools by description — a well-designed five-tool surface outperforms a twenty-tool REST dump |
| 2 | [Official Plugins](../approaches/official-plugins.md) | some | Ready to implement | MCP's protocol surface is undifferentiated work; a guided workflow encodes current conventions so effort goes into tool semantics |
| 3 | [MCP Context](../approaches/mcp-context.md) | some | Haven't used MCP as a consumer yet | Interface intuition comes from the consumer side — a week consuming MCP beats a month producing it blind |
| 4 | [Headless Mode](../approaches/headless-mode.md) | some | Server built, needs regression testing | Whether a model can use your tools correctly is the real acceptance criterion — headless runs make it executable |
