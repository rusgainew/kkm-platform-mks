'use client';

import { useState } from 'react';
import type { ApiUser } from '@/lib/api/users';
import { assignUserRole } from '@/lib/api/users';
import { Modal, ModalButton, ErrorMessage } from '@/components/ui';
import { Select } from '@/components/ui';

interface ChangeRoleModalProps {
  user: ApiUser | null;
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

const ROLE_OPTIONS = [
  { value: 'manager', label: 'Менеджер' },
  { value: 'admin', label: 'Администратор' },
  { value: 'cashier', label: 'Кассир' },
];

export default function ChangeRoleModal({ user, isOpen, onClose, onSuccess }: ChangeRoleModalProps) {
  const [selectedRoleValue, setSelectedRoleValue] = useState<string>('cashier');
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user || !user.id) return;

    setError(null);
    setIsLoading(true);
    try {
      await assignUserRole(user.id, {
        role: selectedRoleValue as 'user' | 'manager',
      });
      onSuccess?.();
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to change role');
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Изменить роль"
      size="md"
      footer={
        <>
          <ModalButton onClick={onClose} variant="secondary" disabled={isLoading}>
            Отмена
          </ModalButton>
          <ModalButton
            onClick={handleSubmit}
            variant="primary"
            loading={isLoading}
            disabled={!selectedRoleValue || selectedRoleValue === user?.role}
          >
            Сохранить
          </ModalButton>
        </>
      }
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <ErrorMessage message={error} />

        <div>
          <p className="text-sm text-gray-400 mb-2">Пользователь:</p>
          <p className="text-white font-medium">
            {user?.first_name} {user?.last_name}
          </p>
          <p className="text-sm text-gray-400">{user?.email}</p>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Новая роль
          </label>
          <Select
            value={isOpen && user ? (selectedRoleValue || user.role || 'user') : 'user'}
            onChange={(e) => setSelectedRoleValue(e.target.value)}
            disabled={isLoading}
          >
            {ROLE_OPTIONS.map((role) => (
              <option key={role.value} value={role.value}>
                {role.label}
              </option>
            ))}
          </Select>
        </div>

        {selectedRoleValue && user && selectedRoleValue !== user?.role && (
          <div className="p-3 bg-blue-900/20 border border-blue-800 rounded-lg text-sm text-blue-400">
            Роль будет изменена с <strong>{user?.role}</strong> на <strong>{selectedRoleValue}</strong>
          </div>
        )}
      </form>
    </Modal>
  );
}
