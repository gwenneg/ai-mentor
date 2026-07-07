# Custom Skills (Slash Commands)
*Last verified: 2026-07-06*

## What It Is

Custom Skills let you package a repeatable workflow into a slash command that anyone on your team can invoke. You write a markdown file that describes what the skill does, include any scripts or templates it needs, and from that point forward the workflow is a single `/command` instead of a multi-step manual process. Skills live in `.claude/skills/<name>/SKILL.md` (or `~/.claude/skills/` for personal, cross-project skills), can accept arguments, run bundled scripts, and reference local documentation — and Claude can also invoke them automatically when a request matches their description.

## Why It Works

Skills capture the judgment calls, file conventions, and sequencing a senior engineer would follow — and make that knowledge executable by anyone as a single command.

## When to Use It

- A workflow you or your team performs repeatedly (weekly migrations, component scaffolding, release prep)
- A multi-step process that involves reading context, generating multiple files, and following project conventions
- Onboarding tasks where a new team member should be able to do the right thing without tribal knowledge
- Processes that combine AI judgment with deterministic steps (generate code, then run a specific formatter, then update a manifest)

## When NOT to Use It

- One-off tasks you will never repeat — the overhead of writing the skill exceeds the benefit
- Fully deterministic workflows with no AI judgment needed — a shell script or Makefile target is simpler and faster

## How It Works

### Basic (Beginner)

1. Create the directory `.claude/skills/<skill-name>/` in your project
2. Write a `SKILL.md` file with a one-line `description:` in YAML frontmatter (Claude uses it to decide when to load the skill automatically), followed by the steps to follow and what inputs it expects
3. Use `$ARGUMENTS` in the skill file to accept user input (e.g., `/create-migration add-user-roles-table`)
4. Optionally add bundled scripts (shell, Python, etc.) in the same directory that the skill references for deterministic steps
5. Invoke the skill with `/<skill-name>` — Claude reads the SKILL.md, follows the instructions, and executes the workflow

Example SKILL.md for a migration generator:
```markdown
---
description: Create a database migration from a short description. Use when adding or changing tables.
---

Create a new database migration for: $ARGUMENTS

Steps:
1. Run `${CLAUDE_SKILL_DIR}/generate.sh $ARGUMENTS` to create the migration file
2. Read the schema in `db/schema.prisma` for current table definitions
3. Write the migration SQL following our conventions in `docs/migration-guide.md`
4. Add a corresponding test in `tests/migrations/`
```

### Composing with Other Approaches (Intermediate)

- **Skills plus hooks**: Create a skill that scaffolds a new API endpoint, then add a PostToolUse hook that automatically runs the linter and type checker on every file the skill generates. The skill handles creation; the hook enforces quality.
- **Skills plus subagent delegation**: A complex skill can spawn subagents for independent subtasks. A `/new-service` skill might spawn one agent to generate the service code, another to write tests, and a third to update the CI configuration — all in parallel.
- **Skills plus plan mode**: For large skills, start with plan mode to outline what will be generated, get user confirmation, then execute. This prevents surprises when a skill creates 10+ files.

### Advanced Patterns

- **Forked context with `context: fork`**: Run a skill in an isolated subagent context so it does not pollute your main session's context window. The forked skill does not see your conversation history, so its instructions must be self-contained; it executes, returns a summary, and your main session stays clean. Useful for skills that read many files.
- **Skills with reference docs**: Bundle project-specific reference material (API schemas, style guides, architecture docs) in the skill directory. The skill instructions tell Claude to read these before generating code, ensuring output matches your conventions without relying on the main context.
- **Argument parsing patterns**: Use structured argument formats in your SKILL.md to handle multiple parameters: `/create-endpoint POST /api/users CreateUserRequest`. Use positional substitutions — `$0` for the method, `$1` for the path, `$2` for the request type — instead of asking Claude to parse the raw `$ARGUMENTS` string.

## Common Pitfalls

- **Skills that are too vague**: A SKILL.md that says "create a component" without specifying file locations, naming conventions, or what files to read for context will produce inconsistent results. Be specific about the steps and reference files.
- **Skipping the deterministic parts**: If a step should always produce the same output (timestamps, UUIDs, boilerplate), use a bundled script instead of asking the AI to generate it. AI adds variance where you want consistency.
- **Forgetting to version the skill**: Skills live in your repo. When your project conventions change (new test framework, different directory structure), update the skill or it will generate outdated scaffolding.
- **Overloading a single skill**: A skill that handles creation, updating, deletion, and listing is four skills pretending to be one. Keep each skill focused on one workflow.

## Sources

- [Claude Code Skills](https://code.claude.com/docs/en/skills) — Official docs for creating custom skills in .claude/skills/ directories
- [Skill authoring best practices](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices) — Official guidance on writing effective SKILL.md instructions across Claude products
