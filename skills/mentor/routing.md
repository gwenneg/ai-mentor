# Goal Routing Table
*Last verified: 2026-07-03*

The mentor's per-goal judgment layer: for every goal in SKILL.md's Phase 1 classification table, the ranked approaches with their curated trigger ("Best when") and rationale ("Why it fits"), plus the hidden gem. Rankings, gems, and rationales are editorial judgment — they change when we change our minds, not when the product changes; verifiable product claims live in the approach files. "Setup" is what the approach requires (the skill renders it as no/some/involved setup), not a statement about the user.

Each section carries an **Exemplar move** — a fictional prompt showing the *shape* of a well-grounded move for that goal (named files, embedded values, real commands). Never show it verbatim: rewrite it against the actual repo, or keep it portable when the problem targets another repo.

Condensed 2026-07-03 from the former `goals/` essays; the full prose is recoverable in `archive/goals/` and git history.

## accessibility

**Hidden gem:** Hooks — running the a11y scanner after every component edit catches regressions the moment they're introduced, not at audit time.

**Exemplar move:** Connect to the browser at localhost:3000/settings, record tab order; verify the Delete Account modal in src/components/Settings/DeleteAccountModal.tsx traps focus and Escape returns focus to the trigger.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Browser Integration](approaches/browser-integration.md) | Advanced | Verifying tab order, focus traps, or screen reader behavior | Tests the actual user experience — tab order, focus management, live regions — where the hardest a11y bugs hide |
| 2 | [Deep Research](approaches/deep-research.md) | Beginner | Building a new component, unsure which ARIA pattern applies | ARIA standards are precise but sprawling; one wrong attribute can make a component opaque to assistive technology |
| 3 | [Autonomous Loops](approaches/autonomous-loops.md) | Intermediate | Scanner reported dozens of violations to fix | Scanner violations are objective and verifiable — ideal exit conditions for fix-scan-repeat loops with no judgment calls |
| 4 | [Built-In Review Skills](approaches/built-in-review-skills.md) | Beginner | Reviewing a PR for accessibility regressions | Most violations are introduced one PR at a time; catching them at review is vastly cheaper than at audit |
| 5 | [Hooks](approaches/hooks-as-workflow.md) | Intermediate | Want every component edit automatically checked for a11y | A11y fixes have unintended side effects; instant scanner feedback catches regressions before they compound |

## api-design

**Hidden gem:** Built-In Review Skills — pointing `/code-review` at only the API surface catches breaking changes and naming drift that no test suite will.

**Exemplar move:** Enter plan mode. Design a REST API for multi-tenant project management (projects, tasks, comments, members): URL structure, schemas, error codes, versioning, tenant scoping, pagination for 10K+ tasks. No code.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](approaches/plan-mode.md) | Beginner | Starting a new API or major version from scratch | An API is a promise to consumers; planning the contract first finds design mistakes while fixing them is free |
| 2 | [Deep Research](approaches/deep-research.md) | Beginner | Unsure about conventions like pagination, versioning, error format | Studying proven APIs gives rationale for your choices and avoids pitfalls others already documented |
| 3 | [MCP Context](approaches/mcp-context.md) | Intermediate | Requirements scattered across specs, tickets, and existing APIs | Pulling the existing API surface alongside requirements keeps new endpoints consistent — same naming, errors, pagination |
| 4 | [Built-in Review Skills](approaches/built-in-review-skills.md) | Beginner | Extending an existing API and need to stay consistent | Misspelled fields and inconsistent errors break no test but frustrate every consumer; review catches surface issues type checkers miss |
| 5 | [Custom Skills](approaches/custom-skills.md) | Intermediate | Your team creates new endpoints frequently | API consistency comes from repetition, and skills automate repetition — every endpoint scaffolded from the same template |

## building-agents

**Hidden gem:** Custom Agents — prototyping your agent as a ten-line `.claude/agents/` file answers most of the design questions (tools, model, instructions) before you write a single line of SDK code.

**Exemplar move:** Enter plan mode. Design a support-inbox triage agent — classify severity, draft responses, escalate billing/data-loss to humans: what tools, forbidden actions, escalation boundaries, minimal first version worth shipping?

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](approaches/plan-mode.md) | Beginner | Unsure what the agent should even do | An agent is a policy wrapped around a model — explicit designs fail predictably, emergent ones fail creatively |
| 2 | [Custom Agents](approaches/custom-agents.md) | Beginner | Want a working prototype today | Editing markdown and re-running is the cheapest iteration loop — converge on instructions and tool surface in an afternoon |
| 3 | [Official Plugins](approaches/official-plugins.md) | Intermediate | Ready to build a standalone product | Agent infrastructure is undifferentiated heavy lifting; a production-tested engine puts your effort into the agent's judgment |
| 4 | [Subagent Delegation](approaches/subagent-delegation.md) | Advanced | Agent needs multiple cooperating workers | Multi-agent failure modes show up in a simulation you can run today, not after building the message bus |
| 5 | [Permissions & Safe Autonomy](approaches/safe-autonomy.md) | Intermediate | Need to constrain what the agent can do | Capability is easy and trust hard to win back — agents earn adoption through provable boundaries, not demos |

