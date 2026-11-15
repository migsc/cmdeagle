# todo

A to-do list for the terminal.

```sh
todo                         # interactive list on a TTY
todo add "renew passport" --due 2026-12-01
todo ls --json
todo done 3 4
todo clear --dry-run         # shows what would be deleted
```

## What it proves

- Saved state in a declared file, with stable item IDs.
- Scripted and interactive use with the same capabilities: every TUI key
  maps to a command, so an agent can do anything the TUI can.
- Destructive commands (`rm`, `clear`) with required dry runs and a confirm
  prompt that agents answer with `--yes`.
- Components: `List`, `Checkbox`, `Prompt`.
- Later: concurrent edits from the TUI and an agent through the optional
  separate process, and a Todoist backend behind the same commands.

## Fixtures

- `fixtures/todos.json`: a sample state file.
