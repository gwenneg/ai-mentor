# Building LLM Features
*Last reviewed: 2026-07-02*

## When You're Here

You're adding AI to your product: a summarizer, a support chatbot, extraction from documents, a natural-language search box. This is ordinary feature work with three unusual properties — the core logic is a prompt, correctness is a distribution rather than a boolean, and cost scales with usage in ways your CFO will notice. The engineering discipline that handles all three is evals: a test suite for model behavior that turns prompt changes from vibes into engineering.

Claude Code is both your assistant and your reference here — the bundled `/claude-api` skill loads current API documentation for your language, so the code it writes against the Claude API uses today's SDK surface instead of a remembered one.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Choosing models, patterns, or architecture | Deep research | Pricing, capabilities, and patterns change fast — research beats recall |
| Designing the feature before building | Plan mode | Prompt strategy, eval plan, fallbacks, and cost model in one review |
| Prompt exists, quality is inconsistent | Autonomous loops | Point the loop at your eval suite and iterate until the pass rate holds |
| Need eval cases and judged outputs at scale | Fan-out workflows | Parallel generation and judge panels build eval sets fast |
| Writing the integration code | Deep research | Run /claude-api first so the SDK usage matches current docs |

**Hidden gem:** Autonomous Loops — point a loop at your eval suite and prompt engineering becomes test-driven development: `/goal 90% of eval cases pass` turns the squishiest part of LLM work into the most mechanical.

## Approaches (Ranked)

### 1. Plan Mode — Design the feature like an engineer, not a demo
**Level:** Beginner

LLM features that ship as demos die in production: no eval plan, no fallback for bad outputs, no cost model. Plan mode forces the full design — prompt strategy, structured-output shape, failure handling, latency budget, eval criteria — before the first API call. The plan's most important output is the definition of "good enough," because without it you'll tune forever.

**Try it now:**
> Enter plan mode. We're adding AI-generated summaries to support tickets in our dashboard. Design the feature: prompt strategy (what context goes in, what shape comes out), how we handle a bad or refused summary in the UI, eval plan (what does a good summary mean for us, measured how?), model tier and cost per 1K tickets, and latency budget. Don't write code — give me the design to review.

**Why this works:** The hard part of LLM features isn't calling the API — it's defining measurable quality and graceful failure, and both are design decisions that code inherits rather than fixes.

**Pros:**
- Forces the eval criteria conversation before quality debates get subjective
- Failure UX designed upfront instead of patched after the first bad output ships
- Produces a cost model before the invoice does

**Cons:**
- Real user inputs will surprise the design — budget for revision after first contact

**Deeper:** See `approaches/plan-mode.md`

---

### 2. Deep Research — Ground decisions in current facts, not recalled ones
**Level:** Beginner

Model lineups, context windows, pricing, and API patterns change monthly — training-data recall is reliably stale, which makes research the difference between designing against today's platform and last year's. Research the model choice and pattern questions; then, when writing the integration, run the bundled `/claude-api` skill so the code targets current SDK reference material for your language.

**Try it now:**
> /deep-research We're building ticket summarization: ~40K tickets/month, needs structured JSON output (summary, sentiment, category), latency under 3s, cost matters. Compare current Claude models for this workload — pricing per ticket at our volume, structured-output support, and prompt-caching economics if we include our category taxonomy in every request. Cite current official docs.

**Why this works:** LLM platform facts have a shelf life of weeks — decisions grounded in current documentation avoid building on a model tier that's been superseded or a price that changed.

**Pros:**
- Current pricing and capability data, cited
- Surfaces platform features you didn't know to ask about (caching, batching)
- `/claude-api` keeps the actual integration code on current SDK surface

**Cons:**
- Research answers go stale too — re-verify before major scaling decisions

**Deeper:** See `approaches/deep-research.md`

---

### 3. Autonomous Loops — Test-driven prompt engineering
**Level:** Intermediate

Once an eval suite exists, prompt tuning is a convergence problem — exactly what loops are for. Set the goal as an eval pass rate, and Claude iterates: run the evals, read the failures, adjust the prompt, re-run. What used to be an afternoon of eyeballing outputs becomes a loop that grinds to the target while you review the diffs.

**Try it now:**
> /goal: at least 18 of the 20 cases in evals/summarize-ticket/ pass. Run them with `npm run eval:summaries`. On failures, read the failing case's expected properties, adjust the prompt in src/ai/prompts/summarize.ts (never the eval cases), and re-run. Show me the prompt diff after each iteration.

**Why this works:** Prompts regress invisibly because output quality is a distribution — an eval suite plus a loop pins the distribution to a number and pushes it in one direction.

**Pros:**
- Converts prompt tuning from taste into iteration against a target
- Every future prompt change re-runs the same suite — regressions surface immediately
- The "never edit the eval cases" rule keeps the loop honest

**Cons:**
- Only as good as the eval suite — 20 weak cases produce confident mediocrity
- Pass-rate optimization can overfit phrasing; keep a held-out set

**Deeper:** See `approaches/autonomous-loops.md`

---

### 4. Fan-Out Workflows — Build eval sets and judge outputs at scale
**Level:** Advanced

The eval suite everyone skips building is embarrassingly parallel to build: fan out agents to draft candidate cases from real (sanitized) production inputs, fan out a judge panel to grade current outputs from multiple lenses, aggregate. What's tedious sequentially — reading fifty tickets and writing expected properties — is an hour of orchestrated parallel work.

**Try it now:**
> Fan out over the 50 sanitized tickets in data/eval-seed/: for each, one agent drafts an eval case (input, expected summary properties, category label) and a second agent independently checks the case is unambiguous and the label defensible. Collect the cases that pass both, flag disagreements for my review, and write the survivors to evals/summarize-ticket/.

**Why this works:** Eval quality comes from volume and independence — parallel drafting with adversarial checking produces in an hour the case coverage that manual curation abandons by case twelve.

**Pros:**
- Turns "we should have evals" into an afternoon instead of a backlog item
- Independent check agents catch ambiguous cases before they pollute the suite
- Judge panels give multi-lens quality reads on current outputs

**Cons:**
- Generated cases need human spot-review — disagreement flags are where to look
- Meaningful token cost; scale the fan-out to the seed data you trust

**Deeper:** See `approaches/fan-out-workflows.md`

---

### 5. Headless Mode — LLM quality gates in CI
**Level:** Intermediate

Prompts live in your repo; their regressions should be caught where code regressions are — CI. A headless run executes the eval suite on every PR that touches prompt files and fails the build when the pass rate drops, making prompt quality a merge requirement instead of a post-deploy discovery.

**Try it now:**
> Add a CI job that runs when files under src/ai/prompts/ change: execute `npm run eval:summaries -- --json`, parse the pass rate, and fail the job below 90%. Post the failing case IDs and their diffs as a PR comment so the author sees exactly which behaviors changed.

**Why this works:** A prompt edit is a behavior change shipped without a type error — CI evals give it the same safety net every other behavior change gets.

**Pros:**
- Prompt regressions blocked at merge time, not found by users
- PR comments make eval failures reviewable like test failures
- Cheap to scope: only runs when prompt files change

**Cons:**
- Eval flakiness will fail builds — assert on outcome properties, not exact strings
- Per-PR token cost; keep the CI subset lean and run the full suite nightly

**Deeper:** See `approaches/headless-mode.md`
