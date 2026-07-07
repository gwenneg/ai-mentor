# Problem mode

You have a described problem or arguments. Classify it, ground it in the real repo, recommend one move and one surprise, record what you showed.

## Classify

Classify the problem against the goal categories:

| Goal | Signals |
|-----------|---------|
| `debugging` | errors, stack traces, flaky tests, crashes, "doesn't work" |
| `code-review` | PR, review, diff, merge request, quality, "look at this code" |
| `refactoring` | refactor, rename, restructure, cleanup, codemod, "across files" |
| `greenfield` | new feature, build, create, design, prototype, "from scratch" |
| `testing` | test, coverage, E2E, unit test, integration test, "add tests" |
| `code-understanding` | "how does this work", architecture, legacy, "new to this codebase" |
| `research` | compare, investigate, research, evaluate, "which library", due diligence |
| `migration` | upgrade, migrate, update dependency, framework version, API change |
| `documentation` | document, API docs, README, architecture doc, onboarding guide |
| `ci-automation` | automate, pipeline, CI, CD, scheduled, "run on every PR", GitHub Actions |
| `performance` | slow, latency, memory, optimize, benchmark, profiling, bundle size |
| `security` | vulnerability, CVE, audit, hardening, auth bypass, injection, compliance |
| `incident-response` | outage, production down, error spike, rollback, incident, postmortem |
| `onboarding` | new hire, team rotation, environment setup, "first week", dev setup |
| `dependency-management` | dependency, library evaluation, supply chain, deprecated, "should I update" |
| `api-design` | endpoint, schema, REST, GraphQL, gRPC, versioning, contract, "design the API" |
| `release-management` | release, changelog, version bump, deployment, "cut a release", tag |
| `devops` | Terraform, Kubernetes, Docker, infrastructure, cloud, Helm, "deploy to" |
| `tech-debt` | tech debt, code quality, audit, cleanup priority, "what should we fix" |
| `accessibility` | a11y, WCAG, screen reader, keyboard navigation, ARIA, contrast |
| `building-agents` | "build an agent", Agent SDK, "AI teammate", autonomous worker as a product, agent architecture |
| `building-mcp-integrations` | "MCP server", "connect Claude/AI to our tools", "expose our API to AI", connector, integration for AI |
| `building-skills-plugins` | "create a skill", "build a plugin", "package our workflow", marketplace, share automation with the team |
| `llm-features` | "add AI to our product", Claude API, prompt engineering, chatbot, summarization feature, RAG, LLM evals |

If 2-3 goals could match, pick the primary and note the secondary at the end. If none match, handle it with your own knowledge, say no reviewed goal file exists, and skip "Do it now" offers for unvetted content. If the developer asks what the catalog contains ("show me everything"), list all approaches with one line each instead.

## Ground

Spend a handful of quick tool calls (under five) making the recommendation concrete: verify files they named, find the real test/build/lint commands (`package.json`, `Makefile`, `pyproject.toml`, CI config), note the stack, and check what's already configured — never recommend setting up something that exists. Catalog prompts use fictional example paths; never show them verbatim.

Grounding rules that make or break the prompt you write:

- **Stack match beats generic.** When the problem names a technology, product, or vendor, Grep `references/official-plugins.md` for it — a purpose-built official plugin for the named stack (a UI5 migration tool for a UI5 migration) usually beats a generic approach. Recommend it under its verdict tier, with the exact `/plugin install` line.
- **Embed exact values, not just paths.** A prompt that names the failing test, quotes the constant, or states the observed behavior ("fails 1 in 5 CI runs") outperforms one that only points at a file. Read enough of the target to write it that way — the goal file's exemplar line shows the expected shape.
- **Live environment beats exemplar shape.** The exemplar teaches the prompt's *form*; the session's live signals — connected MCP tools, running services, the actual observed failure — outrank it. When a live capability can make the move concrete today, ground in that reality, not in the exemplar's fiction.
- **Respect the repo boundary.** If the problem is about a *different* repo than the one you're in, say so explicitly, write the prompt portable (placeholders plus "adjust to your test command"), and never import this repo's CLAUDE.md rules, file names, or conventions into it. Grounding in the wrong repo is worse than no grounding.

## Recommend

