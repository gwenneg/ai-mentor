# /schedule
*Last verified: 2026-07-06*

kind: builtin-command
goals: ci-automation, dependency-management, release-management
best_when: recurring maintenance should happen on Anthropic-run cloud infrastructure with the laptop closed
composes_with: scheduled-agents, headless-mode
exemplar: /schedule daily Renovate-PR triage at 6am
session_signal: user ran /schedule in this conversation, or routines exist for their account
pitfalls:
- Research preview: needs a claude.ai subscription login; minimum interval one hour; runs clone from GitHub, never from local disk.
- A green run status means the session completed, not that the task succeeded — read the transcript.
source: https://code.claude.com/docs/en/routines
