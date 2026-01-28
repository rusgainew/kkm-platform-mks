'use client';

import React, { useState } from 'react';
import { CheckCircle2 } from 'lucide-react';
import { useResetPasswordMutation } from '@/lib/hooks/useAuthApi';
import { Modal, ModalButton, ErrorMessage, PasswordInput, Input } from '@/components/ui';

interface ResetPasswordModalProps {
  isOpen: boolean;
  onClose: () => void;
  resetToken?: string;
  onSuccess?: () => void;
}

export default function ResetPasswordModal({
  isOpen,
  onClose,
  resetToken = '',
  onSuccess,
}: ResetPasswordModalProps) {
  const [token, setToken] = useState(resetToken);
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showSuccess, setShowSuccess] = useState(false);
  const resetPasswordMutation = useResetPasswordMutation();

  const passwordsMatch = newPassword === confirmPassword && newPassword.length >= 8;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!passwordsMatch) {
      return;
    }

    try {
      await resetPasswordMutation.mutateAsync({
        token,
        new_password: newPassword,
      });
      setShowSuccess(true);

      setTimeout(() => {
        setToken('');
        setNewPassword('');
        setConfirmPassword('');
        setShowSuccess(false);
        onSuccess?.();
        onClose();
      }, 3000);
    } catch (error) {
      console.error('Failed to reset password:', error);
    }
  };

  if (!isOpen) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Сброс пароля"
      size="md"
      footer={
        !showSuccess ? (
          <>
            <ModalButton onClick={onClose} variant="secondary" disabled={resetPasswordMutation.isPending}>
              Отмена
            </ModalButton>
            <ModalButton 
              type="submit" 
              form="reset-password-form" 
              variant="primary" 
              loading={resetPasswordMutation.isPending}
              disabled={!passwordsMatch || !token.trim()}
            >
              Сбросить пароль
            </ModalButton>
          </>
        ) : null
      }
    >
          {showSuccess ? (
            <div className="text-center space-y-4">
              <div className="flex justify-center">
                <CheckCircle2 className="w-16 h-16 text-emerald-400" />
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white mb-2">Пароль изменён!</h3>
                <p className="text-gray-400 text-sm">
                  Ваш пароль успешно сброшен. Используйте новый пароль для входа.
                </p>
              </div>
            </div>
          ) : (
            <form id="reset-password-form" onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label htmlFor="token" className="block text-sm font-medium text-gray-300 mb-2">
                  Код восстановления
                </label>
                <Input
                  id="token"
                  type="text"
                  value={token}
                  onChange={(e) => setToken(e.target.value)}
                  placeholder="Вставьте код из письма"
                  className="text-sm font-mono"
                  required
                  disabled={resetPasswordMutation.isPending}
                />
              </div>

              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-300 mb-2">
                  Новый пароль
                </label>
                <PasswordInput
                  id="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  placeholder="Минимум 8 символов"
                  required
                  disabled={resetPasswordMutation.isPending}
                />
              </div>

              <div>
                <label htmlFor="confirm" className="block text-sm font-medium text-gray-300 mb-2">
                  Подтвердите пароль
                </label>
                <PasswordInput
                  id="confirm"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  placeholder="Повторите новый пароль"
                  className={confirmPassword && !passwordsMatch ? 'border-red-700 focus:border-red-500' : ''}
                  required
                  disabled={resetPasswordMutation.isPending}
                />
                {confirmPassword && !passwordsMatch && (
                  <p className="text-red-400 text-xs mt-1">
                    Пароли не совпадают или пароль короче 8 символов
                  </p>
                )}
              </div>

              <ErrorMessage 
                message={
                  resetPasswordMutation.isError
                    ? resetPasswordMutation.error instanceof Error
                      ? resetPasswordMutation.error.message
                      : 'Не удалось сбросить пароль. Проверьте код восстановления.'
                    : null
                } 
              />
            </form>
          )}
    </Modal>
  );
}
