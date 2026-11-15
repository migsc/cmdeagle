import type { Color, Line } from "./draw.ts";

const CODES: Record<Color, number> = { red: 31, green: 32, yellow: 33, blue: 34, magenta: 35, cyan: 36, gray: 90 };

/** Renders drawn lines as text, with ANSI color only when `color` is true. */
export function toText(lines: Line[], color: boolean): string {
  return lines
    .map((line) =>
      line
        .map((s) => {
          if (!color || (!s.color && !s.dim && !s.bold)) return s.text;
          const codes = [s.bold ? 1 : 0, s.dim ? 2 : 0, s.color ? CODES[s.color] : 0].filter(Boolean);
          return `\x1b[${codes.join(";")}m${s.text}\x1b[0m`;
        })
        .join("")
        .trimEnd(),
    )
    .join("\n");
}
