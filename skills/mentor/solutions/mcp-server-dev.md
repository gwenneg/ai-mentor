# mcp-server-dev
*Last verified: 2026-07-03*

kind: plugin
goals: building-mcp-integrations
best_when: writing an MCP server and wanting current SDK idioms instead of training-data guesses
composes_with: mcp-context, custom-plugins
install: /plugin install mcp-server-dev@claude-plugins-official
facts: Guided MCP server design and implementation. Hands-on: produced a syntax-clean stdio server with current SDK idioms (registerTool, zod validation, stdout hygiene) plus both config snippets. The SDK-idiom guidance is the value over base Claude.
session_signal: mcp-server-dev is installed (its skills/commands are visible in the session) or its commands run in this conversation
pitfalls:
- Guides design and code; it does not run or test the server — verify with a real client connection.
source: https://github.com/anthropics/claude-plugins-official
