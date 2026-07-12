---
name: ai-mentor-update
description: >-
  Maintenance skill for the ai-mentor plugin. Audits structure, verifies
  content accuracy against current tool docs, and detects recent tool
  changes. Each step is optional — the user picks what to run, or passes
  --auto for a non-interactive run (CI).
argument-hint: "[--auto 4,5 --files 5]"
disable-model-invocation: true
---

# AI Mentor Update

You are a content maintenance tool for the ai-mentor plugin. Your job is to keep the goal and approach files structurally consistent, factually accurate, and up to date.

All paths below are relative to the repo root.

---

## Modes

**Interactive (default):** follow the steps as written, asking the user at each decision point.

**Non-interactive (`--auto`):** when `$ARGUMENTS` contains `--auto`, never ask a question and never wait for input — there is no user present (this mode runs headless in CI). Parse from the arguments:

- `--auto <steps>` — comma-separated step numbers to run, using the Step headings below (2 = structural audit, 3 = content verification, 4 = process new changelogs, 5 = plugin catalog sync), e.g. `--auto 4,5`
- `--files N` — Step 3 scope: process the N oldest-verified files (default: 5). Promoted records flagged by `go -C tools/catalog-drift run .` as upstream-changed (its `!` lines) are ALWAYS in scope, in addition to the N oldest — they are why step 3 was triggered

The routine weekly run is `--auto 4,5`: process new changelogs and sync the plugin catalog. Steps 2 and 3 are deeper audits for occasional use — Step 2 largely duplicates the CI gate, and Step 3 re-verifies old content that changelog processing doesn't touch.

Auto-mode overrides, in addition to skipping every question below:

- **Step 2**: apply only unambiguous structural fixes (broken separator, wrong field order); report anything requiring judgment instead of fixing it. Note: CI also runs the structural audit (`go -C tools/structural-audit run .`) as a deterministic gate — prefer reporting over creative fixing.
- **Step 3**: process the N oldest-verified files — plus every promoted record `catalog-drift` flags as upstream-changed — with no per-file pause. Apply only changes that meet the "Recommended changes" bar (official-tier source + direct quote); list everything else under "needs manual verification" in the report without applying it.
- **Step 4**: process every digest not yet in the ledger; same evidence bar as Step 3; always append the ledger row for each processed digest.
- **Step 5**: apply additions and removals directly — the GitHub API response is authoritative.
- **Never** run `git commit`, `git push`, or create branches — the calling workflow owns git.
- **Final output**: end with a single markdown report (the caller uses it as a PR body) with two sections: **Changes applied** (file, change, source URL, supporting quote) and **Not applied** (finding + why it needs a human). If nothing changed, say so explicitly.

---

## Step 1 — Setup

Run `date +%Y-%m-%d` via bash to establish today's date.

Then ask the user which steps to run. The menu numbers are the Step headings below — the same numbers `--auto` uses, so "4" means the same thing interactively and in CI:

> Which maintenance steps do you want to run?
>
> 2. **Structural audit** — check all files against templates, cross-references, staleness
> 3. **Content verification** — web-search claims in files against current tool docs
> 4. **Process new changelogs** — incorporate what's-new digests not yet listed in `processed-changelogs.md`
> 5. **Plugin catalog sync** — check `claude-plugins-official` for new or removed plugins not yet reflected in `marketplace.md` or the promoted `approaches/tools/`
>
> You can pick any combination (e.g. "all", "2 and 4", "just 5"). The routine pass is 4 and 5; the others are occasional deep audits.

Wait for the user's response, then run the selected steps in order. *(Auto mode: skip the question — the steps come from `--auto`.)*

---

## Step 2 — Structural audit

*Skip this step if the user did not select it.*

Read `templates/technique.md` (technique deep-dives — prose files under `approaches/techniques/`) and `templates/record.md` (flat records — pure YAML-frontmatter files under `approaches/tools/`) from this skill's directory for the two approach-file structures; the playbooks-table structure is specified inline below.

### Playbook files

For each per-goal file in `skills/mentor/playbooks/`:

