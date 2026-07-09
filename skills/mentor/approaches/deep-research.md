# Deep Research
*Last verified: 2026-07-06*

## What It Is

Deep Research lets you ask a complex question and get back a thoroughly sourced report instead of a single-perspective answer. The tool fans out multiple web searches in parallel, fetches and reads the actual source pages, has independent agents try to disprove each claim, and then synthesizes everything into a cited report with confidence levels. It is the difference between asking a colleague for a quick opinion and commissioning a proper investigation.

## Why It Works

Deep Research applies the same adversarial verification that good journalism and peer review use: multiple independent sources must agree before a claim is treated as established.

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
2. Type: `/deep-research Why does React 19 drop support for defaultProps on function components, and what is the recommended migration path?` and approve the run when Claude Code asks. Deep Research is a bundled dynamic workflow — available on all paid plans (on Pro, turn on Dynamic workflows in `/config`) — and it needs the WebSearch tool available
3. The run continues in the background while your session stays free: Claude fans out multiple parallel web searches with different query variations, then fetches and reads the actual content of the most relevant pages — not just snippet previews. Check progress anytime with `/workflows`
4. Independent verification agents adversarially cross-check each claim, and claims that do not survive verification are filtered out of the report
5. Claude synthesizes a cited report with confidence levels and source URLs for each finding, so you can verify any claim that matters to your decision

### Composing with Other Approaches (Intermediate)

- **Deep Research then Plan Mode**: Research a migration path first, then switch to Plan Mode to design the implementation strategy based on what you learned. The research grounds the plan in reality rather than the AI's training data.
- **Deep Research then Autonomous Loop**: Research the correct API changes for a library upgrade, then run `/goal all tests pass` to let Claude apply the migration mechanically. The research context stays in the session, so the loop has accurate information about what changed and why.
- **Deep Research plus MCP Context**: During an incident, pull the error timeline from your observability stack via MCP Context, then run Deep Research on the suspect dependency — known CVEs, open issues, available patches — to check whether the root cause is a known upstream problem before debugging your own code.

### Advanced Patterns

- **Scoped questions get better results**: "What are the performance implications of enabling React Compiler with existing useMemo calls in a Next.js 15 app?" will produce a more useful report than "Tell me about React Compiler." Specific questions constrain the search space and produce higher-signal results.
- **Asking for comparison tables**: Request output in a specific format: "Compare Drizzle ORM and Prisma for a serverless PostgreSQL deployment. Include cold start impact, migration tooling, and type safety. Format as a comparison table." The synthesis step will structure the report accordingly.
- **Pre-research for unfamiliar domains**: Before starting work in a codebase area you have never touched, run Deep Research on the underlying technology. "How does gRPC streaming handle backpressure in Go?" gives you foundational understanding that makes your code review or implementation significantly better informed.

## Common Pitfalls

- **Underspecified questions**: "What database should I use?" is too broad. The tool will produce a generic survey. Add constraints: "What database should I use for a read-heavy analytics workload with 500M rows, deployed on AWS, team familiar with SQL?" Constraints produce actionable answers.
- **Treating the report as gospel**: Deep Research is thorough, but it is still synthesizing web sources that may themselves be outdated or wrong. Treat confidence levels seriously — "high confidence" with three agreeing sources is more reliable than "moderate confidence" with one source.
- **Using it for speed-sensitive decisions**: Deep Research takes longer than a simple prompt because it runs multiple searches, fetches pages, and verifies claims. If you need a quick answer in 10 seconds, use a regular prompt. Deep Research is for decisions that are worth spending a minute or two to get right.
- **Asking multiple unrelated questions in one prompt**: Each Deep Research session works best with a single focused question. If you need to research both "best ORM for our stack" and "CI provider comparison," run them as separate sessions. A combined prompt dilutes the search queries and produces a shallower report on both topics.

## Sources

- [Dynamic workflows](https://code.claude.com/docs/en/workflows) — Official docs for the bundled /deep-research workflow: availability, approval, and monitoring runs with /workflows

## Signals

- Setup: —
- Session: `/deep-research` in the transcript
