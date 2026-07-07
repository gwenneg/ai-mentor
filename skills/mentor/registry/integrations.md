# Capability Registry — integrations & docs
*Last verified: 2026-07-06*

The `kind: integration` and `kind: doc` slice of the capability registry (NEXT_VERSION.md D1/D6): recommendable things that are neither marketplace plugins nor techniques — a GitHub Action, an SDK, an official doc page that IS the right answer. Records carry enough dated facts inline that the mentor never needs to fetch a URL mid-invocation (a runtime web fetch is a permission prompt and a latency cliff); the URL is for the human. Facts below are drawn from the approach files' verified content and carry that verification date. These records are teachable by default — same ignorance-map citizenship as techniques and built-in commands; profile rows use the record `id`.

## claude-code-action

id: claude-code-action
kind: integration
goals: ci-automation, code-review
best_when: a repo wants Claude reviews or automation on every PR without hand-rolling workflow YAML
setup: some
composes_with: headless-mode, built-in-review-skills
facts: The maintained GitHub Action for running Claude Code in CI — packages checkout, CLI install, auth, and prompt execution; recommended over hand-rolled `curl | bash` workflow YAML for the standard PR-review case (hand-roll only for full pipeline control). Findings can post as inline PR comments alongside the diff.
session_signal: a workflow under .github/workflows/ uses anthropics/claude-code-action
pitfalls:
- It runs with the repository's ANTHROPIC_API_KEY secret — API billing, not a claude.ai subscription.
source: https://code.claude.com/docs/en/github-actions

## agent-sdk

id: agent-sdk
kind: doc
goals: building-agents, llm-features
best_when: someone is building an autonomous AI product or teammate, not just using Claude Code interactively
setup: involved
composes_with: headless-mode, custom-agents
facts: The Claude Agent SDK is the supported path for building custom agents as products — programmatic sessions, custom tool definitions, and agent loops outside the terminal. It is a different altitude than custom agent definitions (`.claude/agents/*.md`, which configure subagents inside Claude Code): the SDK builds standalone agent applications. Headless mode (`claude -p`) covers the simpler "script Claude in a pipeline" case without the SDK.
session_signal: the repo imports @anthropic-ai/claude-agent-sdk or discusses building an agent product
pitfalls:
- Reaching for the SDK when headless mode or a custom skill would do — the SDK is for products, not automation glue.
source: https://code.claude.com/docs/en/agent-sdk/overview

## gitlab-ci-integration

id: gitlab-ci-integration
kind: integration
goals: ci-automation
best_when: the team is on GitLab and wants Claude in the pipeline without GitHub-specific tooling
setup: some
composes_with: headless-mode
facts: GitLab CI/CD integration is beta and maintained by GitLab: one `.gitlab-ci.yml` job that installs the CLI and runs `claude -p` with a masked `ANTHROPIC_API_KEY` variable. There is no GitLab equivalent of claude-code-action — the job is the integration.
session_signal: the repo has a .gitlab-ci.yml
pitfalls:
- Beta status and GitLab ownership mean its cadence follows GitLab releases, not Claude Code's.
source: https://code.claude.com/docs/en/gitlab-ci-cd
