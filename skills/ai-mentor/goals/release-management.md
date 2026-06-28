# Release Management
*Last reviewed: 2026-06-28*

## When You're Here

You're preparing to ship. Maybe you need to cut a release, write changelog entries that actually help users, bump versions across multiple packages, validate that nothing critical slipped through, or coordinate a deployment that involves database migrations and feature flags. The stakes are higher than a regular PR — a bad release affects every user, not just one feature.

This is distinct from CI/CD automation (which is about pipeline mechanics) and code review (which catches issues during development). Release management is about the decision-making and preparation process itself: what's in the release, is it ready, and how do we get it out safely.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Need categorized release notes from commit history | Custom skills | One command generates a changelog from commits since the last tag |
| Validating release readiness in CI before cutting a tag | Headless mode | Non-interactive checks for unreleased migrations, version mismatches, TODOs |
| Complex release with multiple services, migrations, and rollback steps | Plan mode | Maps deployment order and rollback strategy before you start |
| Final quality check on the release diff | Built-in review skills | Catches last-minute bugs and security issues in the release branch |
| Want automatic validation on every edit during release prep | Hooks | PostToolUse checks run version consistency and changelog format after each change |

## Approaches (Ranked)

### 1. Custom Skills — One-command release workflows
**Level:** Intermediate | **Tools:** Claude Code

Your release process has steps that repeat every time: read commits since the last tag, categorize them (features, fixes, breaking changes), generate a changelog entry, bump the version, and create the release PR. A `/release-notes` skill reads your commit history and produces a categorized changelog. A `/prepare-release` skill bumps versions across `package.json`, `CHANGELOG.md`, and any other version-bearing files, then opens the PR. Encode your release process once, run it every time.

**Try it now:**
> Create a custom skill at `.claude/skills/release-notes.md`. When invoked with `/release-notes`, it should: (1) find the most recent git tag, (2) read all commits since that tag, (3) categorize them into "Features," "Bug Fixes," "Breaking Changes," and "Internal" based on conventional commit prefixes, (4) generate a markdown changelog entry with the next version number and today's date, and (5) prepend it to `CHANGELOG.md`. If any commits have a `BREAKING CHANGE` footer, add a migration note at the top.

**Why this works:** Release preparation is procedural — the same steps in the same order, every time. Custom skills eliminate the risk of forgetting a step or doing them out of order, and they ensure every release gets the same quality of changelog regardless of who cuts it.

**Pros:**
- Eliminates manual changelog writing and version-bump errors
- Same workflow whether triggered interactively or in a CI pipeline
- Versioned in the repo, so the release process itself is reviewable

**Cons:**
- Skill definitions need updates when your release process changes
- Unusual releases (hotfixes, pre-releases) may need separate skills

**Deeper:** See `approaches/custom-skills.md`

---

### 2. Headless Mode — Automated release validation in CI
**Level:** Intermediate | **Tools:** Claude Code

Run `claude -p` as a CI step before tagging a release. The prompt checks for release blockers that static analysis misses: unreleased database migration scripts sitting in the repo, version strings that weren't bumped, TODO/FIXME comments in release-critical paths, changelog entries that reference the wrong version, or dependencies pinned to pre-release versions. If any check fails, the CI step exits non-zero and blocks the release.

**Try it now:**
> Write a CI step that runs `claude -p` with `--output-format json` before every release tag. The prompt should: (1) check that the version in `package.json` matches the tag being created, (2) verify `CHANGELOG.md` has an entry for this version, (3) scan `src/` for any `TODO` or `FIXME` comments containing the word "release" or "before shipping," (4) check `db/migrations/` for any migration files not referenced in the latest release notes, and (5) output a JSON object with `{"ready": boolean, "blockers": [...]}`. Block the tag if `ready` is false.

**Why this works:** Humans forget pre-release checks under deadline pressure. A headless validation step catches the blockers that slip through when you're rushing to ship on a Friday afternoon — version mismatches, forgotten migrations, debug code left in production paths.

**Pros:**
- Catches release blockers that linters and tests don't cover
- Runs automatically, no human discipline required
- Structured JSON output integrates with release tooling and dashboards

**Cons:**
- Prompt must anticipate all relevant checks — it can't ask follow-up questions
- Adds CI time to the release pipeline

