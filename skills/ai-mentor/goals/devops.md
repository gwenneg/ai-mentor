# DevOps & Infrastructure
*Last reviewed: 2026-06-28*

## When You're Here

You're writing or modifying the definitions that describe your infrastructure: Terraform modules, Kubernetes manifests, Helm charts, Dockerfiles, Ansible playbooks, or cloud provider resource configurations. The work is precise and unforgiving — a typo in a security group rule can expose a database to the internet, and a misconfigured resource limit can take down a production cluster. Unlike application code where tests catch most mistakes before deploy, infrastructure changes often can't be fully validated until they hit a real environment.

This is distinct from CI/CD automation (which is about pipeline workflows). DevOps & Infrastructure is about the infrastructure definitions themselves — the Terraform that provisions your VPC, the Helm chart that deploys your microservice, the Ansible playbook that configures your servers.

## Quick Decision Guide

| Situation | Best approach | Why |
|-----------|--------------|-----|
| About to change a resource that other services depend on | Plan mode | Maps the blast radius before you touch anything |
| Unfamiliar with a cloud provider's limits or pricing model | Deep research | Surfaces gotchas and best practices before you provision |
| Want to test Terraform or Helm changes without affecting your branch | Worktree isolation | Run `terraform plan` or `helm template` in a disposable copy |
| Debugging a CrashLoopBackOff or Terraform apply error | Autonomous loops | Set the success condition and let Claude iterate through fixes |
| Need to cross-reference live cloud state with proposed changes | MCP context | Pull real-time infrastructure data into the conversation |

## Approaches (Ranked)

### 1. Plan Mode — Map infrastructure dependencies before changing anything
**Level:** Beginner | **Tools:** Any

Infrastructure changes cascade unpredictably. Renaming a security group can break five services that reference it. Changing a subnet CIDR forces recreation of every resource inside it. Plan mode forces you to map these dependencies before making the first edit — identify every resource that references the one you're changing, trace the dependency graph, and sequence the work so nothing breaks mid-apply.

**Try it now:**
> Enter plan mode. I need to split our monolithic `main.tf` into per-service Terraform modules. The file defines a VPC, three subnets, two security groups, an RDS instance, an ECS cluster, and an ALB. Map every cross-resource reference (e.g., which security groups are used by which services, which subnets the RDS and ECS resources live in). Then give me a safe extraction order — which resources can be moved to their own module first without breaking references to the others.

**Why this works:** Infrastructure dependency graphs are implicit — they live in string references, data source lookups, and variable pass-throughs that are easy to miss by reading individual files. Plan mode makes the graph explicit before you start refactoring, preventing the "I moved this resource and now `terraform plan` wants to destroy and recreate half my stack" problem.

**Pros:**
- Reveals implicit dependencies that grep alone misses
- Prevents destructive plan diffs from unexpected resource recreation
- Creates a sequenced work plan you can execute incrementally

**Cons:**
- Feels slow when you want to "just make the change" — but infrastructure mistakes are expensive
- Requires discipline to map before editing

**Deeper:** See `approaches/plan-mode.md`

---

### 2. Deep Research — Learn cloud provider best practices and gotchas
**Level:** Beginner | **Tools:** Claude Code

AWS, GCP, and Azure each have thousands of resource types, each with specific limits, pricing models, and undocumented behaviors. Deep research fans out across official documentation, community forums, and post-mortems to surface the gotchas before you provision. The difference between a $50/month and $5,000/month architecture is often a single configuration choice that only appears in a pricing FAQ.

**Try it now:**
> /deep-research We're deploying a new service on AWS ECS Fargate behind an ALB. The service handles WebSocket connections that can last up to 4 hours. Before I write the Terraform, I need to know: ALB idle timeout limits for WebSocket, Fargate task networking requirements for long-lived connections, NAT gateway costs for outbound traffic, and whether we should use ECS service connect or a service mesh. Also flag any hard limits that would force an architecture change.

**Why this works:** Infrastructure provisioning mistakes are expensive to fix — you can't just "refactor" a VPC CIDR range or undo a month of NAT gateway charges. Research before provisioning catches the limits, pricing traps, and architectural constraints that would otherwise surface as production incidents or surprise bills.

**Pros:**
- Surfaces hard limits and pricing gotchas before you commit to an architecture
- Cross-references official docs with real-world community experience
- Catches provider-specific best practices that differ from generic advice

**Cons:**
- Cloud provider docs change frequently — verify critical findings against the current console
- Cannot access your specific account quotas or pricing agreements

**Deeper:** See `approaches/deep-research.md`

---

