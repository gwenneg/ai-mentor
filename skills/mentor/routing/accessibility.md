# accessibility
*Last verified: 2026-07-03*

**Hidden gem:** Hooks — running the a11y scanner after every component edit catches regressions the moment they're introduced, not at audit time.

**Exemplar move:** Connect to the browser at localhost:3000/settings, record tab order; verify the Delete Account modal in src/components/Settings/DeleteAccountModal.tsx traps focus and Escape returns focus to the trigger.

**Plugins:** none mapped for this goal yet.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Browser Integration](../approaches/browser-integration.md) | involved | Verifying tab order, focus traps, or screen reader behavior | Tests the actual user experience — tab order, focus management, live regions — where the hardest a11y bugs hide |
| 2 | [Deep Research](../approaches/deep-research.md) | none | Building a new component, unsure which ARIA pattern applies | ARIA standards are precise but sprawling; one wrong attribute can make a component opaque to assistive technology |
| 3 | [Autonomous Loops](../approaches/autonomous-loops.md) | some | Scanner reported dozens of violations to fix | Scanner violations are objective and verifiable — ideal exit conditions for fix-scan-repeat loops with no judgment calls |
| 4 | [Hooks](../approaches/hooks-as-workflow.md) | some | Want every component edit automatically checked for a11y | A11y fixes have unintended side effects; instant scanner feedback catches regressions before they compound |