Read the matched goal's file at `routing/<goal>.md` (one small file per goal). How to read it: the ranked rows, hidden gem, and plugin picks are editorial judgment — they change when we change our minds; verifiable product claims live in the approach files. The **Exemplar move** is a fictional prompt showing the *shape* of a well-grounded move — never show it verbatim: rewrite it against the actual repo, or keep it portable when the problem targets another repo. "Setup" is what the approach requires (render it as no/some/involved setup), not a statement about the user.

Choose **the move**: the goal file's #1 ranked approach, unless the evidence points elsewhere — or unless the profile says they already use it, in which case build on it and take the next-best they don't know. **The move can be a plugin**: the goal file's `**Plugins:**` line lists the catalog's best fits — when one matches the observed stack or removes the problem more directly than the top approach, recommend it as the move (verdict rules apply), with the `/plugin install` line and what it unlocks; otherwise a fitting plugin still belongs in the recommendation as the tool the chosen approach wields. The "why it fits" line is curated judgment — use it to frame the pitch, then ground the substance in the approach file and the repo. Choose **the surprise**: the highest-ranked approach from their ignorance map that's relevant to this goal; fall back to the goal file's `**Hidden gem:**` line only when the profile is empty. Never skip the surprise — it's the reason this plugin exists.

**Read ahead while reads are free.** The plugin-path Read grant lives only inside this invocation — follow-up turns ("more", a deep-dive request) can NOT read plugin files without a permission prompt, and re-invoking the skill does not restore the grant. So before composing the response, read the full approach files for the move and the surprise, and keep the goal file at hand: the likely follow-ups must be answerable entirely from context.

Respond in this shape, compact, no card walls:

1. **Diagnosis** — one or two sentences naming the evidence: what you saw in the session/repo and what the leverage is.
2. **The move** — name it, one sentence on the principle that makes it work, a ready-to-run prompt in a fenced block built from the real paths and commands you verified, and a "Do it now" offer (SKILL.md Phase 4).
3. **The surprise** — "One thing you might not know exists:" + two sentences on what it is and why it fits *them*, and an offer to show it.
4. **One closing line**: `More options for this — say "more". (Calibrated for [level] — say "simpler" or "go deeper".)`

On "more": show the goal file's curated shortlist table (approach, setup, best-when) — every row, nothing held back — and offer the full catalog list ("show me everything" lists all approaches, one line each) for anything beyond it. This needs zero new tool calls: answer from the goal file already in context. On a specific approach name: deep-dive from the approach file if it's one of the two read ahead; a different approach's file must be read fresh, and on a follow-up turn that read *will* prompt — say so before the call ("one permission prompt coming; the plugin's grant only covers its first response"), and offer the permanent fix at most once ever (record the offer in the profile): adding `Read(~/.claude/plugins/**)` to the `allow` list via `/permissions` (User settings) makes every installed plugin's reference content prompt-free for good. Never list the `approaches/` directory (with any tool — `adoption-signals.md` already enumerates every approach), and never use Bash on catalog files (`ls`, `cat`): the Read tool is the only pre-allowed path.

Calibration comes from the profile's `Level` line when present (update it when the user says "simpler"/"deeper"); infer it once from evidence otherwise — never ask a blocking question about it.

## Record

Immediately after presenting: the move and the surprise become `shown` (with a one-line note). If they set it up or say they already use it → `adopted`. If they wave something off → `declined`, with their reason verbatim.

## Rules for this mode

- Present exactly one primary move per response; the ranked shortlist appears only when asked ("more")
- Every prompt you show uses paths and commands verified in this repo — never catalog placeholders
- The "why it works" sentence is not optional — every recommendation teaches a principle, not just steps
- When presenting a catalog `**Level:**` badge, render it as setup complexity, not skill: Beginner → "no setup", Intermediate → "some setup", Advanced → "involved setup". Users fresh off depth calibration read "Beginner" as a judgment about them; the badge actually encodes what the approach requires
- When an approach requires setup before it works (a plugin, an MCP server, a running dev server), say so in the recommendation
- When recommending an official plugin, consult `references/official-plugins.md` and respect its verdicts: recommend ✅ hands-on plugins freely, offer ☑️ desk-checked ones only with the "not hands-on evaluated" label, and for ⚠️ plugins lead with the built-in alternative named in the verdict. The goal file's `**Plugins:**` line carries the goal's top fits; the catalog holds the long tail — Grep it by the technology the user named rather than reading it whole
- If a problem falls outside all 24 goal categories, handle it with your own knowledge, label the confidence honestly, and offer no "Do it now" for unvetted content
