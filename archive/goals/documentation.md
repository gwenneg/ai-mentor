# Documentation
*Last verified: 2026-06-27*

## When You're Here

You need to write or improve documentation — API references, architecture decision records, onboarding guides, or runbooks. Maybe the existing docs are outdated and misleading, maybe there are no docs at all and the team's knowledge lives in Slack threads and tribal memory. Writing good docs is one of the highest-leverage things an engineer can do, and also one of the things engineers most consistently put off.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Existing docs, specs, and decisions are scattered | MCP context | Pull everything into one place for Claude |
| Writing docs for an unfamiliar domain or standard | Deep research | Get the conventions right before writing |
| Large doc set with multiple sections and audiences | Plan mode | Structure before prose prevents rewrites |
| Updating API docs after code changes | MCP context | Pull the actual code and existing docs together |
| Creating a new onboarding guide from scratch | Deep research + Plan mode | Research what good onboarding looks like, then outline |
| Documenting multiple modules or services in parallel | Beyond the Catalog (see /ai-mentor) | Fan-out or subagent approaches can parallelize doc generation |
| Docs need diagrams, dashboards, or sharing beyond the repo | Visual artifacts | A rendered page communicates structure that linear text can't |

**Hidden gem:** Custom Skills — a `/gen-api-doc` command makes regenerating docs cheaper than letting them drift.

## Approaches (Ranked)

### 1. MCP Context — Pull existing docs, API specs, and design decisions
**Level:** Intermediate

The biggest obstacle to good documentation is scattered context. The API behavior lives in code, the design rationale lives in a Notion page, the edge cases live in Jira tickets, and the deployment requirements live in someone's head. MCP context lets you pull all of these sources into Claude's working memory so it can synthesize accurate, complete documentation instead of guessing.

**Try it now:**
> Connect to our GitHub repo and read `src/api/routes/payments.ts` and `src/api/middleware/auth.ts`. Also fetch the OpenAPI spec at `docs/openapi.yaml` and the original design doc from Notion (linked in `docs/DESIGN_DECISIONS.md`). Using all of this context, generate API reference documentation for the payments endpoints, including authentication requirements, request/response schemas with examples, error codes, and rate limits.

**Why this works:** Documentation quality is directly proportional to context quality. AI can write excellent prose, but only if it has access to the truth. MCP bridges the gap between where knowledge lives (code, tickets, wikis) and where it needs to go (documentation).

**Pros:**
- Produces documentation that reflects actual behavior, not assumptions
- Can cross-reference multiple sources to catch inconsistencies
- Keeps docs in sync with code by re-running with updated context

**Cons:**
- Requires MCP connections to be configured for your tools (GitHub, Notion, Jira, etc.)
- Garbage in, garbage out — if the source material is wrong, the docs will be too

**Deeper:** See `approaches/mcp-context.md`

---

### 2. Deep Research — Research documentation standards and best practices
**Level:** Beginner

Before writing a single word, understand what good looks like for your type of documentation. Deep research can survey how leading projects document their APIs, what an effective architecture decision record contains, or what onboarding guides actually help new engineers ramp up. This prevents the common failure mode of writing docs that are technically accurate but structurally useless — walls of text that no one reads because they weren't designed for the reader.

**Try it now:**
> /deep-research What are the best practices for writing API reference documentation for a REST API? I want to document our internal payments API (15 endpoints, auth via OAuth2, versioned via URL path). Compare approaches from Stripe, Twilio, and GitHub's API docs. What structure, level of detail, and examples make API docs actually useful vs. just checkbox compliance?

**Why this works:** Documentation is a communication design problem, not just a writing task. Researching proven patterns gives you a structural blueprint so you spend your effort on content, not on figuring out how to organize it.

**Pros:**
- Prevents "blank page paralysis" by giving you a proven structure to follow
- Surfaces best practices you wouldn't discover by staring at your own codebase
- Helps you match industry standards your users already expect

**Cons:**
- Can lead to over-engineering simple docs — a README doesn't need the Stripe treatment
- Research takes 5-10 minutes; skip it for small, internal-only docs

**Deeper:** See `approaches/deep-research.md`

---

### 3. Plan Mode — Structured outline before writing
**Level:** Beginner

For documentation larger than a single page — onboarding guides, architecture overviews, runbook collections — writing without an outline leads to rambling, redundancy, and gaps. Plan mode produces a detailed outline that defines the audience, scope, structure, and key points for each section before any prose is written. This is especially valuable when multiple people will contribute to the docs.

