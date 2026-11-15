import { ExitCode } from "../../../packages/protocol/src/index.ts";
import type { CoreEvent, ErrorEvent } from "../../../packages/protocol/src/index.ts";
import type { ChartResult, ChartType, PlotCommand, PlotInput, Point, Series } from "./contract.ts";
import { compact } from "./format.ts";
import { detectFormat, parseCsv, parseJson, type Table } from "./parse.ts";

export const CHART_TYPES: readonly ChartType[] = ["bar", "line", "sparkline", "histogram", "scatter"];

/** Where input comes from. Injected so tests and the TUI can reuse data already read. */
export interface PlotIO {
  read(file: string | undefined): Promise<{ text: string; name?: string }>;
}

export type Outcome = { ok: true; chart: ChartResult } | ErrorEvent;

export function createPlot(io: PlotIO): PlotCommand {
  return { name: "plot", run: (input) => plot(input, io) };
}

async function* plot(input: PlotInput, io: PlotIO): AsyncGenerator<CoreEvent<ChartResult>> {
  const outcome = await plotOnce(input, io);
  yield outcome.ok ? { type: "result", ok: true, data: outcome.chart } : outcome;
}

export async function plotOnce(input: PlotInput, io: PlotIO): Promise<Outcome> {
  const file = input.file === "-" ? undefined : input.file;
  let source: { text: string; name?: string };
  try {
    source = await io.read(file);
  } catch (error) {
    return fail("E_UNREADABLE_INPUT", `can't read ${file ?? "stdin"}: ${reason(error)}`, ExitCode.notFound,
      file ? "check the path, or pipe the data in" : "pass a file, or pipe CSV/JSON into chart");
  }
  if (!source.text.trim()) {
    return fail("E_EMPTY_INPUT", "the input is empty", ExitCode.invalidInput, "pass a CSV or JSON file with at least one row");
  }
  const format = detectFormat(source.text, input.format, source.name);
  let table: Table;
  try {
    table = format === "json" ? parseJson(source.text) : parseCsv(source.text);
  } catch (error) {
    return fail("E_UNREADABLE_INPUT", `couldn't read it as ${format.toUpperCase()}: ${reason(error)}`, ExitCode.invalidInput,
      `check the file, or force a format with --format ${format === "json" ? "csv" : "json"}`);
  }
  if (table.rows.length === 0) {
    return fail("E_EMPTY_INPUT", "there are no rows to plot", ExitCode.invalidInput, "the first line should be a header, followed by data rows");
  }
  return buildChart(input, table);
}

export function buildChart(input: PlotInput, table: Table): Outcome {
  const numeric = new Set(table.columns.filter((c) => isNumericColumn(table, c)));
  const columns = table.columns.map((name) => ({ name, numeric: numeric.has(name) }));

  for (const name of [input.x, ...input.y]) {
    if (name !== undefined && !table.columns.includes(name)) return notFound(name, table.columns);
  }
  for (const name of input.y) {
    if (!numeric.has(name)) {
      return fail("E_NOT_NUMERIC", `column "${name}" isn't numeric`, ExitCode.invalidInput,
        numeric.size ? `numeric columns: ${[...numeric].join(", ")}` : "no column here is entirely numbers");
    }
  }

  let type = input.type;
  let reasonText = type === "auto" ? "" : "chosen with --type";
  let x = input.x ?? table.columns.find((c) => !numeric.has(c));
  let y = input.y.length ? [...input.y] : [...numeric].filter((c) => c !== x);
  if (y.length === 0) {
    return fail("E_NOT_NUMERIC", "there's no numeric column to plot", ExitCode.invalidInput,
      `columns: ${table.columns.join(", ")}`);
  }

  // All-numeric data with no x column: one column is a distribution, two are a relationship.
  if (x === undefined) {
    if (type === "auto" && y.length === 1) {
      type = "histogram";
      reasonText = "one numeric column and no x axis, so this shows how its values are distributed";
    } else if ((type === "auto" && y.length === 2 && input.y.length === 0) || (type === "scatter" && y.length >= 2)) {
      if (type === "auto") reasonText = "two numeric columns, so a scatter shows how they relate";
      type = "scatter";
      [x, ...y] = y as [string, ...string[]];
    }
  }

  const xKind: ChartResult["x"]["kind"] =
    x === undefined || numeric.has(x) ? "number" : isTimeColumn(table, x) ? "time" : "category";

  if (type === "auto") {
    if (xKind === "time") {
      type = "line";
      reasonText = `${x} looks like dates, and lines show change over time`;
    } else if (xKind === "category") {
      type = "bar";
      reasonText = `${x} holds categories, and bars compare them`;
    } else if (x === undefined || isIncreasing(table, x)) {
      type = "line";
      reasonText = x === undefined ? "there's no x column, so rows are plotted in order" : `${x} increases steadily, so a line shows the trend`;
    } else {
      type = "scatter";
      reasonText = `${x} and ${y[0]} are both numeric, so a scatter shows how they relate`;
    }
  }

  if (type === "scatter" && (x === undefined || xKind !== "number")) {
    return fail("E_SCATTER_NEEDS_NUMERIC_X", "a scatter needs a numeric x column", ExitCode.invalidInput,
      `pass --x with one of: ${[...numeric].join(", ") || "(no numeric columns)"}`);
  }

  const base = { title: input.title, rows: table.rows.length, columns };

  if (type === "histogram") {
    const field = y[0]!;
    const values = table.rows.map((r) => r[field]).filter((v): v is number => typeof v === "number");
    if (y.length > 1) reasonText += `; only the first column (${field}) is binned`;
    return {
      ok: true,
      chart: { ...base, type, reason: reasonText, x: { field, kind: "number" }, series: [histogram(field, values, input.bins)] },
    };
  }

  const series = y.map((field) => {
    const points: Point[] = [];
    table.rows.forEach((row, i) => {
      const xv = x === undefined ? i + 1 : row[x];
      const yv = row[field];
      if (xv === null || xv === undefined || typeof yv !== "number") return;
      points.push({ x: xv, y: yv, row: i });
    });
    if (type === "line" && xKind === "number") points.sort((a, b) => Number(a.x) - Number(b.x));
    if (xKind === "time") points.sort((a, b) => time(a.x) - time(b.x));
    return { field, points, stats: stats(points.map((p) => p.y)) };
  });

  return {
    ok: true,
    chart: { ...base, type, reason: reasonText, x: { field: x ?? "#", kind: xKind }, series },
  };
}

