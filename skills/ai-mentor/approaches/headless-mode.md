# Headless Mode
*Last reviewed: 2026-06-27*

## What It Is

Headless Mode runs Claude without a terminal interface and without human interaction. You pass a prompt on the command line with the `-p` flag, Claude executes it, and outputs the result to stdout. There is no interactive session, no back-and-forth conversation — just input, execution, output. This is the bridge between "AI tool I use at my desk" and "AI step in my CI pipeline."

## Why It Works

Automation requires predictability and composability. Interactive tools are powerful for exploration, but they cannot be scheduled, piped, or embedded in scripts. Headless Mode turns Claude from an interactive assistant into a Unix-style command-line tool that follows the standard input/output contract. This means it composes with everything else in your toolchain: shell scripts, CI/CD pipelines, cron jobs, GitHub Actions, `xargs`, `jq`, and any other tool that speaks stdin/stdout. The same capability that helps you debug at your desk can now run unattended at 3 AM when a PR is opened.

## When to Use It

- Automated PR review: a GitHub Action runs Claude on every pull request to post a review comment
- Issue triage: Claude reads new issues, classifies them by component, and applies labels
- Scheduled code generation: a nightly job regenerates API client code from an updated OpenAPI spec
- Batch processing: run Claude on each file matching a pattern, like translating 50 i18n resource files

## When NOT to Use It

- Exploratory debugging where you need to steer the AI based on intermediate findings — headless mode gives you one shot, not a conversation
- Tasks that require subjective judgment calls you want to make yourself — headless mode will make those calls autonomously with no opportunity for you to intervene

## How It Works

### Basic (Beginner)

1. Run a simple prompt: `claude -p "Explain what src/auth/middleware.ts does"`
2. Claude reads the file, generates a response, and prints it to stdout
3. Get structured output: `claude -p "List all exported functions in src/api/" --output-format json`
4. The JSON output can be piped to `jq`, saved to a file, or consumed by another program
5. For scripts that need specific tools: `claude -p "Run the tests" --allowedTools "Bash(npm test)"` grants only the permissions your task needs

### Composing with Other Approaches (Intermediate)

- **Headless plus Autonomous Loops**: Run `claude -p "/goal all tests in tests/api/ pass"` in a CI job after a dependency update. The loop iterates until green, then the pipeline continues.
- **Headless in GitHub Actions**: Trigger Claude on `pull_request` events to review the diff, post comments, and optionally push fix commits — all without a human in the loop.
- **Headless for batch operations**: Use a shell loop to run Claude on multiple files: `for f in src/api/*.ts; do claude -p "Add JSDoc comments to all exported functions in $f" --allowedTools "Edit"; done`

### Advanced Patterns

- **`--bare` for reproducible runs**: The `--bare` flag tells Claude to ignore `CLAUDE.md` files, memory, and user preferences. This ensures the same input always produces the same behavior, which is critical for CI where reproducibility matters.
- **Session continuity**: Use `--continue` to resume the most recent session, or `--resume <session-id>` to pick up a specific session. This allows multi-step pipelines where step 1 runs Claude, a human reviews the output, and step 2 continues the same session with additional instructions.
- **Structured streaming**: Use `--output-format stream-json` to get a stream of JSON events as Claude works. Each event includes the type (text, tool_use, tool_result), allowing your pipeline to react to intermediate steps — for example, logging tool invocations in real time.
- **MCP server integration**: Headless Claude can connect to MCP servers with `--mcp-config`, giving it access to databases, APIs, or custom tools without any interactive setup: `claude -p "Query the staging database for users created today" --mcp-config mcp-servers.json`

## Tool Support

| Tool | Support | Notes |
|------|---------|-------|
| Claude Code | Native | `-p` flag, `--output-format`, `--allowedTools`, `--bare`, `--continue`, `--resume` |
| OpenCode | None | No headless or scriptable mode |
| Cursor | None | IDE-only; no CLI interface for scripting |
| aider | Partial | Can run non-interactively with `--message` flag, but limited output formatting |

## Common Pitfalls

- **Forgetting `--allowedTools`**: Without explicit tool permissions, headless Claude may lack the ability to run commands or edit files, and there is no human present to approve permission requests. Always specify which tools the headless run needs.
- **No human safety net**: In interactive mode, you see what Claude is doing and can stop it. In headless mode, it runs to completion. Scope permissions tightly and test your headless commands interactively first before putting them in a pipeline.
- **Ignoring exit codes**: Headless Claude returns an exit code. Check it in your scripts. A non-zero exit means something went wrong — a failed tool call, a prompt that could not be completed, or a permission denial. Do not pipe the output blindly into the next step.
- **Overloading a single prompt**: Headless mode works best with focused, single-purpose prompts. "Review this PR, fix any issues, update the changelog, and bump the version" should be four separate headless invocations, not one mega-prompt that might lose track of subtasks.

## Real-World Example

Your team wants automated PR reviews on every pull request. You create a GitHub Actions workflow:

```yaml
# .github/workflows/claude-review.yml
name: Claude PR Review
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Run Claude Review
        run: |
          claude -p "Review the changes in this PR. Focus on correctness bugs,
            security issues, and missing error handling. Post your findings as
            a PR comment." \
            --allowedTools "Bash(gh pr comment:*)" \
            --output-format json \
            --bare
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
```

Every time a PR is opened or updated, Claude checks out the code, reads the diff, and posts a review comment via `gh pr comment`. The `--bare` flag ensures no local `CLAUDE.md` or memory affects the review. `--allowedTools` restricts Claude to only posting comments — it cannot edit code, run tests, or access anything else. The `--output-format json` lets you parse the result in subsequent pipeline steps if needed.

The first week, the team notices Claude is flagging style issues they do not care about. They add a one-line instruction to the prompt: "Ignore formatting and style. Focus only on logic errors and security." The false-positive rate drops, and the reviews become a net time saver — catching a null-pointer dereference in `src/handlers/webhook.ts` that two human reviewers had missed.

## Sources

- [Claude Code Headless Mode](https://docs.anthropic.com/en/docs/claude-code/sdk/sdk-headless) — Official docs for running Claude Code non-interactively
- [Claude Code GitHub Actions](https://docs.anthropic.com/en/docs/claude-code/github-actions) — CI/CD integration with GitHub Actions
- [Claude Code CLI Reference](https://docs.anthropic.com/en/docs/claude-code/cli-reference) — CLI reference covering -p flag and non-interactive options
