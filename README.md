# ai-mentor

Learn the Claude Code capabilities you didn't know existed. ai-mentor is a discovery mentor: it reads your setup and your session, remembers what you already use, and teaches you the highest-leverage thing you're missing — one demonstrated move at a time, grounded in your actual repo.

## What it does

- **Finds your unknown unknowns** — computes the gap between what Claude Code offers and what you demonstrably use (your hooks, agents, plugins, MCP config, session habits), and teaches from the top of that gap
- **Diagnoses instead of interrogating** — reads the conversation and your repo before recommending; no questionnaires
- **One move, demonstrated** — a diagnosis, one recommended move with a ready-to-paste prompt built from your real paths and commands, and an offer to set it up on the spot (say "more" for the full ranked list)
- **One personalized surprise, every time** — each interaction carries a capability you didn't know to ask about, chosen for *you*, not from a static list
- **Never repeats itself** — a local profile records what you've been shown, adopted, or declined; nothing is taught twice, and a "no" is remembered
- **Teaches the why** — every recommendation names the principle that makes it work, backed by verified official sources

## The profile

The mentor keeps one small markdown file at `~/.ai-mentor/profile.md`: one line per capability (shown / adopted / declined). It's machine-local, never committed, never uploaded, and yours to edit or delete — a hand edit always wins over anything the mentor inferred. It requires no setup and no permission prompts.

## Install

ai-mentor is distributed through a Claude Code [plugin marketplace](https://code.claude.com/docs/en/plugin-marketplaces) — no cloning or file editing. In Claude Code:

```
/plugin marketplace add gwenneg/claude-ichiba
/plugin install ai-mentor@claude-ichiba
/reload-plugins
```

## Staying up to date

Third-party marketplaces have auto-update disabled by default ([plugin docs](https://code.claude.com/docs/en/discover-plugins#configure-auto-updates)). Either enable it once — `/plugin` → **Marketplaces** → `claude-ichiba` → **Enable auto-update** — and Claude Code will refresh at startup and notify you when the plugin updates, or update manually:

```
/plugin marketplace update claude-ichiba
/reload-plugins
```

Releases are pinned to an immutable commit SHA so work-in-progress on `main` never reaches you.

## Usage

### Problem mode — you have something to solve

```
/ai-mentor:mentor debug a flaky test
/ai-mentor:mentor review a large PR
/ai-mentor:mentor refactor authentication across 30 files
```

You get a diagnosis grounded in your repo, one recommended move with a ready-to-run prompt, and one capability you probably didn't know about. Say "more" for the full ranked list, or name any approach to go deeper.

### Growth mode — teach me something

```
/ai-mentor:mentor
```

No problem needed. The mentor checks your setup and profile, then teaches the single highest-leverage capability you're not using — or follows up on the last thing it showed you, or surfaces what shipped since you last checked.

### Auto-triggered

Ask anything mentor-shaped in a normal session — "what's the best way to use AI for this?" — and Claude invokes the mentor itself. If it notices you grinding through something a capability would remove, it may offer to help: permission first, at most once per session, and a "no thanks" is remembered.

## Problem categories

| Category | Example problems |
|----------|-----------------|
| Debugging | Stack traces, flaky tests, runtime errors |
| Code review | PR review, security review, quality gates |
| Refactoring | Cross-file changes, codemods, tech debt cleanup |
| Greenfield | New features, design decisions, prototyping |
| Testing | Test generation, coverage gaps, E2E testing |
| Code understanding | Architecture discovery, "how does this work?" |
| Research | Comparing approaches, technical due diligence |
| Migration | Framework upgrades, API migrations, dependency updates |
| Documentation | API docs, architecture docs, onboarding guides |
| CI/CD | Automated reviews, scheduled tasks, pipeline integration |
| Performance | Profiling, latency, memory, bundle size optimization |
| Security | Vulnerability audits, hardening, compliance |
| Incident response | Production outages, error spikes, emergency fixes |
| Onboarding | New hire setup, team rotation, environment setup |
| Dependency management | Library evaluation, updates, supply chain health |
| API design | Endpoint design, schemas, versioning, contracts |
| Release management | Changelogs, version bumps, deployment coordination |
| DevOps | Terraform, Kubernetes, Docker, cloud infrastructure |
| Tech debt | Code quality audits, cleanup prioritization |
| Accessibility | WCAG compliance, screen readers, keyboard navigation |
| Building AI agents | Agent design, prototyping, Agent SDK products |
| Building MCP integrations | MCP servers, exposing internal tools to AI |
| Building skills & plugins | Packaging team workflows, marketplace distribution |
| Building LLM features | AI product features, prompt engineering, evals |

## AI workflow approaches

Each approach file covers: what it is, why it works, when to use it (and when not to), beginner through advanced usage, common pitfalls, a real-world example, and verified official sources.

Approaches include: Plan Mode, Autonomous Loops, Subagent Delegation, Worktree Isolation, Fan-Out Workflows, Deep Research, Browser Integration, Headless Mode, MCP Context, Checkpoints & Rewind, Model & Effort Selection, LSP Self-Correction, Built-in Review Skills, Hooks as Workflow, Custom Skills, Official Plugins, Custom Plugins, Custom Agents, Project Memory & Context Docs, Session & Context Management, Background Agents, Scheduled & Recurring Agents, Cloud Sessions & Remote Work, Permissions & Safe Autonomy, Visual Artifacts, and Channels.

## License

Apache-2.0
