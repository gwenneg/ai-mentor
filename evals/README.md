# ai-mentor evals

Benchmark the skill against its baseline alternative (asking Claude directly, no skill) and catch regressions in the discovery-first behavior before release.

## Why this exists

Two claims need numbers and a harness to stay true:

- **The differentiation claim**: verified, grounded, profile-personalized recommendations beat generated ones. Headline metric: "X% of unassisted answers contained a command that doesn't exist; the skill's contained none."
- **The discovery claim**: the skill teaches what the user *doesn't* know and never repeats what they do. This is what the profile-fixture cases (Groups B and C in `cases.md`) protect.

## Method

For each case, run in a fixture repo (any real project checkout works; one with a `package.json`, tests, and some `.claude/` config exercises grounding and signal detection best). Groups B and C additionally control the profile state: write the fixture `~/.ai-mentor/profile.md` before the run, inspect it after, and remove it between cases — the profile is global to the machine, so eval runs pollute a real one. Back up any real profile first.

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

- **Groups A-C run in CI and gate** (`.github/workflows/evals.yml`, via `tools/eval-runner`): automatically on the standing release PR (`release/next`) — so every state that could be tagged gets a full, blocking run — and on manual dispatch for feature branches (`gh workflow run evals.yml -f cases=A19,A20`, add `-f gate=false` for report-only). What gates are the deterministic expectations in `cases.md` (classification, fabrication, grounding, output shape, never-repeat, profile hygiene); any fuzzy quality rubric added later stays advisory — it informs, never blocks. Results post to the job summary and as a PR comment either way. Each case runs with an isolated `HOME`, so profile fixtures never touch a real profile.
- Group D (trigger calibration) stays interactive-only: run it by hand whenever the frontmatter `description` changes (it's the trigger's tuning dial).
- Re-run the baseline arm when a new Claude model ships — the baseline improves over time, and the differentiation claim should be re-measured, not assumed. The baseline arm is not part of the CI run; it measures the differentiation claim, not regressions.

The skill-creator plugin (`/plugin install skill-creator@claude-plugins-official`) provides tooling for skill evals, including trigger-accuracy testing, and can replace this manual harness once the case set stabilizes.
