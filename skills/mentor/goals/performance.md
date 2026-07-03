# Performance Optimization
*Last verified: 2026-06-28*

## When You're Here

Your code works — it's just slow. Maybe an API endpoint takes three seconds to respond, a page load triggers a Lighthouse score in the red, a database query scans every row instead of hitting an index, or your JavaScript bundle has ballooned past a megabyte. Performance optimization is distinct from debugging (the behavior is correct) and from refactoring (the goal is speed, not structure). The two disciplines overlap, but the success criteria are different: you're done when the numbers hit the target, not when the code looks cleaner.

The biggest trap in performance work is premature optimization — spending hours micro-optimizing a function that accounts for 2% of latency while ignoring the unindexed query that accounts for 80%. AI workflows help you avoid this by profiling first, setting measurable targets, and iterating against real benchmarks. Every approach below is built around the same principle: measure, change, measure again.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Not sure where the bottleneck actually is | Plan mode | Profile and analyze before changing anything |
| Clear target like "under 200ms" or "below 500KB" | Autonomous loops | Measurable goals are perfect for iterate-until-met |
| Need framework-specific optimization techniques | Deep research | Find proven patterns instead of reinventing them |
| Frontend performance — Lighthouse, Core Web Vitals | Browser integration | Measures real rendering performance in a real browser |
| Want to see performance impact of every edit | Hooks | Auto-run benchmarks so regressions are caught immediately |

**Hidden gem:** Hooks — benchmarking after every single edit turns optimization from guesswork into a measured experiment per change.

## Approaches (Ranked)

### 1. Plan Mode — Profile first, optimize second
**Level:** Beginner

The fastest way to waste a day is to optimize the wrong thing. Plan mode forces you to profile, identify the actual bottlenecks, and rank them by impact before writing a single line of optimized code. Most performance problems aren't where developers think they are — the slow endpoint isn't slow because of your business logic, it's slow because of three N+1 queries buried in a serializer.

**Try it now:**
> Enter plan mode. Our `/api/v2/orders` endpoint averages 2.4 seconds. I've added basic timing logs but can't tell what's slow. The handler is in `src/controllers/orders_controller.py`, it calls `OrderService.list_with_details()` in `src/services/order_service.py`, which queries through `OrderRepository` in `src/repositories/order_repo.py`. Analyze the data flow, identify likely bottlenecks (N+1 queries, missing indexes, unnecessary serialization), and rank them by probable impact. Don't change code yet.

**Why this works:** Performance budgets are finite — you get the most gain from fixing the biggest bottleneck first. Without profiling, developers optimize by intuition, which is notoriously unreliable. Structured analysis ensures effort goes where the data says it should.

**Pros:**
- Prevents wasted effort on low-impact optimizations
- Builds a prioritized optimization roadmap
- Works even without profiling tools — Claude can reason about algorithmic complexity from code

**Cons:**
- Requires discipline to resist "just fixing" an obvious inefficiency before finishing analysis
- Static analysis can miss runtime-specific bottlenecks like cache miss rates

**Deeper:** See `approaches/plan-mode.md`

---

### 2. Autonomous Loops — Set a target, iterate until met
**Level:** Intermediate

Performance optimization has something most coding tasks lack: a clear, measurable finish line. "Make this endpoint respond under 200ms" or "reduce bundle size below 500KB" are perfect autonomous loop goals. Claude can make a change, run the benchmark, check the number, and decide whether to keep going or try a different approach — all without you watching.

**Try it now:**
> /goal: The endpoint at `src/api/routes/search.py::search_products` currently averages 1.8 seconds. Get it under 300ms. The benchmark command is `python -m pytest tests/benchmarks/test_search_perf.py -v`. Start by analyzing the query in `src/repositories/product_repo.py::search()`, then optimize iteratively — try indexing changes, query restructuring, or adding a cache layer. Run the benchmark after each change.

**Why this works:** Performance tuning is inherently iterative — each optimization reveals the next bottleneck. Autonomous loops handle this cycle naturally, and the measurable target prevents both premature stopping ("good enough") and over-optimization ("but I can shave off 2 more milliseconds").

**Pros:**
- Measurable success criteria prevent ambiguity
- Handles multi-step optimization chains without hand-holding
- Each iteration is verified against the actual benchmark

**Cons:**
- Can over-optimize for the benchmark while missing real-world patterns
- Without profiling context, may try brute-force approaches before surgical ones

**Deeper:** See `approaches/autonomous-loops.md`

---

### 3. Deep Research — Find proven optimization patterns
**Level:** Beginner

Before inventing a custom caching layer or hand-rolling a query optimizer, check whether your framework already has a built-in solution. Deep research finds battle-tested optimization patterns — Django's `select_related` and `prefetch_related`, React's `useMemo` and code splitting, PostgreSQL's partial indexes and query plan hints. The best optimization is often one someone else already wrote.

