# dependency-management
*Last verified: 2026-07-03*

**Hidden gem:** Worktree Isolation — trying the upgrade in a throwaway copy gives you a real damage report before you commit to anything.

**Exemplar move:** /deep-research Compare zod vs joi for our Node.js API (src/validators/): TypeScript integration, bundle size, performance, maintenance activity, breaking-change history, CVEs — must stay maintained 3+ years.

**Plugins:** `sonatype-guide` ☑️ vulnerability and version analysis · `ai-plugins` ☑️ (Endor Labs) supply-chain scanning.

**Built-ins:** `/deep-research` — due diligence on a candidate dependency; `/schedule` — nightly dependency-PR triage. Facts and pitfalls per command: its `solutions/<id>.md` record.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Deep Research](../solutions/deep-research.md) | Evaluating a library you haven't used before | Adoption is a long-term bet; research automates the due diligence — maintenance health, CVEs, licenses — most developers skip |
| 2 | [Plan Mode](../solutions/plan-mode.md) | Understanding what depends on a package before removing it | Dependency changes are graph operations — mapping the full graph first catches transitive breakage before it happens |
| 3 | [MCP Context](../solutions/mcp-context.md) | Your org has approved dependency lists or security policies | The fastest decision is one someone already made — internal sources prevent duplicating or contradicting prior evaluations |
| 4 | [Worktree Isolation](../solutions/worktree-isolation.md) | Want to test a major upgrade without risking your branch | A disposable environment gives a realistic damage assessment instead of guessing from the changelog |
| 5 | [Scheduled & Recurring Agents](../solutions/scheduled-agents.md) | Dependency triage is important but never urgent | Recurring maintenance survives only when it stops depending on human initiative — a schedule converts "someone should" into "it happened" |
