# Research & Investigation
*Last verified: 2026-06-27*

## When You're Here

You need to make a decision and the cost of getting it wrong is high. Maybe you're choosing between message brokers for a new service, evaluating whether to build or buy an auth solution, or investigating whether a migration to a new framework is worth the effort. Research isn't about finding one right answer — it's about mapping the tradeoff space so you can make a defensible choice. AI turns what used to be days of tab-hopping into a structured investigation with cited sources.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Comparing libraries, frameworks, or SaaS tools | Deep research | Fan-out searches surface tradeoffs across docs, issues, and benchmarks |
| Need a structured evaluation framework first | Plan mode | Define your criteria before gathering evidence to avoid confirmation bias |
| Evaluating tools with interactive demos or UIs | Browser integration | Reading docs isn't enough — test the actual experience |
| Past decisions or prior art exist inside the company | MCP context | Check what was already decided before researching from scratch |
| Different phases need different reasoning depth | Model & effort selection | Opus at high effort for analysis, Haiku subagents for summarization |

**Hidden gem:** Plan Mode — defining evaluation criteria before gathering evidence is the only reliable defense against confirmation bias, and nobody thinks of plan mode for research.

## Approaches (Ranked)

### 1. Deep Research — Fan out, verify, synthesize
**Level:** Beginner

Deep research is purpose-built for investigation. It fans out across multiple web sources — documentation, GitHub issues, benchmark reports, community discussions, release changelogs — then cross-references findings to filter out outdated or biased information. The output is a cited report with sources you can verify, not an opinion dressed up as fact.

**Try it now:**
> /deep-research We're evaluating message brokers for our event-driven architecture. Compare RabbitMQ, Apache Kafka, and NATS JetStream for our use case: ~50K events/sec peak, exactly-once delivery required for payment events, at-least-once acceptable for analytics, team has Go and Python services. Compare on: throughput, delivery guarantees, operational complexity, client library maturity in Go and Python, and managed hosting options on AWS.

**Why this works:** Good technical decisions require evidence from multiple angles — performance, operational burden, ecosystem maturity, community health. Deep research automates the "open 30 browser tabs" workflow and synthesizes findings into a structure you can present to your team.

**Pros:**
- Produces cited reports with verifiable sources
- Cross-references multiple sources to filter out bias and outdated information
- Covers dimensions you might not think to check (licensing, maintenance cadence, CVE history)

**Cons:**
- Quality depends on publicly available information — poorly documented tools get thin coverage
- Can overwhelm with information if the scope is too broad — be specific about your constraints

**Deeper:** See `approaches/deep-research.md`

---

### 2. Plan Mode — Define the question before seeking the answer
**Level:** Beginner

Before researching, you need to know what you're researching. Plan mode helps you define evaluation criteria, weight them by importance, and structure the investigation so you're gathering evidence against a framework rather than collecting random facts. This prevents the most common research failure: confirmation bias, where you unconsciously seek evidence for the option you already prefer.

**Try it now:**
> Enter plan mode. We need to decide whether to build a custom permissions system or adopt an existing solution like Oso, Cerbos, or OpenFGA. Before we research any of them, help me define the evaluation framework. What dimensions matter for our use case? We have: 12 microservices in Go, ~200 permission rules, multi-tenant with tenant-scoped roles, and need sub-10ms policy evaluation. Rank the evaluation criteria by importance and define what "good enough" looks like for each.

**Why this works:** Research without a framework is just browsing. By defining evaluation criteria first, you create a decision matrix that forces objective comparison. This is how staff engineers make architecture decisions — criteria first, evidence second, decision last.

**Pros:**
- Prevents confirmation bias by defining criteria before gathering evidence
- Creates a shareable decision framework the team can weigh in on
- Makes the final decision defensible because the reasoning is explicit

**Cons:**
- Adds a step before the "real" research begins — requires patience
- Criteria can be hard to define without some preliminary exploration