## building-mcp-integrations

**Hidden gem:** Headless Mode — a shell script of `claude -p` calls against your server is an MCP integration test suite: repeatable, CI-friendly, and it tests the thing that actually matters (can a model use your tools correctly?).

**Exemplar move:** Enter plan mode. Design MCP tool surface for incident management (search, read detail, comment, change status, page on-call): tools vs. left out, descriptions, read-only ones, damage potential.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](approaches/plan-mode.md) | Beginner | Not sure what the server should expose | Models select tools by description — a well-designed five-tool surface outperforms a twenty-tool REST dump |
| 2 | [Official Plugins](approaches/official-plugins.md) | Intermediate | Ready to implement | MCP's protocol surface is undifferentiated work; a guided workflow encodes current conventions so effort goes into tool semantics |
| 3 | [MCP Context](approaches/mcp-context.md) | Beginner | Haven't used MCP as a consumer yet | Interface intuition comes from the consumer side — a week consuming MCP beats a month producing it blind |
| 4 | [Deep Research](approaches/deep-research.md) | Beginner | Unfamiliar with the protocol's details | Protocols punish improvisation — the design constraints you don't know become the rewrite you do six weeks in |
| 5 | [Headless Mode](approaches/headless-mode.md) | Intermediate | Server built, needs regression testing | Whether a model can use your tools correctly is the real acceptance criterion — headless runs make it executable |

## building-skills-plugins

**Hidden gem:** Hooks as Workflow — the `hookify` official plugin mines your own conversation history for repeated patterns and turns them into hooks: your plugin's best features are often already hiding in what you keep asking Claude to do.

**Exemplar move:** Create .claude/skills/release-notes/SKILL.md: find latest git tag, categorize commits since by conventional-commit prefix into Features/Fixes/Breaking, prepend dated entry to CHANGELOG.md; run it and iterate.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Custom Skills](approaches/custom-skills.md) | Beginner | A workflow you repeat and want to package | Skills fail on unclear instructions, and clarity is discovered by iteration — the standalone stage keeps iteration cost near zero |
| 2 | [Custom Plugins](approaches/custom-plugins.md) | Intermediate | Skill works; time to share or distribute | Distribution is where team knowledge compounds — one engineer's workflow becomes everyone's default, shipped deliberately via versioning |
| 3 | [Hooks as Workflow](approaches/hooks-as-workflow.md) | Intermediate | The behavior should be automatic, not invoked | Skills require someone to invoke them; hooks fire regardless — the workflow works even for teammates who skip the README |
| 4 | [Custom Agents](approaches/custom-agents.md) | Advanced | Plugin should ship a specialist's judgment, not just procedures | Judgment scales worse than procedure — an agent definition runs one specialist's standards on every team, every time |
| 5 | [Autonomous Loops](approaches/autonomous-loops.md) | Advanced | Want to measure whether the skill actually works | Skill instructions are prompts and prompts regress silently — an eval set makes quality a number that goes up |

## ci-automation

**Hidden gem:** Custom Agents — project-specific reviewer agents catch the issues only someone who knows your codebase would notice, automatically on every PR.

**Exemplar move:** Write a GitHub Actions step running `claude -p` with `--output-format json` on every PR to main: find newly added TODO/FIXME/HACK comments, output JSON, post via `gh pr comment`.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Headless Mode](approaches/headless-mode.md) | Intermediate | Running Claude in GitHub Actions or GitLab CI | CI systems are stdin/stdout/exit-code machines; headless mode adapts Claude to that interface with parseable JSON output |
| 2 | [Fan-Out Workflows](approaches/fan-out-workflows.md) | Advanced | Multi-step pipeline with verification between stages | Decomposes pipelines into independently verifiable stages — parallelism for speed, explicit gates stop bad state propagating |
| 3 | [Subagent Delegation](approaches/subagent-delegation.md) | Advanced | Automated PR review on every push | Specialization beats generalization — a security-focused prompt catches more than a general review pass |
| 4 | [Hooks](approaches/hooks-as-workflow.md) | Intermediate | Want to enforce rules automatically on every edit | Automation eliminates the gap between intention and execution — the right thing happens every time, no discipline needed |
| 5 | [MCP Context](approaches/mcp-context.md) | Intermediate | Pipeline needs data from GitHub, monitoring, or wikis | Automated AI's value scales with context breadth — cross-system correlation catches what code-only review never would |
| 6 | [Custom Skills](approaches/custom-skills.md) | Advanced | CI triggers the same Claude workflows repeatedly | Skills give CI workflows a versioned, reviewable definition so behavior is consistent for humans and pipelines alike |
| 7 | [Custom Agents](approaches/custom-agents.md) | Advanced | CI needs domain-specific review of your codebase's risks | Custom agents encode project-specific risks, catching issues only an insider would notice — automatically on every PR |
| 8 | [Channels](approaches/channels.md) | Advanced | Want CI failures pushed into your live session | Reacting to CI needs the context that produced the change; pushing events there beats cold-starting an agent |

## code-review

**Hidden gem:** MCP Context — reviewing against the ticket's acceptance criteria instead of just the diff catches the most expensive bugs: code that works but solves the wrong problem.

