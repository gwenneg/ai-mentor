---
kind: plugin
last_verified: 2026-07-03
composes_with: [mcp-context, deep-research]
install: /plugin install context7@claude-plugins-official
facts: "Pulls version-pinned documentation for any library on demand (Upstash-maintained MCP server). Hands-on: returned real Express v5 docs through the MCP tool, no account needed."
session_signal: "context7 is installed (its MCP tools are visible in the session) or its doc lookups run in this conversation"
source: https://github.com/anthropics/claude-plugins-official
pitfalls:
  - "Headless callers must allowlist the whole MCP server — an unallowlisted `claude -p` run silently gets no docs."
  - "Third-party (Upstash) — its release cadence is not Claude Code's; re-verify on the evaluation pass, not just the manifest sync."
---
