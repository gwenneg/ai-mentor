# ci-automation
*Last verified: 2026-07-03*

**Hidden gem:** Custom Agents — project-specific reviewer agents catch the issues only someone who knows your codebase would notice, automatically on every PR.

**Exemplar move:** Write a GitHub Actions step running `claude -p` with `--output-format json` on every PR to main: find newly added TODO/FIXME/HACK comments, output JSON, post via `gh pr comment`.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Headless Mode](../approaches/techniques/headless-mode.md) | Running Claude in GitHub Actions or GitLab CI | CI systems are stdin/stdout/exit-code machines; headless mode adapts Claude to that interface with parseable JSON output |
| 2 | [Fan-Out Workflows](../approaches/techniques/fan-out-workflows.md) | Multi-step pipeline with verification between stages | Decomposes pipelines into independently verifiable stages — parallelism for speed, explicit gates stop bad state propagating |
| 3 | [Subagent Delegation](../approaches/techniques/subagent-delegation.md) | Automated PR review on every push | Specialization beats generalization — a security-focused prompt catches more than a general review pass |
| 4 | [Custom Agents](../approaches/techniques/custom-agents.md) | CI needs domain-specific review of your codebase's risks | Custom agents encode project-specific risks, catching issues only an insider would notice — automatically on every PR |
| 5 | [hookify](../approaches/tools/hookify.md) | A conversation pattern should become a real hook without hand-writing JSON | Automation you describe in plain language and can verify firing beats settings syntax you look up every time |
| 6 | [claude-code-action](../approaches/tools/claude-code-action.md) | A repo wants Claude reviews or automation on every PR | The maintained GitHub Action packages checkout, auth, and prompt execution — no hand-rolled workflow YAML to maintain |
| 7 | [gitlab-ci-integration](../approaches/tools/gitlab-ci-integration.md) | The team is on GitLab and wants Claude in the pipeline | GitLab's maintained job template is the supported path — one .gitlab-ci.yml job, no GitHub-specific tooling |
