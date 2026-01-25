'use client';

import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { 
  ArrowLeft, 
  Mail, 
  Phone, 
  Calendar, 
  Shield, 
  Activity,
  Edit2,
  Trash2,
  UserCog
} from 'lucide-react';
import Link from 'next/link';
import { getUserById } from '@/lib/api/users';
import type { ApiUser } from '@/lib/api/users';

export default function UserDetailPage() {
  const params = useParams();
  const router = useRouter();
  const userId = params.id as string;
  
  const [user, setUser] = useState<ApiUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const userData = await getUserById(userId);
        setUser(userData);
      } catch (err) {
        console.error('[UserDetail] Error fetching user:', err);
        setError(err instanceof Error ? err.message : 'Failed to load user');
      } finally {
        setIsLoading(false);
      }
    };

    fetchUser();
  }, [userId]);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center">
        <div className="text-white">Загрузка...</div>
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 flex items-center justify-center">
        <div className="text-red-400">Ошибка: {error || 'Пользователь не найден'}</div>
      </div>
    );
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

  const getStatusBadgeColor = (status: string) => {
    const colors: Record<string, string> = {
      active: 'bg-green-900/20 text-green-300 border border-green-800',
      inactive: 'bg-gray-800/50 text-gray-400 border border-gray-700',
      suspended: 'bg-red-900/20 text-red-300 border border-red-800',
    };
    return colors[status] || colors.inactive;
  };

  const formatDate = (dateString: string | number | undefined) => {
    if (!dateString) return 'Нет данных';
    try {
      const date = typeof dateString === 'number' ? new Date(dateString * 1000) : new Date(dateString);
      return date.toLocaleString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
      });
    } catch {
      return String(dateString);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-6 flex items-center justify-between">
          <button
            onClick={() => router.back()}
            className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
          >
            <ArrowLeft className="w-5 h-5" />
            Назад
          </button>
          
          <div className="flex gap-2">
            <Link
              href={`/users/${userId}/edit`}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors flex items-center gap-2"
            >
              <Edit2 className="w-4 h-4" />
              Редактировать
            </Link>
            <Link
              href={`/users/${userId}/change-role`}
              className="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-colors flex items-center gap-2"
            >
              <UserCog className="w-4 h-4" />
              Роль
            </Link>
            <Link
              href={`/users/${userId}/delete`}
              className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors flex items-center gap-2"
            >
              <Trash2 className="w-4 h-4" />
              Удалить
            </Link>
          </div>
        </div>

        {/* Main Card */}
        <div className="bg-gray-800/50 backdrop-blur-sm rounded-xl border border-gray-700 p-8">
          {/* User Header */}
          <div className="flex items-center gap-6 mb-8 pb-8 border-b border-gray-700">
            <div className="w-24 h-24 rounded-full bg-gradient-to-br from-blue-500 to-emerald-500 flex items-center justify-center text-white text-3xl font-bold">
              {user.first_name.charAt(0)}{user.last_name.charAt(0)}
            </div>
            
            <div className="flex-1">
              <h1 className="text-3xl font-bold text-white mb-2">
                {user.first_name} {user.last_name}
              </h1>
              <div className="flex items-center gap-3">
                <span className={`inline-flex px-3 py-1 rounded-full text-sm font-semibold ${getRoleBadgeColor(user.role)}`}>
                  {user.role}
                </span>
                <span className={`inline-flex px-3 py-1 rounded-full text-sm font-semibold ${getStatusBadgeColor(user.status)}`}>
                  {user.status}
                </span>
              </div>
            </div>
          </div>

          {/* User Details */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Email */}
            <div className="flex items-start gap-3">
              <div className="p-3 bg-blue-900/20 rounded-lg">
                <Mail className="w-5 h-5 text-blue-400" />
              </div>
              <div>
                <p className="text-sm text-gray-400 mb-1">Email</p>
                <p className="text-white font-medium">{user.email}</p>
              </div>
            </div>

            {/* Phone */}
            <div className="flex items-start gap-3">
              <div className="p-3 bg-green-900/20 rounded-lg">
                <Phone className="w-5 h-5 text-green-400" />
              </div>
              <div>
                <p className="text-sm text-gray-400 mb-1">Телефон</p>
                <p className="text-white font-medium">{user.phone || 'Не указан'}</p>
              </div>
            </div>

            {/* Role */}
            <div className="flex items-start gap-3">
              <div className="p-3 bg-purple-900/20 rounded-lg">
                <Shield className="w-5 h-5 text-purple-400" />
              </div>
              <div>
                <p className="text-sm text-gray-400 mb-1">Роль</p>
                <p className="text-white font-medium capitalize">{user.role}</p>
              </div>
            </div>

            {/* Status */}
            <div className="flex items-start gap-3">
              <div className="p-3 bg-emerald-900/20 rounded-lg">
                <Activity className="w-5 h-5 text-emerald-400" />
              </div>
              <div>
                <p className="text-sm text-gray-400 mb-1">Статус</p>
                <p className="text-white font-medium capitalize">{user.status}</p>
              </div>
            </div>

            {/* Created At */}
            <div className="flex items-start gap-3">
              <div className="p-3 bg-gray-700/50 rounded-lg">
                <Calendar className="w-5 h-5 text-gray-400" />
              </div>
              <div>
                <p className="text-sm text-gray-400 mb-1">Создан</p>
                <p className="text-white font-medium">{formatDate(user.created_at)}</p>
              </div>
            </div>

            {/* Updated At */}
            <div className="flex items-start gap-3">
              <div className="p-3 bg-gray-700/50 rounded-lg">
                <Calendar className="w-5 h-5 text-gray-400" />
              </div>
              <div>
                <p className="text-sm text-gray-400 mb-1">Обновлен</p>
                <p className="text-white font-medium">{formatDate(user.updated_at)}</p>
              </div>
            </div>

            {/* Last Login */}
            {user.last_login_at && (
              <div className="flex items-start gap-3 md:col-span-2">
                <div className="p-3 bg-yellow-900/20 rounded-lg">
                  <Activity className="w-5 h-5 text-yellow-400" />
                </div>
                <div>
                  <p className="text-sm text-gray-400 mb-1">Последний вход</p>
                  <p className="text-white font-medium">{formatDate(user.last_login_at)}</p>
                </div>
              </div>
            )}

            {/* User ID */}
            <div className="flex items-start gap-3 md:col-span-2">
              <div>
                <p className="text-sm text-gray-400 mb-1">ID пользователя</p>
                <p className="text-gray-500 font-mono text-sm">{user.user_id}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
