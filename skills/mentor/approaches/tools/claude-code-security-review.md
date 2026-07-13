---
kind: integration
last_verified: 2026-07-13
composes_with:
  - claude-code-action
  - built-in-review-skills
facts: "Anthropic's dedicated security GitHub Action (`anthropics/claude-code-security-review`): on every PR it runs the same diff-aware semantic analysis as the built-in /security-review and posts findings as inline comments on the specific lines, with fix recommendations. Requires an Anthropic API key secret enabled for both the Claude API and Claude Code. Its false-positive filter deliberately drops DoS, rate-limiting, resource-exhaustion, generic input-validation, and open-redirect findings."
session_signal: "a workflow under .github/workflows/ uses anthropics/claude-code-security-review"
source: https://github.com/anthropics/claude-code-security-review
pitfalls:
  - "Not hardened against prompt injection — only run it on trusted PRs; enable GitHub's 'Require approval for all external contributors' so workflows run only after maintainer review."
  - "Semantic diff review, not SCA: it does not check dependencies against CVE/advisory databases — pair with a dependency scanner for known-vulnerability detection."
  - "The README example pins `@main`; pin a released version for a reviewable supply chain."
---
