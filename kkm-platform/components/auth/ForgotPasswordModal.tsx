'use client';

import React, { useState } from 'react';
import { Mail, CheckCircle2 } from 'lucide-react';
import { useForgotPasswordMutation } from '@/lib/hooks/useAuthApi';
import { Modal, ModalButton, ErrorMessage, Input } from '@/components/ui';

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
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Восстановление пароля"
      size="md"
      footer={
        !showSuccess ? (
          <>
            <ModalButton onClick={onClose} variant="secondary" disabled={forgotPasswordMutation.isPending}>
              Отмена
            </ModalButton>
            <ModalButton 
              type="submit" 
              form="forgot-password-form" 
              variant="primary" 
              loading={forgotPasswordMutation.isPending}
              disabled={!email.trim()}
            >
              Отправить
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
            <form id="forgot-password-form" onSubmit={handleSubmit} className="space-y-4">
              <div>
                <p className="text-gray-400 text-sm mb-4">
                  Введите адрес электронной почты, связанный с вашим аккаунтом. Мы отправим вам ссылку для сброса пароля.
                </p>
              </div>

              <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-300 mb-2">
                  Email адрес
                </label>
                <div className="relative">
                  <Mail className="w-5 h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
                  <Input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="example@company.com"
                    className="pl-10"
                    required
                    disabled={forgotPasswordMutation.isPending}
                  />
                </div>
              </div>

              <ErrorMessage 
                message={
                  forgotPasswordMutation.isError
                    ? forgotPasswordMutation.error instanceof Error
                      ? forgotPasswordMutation.error.message
                      : 'Не удалось отправить письмо. Попробуйте позже.'
                    : null
                } 
              />
            </form>
          )}
    </Modal>
  );
}
