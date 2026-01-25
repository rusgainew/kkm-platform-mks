'use client';

import { AlertCircle, AlertTriangle } from 'lucide-react';
import { InventoryItem } from '@/types';

interface InventoryAlertsProps {
  items: InventoryItem[];
  title?: string;
}

const statusConfig = {
  critical: {
    color: 'bg-red-500/10 border-red-500/50',
    badge: 'bg-red-500/20 text-red-400 border border-red-500/50',
    label: 'Критично',
  },
  low: {
    color: 'bg-yellow-500/10 border-yellow-500/50',
    badge: 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/50',
    label: 'Низко',
  },
  ok: {
    color: 'bg-emerald-500/10 border-emerald-500/50',
    badge: 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/50',
    label: 'В норме',
  },
};

export default function InventoryAlerts({ items, title = 'Товары заканчиваются' }: InventoryAlertsProps) {
  const alertItems = items.filter((item) => item.status !== 'ok');

  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 p-6 shadow-xl">
      <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
        <span className="w-1 h-6 bg-linear-to-b from-blue-400 to-emerald-400 rounded-full"></span>
        {title}
      </h3>

      {alertItems.length === 0 ? (
        <div className="text-center py-8">
          <div className="text-green-600 mb-2">✓</div>
          <p className="text-sm text-gray-600">Все товары в норме</p>
        </div>
      ) : (
        <div className="space-y-3">
          {alertItems.map((item) => {
            const config = statusConfig[item.status];
            const Icon = item.status === 'critical' ? AlertCircle : AlertTriangle;

            return (
              <div key={item.id} className={`rounded-lg border p-4 ${config.color} backdrop-blur-sm`}>
                <div className="flex items-start gap-3">
                  <Icon className="w-5 h-5 text-gray-300 shrink-0 mt-0.5" />
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <p className="font-semibold text-white">{item.name}</p>
                      <span className={`text-xs font-semibold px-2 py-1 rounded-full ${config.badge}`}>
                        {config.label}
                      </span>
                    </div>
                    <p className="text-sm text-gray-400">
                      На складе:{' '}
                      <span className="font-semibold text-white">{item.currentStock}</span> шт. (
                      минимум: {item.minThreshold})
                    </p>
                    <div className="w-full bg-gray-300 rounded-full h-2 mt-2">
                      <div
                        className={`h-2 rounded-full ${
                          item.status === 'critical'
                            ? 'bg-red-600'
                            : item.status === 'low'
                              ? 'bg-yellow-600'
                              : 'bg-green-600'
                        }`}
                        style={{
                          width: `${Math.min((item.currentStock / item.minThreshold) * 100, 100)}%`,
                        }}
                      />
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
