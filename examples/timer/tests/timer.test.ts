import { describe, expect, test } from "bun:test";
import type { Context } from "../../../packages/protocol/src/index.ts";
import { TimerControls, createStart, formatClock, parseDuration } from "../core/index.ts";

const ctx = (signal = new AbortController().signal): Context => ({
  ask: () => Promise.reject(new Error("unused")),
  dryRun: false,
  signal,
  now: () => new Date(),
});

describe("durations", () => {
  test.each([
    ["25m", 1_500_000],
    ["1h30m", 5_400_000],
    ["90s", 90_000],
    ["10", 600_000],
    ["1.5", 90_000],
  ])("%s", (text, ms) => expect(parseDuration(text)).toBe(ms));

  test.each(["", "0m", "abc", "5x", "m"])("rejects %p", (text) => expect(parseDuration(text)).toBeUndefined());

  test("clock format", () => {
    expect(formatClock(1_500_000)).toBe("25:00");
    expect(formatClock(5_400_000)).toBe("1:30:00");
    expect(formatClock(59_990)).toBe("01:00");
  });
});

describe("core", () => {
  test("ticks each second and completes", async () => {
    const events = [];
    for await (const e of createStart(new TimerControls()).run({ durationMs: 1200, notify: false, detach: false }, ctx())) events.push(e);
    expect(events.at(0)).toMatchObject({ type: "timer.tick", remainingMs: 1200 });
    expect(events.filter((e) => e.type === "timer.tick").length).toBeGreaterThanOrEqual(3);
    expect(events.at(-1)).toMatchObject({ type: "result", data: { completed: true, durationMs: 1200 } });
  });

  test("pause holds the remaining time, resume continues", async () => {
    const controls = new TimerControls();
    const types: string[] = [];
    const run = (async () => {
      for await (const e of createStart(controls).run({ durationMs: 600, notify: false, detach: false }, ctx())) types.push(e.type);
    })();
    await Bun.sleep(100);
    controls.toggle();
    await Bun.sleep(700); // longer than the whole timer: it must not finish while paused
    expect(types).not.toContain("result");
    controls.toggle();
    await run;
    expect(types).toEqual(expect.arrayContaining(["timer.paused", "timer.resumed", "result"]));
  });

  test("abort ends early with completed: false", async () => {
    const abort = new AbortController();
    setTimeout(() => abort.abort(), 100);
    let last;
    for await (const e of createStart(new TimerControls()).run({ durationMs: 60_000, notify: false, detach: false }, ctx(abort.signal))) last = e;
    expect(last).toMatchObject({ type: "result", data: { completed: false } });
  });
});

describe("cli", () => {
  const cli = `${import.meta.dir}/../cli.tsx`;

  test("--json streams ticks then one result", async () => {
    const proc = Bun.spawn(["bun", cli, "1s", "--json", "--no-notify", "--label", "test"], { stdout: "pipe" });
    const lines = (await new Response(proc.stdout).text()).trim().split("\n").map((l) => JSON.parse(l));
    expect(await proc.exited).toBe(0);
    expect(lines.at(-1)).toMatchObject({ type: "result", data: { completed: true, label: "test" } });
  });

  test("Ctrl-C exits 130 with completed: false", async () => {
    const proc = Bun.spawn(["bun", cli, "1m", "--json", "--no-notify"], { stdout: "pipe" });
    await Bun.sleep(400);
    proc.kill("SIGINT");
    const lines = (await new Response(proc.stdout).text()).trim().split("\n").map((l) => JSON.parse(l));
    expect(await proc.exited).toBe(130);
    expect(lines.at(-1)).toMatchObject({ type: "result", data: { completed: false } });
  });

  test("a bad duration explains the accepted forms", async () => {
    const proc = Bun.spawn(["bun", cli, "soon"], { stderr: "pipe" });
    const stderr = await new Response(proc.stderr).text();
    expect(await proc.exited).toBe(2);
    expect(stderr).toContain("use forms like 25m");
  });
});
