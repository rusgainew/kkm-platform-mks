'use client';

import React, { useState } from 'react';
import { Upload, X, File } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  createDocument,
  updateDocument,
  uploadDocumentFile,
} from '@/lib/api/documents';
import type {
  Document,
  DocumentStatus,
  CreateDocumentRequest,
  UpdateDocumentRequest,
} from '@/types/entities';

interface DocumentFormProps {
  initialData?: Document | null;
  organizationId?: string;
  createdBy?: string;
  onSuccess?: () => void;
  onCancel?: () => void;
}

interface FormData {
  title: string;
  content: string;
  status: DocumentStatus;
  assigned_to: string;
  organization_id: string;
  created_by: string;
  entries: Array<{ key: string; value: string }>;
}

interface UploadedFile {
  file: File;
  preview?: string;
}

// Validation
function isValidTitle(title: string): boolean {
  return title.trim().length >= 3 && title.length <= 200;
}

function isValidContent(content: string): boolean {
  return content.trim().length >= 10;
}

// File validation
const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB
const ALLOWED_FILE_TYPES = [
  'application/pdf',
  'image/jpeg',
  'image/png',
  'image/gif',
  'application/msword',
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
  'application/vnd.ms-excel',
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
];

function validateFile(file: File): string | null {
  if (file.size > MAX_FILE_SIZE) {
    return `Файл "${file.name}" слишком большой. Максимум 10MB`;
  }
  if (!ALLOWED_FILE_TYPES.includes(file.type)) {
    return `Файл "${file.name}" имеет неподдерживаемый формат`;
  }
  return null;
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
}

