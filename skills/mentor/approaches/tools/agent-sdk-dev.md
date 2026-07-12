---
kind: plugin
last_verified: 2026-07-12
composes_with:
  - custom-agents
  - headless-mode
install: /plugin install agent-sdk-dev@claude-plugins-official
facts: "Scaffolds Agent SDK projects and validates against best practices. Hands-on: sane strict-TS scaffold with streaming `query()`."
session_signal: "agent-sdk-dev is installed (its skills/commands are visible in the session) or its commands run in this conversation"
source: https://github.com/anthropics/claude-plugins-official
pitfalls:
  - "Pins dependencies to `latest` when the npm registry is unreachable, and its verifier agents only work after `npm install`."
---
