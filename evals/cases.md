# Eval cases

Problem statements with expected classifications. Phrasings deliberately vary in specificity and vocabulary — some name the goal outright, some only describe symptoms.

| ID | Problem statement | Expected goal | Notes |
|----|-------------------|--------------|-------|
| 01 | `debug a flaky test that only fails in CI` | debugging | The README's canonical example |
| 02 | `our checkout endpoint got slow after the last release` | performance | Symptom wording, no "performance" keyword |
| 03 | `refactor authentication across 30 files` | refactoring | Cross-file scale signal |
| 04 | `we need to move from Vue 2 to Vue 3` | migration | No "migrate" keyword |
| 05 | `review a large PR that touches billing` | code-review | |
| 06 | `I just joined this team and the codebase is huge` | code-understanding or onboarding | Either accepted; must not block on clarification for more than one round |
| 07 | `should we use Prisma or Drizzle?` | research or dependency-management | Either accepted |
| 08 | `add tests before I dare touch this legacy module` | testing | "Before refactoring" phrasing must not misroute to refactoring |
| 09 | `production error rates spiked 20 minutes ago` | incident-response | Must not route to debugging |
| 10 | `run a code review automatically on every PR` | ci-automation | |
| 11 | `our screen reader users can't complete signup` | accessibility | |
| 12 | `I want to be more productive with AI tools` | (no single goal) | Guided flow: must ask what they're working on, not dump the catalog |
| 13 | `what approaches do you know?` | (catalog browse) | Must list approaches directly, no classification |
| 14 | `help me write a poem about my cat` | (out of scope) | Must decline the catalog gracefully, no forced classification |
| 15 | `I want to build an AI agent that triages our support tickets` | building-agents | Must not route to greenfield despite "build" |
| 16 | `expose our internal ticket API to Claude` | building-mcp-integrations | |
| 17 | `package our release workflow so the whole team can use it` | building-skills-plugins | |
| 18 | `add an AI summary box to our dashboard` | llm-features | Must not route to greenfield or ui work |
| 19 | `my long session keeps getting dumber` | (no dedicated goal) | Should surface session-context-management via the relevant goal or beyond-catalog, not misclassify |

## Per-case regression expectations (all classified cases)

- Safe pick card and surprising pick card both present; surprising pick matches the goal file's Hidden gem unless the output justifies the substitution
- Every "Try it now" references at least one path or command that exists in the fixture repo
- Every catalog card has a "Do it now" line; no beyond-catalog item has one
- First response ends with the depth recalibration line instead of a blocking question
