# ci-automation
*Last verified: 2026-07-03*

**Hidden gem:** Custom Agents — project-specific reviewer agents catch the issues only someone who knows your codebase would notice, automatically on every PR.

**Exemplar move:** Write a GitHub Actions step running `claude -p` with `--output-format json` on every PR to main: find newly added TODO/FIXME/HACK comments, output JSON, post via `gh pr comment`.

**Plugins:** `hookify` ✅ hooks from conversation patterns · `gitlab` ☑️ MRs and pipelines · vendor CI: `buildkite`, `mergify`, `teamcity-cli` (all ☑️).

**Built-ins:** `/loop` — recur on a time interval in this session; `/schedule` — recurring cloud runs that survive a closed laptop. Facts and pitfalls per command: `registry/builtin-commands.md`.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Headless Mode](../approaches/headless-mode.md) | Intermediate | Running Claude in GitHub Actions or GitLab CI | CI systems are stdin/stdout/exit-code machines; headless mode adapts Claude to that interface with parseable JSON output |
| 2 | [Fan-Out Workflows](../approaches/fan-out-workflows.md) | Advanced | Multi-step pipeline with verification between stages | Decomposes pipelines into independently verifiable stages — parallelism for speed, explicit gates stop bad state propagating |
| 3 | [Subagent Delegation](../approaches/subagent-delegation.md) | Advanced | Automated PR review on every push | Specialization beats generalization — a security-focused prompt catches more than a general review pass |
| 4 | [Custom Agents](../approaches/custom-agents.md) | Advanced | CI needs domain-specific review of your codebase's risks | Custom agents encode project-specific risks, catching issues only an insider would notice — automatically on every PR |
