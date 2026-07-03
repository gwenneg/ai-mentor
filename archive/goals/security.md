# Security Hardening
*Last verified: 2026-06-28*

## When You're Here

You need to proactively strengthen your application's security posture. Maybe a penetration test is scheduled next month, a CVE just dropped for a dependency you rely on, or a compliance requirement landed on your sprint. This isn't general code review — it's targeted work: auditing for injection vectors, verifying auth flows, scanning for exposed secrets, checking dependency chains, and ensuring security headers and encryption are configured correctly.

Security hardening is uniquely suited to AI workflows because the work is systematic and exhaustive. A human reviewer might check the three auth endpoints they remember, but AI will check every route. The challenge isn't creativity — it's coverage. The approaches below are ranked by how quickly they get you from "I need to harden this" to "I have findings I can act on."

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Quick security scan before a release or audit | Built-in review skills | One command, immediate vulnerability report |
| Large codebase with multiple vulnerability classes to check | Subagent delegation | Parallel specialists cover injection, auth, secrets simultaneously |
| New CVE announced for a dependency in your stack | Deep research | Research the vulnerability landscape before patching |
| Recurring security patterns unique to your project | Custom agents | Encode your auth middleware, ORM quirks, and sensitive fields once |
| Need to check against compliance frameworks or threat models | MCP context | Pull OWASP checklists, SOC 2 controls, or internal security policies |
| Protect security-critical files from accidental modification | Hooks | Block edits to auth config, encryption keys, security headers without review |

**Hidden gem:** Hooks — a PreToolUse guard on auth configs and crypto files prevents the accidental security regressions that no scanner catches.

## Approaches (Ranked)

### 1. Built-In Review Skills — First line of defense, zero setup
**Level:** Beginner

The `/security-review` skill analyzes your current branch diff for common vulnerability classes: injection, auth bypass, secrets exposure, insecure deserialization, and more. It runs immediately with no configuration and produces a severity-ranked report. For a quick pre-audit check or a sanity scan before merging a sensitive PR, this is where you start.

**Try it now:**
> Run `/security-review` on the current branch. Pay special attention to the auth middleware in `src/middleware/auth.ts` and any raw database queries in `src/services/`. I'm preparing for a security audit next week and need to know what's exposed.

**Why this works:** Security vulnerabilities follow well-known patterns — SQL injection, XSS, CSRF, hardcoded secrets. Built-in review skills encode these patterns and apply them exhaustively to every changed line, catching issues that slip past manual review.

**Pros:**
- Immediate results with zero configuration
- Covers OWASP Top 10 categories out of the box
- Consistent — checks every file, not just the ones you remember

**Cons:**
- Limited to the diff on the current branch — won't audit unchanged legacy code
- May miss project-specific vulnerability patterns unique to your architecture

**Deeper:** See `approaches/built-in-review-skills.md`

---

### 2. Subagent Delegation — Parallel reviewers by vulnerability class
**Level:** Advanced

A single security pass tries to hold injection, authentication, authorization, secrets, and cryptography in its head simultaneously. Subagent delegation splits these into focused reviewers running in parallel. One agent audits every database query for injection. Another traces every auth flow for bypass opportunities. A third scans for hardcoded secrets and leaked credentials. Specialized focus catches what a general sweep misses.

**Try it now:**
> Review the `src/` directory using four parallel security subagents: (1) an injection reviewer checking all database queries in `src/repositories/` for SQL injection, especially any raw queries bypassing the ORM, (2) an auth reviewer tracing every route in `src/api/routes/` to verify `requireAuth` middleware is applied and role checks are correct, (3) a secrets reviewer scanning the entire repo for hardcoded API keys, connection strings, and tokens, (4) a cryptography reviewer checking `src/utils/crypto.ts` for weak algorithms, improper IV usage, and insecure key storage. Consolidate findings by severity.

**Why this works:** Security auditing is inherently multi-dimensional. Injection expertise is different from cryptography expertise. By giving each agent a single concern, you get deeper analysis in each area without the attention dilution of a monolithic scan.

**Pros:**
- Deeper coverage across multiple vulnerability classes simultaneously
- Each reviewer maintains focus on its specialty
- Findings arrive in parallel, reducing total audit time

**Cons:**
- Requires careful scoping of each agent's responsibility to avoid gaps
- More token-intensive than a single-pass review

**Deeper:** See `approaches/subagent-delegation.md`

---

### 3. Deep Research — Know the threat landscape before you harden
**Level:** Beginner

Before patching, understand what you're patching against. Deep research fans out across CVE databases, security advisories, GitHub issues, and best-practice guides to build a picture of the vulnerabilities relevant to your specific stack. When a new CVE drops or a compliance requirement arrives, research first — then harden with confidence.

**Try it now:**
> /deep-research CVE-2024-21626 affects our container runtime. We run Node.js services on `containerd` 1.6.x in production. What are the actual exploitation conditions, which versions are patched, and what mitigations should we apply if we can't upgrade immediately? Also check if any of our base images in `Dockerfile` and `docker-compose.yml` are affected.

**Why this works:** Security hardening without context is guesswork. A CVE with a CVSS score of 9.8 might be irrelevant to your deployment configuration. Deep research gives you the specifics — affected versions, exploitation prerequisites, available patches — so you harden what matters instead of everything.

**Pros:**
- Surfaces relevant CVEs, advisories, and patches for your exact stack
- Cross-references multiple sources to avoid acting on incomplete information
- Saves hours of manual advisory reading and version-matching

