# workspace

See, name, and clear Hyprland workspaces. Combines the ideas in the personal
`workspace-clear` and `workspace-label` scripts.

```sh
workspace ls --json
workspace label 3 "cmdeagle"
workspace clear 3 --dry-run     # which windows close, which stay, and why
workspace pick                  # interactive switcher
```

## What it proves

- A truly destructive command (closing windows, ending processes) with a
  required dry run that explains what stays and why.
- A `desktop` side effect for focus changes. The generated skill tells an
  agent to ask before stealing focus, which is a real rule on this desktop.
- A required external executable (`hyprctl`), checked up front with a clear error.
- Tests run against saved `hyprctl -j` output through an adapter, never
  against a live desktop.
- Components: `List`, `Table`.

## Fixtures

- `fixtures/hyprctl-clients.json` and `fixtures/hyprctl-workspaces.json`:
  **made-up** output in the shape of `hyprctl -j clients` and
  `hyprctl -j workspaces`.
