# testing
*Last verified: 2026-07-03*

**Hidden gem:** Hooks — a PreToolUse hook that blocks edits to fixtures stops the AI from "fixing" a failing test by changing the expected output.

**Exemplar move:** /goal Raise src/services/order-processing/ coverage from 52% to 80% — run `npx jest --coverage --collectCoverageFrom='src/services/order-processing/**/*.ts'`, target untested branches and error paths, keep existing tests passing.

**Plugins:** `playwright` ☑️ browser E2E automation · `fakechat` ☑️ channel-flow testing.

**Built-ins:** `/verify` — confirm the feature works end-to-end, not just green tests; `/run` — launch the app and see the change working; `/goal` — iterate unattended until the suite is green. Facts and pitfalls per command: `registry/builtin-commands.md`.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Autonomous Loops](../approaches/autonomous-loops.md) | Intermediate | Coverage must reach a specific threshold to merge | Coverage improvement is iterative optimization — find the gap, write a test, measure, repeat without getting bored |
| 2 | [Fan-out Workflows](../approaches/fan-out-workflows.md) | Advanced | Adding tests across many modules to hit a coverage target | Test generation is embarrassingly parallel — modules don't depend on each other, so fan-out cuts wall-clock time |
| 3 | [Worktree Isolation](../approaches/worktree-isolation.md) | Intermediate | Need to run tests without affecting your working branch | Test reliability depends on environment consistency — a worktree gives CI's isolation without the pipeline wait |
| 4 | [Hooks](../approaches/hooks-as-workflow.md) | Intermediate | Auto-run tests per edit and protect fixtures from tampering | Hooks automate the run step, and fixture protection stops tests that pass by matching buggy behavior |
