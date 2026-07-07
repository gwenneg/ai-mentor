# refactoring
*Last verified: 2026-07-03*

**Hidden gem:** Checkpoints & Rewind — knowing any restructuring is instantly reversible changes which refactors you dare to attempt.

**Exemplar move:** Enter plan mode. Extract auth from ~800-line src/controllers/user_controller.rb into AuthService — which methods move, new interfaces, which callers update, change order keeping tests passing each step.

**Plugins:** none recommended — the ⚠️-flagged `code-simplifier` duplicates the built-in /simplify; lead with the built-in.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | Beginner | Complex refactor that could go sideways | Refactoring is dependency management — planning reveals the graph before cutting instead of one crash at a time |
| 2 | [Subagent Delegation](../approaches/subagent-delegation.md) | Advanced | Changes spanning many files independently | Independent subtasks parallelize safely with roughly linear speedup — where AI refactoring dramatically outperforms manual work |
| 3 | [Autonomous Loops](../approaches/autonomous-loops.md) | Intermediate | Refactor with comprehensive test coverage | A known-good end state with an unknown number of steps — loops grind the mechanical fixes without intervention |
| 4 | [Checkpoints & Rewind](../approaches/checkpoints-rewind.md) | Beginner | Risky refactor you might need to undo | Fear of irreversibility makes developers conservative — cheap, reliable undo enables bolder, better designs |
