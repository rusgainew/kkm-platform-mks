'use client';

import React from 'react';
import Link from 'next/link';
import { useAudit, AuditAction, AuditResource } from '@/lib/hooks/useAudit';

/**
 * Пример интеграции Audit логирования в компонент
 * Логирует все операции CRUD
 */
function InvoiceActionsWithAudit() {
  const { logAction, logChange } = useAudit({
    userId: 'user-123', // В реальном приложении получать из auth
    userName: 'Иван Петров',
    userEmail: 'ivan@example.com',
  });

  const handleCreateInvoice = async () => {
    await logAction(
      AuditAction.CREATE,
      AuditResource.INVOICE,
      'inv-new-123',
      { status: 'draft', amount: 10000 }
    );
    // создание счета...
  };

  const handleUpdateInvoice = async (invoiceId: string) => {
    const oldData = { status: 'draft', amount: 10000 };
    const newData = { status: 'issued', amount: 10000 };

    await logChange(
      AuditAction.UPDATE,
      AuditResource.INVOICE,
      invoiceId,
      oldData,
      newData,
      { reason: 'Счет отправлен клиенту' }
    );
    // обновление счета...
  };

  const handleDeleteInvoice = async (invoiceId: string) => {
    await logAction(
      AuditAction.DELETE,
      AuditResource.INVOICE,
      invoiceId,
      { deletedAt: new Date().toISOString() }
    );
    // удаление счета...
  };

  const handleExportInvoices = async () => {
    await logAction(
      AuditAction.EXPORT,
      AuditResource.INVOICE,
      'batch-export',
      { count: 42, format: 'PDF' }
    );
    // экспорт...
  };

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-bold">Управление счетами</h2>
      <p className="text-gray-600">
        Все операции логируются в журнал аудита (просмотрите в {' '}
        <Link href="/admin/audit-logs" className="text-blue-600 hover:underline">
          журнале аудита
        </Link>
        )
      </p>

      <div className="flex gap-2 flex-wrap">
        <button
          onClick={handleCreateInvoice}
          className="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600"
        >
          ➕ Создать счет
        </button>
        <button
          onClick={() => handleUpdateInvoice('inv-123')}
          className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
        >
          ✏️ Обновить счет
        </button>
        <button
          onClick={() => handleDeleteInvoice('inv-123')}
          className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600"
        >
          🗑️ Удалить счет
        </button>
        <button
          onClick={handleExportInvoices}
          className="px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600"
        >
          📤 Экспортировать
        </button>
      </div>
    </div>
  );
}

/**
 * Пример использования useAuditTimer для отслеживания длительности операций
 */
function AdvancedAuditExample() {
  const { log } = useAudit({
    userId: 'user-456',
    userName: 'Мария Сидорова',
    userEmail: 'maria@example.com',
  });

  const handleComplexOperation = async () => {
    const startTime = performance.now();

    try {
      // Имитируем сложную операцию
      await new Promise((resolve) => setTimeout(resolve, 2000));

      const duration = performance.now() - startTime;

      await log({
        action: AuditAction.IMPORT,
        resource: AuditResource.PRODUCT,
        resourceId: 'batch-import-789',
        status: 'success',
        userId: 'user-456',
        userName: 'Мария Сидорова',
        userEmail: 'maria@example.com',
        duration: Math.round(duration),
        metadata: {
          itemsProcessed: 1542,
          errorCount: 3,
          source: 'csv_upload',
        },
      });
    } catch (error) {
      const duration = performance.now() - startTime;

      await log({
        action: AuditAction.IMPORT,
        resource: AuditResource.PRODUCT,
        resourceId: 'batch-import-failed',
        status: 'failure',
        error: error instanceof Error ? error.message : 'Unknown error',
        userId: 'user-456',
        userName: 'Мария Сидорова',
        userEmail: 'maria@example.com',
        duration: Math.round(duration),
        metadata: {
          source: 'csv_upload',
        },
      });

      throw error;
    }
  };

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-bold">Импорт товаров</h2>
      <p className="text-gray-600">
        Операция логируется с отслеживанием времени выполнения
      </p>

      <button
        onClick={handleComplexOperation}
        className="px-4 py-2 bg-indigo-500 text-white rounded-lg hover:bg-indigo-600"
      >
        📥 Импортировать товары (2 сек)
      </button>
    </div>
  );
}

export default function AuditDemoPage() {
  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-gray-900 mb-2">🔐 Демо Audit Логирования</h1>
        <p className="text-gray-600 mb-8">
          Примеры использования audit логирования в приложении
        </p>

        <div className="bg-white rounded-lg shadow p-6 space-y-8">
          <InvoiceActionsWithAudit />
          <hr />
          <AdvancedAuditExample />
        </div>

        <div className="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-6">
          <h3 className="font-semibold text-blue-900 mb-2">ℹ️ Как использовать</h3>
          <ul className="text-sm text-blue-800 space-y-1">
            <li>✅ Нажимайте на кнопки выше для логирования действий</li>
            <li>✅ Все логи сохраняются в IndexedDB + localStorage</li>
            <li>✅ Просмотрите все логи в <Link href="/admin/audit-logs" className="text-blue-600 font-semibold hover:underline">журнале аудита</Link></li>
            <li>✅ Экспортируйте логи в JSON или CSV</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
