# API Design
*Last verified: 2026-06-28*

## When You're Here

You're designing an API — REST, GraphQL, or gRPC — and the critical decisions aren't about handler code, they're about the contract. Endpoint naming, resource modeling, versioning strategy, error shapes, pagination, and how this API evolves over the next two years without breaking the consumers who depend on it today. These decisions are cheap to change in a design doc and ruinously expensive to change after clients are in production.

API design is distinct from building a feature (which is broader) and from writing documentation (which describes an API that already exists). This is about making the design decisions themselves: what resources exist, how they relate, what operations are allowed, and what the response envelope looks like before anyone writes a handler.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Starting a new API or major version from scratch | Plan mode | Map every resource, relationship, and error contract before writing handlers |
| Unsure about conventions (pagination, versioning, error format) | Deep research | Learn from APIs that millions of developers already use |
| Extending an existing API and need to stay consistent | Built-in review skills | Catch naming inconsistencies and breaking changes before merge |
| Requirements scattered across specs, tickets, and existing APIs | MCP context | Pull everything together so the design reflects reality |
| Your team creates new endpoints frequently | Custom skills | Encode your conventions into a repeatable scaffold |

**Hidden gem:** Built-in Review Skills — pointing `/code-review` at only the API surface catches breaking changes and naming drift that no test suite will.

## Approaches (Ranked)

### 1. Plan Mode — Design the contract before writing handlers
**Level:** Beginner

API design is contract design. Plan mode forces you to map out every endpoint, request/response schema, versioning strategy, error code, and pagination approach before writing a single handler. Claude walks through the resource model, identifies relationships you haven't considered, and surfaces edge cases — what happens on a partial update? How do you handle soft deletes in list responses? What's the pagination strategy when the underlying data changes between pages?

**Try it now:**
> Enter plan mode. I'm designing a REST API for a multi-tenant project management system. Resources: projects, tasks, comments, and team members. I need endpoints for CRUD on each resource, plus task assignment and status transitions. Design the full API surface: URL structure, HTTP methods, request/response schemas, error codes, and versioning strategy. Consider: how do we scope all queries to a tenant? What's the pagination approach for task lists that can have 10K+ items? How do we handle bulk operations (assign 50 tasks at once)? Don't write code — give me the API contract to review.

**Why this works:** An API is a promise to your consumers. Changing a promise after people depend on it requires deprecation cycles, migration guides, and apology emails. Planning the contract first lets you find the design mistakes while fixing them is still free.

**Pros:**
- Surfaces resource relationship problems before they're encoded in URLs
- Forces explicit decisions on versioning, pagination, and error shapes
- Produces a reviewable contract that backend and frontend teams can align on

**Cons:**
- Feels slow when you think you already know the API shape
- Plans need revisiting as you discover implementation constraints

**Deeper:** See `approaches/plan-mode.md`

---

### 2. Deep Research — Learn from APIs developers actually love
**Level:** Beginner

Before inventing your own conventions, study the APIs that millions of developers use daily. Deep research surveys how Stripe handles idempotency keys, how GitHub structures nested resources, how Twilio versions its API, and how GraphQL schemas handle federation. It also covers standards you should know about — JSON:API, OpenAPI, gRPC service definition best practices, REST maturity levels — so you're making informed choices rather than guessing.

**Try it now:**
> /deep-research I'm designing a REST API that needs to support versioning, and I'm torn between URL path versioning (`/v2/users`), header versioning (`Accept: application/vnd.api.v2+json`), and query parameter versioning. Compare how Stripe, GitHub, Twilio, and Google Cloud handle API versioning. What are the tradeoffs for developer experience, caching, and long-term maintenance? I also need a strategy for deprecating old versions — how do these companies communicate and enforce sunset timelines?

**Why this works:** API design is full of decisions that seem arbitrary until you hit scale. Researching how proven APIs handle these decisions gives you rationale for your choices and helps you avoid pitfalls that others have already discovered and documented.

**Pros:**
- Prevents reinventing conventions that already have industry consensus
- Surfaces tradeoffs you wouldn't discover until consumers complain
- Gives you a decision record backed by real-world precedent

**Cons:**
- Can lead to over-engineering a simple internal API with patterns meant for public APIs
- Not every convention from Stripe applies to your domain

**Deeper:** See `approaches/deep-research.md`

---

### 3. MCP Context — Pull requirements from specs and existing APIs
**Level:** Intermediate

