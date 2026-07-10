# refactoring
*Last verified: 2026-07-03*

**Hidden gem:** Checkpoints & Rewind — knowing any restructuring is instantly reversible changes which refactors you dare to attempt.

**Exemplar move:** Enter plan mode. Extract auth from ~800-line src/controllers/user_controller.rb into AuthService — which methods move, new interfaces, which callers update, change order keeping tests passing each step.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Plan Mode](../solutions/plan-mode.md) | Complex refactor that could go sideways | Refactoring is dependency management — planning reveals the graph before cutting instead of one crash at a time |
| 2 | [Subagent Delegation](../solutions/subagent-delegation.md) | Changes spanning many files independently | Independent subtasks parallelize safely with roughly linear speedup — where AI refactoring dramatically outperforms manual work |
| 3 | [Autonomous Loops](../solutions/autonomous-loops.md) | Refactor with comprehensive test coverage | A known-good end state with an unknown number of steps — loops grind the mechanical fixes without intervention |
| 4 | [Checkpoints & Rewind](../solutions/checkpoints-rewind.md) | Risky refactor you might need to undo | Fear of irreversibility makes developers conservative — cheap, reliable undo enables bolder, better designs |
