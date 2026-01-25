'use client';

import React from 'react';
import { X, Calendar, User } from 'lucide-react';
import { Document } from '@/lib/api/documents';
import { STATUS_LABELS, STATUS_COLORS } from '@/lib/documents/constants';
import { formatDate } from '@/lib/documents/formatting';

interface ViewDocumentModalProps {
  document: Document;
  onClose: () => void;
}

export function ViewDocumentModal({ document, onClose }: ViewDocumentModalProps) {
  const statusColor = STATUS_COLORS[document.status as keyof typeof STATUS_COLORS] || STATUS_COLORS.draft;
  const statusLabel = STATUS_LABELS[document.status as keyof typeof STATUS_LABELS] || document.status;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-gray-900 border border-gray-800 rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="sticky top-0 flex items-center justify-between p-6 border-b border-gray-800 bg-gray-900">
          <h2 className="text-xl font-bold text-white">{document.title}</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-800 rounded text-gray-400 hover:text-white transition"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Status */}
          <div className="flex items-center gap-4">
            <div>
              <p className="text-xs text-gray-500 uppercase mb-1">Статус</p>
              <span className={`inline-flex px-3 py-1 rounded-full text-sm font-semibold ${statusColor}`}>
                {statusLabel}
              </span>
            </div>
            <div className="flex-1"></div>
          </div>

          {/* Meta information */}
          <div className="grid grid-cols-2 gap-6">
            <div>
              <p className="text-xs text-gray-500 uppercase mb-2">Документ ID</p>
              <p className="text-white font-mono text-sm">{document.id}</p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase mb-2">Организация</p>
              <p className="text-white">{document.organization_id}</p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase mb-2 flex items-center gap-2">
                <User className="w-3 h-3" />
                Автор
              </p>
              <p className="text-white">{document.created_by || 'N/A'}</p>
            </div>
            <div>
              <p className="text-xs text-gray-500 uppercase mb-2 flex items-center gap-2">
                <Calendar className="w-3 h-3" />
                Создан
              </p>
              <p className="text-white">{formatDate(document.created_at)}</p>
            </div>
          </div>

          {/* Content */}
          {document.content && (
            <div className="border-t border-gray-800 pt-6">
              <p className="text-xs text-gray-500 uppercase mb-3">Содержание</p>
              <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-4 text-gray-300 text-sm whitespace-pre-wrap break-words">
                {document.content}
              </div>
            </div>
          )}

          {/* Timestamps */}
          {document.status_changed_at && (
            <div className="border-t border-gray-800 pt-6">
              <div className="grid grid-cols-2 gap-4 text-xs">
                <div>
                  <p className="text-gray-500 uppercase mb-1">Последнее обновление</p>
                  <p className="text-gray-300">{formatDate(document.updated_at)}</p>
                </div>
                <div>
                  <p className="text-gray-500 uppercase mb-1">Статус изменён</p>
                  <p className="text-gray-300">{formatDate(document.status_changed_at)}</p>
                </div>
              </div>
            </div>
          )}

          {/* Close button */}
          <div className="flex justify-end pt-4 border-t border-gray-800">
            <button
              onClick={onClose}
              className="px-4 py-2 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg transition"
            >
              Закрыть
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