**Deeper:** See `approaches/plan-mode.md`

---

### 3. Browser Integration — Test it, don't just read about it
**Level:** Advanced

Documentation describes the ideal. Demos reveal the reality. Browser integration lets Claude navigate to interactive demos, documentation sites, and admin UIs so you can evaluate the actual developer experience. This is critical for tools where the UX matters — dashboards, CLI tools with web UIs, API explorers, and configuration interfaces.

**Try it now:**
> Connect to the browser and compare the API documentation experience for three auth providers: navigate to `docs.auth0.com`, `clerk.com/docs`, and `supabase.com/docs/guides/auth`. For each, check: is the quickstart for Go clear and complete? Are the API references searchable? Is there a working example for multi-tenant role-based access? Screenshot each provider's quickstart page and rate the documentation quality.

**Why this works:** Developer tools are chosen partly on capability and partly on experience. A library with perfect features but confusing docs will slow your team down more than a slightly less capable one with excellent documentation. Browser integration lets you evaluate the experience, not just the feature list.

**Pros:**
- Evaluates the real developer experience, not marketing claims
- Can test interactive demos and playgrounds
- Screenshots create shareable evidence for team discussions

**Cons:**
- Requires browser MCP setup
- Slower than reading cached documentation
- Only useful for tools with web-accessible documentation or demos

**Deeper:** See `approaches/browser-integration.md`

---

### 4. MCP Context — Check internal prior art first
**Level:** Intermediate

Before researching externally, check what your organization already knows. MCP context servers can pull past architecture decision records, previous spikes, Slack discussions, and internal wiki pages. Maybe another team already evaluated these options. Maybe there's an approved vendor list. Starting with internal knowledge prevents you from duplicating work and ensures your recommendation aligns with organizational constraints.

**Try it now:**
> Pull architecture decision records from Confluence space "Platform Decisions" and search Slack channel #platform-engineering for discussions about message brokers or event streaming from the last 6 months. Has any team already evaluated Kafka vs NATS? Are there internal guidelines or approved infrastructure that would constrain our choice?

**Why this works:** Organizations accumulate institutional knowledge that's invisible if you only search the public web. Checking internal sources first can save days of research — and prevents the embarrassment of recommending something another team already tried and rejected.

**Pros:**
- Prevents duplicated research effort across teams
- Surfaces organizational constraints (approved vendors, existing infrastructure)
- Incorporates lessons learned from past decisions

**Cons:**
- Requires MCP server setup for internal tools (Confluence, Slack, etc.)
- Internal documentation may be outdated or incomplete
- Not useful if your organization doesn't document decisions

**Deeper:** See `approaches/mcp-context.md`

---

### 5. Model & Effort Selection — Match reasoning depth to each research phase
**Level:** Advanced

Research has phases with very different reasoning demands: analyzing architectural tradeoffs needs maximum depth, while summarizing twenty GitHub issues is volume work. Model & effort selection lets you run the analysis on Opus at high effort and delegate the summarization to cheap, fast Haiku subagents — deeper conclusions where it matters, lower cost where it doesn't.

**Try it now:**
> Switch to Opus and set /effort high. Analyze the architectural tradeoffs between event sourcing and traditional CRUD for our order management system (~10K orders/day, 3-person team). Then spawn Haiku subagents to summarize the top 5 GitHub issues for each major event sourcing library in our language, and fold their summaries into the final recommendation.

**Why this works:** Reasoning depth is a budget. Spending it uniformly means the hard analysis gets less than it needs while routine summarization gets more than it can use. Allocating depth per phase gets better conclusions and a faster, cheaper session at the same time.

**Pros:**
- Deep reasoning goes where the decision actually gets made
- Haiku subagents parallelize the reading grunt work at low cost
- Second-pass review at high effort catches overconfident claims

**Cons:**
- Requires judgment about which phase deserves which depth
- Diminishing returns for straightforward research questions

**Deeper:** See `approaches/model-effort-selection.md`
