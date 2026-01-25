import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { TrendingUp } from 'lucide-react';

interface RevenueTrendData {
  date: string;
  revenue: number;
  transactions: number;
}

interface RevenueTrendChartProps {
  data: RevenueTrendData[];
  title?: string;
}

export default function RevenueTrendChart({
  data,
  title = 'Тренд выручки за неделю',
}: RevenueTrendChartProps) {
  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 p-6 shadow-xl">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-semibold text-white flex items-center gap-2">
          <span className="w-1 h-6 bg-linear-to-b from-blue-400 to-emerald-400 rounded-full"></span>
          {title}
        </h3>
        <div className="flex items-center gap-2 px-3 py-1 bg-emerald-500/20 border border-emerald-500/50 rounded-full">
          <TrendingUp className="w-4 h-4 text-emerald-400" />
          <span className="text-sm font-semibold text-emerald-400">+12.5%</span>
        </div>
      </div>

      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data} margin={{ top: 10, right: 30, left: 0, bottom: 10 }}>
          <defs>
            <linearGradient id="revenueGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="#60A5FA" stopOpacity={0.8} />
              <stop offset="100%" stopColor="#60A5FA" stopOpacity={0.1} />
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
          <XAxis dataKey="date" stroke="#9CA3AF" style={{ fontSize: '12px' }} />
          <YAxis stroke="#9CA3AF" style={{ fontSize: '12px' }} />
          <Tooltip
            contentStyle={{
              backgroundColor: '#1F2937',
              border: '1px solid #374151',
              borderRadius: '8px',
              color: '#F3F4F6',
            }}
          />
          <Legend wrapperStyle={{ color: '#9CA3AF' }} />
          <Line
            type="monotone"
            dataKey="revenue"
            stroke="#60A5FA"
            strokeWidth={3}
            dot={{ fill: '#60A5FA', r: 5 }}
            activeDot={{ r: 7, fill: '#3B82F6' }}
            name="Выручка (₽)"
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
