# security
*Last verified: 2026-07-03*

**Hidden gem:** Hooks — a PreToolUse guard on auth configs and crypto files prevents the accidental security regressions that no scanner catches.

**Exemplar move:** Run /security-review on the current branch — special attention to auth middleware in src/middleware/auth.ts and raw database queries in src/services/; security audit next week.

**Plugins:** `security-guidance` ✅ per-edit security hooks · `semgrep` ☑️ scanning · `42crunch-api-security-testing` ☑️ API security · `auth0` ☑️ authn/authz · `vanta` ☑️ compliance — 9 more in the catalog.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Built-In Review Skills](../approaches/built-in-review-skills.md) | Beginner | Quick security scan before a release or audit | Vulnerabilities follow well-known patterns — encoded checks applied exhaustively to every changed line beat manual review |
| 2 | [Subagent Delegation](../approaches/subagent-delegation.md) | Advanced | Large codebase with multiple vulnerability classes to check | Security auditing is multi-dimensional — one concern per agent gives deeper analysis without attention dilution |
| 3 | [Deep Research](../approaches/deep-research.md) | Beginner | New CVE announced for a dependency in your stack | Hardening without context is guesswork — affected versions and exploitation prerequisites let you patch what matters |
| 4 | [Hooks](../approaches/hooks-as-workflow.md) | Intermediate | Protect security-critical files from accidental modification | Most security regressions are accidental — a speed bump forces conscious acknowledgment before touching critical code |