**Cons:**
- Only as current as publicly available advisory databases
- Research alone doesn't fix anything — you still need to apply findings

**Deeper:** See `approaches/deep-research.md`

---

### 4. Custom Agents — Your project's security reviewer, permanently defined
**Level:** Advanced

Every project has security patterns that generic tools don't know about: your specific ORM and where raw queries are allowed, your auth middleware chain and which routes are intentionally public, your sensitive data fields that must never appear in logs. A custom agent in `.claude/agents/` encodes this knowledge once and applies it consistently on every review.

**Try it now:**
> Create a custom security agent at `.claude/agents/security-auditor.md`. It should know: (1) we use Prisma — flag any `$queryRaw` or `$executeRaw` usage, (2) all routes in `src/api/` must use `authMiddleware` except `/health` and `/metrics`, (3) PII fields (`email`, `ssn`, `dateOfBirth`) must never appear in log statements or error responses, (4) JWT tokens must be verified with `RS256` — flag any `HS256` usage, (5) all user input in `req.body` and `req.params` must pass through `src/validation/schemas/` before use. Then run it against the full `src/` directory.

**Why this works:** Generic security scanners produce noise because they don't know your codebase's intentional patterns. A custom agent that knows your auth middleware is called `authMiddleware`, your ORM is Prisma, and your PII fields have specific names produces high-signal findings with fewer false positives.

**Pros:**
- Encodes project-specific security knowledge permanently
- Dramatically reduces false positives compared to generic scanning
- New team members get the same security review quality as veterans

**Cons:**
- Agent definitions require maintenance as security patterns evolve
- Initial setup takes time to enumerate all project-specific rules

**Deeper:** See `approaches/custom-agents.md`

---

### 5. MCP Context — Pull compliance requirements and threat models into the review
**Level:** Intermediate

Security hardening doesn't happen in a vacuum. You're hardening against specific threats documented in threat models, and toward specific controls required by compliance frameworks. MCP context connects Claude to your security tools — Snyk for dependency vulnerabilities, your compliance tracker for SOC 2 controls, your internal wiki for threat models — so findings map directly to requirements.

**Try it now:**
> Connect to our Confluence workspace and pull the threat model document for the payments service. Then audit `src/services/payment/` against each identified threat. For each threat in the model, tell me whether our current implementation mitigates it, partially addresses it, or leaves it unmitigated. Cross-reference with our SOC 2 control spreadsheet in Google Sheets to flag any compliance gaps.

**Why this works:** Auditors don't ask "is your code secure?" — they ask "do you meet control AC-3.2?" MCP context lets you review code against the specific controls and threats that matter for your compliance posture, producing findings that map directly to audit requirements.

**Pros:**
- Findings map directly to compliance controls and threat model entries
- Connects code review to organizational security requirements
- Reduces the gap between engineering work and audit evidence

**Cons:**
- Requires MCP server setup for your security and compliance tools
- Compliance documents must be accessible and well-structured

**Deeper:** See `approaches/mcp-context.md`

---

### 6. Hooks — Guard security-critical files from unreviewed changes
**Level:** Intermediate

Some files are too sensitive for casual edits: auth middleware, encryption configurations, security headers, CORS policies, CSP directives. A PreToolUse hook intercepts edits to these files and blocks them unless explicitly approved, creating a guardrail that prevents accidental weakening of security controls during routine development.

**Try it now:**
> Set up a PreToolUse hook that triggers when any file matching `src/middleware/auth*`, `src/config/security*`, `src/utils/crypto*`, or `*.env*` is about to be edited. The hook should print a warning: "This file is security-critical. Confirm the edit maintains existing security guarantees." and require explicit approval before proceeding.

**Why this works:** Most security regressions aren't malicious — they're accidental. A developer refactoring middleware removes a CSRF check they didn't notice. A config change loosens CORS without realizing it. Hooks create a speed bump that forces conscious acknowledgment before touching security-critical code.

**Pros:**
- Prevents accidental weakening of security controls
- Zero ongoing effort after initial setup
- Works as a teaching tool — developers learn which files are security-sensitive

**Cons:**
- Only protects files you think to guard — new security-critical files need manual addition
- Can become annoying if the file list is too broad

**Deeper:** See `approaches/hooks-as-workflow.md`

---

### 7. Scheduled & Recurring Agents — Security review that never skips a week
**Level:** Intermediate

Security posture decays between audits. A routine with two triggers keeps it continuous: a weekly scheduled run audits the highest-risk directories against your checklist, and a GitHub trigger runs a security-focused review on every PR that touches auth paths (filter: head branch or title contains `auth`). Findings arrive as session transcripts and PR comments — no one has to remember to look.

**Try it now:**
> /schedule weekly on Monday at 7am: audit src/api/ and src/auth/ against our security checklist — raw SQL usage, routes missing authMiddleware, secrets in code, weak crypto parameters. Open an issue titled "Weekly security audit YYYY-MM-DD" listing findings by severity with file paths.

**Why this works:** Attackers don't wait for your audit cadence; scheduled review converts security from a quarterly event into a standing property of the codebase.

**Pros:**
- Combines calendar cadence with event triggers (every auth-touching PR)
- Findings land as issues/PR comments inside your existing workflow
- Runs even during crunch weeks when manual audits get skipped

**Cons:**
- A routine's checklist goes stale — review its prompt when your threat model changes
- Autonomous runs can't ask clarifying questions; expect some noise to triage

**Deeper:** See `approaches/scheduled-agents.md`
