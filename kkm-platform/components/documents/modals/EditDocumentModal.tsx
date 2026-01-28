'use client';

import React, { useState } from 'react';
import toast from 'react-hot-toast';
import { Document, updateDocument } from '@/lib/api/documents';
import { Modal, ModalButton, ErrorMessage, Input, Textarea } from '@/components/ui';

interface EditDocumentModalProps {
  document: Document;
  onClose: () => void;
  onSuccess: () => void;
}

export function EditDocumentModal({ document, onClose, onSuccess }: EditDocumentModalProps) {
  const [title, setTitle] = useState(document.title);
  const [content, setContent] = useState(document.content || '');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      await updateDocument(document.id, {
        title,
        content,
      });

      toast.success('Документ успешно обновлен');
      onSuccess();
      onClose();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Ошибка при сохранении документа';
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
      title="Редактировать документ"
      size="lg"
      footer={
        <>
          <ModalButton onClick={onClose} variant="secondary" disabled={loading}>
            Отмена
          </ModalButton>
          <ModalButton type="submit" form="edit-document-form" variant="primary" loading={loading}>
            Сохранить
          </ModalButton>
        </>
      }
    >
      <form id="edit-document-form" onSubmit={handleSubmit} className="space-y-4">
        <ErrorMessage message={error} />

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Название документа
          </label>
          <Input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Содержание
          </label>
          <Textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            rows={6}
          />
        </div>
      </form>
    </Modal>
  );
}
