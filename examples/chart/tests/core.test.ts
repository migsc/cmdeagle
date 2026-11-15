import { describe, expect, test } from "bun:test";
import type { PlotInput } from "../core/contract.ts";
import { plotOnce, type Outcome, type PlotIO } from "../core/index.ts";
import { parseCsv } from "../core/parse.ts";
import { draw, lineText } from "../ui/draw.ts";

const fixture = (name: string): PlotIO => ({
  read: async () => ({ text: await Bun.file(`${import.meta.dir}/../fixtures/${name}`).text(), name }),
});
const inline = (text: string): PlotIO => ({ read: async () => ({ text }) });

const input = (overrides: Partial<PlotInput> = {}): PlotInput => ({ type: "auto", y: [], format: "auto", ...overrides });

function chartOf(outcome: Outcome) {
  if (!outcome.ok) throw new Error(`expected a chart, got ${outcome.code}: ${outcome.message}`);
  return outcome.chart;
}

describe("parsing", () => {
  test("quoted CSV fields keep commas, quotes and numbers", () => {
    const table = parseCsv('name,note,value\n"Smith, J","said ""hi""",3.5\n');
    expect(table.rows).toEqual([{ name: "Smith, J", note: 'said "hi"', value: 3.5 }]);
  });
});

describe("auto chart type", () => {
  test("dates on x make a line", async () => {
    const chart = chartOf(await plotOnce(input(), fixture("sales.csv")));
    expect(chart.type).toBe("line");
    expect(chart.x).toEqual({ field: "month", kind: "time" });
    expect(chart.series.map((s) => s.field)).toEqual(["revenue", "cost"]);
  });

  test("categories make bars", async () => {
    const chart = chartOf(await plotOnce(input({ x: "service", y: ["p95_ms"] }), fixture("metrics.json")));
    expect(chart.type).toBe("bar");
    expect(chart.series[0]!.stats.max).toBe(910);
  });

  test("one numeric column makes a histogram", async () => {
    const chart = chartOf(await plotOnce(input(), fixture("latency.csv")));
    expect(chart.type).toBe("histogram");
    expect(chart.series[0]!.points.reduce((sum, p) => sum + p.y, 0)).toBe(200);
  });

  test("two numeric columns make a scatter", async () => {
    const chart = chartOf(await plotOnce(input(), inline("a,b\n1,2\n2,4\n3,1\n")));
    expect(chart.type).toBe("scatter");
    expect(chart.x.field).toBe("a");
  });
});

describe("errors an agent can repair", () => {
  test("a misspelled column suggests the closest one", async () => {
    const outcome = await plotOnce(input({ y: ["revnue"] }), fixture("sales.csv"));
    expect(outcome).toMatchObject({ ok: false, code: "E_COLUMN_NOT_FOUND", exitCode: 2 });
    expect(!outcome.ok && outcome.hint).toContain('did you mean "revenue"');
  });

  test("plotting a text column is refused", async () => {
    const outcome = await plotOnce(input({ y: ["month"] }), fixture("sales.csv"));
    expect(outcome).toMatchObject({ ok: false, code: "E_NOT_NUMERIC" });
  });

  test("a scatter needs a numeric x", async () => {
    const outcome = await plotOnce(input({ type: "scatter", x: "month" }), fixture("sales.csv"));
    expect(outcome).toMatchObject({ ok: false, code: "E_SCATTER_NEEDS_NUMERIC_X" });
  });

  test("a missing file is reported, not thrown", async () => {
    const missing: PlotIO = { read: () => Bun.file("/nonexistent.csv").text().then((text) => ({ text })) };
    expect(await plotOnce(input({ file: "/nonexistent.csv" }), missing)).toMatchObject({ code: "E_UNREADABLE_INPUT", exitCode: 3 });
  });
});

describe("review findings", () => {
  test("an unclosed quote is a parse error, not a merged row", async () => {
    const outcome = await plotOnce(input(), inline('value,category\n1,"first\n2,second\n'));
    expect(outcome).toMatchObject({ ok: false, code: "E_UNREADABLE_INPUT" });
  });

  test("dates are plotted in time order", async () => {
    const chart = chartOf(await plotOnce(input(), inline("date,value\n2026-03-01,30\n2026-01-01,10\n2026-02-01,20\n")));
    expect(chart.series[0]!.points.map((p) => p.y)).toEqual([10, 20, 30]);
  });

  test("a category missing from the first series still gets its bar", async () => {
    const chart = chartOf(await plotOnce(input({ type: "bar" }), inline("category,a,b\nfirst,1,\nsecond,,20\nthird,3,30\n")));
    const text = draw(chart, 60).map(lineText).join("\n");
    expect(text).toContain("second");
    expect(text).toMatch(/second\s+–\n\s+[█▏▎▍▌▋▊▉]+ 20/);
  });

  test("histograms with 10+ buckets stay in order", async () => {
    const chart = chartOf(await plotOnce(input({ bins: 12 }), inline("value\n0\n11\n")));
    const starts = draw(chart, 50).map(lineText).filter((l) => /^\d/.test(l)).map((l) => Number.parseFloat(l));
    expect(starts).toEqual([...starts].sort((a, b) => a - b));
    expect(starts.length).toBe(12);
  });

  test("a long label never pushes a bar past the width", async () => {
    const chart = chartOf(await plotOnce(input(), inline("category,value\nabcdefghijklmnopqrstuvwx,-120\nother,40\n")));
    for (const line of draw(chart, 20)) expect(lineText(line).length).toBeLessThanOrEqual(20);
  });

  test("negative values draw bars left of a zero axis", async () => {
    const chart = chartOf(await plotOnce(input(), inline("category,value\na,-5\nb,10\n")));
    const [a, b] = draw(chart, 40).map(lineText).filter((l) => l.startsWith("a ") || l.startsWith("b "));
    expect(a).toMatch(/█+│/);
    expect(b).toMatch(/│█+/);
  });
});

describe("drawing", () => {
  test("repeated category labels keep their own values", async () => {
    const chart = chartOf(await plotOnce(input(), inline("name,mb\nclaude,375\ncodex,373\nclaude,120\n")));
    const text = draw(chart, 60).map(lineText);
    expect(text.filter((l) => l.startsWith("claude")).map((l) => l.trim().split(" ").at(-1))).toEqual(["375", "120"]);
  });

  test.each(["bar", "line", "sparkline", "histogram", "scatter"] as const)("%s fits the requested width", async (type) => {
    const io = type === "scatter" ? inline("a,b\n1,2\n2,4\n3,1\n") : fixture("sales.csv");
    const chart = chartOf(await plotOnce(input({ type, y: type === "scatter" ? [] : ["revenue"] }), io));
    for (const line of draw(chart, 60)) expect(lineText(line).length).toBeLessThanOrEqual(60);
  });
});
