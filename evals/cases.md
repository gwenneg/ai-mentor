# Eval cases

Cases for the discovery-first skill, in four groups: classification (problem mode routes correctly), output shape (the response is a diagnosis + one move + one surprise, not a menu), profile behavior (the never-repeat rule holds), and trigger calibration (the skill fires on mentor-shaped questions and stays quiet otherwise).

## Group A — Classification (problem mode)

Run as `/ai-mentor:mentor <statement>` in the fixture repo. Phrasings deliberately vary in specificity and vocabulary.

| ID | Problem statement | Expected goal | Notes |
|----|-------------------|--------------|-------|
| A01 | `debug a flaky test that only fails in CI` | debugging | The README's canonical example |
| A02 | `our checkout endpoint got slow after the last release` | performance | Symptom wording, no "performance" keyword |
| A03 | `refactor authentication across 30 files` | refactoring | Cross-file scale signal |
| A04 | `we need to move from Vue 2 to Vue 3` | migration | No "migrate" keyword |
| A05 | `review a large PR that touches billing` | code-review | |
| A06 | `I just joined this team and the codebase is huge` | code-understanding or onboarding | Either accepted; must not ask more than one light question |
| A07 | `should we use Prisma or Drizzle?` | research or dependency-management | Either accepted |
| A08 | `add tests before I dare touch this legacy module` | testing | "Before refactoring" phrasing must not misroute to refactoring |
| A09 | `production error rates spiked 20 minutes ago` | incident-response | Must not route to debugging |
| A10 | `run a code review automatically on every PR` | ci-automation | |
| A11 | `our screen reader users can't complete signup` | accessibility | |
| A12 | `what approaches do you know?` | (catalog browse) | Must list approaches directly, no classification |
| A13 | `help me write a poem about my cat` | (out of scope) | Must decline gracefully, no forced classification |
| A14 | `I want to build an AI agent that triages our support tickets` | building-agents | Must not route to greenfield despite "build" |
| A15 | `expose our internal ticket API to Claude` | building-mcp-integrations | |
| A16 | `package our release workflow so the whole team can use it` | building-skills-plugins | |
| A17 | `add an AI summary box to our dashboard` | llm-features | Must not route to greenfield or UI work |
| A18 | `my long session keeps getting dumber` | (no dedicated goal) | Should surface session-context-management, not misclassify |
| A19 | `migrate our legacy COBOL billing system to Java` | migration | Response must surface `code-modernization` (✅, from the routing section's Plugins line) as the move or its tool |
| A20 | `convert our SAPUI5 app from JavaScript to TypeScript` | migration | Stack-match rule: must surface `ui5-typescript-conversion` with the "not hands-on evaluated" label |

### Group A output-shape expectations (every classified case)

- Opens with the one-sentence Phase 0 announcement, then a diagnosis naming observed evidence — never a questionnaire
- Exactly **one** primary move, with a fenced prompt using at least one real path or command from the fixture repo — unless the problem targets a different repo than the fixture (e.g. A20's SAPUI5 app in a non-UI5 fixture): then SKILL.md's repo-boundary rule requires a *portable* prompt instead, which must not import fixture-repo paths or conventions
- Exactly **one** surprising pick, labeled as such, drawn from capabilities the profile doesn't mark known
- Ends with the single closing line (more options + calibration offer); the ranked list appears only after replying "more"
- No safe/surprising *card wall*: response is prose + one fenced prompt, not 3-5 formatted cards
- When a catalog plugin matches the goal or named stack, it appears with its tier label; a ⚠️ plugin never appears without its built-in alternative
- Zero permission prompts during the run

## Group B — Growth mode (bare invocation)

Run as `/ai-mentor:mentor` with a controlled `~/.ai-mentor/profile.md` fixture (set up before each case, removed after).

| ID | Profile fixture | Expected behavior |
|----|----------------|-------------------|
| B01 | No profile file | First-meeting announcement (names the profile path once); teaches ONE capability from the ignorance map; creates the profile with correct schema |
| B02 | One `shown` row from a past date | Opens by following up on the shown capability ("did it stick?") before teaching anything new |
| B03 | A `declined` row (e.g. fan-out-workflows, "too token-heavy") | The declined capability is never offered; no reference to it |
| B04 | Empty profile, but fixture repo has hooks configured in `.claude/settings.json` | hooks-as-workflow is silently recorded `adopted`, not taught; the lesson picks something else |
| B05 | `Last new-capability check` older than the newest ledger week | Opens with what's-new since that week, then updates the anchor |
| B06 | Profile marks all 26 approaches adopted/declined | Honest empty-map answer ("you're using everything I'd recommend"), offers the catalog list, invents nothing |

## Group C — Never-repeat under problem mode

| ID | Setup | Expected behavior |
|----|-------|-------------------|
| C01 | Profile marks the matched goal's #1 approach `adopted`; run a Group A case for that goal | The move builds on the adopted approach or picks the next-best; it is NOT re-taught from scratch |
| C02 | Run the same Group A case twice in a row (same profile) | Second run's surprising pick differs from the first (first is now `shown`) |
| C03 | After any Group A run | Profile contains new `shown` rows for the move and surprise, dated today, with one-line notes |

## Group D — Trigger calibration (interactive only)

**These cases cannot run headless**: model-triggered invocation never fires in `-p` mode (verified 2026-07-03; one-shot bias). Run them by typing into an interactive session with the plugin installed, without any slash command.

| ID | Prompt (typed, no slash command) | Expected |
|----|----------------------------------|----------|
| D01 | `what's the best way to use AI to add tests to this codebase?` | Skill fires |
| D02 | `I don't know how to use AI for refactoring something this big` | Skill fires (paraphrase, no trigger keyword) |
| D03 | `is there a smarter way to do this with Claude?` | Skill fires |
| D04 | `fix the failing test in tests/api.test.ts` | Skill does NOT fire — ordinary task, just do it |
| D05 | `what does this regex do?` | Skill does NOT fire — knowledge question |
| D06 | `write a README for this project` | Skill does NOT fire — ordinary task |

Score as precision/recall over the should-fire (D01-D03) and shouldn't-fire (D04-D06) sets. A miss on D01-D03 means the description is too shy; a fire on D04-D06 means too eager — the frontmatter `description` is the tuning dial.
