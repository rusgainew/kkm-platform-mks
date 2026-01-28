'use client';

import { useState } from 'react';
import type { ApiUser } from '@/lib/api/users';
import { updateUser } from '@/lib/api/users';
import { Modal, ModalButton, ErrorMessage, FormField, Input } from '@/components/ui';

interface EditUserModalProps {
  user: ApiUser | null;
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

export default function EditUserModal({ user, isOpen, onClose, onSuccess }: EditUserModalProps) {
  const [firstNameInput, setFirstNameInput] = useState('');
  const [lastNameInput, setLastNameInput] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user?.id) return;

    setError(null);
    setIsLoading(true);
    try {
      await updateUser(user.id, {
        first_name: firstNameInput || undefined,
        last_name: lastNameInput || undefined,
      });
      onSuccess?.();
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update user');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Редактировать пользователя"
      footer={
        <>
          <ModalButton variant="secondary" onClick={onClose} disabled={isLoading}>
            Отмена
          </ModalButton>
          <ModalButton variant="primary" onClick={handleSubmit} isLoading={isLoading}>
            Сохранить
          </ModalButton>
        </>
      }
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <ErrorMessage message={error} />

        <FormField label="Email">
          <Input type="email" value={user?.email || ''} disabled />
        </FormField>

        <FormField label="Имя">
          <Input
            type="text"
            value={isOpen && user ? (firstNameInput || user.first_name || '') : ''}
            onChange={(e) => setFirstNameInput(e.target.value)}
          />
        </FormField>

        <FormField label="Фамилия">
          <Input
            type="text"
            value={isOpen && user ? (lastNameInput || user.last_name || '') : ''}
            onChange={(e) => setLastNameInput(e.target.value)}
          />
        </FormField>
      </form>
    </Modal>
  );
}
