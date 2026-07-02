# Checkpoints and Rewind
*Last reviewed: 2026-07-02*

## What It Is

Checkpoints and Rewind gives you the ability to undo AI-generated code changes without losing the conversation that produced them. Every time Claude edits a file, it automatically saves a snapshot of the previous state. You can rewind to any earlier snapshot, restoring your files to exactly how they were at that point, while keeping everything Claude learned during the conversation. You can also fork the conversation from a checkpoint to explore a completely different direction.

## Why It Works

Experimentation is how good software gets built, but experimentation requires cheap failure. In traditional development, undoing a series of changes means manual `git reset`, `git stash`, or hunting through reflog entries — friction that discourages bold attempts. Checkpoints make reverting trivially cheap, which changes your behavior: you try riskier refactors, explore more alternatives, and say "let's try it and see" instead of spending twenty minutes debating whether an approach will work. The key insight is separating the state of the code from the state of the conversation — you can revert what the AI did while preserving what the AI learned.

## When to Use It

- Experimental refactoring where you want to try an approach and revert if it does not work out
- Comparing two implementation strategies by rewinding and trying the alternative
- Recovering from a wrong turn when Claude misunderstands your intent and edits the wrong files
- Iterative design where you build, evaluate, tear down, and rebuild with new constraints

## When NOT to Use It

- When you have already committed the changes — checkpoints track uncommitted file edits, not git history
- For long-term version control — while checkpoints do persist across sessions, they are automatically cleaned up after 30 days and are not a substitute for git commits
- When the task is straightforward and well-understood — if you know exactly what you want, there is nothing to experiment with

## How It Works

### Basic (Beginner)

1. Start working with Claude normally. Every file edit creates an automatic checkpoint — no action needed from you.
2. After several edits, you realize the approach is wrong. Type `/rewind`, use the undo keyboard shortcut, or press Escape twice (when the prompt input is empty) to open the rewind menu.
3. The rewind menu offers three restore options: "Restore code and conversation" (full rewind), "Restore conversation" (keep current code), or "Restore code" (keep conversation). Choose the option that fits your situation.
4. Ask Claude to try a different approach: "That approach had too much coupling. Let's use an event-based design instead." Claude now has the context of the failed attempt to inform the new one.
5. If the new approach works, continue normally. The earlier checkpoint is still available if you need it.

### Composing with Other Approaches (Intermediate)

- **Checkpoints plus Plan Mode**: Use Plan Mode to outline an approach, execute it, then evaluate the result. If the execution reveals a flaw in the plan, rewind the code but keep the plan context. Refine the plan and re-execute — each iteration is informed by the previous attempt.
- **Checkpoints plus worktree isolation**: In a worktree, checkpoints give you two levels of undo: rewind within the session for fine-grained rollback, or discard the entire worktree for a complete reset. Use checkpoints for "that last edit was wrong" and worktree discard for "this whole approach was wrong."
- **Checkpoint forking for A/B comparison**: After reaching a checkpoint, run `/branch try-alternative` to create a copy of the conversation and switch into it, leaving the original intact (from the CLI: `claude --continue --fork-session`). Both branches share the same starting point but diverge in implementation; forks are grouped under their root session in the `/resume` picker, so comparing results is a session switch away.

### Advanced Patterns

- **Intentional exploration loops**: Deliberately ask Claude to try an approach you suspect might fail, because the failure will reveal constraints you have not articulated yet. Rewind after the attempt and use the new understanding in your revised prompt. This is faster than trying to specify every constraint upfront.
- **Progressive refinement**: Build a first draft, evaluate it, rewind to a checkpoint partway through (not all the way back), and refine from the midpoint. This preserves the parts that worked while revising the parts that did not.
- **Checkpoint as decision record**: Before rewinding, ask Claude to summarize what it tried and why it did not work. This creates a natural record of rejected approaches, which is valuable when a teammate later asks "why didn't you just do X?"

## Common Pitfalls

- **Waiting too long to rewind**: If you let Claude make twenty edits before realizing the approach is wrong, rewinding requires stepping back through many checkpoints. Evaluate early and often — rewind after three or four edits if the direction feels off.
- **Confusing checkpoints with git commits**: Checkpoints are file snapshots that persist across sessions (auto-cleaned after 30 days). They do not create git commits and do not appear in `git log`. If you want to preserve a state permanently, commit it.
- **Losing the lesson when rewinding**: The whole point of rewind is that the conversation context survives. After rewinding, explicitly reference what went wrong: "The last approach coupled OrderService to PaymentGateway. Let's keep them decoupled this time." This gives Claude actionable constraints for the next attempt.
- **Over-reliance on rewind instead of planning**: If you find yourself rewinding five times in a row, stop and switch to Plan Mode. Repeated rewinds often mean the problem needs more upfront analysis, not more trial and error.

## Real-World Example

You are extracting a `NotificationService` from a monolithic `UserController` in a Rails application. You ask Claude to move all notification logic into a new service class.

Claude creates `app/services/notification_service.rb`, moves six methods out of `app/controllers/user_controller.rb`, and updates the controller to delegate. You run `bundle exec rspec spec/controllers/user_controller_spec.rb` and 4 of 22 tests fail — the extraction broke the test setup because the tests were stubbing methods that no longer exist on the controller.

You realize Claude moved the methods but kept the same method signatures, creating a service that is still tightly coupled to controller concerns (it references `params` and `current_user` directly). Rather than asking Claude to patch this incrementally, you rewind:

```
> /rewind
> That extraction kept controller dependencies in the service. Let's try again:
  extract notification logic into NotificationService, but give it a clean
  interface that accepts explicit arguments (user_id, event_type, payload)
  instead of accessing controller state. Update the controller to build
  those arguments and pass them in.
```

Claude rebuilds `notification_service.rb` with a proper interface: `NotificationService.new.send_notification(user_id: user.id, event_type: :welcome, payload: { name: user.name })`. It updates the five controller call sites to build these arguments and pass them in, and adjusts the test stubs to match.

All 22 tests pass. The rewind cost ten seconds; the flawed first attempt taught both you and Claude exactly what "clean extraction" meant in this context.

## Sources

- [Claude Code Interactive Mode](https://code.claude.com/docs/en/interactive-mode) — Official docs covering /rewind and checkpoint options
- [Claude Code IDE Integrations](https://code.claude.com/docs/en/ide-integrations) — VS Code checkpoint UI and rewind controls
