/**
 * Типы для Rate Limiting UI
 */

export interface RateLimitInfo {
  endpoint: string;
  limit: number; // requests per window
  remaining: number;
  reset: number; // timestamp
  windowMs: number;
}

export interface RateLimitStatus {
  endpoint: string;
  requestsUsed: number;
  requestsLimit: number;
  percentageUsed: number;
  resetAt: Date;
  isNearLimit: boolean; // > 80%
  isLimited: boolean; // 100%
}

export interface RateLimitConfig {
  enabled: boolean;
  requests: number;
  windowMs: number;
  endpoints?: {
    [key: string]: {
      requests: number;
      windowMs: number;
    };
  };
}

export interface RateLimitHistoryEntry {
  timestamp: Date;
  endpoint: string;
  status: "success" | "limited";
  remainingRequests: number;
}
