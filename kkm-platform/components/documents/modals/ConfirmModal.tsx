'use client';

import React, { useState } from 'react';
import { AlertTriangle } from 'lucide-react';
import { Modal, ModalButton, ErrorMessage } from '@/components/ui';

interface ConfirmModalProps {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  isDangerous?: boolean;
  onConfirm: () => Promise<void>;
  onCancel: () => void;
}

export function ConfirmModal({
  title,
  message,
  confirmText = 'Подтвердить',
  cancelText = 'Отмена',
  isDangerous = false,
  onConfirm,
  onCancel,
}: ConfirmModalProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleConfirm = async () => {
    setLoading(true);
    setError('');
    try {
      await onConfirm();
      onCancel();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при выполнении операции');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      isOpen={true}
      onClose={onCancel}
      title={
        <div className="flex items-center gap-3">
          {isDangerous && <AlertTriangle className="w-5 h-5 text-red-500" />}
          <span>{title}</span>
        </div>
      }
      size="md"
      footer={
        <>
          <ModalButton onClick={onCancel} variant="secondary" disabled={loading}>
            {cancelText}
          </ModalButton>
          <ModalButton 
            onClick={handleConfirm} 
            variant={isDangerous ? 'danger' : 'primary'} 
            loading={loading}
          >
            {confirmText}
          </ModalButton>
        </>
      }
    >
      <div className="space-y-4">
        <ErrorMessage message={error} />
        <p className="text-gray-300">{message}</p>
      </div>
    </Modal>
  );
}