**Exemplar move:** Run `/code-review --effort high` on the branch diff — focus on src/services/billing/ discount stacking edge cases and Stripe webhook failure-mode error handling.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Built-In Review Skills](approaches/built-in-review-skills.md) | Beginner | Quick review of a focused diff or your own pre-PR code | Most review value is systematic checking — codified reviewer instincts applied consistently to every line |
| 2 | [Subagent Delegation](approaches/subagent-delegation.md) | Advanced | Large PR touching security, perf, and correctness | Attention is finite — specialized agents keep focus per concern while parallelism keeps wall-clock time low |
| 3 | [Fan-Out Workflows](approaches/fan-out-workflows.md) | Advanced | Critical change needing adversarial verification | Automated review's biggest problem is false positives; adversarial verification mimics human pushback for higher-signal results |
| 4 | [MCP Context](approaches/mcp-context.md) | Intermediate | PR implements a design doc or addresses an issue | The most expensive bugs are specification bugs — grounding review in requirements catches code solving the wrong problem |
| 5 | [Custom Agents](approaches/custom-agents.md) | Advanced | Same project-specific review concerns recur on every PR | Review quality depends on knowing the codebase's patterns and risks; agents encode that knowledge once, applied consistently |
| 6 | [Official Plugins](approaches/official-plugins.md) | Intermediate | Want structured review without building your own pipeline | Plugins encode community best practices — structured output, severity ranking, false-positive filtering — covering most needs immediately |
| 7 | [Cloud Sessions](approaches/cloud-sessions.md) | Intermediate | Deep review off your machine, or PRs that fix themselves | Most post-review churn is mechanical; a watching cloud agent handles it so human reviewers focus on design |

## code-understanding

**Hidden gem:** LSP Self-Correction — compiler-backed go-to-definition and find-references beat text search for tracing how components actually connect.

**Exemplar move:** Enter plan mode. Trace one payment end-to-end through services/payment-gateway/ — entry points, validation, fraud checks, processor integration, persistence, retry logic; produce an architecture summary with dependency diagram.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](approaches/plan-mode.md) | Beginner | Onboarding to a new team or project | Understanding is about building the right mental model — learn the shape of the system first, details later |
| 2 | [MCP Context](approaches/mcp-context.md) | Intermediate | Architecture docs exist but are scattered or outdated | Code says what, docs say why; MCP brings both into one conversation and flags where docs drifted |
| 3 | [Deep Research](approaches/deep-research.md) | Beginner | Codebase uses unfamiliar frameworks or patterns | Frameworks encode decisions into conventions — learning the grammar lets you read the codebase's sentences fluently |
| 4 | [LSP Self-Correction](approaches/lsp-self-correction.md) | Intermediate | Tracing how components connect across a large codebase | LSP gives compiler-precise answers to who-calls-what, where text search misses indirect calls and aliased imports |
| 5 | [Project Memory & Context Docs](approaches/project-memory.md) | Beginner | Want the map you built to persist across sessions | Exploration output is knowledge — storing it where every session reads converts one-off investigation into permanent capability |
| 6 | [Session & Context Management](approaches/session-context-management.md) | Beginner | Long exploration is saturating the context window | Exploration quality degrades silently as context fills; curating the window keeps reasoning over conclusions, not noise |
| 7 | [Visual Artifacts](approaches/visual-artifacts.md) | Beginner | Sharing what you learned with the team | Codebase understanding is a graph — a rendered diagram at a stable URL turns private investigation into team knowledge |

## debugging

**Hidden gem:** Hooks — wiring the failing test to run after every single edit is the tightest feedback loop in debugging, and almost nobody thinks of hooks as a debugging tool.

**Exemplar move:** Enter plan mode. Job scheduler double-processes jobs under load — analyze the concurrency model in src/scheduler/worker_pool.go and src/scheduler/job_queue.go, give ranked hypotheses, don't fix anything yet.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](approaches/plan-mode.md) | Beginner | Complex bug with multiple possible causes | Bugs survive because developers jump to the first plausible explanation; enumerating all possibilities cuts time to resolution |
| 2 | [Autonomous Loops](approaches/autonomous-loops.md) | Intermediate | Test failures after a refactor | Debugging is tight change-test-observe iteration; AI handles the mechanical grind while you think about architecture |
| 3 | [Worktree Isolation](approaches/worktree-isolation.md) | Intermediate | Bug might be your recent changes mixed with others | Eliminating local-state variables reduces the search space to only the change that matters — the scientific method for code |
| 4 | [Deep Research](approaches/deep-research.md) | Beginner | Error message matches a known library issue | Many bugs are solved by knowing the right search terms — a thorough, cross-referenced "Google the error" workflow |
| 5 | [Browser Integration](approaches/browser-integration.md) | Advanced | UI renders incorrectly but logic seems right | UI bugs live in the gap between data and rendering; visual debugging observes the actual output |
| 6 | [Hooks](approaches/hooks-as-workflow.md) | Intermediate | A failing test should run after every single edit | Each edit is an experiment and the test result its observation; hooks remove the delay between them |
| 7 | [Checkpoints & Rewind](approaches/checkpoints-rewind.md) | Beginner | Want to try multiple fix hypotheses without losing progress | Checkpoints make each hypothesis zero-cost to try and abandon, so you test more hypotheses faster |
| 8 | [Headless Mode](approaches/headless-mode.md) | Intermediate | Failure only reproduces in CI, never locally | Environment-dependent bugs must be debugged in the environment that produces them — run the investigation inside the CI job itself |