Good API design doesn't happen in a vacuum. You need the product requirements (what operations do users need?), the existing API surface (what conventions are already established?), and any OpenAPI specs or protobuf definitions that define the current contract. MCP context servers let Claude pull all of this into its working memory — your existing OpenAPI spec, the product requirements from Linear, the frontend team's wishlist from a Notion doc — so the new API design is grounded in what actually exists and what's actually needed.

**Try it now:**
> Connect to our GitHub repo and read the existing OpenAPI spec at `docs/openapi.yaml` and the current route handlers in `src/api/routes/`. Also pull the product requirements from LINEAR-2847 "Billing API v2." I need to design new billing endpoints that follow the same conventions as our existing API (naming, error format, pagination style) while meeting the requirements in the ticket. Identify any conflicts between what the ticket asks for and what our current conventions support.

**Why this works:** API inconsistency is the #1 complaint from API consumers. By pulling in your existing API surface alongside new requirements, Claude can ensure the new endpoints feel like they belong in the same API — same naming conventions, same error envelope, same pagination style.

**Pros:**
- Ensures new endpoints match existing API conventions automatically
- Catches conflicts between requirements and current API patterns early
- Grounds the design in real specs rather than assumptions about what exists

**Cons:**
- Requires MCP connections to your spec hosting and project management tools
- Existing API conventions may themselves be inconsistent, propagating bad patterns

**Deeper:** See `approaches/mcp-context.md`

---

### 4. Built-in Review Skills — Review API changes for consistency and breaking changes
**Level:** Beginner

Once you have API code, `/code-review` can focus specifically on the API surface: are field names consistent with the rest of the API? Are there error cases the handler doesn't cover? Did you add a required field to an existing response (breaking change)? Are query parameters documented? This is especially valuable before merging, when catching a naming inconsistency is a one-line fix rather than a deprecation cycle.

**Try it now:**
> /code-review Focus on the API surface of this diff. Check for: naming consistency with our existing endpoints in `src/api/routes/`, missing error responses (what if the resource doesn't exist? what if the user lacks permission?), undocumented query parameters, response fields that don't match our OpenAPI spec at `docs/openapi.yaml`, and any breaking changes to existing response shapes. Ignore implementation details — I only care about what consumers see.

**Why this works:** API bugs are different from implementation bugs. A misspelled field name or an inconsistent error format won't break any test, but it will frustrate every consumer who integrates with your API. Review skills catch these surface-level issues that are invisible to type checkers and test suites.

**Pros:**
- Catches breaking changes before they reach consumers
- Enforces naming and format consistency across endpoints
- Quick feedback loop — runs in seconds on a diff

**Cons:**
- Only catches issues visible in the diff, not systemic design problems
- Needs existing API context to compare against for consistency checks

**Deeper:** See `approaches/built-in-review-skills.md`

---

### 5. Custom Skills — Repeatable API scaffolding for your project
**Level:** Intermediate

If your team adds new endpoints regularly, a custom skill encodes your API conventions into a single command. `/new-endpoint POST /api/v1/invoices` generates the route handler, request validator, response serializer, route registration, OpenAPI spec entry, and test stub — all following your team's established patterns. No more copying an existing endpoint and doing find-and-replace, no more forgetting to register the route or update the spec.

**Try it now:**
> Create a custom skill at `.claude/skills/new-endpoint.md`. When I run `/new-endpoint POST /api/v1/invoices`, it should: create a handler in `src/api/routes/invoices.ts` following the pattern in `src/api/routes/payments.ts`, add a Zod validator in `src/api/validators/invoices.ts`, register the route in `src/api/router.ts`, add the endpoint to `docs/openapi.yaml`, and create a test stub in `tests/api/invoices.test.ts`. Follow our existing naming conventions and error handling patterns.

**Why this works:** API consistency comes from repetition, and repetition is what skills automate. When every endpoint is scaffolded from the same template, consumers get a predictable API surface and new team members produce endpoints that look like they were written by the same person.

**Pros:**
- Every new endpoint follows the same conventions automatically
- Eliminates forgotten steps (route registration, spec updates, test stubs)
- New team members produce consistent API surfaces from day one

**Cons:**
- Upfront effort to define the skill and keep it current with evolving conventions
- May need multiple skills for different endpoint patterns (CRUD vs. RPC-style actions)

**Deeper:** See `approaches/custom-skills.md`
