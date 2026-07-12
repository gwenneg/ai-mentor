# LSP Self-Correction
*Last verified: 2026-07-12*

## What It Is

LSP Self-Correction feeds real-time compiler and type-checker diagnostics directly into the AI agent's editing loop. When the AI generates code that has a type error, a missing import, or an undefined variable, the Language Server Protocol catches the mistake immediately and feeds it back to the agent, which fixes it in the same turn — automatically, without you intervening. The result is code that compiles on the first try far more often than code generated without this feedback loop.

## Why It Works

Inserting the compiler into the generation loop replaces "generate, copy, paste, compile, read error, re-prompt" with "generate, diagnose, fix" — and it happens inside a single AI turn.

## When to Use It

- Generating code in strongly typed languages (TypeScript, Go, Rust, Java) where type errors are common in AI output
- Writing tests that need to match exact function signatures, return types, and error types
- Refactoring across multiple files where a rename or type change can cascade into dozens of compile errors
- Working in unfamiliar codebases where the AI is likely to guess wrong about types and interfaces
- Tracing how components connect across a large codebase — the same plugins give the agent go-to-definition, find-references, and call hierarchies, which beat text search on indirect calls and aliased imports

## When NOT to Use It

- In untyped or sparsely annotated codebases where the language server has little to report — official plugins exist for Python (`pyright-lsp`), Ruby, and PHP, but without type annotations the loop catches far less than in TypeScript, Go, or Rust
- When the LSP server itself is misconfigured or produces excessive false positives, which would send the agent into pointless fix loops
- For quick questions, code explanations, or tasks that do not involve generating or editing code — there are no diagnostics to feed back if no code is being written

## How It Works

### Basic (Beginner)

1. Install the code intelligence (LSP) plugin for your language from the official marketplace — e.g. `/plugin install gopls-lsp@claude-plugins-official` for Go, or `typescript-lsp`, `pyright-lsp`, `rust-analyzer-lsp`, and others for their languages. The language server binary itself must be installed on your machine.
2. Ask the agent to generate or modify code: "Add a `CalculateDiscount` method to the `Order` struct that returns the discounted total."
3. The agent writes the code. The LSP immediately analyzes it and reports any diagnostics — type mismatches, missing imports, undefined references.
4. If there are errors, the agent sees them in its context and generates a fix in the same response cycle. This may repeat for two or three iterations until the diagnostics are clean.
5. You receive code that already compiles. No manual error-fixing step needed.

### Composing with Other Approaches (Intermediate)

- **LSP Self-Correction plus Autonomous Loops**: Ask the agent to write tests for a module and iterate until the suite passes. The LSP catches type errors in test assertions (wrong return type, missing struct field), and the agent corrects them before you ever run the test suite. When you do run the tests, failures are logic errors, not type errors.
- **LSP Self-Correction plus Plan Mode**: Plan a cross-file rename or signature change in plan mode, then implement — the LSP reports every downstream compile error across the project. The agent walks through each diagnostic and fixes them in sequence, handling cascade effects you might miss manually.
- **LSP Self-Correction plus Built-In Review Skills**: After the agent finishes a task, check for zero LSP diagnostics, then run `/code-review` for logic-level issues. This is a lightweight verification step that catches a category of errors without running the full test suite.

### Advanced Patterns

- **Diagnostic-driven exploration**: When working in an unfamiliar codebase, intentionally write code that you expect to fail type-checking. The LSP errors reveal the actual types, interfaces, and constraints the codebase uses. The agent learns the codebase's type landscape through the errors and produces correct code on the next attempt.
- **Layered feedback loops**: Combine LSP diagnostics (type errors) with linter output (style violations) and test results (logic errors) in a single agent loop. Each layer catches a different category of mistake, and the agent addresses all three before presenting the final result.
- **Suppressing noisy diagnostics**: If the LSP produces warnings that are not actionable (deprecated API notices, optional lint suggestions), configure the LSP server to suppress them. Noisy diagnostics waste agent context and can cause unnecessary edits.

## Common Pitfalls

- **Infinite correction loops**: If the LSP reports an error that the agent cannot fix (a misconfigured build, a missing dependency), the agent may loop: fix, re-diagnose, "fix" again, re-diagnose. Set a maximum iteration count or interrupt when you see the same error repeated three times.
- **Over-trusting zero diagnostics**: Zero LSP errors means the code compiles, not that it is correct. Type-safe code can still have logic errors, performance problems, and security vulnerabilities. LSP Self-Correction handles syntactic and type-level correctness; you still need tests for behavioral correctness.
- **Missing LSP server installation**: The self-correction loop only works if the LSP server for your language is installed and reachable. If you are working in Rust but `rust-analyzer` is not installed, there are no diagnostics to feed back. Verify your LSP setup if you notice the agent producing code with obvious type errors.
- **Conflicting LSP configurations**: Project-level LSP settings (like `tsconfig.json` paths or Go build tags) can cause the LSP to report errors that do not match your actual build. Ensure the LSP configuration matches your build system.

## Sources

- [Best practices for Claude Code](https://code.claude.com/docs/en/best-practices) — official guide on giving Claude checks it can run; recommends code intelligence plugins for typed languages
- [Discover and install prebuilt plugins](https://code.claude.com/docs/en/discover-plugins#code-intelligence) — official reference for code intelligence (LSP) plugins: plugin names, required binaries, automatic diagnostics after edits, and code navigation

## Signals

- Setup: An LSP plugin installed for the project language
- Session: Mentions go-to-definition / diagnostics-driven fixes