## dependency-management

**Hidden gem:** Worktree Isolation — trying the upgrade in a throwaway copy gives you a real damage report before you commit to anything.

**Exemplar move:** /deep-research Compare zod vs joi for our Node.js API (src/validators/): TypeScript integration, bundle size, performance, maintenance activity, breaking-change history, CVEs — must stay maintained 3+ years.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Deep Research](approaches/deep-research.md) | Beginner | Evaluating a library you haven't used before | Adoption is a long-term bet; research automates the due diligence — maintenance health, CVEs, licenses — most developers skip |
| 2 | [Plan Mode](approaches/plan-mode.md) | Beginner | Understanding what depends on a package before removing it | Dependency changes are graph operations — mapping the full graph first catches transitive breakage before it happens |
| 3 | [MCP Context](approaches/mcp-context.md) | Intermediate | Your org has approved dependency lists or security policies | The fastest decision is one someone already made — internal sources prevent duplicating or contradicting prior evaluations |
| 4 | [Autonomous Loops](approaches/autonomous-loops.md) | Intermediate | Bumping a dependency version with cascading breakage | Upgrades are convergent — each fix brings you closer to green, and the test suite measures the remaining distance |
| 5 | [Worktree Isolation](approaches/worktree-isolation.md) | Intermediate | Want to test a major upgrade without risking your branch | A disposable environment gives a realistic damage assessment instead of guessing from the changelog |
| 6 | [Scheduled & Recurring Agents](approaches/scheduled-agents.md) | Intermediate | Dependency triage is important but never urgent | Recurring maintenance survives only when it stops depending on human initiative — a schedule converts "someone should" into "it happened" |

## devops

**Hidden gem:** Worktree Isolation — rendering `terraform plan` or `helm template` in a disposable copy lets you evaluate risky infra changes with zero blast radius.

**Exemplar move:** Enter plan mode. Split monolithic main.tf (VPC, three subnets, security groups, RDS, ECS, ALB) into per-service modules — map every cross-resource reference, give a safe extraction order.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](approaches/plan-mode.md) | Beginner | About to change a resource that other services depend on | Infra dependency graphs hide in string references; making them explicit prevents destroy-and-recreate surprises mid-apply |
| 2 | [Deep Research](approaches/deep-research.md) | Beginner | Unfamiliar with a cloud provider's limits or pricing model | Provisioning mistakes are expensive to fix; research catches limits and pricing traps before they become incidents or bills |
| 3 | [Worktree Isolation](approaches/worktree-isolation.md) | Intermediate | Want to test Terraform or Helm changes without touching your branch | Throwaway safety matters most for infrastructure — worst case is deleting a worktree, not filing an incident report |
| 4 | [Autonomous Loops](approaches/autonomous-loops.md) | Intermediate | Debugging a CrashLoopBackOff or Terraform apply error | Infra debugging is convergent — each fix eliminates one error and the tooling says exactly what's wrong next |
| 5 | [MCP Context](approaches/mcp-context.md) | Intermediate | Need to cross-reference live cloud state with proposed changes | Infrastructure-as-code is half the picture; live state lets Claude reason about real utilization, costs, and conflicts |
| 6 | [Permissions & Safe Autonomy](approaches/safe-autonomy.md) | Intermediate | Iterating on infra where an AI mistake costs the most | Guardrails beat vigilance — a deny rule binds no matter what the model attempts, unlike CLAUDE.md guidance |

## documentation

**Hidden gem:** Custom Skills — a `/gen-api-doc` command makes regenerating docs cheaper than letting them drift.

**Exemplar move:** Read src/api/routes/payments.ts, src/api/middleware/auth.ts, docs/openapi.yaml, and the Notion design doc linked in docs/DESIGN_DECISIONS.md; generate payments API reference with auth, schemas, error codes, rate limits.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [MCP Context](approaches/mcp-context.md) | Intermediate | Existing docs, specs, and decisions are scattered | Documentation quality is proportional to context quality — MCP bridges where knowledge lives and where it needs to go |
| 2 | [Deep Research](approaches/deep-research.md) | Beginner | Writing docs for an unfamiliar domain or standard | Documentation is a communication design problem; proven patterns give a structural blueprint so effort goes to content |
| 3 | [Plan Mode](approaches/plan-mode.md) | Beginner | Large doc set with multiple sections and audiences | Documentation is information architecture — an outline gives every fact one home and writers clear coverage without overlap |
| 4 | [Custom Skills](approaches/custom-skills.md) | Intermediate | You regenerate the same type of documentation repeatedly | Same-structure docs are a solved problem once encoded — consistent, complete output from source code every time |
| 5 | [Visual Artifacts](approaches/visual-artifacts.md) | Beginner | Docs need diagrams, dashboards, or sharing beyond the repo | Diagrams carry relationships linear Markdown flattens, and a stable URL turns docs into a link people actually open |

