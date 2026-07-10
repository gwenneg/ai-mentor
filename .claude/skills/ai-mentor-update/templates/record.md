# Record Template

Use this template when adding a **flat record** to `approaches/` — a file WITH a `kind:` line (promoted plugin, integration, or doc), documenting an external artifact whose depth lives at its source. For technique deep-dives (files without a `kind:` line), use `templates/approach.md` instead.

A record is a verified fact sheet, deliberately short: every line is a claim we can stand behind (verified against official docs, or our own hands-on evaluation). Do not pad it toward the deep-dive shape — if we ever own enough pedagogy about a capability to teach it ourselves, that is the signal to grow it into a technique deep-dive (the deep-research precedent), not to bloat the record.

---

## Required Structure

```markdown
# [id — must equal the filename]
*Last verified: YYYY-MM-DD*

kind: plugin | integration | doc
composes_with: [approach ids this pairs with, comma-separated]
install: /plugin install [id]@claude-plugins-official
facts: [What it is and what the evaluation showed, 1-3 sentences.
Hands-on evidence quoted or tightly paraphrased — never smoothed over.]
session_signal: [conversation evidence that the user already uses this —
installed and visible in the session, or its commands/tools ran]
pitfalls:
- [Sharp edges from the evaluation or the official docs, verbatim-honest.]
source: [canonical URL — the marketplace repo for plugins, official docs otherwise]
```

## Rules

- **The filename is the `id`** — there is no `id:` field; one capability, one file, by construction
- **No `goals:` or `best_when:` fields** — both derive from the record's ranked rows in `playbooks/` (the generator errors on inline copies). Every record must be a ranked row in at least one playbook
- **`kind:` is a semantic label, never a routing tier** — records are reached through the ranking exactly like techniques
- **`install:` is for plugins only** — integrations/docs describe their concrete setup step inside `facts:` or `pitfalls:` instead
- **Promoted plugins must not retain a `marketplace.md` row** (the audit enforces this), and their `facts:` must carry the hands-on evaluation date and findings
- **`*Last verified:*` moves only when the record's claims are re-verified** against the source — never on mechanical edits
- **Regenerate the index** after adding or editing a record: `go -C tools/approaches-index run .`