export default function DocumentForm({
  initialData,
  organizationId,
  createdBy,
  onSuccess,
  onCancel,
}: DocumentFormProps) {
  const queryClient = useQueryClient();
  const isEditMode = !!initialData?.id;

  const [formData, setFormData] = useState<FormData>({
    title: initialData?.title || '',
    content: initialData?.content || '',
    status: initialData?.status || 'draft',
    assigned_to: initialData?.assigned_to || '',
    organization_id: initialData?.organization_id || organizationId || '',
    created_by: initialData?.created_by || createdBy || '',
    entries: initialData?.entries || [],
  });

  const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([]);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [touched, setTouched] = useState<Record<string, boolean>>({});

  // Validation
  const validateForm = React.useCallback((): Record<string, string> => {
    const newErrors: Record<string, string> = {};

    if (!formData.title.trim()) {
      newErrors.title = 'Название документа обязательно';
    } else if (!isValidTitle(formData.title)) {
      newErrors.title = 'Название должно содержать от 3 до 200 символов';
    }

    if (!formData.content.trim()) {
      newErrors.content = 'Содержимое документа обязательно';
    } else if (!isValidContent(formData.content)) {
      newErrors.content = 'Содержимое должно содержать минимум 10 символов';
    }

    if (!formData.organization_id) {
      newErrors.organization_id = 'Необходимо указать организацию';
    }

    if (!formData.created_by && !isEditMode) {
      newErrors.created_by = 'Необходимо указать создателя';
    }

    return newErrors;
  }, [formData, isEditMode]);

  // Mutations
  const createMutation = useMutation({
    mutationFn: (data: CreateDocumentRequest) => createDocument(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['documents'] });
      onSuccess?.();
    },
  });

  const updateMutation = useMutation({
    mutationFn: (data: UpdateDocumentRequest) =>
      updateDocument(initialData!.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['documents'] });
      queryClient.invalidateQueries({
        queryKey: ['document', initialData!.id],
      });
      onSuccess?.();
    },
  });

  const uploadFileMutation = useMutation({
    mutationFn: ({ documentId, file }: { documentId: string; file: File }) =>
      uploadDocumentFile(documentId, file),
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Mark all fields as touched
    const allFields = Object.keys(formData).reduce(
      (acc, key) => ({ ...acc, [key]: true }),
      {}
    );
    setTouched(allFields);

    const validationErrors = validateForm();
    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }

    try {
      if (isEditMode) {
        // Update
        const updateData: UpdateDocumentRequest = {};
        if (formData.title !== initialData.title) {
          updateData.title = formData.title;
        }
        if (formData.content !== initialData.content) {
          updateData.content = formData.content;
        }
        if (formData.status !== initialData.status) {
          updateData.status = formData.status;
        }
        if (formData.assigned_to !== initialData.assigned_to) {
          updateData.assigned_to = formData.assigned_to;
        }
        if (formData.entries !== initialData.entries) {
          updateData.entries = formData.entries;
        }

        await updateMutation.mutateAsync(updateData);
      } else {
        // Create
        const createData: CreateDocumentRequest = {
          title: formData.title,
          content: formData.content,
          organization_id: formData.organization_id,
          assigned_to: formData.assigned_to || undefined,
          entries: formData.entries.length > 0 ? formData.entries : undefined,
        };

        const result = await createMutation.mutateAsync(createData);

        // Upload files if any
        if (uploadedFiles.length > 0 && result) {
          for (const { file } of uploadedFiles) {
            await uploadFileMutation.mutateAsync({
              documentId: result.id,
              file,
            });
          }
        }
      }
    } catch (error) {
      console.error('Ошибка отправки формы:', error);
      setErrors({
        submit:
          error instanceof Error ? error.message : 'Ошибка сохранения данных',
      });
    }
  };

  const handleBlur = (field: keyof FormData) => {
    setTouched((prev) => ({ ...prev, [field]: true }));
    const validationErrors = validateForm();
    setErrors((prev) => ({
      ...prev,
      [field]: validationErrors[field] || '',
    }));
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);
    const validFiles: UploadedFile[] = [];
    const fileErrors: string[] = [];

    files.forEach((file) => {
      const error = validateFile(file);
      if (error) {
        fileErrors.push(error);
      } else {
        validFiles.push({ file });
      }
    });

    if (fileErrors.length > 0) {
      setErrors((prev) => ({ ...prev, files: fileErrors.join('; ') }));
    } else {
      setErrors((prev) => ({ ...prev, files: '' }));
    }

    setUploadedFiles((prev) => [...prev, ...validFiles]);
  };

  const removeFile = (index: number) => {
    setUploadedFiles((prev) => prev.filter((_, i) => i !== index));
  };

  const addEntry = () => {
    setFormData((prev) => ({
      ...prev,
      entries: [...prev.entries, { key: '', value: '' }],
    }));
  };

  const removeEntry = (index: number) => {
    setFormData((prev) => ({
      ...prev,
      entries: prev.entries.filter((_, i) => i !== index),
    }));
  };

  const updateEntry = (
    index: number,
    field: 'key' | 'value',
    value: string
  ) => {
    setFormData((prev) => ({
      ...prev,
      entries: prev.entries.map((entry, i) =>
        i === index ? { ...entry, [field]: value } : entry
      ),
    }));
  };

  const isLoading = createMutation.isPending || updateMutation.isPending;
  const submitError = createMutation.error || updateMutation.error;

  return (
    <form
      onSubmit={handleSubmit}
      className="bg-gray-900 rounded-lg border border-gray-800 p-6 space-y-6"
    >
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-white">
          {isEditMode ? 'Редактировать документ' : 'Создать документ'}
        </h2>
      </div>

      {/* Error messages */}
      {(errors.submit || submitError) && (
        <div className="p-4 bg-red-900/20 border border-red-800 text-red-300 rounded-lg">
          {errors.submit ||
            (submitError instanceof Error
              ? submitError.message
              : 'Ошибка сохранения')}
        </div>
      )}

      {/* Basic Info */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-gray-300">
          Основная информация
        </h3>

        {/* Title */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Название документа <span className="text-red-400">*</span>
          </label>
          <input
            type="text"
            value={formData.title}
            onChange={(e) =>
              setFormData({ ...formData, title: e.target.value })
            }
            onBlur={() => handleBlur('title')}
            className={`w-full px-4 py-2 rounded-lg bg-gray-800 border ${
              touched.title && errors.title
                ? 'border-red-500'
                : 'border-gray-700'
            } text-white focus:border-green-500 outline-none`}
            placeholder="Договор поставки №123"
          />
          {touched.title && errors.title && (
            <p className="text-red-400 text-sm mt-1">{errors.title}</p>
          )}
          <p className="text-gray-400 text-xs mt-1">
            Длина: {formData.title.length}/200
          </p>
        </div>

        {/* Content */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Содержимое <span className="text-red-400">*</span>
          </label>
          <textarea
            value={formData.content}
            onChange={(e) =>
              setFormData({ ...formData, content: e.target.value })
            }
            onBlur={() => handleBlur('content')}
            rows={6}
            className={`w-full px-4 py-2 rounded-lg bg-gray-800 border ${
              touched.content && errors.content
                ? 'border-red-500'
                : 'border-gray-700'
            } text-white focus:border-green-500 outline-none resize-none`}
            placeholder="Описание документа, основные положения..."
          />
          {touched.content && errors.content && (
            <p className="text-red-400 text-sm mt-1">{errors.content}</p>
          )}
          <p className="text-gray-400 text-xs mt-1">
            Минимум 10 символов, текущая длина: {formData.content.length}
          </p>
        </div>

        {/* Status */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Статус <span className="text-red-400">*</span>
          </label>
          <select
            value={formData.status}
            onChange={(e) =>
              setFormData({
                ...formData,
                status: e.target.value as DocumentStatus,
              })
            }
            className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-green-500 outline-none"
          >
            <option value="draft">📝 Черновик</option>
            <option value="pending">⏳ На рассмотрении</option>
            <option value="approved">✅ Утвержден</option>
            <option value="rejected">❌ Отклонен</option>
            <option value="archived">📦 В архиве</option>
          </select>
        </div>

        {/* Assigned To */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Назначено
          </label>
          <input
            type="text"
            value={formData.assigned_to}
            onChange={(e) =>
              setFormData({ ...formData, assigned_to: e.target.value })
            }
            className="w-full px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-green-500 outline-none"
            placeholder="User ID или email"
          />
          <p className="text-gray-400 text-xs mt-1">
            Оставьте пустым, если не нужно назначать
          </p>
        </div>
      </div>

      {/* File Upload */}
      {!isEditMode && (
        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-gray-300">
            Файлы документов
          </h3>

          <div className="border-2 border-dashed border-gray-700 rounded-lg p-6 text-center hover:border-gray-600 transition-colors">
            <Upload className="mx-auto mb-3 text-gray-400" size={40} />
            <p className="text-gray-300 mb-2">
              Загрузите файлы документов (PDF, Word, Excel, изображения)
            </p>
            <p className="text-gray-400 text-sm mb-4">
              Максимальный размер файла: 10MB
            </p>
            <input
              type="file"
              multiple
              accept=".pdf,.doc,.docx,.xls,.xlsx,.jpg,.jpeg,.png,.gif"
              onChange={handleFileSelect}
              className="hidden"
              id="file-upload"
            />
            <label
              htmlFor="file-upload"
              className="inline-block px-6 py-2 bg-blue-600 text-white rounded-lg cursor-pointer hover:bg-blue-700 transition-colors"
            >
              Выбрать файлы
            </label>
          </div>

          {errors.files && (
            <p className="text-red-400 text-sm">{errors.files}</p>
          )}

          {uploadedFiles.length > 0 && (
            <div className="space-y-2">
              <p className="text-sm font-medium text-gray-300">
                Загружено файлов: {uploadedFiles.length}
              </p>
              {uploadedFiles.map((uploaded, index) => (
                <div
                  key={index}
                  className="flex items-center justify-between p-3 bg-gray-800 rounded-lg"
                >
                  <div className="flex items-center gap-3">
                    <File className="text-gray-400" size={20} />
                    <div>
                      <p className="text-white text-sm">
                        {uploaded.file.name}
                      </p>
                      <p className="text-gray-400 text-xs">
                        {formatFileSize(uploaded.file.size)}
                      </p>
                    </div>
                  </div>
                  <button
                    type="button"
                    onClick={() => removeFile(index)}
                    className="p-1 text-red-400 hover:text-red-300 transition-colors"
                  >
                    <X size={20} />
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Metadata Entries */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-gray-300">
            Дополнительные поля (метаданные)
          </h3>
          <button
            type="button"
            onClick={addEntry}
            className="px-3 py-1 bg-gray-700 text-white rounded text-sm hover:bg-gray-600 transition-colors"
          >
            + Добавить поле
          </button>
        </div>

        {formData.entries.length === 0 && (
          <p className="text-gray-400 text-sm">
            Нет дополнительных полей. Нажмите &quot;Добавить поле&quot; для создания.
          </p>
        )}

        {formData.entries.map((entry, index) => (
          <div key={index} className="flex gap-3">
            <input
              type="text"
              value={entry.key}
              onChange={(e) => updateEntry(index, 'key', e.target.value)}
              placeholder="Ключ (например, invoice_number)"
              className="flex-1 px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-green-500 outline-none"
            />
            <input
              type="text"
              value={entry.value}
              onChange={(e) => updateEntry(index, 'value', e.target.value)}
              placeholder="Значение"
              className="flex-1 px-4 py-2 rounded-lg bg-gray-800 border border-gray-700 text-white focus:border-green-500 outline-none"
            />
            <button
              type="button"
              onClick={() => removeEntry(index)}
              className="p-2 text-red-400 hover:text-red-300 transition-colors"
            >
              <X size={20} />
            </button>
          </div>
        ))}
      </div>

      {/* Buttons */}
      <div className="flex gap-3">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            disabled={isLoading}
            className="px-6 py-2 bg-gray-700 text-white rounded-lg font-medium hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Отмена
          </button>
        )}
        <Button
          type="submit"
          disabled={isLoading || Object.keys(errors).some((k) => k !== 'files' && errors[k])}
          loading={isLoading}
          variant="success"
          className="flex-1"
        >
          {isEditMode ? 'Обновить документ' : 'Создать документ'}
        </Button>
      </div>
    </form>
  );
}
