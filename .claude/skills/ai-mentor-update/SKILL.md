---
name: ai-mentor-update
description: >-
  Maintenance skill for the ai-mentor plugin. Audits structure, verifies
  content accuracy against current tool docs, and detects recent tool
  changes. Each step is optional — the user picks what to run.
disable-model-invocation: true
---

# AI Mentor Update

You are a content maintenance tool for the ai-mentor plugin. Your job is to keep the goal and approach files structurally consistent, factually accurate, and up to date.

All paths below are relative to the repo root.

---

## Step 1 — Setup

Run `date +%Y-%m-%d` via bash to establish today's date.

Then ask the user which steps to run:

> Which maintenance steps do you want to run?
>
> 1. **Structural audit** — check all files against templates, cross-references, staleness
> 2. **Content verification** — web-search claims in files against current tool docs
> 3. **Recent tool changes** — search changelogs for new features and breaking changes
> 4. **Plugin catalog sync** — check `claude-plugins-official` for new or removed plugins not yet reflected in `approaches/plugins.md`
>
> You can pick any combination (e.g. "all", "1 and 3", "just 4").

Wait for the user's response, then run the selected steps in order.

---

## Step 2 — Structural audit

*Skip this step if the user did not select it.*

Read `templates/goal.md` and `templates/approach.md` from this skill's directory to know the required structure.

### Goal files

For each `.md` file in `skills/ai-mentor/goals/`:

**Section order check** — verify these sections exist in this order:
1. `# [Title]`
2. `*Last reviewed: YYYY-MM-DD*` (line 2)
3. `## When You're Here`
4. `## Quick Decision Guide` (with a 3-column table)
5. `## Approaches (Ranked)`

**Approach entry check** — each `### N. Name — pitch` entry must have exactly these fields in order:
1. `**Level:**` and `**Tools:**` on the same line
2. Description paragraph
3. `**Try it now:**` with a blockquote
4. `**Why this works:**`
5. `**Pros:**` (bullet list)
6. `**Cons:**` (bullet list)
7. `**Deeper:** See \`approaches/<name>.md\``

Flag any extra fields (`Also try`, `Tip`, `Real-world example`, `When to combine`, or anything else not in this list).

**Numbering check** — approach numbers must be sequential starting from 1.

**Separator check** — `---` must appear between approach entries but not after the last one.

### Approach files

For each `.md` file in `skills/ai-mentor/approaches/`:

**Section order check** — verify these sections exist in this order:
1. `# [Title]`
2. `*Last reviewed: YYYY-MM-DD*` (line 2)
3. `## What It Is`
4. `## Why It Works`
5. `## When to Use It`
6. `## When NOT to Use It`
7. `## How It Works`
   - `### Basic (Beginner)`
   - `### Composing with Other Approaches (Intermediate)`
   - `### Advanced Patterns`
8. `## Tool Support` (table with Claude Code, OpenCode, Cursor, aider)
9. `## Common Pitfalls`
10. `## Real-World Example`
11. `## Sources` (1-3 entries, each a markdown link with a one-line description)

**Line count check** — flag files outside the 80-120 line range.

**No sub-sections in Basic** — `### Basic (Beginner)` should not contain bold sub-headers acting as sub-sections.

### Cross-references

- Every `Deeper: See \`approaches/<name>.md\`` reference in goal files must point to an existing file.
- Every approach file must be referenced by at least one goal file.
- Report orphan and missing approach files.

### Staleness

- Parse `*Last reviewed: YYYY-MM-DD*` from every file.
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

If there are structural issues, ask the user whether to fix them now or continue. Apply confirmed fixes before proceeding.

---

## Step 3 — Content verification

*Skip this step if the user did not select it.*

Ask the user:

> Verify all files or a specific one?
>
> - **All files** — check every goal and approach file (oldest-reviewed first)
> - **Specific file** — enter a path, e.g. `approaches/plan-mode.md` or `goals/debugging.md`

Wait for the user's response.

For each file in scope, use web search to verify claims against current tool documentation. Target official sources: tool documentation, changelogs, GitHub releases, official blogs.

### For approach files, check:

- **Feature accuracy**: Does the tool actually support what's described? Has the feature been renamed, changed, or removed?
- **Command syntax**: Are CLI commands, flags, and slash commands still correct?
- **Tool Support table**: Is each tool's support level (Native/Partial/None) still accurate for Claude Code, OpenCode, Cursor, and aider?
- **Missing features**: Are there significant new features related to this approach that the file doesn't mention?
- **"How It Works" accuracy**: Are the step-by-step instructions still correct?

### For goal files, also check:

- **Approach rankings**: Is the most broadly useful approach ranked first?
- **Missing approaches**: Cross-check against all approach files — is any relevant approach missing from this goal?
- **"Try it now" prompts**: Do they use current syntax and realistic file paths?
- **Quick Decision Guide**: Does the table cover the main scenarios for this goal?
- **Misplaced approaches**: Are any listed approaches a poor fit for this goal?

### Output

For each file with issues:

```
### [filename] — NEEDS UPDATE
- [what's wrong] → [proposed fix] (source: [URL])
```

Ask the user which fixes to apply. For each confirmed fix, edit the file and update its `*Last reviewed*` date to today.

If processing all files, ask after each file whether to continue to the next one or stop.

---

## Step 4 — Detect recent tool changes

*Skip this step if the user did not select it.*

Ask the user:

> How far back should I search? (default: 30 days)

Wait for the user's response. Use the provided number or default to 30.

Search for changes to AI coding tools published within that window.

| Tool | What to search for |
|------|--------------------|
| Claude Code | Changelog, release notes, new features, CLI changes, new slash commands, hooks updates, agent/workflow changes |
| OpenCode | GitHub releases, new features, plugin system changes, LSP updates |
| Cursor | Changelog, new AI features, MCP support changes |
| aider | GitHub releases, new features, model support changes |

For each change found, identify which approach and goal files it affects:

- A new feature → may need a new approach file or updates to existing ones
- A changed command or flag → "Try it now" prompts and "Basic (Beginner)" sections may be wrong
- A tool gaining/losing support for a feature → Tool Support tables need updating
- A renamed or removed feature → content may reference something that no longer exists
- A new capability category → may warrant a new goal file

### Output

```
## Recent Tool Changes

**Lookback window:** [N] days ([start date] → [today])

### [Tool Name]
- [change description] → affects: [list of files]

### Suggested Actions
- [ ] Update [file] — [what to change]
- [ ] Consider new approach file for [new feature]
- [ ] Consider new goal file for [new capability]
```

Present the suggested actions and ask the user which ones to apply. For confirmed updates to existing files, make the edits and update `*Last reviewed*` dates. For new files, scaffold them using the templates from this skill's `templates/` directory.

---

## Step 5 — Plugin catalog sync

*Skip this step if the user did not select it.*

Fetch the current plugin list from `anthropics/claude-plugins-official` using the GitHub API:

```
gh api repos/anthropics/claude-plugins-official/contents/plugins | python3 -c "import json,sys; [print(d['name']) for d in json.load(sys.stdin) if d['type']=='dir']"
gh api repos/anthropics/claude-plugins-official/contents/external_plugins | python3 -c "import json,sys; [print(d['name']) for d in json.load(sys.stdin) if d['type']=='dir']"
```

Extract the plugin names currently documented in `skills/ai-mentor/approaches/plugins.md` by reading the file and collecting all backtick-wrapped names in the `## Official Claude Code Plugins` section.

Compare the two lists:

- **New plugins** — present in the repo but not mentioned in `approaches/plugins.md`
- **Removed plugins** — mentioned in `approaches/plugins.md` but no longer in the repo

For each new plugin, fetch its description:

```
gh api repos/anthropics/claude-plugins-official/contents/plugins/<name>/.claude-plugin/plugin.json
```

Decode the base64 content and extract the `description` field.

### Output

```
## Plugin Catalog Sync

**Source:** anthropics/claude-plugins-official (fetched today)

### New plugins (not yet in approaches/plugins.md)
- `<name>` (Anthropic-built / External) — <description>
  → Suggested table row: | `<name>` | <short description> | `goals/<goal>.md` |

### Removed plugins (in approaches/plugins.md but no longer in repo)
- `<name>` — remove from the relevant table

### No changes
(if lists match)
```

Ask the user which additions and removals to apply. For confirmed changes, edit `approaches/plugins.md` and update its `*Last reviewed*` date to today.

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
- Update `*Last reviewed*` dates on every file that gets modified
- Use the templates as the source of truth for structural requirements
- When reporting issues, include the file path and line number when possible
