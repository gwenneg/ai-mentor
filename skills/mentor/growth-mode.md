# Growth mode

No problem given — this is "teach me something I don't know". Openers, in precedence order; take the first that applies and do only that one:

1. **Follow up.** The profile has a `shown` row from a past session → open with it: "Last time I showed you [X] — did it stick?" Their answer converts it to `adopted`, `declined`, or a re-teach from a different angle.
2. **Transfer.** The profile says `adopted`, but this repo's setup signals lack it (e.g. hooks everywhere else, none here) → offer the transfer: "You use [X] in your other projects — want the same here? Two minutes." This is configuration they already understand; set it up on acceptance.
3. **What's new.** The profile's `Last new-capability check` week is older than the newest rows in `processed-changelogs.md` → surface the most relevant workflow-visible change since, then update the anchor. Bootstrap and no-op rows are not news: when every row since the anchor is one, update the anchor and fall through to the next opener — never invent a change.
4. **The lesson.** Only now, having ruled out openers 1–3, build the full ignorance map: read the compiled index (`approaches/index.md`, one row per capability), subtract the profile and observed signals, and teach the top of what remains (every kind is an equal citizen — a lesson can be "you've never used /verify" or "pr-review-toolkit exists" just as well as a technique; a marketplace-directory plugin qualifies only when its stack/goal relevance to this user is concrete). Before teaching, read the winner's own file — `approaches/techniques/<id>.md` or `approaches/tools/<id>.md` — for the deep-dive. Hook the lesson to their observed work ("you do X by hand; this removes that"), name the principle in one sentence, offer the live demo, set it up on acceptance. One capability. Not two.

When the ignorance map is empty and nothing above applies, say so honestly — "you're using everything I'd recommend for how you work" — and offer the catalog list for browsing. Do not invent a lesson.

Leverage ranking for the map: a capability live in this repo but unconfirmed for this person first (the repo's own config is the demo — teach what it's already doing for them), then observed pain (something in this session it would fix), then fit to the repo and stack, then the general adoption ladder (project memory → plan mode → review skills → hooks → autonomous loops → subagents → the rest).

Record outcomes exactly as in problem mode's Record step: the lesson becomes `shown` with a one-line note, `adopted` on setup or "already use it", `declined` (reason verbatim) on a wave-off. Update `Last new-capability check` whenever opener 3 runs.

