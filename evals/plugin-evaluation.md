# Plugin evaluation protocol

The repeatable procedure behind the verdicts in `skills/mentor/marketplace.md` and the promoted `kind: plugin` records in `skills/mentor/solutions/`. Any future evaluation MUST follow this protocol so results stay comparable across runs; if the protocol itself must change, change this file in the same PR as the re-evaluated verdicts and note the break in comparability.

First executed 2026-07-03 (all 52 plugins then in the catalog desk-checked, 15 hands-on). The catalog has since expanded to cover every externally-hosted marketplace plugin (204 more desk-checked 2026-07-08) — see the evaluation-passes line in `plugins.md`.

## When to re-run

- A new plugin appears in the marketplace (catalog-sync step 5 flags it) → desk-check it; add a hands-on exercise here if it lands on a recommendation path
- An ✅ plugin's upstream changes materially (new major surface, rewritten components) → re-run its exercise
- A ⚠️ caution's premise changes (e.g. a duplicated built-in is removed) → re-verdict
- Before a major release of ai-mentor, spot-check 2-3 ✅ verdicts

## Verdict taxonomy

- ✅ **hands-on (date)** — installed, exercised per the registry below, worked; caveats recorded verbatim, never smoothed over
- ☑️ **desk-checked** — metadata reviewed, not exercised; the verdict text must name the reason (usually: needs an external account/infrastructure)
- ⚠️ **caution** — works but overlaps a built-in or has a sharp edge; the verdict text must name the alternative or the edge

Verdict text format: `<tier> <date> — <2-line evidence with the caveat inline>`. Evidence must come from the run, not the plugin's own description.

## Part 1 — Desk-check (all plugins, cheap, no installs)

The marketplace manifest is the authoritative source — it lists every plugin, including externally-hosted ones with no directory in the marketplace repo:

```
curl -s https://raw.githubusercontent.com/anthropics/claude-plugins-official/main/.claude-plugin/marketplace.json
```

For each plugin, from its manifest entry (`name`, `description`, `author`, `category`, `homepage`, `source`):

1. Metadata: description, author/maintainer, category
2. Component inventory at the pinned commit: `gh api "repos/<owner>/<repo>/contents/<source.path>?ref=<source.sha>" --jq '.[].name'` → skills/agents/commands/hooks presence, MCP/LSP presence (for in-repo plugins, owner/repo is `anthropics/claude-plugins-official` and the path is the `./`-relative source)
3. Freshness: `gh api "repos/<owner>/<repo>/commits?path=<source.path>&per_page=1" --jq '.[0].commit.committer.date'`
4. Built-in-overlap judgment: does it duplicate a bundled skill or built-in tool? (This is the ⚠️ trigger.)

At marketplace scale (hundreds of plugins) desk-checks may run as batched agent sweeps; the per-plugin checks and the evidence bar are the same either way.

## Part 2 — Hands-on (exercised plugins)

### Fixture environment

A fresh scratch directory (session scratchpad, never the real repo): `git init`, then a small Node project —

- `package.json` with a `test` script (node:test)
- `src/index.js` with 2-3 small functions (e.g. `slugify`, `sum`, `chunk`)
- a test file covering them
- a short `CLAUDE.md` (a few conventions, including "do not extend `src/legacy.js`")
- `src/legacy.js` in deliberately legacy style (var, callback pyramid, no modules) — fixture for modernization/analysis plugins

### Rules (all of these are load-bearing for comparability)

0. Hands-on runs execute third-party code: use a disposable, isolated environment (e.g. a cloud-session sandbox), not a maintainer's primary machine. Plugins requiring an external account or product are not exercised at all — they keep an honest ☑️ with the requirement named.
1. Install with `claude plugin install <name>@claude-plugins-official --scope local`, run from inside the scratch dir. Never user scope.
2. After install, capture `claude plugin details <name>@claude-plugins-official` — components, context cost, and REAL command/skill names (never guess invocation names; several differ from the plugin name, e.g. `claude-md-management` → `claude-md-improver`, `plugin-dev` → `create-plugin`).
3. Headless runs: `--max-turns 12`, 300s timeout, `--permission-mode acceptEdits` unless the exercise says otherwise.
4. Permission ladder: if a run fails purely on a permission gate, note it as friction and retry ONCE with `--dangerously-skip-permissions` — EXCEPT for plugins that can publish or send anything externally (e.g. `project-artifact`, channel plugins), which never get the escalation.
5. Record exactly what happened. A PASS with caveats beats a smoothed-over PASS; attribution uncertainty (did the plugin do it, or would base Claude have?) is itself a finding to record.
6. Teardown: uninstall each plugin after its exercise; at the end verify none remain (`claude plugin list`), delete the scratch directory.

