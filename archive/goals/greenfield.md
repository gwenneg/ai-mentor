# Greenfield / Feature Development
*Last verified: 2026-06-27*

## When You're Here

You're building something new. Maybe it's a feature from scratch, a new service, or a proof of concept that needs to become production-ready. The blank canvas is exciting but also where most AI-assisted mistakes happen — generating code without a plan leads to architectures you'll regret in three months. The key is using AI to think first, then build fast.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Starting a new service or module from scratch | Plan mode | Architecture decisions made early compound; get them right |
| Need to choose between libraries or frameworks | Deep research | Avoid picking a dependency you'll rip out in 6 months |
| Feature is well-defined and you want speed | Autonomous loops | Set completion criteria and let Claude build iteratively |
| Building a user-facing feature with UI | Browser integration | Visual feedback catches layout and UX issues during development |
| Requirements live in Jira/Linear/Notion | MCP context | Pull specs directly so nothing gets lost in translation |

**Hidden gem:** Custom Skills — encoding your conventions as a scaffold command before writing feature #2 pays off for every feature after it.

## Approaches (Ranked)

### 1. Plan Mode — Architect before you code
**Level:** Beginner

When building something new, the most expensive mistakes are structural. Plan mode forces you to define the architecture — data models, API boundaries, module responsibilities — before a single line of production code exists. Claude maps out the system, identifies edge cases you haven't considered, and produces a blueprint you can critique before committing to it.

**Try it now:**
> Enter plan mode. I'm building a notification service for our platform. It needs to support email, Slack, and in-app notifications with user-configurable preferences. The service will consume events from a Kafka topic. Design the module structure, data models, and API surface. Consider: how do we handle delivery failures and retries? What about rate limiting per channel? Don't write code — give me the architecture to review.

**Why this works:** In greenfield development, code is cheap but architecture is expensive to change. Planning forces you to make structural decisions explicitly rather than letting them emerge accidentally from whichever file Claude writes first.

**Pros:**
- Prevents the "rewrite it next quarter" outcome
- Surfaces edge cases and integration points early
- Creates a shared reference for the whole team

**Cons:**
- Feels like overhead when you just want to start coding
- Plans can become stale if requirements shift mid-build

**Deeper:** See `approaches/plan-mode.md`

---

### 2. Deep Research — Know the landscape before you build on it
**Level:** Beginner

Before committing to a framework, library, or architectural pattern, deep research fans out across documentation, GitHub issues, benchmarks, and community discussions to surface the tradeoffs you'd otherwise discover six months in. This is especially valuable for decisions that are costly to reverse — ORMs, message brokers, auth providers.

**Try it now:**
> /deep-research I'm building a real-time collaboration feature (think Google Docs-style concurrent editing) for our React app. Compare CRDTs vs OT for conflict resolution. Which libraries are production-ready in the JS/TS ecosystem? I need something that works with our PostgreSQL backend and can handle ~500 concurrent users per document.

**Why this works:** Technology choices in greenfield projects have compounding consequences. Spending 20 minutes on research now avoids the "we chose the wrong library" realization after 10,000 lines of integration code.

**Pros:**
- Surfaces maintenance risks, performance limits, and licensing issues early
- Cross-references multiple sources so you're not relying on one blog post
- Produces a decision record you can share with your team

**Cons:**
- Can lead to analysis paralysis if you research instead of building
- Research quality depends on how well the ecosystem is documented

**Deeper:** See `approaches/deep-research.md`

---

### 3. Autonomous Loops — Define "done" and let Claude build
**Level:** Intermediate

Once you have a plan and have chosen your tools, autonomous loops let you set a feature completion goal and let Claude iterate toward it. You define the acceptance criteria — "the endpoint handles these request shapes, validates input, returns these response codes, and all tests pass" — and Claude writes code, runs tests, fixes failures, and repeats until the criteria are met.

**Try it now:**
> /goal: Implement the `POST /api/v1/notifications/preferences` endpoint in `src/routes/notifications.ts`. It should accept a JSON body with `{ userId, channels: { email: bool, slack: bool, inApp: bool }, quietHours: { start: "HH:MM", end: "HH:MM" } }`. Validate all fields, store in the `notification_preferences` table, return 201 on create and 200 on update. Write integration tests in `tests/notifications.test.ts` covering valid input, missing fields, and invalid time formats. All tests must pass.

**Why this works:** Feature development involves tight build-test-fix cycles. Autonomous loops excel at this mechanical iteration, letting you focus on reviewing the output rather than typing every line. The key is setting clear, testable criteria so Claude knows when it's done.

**Pros:**
- Dramatically accelerates boilerplate-heavy feature work
- Self-verifies through test execution
- Frees you to review and think architecturally while Claude codes

**Cons:**
- Without a prior plan, can produce technically correct but architecturally poor code
- May over-engineer or under-engineer without clear scope boundaries

**Deeper:** See `approaches/autonomous-loops.md`

---

### 4. Browser Integration — Build what users actually see
**Level:** Advanced

