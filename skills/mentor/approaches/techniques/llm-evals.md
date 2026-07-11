# LLM Evals
*Last verified: 2026-07-11*

## What It Is

An eval suite is a test suite for AI behavior: each case pairs a real input (plus the state the AI starts in) with an expected behavior, runs it against your LLM feature, agent, or skill, and grades the result automatically — plain code where the check is mechanical, an LLM judge where the expectation is semantic. It is how you find out that a prompt change made things worse before your users do.

## Why It Works

LLM output quality is a distribution, not a property — a change's effect is invisible in any single output, and only a graded sample of that distribution turns "feels better" into a number you can gate a release on.

## When to Use It

- An LLM feature is live and prompt tweaks keep shipping regressions that users notice before you do
- You are building an agent and need reliability evidence — a demo that works once says nothing about run 8 of 8
- Before a model swap or upgrade: run old and new on the same suite and compare per case, not on vibes
- A Claude Code skill or plugin matters enough that edits must not silently change its behavior
- Before optimizing cost: proof that the cheaper model tier holds quality on your specific task

## When NOT to Use It

- Exploratory or one-off prompting — a suite protects repeated behavior; build it once a behavior is worth keeping, not before
- The deterministic code around the LLM call — parsing, retries, schema validation — belongs in ordinary unit tests, which are cheaper and exact
- When you cannot state any success criterion yet — that is a design gap, not a tooling gap: define what "good" means first (specific and measurable, never "good performance")

## How It Works

### Basic (Beginner)

1. Collect 20–50 real inputs from logs, support tickets, or your own usage — volume beats polish: many automatically graded cases outperform a few hand-graded ones
2. Write the expected behavior next to each input, specific and checkable: "classifies as billing", "names the 14-day refund window", "refuses and says why" — never "answers well"
3. Grade mechanically wherever possible (exact match, contains, valid JSON, correct label); reserve an LLM judge for semantic expectations, and give it the rubric AND a reference answer so it grades against ground truth instead of its own taste:

```bash
claude -p "Rubric: the response must name the SAVE20 code and the 14-day window.
Reference answer: <known-good answer>.
Response to grade: <the output>.
Reply with STRICT JSON only: {\"pass\": bool, \"reason\": string}" --model claude-sonnet-5
```

4. Run the suite headless on every prompt change and gate merges on the mechanical checks; keep judge-scored quality advisory until the judge has earned trust
5. Re-run a failing case before believing it — one run of a nondeterministic subject proves nothing

### Composing with Other Approaches (Intermediate)

- **LLM Evals plus Plan Mode**: design the success criteria and the eval plan as part of the feature design, before writing the feature — the Anthropic guide's ordering — so quality is measurable from the first prompt draft
- **LLM Evals plus Autonomous Loops**: `/goal 90% of eval cases pass` turns prompt engineering into test-driven development — the loop tweaks, runs the suite, reads the failures, and tweaks again
- **LLM Evals plus Headless Mode**: `claude -p` runs both the subject and the judge in CI on every PR that touches a prompt; the suite's exit code is the merge gate

### Advanced Patterns

- **Statistics, not single numbers**: report suite scores with error bars (mean ± SEM), resample nondeterministic cases several times and average, and compare two prompt or model versions with paired per-case differences — pairing cancels case-difficulty variance and detects much smaller regressions from the same runs
- **pass^k for agents**: score an agent on passing ALL k repeated trials of a case, not one; agent pass rates fall steadily as k grows, so a consistency score is the honest reliability number — gate agents on it, not on best-of demos
- **Judge discipline**: pin the judge model explicitly (some frameworks silently pick whatever the available credentials allow), inject ground truth the judge cannot invent (real file lists, real API names, the reference answer), and keep a small human-labeled anchor set — anchor scores moving while the subject is unchanged means the judge drifted, not your feature
- **Stateful subjects need seeded conditions**: when behavior depends on files, profiles, or repo state, every case must construct that state (temp HOME, per-case fixture copies) before the run — no mainstream eval framework models this natively, so it is your harness's job regardless of the tool you pick

## Common Pitfalls

- **Single-run verdicts**: on a nondeterministic subject, flakes read as regressions and burn debugging time on noise. Confirm failures with repeat runs, and flag cases that disagree with themselves as flaky rather than red.
- **An unpinned judge**: if the grading model changes under you, scores move without any change to the subject — every red becomes ambiguous. Pin the judge's model version and treat judge upgrades as suite migrations.
- **Optimizing to the judge**: a suite that never changes becomes the target instead of the proxy — outputs drift toward the judge's quirks. Refresh cases from real traffic and audit a sample of judge verdicts by hand now and then.
- **Tiny artisanal suites**: ten lovingly hand-graded cases miss the edge cases where LLM features actually fail. Prioritize case volume with automated grading, and mine production inputs rather than inventing polite ones.

## Sources

- [Define success criteria and build evaluations](https://platform.claude.com/docs/en/test-and-evaluate/develop-tests) — official guide: criteria design (specific/measurable), eval methods from exact-match to LLM-graded, grading code examples
- [A statistical approach to model evaluations](https://www.anthropic.com/research/statistical-approach-to-model-evals) — Anthropic research: SEM/confidence intervals on eval scores, resampling nondeterministic evals, paired-difference comparisons

## Signals

- Setup: an `evals/` or `eval/` directory, `promptfooconfig.yaml`, `*.eval.ts`, or Braintrust/LangSmith config in the repo
- Session: "the prompt got worse", comparing models by feel, prompt tweaks shipped without any regression check
