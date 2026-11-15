# sessions

Search and browse coding-agent transcripts (Codex and Claude Code JSONL).

```sh
sessions                                  # fuzzy finder with a preview pane
sessions search "hyprctl reload" --agent codex --json
sessions show 2026-11-14-a1b2 --role user
```

Answers "did we already do this?" before an agent redoes work.

## What it proves

- A read-only tool: only `reads` effects, so it's always safe to run.
- Streaming large files with progress events.
- The richest TUI in the set: a fuzzy finder with a live preview.
- Per-agent format adapters behind one result shape.
- Components: `List`, `Table`, a preview pane.

## Fixtures

- `fixtures/codex/` and `fixtures/claude/`: **made-up** transcripts in a
  simplified form of each format. Real transcripts never go in this public
  repo. Check the real formats before writing the adapters.
