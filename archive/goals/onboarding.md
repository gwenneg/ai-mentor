# Onboarding & Environment Setup
*Last verified: 2026-06-28*

## When You're Here

You're starting from zero in a new environment. Maybe you just joined a company and your laptop is fresh out of the box. Maybe you're rotating to a different team and their stack is completely foreign. Or maybe you're inheriting a project from a departing colleague, and the handoff was a 30-minute call and a "good luck." Whatever the trigger, you need to go from "I don't even know where the repo is" to "I can ship a change" — and the clock is ticking because the team expects you to be productive.

Onboarding is broader than code understanding. Code understanding assumes you're already set up and just need to read the codebase. Onboarding includes everything before that: getting your machine configured, understanding team processes and rituals, learning who owns what services, figuring out how the deployment pipeline works, getting access to the right Slack channels and wikis, and building the mental map of how the team actually operates versus what the documentation says. The biggest time sink isn't any single step — it's the constant context-switching between 10 different sources to find the next piece of information you need.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Team knowledge is scattered across Confluence, Slack, Notion, and READMEs | MCP context | Pulls all sources into one conversation instead of 10 tabs |
| Need a systematic understanding of the codebase and architecture | Plan mode | Structured exploration builds a reusable mental map |
| Codebase uses frameworks or patterns you've never worked with | Deep research | Learn the conventions so you can distinguish scaffolding from business logic |
| Local dev setup takes a half-day of manual steps | Custom skills | A `/setup-dev` skill turns the checklist into one command |
| Team maintains shared onboarding workflows and conventions | Custom plugins | Install the team's plugin to get their accumulated knowledge immediately |

**Hidden gem:** Custom Skills — a `/setup-dev` skill is executable documentation: it can't silently go stale the way a wiki page does.

## Approaches (Ranked)

### 1. MCP Context — Pull from wikis, Slack, and team docs
**Level:** Intermediate

Onboarding knowledge lives in scattered places: the getting-started guide is in Confluence, the environment variables are documented in a pinned Slack message from 2024, the architecture diagram is in a Notion page that three people have bookmarked, and the actual deployment process is in a README that hasn't been updated since the team migrated to Kubernetes. MCP context servers let Claude pull all of these into one conversation, so you can ask questions that span multiple sources without hunting across tabs.

**Try it now:**
> Connect to our Confluence space "Engineering Onboarding" and the Slack channel #platform-team. Pull the getting-started guide, the local development setup instructions, and any pinned messages about environment configuration. Cross-reference them with the actual `README.md` and `docker-compose.yml` in the repo root. Which setup steps in the docs are still accurate? Which are outdated? Produce a single, consolidated setup guide that reflects the current state of the codebase.

**Why this works:** The biggest onboarding bottleneck is not understanding information — it's finding it. New hires waste hours navigating between documentation systems, following outdated links, and asking colleagues "where is the doc for X?" MCP eliminates the search problem by bringing every source into one queryable context, so you spend your time learning instead of hunting.

**Pros:**
- Consolidates scattered knowledge into a single conversation
- Surfaces outdated documentation by cross-referencing docs with code
- Answers span multiple sources without manual tab-switching

**Cons:**
- Requires MCP server setup for each documentation source
- Only surfaces knowledge that was written down somewhere — misses pure tribal knowledge
- May hit context limits when pulling from many large documentation sources

**Deeper:** See `approaches/mcp-context.md`

---

### 2. Plan Mode — Structured exploration roadmap
**Level:** Beginner

Instead of randomly clicking through files hoping to stumble on the important ones, plan mode builds a systematic exploration roadmap. Start from the entry points: where does the application boot? What are the main route handlers or event consumers? How does data flow from an API request through the service layers to the database? This structured approach produces an architecture summary you can refer back to for weeks — a map that makes every subsequent task faster.

**Try it now:**
> Enter plan mode. I just joined this team and need to understand how this service works end-to-end. Start with the entry point — find `main` or the application bootstrap file. Map the top-level directory structure and explain what each major directory owns. Trace one core user-facing request from the HTTP handler through the service layer to the database. Identify the deployment configuration (Dockerfile, Kubernetes manifests, CI pipeline). Who are the service's upstream and downstream dependencies? Produce an architecture summary I can use as a reference for my first month.

**Why this works:** Experienced engineers onboard fast because they have a system: they learn the shape of the architecture first, then fill in module-level details as tickets demand. Plan mode replicates this expert strategy — it prioritizes the connections between components (who calls whom, what data flows where) over the implementation details within any single file.

**Pros:**
- Builds a reusable mental model, not just one-off answers
- Identifies critical paths and service boundaries first
- Produces documentation you can share with the next person who onboards

**Cons:**
- The architecture summary is a snapshot — it drifts as the code evolves
- Misses operational knowledge like "this service is flaky on Mondays" that lives in team memory

**Deeper:** See `approaches/plan-mode.md`

---

### 3. Deep Research — Learn the tech stack before reading the code
**Level:** Beginner

If the team uses frameworks, patterns, or infrastructure you haven't worked with before, diving into the codebase first is working backwards. You'll confuse framework conventions with business logic and mistake boilerplate for intentional design decisions. Deep research lets you learn the idioms — how NestJS organizes modules, what Terraform state files are for, why the team uses event sourcing — so that when you read the code, you can instantly recognize what's scaffolding versus what's the team's actual work.