**Try it now:**
> Enter plan mode. I need to create an onboarding guide for new engineers joining our platform team. The tech stack is Go microservices, PostgreSQL, Kafka, Kubernetes on AWS. They need to set up their local dev environment, understand our service architecture (12 services), learn our deployment pipeline, and know who to ask about what. Create a detailed outline with sections, subsections, and 2-3 bullet points per subsection describing what content goes there. Flag any sections where I'll need input from specific team members.

**Why this works:** Documentation is information architecture. A well-structured outline ensures every piece of information has exactly one home, readers can navigate to what they need, and writers know exactly what to cover without overlap or gaps.

**Pros:**
- Prevents the "stream of consciousness" anti-pattern in long docs
- Makes it easy to divide writing work across team members
- Catches structural problems (missing sections, wrong audience assumptions) cheaply

**Cons:**
- Adds a planning step that feels unnecessary for short docs
- The outline itself needs review — a bad structure propagates into bad docs

**Deeper:** See `approaches/plan-mode.md`

---

### 4. Custom Skills — Repeatable doc generation on demand
**Level:** Intermediate

When you generate the same type of documentation repeatedly — API references after endpoint changes, changelog entries for releases, module READMEs for new packages — a custom skill turns the pattern into a single command. `/gen-api-doc payments` reads the source code, pulls the OpenAPI spec, and produces a formatted API reference. `/update-changelog` reads recent commits and generates a categorized changelog entry. The skill encodes your team's documentation standards so every output is consistent.

**Try it now:**
> Create a custom skill at `.claude/skills/gen-api-doc.md`. When I run `/gen-api-doc billing`, it should: read all route handlers in `src/api/routes/billing/`, extract endpoint paths, methods, request/response schemas, and auth requirements, then generate a markdown API reference following the format in `docs/api/payments.md`. Include request examples, response examples, and error codes for each endpoint.

**Why this works:** Documentation that requires the same structure every time is a solved problem once you encode the pattern. Custom skills eliminate the "copy the last API doc and update it" workflow, producing complete, consistent output from the source code every time.

**Pros:**
- Consistent doc format across the entire project
- One command produces a complete, structured document
- New endpoints get documented the same way as existing ones

**Cons:**
- Requires upfront effort to define the skill and template
- Skills need updating when documentation standards evolve

**Deeper:** See `approaches/custom-skills.md`

---

### 5. Visual Artifacts — Publish docs as a rendered, shareable page
**Level:** Beginner

Some documentation is spatial: architecture overviews, onboarding maps, request-flow diagrams. The built-in Artifact tool renders an HTML or Markdown page to a private claude.ai URL you can open in a browser, iterate on conversationally, and share with teammates who were never in the session. The repo keeps the Markdown source of truth; the artifact is the readable, shareable view of it.

**Try it now:**
> Read `docs/architecture.md` and the service directories under `src/services/`. Publish an artifact: an architecture overview with a dependency diagram at the top, one card per service (purpose, entry points, key files), and a "gotchas" section. Keep the Markdown source in `docs/` — the artifact is the rendered view I'll share with the team.

**Why this works:** A document's format should match its structure. Diagrams and cards carry relationships that linear Markdown flattens, and a stable URL turns documentation from a file people must find into a link people actually open.

**Pros:**
- Layout and diagrams communicate structure that prose can't
- A stable link survives the session and is shareable with the whole team
- Iterating is conversational — the same URL updates on each revision

**Cons:**
- Publishing uploads content to claude.ai hosting — check policy for sensitive material
- The page is a view, not the source of truth — keep the Markdown in the repo

**Deeper:** See `approaches/visual-artifacts.md`

---

## When to Skip the Approaches

Unlike debugging or migration, documentation is often a category where the simplest approach works well. If you just need to document a function, a module, or a small feature, skip the ranked approaches above and simply ask Claude with good context:

> Read `src/billing/invoice_generator.py` and write a docstring for each public method. Include parameter types, return types, a one-line summary, and one usage example per method. Follow Google's Python docstring style.

> Read `src/middleware/rate_limiter.go` and the existing comment block at the top. Rewrite it as a package-level doc comment that explains what the middleware does, its configuration options, and one example of how to mount it on a router.

The ranked approaches above are most valuable when:
- Your documentation task spans multiple files or systems
- You're writing for an audience whose needs you don't fully understand yet
- The information you need to document is scattered across code, tickets, and tribal knowledge
- You need to match an external standard or convention you haven't used before

For small, well-scoped documentation tasks, direct prompting with the relevant code in context is the fastest path. The art is knowing which kind of task you're facing.
