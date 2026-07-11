# Record Template

Use this template when adding a **flat record** to `approaches/tools/` — a pure YAML-frontmatter file documenting an external artifact (promoted plugin, integration, or doc) whose depth lives at its source. For technique deep-dives (prose files under `approaches/techniques/`), use `templates/technique.md` instead.

A record is a verified fact sheet, deliberately short: every value is a claim we can stand behind (verified against official docs, or our own hands-on evaluation). Do not pad it toward the deep-dive shape — if we ever own enough pedagogy about a capability to teach it ourselves, that is the signal to grow it into a technique deep-dive (the deep-research precedent), not to bloat the record.

---

## Required Structure

The file is a single YAML frontmatter block and nothing else — no title, no body (the filename is the id):

```yaml
---
kind: plugin | integration | doc
last_verified: YYYY-MM-DD
composes_with:
  - approach-ids
  - this-pairs-with
install: /plugin install <id>@claude-plugins-official
facts: "What it is and what the evaluation showed, 1-3 sentences. Hands-on evidence quoted or tightly paraphrased — never smoothed over."
session_signal: "Conversation evidence that the user already uses this — installed and visible in the session, or its commands/tools ran."
source: <canonical URL — the marketplace repo for plugins, official docs otherwise>
pitfalls:
  - "Sharp edges from the evaluation or the official docs, verbatim-honest."
---
```

## Rules

- **Valid YAML, always** — free-text values (`facts`, `session_signal`, `pitfalls` items) are double-quoted because they routinely contain `": "`, which breaks plain scalars; bare identifiers, dates, lists, and URLs stay unquoted
- **The filename is the `id`** — there is no `id:` field; one capability, one file, and the generator rejects an id that exists in both subfolders
- **No `goals:` or `best_when:` fields** — both derive from the record's ranked rows in `playbooks/` (the generator errors on inline copies). Every record must be a ranked row in at least one playbook
- **`kind:` is a semantic label, never a routing tier** — records are reached through the ranking exactly like techniques
- **`install:` is for plugins only** — integrations/docs describe their concrete setup step inside `facts:` or `pitfalls:` instead
- **Promoted plugins must not retain a `marketplace.md` row** (the audit enforces this), and their `facts:` must carry the hands-on evaluation date and findings
- **`last_verified:` moves only when the record's claims are re-verified** against the source — never on mechanical edits
- **Regenerate the index** after adding or editing a record: `go -C tools/approaches-index run .`