**Try it now:**
> /deep-research We're using Next.js 14 with App Router and our Lighthouse performance score is 42. The main issues are large JavaScript bundles (1.2MB total), slow Largest Contentful Paint (4.1s), and high Cumulative Layout Shift (0.28). What are the proven Next.js 14 optimization techniques for each of these metrics? Include specific configuration examples and any relevant next.config.js settings.

**Why this works:** Frameworks evolve faster than developers can track. The optimization techniques that worked in React 17 may be anti-patterns in React 19. Research ensures you're applying current best practices, not outdated ones from a two-year-old blog post.

**Pros:**
- Surfaces framework-native solutions you might not know exist
- Cross-references multiple sources to verify advice is current
- Finds specific configuration and code examples

**Cons:**
- Only helps when the optimization is a known pattern, not a novel bottleneck
- Research results need validation against your specific setup

**Deeper:** See `approaches/deep-research.md`

---

### 4. Browser Integration — Measure real user experience
**Level:** Advanced

Lighthouse scores, Core Web Vitals, rendering waterfalls, network timing — these are things you can only measure in a real browser. Browser integration lets Claude run Lighthouse audits, inspect render performance, analyze network requests, and measure layout shifts directly. This is essential for frontend performance work where the bottleneck is in rendering, not in your server code.

**Try it now:**
> Connect to the browser and navigate to `localhost:3000/products`. Run a Lighthouse performance audit and analyze the results. Focus on the three worst-scoring metrics. Then inspect the network waterfall — are there render-blocking resources or unnecessary sequential requests? Give me a prioritized list of fixes with estimated impact on the overall score.

**Why this works:** Frontend performance is multi-dimensional — a fast API doesn't matter if the browser spends two seconds parsing a massive JavaScript bundle or re-rendering a poorly memoized component tree. Browser tools measure what the user actually experiences, not what your server logs say.

**Pros:**
- Measures real browser rendering, not theoretical performance
- Lighthouse provides standardized, actionable scoring
- Catches issues invisible to server-side profiling — layout shifts, render blocking, unused JavaScript

**Cons:**
- Requires browser MCP setup and a running local dev server
- Lighthouse scores can be noisy between runs — average multiple measurements
- Only applicable to frontend performance

**Deeper:** See `approaches/browser-integration.md`

---

### 5. Hooks — Auto-run benchmarks after every edit
**Level:** Intermediate

When optimizing a hot path, you need to know whether each change made things faster or slower. A PostToolUse hook runs your benchmark suite after every file edit, so Claude sees the performance impact immediately. This creates the tightest possible feedback loop: edit, measure, adjust. No forgotten benchmark runs, no guessing whether the last change helped.

**Try it now:**
> Set up a PostToolUse hook that runs `go test -bench=BenchmarkSearchQuery -benchmem -count=3 ./internal/search/` after every file edit. Then optimize `internal/search/query_builder.go` — the `BuildFilteredQuery` function is allocating excessively. Make changes one at a time so we can see the benchmark impact of each optimization.

**Why this works:** Performance optimization without measurement is guesswork. Hooks ensure every single change is measured automatically, turning optimization into a data-driven process. Claude can see "that change reduced allocations by 40% but didn't improve latency" and adjust its strategy accordingly.

**Pros:**
- Every edit is immediately benchmarked — no forgotten measurements
- Claude self-corrects when a change hurts performance
- Creates a clear record of what helped and what didn't

**Cons:**
- Benchmark suites that take more than a few seconds create painful edit latency
- Micro-benchmarks may not reflect real-world performance patterns

**Deeper:** See `approaches/hooks-as-workflow.md`

---

### 6. Built-In Review Skills — Catch performance regressions in review
**Level:** Beginner

The cheapest performance optimization is preventing regressions. `/code-review` with a performance focus catches N+1 queries, unnecessary allocations, missing database indexes, unoptimized loops, and synchronous calls that should be async — before the code reaches production. This is especially valuable after a refactor, when structural changes can accidentally introduce performance problems.

**Try it now:**
> /code-review Focus on performance: look for N+1 query patterns, missing database indexes on filtered or joined columns, unnecessary object allocations in hot paths, synchronous I/O that could be async, and any O(n^2) or worse algorithms operating on potentially large collections.

**Why this works:** Performance regressions are insidious because they rarely cause test failures — the code works correctly, just slower. Code review with an explicit performance lens catches patterns that human reviewers often miss because they're focused on correctness and readability.

**Pros:**
- Catches regressions before they reach production
- Low effort — a single command on your existing diff
- Finds structural patterns (N+1, missing index) not just micro-optimizations

**Cons:**
- Static analysis only — cannot detect runtime performance characteristics
- May flag patterns that are acceptable at your current scale

**Deeper:** See `approaches/built-in-review-skills.md`
