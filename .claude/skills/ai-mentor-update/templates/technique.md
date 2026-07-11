# Technique Deep-Dive Template

Use this template when adding a **technique deep-dive** to `approaches/` — a file WITHOUT a `kind:` line, documenting a methodology we author and own. For the other format in the same directory (flat records: promoted plugins, integrations, docs — files WITH a `kind:` line), use `templates/record.md` instead. The two formats are deliberate: a deep-dive carries our own pedagogy; a record carries only verified facts about an external artifact, because the verifiability bar caps what we can honestly write about things we don't maintain.

Each technique file is the deep-dive reference for a single AI workflow technique. Developers land here from a goal file's "Deeper" link when they want the full picture: what it is, how to use it at every skill level, how to compose it with other approaches, and what to watch out for.

---

## Required Structure

```markdown
# [Approach Name]
*Last verified: YYYY-MM-DD*

## What It Is

[2-3 sentences in plain language. No jargon. A developer who has
never heard of this approach should understand what it does after
reading this paragraph.]

## Why It Works

[ONE sentence stating the PRINCIPLE — the lasting insight, not the
mechanics. The model consuming this file already knows the generic
argument for the technique; one sharp sentence beats four soft ones.]

## When to Use It

- [Specific scenario 1]
- [Specific scenario 2]
- [Specific scenario 3]

[3-5 bullets. Each describes a concrete situation where this
approach is the right choice. Start with the verb or the signal
that tells the developer "this is for you."]

## When NOT to Use It

- [Anti-pattern 1]
- [Anti-pattern 2]

[2-3 bullets. Be honest about where this approach is overkill,
too slow, or the wrong tool for the job.]

## How It Works

### Basic (Beginner)

[Step-by-step walkthrough. 3-5 numbered steps. Include at least
one concrete example with a code block or command. A developer
should be able to follow these steps and get a result on their
first try.]

### Composing with Other Approaches (Intermediate)

[2-3 bullet points, each describing a specific combination with
another approach. Format: **This approach plus [other approach]**:
description of how they combine and what the combination achieves.]

### Advanced Patterns

[2-3 bullet points describing power-user techniques, edge cases,
or non-obvious applications. These are for developers who already
use the basic approach and want to push further.]

## Common Pitfalls

- **[Pitfall name]**: [1-2 sentences describing the mistake and
  how to avoid it. Bold the pitfall name, then explain.]

[2-4 pitfalls. These should be mistakes real developers make,
not theoretical warnings.]

## Real-World Example

[OPTIONAL — include only when the example embeds exact syntax
(a settings JSON, a workflow YAML, an orchestration script) that
appears nowhere else in the file. A narrative that only re-uses
commands already shown above is filler for the model that consumes
this file: omit the section entirely in that case.]

## Sources

- [Source title](URL) — [1-line description of what this covers]
- [Source title](URL) — [1-line description of what this covers]

[At least one official external source. Prefer official documentation
(Anthropic docs, tool docs, maintainer sites) over blog posts
or community content. These are the references the content was
verified against.]

## Signals

- Setup: [observable disk/config evidence the engineer already uses this — a glob or grep against the project or ~/.claude; use — if none exists]
- Session: [conversation evidence — a command in the transcript, a described habit]
```

## Rules

- **At least 40 lines** per file — long enough to be useful; keep it focused so it fits in a context window alongside other files
- **"Why It Works" is exactly one sentence** — not "it saves time" but the specific cognitive or process insight that makes the approach effective
- **Real-World Example is optional and must earn its place** — include it only when it embeds exact syntax shown nowhere else in the file; otherwise omit the section
- **Common Pitfalls come from experience** — describe mistakes developers actually make, not obvious warnings
- **`## Sources` lists at least one official external link** — prefer official docs (Anthropic, tool maintainers) over community content. Each entry: `[Title](URL) — one-line description`
- **No extra sections** beyond the template — keep the structure consistent across all approach files
- **`## Signals` is required and machine-consumed** — the solutions-index generator compiles it into `approaches/index.md`; both lines must be present (`—` for a tier with no signal). Regenerate the index after editing it
- **No sub-sections within "Basic (Beginner)"** — if the approach has distinct modes (e.g., "using" vs. "setting up"), pick the most common one for Basic and put variations in Advanced Patterns
- **`*Last verified: YYYY-MM-DD*`** on line 2 — moves ONLY when the file's claims are verified against current official docs, or at creation from verified sources; never on mechanical edits
- **Composing patterns reference other approaches by name** — this helps developers discover related techniques
