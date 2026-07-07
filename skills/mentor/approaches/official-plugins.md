# Official Plugins
*Last verified: 2026-07-06*

## What It Is

Official Plugins are ready-made packages from Anthropic's curated marketplace — a few dozen Anthropic-built plugins plus a large and fast-growing set of partner-maintained integrations (GitHub, GitLab, Figma, Datadog, Stripe, and hundreds more) — that bundle skills, hooks, MCP servers, and agent definitions into a single installable unit. Instead of building an automation yourself, you install one that has already been built, tested, and maintained. Think of it as the package manager layer for your AI workflow: the same way you install an npm package instead of writing your own HTTP client.

## Why It Works

Most development workflows are not unique. Code review, commit message generation, feature scaffolding, language-server integration — thousands of teams do these the same way, and the official marketplace captures that shared knowledge into packages curated by Anthropic. Installing beats building on three axes at once: zero development time, someone else's testing, and someone else's maintenance as Claude Code evolves. The best automation is the kind you do not have to build.

## When to Use It

- You need a capability that is common across many projects (code review, commit workflows, feature development guides, LSP integration, guided MCP server and Agent SDK development)
- Your task centers on a specific product or stack — clouds, databases, observability, design tools, frameworks — where a vendor-built integration ships the connection ready-made
- You want to adopt best practices without designing them from scratch
- You are setting up a new project and want productive AI workflows immediately
- You are about to build something custom — check first whether it already exists

## When NOT to Use It

- Your workflow is highly specific to your project's domain and no existing plugin comes close — package your own instead (see Custom Plugins)
- You need fine-grained control over every step — plugins are opinionated by design, and fighting their opinions costs more than building from scratch

## How It Works

### Basic (Beginner)

1. Find the right plugin: browse the marketplace interactively with the `/plugin` UI; the bundled catalog at `references/official-plugins.md` maps the plugins evaluated so far to the goal each one serves.
2. Install it: `/plugin install <name>@claude-plugins-official` (nothing is installed by default).
3. Run `/reload-plugins` (or restart Claude Code) to activate it — the plugin's skills, hooks, agents, and MCP servers then register.
4. Invoke its skills via namespaced commands — plugin skills are always prefixed with the plugin name, e.g. `/commit-commands:commit`.
5. Hooks and MCP servers activate in the background with no manual invocation.

### Composing with Other Approaches (Intermediate)

- **Official plugins plus custom skills**: install the plugin for the general case, then add a project-specific skill on top — e.g. the `pr-review-toolkit` plugin for general review, plus your own skill carrying the company security checklist.
- **Official plugins plus hooks**: a plugin might provide a `/deploy` skill but not enforce pre-deploy checks. Layer a PreToolUse hook that runs your integration tests before any deploy command — your safety net over their convenience.
- **Official plugins plus custom agents**: use a plugin's MCP server for data access, and define your own agent that wields it with project-specific judgment. The plugin provides the plumbing; your agent provides the policy.

### Advanced Patterns

- **Evaluate before building**: make "check the marketplace first" a team habit — the most common waste is spending a day building a custom skill that a mature plugin already covers.
- **Beyond the official marketplace**: the community marketplace (`/plugin marketplace add anthropics/claude-plugins-community`) carries third-party plugins that have passed Anthropic's automated validation and safety screening, each pinned to a specific commit SHA in the catalog — a second place to look before building.
- **Deliberate updates**: official Anthropic marketplaces auto-update at startup by default. Where you need control, disable auto-update in the `/plugin` Marketplaces tab, refresh on your own schedule with `/plugin marketplace update <marketplace>`, and test updated hooks before rolling them across the team.

## Common Pitfalls

- **Installing without reading**: plugins add hooks and skills that run automatically. A poorly written plugin can slow down your workflow or make unwanted changes. Read the plugin's documentation and review its hooks before installing.
- **Plugin sprawl**: many plugins with overlapping functionality lead to conflicts — two plugins formatting the same file, competing commit hooks. Audit your installed plugins periodically and remove redundancies.
- **Assuming plugins are maintained**: partner-maintained plugins evolve on the partner's schedule. Disable marketplace auto-update where stability matters, and test updates before adopting them across your team.
- **Expecting un-namespaced commands**: installing `commit-commands` gives you `/commit-commands:commit`, not `/commit`. Plugin skills always carry the plugin prefix — factor that into your muscle memory and docs.

## Real-World Example

A team of six is starting a new TypeScript microservice. The tech lead wants productive AI workflows from day one without spending a week configuring everything.

She runs:
```
claude plugin install commit-commands@claude-plugins-official
claude plugin install pr-review-toolkit@claude-plugins-official
claude plugin install feature-dev@claude-plugins-official
```

In 30 seconds, the team has:
- `/commit-commands:commit` — reads the diff, generates a conventional commit message following the project's type prefixes, and handles the git workflow
- `pr-review-toolkit` — six specialized review agents covering correctness, test thoroughness, error handling, and type design, triggered by asking for a review (e.g. "review this PR's error handling")
- `/feature-dev:feature-dev` — guides developers through feature implementation with a structured plan-then-execute workflow

Two weeks later, a developer notices the `pr-review-toolkit` plugin does not check for their company's required error handling patterns. Instead of replacing the plugin, the team creates a custom `/error-review` skill focused on their conventions. The plugin handles the general case; the custom skill handles the specific one — and nobody spent a day rebuilding what the marketplace already provided.

## Sources

- [Discover and install plugins](https://code.claude.com/docs/en/discover-plugins) — Official docs for browsing marketplaces and installing plugins
- [Create plugins](https://code.claude.com/docs/en/plugins) — Official plugin documentation including the official and community marketplaces
