# Accessibility
*Last reviewed: 2026-07-02*

## When You're Here

You need to make your application usable by everyone, including people who navigate with keyboards, use screen readers, or rely on high contrast and magnification. Maybe an accessibility audit is coming up and you need to remediate findings. Maybe a scanner like axe or pa11y just dumped fifty violations on your dashboard. Maybe you're building a new component and want to get ARIA attributes and keyboard behavior right from the start. Or maybe a user filed a bug: "I can't tab to the submit button."

Accessibility work is uniquely suited to AI workflows because it combines two things that are hard for humans and easy for AI: exhaustive pattern-matching (finding every `<div onClick>` that should be a `<button>`) and standards lookup (knowing that a combobox needs `aria-expanded`, `aria-activedescendant`, and specific key handlers). The approaches below are ranked by how directly they catch real accessibility problems — because the gap between "passes a linter" and "works with a screen reader" is where most a11y bugs live.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Need to verify tab order, focus traps, or screen reader behavior | Browser integration | A11y is interactive — you have to navigate the page to find the bugs |
| Building a new component and unsure which ARIA pattern applies | Deep research | WCAG has specific patterns for each widget; the wrong one creates worse a11y |
| Scanner reported dozens of violations to fix | Autonomous loops | Mechanical fix-scan-fix cycles until zero violations remain |
| Reviewing a PR for accessibility regressions | Built-in review skills | Catches missing alt text, non-semantic HTML, and ARIA mistakes in the diff |
| Want every component edit automatically checked for a11y | Hooks | PostToolUse hook runs the scanner after each edit — instant feedback |

**Hidden gem:** Hooks — running the a11y scanner after every component edit catches regressions the moment they're introduced, not at audit time.

## Approaches (Ranked)

### 1. Browser Integration — Test with real assistive technology behavior
**Level:** Advanced

Accessibility bugs are interactive. A missing `aria-label` won't show up in a grep, but it will leave a screen reader user hearing "button" with no context. A focus trap that looks correct in code might let focus escape to the background when you actually press Tab. Browser integration lets Claude navigate your running app in a real browser — checking tab order, verifying that ARIA labels render correctly, confirming that modals trap focus, and testing that dynamic content announcements reach the accessibility tree.

**Try it now:**
> Connect to the browser and navigate to `localhost:3000/settings`. Tab through the entire page and record the focus order. Check whether the "Delete Account" modal at `src/components/Settings/DeleteAccountModal.tsx` traps focus when opened — Tab should cycle within the modal, not escape to the page behind it. Verify that the close button has an accessible name and that pressing Escape dismisses the modal and returns focus to the trigger button.

**Why this works:** WCAG compliance is necessary but not sufficient. A page can pass every automated rule and still be unusable with a keyboard. Browser integration tests the actual user experience — tab order, focus management, live region announcements — which is where the hardest accessibility bugs hide.

**Pros:**
- Catches focus management and keyboard navigation bugs that no static tool can find
- Tests real screen reader semantics via the accessibility tree
- Validates dynamic behavior: modals, dropdowns, toast notifications

**Cons:**
- Requires browser MCP setup and a running dev server
- Slower than static analysis — best used for targeted checks, not full-page sweeps
- Cannot fully replicate every screen reader's interpretation

**Deeper:** See `approaches/browser-integration.md`

---

### 2. Deep Research — Learn WCAG standards for your specific UI patterns
**Level:** Beginner

WCAG is not a single checklist — it is an extensive set of success criteria, and the correct implementation depends on the component. A combobox requires `aria-expanded`, `aria-activedescendant`, and a specific set of keyboard interactions. A data table with sortable columns needs `aria-sort` and live region announcements when the sort changes. Implementing the wrong ARIA pattern is often worse than no ARIA at all, because it gives screen readers false signals. Research the standard for your specific widget before writing a line of code.

**Try it now:**
> /deep-research I need to build an accessible autocomplete combobox for `src/components/Search/SearchInput.tsx`. What does the WAI-ARIA Authoring Practices combobox pattern require? Cover the required roles, states, and properties (aria-expanded, aria-autocomplete, aria-activedescendant), the full keyboard interaction model (arrow keys, Enter, Escape, Home, End), and how the listbox popup should be announced to screen readers. Include common mistakes that break screen reader compatibility.

**Why this works:** Accessibility standards are precise but sprawling. The WAI-ARIA Authoring Practices alone cover dozens of widget patterns, each with specific keyboard contracts and state management requirements. Getting one attribute wrong can make a component completely opaque to assistive technology. Research prevents building on wrong assumptions.

**Pros:**
- Surfaces the exact ARIA pattern for your component, not generic advice
- Covers keyboard interaction contracts that are easy to miss
- Identifies common implementation mistakes before you make them

