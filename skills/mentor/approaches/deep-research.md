# Deep Research
*Last reviewed: 2026-07-01*

## What It Is

Deep Research lets you ask a complex question and get back a thoroughly sourced report instead of a single-perspective answer. The tool fans out multiple web searches in parallel, fetches and reads the actual source pages, has independent agents try to disprove each claim, and then synthesizes everything into a cited report with confidence levels. It is the difference between asking a colleague for a quick opinion and commissioning a proper investigation.

## Why It Works

Single-source answers are fragile. When you search one page and take its answer at face value, you inherit that source's biases, omissions, and potential inaccuracies. Deep Research applies the same adversarial verification that good journalism and peer review use: multiple independent sources must agree before a claim is treated as established. The fan-out search pattern also surfaces information you would not find in a single query, because different search terms hit different parts of the web. The adversarial step is particularly valuable: rather than just finding supporting evidence, independent agents actively try to disprove each claim, which catches errors that consensus-seeking approaches miss. The result is a report you can actually make decisions from, not a guess you have to double-check yourself.

## When to Use It

- Investigating a library or framework issue where the answer might be spread across GitHub issues, Stack Overflow, blog posts, and documentation
- Comparing frameworks or tools for a technical decision (e.g., "Vite vs. Turbopack for our monorepo build")
- Understanding a breaking change in a dependency upgrade where the migration guide is incomplete
- Due diligence on a new dependency — license, maintenance health, known security issues, community adoption
- Investigating production incidents where the root cause might be a known issue in an upstream library or infrastructure component

## When NOT to Use It

- Questions with simple, well-documented answers — "How do I add a column in PostgreSQL?" does not need multi-source verification
- Questions about your own codebase — Deep Research searches the web, not your repo. Use Plan Mode or standard prompts for internal code questions
- Time-critical decisions where you need an answer in seconds — the research process takes one to several minutes depending on the question's breadth

## How It Works

### Basic (Beginner)

1. First, check whether your question is specific enough. If it is too broad, Claude will ask clarifying questions before starting the research.
2. Start Claude Code and type: `/deep-research Why does React 19 drop support for defaultProps on function components, and what is the recommended migration path?`
3. Claude fans out approximately 5 web searches with different query variations simultaneously
4. It fetches and reads the actual content of the most relevant pages — not just snippet previews
5. Independent verification agents check key claims against other sources
6. Claude synthesizes a report with inline citations, flags any conflicting information, and assigns confidence levels to each finding
7. The final report includes source URLs so you can verify any claim that matters to your decision

### Composing with Other Approaches (Intermediate)

- **Deep Research then Plan Mode**: Research a migration path first, then switch to Plan Mode to design the implementation strategy based on what you learned. The research grounds the plan in reality rather than the AI's training data.
- **Deep Research then Autonomous Loop**: Research the correct API changes for a library upgrade, then run `/goal all tests pass` to let Claude apply the migration mechanically. The research context stays in the session, so the loop has accurate information about what changed and why.
- **Deep Research for incident response**: When a production dependency has a CVE, run Deep Research to understand the vulnerability, its exploitability, available patches, and workarounds — all in one step instead of manually triaging GitHub advisories.

### Advanced Patterns

- **Scoped questions get better results**: "What are the performance implications of enabling React Compiler with existing useMemo calls in a Next.js 15 app?" will produce a more useful report than "Tell me about React Compiler." Specific questions constrain the search space and produce higher-signal results.
- **Asking for comparison tables**: Request output in a specific format: "Compare Drizzle ORM and Prisma for a serverless PostgreSQL deployment. Include cold start impact, migration tooling, and type safety. Format as a comparison table." The synthesis step will structure the report accordingly.
- **Chaining research sessions**: If the first report reveals a subtopic you need to dig into, run a follow-up Deep Research with a narrower question. Each session is independent, so the second query can target exactly the gap the first report exposed.
- **Pre-research for unfamiliar domains**: Before starting work in a codebase area you have never touched, run Deep Research on the underlying technology. "How does gRPC streaming handle backpressure in Go?" gives you foundational understanding that makes your code review or implementation significantly better informed.

## Common Pitfalls

- **Underspecified questions**: "What database should I use?" is too broad. The tool will produce a generic survey. Add constraints: "What database should I use for a read-heavy analytics workload with 500M rows, deployed on AWS, team familiar with SQL?" Constraints produce actionable answers.
- **Treating the report as gospel**: Deep Research is thorough, but it is still synthesizing web sources that may themselves be outdated or wrong. Treat confidence levels seriously — "high confidence" with three agreeing sources is more reliable than "moderate confidence" with one source.
- **Using it for speed-sensitive decisions**: Deep Research takes longer than a simple prompt because it runs multiple searches, fetches pages, and verifies claims. If you need a quick answer in 10 seconds, use a regular prompt. Deep Research is for decisions that are worth spending a minute or two to get right.
- **Asking multiple unrelated questions in one prompt**: Each Deep Research session works best with a single focused question. If you need to research both "best ORM for our stack" and "CI provider comparison," run them as separate sessions. A combined prompt dilutes the search queries and produces a shallower report on both topics.

## Real-World Example

**Scenario**: Choosing a replacement for a deprecated library under time pressure.

Your team is considering replacing Moment.js with a lighter alternative in a large application. You need to make a recommendation to the team lead by end of day.

```
/deep-research Compare date-fns, Day.js, and Luxon as Moment.js replacements.
  Consider: bundle size, tree-shaking support, timezone handling, locale support,
  TypeScript types, active maintenance status, and migration effort from Moment.js.
  Our app uses Moment's timezone features heavily in src/scheduling/.
```

Claude runs parallel searches for each library, fetches their documentation pages, npm download stats, GitHub issue trackers, and recent blog posts comparing them. It finds that one blog post claims Day.js has "full timezone support" but a verification agent discovers this requires the `dayjs/plugin/timezone` plugin which has an open issue about DST edge cases in recurring events. The final report flags this with "moderate confidence" and cites the specific GitHub issue.

The report arrives as a structured comparison with a recommendation section: Day.js for most use cases, but Luxon for your project specifically because of the heavy timezone reliance in `src/scheduling/`. Each claim has a citation linking to the original source. You forward the report to your team lead with minor edits, and the decision is made in the same meeting — saving several hours of manual research across documentation sites, npm pages, and GitHub issue trackers.

## Sources

- [Claude Code Skills](https://code.claude.com/docs/en/skills) — Official docs for skills and slash commands including /deep-research
- [Claude Code Expertise](https://www.anthropic.com/research/claude-code-expertise) — Anthropic research on how Claude Code is used in practice
