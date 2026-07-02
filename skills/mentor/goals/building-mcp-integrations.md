# Building MCP Integrations
*Last reviewed: 2026-07-02*

## When You're Here

You want to connect AI to systems it can't currently reach — your internal ticket tracker, a proprietary database, a vendor API, your company's knowledge base. The Model Context Protocol (MCP) is the standard way to do that: you build a server that exposes tools and data, and every MCP-capable client (Claude Code, claude.ai, and a growing ecosystem) can use it. The engineering questions are about interface design: which operations to expose as tools, how to describe them so a model uses them correctly, and how to keep credentials and blast radius under control.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Not sure what the server should expose | Plan mode | Design the tool surface before writing handlers |
| Ready to implement | Official plugins | `mcp-server-dev` provides a guided server design and build workflow |
| Haven't used MCP as a consumer yet | MCP context | Using existing servers first teaches you what good tools feel like |
| Unfamiliar with the protocol's details | Deep research | Spec, transports, and auth patterns are documented — study before building |
| Server built, needs regression testing | Headless mode | Scripted `claude -p` calls exercise your tools repeatably |

**Hidden gem:** Headless Mode — a shell script of `claude -p` calls against your server is an MCP integration test suite: repeatable, CI-friendly, and it tests the thing that actually matters (can a model use your tools correctly?).

## Approaches (Ranked)

### 1. Plan Mode — Design the tool surface first
**Level:** Beginner

An MCP server is an API designed for a model as the consumer, and models are unusual consumers: they choose tools by reading descriptions, they compose calls in ways you didn't script, and every registered tool costs context. Plan mode turns this into deliberate interface design — which operations, what granularity, what the descriptions must say, what stays read-only.

**Try it now:**
> Enter plan mode. I'm building an MCP server for our internal incident-management system. Candidate operations: search incidents, read incident detail, add a timeline comment, change status, page an on-call. Design the tool surface: which of these should be tools vs. left out, what should each tool's description say so a model picks the right one, which need to be read-only, and where's the damage potential if a model calls something with bad arguments?

**Why this works:** Models select tools by description and affordance — a well-designed five-tool surface outperforms a twenty-tool dump of your REST API, and the design mistakes are cheapest before the first handler exists.

**Pros:**
- Forces the read-only vs. write boundary decision upfront
- Tool descriptions get designed, not retrofitted after confusing the model
- Keeps the surface small — every tool costs context in every client

**Cons:**
- Some usage patterns only reveal themselves once a real model uses the server — expect one revision

**Deeper:** See `approaches/plan-mode.md`

---

### 2. Official Plugins — Build it with mcp-server-dev
**Level:** Intermediate

The `mcp-server-dev` official plugin provides a guided MCP server design and implementation workflow — protocol details, transport setup, and tool registration handled with current best practices instead of reverse-engineered from examples. Pair it with your Phase-1 design and the implementation becomes the easy part.

**Try it now:**
> Install with: /plugin install mcp-server-dev@claude-plugins-official — then use its workflow to implement the incident-management server we designed: five tools, search and read as read-only, comment-adding gated, in TypeScript with stdio transport for local use first.

**Why this works:** MCP has real protocol surface (transports, schemas, capability negotiation) that's undifferentiated work — a guided workflow encodes the current conventions so your effort goes into the tools' semantics.

**Pros:**
- Official, maintained guidance that tracks protocol evolution
- Covers the full path from design through working server
- Complements the `mcp-builder` skill in Anthropic's open skills repo

**Cons:**
- Guided scaffolding still leaves auth and deployment decisions to you

**Deeper:** See `approaches/official-plugins.md`

---

### 3. MCP Context — Consume before you produce
**Level:** Beginner

The fastest way to learn what makes an MCP server good is to use several. Connect existing servers (GitHub, your database, Slack) to your daily workflow for a week: notice which tool descriptions lead the model to the right call, where over-broad tools produce flailing, and how context cost shows up. Also check whether the server you're about to build already exists — the official reference repo and marketplace cover a lot of ground.

**Try it now:**
> Before I build our incident-tool MCP server, connect the GitHub MCP server to this session and watch how I use it: give me a task that needs issue search and PR reads, then critique the experience with me — which tool descriptions worked, where did the model pick the wrong tool, and what does that teach us for our own server's design?

**Why this works:** Interface intuition comes from the consumer side — a week of consuming MCP produces better tool design than a month of producing it blind.

**Pros:**
- Free lessons from servers that already solved your design questions
- Surfaces the "does this already exist?" answer before you build
- You experience context cost and tool-selection failures firsthand

**Cons:**
- Existing servers can anchor you to their patterns even where your domain differs

**Deeper:** See `approaches/mcp-context.md`

---

### 4. Deep Research — Learn the protocol before fighting it
**Level:** Beginner

MCP is a real protocol with a spec: transports (stdio for local, HTTP for remote), tool/resource/prompt primitives, auth patterns, capability negotiation. An hour of research on the spec and a few well-built open-source servers prevents the classic mistakes — wrong transport choice, resources misused as tools, auth bolted on late.

**Try it now:**
> /deep-research I'm building my first MCP server (internal incident-management API, needs auth, will eventually be shared across the org). Current best practices: when to choose stdio vs HTTP transport, how existing production servers handle per-user authentication, tools vs resources for read-mostly data, and what the official spec says about tool descriptions and annotations. Cite the spec and 2-3 well-regarded open-source servers I should read.

**Why this works:** Protocols punish improvisation — the design constraints you don't know about become the rewrite you do six weeks in.

**Pros:**
- Transport and auth decisions made with the spec, not against it
- Good open-source servers double as reference implementations
- Surfaces ecosystem conventions your users will expect

**Cons:**
- The ecosystem moves fast — verify findings against the current spec version

**Deeper:** See `approaches/deep-research.md`

---

### 5. Headless Mode — Regression tests for your tool surface
**Level:** Intermediate

The failure mode that matters for an MCP server isn't "returns wrong JSON" — it's "a model misuses the tools." Headless mode turns that into a testable property: scripted `claude -p` runs with your server attached, exercising realistic tasks, with assertions on the outcomes. Run it in CI and your server's usability by models stops regressing silently.

**Try it now:**
> Write a test script for our incident MCP server: five `claude -p` invocations with `--mcp-config test-mcp.json`, each a realistic task ("find open sev-1 incidents about the payments service and summarize the latest timeline entry"), with `--allowedTools` scoped to the server's tools. Assert on the JSON output: right tool called, no invented incident IDs, summary mentions the actual latest entry.

**Why this works:** "Can a model use this correctly?" is your server's real acceptance criterion — headless runs make it executable instead of anecdotal.

**Pros:**
- Catches description regressions that unit tests can't see
- CI-friendly: runs on every change to tool schemas or descriptions
- Doubles as living documentation of intended usage

**Cons:**
- Model behavior is stochastic — assert on outcomes, not exact call sequences
- Each test run costs tokens; keep the suite focused

**Deeper:** See `approaches/headless-mode.md`
