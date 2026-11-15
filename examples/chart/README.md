# chart

CSV or JSON in, a chart out.

```sh
chart fixtures/sales.csv --x month --y revenue,cost --type line
cat fixtures/metrics.json | chart --x service --y p95_ms --type bar
chart fixtures/latency.csv --y ms --type histogram --json
chart fixtures/sales.csv --y revnue      # error: did you mean "revenue"?
```

## Try it (implemented)

```sh
cd examples/chart
bun install
bun cli.tsx fixtures/sales.csv -i     # run from source
bun test                              # 19 tests
bun run build                         # dist/chart, standalone
bun run install-bin                   # build + copy to ~/.local/bin/chart
```

## What it proves

- Input from a file or stdin, with CSV/JSON detection.
- Width-aware rendering, and a plain-text fallback that still reads well.
- `--json` returns the normalized chart definition and stats.
- Helpful validation errors (`E_COLUMN_NOT_FOUND` with the closest match).
- Hand-written skill guidance (`skill/plotting.md`) that teaches an agent
  how to plot, not just which flags exist.
- The first reusable components: `BarChart`, `LineChart`, `Sparkline`,
  `Histogram`, `Scatter`. Each renders from `ChartResult` as Ink, as
  plain text, or as data.

## Fixtures

- `fixtures/sales.csv`: monthly revenue and cost (line, bar).
- `fixtures/metrics.json`: per-service latency (bar).
- `fixtures/latency.csv`: request latencies (histogram).
