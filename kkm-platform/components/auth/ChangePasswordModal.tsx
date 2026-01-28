'use client';

import React, { useState } from 'react';
import { CheckCircle2 } from 'lucide-react';
import { useChangePasswordMutation } from '@/lib/hooks/useAuthApi';
import { useAuthStore } from '@/store/authStore';
import { Modal, ModalButton, ErrorMessage, PasswordInput } from '@/components/ui';

interface ChangePasswordModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

export default function ChangePasswordModal({
  isOpen,
  onClose,
  onSuccess,
}: ChangePasswordModalProps) {
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showSuccess, setShowSuccess] = useState(false);
  const changePasswordMutation = useChangePasswordMutation();
  const user = useAuthStore((state) => state.user);

  const passwordsMatch = newPassword === confirmPassword && newPassword.length >= 8;
  const formValid =
    oldPassword.length >= 8 &&
    passwordsMatch &&
    newPassword !== oldPassword;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formValid || !user) {
      return;
    }

    try {
      await changePasswordMutation.mutateAsync({
        user_id: user?.id || '',
        old_password: oldPassword,
        new_password: newPassword,
      });
      setShowSuccess(true);

      setTimeout(() => {
        setOldPassword('');
        setNewPassword('');
        setConfirmPassword('');
        setShowSuccess(false);
        onSuccess?.();
        onClose();
      }, 2000);
    } catch (error) {
      console.error('Failed to change password:', error);
    }
  };

  if (!isOpen) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Изменить пароль"
      size="md"
      footer={
        !showSuccess ? (
          <>
            <ModalButton onClick={onClose} variant="secondary" disabled={changePasswordMutation.isPending}>
              Отмена
            </ModalButton>
            <ModalButton 
              type="submit" 
              form="change-password-form" 
              variant="primary" 
              loading={changePasswordMutation.isPending}
              disabled={!formValid}
            >
              Изменить пароль
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
                  Ваш пароль успешно изменён. Для защиты вашего аккаунта используйте новый пароль при следующем входе.
                </p>
              </div>
            </div>
          ) : (
            <form id="change-password-form" onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label htmlFor="oldPassword" className="block text-sm font-medium text-gray-300 mb-2">
                  Текущий пароль
                </label>
                <PasswordInput
                  id="oldPassword"
                  value={oldPassword}
                  onChange={(e) => setOldPassword(e.target.value)}
                  placeholder="Введите текущий пароль"
                  required
                  disabled={changePasswordMutation.isPending}
                />
              </div>

              <div>
                <label htmlFor="newPassword" className="block text-sm font-medium text-gray-300 mb-2">
                  Новый пароль
                </label>
                <PasswordInput
                  id="newPassword"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  placeholder="Минимум 8 символов"
                  required
                  disabled={changePasswordMutation.isPending}
                />
              </div>

              <div>
                <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-300 mb-2">
                  Подтвердите пароль
                </label>
                <PasswordInput
                  id="confirmPassword"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  placeholder="Повторите новый пароль"
                  className={confirmPassword && !passwordsMatch ? 'border-red-700 focus:border-red-500' : ''}
                  required
                  disabled={changePasswordMutation.isPending}
                />
                {confirmPassword && !passwordsMatch && (
                  <p className="text-red-400 text-xs mt-1">
                    Пароли не совпадают или пароль короче 8 символов
                  </p>
                )}
              </div>

              <div className="p-3 bg-blue-900/20 border border-blue-800 rounded-lg">
                <p className="text-blue-300 text-xs">
                  <strong>Требования:</strong> Новый пароль должен отличаться от текущего и содержать минимум 8 символов.
                </p>
              </div>

              <ErrorMessage 
                message={
                  changePasswordMutation.isError
                    ? changePasswordMutation.error instanceof Error
                      ? changePasswordMutation.error.message
                      : 'Не удалось изменить пароль. Проверьте правильность текущего пароля.'
                    : null
                } 
              />
            </form>
          )}
    </Modal>
  );
}
