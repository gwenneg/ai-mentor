# code-understanding
*Last verified: 2026-07-03*

**Hidden gem:** LSP Self-Correction — compiler-backed go-to-definition and find-references beat text search for tracing how components actually connect.

**Exemplar move:** Enter plan mode. Trace one payment end-to-end through services/payment-gateway/ — entry points, validation, fraud checks, processor integration, persistence, retry logic; produce an architecture summary with dependency diagram.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Plan Mode](../approaches/techniques/plan-mode.md) | Onboarding to a new team or project | Understanding is about building the right mental model — learn the shape of the system first, details later |
| 2 | [MCP Context](../approaches/techniques/mcp-context.md) | Architecture docs exist but are scattered or outdated | Code says what, docs say why; MCP brings both into one conversation and flags where docs drifted |
| 3 | [Deep Research](../approaches/techniques/deep-research.md) | Codebase uses unfamiliar frameworks or patterns | Frameworks encode decisions into conventions — learning the grammar lets you read the codebase's sentences fluently |
| 4 | [LSP Self-Correction](../approaches/techniques/lsp-self-correction.md) | Tracing how components connect across a large codebase | LSP gives compiler-precise answers to who-calls-what, where text search misses indirect calls and aliased imports |
| 5 | [Project Memory & Context Docs](../approaches/techniques/project-memory.md) | Want the map you built to persist across sessions | Exploration output is knowledge — storing it where every session reads converts one-off investigation into permanent capability |
| 6 | [Session & Context Management](../approaches/techniques/session-context-management.md) | Long exploration is saturating the context window | Exploration quality degrades silently as context fills; curating the window keeps reasoning over conclusions, not noise |
| 7 | [context7](../approaches/tools/context7.md) | Working against a library version your training data predates | Version-pinned docs on demand beat guessing from stale memory — evaluation returned real Express v5 docs |
