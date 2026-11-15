import type { Command } from "../../../packages/protocol/src/index.ts";

export type Agent = "codex" | "claude";

export interface Session {
  id: string;
  agent: Agent;
  startedAt: string;
  cwd?: string;
  /** First user prompt, trimmed. */
  title: string;
  path: string;
}

export interface Hit {
  session: Session;
  role: "user" | "assistant";
  at: string;
  snippet: string;
}

export interface SearchInput {
  query: string;
  agent: Agent | "any";
  since?: string;
  limit: number;
}

export interface SearchResult {
  hits: Hit[];
  scannedFiles: number;
}

export interface Turn {
  role: "user" | "assistant";
  at: string;
  text: string;
}

export interface ShowResult {
  session: Session;
  turns: Turn[];
}

/** Emits progress while streaming large JSONL files. */
export type SearchCommand = Command<SearchInput, SearchResult>;
export type BrowseCommand = Command<Partial<SearchInput>, SearchResult>;
/** Error: E_SESSION_NOT_FOUND (3). */
export type ShowCommand = Command<{ id: string; role: "all" | "user" | "assistant" }, ShowResult>;
