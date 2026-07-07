# Fan-Out Workflows (Orchestration)
*Last verified: 2026-07-06*

## What It Is

Fan-Out Workflows let you orchestrate tens to hundreds of AI agents through a JavaScript runtime built into Claude Code. You describe the task, Claude writes a short script that defines phases, parallelism, and verification steps, and the runtime manages spawning agents, collecting results, and feeding outputs between stages. Think of it as a build system for AI tasks — you define the dependency graph, and the runtime handles execution.

## Why It Works

Embarrassingly parallel problems decompose into independent units processed in parallel by fresh-context agents — and verification scales the same way, with skeptic agents challenging findings before they reach you.

## When to Use It

- Large-scale migration (rename a pattern across 200 files, update API versions across services)
- Exhaustive code review where you want each file reviewed independently by a fresh-context agent
- Multi-rule compliance checking (security policies, style guidelines, deprecation warnings)
- Generating structured output at scale (extracting metadata from hundreds of source files)

## When NOT to Use It

- Tasks with fewer than 5-10 independent units — the orchestration overhead is not worth it for small batches
- Highly interdependent changes where each step depends on the output of the previous step — use sequential subagents instead
- When you need interactive human feedback during execution — fan-out workflows run autonomously

## How It Works

### Basic (Beginner)

1. Ask for a workflow in your prompt — say "use a workflow to ..." or include the keyword `ultracode` — and Claude writes the orchestration script (paid plans; on Pro, enable Dynamic workflows in `/config`)
2. In the script, `agent()` spawns individual workers with specific prompts
3. `parallel()` runs a batch and waits for all results; `pipeline()` runs one agent per item in a list
4. Approve the run when prompted, then track progress with `/workflows` while your session stays responsive — and if the run is worth repeating, save it from that view as a reusable command (stored in `.claude/workflows/` for the project, or `~/.claude/workflows/` for yourself)

Example — review three services in parallel:
```javascript
const services = ["auth", "billing", "notifications"];
const reviews = await parallel(services.map(svc =>
  agent(`Review services/${svc}/ for error handling gaps. Report findings as JSON.`)
));
```

5. Each agent gets a clean context window, reviews its service, and returns findings. `parallel()` waits for all three before continuing.

### Composing with Other Approaches (Intermediate)

- **Fan-out with worktree isolation**: Ask for each agent to work in its own isolated worktree copy so parallel edits don't conflict. Useful for migration tasks where each agent modifies files.
- **Plan Mode then fan-out**: Use Plan Mode to identify the 50 files that need migration, then fan out agents to migrate each file independently.
- **Fan-out then review skills**: After agents make parallel changes across worktrees, run `/code-review` on the combined diff to catch cross-cutting issues that individual agents missed.

### Advanced Patterns

- **Adversarial verification**: For each finding from a reviewer agent, spawn N "skeptic" agents that try to disprove it. Only findings that survive skeptic challenge are reported. This dramatically reduces false positives.
```javascript
const verified = await parallel(findings.map(f =>
  agent(`Challenge this finding. Is it a real bug or a false positive? ${f.description}`)
));
```
- **Per-item pipelines**: Use `pipeline(items, fn)` to run one agent per item in a list; chain multiple `agent()` calls inside `fn` when an item needs staged processing.
- **Structured output with schema**: Pass a `schema` option to `agent()` to get validated JSON output. The runtime validates the output against the schema, so you get parseable, typed results.
- **Budget-aware loops**: Use `loop-until-dry` patterns where agents keep finding issues until a round produces zero new findings, with a budget cap to prevent runaway costs.

## Common Pitfalls

- **Unbounded fan-out**: Spawning 500 agents without a concurrency limit can exhaust API rate limits and produce unreliable results. The runtime caps runs at 16 concurrent agents (fewer on limited-CPU machines) and 1,000 agents per run; validate on a small slice before running the full set.
- **Ignoring pipeline vs parallel**: `pipeline()` runs one agent per item in a list; `parallel()` runs a batch and waits for all of it before continuing. Pick the shape that matches whether results are consumed per item or as a whole batch.
- **Skipping verification**: Raw fan-out results have a false-positive rate. Always add a verification phase — even a simple "does this finding still hold?" recheck — before presenting results to a human.
- **Overly broad agent prompts**: If each agent gets the same vague prompt, you get 50 copies of shallow analysis. Give each agent specific scope: the file, the rule, the exact question to answer.

## Real-World Example

Your team is migrating from `moment.js` to `dayjs` across a monorepo with 180 files that import moment. A single agent would lose context by file 30. Instead:

```javascript
const files = await agent("List all files importing 'moment'. Return as JSON array.");

const migrations = await parallel(files.map(f =>
  agent(`Migrate ${f} from moment.js to dayjs in your own isolated worktree copy.
         Replace all moment() calls with dayjs(), update duration APIs, and handle
         timezone calls with dayjs/plugin/timezone. Return the file path and a
         summary of changes.`, { label: f })
));

const verified = await parallel(migrations.map(m =>
  agent(`Verify the migration in ${m.file}. Check that: (1) no moment imports remain,
         (2) timezone handling uses dayjs/plugin/timezone, (3) format strings are
         dayjs-compatible. Report pass/fail with details.`)
));

return `Migrated ${migrations.length} files. ${verified.filter(v => v.pass).length} passed verification.`;
```

The runtime spawns agents in batches, each migrating one file in its own worktree. A second wave verifies each migration. The whole operation takes minutes instead of hours, and each agent works with full context on a single file rather than degraded context across 180.

## Sources

- [Dynamic Workflows](https://code.claude.com/docs/en/workflows) — Official docs on orchestrating subagents at scale with workflow scripts
- [Building a C Compiler with Claude](https://www.anthropic.com/engineering/building-c-compiler) — Anthropic engineering blog demonstrating parallel fan-out patterns
