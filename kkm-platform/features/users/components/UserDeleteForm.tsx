'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowLeft, Loader2, AlertCircle, CheckCircle, Trash2 } from 'lucide-react';
import type { ApiUser } from '@/lib/api/users';
import { deleteUser } from '@/lib/api/users';

interface UserDeleteFormProps {
  user: ApiUser;
}

export default function UserDeleteForm({ user }: UserDeleteFormProps) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState('');
  const [isConfirmed, setIsConfirmed] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!isConfirmed) {
      setError('Вы должны подтвердить удаление пользователя');
      return;
    }

    setIsLoading(true);
    setError(null);
    setSuccessMessage('');

    try {
      await deleteUser(user.id);
      setSuccessMessage('Пользователь успешно удалён');

      // Redirect back to users list after 2 seconds
      setTimeout(() => {
        router.push('/users');
      }, 2000);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : 'Ошибка при удалении пользователя';
      setError(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-900 to-black">
      {/* Header */}
      <div className="bg-gray-800/50 border-b border-gray-700 sticky top-0 z-10">
        <div className="max-w-2xl mx-auto px-4 py-4">
          <button
            onClick={() => router.back()}
            className="flex items-center gap-2 text-gray-400 hover:text-white mb-4 transition-colors"
          >
            <ArrowLeft size={20} />
            Вернуться
          </button>
          <div>
            <h1 className="text-3xl font-bold text-white">Удаление пользователя</h1>
            <p className="text-gray-400 mt-2">
              Email: <span className="text-gray-300">{user.email}</span>
            </p>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-2xl mx-auto px-4 py-8">
        <form onSubmit={handleSubmit} className="space-y-8">
          {/* Success Message */}
          {successMessage && (
            <div className="bg-emerald-900/20 border border-emerald-800 rounded-lg p-4 flex items-start gap-3">
              <CheckCircle size={20} className="text-emerald-400 flex-shrink-0 mt-0.5" />
              <div>
                <p className="text-emerald-300 font-medium">{successMessage}</p>
                <p className="text-emerald-300/70 text-sm mt-1">
                  Перенаправление на список пользователей...
                </p>
              </div>
            </div>
          )}

          {/* Error Message */}
          {error && (
            <div className="bg-red-900/20 border border-red-800 rounded-lg p-4 flex items-start gap-3">
              <AlertCircle size={20} className="text-red-400 flex-shrink-0 mt-0.5" />
              <p className="text-red-300">{error}</p>
            </div>
          )}

          {/* Warning Alert */}
          <div className="bg-red-900/30 border border-red-800 rounded-lg p-6 flex items-start gap-4">
            <Trash2 size={24} className="text-red-400 flex-shrink-0 mt-1" />
            <div className="flex-1">
              <h3 className="text-lg font-semibold text-red-300 mb-2">Внимание!</h3>
              <p className="text-red-300/90 mb-4">
                Вы собираетесь удалить пользователя <span className="font-semibold">{user.email}</span>.
              </p>
              <p className="text-red-300/70 text-sm">
                Это действие <span className="font-semibold">необратимо</span> и удалит все данные, связанные с этим пользователем.
              </p>
            </div>
          </div>

          {/* User Info */}
          <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6">
            <h3 className="text-sm font-medium text-gray-400 mb-4">Информация о пользователе</h3>
            <div className="space-y-3 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-400">Email:</span>
                <span className="text-gray-300">{user.email}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-400">Имя:</span>
                <span className="text-gray-300">{user.first_name || '-'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-400">Фамилия:</span>
                <span className="text-gray-300">{user.last_name || '-'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-400">Роль:</span>
                <span className="text-gray-300">
                  {user.role === 'admin'
                    ? 'Администратор'
                    : user.role === 'manager'
                    ? 'Менеджер'
                    : user.role === 'cashier'
                    ? 'Кассир'
                    : user.role}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-400">ID пользователя:</span>
                <span className="text-gray-300 font-mono text-xs">{user.id}</span>
              </div>
              {user.created_at && (
                <div className="flex justify-between">
                  <span className="text-gray-400">Дата создания:</span>
                  <span className="text-gray-300">
                    {new Date(user.created_at).toLocaleDateString('ru-RU')}
                  </span>
                </div>
              )}
            </div>
          </div>

          {/* Confirmation Checkbox */}
          <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6">
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={isConfirmed}
                onChange={(e) => setIsConfirmed(e.target.checked)}
                disabled={isLoading}
                className="w-5 h-5 rounded border-gray-600 bg-gray-900 text-red-600 focus:ring-red-600 cursor-pointer"
              />
              <span className="text-sm text-gray-300">
                Я понимаю, что это действие <span className="font-semibold">необратимо</span> и подтверждаю удаление пользователя{' '}
                <span className="font-semibold text-red-400">{user.email}</span>
              </span>
            </label>
          </div>

          {/* Actions */}
          <div className="flex gap-4">
            <button
              type="submit"
              disabled={isLoading || !isConfirmed}
              className={`flex-1 px-6 py-3 rounded-lg font-medium transition-all flex items-center justify-center gap-2 ${
                isLoading || !isConfirmed
                  ? 'bg-gray-700 text-gray-400 cursor-not-allowed'
                  : 'bg-red-600 text-white hover:bg-red-700 active:bg-red-800'
              }`}
            >
              {isLoading ? (
                <>
                  <Loader2 size={18} className="animate-spin" />
                  Удаление...
                </>
              ) : (
                <>
                  <Trash2 size={18} />
                  Удалить пользователя
                </>
              )}
            </button>
            <button
              type="button"
              onClick={() => router.push('/users')}
              disabled={isLoading}
              className="px-6 py-3 rounded-lg font-medium bg-gray-700 text-gray-200 hover:bg-gray-600 transition-colors disabled:opacity-50"
            >
              Отмена
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
