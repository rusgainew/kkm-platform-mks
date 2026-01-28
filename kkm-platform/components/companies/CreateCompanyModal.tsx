/**
 * Create Company Modal - Using same style as CreateUserModal
 */

'use client';

import React, { useState, memo } from 'react';
import { useCreateCompanyMutation } from '@/lib/hooks/useCompaniesApi';
import { Modal, ModalButton, ErrorMessage, Input, Textarea } from '@/components/ui';

interface CreateCompanyModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

const CreateCompanyModal = memo(function CreateCompanyModal({ isOpen, onClose, onSuccess }: CreateCompanyModalProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [error, setError] = useState('');

  const createMutation = useCreateCompanyMutation();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!name.trim()) {
      setError('Название компании обязательно');
      return;
    }

    try {
      await createMutation.mutateAsync({
        name: name.trim(),
        tin: '',
        address: '',
        phone: '',
        email: '',
        description: description.trim(),
      });

      setName('');
      setDescription('');
      onSuccess?.();
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при создании компании');
    }
  };

  if (!isOpen) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Создать компанию"
      size="md"
      footer={
        <>
          <ModalButton onClick={onClose} variant="secondary" disabled={createMutation.isPending}>
            Отмена
          </ModalButton>
          <ModalButton type="submit" form="create-company-form" variant="primary" loading={createMutation.isPending}>
            Создать
          </ModalButton>
        </>
      }
    >
      <form id="create-company-form" onSubmit={handleSubmit} className="space-y-4">
        <ErrorMessage message={error} />

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Название компании
          </label>
          <Input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="ООО Компания"
            disabled={createMutation.isPending}
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Описание
          </label>
          <Textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Описание компании..."
            rows={3}
            disabled={createMutation.isPending}
          />
        </div>
      </form>
    </Modal>
  );
});

export default CreateCompanyModal;
