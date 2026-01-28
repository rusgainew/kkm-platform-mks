'use client';

import React, { useState } from 'react';
import toast from 'react-hot-toast';
import { rejectDocument } from '@/lib/api/documents';
import { Modal, ModalButton, ErrorMessage, Textarea } from '@/components/ui';

interface RejectDocumentModalProps {
  documentId: string;
  onClose: () => void;
  onSuccess: () => void;
}

export function RejectDocumentModal({ documentId, onClose, onSuccess }: RejectDocumentModalProps) {
  const [reason, setReason] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      await rejectDocument(documentId, reason);
      toast.success('Документ отклонен');
      onSuccess();
      onClose();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Ошибка при отклонении документа';
      setError(message);
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      isOpen={true}
      onClose={onClose}
      title="Отклонить документ"
      size="md"
      footer={
        <>
          <ModalButton onClick={onClose} variant="secondary" disabled={loading}>
            Отмена
          </ModalButton>
          <ModalButton 
            type="submit" 
            form="reject-document-form" 
            variant="danger" 
            loading={loading}
            disabled={!reason.trim()}
          >
            Отклонить
          </ModalButton>
        </>
      }
    >
      <form id="reject-document-form" onSubmit={handleSubmit} className="space-y-4">
        <ErrorMessage message={error} />

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Причина отклонения
          </label>
          <Textarea
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            rows={4}
            placeholder="Укажите причину отклонения документа (обязательно)"
            required
          />
        </div>
      </form>
    </Modal>
  );
}
