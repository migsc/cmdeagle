# t3ctl

Rebuild a real, working personal tool (the T3 Code controller) with cmdeagle.

## Two phases

1. **Wrap.** The spec runs the existing Python script through `exec`. This
   tests the "any language" claim: validation, help, effects, and the
   generated skill all come from cmdeagle while the logic stays in Python.
2. **Port.** Move commands into TypeScript cores one at a time, keeping the
   same JSON output, so the switch is invisible to agents.

## What it proves

- Wrapping an executable in another language with separate argv (no shell strings).
- Loopback-only `network` effects.
- A new `cost` effect: `send` spends provider quota or money and never runs
  without an explicit confirmation flag that an agent must get from the user.
- Rebuilding something real, not a demo.

## Setup

`legacy/t3ctl` is not committed. Link or copy the script locally
(`ln -s ~/.local/bin/t3ctl legacy/t3ctl`).

## Fixtures

- `fixtures/status.json`: **made-up** status output. No real project or
  thread IDs go in this public repo.
