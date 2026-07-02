# ai-mentor evals

Benchmark the skill against its baseline alternative (asking Claude directly, no skill) and catch prompt regressions before release.

## Why this exists

The skill's differentiation claim — verified recommendations beat generated ones — needs a number to be communicable and a harness to stay true. This eval produces both:

- **The blog post number**: "X% of unassisted answers contained a command that doesn't exist; the skill's contained none."
- **The regression suite**: every SKILL.md change gets validated against the same cases before release.

## Method

For each case in `cases.md`, run both arms in a fixture repo (any real project checkout works; a dedicated fixture repo with a `package.json`, tests, and `.claude/` config exercises the grounding phase best):

```bash
# Arm A — the skill
claude -p "/ai-mentor:mentor <problem statement>" --output-format json > results/<id>-skill.json

# Arm B — baseline, no skill
claude -p "What's the best way to use AI tooling for this: <problem statement>" \
  --output-format json > results/<id>-baseline.json
```

## Scoring

Score each output on these checks (manual first; automate the mechanical ones over time):

| Check | Applies to | Pass condition |
|-------|-----------|----------------|
| Classification | skill | Matches the expected goal in `cases.md` |
| No fabrication | both | Every command, flag, and slash command exists in current Claude Code docs |
| Grounding | skill | "Try it now" prompts reference real paths/commands from the fixture repo, not invented ones |
| Safe + surprising | skill | Both picks present and labeled; surprising pick differs from what the baseline arm led with |
| Do it now | skill | Every catalog card ends with an executable offer or a copy-ready command |
| Beyond-catalog discipline | skill | Unvetted suggestions carry no "Try it now"/"Do it now" |

The **no-fabrication comparison between arms** is the headline metric. The rest are regression checks on the skill's own behavior.

## Cadence

- Run the full case set before every release tag.
- Re-run the baseline arm when a new Claude model ships — the baseline improves over time, and the differentiation claim should be re-measured, not assumed.

The skill-creator plugin (`/plugin install skill-creator@claude-plugins-official`) provides tooling for skill evals and can replace this manual harness once the case set stabilizes.
