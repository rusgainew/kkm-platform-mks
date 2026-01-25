import { useInvoiceStore } from '@/store';
import { useCallback } from 'react';

interface InvoiceFilter {
  status?: string;
  companyId?: string;
  search?: string;
  dateFrom?: string;
  dateTo?: string;
}

/**
 * useInvoices Hook - управление счетами
 * 
 * Примеры использования:
 * const { invoices, fetchInvoices, createInvoice, updateInvoice } = useInvoices();
 * 
 * // Fetch with filters
 * await fetchInvoices({ status: 'pending', dateFrom: '2026-01-01' });
 * 
 * // Create
 * await createInvoice({ number, companyId, items, totalAmount });
 * 
 * // Update
 * await updateInvoice(invoiceId, { status });
 * 
 * // Statistics
 * const { totalRevenue, paidAmount, pendingAmount } = useInvoices();
 */
export const useInvoices = () => {
  const {
    invoices,
    filteredInvoices,
    totalRevenue,
    paidAmount,
    pendingAmount,
    loading,
    error,
    filters,
    pagination,
    fetchInvoices,
    fetchInvoiceById,
    createInvoice,
    updateInvoice,
    deleteInvoice,
    updateInvoiceStatus,
    setFilters,
    setPagination,
    clearError,
    getInvoicesByStatus,
  } = useInvoiceStore();

  const handleFetchInvoices = useCallback(
    async (newFilters?: InvoiceFilter) => {
      try {
        await fetchInvoices(newFilters);
        return true;
      } catch (error) {
        return false;
      }
    },
    [fetchInvoices]
  );

  const handleCreateInvoice = useCallback(
    async (data: any) => {
      try {
        await createInvoice(data);
        return true;
      } catch (error) {
        return false;
      }
    },
    [createInvoice]
  );

  const handleUpdateInvoice = useCallback(
    async (id: string, data: any) => {
      try {
        await updateInvoice(id, data);
        return true;
      } catch (error) {
        return false;
      }
    },
    [updateInvoice]
  );

  const handleDeleteInvoice = useCallback(
    async (id: string) => {
      try {
        await deleteInvoice(id);
        return true;
      } catch (error) {
        return false;
      }
    },
    [deleteInvoice]
  );

  const handleUpdateStatus = useCallback(
    async (id: string, status: any) => {
      try {
        await updateInvoiceStatus(id, status);
        return true;
      } catch (error) {
        return false;
      }
    },
    [updateInvoiceStatus]
  );

  return {
    invoices,
    filteredInvoices,
    totalRevenue,
    paidAmount,
    pendingAmount,
    isLoading: loading,
    error,
    filters,
    pagination,
    fetchInvoices: handleFetchInvoices,
    fetchInvoiceById,
    createInvoice: handleCreateInvoice,
    updateInvoice: handleUpdateInvoice,
    deleteInvoice: handleDeleteInvoice,
    updateInvoiceStatus: handleUpdateStatus,
    setFilters,
    setPagination,
    clearError,
    getInvoicesByStatus,
  };
};

export default useInvoices;
