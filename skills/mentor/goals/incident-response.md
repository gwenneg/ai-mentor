# Incident Response
*Last verified: 2026-06-28*

## When You're Here

Something is on fire in production. Maybe error rates are spiking, a service is completely down, data is coming back wrong, or your security team just flagged suspicious activity. Incident response overlaps with debugging, but the context changes everything: you're operating on live systems, under time pressure, and your users are actively affected. Every minute matters.

The instinct is to start reverting deploys or restarting services. AI doesn't replace your judgment here — it accelerates the data gathering, hypothesis testing, and coordination that eat up most of your incident timeline. The goal is to reduce mean time to resolution without introducing new damage along the way.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| Error spike correlating with a recent deploy | MCP context | Cross-correlate deploy timestamps with error rates across monitoring systems |
| Service down and you don't know where to start | Plan mode | Structured triage prevents thrashing — assess blast radius before acting |
| Outage symptoms match something you've seen in a dependency | Deep research | Check if this is a known upstream issue before debugging your own code |
| Need to query production state without manual SSH sessions | Headless mode | Non-interactive diagnostic scripts run safely and produce structured output |
| Applying a hotfix and worried about making it worse | Checkpoints & rewind | Checkpoint before the fix — instant rollback if things go sideways |

**Hidden gem:** Deep Research — two minutes checking whether it's a known upstream outage can save an hour of debugging code that was never the problem.

## Approaches (Ranked)

### 1. MCP Context — Pull from monitoring, logs, and alerting in real time
**Level:** Intermediate

During an incident, the critical bottleneck is context. You need error rates from Grafana, log traces from Datadog, recent deploys from your CD pipeline, and the alert history from PagerDuty — all at once. MCP integrations let Claude pull from these systems directly, cross-correlating data that would take you minutes to assemble manually across browser tabs.

**Try it now:**
> Connect to our Grafana and Datadog MCP servers. The `payment-service` is returning 503s to about 12% of requests starting ~15 minutes ago. Pull the error rate timeline from Grafana, the most recent deploy from our CD pipeline, and the top error messages from Datadog logs for `service:payment-service` in the last 30 minutes. Correlate the error spike with the deploy timestamp and tell me if they align.

**Why this works:** Incidents are data problems. The root cause is usually visible across multiple monitoring systems, but humans are slow at tab-switching and mental correlation. MCP turns Claude into a unified query layer across your entire observability stack, collapsing the data-gathering phase from minutes to seconds.

**Pros:**
- Queries multiple monitoring systems simultaneously
- Cross-correlates data that's hard to compare manually
- Surfaces patterns across deploys, traffic, and errors

**Cons:**
- Requires MCP servers configured for your monitoring stack
- Only as good as the data your monitoring captures
- Initial setup takes time — but pays off across every future incident

**Deeper:** See `approaches/mcp-context.md`

---

### 2. Deep Research — Check if this is a known upstream issue
**Level:** Beginner

Before spending an hour debugging your own code, spend two minutes checking whether your cloud provider is having an outage, whether your database vendor shipped a buggy release, or whether that npm package you depend on has a known issue matching your symptoms. Deep research fans out across status pages, GitHub issues, and community forums to find matches fast.

**Try it now:**
> /deep-research Our Redis cluster on AWS ElastiCache started throwing `CLUSTERDOWN` errors about 20 minutes ago. We haven't changed any Redis configuration or deployed anything that touches Redis. Is there a current AWS ElastiCache incident in us-east-1? Are there known issues with ElastiCache Redis 7.1 that match this error pattern? Check the AWS status page, Redis GitHub issues, and relevant community threads.

**Why this works:** A significant percentage of production incidents originate outside your codebase — in dependencies, cloud infrastructure, or third-party APIs. Searching for known issues first can save hours of debugging code that was never the problem.

**Pros:**
- Can resolve an incident in minutes if it's a known upstream issue
- Cross-references multiple sources for reliability
- Costs almost nothing to try first

**Cons:**
- Only helps when the problem is external to your code
- Results depend on how quickly upstream providers acknowledge issues

**Deeper:** See `approaches/deep-research.md`

---

### 3. Plan Mode — Structured triage before taking action
**Level:** Beginner

Under pressure, your instinct is to start reverting deploys, restarting pods, or scaling up resources. Plan mode forces you to pause and think: What's the blast radius? Which users are affected? What's the least-risky mitigation path? What should we not touch? This structured triage prevents the common anti-pattern of making an incident worse by applying the wrong fix to the wrong component.

**Try it now:**
> Enter plan mode. We have a production incident: the `order-service` is returning stale data — customers are seeing order statuses from hours ago. The service reads from a PostgreSQL replica and writes to the primary. Recent changes include a deploy 2 hours ago that modified the replication config in `infra/terraform/rds.tf` and a schema migration in `migrations/0087_add_order_metadata.sql`. Before I touch anything, help me triage: what's the blast radius, what are the possible causes ranked by likelihood, and what's the safest first action?

