import { useEffect, useState } from "react";
import { Box, Text, render, useApp, useInput, useWindowSize } from "ink";
import type { CoreEvent } from "../../../packages/protocol/src/index.ts";
import type { StartInput, TimerEvent, TimerResult } from "../core/contract.ts";
import { formatClock, type TimerControls } from "../core/index.ts";

// 3×5 pixel font; each lit pixel draws as two block characters.
const FONT: Record<string, string[]> = {
  "0": ["111", "101", "101", "101", "111"],
  "1": ["010", "110", "010", "010", "111"],
  "2": ["111", "001", "111", "100", "111"],
  "3": ["111", "001", "111", "001", "111"],
  "4": ["101", "101", "111", "001", "001"],
  "5": ["111", "100", "111", "001", "111"],
  "6": ["111", "100", "111", "101", "111"],
  "7": ["111", "001", "010", "010", "010"],
  "8": ["111", "101", "111", "101", "111"],
  "9": ["111", "101", "111", "001", "111"],
  ":": ["0", "1", "0", "1", "0"],
};

function bigText(text: string): string[] {
  return [0, 1, 2, 3, 4].map((row) =>
    [...text].map((ch) => (FONT[ch]?.[row] ?? "").replace(/1/g, "██").replace(/0/g, "  ")).join("  "),
  );
}

interface Props {
  input: StartInput;
  events: AsyncIterable<CoreEvent<TimerResult, TimerEvent>>;
  controls: TimerControls;
  /** Set when Ctrl-C ends the timer, so the CLI can exit 130. */
  cancelled: { current: boolean };
  onResult(result: TimerResult): void;
}

function TimerView({ input, events, controls, cancelled, onResult }: Props) {
  const { exit } = useApp();
  const { columns } = useWindowSize();
  const [remaining, setRemaining] = useState(input.durationMs);
  const [paused, setPaused] = useState(false);
  const [result, setResult] = useState<TimerResult>();

  useEffect(() => {
    void (async () => {
      for await (const event of events) {
        if (event.type === "timer.tick") setRemaining(event.remainingMs);
        else if (event.type === "timer.paused") setPaused(true);
        else if (event.type === "timer.resumed") setPaused(false);
        else if (event.type === "result") {
          setResult(event.data);
          onResult(event.data);
        }
      }
    })();
  }, []);

  useEffect(() => {
    if (result) setTimeout(exit, 20);
  }, [result]);

  useInput((ch, key) => {
    if (result) return;
    if (ch === " " || ch === "p") controls.toggle();
    if (ch === "q" || key.escape) controls.stop();
    if (key.ctrl && ch === "c") {
      cancelled.current = true;
      controls.stop();
    }
  });

  const clock = formatClock(result?.completed ? 0 : remaining);
  const accent = result ? (result.completed ? "green" : "gray") : paused ? "yellow" : "green";
  const barWidth = Math.max(10, Math.min(48, columns - 8));
  const done = Math.round((1 - remaining / input.durationMs) * barWidth);
  const status = result
    ? result.completed ? "✓ done" : "stopped"
    : paused ? "paused · space to resume · q to stop" : "space to pause · q to stop";

  return (
    <Box flexDirection="column" paddingX={2} paddingY={1}>
      <Text bold color={accent}>{input.label ?? "timer"}</Text>
      <Text> </Text>
      {columns >= 44 ? (
        bigText(clock).map((row, i) => <Text key={i} color={accent}>{row}</Text>)
      ) : (
        <Text bold color={accent}>{clock}</Text>
      )}
      <Text> </Text>
      <Text>
        <Text color={accent}>{"█".repeat(result?.completed ? barWidth : done)}</Text>
        <Text dimColor>{"░".repeat(result?.completed ? 0 : barWidth - done)}</Text>
      </Text>
      <Text dimColor>{status}</Text>
    </Box>
  );
}

export async function show(props: Props): Promise<void> {
  const app = render(<TimerView {...props} />, { exitOnCtrlC: false });
  await app.waitUntilExit();
}
