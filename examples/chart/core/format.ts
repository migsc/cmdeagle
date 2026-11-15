/** Short, readable numbers for labels: 12400 → 12.4k, 0.25 → 0.25. */
export function compact(n: number): string {
  const abs = Math.abs(n);
  // Thousands stay exact below 10,000 (8100 reads better than 8.1k).
  const units: [threshold: number, divisor: number, suffix: string][] = [[1e9, 1e9, "B"], [1e6, 1e6, "M"], [1e4, 1e3, "k"]];
  for (const [threshold, divisor, suffix] of units) {
    if (abs >= threshold) return trim((n / divisor).toFixed(1)) + suffix;
  }
  if (Number.isInteger(n)) return String(n);
  return trim(n.toFixed(abs < 1 ? 3 : 2));
}

function trim(s: string): string {
  return s.includes(".") ? s.replace(/\.?0+$/, "") : s;
}