**Try it now:**
> /deep-research This project uses a hexagonal architecture (ports and adapters) with Go. I'm seeing directories like `internal/domain/`, `internal/ports/`, `internal/adapters/`, and `cmd/`. Explain the hexagonal architecture pattern as typically implemented in Go: what goes in each layer, how do ports and adapters connect, and what are the dependency rules? Include Go-specific conventions so I can read this codebase knowing which code is structural and which is business logic.

**Why this works:** Frameworks and architecture patterns encode decisions into structure. If you don't understand those decisions, every directory looks arbitrary and every file looks equally important. Deep research gives you the vocabulary and mental model of the stack so you can read the codebase fluently from day one.

**Pros:**
- Dramatically accelerates onboarding to unfamiliar tech stacks
- Prevents misreading framework conventions as business logic
- Covers official docs, community patterns, and real-world usage

**Cons:**
- Doesn't cover the team's specific customizations or deviations from standard patterns
- May surface outdated practices for rapidly evolving frameworks

**Deeper:** See `approaches/deep-research.md`

---

### 4. Custom Skills — Automate the setup checklist
**Level:** Intermediate

Every team has a setup checklist: clone the repo, install dependencies, configure the local database, set environment variables, seed test data, run smoke tests, verify the dev server starts. It takes a new hire half a day, and something always goes wrong because step 7 assumed you did step 3 differently than the docs described. A custom `/setup-dev` skill encodes this entire checklist as one command — it installs dependencies, configures services, handles platform differences, and runs verification checks, telling you exactly what failed and why instead of leaving you stuck.

**Try it now:**
> Create a custom skill at `.claude/skills/setup-dev.md`. When invoked with `/setup-dev`, it should: (1) check prerequisites — Node.js >= 20, Docker running, `psql` client installed, (2) copy `.env.example` to `.env` and fill in default local development values, (3) run `docker-compose up -d` to start PostgreSQL and Redis, (4) run `npm install`, (5) run `npm run db:migrate` and `npm run db:seed`, (6) run `npm run dev` and verify the health endpoint at `localhost:3000/health` responds with 200. If any step fails, explain what went wrong and how to fix it. Output a summary of what was set up and what to do next.

**Why this works:** Setup checklists break because they're written once and maintained never. A skill is executable documentation — if a step breaks, you fix the skill and every future onboarder benefits. It also handles the conditional logic ("if macOS, use brew; if Linux, use apt") that static docs handle poorly.

**Pros:**
- Turns a half-day manual process into a single command
- Self-verifying — each step confirms success before proceeding
- Executable documentation that stays current because it breaks visibly when outdated

**Cons:**
- Requires upfront investment to encode the full setup process
- Platform-specific edge cases (M1 Mac, WSL, corporate proxies) need explicit handling

**Deeper:** See `approaches/custom-skills.md`

---

### 5. Custom Plugins — Install team-specific onboarding workflows
**Level:** Intermediate

If the team maintains a Claude Code plugin with onboarding skills, project conventions, and common tasks, installing it gives you the team's accumulated workflow knowledge immediately. Instead of discovering that the team has a special deployment process or a non-obvious testing convention, the plugin surfaces these as available skills and contextual guidance. It's the difference between reading a wiki page about how the team works and having a colleague sitting next to you who knows.

**Try it now:**
> Install our team's Claude Code plugin: `claude plugin install @myorg/platform-team-plugin`. Then run `/team-overview` to see the service architecture, team members and their areas of ownership, on-call rotation, and links to key dashboards. Check what other skills the plugin provides — are there skills for common tasks like creating a new API endpoint, running the integration test suite, or deploying to staging?

**Why this works:** Plugins encode team-specific knowledge in a form that's immediately usable. Unlike documentation that you have to find, read, and interpret, a plugin's skills are discoverable via tab-completion and executable on the spot. The team's best practices become your default workflow from day one.

**Pros:**
- Instant access to team-specific workflows and conventions
- Discoverable — tab-completion shows available skills without reading docs
- Maintained by the team, so it evolves with the team's actual practices

**Cons:**
- Only useful if the team has invested in building and maintaining the plugin
- Plugin quality varies — a poorly maintained plugin is worse than no plugin
- Requires the team to adopt Claude Code plugin conventions

**Deeper:** See `approaches/custom-plugins.md`

---

### 6. Project Memory & Context Docs — Make your ramp-up inheritable
**Level:** Beginner

Everything you learn during onboarding — build quirks, conventions, "who owns what" — evaporates unless it's written where the AI reads it. Run `/init` on day one to generate a starter `CLAUDE.md`, then append every correction you find yourself repeating. Auto memory captures the rest on its own. By week two, your sessions start pre-loaded with everything week one taught you — and the next person to onboard inherits it from `git pull`.

**Try it now:**
> Run `/init` to generate a CLAUDE.md for this project. Then review it with me: check the build and test commands it discovered actually work, and add the three things I've had to explain twice this week — the migration script wrapper, the append-only tables, and the integration tests needing Docker.

**Why this works:** Onboarding knowledge has a half-life of one session unless persisted. CLAUDE.md converts per-session re-explaining into a one-time write that compounds for you, the AI, and every future teammate.

**Pros:**
- Day-one `/init` gives immediate structure to your ramp-up notes
- Auto memory accumulates learnings with zero manual effort
- The next new hire starts with your accumulated context for free

**Cons:**
- Needs pruning discipline — a bloated CLAUDE.md reduces adherence
- Personal preferences must go in `CLAUDE.local.md`, not the shared file

**Deeper:** See `approaches/project-memory.md`
