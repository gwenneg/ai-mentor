# Custom Skills (Slash Commands)
*Last reviewed: 2026-07-01*

## What It Is

Custom Skills let you package a repeatable workflow into a slash command that anyone on your team can invoke. You write a markdown file that describes what the skill does, include any scripts or templates it needs, and from that point forward the workflow is a single `/command` instead of a multi-step manual process. Skills live in `.claude/skills/<name>/SKILL.md` and can accept arguments, run bundled scripts, and reference local documentation.

## Why It Works

Most engineering teams have workflows that are done the same way every time but never get automated because they fall in the gap between "too complex for a shell script" and "not worth building a full tool for." Skills close that gap. They capture the judgment calls, file conventions, and sequencing that a senior engineer would follow — and make that knowledge executable. The next person who needs to do the task does not need to know the steps; they just invoke the command.

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
2. Write a `SKILL.md` file that describes what the skill does, what inputs it expects, and what steps to follow
3. Use `$ARGUMENTS` in the skill file to accept user input (e.g., `/create-migration add-user-roles-table`)
4. Optionally add bundled scripts (shell, Python, etc.) in the same directory that the skill references for deterministic steps
5. Invoke the skill with `/<skill-name>` — Claude reads the SKILL.md, follows the instructions, and executes the workflow

Example SKILL.md for a migration generator:
```markdown
Create a new database migration for: $ARGUMENTS

Steps:
1. Run `./skills/create-migration/generate.sh $ARGUMENTS` to create the migration file
2. Read the schema in `db/schema.prisma` for current table definitions
3. Write the migration SQL following our conventions in `docs/migration-guide.md`
4. Add a corresponding test in `tests/migrations/`
```

### Composing with Other Approaches (Intermediate)

- **Skills plus hooks**: Create a skill that scaffolds a new API endpoint, then add a PostToolUse hook that automatically runs the linter and type checker on every file the skill generates. The skill handles creation; the hook enforces quality.
- **Skills plus subagent delegation**: A complex skill can spawn subagents for independent subtasks. A `/new-service` skill might spawn one agent to generate the service code, another to write tests, and a third to update the CI configuration — all in parallel.
- **Skills plus plan mode**: For large skills, start with plan mode to outline what will be generated, get user confirmation, then execute. This prevents surprises when a skill creates 10+ files.

### Advanced Patterns

- **Forked context with `context: fork`**: Run a skill in an isolated subagent context so it does not pollute your main session's context window. The skill executes, returns a summary, and your main session stays clean. Useful for skills that read many files.
- **Skills with reference docs**: Bundle project-specific reference material (API schemas, style guides, architecture docs) in the skill directory. The skill instructions tell Claude to read these before generating code, ensuring output matches your conventions without relying on the main context.
- **Argument parsing patterns**: Use structured argument formats in your SKILL.md to handle multiple parameters: `/create-endpoint POST /api/users CreateUserRequest`. The skill parses `$ARGUMENTS` into method, path, and request type.

## Common Pitfalls

- **Skills that are too vague**: A SKILL.md that says "create a component" without specifying file locations, naming conventions, or what files to read for context will produce inconsistent results. Be specific about the steps and reference files.
- **Skipping the deterministic parts**: If a step should always produce the same output (timestamps, UUIDs, boilerplate), use a bundled script instead of asking the AI to generate it. AI adds variance where you want consistency.
- **Forgetting to version the skill**: Skills live in your repo. When your project conventions change (new test framework, different directory structure), update the skill or it will generate outdated scaffolding.
- **Overloading a single skill**: A skill that handles creation, updating, deletion, and listing is four skills pretending to be one. Keep each skill focused on one workflow.

## Real-World Example

Your team builds a Django REST API. Every new endpoint requires five files: a serializer, a view, a URL route entry, a test file, and an OpenAPI schema fragment. Developers forget steps, name files inconsistently, and skip the schema update.

You create `/create-endpoint`:

```
> /create-endpoint PATCH /api/v2/orders/{id}/cancel CancelOrderRequest
```

The skill reads `$ARGUMENTS`, splits them into method, path, and request type, then:
1. Runs `./skills/create-endpoint/scaffold.sh` to create empty files with correct names in the right directories
2. Reads `api/serializers/order_serializer.py` to understand existing patterns
3. Generates `api/serializers/cancel_order_serializer.py` matching the existing style
4. Generates the view in `api/views/cancel_order_view.py` with proper permission classes
5. Adds the URL route to `api/urls/v2.py`
6. Creates `tests/api/test_cancel_order.py` with happy path, validation error, and permission denied test cases
7. Appends the endpoint definition to `openapi/paths/orders.yaml`

What used to take 25-30 minutes and a review for forgotten steps now takes 2 minutes and produces consistent output every time.

## Sources

- [Claude Code Skills](https://code.claude.com/docs/en/skills) — Official docs for creating custom skills in .claude/skills/ directories
- [Claude Code Common Workflows](https://code.claude.com/docs/en/common-workflows) — Common workflow tutorials including skill creation patterns
