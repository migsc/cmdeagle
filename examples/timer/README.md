# timer

A focus timer: `timer 25m --label focus`.

```sh
timer 25m --label focus      # live countdown; space pauses, q stops
timer 90s --json             # one JSON event per line while it runs
timer 50m --detach           # runs in the background
timer status --json          # time left, from any terminal or agent
timer stop
```

## Try it (implemented)

```sh
cd examples/timer
bun install
bun cli.tsx 25s --label focus         # run from source
bun test                              # 17 tests
bun run build                         # dist/timer, standalone
bun run install-bin                   # build + copy to ~/.local/bin/timer
```

`--detach`, `status`, and `stop` are not built yet; they report `E_NOT_IMPLEMENTED`.

## What it proves

- Long-running, live-updating output in the TUI (`Countdown`, `Progress`).
- NDJSON event streaming when piped: `timer.tick` events, then one `result`.
- Signals: Ctrl-C ends with `completed: false` and exit 130.
- A `desktop` side effect (the notification), conditional on `--notify`.
- The optional separate process: `--detach` plus `status`/`stop` against it.
- Smallest-tool benchmarks: startup time and compiled binary size.

## Fixtures

- `fixtures/run-3s.ndjson`: a recorded event stream for TUI snapshot tests.

Later, this replaces `~/.local/bin/timer-focus`. A v1 `.cmd.yaml` version will be added to test the converter.
