# Visual Artifacts
*Last verified: 2026-07-12*

## What It Is

Visual Artifacts turn Claude Code's terminal output into rendered, shareable web pages. Claude writes an HTML or Markdown file and publishes it with the built-in Artifact tool to a private page hosted on claude.ai — a link you can open in a browser, iterate on, and choose to share with teammates. Instead of scrolling back through a wall of terminal text, you get a real document: an architecture diagram, a review-findings dashboard, an interactive comparison table, a UI mockup — with the bundled `dataviz` and `artifact-design` skills guiding Claude toward charts and layouts that are actually readable, not just technically rendered.

## Why It Works

Some information is spatial, not linear: layout carries meaning terminal text cannot express, and a stable URL turns the output from a transcript into a deliverable.

## When to Use It

- Architecture summaries and dependency maps after a codebase exploration — diagrams beat prose for structure
- Review or audit reports with many findings — severity, file, and status scan better as a dashboard than as a list
- Comparisons and decision documents — libraries, migration options, benchmark results side by side
- UI mockups during greenfield design — react to a rendered page instead of imagining one from a description
- Anything you plan to share with someone who was not in the session

## When NOT to Use It

- Quick answers and small results — a page for three findings is ceremony, not clarity
- Content that must live in the repo — commit the Markdown; publish the artifact as a *view* of it, not the source of truth
- Sessions that can't publish — headless and CI contexts where nobody clicks links, and setups on Amazon Bedrock, Google Cloud's Agent Platform, or Microsoft Foundry; artifacts require a paid claude.ai plan (Pro, Max, Team, or Enterprise) signed in with `/login` — on Enterprise an Owner must enable them first, and org data policies (CMEK, HIPAA, Zero Data Retention) disable them entirely — and Claude writes a local HTML file instead

## How It Works

### Basic (Beginner)

1. Do the underlying work first — the exploration, the review, the comparison. The artifact is a presentation of results, not a substitute for them.
2. Ask Claude to render it: "Publish this architecture summary as an artifact — one section per service, with a dependency diagram at the top."
3. Claude writes a self-contained HTML file, asks permission the first time it publishes, and prints a `claude.ai` link — your browser opens to the page automatically.
4. Iterate in conversation: "make the timeline horizontal", "collapse the low-severity findings". Claude edits the file and redeploys to the same URL.
5. Share the link when you're happy. Artifacts are private to you by default, and the URL alone grants nothing — viewers must be signed in to claude.ai. On Team and Enterprise plans you can share with specific people or everyone in your organization, never outside it (there is no public option); on Pro and Max plans artifacts stay private to you entirely.

### Composing with Other Approaches (Intermediate)

- **Visual artifacts plus Plan Mode**: After a structured exploration produces an architecture summary, render it as a page with a diagram. The map you built for yourself becomes onboarding material for the next engineer.
- **Visual artifacts plus Built-in Review Skills**: Run `/code-review` or a security audit, then publish the findings as a dashboard grouped by severity — far easier to triage in a team meeting than raw terminal output.
- **Visual artifacts plus Fan-Out Workflows**: When parallel agents each audit one module, have the final step aggregate their reports into a single rendered scorecard instead of concatenated text.

### Advanced Patterns

- **Living status pages**: Keep one artifact current across sessions — asking to "refresh the artifact" re-gathers live state and redeploys to the same URL, so the team's bookmark stays valid while the content tracks reality. A *new* session needs the artifact's URL to update it — without the URL it mints a new page. The official `project-artifact` plugin packages this pattern, remembering the project's sources and published URL between sessions.
- **Interactive deliverables**: Artifacts execute inline JavaScript, so a comparison can have sortable columns and a dependency map can have clickable nodes. Ask for interactivity when the data is bigger than one screen.
- **Design iteration on mockups**: For greenfield UI work, generate two or three visual directions as separate artifacts and A/B them with stakeholders before writing any application code.

## Common Pitfalls

- **Treating the artifact as the source of truth**: The page is a rendered view. If the content matters, commit it to the repo as Markdown and regenerate the artifact from it — otherwise the canonical copy lives outside version control.
- **Linking external assets**: Artifact pages block requests to external hosts (CDN scripts, remote images, web fonts). Everything must be inlined; a page that looks fine in a local preview can silently break when published.
- **Publishing a wall of text**: Rendering prose to HTML does not make it scannable. The value comes from structure — tables, diagrams, severity groupings — so ask for the layout, not just "make it a page."
- **Forgetting the content leaves the machine**: Publishing uploads the page to claude.ai hosting. For proprietary code excerpts or unreleased plans, confirm that is acceptable before rendering.

## Sources

- [Claude Code Artifacts](https://code.claude.com/docs/en/artifacts) — Official docs for creating and sharing artifacts from Claude Code
- [What are Artifacts?](https://support.claude.com/en/articles/9487310-what-are-artifacts-and-how-do-i-use-them) — Artifact hosting, privacy, and sharing model

## Signals

- Setup: —
- Session: Publishes artifacts; asks for rendered/shareable pages
