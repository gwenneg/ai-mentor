# Processed Changelogs
*Updated: 2026-07-20*

Tracks which official what's-new digests (https://code.claude.com/docs/en/whats-new/) have been incorporated into this catalog. The weekly digest slug (e.g. `2026-w26`) is the stable unit of processing: maintenance Step 4 fetches the digest index, processes any week not listed below (oldest first), and appends one row per digest with what was done. A digest is never processed twice, and a gap in this table is by definition unprocessed work.

The catalog's initial completeness was established by a one-time bootstrap on 2026-07-02: every page of the official docs index was classified and every coverage claim verified against the live docs. From that point forward, this ledger is the maintenance trail.

| Week | Processed | Outcome |
|------|-----------|---------|
| [2026-w13](https://code.claude.com/docs/en/whats-new/2026-w13.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs |
| [2026-w14](https://code.claude.com/docs/en/whats-new/2026-w14.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs |
| [2026-w15](https://code.claude.com/docs/en/whats-new/2026-w15.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs |
| [2026-w16](https://code.claude.com/docs/en/whats-new/2026-w16.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs |
| [2026-w17](https://code.claude.com/docs/en/whats-new/2026-w17.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs |
| [2026-w18](https://code.claude.com/docs/en/whats-new/2026-w18.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs; resume-by-PR-URL / `--from-pr` backfilled into session-context-management on 2026-07-13 |
| [2026-w19](https://code.claude.com/docs/en/whats-new/2026-w19.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs |
| [2026-w20](https://code.claude.com/docs/en/whats-new/2026-w20.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs |
| [2026-w21](https://code.claude.com/docs/en/whats-new/2026-w21.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs; `/usage` per-category attribution backfilled into model-effort-selection on 2026-07-13 |
| [2026-w22](https://code.claude.com/docs/en/whats-new/2026-w22.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs |
| [2026-w23](https://code.claude.com/docs/en/whats-new/2026-w23.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs |
| [2026-w24](https://code.claude.com/docs/en/whats-new/2026-w24.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs; `/cd` and `--safe-mode` backfilled into session-context-management on 2026-07-13 |
| [2026-w25](https://code.claude.com/docs/en/whats-new/2026-w25.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs |
| [2026-w26](https://code.claude.com/docs/en/whats-new/2026-w26.md) | 2026-07-02 | Initial bootstrap — catalog built and verified against live docs; shell-mode-responds-to-output backfilled into session-context-management on 2026-07-13 |
| [2026-w27](https://code.claude.com/docs/en/whats-new/2026-w27.md) | 2026-07-12 | Updated approaches/techniques/subagent-delegation.md (background-by-default, `background` frontmatter) and model-effort-selection.md (Sonnet 5 raised the everyday tier's ceiling); Chrome GA, Artifacts GA, dataviz, draft-PR automation already covered |
| [2026-w28](https://code.claude.com/docs/en/whats-new/2026-w28.md) | 2026-07-12 | Updated approaches/techniques/browser-integration.md (Desktop in-app browser) and session-context-management.md (/doctor checkup folded in); remaining items UX/infra — no action |
| [2026-w29](https://code.claude.com/docs/en/whats-new/2026-w29.md) | 2026-07-20 | Updated visual-artifacts.md (MCP-connector-backed dashboards, corrected public-sharing model, editor roles), subagent-delegation.md and background-agents.md (`/fork` now spawns a background session, in-session behavior renamed `/subtask`, session-wide WebSearch/subagent caps), safe-autonomy.md (auto mode needs no opt-in on Bedrock/GCP/Foundry, "always allow" persists across worktrees), headless-mode.md (`--forward-subagent-text`), and mcp-context.md (MCP calls auto-background after 2 minutes); screen reader mode, corporate launcher, vim remaps, and other admin/UX items — no action |
