import type { Command } from "../../../packages/protocol/src/index.ts";

export interface Window {
  address: string;
  class: string;
  title: string;
  pid: number;
}

export interface Workspace {
  id: number;
  label?: string;
  active: boolean;
  windows: Window[];
}

export interface ListResult {
  workspaces: Workspace[];
}

export interface LabelResult {
  id: number;
  label?: string;
  previous?: string;
}

export interface ClearResult {
  id: number;
  dryRun: boolean;
  close: Window[];
  /** Windows left open because the same app also has windows on another workspace. */
  keep: (Window & { reason: string })[];
  /** Same-user processes ended after their windows closed. */
  terminated: { pid: number; command: string }[];
}

export type ListCommand = Command<Record<string, never>, ListResult>;
export type LabelCommand = Command<{ id?: number; name?: string }, LabelResult>;
/** Asks before closing (flag: --yes). Error: E_HYPRCTL_UNAVAILABLE (10). */
export type ClearCommand = Command<{ id?: number }, ClearResult>;
/** Takes desktop focus; the generated skill tells agents to ask first. */
export type PickCommand = Command<Record<string, never>, { id: number }>;
