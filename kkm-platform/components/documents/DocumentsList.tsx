'use client';

import React, { useState, useMemo } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { FileText, Search, Plus } from 'lucide-react';
import { useDocuments } from '@/hooks/useDocuments';
import { Document, sendDocument, approveDocument, archiveDocument } from '@/lib/api/documents';
import CreateDocumentForm from './CreateDocumentForm';
import { DocumentRow } from './DocumentRow';
import Pagination from '@/components/ui/Pagination';
import { 
  EditDocumentModal, 
  ConfirmModal, 
  RejectDocumentModal, 
  ViewDocumentModal 
} from './modals';

const ITEMS_PER_PAGE = 10;

export default function DocumentsList() {
  const queryClient = useQueryClient();
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [currentPage, setCurrentPage] = useState(1);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  // Modal states
  const [viewingDocument, setViewingDocument] = useState<Document | null>(null);
  const [editingDocument, setEditingDocument] = useState<Document | null>(null);
  const [rejectingDocument, setRejectingDocument] = useState<Document | null>(null);
  const [operationModal, setOperationModal] = useState<{
    type: 'send' | 'approve' | 'archive' | 'delete' | null;
    document: Document | null;
  }>({ type: null, document: null });

  // Use custom hook for document loading
  const { documents, pagination, isLoading, error } = useDocuments({
    page: currentPage,
    pageSize: ITEMS_PER_PAGE,
    status: statusFilter || undefined,
    refreshTrigger,
  });

  // Filter documents based on search
  const filteredDocuments = useMemo(() => {
    return documents.filter((doc) => {
      const matchesSearch =
        doc.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        doc.id.toLowerCase().includes(searchQuery.toLowerCase());
      return matchesSearch;
    });
  }, [documents, searchQuery]);

  // Handlers for operations
  const handleRefresh = () => {
    setRefreshTrigger(prev => prev + 1);
  };

  const handleSendDocument = async () => {
    if (!operationModal.document) return;
    try {
      await sendDocument(operationModal.document.id);
      toast.success(`Документ "${operationModal.document.title}" отправлен на одобрение`);
      handleRefresh();
      setOperationModal({ type: null, document: null });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Ошибка при отправке документа';
      toast.error(message);
      throw error;
    }
  };

  const handleApproveDocument = async () => {
    if (!operationModal.document) return;
    try {
      await approveDocument(operationModal.document.id);
      toast.success(`Документ "${operationModal.document.title}" одобрен`);
      handleRefresh();
      setOperationModal({ type: null, document: null });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Ошибка при одобрении документа';
      toast.error(message);
      throw error;
    }
  };

  const handleArchiveDocument = async () => {
    if (!operationModal.document) return;
    try {
      await archiveDocument(operationModal.document.id);
      toast.success(`Документ "${operationModal.document.title}" архивирован`);
      setCurrentPage(1);  // Reset to first page
      handleRefresh();
      setOperationModal({ type: null, document: null });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Ошибка при архивировании документа';
      toast.error(message);
      throw error;
    }
  };

  const handleDeleteDocument = async () => {
    if (!operationModal.document) return;
    try {
      // Delete is handled by archive API
      await archiveDocument(operationModal.document.id);
      toast.success(`Документ "${operationModal.document.title}" удалён`);
      setCurrentPage(1);  // Reset to first page
      handleRefresh();
      setOperationModal({ type: null, document: null });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Ошибка при удалении документа';
      toast.error(message);
      throw error;
    }
  };

  if (isLoading) {
    return (
      <div className="p-8">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-700 rounded w-1/4"></div>
          <div className="h-12 bg-gray-700 rounded"></div>
          <div className="space-y-3">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-16 bg-gray-700 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-8">
        <div className="bg-red-500/10 border border-red-500 rounded-lg p-4 text-red-400">
          <p className="font-semibold">Ошибка</p>
          <p className="text-sm mt-2">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-8 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white flex items-center gap-3">
            <FileText className="w-8 h-8 text-blue-400" />
            Документы
          </h1>
          <p className="text-gray-400 text-sm mt-1">
            Всего документов: {pagination.totalCount}
          </p>
        </div>
        <button 
          onClick={() => setShowCreateForm(true)}
          className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition">
          <Plus className="w-4 h-4" />
          Новый документ
        </button>
      </div>

      {/* Search and Filters */}
      <div className="space-y-4 bg-gray-900/50 border border-gray-800 rounded-lg p-4">
        <div className="flex gap-4">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-3 w-4 h-4 text-gray-500" />
            <input
              type="text"
              placeholder="Поиск по названию или номеру..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500"
            />
          </div>
        </div>

        <div className="flex gap-4">
          <select
            value={statusFilter}
            onChange={(e) => {
              setStatusFilter(e.target.value);
              setCurrentPage(1);
            }}
            className="px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-gray-300 text-sm focus:outline-none focus:border-blue-500"
          >
            <option value="">Все статусы</option>
            <option value="draft">Черновик</option>
            <option value="sent">Отправлено</option>
            <option value="approved">Одобрено</option>
            <option value="rejected">Отклонено</option>
            <option value="archived">В архиве</option>
          </select>
        </div>
      </div>

      {/* Documents Table */}
      {filteredDocuments.length > 0 ? (
        <>
          <div className="bg-gray-900/50 border border-gray-800 rounded-lg overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-800 bg-gray-800/50">
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                      Документ
                    </th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                      Организация
                    </th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                      Статус
                    </th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                      Автор
                    </th>
                    <th className="px-6 py-4 text-left text-sm font-semibold text-gray-300">
                      Дата
                    </th>
                    <th className="px-6 py-4 text-center text-sm font-semibold text-gray-300">
                      Действия
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-800">
                  {filteredDocuments.map((doc) => (
                    <DocumentRow
                      key={doc.id}
                      document={doc}
                      onView={(doc) => setViewingDocument(doc)}
                      onDownload={(doc) => console.log('Загрузка документа:', doc)}
                      onEdit={(doc) => setEditingDocument(doc)}
                      onDelete={(doc) => setOperationModal({ type: 'delete', document: doc })}
                      onArchive={(doc) => setOperationModal({ type: 'archive', document: doc })}
                      onSend={(doc) => setOperationModal({ type: 'send', document: doc })}
                      onApprove={(doc) => setOperationModal({ type: 'approve', document: doc })}
                      onReject={(doc) => setRejectingDocument(doc)}
                    />
                  ))}
                </tbody>
              </table>
            </div>
          </div>
          
          {/* Pagination */}
          <Pagination
            currentPage={pagination.currentPage}
            totalPages={pagination.totalPages}
            onPageChange={setCurrentPage}
            isLoading={isLoading}
          />
        </>
      ) : (
        <div className="bg-gray-900/50 border border-gray-800 rounded-lg p-12 text-center">
          <FileText className="w-12 h-12 text-gray-600 mx-auto mb-4" />
          <p className="text-gray-400 text-lg">Документы не найдены</p>
          <p className="text-gray-500 text-sm mt-2">
            {searchQuery ? 'Попробуйте изменить критерии поиска' : 'Начните с создания нового документа'}
          </p>
        </div>
      )}

      {/* Create Document Form Modal */}
      {showCreateForm && (
        <CreateDocumentForm
          onClose={() => setShowCreateForm(false)}
          onSuccess={() => {
            // Trigger documents refresh by incrementing refreshTrigger
            setRefreshTrigger(prev => prev + 1);
            setCurrentPage(1);
            setStatusFilter('');
          }}
        />
      )}

      {/* View Document Modal */}
      {viewingDocument && (
        <ViewDocumentModal
          document={viewingDocument}
          onClose={() => setViewingDocument(null)}
        />
      )}

      {/* Edit Document Modal */}
      {editingDocument && (
        <EditDocumentModal
          document={editingDocument}
          onClose={() => setEditingDocument(null)}
          onSuccess={handleRefresh}
        />
      )}

      {/* Reject Document Modal */}
      {rejectingDocument && (
        <RejectDocumentModal
          documentId={rejectingDocument.id}
          onClose={() => setRejectingDocument(null)}
          onSuccess={handleRefresh}
        />
      )}

      {/* Operation Confirmation Modals */}
      {operationModal.type === 'send' && operationModal.document && (
        <ConfirmModal
          title="Отправить документ"
          message={`Вы уверены, что хотите отправить документ "${operationModal.document.title}" на одобрение?`}
          confirmText="Отправить"
          onConfirm={handleSendDocument}
          onCancel={() => setOperationModal({ type: null, document: null })}
        />
      )}

      {operationModal.type === 'approve' && operationModal.document && (
        <ConfirmModal
          title="Одобрить документ"
          message={`Вы уверены, что хотите одобрить документ "${operationModal.document.title}"?`}
          confirmText="Одобрить"
          onConfirm={handleApproveDocument}
          onCancel={() => setOperationModal({ type: null, document: null })}
        />
      )}

      {operationModal.type === 'archive' && operationModal.document && (
        <ConfirmModal
          title="Архивировать документ"
          message={`Вы уверены, что хотите архивировать документ "${operationModal.document.title}"?`}
          confirmText="В архив"
          onConfirm={handleArchiveDocument}
          onCancel={() => setOperationModal({ type: null, document: null })}
        />
      )}

      {operationModal.type === 'delete' && operationModal.document && (
        <ConfirmModal
          title="Удалить документ"
          message={`Вы уверены, что хотите удалить документ "${operationModal.document.title}"? Это действие нельзя отменить.`}
          confirmText="Удалить"
          isDangerous={true}
          onConfirm={handleDeleteDocument}
          onCancel={() => setOperationModal({ type: null, document: null })}
        />
      )}
    </div>
  );
}
