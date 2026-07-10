# building-skills-plugins
*Last verified: 2026-07-03*

**Hidden gem:** Hooks as Workflow — the `hookify` official plugin mines your own conversation history for repeated patterns and turns them into hooks: your plugin's best features are often already hiding in what you keep asking Claude to do.

**Exemplar move:** Create .claude/skills/release-notes/SKILL.md: find latest git tag, categorize commits since by conventional-commit prefix into Features/Fixes/Breaking, prepend dated entry to CHANGELOG.md; run it and iterate.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Custom Skills](../approaches/custom-skills.md) | A workflow you repeat and want to package | Skills fail on unclear instructions, and clarity is discovered by iteration — the standalone stage keeps iteration cost near zero |
| 2 | [Custom Plugins](../approaches/custom-plugins.md) | Skill works; time to share or distribute | Distribution is where team knowledge compounds — one engineer's workflow becomes everyone's default, shipped deliberately via versioning |
| 3 | [Hooks as Workflow](../approaches/hooks-as-workflow.md) | The behavior should be automatic, not invoked | Skills require someone to invoke them; hooks fire regardless — the workflow works even for teammates who skip the README |
| 4 | [hookify](../approaches/hookify.md) | Turning a repeated instruction into an enforced hook | The gap between 'I keep telling Claude X' and 'X is enforced' is one generated hook |
| 5 | [plugin-dev](../approaches/plugin-dev.md) | Building a plugin with scaffolding, validation, and review in one flow | The validator catches manifest mistakes before `claude plugin validate` does, with a reviewer that critiques honestly |
