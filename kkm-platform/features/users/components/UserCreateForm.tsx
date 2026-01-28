'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Mail, Lock, User, ArrowLeft, AlertCircle, CheckCircle, Shield, Eye, EyeOff } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { createUser } from '@/lib/api/users';
import type { RegisterRequest } from '@/lib/api/users';

const ROLES = [
  { value: 'cashier', label: 'Кассир', description: 'Может работать с кассой' },
  { value: 'manager', label: 'Менеджер', description: 'Может управлять магазином' },
  { value: 'admin', label: 'Администратор', description: 'Полный доступ к системе' },
];

const PASSWORD_REQUIREMENTS = [
  { label: '8 символов', regex: /.{8,}/ },
  { label: 'Заглавная буква', regex: /[A-Z]/ },
  { label: 'Строчная буква', regex: /[a-z]/ },
  { label: 'Цифра', regex: /[0-9]/ },
  { label: 'Спецсимвол', regex: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/ },
];

// Parse error response from API
const parseApiError = (error: unknown): string => {
  if (error instanceof Error) {
    const message = error.message.toLowerCase();
    
    // Handle specific error messages
    if (message.includes('404') || message.includes('not found'))
      return 'Сервис недоступен. Пожалуйста, свяжитесь с администратором';
    if (message.includes('email')) return 'Email уже зарегистрирован в системе';
    if (message.includes('password')) return 'Пароль не соответствует требованиям безопасности';
    if (message.includes('permission') || message.includes('forbidden')) 
      return 'У вас нет прав для создания пользователей';
    if (message.includes('unauthorized') || message.includes('unauthenticated'))
      return 'Вы не авторизованы. Пожалуйста, перезагрузитесь';
    if (message.includes('network') || message.includes('fetch'))
      return 'Ошибка сети. Проверьте подключение к интернету';
    if (message.includes('timeout'))
      return 'Запрос истёк по времени. Попробуйте ещё раз';
    
    // Return original message if no match
    return error.message;
  }
  
  if (typeof error === 'string') return error;
  if (typeof error === 'object' && error !== null) {
    const err = error as any;
    if (err.message) return err.message;
    if (err.error) return err.error;
    if (err.detail) return err.detail;
  }
  
  return 'Неизвестная ошибка при создании пользователя';
};

interface FormErrors {
  email?: string;
  password?: string;
  confirmPassword?: string;
  firstName?: string;
  lastName?: string;
  role?: string;
  submit?: string;
}

