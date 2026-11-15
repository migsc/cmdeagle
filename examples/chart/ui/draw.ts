/**
 * Chart drawing shared by every renderer. It produces lines of styled
 * segments; the ANSI renderer and the Ink view only decide how to color them.
 * This is the seed of the future @cmdeagle/ui chart components.
 */
import type { ChartResult, Series } from "../core/contract.ts";
import { compact } from "../core/format.ts";

export type Color = "green" | "cyan" | "yellow" | "magenta" | "blue" | "red" | "gray";

export interface Segment {
  text: string;
  color?: Color;
  dim?: boolean;
  bold?: boolean;
}

export type Line = Segment[];

export const PALETTE: readonly Color[] = ["green", "cyan", "yellow", "magenta", "blue", "red"];

const color = (i: number): Color => PALETTE[i % PALETTE.length]!;

export function draw(chart: ChartResult, width: number, height = 12): Line[] {
  const lines: Line[] = [];
  if (chart.title) lines.push([{ text: chart.title, bold: true }], []);
  if (chart.series.length > 1 && chart.type !== "sparkline") lines.push(legend(chart.series), []);
  switch (chart.type) {
    case "bar":
    case "histogram":
      lines.push(...bars(chart, width));
      break;
    case "sparkline":
      lines.push(...sparklines(chart, width));
      break;
    case "line":
    case "scatter":
      lines.push(...plotArea(chart, width, height));
      break;
  }
  lines.push([], ...footer(chart, width));
  return lines;
}

function legend(series: Series[]): Line {
  return series.flatMap((s, i) => [{ text: "■ ", color: color(i) }, { text: `${s.field}   ` }]);
}

function footer(chart: ChartResult, width: number): Line[] {
  const y = chart.type === "histogram" ? `${chart.x.field} (count)` : chart.series.map((s) => s.field).join(", ");
  const x = chart.type === "histogram" ? "" : ` · x: ${chart.x.field}`;
  return [
    [{ text: truncate(`${chart.type}${x} · y: ${y} · ${chart.rows} rows`, width), dim: true }],
    ...wrap(chart.reason, width).map((text) => [{ text, dim: true }]),
  ];
}

function wrap(text: string, width: number): string[] {
  const lines: string[] = [];
  let current = "";
  for (const word of text.split(" ")) {
    if (current && current.length + 1 + word.length > width) {
      lines.push(current);
      current = word;
    } else {
      current = current ? `${current} ${word}` : word;
    }
  }
  if (current) lines.push(current);
  return lines;
}

const EIGHTHS = ["", "▏", "▎", "▍", "▌", "▋", "▊", "▉"];

function bar(length: number): string {
  const full = Math.floor(length);
  const part = Math.round((length - full) * 8);
  return "█".repeat(full + (part === 8 ? 1 : 0)) + (part === 8 ? "" : EIGHTHS[part]);
}

function bars(chart: ChartResult, width: number): Line[] {
  // Group by source row so a value missing from one series, or a repeated
  // label (two "claude" processes), never hides or merges a row.
  const groups = new Map<number, { label: string; values: (number | undefined)[] }>();
  chart.series.forEach((s, i) => {
    s.points.forEach((p, n) => {
      const key = p.row ?? n;
      const group = groups.get(key) ?? { label: String(p.x), values: [] };
      group.values[i] = p.y;
      groups.set(key, group);
    });
  });
  const rows = [...groups.entries()].sort(([a], [b]) => a - b).map(([, g]) => g);
  if (rows.length === 0) return [];

  const values = chart.series.flatMap((s) => s.points.map((p) => p.y));
  const max = Math.max(0, ...values);
  const min = Math.min(0, ...values);
  const span = max - min || 1;
  const valueWidth = Math.max(...values.map((v) => compact(v).length));
  const fixed = valueWidth + 3 + (min < 0 ? 1 : 0);
  const labelWidth = Math.max(1, Math.min(24, Math.max(...rows.map((r) => r.label.length)), width - fixed - 4));
  const barWidth = Math.max(4, width - labelWidth - fixed);
  // With negatives, bars grow left and right from a zero axis.
  const left = min < 0 ? Math.round((-min / span) * barWidth) : 0;
  const right = barWidth - left;

  const lines: Line[] = [];
  for (const { label, values: rowValues } of rows) {
    chart.series.forEach((_, i) => {
      const v = rowValues[i];
      const name = i === 0 ? truncate(label, labelWidth).padEnd(labelWidth) : " ".repeat(labelWidth);
      const negative = v !== undefined && v < 0 ? "█".repeat(Math.round((-v / span) * barWidth)).padStart(left) : " ".repeat(left);
      const positive = v !== undefined && v > 0 ? bar((v / span) * barWidth) : "";
      lines.push([
        { text: `${name} `, dim: i > 0 },
        ...(min < 0 ? [{ text: negative, color: color(i) }, { text: "│", dim: true }] : []),
        { text: positive.padEnd(min < 0 ? right : 0), color: color(i) },
        { text: ` ${v === undefined ? "–" : compact(v)}`, dim: true },
      ]);
    });
  }
  return lines;
}

const TICKS = "▁▂▃▄▅▆▇█";

function sparklines(chart: ChartResult, width: number): Line[] {
  const labelWidth = Math.max(...chart.series.map((s) => s.field.length));
  return chart.series.map((s, i) => {
    const summary = ` min ${compact(s.stats.min)}  max ${compact(s.stats.max)}  last ${compact(s.stats.last)}`;
    const room = Math.max(4, width - labelWidth - summary.length - 2);
    const ys = resample(s.points.map((p) => p.y), room);
    const lo = Math.min(...ys);
    const span = Math.max(...ys) - lo || 1;
    const spark = ys.map((v) => TICKS[Math.round(((v - lo) / span) * (TICKS.length - 1))]).join("");
    return [{ text: `${s.field.padEnd(labelWidth)} ` }, { text: spark, color: color(i) }, { text: summary, dim: true }];
  });
}

