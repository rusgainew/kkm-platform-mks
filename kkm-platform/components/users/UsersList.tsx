'use client';

import React, { useState, useEffect, useMemo } from 'react';
import {
  User,
  Mail,
  Phone,
  Search,
  Edit2,
  Trash2,
  AlertCircle,
  Loader2,
  Plus,
} from 'lucide-react';
import type { ApiUser } from '@/lib/api/users';
import { listUsersQuery } from '@/lib/api/users';
import { formatDateTime } from '@/lib/utils/dateFormatter';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import ExportUsersButton from '@/features/users/components/ExportUsersButton';
import { useAuthStore } from '@/store/authStore';

interface User {
  id: string;
  name: string;
  email: string;
  phone: string;
  role: 'admin' | 'manager' | 'cashier' | 'viewer';
  store: string;
  status: 'active' | 'inactive' | 'suspended';
  lastLogin: string;
  createdAt: string;
}

const getRoleBadgeColor = (role: string) => {
  const colors: Record<string, string> = {
    admin: 'bg-red-900/20 text-red-300 border border-red-800',
    manager: 'bg-blue-900/20 text-blue-300 border border-blue-800',
    cashier: 'bg-green-900/20 text-green-300 border border-green-800',
    viewer: 'bg-gray-800/50 text-gray-300 border border-gray-700',
  };
  return colors[role] || colors.viewer;
};

const getRoleLabel = (role: string) => {
  const labels: Record<string, string> = {
    admin: 'Администратор',
    manager: 'Менеджер',
    cashier: 'Кассир',
    viewer: 'Просмотр',
  };
  return labels[role] || role;
};

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
    active: 'Активен',
    inactive: 'Неактивен',
    suspended: 'Заблокирован',
  };
  return labels[status] || status;
};

