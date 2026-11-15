import type { Command } from "../../../packages/protocol/src/index.ts";

export type ChartType = "bar" | "line" | "sparkline" | "histogram" | "scatter";

export interface PlotInput {
  /** Undefined or "-" means stdin. */
  file?: string;
  type: ChartType | "auto";
  x?: string;
  y: string[];
  format: "auto" | "csv" | "json";
  bins?: number;
  width?: number;
  title?: string;
}

export interface Point {
  x: string | number;
  y: number;
  /** 0-based source row. Absent for histogram buckets. */
  row?: number;
}

export interface Series {
  field: string;
  points: Point[];
  stats: { min: number; max: number; mean: number; last: number };
}

/**
 * The normalized chart definition. Ink, plain text, and `--json` all render
 * from this, and the same shape will feed reusable chart components.
 */
export interface ChartResult {
  type: ChartType;
  /** Why this type was chosen, especially when `--type auto`. */
  reason: string;
  title?: string;
  x: { field: string; kind: "category" | "time" | "number" };
  series: Series[];
  rows: number;
  /** Every column in the input, so agents and the TUI can pick others. */
  columns: { name: string; numeric: boolean }[];
}

/**
 * Errors: E_COLUMN_NOT_FOUND (2, hint: closest column name),
 * E_NOT_NUMERIC (2), E_EMPTY_INPUT (2), E_UNREADABLE_INPUT (3).
 */
export type PlotCommand = Command<PlotInput, ChartResult>;
