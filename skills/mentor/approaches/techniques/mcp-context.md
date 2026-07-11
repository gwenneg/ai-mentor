# MCP Context (Model Context Protocol)
*Last verified: 2026-07-06*

## What It Is

MCP Context lets Claude Code pull information from external systems — issue trackers, design documents, databases, monitoring dashboards, Slack threads — directly into the conversation. Instead of copy-pasting context into your prompt, you configure connections once and the AI queries them as needed. It turns Claude from a tool that only sees files on disk into one that can reach across your entire development ecosystem.

## Why It Works

Most coding mistakes happen not because the AI cannot write code, but because it lacks the requirements, constraints, or domain knowledge a human developer would check before writing — MCP gives it the same sources of truth you would consult.

## When to Use It

- Reviewing code against acceptance criteria pulled from your issue tracker (Jira, Linear, GitHub Issues)
- Understanding unfamiliar code by pulling architecture docs or design specs from Confluence or Notion
- Correlating a code change with production behavior by querying metrics, logs, or error tracking
- Writing database migrations by inspecting the live schema rather than relying on outdated documentation
- Pulling Slack thread context into a debugging session so the AI understands the reported symptoms

## When NOT to Use It

- When the context you need is already in the repository (READMEs, inline comments, test fixtures) — file reads are simpler and faster
- When the external system requires complex authentication that is not already configured — fighting with OAuth flows mid-session derails the actual work

## How It Works

### Basic (Beginner)

1. Identify the external context your task needs. For example: "I need to implement the feature described in PROJ-1234."
2. Ensure the relevant MCP server is configured in your `.mcp.json` (project-scoped) or `~/.claude.json` (user-scoped) — your project may already have this set up. If not, add it once with `claude mcp add --transport http <name> <server-url>`, or ask Claude: "Help me set up the Jira MCP server so I can pull ticket details directly." Check connection status anytime with `/mcp`, which also handles OAuth sign-in for remote servers.
3. Ask Claude to pull the context: "Read the requirements from PROJ-1234 and summarize what needs to change."
4. Claude calls the server's issue-lookup tool, retrieves the ticket, and now has the acceptance criteria in its context window.
5. Continue your task with the AI grounded in real requirements: "Now implement the changes described in that ticket."

### Composing with Other Approaches (Intermediate)

- **MCP Context plus Plan Mode**: Pull requirements from your issue tracker, then enter Plan Mode. Claude proposes an implementation plan grounded in the actual acceptance criteria rather than your paraphrased summary. Review the plan against the original ticket before approving.
- **MCP Context plus Subagent Delegation**: A parent agent pulls the architecture overview from Confluence, distills the relevant constraints, and includes them in each subagent's task prompt. Subagents run in their own context windows and see only what the prompt passes along — the parent fetches once instead of every subagent re-querying Confluence.
- **MCP Context plus Built-in Review Skills**: After making changes, run `/code-review`, then ask Claude to pull the original requirements again and verify that the implementation satisfies each acceptance criterion. This is a lightweight requirements traceability check.

### Advanced Patterns

- **MCP Tool Search for large tool sets**: Without tool search, dozens of MCP servers' tool schemas would consume a significant portion of the context window. MCP Tool Search — on by default, controlled via the `ENABLE_TOOL_SEARCH` environment variable — loads only tool names at session start and defers full schemas until Claude actually needs them, so you can register hundreds of tools without paying the context cost upfront.
- **Cross-system correlation**: Pull a Sentry error report, then query your database schema, then read the relevant code — all in one conversation. The AI connects the dots across systems that normally require three browser tabs and manual context-switching.
- **Production-informed refactoring**: Query your monitoring system for the slowest API endpoints, then ask Claude to profile and optimize the relevant code paths. The refactoring is driven by real data, not guesswork.

## Common Pitfalls

- **Over-fetching context**: Pulling an entire Confluence space into the conversation wastes context window tokens. Be specific: fetch one document, one ticket, one query result. You can always fetch more if needed.
- **Skipping server setup when it would solve the problem**: If you find yourself repeatedly pasting context from an external system, that is the signal to configure an MCP server. Do the setup once, then focus on workflow.
- **Stale external data**: MCP queries return point-in-time snapshots. If your Jira ticket gets updated mid-session, Claude still has the old version. Re-fetch if requirements may have changed.
- **Sensitive data exposure**: MCP servers can expose production databases, customer data, or internal communications. Understand what data each MCP tool can access and whether your organization permits it in AI contexts.

## Sources

- [Claude Code MCP](https://code.claude.com/docs/en/mcp) — Official docs for connecting Claude Code to external tools via MCP
- [Model Context Protocol](https://modelcontextprotocol.io/docs/getting-started/intro) — General MCP protocol overview and architecture
- [MCP Reference Servers](https://github.com/modelcontextprotocol/servers) — Official repository of MCP server implementations

## Signals

- Setup: `.mcp.json` exists, or MCP servers configured in settings
- Session: Uses MCP-backed tools; mentions connecting external systems
