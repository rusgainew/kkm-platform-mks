'use client';

import { useState } from 'react';
import { useInvoicesQuery, useDeleteInvoiceMutation, useInvoiceStatusMutation } from '@/lib/hooks/useInvoicesQuery';
import { useAdvancedTable } from '@/hooks/useAdvancedTable';
import { AdvancedFilter, FilterConfig } from '@/components/filters/AdvancedFilter';
import { BulkActions, BulkActionConfig } from '@/components/bulk/BulkActions';
import { Trash2, CheckCircle } from 'lucide-react';
import { showToast } from '@/lib/toast';

export function AdvancedInvoiceTable() {
  const { data: invoices = [], isLoading } = useInvoicesQuery();
  const deleteInvoice = useDeleteInvoiceMutation();
  const updateStatus = useInvoiceStatusMutation();

  const {
    filteredData,
    selectedIds,
    selectedCount,
    isAllSelected,
    filters,
    setFilters,
    toggleSelection,
    selectAll,
    clearSelection,
  } = useAdvancedTable({
    data: invoices,
    filterFn: (invoice, filters) => {
      if (filters.status && invoice.status !== filters.status) return false;
      if (filters.dateRange?.from && new Date(invoice.issueDate) < new Date(filters.dateRange.from)) return false;
      if (filters.dateRange?.to && new Date(invoice.issueDate) > new Date(filters.dateRange.to)) return false;
      if (filters.amount && invoice.totalAmount < parseInt(filters.amount)) return false;
      return true;
    },
    searchFields: ['invoiceNumber'],
  });

  const handleBulkDelete = async (ids: string[]) => {
    for (const id of ids) {
      await deleteInvoice.mutateAsync(id);
    }
    clearSelection();
    showToast.success(`Удалено ${ids.length} счетов`);
  };

  const handleBulkStatusUpdate = async (ids: string[]) => {
    for (const id of ids) {
      await updateStatus.mutateAsync({ id, status: 'paid' });
    }
    clearSelection();
    showToast.success(`Обновлено ${ids.length} счетов`);
  };

  const filterConfigs: FilterConfig[] = [
    {
      id: 'status',
      label: 'Статус',
      type: 'select',
      options: [
        { value: 'draft', label: 'Черновик' },
        { value: 'sent', label: 'Отправлен' },
        { value: 'paid', label: 'Оплачен' },
        { value: 'overdue', label: 'Просрочен' },
        { value: 'cancelled', label: 'Отменен' },
      ],
      value: filters.status || null,
      onChange: (value) => setFilters({ ...filters, status: value }),
    },
    {
      id: 'dateRange',
      label: 'Диапазон дат',
      type: 'dateRange',
      value: filters.dateRange || { from: null, to: null },
      onChange: (value) => setFilters({ ...filters, dateRange: value }),
    },
    {
      id: 'amount',
      label: 'Минимальная сумма ($)',
      type: 'text',
      placeholder: '100',
      value: filters.amount || '',
      onChange: (value) => setFilters({ ...filters, amount: value }),
    },
  ];

  const bulkActions: BulkActionConfig[] = [
    {
      id: 'mark-paid',
      label: 'Отметить как оплачено',
      icon: <CheckCircle className="w-4 h-4" />,
      action: handleBulkStatusUpdate,
      color: 'green',
      confirm: true,
    },
    {
      id: 'delete',
      label: 'Удалить',
      icon: <Trash2 className="w-4 h-4" />,
      action: handleBulkDelete,
      color: 'red',
      confirm: true,
    },
  ];

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Фильтры */}
      <AdvancedFilter
        filters={filterConfigs}
        onReset={() => setFilters({})}
      />

      {/* Bulk Actions */}
      <BulkActions
        selectedIds={selectedIds}
        actions={bulkActions}
        onClear={clearSelection}
      />

      {/* Таблица */}
      <div className="overflow-x-auto border rounded-lg">
        <table className="w-full">
          <thead className="bg-gray-100">
            <tr>
              <th className="px-4 py-3 text-left">
                <input
                  type="checkbox"
                  checked={isAllSelected}
                  onChange={isAllSelected ? clearSelection : selectAll}
                  className="rounded border-gray-300"
                />
              </th>
              <th className="px-4 py-3 text-left font-semibold">Номер</th>
              <th className="px-4 py-3 text-left font-semibold">Сумма</th>
              <th className="px-4 py-3 text-left font-semibold">Статус</th>
              <th className="px-4 py-3 text-left font-semibold">Дата</th>
            </tr>
          </thead>
          <tbody>
            {filteredData.map((invoice) => (
              <tr key={invoice.id} className="border-t hover:bg-gray-50">
                <td className="px-4 py-3">
                  <input
                    type="checkbox"
                    checked={selectedIds.includes(invoice.id)}
                    onChange={() => toggleSelection(invoice.id)}
                    className="rounded border-gray-300"
                  />
                </td>
                <td className="px-4 py-3 font-semibold">{invoice.invoiceNumber}</td>
                <td className="px-4 py-3">${invoice.totalAmount.toFixed(2)}</td>
                <td className="px-4 py-3">
                  <span className={`px-2 py-1 rounded-full text-xs font-semibold ${
                    invoice.status === 'paid' ? 'bg-green-100 text-green-800' :
                    invoice.status === 'sent' ? 'bg-blue-100 text-blue-800' :
                    invoice.status === 'overdue' ? 'bg-red-100 text-red-800' :
                    'bg-gray-100 text-gray-800'
                  }`}>
                    {invoice.status}
                  </span>
                </td>
                <td className="px-4 py-3">{new Date(invoice.issueDate).toLocaleDateString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {filteredData.length === 0 && (
        <div className="text-center py-8 text-gray-500">
          Нет счетов, соответствующих фильтрам
        </div>
      )}

      {/* Статистика */}
      <div className="text-sm text-gray-600">
        Показано: {filteredData.length} из {invoices.length} счетов
      </div>
    </div>
  );
}
