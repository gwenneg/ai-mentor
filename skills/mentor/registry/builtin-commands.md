# Capability Registry — built-in commands
*Last verified: 2026-07-06*

The built-in-command slice of the capability registry: every recommendable built-in command as one flat record. These are first-class mentoring units — they get profile rows under their `id`, direct routing via the `**Built-ins:**` line in `routing/<goal>.md`, and a place in the ignorance map. Goal membership is a slug list, **never a numeric fit score**; the model orders candidates at answer time from live evidence. `best_when` is one auditable sentence. Facts here are verified against the official docs on the `Last verified` date — syntax is quoted exactly; paraphrase it, never reconstruct it from training data.

Record fields: `id` (profile row key), `kind`, `goals` (slugs from the classification table), `best_when`, `setup` (none | some | involved), `composes_with` (approach basenames or registry ids), `exemplar` (copy-ready line), `session_signal` (what marks it adopted from conversation evidence), `pitfalls`, `source`.

## code-review

id: code-review
kind: builtin-command
goals: code-review, tech-debt
best_when: a diff or PR needs a structured correctness pass before humans read it
setup: none
composes_with: verify, built-in-review-skills, subagent-delegation
exemplar: /code-review high src/api/
session_signal: user ran /code-review or /review in this conversation
pitfalls:
- Effort levels trade coverage for confidence: low/medium give fewer, higher-confidence findings; high through max widen coverage and may include uncertain ones. Pick by stakes, not habit.
- `--fix` applies findings to the working tree and `--comment` posts them as inline PR comments — say which one the situation wants.
source: https://code.claude.com/docs/en/commands

## security-review

id: security-review
kind: builtin-command
goals: security, code-review
best_when: pending changes touch auth, input handling, or anything a CVE could grow from
setup: none
composes_with: code-review, built-in-review-skills
exemplar: /security-review
session_signal: user ran /security-review in this conversation
pitfalls:
- It reviews the pending changes on the current branch — run it before merging, not after.
- A clean pass is one lens, not a security audit; critical paths still need a human pass.
source: https://code.claude.com/docs/en/commands

## simplify

id: simplify
kind: builtin-command
goals: refactoring, tech-debt
best_when: working code needs reuse, simplification, and efficiency cleanup — quality only, not bug hunting
setup: none
composes_with: code-review, built-in-review-skills
exemplar: /simplify
session_signal: user ran /simplify in this conversation
pitfalls:
- Run /code-review first: simplification can consolidate duplicated code that still contains a logic error, leaving the bug in one tidier place.
source: https://code.claude.com/docs/en/commands

## verify

id: verify
kind: builtin-command
goals: testing, code-review
best_when: tests and typecheck pass but nobody has watched the change actually work end-to-end
setup: none
composes_with: code-review, run, built-in-review-skills
exemplar: /verify
session_signal: user ran /verify in this conversation
pitfalls:
- Don't run it on diffs with no runtime surface (docs-only, test-only) — there is no behavior to observe.
source: https://code.claude.com/docs/en/commands

## run

id: run
kind: builtin-command
goals: testing, debugging
best_when: you want to see the change working in the real app, not infer it from green tests
setup: none
composes_with: verify, browser-integration
exemplar: /run
session_signal: user ran /run in this conversation
pitfalls:
- It launches and drives the project's app; for a library with no runnable surface it has nothing to drive.
source: https://code.claude.com/docs/en/commands

## goal

id: goal
kind: builtin-command
goals: debugging, testing, migration
best_when: the finish line is machine-verifiable (tests green, zero type errors) and human supervision of each iteration adds nothing
setup: none
composes_with: autonomous-loops, safe-autonomy, worktree-isolation
exemplar: /goal all tests in tests/payment/ pass
session_signal: user ran /goal in this conversation
pitfalls:
- The evaluator optimizes for the stated condition — a "tests pass" goal can be satisfied by deleting the test. Review the diff (`git diff --stat` first).
- Vague conditions ("improve performance") give the evaluator nothing checkable; check progress with `/goal`, stop with `/goal clear`.
source: https://code.claude.com/docs/en/goal

## loop

id: loop
kind: builtin-command
goals: ci-automation, incident-response
best_when: work should recur on a time interval — polling a deploy, re-checking a queue — rather than converge on a condition
setup: none
composes_with: goal, autonomous-loops, scheduled-agents
exemplar: /loop 5m check the deploy status and summarize changes
session_signal: user ran /loop in this conversation
pitfalls:
- `/loop` re-runs on an interval and stops when you stop it or Claude decides the work is done; `/goal` runs to a condition. Polling wants /loop, convergence wants /goal.
- It runs in the open session on this machine — for work that must survive a closed laptop, use /schedule instead.
source: https://code.claude.com/docs/en/commands

## deep-research

id: deep-research
kind: builtin-command
goals: research, dependency-management, incident-response
best_when: a decision needs a multi-source, adversarially verified answer rather than one page's opinion
setup: some
composes_with: deep-research, plan-mode
exemplar: /deep-research Compare date-fns and Luxon as Moment.js replacements for timezone-heavy scheduling code
session_signal: user ran /deep-research in this conversation
pitfalls:
- It is a bundled dynamic workflow on paid plans (on Pro, enable Dynamic workflows in /config); it runs in the background — check progress with /workflows.
- Underspecified questions produce generic surveys; add constraints (stack, scale, region) before running.
source: https://code.claude.com/docs/en/workflows

## schedule

id: schedule
kind: builtin-command
goals: ci-automation, dependency-management, release-management
best_when: recurring maintenance should happen on Anthropic-run cloud infrastructure with the laptop closed
setup: some
composes_with: scheduled-agents, headless-mode
exemplar: /schedule daily Renovate-PR triage at 6am
session_signal: user ran /schedule in this conversation, or routines exist for their account
pitfalls:
- Research preview: needs a claude.ai subscription login; minimum interval one hour; runs clone from GitHub, never from local disk.
- A green run status means the session completed, not that the task succeeded — read the transcript.
source: https://code.claude.com/docs/en/routines

## init

id: init
kind: builtin-command
goals: onboarding, documentation, code-understanding
best_when: a repo has no CLAUDE.md and every session re-explains the same build commands and conventions
setup: none
composes_with: project-memory
exemplar: /init
session_signal: a CLAUDE.md generated by /init exists in the repo
pitfalls:
- If a CLAUDE.md already exists it suggests improvements instead of overwriting — it is safe to run on an established repo.
source: https://code.claude.com/docs/en/memory
