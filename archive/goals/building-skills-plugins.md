# Building Skills & Plugins
*Last verified: 2026-07-02*

## When You're Here

Your team has workflows worth packaging: the release checklist everyone half-remembers, the scaffolding conventions, the review checklist that lives in one senior engineer's head. Skills turn those into invocable commands; plugins bundle skills, hooks, agents, and MCP config into a versioned package anyone can install. This is how team knowledge stops being tribal — and, if you publish, how it reaches the wider ecosystem through marketplaces.

The craft questions: what deserves to be a skill versus a script versus a hook, how to write instructions a model executes reliably, and how to keep the package maintained once teammates depend on it.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| A workflow you repeat and want to package | Custom skills | Start standalone in `.claude/skills/` and iterate at markdown speed |
| Skill works; time to share or distribute | Custom plugins | `plugin-dev` guides packaging; marketplaces handle distribution |
| The behavior should be automatic, not invoked | Hooks as workflow | Enforcement and automation belong in hooks, not skill instructions |
| Want to measure whether the skill actually works | Official plugins | `skill-creator` creates, improves, and benchmarks skills with evals |
| Unsure if it's a skill, a hook, or a script | Plan mode | Design the extension surface before writing any of them |

**Hidden gem:** Hooks as Workflow — the `hookify` official plugin mines your own conversation history for repeated patterns and turns them into hooks: your plugin's best features are often already hiding in what you keep asking Claude to do.

## Approaches (Ranked)

### 1. Custom Skills — Start standalone, iterate fast
**Level:** Beginner

Every good plugin starts as a working skill in `.claude/skills/<name>/SKILL.md`. Build it there first: instructions in markdown, `$ARGUMENTS` for input, bundled scripts for deterministic steps. Iterate at conversation speed with `/reload-skills` until teammates use it without asking questions — then it's ready to package.

**Try it now:**
> Create `.claude/skills/release-notes/SKILL.md`: when invoked, find the latest git tag, read commits since, categorize by conventional-commit prefix into Features/Fixes/Breaking, and prepend a dated entry to CHANGELOG.md. Then run it on this repo and let's iterate on the output format until it matches our style.

**Why this works:** Skills fail on unclear instructions, and instruction clarity is discovered by iteration — the standalone stage keeps iteration cost near zero before packaging adds ceremony.

**Pros:**
- Working prototype in minutes; refine with `/reload-skills`
- The description field teaches you skill-triggering behavior early
- Bundled scripts pin down the deterministic steps

**Cons:**
- Project-local until packaged — sharing means copy-paste at this stage
- A vague description makes the skill fire (or not fire) at the wrong times

**Deeper:** See `approaches/custom-skills.md`

---

### 2. Custom Plugins — Package and distribute with plugin-dev
**Level:** Intermediate

Packaging is a solved problem: the `plugin-dev` official plugin walks an 8-phase workflow — manifest, structure, validation, review — and `claude plugin validate` catches what you miss. Distribution goes through marketplaces: a git repo for your team, or the community marketplace (`anthropics/claude-plugins-community`) for the world.

**Try it now:**
> Install with: /plugin install plugin-dev@claude-plugins-official — then use it to package our release-notes skill and the two release hooks into a plugin called `release-toolkit`, with a manifest, a README, and validation passing. Target: our team installs it from our internal git repo with /plugin marketplace add.

**Why this works:** Distribution is where team knowledge compounds — one engineer's workflow becomes everyone's default, and version pinning means improvements ship deliberately instead of by rumor.

**Pros:**
- Guided packaging with a validator and reviewer built in
- Marketplace distribution: internal repos for teams, community for public
- Versioning (explicit `version` or commit SHA) makes updates predictable

**Cons:**
- Namespaced invocation (`/plugin-name:skill`) — plan skill names accordingly
- Published plugins create maintenance expectations; unmaintained ones rot visibly

**Deeper:** See `approaches/custom-plugins.md`

---

### 3. Hooks as Workflow — Package enforcement, not just commands
**Level:** Intermediate

The strongest plugins pair skills (what users invoke) with hooks (what happens automatically). A `hooks/hooks.json` in your plugin enforces the workflow the skill enables — format-on-edit, protected paths, post-edit checks. And when you're unsure what deserves automating, `hookify` generates hooks from patterns in your own conversation history.

**Try it now:**
> Install with: /plugin install hookify@claude-plugins-official — then have it analyze my recent sessions in this repo for repeated corrections worth automating. I suspect I keep asking for the same formatting fix and the same "don't touch the generated files" warning — turn those into hooks I can ship inside our team plugin.

**Why this works:** Skills require someone to remember to invoke them; hooks fire regardless — packaging both means the workflow works even for the teammate who never reads the README.

**Pros:**
- Hooks make plugin behavior automatic instead of opt-in
- `hookify` mines real usage instead of guessing what to automate
- Plugin hooks version and ship with the package

**Cons:**
- Auto-firing behavior needs care — a noisy hook gets the whole plugin disabled
- Hook bugs affect every edit; test before shipping to the team

**Deeper:** See `approaches/hooks-as-workflow.md`

---

### 4. Custom Agents — Ship expertise, not just automation
**Level:** Advanced

Plugins can bundle agent definitions in `agents/` — a reviewer that knows your security rules, a migration specialist with your ORM conventions. Where skills package procedures, agents package judgment: install the plugin, get the specialist.

**Try it now:**
> Add an agent to our release-toolkit plugin: `agents/release-reviewer.md`, model opus, read-only tools. Its brief: review the release diff against our ship-blockers list — unfinished feature flags, schema migrations without rollback scripts, TODO(release) comments. Then test it against last week's release branch and compare its findings to what we caught manually.

**Why this works:** Judgment scales worse than procedure — an agent definition is how one specialist's review standards run on every team, every time, without the specialist in the room.

**Pros:**
- Installable expertise with consistent quality
- Tool and model constraints ship inside the definition
- Composes with the plugin's skills (a skill can invoke the agent)

**Cons:**
- Agent instructions encode conventions that drift — review them each release
- Judgment-quality claims deserve evals before teammates trust them

**Deeper:** See `approaches/custom-agents.md`

---

### 5. Autonomous Loops — Eval your skill until it actually works
**Level:** Advanced

"It worked when I tried it" is not a quality bar for something a team depends on. The `skill-creator` official plugin measures skill performance with evals; pair that with a goal-driven loop — run the eval set, adjust instructions, re-run — and skill quality becomes a number that goes up instead of a vibe.

**Try it now:**
> Install with: /plugin install skill-creator@claude-plugins-official — then build an eval set of eight realistic invocations for our release-notes skill (empty changelog, no tags yet, breaking changes present, merge commits). /goal: all eight eval cases produce correctly categorized output. Iterate on the SKILL.md instructions until the goal holds.

**Why this works:** Skill instructions are prompts, and prompts regress silently — an eval set turns "did my edit break it?" into a question with an answer.

**Pros:**
- Objective quality bar before teammates depend on the skill
- Eval cases become the regression suite for future edits
- The loop grinds instruction-tuning you'd never do manually

**Cons:**
- Building a good eval set takes real effort — start with the failure cases you've already seen
- Token cost per iteration; keep eval sets focused

**Deeper:** See `approaches/autonomous-loops.md`