**Why this works:** Incident response has an asymmetry: a wrong action costs more than a delayed action. Plan mode forces a structured assessment that prevents you from reverting a deploy that wasn't the cause, or restarting a service that's masking a deeper issue. The few minutes spent planning consistently save more time than they cost.

**Pros:**
- Prevents knee-jerk reactions that worsen the incident
- Creates a triage record useful for postmortems
- Works with any tool — no setup required

**Cons:**
- Feels agonizingly slow when systems are down
- Requires discipline to follow the plan instead of freelancing

**Deeper:** See `approaches/plan-mode.md`

---

### 4. Headless Mode — Run diagnostic scripts against production
**Level:** Intermediate

Use `claude -p` with specific diagnostic prompts to query production state, analyze log files, or generate structured incident reports without an interactive session. This is ideal for running the same diagnostic across multiple services, piping production data through Claude for analysis, or integrating AI diagnostics into your incident runbooks.

**Try it now:**
> `claude -p "Analyze this production log excerpt and identify the root cause of the 503 errors. Group errors by type, show the timeline of when each error type started, and flag any errors that correlate with the deploy at 14:32 UTC." < /tmp/payment-service-logs-last-hour.txt`

**Why this works:** During incidents, you often need to run the same analysis repeatedly — on different services, different time windows, or different log sources. Headless mode makes these analyses scriptable and repeatable, turning ad-hoc debugging into a systematic diagnostic process.

**Pros:**
- Scriptable — can be embedded in incident runbooks
- Processes large log files or diagnostic output non-interactively
- Produces structured output suitable for incident channels

**Cons:**
- No interactive follow-up — each invocation is standalone
- Requires well-crafted prompts to get useful output
- Not suited for exploratory debugging where you need back-and-forth

**Deeper:** See `approaches/headless-mode.md`

---

### 5. Checkpoints & Rewind — Try a hotfix safely
**Level:** Beginner

When you've identified a probable cause and need to ship a hotfix fast, checkpoint your current state first. If the hotfix doesn't work — or makes things worse — you rewind instantly instead of scrambling to undo changes under pressure. This is especially valuable during incidents because the cost of a bad fix is amplified when systems are already degraded.

**Try it now:**
> Checkpoint here before I apply a hotfix. I think the 503s in `payment-service` are caused by a connection pool exhaustion bug introduced in the last deploy. I'm going to modify `src/db/connection_pool.py` to increase `max_connections` from 10 to 50 and add connection timeout handling. If the fix doesn't resolve the 503s after deploying to staging, I want to rewind and try a different approach.

**Why this works:** Incidents create pressure to move fast, and fast movement creates mistakes. Checkpoints make hotfix attempts zero-risk — you can try the most likely fix immediately, and if it's wrong, you're back to a clean state in seconds instead of trying to manually undo changes while your pager keeps firing.

**Pros:**
- Zero-risk hotfix attempts — rewind instantly if the fix is wrong
- Preserves conversation context so you don't lose your triage analysis
- No git stash juggling under pressure

**Cons:**
- Only covers local code changes — doesn't undo deployed changes
- Not a substitute for a proper rollback strategy in your CD pipeline

**Deeper:** See `approaches/checkpoints-rewind.md`

---

### 6. Subagent Delegation — Investigate multiple hypotheses in parallel
**Level:** Advanced

When an incident could be caused by the database layer, the API layer, or the infrastructure, investigating sequentially wastes precious minutes. Spawn subagents to check each hypothesis simultaneously: one agent queries database connection metrics, another analyzes API error patterns, and a third checks infrastructure health. You synthesize their findings to converge on the root cause faster.

**Try it now:**
> Spawn three investigation agents in parallel. Agent 1: check `src/db/` for connection pool leaks or slow query patterns — look at connection lifecycle in `connection_pool.py` and recent query changes in `repositories/`. Agent 2: analyze the API layer in `src/api/handlers/` for error handling regressions — focus on timeout handling and retry logic. Agent 3: review the infrastructure config in `infra/terraform/` and `k8s/deployments/` for recent resource limit changes or scaling policy modifications. Report back what each finds.

**Why this works:** Incidents have a critical property: investigation time is linear but hypotheses are independent. Subagents convert sequential investigation into parallel investigation, cutting the diagnosis phase by the number of hypotheses you can test simultaneously.

**Pros:**
- Investigates multiple root cause hypotheses at once
- Each agent works independently with full context
- Dramatically reduces time-to-diagnosis for complex incidents

**Cons:**
- Requires clear hypothesis framing to avoid agents duplicating effort
- Advanced feature — adds complexity during an already stressful situation
- Synthesis step still requires human judgment

**Deeper:** See `approaches/subagent-delegation.md`
