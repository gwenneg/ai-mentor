# documentation
*Last verified: 2026-07-12*

**Hidden gem:** Custom Skills — a `/gen-api-doc` command makes regenerating docs cheaper than letting them drift.

**Exemplar move:** Read src/api/routes/payments.ts, src/api/middleware/auth.ts, docs/openapi.yaml, and the Notion design doc linked in docs/DESIGN_DECISIONS.md; generate payments API reference with auth, schemas, error codes, rate limits.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [MCP Context](../approaches/techniques/mcp-context.md) | Existing docs, specs, and decisions are scattered | Documentation quality is proportional to context quality — MCP bridges where knowledge lives and where it needs to go |
| 2 | [Deep Research](../approaches/techniques/deep-research.md) | Writing docs for an unfamiliar domain or standard | Documentation is a communication design problem; proven patterns give a structural blueprint so effort goes to content |
| 3 | [Plan Mode](../approaches/techniques/plan-mode.md) | Large doc set with multiple sections and audiences | Documentation is information architecture — an outline gives every fact one home and writers clear coverage without overlap |
| 4 | [Custom Skills](../approaches/techniques/custom-skills.md) | You regenerate the same type of documentation repeatedly | Same-structure docs are a solved problem once encoded — consistent, complete output from source code every time |
| 5 | [Visual Artifacts](../approaches/techniques/visual-artifacts.md) | Docs need diagrams, dashboards, or sharing beyond the repo | Diagrams carry relationships linear Markdown flattens, and a stable URL turns docs into a link people actually open |
| 6 | [claude-md-management](../approaches/tools/claude-md-management.md) | A CLAUDE.md exists but nobody knows if it is still good | A scored audit against the real codebase turns 'our CLAUDE.md is probably stale' into a fix list |
| 7 | [project-artifact](../approaches/tools/project-artifact.md) | Project status should live on a shareable page | A living status page with honest unverified-state markings replaces the weekly re-explanation meeting |