For user-facing features, browser integration lets Claude render the UI as it builds, catching visual issues that no unit test would find. Claude can navigate your running app, verify that components render correctly, test responsive layouts, and confirm that user flows work end-to-end. This is the difference between "the code works" and "the feature works."

**Try it now:**
> Connect to the browser at `localhost:3000`. I just implemented the notification preferences page at `/settings/notifications`. Check that: the toggle switches for each channel render correctly, the quiet hours time pickers accept valid times and reject invalid ones, the save button shows a loading state during submission, and the success toast appears after saving. Screenshot any visual issues you find.

**Why this works:** Users interact with pixels, not code. Browser integration closes the gap between what your code does and what your users experience, catching CSS conflicts, z-index issues, and interaction bugs that only manifest in a real browser.

**Pros:**
- Catches visual regressions during development, not after deploy
- Can verify entire user flows, not just individual components
- Screenshots create a visual record of what was tested

**Cons:**
- Requires browser MCP setup and a running dev server
- Slower iteration cycle than code-only development
- Only relevant for features with a UI component

**Deeper:** See `approaches/browser-integration.md`

---

### 5. MCP Context — Build from the spec, not from memory
**Level:** Intermediate

Requirements scattered across Jira tickets, Notion docs, and Slack threads are the #1 cause of features that miss the mark. MCP context servers let Claude pull requirements directly from your project management tools, ensuring that acceptance criteria, design mockups, and stakeholder decisions are part of the prompt — not paraphrased from your memory of a meeting last Tuesday.

**Try it now:**
> Pull the requirements from LINEAR-4521 "User notification preferences" and the design spec from our Notion page "Notifications V2 Design." Cross-reference them: are there any requirements in the ticket that the design doesn't cover? Any design elements that aren't mentioned in the ticket? List the gaps before we start building.

**Why this works:** The biggest risk in greenfield development isn't bad code — it's building the wrong thing. By pulling requirements directly from the source of truth, you eliminate the telephone game between PM, designer, and developer.

**Pros:**
- Ensures nothing gets lost between spec and implementation
- Catches requirement gaps before you write code
- Keeps the AI grounded in actual requirements, not hallucinated ones

**Cons:**
- Requires MCP server setup for each tool (Jira, Linear, Notion, etc.)
- Only as good as the quality of the source documents

**Deeper:** See `approaches/mcp-context.md`

---

### 6. Custom Skills — Repeatable scaffolding for your project patterns
**Level:** Advanced

If your team creates the same structure every time — a new React component with a test file, a Storybook story, and an index barrel export, or a new API endpoint with handler, validator, and route registration — a custom skill turns that pattern into a single command. `/scaffold component UserProfile` or `/new-endpoint POST /api/v1/invoices` produces the right files in the right places, following your team's conventions.

**Try it now:**
> Create a custom skill at `.claude/skills/scaffold-component.md` that scaffolds a React component. When I run `/scaffold-component Dashboard`, it should create: `src/components/Dashboard/Dashboard.tsx` with a typed functional component, `src/components/Dashboard/Dashboard.test.tsx` with a render smoke test, `src/components/Dashboard/Dashboard.stories.tsx` with a default Storybook story, and `src/components/Dashboard/index.ts` re-exporting the component. Follow our existing component patterns in `src/components/Button/`.

**Why this works:** Greenfield development has a bootstrapping problem: the first 20% of every new module is boilerplate that follows a known pattern. Custom skills encode that pattern once and eliminate the repetitive setup, letting you jump straight to the interesting work.

**Pros:**
- Enforces project conventions on every new component or module
- Eliminates the "copy an existing one and rename everything" workflow
- New team members produce consistent structure from day one

**Cons:**
- Requires upfront investment to define the skill
- Skills need maintenance when project conventions change

**Deeper:** See `approaches/custom-skills.md`

---

### 7. Official Plugins — Install proven workflows instead of building from scratch
**Level:** Intermediate

Before building a custom workflow for feature development, check whether a plugin already handles it. Plugins like `feature-dev` provide structured feature development workflows (plan, implement, test, review). `frontend-design` integrates visual design feedback into the build cycle. Installing a plugin takes minutes and gives you a battle-tested workflow that someone else has already debugged.

**Try it now:**
> Install the `feature-dev` plugin for this project. I want to use it to build the new user notification preferences feature — it should guide me through planning, implementation, and testing in a structured way. Show me what commands it adds and how to kick off the workflow.

**Why this works:** Plugin authors have already solved the workflow orchestration problems you would encounter building from scratch. By standing on their work, you skip the iteration cycle of designing, testing, and refining a custom workflow — and you benefit from improvements they push in future versions.

**Pros:**
- Production-ready workflows with zero development effort
- Community-maintained and regularly updated
- Can be combined with custom skills for project-specific extensions

**Cons:**
- Plugin conventions may not match your team's preferences exactly
- Dependency on external maintainers for updates and bug fixes

**Deeper:** See `approaches/official-plugins.md`
