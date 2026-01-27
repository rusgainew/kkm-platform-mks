'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowLeft, Loader2, AlertCircle, CheckCircle, Shield } from 'lucide-react';
import type { ApiUser } from '@/lib/api/users';
import { updateUser } from '@/lib/api/users';

const ROLES = [
  { value: 'cashier', label: 'Кассир', description: 'Может работать с кассой' },
  { value: 'manager', label: 'Менеджер', description: 'Может управлять магазином' },
  { value: 'admin', label: 'Администратор', description: 'Полный доступ к системе' },
];

interface FormErrors {
  firstName?: string;
  lastName?: string;
  role?: string;
  submit?: string;
}

interface UserEditFormProps {
  user: ApiUser;
}

export default function UserEditForm({ user }: UserEditFormProps) {
  const router = useRouter();
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [role, setRole] = useState<'cashier' | 'manager' | 'admin'>('cashier');
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<FormErrors>({});
  const [successMessage, setSuccessMessage] = useState('');
  const [hasChanges, setHasChanges] = useState(false);

  // Initialize form with user data
  useEffect(() => {
    if (user) {
      setFirstName(user.first_name || '');
      setLastName(user.last_name || '');
      setRole((user.role as 'cashier' | 'manager' | 'admin') || 'cashier');
    }
  }, [user]);

  // Track if form has been modified
  useEffect(() => {
    const hasChanged =
      firstName !== (user.first_name || '') ||
      lastName !== (user.last_name || '') ||
      role !== (user.role as 'cashier' | 'manager' | 'admin');
    setHasChanges(hasChanged);
  }, [firstName, lastName, role, user]);

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};

    if (firstName.trim().length === 0) {
      newErrors.firstName = 'Имя обязательно';
    } else if (firstName.trim().length < 2) {
      newErrors.firstName = 'Имя должно быть не менее 2 символов';
    }

    if (lastName.trim().length === 0) {
      newErrors.lastName = 'Фамилия обязательна';
    } else if (lastName.trim().length < 2) {
      newErrors.lastName = 'Фамилия должна быть не менее 2 символов';
    }

    if (!role) {
      newErrors.role = 'Выберите роль';
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
      const updateData: Record<string, any> = {};

      if (firstName !== (user.first_name || '')) {
        updateData.first_name = firstName.trim();
      }
      if (lastName !== (user.last_name || '')) {
        updateData.last_name = lastName.trim();
      }
      if (role !== (user.role as 'cashier' | 'manager' | 'admin')) {
        updateData.role = role;
      }

      if (Object.keys(updateData).length === 0) {
        setErrors({ submit: 'Нет изменений для сохранения' });
        setIsLoading(false);
        return;
      }

      if (!user.id) {
        setErrors({ submit: 'ID пользователя отсутствует' });
        setIsLoading(false);
        return;
      }

      await updateUser(user.id, updateData);

      setSuccessMessage('Пользователь успешно обновлён');
      setHasChanges(false);

      // Redirect back to users list after 2 seconds
      setTimeout(() => {
        router.push('/users');
      }, 2000);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : 'Ошибка при обновлении пользователя';
      setErrors({ submit: message });
    } finally {
      setIsLoading(false);
    }
  };

  const getRoleDescription = (selectedRole: string) => {
    return ROLES.find((r) => r.value === selectedRole)?.description || '';
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
            <h1 className="text-3xl font-bold text-white">Редактирование пользователя</h1>
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
          {errors.submit && (
            <div className="bg-red-900/20 border border-red-800 rounded-lg p-4 flex items-start gap-3">
              <AlertCircle size={20} className="text-red-400 flex-shrink-0 mt-0.5" />
              <p className="text-red-300">{errors.submit}</p>
            </div>
          )}

          {/* Form Fields */}
          <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6 space-y-6">
            {/* First Name */}
            <div>
              <label htmlFor="firstName" className="block text-sm font-medium text-gray-300 mb-2">
                Имя
              </label>
              <input
                id="firstName"
                type="text"
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
                placeholder="Иван"
                disabled={isLoading}
                className={`w-full px-4 py-3 rounded-lg bg-gray-900/50 border transition-colors outline-none ${
                  errors.firstName
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white placeholder-gray-500 disabled:opacity-50`}
              />
              {errors.firstName && (
                <p className="text-red-400 text-sm mt-1 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.firstName}
                </p>
              )}
            </div>

            {/* Last Name */}
            <div>
              <label htmlFor="lastName" className="block text-sm font-medium text-gray-300 mb-2">
                Фамилия
              </label>
              <input
                id="lastName"
                type="text"
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
                placeholder="Петров"
                disabled={isLoading}
                className={`w-full px-4 py-3 rounded-lg bg-gray-900/50 border transition-colors outline-none ${
                  errors.lastName
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white placeholder-gray-500 disabled:opacity-50`}
              />
              {errors.lastName && (
                <p className="text-red-400 text-sm mt-1 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.lastName}
                </p>
              )}
            </div>

            {/* Role Selection */}
            <div>
              <label htmlFor="role" className="block text-sm font-medium text-gray-300 mb-2">
                Роль
              </label>
              <select
                id="role"
                value={role}
                onChange={(e) => setRole(e.target.value as 'cashier' | 'manager' | 'admin')}
                disabled={isLoading}
                className={`w-full px-4 py-3 rounded-lg bg-gray-900/50 border transition-colors outline-none ${
                  errors.role
                    ? 'border-red-600 focus:border-red-500'
                    : 'border-gray-600 focus:border-blue-500'
                } text-white disabled:opacity-50`}
              >
                <option value="">Выберите роль</option>
                {ROLES.map((r) => (
                  <option key={r.value} value={r.value}>
                    {r.label}
                  </option>
                ))}
              </select>

              {/* Role Description */}
              {role && (
                <div className="mt-3 flex items-start gap-2 p-3 bg-blue-900/20 border border-blue-800/50 rounded">
                  <Shield size={16} className="text-blue-400 mt-0.5 flex-shrink-0" />
                  <p className="text-sm text-blue-300">{getRoleDescription(role)}</p>
                </div>
              )}

              {errors.role && (
                <p className="text-red-400 text-sm mt-1 flex items-center gap-1">
                  <AlertCircle size={14} />
                  {errors.role}
                </p>
              )}
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-4">
            <button
              type="submit"
              disabled={isLoading || !hasChanges}
              className={`flex-1 px-6 py-3 rounded-lg font-medium transition-all flex items-center justify-center gap-2 ${
                isLoading || !hasChanges
                  ? 'bg-gray-700 text-gray-400 cursor-not-allowed'
                  : 'bg-blue-600 text-white hover:bg-blue-700 active:bg-blue-800'
              }`}
            >
              {isLoading ? (
                <>
                  <Loader2 size={18} className="animate-spin" />
                  Сохранение...
                </>
              ) : (
                'Сохранить изменения'
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

        {/* User Info Card */}
        <div className="mt-8 bg-gray-800/50 border border-gray-700 rounded-lg p-6">
          <h3 className="text-sm font-medium text-gray-400 mb-4">Информация о пользователе</h3>
          <div className="space-y-3 text-sm">
            <div className="flex justify-between">
              <span className="text-gray-400">Email:</span>
              <span className="text-gray-300">{user.email}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">ID пользователя:</span>
              <span className="text-gray-300 font-mono">{user.id}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Текущая роль:</span>
              <span className="text-gray-300">{ROLES.find((r) => r.value === user.role)?.label}</span>
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
      </div>
    </div>
  );
}
