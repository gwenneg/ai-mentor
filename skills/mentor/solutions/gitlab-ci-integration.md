# gitlab-ci-integration
*Last verified: 2026-07-06*

kind: integration
goals: ci-automation
best_when: the team is on GitLab and wants Claude in the pipeline without GitHub-specific tooling
composes_with: headless-mode
facts: GitLab CI/CD integration is beta and maintained by GitLab: one `.gitlab-ci.yml` job that installs the CLI and runs `claude -p` with a masked `ANTHROPIC_API_KEY` variable. There is no GitLab equivalent of claude-code-action — the job is the integration.
session_signal: the repo has a .gitlab-ci.yml
pitfalls:
- Beta status and GitLab ownership mean its cadence follows GitLab releases, not Claude Code's.
source: https://code.claude.com/docs/en/gitlab-ci-cd
