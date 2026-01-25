"use client";

/**
 * Audit Logging Hook
 * Отслеживает все действия пользователей в приложении
 * - Create, Update, Delete, View действия
 * - Сохранение старых и новых значений
 * - IP адрес, User Agent, Timestamp
 * - Интеграция с Zustand store
 */

import { useCallback, useEffect, useState } from "react";
import { useAuditStore } from "@/store/audit";

export enum AuditAction {
  CREATE = "create",
  READ = "read",
  UPDATE = "update",
  DELETE = "delete",
  EXPORT = "export",
  IMPORT = "import",
  LOGIN = "login",
  LOGOUT = "logout",
  PERMISSION_CHANGE = "permission_change",
  ROLE_CHANGE = "role_change",
}

export enum AuditResource {
  INVOICE = "invoice",
  PRODUCT = "product",
  COMPANY = "company",
  USER = "user",
  SETTINGS = "settings",
  REPORT = "report",
  DOCUMENT = "document",
}

export interface AuditLogEntry {
  id: string;
  userId: string;
  userName: string;
  userEmail: string;
  action: AuditAction;
  resource: AuditResource;
  resourceId: string;
  resourceName?: string;
  oldValue?: Record<string, any>;
  newValue?: Record<string, any>;
  status: "success" | "failure";
  error?: string;
  ipAddress: string;
  userAgent: string;
  timestamp: number;
  duration?: number; // milliseconds
  metadata?: Record<string, any>;
}

interface UseAuditOptions {
  userId?: string;
  userName?: string;
  userEmail?: string;
}

interface AuditLoggerMethods {
  log: (
    entry: Omit<AuditLogEntry, "id" | "timestamp" | "ipAddress" | "userAgent">
  ) => Promise<void>;
  logAction: (
    action: AuditAction,
    resource: AuditResource,
    resourceId: string,
    metadata?: Record<string, any>
  ) => Promise<void>;
  logChange: (
    action: AuditAction,
    resource: AuditResource,
    resourceId: string,
    oldValue: Record<string, any>,
    newValue: Record<string, any>,
    metadata?: Record<string, any>
  ) => Promise<void>;
}

export function useAudit(options?: UseAuditOptions): AuditLoggerMethods {
  const { addLog } = useAuditStore();
  const [ipAddress, setIpAddress] = useState<string>("unknown");

  // Получаем IP адрес при монтировании
  useEffect(() => {
    const getIP = async () => {
      try {
        // Попробуем получить IP из API (если доступна)
        const response = await fetch("/api/ip", { method: "GET" });
        if (response.ok) {
          const data = await response.json();
          setIpAddress(data.ip);
        }
      } catch {
        // Fallback: используем значение по умолчанию
        setIpAddress("unknown");
      }
    };

    getIP();
  }, []);

  const log = useCallback(
    async (
      entry: Omit<AuditLogEntry, "id" | "timestamp" | "ipAddress" | "userAgent">
    ) => {
      const logEntry: AuditLogEntry = {
        ...entry,
        id: `audit-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        timestamp: Date.now(),
        ipAddress,
        userAgent:
          typeof navigator !== "undefined" ? navigator.userAgent : "unknown",
        userId: options?.userId || entry.userId,
        userName: options?.userName || entry.userName,
        userEmail: options?.userEmail || entry.userEmail,
      };

      try {
        addLog(logEntry);
        console.log("[Audit] Log entry added:", logEntry);
      } catch (error) {
        console.error("[Audit] Failed to add log entry:", error);
      }
    },
    [addLog, ipAddress, options]
  );

  const logAction = useCallback(
    async (
      action: AuditAction,
      resource: AuditResource,
      resourceId: string,
      metadata?: Record<string, any>
    ) => {
      await log({
        action,
        resource,
        resourceId,
        status: "success",
        userId: options?.userId || "anonymous",
        userName: options?.userName || "Anonymous",
        userEmail: options?.userEmail || "unknown@example.com",
        metadata,
      });
    },
    [log, options]
  );

  const logChange = useCallback(
    async (
      action: AuditAction,
      resource: AuditResource,
      resourceId: string,
      oldValue: Record<string, any>,
      newValue: Record<string, any>,
      metadata?: Record<string, any>
    ) => {
      // Получаем только измененные поля
      const changes: Record<string, any> = {};
      for (const key in newValue) {
        if (oldValue[key] !== newValue[key]) {
          changes[key] = {
            from: oldValue[key],
            to: newValue[key],
          };
        }
      }

      await log({
        action,
        resource,
        resourceId,
        oldValue,
        newValue: changes,
        status: "success",
        userId: options?.userId || "anonymous",
        userName: options?.userName || "Anonymous",
        userEmail: options?.userEmail || "unknown@example.com",
        metadata: {
          ...metadata,
          changedFields: Object.keys(changes),
        },
      });
    },
    [log, options]
  );

  return {
    log,
    logAction,
    logChange,
  };
}

/**
 * Хук для логирования с таймером (для отслеживания длительности операций)
 */
export function useAuditTimer(
  userId?: string,
  userName?: string,
  userEmail?: string
) {
  const { log } = useAudit({ userId, userName, userEmail });

  return useCallback(
    async <T>(
      action: AuditAction,
      resource: AuditResource,
      resourceId: string,
      operation: () => Promise<T>,
      metadata?: Record<string, any>
    ): Promise<T> => {
      const startTime = performance.now();

      try {
        const result = await operation();
        const duration = performance.now() - startTime;

        await log({
          action,
          resource,
          resourceId,
          status: "success",
          userId: userId || "anonymous",
          userName: userName || "Anonymous",
          userEmail: userEmail || "unknown@example.com",
          duration: Math.round(duration),
          metadata,
        });

        return result;
      } catch (error) {
        const duration = performance.now() - startTime;

        await log({
          action,
          resource,
          resourceId,
          status: "failure",
          error: error instanceof Error ? error.message : String(error),
          userId: userId || "anonymous",
          userName: userName || "Anonymous",
          userEmail: userEmail || "unknown@example.com",
          duration: Math.round(duration),
          metadata,
        });

        throw error;
      }
    },
    [log]
  );
}
