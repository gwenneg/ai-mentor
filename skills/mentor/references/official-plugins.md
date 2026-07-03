# Official Claude Code Plugins Catalog
*Last synced: 2026-07-02 · Source: [`anthropics/claude-plugins-official`](https://github.com/anthropics/claude-plugins-official)*

<!-- TODO(evaluation): descriptions below are self-reported by each plugin — no value assessment has been done yet.
     Plan: evaluate each plugin hands-on for its mapped goal (works out of the box, doesn't duplicate a built-in
     feature, maintenance/provenance, context cost of its hooks/skills) and add a Recommended tier. The mentor
     skill should then suggest only evaluated-valuable plugins and treat the rest as listed-for-completeness. -->

All plugins below are in the official marketplace and installable via `/plugin install <name>@claude-plugins-official`. None are installed by default. The repo contains 37 Anthropic-built plugins and 15 externally-maintained plugins; this catalog is the full list, kept in sync by the maintenance skill's catalog-sync step. Scope decision (2026-07-03): externally-maintained plugins listed in the official marketplace ARE in scope — "official" means Anthropic-curated, not Anthropic-authored — and both directories stay synced. Goal paths are relative to the skill root.

## Anthropic-built plugins

### Dev workflow

| Plugin | What it does | Relevant goal |
|--------|-------------|--------------|
| `code-review` | Multi-agent PR review with confidence scoring | `goals/code-review.md` |
| `pr-review-toolkit` | 6-agent review covering comments, tests, types, error handling, simplification | `goals/code-review.md` |
| `commit-commands` | `/commit`, `/commit-push-pr`, `/clean_gone` git workflow commands | `goals/release-management.md` |
| `feature-dev` | 7-phase guided feature development with codebase explorer and architect agents | `goals/greenfield.md` |
| `code-modernization` | Structured migration of legacy codebases (COBOL, legacy Java/C++, monoliths) | `goals/migration.md` |
| `code-simplifier` | Agent for clarity and maintainability refactors | `goals/refactoring.md` |
| `frontend-design` | Auto-invoked skill for bold, production-grade UI design | `goals/greenfield.md` |
| `security-guidance` | PreToolUse hook flagging 9 security patterns on every file edit | `goals/security.md` |
| `ralph-loop` | Continuous while-true agent loops ("Ralph Wiggum technique") — re-runs the same prompt until task completion | `goals/migration.md` |
| `playground` | Interactive single-file HTML playgrounds with visual controls and live preview | `goals/greenfield.md` |

### Hooks & output styles

| Plugin | What it does | Relevant goal |
|--------|-------------|--------------|
| `hookify` | Creates hooks from conversation patterns or explicit rules | `goals/ci-automation.md` |
| `explanatory-output-style` | SessionStart hook injecting educational insights about implementation choices | `goals/onboarding.md` |
| `learning-output-style` | Prompts users to write meaningful code at decision points | `goals/onboarding.md` |

### Plugin & SDK development

| Plugin | What it does | Relevant goal |
|--------|-------------|--------------|
| `plugin-dev` | 8-phase guided workflow for building plugins, with validator and reviewer agents | `goals/building-skills-plugins.md` |
| `mcp-server-dev` | Guided MCP server design and implementation | `goals/building-mcp-integrations.md` |
| `agent-sdk-dev` | Scaffolds Agent SDK projects, validates against best practices | `goals/building-agents.md` |
| `skill-creator` | Creates and improves skills, measures skill performance | `goals/building-skills-plugins.md` |
| `mcp-tunnels` | Connects Claude to a private MCP server through an Anthropic MCP tunnel (Docker Compose quickstart: certificates, proxy, cloudflared) | `goals/building-mcp-integrations.md` |
| `example-plugin` | Reference plugin demonstrating every extension surface: commands, agents, skills, hooks, and MCP servers | `goals/building-skills-plugins.md` |

### Project & session management

| Plugin | What it does | Relevant goal |
|--------|-------------|--------------|
| `claude-code-setup` | Analyzes a codebase and recommends tailored Claude Code automations | `goals/onboarding.md` |
| `claude-md-management` | Audits and maintains CLAUDE.md files | `goals/documentation.md` |
| `session-report` | Generates an HTML report of session token usage and cache efficiency | `goals/devops.md` |
| `project-artifact` | Publishes a living project status page with workstreams and decisions | `goals/documentation.md` |

### Language servers (LSPs)

Drop-in LSP integrations for code intelligence: `clangd-lsp` (C/C++), `csharp-lsp`, `gopls-lsp` (Go), `jdtls-lsp` (Java), `kotlin-lsp`, `lua-lsp`, `php-lsp`, `pyright-lsp` (Python), `ruby-lsp`, `rust-analyzer-lsp`, `swift-lsp`, `typescript-lsp`.

### Specialty

Rarely relevant to everyday engineering, listed for completeness: `math-olympiad` (competition math solving with adversarial proof verification) and `cwc-makers` (onboarding for the Code-with-Claude Makers Cardputer hardware kit).

## External plugins (partner-maintained)

| Plugin | What it does | Relevant goal |
|--------|-------------|--------------|
| `asana` | Create and manage tasks, search projects, update assignments | `goals/devops.md` |
| `context7` | Pulls version-specific documentation for any library on demand | `goals/code-understanding.md` |
| `discord` | Discord messaging bridge for Claude Code with access control | `goals/devops.md` |
| `fakechat` | Localhost iMessage-style chat for testing notification flows | `goals/testing.md` |
| `firebase` | Firestore, auth, cloud functions, and hosting via Firebase MCP | `goals/devops.md` |
| `github` | Official GitHub MCP: issues, PRs, code review, repo management | `goals/code-review.md` |
| `gitlab` | GitLab MCP: merge requests, CI/CD, pipelines, issues | `goals/ci-automation.md` |
| `greptile` | AI PR review agent for GitHub and GitLab | `goals/code-review.md` |
| `imessage` | iMessage bridge for Claude Code (reads `chat.db`, sends via AppleScript) | `goals/devops.md` |
| `laravel-boost` | Laravel development toolkit MCP server | `goals/greenfield.md` |
| `linear` | Linear issue tracking: create issues, manage projects, search | `goals/devops.md` |
| `playwright` | Browser automation and E2E testing MCP server by Microsoft | `goals/testing.md` |
| `serena` | Semantic code analysis MCP for refactoring and code understanding | `goals/code-understanding.md` |
| `telegram` | Telegram messaging bridge for Claude Code with access control | `goals/devops.md` |
| `terraform` | Terraform MCP for IaC registry integration and module management | `goals/devops.md` |