**Cons:**
- Research alone does not fix anything — you still need to implement and test
- Standards can be ambiguous; different screen readers may interpret them differently

**Deeper:** See `approaches/deep-research.md`

---

### 3. Autonomous Loops — Fix violations until the scanner passes
**Level:** Intermediate

When axe or pa11y dumps a list of violations, the fix-scan-fix cycle is mechanical: add the missing alt text, re-scan, fix the next color contrast issue, re-scan, add the missing form label, re-scan. This is exactly the kind of work where autonomous loops excel. Set "zero violations" as the exit condition and let Claude grind through the list — each fix is verified by the scanner before moving to the next one.

**Try it now:**
> Run `npx pa11y http://localhost:3000/dashboard` and fix every reported violation. The main offenders are likely in `src/components/Dashboard/MetricsCard.tsx` (missing alt text on icons), `src/components/Dashboard/FilterBar.tsx` (form inputs without labels), and `src/components/shared/Badge.tsx` (insufficient color contrast). After each fix, re-run pa11y to confirm the violation count decreases. Stop when zero violations remain.

**Why this works:** Accessibility scanner violations are objective and verifiable — either the violation is present or it is not. This makes them ideal exit conditions for autonomous loops. Claude does not need to make judgment calls; it fixes, scans, and repeats until the tool reports success.

**Pros:**
- Clears scanner backlogs that would take hours of manual work
- Every fix is verified by the scanner before moving on
- Handles mechanical fixes (alt text, labels, contrast) without supervision

**Cons:**
- Scanners only catch ~30-40% of WCAG issues — passing a scan does not mean accessible
- May apply superficial fixes (empty alt text, hidden labels) that satisfy the tool but not users
- Needs a running dev server for the scanner to test against

**Deeper:** See `approaches/autonomous-loops.md`

---

### 4. Built-In Review Skills — Catch a11y issues in code review
**Level:** Beginner

The `/code-review` skill can catch accessibility regressions in a PR diff before they reach production: images without `alt` attributes, `<div>` elements with click handlers instead of `<button>`, color used as the only state indicator, missing `aria-label` on icon-only buttons, and non-semantic heading hierarchies. It will not catch everything a scanner or manual test would find, but it catches the common patterns that account for the majority of introduced violations.

**Try it now:**
> Run `/code-review` on the current branch with an accessibility focus. Pay attention to the new components in `src/components/Onboarding/` — check for missing alt text on images, interactive elements that are not keyboard-focusable, form inputs without associated labels, and any use of `tabindex` values greater than 0. Flag any `<div onClick>` or `<span onClick>` that should be a `<button>`.

**Why this works:** Most accessibility violations are introduced one PR at a time — a developer adds an image and forgets `alt`, wraps a `<div>` in an `onClick` instead of using a `<button>`, or uses color alone to indicate an error state. Catching these at review time is vastly cheaper than catching them in an audit.

**Pros:**
- Zero setup — works immediately on any branch
- Catches the most common violation patterns in the diff
- Integrates into existing code review workflow

**Cons:**
- Limited to the current diff — will not audit unchanged code
- Cannot verify interactive behavior like keyboard navigation or focus management
- May miss project-specific accessibility patterns

**Deeper:** See `approaches/built-in-review-skills.md`

---

### 5. Hooks — Auto-run a11y checks after every component edit
**Level:** Intermediate

A PostToolUse hook runs your accessibility scanner automatically after every file edit, giving Claude instant feedback on whether the change introduced or fixed a violation. This creates a tight feedback loop: Claude edits a component, the hook runs axe or pa11y on the affected page, Claude sees the results, and adjusts its next edit accordingly. Instead of batching fixes and checking later, every single change is validated against the accessibility standard in real time.

**Try it now:**
> Set up a PostToolUse hook that runs `npx axe http://localhost:3000/components/preview` after every edit to files in `src/components/`. Then refactor `src/components/Navigation/Sidebar.tsx` to use semantic `<nav>` and `<ul>` elements instead of nested `<div>` tags, add `aria-current="page"` to the active link, and ensure the entire sidebar is keyboard-navigable. The hook will confirm each change reduces violations.

**Why this works:** Accessibility fixes often have unintended side effects — adding an `aria-label` to one element can change how a screen reader interprets its children. Instant scanner feedback after every edit catches these regressions immediately, before they compound into a confusing debugging session.

**Pros:**
- Every edit is validated against a11y rules automatically
- Catches regressions the moment they are introduced
- Creates a tight feedback loop that accelerates fix iteration

**Cons:**
- Scanner execution adds latency to every edit cycle — keep the scope narrow
- Only catches violations the scanner can detect — manual testing still needed for UX
- Requires a running dev server for the scanner to test against

**Deeper:** See `approaches/hooks-as-workflow.md`
