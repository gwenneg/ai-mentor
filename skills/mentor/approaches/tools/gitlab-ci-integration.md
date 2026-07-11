---
kind: integration
last_verified: 2026-07-06
composes_with:
  - headless-mode
facts: "GitLab CI/CD integration is beta and maintained by GitLab: one `.gitlab-ci.yml` job that installs the CLI and runs `claude -p` with a masked `ANTHROPIC_API_KEY` variable. There is no GitLab equivalent of claude-code-action — the job is the integration."
session_signal: "the repo has a .gitlab-ci.yml"
source: https://code.claude.com/docs/en/gitlab-ci-cd
pitfalls:
  - "Beta status and GitLab ownership mean its cadence follows GitLab releases, not Claude Code's."
---
