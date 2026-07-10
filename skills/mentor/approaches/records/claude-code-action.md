---
kind: integration
last_verified: 2026-07-06
composes_with: [headless-mode, built-in-review-skills]
facts: "The maintained GitHub Action for running Claude Code in CI — packages checkout, CLI install, auth, and prompt execution; recommended over hand-rolled `curl | bash` workflow YAML for the standard PR-review case (hand-roll only for full pipeline control). Findings can post as inline PR comments alongside the diff."
session_signal: "a workflow under .github/workflows/ uses anthropics/claude-code-action"
source: https://code.claude.com/docs/en/github-actions
pitfalls:
  - "It runs with the repository's ANTHROPIC_API_KEY secret — API billing, not a claude.ai subscription."
---
