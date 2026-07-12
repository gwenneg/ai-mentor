# devops
*Last verified: 2026-07-12*

**Hidden gem:** Worktree Isolation — rendering `terraform plan` or `helm template` in a disposable copy lets you evaluate risky infra changes with zero blast radius.

**Exemplar move:** Enter plan mode. Split monolithic main.tf (VPC, three subnets, security groups, RDS, ECS, ALB) into per-service modules — map every cross-resource reference, give a safe extraction order.

| # | Approach | Best when | Why it fits |
|---|----------|-----------|-------------|
| 1 | [Plan Mode](../approaches/techniques/plan-mode.md) | About to change a resource that other services depend on | Infra dependency graphs hide in string references; making them explicit prevents destroy-and-recreate surprises mid-apply |
| 2 | [Deep Research](../approaches/techniques/deep-research.md) | Unfamiliar with a cloud provider's limits or pricing model | Provisioning mistakes are expensive to fix; research catches limits and pricing traps before they become incidents or bills |
| 3 | [Worktree Isolation](../approaches/techniques/worktree-isolation.md) | Want to test Terraform or Helm changes without touching your branch | Throwaway safety matters most for infrastructure — worst case is deleting a worktree, not filing an incident report |
| 4 | [session-report](../approaches/tools/session-report.md) | Token usage and cache efficiency need real numbers | Optimization without measurement is guesswork — a usage report turns cost intuitions into data |
