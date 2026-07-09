# building-skills-plugins
*Last verified: 2026-07-03*

**Hidden gem:** Hooks as Workflow — the `hookify` official plugin mines your own conversation history for repeated patterns and turns them into hooks: your plugin's best features are often already hiding in what you keep asking Claude to do.

**Exemplar move:** Create .claude/skills/release-notes/SKILL.md: find latest git tag, categorize commits since by conventional-commit prefix into Features/Fixes/Breaking, prepend dated entry to CHANGELOG.md; run it and iterate.

**Plugins:** `plugin-dev` ✅ guided plugin building · `skill-creator` ☑️ skill authoring and evals.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Custom Skills](../approaches/custom-skills.md) | none | A workflow you repeat and want to package | Skills fail on unclear instructions, and clarity is discovered by iteration — the standalone stage keeps iteration cost near zero |
| 2 | [Custom Plugins](../approaches/custom-plugins.md) | some | Skill works; time to share or distribute | Distribution is where team knowledge compounds — one engineer's workflow becomes everyone's default, shipped deliberately via versioning |
| 3 | [Hooks as Workflow](../approaches/hooks-as-workflow.md) | some | The behavior should be automatic, not invoked | Skills require someone to invoke them; hooks fire regardless — the workflow works even for teammates who skip the README |