function histogram(field: string, values: number[], bins?: number): Series {
  const count = bins ?? Math.max(2, Math.min(20, Math.ceil(Math.log2(values.length) + 1)));
  const min = Math.min(...values);
  const max = Math.max(...values);
  const width = (max - min) / count || 1;
  const counts = new Array<number>(count).fill(0);
  for (const v of values) counts[Math.min(count - 1, Math.floor((v - min) / width))]! += 1;
  const points = counts.map((c, i) => ({ x: `${compact(min + i * width)}–${compact(min + (i + 1) * width)}`, y: c }));
  return { field, points, stats: stats(values) };
}

function stats(values: number[]): Series["stats"] {
  if (values.length === 0) return { min: 0, max: 0, mean: 0, last: 0 };
  const sum = values.reduce((a, b) => a + b, 0);
  return { min: Math.min(...values), max: Math.max(...values), mean: sum / values.length, last: values[values.length - 1]! };
}

function isNumericColumn(table: Table, column: string): boolean {
  let seen = false;
  for (const row of table.rows) {
    const v = row[column];
    if (v === null || v === undefined) continue;
    if (typeof v !== "number") return false;
    seen = true;
  }
  return seen;
}

function time(value: string | number): number {
  return Date.parse(String(value).replace(/\//g, "-"));
}

const DATE = /^\d{4}[-/]\d{2}([-/]\d{2})?([T ][\d:.]+Z?)?$/;

function isTimeColumn(table: Table, column: string): boolean {
  const values = table.rows.map((r) => r[column]).filter((v) => v !== null);
  return values.length > 0 && values.every((v) => typeof v === "string" && DATE.test(v));
}

function isIncreasing(table: Table, column: string): boolean {
  const values = table.rows.map((r) => r[column]).filter((v): v is number => typeof v === "number");
  return values.every((v, i) => i === 0 || v > values[i - 1]!);
}

function notFound(name: string, columns: string[]): ErrorEvent {
  const best = columns
    .map((c) => ({ c, d: distance(name.toLowerCase(), c.toLowerCase()) }))
    .sort((a, b) => a.d - b.d)[0];
  const close = best && best.d <= Math.max(2, Math.floor(name.length / 3));
  return fail("E_COLUMN_NOT_FOUND", `column "${name}" not found`, ExitCode.invalidInput,
    `${close ? `did you mean "${best.c}"? ` : ""}columns: ${columns.join(", ")}`);
}

function distance(a: string, b: string): number {
  let prev = Array.from({ length: b.length + 1 }, (_, i) => i);
  for (let i = 1; i <= a.length; i++) {
    const cur = [i];
    for (let j = 1; j <= b.length; j++) {
      cur[j] = Math.min(prev[j]! + 1, cur[j - 1]! + 1, prev[j - 1]! + (a[i - 1] === b[j - 1] ? 0 : 1));
    }
    prev = cur;
  }
  return prev[b.length]!;
}

export function fail(code: string, message: string, exitCode: number, hint?: string): ErrorEvent {
  return { type: "error", ok: false, code, message, hint, exitCode };
}

function reason(error: unknown): string {
  if (error instanceof Error) return (error as NodeJS.ErrnoException).code === "ENOENT" ? "no such file" : error.message;
  return String(error);
}
