import type { Command } from "../../../packages/protocol/src/index.ts";

export interface Todo {
  id: number;
  text: string;
  done: boolean;
  createdAt: string;
  doneAt?: string;
  due?: string;
}

export interface ListResult {
  items: Todo[];
  open: number;
  finished: number;
}

export interface ChangeResult {
  /** Items as they are after the change, or would be in a dry run. */
  changed: Todo[];
  removed: Todo[];
  dryRun: boolean;
}

export type BrowseCommand = Command<{ all: boolean }, ListResult>;
export type AddCommand = Command<{ text: string; due?: string }, ChangeResult>;
export type ListCommand = Command<{ all: boolean }, ListResult>;
export type DoneCommand = Command<{ ids: number[] }, ChangeResult>;
export type EditCommand = Command<{ id: number; text: string }, ChangeResult>;
/** Asks for confirmation (flag: --yes) unless --dry-run. Error: E_TODO_NOT_FOUND (3). */
export type RemoveCommand = Command<{ ids: number[] }, ChangeResult>;
export type ClearCommand = Command<Record<string, never>, ChangeResult>;