## greenfield

**Hidden gem:** Custom Skills — encoding your conventions as a scaffold command before writing feature #2 pays off for every feature after it.

**Exemplar move:** Enter plan mode. Design a notification service (email, Slack, in-app; consumes Kafka events) — module structure, data models, API surface, delivery retries, per-channel rate limiting. Architecture only, no code.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](approaches/plan-mode.md) | Beginner | Starting a new service or module from scratch | In greenfield work code is cheap but architecture is expensive; planning makes structural decisions explicit, not accidental |
| 2 | [Deep Research](approaches/deep-research.md) | Beginner | Need to choose between libraries or frameworks | Technology choices compound — twenty minutes of research avoids the wrong-library realization after 10,000 lines of integration |
| 3 | [Autonomous Loops](approaches/autonomous-loops.md) | Intermediate | Feature is well-defined and you want speed | Loops excel at mechanical build-test-fix iteration — clear, testable criteria tell Claude when it's done |
| 4 | [Browser Integration](approaches/browser-integration.md) | Advanced | Building a user-facing feature with UI | Users interact with pixels, not code — browser feedback catches CSS conflicts and interaction bugs only a real browser shows |
| 5 | [MCP Context](approaches/mcp-context.md) | Intermediate | Requirements live in Jira, Linear, or Notion | The biggest greenfield risk is building the wrong thing; pulling requirements from the source eliminates the telephone game |
| 6 | [Custom Skills](approaches/custom-skills.md) | Advanced | Your team scaffolds the same structure for every new module | The first 20% of every module is patterned boilerplate — skills encode it once and skip to the interesting work |
| 7 | [Official Plugins](approaches/official-plugins.md) | Intermediate | A plugin already covers your feature workflow | Plugin authors already solved the orchestration problems — standing on their work skips the design-test-refine cycle |

## incident-response

**Hidden gem:** Deep Research — two minutes checking whether it's a known upstream outage can save an hour of debugging code that was never the problem.

**Exemplar move:** Connect to Grafana and Datadog MCP servers. payment-service returns 503s on ~12% of requests — pull error timeline, latest deploy, top service:payment-service log errors; correlate spike with deploy timestamp.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [MCP Context](approaches/mcp-context.md) | Intermediate | Error spike correlating with a recent deploy | Incidents are data problems — MCP makes Claude a unified query layer across your observability stack, collapsing data-gathering |
| 2 | [Deep Research](approaches/deep-research.md) | Beginner | Outage symptoms match something you've seen in a dependency | Many incidents originate outside your codebase; checking known issues first saves hours debugging code that wasn't the problem |
| 3 | [Plan Mode](approaches/plan-mode.md) | Beginner | Service down and you don't know where to start | In incidents a wrong action costs more than a delayed one — structured triage prevents fixing the wrong component |
| 4 | [Headless Mode](approaches/headless-mode.md) | Intermediate | Need to query production state without manual SSH sessions | Incident analyses repeat across services and time windows; headless makes them scriptable, repeatable, and systematic |
| 5 | [Checkpoints & Rewind](approaches/checkpoints-rewind.md) | Beginner | Applying a hotfix and worried about making it worse | Pressure creates mistakes — checkpoints make hotfix attempts zero-risk, back to a clean state in seconds |
| 6 | [Subagent Delegation](approaches/subagent-delegation.md) | Advanced | Multiple plausible root causes to investigate at once | Hypotheses are independent — parallel investigation cuts diagnosis time by the number tested simultaneously |
| 7 | [Channels](approaches/channels.md) | Advanced | Alerts should land in the session with context already loaded | Triage speed is dominated by context assembly; channels deliver the event where that assembly already happened |

## llm-features

**Hidden gem:** Autonomous Loops — point a loop at your eval suite and prompt engineering becomes test-driven development: `/goal 90% of eval cases pass` turns the squishiest part of LLM work into the most mechanical.

**Exemplar move:** Enter plan mode. Design AI-generated ticket summaries: prompt strategy, UI handling of bad/refused summaries, eval plan, model tier and cost per 1K tickets, latency budget. No code.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](approaches/plan-mode.md) | Beginner | Designing the feature before building | The hard part isn't calling the API — measurable quality and graceful failure are design decisions code inherits |
| 2 | [Deep Research](approaches/deep-research.md) | Beginner | Choosing models, patterns, or architecture | LLM platform facts have a shelf life of weeks — current documentation beats stale training-data recall |
| 3 | [Autonomous Loops](approaches/autonomous-loops.md) | Intermediate | Prompt exists, quality is inconsistent | Prompts regress invisibly because quality is a distribution; evals plus a loop pin it to a number and push it up |
| 4 | [Fan-Out Workflows](approaches/fan-out-workflows.md) | Advanced | Need eval cases and judged outputs at scale | Eval quality comes from volume and independence — parallel drafting with adversarial checking builds coverage manual curation abandons |
| 5 | [Headless Mode](approaches/headless-mode.md) | Intermediate | Prompt regressions should be caught in CI at merge time | A prompt edit is a behavior change shipped without a type error — CI evals give it the same safety net |

## migration

