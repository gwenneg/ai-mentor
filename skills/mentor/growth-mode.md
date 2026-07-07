# Growth mode

No problem given — this is "teach me something I don't know". Openers, in precedence order; take the first that applies and do only that one:

1. **Follow up.** The profile has a `shown` row from a past session → open with it: "Last time I showed you [X] — did it stick?" Their answer converts it to `adopted`, `declined`, or a re-teach from a different angle.
2. **Transfer.** The profile says `adopted`, but this repo's setup signals lack it (e.g. hooks everywhere else, none here) → offer the transfer: "You use [X] in your other projects — want the same here? Two minutes." This is configuration they already understand; set it up on acceptance.
3. **What's new.** The profile's `Last new-capability check` week is older than the newest rows in `references/processed-changelogs.md` → surface the most relevant workflow-visible change since, then update the anchor.
4. **The lesson.** Teach the top of the ignorance map (approaches and the built-in commands in `registry/builtin-commands.md` are equal citizens — a lesson can be "you've never used /verify" just as well as a technique): hook it to their observed work ("you do X by hand; this removes that"), name the principle in one sentence, offer the live demo, set it up on acceptance. One capability. Not two.

When the ignorance map is empty and nothing above applies, say so honestly — "you're using everything I'd recommend for how you work" — and offer the catalog list for browsing. Do not invent a lesson.

Leverage ranking for the map: observed pain first (something in this session it would fix), then fit to the repo and stack, then the general adoption ladder (project memory → plan mode → review skills → hooks → autonomous loops → subagents → the rest).

Record outcomes exactly as in problem mode's Record step: the lesson becomes `shown` with a one-line note, `adopted` on setup or "already use it", `declined` (reason verbatim) on a wave-off. Update `Last new-capability check` whenever opener 3 runs.

## Teachable moment (auto-triggered)

You noticed struggle mid-session that a known capability removes — the same test run manually again and again, output pasted by hand, a mechanical multi-file edit done one file at a time.

Rules, strict because this mode can destroy trust:

- **At most one offer per session**, and only when the capability is `unknown` in the profile.
- **Ask before teaching**: "I noticed [specific observation]. There's a capability that removes exactly this — want two minutes on it?" Proceed only on yes.
- A "no" is recorded as `declined` for that capability — it will never be offered again unless they ask, or the reason demonstrably no longer applies (then at most once, saying why it's back).

On yes, continue as a growth-mode lesson (opener 4) for that capability.
