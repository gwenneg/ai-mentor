# release-management
*Last verified: 2026-07-03*

**Hidden gem:** Headless Mode — a headless pre-tag validation step catches the version mismatches and forgotten migrations humans miss on Friday afternoons.

**Exemplar move:** Create .claude/skills/release-notes/SKILL.md: /release-notes finds the latest git tag, categorizes commits by conventional prefixes, prepends a versioned changelog entry to CHANGELOG.md, adds migration notes for BREAKING CHANGE footers.

**Plugins:** `commit-commands` ✅ commit/PR workflow commands · `confidence` ☑️ feature flags and rollouts.

**Built-ins:** `/schedule` — scheduled release chores on cloud infrastructure. Facts and pitfalls per command: its `solutions/<id>.md` record.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Custom Skills](../solutions/custom-skills.md) | Need categorized release notes from commit history | Release prep is procedural — same steps, same order — so a skill removes forgotten or misordered steps |
| 2 | [Headless Mode](../solutions/headless-mode.md) | Validating release readiness in CI before cutting a tag | Humans forget pre-release checks under deadline pressure; headless validation catches blockers when you're rushing to ship |
| 3 | [Plan Mode](../solutions/plan-mode.md) | Complex release with multiple services, migrations, and rollback steps | Complex releases fail on ordering and improvised rollbacks — planning the full sequence turns coordination into a checklist |
