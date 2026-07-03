# Model & Effort Selection
*Last verified: 2026-07-01*

## What It Is

Model & Effort Selection lets you match the model and the reasoning depth to what each task needs, instead of running everything at one setting. Claude Code ships multiple models (Opus for deep reasoning, Sonnet as the balanced default, Haiku for fast, cheap work) plus an adjustable effort level that controls how much thinking the model spends per response. You choose per task, per session, or per agent — and switching takes seconds.

## Why It Works

No single model configuration is best at everything. Deep reasoning models at high effort excel at architecture and debugging but are slower and more expensive than routine work justifies. Fast models handle boilerplate and mechanical edits efficiently but miss subtleties in multi-file reasoning. Matching the model to the task is the same engineering principle as choosing between a profiler, a debugger, and a linter — the right tool for the job. Because subagents can each run their own model, you can compose both in one workflow: a strong model orchestrates while cheap models handle the volume.

## When to Use It

- Optimizing costs by running routine tasks on cheaper models and reserving Opus or high effort for complex reasoning
- Long sessions where most turns are mechanical — drop the effort level, then raise it for the hard debugging moment
- Delegating bulk work to subagents: Haiku agents for repetitive checks, Opus agents for security review
- Architecture decisions, tricky debugging, or security analysis where maximum reasoning depth pays for itself

## When NOT to Use It

- When you are just getting started — the defaults are well-chosen; master the workflow before optimizing cost
- When the task is simple enough that model choice does not meaningfully affect the result
- When constant switching costs you more in attention than it saves in tokens

## How It Works

### Basic (Beginner)

1. Check your current model with `/model` — it lists available models and saves your choice as the default for new sessions.
2. Switch based on the task ahead: pick Opus before a gnarly debugging session, Haiku for a batch of mechanical renames.
3. Adjust reasoning depth with `/effort` (`low`, `medium`, `high`, `xhigh`, `max`, or `auto`) — lower effort responds faster and costs less, higher effort thinks longer on hard problems.
4. On Opus, toggle `/fast` for fast mode: the same model with faster output, useful for interactive back-and-forth.
5. Continue working. The conversation context carries over — the model you switch to sees everything discussed so far.

### Composing with Other Approaches (Intermediate)

- **Model selection plus subagent delegation**: Run the main session on a strong model and delegate volume work to Haiku subagents — summarizing files, checking conventions, scanning for patterns. The orchestrator keeps the judgment; the cheap agents do the reading.
- **Effort plus Plan Mode**: Plan a risky migration at `xhigh` effort where reasoning quality matters most, then execute the approved plan at normal effort — the thinking-intensive phase is the plan, not the edits.
- **Model tiers plus custom agents**: Define the same reviewer agent twice in `.claude/agents/` — a `quick-review` on Haiku for draft PRs and a `deep-review` on Opus for release candidates — and pick by stakes.

### Advanced Patterns

- **Per-agent model overrides**: Custom agent definitions accept a `model` field in frontmatter (`opus`, `sonnet`, `haiku`), and ad-hoc subagents can override the model per invocation. A fan-out workflow can run its finder agents cheap and its verifier agents strong.
- **Effort-tiered workflows**: In orchestrated workflows, set the effort per stage — low for mechanical extraction stages, max for the adversarial verification stage where missed reasoning means false positives slip through.
- **Cost-aware budgeting**: Watch `/usage` during long sessions. If a session is burning budget on routine turns, drop the model or effort and save the expensive configuration for the tasks where quality measurably differs.

## Common Pitfalls

- **Switching models too frequently**: Every switch has a cognitive cost — you recalibrate expectations for the new model's style and capability. Switch when the task genuinely demands it, not out of curiosity.
- **Running everything at max effort**: High effort on trivial tasks buys latency, not quality. Reserve `xhigh` and `max` for problems where reasoning depth is the bottleneck.
- **Using a fast model for subtle work**: Haiku is excellent at well-specified mechanical tasks and weaker at multi-file reasoning and nuanced review. If a cheap agent's output needs constant correction, the time cost exceeds the token savings.
- **Forgetting the setting persists**: `/model` saves your choice as the default for new sessions. If you dropped to a fast model for a batch job, remember to switch back before your next deep session.

## Real-World Example

You are building a new GraphQL API layer for an existing REST backend. The work has three distinct phases, each suited to a different configuration.

First, you need to design the schema. You switch to Opus at high effort:

```
/model opus
/effort high
> Analyze the REST endpoints in src/api/routes/ and propose a GraphQL schema
  that consolidates the N+1 query patterns in the order and inventory endpoints.
```

Claude produces a schema in `src/graphql/schema.graphql` with thoughtful type relationships and resolver structure. This took careful reasoning about data dependencies.

Second, you need to generate 14 resolver files — mostly mechanical translation from REST handlers. You delegate to cheap subagents:

```
> Spawn one Haiku subagent per type in schema.graphql. Each generates the
  resolver for its type, following the pattern in src/graphql/resolvers/user.ts
  and mapping to the corresponding REST client in src/api/clients/.
```

The agents generate all 14 files in parallel in under two minutes at a fraction of the cost.

Third, you need DataLoader batching logic to solve the N+1 queries — the subtle part. You return to Opus at `xhigh` effort for the batching design, then run `npm test` — all 31 new tests pass. Total session cost: well under half of what running everything on the deep-reasoning configuration would have been.

## Sources

- [Claude Code Model Configuration](https://code.claude.com/docs/en/model-config) — Official docs covering /model, available models, and effort levels
- [Claude Code Fast Mode](https://code.claude.com/docs/en/fast-mode) — Official docs for fast mode on Opus
