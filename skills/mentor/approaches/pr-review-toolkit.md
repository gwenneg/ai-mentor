# pr-review-toolkit
*Last verified: 2026-07-03*

kind: plugin
composes_with: built-in-review-skills, subagent-delegation
install: /plugin install pr-review-toolkit@claude-plugins-official
facts: 6-agent review covering comments, tests, types, error handling, and simplification. Hands-on: found a planted off-by-one at the exact line with a verified repro, and flagged the deliberate test-coverage gap.
session_signal: pr-review-toolkit is installed (its skills/commands are visible in the session) or its commands run in this conversation
pitfalls:
- Token-hungry: ~2k always-on plus multiple subagents per review — reserve it for PRs that earn the spend.
- Overlaps the built-in /code-review; its additions are the comment, test-coverage, and type-design angles.
source: https://github.com/anthropics/claude-plugins-official
