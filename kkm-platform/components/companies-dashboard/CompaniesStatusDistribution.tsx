/**
 * Companies Status Distribution - Pie chart of company statuses
 */

'use client';

import React from 'react';
import { PieChart, Pie, Cell, Legend, Tooltip, ResponsiveContainer } from 'recharts';
import { AlertCircle } from 'lucide-react';

interface CompaniesStatusDistributionProps {
  active: number;
  inactive: number;
  suspended: number;
}

export default function CompaniesStatusDistribution({
  active,
  inactive,
  suspended,
}: CompaniesStatusDistributionProps) {
  const data = [
    { name: 'Активных', value: active, color: '#10b981' },
    { name: 'Неактивных', value: inactive, color: '#6b7280' },
    { name: 'Заблокированных', value: suspended, color: '#f97316' },
  ];

  const total = active + inactive + suspended;

  return (
    <div className="bg-gray-800/30 border border-gray-700 rounded-xl p-6 backdrop-blur-sm">
      <div className="flex items-center gap-3 mb-6">
        <div className="p-2 bg-purple-900/30 rounded-lg">
          <AlertCircle className="w-5 h-5 text-purple-400" />
        </div>
        <div>
          <h3 className="text-lg font-bold text-white">Статус компаний</h3>
          <p className="text-gray-400 text-sm">Всего: {total} компаний</p>
        </div>
      </div>

      <ResponsiveContainer width="100%" height={250}>
        <PieChart>
          <Pie
            data={data}
            cx="50%"
            cy="50%"
            innerRadius={80}
            outerRadius={100}
            paddingAngle={2}
            dataKey="value"
          >
            {data.map((entry, index) => (
              <Cell key={`cell-${index}`} fill={entry.color} />
            ))}
          </Pie>
          <Tooltip
            contentStyle={{
              backgroundColor: '#111827',
              border: '1px solid #374151',
              borderRadius: '8px',
            }}
          />
          <Legend />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
}
