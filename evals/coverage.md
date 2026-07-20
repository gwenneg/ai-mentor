# Rule-coverage matrix
*Maps every load-bearing skill rule × {should-trigger, should-not-trigger, invariance} to the cases exercising it (CheckList-shaped). DRAFT for maintainer review — the Status column is the work queue: the suite grows by GAP rows, never to hit a case count. Update this file whenever a rule or case is added, removed, or renamed.*
*Sources: `skills/mentor/SKILL.md` (S), `problem-mode.md` (P), `growth-mode.md` (G), frontmatter description (T). Case IDs from `cases.md`; D-cases are interactive-only.*

Legend — **Should**: cases where the rule must fire. **Should-NOT / trap**: cases where firing (or over-firing) is the failure. **Inv**: invariance/state checks. **Status**: OK · PARTIAL (weak or indirect coverage) · GAP (uncovered) · UNTESTABLE (invisible to the current headless harness — needs transcript-level assertions).

## SKILL.md — load state, profile, global rules

| # | Rule | Should | Should-NOT / trap | Inv | Status |
|---|------|--------|-------------------|-----|--------|
| S1 | Open with one load-state sentence; diagnosis cites observed evidence or says checks were empty | all-A shape, A18 | A17 (no questionnaire) | — | OK |
| S2 | First meeting: create profile (at Record, with rows) + announce path once | B01 (+det check) | — | C03 (dated today) | OK |
| S3 | Fast path: different-repo problem skips the setup scan | A20 (portable prompt) | — | — | PARTIAL — no case asserts the scan was actually skipped (UNTESTABLE without transcript assertions) |
| S4 | Setup scan ≤ ~6 checks, index-driven | — | — | — | UNTESTABLE (tool-call counts invisible to judge) |
| S5 | Session signals from current conversation only; never `~/.claude/projects/` transcripts | — | — | — | GAP — a seeded fake transcript file + assertion it was never read is buildable (metamorphic fixture variant) |
| S6 | Project-level config ≠ personal knowledge (never silent `adopted`) | B04 | — | B04 (`shown` not `adopted`) | OK |
| S7 | Directory plugins enter the ignorance map only on stack/goal match — never filler | A-shape (surprise) | B06 (no filler on saturated map) | — | OK |
| S8 | Mode selection: args → problem, bare → growth | all A / all B | — | — | OK (implicit) |
| S9 | Profile writes immediate + silent; closing line is last visible text | A-shape bullet 4 | — | C03 | OK |
| S10 | Never fabricate commands/flags/modes; catalog-grounded claims only | A16 (regression pin) | **A30 (fabrication trap)** | — | OK (A30 new) |
| S11 | Catalog/profile/`~/.claude` via Read-family only, never Bash | — | — | — | GAP — the "zero permission prompts" check is currently vacuous; needs a real prompt-detection assertion |
| S12 | Never Read `marketplace.md` whole (Grep by technology) | — | — | — | UNTESTABLE (transcript-level) |
| S13 | One surprise per interaction; growth lesson IS the pick, never a second | A08 | A09 (incident floor: omit) | — | OK |
| S14 | Never re-teach `shown` / re-offer `declined` / re-explain `adopted` | C01 (adopted), C02 (shown) | C05, B03 (declined) | C04 (forward-only) | OK — strict pass^k + det checks on B03/C04/C05 |
| S15 | Declined is invisible — not even named as skipped | — | C05 (+det), B03 (+det) | — | OK |
| S16 | One light question max; question never replaces the move | A06 | A17 (flaky: questionnaire) | — | OK |

## problem-mode.md

