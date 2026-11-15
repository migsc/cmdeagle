# cmdeagle v2 examples

Real CLIs and TUIs that serve as test subjects. If an example can't do something
the v2 site claims, the claim isn't true yet. Everything here is a **draft**
until the runtime spike lands. Specs and contracts describe intent.

## Rules

- **TypeScript only.** No `.js`/`.mjs` source. Bun installs, runs, tests, and
  compiles (`bun build --compile`).
- **Headless cores.** Each command is `input → events` (see
  [`packages/protocol`](../packages/protocol/src/index.ts)). Ink views live in
  `ui/` and only render events. Importing Ink or React from `core/` is a lint error.
- **Made-up data only.** This repo is public. Fixtures never contain real
  notes, transcripts, desktop state, account details, or IDs.
- **Tests never touch the network or a live desktop.** They use `fixtures/`.

## Layout

```
<name>/
  cmdeagle.yaml   draft spec: commands, flags, side effects
  core/           headless commands; contract.ts holds the typed input/events/result
  ui/             Ink views (arrives with implementation)
  fixtures/       fake data and recorded event streams (.ndjson)
  tests/          core via JSON, ui via frame snapshots, a few PTY end-to-end runs
```

## What each example proves

| Claim | timer | weather | todo | chart | backup | sessions | vault | workspace | t3ctl |
|---|---|---|---|---|---|---|---|---|---|
| Args/flags validation, help, errors | ● | ● | ● | ● | ● | ● | ● | ● | ● |
| Subcommands | ● | | ● | | | ● | ● | ● | ● |
| Interactive TUI | ● | | ● | ● | | ● | | ● | |
| Live-updating output | ● | | | | ● | ● | | | |
| NDJSON event stream when piped | ● | | | | ● | ● | | | |
| Prompts answered by flags | | | ● | | ● | | | ● | ● |
| stdin input | | | | ● | | | ● | | |
| Saved state | ● | ● | ● | | | | | ● | |
| `reads` | ● | | ● | ● | ● | ● | ● | ● | ● |
| `writes` | ● | ● | ● | | ● | | ● | ● | ● |
| `network` | | ● | | | | | | | ● |
| `destructive` + required dry run | | | ● | | ● | | | ● | |
| `desktop` | ● | | | | | | | ● | |
| `cost` | | | | | | | | | ● |
| Required external executable | | | | | | | | ● | ● |
| Wraps another language (`exec`) | | | | | | | | | ● |
| Optional separate process (`serve`) | ● | | ○ | | | | | | |
| Hand-written skill guidance | | | | ● | | | | | |
| v1 `.cmd.yaml` converter | ○ | | | | | | | | |

● planned · ○ later

## Shared components

Components move into [`packages/ui`](../packages/ui/README.md) once a second example needs them.

| Component | First needed by |
|---|---|
| `BarChart`, `LineChart`, `Sparkline`, `Histogram`, `Scatter` | chart (Sparkline also weather) |
| `Countdown`, `Progress` | timer (Progress also backup, sessions) |
| `List`, `Checkbox`, `Prompt` | todo (List also sessions, vault, workspace) |
| `Table`, `KeyValue` | weather (Table also backup, sessions, workspace) |
| `Diff` | vault |

## Outgrowing this folder

When an example gets real users or its own release cadence, split it into its
own repo with `git subtree split --prefix examples/<name>`. cmdeagle keeps a
pinned copy as a regression test.

## Status

**chart** and **timer** are implemented (hand-written; the builder doesn't exist yet).
Each is its own Bun package: `bun install`, `bun test`, `bun run build`. The others are
draft specs and contracts. Type-check those with
`deno check packages/protocol/src/index.ts examples/*/core/contract.ts`.

`stubs/react-devtools-core` is a TypeScript stand-in for Ink's optional devtools
import. Without it, `bun build --compile` binaries crash at startup.
