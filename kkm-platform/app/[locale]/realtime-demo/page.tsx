'use client';

import React from 'react';
import { RealtimeProvider } from '@/components/realtime/RealtimeProvider';
import { InvoicesListWithRealtime } from '@/components/realtime/InvoicesListWithRealtime';
import { CatalogListWithRealtime } from '@/components/realtime/CatalogListWithRealtime';

export default function RealtimeDemoPage() {
  return (
    <RealtimeProvider wsUrl={process.env.NEXT_PUBLIC_WS_URL || 'localhost:8080/ws'}>
      <div className="min-h-screen bg-gray-50 p-8">
        <div className="max-w-7xl mx-auto">
          {/* Заголовок */}
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-gray-900 mb-2">
              🚀 Реал-тайм платформа
            </h1>
            <p className="text-gray-600">
              Live синхронизация счетов и товаров через WebSocket
            </p>
          </div>

          {/* Вкладки */}
          <div className="space-y-8">
            {/* Секция счетов */}
            <div className="bg-white rounded-lg shadow p-6">
              <InvoicesListWithRealtime />
            </div>

            {/* Секция товаров */}
            <div className="bg-white rounded-lg shadow p-6">
              <CatalogListWithRealtime />
            </div>

            {/* Информация */}
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-6">
              <h3 className="font-semibold text-blue-900 mb-2">ℹ️ О реал-тайм обновлениях</h3>
              <ul className="text-sm text-blue-800 space-y-1">
                <li>✅ Автоматическое переподключение при разрыве связи</li>
                <li>✅ Heartbeat пинги каждые 30 секунд</li>
                <li>✅ Toast уведомления о новых счетах и товарах</li>
                <li>✅ Live отслеживание изменений статуса и остатков</li>
                <li>✅ Экспоненциальная задержка при переподключении</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </RealtimeProvider>
  );
}