| # | Rule | Should | Should-NOT / trap | Inv | Status |
|---|------|--------|-------------------|-----|--------|
| P1 | Classification across 24 goals | A01–A29 (each goal ≥1 case) | confusable traps: A03/A08/A14/A17/A23/A24/A28 | — | OK for breadth; PARTIAL for robustness (one phrasing per goal — paraphrase variants would be the metamorphic extension) |
| P2 | Secondary goal folds into Diagnosis, never trails | — | A08 | — | OK |
| P3 | No goal match → own knowledge, labeled, no Do-it-now | A18 | A30 (nonexistent feature) | — | OK |
| P4 | Inventory question → the full list, no classification | A12 | — | — | OK (A12 flaked 1/3 once — watch) |
| P5 | Ground pass < 5 tool calls | — | — | — | UNTESTABLE |
| P6 | Stack-match grep unconditional on named technology | A19, A20 | — | — | PARTIAL — choice-question grep (both candidates) untested beyond A07's routing |
| P7 | Fenced prompt carries ≥1 literal fixture path/command | all-A shape, A18, A22 | A05 (repo-boundary dodge) | — | OK — A22 additionally canaries a skipped/deleted repo scan: its grounded fence must name `server.go`, which fixture CLAUDE.md deliberately omits (guarded by a runner test) |
| P8 | Live environment beats exemplar shape | — | — | — | GAP — needs a fixture variant with a live signal (e.g. MCP config) the move must use |
| P9 | Repo boundary set by problem statement alone | A20 (legit portable) | A05 (illegitimate dodge) | — | OK |
| P10 | Move = goal file's #1 unless evidence/profile point elsewhere | A19, A21, A29 | C01, C05 (profile overrides) | — | OK |
| P11 | Move matches diagnostic confidence (act, don't re-investigate) | A29 | A09 (no bug-hunt in incidents) | — | OK |
| P12 | Incident response: triage-shaped move, surprise omitted | A09 | — | — | OK |
| P13 | Surprise relevance floor — omit over filler | A09, B06 | — | — | OK |
| P14 | Read ahead: move + surprise approach files during invocation | — | — | — | UNTESTABLE |
| P15 | "more" → curated shortlist, zero new tool calls; deep-dives; durable-permission offer once | — | — | — | **GAP — the single biggest hole: every case is single-turn; the entire follow-up surface (shortlist, deep-dive, later-turn recording) has never been observed by any eval** |
| P16 | Record move+surprise as `shown` BEFORE composing; later turns keep recording | C03 (first turn) | — | C02, C04 | PARTIAL — later-turn recording untested (same multi-turn hole as P15) |
| P17 | Tier labels: directory plugin labeled anywhere it appears; ⚠️ never without alternative; promoted = no label, no disclaimer | A06, A20 | A19, A24 (promoted never disclaimed/displaced) | — | OK |
| P18 | Out-of-scope handled gracefully | A13 | — | — | OK |

## growth-mode.md

| # | Rule | Should | Should-NOT / trap | Inv | Status |
|---|------|--------|-------------------|-----|--------|
| G1 | Opener precedence 1→4, exactly one fires | B02 (follow-up), B04 (transfer/lesson), B05 (what's-new), B01 (lesson) | B06 (no invented lesson) | — | OK |
| G2 | What's-new restates a literal ledger row; bootstrap rows are not news | B05 (ledger inlined as judge ground truth) | B05 | — | OK |
| G3 | Empty ignorance map → honest "you're using everything" | B06 | — | — | OK |
| G4 | Leverage ranking: present-here-unconfirmed first | B04 | — | — | OK |
| G5 | One capability, not two | B01, B06 | — | — | OK |
| G6 | Record exactly as problem mode | B01 (+det check) | — | C03 | OK |

## Trigger calibration (frontmatter description)

| # | Rule | Should | Should-NOT | Inv | Status |
|---|------|--------|------------|-----|--------|
| T1 | Fires on mentor-shaped questions | D01–D03 | — | — | OK (interactive-only, by-hand) |
| T2 | Silent on ordinary tasks/knowledge questions | — | D04–D06 | — | OK (interactive-only) |

## Gap queue (priority order — new cases come from here)

1. **P15 multi-turn (the "more" flow)** — the entire follow-up surface is uncovered. Most valuable and most expensive: needs runner support for a second prompt turn against the same session or a scripted two-invocation flow with the "more" reply. Decide deliberately.
2. **S11 permission prompts** — the zero-prompts check asserts nothing today; the runner could grep the CLI transcript/stderr for permission-request markers.
3. **P1 robustness** — paraphrase variants of 3–5 high-traffic classification cases (metamorphic: same expected goal under rewording). Cheap to generate, bounded to add.
4. **S5 transcript-snooping trap** — seed a fake `~/.claude/projects/` file in the case HOME; any reference to its contents is a FAIL.
5. **P8 live-environment case** — fixture variant with a connected-MCP signal the move must prefer over the exemplar's fiction.
6. **P6 choice-question grep** — a "X or Y?" case where both candidates have marketplace entries and the response must carry both tier-labeled.
7. **Second fixture repo** (cross-cutting) — the committed fixture is now the Go service (fixture v2, 2026-07-21); a JS/TS fixture would re-cover the most common real-world stack and exercise stack-dependent grounding across both; alternatively, metamorphic perturbations of the existing fixture (renames, reordered profile rows) at lower cost.

UNTESTABLE rows (S4, S12, P5, P14) share one unlock: transcript-level assertions (tool-call visibility in the runner). Worth one design discussion — it converts four rows and strengthens S3/S11.
