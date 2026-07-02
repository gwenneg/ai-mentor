# Goal Template

Use this template when adding a new goal file to `goals/`.

Each goal file maps a developer's engineering goal to ranked AI workflow approaches. The developer reads the file to understand which approaches fit their situation, with enough detail to try each one immediately.

---

## Required Structure

```markdown
# [Goal Name]
*Last reviewed: YYYY-MM-DD*

## When You're Here

[1-2 paragraphs describing the situation. What brings an engineer here?
What is the core challenge? Write in second person ("You're...").
Keep it concrete — name specific symptoms or triggers.]

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| [specific scenario] | [approach name] | [one-line reason] |
| [specific scenario] | [approach name] | [one-line reason] |
| [specific scenario] | [approach name] | [one-line reason] |

[3-5 rows. Each row is a quick shortcut for developers who already
know their specific situation and don't need to read the full list.]

**Hidden gem:** [Approach Name] — [one sentence: the non-obvious fit
for this goal and why most developers miss it.]

## Approaches (Ranked)

[Repeat the following block for each approach, separated by `---`]

### N. [Approach Name] — [one-line pitch]
**Level:** Beginner|Intermediate|Advanced

[2-3 sentences: what this approach does for THIS specific goal.
Not a generic description of the approach — explain how it helps
with the goal at hand.]

**Try it now:**
> [A concrete prompt the developer can paste into Claude Code.
>  Must be specific to a realistic scenario for this goal.
>  Include file paths, module names, and concrete details.]

**Why this works:** [1-2 sentences explaining the underlying principle.
This is the educational content — teach something lasting, not just
mechanics.]

**Pros:**
- [2-3 short bullets]

**Cons:**
- [1-2 short bullets — be honest about limitations]

**Deeper:** See `approaches/<approach-name>.md`

---
```

## Rules

- **Rank approaches** from most broadly useful to most specialized for this goal
- **Hidden gem names one approach from this file's ranked list** — the curated surprising pick the skill always surfaces alongside the #1 ranked approach. Prefer an approach ranked 4th or lower: the point is surfacing what users would not find on their own
- **3-7 approaches per goal** — fewer than 3 means the goal is too narrow; more than 7 is overwhelming
- **"Try it now" prompts must be specific** — include file paths, module names, concrete details. Never use generic placeholders like `<your-file>`
- **"Why this works" is not optional** — every approach must teach something
- **Level badges** reflect the approach complexity, not the developer's skill:
  - Beginner = works with basic Claude usage
  - Intermediate = requires specific setup or tool configuration
  - Advanced = requires composing multiple tools or deep expertise
- **No extra fields** in approach entries — stick to the standard set (Level, Try it now, Why this works, Pros, Cons, Deeper). Additional context belongs in the approach file itself
- **`---` separator** between approach entries, but not after the last one
- **`*Last reviewed: YYYY-MM-DD*`** on line 2 — update whenever the file is reviewed or modified
