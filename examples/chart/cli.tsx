#!/usr/bin/env bun
/**
 * Hand-written stand-in for what cmdeagle will generate from cmdeagle.yaml:
 * flag parsing, validation, and choosing a renderer for the headless core.
 */
import { openSync } from "node:fs";
import { ReadStream } from "node:tty";
import { parseArgs } from "node:util";
import { ExitCode, type ErrorEvent } from "../../packages/protocol/src/index.ts";
import type { ChartType, PlotInput } from "./core/contract.ts";
import { CHART_TYPES, fail, plotOnce, type PlotIO } from "./core/index.ts";
import { toText } from "./ui/ansi.ts";
import { draw } from "./ui/draw.ts";

const HELP = `chart: turn CSV or JSON into a chart, right in the terminal.

Usage
  chart [file] [options]
  cat data.csv | chart [options]

Options
  -t, --type <type>     auto, bar, line, sparkline, histogram, scatter (default: auto)
      --x <column>      Column for the x axis or categories
      --y <columns>     Numeric columns to plot, comma-separated
  -f, --format <fmt>    auto, csv, json (default: auto)
      --bins <n>        Histogram bucket count
  -w, --width <n>       Override the terminal width
      --height <n>      Plot height in rows (default: 12)
      --title <text>    Title above the chart
  -i, --interactive     Explore the data with an interactive chart
      --json            Print the chart definition as JSON instead of drawing it
      --plain           No colors, even in a terminal
  -h, --help            Show this help

Examples
  chart sales.csv --x month --y revenue,cost
  chart metrics.json --x service --y p95_ms --type bar
  chart latency.csv --type histogram --bins 12
  chart sales.csv -i`;

async function main(argv: string[]): Promise<number> {
  const wantsJson = argv.includes("--json");
  let parsed;
  try {
    parsed = parseArgs({
      args: argv,
      allowPositionals: true,
      options: {
        type: { type: "string", short: "t", default: "auto" },
        x: { type: "string" },
        y: { type: "string", multiple: true },
        format: { type: "string", short: "f", default: "auto" },
        bins: { type: "string" },
        width: { type: "string", short: "w" },
        height: { type: "string" },
        title: { type: "string" },
        interactive: { type: "boolean", short: "i" },
        json: { type: "boolean" },
        plain: { type: "boolean" },
        help: { type: "boolean", short: "h" },
      },
    });
  } catch (error) {
    return report(fail("E_INVALID_FLAG", (error as Error).message, ExitCode.invalidInput, "run chart --help"), wantsJson);
  }
  const { values, positionals } = parsed;
  if (values.help) {
    console.log(HELP);
    return ExitCode.ok;
  }

  const type = values.type as string;
  if (type !== "auto" && !CHART_TYPES.includes(type as ChartType)) {
    return report(fail("E_INVALID_FLAG", `unknown chart type "${type}"`, ExitCode.invalidInput, `use one of: auto, ${CHART_TYPES.join(", ")}`), wantsJson);
  }
  const format = values.format as string;
  if (!["auto", "csv", "json"].includes(format)) {
    return report(fail("E_INVALID_FLAG", `unknown format "${format}"`, ExitCode.invalidInput, "use auto, csv, or json"), wantsJson);
  }
  const bins = integer(values.bins, "--bins", 2);
  if (typeof bins === "object") return report(bins, wantsJson);
  const width = integer(values.width, "--width", 20);
  if (typeof width === "object") return report(width, wantsJson);
  const height = integer(values.height, "--height", 4);
  if (typeof height === "object") return report(height, wantsJson);
  if (positionals.length > 1) {
    return report(fail("E_INVALID_ARGS", "chart takes one file at a time", ExitCode.invalidInput, "pass a single CSV or JSON file"), wantsJson);
  }
  const file = positionals[0];
  if (file === undefined && process.stdin.isTTY) {
    return report(fail("E_NO_INPUT", "no data to plot", ExitCode.invalidInput, "pass a file (chart sales.csv) or pipe data in (cat sales.csv | chart)"), wantsJson);
  }

  const input: PlotInput = {
    file,
    type: type as PlotInput["type"],
    x: values.x,
    y: (values.y ?? []).flatMap((v) => v.split(",")).map((v) => v.trim()).filter(Boolean),
    format: format as PlotInput["format"],
    bins,
    width,
    title: values.title,
  };

  // Read once, so the interactive view can re-plot stdin data.
  let cached: Promise<{ text: string; name?: string }> | undefined;
  const io: PlotIO = {
    read: (path) => (cached ??= path ? Bun.file(path).text().then((text) => ({ text, name: path })) : Bun.stdin.text().then((text) => ({ text }))),
  };

  const outcome = await plotOnce(input, io);

  if (values.json) {
    console.log(JSON.stringify(outcome.ok ? { type: "result", ok: true, data: outcome.chart } : outcome));
    return outcome.ok ? ExitCode.ok : outcome.exitCode;
  }

  if (values.interactive) {
    if (!process.stdout.isTTY) {
      return report(fail("E_NOT_A_TERMINAL", "--interactive needs a terminal", ExitCode.invalidInput, "drop --interactive when piping the output"), false);
    }
    // Keys come from the terminal even when the data arrived on stdin.
    const keys = process.stdin.isTTY ? undefined : new ReadStream(openSync("/dev/tty", "r"));
    const { explore } = await import("./ui/explorer.tsx");
    await explore({ input, first: outcome, run: (next) => plotOnce(next, io) }, keys);
    return ExitCode.ok;
  }

  if (!outcome.ok) return report(outcome, false);
  const color = Boolean(process.stdout.isTTY) && !values.plain && !process.env.NO_COLOR;
  const columns = input.width ?? (process.stdout.isTTY ? process.stdout.columns - 1 : 80);
  console.log(toText(draw(outcome.chart, columns, height ?? 12), color));
  return ExitCode.ok;
}

function integer(raw: string | undefined, flag: string, min: number): number | ErrorEvent | undefined {
  if (raw === undefined) return undefined;
  const value = Number(raw);
  if (!Number.isInteger(value) || value < min) {
    return fail("E_INVALID_FLAG", `${flag} must be a whole number of at least ${min}`, ExitCode.invalidInput, `e.g. ${flag} ${min * 4}`);
  }
  return value;
}

function report(error: ErrorEvent, json: boolean): number {
  if (json) {
    console.log(JSON.stringify(error));
  } else {
    const tty = process.stderr.isTTY && !process.env.NO_COLOR;
    const paint = (code: number, text: string) => (tty ? `\x1b[${code}m${text}\x1b[0m` : text);
    console.error(`${paint(31, "chart:")} ${error.message}`);
    if (error.hint) console.error(`${paint(33, "hint:")} ${error.hint}`);
  }
  return error.exitCode;
}

process.exit(await main(process.argv.slice(2)));
