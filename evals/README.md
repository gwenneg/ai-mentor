# ai-mentor evals

Benchmark the skill against its baseline alternative (asking Claude directly, no skill) and catch regressions in the discovery-first behavior before release.

Group key — **A**: problem mode (classifies correctly + output shape). **B**: growth mode (bare invocation picks the right opener). **C**: never-repeat (the profile promise holds). **D**: trigger calibration (fires on mentor-shaped questions, silent otherwise; interactive-only).

## Why this exists

Two claims need numbers and a harness to stay true:

- **The differentiation claim**: verified, grounded, profile-personalized recommendations beat generated ones. Headline metric: "X% of unassisted answers contained a command that doesn't exist; the skill's contained none."
- **The discovery claim**: the skill teaches what the user *doesn't* know and never repeats what they do. This is what the profile-fixture cases (Groups B and C in `cases.md`) protect.

## Method

For each case, run in a fixture repo — the committed one at `evals/fixture/` is a small Go HTTP service (stdlib-only, `go test ./...`, real routes in `server.go` deliberately unmentioned by its CLAUDE.md so grounding can require an actual scan). Any real project checkout also works; one with a build/test config and some `.claude/` config exercises grounding and signal detection best. Groups B and C additionally control the profile state: write the fixture `~/.ai-mentor/profile.md` before the run, inspect it after, and remove it between cases — the profile is global to the machine, so eval runs pollute a real one. Back up any real profile first.

```bash
# Arm A — the skill (Groups A-C)
claude -p "/ai-mentor:mentor <problem statement>" --output-format json > results/<id>-skill.json

# Arm B — baseline, no skill (Group A only)
claude -p "What's the best way to use AI tooling for this: <problem statement>" \
  --output-format json > results/<id>-baseline.json
```

**Group D (trigger calibration) is interactive-only.** Model-triggered invocation never fires in `-p` mode — verified empirically (3/3 misses, 2026-07-03) — so headless trigger results are meaningless. Type the Group D prompts into a live session and record fire/no-fire by hand.

## Scoring

| Check | Applies to | Pass condition |
|-------|-----------|----------------|
| Classification | skill, Group A | Matches the expected goal in `cases.md` |
| No fabrication | both arms | Every command, flag, and slash command exists in current official docs |
| Grounding | skill | The move's fenced prompt references real paths/commands from the fixture repo |
| One move + one surprise | skill | Exactly one primary move and one labeled surprising pick; ranked list only after "more" |
| Personalized surprise | skill | The surprising pick is not marked adopted/shown/declined in the profile fixture |
| Never-repeat | skill, Groups B-C | Shown not re-taught, declined never re-offered, adopted built upon |
| Profile hygiene | skill | Profile updated in-flow with correct schema; zero permission prompts in the run |
| Trigger precision/recall | skill, Group D | Fires on D01-D03, silent on D04-D06 |

The **no-fabrication comparison between arms** remains the headline metric. The rest are regression checks on the skill's own behavior.

## Cadence

