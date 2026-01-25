'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowLeft, Loader2, AlertCircle, CheckCircle, Shield } from 'lucide-react';
import type { ApiUser } from '@/lib/api/users';
import { assignUserRole } from '@/lib/api/users';

const ROLES = [
  { value: 'cashier', label: 'Кассир', description: 'Может работать с кассой' },
  { value: 'manager', label: 'Менеджер', description: 'Может управлять магазином' },
  { value: 'admin', label: 'Администратор', description: 'Полный доступ к системе' },
];

interface FormErrors {
  role?: string;
  submit?: string;
}

interface UserRoleFormProps {
  user: ApiUser;
}

export default function UserRoleForm({ user }: UserRoleFormProps) {
  const router = useRouter();
  const [newRole, setNewRole] = useState<'cashier' | 'manager' | 'admin'>(
    (user.role as 'cashier' | 'manager' | 'admin') || 'cashier'
  );
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<FormErrors>({});
  const [successMessage, setSuccessMessage] = useState('');

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};

    if (!newRole) {
      newErrors.role = 'Выберите новую роль';
    }

    if (newRole === user.role) {
      newErrors.submit = 'Новая роль совпадает с текущей';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    setErrors({});
    setSuccessMessage('');

    try {
      await assignUserRole(user.id, {
        role: newRole as 'user' | 'manager',
      });

      setSuccessMessage('Роль пользователя успешно изменена');

      // Redirect back to users list after 2 seconds
      setTimeout(() => {
        router.push('/users');
      }, 2000);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : 'Ошибка при изменении роли пользователя';
      setErrors({ submit: message });
    } finally {
      setIsLoading(false);
    }
  };

  const getRoleDescription = (selectedRole: string) => {
    return ROLES.find((r) => r.value === selectedRole)?.description || '';
  };

  const getCurrentRoleLabel = () => {
    return ROLES.find((r) => r.value === user.role)?.label || user.role;
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
            <h1 className="text-3xl font-bold text-white">Изменение роли пользователя</h1>
            <p className="text-gray-400 mt-2">
              <span className="text-gray-300">{user.email}</span>
              {' | Текущая роль: '}
              <span className="text-blue-400 font-semibold">{getCurrentRoleLabel()}</span>
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
          {errors.submit && (
            <div className="bg-red-900/20 border border-red-800 rounded-lg p-4 flex items-start gap-3">
              <AlertCircle size={20} className="text-red-400 flex-shrink-0 mt-0.5" />
              <p className="text-red-300">{errors.submit}</p>
            </div>
          )}

          {/* Current Role Info */}
          <div className="bg-blue-900/20 border border-blue-800 rounded-lg p-4">
            <p className="text-sm text-blue-300">
              <span className="font-semibold">Текущая роль:</span> {getCurrentRoleLabel()}
            </p>
          </div>

          {/* Role Selection */}
          <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6 space-y-6">
            <div>
              <label htmlFor="role" className="block text-sm font-medium text-gray-300 mb-2">
                Новая роль
              </label>
              <select
                id="role"
                value={newRole}
                onChange={(e) => setNewRole(e.target.value as 'cashier' | 'manager' | 'admin')}
                disabled={isLoading}
                className={`w-full px-4 py-3 rounded-lg bg-gray-900/50 border transition-colors outline-none ${
                  errors.role
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white disabled:opacity-50`}
              >
                {ROLES.map((r) => (
                  <option key={r.value} value={r.value}>
                    {r.label}
                  </option>
                ))}
              </select>

              {/* Role Description */}
              {newRole && (
                <div className="mt-3 flex items-start gap-2 p-3 bg-purple-900/20 border border-purple-800/50 rounded">
                  <Shield size={16} className="text-purple-400 mt-0.5 flex-shrink-0" />
                  <div>
                    <p className="text-sm font-medium text-purple-300">
                      {ROLES.find((r) => r.value === newRole)?.label}
                    </p>
                    <p className="text-sm text-purple-300/70 mt-1">{getRoleDescription(newRole)}</p>
                  </div>
                </div>
              )}

              {errors.role && (
                <p className="text-red-400 text-sm mt-1 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.role}
                </p>
              )}
            </div>

            {/* Role Comparison */}
            <div className="border-t border-gray-700 pt-6">
              <h3 className="text-sm font-medium text-gray-300 mb-4">Все доступные роли</h3>
              <div className="space-y-3">
                {ROLES.map((role) => (
                  <div
                    key={role.value}
                    className={`p-3 rounded-lg border transition-all ${
                      newRole === role.value
                        ? 'bg-purple-900/30 border-purple-600'
                        : user.role === role.value
                        ? 'bg-blue-900/30 border-blue-600'
                        : 'bg-gray-900/30 border-gray-700'
                    }`}
                  >
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="font-medium text-gray-300">{role.label}</p>
                        <p className="text-xs text-gray-400 mt-1">{role.description}</p>
                      </div>
                      {user.role === role.value && (
                        <span className="text-xs bg-blue-600/50 text-blue-300 px-2 py-1 rounded">
                          Текущая
                        </span>
                      )}
                      {newRole === role.value && user.role !== role.value && (
                        <span className="text-xs bg-purple-600/50 text-purple-300 px-2 py-1 rounded">
                          Новая
                        </span>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-4">
            <button
              type="submit"
              disabled={isLoading || newRole === user.role}
              className={`flex-1 px-6 py-3 rounded-lg font-medium transition-all flex items-center justify-center gap-2 ${
                isLoading || newRole === user.role
                  ? 'bg-gray-700 text-gray-400 cursor-not-allowed'
                  : 'bg-purple-600 text-white hover:bg-purple-700 active:bg-purple-800'
              }`}
            >
              {isLoading ? (
                <>
                  <Loader2 size={18} className="animate-spin" />
                  Изменение...
                </>
              ) : (
                'Изменить роль'
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
