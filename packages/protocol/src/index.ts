/**
 * Draft cmdeagle v2 command protocol.
 *
 * A command core is headless: it takes typed input and yields events. Renderers
 * (Ink TUI, plain text, JSON/NDJSON, and optionally a socket or MCP server) are
 * clients of the same stream, so the TUI can never do something the CLI and
 * agents cannot. Core code must not import Ink or React.
 */

export const PROTOCOL_VERSION = 1;

/** Side-effect kinds a command declares in cmdeagle.yaml. */
export type EffectKind =
  | "reads"
  | "writes"
  | "network"
  | "destructive"
  | "desktop"
  | "cost";

export interface ProgressEvent {
  type: "progress";
  step: string;
  label?: string;
  done?: number;
  total?: number;
}

export interface LogEvent {
  type: "log";
  level: "info" | "warn";
  message: string;
}

/**
 * A question for the user. A TUI answers it from the keyboard; agents and
 * scripts answer it up front with `flag`.
 */
export type InputRequest =
  | { kind: "confirm"; id: string; message: string; flag: string }
  | { kind: "text"; id: string; message: string; flag: string; placeholder?: string }
  | { kind: "select"; id: string; message: string; flag: string; choices: readonly string[] };

export type InputEvent = { type: "input" } & InputRequest;

export interface ResultEvent<R> {
  type: "result";
  ok: true;
  data: R;
}

export interface ErrorEvent {
  type: "error";
  ok: false;
  /** Stable, documented code such as `E_COLUMN_NOT_FOUND`. */
  code: string;
  message: string;
  /** One actionable next step. Agents repair against this. */
  hint?: string;
  exitCode: number;
}

/** Command-specific events are namespaced, e.g. `timer.tick`. */
export interface DomainEvent {
  type: `${string}.${string}`;
}

export type CoreEvent<R, D extends DomainEvent = never> =
  | ProgressEvent
  | LogEvent
  | InputEvent
  | ResultEvent<R>
  | ErrorEvent
  | D;

export type Answer<Q extends InputRequest> = Q extends { kind: "confirm" } ? boolean : string;

export interface Context {
  /**
   * Emits an `input` event and waits for the answer. Resolves from the
   * keyboard in a TUI or from flags elsewhere; a renderer with no way to
   * answer fails with `E_INPUT_REQUIRED` naming the flag.
   */
  ask<Q extends InputRequest>(request: Q): Promise<Answer<Q>>;
  readonly dryRun: boolean;
  readonly signal: AbortSignal;
  now(): Date;
}

export interface Command<I, R, D extends DomainEvent = never> {
  readonly name: string;
  run(input: I, ctx: Context): AsyncIterable<CoreEvent<R, D>>;
}

/** Exit codes shared by every cmdeagle tool. Tools add their own from 10 up. */
export const ExitCode = {
  ok: 0,
  failure: 1,
  invalidInput: 2,
  notFound: 3,
  inputRequired: 4,
  cancelled: 130,
} as const;
