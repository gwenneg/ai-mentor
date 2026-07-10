# ci-automation
*Last verified: 2026-07-03*

**Hidden gem:** Custom Agents — project-specific reviewer agents catch the issues only someone who knows your codebase would notice, automatically on every PR.

**Exemplar move:** Write a GitHub Actions step running `claude -p` with `--output-format json` on every PR to main: find newly added TODO/FIXME/HACK comments, output JSON, post via `gh pr comment`.

**Plugins:** `hookify` ✅ hooks from conversation patterns · `gitlab` ☑️ MRs and pipelines · vendor CI: `buildkite`, `mergify`, `teamcity-cli` (all ☑️).

**Integrations:** `claude-code-action` — the maintained GitHub Action; packages checkout, auth, and prompt execution for PR review/automation without hand-rolled YAML · `gitlab-ci-integration` — GitLab's beta `.gitlab-ci.yml` job for teams not on GitHub. Facts and pitfalls: its `solutions/<id>.md` record.

**Built-ins:** `/loop` — recur on a time interval in this session; `/schedule` — recurring cloud runs that survive a closed laptop. Facts and pitfalls per command: its `solutions/<id>.md` record.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Headless Mode](../solutions/headless-mode.md) | Running Claude in GitHub Actions or GitLab CI | CI systems are stdin/stdout/exit-code machines; headless mode adapts Claude to that interface with parseable JSON output |
| 2 | [Fan-Out Workflows](../solutions/fan-out-workflows.md) | Multi-step pipeline with verification between stages | Decomposes pipelines into independently verifiable stages — parallelism for speed, explicit gates stop bad state propagating |
| 3 | [Subagent Delegation](../solutions/subagent-delegation.md) | Automated PR review on every push | Specialization beats generalization — a security-focused prompt catches more than a general review pass |
| 4 | [Custom Agents](../solutions/custom-agents.md) | CI needs domain-specific review of your codebase's risks | Custom agents encode project-specific risks, catching issues only an insider would notice — automatically on every PR |
