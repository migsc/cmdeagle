import { useState } from "react";
import { Box, Text, render, useApp, useInput, useWindowSize } from "ink";
import type { ChartResult, ChartType, PlotInput } from "../core/contract.ts";
import { CHART_TYPES, type Outcome } from "../core/index.ts";
import { draw, type Line } from "./draw.ts";

interface Props {
  input: PlotInput;
  first: Outcome;
  /** Re-runs the headless core. Every key here maps to a flag the CLI also accepts. */
  run(input: PlotInput): Promise<Outcome>;
}

function Explorer({ input, first, run }: Props) {
  const { exit } = useApp();
  const { columns: width, rows: height } = useWindowSize();
  const [current, setCurrent] = useState(input);
  const [outcome, setOutcome] = useState(first);
  const [lastChart, setLastChart] = useState<ChartResult | undefined>(first.ok ? first.chart : undefined);

  const change = (next: PlotInput) => {
    setCurrent(next);
    void run(next).then((result) => {
      setOutcome(result);
      if (result.ok) setLastChart(result.chart);
    });
  };

  const shown: ChartType | "auto" = outcome.ok ? outcome.chart.type : current.type;
  const yChoices: string[][] = [
    [],
    ...(lastChart?.columns.filter((c) => c.numeric && c.name !== lastChart.x.field).map((c) => [c.name]) ?? []),
  ];

  useInput((ch, key) => {
    if (ch === "q" || key.escape || (key.ctrl && ch === "c")) return exit();
    if (ch === "t" || key.rightArrow || key.leftArrow) {
      const step = key.leftArrow ? -1 : 1;
      const at = shown === "auto" ? -1 : CHART_TYPES.indexOf(shown);
      change({ ...current, type: CHART_TYPES[(at + step + CHART_TYPES.length) % CHART_TYPES.length]! });
    }
    if (ch === "a") change({ ...current, type: "auto" });
    if (ch === "y") {
      const at = yChoices.findIndex((choice) => choice.join() === current.y.join());
      change({ ...current, y: yChoices[(at + 1) % yChoices.length]! });
    }
  });

  const flags = [
    current.type === "auto" ? "--type auto" : `--type ${current.type}`,
    current.y.length ? `--y ${current.y.join(",")}` : "",
  ].filter(Boolean).join(" ");

  return (
    <Box flexDirection="column" paddingX={1}>
      <Box justifyContent="space-between">
        <Text bold color="green">chart</Text>
        <Text dimColor>{flags}</Text>
      </Box>
      <Text> </Text>
      {outcome.ok ? (
        <Lines lines={draw(outcome.chart, width - 3, Math.max(6, Math.min(24, height - 12)))} />
      ) : (
        <Box flexDirection="column">
          <Text color="red">{outcome.message}</Text>
          {outcome.hint ? <Text color="yellow">hint: {outcome.hint}</Text> : null}
        </Box>
      )}
      <Text> </Text>
      <Text dimColor>←/→ or t  chart type · a  auto · y  columns · q  quit</Text>
    </Box>
  );
}

function Lines({ lines }: { lines: Line[] }) {
  return (
    <Box flexDirection="column">
      {lines.map((line, i) => (
        <Text key={i} wrap="truncate">
          {line.length === 0
            ? " "
            : line.map((s, j) => (
                <Text key={j} color={s.color} dimColor={s.dim} bold={s.bold}>
                  {s.text}
                </Text>
              ))}
        </Text>
      ))}
    </Box>
  );
}

export async function explore(props: Props, stdin?: NodeJS.ReadableStream): Promise<void> {
  const app = render(<Explorer {...props} />, { exitOnCtrlC: false, stdin });
  await app.waitUntilExit();
}
