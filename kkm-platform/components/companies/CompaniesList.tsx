/**
 * Список компаний с поддержкой:
 * - Поиска по названию и описанию
 * - Фильтрации по статусу (активные, неактивные, заблокированные)
 * - Сортировки (по дате, названию, членам)
 * - CRUD операций (создание, просмотр, редактирование, удаление)
 */

'use client';

import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { Building2, Mail, MapPin, MoreVertical, Search, Edit2, Trash2, AlertCircle, Loader2, Plus, Filter, ChevronDown } from 'lucide-react';
import { useListCompaniesQuery, useDeleteCompanyMutation } from '@/lib/hooks/useCompaniesApi';
import type { Company } from '@/types/entities';
import CompanyInfoModal from './CompanyInfoModal';
import EditCompanyModal from './EditCompanyModal';
import CreateCompanyModal from './CreateCompanyModal';
import { useApiToken } from '@/lib/hooks/useApiToken';
import { formatDate } from '@/lib/utils/dateFormatter';

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

export default function CompaniesList() {
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | 'active' | 'inactive' | 'suspended'>('all');
  const [sortBy, setSortBy] = useState<'name' | 'created' | 'members'>('created');
  const [selectedCompany, setSelectedCompany] = useState<Company | null>(null);
  const [infoModalOpen, setInfoModalOpen] = useState(false);
  const [editModalOpen, setEditModalOpen] = useState(false);
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [deleteConfirmId, setDeleteConfirmId] = useState<string | null>(null);
  const [showFilterMenu, setShowFilterMenu] = useState(false);

  const { data: companies = [], isLoading, error, refetch } = useListCompaniesQuery();
  const deleteMutation = useDeleteCompanyMutation();
  const token = useApiToken();

  // Мемоизация отфильтрованных и отсортированных компаний
  const filteredCompanies = useMemo(() => {
    let result = companies.filter((company: Company) => {
      const matchesSearch = company.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        (company.description?.toLowerCase().includes(searchQuery.toLowerCase()) ?? false);
      const matchesStatus = statusFilter === 'all' || company.status === statusFilter;
      return matchesSearch && matchesStatus;
    });

    // Сортировка
    return [...result].sort((a, b) => {
      switch (sortBy) {
        case 'name':
          return a.name.localeCompare(b.name);
        case 'created':
          return (b.created_at || 0) - (a.created_at || 0);
        case 'members':
          return (b.member_count || 0) - (a.member_count || 0);
        default:
          return 0;
      }
    });
  }, [companies, searchQuery, statusFilter, sortBy]);

  // Оптимизированные обработчики событий
  const handleViewInfo = useCallback((company: Company) => {
    setSelectedCompany(company);
    setInfoModalOpen(true);
  }, []);

  const handleEditCompany = useCallback((company: Company) => {
    setSelectedCompany(company);
    setEditModalOpen(true);
  }, []);

  const handleDeleteCompany = useCallback(async (id: string) => {
    try {
      await deleteMutation.mutateAsync(id);
      setDeleteConfirmId(null);
      refetch();
    } catch (err) {
      console.error('Error deleting company:', err);
    }
  }, [deleteMutation, refetch]);

  const handleCreateSuccess = useCallback(() => {
    setCreateModalOpen(false);
    refetch();
  }, [refetch]);

  const handleEditSuccess = useCallback(() => {
    setEditModalOpen(false);
    refetch();
  }, [refetch]);

  if (isLoading) {
    return (
      <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden p-8">
        <div className="flex items-center justify-center gap-3">
          <Loader2 className="w-6 h-6 text-blue-400 animate-spin" />
          <p className="text-gray-400">Загрузка компаний...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-900/20 border border-red-800 rounded-lg p-4 flex items-start gap-3">
        <AlertCircle className="w-5 h-5 text-red-400 mt-0.5 shrink-0" />
        <div>
          <h3 className="text-red-300 font-medium">Ошибка при загрузке компаний</h3>
          <p className="text-red-200 text-sm mt-1">{error instanceof Error ? error.message : 'Unknown error'}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
      {/* Header with Search and Create Button */}
      <div className="bg-gray-800/50 border-b border-gray-700 p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold text-white">Компании</h2>
          <button
            onClick={() => setCreateModalOpen(true)}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
          >
            <Plus size={18} />
            Создать
          </button>
        </div>

        {/* Search and Filters */}
        <div className="flex gap-3">
          {/* Search */}
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-3 text-gray-500" size={18} />
            <input
              type="text"
              placeholder="Поиск по названию или описанию..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500"
            />
          </div>

          {/* Status Filter */}
          <div className="relative">
            <button
              onClick={() => setShowFilterMenu(!showFilterMenu)}
              className="flex items-center gap-2 px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-gray-300 hover:bg-gray-700 transition"
            >
              <Filter size={18} />
              <span className="text-sm">
                {statusFilter === 'all' ? 'Все статусы' : 
                  statusFilter === 'active' ? 'Активные' :
                  statusFilter === 'inactive' ? 'Неактивные' :
                  'Заблокированные'}
              </span>
              <ChevronDown size={16} />
            </button>

            {showFilterMenu && (
              <div className="absolute right-0 mt-1 w-48 bg-gray-900 border border-gray-700 rounded-lg shadow-lg z-20">
                {[
                  { value: 'all' as const, label: 'Все статусы' },
                  { value: 'active' as const, label: 'Активные' },
                  { value: 'inactive' as const, label: 'Неактивные' },
                  { value: 'suspended' as const, label: 'Заблокированные' },
                ].map((option) => (
                  <button
                    key={option.value}
                    onClick={() => {
                      setStatusFilter(option.value);
                      setShowFilterMenu(false);
                    }}
                    className={`w-full text-left px-4 py-2 hover:bg-gray-800 transition ${
                      statusFilter === option.value ? 'text-blue-400 bg-gray-800' : 'text-gray-300'
                    }`}
                  >
                    {option.label}
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Sort */}
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value as 'name' | 'created' | 'members')}
            className="px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-gray-300 focus:outline-none focus:border-blue-500 text-sm"
          >
            <option value="created">По дате создания</option>
            <option value="name">По названию</option>
            <option value="members">По членам</option>
          </select>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        {filteredCompanies.length === 0 ? (
          <div className="p-8 text-center">
            <Building2 size={48} className="mx-auto text-gray-600 mb-3" />
            <p className="text-gray-400">Компаний не найдено</p>
          </div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-800/50 border-b border-gray-700">
              <tr>
                <th className="px-6 py-4 text-left text-sm font-medium text-gray-300">Название</th>
                <th className="px-6 py-4 text-left text-sm font-medium text-gray-300">Описание</th>
                <th className="px-6 py-4 text-left text-sm font-medium text-gray-300">Членов</th>
                <th className="px-6 py-4 text-left text-sm font-medium text-gray-300">Статус</th>
                <th className="px-6 py-4 text-left text-sm font-medium text-gray-300">Создана</th>
                <th className="px-6 py-4 text-center text-sm font-medium text-gray-300">Действия</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-700">
              {filteredCompanies.map((company) => (
                <tr key={company.id} className="hover:bg-gray-800/30 transition">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-blue-900/30 rounded-lg flex items-center justify-center">
                        <Building2 size={20} className="text-blue-400" />
                      </div>
                      <div>
                        <p className="text-white font-medium">{company.name}</p>
                        <p className="text-gray-400 text-sm">{company.owner_id}</p>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <p className="text-gray-400 text-sm truncate max-w-xs">{company.description || '—'}</p>
                  </td>
                  <td className="px-6 py-4">
                    <p className="text-gray-300">{company.member_count}</p>
                  </td>
                  <td className="px-6 py-4">
                    <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusBadgeColor(company.status || 'active')}`}>
                      {getStatusLabel(company.status || 'active')}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <p className="text-gray-400 text-sm">
                      {formatDate(company.created_at)}
                    </p>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex items-center justify-center gap-2">
                      <button
                        onClick={() => handleViewInfo(company)}
                        className="p-2 hover:bg-gray-700 rounded transition"
                        title="Просмотр информации"
                      >
                        <Building2 size={18} className="text-gray-400" />
                      </button>
                      <button
                        onClick={() => handleEditCompany(company)}
                        className="p-2 hover:bg-gray-700 rounded transition"
                        title="Редактировать"
                      >
                        <Edit2 size={18} className="text-gray-400" />
                      </button>
                      <button
                        onClick={() => setDeleteConfirmId(company.id || company.company_id || null)}
                        className="p-2 hover:bg-gray-700 rounded transition"
                        title="Удалить"
                      >
                        <Trash2 size={18} className="text-red-400" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Delete Confirmation Modal */}
      {deleteConfirmId && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-gray-900 rounded-lg p-6 border border-gray-700 max-w-sm">
            <div className="flex items-center gap-3 mb-4">
              <div className="w-12 h-12 bg-red-900/20 rounded-lg flex items-center justify-center">
                <AlertCircle className="w-6 h-6 text-red-400" />
              </div>
              <h3 className="text-lg font-bold text-white">Удалить компанию?</h3>
            </div>
            <p className="text-gray-400 mb-6">Это действие невозможно отменить.</p>
            <div className="flex gap-3">
              <button
                onClick={() => setDeleteConfirmId(null)}
                className="flex-1 px-4 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition"
              >
                Отмена
              </button>
              <button
                onClick={() => handleDeleteCompany(deleteConfirmId)}
                disabled={deleteMutation.isPending}
                className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition disabled:opacity-50"
              >
                {deleteMutation.isPending ? 'Удаление...' : 'Удалить'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Modals */}
      <CompanyInfoModal
        isOpen={infoModalOpen}
        company={selectedCompany}
        onClose={() => setInfoModalOpen(false)}
      />

      <EditCompanyModal
        isOpen={editModalOpen}
        company={selectedCompany}
        onClose={() => setEditModalOpen(false)}
        onSuccess={handleEditSuccess}
      />

      <CreateCompanyModal
        isOpen={createModalOpen}
        onClose={() => setCreateModalOpen(false)}
        onSuccess={handleCreateSuccess}
      />
    </div>
  );
}
