import type { Command } from "../../../packages/protocol/src/index.ts";

export interface NoteRef {
  /** Path relative to the vault root. */
  path: string;
  title: string;
}

export interface SearchResult {
  notes: (NoteRef & { snippet: string; tags: string[] })[];
}

/** One file the command creates or edits. `before` is absent for new files. */
export interface FileDiff {
  path: string;
  before?: string;
  after: string;
}

export interface WriteResult {
  dryRun: boolean;
  note: NoteRef;
  diffs: FileDiff[];
}

export type SearchCommand = Command<{ query: string; limit: number }, SearchResult>;
/** Error: E_NOTE_EXISTS (10, hint: pick another title or use --folder). */
export type CreateCommand = Command<
  { title: string; folder: string; body?: string; linkDaily: boolean },
  WriteResult
>;
/** Errors: E_NOTE_NOT_FOUND (3), E_ALREADY_LINKED (11). */
export type LinkDailyCommand = Command<{ note: string }, WriteResult>;
