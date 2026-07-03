# Official Claude Code Plugins Catalog
*Last synced: 2026-07-02 · Source: [`anthropics/claude-plugins-official`](https://github.com/anthropics/claude-plugins-official) · Evaluation pass: 2026-07-03 (all 52 desk-checked; 15 exercised hands-on)*

All plugins below are in the official marketplace and installable via `/plugin install <name>@claude-plugins-official`. None are installed by default. The repo contains 37 Anthropic-built plugins and 15 externally-maintained plugins; this catalog is the full list, kept in sync by the maintenance skill's catalog-sync step. Scope decision (2026-07-03): externally-maintained plugins listed in the official marketplace ARE in scope — "official" means Anthropic-curated, not Anthropic-authored — and both directories stay synced. The goal column names the routing-table section (`routing.md`) each plugin maps to.

Verdicts are produced by the repeatable protocol in `evals/plugin-evaluation.md` — same fixture, same per-plugin exercises, same criteria on every run, so evaluations stay comparable over time.

**Verdict legend** — every plugin carries one:

- ✅ **hands-on (date)** — installed, exercised against its mapped goal, and it worked; caveats noted verbatim
- ☑️ **desk-checked** — manifest, components, freshness, and provenance reviewed (2026-07-03); not exercised. For MCP integrations this usually means hands-on needs an external account or infrastructure we don't have — an honest label, not a defect
- ⚠️ **caution** — works, but overlaps a built-in feature or has a sharp edge; lead with the alternative named

The mentor recommends ✅ plugins freely, offers ☑️ ones with the "not hands-on evaluated" label, and never presents a ⚠️ without its caveat.

## Anthropic-built plugins

### Dev workflow

| Plugin | What it does | Relevant goal | Verdict |
|--------|-------------|--------------|---------|
| `security-guidance` | Per-edit security hooks + Stop-time LLM diff review (12 hooks, ~0 always-on tokens) | `security` | ✅ hands-on 2026-07-03 — injection attempt produced hardened parameterized code; invisible when quiet (complements on-demand `/security-review`) |
| `hookify` | Creates hooks from conversation patterns or explicit rules | `ci-automation` | ✅ hands-on 2026-07-03 — generated a working PostToolUse hook, verified firing; headless caveat: can't write settings files non-interactively |
| `feature-dev` | 7-phase guided feature development with explorer/architect/reviewer agents | `greenfield` | ✅ hands-on 2026-07-03 (start verified) — phased flow engages correctly, scales down sensibly on small repos; overlaps plan mode, packaged as one pipeline |
| `commit-commands` | `/commit`, `/commit-push-pr`, `/clean_gone` git workflow commands | `release-management` | ✅ hands-on 2026-07-03 — flawless first try; ⚠️ mostly duplicates native committing — value is team commit conventions and `clean_gone` |
| `code-review` | Multi-agent PR review with confidence scoring | `code-review` | ⚠️ duplicates the built-in `/code-review`, `/review`, and `/code-review ultra` — recommend the built-ins first |
| `pr-review-toolkit` | 6-agent review covering comments, tests, types, error handling, simplification | `code-review` | ✅ hands-on 2026-07-03 — found a planted off-by-one at the exact line with a verified repro and flagged the deliberate test gap; token-hungry (~2k always-on + multiple subagents); overlaps built-in reviews but adds comment/test-coverage/type-design angles |
| `code-modernization` | Structured migration of legacy codebases (COBOL, legacy Java/C++, monoliths) | `migration` | ✅ hands-on 2026-07-03 (start verified) — preflight phase engages coherently; needs a generous turn budget and Bash allowlisting (its multi-command probes fragment under default permissions); biggest component surface in the catalog |
| `code-simplifier` | Agent for clarity and maintainability refactors | `refactoring` | ⚠️ overlaps the built-in `/simplify` skill — recommend the built-in first |
| `frontend-design` | Auto-invoked skill for bold, production-grade UI design | `greenfield` | ✅ hands-on 2026-07-03 — auto-engaged (invocation observed directly in transcript) and produced a branded page in 4 turns; caveat: its "self-contained" output included a Google Fonts link |
| `ralph-loop` | Continuous while-true agent loops re-running the same prompt until completion | `migration` | ⚠️ overlaps the built-in `/loop` and `/goal` — recommend the built-ins first |
| `playground` | Interactive single-file HTML playgrounds with visual controls and live preview | `greenfield` | ☑️ desk-checked — partially overlaps the built-in Artifact tool for shareable pages |

### Hooks & output styles

| Plugin | What it does | Relevant goal | Verdict |
|--------|-------------|--------------|---------|
| `explanatory-output-style` | SessionStart hook injecting educational insights about implementation choices | `onboarding` | ☑️ desk-checked — mimics a deprecated output style; niche |
| `learning-output-style` | Prompts users to write meaningful code at decision points | `onboarding` | ☑️ desk-checked — mimics an unshipped output style; niche |

### Plugin & SDK development

| Plugin | What it does | Relevant goal | Verdict |
|--------|-------------|--------------|---------|
| `skill-creator` | Creates and improves skills, measures skill performance | `building-skills-plugins` | ☑️ desk-checked — in active daily use by this catalog's maintainer, which is stronger evidence than most desk checks |
| `plugin-dev` | 8-phase guided workflow for building plugins, with validator and reviewer agents | `building-skills-plugins` | ✅ hands-on 2026-07-03 — scaffolded a plugin that passed `claude plugin validate` and self-reviewed honestly; entry point is `create-plugin`; heaviest always-on context of the evaluated set (~2.3k tokens) |
| `mcp-server-dev` | Guided MCP server design and implementation | `building-mcp-integrations` | ✅ hands-on 2026-07-03 — produced a syntax-clean stdio server with current SDK idioms (registerTool, zod validation, stdout hygiene) plus both config snippets; the SDK-idiom guidance is the value over base Claude |
| `agent-sdk-dev` | Scaffolds Agent SDK projects, validates against best practices | `building-agents` | ✅ hands-on 2026-07-03 — sane strict-TS scaffold with streaming `query()`; pins deps to `latest` when the registry is unreachable, and its verifier agents only work after `npm install` |
| `mcp-tunnels` | Connects Claude to a private MCP server through an Anthropic MCP tunnel | `building-mcp-integrations` | ☑️ desk-checked — needs Docker Compose infrastructure to exercise |
| `example-plugin` | Reference plugin demonstrating every extension surface | `building-skills-plugins` | ☑️ desk-checked — reference material, not a workflow tool |

### Project & session management

| Plugin | What it does | Relevant goal | Verdict |
|--------|-------------|--------------|---------|
| `claude-md-management` | Audits and maintains CLAUDE.md files | `documentation` | ✅ hands-on 2026-07-03 — scored audit (rubric + real gaps found, cross-checked against the codebase); note the skill is invoked as `claude-md-improver` |
| `claude-code-setup` | Analyzes a codebase and recommends tailored Claude Code automations | `onboarding` | ✅ hands-on 2026-07-03 — recommendations were concretely repo-tailored (justified each hook from real project facts, declined unjustified MCP servers); conceptually overlaps this plugin's own growth mode |
| `session-report` | Generates an HTML report of session token usage and cache efficiency | `devops` | ✅ hands-on 2026-07-03 — self-contained HTML with real usage numbers; cheapest always-on cost (~70 tokens) but needs >12 turns and default permissions block its bundled analyzer; reports a 7-day window, not strictly the current session |
| `project-artifact` | Publishes a living project status page with workstreams and decisions | `documentation` | ✅ hands-on 2026-07-03 — produced a project-specific tabbed status page with honest unverified-state markings; publishing needs an interactive claude.ai session (headless falls back to a local HTML file + refresh config) |

### Language servers (LSPs)

Drop-in LSP integrations for code intelligence: `clangd-lsp` (C/C++), `csharp-lsp`, `gopls-lsp` (Go), `jdtls-lsp` (Java), `kotlin-lsp`, `lua-lsp`, `php-lsp`, `pyright-lsp` (Python), `ruby-lsp`, `rust-analyzer-lsp`, `swift-lsp`, `typescript-lsp`.

☑️ desk-checked as a family — uniform official wrappers around standard language servers, low risk; each requires its language-server binary on `$PATH` (the plugin errors visibly if missing). Recommend freely when the user's language matches and the binary exists or is easily installed.

### Specialty

Rarely relevant to everyday engineering, listed for completeness (☑️ desk-checked): `math-olympiad` (competition math solving with adversarial proof verification) and `cwc-makers` (onboarding for the Code-with-Claude Makers Cardputer hardware kit).

## External plugins (partner-maintained)

Hands-on evaluation of most integrations requires accounts or infrastructure (Slack workspaces, Figma files, cloud projects); those carry ☑️ with that caveat rather than a pretend verdict.

| Plugin | What it does | Relevant goal | Verdict |
|--------|-------------|--------------|---------|
| `context7` (Upstash) | Pulls version-pinned documentation for any library on demand | `code-understanding` | ✅ hands-on 2026-07-03 — returned real Express v5 docs, no account needed; headless callers must allowlist the MCP server |
| `github` (GitHub) | Official GitHub MCP: issues, PRs, code review, repo management | `code-review` | ☑️ desk-checked — first-party GitHub; needs repo auth to exercise |
| `gitlab` (GitLab) | GitLab MCP: merge requests, CI/CD, pipelines, issues | `ci-automation` | ☑️ desk-checked — first-party GitLab; needs instance auth |
| `playwright` (Microsoft) | Browser automation and E2E testing MCP server | `testing` | ☑️ desk-checked — first-party Microsoft; needs browser install; note the built-in Chrome integration covers some of this |
| `serena` (Oraios) | Semantic code analysis MCP for refactoring and code understanding | `code-understanding` | ☑️ desk-checked — note built-in LSP plugins cover much of the navigation ground |
| `greptile` (Greptile) | AI PR review agent for GitHub and GitLab | `code-review` | ☑️ desk-checked — needs a Greptile account; overlaps built-in review skills |
| `linear` (Linear) | Linear issue tracking: create issues, manage projects, search | `devops` | ☑️ desk-checked — needs workspace auth |
| `asana` (Asana) | Create and manage tasks, search projects, update assignments | `devops` | ☑️ desk-checked — needs workspace auth |
| `firebase` (Google) | Firestore, auth, cloud functions, and hosting via Firebase MCP | `devops` | ☑️ desk-checked — needs a Firebase project |
| `terraform` (HashiCorp) | Terraform MCP for IaC registry integration and module management | `devops` | ☑️ desk-checked — first-party HashiCorp, fresh (2026-03) |
| `laravel-boost` (Laravel) | Laravel development toolkit MCP server | `greenfield` | ☑️ desk-checked — first-party Laravel; needs a Laravel app |
| `telegram` | Telegram messaging bridge with access control (channels) | `devops` | ☑️ desk-checked — fresh (2026-04); covered by `approaches/channels.md`; needs a bot token |
| `discord` | Discord messaging bridge with access control (channels) | `devops` | ☑️ desk-checked — fresh (2026-04); covered by `approaches/channels.md`; needs a bot |
| `imessage` | iMessage bridge (reads `chat.db`, sends via AppleScript; channels) | `devops` | ☑️ desk-checked — macOS only, Full Disk Access required; covered by `approaches/channels.md` |
| `fakechat` | Localhost chat UI for testing channel flows — no tokens, no access control | `testing` | ☑️ desk-checked — the intended channels demo; requires Bun |
