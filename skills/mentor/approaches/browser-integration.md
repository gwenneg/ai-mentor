# Browser Integration
*Last verified: 2026-07-06*

## What It Is

Browser Integration connects Claude Code to your real Chrome or Edge browser through the Claude in Chrome extension — navigating pages, clicking buttons, filling forms, reading the DOM, and watching the console, in a visible window that shares your existing login state. Start a session with `claude --chrome` (or run `/chrome` mid-session) and describe browser tasks in plain language. A separate, heavier capability — computer use — extends the same idea beyond the browser to native macOS apps.

## Why It Works

Code and its visual output are two different things, and bugs frequently hide in the gap between them. A developer debugging a CSS layout issue does not just read the stylesheet — they open the browser, inspect the element, and look at the computed styles; Browser Integration gives the AI the same ability. By interacting with the running application, the AI can verify its fixes immediately rather than reasoning abstractly about what a CSS change might look like. This closes the feedback loop between code change and visual result, which is particularly valuable for frontend work where the rendered output is the product.

## When to Use It

- Debugging CSS and layout issues where you need to see what the page actually renders, not just what the code says it should render
- Validating that a UI change works end-to-end — form submission, navigation, error states, responsive breakpoints
- Investigating console errors, network failures, or runtime exceptions that only appear in the browser
- Working inside authenticated web apps (your CRM, Google Docs, an internal dashboard) — the browser shares your login state, so no API connector is needed

## When NOT to Use It

- Backend-only changes where there is no visual component — a database query fix does not need a browser
- Comprehensive E2E test suites or load testing — this is for interactive debugging and ad-hoc verification, not a replacement for Playwright or Cypress test infrastructure
- Unsupported environments — the extension supports Chrome and Edge only (not Brave, Arc, or WSL), and the integration requires a direct Anthropic plan (Pro, Max, Team, or Enterprise): it is unavailable through Amazon Bedrock, Google Cloud's Agent Platform (formerly Vertex AI), or Microsoft Foundry

## How It Works

### Basic (Beginner)

1. Start your application locally (e.g., `npm run dev` on `localhost:3000`)
2. Install the [Claude in Chrome extension](https://chromewebstore.google.com/detail/claude/fcoeoabgfenejglbffodgkkbkcdhcgfn) from the Chrome Web Store (works in Chrome and Edge)
3. Launch Claude Code with `claude --chrome`, or run `/chrome` in an existing session, then ask in plain language: "Open localhost:3000/settings and check whether the save button is disabled when no fields have changed"
4. Claude opens a new tab, navigates, inspects the DOM and console, and reports what it finds — pausing and handing control to you at login pages and CAPTCHAs. If there is a bug, Claude edits the source, waits for hot reload, and checks the browser again to verify the fix
5. Run `/chrome` at any time to check connection status, manage permissions, or reconnect; selecting "Enabled by default" removes the need for the flag, at the cost of browser tools always loaded in context

### Composing with Other Approaches (Intermediate)

- **Browser plus Autonomous Loops**: Set a goal like `/goal the signup form submits successfully and shows the confirmation page`, then let the AI iterate — editing code, checking the browser, fixing errors — until the flow works end-to-end.
- **Browser plus Plan Mode**: Have Claude inspect the running app, catalog all the visual issues on a page, and propose a prioritized fix plan before changing any code. Read-only browser calls (reading the page or console) run without permission prompts in plan mode; state-changing calls like clicks and navigation still ask.
- **Browser plus Worktree Isolation**: Check out a teammate's PR into a separate worktree, start its dev server on a second port, and have Claude verify the flow in the browser — recording it as a GIF (a built-in capability) for the PR thread — without touching your own working tree.

### Advanced Patterns

- **Computer use for native applications**: A separate built-in MCP server, `computer-use`, extends screen control past the browser — Claude can launch an Electron app, click through the iOS Simulator, or drive GUI-only tools, seeing the screen and moving the mouse. Enable it via `/mcp`; it is a research preview gated to macOS, Pro/Max plans, Claude Code v2.1.85+, and interactive sessions, with per-app approval each session. It is the slowest, broadest tool, so Claude reaches for it only when Bash, MCP, and the browser cannot do the job.
- **Console-driven debugging**: Ask Claude to watch the console while navigating the app, telling it which patterns to look for — logs are verbose. It can correlate a JavaScript `TypeError` with the specific component and line of code, then fix it in one step.
- **Multi-viewport testing**: Ask Claude to check a component at different widths: "Check the navigation menu at 1280px, 768px, and 375px widths and report any layout breaks." Window management is part of the browser tool set, so Claude resizes the window and screenshots each state.

## Common Pitfalls

- **Forgetting to start the dev server**: The browser needs a running application to navigate to. If Claude reports "connection refused," your dev server is not running. Start it in a separate terminal before using browser features.
- **Forgetting the browser acts as you**: Claude shares your browser's login state and can access any site you're signed into. That is the feature — but scope your asks accordingly, and manage site-level permissions in the Chrome extension settings for sites Claude should not touch.
- **Relying on screenshots for pixel-precision**: Browser Integration is excellent for "is this element visible and functional" and poor for "is this exactly 16px from the left edge." For pixel-precise visual regression testing, use dedicated tools like Chromatic or Percy.
- **Stale connection after idle periods**: The extension's service worker can go idle during long sessions, and browser commands stop responding. Run `/chrome` and select "Reconnect extension" rather than restarting everything.

## Real-World Example

A user reports that the date picker in your scheduling form shows the wrong month when opened for the second time. You can reproduce it locally but cannot figure out why from the code alone.

```
claude --chrome
> Open localhost:3000/schedule/new. Open the date picker, select June 15,
  close it, then open it again. Tell me what month is showing and check
  the console for errors.
```

Claude opens a tab, navigates to the page, clicks the date input, selects June 15, closes the picker, then reopens it. It reports: "The picker shows January 1970 on second open. Console shows: `Warning: Invalid time value` from `DatePicker.tsx:42`."

Claude reads `src/components/DatePicker.tsx` and finds that on line 42, `new Date(selectedDate)` is called where `selectedDate` is a Unix timestamp in seconds, but the `Date` constructor expects milliseconds. The first open works because the default is `new Date()` (current date), but after selection, the stored timestamp is interpreted as milliseconds, producing a date in January 1970.

Claude changes `new Date(selectedDate)` to `new Date(selectedDate * 1000)`, the hot reload triggers, and it reopens the date picker to confirm it now shows June correctly. It also checks the console — the `Invalid time value` warning is gone.

Total debugging time: under two minutes. Without the browser, this would have required reading the component code, mentally tracing the date flow through two state updates, guessing at the constructor behavior, and hoping the fix was correct — a process that took the original developer 45 minutes before they filed the bug report.

## Sources

- [Use Claude Code with Chrome](https://code.claude.com/docs/en/chrome) — Official docs: extension setup, `--chrome` and `/chrome`, capabilities, plan-mode permission behavior, troubleshooting
- [Let Claude use your computer from the CLI](https://code.claude.com/docs/en/computer-use) — Official docs for the built-in `computer-use` MCP server: enablement, platform and plan gates, safety model
