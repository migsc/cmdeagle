import type { Command } from "../../../packages/protocol/src/index.ts";

export interface ForecastInput {
  location: string;
  days: number;
  units: "c" | "f";
  cached: boolean;
}

export interface Place {
  name: string;
  region?: string;
  country: string;
  latitude: number;
  longitude: number;
}

export interface Day {
  date: string;
  high: number;
  low: number;
  precipitationChance: number;
  /** WMO weather code, as returned by Open-Meteo. */
  code: number;
  summary: string;
}

export interface ForecastResult {
  place: Place;
  units: "c" | "f";
  fetchedAt: string;
  fromCache: boolean;
  days: Day[];
}

/** Errors: E_LOCATION_NOT_FOUND (3), E_OFFLINE (10, hint: --cached), E_NO_CACHE (11). */
export type ForecastCommand = Command<ForecastInput, ForecastResult>;
