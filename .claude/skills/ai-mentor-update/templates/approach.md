# Approach Template

Use this template when adding a new approach file to `approaches/`.

Each approach file is the deep-dive reference for a single AI workflow technique. Developers land here from a goal file's "Deeper" link when they want the full picture: what it is, how to use it at every skill level, how to compose it with other approaches, and what to watch out for.

---

## Required Structure

```markdown
# [Approach Name]
*Last reviewed: YYYY-MM-DD*

## What It Is

[2-3 sentences in plain language. No jargon. A developer who has
never heard of this approach should understand what it does after
reading this paragraph.]

## Why It Works

[2-4 sentences explaining the PRINCIPLE — the lasting insight, not
the mechanics. Why does this approach produce better results than
doing the task manually? What cognitive or process bottleneck does
it address?]

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

## Tool Support

| Tool | Support | Notes |
|------|---------|-------|
| Claude Code | Native/Partial/None | [specific details] |
| OpenCode | Native/Partial/None | [specific details] |
| Cursor | Native/Partial/None | [specific details] |
| aider | Native/Partial/None | [specific details] |

[Always include all four tools. Use "Native" for built-in support,
"Partial" for limited or plugin-based support, "None" for no support.]

## Common Pitfalls

- **[Pitfall name]**: [1-2 sentences describing the mistake and
  how to avoid it. Bold the pitfall name, then explain.]

[2-4 pitfalls. These should be mistakes real developers make,
not theoretical warnings.]

## Real-World Example

[A concrete, end-to-end scenario with specific details — file
names, module names, realistic problem descriptions. Include a
code block showing the prompt or command. Follow with a narrative
describing what happens and the outcome. The example should feel
like something that actually happened, not a contrived demo.]

## Sources

- [Source title](URL) — [1-line description of what this covers]
- [Source title](URL) — [1-line description of what this covers]

[1-3 official external sources. Prefer official documentation
(Anthropic docs, tool docs, maintainer sites) over blog posts
or community content. These are the references the content was
verified against.]
```

## Rules

- **80-120 lines** per file — long enough to be useful, short enough to fit in a context window alongside other files
- **"Why It Works" teaches a principle** — not "it saves time" but the specific cognitive or process insight that makes the approach effective
- **Real-World Example must feel real** — specific file names, realistic module names, a plausible problem with a concrete resolution
- **Tool Support table is always present** — even if most tools show "None," this sets expectations
- **Common Pitfalls come from experience** — describe mistakes developers actually make, not obvious warnings
- **`## Sources` lists 1-3 official external links** — prefer official docs (Anthropic, tool maintainers) over community content. Each entry: `[Title](URL) — one-line description`
- **No extra sections** beyond the template — keep the structure consistent across all approach files
- **No sub-sections within "Basic (Beginner)"** — if the approach has distinct modes (e.g., "using" vs. "setting up"), pick the most common one for Basic and put variations in Advanced Patterns
- **`*Last reviewed: YYYY-MM-DD*`** on line 2 — update whenever the file is reviewed or modified
- **Composing patterns reference other approaches by name** — this helps developers discover related techniques
