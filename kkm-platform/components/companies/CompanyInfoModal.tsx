/**
 * Company Info Modal - Using same style as UserInfoModal
 */

'use client';

import React, { memo } from 'react';
import { Company } from '@/types/company';
import { Building2, Calendar, Users, Tag } from 'lucide-react';
import { formatDate } from '@/lib/utils/dateFormatter';

interface CompanyInfoModalProps {
  isOpen: boolean;
  company: Company | null;
  onClose: () => void;
}

const CompanyInfoModal = memo(function CompanyInfoModal({ isOpen, company, onClose }: CompanyInfoModalProps) {
  if (!isOpen || !company) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-gray-900 rounded-lg border border-gray-700 max-w-md w-full mx-4">
        {/* Header */}
        <div className="border-b border-gray-700 px-6 py-4">
          <h2 className="text-lg font-bold text-white">Информация о компании</h2>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Company Name */}
          <div>
            <div className="flex items-center gap-2 text-gray-400 text-sm font-medium mb-2">
              <Building2 size={16} />
              Название
            </div>
            <p className="text-white font-medium">{company.name}</p>
          </div>

          {/* Description */}
          {company.description && (
            <div>
              <div className="text-gray-400 text-sm font-medium mb-2">Описание</div>
              <p className="text-gray-300">{company.description}</p>
            </div>
          )}

          {/* Owner ID */}
          <div>
            <div className="text-gray-400 text-sm font-medium mb-2">Владелец</div>
            <p className="text-gray-300 text-sm font-mono">{company.owner_id}</p>
          </div>

          {/* Members */}
          <div>
            <div className="flex items-center gap-2 text-gray-400 text-sm font-medium mb-2">
              <Users size={16} />
              Членов
            </div>
            <p className="text-white font-medium">{company.member_count || 0}</p>
          </div>

          {/* Status */}
          <div>
            <div className="flex items-center gap-2 text-gray-400 text-sm font-medium mb-2">
              <Tag size={16} />
              Статус
            </div>
            <span className={`px-3 py-1 rounded-full text-xs font-medium inline-block ${
              company.status === 'active'
                ? 'bg-emerald-900/20 text-emerald-300 border border-emerald-800'
                : company.status === 'suspended'
                ? 'bg-orange-900/20 text-orange-300 border border-orange-800'
                : 'bg-gray-800/50 text-gray-300 border border-gray-700'
            }`}>
              {company.status === 'active' ? 'Активна' : 
               company.status === 'suspended' ? 'Заблокирована' : 'Неактивна'}
            </span>
          </div>

          {/* Created Date */}
          <div>
            <div className="flex items-center gap-2 text-gray-400 text-sm font-medium mb-2">
              <Calendar size={16} />
              Создана
            </div>
            <p className="text-gray-300 text-sm">{formatDate(company.created_at)}</p>
          </div>

          {/* Updated Date */}
          <div>
            <div className="flex items-center gap-2 text-gray-400 text-sm font-medium mb-2">
              <Calendar size={16} />
              Обновлена
            </div>
            <p className="text-gray-300 text-sm">{formatDate(company.updated_at)}</p>
          </div>
        </div>

        {/* Footer */}
        <div className="border-t border-gray-700 px-6 py-4">
          <button
            onClick={onClose}
            className="w-full px-4 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition"
          >
            Закрыть
          </button>
        </div>
      </div>
    </div>
  );
});

export default CompanyInfoModal;
