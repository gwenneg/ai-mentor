# testing
*Last verified: 2026-07-12*

**Hidden gem:** Hooks — a PreToolUse hook that blocks edits to fixtures stops the AI from "fixing" a failing test by changing the expected output.

**Exemplar move:** /goal Raise src/services/order-processing/ coverage from 52% to 80% — run `npx jest --coverage --collectCoverageFrom='src/services/order-processing/**/*.ts'`, target untested branches and error paths, keep existing tests passing.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Autonomous Loops](../approaches/techniques/autonomous-loops.md) | Coverage must reach a specific threshold to merge | Coverage improvement is iterative optimization — find the gap, write a test, measure, repeat without getting bored |
| 2 | [Fan-out Workflows](../approaches/techniques/fan-out-workflows.md) | Adding tests across many modules to hit a coverage target | Test generation is embarrassingly parallel — modules don't depend on each other, so fan-out cuts wall-clock time |
| 3 | [Worktree Isolation](../approaches/techniques/worktree-isolation.md) | Need to run tests without affecting your working branch | Test reliability depends on environment consistency — a worktree gives CI's isolation without the pipeline wait |
| 4 | [Hooks](../approaches/techniques/hooks-as-workflow.md) | Auto-run tests per edit and protect fixtures from tampering | Hooks automate the run step, and fixture protection stops tests that pass by matching buggy behavior |
| 5 | [Built-In Review Skills](../approaches/techniques/built-in-review-skills.md) | Tests are green but nobody has watched the feature actually work | Green tests prove the cases you wrote; `/verify` drives the real flow end-to-end and `/run` shows the change working — behavior a diff reader can't see |
