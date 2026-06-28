# Code Understanding & Exploration
*Last reviewed: 2026-06-27*

## When You're Here

You're staring at a codebase you didn't write and need to make sense of it fast. Maybe you just joined a team and have a ticket due in three days. Maybe you're reviewing a PR that touches a module you've never seen. Or maybe you're the one who wrote it two years ago and your own abstractions have become a mystery. Understanding code is the prerequisite for changing it safely, and AI can compress weeks of passive absorption into hours of active exploration.

The traditional approach — read files, ask a colleague, read more files — works but doesn't scale. A mid-size service has hundreds of files, and the important relationships between them aren't visible from any single file. AI excels here because it can hold multiple files in context simultaneously and trace connections that would take you hours of tab-switching to piece together.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Onboarding to a new team or project | Plan mode | Structured exploration builds a mental map methodically |
| Architecture docs exist but are scattered or outdated | MCP context | Pull docs from Notion/Confluence and cross-reference with actual code |
| Codebase uses unfamiliar frameworks or patterns | Deep research | Understand the framework's conventions before reading the code |
| Reviewing a PR in an unfamiliar module | Plan mode | Map the module's structure before judging the change |
| Legacy code with no documentation or original authors | Plan mode + MCP context | Combine code analysis with whatever written history exists |
| Large codebase — need to trace how components connect | LSP self-correction | Compiler-precise "go to definition" and "find references" |

## Approaches (Ranked)

### 1. Plan Mode — Explore with a map, not a flashlight
**Level:** Beginner | **Tools:** Any

When you're new to a codebase, the temptation is to start reading files randomly — `main.go`, then whatever it imports, then whatever catches your eye. Plan mode imposes structure: start with entry points, trace data flow, identify module boundaries, and build a layered understanding. Claude reads the code systematically and produces an architecture summary you can refer back to — a map that makes every future exploration faster.

**Try it now:**
> Enter plan mode. I just joined the team and need to understand the `services/payment-gateway/` module. Start from the entry points — find the main route handlers or event consumers. Trace a single payment through the system end-to-end: from the initial API call through validation, fraud checks, processor integration, and persistence. Identify the key abstractions, external service calls, error handling patterns, and retry logic. What are the failure modes? Produce an architecture summary with a dependency diagram I can reference during code reviews.

**Why this works:** Understanding a codebase is not about reading every file — it's about building the right mental model. Structured exploration prioritizes the connections between components (who calls whom, what data flows where) over the details within them. This is how experienced engineers onboard quickly: they learn the shape of the system first, then fill in details as needed.

**Pros:**
- Builds a reusable mental model, not just one-off answers
- Identifies the critical paths through the system first
- Produces documentation you can share with the next new hire

**Cons:**
- The architecture summary is a snapshot — it drifts as the code evolves
- May miss undocumented conventions that only exist in tribal knowledge

**Deeper:** See `approaches/plan-mode.md`

---

### 2. MCP Context — Reunite code with its documentation
**Level:** Intermediate | **Tools:** Claude Code (with MCP)

Most mature codebases have documentation — it's just scattered across Confluence pages nobody bookmarked, Notion docs from a team that reorged, Google Docs shared in a Slack thread you weren't in, and README files that haven't been updated since 2022. MCP context servers let Claude pull these documents alongside the code, cross-referencing design decisions with their implementation and flagging where the docs have drifted from reality. This turns "archaeology" into "analysis."

**Try it now:**
> Pull the architecture decision records from our Confluence space "Platform Architecture" and the original design doc for the event-driven migration from Notion page "Event Bus RFC." Cross-reference them with the current implementation in `src/events/`. Which decisions from the RFC were actually implemented? Which were changed? Are there any patterns in the code that aren't documented anywhere? Produce a "documentation drift" report I can share with the team.

**Why this works:** Code tells you what the system does. Documentation tells you why it was built that way. Understanding requires both. MCP context bridges the gap by pulling documentation into the same conversation as the code, so Claude can answer "why is it built this way?" and "what alternatives were considered?" instead of just "what does this function do?"

**Pros:**
- Surfaces design intent that's invisible in code alone
- Identifies where documentation has drifted from implementation
- Brings scattered knowledge into a single, queryable context

**Cons:**
- Requires MCP server setup for each documentation source
- Only as useful as the existing documentation — can't surface what was never written
- May hit context limits with large documentation sets

**Deeper:** See `approaches/mcp-context.md`

---

### 3. Deep Research — Understand the framework before the code
**Level:** Beginner | **Tools:** Claude Code

When a codebase uses frameworks, patterns, or conventions you're unfamiliar with, reading the application code first is working backwards. You'll confuse framework conventions with business logic and miss the architectural intent entirely. Deep research lets you understand the framework's idioms — dependency injection containers, middleware chains, ORM conventions, decorator patterns — so that when you read the application code, you instantly recognize what's scaffolding versus what's the team's actual work.

**Try it now:**
> /deep-research This codebase uses NestJS with CQRS (command-query responsibility segregation) and the `@nestjs/cqrs` module. I'm seeing classes like `CreateOrderCommand`, `OrderCreatedEvent`, and `GetOrderQuery`. Explain the CQRS pattern as implemented in NestJS: how do commands, events, queries, and their handlers connect? What's the lifecycle of a request through this pattern? Include the NestJS-specific conventions so I can read this codebase fluently.

**Why this works:** Frameworks encode architectural decisions into conventions. If you don't understand the conventions, every file looks arbitrary. Deep research gives you the "grammar" of the framework so you can read the codebase's "sentences" fluently.

**Pros:**
- Accelerates onboarding to unfamiliar tech stacks dramatically
- Distinguishes framework patterns from application-specific logic
- Sources from official docs, community guides, and real-world examples

**Cons:**
- Doesn't cover project-specific customizations of the framework
- May surface outdated information for rapidly evolving frameworks

**Deeper:** See `approaches/deep-research.md`

---

### 4. LSP Self-Correction — Navigate definitions and references like an IDE
**Level:** Intermediate | **Tools:** Claude Code / OpenCode

When exploring unfamiliar code, the most powerful operations are "go to definition" and "find all references." LSP self-correction gives Claude access to the same navigation a developer uses in their IDE — tracing function calls to their implementations, finding every caller of a method, and resolving type hierarchies. This is especially valuable for understanding how components connect across a large codebase without reading every file.

**Try it now:**
> I need to understand how authentication works in this codebase. Start from the `requireAuth` middleware in `src/middleware/auth.ts`. Use LSP to find every route that uses this middleware, then go to the definition of the `validateToken()` function it calls. Trace the token validation through to the session store. Map the full auth chain from HTTP request to validated user object.

**Why this works:** Code understanding is fundamentally about tracing connections — who calls what, where does data come from, what types flow through the system. LSP provides precise, compiler-backed answers to these questions instead of relying on text search, which misses indirect calls and aliased imports.

**Pros:**
- Compiler-precise navigation — no false positives from text search
- Traces type hierarchies and interface implementations automatically
- Works across files and modules without reading everything

**Cons:**
- Requires LSP server to be configured and running
- Less useful for dynamically typed languages where LSP has limited information

**Deeper:** See `approaches/lsp-self-correction.md`
