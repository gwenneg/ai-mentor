# Custom Plugins
*Last verified: 2026-07-06*

## What It Is

Custom Plugins are how you package your own AI workflow — skills, hooks, agent definitions, MCP configuration — into a single versioned unit that anyone can install. What starts as personal configuration in one project's `.claude/` directory becomes `/plugin install your-toolkit` for your whole team, or for the community through a marketplace. Building a plugin is a different activity from using one: you are the author and maintainer, not the consumer.

## Why It Works

Team knowledge compounds only when it's distributable: versioned releases ship as decisions rather than drift, survive laptop changes, and improve in one place for everyone at once.

## When to Use It

- The same skills or hooks are being copy-pasted between projects or teammates
- Onboarding should be "install one package," not "follow this twelve-step config doc"
- You want your team's conventions (review checklists, scaffolding, guardrails) enforced identically everywhere
- A workflow you built is good enough that people outside your team ask for it

## When NOT to Use It

- The configuration serves a single project or just you — standalone `.claude/` files do the job with zero packaging ceremony
- The workflow is still changing daily — package after it stabilizes, or you'll ship churn to your consumers
- You're not prepared to maintain it — an unmaintained plugin that teammates depend on is worse than a wiki page

## How It Works

### Basic (Beginner)

1. Build and iterate the components standalone first — skills in `.claude/skills/`, hooks in settings — until they work reliably.
2. Create the plugin directory with a manifest: `.claude-plugin/plugin.json` holding `name` and `description` (only `plugin.json` goes inside `.claude-plugin/`).
3. Move components to the plugin root: `skills/`, `hooks/hooks.json`, `agents/`, `.mcp.json` — each at the top level of the plugin directory.
4. Test locally: `claude --plugin-dir ./my-plugin`, then `/reload-plugins` as you iterate. Remember your skills are now namespaced: `/my-plugin:release-notes`.
5. Run `claude plugin validate .` and fix what it flags before sharing anything.

### Composing with Other Approaches (Intermediate)

- **Custom plugins plus custom skills**: the natural pipeline — skills prove themselves standalone, then graduate into the plugin. Keep iterating in-place; the plugin is the distribution channel, not the workshop.
- **Custom plugins plus hooks**: the strongest plugins pair invocable skills with automatic hooks (format-on-edit, protected paths). The `hookify` official plugin can mine your conversation history for patterns worth packaging.
- **Custom plugins plus custom agents**: bundle agent definitions in `agents/` to ship judgment, not just procedure — install the plugin, get the team's security reviewer.

### Advanced Patterns

- **Internal marketplace**: a git repo with a marketplace manifest is a private plugin registry — teammates run `/plugin marketplace add your-org/your-marketplace` once and install from it like any marketplace. Keep it in a private repo for internal-only distribution.
- **Community distribution**: submit for review via the in-app form at [platform.claude.com/plugins/submit](https://platform.claude.com/plugins/submit); approved plugins are pinned to a commit SHA in the public `anthropics/claude-plugins-community` catalog and CI bumps the pin as you push.
- **Dependency constraints** (v2.1.110+): declare `dependencies` in `plugin.json` — a bare name tracks the marketplace's latest, while `{ "name": "secrets-vault", "version": "~2.1.0" }` pins a semver range so an upstream breaking release can't move the dependency under you. Resolution works off git tags named `{plugin}--v{version}`; create them with `claude plugin tag --push`. Cross-marketplace dependencies are blocked unless the root marketplace allowlists the source via `allowCrossMarketplaceDependenciesOn`.

## Common Pitfalls

- **Components inside `.claude-plugin/`**: the classic structure mistake — only `plugin.json` lives there; `skills/`, `hooks/`, and `agents/` belong at the plugin root or they silently don't load.
- **Forgetting namespacing when naming**: your skill invocations become `/plugin-name:skill-name`. A plugin named `ai-mentor` with a skill named `ai-mentor` produces the awkward `/ai-mentor:ai-mentor` — name the skill for how the combination reads.
- **Shipping noisy hooks**: a hook that misfires on every edit gets your whole plugin disabled by annoyed teammates. Test hooks thoroughly before they ship to anyone.
- **Ambient versioning**: leaving `version` unset means every push to your marketplace repo is an update. Fine for a personal plugin; for a team, explicit versions make updates reviewable events, and `displayName` gives the listing a human-readable name without changing the install name.

## Sources

- [Create plugins](https://code.claude.com/docs/en/plugins) — Official guide for building plugins, structure, testing, and distribution
- [Plugin marketplaces](https://code.claude.com/docs/en/plugin-marketplaces) — Creating and distributing marketplaces, including private and community options
