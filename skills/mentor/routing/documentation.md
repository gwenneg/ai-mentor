# documentation
*Last verified: 2026-07-03*

**Hidden gem:** Custom Skills — a `/gen-api-doc` command makes regenerating docs cheaper than letting them drift.

**Exemplar move:** Read src/api/routes/payments.ts, src/api/middleware/auth.ts, docs/openapi.yaml, and the Notion design doc linked in docs/DESIGN_DECISIONS.md; generate payments API reference with auth, schemas, error codes, rate limits.

**Plugins:** `claude-md-management` ✅ CLAUDE.md audits · `project-artifact` ✅ living status pages · `mintlify` ☑️ docs sites · `notion`/`atlassian` ☑️ knowledge bases.

**Built-ins:** `/init` — bootstrap CLAUDE.md so sessions stop re-learning the repo. Facts and pitfalls per command: `registry/builtin-commands.md`.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [MCP Context](../approaches/mcp-context.md) | Intermediate | Existing docs, specs, and decisions are scattered | Documentation quality is proportional to context quality — MCP bridges where knowledge lives and where it needs to go |
| 2 | [Deep Research](../approaches/deep-research.md) | Beginner | Writing docs for an unfamiliar domain or standard | Documentation is a communication design problem; proven patterns give a structural blueprint so effort goes to content |
| 3 | [Plan Mode](../approaches/plan-mode.md) | Beginner | Large doc set with multiple sections and audiences | Documentation is information architecture — an outline gives every fact one home and writers clear coverage without overlap |
| 4 | [Custom Skills](../approaches/custom-skills.md) | Beginner | You regenerate the same type of documentation repeatedly | Same-structure docs are a solved problem once encoded — consistent, complete output from source code every time |
| 5 | [Visual Artifacts](../approaches/visual-artifacts.md) | Beginner | Docs need diagrams, dashboards, or sharing beyond the repo | Diagrams carry relationships linear Markdown flattens, and a stable URL turns docs into a link people actually open |
