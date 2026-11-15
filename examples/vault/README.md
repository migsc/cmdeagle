# vault

Find and write notes in an Obsidian vault.

```sh
vault search "cmdeagle" --json
vault new "Weekly review" --folder Projects --link-daily --dry-run
echo "Notes from today" | vault new "Standup" --body - --link-daily
```

## What it proves

- Writes shown as file diffs on a required dry run, before anything changes.
- Config from an environment variable (`VAULT_ROOT`) with a clear error when missing.
- stdin input for note bodies.
- Agent-safe edits: the agent shows you the diff, then applies it.
- Components: `Diff`, `List`.

## Fixtures

- `fixtures/vault/`: a tiny **made-up** vault with a daily note, a project
  note, and an inbox note.
