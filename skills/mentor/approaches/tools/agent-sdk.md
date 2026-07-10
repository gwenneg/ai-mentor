---
kind: doc
last_verified: 2026-07-06
composes_with: [headless-mode, custom-agents]
facts: "The Claude Agent SDK is the supported path for building custom agents as products — programmatic sessions, custom tool definitions, and agent loops outside the terminal. It is a different altitude than custom agent definitions (`.claude/agents/*.md`, which configure subagents inside Claude Code): the SDK builds standalone agent applications. Headless mode (`claude -p`) covers the simpler \"script Claude in a pipeline\" case without the SDK."
session_signal: "the repo imports @anthropic-ai/claude-agent-sdk or discusses building an agent product"
source: https://code.claude.com/docs/en/agent-sdk/overview
pitfalls:
  - "Reaching for the SDK when headless mode or a custom skill would do — the SDK is for products, not automation glue."
---
