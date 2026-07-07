# Mentor Profile Schema
*Last reviewed: 2026-07-03*

The mentor's memory of one engineer: what they know, what they've been shown, what they adopted, what they declined. This is what makes discovery work — the ignorance map is "the catalog minus this profile minus live signals," and the never-repeat rule depends on it.

## Where it lives

A single user-level file: `~/.ai-mentor/profile.md`, created on the mentor's first invocation.

It is deliberately **person-level, not per-repo**: knowledge belongs to the engineer, and someone working across ten repositories must never be re-taught in repo #7 what they adopted in repo #2. Per-repo state needs no storage at all — setup signals are re-read from the current repo's disk every invocation (see `adoption-signals.md`), and combining person-level knowledge with this-repo disk state is what enables the transfer move: "you use hooks in your other projects but not here — want the same test hook wired up?" Rare repo-scoped nuances ("declined loops for repo X — slow CI") go in the Note column with the repo named.

The file is machine-local (never committed, never leaves the machine) and plain markdown the engineer can open, edit, or delete at any time — the skill should mention the path when it first creates it. A hand-edit by the user is authoritative and overrides anything the mentor inferred.

## Permissions: zero prompts, via skill frontmatter

The location is `~/.ai-mentor/`, **not** `~/.claude/`, for a reason verified empirically (2026-07-03, Claude Code v2.x): *writes* to files under `~/.claude/` are treated as sensitive, and that built-in protection **overrides even an exactly-matching allow rule at any settings level** — every profile write there would prompt, forever. (*Reads* under `~/.claude/` respect allow rules normally; the protection is edit-specific.) Outside `~/.claude/`, standard allow rules silence prompts completely.

The whole permission story fits in the mentor's own frontmatter — verified to work, including auto-creating the missing `~/.ai-mentor/` directory on the first write:

```yaml
allowed-tools:
  - Read(~/.ai-mentor/**)
  - Edit(~/.ai-mentor/**)
```

No settings.json change, no consent prompt, no setup step: the grant is scoped to this skill's execution only, which is least-privilege by construction. For transparency, the mentor still announces the file on first creation ("I keep a profile at `~/.ai-mentor/profile.md` so I never re-teach you things — it's yours to edit or delete").

Facts future maintainers must not "fix" (each verified by test, not assumption):

- The rule family must be `Edit(...)`: path-scoped `Edit` rules cover all file-editing tools including Write, whereas a path-scoped `Write(...)` rule does not match and is silently ineffective.
- The `~/` anchor works in frontmatter and settings rules alike and keeps rules portable across users — do not expand it to an absolute path. Conversely, a rule built from `${CLAUDE_PLUGIN_ROOT}` needs a leading `/` (yielding `//abs/path`), because a single leading slash is project-root-relative.
- The *tool calls* must use the literal `~/.ai-mentor/...` path too, never an absolute home path inferred from repo paths in context (verified 2026-07-07: in an isolated-HOME session the model guessed `/Users/<name>/...` from the plugin path, the `~`-anchored grant expanded against the real `$HOME` and did not match, and every profile read/write was denied). The file tools expand `~` against the session's `$HOME`, which is always the profile's true location.
- Never use Bash `mkdir` for the profile directory — the Write tool creates missing parents under the same Edit rule; a `mkdir` would trigger a separate Bash prompt.
- Write the profile immediately when a status changes, within the mentor's own flow — never defer to "session end". Writes are silent, so there is nothing to batch, and the frontmatter grant is only guaranteed active while the skill is executing; a deferred flush many turns later would be betting on unverified permission-scope semantics.

## Format

```markdown
# Mentor Profile
*Updated: YYYY-MM-DD*

Level: getting-started | comfortable | advanced — one line of evidence for the calibration
Last new-capability check: <what's-new week slug, e.g. 2026-w26 — always a week slug; on profile creation set it to the current ISO week (first meeting is the baseline), never "never">

| Capability | Status | Date | Note |
|------------|--------|------|------|
| hooks-as-workflow | adopted | 2026-07-03 | PostToolUse test hook in settings, fired 40×/week |
| autonomous-loops | shown | 2026-07-03 | Demoed /loop on the flaky retry test |
| fan-out-workflows | declined | 2026-07-03 | "Too token-heavy for us" |
```

`Capability` is a registry id: an approach file basename (kind: technique — the enumerable set in `adoption-signals.md`), a built-in command or integration id from `registry/`, or a marketplace plugin name from `references/official-plugins.md`. All kinds share the table and the statuses; ids are additive across versions and existing rows are never orphaned by schema changes — a user's hand-edited row always stays valid. A capability with no row is **unknown** — a discovery candidate (plugins only when stack/goal-relevant).

## Statuses and transitions

- **unknown → shown**: the mentor taught or demonstrated it this session.
- **shown → adopted**: a positive signal is later observed (see `adoption-signals.md`) or the user confirms they use it.
- **unknown → adopted**: a signal shows they already knew it — record silently, never "teach" it.
- **any → declined**: the user waved it off. Record the reason verbatim if given. Never re-offer a declined capability unless the user asks, or the reason no longer applies (e.g. declined for a missing plan feature they now have) — and then at most once, naming why it's being raised again.

Transitions move forward only; the exceptions are user edits (always win) and evidence contradicting the profile (a "declined" capability now configured in settings → flip to adopted).

## Rules for the mentor

1. **Evidence beats memory.** Re-check the cheap setup/repo signals every invocation; the profile accumulates only what can't be re-checked (session-tier signals, shown/declined history).
2. **Never repeat.** Don't re-teach `shown`, don't re-offer `declined`, don't explain `adopted`. For a discovery product, re-showing a known capability is the primary failure mode.
3. **Follow up on `shown`.** A capability shown but not adopted after a while is the opening move of the next session: "Last time I showed you X — did it stick?" The answer converts it to adopted, declined, or a re-teach with a different angle.
4. **`Last new-capability check`** anchors "new since you last checked": anything landed in `references/processed-changelogs.md` after that week is fresh discovery material for this user, regardless of their level. Update it whenever new capabilities are surfaced.
5. **Keep it small.** One line per capability, one profile per engineer per machine. This is a ledger, not a diary — details belong in the Note column or nowhere.