**Hidden gem:** Worktree Isolation — running the upgrade in a disposable copy first turns "should we migrate?" from changelog speculation into a concrete damage report.

**Exemplar move:** Enter plan mode. Migrate React Router v5 to v6: 34 routes in src/routes/, useHistory in ~20 components, guards in src/auth/ProtectedRoute.tsx — map changes, categorize, order to stay functional.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](approaches/plan-mode.md) | Beginner | Large migration touching dozens of files | Migrations are dependency graphs, not task lists — planning reveals the structure so you traverse it in the right order |
| 2 | [Subagent Delegation](approaches/subagent-delegation.md) | Advanced | Migration across independent modules or services | Migrations are embarrassingly parallel when modules don't share state — subagents do in minutes what takes sequential hours |
| 3 | [Autonomous Loops](approaches/autonomous-loops.md) | Intermediate | Upgrade where "it compiles and tests pass" is the goal | Mechanical migrations have a clear convergence criterion — AI doesn't get bored or lose focus on iteration 47 |
| 4 | [Worktree Isolation](approaches/worktree-isolation.md) | Intermediate | Risky upgrade you want to test without polluting main | A throwaway environment changes how you approach risk — bolder strategies, faster discovery, nothing at stake |
| 5 | [Deep Research](approaches/deep-research.md) | Beginner | Unfamiliar framework with unclear migration path | Maintainer guides plus community posts on undocumented edge cases synthesize into a single breaking-changes briefing |
| 6 | [Checkpoints & Rewind](approaches/checkpoints-rewind.md) | Beginner | Multi-step migration where any step might break things | Checkpoints turn an irreversible process reversible — a wrong step costs seconds, not hours of untangling |
| 7 | [Custom Agents](approaches/custom-agents.md) | Advanced | Your team runs migrations regularly with project-specific rules | Migration rules are tribal knowledge — a custom agent captures them so every migration meets the same standards |

## onboarding

**Hidden gem:** Custom Skills — a `/setup-dev` skill is executable documentation: it can't silently go stale the way a wiki page does.

**Exemplar move:** Pull the Confluence "Engineering Onboarding" guide and #platform-team pinned messages, cross-reference with README.md and docker-compose.yml — flag outdated steps, produce one consolidated setup guide.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [MCP Context](approaches/mcp-context.md) | Intermediate | Team knowledge scattered across Confluence, Slack, Notion, and READMEs | The onboarding bottleneck is finding information, not understanding it — MCP brings every source into one queryable context |
| 2 | [Plan Mode](approaches/plan-mode.md) | Beginner | Need a systematic understanding of the codebase and architecture | Replicates how experts onboard — learn the architecture's shape first, fill in module details as tickets demand |
| 3 | [Deep Research](approaches/deep-research.md) | Beginner | Codebase uses frameworks or patterns you've never worked with | Patterns encode decisions into structure; learning the stack's vocabulary lets you read the codebase fluently from day one |
| 4 | [Custom Skills](approaches/custom-skills.md) | Intermediate | Local dev setup takes a half-day of manual steps | A skill is executable documentation — when a step breaks you fix the skill and every future onboarder benefits |
| 5 | [Custom Plugins](approaches/custom-plugins.md) | Intermediate | Team maintains shared onboarding workflows and conventions | Plugin skills are discoverable and executable on the spot — team best practices become your default workflow from day one |
| 6 | [Project Memory & Context Docs](approaches/project-memory.md) | Beginner | Making what you learn while ramping up inheritable | Onboarding knowledge has a one-session half-life unless persisted — CLAUDE.md compounds it for you and future teammates |

## performance

**Hidden gem:** Hooks — benchmarking after every single edit turns optimization from guesswork into a measured experiment per change.

**Exemplar move:** Enter plan mode. /api/v2/orders averages 2.4s — trace src/controllers/orders_controller.py through OrderService.list_with_details() and OrderRepository, rank bottlenecks (N+1 queries, missing indexes), no code changes yet.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](approaches/plan-mode.md) | Beginner | Not sure where the bottleneck actually is | Performance budgets are finite and intuition unreliable — analysis sends effort where the data says it should go |
| 2 | [Autonomous Loops](approaches/autonomous-loops.md) | Intermediate | Clear target like "under 200ms" or "below 500KB" | Tuning is inherently iterative, and a measurable target prevents both premature stopping and over-optimization |
| 3 | [Deep Research](approaches/deep-research.md) | Beginner | Need framework-specific optimization techniques | Frameworks evolve faster than developers track — research applies current best practices, not two-year-old blog advice |
| 4 | [Browser Integration](approaches/browser-integration.md) | Advanced | Frontend performance — Lighthouse, Core Web Vitals | Browser tools measure what the user actually experiences — parsing, rendering, layout shifts — not what server logs say |
| 5 | [Hooks](approaches/hooks-as-workflow.md) | Intermediate | Want to see performance impact of every edit | Optimization without measurement is guesswork; auto-benchmarks make every change a data point Claude adjusts to |
| 6 | [Built-In Review Skills](approaches/built-in-review-skills.md) | Beginner | Preventing performance regressions from reaching production | Regressions rarely fail tests — a performance-lens review catches patterns humans miss while focused on correctness |

## refactoring