export default function UsersList() {
  const router = useRouter();
  const currentUser = useAuthStore((s) => s.user);
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [roleFilter, setRoleFilter] = useState<string>('all');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [apiUsers, setApiUsers] = useState<ApiUser[]>([]);

  // Проверяем роль пользователя
  useEffect(() => {
    if (currentUser) {
      const userRole = currentUser.role;
      // Если роль не admin и не manager - перенаправляем на страницу профиля
      if (userRole !== 'admin' && userRole !== 'manager') {
        console.log('[UsersList] User role is not admin/manager, redirecting to profile');
        router.replace(`/users/${currentUser.id}`);
      }
    }
  }, [currentUser, router]);

  // Load users from API
  useEffect(() => {
    const loadUsers = async () => {
      try {
        setIsLoading(true);
        setError(null);
        
        // Call API to fetch users
        const response = await listUsersQuery();
        
        // Handle different response formats from API
        let apiUsersData: ApiUser[] = [];
        
        // Check if response has a 'users' field with array
        if (response && Array.isArray(response.users)) {
          apiUsersData = response.users;
        }
        // Check if response is an array
        else if (Array.isArray(response)) {
          apiUsersData = response;
        }
        else {
          console.warn('[UsersList] Unexpected API response format:', response);
          apiUsersData = [];
        }
        
        console.log('[UsersList] Loaded users count:', apiUsersData.length);
        console.log('[UsersList] Users data:', apiUsersData.map(u => ({ user_id: u.user_id, name: `${u.first_name} ${u.last_name}` })));
        
        // Debug: log full structure of first user
        if (apiUsersData.length > 0) {
          console.log('[UsersList] Full structure of first user:', apiUsersData[0]);
          console.log('[UsersList] Available keys:', Object.keys(apiUsersData[0]));
        }
        
        setApiUsers(apiUsersData);
        
        // Transform API users to component users format
        const transformedUsers: User[] = apiUsersData
          .map((apiUser: ApiUser, index: number) => {
            // Fallback ID: use user_id from API
            const userId = apiUser.user_id || apiUser.email || `temp-${index}-${Date.now()}`;
            
            if (!apiUser.user_id) {
              console.warn('[UsersList] User missing user_id, using fallback:', {
                fallbackId: userId,
                userData: apiUser
              });
            }
            
            return {
              id: userId,
              name: `${apiUser.first_name} ${apiUser.last_name}`.trim(),
              email: apiUser.email,
              phone: apiUser.phone || '+7 (999) 000-00-00',
              role: (apiUser.role as 'admin' | 'manager' | 'cashier' | 'viewer') || 'viewer',
              store: 'Магазин',
              status: (apiUser.status === 'active' ? 'active' : 'inactive') as 'active' | 'inactive' | 'suspended',
              lastLogin: formatDateTime((apiUser.last_login_at || apiUser.updated_at) as string | number),
              createdAt: formatDateTime(apiUser.created_at as string | number),
            };
          });
        
        // Check for duplicate IDs and log warning
        const ids = new Set<string>();
        const duplicates: string[] = [];
        transformedUsers.forEach(user => {
          if (ids.has(user.id)) {
            duplicates.push(user.id);
          }
          ids.add(user.id);
        });
        
        if (duplicates.length > 0) {
          console.warn('[UsersList] Duplicate user IDs detected:', duplicates);
        }
        
        setUsers(transformedUsers);
      } catch (err) {
        let errorMessage = 'Ошибка загрузки пользователей';
        
        if (err instanceof Error) {
          const msg = err.message.toLowerCase();
          
          // Если недостаточно прав - перенаправляем на страницу профиля
          if (msg.includes('insufficient permissions') || msg.includes('forbidden') || msg.includes('403')) {
            console.log('[UsersList] Insufficient permissions, redirecting to profile');
            if (currentUser?.id) {
              router.replace(`/users/${currentUser.id}`);
              return;
            }
            errorMessage = 'У вас нет прав для просмотра списка пользователей';
          } else {
            errorMessage = err.message;
          }
        } else if (typeof err === 'string') {
          errorMessage = err;
        } else if (err && typeof err === 'object' && 'message' in err) {
          errorMessage = String((err as Record<string, unknown>).message);
        }
        
        console.error('[UsersList] Error loading users:', errorMessage);
        setError(errorMessage);
      } finally {
        setIsLoading(false);
      }
    };

    loadUsers();
  }, []);

  // Оптимизированный обработчик refresh

  // Мемоизация отфильтрованных пользователей
  const filteredUsers = useMemo(() => {
    return users.filter((user) => {
      const matchesSearch =
        user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        user.email.toLowerCase().includes(searchQuery.toLowerCase());
      const matchesRole = roleFilter === 'all' || user.role === roleFilter;
      const matchesStatus = statusFilter === 'all' || user.status === statusFilter;

      return matchesSearch && matchesRole && matchesStatus;
    });
  }, [users, searchQuery, roleFilter, statusFilter]);

  if (isLoading) {
    return (
      <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden p-8">
        <div className="flex items-center justify-center gap-3">
          <Loader2 className="w-6 h-6 text-blue-400 animate-spin" />
          <p className="text-gray-400">Загрузка пользователей...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
      {/* Error Alert */}
      {error && (
        <div className="p-4 bg-red-900/20 border-b border-red-800 flex items-center gap-3">
          <AlertCircle className="w-5 h-5 text-red-400 shrink-0" />
          <div>
            <p className="text-red-300 font-semibold">Ошибка загрузки</p>
            <p className="text-red-400 text-sm">{error}</p>
          </div>
        </div>
      )}

      {/* Header with Controls */}
      <div className="p-6 border-b border-gray-800">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-2xl font-bold text-white flex items-center gap-2">
            <User className="w-6 h-6 text-blue-400" />
            Управление пользователями
          </h2>
          <div className="flex items-center gap-2">
            <ExportUsersButton users={apiUsers} isLoading={isLoading} />
            <Link
              href="/users/import"
              className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-semibold transition-all"
            >
              📥 Импорт
            </Link>
            <Link
              href="/users/create"
              className="flex items-center gap-2 px-4 py-2 bg-linear-to-br from-blue-500 to-emerald-500 hover:shadow-lg hover:shadow-blue-500/30 text-white rounded-lg font-semibold transition-all"
            >
              <Plus className="w-5 h-5" />
              Создать пользователя
            </Link>
          </div>
        </div>

        {/* Filters */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {/* Search */}
          <div className="relative">
            <Search className="w-5 h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
            <input
              type="text"
              placeholder="Поиск по имени или email..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500 transition-colors"
            />
          </div>

          {/* Role Filter */}
          <select
            value={roleFilter}
            onChange={(e) => setRoleFilter(e.target.value)}
            className="px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-blue-500 transition-colors"
          >
            <option value="all">Все роли</option>
            <option value="admin">Администратор</option>
            <option value="manager">Менеджер</option>
            <option value="cashier">Кассир</option>
            <option value="viewer">Просмотр</option>
          </select>

          {/* Status Filter */}
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-blue-500 transition-colors"
          >
            <option value="all">Все статусы</option>
            <option value="active">Активен</option>
            <option value="inactive">Неактивен</option>
            <option value="suspended">Заблокирован</option>
          </select>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-800/50 border-b border-gray-700">
            <tr>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                Пользователь
              </th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                Контакты
              </th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                Роль
              </th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                Магазин
              </th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                Статус
              </th>
              <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                Последний вход
              </th>
              <th className="px-6 py-4 text-center text-sm font-semibold text-gray-300">
                Действия
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800">
            {filteredUsers.length > 0 ? (
              filteredUsers.map((user) => (
                <tr
                  key={user.id}
                  onClick={() => router.push(`/users/${user.id}`)}
                  className="hover:bg-gray-800/50 transition-colors cursor-pointer"
                >
                  {/* Name */}
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-full bg-linear-to-br from-blue-500 to-emerald-500 flex items-center justify-center text-white font-bold shrink-0">
                        {user.name.charAt(0)}
                      </div>
                      <div>
                        <p className="text-white font-medium">{user.name}</p>
                        <p className="text-xs text-gray-400">ID: {user.id}</p>
                      </div>
                    </div>
                  </td>

                  {/* Contacts */}
                  <td className="px-6 py-4">
                    <div className="space-y-1">
                      <div className="flex items-center gap-2 text-sm text-gray-300">
                        <Mail className="w-4 h-4 text-gray-500" />
                        {user.email}
                      </div>
                      <div className="flex items-center gap-2 text-sm text-gray-400">
                        <Phone className="w-4 h-4 text-gray-500" />
                        {user.phone}
                      </div>
                    </div>
                  </td>

                  {/* Role */}
                  <td className="px-6 py-4">
                    <span
                      className={`inline-flex px-3 py-1 rounded-full text-xs font-semibold ${getRoleBadgeColor(
                        user.role
                      )}`}
                    >
                      {getRoleLabel(user.role)}
                    </span>
                  </td>

                  {/* Store */}
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-2 text-sm text-gray-300">
                      {user.store}
                    </div>
                  </td>

                  {/* Status */}
                  <td className="px-6 py-4">
                    <span
                      className={`inline-flex px-3 py-1 rounded-full text-xs font-semibold ${getStatusBadgeColor(
                        user.status
                      )}`}
                    >
                      {getStatusLabel(user.status)}
                    </span>
                  </td>

                  {/* Last Login */}
                  <td className="px-6 py-4">
                    <p className="text-sm text-gray-400">{user.lastLogin}</p>
                  </td>

                  {/* Actions */}
                  <td className="px-6 py-4" onClick={(e) => e.stopPropagation()}>
                    <div className="flex items-center justify-center gap-2">
                      <Link
                        key={`edit-${user.id}`}
                        href={`/users/${user.id}/edit`}
                        className="p-2 hover:bg-blue-600/20 text-blue-400 rounded-lg transition-colors"
                        title="Редактировать"
                      >
                        <Edit2 className="w-4 h-4" />
                      </Link>
                      <Link
                        key={`role-${user.id}`}
                        href={`/users/${user.id}/change-role`}
                        className="p-2 hover:bg-purple-600/20 text-purple-400 rounded-lg transition-colors"
                        title="Изменить роль"
                      >
                        <User className="w-4 h-4" />
                      </Link>
                      <Link
                        key={`delete-${user.id}`}
                        href={`/users/${user.id}/delete`}
                        className="p-2 hover:bg-red-600/20 text-red-400 rounded-lg transition-colors"
                        title="Удалить"
                      >
                        <Trash2 className="w-4 h-4" />
                      </Link>
                    </div>
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={7} className="px-6 py-8 text-center text-gray-400">
                  Пользователи не найдены
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Footer with Stats */}
      <div className="px-6 py-4 bg-gray-800/50 border-t border-gray-800">
        <p className="text-sm text-gray-400">
          Показано <span className="text-white font-semibold">{filteredUsers.length}</span> из{' '}
          <span className="text-white font-semibold">{users.length}</span> пользователей
        </p>
      </div>
    </div>
  );
}
