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
- `--files N` — Step 3 scope: process the N oldest-verified files (default: 5)

The routine weekly run is `--auto 4,5`: process new changelogs and sync the plugin catalog. Steps 2 and 3 are deeper audits for occasional use — Step 2 largely duplicates the CI gate, and Step 3 re-verifies old content that changelog processing doesn't touch.

Auto-mode overrides, in addition to skipping every question below:

- **Step 2**: apply only unambiguous structural fixes (broken separator, wrong field order); report anything requiring judgment instead of fixing it. Note: CI also runs the structural audit (`go -C tools/structural-audit run .`) as a deterministic gate — prefer reporting over creative fixing.
- **Step 3**: process the N oldest-verified files with no per-file pause. Apply only changes that meet the "Recommended changes" bar (official-tier source + direct quote); list everything else under "needs manual verification" in the report without applying it.
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
> 5. **Plugin catalog sync** — check `claude-plugins-official` for new or removed plugins not yet reflected in `plugins.md`
>
> You can pick any combination (e.g. "all", "2 and 4", "just 5"). The routine pass is 4 and 5; the others are occasional deep audits.

Wait for the user's response, then run the selected steps in order. *(Auto mode: skip the question — the steps come from `--auto`.)*

---

## Step 2 — Structural audit

*Skip this step if the user did not select it.*

Read `templates/approach.md` from this skill's directory for the technique-file structure; the problems-table structure is specified inline below.

### Problem files

For each per-goal file in `skills/mentor/problems/`:

- `*Last verified: YYYY-MM-DD*` on line 2
- One `problems/<slug>.md` file per goal in problem-mode.md's classification table (and no extras)
- Each file: a `**Hidden gem:**` line naming an approach that appears in its rows, and a ranked shortlist with sequential numbering, at least 3 rows, and approach links that resolve
- Row content: "Best when" is one short clause; "Why it fits" is one sentence of goal-specific judgment — flag rows that have drifted into generic filler
- The shortlist is curated, not exhaustive: top picks plus the hidden gem. Every technique file must still be ranked by at least one problems file (the audit's orphan check)

### Technique files

For each technique file in `skills/mentor/solutions/` (the deep-dives — any solution file *without* a `kind:` line):

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

For the flat record files in `skills/mentor/solutions/` (a solution file *with* a `kind:` line — built-in command, integration, or doc; the filename is its id) and the external plugin catalog (`plugins.md`):

- `*Last verified: YYYY-MM-DD*` on line 2 (`*Last synced*` for the plugin catalog)
- `solutions/index.md` is **generated** — never hand-edit it. After changing problems rows, a technique's `## Signals` section, or a record file, run `go -C tools/solutions-index run .` and commit the regenerated file (CI fails on a stale index)
- Every `kind: builtin-command` record is referenced by at least one `**Built-ins:**` line, every `kind: integration`/`doc` record by an `**Integrations:**` line, and every Built-ins/Integrations token resolves to a `solutions/<id>.md` file
- Every `goals` slug in any record resolves to an existing `problems/<goal>.md`

The Go audit (`go -C tools/structural-audit run .`) enforces all of this deterministically — run it first; this checklist explains its failures.

### Cross-references

- Every `Deeper: See \`solutions/<name>.md\`` reference in goal files must point to an existing file.
- Every approach file must be referenced by at least one goal file.
- Report orphan and missing approach files.

### Staleness

- Parse `*Last verified: YYYY-MM-DD*` from every file.
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
> - **Specific file** — enter a path, e.g. `solutions/plan-mode.md`, or `problems/<goal>.md` for one goal's rankings

Wait for the user's response. *(Auto mode: skip the question — process the `--files` N oldest-verified files.)*

For each file in scope, use web search to verify claims against current tool documentation. Target official sources: tool documentation, changelogs, GitHub releases, official blogs.

### For approach files, check:

- **Feature accuracy**: Does Claude Code actually support what's described? Has the feature been renamed, changed, or removed?
- **Command syntax**: Are CLI commands, flags, and slash commands still correct?
- **Missing features**: Are there significant new Claude Code features related to this approach that the file doesn't mention?
- **"How It Works" accuracy**: Are the step-by-step instructions still correct?

### For the goal routing files (`problems/<goal>.md`), also check:

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
   - **A new workflow-relevant capability** → add it to the closest approach file, or scaffold a new approach from the templates if it is a distinct recommendable technique. If it is not worth covering, say why in the ledger row. Keep the registry consistent — CI fails otherwise: a new technique needs at least one problems row and a `## Signals` section; a new built-in command or integration gets its own `solutions/<id>.md` record file, referenced from the relevant `**Built-ins:**` / `**Integrations:**` line. Then regenerate the compiled index (`go -C tools/solutions-index run .`) and commit it alongside the change.
   - **UX, enterprise-admin, install, or surface changes** → no action; the catalog is workflow-focused.

4. Append one row per digest to the ledger — slug, today's date, one-line outcome ("updated solutions/x.md and problems/debugging.md", "no workflow-relevant changes", ...) — and update the ledger's `*Updated*` date. Every processed digest gets a row, including no-op weeks; a gap in the ledger means unprocessed work.

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

Extract the plugin names currently documented in `skills/mentor/plugins.md` by reading the file and collecting all backtick-wrapped names in its tables.

Compare the two lists:

- **New plugins** — present in the repo but not mentioned in `plugins.md`
- **Removed plugins** — mentioned in `plugins.md` but no longer in the manifest

For each new plugin, take its `description` (and `author`, to label Anthropic-built vs external) from the same manifest JSON — no per-plugin fetch needed. If a batch of new plugins is very large (e.g. a marketplace expansion), still list every name in the report, but it is acceptable to add table rows in slices across runs, oldest-known first, noting the remaining backlog count in the report.

### Output

```
## Plugin Catalog Sync

**Source:** anthropics/claude-plugins-official (fetched today)

### New plugins (not yet in plugins.md)
- `<name>` (Anthropic-built / External) — <description>
  → Suggested table row: | `<name>` | <short description> | `<goal slug>` | ☑️ desk-checked — <reason> |

### Removed plugins (in plugins.md but no longer in the manifest)
- `<name>` — remove from the relevant table

### No changes
(if lists match)
```

Ask the user which additions and removals to apply *(auto mode: apply all — the API is authoritative)*. For confirmed changes, edit `plugins.md` and update its `*Last synced*` date to today.

Then reconcile the `**Plugins:**` lines in `problems/<goal>.md`: a removed plugin's token is deleted from any line naming it (mechanical — apply directly); a new plugin is only *suggested* for a goal's line when it clearly beats the current picks (editorial — list under suggested actions, never auto-applied). The structural audit fails on routing plugin tokens missing from the catalog, so removals must not be skipped.

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
