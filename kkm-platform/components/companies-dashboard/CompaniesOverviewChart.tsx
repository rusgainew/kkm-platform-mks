/**
 * Companies Overview Chart - Shows top companies
 */

'use client';

import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { Building2 } from 'lucide-react';

interface TopCompany {
  name: string;
  members: number;
  revenue?: number;
}

interface CompaniesOverviewChartProps {
  data: TopCompany[];
}

export default function CompaniesOverviewChart({ data }: CompaniesOverviewChartProps) {
  const chartData = data.map((company) => ({
    name: company.name.length > 15 ? company.name.substring(0, 15) + '...' : company.name,
    members: company.members,
    revenue: company.revenue || 0,
  }));

  return (
    <div className="bg-gray-800/30 border border-gray-700 rounded-xl p-6 backdrop-blur-sm">
      <div className="flex items-center gap-3 mb-6">
        <div className="p-2 bg-blue-900/30 rounded-lg">
          <Building2 className="w-5 h-5 text-blue-400" />
        </div>
        <div>
          <h3 className="text-lg font-bold text-white">Топ компании</h3>
          <p className="text-gray-400 text-sm">По количеству членов</p>
        </div>
      </div>

      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={chartData}>
          <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
          <XAxis dataKey="name" stroke="#9CA3AF" />
          <YAxis stroke="#9CA3AF" />
          <Tooltip
            contentStyle={{
              backgroundColor: '#111827',
              border: '1px solid #374151',
              borderRadius: '8px',
            }}
            cursor={{ fill: 'rgba(59, 130, 246, 0.1)' }}
          />
          <Legend />
          <Bar dataKey="members" fill="#3b82f6" name="Членов" radius={[8, 8, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
