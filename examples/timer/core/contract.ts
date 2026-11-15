import type { Command } from "../../../packages/protocol/src/index.ts";

export interface StartInput {
  durationMs: number;
  label?: string;
  notify: boolean;
  detach: boolean;
}

export type TimerEvent =
  | { type: "timer.tick"; remainingMs: number; elapsedMs: number }
  | { type: "timer.paused"; remainingMs: number }
  | { type: "timer.resumed"; remainingMs: number };

export interface TimerResult {
  label?: string;
  durationMs: number;
  /** False when stopped early or cancelled with Ctrl-C (exit 130). */
  completed: boolean;
  startedAt: string;
  endedAt: string;
}

export interface StatusResult {
  running: boolean;
  label?: string;
  remainingMs?: number;
  paused?: boolean;
}

export type StartCommand = Command<StartInput, TimerResult, TimerEvent>;
export type StatusCommand = Command<Record<string, never>, StatusResult>;
export type StopCommand = Command<Record<string, never>, TimerResult>;
