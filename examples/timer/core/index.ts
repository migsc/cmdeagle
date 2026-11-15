import type { CoreEvent } from "../../../packages/protocol/src/index.ts";
import type { StartCommand, TimerEvent, TimerResult } from "./contract.ts";

/**
 * How a renderer steers a running timer. The TUI maps keys to these; the CLI
 * maps SIGINT to `stop()`. A future `timer serve` would expose the same calls.
 */
export class TimerControls {
  #paused = false;
  #stopped = false;
  #wake: (() => void) | undefined;

  get paused(): boolean {
    return this.#paused;
  }

  get stopped(): boolean {
    return this.#stopped;
  }

  toggle(): void {
    this.#paused = !this.#paused;
    this.#wake?.();
  }

  stop(): void {
    this.#stopped = true;
    this.#wake?.();
  }

  /** Sleeps until `ms` passes or a control call wakes it early. */
  sleep(ms: number): Promise<void> {
    return new Promise((resolve) => {
      const done = () => {
        clearTimeout(timer);
        this.#wake = undefined;
        resolve();
      };
      const timer = setTimeout(done, ms);
      this.#wake = done;
    });
  }
}

/** Seconds to display for a remaining time, tolerant of timers firing a few ms early. */
export function displaySeconds(remainingMs: number): number {
  return Math.max(0, Math.ceil((remainingMs - 50) / 1000));
}

export function formatClock(remainingMs: number): string {
  const total = displaySeconds(remainingMs);
  const h = Math.floor(total / 3600);
  const m = Math.floor((total % 3600) / 60);
  const s = total % 60;
  const mm = String(m).padStart(h ? 2 : 1, "0");
  return h ? `${h}:${mm}:${String(s).padStart(2, "0")}` : `${mm.padStart(2, "0")}:${String(s).padStart(2, "0")}`;
}

/** Parses 25m, 1h30m, 90s, 1h, or a bare number of minutes. */
export function parseDuration(text: string): number | undefined {
  if (/^\d+(\.\d+)?$/.test(text)) return Math.round(Number(text) * 60_000);
  const match = /^(?:(\d+)h)?(?:(\d+)m)?(?:(\d+)s)?$/.exec(text);
  if (!match || text === "") return undefined;
  const [, h = "0", m = "0", s = "0"] = match;
  const ms = (Number(h) * 3600 + Number(m) * 60 + Number(s)) * 1000;
  return ms > 0 ? ms : undefined;
}

export function createStart(controls: TimerControls, clock: () => number = () => performance.now()): StartCommand {
  return {
    name: "start",
    async *run(input, ctx): AsyncGenerator<CoreEvent<TimerResult, TimerEvent>> {
      ctx.signal.addEventListener("abort", () => controls.stop(), { once: true });
      const startedAt = ctx.now().toISOString();
      let remaining = input.durationMs;
      let last = clock();
      let paused = false;

      yield { type: "timer.tick", remainingMs: remaining, elapsedMs: 0 };
      while (remaining > 0 && !controls.stopped) {
        // Wake on each whole-second boundary so the display changes exactly then.
        let wait = remaining % 1000 || 1000;
        if (wait < 50) wait += 1000;
        await controls.sleep(controls.paused ? 60_000 : Math.min(wait, remaining));
        const now = clock();
        if (!paused) remaining = Math.max(0, remaining - (now - last));
        last = now;
        if (controls.paused !== paused) {
          paused = controls.paused;
          const remainingMs = Math.round(remaining);
          yield paused ? { type: "timer.paused", remainingMs } : { type: "timer.resumed", remainingMs };
        }
        if (!paused && !controls.stopped) {
          yield { type: "timer.tick", remainingMs: Math.round(remaining), elapsedMs: Math.round(input.durationMs - remaining) };
        }
      }

      const data: TimerResult = {
        label: input.label,
        durationMs: input.durationMs,
        completed: remaining <= 0,
        startedAt,
        endedAt: ctx.now().toISOString(),
      };
      yield { type: "result", ok: true, data };
    },
  };
}
