# devops
*Last verified: 2026-07-03*

**Hidden gem:** Worktree Isolation — rendering `terraform plan` or `helm template` in a disposable copy lets you evaluate risky infra changes with zero blast radius.

**Exemplar move:** Enter plan mode. Split monolithic main.tf (VPC, three subnets, security groups, RDS, ECS, ALB) into per-service modules — map every cross-resource reference, give a safe extraction order.

**Plugins:** `terraform` ☑️ IaC · `firebase` ☑️ · `linear`/`asana` ☑️ trackers · `session-report` ✅ usage reports — ~75 more vendor integrations (clouds, databases, observability, messaging) in the catalog; grep by vendor.

| # | Approach | Setup | Best when | Why it fits |
|---|----------|-------|-----------|-------------|
| 1 | [Plan Mode](../approaches/plan-mode.md) | Beginner | About to change a resource that other services depend on | Infra dependency graphs hide in string references; making them explicit prevents destroy-and-recreate surprises mid-apply |
| 2 | [Deep Research](../approaches/deep-research.md) | Beginner | Unfamiliar with a cloud provider's limits or pricing model | Provisioning mistakes are expensive to fix; research catches limits and pricing traps before they become incidents or bills |
| 3 | [Worktree Isolation](../approaches/worktree-isolation.md) | Intermediate | Want to test Terraform or Helm changes without touching your branch | Throwaway safety matters most for infrastructure — worst case is deleting a worktree, not filing an incident report |
