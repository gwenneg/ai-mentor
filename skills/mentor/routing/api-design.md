# api-design
*Last verified: 2026-07-03*

**Hidden gem:** Built-In Review Skills — pointing `/code-review` at only the API surface catches breaking changes and naming drift that no test suite will.

**Exemplar move:** Enter plan mode. Design a REST API for multi-tenant project management (projects, tasks, comments, members): URL structure, schemas, error codes, versioning, tenant scoping, pagination for 10K+ tasks. No code.

**Plugins:** `apollo-skills` ☑️ GraphQL (client, server, federation) · `postman` ☑️ collections, tests, and mocks.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | Beginner | Starting a new API or major version from scratch | An API is a promise to consumers; planning the contract first finds design mistakes while fixing them is free |
| 2 | [Deep Research](../approaches/deep-research.md) | Beginner | Unsure about conventions like pagination, versioning, error format | Studying proven APIs gives rationale for your choices and avoids pitfalls others already documented |
| 3 | [MCP Context](../approaches/mcp-context.md) | Intermediate | Requirements scattered across specs, tickets, and existing APIs | Pulling the existing API surface alongside requirements keeps new endpoints consistent — same naming, errors, pagination |
| 4 | [Built-in Review Skills](../approaches/built-in-review-skills.md) | Beginner | Extending an existing API and need to stay consistent | Misspelled fields and inconsistent errors break no test but frustrate every consumer; review catches surface issues type checkers miss |
