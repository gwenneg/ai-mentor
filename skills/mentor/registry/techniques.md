# Capability Registry — techniques
*Last verified: 2026-07-06*

Every catalog approach as one flat record — the technique slice of the kind-aware capability registry (built-in commands live in `builtin-commands.md`, marketplace plugins in `references/official-plugins.md`, which IS the registry's plugin slice). Records are generated from the routing files and approach files and kept in lockstep by the structural audit: `goals` mirrors routing membership exactly, `id` is the approach basename, and the deep-dive prose stays in `approaches/<id>.md` — this file adds the machine-readable goal membership and one-line trigger, nothing else. `best_when` is the highest-ranked routing row's trigger. No numeric fit scores, ever.

## autonomous-loops

id: autonomous-loops
kind: technique
goals: accessibility, debugging, greenfield, llm-features, migration, performance, refactoring, testing
best_when: coverage must reach a specific threshold to merge
setup: some
composes_with: plan-mode, worktree-isolation
source: approaches/autonomous-loops.md

## background-agents

id: background-agents
kind: technique
goals: tech-debt
best_when: audit takes an hour but your attention shouldn't
setup: some
composes_with: autonomous-loops, worktree-isolation, session-context-management
source: approaches/background-agents.md

## browser-integration

id: browser-integration
kind: technique
goals: accessibility, research
best_when: verifying tab order, focus traps, or screen reader behavior
setup: involved
composes_with: autonomous-loops, plan-mode, worktree-isolation
source: approaches/browser-integration.md

## built-in-review-skills

id: built-in-review-skills
kind: technique
goals: api-design, code-review, security, tech-debt
best_when: quick review of a focused diff or your own pre-PR code
setup: none
composes_with: plan-mode
source: approaches/built-in-review-skills.md

## channels

id: channels
kind: technique
goals: incident-response
best_when: alerts should land in the session with context already loaded
setup: involved
composes_with: background-agents, safe-autonomy, mcp-context
source: approaches/channels.md

## checkpoints-rewind

id: checkpoints-rewind
kind: technique
goals: refactoring
best_when: risky refactor you might need to undo
setup: none
composes_with: plan-mode, worktree-isolation, session-context-management
source: approaches/checkpoints-rewind.md

## cloud-sessions

id: cloud-sessions
kind: technique
goals: code-review
best_when: deep review off your machine, or PRs that fix themselves
setup: some
composes_with: plan-mode
source: approaches/cloud-sessions.md

## custom-agents

id: custom-agents
kind: technique
goals: building-agents, ci-automation, tech-debt
best_when: want a working prototype today
setup: none
composes_with: worktree-isolation, hooks-as-workflow
source: approaches/custom-agents.md

## custom-plugins

id: custom-plugins
kind: technique
goals: building-skills-plugins
best_when: skill works; time to share or distribute
setup: some
composes_with: custom-skills, hooks-as-workflow, custom-agents
source: approaches/custom-plugins.md

## custom-skills

id: custom-skills
kind: technique
goals: building-skills-plugins, documentation, greenfield, onboarding, release-management
best_when: a workflow you repeat and want to package
setup: none
composes_with: hooks-as-workflow, subagent-delegation, plan-mode
source: approaches/custom-skills.md

## deep-research

id: deep-research
kind: technique
goals: accessibility, api-design, code-understanding, dependency-management, devops, documentation, greenfield, incident-response, llm-features, onboarding, performance, research, security
best_when: evaluating a library you haven't used before
setup: none
composes_with: plan-mode, autonomous-loops, mcp-context
source: approaches/deep-research.md

## fan-out-workflows

id: fan-out-workflows
kind: technique
goals: ci-automation, code-review, testing
best_when: multi-step pipeline with verification between stages
setup: involved
composes_with: worktree-isolation, plan-mode, built-in-review-skills
source: approaches/fan-out-workflows.md

## headless-mode

id: headless-mode
kind: technique
goals: building-mcp-integrations, ci-automation, release-management
best_when: running Claude in GitHub Actions or GitLab CI
setup: some
composes_with: autonomous-loops, built-in-review-skills, fan-out-workflows
source: approaches/headless-mode.md

## hooks-as-workflow

id: hooks-as-workflow
kind: technique
goals: accessibility, building-skills-plugins, debugging, performance, security, testing
best_when: the behavior should be automatic, not invoked
setup: some
composes_with: autonomous-loops, worktree-isolation, subagent-delegation
source: approaches/hooks-as-workflow.md

## lsp-self-correction

id: lsp-self-correction
kind: technique
goals: code-understanding
best_when: tracing how components connect across a large codebase
setup: some
composes_with: autonomous-loops, plan-mode, built-in-review-skills
source: approaches/lsp-self-correction.md

## mcp-context

id: mcp-context
kind: technique
goals: api-design, building-mcp-integrations, code-review, code-understanding, dependency-management, documentation, incident-response, onboarding
best_when: existing docs, specs, and decisions are scattered
setup: some
composes_with: plan-mode, subagent-delegation, built-in-review-skills
source: approaches/mcp-context.md

## model-effort-selection

id: model-effort-selection
kind: technique
goals: research
best_when: different phases need different reasoning depth
setup: involved
composes_with: subagent-delegation, plan-mode, custom-agents
source: approaches/model-effort-selection.md

## official-plugins

id: official-plugins
kind: technique
goals: building-agents, building-mcp-integrations
best_when: ready to implement
setup: some
composes_with: custom-skills, hooks-as-workflow, custom-agents
source: approaches/official-plugins.md

## plan-mode

id: plan-mode
kind: technique
goals: api-design, building-agents, building-mcp-integrations, code-understanding, debugging, dependency-management, devops, documentation, greenfield, incident-response, llm-features, migration, onboarding, performance, refactoring, release-management, research, tech-debt
best_when: starting a new API or major version from scratch
setup: none
composes_with: subagent-delegation, built-in-review-skills, worktree-isolation
source: approaches/plan-mode.md

## project-memory

id: project-memory
kind: technique
goals: code-understanding
best_when: want the map you built to persist across sessions
setup: none
composes_with: plan-mode, hooks-as-workflow, subagent-delegation
source: approaches/project-memory.md

## safe-autonomy

id: safe-autonomy
kind: technique
goals: building-agents
best_when: need to constrain what the agent can do
setup: some
composes_with: autonomous-loops, hooks-as-workflow
source: approaches/safe-autonomy.md

## scheduled-agents

id: scheduled-agents
kind: technique
goals: dependency-management
best_when: dependency triage is important but never urgent
setup: some
composes_with: headless-mode, cloud-sessions
source: approaches/scheduled-agents.md

## session-context-management

id: session-context-management
kind: technique
goals: code-understanding
best_when: long exploration is saturating the context window
setup: none
composes_with: project-memory, checkpoints-rewind, subagent-delegation
source: approaches/session-context-management.md

## subagent-delegation

id: subagent-delegation
kind: technique
goals: ci-automation, code-review, migration, refactoring, security, tech-debt
best_when: multi-dimensional audit across the whole codebase
setup: involved
composes_with: worktree-isolation, plan-mode, built-in-review-skills
source: approaches/subagent-delegation.md

## visual-artifacts

id: visual-artifacts
kind: technique
goals: documentation
best_when: docs need diagrams, dashboards, or sharing beyond the repo
setup: none
composes_with: plan-mode, built-in-review-skills, fan-out-workflows
source: approaches/visual-artifacts.md

## worktree-isolation

id: worktree-isolation
kind: technique
goals: debugging, dependency-management, devops, migration, testing
best_when: bug might be your recent changes mixed with others
setup: some
composes_with: plan-mode, subagent-delegation, checkpoints-rewind
source: approaches/worktree-isolation.md
