import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import { AuditLogEntry } from "@/lib/hooks/useAudit";

interface AuditStoreState {
  logs: AuditLogEntry[];
  addLog: (log: AuditLogEntry) => void;
  getLogs: () => AuditLogEntry[];
  getLogsByUser: (userId: string) => AuditLogEntry[];
  getLogsByResource: (resource: string, resourceId?: string) => AuditLogEntry[];
  getLogsByAction: (action: string) => AuditLogEntry[];
  getLogsByDateRange: (startDate: number, endDate: number) => AuditLogEntry[];
  clearOldLogs: (daysToKeep: number) => void;
  exportLogs: (format: "json" | "csv") => string;
  searchLogs: (query: string) => AuditLogEntry[];
  getStatistics: () => {
    totalActions: number;
    actionsByType: Record<string, number>;
    actionsByUser: Record<string, number>;
    actionsByResource: Record<string, number>;
    failureCount: number;
    successCount: number;
  };
}

const MAX_LOGS_IN_MEMORY = 10000; // Максимум логов в памяти

/**
 * IndexedDB для персистентного хранилища логов
 */
const idbStorage = {
  getItem: async (_name?: string): Promise<string | null> => {
    if (typeof window === "undefined") return null;

    return new Promise((resolve) => {
      const request = indexedDB.open("kkm-audit", 1);

      request.onupgradeneeded = () => {
        request.result.createObjectStore("logs", { keyPath: "id" });
      };

      request.onsuccess = () => {
        const db = request.result;
        const transaction = db.transaction("logs", "readonly");
        const store = transaction.objectStore("logs");
        const getAllRequest = store.getAll();

        getAllRequest.onsuccess = () => {
          resolve(JSON.stringify({ logs: getAllRequest.result }));
        };

        getAllRequest.onerror = () => {
          resolve(null);
        };
      };

      request.onerror = () => {
        resolve(null);
      };
    });
  },

  setItem: async (_name: string, value: string) => {
    if (typeof window === "undefined") return;

    const data = JSON.parse(value);
    const request = indexedDB.open("kkm-audit", 1);

    request.onsuccess = () => {
      const db = request.result;
      const transaction = db.transaction("logs", "readwrite");
      const store = transaction.objectStore("logs");

      // Очистить старые и добавить новые
      store.clear();
      data.logs?.forEach((log: AuditLogEntry) => {
        store.add(log);
      });
    };
  },

  removeItem: async (_name?: string) => {
    if (typeof window === "undefined") return;

    const request = indexedDB.open("kkm-audit", 1);

    request.onsuccess = () => {
      const db = request.result;
      const transaction = db.transaction("logs", "readwrite");
      const store = transaction.objectStore("logs");
      store.clear();
    };
  },
};

export const useAuditStore = create<AuditStoreState>()(
  persist(
    (set, get) => ({
      logs: [],

      addLog: (log: AuditLogEntry) => {
        set((state) => {
          const newLogs = [log, ...state.logs];
          // Ограничиваем количество логов в памяти
          if (newLogs.length > MAX_LOGS_IN_MEMORY) {
            newLogs.splice(MAX_LOGS_IN_MEMORY);
          }
          return { logs: newLogs };
        });
      },

      getLogs: () => get().logs,

      getLogsByUser: (userId: string) => {
        return get().logs.filter((log) => log.userId === userId);
      },

      getLogsByResource: (resource: string, resourceId?: string) => {
        return get().logs.filter(
          (log) =>
            log.resource === resource &&
            (!resourceId || log.resourceId === resourceId),
        );
      },

      getLogsByAction: (action: string) => {
        return get().logs.filter((log) => log.action === action);
      },

      getLogsByDateRange: (startDate: number, endDate: number) => {
        return get().logs.filter(
          (log) => log.timestamp >= startDate && log.timestamp <= endDate,
        );
      },

      clearOldLogs: (daysToKeep: number) => {
        const cutoffTime = Date.now() - daysToKeep * 24 * 60 * 60 * 1000;
        set((state) => ({
          logs: state.logs.filter((log) => log.timestamp >= cutoffTime),
        }));
      },

      searchLogs: (query: string) => {
        const lowerQuery = query.toLowerCase();
        return get().logs.filter(
          (log) =>
            log.resourceName?.toLowerCase().includes(lowerQuery) ||
            log.userName.toLowerCase().includes(lowerQuery) ||
            log.userEmail.toLowerCase().includes(lowerQuery) ||
            log.resourceId.toLowerCase().includes(lowerQuery) ||
            log.action.toLowerCase().includes(lowerQuery),
        );
      },

      exportLogs: (format: "json" | "csv") => {
        const logs = get().logs;

        if (format === "json") {
          return JSON.stringify(logs, null, 2);
        } else if (format === "csv") {
          const headers = [
            "ID",
            "User",
            "Email",
            "Action",
            "Resource",
            "Resource ID",
            "Status",
            "IP Address",
            "Timestamp",
            "Duration (ms)",
            "Error",
          ];

          const rows = logs.map((log) => [
            log.id,
            log.userName,
            log.userEmail,
            log.action,
            log.resource,
            log.resourceId,
            log.status,
            log.ipAddress,
            new Date(log.timestamp).toISOString(),
            log.duration?.toString() || "-",
            log.error || "-",
          ]);

          const csv = [
            headers.join(","),
            ...rows.map((row) => row.map((cell) => `"${cell}"`).join(",")),
          ].join("\n");

          return csv;
        }

        return "";
      },

      getStatistics: () => {
        const logs = get().logs;

        const stats = {
          totalActions: logs.length,
          actionsByType: {} as Record<string, number>,
          actionsByUser: {} as Record<string, number>,
          actionsByResource: {} as Record<string, number>,
          failureCount: 0,
          successCount: 0,
        };

        logs.forEach((log) => {
          // По типам действий
          stats.actionsByType[log.action] =
            (stats.actionsByType[log.action] || 0) + 1;

          // По пользователям
          stats.actionsByUser[log.userName] =
            (stats.actionsByUser[log.userName] || 0) + 1;

          // По ресурсам
          stats.actionsByResource[log.resource] =
            (stats.actionsByResource[log.resource] || 0) + 1;

          // Успехи и ошибки
          if (log.status === "success") {
            stats.successCount++;
          } else {
            stats.failureCount++;
          }
        });

        return stats;
      },
    }),
    {
      name: "audit-logs",
      storage: createJSONStorage(() => ({
        getItem: async (name) => {
          // Сначала пытаемся получить из IndexedDB, потом из localStorage
          try {
            const idbData = await idbStorage.getItem(name);
            if (idbData) return idbData;
          } catch (e) {
            console.error("[Audit Store] IndexedDB error:", e);
          }

          // Fallback на localStorage
          if (typeof window !== "undefined") {
            return localStorage.getItem(name);
          }
          return null;
        },
        setItem: async (name, value) => {
          try {
            await idbStorage.setItem(name, value);
          } catch (e) {
            console.error("[Audit Store] IndexedDB error:", e);
          }

          // Fallback на localStorage
          if (typeof window !== "undefined") {
            try {
              localStorage.setItem(name, value);
            } catch (e) {
              console.error("[Audit Store] localStorage error:", e);
            }
          }
        },
        removeItem: async (name) => {
          try {
            await idbStorage.removeItem(name);
          } catch (e) {
            console.error("[Audit Store] IndexedDB error:", e);
          }

          if (typeof window !== "undefined") {
            localStorage.removeItem(name);
          }
        },
      })),
      partialize: (state) => ({
        logs: state.logs.slice(0, 5000), // Сохраняем только последние 5000
      }),
    },
  ),
);
