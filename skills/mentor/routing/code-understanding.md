# code-understanding
*Last verified: 2026-07-03*

**Hidden gem:** LSP Self-Correction — compiler-backed go-to-definition and find-references beat text search for tracing how components actually connect.

**Exemplar move:** Enter plan mode. Trace one payment end-to-end through services/payment-gateway/ — entry points, validation, fraud checks, processor integration, persistence, retry logic; produce an architecture summary with dependency diagram.

**Plugins:** `context7` ✅ version-pinned library docs · `sourcegraph` ☑️ cross-repo search · `serena` ☑️ semantic analysis · `lumen` ☑️ local semantic search.

**Built-ins:** `/init` — capture what you learn into a starter CLAUDE.md. Facts and pitfalls per command: `registry/builtin-commands.md`.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | Beginner | Onboarding to a new team or project | Understanding is about building the right mental model — learn the shape of the system first, details later |
| 2 | [MCP Context](../approaches/mcp-context.md) | Intermediate | Architecture docs exist but are scattered or outdated | Code says what, docs say why; MCP brings both into one conversation and flags where docs drifted |
| 3 | [Deep Research](../approaches/deep-research.md) | Beginner | Codebase uses unfamiliar frameworks or patterns | Frameworks encode decisions into conventions — learning the grammar lets you read the codebase's sentences fluently |
| 4 | [LSP Self-Correction](../approaches/lsp-self-correction.md) | Intermediate | Tracing how components connect across a large codebase | LSP gives compiler-precise answers to who-calls-what, where text search misses indirect calls and aliased imports |
| 5 | [Project Memory & Context Docs](../approaches/project-memory.md) | Beginner | Want the map you built to persist across sessions | Exploration output is knowledge — storing it where every session reads converts one-off investigation into permanent capability |
| 6 | [Session & Context Management](../approaches/session-context-management.md) | Beginner | Long exploration is saturating the context window | Exploration quality degrades silently as context fills; curating the window keeps reasoning over conclusions, not noise |