### 3. Worktree Isolation — Test infrastructure changes safely
**Level:** Intermediate | **Tools:** Claude Code

Infrastructure changes are uniquely scary because `terraform apply` is irreversible in ways that `git revert` is not. Worktree isolation lets you try Terraform or Helm chart changes in a disposable copy of your repo: run `terraform plan` to see the execution plan, run `helm template` to render the manifests, diff the output against the current state — all without touching your working branch or risking an accidental apply from dirty state.

**Try it now:**
> Create a worktree from main. In it, refactor our Kubernetes deployment at `k8s/deployments/api-server.yaml` to add resource limits (256Mi memory, 250m CPU), a readiness probe on `/healthz`, and a PodDisruptionBudget with `minAvailable: 1`. Run `helm template ./charts/api-server` and diff the rendered output against the current version. Flag any issues — missing labels, selector mismatches, or values that would prevent a rolling update.

**Why this works:** The psychological safety of a throwaway environment is especially valuable for infrastructure. You'll try the risky refactor you've been avoiding when the worst case is deleting a worktree, not filing an incident report.

**Pros:**
- Zero risk to your working branch or current Terraform state
- Run `terraform plan` or `helm template` against proposed changes safely
- Evaluate the full diff before deciding whether to merge

**Cons:**
- Cannot test against real cloud APIs without configuring provider credentials in the worktree
- Only validates syntax and plan output, not actual deployment behavior

**Deeper:** See `approaches/worktree-isolation.md`

---

### 4. Autonomous Loops — Iterate until the deployment works
**Level:** Intermediate | **Tools:** Claude Code

Infrastructure configuration is full of tight feedback loops: `terraform plan` shows an error, you fix a reference, plan again, hit a new error, fix it, repeat. Autonomous loops excel here. Set "kubectl get pods shows all containers Ready" or "terraform plan shows no errors and no unexpected resource changes" as the success condition, and let Claude grind through the configuration issues — fixing provider version constraints, resolving circular dependencies, and correcting resource attribute types one by one.

**Try it now:**
> The Terraform configuration in `infra/environments/staging/` is failing on `terraform plan` with multiple errors after I upgraded the AWS provider from 4.x to 5.x. Run `terraform init -upgrade` and then `terraform plan`. For each error, fix the deprecated attribute or changed resource argument according to the AWS provider 5.x upgrade guide. Keep running `terraform plan` after each fix until the plan completes with zero errors and no unexpected resource destructions. Show me the final plan summary.

**Why this works:** Infrastructure debugging is convergent — each fix eliminates one error, and the tooling (`terraform plan`, `kubectl describe`, `helm lint`) tells you exactly what's wrong next. AI handles these mechanical iteration cycles without losing focus or patience.

**Pros:**
- Handles cascading configuration fixes without supervision
- Self-verifies using infrastructure tooling output after each change
- Especially effective for provider version upgrades with many small breaking changes

**Cons:**
- Can make fixes that pass validation but aren't idiomatic — review the diff carefully
- Needs clear, measurable success criteria to avoid spinning

**Deeper:** See `approaches/autonomous-loops.md`

---

### 5. MCP Context — Pull cloud state and monitoring into the conversation
**Level:** Intermediate | **Tools:** Claude Code/OpenCode with MCP

Infrastructure decisions improve dramatically when Claude can see your actual cloud state — not just the Terraform files, but the running resources, current costs, monitoring dashboards, and recent alerts. MCP servers for AWS, Kubernetes, and observability platforms let Claude cross-reference your proposed changes against reality: "you're adding a new NAT gateway, but your existing one is only at 12% utilization" or "this security group change would block traffic from the monitoring subnet."

**Try it now:**
> Connect to our Kubernetes cluster and run `kubectl get pods -A` to see current workload distribution. I'm planning to add node affinity rules to `k8s/deployments/` to separate our data-processing pods from API pods. Before I write the rules, show me which nodes each pod type currently runs on, what the resource utilization looks like per node, and whether we have enough capacity to enforce strict separation. Flag if any node would become overcommitted.

**Why this works:** Infrastructure-as-code is only half the picture — the other half is the actual running state. MCP context bridges this gap, letting Claude reason about proposed changes in the context of real utilization, costs, and traffic patterns instead of guessing from configuration files alone.

**Pros:**
- Grounds infrastructure decisions in real utilization and cost data
- Catches conflicts between proposed changes and current state
- Enables "what-if" reasoning with actual resource data

**Cons:**
- Requires MCP server setup and appropriate cloud credentials
- Live state queries add latency to the conversation
- Sensitive infrastructure data requires careful access scoping

**Deeper:** See `approaches/mcp-context.md`
