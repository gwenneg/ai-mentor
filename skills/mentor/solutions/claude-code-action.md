# claude-code-action
*Last verified: 2026-07-06*

kind: integration
goals: ci-automation, code-review
best_when: a repo wants Claude reviews or automation on every PR without hand-rolling workflow YAML
composes_with: headless-mode, built-in-review-skills
facts: The maintained GitHub Action for running Claude Code in CI — packages checkout, CLI install, auth, and prompt execution; recommended over hand-rolled `curl | bash` workflow YAML for the standard PR-review case (hand-roll only for full pipeline control). Findings can post as inline PR comments alongside the diff.
session_signal: a workflow under .github/workflows/ uses anthropics/claude-code-action
pitfalls:
- It runs with the repository's ANTHROPIC_API_KEY secret — API billing, not a claude.ai subscription.
source: https://code.claude.com/docs/en/github-actions
