'use client';

import React from 'react';
import { Download, Loader2 } from 'lucide-react';
import { listUsersQuery } from '@/lib/api/users';
import type { ApiUser } from '@/lib/api/users';

interface ExportUsersButtonProps {
  users?: ApiUser[];
  isLoading?: boolean;
}

export default function ExportUsersButton({ users = [], isLoading = false }: ExportUsersButtonProps) {
  const [isExporting, setIsExporting] = React.useState(false);

  const handleExportCSV = async () => {
    try {
      setIsExporting(true);

      // Fetch fresh data if not provided
      let usersToExport = users;
      if (!users || users.length === 0) {
        try {
          const response = await listUsersQuery();
          usersToExport = response.users || [];
        } catch (error) {
          console.error('Ошибка загрузки пользователей для экспорта:', error);
          alert('Ошибка при загрузке пользователей');
          return;
        }
      }

      // Create CSV content
      const headers = ['Email', 'Имя', 'Фамилия', 'Роль', 'Статус', 'Дата создания'];
      const csvContent = [
        headers.join(','),
        ...usersToExport.map(user => {
          const ts = typeof user.created_at === 'string' ? parseInt(user.created_at) : user.created_at;
          const createdDate = new Date(ts * 1000).toLocaleDateString('ru-RU');
          return [
            `"${user.email}"`,
            `"${user.first_name}"`,
            `"${user.last_name}"`,
            user.role,
            user.status,
            createdDate,
          ].join(',');
        }),
      ].join('\n');

      // Download CSV
      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
      const link = document.createElement('a');
      const url = URL.createObjectURL(blob);

      link.setAttribute('href', url);
      link.setAttribute('download', `users_export_${new Date().toISOString().split('T')[0]}.csv`);
      link.style.visibility = 'hidden';

      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (error) {
      console.error('Ошибка экспорта:', error);
      alert('Ошибка при экспорте пользователей');
    } finally {
      setIsExporting(false);
    }
  };

  return (
    <button
      onClick={handleExportCSV}
      disabled={isLoading || isExporting}
      className="flex items-center gap-2 px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-semibold transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
    >
      {isExporting ? (
        <>
          <Loader2 className="w-5 h-5 animate-spin" />
          Экспорт...
        </>
      ) : (
        <>
          <Download className="w-5 h-5" />
          Экспортировать CSV
        </>
      )}
    </button>
  );
}
