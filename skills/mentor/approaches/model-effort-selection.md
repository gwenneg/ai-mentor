# Model & Effort Selection
*Last verified: 2026-07-06*

## What It Is

Model & Effort Selection lets you match the model and the reasoning depth to what each task needs, instead of running everything at one setting. Claude Code ships multiple models (Fable — the Claude 5 flagship — for the hardest, longest-running tasks, Opus for deep reasoning, Sonnet for balanced everyday coding, Haiku for fast, cheap work; availability varies by plan) plus an adjustable effort level that controls how much thinking the model spends per response. You choose per task, per session, or per agent — and switching takes seconds.

## Why It Works

No single model configuration is best at everything; matching the model and effort to the task — and composing them, a strong orchestrator over cheap volume workers — is the right-tool-for-the-job principle applied to reasoning depth.

## When to Use It

- Optimizing costs by running routine tasks on cheaper models and reserving the deep-reasoning models (Fable, Opus) or high effort for complex reasoning
- Long sessions where most turns are mechanical — drop the effort level, then raise it for the hard debugging moment
- Delegating bulk work to subagents: Haiku agents for repetitive checks, Opus agents for security review
- Architecture decisions, tricky debugging, or security analysis where maximum reasoning depth pays for itself

## When NOT to Use It

- When you are just getting started — the defaults are well-chosen; master the workflow before optimizing cost
- When the task is simple enough that model choice does not meaningfully affect the result
- When constant switching costs you more in attention than it saves in tokens

## How It Works

### Basic (Beginner)

1. Check your current model with `/model` — it lists available models and saves your choice as the default for new sessions.
2. Switch based on the task ahead: pick Fable or Opus before a gnarly debugging session, Haiku for a batch of mechanical renames.
3. Adjust reasoning depth with `/effort` (`low`, `medium`, `high`, `xhigh`, `max`, or `auto`) — lower effort responds faster and costs less, higher effort thinks longer on hard problems.
4. On Opus, toggle `/fast` for fast mode: the same model with faster output at a higher cost per token — worth it for interactive back-and-forth where latency matters more than cost.
5. Continue working. The conversation context carries over — the model you switch to sees everything discussed so far, though the first turn after a switch re-reads the history without cached context, so switch at natural task boundaries.

### Composing with Other Approaches (Intermediate)

- **Model selection plus subagent delegation**: Run the main session on a strong model and delegate volume work to Haiku subagents — summarizing files, checking conventions, scanning for patterns. The orchestrator keeps the judgment; the cheap agents do the reading.
- **Effort plus Plan Mode**: Plan a risky migration at `xhigh` effort where reasoning quality matters most, then execute the approved plan at normal effort — the thinking-intensive phase is the plan, not the edits.
- **Model tiers plus custom agents**: Define the same reviewer agent twice in `.claude/agents/` — a `quick-review` on Haiku for draft PRs and a `deep-review` on Opus for release candidates — and pick by stakes.

### Advanced Patterns

- **Per-agent model overrides**: Custom agent definitions accept a `model` field in frontmatter (`opus`, `sonnet`, `haiku`), and ad-hoc subagents can override the model per invocation. A fan-out workflow can run its finder agents cheap and its verifier agents strong.
- **Stage-tiered workflows**: In orchestrated workflows, every agent uses the session's model unless the script routes a stage to a different one — run mechanical extraction stages cheap and the adversarial verification stage strong, where missed reasoning means false positives slip through. Watch `/usage` during long runs and save the expensive configuration for the stages where quality measurably differs.
- **Advisor pairing** (experimental, Anthropic API only, v2.1.98+): instead of switching models, keep a fast main model and set a stronger advisor with `/advisor opus` (or `advisorModel` in settings). Claude consults the advisor at decision points — before committing to an approach, when an error keeps recurring, before declaring a task done — sending it the full conversation. Consultations bill at the advisor model's rates, but pairing a fast main with a strong advisor typically costs less than running the strong model throughout, and toggling `/advisor` doesn't invalidate the prompt cache the way `/model` does.

## Common Pitfalls

- **Switching models too frequently**: Every switch has a cognitive cost — you recalibrate expectations for the new model's style and capability. Switch when the task genuinely demands it, not out of curiosity.
- **Running everything at max effort**: High effort on trivial tasks buys latency, not quality. Reserve `xhigh` and `max` for problems where reasoning depth is the bottleneck.
- **Using a fast model for subtle work**: Haiku is excellent at well-specified mechanical tasks and weaker at multi-file reasoning and nuanced review. If a cheap agent's output needs constant correction, the time cost exceeds the token savings.
- **Forgetting the setting persists**: `/model` saves your choice as the default for new sessions. If you dropped to a fast model for a batch job, remember to switch back before your next deep session.

## Sources

- [Claude Code Model Configuration](https://code.claude.com/docs/en/model-config) — Official docs covering /model, available models, and effort levels
- [Claude Code Fast Mode](https://code.claude.com/docs/en/fast-mode) — Official docs for fast mode on Opus
- [Escalate hard decisions with the advisor tool](https://code.claude.com/docs/en/advisor) — Official docs for /advisor model pairing and billing
