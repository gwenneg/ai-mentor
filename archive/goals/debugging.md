# Debugging & Root Cause Analysis
*Last verified: 2026-06-27*

## When You're Here

Something is broken and you don't know why. Maybe a test started failing after a merge, maybe users are reporting a bug you can't reproduce locally, or maybe you're staring at a stack trace that points everywhere and nowhere. Debugging is the most common place where AI accelerates your workflow — not by guessing, but by systematically narrowing the search space faster than you can alone.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Complex bug with multiple possible causes | Plan mode | Structured analysis prevents chasing red herrings |
| Test failures after a refactor | Autonomous loops | Set "all tests pass" and let Claude iterate |
| Bug might be caused by your recent changes mixed with others | Worktree isolation | Eliminates interference from unrelated work |
| Error message matches a known library issue | Deep research | Fan-out searches find GitHub issues and patches fast |
| UI renders incorrectly but logic seems right | Browser integration | Visual inspection catches what logs cannot |
| Bug only reproduces in CI, not locally | Headless mode | Run Claude against the same environment CI uses |
| Want to try multiple fix hypotheses without losing progress | Checkpoints & rewind | Mark a checkpoint, try a fix, rewind if wrong |

**Hidden gem:** Hooks — wiring the failing test to run after every single edit is the tightest feedback loop in debugging, and almost nobody thinks of hooks as a debugging tool.

## Approaches (Ranked)

### 1. Plan Mode — Think before you grep
**Level:** Beginner

When facing a complex bug, the instinct is to start adding print statements everywhere. Plan mode forces a structured analysis first: understand the symptoms, form hypotheses, rank them by likelihood, then investigate systematically. This prevents the common trap of spending hours chasing a red herring.

**Try it now:**
> Enter plan mode. I have a race condition: our job scheduler sometimes processes the same job twice. It only happens under load. The relevant code is in `src/scheduler/worker_pool.go` and `src/scheduler/job_queue.go`. Analyze the concurrency model and identify where the double-processing could occur. Don't fix anything yet — just give me ranked hypotheses.

**Why this works:** Bugs survive because developers jump to the first plausible explanation. Structured analysis forces you to enumerate all possibilities before committing to one, dramatically reducing mean time to resolution.

**Pros:**
- Prevents wasted effort on wrong hypotheses
- Creates a debugging record you can share with teammates
- Works for any complexity level

**Cons:**
- Feels slow when you're under pressure — but saves time overall
- Requires discipline to not skip ahead

**Deeper:** See `approaches/plan-mode.md`

---

### 2. Autonomous Loops — Set the goal, let Claude grind
**Level:** Intermediate

When a test is failing and you roughly know the area but not the fix, autonomous loops shine. You define the success condition ("all tests in this module pass") and Claude iterates: read code, form a hypothesis, make a change, run tests, observe results, adjust. This is especially powerful for bugs that require multiple small fixes.

**Try it now:**
> /goal: Make all tests pass in `tests/integration/payment_processing/`. Currently `test_refund_partial_amount` and `test_concurrent_charge` are failing. Investigate the failures, fix the root causes, and verify all tests in the module pass. Don't change test assertions unless they're clearly wrong.

**Why this works:** Debugging often requires tight iteration loops — change, test, observe, repeat. AI excels at this mechanical cycle, freeing you to think about architecture while it handles the grind.

**Pros:**
- Handles multi-step debugging without hand-holding
- Verifies its own fixes by running tests
- Great for "fix and check" workflows

**Cons:**
- Can make superficial fixes that mask deeper issues
- Needs clear success criteria to avoid spinning

**Deeper:** See `approaches/autonomous-loops.md`

---

### 3. Worktree Isolation — Debug in a clean room, not a messy workshop
**Level:** Intermediate

Worktree isolation lets you reproduce a bug in a pristine copy of your repo, free from uncommitted experiments and half-finished features. When debugging, the hardest part is often knowing whether the bug is real or an artifact of your local state. An isolated worktree removes that doubt entirely.

**Try it now:**
> I'm seeing a NullPointerException in `UserService.resolvePermissions()` but only after merging `feature/rbac-overhaul`. Create a worktree from main, cherry-pick the merge commit, and run `mvn test -pl services/user-service` to confirm the failure is from that merge — not my local changes.

**Why this works:** Debugging requires controlled conditions. By eliminating variables (your local state, other branches, uncommitted files), you reduce the search space to only the change that matters. This is the scientific method applied to code.

**Pros:**
- Guarantees a clean reproduction environment
- Your main working tree stays untouched
- Can test multiple hypotheses in parallel worktrees

**Cons:**
- Adds overhead for simple bugs where local state isn't a factor

**Deeper:** See `approaches/worktree-isolation.md`

---

