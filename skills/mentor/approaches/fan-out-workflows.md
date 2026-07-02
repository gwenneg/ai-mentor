# Fan-Out Workflows (Orchestration)
*Last reviewed: 2026-07-01*

## What It Is

Fan-Out Workflows let you orchestrate tens to hundreds of AI agents through a JavaScript runtime built into Claude Code. You write a short script that defines phases, parallelism, and verification steps, and the runtime manages spawning agents, collecting results, and feeding outputs between stages. Think of it as a build system for AI tasks — you define the dependency graph, and the runtime handles execution.

## Why It Works

Some problems are embarrassingly parallel: reviewing 50 files, migrating 200 API endpoints, or checking a codebase against 30 security rules. A single agent doing these sequentially would take hours and degrade as its context fills up. Fan-out workflows apply the same principle as MapReduce or CI pipeline matrices — decompose the problem into independent units, process them in parallel with fresh-context agents, then aggregate. The key insight is that verification scales the same way: you can spawn "skeptic" agents to challenge findings, creating an adversarial review loop that catches false positives before they reach you.

## When to Use It

- Large-scale migration (rename a pattern across 200 files, update API versions across services)
- Exhaustive code review where you want each file reviewed independently by a fresh-context agent
- Multi-rule compliance checking (security policies, style guidelines, deprecation warnings)
- Generating structured output at scale (extracting metadata from hundreds of source files)

## When NOT to Use It

- Tasks with fewer than 5-10 independent units — the orchestration overhead is not worth it for small batches
- Highly interdependent changes where each step depends on the output of the previous step — use sequential subagents instead
- Exploratory work where you do not know the structure of the problem yet — plan first, fan out after
- When you need interactive human feedback during execution — fan-out workflows run autonomously

## How It Works

### Basic (Beginner)

1. Define your task list and the work each agent should do
2. Use `agent()` to spawn individual workers with specific prompts
3. Use `parallel()` to run a batch and wait for all results
4. Collect and process the results

Example — review three services in parallel:
```javascript
const services = ["auth", "billing", "notifications"];
const reviews = await parallel(services.map(svc =>
  agent(`Review services/${svc}/ for error handling gaps. Report findings as JSON.`)
));
```

5. Each agent gets a clean context window, reviews its service, and returns findings. `parallel()` waits for all three before continuing.

### Composing with Other Approaches (Intermediate)

- **Fan-out with worktree isolation**: Give each agent `isolation: "worktree"` so they can make edits in parallel. Useful for migration tasks where each agent modifies files.
- **Plan Mode then fan-out**: Use Plan Mode to identify the 50 files that need migration, then fan out agents to migrate each file independently.
- **Fan-out then review skills**: After agents make parallel changes across worktrees, run `/code-review` on the combined diff to catch cross-cutting issues that individual agents missed.

### Advanced Patterns

- **Adversarial verification**: For each finding from a reviewer agent, spawn N "skeptic" agents that try to disprove it. Only findings that survive skeptic challenge are reported. This dramatically reduces false positives.
```javascript
const verified = await parallel(findings.map(f =>
  agent(`Challenge this finding. Is it a real bug or a false positive? ${f.description}`)
));
```
- **Pipeline stages**: Use `pipeline()` for multi-stage processing where items flow through stages independently (no barrier between stages). Item 1 can be in stage 3 while item 5 is still in stage 1.
- **Structured output with schema**: Pass a `schema` option to `agent()` to get validated JSON output. The runtime enforces the schema at the tool-call layer, so you always get parseable, typed results.
- **Budget-aware loops**: Use `loop-until-dry` patterns where agents keep finding issues until a round produces zero new findings, with a budget cap to prevent runaway costs.

## Common Pitfalls

- **Unbounded fan-out**: Spawning 500 agents without a concurrency limit can exhaust API rate limits and produce unreliable results. Start with up to 16 concurrent agents (capped at min(16, cpu_cores - 2)) and scale up after validating.
- **Ignoring pipeline vs parallel**: `pipeline()` lets items flow through stages independently (no barrier). `parallel()` waits for everything to complete before continuing. Using the wrong one causes either unnecessary blocking or premature aggregation.
- **Skipping verification**: Raw fan-out results have a false-positive rate. Always add a verification phase — even a simple "does this finding still hold?" recheck — before presenting results to a human.
- **Overly broad agent prompts**: If each agent gets the same vague prompt, you get 50 copies of shallow analysis. Give each agent specific scope: the file, the rule, the exact question to answer.

## Real-World Example

Your team is migrating from `moment.js` to `dayjs` across a monorepo with 180 files that import moment. A single agent would lose context by file 30. Instead:

```javascript
const files = await agent("List all files importing 'moment'. Return as JSON array.");

const migrations = await parallel(files.map(f =>
  agent({
    prompt: `Migrate ${f} from moment.js to dayjs. Replace all moment() calls
             with dayjs(), update duration APIs, and handle timezone calls with
             dayjs/plugin/timezone. Return the file path and a summary of changes.`,
    isolation: "worktree"
  })
));

const verified = await parallel(migrations.map(m =>
  agent(`Verify the migration in ${m.file}. Check that: (1) no moment imports remain,
         (2) timezone handling uses dayjs/plugin/timezone, (3) format strings are
         dayjs-compatible. Report pass/fail with details.`)
));

log(`Migrated ${migrations.length} files. ${verified.filter(v => v.pass).length} passed verification.`);
```

The runtime spawns agents in batches, each migrating one file in its own worktree. A second wave verifies each migration. The whole operation takes minutes instead of hours, and each agent works with full context on a single file rather than degraded context across 180.

## Sources

- [Claude Code Common Workflows](https://docs.anthropic.com/en/docs/claude-code/common-workflows) — Official docs on workflow patterns and fan-out orchestration
- [Building a C Compiler with Claude](https://www.anthropic.com/engineering/building-c-compiler) — Anthropic engineering blog demonstrating parallel fan-out patterns
