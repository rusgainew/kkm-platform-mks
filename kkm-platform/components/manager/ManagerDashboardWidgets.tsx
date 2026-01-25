import { BarChart3, TrendingUp, AlertCircle } from 'lucide-react';

interface StatCardProps {
  title: string;
  value: string | number;
  change?: string;
  icon: React.ReactNode;
  trend?: 'up' | 'down';
}

export function StatCard({ title, value, change, icon, trend = 'up' }: StatCardProps) {
  return (
    <div className="bg-gray-900/50 border border-gray-800 rounded-xl p-6 hover:border-blue-500/30 transition-all">
      <div className="flex items-start justify-between mb-4">
        <div>
          <p className="text-gray-400 text-sm mb-1">{title}</p>
          <p className="text-2xl font-bold text-white">{value}</p>
        </div>
        <div className="w-10 h-10 text-blue-500/20 bg-blue-500/10 p-2 rounded-lg">
          {icon}
        </div>
      </div>
      {change && (
        <div className={`flex items-center gap-2 text-sm ${trend === 'up' ? 'text-emerald-400' : 'text-red-400'}`}>
          <span>{trend === 'up' ? '↑' : '↓'} {change}</span>
          <span className="text-gray-500">vs вчера</span>
        </div>
      )}
    </div>
  );
}

export function ManagerAlert({ type, title, message }: { type: 'warning' | 'info'; title: string; message: string }) {
  const colors = {
    warning: 'bg-amber-500/10 border-amber-500/30 text-amber-500',
    info: 'bg-blue-500/10 border-blue-500/30 text-blue-500',
  };

  return (
    <div className={`p-4 ${colors[type]} border rounded-lg flex items-start gap-3`}>
      <AlertCircle className="w-5 h-5 shrink-0 mt-0.5" />
      <div>
        <p className="text-white font-medium">{title}</p>
        <p className="text-gray-400 text-sm">{message}</p>
      </div>
    </div>
  );
}

export function PerformanceMetric({ label, value, unit }: { label: string; value: string | number; unit?: string }) {
  return (
    <div className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
      <span className="text-gray-400">{label}</span>
      <span className="text-white font-semibold">
        {value}
        {unit && <span className="text-gray-500 ml-1">{unit}</span>}
      </span>
    </div>
  );
}