### 4. Deep Research — When the bug isn't in your code
**Level:** Beginner

Sometimes the bug is a known issue in a dependency, a breaking change in a new release, or an undocumented edge case in a platform API. Deep research fans out across GitHub issues, Stack Overflow, release notes, and documentation to find whether someone else has already solved your problem.

**Try it now:**
> /deep-research After upgrading to `pydantic` v2.7, our FastAPI app throws `ValidationError` on fields that worked fine in v2.5. The error is on `Optional[list[str]]` fields with default `None`. Is this a known breaking change? What's the fix?

**Why this works:** Not every bug requires reading source code. Many are solved by knowing the right search terms. Deep research automates the "Google the error message" workflow but does it more thoroughly, cross-referencing multiple sources and verifying answers.

**Pros:**
- Finds known issues in minutes instead of hours
- Cross-references multiple sources for reliability
- Surfaces workarounds and patches you'd miss

**Cons:**
- Only helps when the problem exists outside your codebase
- Results depend on how well-documented the issue is

**Deeper:** See `approaches/deep-research.md`

---

### 5. Browser Integration — See what the user sees
**Level:** Advanced

Some bugs only manifest visually: a CSS layout that breaks at certain viewport sizes, a modal that renders behind an overlay, a chart that displays wrong data despite correct API responses. Browser integration lets Claude control Chrome, navigate to the problem, and inspect the actual rendered output.

**Try it now:**
> Connect to the browser and navigate to `localhost:3000/dashboard`. The revenue chart in the top-right widget is showing negative values even though the API response at `/api/v2/metrics/revenue` returns all positive numbers. Inspect the DOM, check what data the chart component receives as props, and find where the sign gets flipped.

**Why this works:** UI bugs live in the gap between data and rendering. Reading code alone often misses CSS interactions, z-index conflicts, and client-side transformation bugs. Visual debugging closes that gap by observing the actual output.

**Pros:**
- Catches visual bugs that are invisible in code review
- Can screenshot and document the issue automatically
- Tests real browser behavior, not assumptions

**Cons:**
- Requires browser MCP setup
- Slower than code-only debugging
- Not useful for backend or logic bugs

**Deeper:** See `approaches/browser-integration.md`

---

### 6. Hooks — Instant feedback loop for every edit
**Level:** Intermediate

When debugging a failing test, the tightest possible feedback loop is: edit code, run the test, see the result. A PostToolUse hook does this automatically — every time Claude edits a file, the hook runs the failing test immediately. No manual re-runs, no forgetting to check. Claude sees the test output and adjusts its next edit accordingly, converging on the fix faster.

**Try it now:**
> Set up a PostToolUse hook that runs `pytest tests/unit/test_payment_service.py::test_refund_partial_amount -x` after every file edit. Then debug why `test_refund_partial_amount` is failing — the refund amount is calculated as negative when the original charge was partially refunded before.

**Why this works:** Debugging is hypothesis testing. Each edit is an experiment, and the test result is the observation. Hooks eliminate the delay between experiment and observation, letting Claude iterate through hypotheses at maximum speed.

**Pros:**
- Zero-effort test execution after every edit
- Claude sees failures immediately and self-corrects
- Prevents the "forgot to run the test" problem

**Cons:**
- Slow test suites make the hook a bottleneck — scope to a single test or small test file
- Only useful when you have a specific test that reproduces the bug

**Deeper:** See `approaches/hooks-as-workflow.md`

---

### 7. Checkpoints & Rewind — Try fixes fearlessly, rewind what doesn't work
**Level:** Beginner

Debugging often means trying a hypothesis: "maybe the bug is in the error handler." You change the code, test it, and discover that wasn't the problem — now you need to undo everything cleanly. Checkpoints let you mark a point before a speculative fix, try it, and rewind to exactly that state if it doesn't pan out. No git stash juggling, no manual undo.

**Try it now:**
> I think the `NullPointerException` in `UserService.resolvePermissions()` is caused by the role cache returning stale entries. Before I change anything, checkpoint here. Then modify `RoleCacheManager.get()` to bypass the cache and hit the database directly. Run the failing test. If it still fails, rewind — the cache wasn't the problem and I want to try a different hypothesis.

**Why this works:** Debugging is hypothesis testing. Each hypothesis requires a code change to test. Checkpoints make each hypothesis zero-cost to try and zero-cost to abandon, so you can test more hypotheses faster without accumulating half-reverted changes.

**Pros:**
- Zero-risk experimentation — try any fix and rewind instantly
- No need for manual git stash or branch juggling
- Preserves full conversation context across rewinds

**Cons:**
- Only useful for speculative changes — if you're confident in the fix, just commit
- Checkpoints are auto-cleaned after 30 days

**Deeper:** See `approaches/checkpoints-rewind.md`
