# debugging
*Last verified: 2026-07-12*

**Hidden gem:** Hooks — wiring the failing test to run after every single edit is the tightest feedback loop in debugging, and almost nobody thinks of hooks as a debugging tool.

**Exemplar move:** Enter plan mode. Job scheduler double-processes jobs under load — analyze the concurrency model in src/scheduler/worker_pool.go and src/scheduler/job_queue.go, give ranked hypotheses, don't fix anything yet.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Plan Mode](../approaches/techniques/plan-mode.md) | Complex bug with multiple possible causes | Bugs survive because developers jump to the first plausible explanation; enumerating all possibilities cuts time to resolution |
| 2 | [Autonomous Loops](../approaches/techniques/autonomous-loops.md) | Test failures after a refactor | Debugging is tight change-test-observe iteration; AI handles the mechanical grind while you think about architecture |
| 3 | [Worktree Isolation](../approaches/techniques/worktree-isolation.md) | Bug might be your recent changes mixed with others | Eliminating local-state variables reduces the search space to only the change that matters — the scientific method for code |
| 4 | [Hooks](../approaches/techniques/hooks-as-workflow.md) | A failing test should run after every single edit | Each edit is an experiment and the test result its observation; hooks remove the delay between them |
| 5 | [Browser Integration](../approaches/techniques/browser-integration.md) | The bug only reproduces in the browser | Reading the console while driving the page correlates the runtime error to its source line — and verifies the fix in the same loop |
