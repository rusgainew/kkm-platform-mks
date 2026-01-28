/**
 * Recent Companies Activity - Table with recent companies
 */

'use client';

import React from 'react';
import { Company } from '@/types/entities';
import { Calendar, Users, Building2 } from 'lucide-react';

interface RecentCompaniesActivityProps {
  companies: Company[];
}

export default function RecentCompaniesActivity({ companies }: RecentCompaniesActivityProps) {
  const getStatusBadgeColor = (status: string) => {
    const colors: Record<string, string> = {
      active: 'bg-emerald-900/20 text-emerald-300 border border-emerald-800',
      inactive: 'bg-gray-800/50 text-gray-300 border border-gray-700',
      suspended: 'bg-orange-900/20 text-orange-300 border border-orange-800',
    };
    return colors[status] || colors.inactive;
  };

  const getStatusLabel = (status: string) => {
    const labels: Record<string, string> = {
      active: 'Активна',
      inactive: 'Неактивна',
      suspended: 'Заблокирована',
    };
    return labels[status] || status;
  };

  return (
    <div className="bg-gray-800/30 border border-gray-700 rounded-xl p-6 backdrop-blur-sm">
      <div className="flex items-center gap-3 mb-6">
        <div className="p-2 bg-green-900/30 rounded-lg">
          <Building2 className="w-5 h-5 text-green-400" />
        </div>
        <div>
          <h3 className="text-lg font-bold text-white">Недавние компании</h3>
          <p className="text-gray-400 text-sm">Последние добавленные</p>
        </div>
      </div>

      <div className="space-y-3">
        {companies.slice(0, 5).map((company) => (
          <div
            key={company.id}
            className="flex items-center justify-between p-3 bg-gray-900/50 rounded-lg hover:bg-gray-900/70 transition"
          >
            <div className="flex-1 min-w-0">
              <p className="text-white font-medium truncate">{company.name}</p>
              <div className="flex items-center gap-4 mt-1 text-gray-400 text-sm">
                <div className="flex items-center gap-1">
                  <Users size={14} />
                  <span>{company.member_count || 0} членов</span>
                </div>
                <div className="flex items-center gap-1">
                  <Calendar size={14} />
                  <span>{new Date(company.created_at * 1000).toLocaleDateString('ru-RU')}</span>
                </div>
              </div>
            </div>
            <span
              className={`px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap ml-4 ${getStatusBadgeColor(
                company.status || 'active'
              )}`}
            >
              {getStatusLabel(company.status || 'active')}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
