# ai-mentor

Describe your engineering problem, get the right AI workflow. Ranked, verified recommendations with ready-to-paste prompts that teach you why each approach works — not just how.

## What it does

- Matches engineering problems to ranked AI workflow approaches
- Adapts recommendations to your experience level (getting started, comfortable, advanced)
- Teaches the *why* behind each approach, not just the mechanics
- Provides concrete "try it now" prompts for every recommendation
- Covers Claude Code (primary), OpenCode, Cursor, and aider
- Every approach backed by verified official sources

## Install

ai-mentor is distributed through a Claude Code [plugin marketplace](https://code.claude.com/docs/en/plugin-marketplaces) — no cloning or file editing. In Claude Code:

```
/plugin marketplace add gwenneg/claude-ichiba
/plugin install ai-mentor@claude-ichiba
/reload-plugins
```

## Staying up to date

Auto-update is off and Claude Code sends no new-version notification — **watch this repo → Releases only**. To update ([plugin docs](https://code.claude.com/docs/en/discover-plugins)):

```
/plugin marketplace update claude-ichiba
/reload-plugins
```

Releases are pinned to an immutable commit SHA so work-in-progress on `main` never reaches you.

## Usage

### Describe your problem

```
/ai-mentor debug a flaky test
/ai-mentor review a large PR
/ai-mentor refactor authentication across 30 files
```

### Guided discovery

```
/ai-mentor
```

The skill asks what you're working on, identifies the best-matching goal, calibrates to your experience level, and presents ranked approaches.

### Auto-triggered

The skill can also be invoked automatically by Claude when it detects you're working on a task that has a known AI workflow approach.

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

## AI workflow approaches

Each approach file covers: what it is, why it works, when to use it (and when not to), beginner through advanced usage, tool support, common pitfalls, a real-world example, and verified official sources.

Approaches include: Plan Mode, Autonomous Loops, Subagent Delegation, Worktree Isolation, Fan-Out Workflows, Deep Research, Browser Integration, Headless Mode, MCP Context, Checkpoints & Rewind, Multi-Provider, LSP Self-Correction, Built-in Review Skills, Hooks as Workflow, Custom Skills, Plugins, and Custom Agents.

## License

Apache-2.0
