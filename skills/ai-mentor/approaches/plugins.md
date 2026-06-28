# Plugins
*Last reviewed: 2026-06-27*

## What It Is

Plugins are distributable packages that bundle skills, hooks, MCP servers, and agent definitions into a single installable unit. Instead of building every automation yourself, you install a plugin that someone else has already built, tested, and maintained. Think of plugins as the package manager layer for your AI coding workflow — the same way you install an npm package instead of writing your own HTTP client.

## Why It Works

Most development workflows are not unique. Code review, commit message generation, test scaffolding, component creation — thousands of teams do these the same way. Plugins capture this shared knowledge into reusable packages so you benefit from collective experience rather than reinventing solutions in isolation. The best automation is the kind you do not have to build.

## When to Use It

- You need a capability that is common across many projects (code review, commit workflows, feature development guides)
- You want to adopt best practices without designing them from scratch
- You are setting up a new project and want productive AI workflows immediately
- You want to distribute your team's internal workflows as a shareable, versioned package

## When NOT to Use It

- Your workflow is highly specific to your project's domain and no existing plugin comes close
- You need fine-grained control over every step — plugins are opinionated by design, and fighting their opinions costs more than building from scratch

## How It Works

### Basic (Beginner)

1. Browse available plugins: check official marketplaces or community registries for plugins that match your needs
2. Install the plugin: run `claude plugin install <plugin-name>` to add it to your project
3. The plugin registers its skills, hooks, and other components automatically
4. Invoke the plugin's skills using the slash commands it provides (e.g., `/commit`, `/review`)
5. The plugin's hooks and MCP servers activate in the background without manual invocation

Example — installing a commit workflow plugin:
```
claude plugin install commit-commands
```
Now `/commit` is available. It reads your staged diff, generates a conventional commit message, lets you edit it, and runs `git commit`.

### Composing with Other Approaches (Intermediate)

- **Plugins plus custom skills**: Install a plugin for the common workflow, then create a custom skill that extends it for your project's needs. For example, install the `code-review` plugin for general review, then build a `/security-review` skill that adds your company's specific security checklist on top.
- **Plugins plus hooks**: A plugin might provide a `/deploy` skill but not enforce pre-deploy checks. Add a PreToolUse hook that runs your integration tests before any deploy command, layering your safety net over the plugin's convenience.
- **Plugins plus agent definitions**: Use a plugin's MCP server for data access, then define a custom agent that uses that server alongside project-specific tools. The plugin provides the plumbing; your agent provides the judgment.

### Advanced Patterns

- **Plugin evaluation before building**: Before writing any custom skill or hook, use the `/plugin` interactive UI to discover if a plugin already exists. A quick browse can save hours of development. Make this a team habit.
- **Forking and customizing plugins**: When a plugin does 80% of what you need, fork it and modify the remaining 20% rather than building from scratch. Most plugins are markdown and scripts — they are straightforward to adapt.
- **Publishing internal plugins**: Package your team's custom skills, hooks, and agents as an internal plugin. New team members install one package and get your entire workflow setup. This scales your team's best practices without relying on onboarding documents that go stale.

## Tool Support

| Tool | Support | Notes |
|------|---------|-------|
| Claude Code | Native | `claude plugin install/uninstall`, official and community marketplaces |
| OpenCode | Native | Full plugin system with npm packages and local loading |
| Cursor | Partial | Plugin and marketplace system since 2026 |
| aider | None | No plugin system |

## Common Pitfalls

- **Installing without reading**: Plugins add hooks and skills that run automatically. A poorly written plugin can slow down your workflow or make unwanted changes. Read the plugin's documentation and review its hooks before installing.
- **Plugin sprawl**: Installing many plugins with overlapping functionality leads to conflicts — two plugins trying to format the same file, or competing commit hooks. Audit your installed plugins periodically and remove redundancies.
- **Assuming plugins are maintained**: Community plugins may go unmaintained. Pin versions for stability and test updates before adopting them across your team.
- **Not checking before building**: The most common waste is spending a day building a custom skill only to discover a mature plugin already does the same thing. Always search first.

## Real-World Example

A team of six is starting a new TypeScript microservice. The tech lead wants productive AI workflows from day one without spending a week configuring everything.

She runs:
```
claude plugin install commit-commands
claude plugin install code-review
claude plugin install feature-dev
```

In 30 seconds, the team has:
- `/commit` — reads the diff, generates a conventional commit message following the project's type prefixes, and handles the git workflow
- `/review` — runs a multi-pass code review checking for correctness bugs, security issues, and style violations
- `/feature` — guides developers through feature implementation with a structured plan-then-execute workflow

Two weeks later, a developer notices the `code-review` plugin does not check for their company's required error handling patterns. Instead of replacing the plugin, the team creates a custom `/error-review` skill that focuses specifically on their error handling conventions. The plugin handles the general case; the custom skill handles the specific one.

The team later packages their custom skills and project-specific hooks into an internal plugin called `@mycompany/service-workflows`. When they spin up their next microservice, one `claude plugin install` command gives the new project the same workflow setup.

## Sources

- [Claude Code Plugins](https://docs.anthropic.com/en/docs/claude-code/plugins) — Official docs for creating and distributing plugins
- [Claude Code Plugins README](https://github.com/anthropics/claude-code/blob/main/plugins/README.md) — Plugin documentation in the official claude-code GitHub repository