- `*Last verified: YYYY-MM-DD*` on line 2
- One `playbooks/<slug>.md` file per goal in problem-mode.md's classification table (and no extras)
- Each file: a `**Hidden gem:**` line naming an approach that appears in its rows, and a ranked shortlist with sequential numbering, at least 3 rows, and approach links that resolve
- Row content: "Best when" is one short clause; "Why it fits" is one sentence of goal-specific judgment — flag rows that have drifted into generic filler
- The shortlist is curated, not exhaustive: top picks plus the hidden gem. Every technique file must still be ranked by at least one playbooks file (the audit's orphan check)

### Technique files

For each technique deep-dive in `skills/mentor/approaches/techniques/`:

**Section order check** — verify these sections exist in this order:
1. `# [Title]`
2. `*Last verified: YYYY-MM-DD*` (line 2)
3. `## What It Is`
4. `## Why It Works`
5. `## When to Use It`
6. `## When NOT to Use It`
7. `## How It Works`
   - `### Basic (Beginner)`
   - `### Composing with Other Approaches (Intermediate)`
   - `### Advanced Patterns`
8. `## Common Pitfalls`
9. `## Real-World Example` — optional; kept only where the example embeds exact syntax not shown elsewhere in the file. When present it sits between Common Pitfalls and Sources.
10. `## Sources` (at least one entry, each a markdown link with a one-line description)
11. `## Signals` (both a `- Setup:` and a `- Session:` line; `—` for a tier with no signal)

**Line count check** — flag files under 40 lines.

**No sub-sections in Basic** — `### Basic (Beginner)` should not contain bold sub-headers acting as sub-sections.

### Record files and the index

For the flat record files in `skills/mentor/approaches/tools/` (pure YAML frontmatter — `kind:` plugin, integration, or doc; the filename is its id) and the marketplace directory (`marketplace.md`):

- Records: valid `---`-delimited YAML with a `last_verified: YYYY-MM-DD` field; free-text values double-quoted. Techniques and playbooks keep the `*Last verified:*` line-2 marker (`*Last synced*` for the plugin catalog)
- `approaches/index.md` is **generated** — never hand-edit it. After changing playbooks rows, a technique's `## Signals` section, or a record file, run `go -C tools/approaches-index run .` and commit the regenerated file (CI fails on a stale index)
- Every record (any kind) is a ranked row in at least one playbooks table and carries NO inline `goals:`/`best_when:` — both derive from its rows. Capability lines (`**Plugins:**`/`**Built-ins:**`/`**Integrations:**`) are forbidden: the ranking is the only routing surface
- Built-in slash commands have no records at all — each lives inside its covering technique deep-dive (`/code-review` et al. in `built-in-review-skills`, `/goal`+`/loop` in `autonomous-loops`, `/schedule` in `scheduled-agents`, `/init` in `project-memory`)
- Every `kind: plugin` record has no `marketplace.md` row (promotion removes the directory row)

The Go audit (`go -C tools/structural-audit run .`) enforces all of this deterministically — run it first; this checklist explains its failures.

### Cross-references

- Every ranked-row link (`../approaches/techniques/<name>.md` or `../approaches/tools/<name>.md`) in goal files must point to an existing file.
- Every approach file must be referenced by at least one goal file.
- Report orphan and missing approach files.

### Staleness

- Parse `*Last verified: YYYY-MM-DD*` (techniques, playbooks) and `last_verified:` (records) from every file.
- Flag files with a date older than 90 days from today.

### Output

Present the audit results:

```
## Audit Results

- Goals: N files, N issues
- Approaches: N files, N issues
- Cross-references: N orphans, N missing
- Staleness: N files older than 90 days

### Issues
[List each issue with file path, issue type, and details]
```

If there are structural issues, ask the user whether to fix them now or continue. Apply confirmed fixes before proceeding. *(Auto mode: apply only unambiguous fixes, report the rest.)*

---

## Step 3 — Content verification

*Skip this step if the user did not select it.*

Ask the user:

> Verify all files or a specific one?
>
> - **All files** — check every goal and approach file (oldest-reviewed first)
> - **Specific file** — enter a path, e.g. `approaches/plan-mode.md`, or `playbooks/<goal>.md` for one goal's rankings

Wait for the user's response. *(Auto mode: skip the question — process the `--files` N oldest-verified files, plus every promoted record `go -C tools/catalog-drift run .` flags as upstream-changed.)*

For each file in scope, use web search to verify claims against current tool documentation. Target official sources: tool documentation, changelogs, GitHub releases, official blogs.

For promoted `kind: plugin` records, run `go -C tools/catalog-drift run .` first: its upstream check compares each record's `last_verified` against the last commit to the plugin's path in the marketplace repo. An unflagged record has had no upstream change since its hands-on evaluation — the evidence stands, and only the manifest description needs a glance. A flagged (`!`) record's hands-on claims predate upstream changes: re-check its facts and pitfalls against the plugin's current components before moving its date.

### For approach files, check:

- **Feature accuracy**: Does Claude Code actually support what's described? Has the feature been renamed, changed, or removed?
- **Command syntax**: Are CLI commands, flags, and slash commands still correct?
- **Missing features**: Are there significant new Claude Code features related to this approach that the file doesn't mention?
- **"How It Works" accuracy**: Are the step-by-step instructions still correct?

### For the goal routing files (`playbooks/<goal>.md`), also check:

- **Rankings**: Is the most broadly useful approach ranked first per goal?
- **Hidden gems**: Does each still name the most non-obvious high-value approach, present in its section's rows?
- **Missing approaches**: Cross-check against all approach files — is any relevant approach missing from a goal's shortlist?
- **Misplaced approaches**: Are any rows a poor fit for their goal?

### Output

For each file with issues:

```
### [filename] — NEEDS UPDATE
- [what's wrong] → [proposed fix] (source: [URL])
```

Ask the user which fixes to apply. For each confirmed fix, edit the file — and because this step verifies the file's claims against current docs, update its `*Last verified*` date to today.

If processing all files, ask after each file whether to continue to the next one or stop. *(Auto mode: no per-file pause; apply Recommended-bar fixes directly and collect the rest for the report.)*

---

## Step 4 — Process new changelogs

*Skip this step if the user did not select it.*

This is the routine maintenance path. New Claude Code capabilities are announced in the official what's-new digests, and `skills/mentor/processed-changelogs.md` is the ledger of which digests have already been incorporated into the catalog. The weekly digest slug (e.g. `2026-w26`) is the stable unit of processing.

1. Fetch the digest index and collect the weekly slugs:

   ```
   curl -s https://code.claude.com/docs/en/whats-new/index.md
   ```

2. Read `skills/mentor/processed-changelogs.md`. Any digest in the index but not in the ledger is unprocessed.

3. For each unprocessed digest, oldest first, fetch `https://code.claude.com/docs/en/whats-new/<slug>.md` and triage each announced change:

   - **A changed command, flag, or behavior** → find the covering files (grep `skills/mentor/` for the feature name and its aliases — check synonyms and spelling variants, e.g. "auto memory" vs "auto-memory") and update them. The digest itself is an official source; quote it as the evidence.
   - **A new workflow-relevant capability** → add it to the closest approach file, or scaffold a new approach from the templates if it is a distinct recommendable technique. If it is not worth covering, say why in the ledger row. Keep the catalog consistent — CI fails otherwise: a new technique needs at least one problems row and a `## Signals` section; a new built-in command folds into its covering technique deep-dive (or a new technique if none covers it); a new integration gets its own `approaches/tools/<id>.md` record file plus a ranked row in at least one playbooks table. Then regenerate the compiled index (`go -C tools/approaches-index run .`) and commit it alongside the change.
   - **UX, enterprise-admin, install, or surface changes** → no action; the catalog is workflow-focused.

4. Append one row per digest to the ledger — slug, today's date, one-line outcome ("updated approaches/x.md and playbooks/debugging.md", "no workflow-relevant changes", ...) — and update the ledger's `*Updated*` date. Every processed digest gets a row, including no-op weeks; a gap in the ledger means unprocessed work.

For breaking changes the digests may not mention (renamed flags, removed features), also skim the release-level changelog at `https://github.com/anthropics/claude-code/blob/main/CHANGELOG.md` for the same period — these matter to "Try it now" prompts even when they are not "notable".

### Output

```
## Changelog Processing

**Digests processed:** [slugs] ([N] were already in the ledger)

### [slug]
- [change description] → [action taken / no action + why] (source: [digest URL])

### Suggested Actions (not auto-applied)
- [ ] Update [file] — [what to change]
- [ ] Consider new approach file for [new feature]
```

Present the suggested actions and ask the user which ones to apply *(auto mode: process every unprocessed digest, apply only changes meeting the "Recommended changes" bar, report the rest — and always append the ledger row; if a change was found but not applied, the row says "see report")*. For updates to existing files, do **not** update `*Last verified*` — a digest-driven edit verifies one claim, not the whole file's contents; only Step 3 (whole-file verification) or verified-at-birth authorship moves the date. For new approach files, scaffold from this skill's `templates/` directory and date them today; registry records have no template — copy the field shape of an existing record in the target registry file.

---

## Step 5 — Plugin catalog sync

*Skip this step if the user did not select it.*

Fetch the current plugin list from the marketplace manifest — the manifest is the authoritative list and includes externally-hosted partner plugins that have no directory in the repo:

```
curl -s https://raw.githubusercontent.com/anthropics/claude-plugins-official/main/.claude-plugin/marketplace.json | python3 -c "import json,sys; [print(p['name']) for p in json.load(sys.stdin)['plugins']]"
```

Extract the plugin names currently documented in TWO places: `skills/mentor/marketplace.md` (backtick-wrapped names in its tables) and the promoted `kind: plugin` records in `skills/mentor/approaches/` (their filenames). The documented set is the union — `tools/catalog-drift` computes exactly this, and additionally flags promoted records whose upstream plugin changed after their `last_verified` date (exit 1 on either kind of drift). Those upstream flags are Step 3 work — report them, never resolve them inside Step 5.

Compare against the manifest:

- **New plugins** — in the manifest but documented in neither place → add a `marketplace.md` row (new plugins always enter through the directory; promotion to `approaches/` is a separate, human editorial decision per the promotion rule in `marketplace.md`'s header)
- **Removed plugins** — documented (directory row or promoted record) but no longer in the manifest. A directory row is removed outright. A **promoted record** leaving the marketplace is NEVER auto-deleted: flag it for a human decision (demote, delete, or keep with a delisted note) — profile rows may point at it
- **Changed metadata on promoted plugins** — if the manifest description of a promoted plugin materially changed, flag its `approaches/tools/<id>.md` for re-verification (its facts are hands-on claims; Step 3 re-verifies them like any other solution file). Upstream *content* changes are caught deterministically by `catalog-drift`'s commit-date check — its `!` flags route to Step 3 the same way

For each new plugin, take its `description` (and `author`, to label Anthropic-built vs external) from the same manifest JSON — no per-plugin fetch needed. If a batch of new plugins is very large (e.g. a marketplace expansion), still list every name in the report, but it is acceptable to add table rows in slices across runs, oldest-known first, noting the remaining backlog count in the report.

### Output

```
## Plugin Catalog Sync

**Source:** anthropics/claude-plugins-official (fetched today)

### New plugins (not yet documented)
- `<name>` (Anthropic-built / External) — <description>
  → Suggested table row: | `<name>` | <short description> | `<goal slug>` | ☑️ desk-checked — <reason> |

### Removed plugins (documented but no longer in the manifest)
- `<name>` — remove from the relevant marketplace.md table
- `<name>` (PROMOTED — human decision needed) — approaches/tools/<name>.md still exists; demote, delete, or keep with a delisted note

### No changes
(if lists match)
```

Ask the user which additions and removals to apply *(auto mode: apply all directory additions/removals — the API is authoritative; promoted-record removals and re-verification flags are always report-only)*. For confirmed changes, edit `marketplace.md` and update its `*Last synced*` date to today. Promoted `approaches/tools/<id>.md` files are maintained by Step 3 (content verification against official docs), same as every other solution file.

Directory plugins are never listed in `playbooks/<goal>.md` (the audit forbids `**Plugins:**` lines), so directory changes need no goal-file reconciliation. If a removed plugin was PROMOTED, its ranked rows in `playbooks/<goal>.md` are part of the human decision flagged above — never auto-deleted.

The evidence rules for this step are lighter than Steps 3 and 4: the GitHub API response is authoritative — no web search needed to verify presence or absence.

---

## Evidence and confidence rules

Every proposed change in Steps 3 and 4 must include inline evidence (Step 5 is exempt — the GitHub API response is the authoritative source): the source URL **and** a direct quote from that source supporting the change. Not "according to the docs" — the actual text. The user must be able to verify any claim in seconds.

### Source tiers

| Tier | Examples | Sufficient alone? |
|------|----------|-------------------|
| **Official** | Anthropic docs, Claude Code changelog, tool's own documentation site, official blog posts, GitHub releases by the tool's maintainers | Yes |
| **Community** | Stack Overflow answers, third-party blog posts, GitHub issues/discussions, forum threads | No — needs a second source or an official source to corroborate |

### Confidence classification

Split findings into two sections:

- **Recommended changes** — backed by at least one official-tier source with a direct quote. These are safe to apply.
- **Needs manual verification** — based only on community-tier sources, or the source is ambiguous, or the finding contradicts what the file says but the evidence isn't conclusive. List the source and quote, but flag clearly that the user should verify before applying.

### What NOT to do

- Never propose a change based on your training data alone — every change must have a web search result backing it.
- Never paraphrase a source and present the paraphrase as a quote — use the actual text.
- If web search returns no results for a claim, report "could not verify" rather than assuming the claim is correct or incorrect.
- Never silently drop a finding because it seems minor — report everything and let the user decide.

---

## General rules

- Always show proposed changes and ask before editing any file
- `*Last verified*` means verified, not edited: it moves only when a file's claims were checked against current official docs (Step 3) or the file was authored from verified sources. Mechanical edits, restructurings, and single-claim updates never touch it
- Use the templates as the source of truth for structural requirements
- When reporting issues, include the file path and line number when possible
