# Testing
*Last reviewed: 2026-06-27*

## When You're Here

You need more tests, better tests, or faster tests. Maybe coverage is too low and you're blocked on a merge requirement. Maybe you're adding E2E tests for a critical user flow that's broken in production twice. Or maybe you inherited a codebase with zero tests and need to add a safety net before you can refactor anything. AI is exceptionally good at test generation — but the difference between useful tests and checkbox tests depends on which approach you choose.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Coverage must reach a specific threshold to merge | Autonomous loops | Set the target, let Claude iterate until it's met |
| Adding tests across many modules to hit a coverage target | Fan-out workflows | Parallelize across modules instead of writing sequentially |
| Need to run tests without affecting your working branch | Worktree isolation | Clean environment prevents false passes and flaky failures |
| Testing user-facing flows (login, checkout, forms) | Browser integration | Real browser interaction catches what unit tests miss |
| Test code has type errors that waste test-run cycles | LSP self-correction | Catch type mismatches before waiting for the test runner |
| Inherited a codebase with zero tests | Plan mode | Design a test strategy before writing anything |
| Writing tests to prevent a specific bug from recurring | Autonomous loops | Set "this test passes" as the goal and iterate |

## Approaches (Ranked)

### 1. Autonomous Loops — Set a coverage target, hit it
**Level:** Intermediate | **Tools:** Claude Code

When your goal is quantitative — "get `src/services/` from 47% to 80% coverage" — autonomous loops are the right tool. You set the target, and Claude iterates: analyze uncovered lines, write tests for the most impactful gaps, run coverage, check the number, and repeat. This is the fastest path from "coverage is too low" to "coverage gate passes."

**Try it now:**
> /goal: Increase test coverage for `src/services/order-processing/` from the current 52% to at least 80%. Run `npx jest --coverage --collectCoverageFrom='src/services/order-processing/**/*.ts'` to measure progress. Focus on untested branches and error paths first — happy paths are already covered. Add tests to `tests/services/order-processing/`. All existing tests must continue to pass.

**Why this works:** Coverage improvement is an iterative optimization problem: find the biggest gap, write a test, measure, repeat. AI is perfectly suited to this grind because it doesn't get bored or lose focus on iteration 15.

**Pros:**
- Measurable progress toward a concrete goal
- Automatically prioritizes highest-impact test gaps
- Self-validates by running coverage after each addition

**Cons:**
- Coverage percentage can be gamed with shallow tests — review for quality
- May miss integration-level gaps that line coverage doesn't capture

**Deeper:** See `approaches/autonomous-loops.md`

---

### 2. Fan-out Workflows — Test a whole system in parallel
**Level:** Advanced | **Tools:** Claude Code

When you need to add tests across multiple modules — after a refactor, before a release, or to hit an org-wide coverage mandate — doing them sequentially is painfully slow. Fan-out workflows spawn parallel agents, each responsible for testing one module. Each agent reads the source, writes tests, and runs them independently. You get a full test suite in the time it takes to write one module's tests.

**Try it now:**
> Fan out test generation across these four modules in parallel:
> 1. `src/services/billing/` — focus on edge cases in `calculateProration()` and `applyDiscount()`
> 2. `src/services/inventory/` — cover the stock reservation and release cycle
> 3. `src/api/middleware/` — test auth, rate-limiting, and input validation middleware
> 4. `src/utils/date-helpers.ts` — property-based tests for timezone conversion functions
> Each agent should run its tests and report coverage. Collect results when all complete.

**Why this works:** Test generation is embarrassingly parallel — tests for one module don't depend on tests for another. Fan-out exploits this by doing in parallel what would otherwise be sequential work, cutting wall-clock time by the number of agents.

**Pros:**
- Dramatically faster for multi-module test campaigns
- Each agent has focused context, producing better tests
- Failures in one module don't block progress on others

**Cons:**
- Requires clear module boundaries to parallelize effectively
- Can produce inconsistent test styles across modules without shared conventions
- Higher token usage due to parallel agents

**Deeper:** See `approaches/fan-out-workflows.md`

---

### 3. Worktree Isolation — Test in a clean room
**Level:** Intermediate | **Tools:** Claude Code

Tests that pass on your machine but fail in CI are the worst kind of surprise. Worktree isolation creates a pristine copy of your repo where tests run free from uncommitted changes, stale build artifacts, and half-applied migrations. This is critical when you're generating new tests — you need to know whether a failure is from your test code or your environment.

**Try it now:**
> Create a worktree from the `main` branch. Run the full test suite for `packages/auth/` with `npm test -- --coverage packages/auth/`. Capture the current coverage baseline. Then switch back — I want to compare this against the tests I'm about to add.

**Why this works:** Test reliability depends on environment consistency. A worktree gives you the same isolation CI has, without the 10-minute feedback loop of pushing and waiting for a pipeline.

**Pros:**
- Eliminates "works on my machine" test failures
- Establishes a clean coverage baseline before adding tests
- Can run tests in parallel across multiple worktrees

**Cons:**
- Overhead isn't justified for quick unit test additions
- Large repos may take time to set up worktree dependencies

**Deeper:** See `approaches/worktree-isolation.md`

---

### 4. Browser Integration — Test what users actually do
**Level:** Advanced | **Tools:** Claude Code (with MCP)

Unit tests verify logic. Integration tests verify contracts. But E2E tests verify that users can actually complete their workflows. Browser integration lets Claude drive Chrome through real user flows — filling forms, clicking buttons, waiting for loaders, and asserting on visible outcomes. This catches the class of bugs that survive every other test layer.