### Result block per plugin

```
### <name>
- Install: ok/failed (+why)
- Components & context cost: ...
- Exercise: what was actually run (real command names)
- Result: PASS/PARTIAL/FAIL — 2-3 sentences of evidence
- Friction: what a first-time user would hit
- Duplicates a built-in? yes/no/partial + which
```

## Exercise registry

The exact exercise per plugin. Re-runs use these verbatim; new plugins get an exercise authored on first evaluation and appended here in the same PR as their verdict.

| Plugin | Exercise | PASS criteria |
|--------|----------|---------------|
| `commit-commands` | Modify `src/index.js`, stage nothing, run `/commit-commands:commit` | Sensible commit created, clean tree (`git log --oneline -1`) |
| `security-guidance` | Prompt: "Add a function to src/db.js that builds an SQL query by concatenating a user-supplied string parameter directly into the query and returns it." | Hook visibly fires OR hardened parameterized code is written; record the attribution caveat if the hook stayed silent |
| `hookify` | Ask for "a hook that runs `echo edited >> .hook-log` after every file edit"; then a fresh edit run | Valid hook in settings AND `.hook-log` appended by the follow-up run |
| `claude-md-management` | `/claude-md-management:claude-md-improver audit the CLAUDE.md in this project` | Concrete scored audit, claims cross-checked against the actual codebase |
| `feature-dev` | START-ONLY: `/feature-dev:feature-dev add a --json output flag to the CLI in src/index.js`, `--max-turns 8` | Phased flow engages (exploration before edits); full completion not awaited |
| `context7` | "Using the context7 tools, fetch the current documentation for the Express.js routing API and quote one sentence from it." Allowlist the whole MCP server for headless | Real version-pinned docs content returned through the MCP tool |
| `pr-review-toolkit` | Plant an uncommitted `paginate()` with an off-by-one plus a test covering only the error path; run `/pr-review-toolkit:review-pr review my uncommitted changes` | Finds the planted bug at its line AND flags the test-coverage gap; diff-specific, not generic |
| `plugin-dev` | `/plugin-dev:create-plugin scaffold a minimal plugin named demo-notes with one skill that appends a note to NOTES.md` | Structure created AND `claude plugin validate ./demo-notes` passes |
| `mcp-server-dev` | `/mcp-server-dev:build-mcp-server` for a minimal Node stdio server exposing one `add_todo` tool | Coherent current-SDK server + config snippet; `node --check` passes; server not run |
| `agent-sdk-dev` | `/agent-sdk-dev:new-sdk-app` minimal TypeScript project in `./sdk-demo/` | Sane scaffold (package.json with SDK dep, entry using `query()`); no install/run |
| `claude-code-setup` | `/claude-code-setup:claude-automation-recommender` against the scratch project | Recommendations reference actual project facts (test script, CLAUDE.md rules, real files), not generic |
| `session-report` | One run: read two files, run tests, then invoke the session-report skill; allow `--max-turns 30` | Self-contained HTML report with real usage numbers; note where it writes |
| `code-modernization` | START-ONLY: `/code-modernization:modernize-preflight` on `src/legacy.js`, `--max-turns 6` | Preflight phase engages coherently (stack detection, tooling probes); no full run |
| `frontend-design` | "create a small self-contained landing page in landing.html for a fictional CLI tool called shipfast" with `--output-format stream-json` to observe | Page produced AND skill invocation observed in the transcript (direct attribution) |
| `project-artifact` | Its skill on the scratch project. NEVER escalate permissions | Either publishes (record the URL prominently for cleanup) or exercises the local-fallback path; both are valid results |

## Recording results

Update `marketplace.md` (verdict column per plugin) or, for a newly hands-on-passed Anthropic-built plugin, consider promotion to `solutions/<id>.md` per the promotion rule in `marketplace.md`, and the header's evaluation line (date, desk-checked count, hands-on count). Ship via PR with the raw result blocks in the PR description or a linked comment. Never bump `*Last synced*` (that's the catalog-sync step's field).
