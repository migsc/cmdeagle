import type { Command } from "../../../packages/protocol/src/index.ts";

/** Phase 2 contracts. Phase 1 passes the Python script's JSON through unchanged. */

export interface Project {
  id: string;
  title: string;
  folder: string;
}

export interface Thread {
  id: string;
  projectId: string;
  title: string;
  status: "idle" | "running" | "awaiting-approval" | "failed";
}

export interface StatusResult {
  server: { running: boolean; origin?: string };
  projects: Project[];
  threads: Thread[];
}

export interface SendResult {
  threadId: string;
  accepted: boolean;
}

export type StatusCommand = Command<Record<string, never>, StatusResult>;
/**
 * A `cost` effect: never runs without --personal-account-confirmed, which an
 * agent must get from the user. Accepted is not finished; poll `status`.
 */
export type SendCommand = Command<{ thread: string; prompt: string }, SendResult>;
