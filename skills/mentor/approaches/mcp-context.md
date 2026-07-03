# MCP Context (Model Context Protocol)
*Last verified: 2026-06-27*

## What It Is

MCP Context lets your AI coding tool pull information from external systems — issue trackers, design documents, databases, monitoring dashboards, Slack threads — directly into the conversation. Instead of copy-pasting context into your prompt, you configure connections once and the AI queries them as needed. It turns your AI assistant from a tool that only sees files on disk into one that can reach across your entire development ecosystem.

## Why It Works

The quality of AI-generated code is directly proportional to the quality of context it receives. Most coding mistakes happen not because the AI cannot write code, but because it lacks the requirements, constraints, or domain knowledge that a human developer would check before writing. MCP closes that gap by giving the AI access to the same sources of truth you would consult: the ticket describing the feature, the architecture decision record explaining why the system is structured this way, the monitoring data showing how the code actually behaves in production.

## When to Use It

- Reviewing code against acceptance criteria pulled from your issue tracker (Jira, Linear, GitHub Issues)
- Understanding unfamiliar code by pulling architecture docs or design specs from Confluence or Notion
- Correlating a code change with production behavior by querying metrics, logs, or error tracking
- Writing database migrations by inspecting the live schema rather than relying on outdated documentation
- Pulling Slack thread context into a debugging session so the AI understands the reported symptoms
- Setting up an MCP server to solve a recurring workflow problem — e.g., you constantly need Jira context, so you configure the Jira MCP server once and use it everywhere

## When NOT to Use It

- When the context you need is already in the repository (READMEs, inline comments, test fixtures) — file reads are simpler and faster
- When the external system requires complex authentication that is not already configured — fighting with OAuth flows mid-session derails the actual work

## How It Works

### Basic (Beginner)

**Using MCP context in your workflow** — when the server is already configured:

1. Identify the external context your task needs. For example: "I need to implement the feature described in PROJ-1234."
2. Ensure the relevant MCP server is configured in your `.mcp.json` (project-scoped) or `~/.claude.json` (user-scoped). Your project may already have this set up.
3. Ask Claude to pull the context: "Read the requirements from PROJ-1234 and summarize what needs to change."
4. Claude calls the MCP tool (e.g., `jira_get_issue`), retrieves the ticket, and now has the acceptance criteria in its context window.
5. Continue your task with the AI grounded in real requirements: "Now implement the changes described in that ticket."

**Setting up an MCP server** — when you identify a gap in your workflow:

If you keep copying context from an external system into your prompts, that is the signal to set up an MCP server. For example: "I review code against Jira requirements every day but paste the ticket text manually each time." The fix is to configure the Jira MCP server once so Claude can fetch tickets directly. Ask Claude: "Help me set up the Jira MCP server in `.mcp.json` so I can pull ticket details directly." The same applies to Confluence, Notion, GitHub, Slack, databases, and monitoring tools — if you are manually bridging information, an MCP server removes that friction permanently.

### Composing with Other Approaches (Intermediate)

- **MCP Context plus Plan Mode**: Pull requirements from your issue tracker, then enter Plan Mode. Claude proposes an implementation plan grounded in the actual acceptance criteria rather than your paraphrased summary. Review the plan against the original ticket before approving.
- **MCP Context plus subagents**: A parent agent pulls the architecture overview from Confluence, then delegates implementation tasks to subagents. Each subagent inherits the architectural context without needing to re-fetch it.
- **MCP Context plus code review**: After making changes, ask Claude to pull the original requirements again and verify that the implementation satisfies each acceptance criterion. This is a lightweight requirements traceability check.

### Advanced Patterns

- **MCP Tool Search for large tool sets**: If your project configures dozens of MCP servers, the tool schemas can consume a significant portion of the context window. Enable MCP Tool Search (controlled via the `ENABLE_TOOL_SEARCH` environment variable (enabled by default)) to defer loading tool schemas until Claude actually needs them. This lets you register hundreds of tools without paying the context cost upfront.
- **Cross-system correlation**: Pull a Sentry error report, then query your database schema, then read the relevant code — all in one conversation. The AI connects the dots across systems that normally require three browser tabs and manual context-switching.
- **Production-informed refactoring**: Query your monitoring system for the slowest API endpoints, then ask Claude to profile and optimize the relevant code paths. The refactoring is driven by real data, not guesswork.

## Common Pitfalls

- **Over-fetching context**: Pulling an entire Confluence space into the conversation wastes context window tokens. Be specific: fetch one document, one ticket, one query result. You can always fetch more if needed.
- **Skipping server setup when it would solve the problem**: If you find yourself repeatedly pasting context from an external system, that is the signal to configure an MCP server. Do the setup once, then focus on workflow.
- **Stale external data**: MCP queries return point-in-time snapshots. If your Jira ticket gets updated mid-session, Claude still has the old version. Re-fetch if requirements may have changed.
- **Sensitive data exposure**: MCP servers can expose production databases, customer data, or internal communications. Understand what data each MCP tool can access and whether your organization permits it in AI contexts.

## Real-World Example

You are assigned BILLING-892: "Customers on the annual plan see incorrect proration when upgrading mid-cycle." Before writing any code, you ground Claude in the real requirements:

```
> Fetch the details of BILLING-892 from Jira.
```

Claude calls `jira_get_issue` and retrieves the ticket, including acceptance criteria: proration should be calculated from the upgrade date to the next renewal date, using the daily rate difference. It also pulls two linked tickets with customer-reported examples showing specific dollar amounts.

```
> Now query the billing database schema for the subscriptions and invoices tables.
```

Claude calls `postgres_query` via the database MCP server and retrieves the column definitions, revealing that `proration_start_date` is stored as a `DATE` but the billing calculation in `src/billing/proration.ts` treats it as a UTC timestamp, causing off-by-one errors at certain times of day.

Claude proposes a fix in `calculateProration()` on line 47 of `src/billing/proration.ts`, normalizing the date comparison to use date-only arithmetic. You verify the fix against the customer examples from the linked tickets — the dollar amounts now match exactly.

You ask Claude to pull the requirements one more time and verify each acceptance criterion against the implementation. All three criteria are satisfied. The entire debugging session used real requirements and real schema, not guesswork about what the code was supposed to do.

## Sources

- [Claude Code MCP](https://code.claude.com/docs/en/mcp) — Official docs for connecting Claude Code to external tools via MCP
- [Model Context Protocol](https://modelcontextprotocol.io/docs/getting-started/intro) — General MCP protocol overview and architecture
- [MCP Reference Servers](https://github.com/modelcontextprotocol/servers) — Official repository of MCP server implementations
