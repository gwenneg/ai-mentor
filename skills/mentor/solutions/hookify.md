# hookify
*Last verified: 2026-07-03*

kind: plugin
goals: ci-automation, building-skills-plugins
best_when: a conversation pattern ("run lint after every edit") should become a real hook without hand-writing settings JSON
composes_with: hooks-as-workflow
install: /plugin install hookify@claude-plugins-official
facts: Creates hooks from conversation patterns or explicit rules. Hands-on: generated a working PostToolUse hook and the firing was verified.
session_signal: hookify is installed (its skills/commands are visible in the session) or its commands run in this conversation
pitfalls:
- Headless caveat: it cannot write settings files non-interactively — hook creation needs an interactive session.
source: https://github.com/anthropics/claude-plugins-official
