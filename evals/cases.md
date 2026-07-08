# Eval cases

Cases for the discovery-first skill, in four groups: **A — classification** (problem mode routes correctly, and every classified response meets the output-shape expectations: a diagnosis + one move + one surprise, not a menu), **B — growth mode** (a bare invocation picks the right opener from the profile), **C — never-repeat** (the profile rule holds under problem mode), and **D — trigger calibration** (the skill fires on mentor-shaped questions and stays quiet otherwise).

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
| A19 | `migrate our legacy COBOL billing system to Java` | migration | Response must surface `code-modernization` (✅, from the goal routing file's Plugins line) as the move or its tool |
| A20 | `convert our SAPUI5 app from JavaScript to TypeScript` | migration | Stack-match rule: must surface `ui5-typescript-conversion` with the "not hands-on evaluated" label |
| A21 | `my tests pass but I'm not convinced the feature really works` | testing | The move must be the `/verify` built-in directly (from the goal file's Built-ins line / registry), with the copy-ready command; a `verify` profile row is recorded |
| A22 | `write API docs for our orders endpoints` | documentation | |
| A23 | `build a discount-code feature from scratch` | greenfield | "From scratch" must not misroute to refactoring despite touching existing checkout code |
| A24 | `check this codebase for injection vulnerabilities before launch` | security | Must not route to code-review; there is no diff, the subject is the codebase |
| A25 | `design a versioning strategy for our public orders API` | api-design | |
| A26 | `containerize this service and deploy it to Kubernetes` | devops | |
| A27 | `cut a release with a changelog users can actually read` | release-management | |
| A28 | `where should we start paying down the mess in this codebase?` | tech-debt | Prioritization phrasing must not misroute to refactoring |

### Group A output-shape expectations (every classified case)

- Opens with a sentence naming what was checked (or is about to be checked) and why — the load-state announcement from SKILL.md, prospective or retrospective. Judged transcripts interleave brief progress narration between tool calls; that narration is acceptable opening material (it is how the announcement reads live), and the diagnosis naming observed evidence must follow it. A questionnaire is never acceptable
- Exactly **one** primary move, with a fenced prompt using at least one real path or command from the fixture repo *inside the fenced block itself* (a setup line for that same move — a `/plugin install` or `claude mcp add` — counts as part of the move, not as a second one). For technology-choice questions about the fixture repo's own future (which ORM, which library), naming the fixture's real stack and test runner IS the grounding — no file path required. The portable-prompt exception applies only when the problem targets a *different* repo than the fixture or names code the fixture does not contain (e.g. A20's SAPUI5 app in a non-UI5 fixture): then the prompt must not import fixture-repo paths or conventions
- Exactly **one** surprising pick, labeled as such, drawn from capabilities the profile doesn't mark known — or zero when the relevance floor applies (incident pressure with a narrow question, or no ignorance-map entry related to the goal/stack); never two, and never filler
- Ends with the single closing line (more options + calibration offer); the ranked list appears only after replying "more". The closing line must be the last user-visible text — trailing recaps or profile-save narration after it violate this
- No safe/surprising *card wall*: response is prose + one fenced prompt, not 3-5 formatted cards
- When a catalog plugin matches the goal or named stack, it appears with its tier label — anywhere in the response, move or surprise; a ⚠️ plugin never appears without its built-in alternative (this alternative rule binds ONLY to plugins the catalog marks ⚠️ — a ☑️ desk-checked plugin needs only its "not hands-on evaluated" label, no alternative required). This rule covers marketplace plugins only: approaches, techniques, and built-in commands carry no tier labels
- Zero permission prompts during the run

## Group B — Growth mode (bare invocation)

Run as `/ai-mentor:mentor` with a controlled `~/.ai-mentor/profile.md` fixture (set up before each case, removed after).

| ID | Profile fixture | Expected behavior |
|----|----------------|-------------------|
| B01 | No profile file | First-meeting announcement (names the profile path once); teaches ONE capability from the ignorance map; creates the profile with correct schema |
| B02 | One `shown` row from a past date | Opens by following up on the shown capability ("did it stick?") before teaching anything new |
| B03 | A `declined` row (e.g. fan-out-workflows, "too token-heavy") | The declined capability is never offered; no reference to it |
| B04 | Empty profile, but fixture repo has hooks configured in `.claude/settings.json` | hooks-as-workflow is silently recorded `adopted`, not taught; the lesson picks something else |
| B05 | Profile with `Last new-capability check: 2026-w20` (older than the newest ledger rows) | Opens with what's-new since that week when a ledger row since carries real content; when every row since is a bootstrap/no-op entry (as in this repo's ledger), simply proceeding with another opener IS the correct fall-through — no acknowledgment of the ledger check is required. The only failure is fabricating a change |
| B06 | Profile marks every approach and every registry record (built-ins, integrations) adopted/declined | Honest empty-map answer ("you're using everything I'd recommend") — plugins may only surface with concrete stack/goal relevance, never as filler; offers the catalog list, invents nothing |

## Group C — Never-repeat under problem mode

| ID | Setup | Expected behavior |
|----|-------|-------------------|
| C01 | Profile marks the matched goal's #1 approach `adopted`; run a Group A case for that goal | The move builds on the adopted approach or picks the next-best; it is NOT re-taught from scratch |
| C02 | Run the same Group A case twice in a row (same profile) | Second run's surprising pick differs from the first (first is now `shown`) |
| C03 | After any Group A run | Profile contains new `shown` rows for the move and surprise, dated today, with one-line notes |
| C04 | Profile has a `declined` row (background-agents) and a `shown` row (plan-mode) from a past date; run a Group A case | After the run the profile STILL contains both seeded rows — existing rows are never dropped (forward-only); new `shown` rows may be added, but nothing pre-existing is lost or rewritten |
| C05 | Profile marks the matched goal's #1 approach (plan-mode, for A01/debugging) `declined`; run that Group A case | Plan mode is never named anywhere in the response — declined is invisible, not even mentioned to say it is being skipped; the move is the next-best debugging approach |

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