export default function UserCreateForm() {
  const router = useRouter();
  
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    confirmPassword: '',
    firstName: '',
    lastName: '',
    role: 'cashier' as 'cashier' | 'manager' | 'admin',
  });

  const [errors, setErrors] = useState<FormErrors>({});
  const [successMessage, setSuccessMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  // Проверка требований пароля
  const validatePassword = (pwd: string) => {
    return PASSWORD_REQUIREMENTS.map(req => ({
      ...req,
      met: req.regex.test(pwd),
    }));
  };

  const passwordRequirements = validatePassword(formData.password);
  const isPasswordValid = passwordRequirements.every(req => req.met);

  // Валидация формы
  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};

    if (!formData.email.trim()) {
      newErrors.email = 'Email обязателен';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = 'Некорректный email';
    }

    if (!formData.firstName.trim()) {
      newErrors.firstName = 'Имя обязательно';
    }

    if (!formData.lastName.trim()) {
      newErrors.lastName = 'Фамилия обязательна';
    }

    if (!formData.role) {
      newErrors.role = 'Роль обязательна';
    }

    if (!formData.password) {
      newErrors.password = 'Пароль обязателен';
    } else if (!isPasswordValid) {
      newErrors.password = 'Пароль не соответствует требованиям';
    }

    if (!formData.confirmPassword) {
      newErrors.confirmPassword = 'Подтверждение пароля обязательно';
    } else if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Пароли не совпадают';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Отправка формы
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSuccessMessage('');

    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    try {
      const payload: RegisterRequest = {
        email: formData.email,
        password: formData.password,
        first_name: formData.firstName,
        last_name: formData.lastName,
      };

      await createUser(payload);

      setSuccessMessage('Пользователь успешно создан! Перенаправляем...');
      
      // Очистить форму
      setFormData({
        email: '',
        password: '',
        confirmPassword: '',
        firstName: '',
        lastName: '',
        role: 'cashier',
      });

      // Перенаправить на список пользователей через 2 секунды
      setTimeout(() => {
        router.push('/users');
      }, 2000);
    } catch (error) {
      const errorMessage = parseApiError(error);
      setErrors({ submit: errorMessage });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-2xl mx-auto">
      {/* Breadcrumb */}
      <button
        onClick={() => router.back()}
        className="flex items-center gap-2 text-blue-400 hover:text-blue-300 mb-6 transition-colors"
      >
        <ArrowLeft className="w-4 h-4" />
        Вернуться к списку пользователей
      </button>

      {/* Card */}
      <div className="bg-gray-900 border border-gray-800 rounded-lg shadow-lg">
        {/* Header */}
        <div className="p-6 border-b border-gray-800">
          <h1 className="text-2xl font-bold text-white">Создание нового пользователя</h1>
          <p className="text-gray-400 text-sm mt-1">
            Заполните форму для создания новой учётной записи
          </p>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          {/* Success Message */}
          {successMessage && (
            <div className="p-4 bg-emerald-900/20 border border-emerald-800 rounded-lg flex items-start gap-3">
              <CheckCircle className="w-5 h-5 text-emerald-400 shrink-0 mt-0.5" />
              <p className="text-emerald-300 text-sm">{successMessage}</p>
            </div>
          )}

          {/* Error Message */}
          {errors.submit && (
            <div className="p-4 bg-red-900/20 border border-red-800 rounded-lg flex items-start gap-3">
              <AlertCircle className="w-5 h-5 text-red-400 shrink-0 mt-0.5" />
              <div>
                <p className="text-red-300 font-semibold text-sm">Ошибка</p>
                <p className="text-red-400 text-sm mt-1">{errors.submit}</p>
              </div>
            </div>
          )}

          {/* Email */}
          <div>
            <label htmlFor="email" className="block text-sm font-semibold text-white mb-2">
              Email
            </label>
            <div className="relative">
              <Mail className="w-5 h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
              <input
                id="email"
                type="email"
                value={formData.email}
                onChange={(e) => {
                  setFormData({ ...formData, email: e.target.value });
                  if (errors.email) setErrors({ ...errors, email: undefined });
                }}
                disabled={isLoading || !!successMessage}
                placeholder="user@example.com"
                className={`w-full pl-10 pr-4 py-2 bg-gray-800 border rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 transition-colors ${
                  errors.email
                    ? 'border-red-700 focus:ring-red-600'
                    : 'border-gray-700 focus:border-blue-600 focus:ring-blue-600'
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              />
            </div>
            {errors.email && (
              <p className="text-red-400 text-sm mt-1 flex items-center gap-1">
                <AlertCircle className="w-4 h-4" />
                {errors.email}
              </p>
            )}
          </div>

          {/* First Name */}
          <div>
            <label htmlFor="firstName" className="block text-sm font-semibold text-white mb-2">
              Имя
            </label>
            <div className="relative">
              <User className="w-5 h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
              <input
                id="firstName"
                type="text"
                value={formData.firstName}
                onChange={(e) => {
                  setFormData({ ...formData, firstName: e.target.value });
                  if (errors.firstName) setErrors({ ...errors, firstName: undefined });
                }}
                disabled={isLoading || !!successMessage}
                placeholder="Иван"
                className={`w-full pl-10 pr-4 py-2 bg-gray-800 border rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 transition-colors ${
                  errors.firstName
                    ? 'border-red-700 focus:ring-red-600'
                    : 'border-gray-700 focus:border-blue-600 focus:ring-blue-600'
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              />
            </div>
            {errors.firstName && (
              <p className="text-red-400 text-sm mt-1 flex items-center gap-1">
                <AlertCircle className="w-4 h-4" />
                {errors.firstName}
              </p>
            )}
          </div>

          {/* Last Name */}
          <div>
            <label htmlFor="lastName" className="block text-sm font-semibold text-white mb-2">
              Фамилия
            </label>
            <div className="relative">
              <User className="w-5 h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
              <input
                id="lastName"
                type="text"
                value={formData.lastName}
                onChange={(e) => {
                  setFormData({ ...formData, lastName: e.target.value });
                  if (errors.lastName) setErrors({ ...errors, lastName: undefined });
                }}
                disabled={isLoading || !!successMessage}
                placeholder="Петров"
                className={`w-full pl-10 pr-4 py-2 bg-gray-800 border rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 transition-colors ${
                  errors.lastName
                    ? 'border-red-700 focus:ring-red-600'
                    : 'border-gray-700 focus:border-blue-600 focus:ring-blue-600'
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              />
            </div>
            {errors.lastName && (
              <p className="text-red-400 text-sm mt-1 flex items-center gap-1">
                <AlertCircle className="w-4 h-4" />
                {errors.lastName}
              </p>
            )}
          </div>

          {/* Role */}
          <div>
            <label htmlFor="role" className="block text-sm font-semibold text-white mb-2">
              Роль
            </label>
            <div className="relative">
              <Shield className="w-5 h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
              <select
                id="role"
                value={formData.role}
                onChange={(e) => {
                  setFormData({ ...formData, role: e.target.value as 'cashier' | 'manager' | 'admin' });
                  if (errors.role) setErrors({ ...errors, role: undefined });
                }}
                disabled={isLoading || !!successMessage}
                className={`w-full pl-10 pr-4 py-2 bg-gray-800 border rounded-lg text-white focus:outline-none focus:ring-2 transition-colors appearance-none ${
                  errors.role
                    ? 'border-red-700 focus:ring-red-600'
                    : 'border-gray-700 focus:border-blue-600 focus:ring-blue-600'
                } disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer`}
              >
                {ROLES.map(role => (
                  <option key={role.value} value={role.value} className="bg-gray-900">
                    {role.label}
                  </option>
                ))}
              </select>
            </div>
            {errors.role && (
              <p className="text-red-400 text-sm mt-1 flex items-center gap-1">
                <AlertCircle className="w-4 h-4" />
                {errors.role}
              </p>
            )}
            {/* Role Description */}
            <p className="text-gray-400 text-xs mt-2">
              {ROLES.find(r => r.value === formData.role)?.description}
            </p>
          </div>

          {/* Password */}
          <div>
            <label htmlFor="password" className="block text-sm font-semibold text-white mb-2">
              Пароль
            </label>
            <div className="relative">
              <Lock className="w-5 h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
              <input
                id="password"
                type={showPassword ? "text" : "password"}
                value={formData.password}
                onChange={(e) => {
                  setFormData({ ...formData, password: e.target.value });
                  if (errors.password) setErrors({ ...errors, password: undefined });
                }}
                disabled={isLoading || !!successMessage}
                placeholder="••••••••"
                className={`w-full pl-10 pr-10 py-2 bg-gray-800 border rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 transition-colors ${
                  errors.password && formData.password
                    ? 'border-red-700 focus:ring-red-600'
                    : 'border-gray-700 focus:border-blue-600 focus:ring-blue-600'
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                disabled={isLoading || !!successMessage}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
                aria-label={showPassword ? "Скрыть пароль" : "Показать пароль"}
              >
                {showPassword ? (
                  <EyeOff className="w-5 h-5" />
                ) : (
                  <Eye className="w-5 h-5" />
                )}
              </button>
            </div>
            {errors.password && (
              <p className="text-red-400 text-sm mt-1 flex items-center gap-1">
                <AlertCircle className="w-4 h-4" />
                {errors.password}
              </p>
            )}

            {/* Password Requirements */}
            {formData.password && (
              <div className="mt-3 p-3 bg-gray-800/50 border border-gray-700 rounded-lg">
                <p className="text-gray-300 text-xs font-semibold mb-2">Требования к паролю:</p>
                <div className="space-y-1">
                  {passwordRequirements.map((req, idx) => (
                    <div key={idx} className="flex items-center gap-2 text-xs">
                      <div
                        className={`w-4 h-4 rounded border flex items-center justify-center ${
                          req.met
                            ? 'bg-emerald-900/30 border-emerald-700'
                            : 'bg-gray-700/30 border-gray-600'
                        }`}
                      >
                        {req.met && <div className="w-2 h-2 bg-emerald-400 rounded-full" />}
                      </div>
                      <span className={req.met ? 'text-emerald-300' : 'text-gray-400'}>
                        {req.label}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Confirm Password */}
          <div>
            <label htmlFor="confirmPassword" className="block text-sm font-semibold text-white mb-2">
              Подтвердите пароль
            </label>
            <div className="relative">
              <Lock className="w-5 h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
              <input
                id="confirmPassword"
                type={showConfirmPassword ? "text" : "password"}
                value={formData.confirmPassword}
                onChange={(e) => {
                  setFormData({ ...formData, confirmPassword: e.target.value });
                  if (errors.confirmPassword) setErrors({ ...errors, confirmPassword: undefined });
                }}
                disabled={isLoading || !!successMessage}
                placeholder="••••••••"
                className={`w-full pl-10 pr-10 py-2 bg-gray-800 border rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 transition-colors ${
                  errors.confirmPassword && formData.confirmPassword
                    ? 'border-red-700 focus:ring-red-600'
                    : 'border-gray-700 focus:border-blue-600 focus:ring-blue-600'
                } disabled:opacity-50 disabled:cursor-not-allowed`}
              />
              <button
                type="button"
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                disabled={isLoading || !!successMessage}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
                aria-label={showConfirmPassword ? "Скрыть пароль" : "Показать пароль"}
              >
                {showConfirmPassword ? (
                  <EyeOff className="w-5 h-5" />
                ) : (
                  <Eye className="w-5 h-5" />
                )}
              </button>
            </div>
            {errors.confirmPassword && (
              <p className="text-red-400 text-sm mt-1 flex items-center gap-1">
                <AlertCircle className="w-4 h-4" />
                {errors.confirmPassword}
              </p>
            )}

            {/* Password Match Indicator */}
            {formData.password && formData.confirmPassword && (
              <div className="mt-2 flex items-center gap-2 text-xs">
                {formData.password === formData.confirmPassword ? (
                  <>
                    <div className="w-4 h-4 rounded-full bg-emerald-900/30 border border-emerald-700 flex items-center justify-center">
                      <div className="w-2 h-2 bg-emerald-400 rounded-full" />
                    </div>
                    <span className="text-emerald-300">Пароли совпадают</span>
                  </>
                ) : (
                  <>
                    <div className="w-4 h-4 rounded-full bg-red-900/30 border border-red-700 flex items-center justify-center">
                      <AlertCircle className="w-3 h-3 text-red-400" />
                    </div>
                    <span className="text-red-300">Пароли не совпадают</span>
                  </>
                )}
              </div>
            )}
          </div>

          {/* Action Buttons */}
          <div className="flex gap-3 pt-4">
            <Button
              type="submit"
              disabled={isLoading || !!successMessage}
              loading={isLoading}
              variant="primary"
              className="flex-1"
            >
              Создать пользователя
            </Button>
            <button
              type="button"
              onClick={() => router.back()}
              disabled={isLoading || !!successMessage}
              className="flex-1 px-4 py-2 bg-gray-800 hover:bg-gray-700 text-white font-semibold rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Отмена
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
