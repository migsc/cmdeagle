import type { Command } from "../../../packages/protocol/src/index.ts";

export interface BackupInput {
  source: string;
  to: string;
  prune: boolean;
}

export interface FileChange {
  path: string;
  bytes: number;
}

export interface BackupResult {
  source: string;
  to: string;
  dryRun: boolean;
  scanned: number;
  added: FileChange[];
  changed: FileChange[];
  /** Empty unless --prune. */
  deleted: FileChange[];
  bytes: number;
  /** The command that applies this plan, present on a dry run. */
  next?: string;
}

/**
 * Asks before deleting (flag: --yes) when --prune applies for real.
 * Errors: E_SOURCE_MISSING (3), E_DEST_UNWRITABLE (10), E_PARTIAL_COPY (11).
 */
export type BackupCommand = Command<BackupInput, BackupResult>;
