"use client";

import { useEffect, useRef, useCallback, useState } from "react";

export interface WebSocketMessage {
  type: "invoice" | "catalog" | "status" | "ping" | "pong";
  action: "create" | "update" | "delete" | "sync";
  data?: Record<string, any>;
  timestamp: number;
  id?: string;
}

interface WebSocketOptions {
  url: string;
  onMessage?: (message: WebSocketMessage) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
  onError?: (error: Event) => void;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
}

interface UseWebSocketReturn {
  isConnected: boolean;
  send: (message: WebSocketMessage) => void;
  subscribe: (
    type: string,
    callback: (data: WebSocketMessage) => void
  ) => () => void;
  disconnect: () => void;
}

export function useWebSocket({
  url,
  onMessage,
  onConnect,
  onDisconnect,
  onError,
  reconnectInterval = 3000,
  maxReconnectAttempts = 5,
}: WebSocketOptions): UseWebSocketReturn {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const subscribersRef = useRef<
    Map<string, Set<(data: WebSocketMessage) => void>>
  >(new Map());
  const [isConnected, setIsConnected] = useState(false);
  const pingIntervalRef = useRef<NodeJS.Timeout | null>(null);

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      return;
    }

    try {
      const wsUrl = url.startsWith("ws") ? url : `ws://${url}`;
      wsRef.current = new WebSocket(wsUrl);

      wsRef.current.onopen = () => {
        console.log("[WebSocket] Connected to", wsUrl);
        setIsConnected(true);
        reconnectAttemptsRef.current = 0;
        onConnect?.();

        // Heartbeat: отправляем ping каждые 30 секунд
        pingIntervalRef.current = setInterval(() => {
          if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(
              JSON.stringify({
                type: "ping",
                action: "sync",
                timestamp: Date.now(),
              })
            );
          }
        }, 30000);
      };

      wsRef.current.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data);
          console.log("[WebSocket] Message received:", message);

          // Handle pong response
          if (message.type === "pong") {
            return;
          }

          // Вызываем глобальный callback
          onMessage?.(message);

          // Вызываем подписчиков по типу
          const subscribers = subscribersRef.current.get(message.type);
          if (subscribers) {
            subscribers.forEach((callback) => callback(message));
          }
        } catch (error) {
          console.error("[WebSocket] Failed to parse message:", error);
        }
      };

      wsRef.current.onerror = (error) => {
        console.error("[WebSocket] Error:", error);
        onError?.(error);
      };

      wsRef.current.onclose = () => {
        console.log("[WebSocket] Disconnected");
        if (pingIntervalRef.current) clearInterval(pingIntervalRef.current);
        setIsConnected(false);
        onDisconnect?.();

        // Автоматический переподключение с экспоненциальной задержкой
        if (reconnectAttemptsRef.current < maxReconnectAttempts) {
          const delay =
            reconnectInterval * Math.pow(2, reconnectAttemptsRef.current);
          console.log(`[WebSocket] Attempting to reconnect in ${delay}ms...`);
          reconnectAttemptsRef.current += 1;

          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, delay);
        } else {
          console.error("[WebSocket] Max reconnection attempts reached");
        }
      };
    } catch (error) {
      console.error("[WebSocket] Connection failed:", error);
      onError?.(error as Event);
    }
  }, [
    url,
    onMessage,
    onConnect,
    onDisconnect,
    onError,
    reconnectInterval,
    maxReconnectAttempts,
  ]);

  const send = useCallback((message: WebSocketMessage) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({
          ...message,
          timestamp: Date.now(),
        })
      );
      console.log("[WebSocket] Message sent:", message);
    } else {
      console.warn("[WebSocket] Connection not ready");
    }
  }, []);

  const subscribe = useCallback(
    (type: string, callback: (data: WebSocketMessage) => void) => {
      if (!subscribersRef.current.has(type)) {
        subscribersRef.current.set(type, new Set());
      }
      subscribersRef.current.get(type)!.add(callback);

      // Return unsubscribe function
      return () => {
        subscribersRef.current.get(type)?.delete(callback);
        if (subscribersRef.current.get(type)?.size === 0) {
          subscribersRef.current.delete(type);
        }
      };
    },
    []
  );

  const disconnect = useCallback(() => {
    if (pingIntervalRef.current) clearInterval(pingIntervalRef.current);
    if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current);

    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }

    setIsConnected(false);
  }, []);

  // Инициализируем подключение
  useEffect(() => {
    connect();

    return () => {
      disconnect();
    };
  }, [connect, disconnect]);

  return {
    isConnected,
    send,
    subscribe,
    disconnect,
  };
}
