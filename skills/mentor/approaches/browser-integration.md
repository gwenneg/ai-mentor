# Browser Integration
*Last verified: 2026-06-27*

## What It Is

Browser Integration lets your AI coding tool control a real web browser — navigating pages, clicking buttons, filling forms, reading the DOM, and checking the console for errors. Instead of guessing what the UI looks like from the source code alone, the AI can actually see and interact with your running application. There are two modes: a Chrome extension for web-app automation, and Computer Use for full screen-level control of any application.

## Why It Works

Code and its visual output are two different things, and bugs frequently hide in the gap between them. A developer debugging a CSS layout issue does not just read the stylesheet — they open the browser, inspect the element, and look at the computed styles. Browser Integration gives the AI the same ability. By interacting with the running application, the AI can verify its fixes immediately rather than reasoning abstractly about what a CSS change might look like. This closes the feedback loop between code change and visual result, which is particularly valuable for frontend work where the rendered output is the product.

## When to Use It

- Debugging CSS and layout issues where you need to see what the page actually renders, not just what the code says it should render
- Validating that a UI change works end-to-end — form submission, navigation, error states, responsive breakpoints
- Investigating console errors, network failures, or runtime exceptions that only appear in the browser
- Automating repetitive manual testing — filling multi-step forms, testing different user roles, checking accessibility

## When NOT to Use It

- Backend-only changes where there is no visual component — a database query fix does not need a browser
- When the application is not running locally or is not accessible from your machine — the browser needs a URL to navigate to
- For comprehensive E2E test suites — Browser Integration is for interactive debugging and ad-hoc verification, not a replacement for Playwright or Cypress test infrastructure
- Performance testing or load testing — the browser integration is single-session and not designed for concurrent load simulation

## How It Works

### Basic (Beginner)

1. Start your application locally (e.g., `npm run dev` on `localhost:3000`)
2. Install the Claude in Chrome extension from the Chrome Web Store and connect it to Claude Code
3. In Claude Code, reference the browser with `@browser`: "Navigate to localhost:3000/settings and check if the save button is disabled when no fields have changed"
4. Claude opens Chrome, navigates to the page, inspects the DOM, and reports what it finds
5. If there is a bug, Claude can edit the source code, wait for hot reload, and check the browser again to verify the fix
6. This edit-reload-verify cycle works without switching windows — Claude handles both the code side and the browser side

### Composing with Other Approaches (Intermediate)

- **Browser plus Autonomous Loop**: Set a goal like `/goal the signup form submits successfully and shows the confirmation page`, then let the AI iterate — editing code, checking the browser, fixing errors — until the flow works end-to-end.
- **Browser plus Plan Mode**: Use Plan Mode to have Claude inspect the running app, catalog all the visual issues on a page, and propose a prioritized fix plan before changing any code.
- **Browser for screenshot documentation**: After fixing a visual bug, ask Claude to screenshot the before and after states. Attach these to your pull request for reviewers who do not want to check out the branch locally.

### Advanced Patterns

- **Computer Use for native applications**: Computer Use mode reads the screen, moves the mouse, and types — not limited to the browser. Use it to test Electron apps, debug desktop UI toolkits, or interact with any application that has a graphical interface.
- **Console-driven debugging**: Ask Claude to open the browser console and monitor for errors while navigating the app. It can correlate a JavaScript `TypeError` in the console with the specific component and line of code, then fix it in one step.
- **Multi-viewport testing**: Ask Claude to check a component at different viewport widths: "Check the navigation menu at 1280px, 768px, and 375px widths and report any layout breaks." The AI resizes the viewport and inspects each state.

## Common Pitfalls

- **Forgetting to start the dev server**: The browser needs a running application to navigate to. If Claude reports "connection refused," your dev server is not running. Start it in a separate terminal before using browser features.
- **Authentication walls**: If your app requires login, Claude needs to authenticate first. Either provide test credentials explicitly, or navigate to a page that does not require auth. Do not assume the browser has an active session.
- **Relying on screenshots for pixel-precision**: Browser Integration is excellent for "is this element visible and functional" and poor for "is this exactly 16px from the left edge." For pixel-precise visual regression testing, use dedicated tools like Chromatic or Percy.
- **Heavy single-page applications**: If your SPA takes 10 seconds to hydrate, browser interactions may time out or report stale DOM state. Wait for the app to fully load before asking Claude to inspect it.

## Real-World Example

A user reports that the date picker in your scheduling form shows the wrong month when opened for the second time. You can reproduce it locally but cannot figure out why from the code alone.

```
claude
> @browser Navigate to localhost:3000/schedule/new. Open the date picker,
  select June 15, close it, then open it again. Tell me what month is
  showing and check the console for errors.
```

Claude opens Chrome, navigates to the page, clicks the date input, selects June 15, closes the picker, then reopens it. It reports: "The picker shows January 1970 on second open. Console shows: `Warning: Invalid time value` from `DatePicker.tsx:42`."

Claude reads `src/components/DatePicker.tsx` and finds that on line 42, `new Date(selectedDate)` is called where `selectedDate` is a Unix timestamp in seconds, but the `Date` constructor expects milliseconds. The first open works because the default is `new Date()` (current date), but after selection, the stored timestamp is interpreted as milliseconds, producing a date in January 1970.

Claude changes `new Date(selectedDate)` to `new Date(selectedDate * 1000)`, the hot reload triggers, and it reopens the date picker to confirm it now shows June correctly. It also checks the console — the `Invalid time value` warning is gone.

Total debugging time: under two minutes. Without the browser, this would have required reading the component code, mentally tracing the date flow through two state updates, guessing at the constructor behavior, and hoping the fix was correct — a process that took the original developer 45 minutes before they filed the bug report.

## Sources

- [Claude Code MCP](https://code.claude.com/docs/en/mcp) — Official docs for configuring MCP servers including Playwright
- [Playwright MCP Server](https://github.com/microsoft/playwright-mcp) — Official Microsoft Playwright MCP server repository
