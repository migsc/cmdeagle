# weather

The forecast at the terminal: `weather Austin --days 3 --units f`.

Uses [Open-Meteo](https://open-meteo.com/), which needs no API key.

## What it proves

- A `network` side effect with named hosts, which permission rules can allow or block.
- A readable offline error with a next step (`E_OFFLINE`, hint: `--cached`).
- Caching as a declared `writes` effect.
- Unit flags and range validation (`--days 1..16`).
- Tests run fully offline from saved responses.
- Components: `Table`, `KeyValue`, `Sparkline` (temperatures across the days).

## Fixtures

- `fixtures/geocode-austin.json` and `fixtures/forecast-austin.json`: responses in
  Open-Meteo's geocoding and forecast shapes, with made-up values. Compare them with
  a real response before relying on exact field names.
