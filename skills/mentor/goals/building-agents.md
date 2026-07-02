# Building AI Agents
*Last reviewed: 2026-07-02*

## When You're Here

You're not using an AI agent — you're building one. Maybe it's an internal tool that triages support tickets, a coding agent embedded in your product, or an autonomous worker for your team's recurring jobs. The questions are different from daily AI-assisted coding: what tools should the agent have, what model tier does it need, how do you constrain what it can do, and how do you know it works before users depend on it?

The good news: you don't have to start from a raw LLM API. Claude Code's own building blocks — agent definitions, the Agent SDK, permission scoping — are a production-tested agent architecture you can prototype in and graduate from.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Unsure what the agent should even do | Plan mode | Design the tool surface, boundaries, and failure modes before code |
| Want a working prototype today | Custom agents | A `.claude/agents/` definition is a running agent in ten lines of markdown |
| Ready to build a standalone product | Official plugins | `agent-sdk-dev` scaffolds Agent SDK projects and validates best practices |
| Agent needs multiple cooperating workers | Subagent delegation | Orchestrator-worker is the proven multi-agent architecture |
| Need to constrain what the agent can do | Permissions & safe autonomy | Tool allowlists and deny rules are your agent's safety spec |

**Hidden gem:** Custom Agents — prototyping your agent as a ten-line `.claude/agents/` file answers most of the design questions (tools, model, instructions) before you write a single line of SDK code.

## Approaches (Ranked)

### 1. Plan Mode — Design the agent before building it
**Level:** Beginner

Agent projects fail on unexamined design questions, not on code: what tools does it need (and which must it never have), what happens when it's uncertain, how does a human intervene? Plan mode turns those into an explicit design review — describe the agent's job and let Claude map the tool surface, the permission boundaries, the escalation paths, and the failure modes before anything is built.

**Try it now:**
> Enter plan mode. I want to build an agent that triages our support inbox: reads new tickets, classifies severity, drafts responses for common issues, and escalates anything involving billing or data loss to a human. Design it: what tools does it need, what should it be forbidden from doing, where are the escalation boundaries, and what's the minimal first version worth shipping?

**Why this works:** An agent is a policy wrapped around a model — and policies designed explicitly fail predictably, while policies that emerge from code fail creatively.

**Pros:**
- Surfaces tool-surface and safety questions when they're cheap to change
- Produces a design doc your team can review before investment
- Identifies the minimal first version instead of the imagined final one

**Cons:**
- Speculative design goes stale fast — validate with a prototype quickly

**Deeper:** See `approaches/plan-mode.md`

---

### 2. Custom Agents — Prototype in-product before writing SDK code
**Level:** Beginner

A file in `.claude/agents/` — frontmatter for model and allowed tools, markdown body for instructions — is a complete, running agent. Prototype your agent idea there first: iterate on its instructions, watch it work on real tasks, discover what tools it actually needs. Most of what you learn transfers directly to an SDK implementation; sometimes you discover the markdown version *is* the product.

**Try it now:**
> Create `.claude/agents/ticket-triager.md`: model haiku, tools Read and Grep only. Instructions: given a support ticket's text, classify it as bug/question/billing/incident, rate severity 1-4 using the rubric I'll paste, and draft a one-paragraph response for severity 3-4 bugs. Then run it against these five real tickets and let's evaluate the outputs together.

**Why this works:** The cheapest agent iteration loop is editing markdown and re-running — you converge on the right instructions and tool surface in an afternoon instead of a sprint of SDK rebuilds.

**Pros:**
- Zero infrastructure — a running agent in minutes
- Tool restrictions and model choice are one-line changes
- Instructions iterate at conversation speed

**Cons:**
- Lives inside Claude Code — shipping to end users still needs the SDK
- No custom UI or programmatic integration at this stage

**Deeper:** See `approaches/custom-agents.md`

---

### 3. Official Plugins — Scaffold the real thing with agent-sdk-dev
**Level:** Intermediate

When the prototype earns productization, the `agent-sdk-dev` official plugin scaffolds Claude Agent SDK projects and validates them against best practices. The SDK gives your agent the same engine Claude Code runs on — tool execution, permissioning, context management — so you build the agent's job, not agent infrastructure.

**Try it now:**
> Install with: /plugin install agent-sdk-dev@claude-plugins-official — then use it to scaffold a TypeScript Agent SDK project for the ticket-triager we prototyped, carrying over its instructions and Read/Grep-only tool policy.

**Why this works:** Agent infrastructure (tool loops, permission enforcement, context handling) is undifferentiated heavy lifting — inheriting a production-tested engine lets your effort go into the judgment your agent encodes.

**Pros:**
- Official scaffolding plus a validator for SDK best practices
- The SDK inherits Claude Code's tool and permission machinery
- Prototype learnings (instructions, tool surface) port directly

**Cons:**
- A real software project with dependencies and deployment — no longer just markdown
- SDK APIs evolve; pin versions and track release notes

**Deeper:** See `approaches/official-plugins.md`

---

### 4. Subagent Delegation — The architecture pattern for multi-agent systems
**Level:** Advanced

If your agent's job decomposes — research then draft then verify, or one worker per data source — the orchestrator-worker pattern from subagent delegation is the architecture to copy: a coordinator with the big picture, workers with clean contexts and narrow briefs, summaries flowing back up. Study how it behaves in Claude Code before you re-implement it in your own system. When you do build it, the Agent SDK's structured outputs (JSON Schema, Zod, or Pydantic) turn each worker's report into validated, type-safe data instead of free text — the handoff contract your orchestrator can actually rely on.

**Try it now:**
> Simulate my planned agent architecture with subagents before I build it: an orchestrator that takes a customer complaint, delegates to one worker to search the docs for relevant policy and another to search recent tickets for precedent, then synthesizes a response. Run it on this example complaint and show me each worker's brief and output — I want to see where the handoffs lose information.

**Why this works:** Multi-agent failure modes (context loss at handoffs, redundant work, contradictory findings) show up in a simulation you can run today, not after you've built the message bus.

**Pros:**
- Validates the decomposition before you commit to it
- Handoff and summary-loss problems surface immediately
- The working simulation doubles as the spec

**Cons:**
- Simulation fidelity has limits — production adds latency, retries, and cost concerns

**Deeper:** See `approaches/subagent-delegation.md`

---

### 5. Permissions & Safe Autonomy — Your agent's safety spec
**Level:** Intermediate

Whatever your agent does, the harder question is what it must never do. Claude Code's permission model — explicit tool allowlists, deny rules that outrank everything, modes from read-only to autonomous — is both the mechanism for your prototype and the vocabulary for your production design. Write the agent's permission policy as deliberately as its instructions.

**Try it now:**
> Draft the permission policy for the ticket-triager agent: it may read ticket text and search the docs; it must never send email, modify tickets, or read customer PII fields. Express that as a Claude Code tool allowlist for the prototype, and as a checklist of enforcement points I need when I port it to the SDK.

**Why this works:** Capability is easy to add and trust is hard to win back — agents earn adoption by having provable boundaries, not impressive demos.

**Pros:**
- Deny-first thinking catches dangerous capabilities at design time
- The prototype's policy is enforced by Claude Code, not by hoping the model behaves
- The policy doubles as the security review artifact for stakeholders

**Cons:**
- Overly tight policies produce an agent that escalates everything — iterate on the boundary

**Deeper:** See `approaches/safe-autonomy.md`
