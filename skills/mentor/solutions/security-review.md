# /security-review
*Last verified: 2026-07-06*

kind: builtin-command
goals: security, code-review
best_when: pending changes touch auth, input handling, or anything a CVE could grow from
composes_with: code-review, built-in-review-skills
exemplar: /security-review
session_signal: user ran /security-review in this conversation
pitfalls:
- It reviews the pending changes on the current branch — run it before merging, not after.
- A clean pass is one lens, not a security audit; critical paths still need a human pass.
source: https://code.claude.com/docs/en/commands
