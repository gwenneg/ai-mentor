# Custom Agent Definitions
*Last verified: 2026-07-06*

## What It Is

Custom Agent Definitions let you create reusable, specialized agent types that live in your project as `.claude/agents/<name>.md` files. Each definition specifies the agent's purpose, which tools it can access, which model it should use, and any rules it must follow. Instead of writing a detailed prompt every time you need a particular kind of analysis or task, you define the agent once and invoke it by name.

## Why It Works

Specialization produces better results than generalization. A general-purpose agent asked to do a security review will do a passable job. An agent explicitly scoped to security — given only Read and Grep tools, running on a high-reasoning model, armed with your company's security checklist — will do a thorough job. Custom agent definitions encode this specialization so it is consistent across invocations and across team members, rather than depending on whoever happens to write the best prompt that day.

## When to Use It

- You have a recurring task that benefits from a specific set of tools, model, or instructions (security review, migration authoring, documentation audit)
- Multiple team members need to perform the same specialized task and you want consistent quality regardless of who invokes it
- You want to restrict an agent's capabilities for safety — read-only agents that cannot edit code, or agents scoped to specific directories
- You need different cost/quality trade-offs for different tasks (haiku for quick checks, opus for deep analysis)

## When NOT to Use It

- One-off tasks where the overhead of defining an agent file is not justified — use ad-hoc subagent delegation instead
- Tasks that change significantly each time — a rigid agent definition will need constant updating

## How It Works

### Basic (Beginner)

1. Create a file at `.claude/agents/<agent-name>.md` in your project, or in `~/.claude/agents/` to make the agent available in all your projects
2. Add frontmatter specifying the model, allowed tools, and any configuration
3. Write the agent's system instructions in the markdown body — its role, rules, and approach
4. Invoke the agent by name from your Claude Code session
5. The agent runs with its defined tools and model, follows its instructions, and returns results

Example agent definition at `.claude/agents/security-reviewer.md`:
```markdown
---
name: security-reviewer
description: Analyze code for security vulnerabilities following OWASP Top 10
model: opus
tools: Read, Grep, Bash
---

You are a security reviewer. Analyze the provided code for:
- Injection vulnerabilities (SQL, command, template)
- Authentication and authorization gaps
- Secrets or credentials in source code
- Unsafe deserialization

Follow the OWASP Top 10 checklist. Report findings with severity, location,
and a recommended fix. If you find nothing, say so — do not invent issues.
```

### Composing with Other Approaches (Intermediate)

- **Custom agents plus worktree isolation**: Define a `refactoring-agent` that uses Edit and Write tools, then set `isolation: worktree` in its frontmatter so it runs in a temporary git worktree. It makes changes on a separate branch, and you review the diff before merging. The agent definition ensures consistent refactoring style; the worktree ensures safety.
- **Custom agents plus skills**: Create a `/security-scan` skill that invokes your `security-reviewer` agent on every file changed in the current branch. The skill handles orchestration (which files, how to report); the agent handles analysis.
- **Custom agents plus hooks**: Define hooks inside the agent's own frontmatter — a `db-reader` agent with a `PreToolUse` hook that runs a validation script can execute queries while the script blocks anything but `SELECT`. The agent's instructions state the rule; the hook enforces it deterministically.

### Advanced Patterns

- **Model-tiered agent families**: Define the same agent at different cost levels. A `quick-review` agent uses haiku and checks for obvious issues in seconds. A `deep-review` agent uses opus and performs thorough analysis. Developers choose based on the stakes — haiku for draft PRs, opus for release candidates.
- **Agent Teams** (experimental — requires `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` env var): Define multiple agents that share a task list and communicate via peer messages. For example, the team lead breaks the work into shared tasks, then a `coder` teammate implements each task and a `reviewer` teammate checks each implementation. They coordinate without a human orchestrating every handoff.
- **Custom agents as background sessions**: To run several custom agents at once (five reviewers on five services), dispatch each as its own background session with `claude --agent <name> --bg "<task>"`, then run `claude agents` to monitor them all from one screen and step in when one needs input. Agent view lists background sessions only — subagents spawned inside a session appear in that session's agent panel, not in agent view.

## Common Pitfalls

- **Overly broad tool access**: An agent defined for code review does not need Edit or Write tools. Giving it those tools risks accidental modifications. Scope tools to the minimum needed for the task.
- **Wrong model for the job**: Using opus for a simple formatting check wastes money and time. Using haiku for a nuanced security review misses subtle vulnerabilities. Match the model to the complexity of the task.
- **Duplicating what subagent delegation does**: If you only need a specialized agent once, use ad-hoc subagent delegation with a detailed prompt. Custom agent definitions pay off when you invoke the same specialization repeatedly.
- **Stale instructions**: When your codebase conventions change (new ORM, new test framework, new security requirements), the agent definition must be updated too. Treat agent files as living documentation that is reviewed alongside your code.

## Real-World Example

A platform team maintains 12 microservices and performs database migrations frequently. Each migration must follow specific rules: use reversible operations, avoid locking tables over 1M rows, include a rollback script, and update the schema registry. Different developers write migrations with varying quality.

The team creates `.claude/agents/migration-helper.md`:

```markdown
---
name: migration-helper
description: Generate safe, reversible PostgreSQL migrations following team conventions
model: opus
tools: Read, Edit, Bash
---

You are a database migration specialist for our PostgreSQL services.

Rules:
- Read the current schema from db/schema.sql before writing any migration
- All migrations must be reversible (include both up and down)
- For tables over 1M rows (check db/table_stats.json), use concurrent index creation and batched updates
- Generate a rollback script in db/rollbacks/
- Update the schema registry entry in services/schema-registry/schemas/
- Follow naming convention: YYYYMMDD_HHMMSS_description.sql
```

Now any developer on the team can invoke this agent when they need a migration. A junior developer who has never written a migration for a high-traffic table gets the same quality output as the senior DBA. The agent reads the table stats, sees that the `orders` table has 4.2M rows, and automatically uses `CREATE INDEX CONCURRENTLY` and batched `ALTER TABLE` operations. It generates the rollback script, updates the schema registry, and follows the naming convention — all without the developer knowing these rules existed.

## Sources

- [Claude Code Sub-Agents](https://code.claude.com/docs/en/sub-agents) — Official docs for creating custom agent definitions in .claude/agents/
- [Claude Code Best Practices](https://code.claude.com/docs/en/best-practices) — Official best-practices guide, including a section on creating custom subagents
