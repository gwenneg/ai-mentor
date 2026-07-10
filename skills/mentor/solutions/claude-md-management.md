# claude-md-management
*Last verified: 2026-07-03*

kind: plugin
goals: documentation
best_when: a CLAUDE.md exists but nobody knows if it is still good — audit it against the real codebase
composes_with: project-memory
install: /plugin install claude-md-management@claude-plugins-official
facts: Audits and maintains CLAUDE.md files. Hands-on: scored audit with a rubric and real gaps found, cross-checked against the codebase.
session_signal: claude-md-management is installed (its skills/commands are visible in the session) or its commands run in this conversation
pitfalls:
- The skill is invoked as `claude-md-improver`, not the plugin name.
source: https://github.com/anthropics/claude-plugins-official
