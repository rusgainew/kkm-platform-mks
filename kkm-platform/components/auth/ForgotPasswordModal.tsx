'use client';

import React, { useState } from 'react';
import { X, Mail, AlertCircle, Loader2, CheckCircle2 } from 'lucide-react';
import { useForgotPasswordMutation } from '@/lib/hooks/useAuthApi';

interface ForgotPasswordModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

export default function ForgotPasswordModal({
  isOpen,
  onClose,
  onSuccess,
}: ForgotPasswordModalProps) {
  const [email, setEmail] = useState('');
  const [showSuccess, setShowSuccess] = useState(false);
  const forgotPasswordMutation = useForgotPasswordMutation();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      await forgotPasswordMutation.mutateAsync({ email });
      setShowSuccess(true);
      
      // Reset form after 3 seconds and close
      setTimeout(() => {
        setEmail('');
        setShowSuccess(false);
        onSuccess?.();
        onClose();
      }, 3000);
    } catch (error) {
      console.error('Ошибка отправки письма для восстановления пароля:', error);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-gray-900 rounded-lg shadow-xl max-w-md w-full mx-4">
        {/* Header */}
        <div className="border-b border-gray-800 p-6 flex items-center justify-between">
          <h2 className="text-xl font-bold text-white">Восстановление пароля</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white transition-colors"
            disabled={forgotPasswordMutation.isPending}
          >
            <X className="w-6 h-6" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6">
          {showSuccess ? (
            <div className="text-center space-y-4">
              <div className="flex justify-center">
                <CheckCircle2 className="w-16 h-16 text-emerald-400" />
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white mb-2">Письмо отправлено!</h3>
                <p className="text-gray-400 text-sm">
                  На адрес <span className="font-medium text-white">{email}</span> отправлено письмо с инструкциями по восстановлению пароля.
                </p>
              </div>
              <p className="text-xs text-gray-500 pt-4">
                Проверьте папку спам, если письмо не пришло в течение 5 минут.
              </p>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <p className="text-gray-400 text-sm mb-4">
                  Введите адрес электронной почты, связанный с вашим аккаунтом. Мы отправим вам ссылку для сброса пароля.
                </p>
              </div>

              {/* Email Input */}
              <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-300 mb-2">
                  Email адрес
                </label>
                <div className="relative">
                  <Mail className="w-5 h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
                  <input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="example@company.com"
                    className="w-full pl-10 pr-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500 transition-colors"
                    required
                    disabled={forgotPasswordMutation.isPending}
                  />
                </div>
              </div>

              {/* Error Message */}
              {forgotPasswordMutation.isError && (
                <div className="p-3 bg-red-900/20 border border-red-800 rounded-lg flex gap-2">
                  <AlertCircle className="w-5 h-5 text-red-400 shrink-0 mt-0.5" />
                  <div>
                    <p className="text-red-300 text-sm font-medium">Ошибка</p>
                    <p className="text-red-400 text-xs">
                      {forgotPasswordMutation.error instanceof Error
                        ? forgotPasswordMutation.error.message
                        : 'Не удалось отправить письмо. Попробуйте позже.'}
                    </p>
                  </div>
                </div>
              )}

              {/* Buttons */}
              <div className="flex gap-2 pt-4">
                <button
                  type="button"
                  onClick={onClose}
                  className="flex-1 px-4 py-2 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                  disabled={forgotPasswordMutation.isPending}
                >
                  Отмена
                </button>
                <button
                  type="submit"
                  className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                  disabled={forgotPasswordMutation.isPending || !email.trim()}
                >
                  {forgotPasswordMutation.isPending ? (
                    <>
                      <Loader2 className="w-4 h-4 animate-spin" />
                      Отправка...
                    </>
                  ) : (
                    'Отправить'
                  )}
                </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}
