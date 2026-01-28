/**
 * Edit Company Modal - Using same style as EditUserModal
 */

'use client';

import React, { useState, useEffect, memo } from 'react';
import type { Company, CompanyStatus } from '@/types/entities';
import { useUpdateCompanyMutation } from '@/lib/hooks/useCompaniesApi';
import { Modal, ModalButton, ErrorMessage, Input, Textarea, Select } from '@/components/ui';

interface EditCompanyModalProps {
  isOpen: boolean;
  company: Company | null;
  onClose: () => void;
  onSuccess?: () => void;
}

const EditCompanyModal = memo(function EditCompanyModal({ isOpen, company, onClose, onSuccess }: EditCompanyModalProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState<'active' | 'inactive' | 'suspended'>('active');
  const [error, setError] = useState('');

  const updateMutation = useUpdateCompanyMutation();

  useEffect(() => {
    if (company) {
      setName(company.name);
      setDescription(company.description ?? '');
      setStatus(company.status as 'active' | 'inactive' | 'suspended');
      setError('');
    }
  }, [company, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!name.trim()) {
      setError('Название компании обязательно');
      return;
    }

    if (!company) return;

    try {
      const trimmedDescription = description.trim();
      const payload: {
        name: string;
        description?: string;
        status: CompanyStatus;
      } = {
        name: name.trim(),
        status,
      };
      
      if (trimmedDescription) {
        payload.description = trimmedDescription;
      }
      
      await updateMutation.mutateAsync({
        id: company.id || company.company_id,
        payload,
      });

      onSuccess?.();
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при обновлении компании');
    }
  };

  if (!isOpen || !company) return null;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Редактировать компанию"
      size="md"
      footer={
        <>
          <ModalButton onClick={onClose} variant="secondary" disabled={updateMutation.isPending}>
            Отмена
          </ModalButton>
          <ModalButton type="submit" form="edit-company-form" variant="primary" loading={updateMutation.isPending}>
            Сохранить
          </ModalButton>
        </>
      }
    >
      <form id="edit-company-form" onSubmit={handleSubmit} className="space-y-4">
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
            disabled={updateMutation.isPending}
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
            disabled={updateMutation.isPending}
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Статус
          </label>
          <Select
            value={status}
            onChange={(e) => setStatus(e.target.value as 'active' | 'inactive' | 'suspended')}
            disabled={updateMutation.isPending}
          >
            <option value="active">Активна</option>
            <option value="inactive">Неактивна</option>
            <option value="suspended">Заблокирована</option>
          </Select>
        </div>
      </form>
    </Modal>
  );
});

export default EditCompanyModal;
