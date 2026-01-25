'use client';

import { useInvoicesQuery, useDeleteInvoiceMutation, useInvoiceStatusMutation, type Invoice } from '@/lib/hooks/useInvoicesQuery';

export function InvoiceListWithReactQuery() {
  const { data: invoices, isLoading, error, refetch } = useInvoicesQuery();
  const deleteInvoice = useDeleteInvoiceMutation();
  const updateStatus = useInvoiceStatusMutation();

  if (isLoading) return <div className="p-4">Загрузка счетов...</div>;
  if (error) return <div className="p-4 text-red-600">Ошибка: {error.message}</div>;

  const handleStatusChange = (invoiceId: string, newStatus: Invoice['status']) => {
    updateStatus.mutate({ id: invoiceId, status: newStatus });
  };

  const handleDelete = (invoiceId: string) => {
    if (confirm('Вы уверены?')) {
      deleteInvoice.mutate(invoiceId);
    }
  };

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Счета (React Query)</h1>
        <button
          onClick={() => refetch()}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          Обновить
        </button>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full border-collapse">
          <thead>
            <tr className="bg-gray-100">
              <th className="border p-2 text-left">Номер</th>
              <th className="border p-2 text-left">Сумма</th>
              <th className="border p-2 text-left">Статус</th>
              <th className="border p-2 text-left">Дата выставления</th>
              <th className="border p-2 text-center">Действия</th>
            </tr>
          </thead>
          <tbody>
            {invoices?.map((invoice: Invoice) => (
              <tr key={invoice.id} className="hover:bg-gray-50">
                <td className="border p-2 font-semibold">{invoice.invoiceNumber}</td>
                <td className="border p-2">${invoice.totalAmount.toFixed(2)}</td>
                <td className="border p-2">
                  <select
                    value={invoice.status}
                    onChange={(e) => handleStatusChange(invoice.id, e.target.value as Invoice['status'])}
                    disabled={updateStatus.isPending}
                    className="px-2 py-1 border rounded"
                  >
                    <option value="draft">Черновик</option>
                    <option value="sent">Отправлен</option>
                    <option value="paid">Оплачен</option>
                    <option value="overdue">Просрочен</option>
                    <option value="cancelled">Отменен</option>
                  </select>
                </td>
                <td className="border p-2">{new Date(invoice.issueDate).toLocaleDateString()}</td>
                <td className="border p-2 text-center">
                  <button
                    onClick={() => handleDelete(invoice.id)}
                    disabled={deleteInvoice.isPending}
                    className="px-3 py-1 text-red-600 hover:bg-red-50 rounded disabled:opacity-50"
                  >
                    Удалить
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {invoices?.length === 0 && <p className="text-center text-gray-500 mt-6">Нет счетов</p>}
    </div>
  );
}