**Hidden gem:** Checkpoints & Rewind — knowing any restructuring is instantly reversible changes which refactors you dare to attempt.

**Exemplar move:** Enter plan mode. Extract auth from ~800-line src/controllers/user_controller.rb into AuthService — which methods move, new interfaces, which callers update, change order keeping tests passing each step.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](approaches/plan-mode.md) | Beginner | Complex refactor that could go sideways | Refactoring is dependency management — planning reveals the graph before cutting instead of one crash at a time |
| 2 | [Subagent Delegation](approaches/subagent-delegation.md) | Advanced | Changes spanning many files independently | Independent subtasks parallelize safely with roughly linear speedup — where AI refactoring dramatically outperforms manual work |
| 3 | [Autonomous Loops](approaches/autonomous-loops.md) | Intermediate | Refactor with comprehensive test coverage | A known-good end state with an unknown number of steps — loops grind the mechanical fixes without intervention |
| 4 | [Worktree Isolation](approaches/worktree-isolation.md) | Intermediate | Multiple agents modifying overlapping areas | Isolation is a prerequisite for safe parallelism — worktrees give filesystem-level separation with zero clone overhead |
| 5 | [Built-In Review Skills](approaches/built-in-review-skills.md) | Beginner | Post-refactor cleanup and polish | Refactoring shifts attention to structure; a fresh cleanup pass catches the local improvements you stopped noticing |
| 6 | [Checkpoints & Rewind](approaches/checkpoints-rewind.md) | Beginner | Risky refactor you might need to undo | Fear of irreversibility makes developers conservative — cheap, reliable undo enables bolder, better designs |
| 7 | [Hooks](approaches/hooks-as-workflow.md) | Intermediate | Keeping dozens of edited files formatted and lint-clean | Auto-formatting every edit keeps the final diff showing only the refactor, not formatting noise |

## release-management

**Hidden gem:** Headless Mode — a headless pre-tag validation step catches the version mismatches and forgotten migrations humans miss on Friday afternoons.

**Exemplar move:** Create .claude/skills/release-notes.md: /release-notes finds the latest git tag, categorizes commits by conventional prefixes, prepends a versioned changelog entry to CHANGELOG.md, adds migration notes for BREAKING CHANGE footers.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Custom Skills](approaches/custom-skills.md) | Intermediate | Need categorized release notes from commit history | Release prep is procedural — same steps, same order — so a skill removes forgotten or misordered steps |
| 2 | [Headless Mode](approaches/headless-mode.md) | Intermediate | Validating release readiness in CI before cutting a tag | Humans forget pre-release checks under deadline pressure; headless validation catches blockers when you're rushing to ship |
| 3 | [Plan Mode](approaches/plan-mode.md) | Beginner | Complex release with multiple services, migrations, and rollback steps | Complex releases fail on ordering and improvised rollbacks — planning the full sequence turns coordination into a checklist |
| 4 | [Built-In Review Skills](approaches/built-in-review-skills.md) | Beginner | Final quality check on the release diff | PRs get reviewed individually but their interactions don't — a full-diff review catches cross-PR inconsistencies |
| 5 | [Hooks](approaches/hooks-as-workflow.md) | Intermediate | Want automatic validation on every edit during release prep | Release prep values consistency over speed — hooks catch version drift and format errors the moment they happen |

## research

**Hidden gem:** Plan Mode — defining evaluation criteria before gathering evidence is the only reliable defense against confirmation bias, and nobody thinks of plan mode for research.

**Exemplar move:** /deep-research Compare RabbitMQ, Kafka, NATS JetStream for ~50K events/sec, exactly-once payment delivery, Go and Python services — throughput, delivery guarantees, operational complexity, client maturity, AWS managed hosting.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Deep Research](approaches/deep-research.md) | Beginner | Comparing libraries, frameworks, or SaaS tools | Automates the 30-browser-tab workflow into a cited, cross-referenced report you can present to your team |
| 2 | [Plan Mode](approaches/plan-mode.md) | Beginner | Need a structured evaluation framework first | Research without a framework is browsing — criteria first, evidence second, decision last is how staff engineers decide |
| 3 | [Browser Integration](approaches/browser-integration.md) | Advanced | Evaluating tools with interactive demos or UIs | Tools are chosen on experience as much as capability — browser testing evaluates the experience, not the feature list |
| 4 | [MCP Context](approaches/mcp-context.md) | Intermediate | Past decisions or prior art exist inside the company | Institutional knowledge is invisible to web search — internal sources save days and prevent recommending what was already rejected |
| 5 | [Model & Effort Selection](approaches/model-effort-selection.md) | Advanced | Different phases need different reasoning depth | Reasoning depth is a budget — allocating it per phase gets better conclusions and a cheaper session at once |

## security

**Hidden gem:** Hooks — a PreToolUse guard on auth configs and crypto files prevents the accidental security regressions that no scanner catches.

