# backup

The reference example from the v2 plan: copy a folder somewhere safe.

```sh
backup ~/Projects --to /mnt/vault --dry-run
backup ~/Projects --to /mnt/vault --dry-run --json
backup ~/Projects --to /mnt/vault --prune      # asks before deleting
```

## What it proves

- The basics: spec, validation, help, readable errors, stable exit codes.
- A required dry run that returns a plan as a preview or as JSON, with
  `next` holding the command that applies it.
- A conditional `destructive` effect (`--prune`) with a confirm prompt.
- A generated `SKILL.md` that makes an agent preview before copying.
- Components: `Progress`, `Table`.

## Fixtures

- `fixtures/source/`: a small source tree.
- `fixtures/dest/`: an older backup with one changed file and one stale file,
  so the plan has additions, changes, and a prune candidate.
