# security
*Last verified: 2026-07-03*

**Hidden gem:** Hooks — a PreToolUse guard on auth configs and crypto files prevents the accidental security regressions that no scanner catches.

**Exemplar move:** Run /security-review on the current branch — special attention to auth middleware in src/middleware/auth.ts and raw database queries in src/services/; security audit next week.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Built-In Review Skills](../solutions/built-in-review-skills.md) | Quick security scan before a release or audit | Vulnerabilities follow well-known patterns — encoded checks applied exhaustively to every changed line beat manual review |
| 2 | [Subagent Delegation](../solutions/subagent-delegation.md) | Large codebase with multiple vulnerability classes to check | Security auditing is multi-dimensional — one concern per agent gives deeper analysis without attention dilution |
| 3 | [Deep Research](../solutions/deep-research.md) | New CVE announced for a dependency in your stack | Hardening without context is guesswork — affected versions and exploitation prerequisites let you patch what matters |
| 4 | [Hooks](../solutions/hooks-as-workflow.md) | Protect security-critical files from accidental modification | Most security regressions are accidental — a speed bump forces conscious acknowledgment before touching critical code |
| 5 | [security-guidance](../solutions/security-guidance.md) | Security review should happen on every edit, automatically | A hook that reviews as code is written catches issues at the cheapest moment — before they exist in a diff |
