# Channels
*Last verified: 2026-07-06*

## What It Is

Channels push events from outside systems into your already-running Claude Code session — chat messages, CI results, monitoring alerts, webhooks — so Claude can react to things that happen while you're not at the terminal. A channel is an MCP server with push capability, installed as a plugin and opted in per session with `claude --channels plugin:<name>@claude-plugins-official`; channels can be two-way, with Claude reading the inbound event and replying through the same channel like a chat bridge. Telegram, Discord, and iMessage plugins ship in the research preview, plus a localhost `fakechat` demo, and you can build your own for systems without one.

## Why It Works

The event arrives in the session that already has your files open and remembers what you were debugging — context a cold-started agent would have to reconstruct, and the reply path makes your phone a steering wheel for work executing on your real machine.

## When to Use It

- Reacting to external events mid-task: CI failures, error-tracker alerts, deploy notifications arriving where the investigation context already lives
- Asking Claude something from your phone via Telegram, Discord, or iMessage while the work runs on your machine
- Long-running local sessions you want steerable from chat without cloud execution
- Incident response where alerts should land in a session already attached to the affected repo

## When NOT to Use It

- Delegating self-contained async work — a cloud session (Claude Code on the web) fits better than keeping a local session alive
- Steering an in-progress session interactively — Remote Control gives you the full session UI; channels give you a message pipe
- Enterprise setups on Amazon Bedrock, Google Cloud's Agent Platform (formerly Vertex AI), or Microsoft Foundry — channels require Anthropic authentication and are unavailable there

## How It Works

### Basic (Beginner)

1. Install a channel plugin: `/plugin install telegram@claude-plugins-official` (Discord works the same way; iMessage is macOS-only and needs Full Disk Access instead of a bot token; the prebuilt plugins require [Bun](https://bun.sh)). Run `/reload-plugins` to activate its configure command.
2. Configure credentials — for Telegram, create a bot via BotFather and run `/telegram:configure <token>`.
3. Restart with the channel enabled: `claude --channels plugin:telegram@claude-plugins-official`. Events only arrive while the session is open; for an always-on setup, run Claude in a persistent terminal or background process.
4. Pair your account: message the bot, get a pairing code, run `/telegram:access pair <code>`, then lock access down with `/telegram:access policy allowlist`.
5. Message the bot from your phone. The event arrives in your session as a channel message; Claude does the work locally and replies through the channel — you see the tool call in the terminal, the answer appears in the chat.

### Composing with Other Approaches (Intermediate)

- **Channels plus background agents**: a persistent background session with a channel attached is a standing worker you can message from anywhere — dispatch from chat in the morning, get the result in the same thread.
- **Channels plus permissions & safe autonomy**: unattended reactions need pre-decided boundaries — allowlist the safe inner loop so a webhook-triggered fix doesn't stall on a prompt nobody is present to approve. Channel servers that declare the permission-relay capability can forward prompts to the chat so you approve remotely.
- **Channels plus MCP context**: wire the error tracker's webhook in as the push trigger and its standard MCP server in as the lookup — the alert lands in the session with the repo open, and Claude queries full stack traces and event history from the same tracker on demand.

### Advanced Patterns

- **Webhook receiver**: build a small channel server that accepts webhooks from CI, deploy pipelines, or monitoring, and forwards them as channel events — the channels reference documents the contract (capability declaration, notification events, reply tools, sender gating).
- **Sender gating as security boundary**: every approved channel maintains a sender allowlist; unknown senders are silently dropped. Anyone on the allowlist of a permission-relaying channel can approve tool use in your session — only allowlist people you'd trust at your keyboard.
- **Organization rollout**: on Team/Enterprise plans channels are blocked until an Owner enables `channelsEnabled`; `allowedChannelPlugins` in managed settings restricts which plugins may register, including internal ones from a private marketplace.

## Common Pitfalls

- **Expecting delivery to a closed session**: events only arrive while the session runs. A channel is not a queue — messages sent while Claude Code is closed don't replay when it opens.
- **Forgetting `--channels` is per-session**: installing the plugin isn't enough, and being in `.mcp.json` isn't either; a server must be named in `--channels` at launch to push messages.
- **Unattended sessions hitting permission prompts**: the session pauses until someone responds. Decide the permission story first (allowlist rules, permission relay, or an isolated environment) before treating a channel session as autonomous.
- **Over-trusting the allowlist**: pairing your own account is the intended default. Adding teammates to a permission-relaying channel gives them approval authority over your session — that's a role, not a convenience.

## Sources

- [Push events into a running session with channels](https://code.claude.com/docs/en/channels) — Official docs: setup, supported plugins, security, enterprise controls
- [Channels reference](https://code.claude.com/docs/en/channels-reference) — The channel contract for building your own: events, reply tools, sender gating, permission relay
