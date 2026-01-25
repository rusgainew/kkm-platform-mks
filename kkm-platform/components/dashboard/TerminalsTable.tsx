'use client';

import React, { memo } from 'react';
import { AlertCircle, Wifi, WifiOff } from 'lucide-react';
import { Terminal } from '@/types';

interface TerminalsTableProps {
  terminals: Terminal[];
  title?: string;
}

const statusConfig = {
  online: { color: 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/50', label: 'Online', icon: Wifi },
  offline: { color: 'bg-red-500/20 text-red-400 border border-red-500/50', label: 'Offline', icon: WifiOff },
  paper_out: { color: 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/50', label: 'Нет бумаги', icon: AlertCircle },
};

const TerminalsTable = memo(function TerminalsTable({ terminals, title = 'Статус кассовых точек' }: TerminalsTableProps) {
  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 shadow-xl">
      <div className="p-6 border-b border-gray-800">
        <h3 className="text-lg font-semibold text-white flex items-center gap-2">
          <span className="w-1 h-6 bg-linear-to-b from-blue-400 to-emerald-400 rounded-full"></span>
          {title}
        </h3>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="bg-gray-800/50 border-b border-gray-700">
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-400 uppercase tracking-wider">ID</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-400 uppercase tracking-wider">Адрес</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-400 uppercase tracking-wider">Статус</th>
              <th className="px-6 py-3 text-left text-xs font-semibold text-gray-400 uppercase tracking-wider">Последняя синхронизация</th>
              <th className="px-6 py-3 text-right text-xs font-semibold text-gray-400 uppercase tracking-wider">Выручка</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800">
            {terminals.map((terminal) => {
              const statusInfo = statusConfig[terminal.status];
              const StatusIcon = statusInfo.icon;

              return (
                <tr key={terminal.id} className="hover:bg-gray-800/50 transition-colors">
                  <td className="px-6 py-4 text-sm font-mono font-medium text-blue-400">{terminal.id}</td>
                  <td className="px-6 py-4 text-sm text-gray-300">{terminal.address}</td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold ${statusInfo.color}`}>
                      <StatusIcon className="w-4 h-4" />
                      {statusInfo.label}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-400">{terminal.lastSync}</td>
                  <td className="px-6 py-4 text-sm font-semibold text-right text-emerald-400">
                    {terminal.revenue.toLocaleString('ru-RU')} ₽
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      <div className="px-6 py-4 bg-gray-800/50 border-t border-gray-800 text-xs text-gray-500">
        Показано {terminals.length} из {terminals.length} кассовых точек
      </div>
    </div>
  );
});

export default TerminalsTable;
