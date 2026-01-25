'use client';

import React, { ReactNode, useEffect } from 'react';
import { useWebSocket } from '@/lib/hooks/useWebSocket';
import { useRealtimeInvoices, useRealtimeCatalog } from '@/store/realtime';
import toast from 'react-hot-toast';

interface RealtimeProviderProps {
  children: ReactNode;
  wsUrl?: string;
  enabled?: boolean;
}

export function RealtimeProvider({
  children,
  wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'localhost:8080/ws',
  enabled = true,
}: RealtimeProviderProps) {
  const syncInvoices = useRealtimeInvoices((state) => state.syncFromMessage);
  const syncCatalog = useRealtimeCatalog((state) => state.syncFromMessage);

  const { isConnected, subscribe } = useWebSocket({
    url: wsUrl,
    onConnect: () => {
      console.log('[Realtime] WebSocket connected');
      toast.success('Подключено к реал-тайм обновлениям', {
        duration: 2000,
        position: 'top-right',
      });
    },
    onDisconnect: () => {
      console.log('[Realtime] WebSocket disconnected');
      toast.error('Потеряна связь с реал-тайм сервером', {
        duration: 3000,
        position: 'top-right',
      });
    },
    onError: (error) => {
      console.error('[Realtime] WebSocket error:', error);
      toast.error('Ошибка подключения к реал-тайм серверу', {
        duration: 3000,
        position: 'top-right',
      });
    },
    onMessage: (message) => {
      // Синхронизируем данные в сторы
      if (message.type === 'invoice') {
        syncInvoices(message);

        // Уведомления для счетов
        if (message.action === 'create') {
          toast.success(`Новый счет #${message.data?.number}`, {
            duration: 3000,
            icon: '📄',
          });
        } else if (message.action === 'update') {
          const status = message.data?.status;
          if (status === 'paid') {
            toast.success(`Счет оплачен!`, {
              duration: 3000,
              icon: '✅',
            });
          } else if (status === 'cancelled') {
            toast.error(`Счет отменен`, {
              duration: 3000,
              icon: '❌',
            });
          }
        }
      } else if (message.type === 'catalog') {
        syncCatalog(message);

        // Уведомления для товаров
        if (message.action === 'create') {
          toast.success(`Добавлен товар: ${message.data?.name}`, {
            duration: 3000,
            icon: '📦',
          });
        } else if (message.action === 'update') {
          const stock = message.data?.stock;
          if (stock !== undefined && stock < 10) {
            toast((t) => (
              <div className="flex items-center gap-2">
                <span>⚠️</span>
                <span>Низкий остаток: {message.data?.name}</span>
              </div>
            ), {
              duration: 4000,
            });
          } else {
            console.log(`Товар обновлен: ${message.data?.name}`);
          }
        }
      }
    },
  });

  useEffect(() => {
    if (!enabled) return;

    // Подписываемся на обновления счетов
    const unsubscribeInvoices = subscribe('invoice', (message) => {
      console.log('[Realtime] Invoice update:', message);
    });

    // Подписываемся на обновления товаров
    const unsubscribeCatalog = subscribe('catalog', (message) => {
      console.log('[Realtime] Catalog update:', message);
    });

    return () => {
      unsubscribeInvoices();
      unsubscribeCatalog();
    };
  }, [enabled, subscribe]);

  return (
    <div className="realtime-provider" data-connected={isConnected}>
      {/* Индикатор подключения (необязательно) */}
      {false && (
        <div className="fixed bottom-4 right-4 text-xs">
          {isConnected ? (
            <div className="flex items-center gap-2 bg-green-100 text-green-700 px-3 py-1 rounded-full">
              <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
              Live
            </div>
          ) : (
            <div className="flex items-center gap-2 bg-gray-100 text-gray-700 px-3 py-1 rounded-full">
              <div className="w-2 h-2 bg-gray-400" />
              Offline
            </div>
          )}
        </div>
      )}
      {children}
    </div>
  );
}
