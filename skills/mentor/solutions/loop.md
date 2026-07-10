# /loop
*Last verified: 2026-07-06*

kind: builtin-command
goals: ci-automation, incident-response
best_when: work should recur on a time interval — polling a deploy, re-checking a queue — rather than converge on a condition
composes_with: goal, autonomous-loops, scheduled-agents
exemplar: /loop 5m check the deploy status and summarize changes
session_signal: user ran /loop in this conversation
pitfalls:
- `/loop` re-runs on an interval and stops when you stop it or Claude decides the work is done; `/goal` runs to a condition. Polling wants /loop, convergence wants /goal.
- It runs in the open session on this machine — for work that must survive a closed laptop, use /schedule instead.
source: https://code.claude.com/docs/en/commands