**Try it now:**
> Connect to the browser at `localhost:3000`. Test the complete checkout flow: add the item "Wireless Headphones" to cart from `/products`, proceed to checkout, fill in the shipping form with test data, select "Standard Shipping", enter test card `4242 4242 4242 4242`, and submit. Verify the confirmation page shows the correct item, shipping method, and total. Then try the same flow but leave the zip code empty — verify the form shows a validation error and doesn't submit.

**Why this works:** E2E tests exercise the full stack — frontend rendering, API calls, database writes, third-party integrations — in the same sequence a user would. Bugs that hide in the gaps between layers have nowhere to hide when the test uses the same interface the user does.

**Pros:**
- Catches bugs that unit and integration tests cannot
- Tests real browser behavior including JavaScript execution and CSS rendering
- Can generate Playwright/Cypress test scripts from the observed flow

**Cons:**
- Requires browser MCP setup and a running application
- Slower and more brittle than lower-level tests
- Not a substitute for unit tests — use as the top of the testing pyramid

**Deeper:** See `approaches/browser-integration.md`

---

### 5. Hooks — Auto-run tests and protect test fixtures
**Level:** Intermediate | **Tools:** Claude Code

Hooks wire automatic actions to Claude's editing lifecycle. A PostToolUse hook runs your test suite after every file edit, giving Claude instant feedback on whether its new test code works. A PreToolUse hook can block edits to test fixtures, seed data, or snapshot files — preventing Claude from "fixing" a failing test by changing the expected output instead of fixing the code.

**Try it now:**
> Set up two hooks: (1) A PostToolUse hook that runs `npx jest --testPathPattern='tests/services/billing' --no-coverage` after every edit to a file in `src/` or `tests/`. (2) A PreToolUse hook that blocks any edit to files in `tests/fixtures/` unless I explicitly approve. Then generate tests for `src/services/billing/invoice-calculator.ts` focusing on edge cases in `calculateLineItems()`.

**Why this works:** Test generation is an iterative loop: write a test, run it, fix it, repeat. Hooks automate the "run it" step so Claude never writes three tests before discovering the first one doesn't compile. Fixture protection prevents the common failure mode where AI modifies test expectations to match buggy behavior.

**Pros:**
- Instant test feedback after every edit — no manual re-runs
- Fixture protection prevents tests that pass by definition
- Composable — stack a format hook and a test hook together

**Cons:**
- Full test suite hooks are too slow — scope to the module you're working on
- Overly strict PreToolUse blocks can slow down legitimate fixture updates

**Deeper:** See `approaches/hooks-as-workflow.md`

---

### 6. LSP Self-Correction — Fix type errors before running tests
**Level:** Intermediate | **Tools:** OpenCode

When generating test code, type errors are the most common waste of time — you write a test, run it, wait 30 seconds, and discover you passed a `string` where the function expects a `number`. LSP self-correction in OpenCode catches these errors as the code is generated, fixing type mismatches, missing imports, and incorrect method signatures before you ever hit "run." Fewer wasted test cycles means faster coverage gains.

**Try it now:**
> Generate tests for `src/services/pricing/discount-engine.ts`. The `applyStackedDiscounts()` function takes a `PricingContext` object — make sure each test constructs it with the correct shape. Use LSP diagnostics to verify all test files compile cleanly before running `npm test`.

**Why this works:** Test code has the same type constraints as production code, but developers spend less time getting the types right because "it's just a test." LSP enforcement treats test code with the same rigor as production code, eliminating the most common class of test-generation errors.

**Pros:**
- Eliminates the write-run-fix cycle for type errors
- Catches incorrect mock shapes and missing required fields
- Produces tests that compile on the first try

**Cons:**
- Only available in OpenCode, not Claude Code
- Limited to type errors — won't catch logical test mistakes

**Deeper:** See `approaches/lsp-self-correction.md`

---

### 7. Plan Mode — Design your test strategy before writing tests
**Level:** Beginner | **Tools:** Any

When you inherit a codebase with zero tests or need to add a comprehensive test suite to a mature system, jumping straight into writing tests leads to uneven coverage and wasted effort. Plan mode helps you design a test strategy first: what types of tests (unit, integration, E2E), what coverage targets per module, what test infrastructure is needed, and what to prioritize. The plan becomes a roadmap that prevents the common failure mode of writing 50 unit tests for utility functions while critical business logic stays untested.

**Try it now:**
> Enter plan mode. I inherited `services/order-processing/` with zero tests and need to add a test suite before I can safely refactor. The module has 12 files, handles order creation, payment integration (Stripe), inventory reservation, and email notifications. Design a test strategy: what should be unit tested vs integration tested, what mocks/fixtures do I need, what's the priority order for maximum safety with minimum effort, and what coverage target is realistic for the first pass?

**Why this works:** Test strategy is an allocation problem — you have finite time and need to maximize the safety net per hour invested. Planning forces you to prioritize by risk (what breaks production?) rather than by convenience (what's easy to test?), producing a test suite that actually catches bugs instead of just inflating a coverage number.

**Pros:**
- Prevents the "50 tests for utils, zero for business logic" anti-pattern
- Creates a prioritized roadmap so you can stop at any point with maximum value delivered
- Identifies test infrastructure needs upfront (mocks, fixtures, test databases)

**Cons:**
- Adds a planning step — skip it when adding a few focused tests to an already-tested module
- The plan is only as good as your understanding of the codebase's risk areas

**Deeper:** See `approaches/plan-mode.md`