- **Groups A-C run in CI and gate** (`.github/workflows/evals.yml`, via `tools/eval-runner`): automatically on the standing release PR (`release/next`) — so every state that could be tagged gets a full, blocking run — and on manual dispatch for feature branches (`gh workflow run evals.yml -f cases=A19,A20`, add `-f gate=false` for report-only). What gates are the deterministic expectations in `cases.md` (classification, fabrication, grounding, output shape, never-repeat, profile hygiene); any fuzzy quality rubric added later stays advisory — it informs, never blocks. Results post to the job summary and as a PR comment either way. **Enforcement lives on `main`, never on `release/next`**: the `evals` check is required by the `main` branch ruleset, which gates merging the release PR (feature PRs satisfy it with a skipped run — the job's `if` limits real runs to `release/next`). Do not add a required-check rule to `release/next` itself: a ruleset's required checks apply to every push to the branch, and since the check can only run after a commit lands, such a rule blocks all pushes — including the release workflow's own version-bump refresh. Cases run concurrently (`-j`, default 3) and each is fully isolated — its own temp `HOME` and its own copy of the fixture repo — so profile fixtures never touch a real profile and parallel cases never observe each other's edits. CI auth accepts either repo secret: `ANTHROPIC_API_KEY` (API billing) or `CLAUDE_CODE_OAUTH_TOKEN` (subscription; generate with `claude setup-token`). When both are set the OAuth token wins, so a stale API key can never shadow a fresh token. CI runs at `-j 9` (the `jobs` dispatch input overrides): subscription limits are a usage pool (5-hour rolling window + weekly cap), not a concurrency cap, so parallelism only compresses wall clock — total spend is identical at any `-j`, and CI draws from the same pool as interactive use either way. The signal that `-j` is too high is ERROR results surviving the runner's retry pass — lower it then.
- **Fix PRs prove their fix in the PR.** A feature PR whose body carries a line-anchored `Evals: <space-separated case ids>` marker (or `Evals: smoke`; optional `Eval-Epochs: N`, default 6) runs exactly those cases on the PR head — gated and pass^k (the runner's `-passk` flag): every targeted case must pass every epoch before merge, because a fix that passes 4 of 6 epochs is not a demonstrated fix. The `evals` required check turns red and blocks the merge otherwise; each run posts its report as a PR comment. PRs without the marker keep today's skipped-run behavior. Declare `Evals: smoke` (or add cases to the list) when a change could ripple beyond its target; the full 3-epoch suite still gates every release on `release/next`, so a missed side effect is caught where it always was — minus the blind merge.
- **Three verdict layers, cheapest first.** (1) *Deterministic pre-checks* (`detChecks` in the runner): mechanical verdicts — profile file created with rows (B01), declined names absent (B03/C05), seeded rows surviving (C04) — decided by plain code before any judge call; a failure skips the judge entirely. (2) *The LLM judge* for everything requiring judgment, against handed ground truth only. (3) *Strict invariants* — cases whose `cases.md` row carries the `[strict]` marker (currently A30, B03, C04, C05) — gate at pass^k: every epoch must pass, never a majority; a promise kept two epochs out of three is a broken promise. Strictness travels with the case row, so a new invariant case declares itself; the deterministic registry fails the run loudly if its IDs drift from `cases.md`.
- **Judge calibration runs on records.** Every CI run writes per-epoch verdict records (`-records`, uploaded as the `eval-records` artifact, 14-day retention): case, verdict, judge's full reply, the mentor's response, the after-run profile. To calibrate the judge: download a run's records, sample ~30 judge-scored entries (mix PASS and FAIL), read each response blind and mark agree/disagree with the judge's verdict — persistent disagreement on a check means the judge prompt (or the case) needs fixing, and no other eval investment is trustworthy until it does. Deterministic FAILs skip the judge, so their records carry no judge verdict — audit them directly from their reason/response instead of scoring agreement. The records are also the raw material for per-case flake history: first-attempt ERROR epochs are preserved, with retries appended as extra records.
- **Coverage is audited, not assumed**: `evals/coverage.md` maps every load-bearing skill rule × {should-trigger, should-not-trigger, invariance} to the cases that exercise it; gaps are named there and drive new cases — the suite grows by coverage holes, never to hit a case count.
- **Two runner knobs keep iteration cheap and gates honest.** `-smoke` runs a curated tier of one case per behavior class (the `smokeCases` list in `tools/eval-runner/main.go` — the runner fails loudly if the list drifts from `cases.md`) at roughly a quarter of a full run's cost; use it per change, and keep the full suite for release gating. `-epochs N` runs every selected case N times and passes it only on a strict majority (`[strict]`-marked cases need every epoch), flagging majority-pass mixed results `FLAKY` in the report — several cases are known to flip between identical runs, so use `-epochs 3` when a verdict matters (confirming a fix) and the default single run while iterating. The release gate runs at 3 epochs automatically (the workflow's `EPOCHS` fallback for non-dispatch runs) — for majority-gated cases a single flaky flip can neither block nor fake a release verdict, while `[strict]` cases are the deliberate exception: there one failed epoch IS the verdict. Dispatch runs keep their chosen `-f epochs` (default 1).
- **Prompt growth is a reviewed number.** PRs touching the always-loaded skill files (SKILL.md, the two mode files) get a token-budget comment (`.github/workflows/skill-budget.yml`, a self-contained curl/jq step against the free `count_tokens` endpoint): per-file base→head deltas plus the cost of the two invocation paths (SKILL.md + one mode file). Comment-only, never a gate — the reviewer decides whether an eval fix earned its tokens. Skipped on `release/next` (its PR carries the eval-report comment; `--edit-last` would collide).
- Group D (trigger calibration) stays interactive-only: run it by hand whenever the frontmatter `description` changes (it's the trigger's tuning dial).
- Re-run the baseline arm when a new Claude model ships — the baseline improves over time, and the differentiation claim should be re-measured, not assumed. The baseline arm is not part of the CI run; it measures the differentiation claim, not regressions.

The skill-creator plugin (`/plugin install skill-creator@claude-plugins-official`) provides tooling for skill evals, including trigger-accuracy testing, and can replace this manual harness once the case set stabilizes.
