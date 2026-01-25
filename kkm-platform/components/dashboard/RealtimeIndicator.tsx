'use client';

import { Wifi, Clock } from 'lucide-react';
import { useDashboardStore } from '@/store/dashboardStore';
import { useEffect, useState } from 'react';

export default function RealtimeIndicator() {
  const lastUpdate = useDashboardStore((state) => state.lastUpdate);
  const [timeAgo, setTimeAgo] = useState<string>('');

  useEffect(() => {
    if (!lastUpdate) return;

    const updateTimeAgo = () => {
      const seconds = Math.floor((Date.now() - lastUpdate.getTime()) / 1000);
      if (seconds < 5) {
        setTimeAgo('только что');
      } else if (seconds < 60) {
        setTimeAgo(`${seconds} сек назад`);
      } else {
        setTimeAgo(`${Math.floor(seconds / 60)} мин назад`);
      }
    };

    updateTimeAgo();
    const interval = setInterval(updateTimeAgo, 1000);

    return () => clearInterval(interval);
  }, [lastUpdate]);

  return (
    <div className="flex items-center gap-3 px-4 py-2 bg-gray-800/50 border border-gray-700 rounded-lg backdrop-blur-sm">
      <div className="flex items-center gap-2">
        <div className="relative">
          <Wifi className="w-4 h-4 text-emerald-400" />
          <span className="absolute -top-0.5 -right-0.5 w-2 h-2 bg-emerald-400 rounded-full animate-pulse"></span>
        </div>
        <span className="text-sm font-semibold text-emerald-400">Live</span>
      </div>
      
      {lastUpdate && (
        <>
          <div className="w-px h-4 bg-gray-700"></div>
          <div className="flex items-center gap-1.5 text-gray-400">
            <Clock className="w-3.5 h-3.5" />
            <span className="text-xs">{timeAgo}</span>
          </div>
        </>
      )}
    </div>
  );
}
