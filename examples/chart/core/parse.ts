export type Cell = string | number | null;
export type Row = Record<string, Cell>;

export interface Table {
  columns: string[];
  rows: Row[];
}

export function detectFormat(text: string, hint: "auto" | "csv" | "json", name?: string): "csv" | "json" {
  if (hint !== "auto") return hint;
  if (name?.endsWith(".json")) return "json";
  if (name?.endsWith(".csv")) return "csv";
  const first = text.trimStart()[0];
  return first === "[" || first === "{" ? "json" : "csv";
}

export function parseCsv(text: string): Table {
  const records = splitCsv(text).filter((r) => r.some((f) => f.trim() !== ""));
  const [header, ...body] = records;
  if (!header) return { columns: [], rows: [] };
  const columns = header.map((h) => h.trim());
  const rows = body.map((record) => Object.fromEntries(columns.map((c, i) => [c, coerce(record[i])])));
  return { columns, rows };
}

/** RFC 4180-style splitting: quoted fields may hold commas, quotes ("") and newlines. */
function splitCsv(text: string): string[][] {
  const records: string[][] = [];
  let record: string[] = [];
  let field = "";
  let quoted = false;
  for (let i = 0; i < text.length; i++) {
    const ch = text[i];
    if (quoted) {
      if (ch !== '"') field += ch;
      else if (text[i + 1] === '"') { field += '"'; i++; }
      else quoted = false;
    } else if (ch === '"') {
      quoted = true;
    } else if (ch === ",") {
      record.push(field);
      field = "";
    } else if (ch === "\n" || ch === "\r") {
      if (ch === "\r" && text[i + 1] === "\n") i++;
      record.push(field);
      records.push(record);
      record = [];
      field = "";
    } else {
      field += ch;
    }
  }
  if (quoted) throw new ParseError("a quoted field is never closed");
  if (field !== "" || record.length > 0) {
    record.push(field);
    records.push(record);
  }
  return records;
}

export class ParseError extends Error {}

/** Accepts an array of objects, or an object whose first array property holds them. */
export function parseJson(text: string): Table {
  const data: unknown = JSON.parse(text);
  const list = Array.isArray(data) ? data : isObject(data) ? Object.values(data).find(Array.isArray) : undefined;
  if (!list) throw new ParseError("expected an array of objects, or an object holding one");
  const columns = new Set<string>();
  const rows: Row[] = [];
  for (const item of list) {
    if (!isObject(item)) continue;
    const row: Row = {};
    for (const [key, value] of Object.entries(item)) {
      columns.add(key);
      row[key] = typeof value === "number" ? value
        : typeof value === "string" ? coerce(value)
        : value === null || value === undefined ? null
        : JSON.stringify(value);
    }
    rows.push(row);
  }
  return { columns: [...columns], rows };
}

function coerce(value: string | undefined): Cell {
  if (value === undefined) return null;
  const trimmed = value.trim();
  if (trimmed === "") return null;
  const n = Number(trimmed);
  return Number.isFinite(n) ? n : trimmed;
}

function isObject(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}