**Deeper:** See `approaches/headless-mode.md`

---

### 3. Plan Mode — Coordinate complex releases
**Level:** Beginner | **Tools:** Any

When a release involves more than "merge and tag" — multiple services that need to deploy in order, database migrations that must run before the new code goes live, feature flags that need toggling at specific stages, or rollback procedures that depend on which step failed — plan mode maps the entire sequence before you start executing. You get a deployment order, rollback strategy, and verification steps for each stage.

**Try it now:**
> Enter plan mode. We're releasing v2.0 which includes: a database migration adding a `preferences` column to the `users` table, a new API endpoint `/api/v2/users/preferences`, a React frontend component that calls this endpoint, and a feature flag `enable_user_preferences` in LaunchDarkly. Plan the deployment sequence: what deploys first, what do we verify at each stage, what's the rollback procedure if the migration succeeds but the API deploy fails, and what order do we enable the feature flag relative to the other steps?

**Why this works:** Complex releases fail when steps execute out of order or rollback procedures are improvised under pressure. Planning the full sequence — including what to do when things go wrong — turns a stressful coordination exercise into a checklist you execute with confidence.

**Pros:**
- Surfaces deployment ordering issues before they become production incidents
- Rollback procedures are defined before you need them, not during the outage
- Works with any AI tool, no special setup required

**Cons:**
- Planning takes time upfront — but prevents costly mid-release surprises

**Deeper:** See `approaches/plan-mode.md`

---

### 4. Built-In Review Skills — Final quality gate before release
**Level:** Beginner | **Tools:** Claude Code

Before cutting the tag, run `/code-review --effort high` on the full release diff (everything since the last release). This catches issues that accumulated across multiple PRs — a function renamed in one PR but still referenced by its old name in another, error handling that's inconsistent across features merged by different developers, or a security fix in one module that wasn't applied to a similar pattern elsewhere. For anything touching authentication, payments, or data access, follow up with `/security-review`.

**Try it now:**
> Run `/code-review --effort high` on the diff between the `v1.9.0` tag and the current `release/v2.0` branch. Pay special attention to `src/services/auth/` and `src/services/billing/` — these changed across multiple PRs and I want to make sure nothing fell through the cracks. Then run `/security-review` specifically on the auth and billing directories.

**Why this works:** Individual PRs get reviewed, but the interactions between PRs merged across a release cycle often don't. A full-diff review at release time catches cross-PR inconsistencies that no individual review would surface.

**Pros:**
- Catches cross-PR inconsistencies that individual reviews miss
- Zero setup — built-in skills work immediately
- Security review provides an extra gate for sensitive code paths

**Cons:**
- Large release diffs produce many findings that need triage
- Cannot catch runtime integration issues, only static code problems

**Deeper:** See `approaches/built-in-review-skills.md`

---

### 5. Hooks — Auto-validate on every change during release prep
**Level:** Intermediate | **Tools:** Claude Code

During release preparation, every edit carries extra risk. A PostToolUse hook runs validation after each file change: verify that version strings stay consistent across `package.json`, `pyproject.toml`, and `src/version.ts`; check that `CHANGELOG.md` still follows your format conventions; run the pre-release test suite if a release-critical file changed. This turns release prep into a guardrailed process where mistakes are caught at edit time, not at tag time.

**Try it now:**
> Set up a PostToolUse hook for release prep. After any edit to files matching `package.json`, `CHANGELOG.md`, or `src/version.ts`, the hook should: (1) verify the version string is identical across all three files, (2) check that the `CHANGELOG.md` entry for the current version has a date and at least one item in each non-empty category, and (3) run `npm run test:release` if it exists. Report any inconsistencies immediately so I can fix them before moving on.

**Why this works:** Release preparation is a sequence of edits where consistency matters more than speed. Hooks enforce that consistency automatically — you can't accidentally bump the version in one file and forget another, because the hook catches it the moment it happens.

**Pros:**
- Catches version drift and format errors at edit time, not at release time
- Zero manual discipline required — validation happens automatically
- Composable with other hooks for layered release-prep protection

**Cons:**
- Hooks that run full test suites slow down the edit-verify cycle
- Must be configured specifically for release prep, then potentially disabled after

**Deeper:** See `approaches/hooks-as-workflow.md`
