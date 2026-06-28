# Multi-Provider Model Selection
*Last reviewed: 2026-06-27*

## What It Is

Multi-Provider Model Selection lets you switch between different AI models during your coding workflow based on what each task needs. Instead of being locked to a single AI provider, you choose the best model for each job: a strong reasoning model for architecture decisions, a fast model for bulk code generation, or a local model for privacy-sensitive work. You make this choice per task, not per project, and switching takes seconds.

## Why It Works

No single AI model is best at everything. Large reasoning models excel at complex architecture and debugging but are slow and expensive for routine code generation. Small, fast models handle boilerplate and repetitive edits efficiently but struggle with multi-file reasoning. Local models keep sensitive code off external servers but may lack the capability for advanced tasks. Multi-provider selection applies the engineering principle of using the right tool for the job — the same way you would choose between a profiler, a debugger, and a linter depending on what you are investigating.

## When to Use It

- Optimizing costs by using cheaper models for straightforward tasks and reserving expensive models for complex reasoning
- Working on privacy-sensitive code (credentials management, proprietary algorithms) where data must not leave your machine
- Leveraging specific model strengths: one model for refactoring Go code, another for generating SQL queries, another for writing documentation
- Teams with mixed subscriptions where some developers have access to certain providers and not others

## When NOT to Use It

- When your organization has standardized on a single provider for compliance or governance reasons
- When you are just getting started with AI coding tools — master one model first before optimizing across multiple
- When the task is simple enough that model choice does not meaningfully affect the result

## How It Works

### Basic (Beginner)

1. Configure your providers in OpenCode's configuration file (`opencode.json` or via environment variables) with API keys for each provider you want to use.
2. Start OpenCode in your project directory.
3. Begin working with the default model. When you hit a task that would benefit from a different model, switch using the model selector keybinding.
4. Select the model suited to your current task: Claude for complex analysis, DeepSeek for fast generation, a local Ollama model for private code.
5. Continue working. The conversation context carries over — the new model sees what the previous model discussed.

### Composing with Other Approaches (Intermediate)

- **Model switching plus iterative refinement**: Use a fast, cheap model to generate a first draft of a module, then switch to a strong reasoning model to review and refine it. The fast model handles the volume, the reasoning model handles the quality.
- **Local model for exploration, cloud model for finalization**: Use a local model via Ollama to rapidly iterate on an approach without API costs or latency concerns. Once you are confident in the direction, switch to a cloud model for the final, polished implementation.
- **Provider-specific strengths for different file types**: Use one model for backend logic where it excels at type reasoning, and switch to another for frontend components where it produces better UI code. Match model strengths to the domain of each file.

### Advanced Patterns

- **Cost-aware workflow budgeting**: Track your API spend across providers and shift to cheaper models as you approach budget limits. Use the expensive model only for the tasks where quality meaningfully differs — architecture, debugging, and security review.
- **GitHub Copilot integration for subscription users**: If your team already pays for GitHub Copilot, OpenCode can use Copilot as a provider. This lets you leverage an existing subscription without additional API costs for routine tasks, reserving direct API access for tasks that need it.
- **Local-first for compliance, cloud for capability**: In regulated environments, default to local models so code never leaves the machine. Escalate to cloud models only for tasks where local capability is insufficient, with explicit developer approval for each escalation.

## Tool Support

| Tool | Support | Notes |
|------|---------|-------|
| Claude Code | Partial | Non-Anthropic models supported via API gateways (LiteLLM, OpenRouter); set ANTHROPIC_BASE_URL to gateway endpoint |
| OpenCode | Native | 75+ providers including Claude, GPT, DeepSeek, Gemini, Groq, Ollama; per-task switching. Note: Anthropic restricted OpenCode from Claude OAuth access in January 2026 — API key access still works but subscription-based access does not |
| Cursor | Native | Built-in support for OpenAI, Anthropic, Google, xAI, and Cursor models with seamless switching |
| aider | Native | Supports many providers; model can be set per session via `--model` flag |

## Common Pitfalls

- **Switching models too frequently**: Every model switch has a cognitive cost — you need to recalibrate your expectations for the new model's style and capability. Switch when the task genuinely demands it, not out of curiosity.
- **Assuming all models handle context the same way**: Different models have different context window sizes, different strengths with long context, and different tendencies with code style. A prompt that works well with Claude may need adjustment for DeepSeek or GPT.
- **Ignoring the cost of local model quality gaps**: Local models are free and private, but if they produce code that requires significant manual correction, the time cost may exceed the API cost of a cloud model. Measure actual productivity, not just API bills.
- **Not validating model output consistency**: If two models contribute to the same file in the same session, the code style may be inconsistent. Run your linter and formatter after model switches to normalize the output.

## Real-World Example

You are building a new GraphQL API layer for an existing REST backend. The work has three distinct phases, each suited to a different model.

First, you need to design the schema. You start OpenCode with Claude, which excels at architectural reasoning:

```
> Analyze the REST endpoints in src/api/routes/ and propose a GraphQL schema
  that consolidates the N+1 query patterns in the order and inventory endpoints.
```

Claude produces a schema in `src/graphql/schema.graphql` with thoughtful type relationships and resolver structure. This took careful reasoning about data dependencies.

Second, you need to generate 14 resolver files — mostly mechanical translation from REST handlers. You switch to DeepSeek, which is fast and cheap for bulk generation:

```
> Generate resolvers for each type in schema.graphql. Follow the pattern
  in src/graphql/resolvers/user.ts as a template. Map each resolver to
  the corresponding REST client in src/api/clients/.
```

DeepSeek generates all 14 files in under a minute at a fraction of the cost.

Third, you need to write data-loader logic to solve the N+1 queries, and this code handles sensitive customer data. You switch to a local Llama model running via Ollama:

```
> Implement DataLoader batching for the order and inventory resolvers.
  The batch keys are customer IDs — this data must not leave the machine.
```

The local model produces working data-loader code. It is not as elegant as what Claude would write, but it is correct, and no customer identifiers were sent to an external API. You run `npm test` — all 31 new tests pass. Total API cost for the session: about 40% of what it would have been using Claude for everything.

## Sources

- [Claude Code Bedrock & Vertex Proxies](https://docs.anthropic.com/en/docs/claude-code/bedrock-vertex-proxies) — Enterprise deployment for Bedrock, Vertex AI, and Foundry providers
- [Claude Code on Amazon Bedrock](https://docs.anthropic.com/en/docs/claude-code/amazon-bedrock) — Dedicated Bedrock setup guide
- [Claude Code Setup](https://docs.anthropic.com/en/docs/claude-code/setup) — Setup guide covering provider environment variables
