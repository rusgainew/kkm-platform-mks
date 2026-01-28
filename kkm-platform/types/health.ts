/**
 * Типы для Health Checks
 */

export type ServiceStatus = "healthy" | "degraded" | "unhealthy" | "unknown";

export interface ServiceHealth {
  name: string;
  status: ServiceStatus;
  version?: string;
  uptime?: number; // seconds
  lastCheck: Date;
  responseTime?: number; // ms
  endpoint: string;
  dependencies?: {
    name: string;
    status: ServiceStatus;
  }[];
  metrics?: {
    cpu?: number; // percentage
    memory?: number; // MB
    requests?: number;
    errors?: number;
  };
  message?: string;
}

export interface HealthCheckHistory {
  timestamp: Date;
  serviceName: string;
  status: ServiceStatus;
  responseTime: number;
}

export interface ServiceGroup {
  name: string;
  services: ServiceHealth[];
}

export interface HealthSummary {
  total: number;
  healthy: number;
  degraded: number;
  unhealthy: number;
  unknown: number;
  averageResponseTime: number;
  uptime: number; // percentage
}