**Exemplar move:** Run /security-review on the current branch — special attention to auth middleware in src/middleware/auth.ts and raw database queries in src/services/; security audit next week.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Built-In Review Skills](approaches/built-in-review-skills.md) | Beginner | Quick security scan before a release or audit | Vulnerabilities follow well-known patterns — encoded checks applied exhaustively to every changed line beat manual review |
| 2 | [Subagent Delegation](approaches/subagent-delegation.md) | Advanced | Large codebase with multiple vulnerability classes to check | Security auditing is multi-dimensional — one concern per agent gives deeper analysis without attention dilution |
| 3 | [Deep Research](approaches/deep-research.md) | Beginner | New CVE announced for a dependency in your stack | Hardening without context is guesswork — affected versions and exploitation prerequisites let you patch what matters |
| 4 | [Custom Agents](approaches/custom-agents.md) | Advanced | Recurring security patterns unique to your project | Generic scanners produce noise; an agent that knows your middleware, ORM, and PII fields yields high-signal findings |
| 5 | [MCP Context](approaches/mcp-context.md) | Intermediate | Need to check against compliance frameworks or threat models | Auditors ask about specific controls — reviewing against them produces findings that map directly to audit requirements |
| 6 | [Hooks](approaches/hooks-as-workflow.md) | Intermediate | Protect security-critical files from accidental modification | Most security regressions are accidental — a speed bump forces conscious acknowledgment before touching critical code |
| 7 | [Scheduled & Recurring Agents](approaches/scheduled-agents.md) | Intermediate | Security review that never skips a week | Attackers don't wait for your audit cadence — scheduled review makes security a standing property of the codebase |

## tech-debt

**Hidden gem:** Custom Agents — encoding your team's own deprecated patterns as a detector turns tribal knowledge into a tracked migration metric.

**Exemplar move:** Spawn four parallel agents auditing src/: duplication clusters, deprecated APIs (TODO/FIXME/@deprecated), test-coverage gaps vs tests/, functions over cyclomatic complexity 10 — consolidate into one prioritized report.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Subagent Delegation](approaches/subagent-delegation.md) | Advanced | Multi-dimensional audit across the whole codebase | Each debt dimension needs a different scanning strategy — specialized agents apply the right one without conflating them |
| 2 | [Built-In Review Skills](approaches/built-in-review-skills.md) | Beginner | Quick sense of the worst offenders, no setup | The most impactful debt is often the most visible — built-in skills catch the low-hanging fruit systematically |
| 3 | [Plan Mode](approaches/plan-mode.md) | Beginner | You have findings and need to prioritize them | Prioritization weighs bug frequency, developer drag, and roadmap fit — plan mode forces explicit trade-off reasoning |
| 4 | [Fan-Out Workflows](approaches/fan-out-workflows.md) | Advanced | Very large codebase with many modules | Debt distribution is uneven — fan-out reveals it quantitatively so cleanup targets the modules that matter most |
| 5 | [Custom Agents](approaches/custom-agents.md) | Advanced | Recurring patterns specific to your project | The most impactful debt is patterns your team decided to leave; custom agents enforce those decisions and track progress |
| 6 | [Background Agents](approaches/background-agents.md) | Intermediate | Audit takes an hour but your attention shouldn't | Audits are attention-light but time-heavy — backgrounding turns them from a thing you do into a thing you review |

## testing

**Hidden gem:** Hooks — a PreToolUse hook that blocks edits to fixtures stops the AI from "fixing" a failing test by changing the expected output.

**Exemplar move:** /goal: Raise src/services/order-processing/ coverage from 52% to 80% — run `npx jest --coverage --collectCoverageFrom='src/services/order-processing/**/*.ts'`, target untested branches and error paths, keep existing tests passing.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Autonomous Loops](approaches/autonomous-loops.md) | Intermediate | Coverage must reach a specific threshold to merge | Coverage improvement is iterative optimization — find the gap, write a test, measure, repeat without getting bored |
| 2 | [Fan-out Workflows](approaches/fan-out-workflows.md) | Advanced | Adding tests across many modules to hit a coverage target | Test generation is embarrassingly parallel — modules don't depend on each other, so fan-out cuts wall-clock time |
| 3 | [Worktree Isolation](approaches/worktree-isolation.md) | Intermediate | Need to run tests without affecting your working branch | Test reliability depends on environment consistency — a worktree gives CI's isolation without the pipeline wait |
| 4 | [Browser Integration](approaches/browser-integration.md) | Advanced | Testing user-facing flows like login, checkout, forms | E2E through the user's own interface leaves nowhere to hide for bugs living in the gaps between layers |
| 5 | [Hooks](approaches/hooks-as-workflow.md) | Intermediate | Auto-run tests per edit and protect fixtures from tampering | Hooks automate the run step, and fixture protection stops tests that pass by matching buggy behavior |
| 6 | [LSP Self-Correction](approaches/lsp-self-correction.md) | Intermediate | Test code has type errors that waste test-run cycles | LSP treats test code with production rigor, eliminating the most common class of test-generation errors before running |
| 7 | [Plan Mode](approaches/plan-mode.md) | Beginner | Inherited a codebase with zero tests | Test strategy is an allocation problem — prioritize by production risk, not by testing convenience |
| 8 | [Built-in Verify & Review Skills](approaches/built-in-review-skills.md) | Beginner | Tests pass but you're not sure the feature actually works | Tests encode expectations; verification observes behavior — driving the real flow closes the gap where mocked bugs live |
