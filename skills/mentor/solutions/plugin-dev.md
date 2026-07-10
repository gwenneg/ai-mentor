# plugin-dev
*Last verified: 2026-07-03*

kind: plugin
composes_with: custom-plugins, custom-skills
install: /plugin install plugin-dev@claude-plugins-official
facts: 8-phase guided workflow for building plugins, with validator and reviewer agents. Hands-on: scaffolded a plugin that passed `claude plugin validate` and self-reviewed honestly. The entry point is `create-plugin`.
session_signal: plugin-dev is installed (its skills/commands are visible in the session) or its commands run in this conversation
pitfalls:
- Heaviest always-on context of the evaluated set (~2.3k tokens) — install it for plugin work, not permanently.
source: https://github.com/anthropics/claude-plugins-official
