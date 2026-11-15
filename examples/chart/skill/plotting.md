## How to choose a chart

- **Change over time** (dates or ordered steps on x): `line`. Several series are fine up to about four.
- **Comparing categories** (names on x, one value each): `bar`. Sort by value unless the order means something.
- **One trend at a glance** inside other output: `sparkline`.
- **Distribution of one numeric column** (latencies, sizes): `histogram`. Leave `--bins` unset first.
- **Relationship between two numeric columns**: `scatter`.

Let `--type auto` decide when unsure. The JSON result's `reason` says why it chose that type.

## Before plotting

- Run with `--json` first when the columns are unknown. Errors list the available columns.
- Pass `--x` and `--y` explicitly for anything you will show the user.
- Do not plot more than one unit on the same chart (e.g. dollars and percentages). Make two charts.
