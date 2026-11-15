#!/usr/bin/env bun
/**
 * Hand-written stand-in for what cmdeagle will generate from cmdeagle.yaml:
 * flag parsing, validation, and choosing a renderer for the headless core.
 */
import { parseArgs } from "node:util";
import { ExitCode, type Context, type ErrorEvent } from "../../packages/protocol/src/index.ts";
import type { StartInput, TimerResult } from "./core/contract.ts";
import { TimerControls, createStart, formatClock, parseDuration } from "./core/index.ts";

const HELP = `timer: a focus timer that stays out of your way.

Usage
  timer [duration] [options]

Duration
  25m, 1h30m, 90s, 1h, or a bare number of minutes (default: 25m)

Options
  -l, --label <text>   Shown in the countdown and the notification
      --no-notify      Skip the desktop notification at the end
      --json           Stream events as JSON lines (ticks, then a result)
      --plain          One line at the start and one at the end, no TUI
  -h, --help           Show this help

Keys
  space or p   pause / resume
  q or esc     stop early
  ctrl-c       cancel (exit 130)

Examples
  timer 25m --label focus
  timer 90s --json | jq -c 'select(.type == "result")'`;

async function main(argv: string[]): Promise<number> {
  const wantsJson = argv.includes("--json");
  let parsed;
  try {
    parsed = parseArgs({
      args: argv,
      allowPositionals: true,
      options: {
        label: { type: "string", short: "l" },
        "no-notify": { type: "boolean" },
        json: { type: "boolean" },
        plain: { type: "boolean" },
        detach: { type: "boolean" },
        help: { type: "boolean", short: "h" },
      },
    });
  } catch (error) {
    return report(fail("E_INVALID_FLAG", (error as Error).message, ExitCode.invalidInput, "run timer --help"), wantsJson);
  }
  const { values, positionals } = parsed;
  if (values.help) {
    console.log(HELP);
    return ExitCode.ok;
  }
  const [first = "25m", ...rest] = positionals;
  if (first === "status" || first === "stop" || values.detach) {
    return report(fail("E_NOT_IMPLEMENTED", "detached timers aren't built yet", ExitCode.failure,
      "run the timer in the foreground for now: timer 25m"), wantsJson);
  }
  if (rest.length > 0) {
    return report(fail("E_INVALID_ARGS", `unexpected argument "${rest[0]}"`, ExitCode.invalidInput,
      'put multi-word labels in quotes: timer 25m --label "deep work"'), wantsJson);
  }
  const durationMs = parseDuration(first);
  if (durationMs === undefined) {
    return report(fail("E_INVALID_DURATION", `"${first}" isn't a duration`, ExitCode.invalidInput,
      "use forms like 25m, 1h30m, 90s, or 10"), wantsJson);
  }

  const input: StartInput = { durationMs, label: values.label, notify: !values["no-notify"], detach: false };
  const controls = new TimerControls();
  const abort = new AbortController();
  const ctx: Context = {
    ask: () => Promise.reject(new Error("timer never asks questions")),
    dryRun: false,
    signal: abort.signal,
    now: () => new Date(),
  };
  const events = createStart(controls).run(input, ctx);
  const cancelled = { current: false };
  let result: TimerResult | undefined;

  if (values.json || values.plain || !process.stdout.isTTY) {
    process.on("SIGINT", () => {
      cancelled.current = true;
      abort.abort();
    });
    if (!values.json) console.log(`timer: ${formatClock(durationMs)}${input.label ? ` · ${input.label}` : ""}`);
    for await (const event of events) {
      if (values.json) console.log(JSON.stringify(event));
      if (event.type === "result") result = event.data;
    }
    if (!values.json && result) {
      console.log(result.completed ? `done${input.label ? `: ${input.label}` : ""}` : `stopped with ${formatClock(result.durationMs - elapsed(result))} left`);
    }
  } else {
    const { show } = await import("./ui/view.tsx");
    await show({ input, events, controls, cancelled, onResult: (r) => (result = r) });
  }

  if (result?.completed && input.notify) await notify(input);
  return cancelled.current ? ExitCode.cancelled : ExitCode.ok;
}

function elapsed(result: TimerResult): number {
  return Date.parse(result.endedAt) - Date.parse(result.startedAt);
}

/** The declared `desktop` effect. Missing notify-send is not an error. */
async function notify(input: StartInput): Promise<void> {
  try {
    const proc = Bun.spawn(["notify-send", "-a", "timer", input.label ?? "Timer", `${formatClock(input.durationMs)} is up`], { stdout: "ignore", stderr: "ignore" });
    await proc.exited;
  } catch {
    // No notification daemon; the terminal already showed the result.
  }
}

function fail(code: string, message: string, exitCode: number, hint?: string): ErrorEvent {
  return { type: "error", ok: false, code, message, hint, exitCode };
}

function report(error: ErrorEvent, json: boolean): number {
  if (json) {
    console.log(JSON.stringify(error));
  } else {
    const tty = process.stderr.isTTY && !process.env.NO_COLOR;
    const paint = (code: number, text: string) => (tty ? `\x1b[${code}m${text}\x1b[0m` : text);
    console.error(`${paint(31, "timer:")} ${error.message}`);
    if (error.hint) console.error(`${paint(33, "hint:")} ${error.hint}`);
  }
  return error.exitCode;
}

process.exit(await main(process.argv.slice(2)));