function resample(values: number[], size: number): number[] {
  if (values.length <= size) return values;
  return Array.from({ length: size }, (_, i) => {
    const from = Math.floor((i * values.length) / size);
    const to = Math.max(from + 1, Math.floor(((i + 1) * values.length) / size));
    const slice = values.slice(from, to);
    return slice.reduce((a, b) => a + b, 0) / slice.length;
  });
}

// Braille cells hold a 2×4 dot grid, so the plot has twice the columns and four times the rows.
const DOTS = [[0x01, 0x08], [0x02, 0x10], [0x04, 0x20], [0x40, 0x80]] as const;

function plotArea(chart: ChartResult, width: number, height: number): Line[] {
  const values = chart.series.flatMap((s) => s.points.map((p) => p.y));
  let lo = Math.min(...values);
  let hi = Math.max(...values);
  if (lo === hi) { lo -= 1; hi += 1; }
  const gutter = Math.max(compact(hi).length, compact(lo).length, compact((hi + lo) / 2).length);
  const cols = Math.max(10, width - gutter - 2);
  const rows = Math.max(4, height);
  const W = cols * 2;
  const H = rows * 4;
  const cells = new Uint8Array(cols * rows);
  const owner = new Int8Array(cols * rows).fill(-1);

  const numericX = chart.x.kind === "number";
  const order = [...new Set(chart.series.flatMap((s) => s.points.map((p) => String(p.x))))];
  if (chart.x.kind === "time") order.sort((a, b) => Date.parse(a.replace(/\//g, "-")) - Date.parse(b.replace(/\//g, "-")));
  const xs = chart.series.flatMap((s) => s.points.map((p) => Number(p.x)));
  const xlo = numericX ? Math.min(...xs) : 0;
  const xspan = numericX ? Math.max(...xs) - xlo || 1 : Math.max(1, order.length - 1);
  const px = (x: string | number) => Math.round(((numericX ? Number(x) - xlo : order.indexOf(String(x))) / xspan) * (W - 1));
  const py = (y: number) => H - 1 - Math.round(((y - lo) / (hi - lo)) * (H - 1));

  const set = (x: number, y: number, s: number) => {
    if (x < 0 || y < 0 || x >= W || y >= H) return;
    const i = (y >> 2) * cols + (x >> 1);
    cells[i]! |= DOTS[y & 3]![x & 1]!;
    owner[i] = s;
  };

  chart.series.forEach((s, si) => {
    let prev: [number, number] | undefined;
    for (const p of s.points) {
      const cur: [number, number] = [px(p.x), py(p.y)];
      if (chart.type === "line" && prev) segment(prev, cur, (x, y) => set(x, y, si));
      else set(cur[0], cur[1], si);
      prev = cur;
    }
  });

  const lines: Line[] = [];
  for (let r = 0; r < rows; r++) {
    const label = r === 0 ? compact(hi) : r === rows - 1 ? compact(lo) : r === Math.floor(rows / 2) ? compact((hi + lo) / 2) : "";
    const line: Line = [{ text: `${label.padStart(gutter)} ${label ? "┤" : "│"}`, dim: true }];
    let run: Segment | undefined;
    for (let c = 0; c < cols; c++) {
      const i = r * cols + c;
      const ch = cells[i] ? String.fromCharCode(0x2800 + cells[i]!) : " ";
      const segColor = owner[i]! >= 0 ? color(owner[i]!) : undefined;
      if (run && run.color === segColor) run.text += ch;
      else line.push((run = { text: ch, color: segColor }));
    }
    lines.push(line);
  }
  lines.push([{ text: `${" ".repeat(gutter)} └${"─".repeat(cols)}`, dim: true }]);
  lines.push([{ text: " ".repeat(gutter + 2) + xLabels(chart, order, xlo, xspan, cols), dim: true }]);
  return lines;
}

function xLabels(chart: ChartResult, order: string[], xlo: number, xspan: number, cols: number): string {
  const numericX = chart.x.kind === "number";
  const label = (t: number) => (numericX ? compact(xlo + t * xspan) : order[Math.round(t * (order.length - 1))] ?? "");
  const first = label(0);
  const last = label(1);
  const middle = cols >= 40 ? label(0.5) : "";
  const out = Array.from({ length: cols }, () => " ");
  const place = (text: string, at: number) => [...text].forEach((ch, i) => { if (at + i >= 0 && at + i < cols) out[at + i] = ch; });
  place(first, 0);
  if (middle && middle !== first && middle !== last) place(middle, Math.floor(cols / 2 - middle.length / 2));
  if (last !== first) place(last, cols - last.length);
  return out.join("");
}

/** Bresenham line between two dot positions. */
function segment([x0, y0]: [number, number], [x1, y1]: [number, number], plot: (x: number, y: number) => void) {
  const dx = Math.abs(x1 - x0);
  const dy = -Math.abs(y1 - y0);
  const sx = x0 < x1 ? 1 : -1;
  const sy = y0 < y1 ? 1 : -1;
  let err = dx + dy;
  for (;;) {
    plot(x0, y0);
    if (x0 === x1 && y0 === y1) return;
    const e2 = 2 * err;
    if (e2 >= dy) { err += dy; x0 += sx; }
    if (e2 <= dx) { err += dx; y0 += sy; }
  }
}

function truncate(text: string, width: number): string {
  return text.length <= width ? text : `${text.slice(0, Math.max(0, width - 1))}…`;
}

export function lineText(line: Line): string {
  return line.map((s) => s.text).join("");
}
