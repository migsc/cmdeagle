import { expect, test } from "bun:test";

const cli = `${import.meta.dir}/../cli.tsx`;
const fixtures = `${import.meta.dir}/../fixtures`;

async function run(args: string[], stdin?: string) {
  const proc = Bun.spawn(["bun", cli, ...args], { stdin: stdin === undefined ? "ignore" : new Blob([stdin]), stdout: "pipe", stderr: "pipe" });
  const [stdout, stderr, code] = await Promise.all([new Response(proc.stdout).text(), new Response(proc.stderr).text(), proc.exited]);
  return { stdout, stderr, code };
}

test("--json prints one result event", async () => {
  const { stdout, code } = await run([`${fixtures}/sales.csv`, "--json"]);
  expect(code).toBe(0);
  expect(JSON.parse(stdout)).toMatchObject({ type: "result", ok: true, data: { type: "line" } });
});

test("errors with --json go to stdout as an error event with the exit code", async () => {
  const { stdout, code } = await run([`${fixtures}/sales.csv`, "--y", "revnue", "--json"]);
  expect(code).toBe(2);
  expect(JSON.parse(stdout)).toMatchObject({ type: "error", code: "E_COLUMN_NOT_FOUND" });
});

test("human errors go to stderr with a hint", async () => {
  const { stdout, stderr, code } = await run([`${fixtures}/sales.csv`, "--type", "pie"]);
  expect(code).toBe(2);
  expect(stdout).toBe("");
  expect(stderr).toContain('unknown chart type "pie"');
  expect(stderr).toContain("hint:");
});

test("piped output has no color codes", async () => {
  const { stdout, code } = await run([`${fixtures}/metrics.json`, "--x", "service", "--y", "p95_ms"]);
  expect(code).toBe(0);
  expect(stdout).not.toContain("\x1b[");
  expect(stdout).toContain("worker");
});

test("reads data from stdin", async () => {
  const { stdout, code } = await run(["--json"], "x,y\na,1\nb,2\n");
  expect(code).toBe(0);
  expect(JSON.parse(stdout).data.type).toBe("bar");
});
