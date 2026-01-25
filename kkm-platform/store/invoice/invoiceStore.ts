import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

export interface InvoiceItem {
  id: string;
  productId: string;
  productName: string;
  quantity: number;
  price: number;
  tax?: number;
  total: number;
}

export interface Invoice {
  id: string;
  number: string;
  companyId: string;
  companyName: string;
  items: InvoiceItem[];
  status: 'draft' | 'pending' | 'paid' | 'overdue' | 'cancelled';
  totalAmount: number;
  taxAmount: number;
  dueDate: string;
  issuedDate: string;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

interface InvoiceFilter {
  status?: string;
  companyId?: string;
  search?: string;
  dateFrom?: string;
  dateTo?: string;
}

interface InvoiceStore {
  // State
  invoices: Invoice[];
  currentInvoice: Invoice | null;
  loading: boolean;
  error: string | null;
  filters: InvoiceFilter;
  pagination: {
    page: number;
    pageSize: number;
    total: number;
  };

  // Computed
  filteredInvoices: Invoice[];
  totalRevenue: number;
  paidAmount: number;
  pendingAmount: number;

  // Actions
  setInvoices: (invoices: Invoice[]) => void;
  setCurrentInvoice: (invoice: Invoice) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setFilters: (filters: InvoiceFilter) => void;
  setPagination: (pagination: Partial<InvoiceStore['pagination']>) => void;
  clearError: () => void;

  // CRUD operations
  fetchInvoices: (filters?: InvoiceFilter) => Promise<void>;
  fetchInvoiceById: (id: string) => Promise<void>;
  createInvoice: (data: Omit<Invoice, 'id' | 'createdAt' | 'updatedAt'>) => Promise<void>;
  updateInvoice: (id: string, data: Partial<Invoice>) => Promise<void>;
  deleteInvoice: (id: string) => Promise<void>;
  updateInvoiceStatus: (id: string, status: Invoice['status']) => Promise<void>;

  // Utility
  getInvoiceById: (id: string) => Invoice | undefined;
  getInvoicesByStatus: (status: Invoice['status']) => Invoice[];
}

const useInvoiceStore = create<InvoiceStore>()(
  devtools(
    (set, get) => ({
      // Initial state
      invoices: [],
      currentInvoice: null,
      loading: false,
      error: null,
      filters: {},
      pagination: {
        page: 1,
        pageSize: 10,
        total: 0,
      },

      // Computed properties
      get filteredInvoices() {
        const { invoices, filters } = get();
        let filtered = invoices;

        if (filters.status) {
          filtered = filtered.filter(i => i.status === filters.status);
        }

        if (filters.companyId) {
          filtered = filtered.filter(i => i.companyId === filters.companyId);
        }

        if (filters.search) {
          const search = filters.search.toLowerCase();
          filtered = filtered.filter(i =>
            i.number.toLowerCase().includes(search) ||
            i.companyName.toLowerCase().includes(search)
          );
        }

        if (filters.dateFrom) {
          filtered = filtered.filter(i => new Date(i.issuedDate) >= new Date(filters.dateFrom!));
        }

        if (filters.dateTo) {
          filtered = filtered.filter(i => new Date(i.issuedDate) <= new Date(filters.dateTo!));
        }

        return filtered;
      },

      get totalRevenue() {
        return get().invoices.reduce((sum, i) => sum + i.totalAmount, 0);
      },

      get paidAmount() {
        return get().invoices
          .filter(i => i.status === 'paid')
          .reduce((sum, i) => sum + i.totalAmount, 0);
      },

      get pendingAmount() {
        return get().invoices
          .filter(i => ['pending', 'overdue'].includes(i.status))
          .reduce((sum, i) => sum + i.totalAmount, 0);
      },

      // Basic actions
      setInvoices: (invoices: Invoice[]) => {
        set({ invoices });
      },

      setCurrentInvoice: (invoice: Invoice) => {
        set({ currentInvoice: invoice });
      },

      setLoading: (loading: boolean) => {
        set({ loading });
      },

      setError: (error: string | null) => {
        set({ error });
      },

      setFilters: (filters: InvoiceFilter) => {
        set({ filters: { ...get().filters, ...filters } });
      },

      setPagination: (pagination: Partial<InvoiceStore['pagination']>) => {
        set({
          pagination: {
            ...get().pagination,
            ...pagination,
          },
        });
      },

      clearError: () => {
        set({ error: null });
      },

      // CRUD operations
      fetchInvoices: async (filters?: InvoiceFilter) => {
        set({ loading: true, error: null });
        try {
          const params = new URLSearchParams();
          if (filters?.status) params.append('status', filters.status);
          if (filters?.companyId) params.append('companyId', filters.companyId);
          if (filters?.search) params.append('search', filters.search);
          if (filters?.dateFrom) params.append('dateFrom', filters.dateFrom);
          if (filters?.dateTo) params.append('dateTo', filters.dateTo);

          const response = await fetch(`/api/v1/invoices?${params.toString()}`);

          if (!response.ok) {
            throw new Error('Failed to fetch invoices');
          }

          const data = await response.json();
          set({
            invoices: data.invoices,
            pagination: {
              page: data.pagination.page,
              pageSize: data.pagination.pageSize,
              total: data.pagination.total,
            },
            loading: false,
          });

          if (filters) {
            set({ filters });
          }
        } catch (error) {
          const message = error instanceof Error ? error.message : 'Unknown error';
          set({ error: message, loading: false });
          throw error;
        }
      },

      fetchInvoiceById: async (id: string) => {
        set({ loading: true, error: null });
        try {
          const response = await fetch(`/api/v1/invoices/${id}`);

          if (!response.ok) {
            throw new Error('Failed to fetch invoice');
          }

          const invoice = await response.json();
          set({ currentInvoice: invoice, loading: false });
        } catch (error) {
          const message = error instanceof Error ? error.message : 'Unknown error';
          set({ error: message, loading: false });
          throw error;
        }
      },

      createInvoice: async (data: Omit<Invoice, 'id' | 'createdAt' | 'updatedAt'>) => {
        set({ loading: true, error: null });
        try {
          const response = await fetch('/api/v1/invoices', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
          });

          if (!response.ok) {
            throw new Error('Failed to create invoice');
          }

          const newInvoice = await response.json();
          set({
            invoices: [newInvoice, ...get().invoices],
            loading: false,
          });
        } catch (error) {
          const message = error instanceof Error ? error.message : 'Unknown error';
          set({ error: message, loading: false });
          throw error;
        }
      },

      updateInvoice: async (id: string, data: Partial<Invoice>) => {
        set({ loading: true, error: null });
        try {
          const response = await fetch(`/api/v1/invoices/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
          });

          if (!response.ok) {
            throw new Error('Failed to update invoice');
          }

          const updatedInvoice = await response.json();
          set({
            invoices: get().invoices.map(i => i.id === id ? updatedInvoice : i),
            currentInvoice: get().currentInvoice?.id === id ? updatedInvoice : get().currentInvoice,
            loading: false,
          });
        } catch (error) {
          const message = error instanceof Error ? error.message : 'Unknown error';
          set({ error: message, loading: false });
          throw error;
        }
      },

      deleteInvoice: async (id: string) => {
        set({ loading: true, error: null });
        try {
          const response = await fetch(`/api/v1/invoices/${id}`, {
            method: 'DELETE',
          });

          if (!response.ok) {
            throw new Error('Failed to delete invoice');
          }

          set({
            invoices: get().invoices.filter(i => i.id !== id),
            currentInvoice: get().currentInvoice?.id === id ? null : get().currentInvoice,
            loading: false,
          });
        } catch (error) {
          const message = error instanceof Error ? error.message : 'Unknown error';
          set({ error: message, loading: false });
          throw error;
        }
      },

      updateInvoiceStatus: async (id: string, status: Invoice['status']) => {
        set({ loading: true, error: null });
        try {
          const response = await fetch(`/api/v1/invoices/${id}/status`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ status }),
          });

          if (!response.ok) {
            throw new Error('Failed to update invoice status');
          }

          const updatedInvoice = await response.json();
          set({
            invoices: get().invoices.map(i => i.id === id ? updatedInvoice : i),
            currentInvoice: get().currentInvoice?.id === id ? updatedInvoice : get().currentInvoice,
            loading: false,
          });
        } catch (error) {
          const message = error instanceof Error ? error.message : 'Unknown error';
          set({ error: message, loading: false });
          throw error;
        }
      },

      // Utilities
      getInvoiceById: (id: string) => {
        return get().invoices.find(i => i.id === id);
      },

      getInvoicesByStatus: (status: Invoice['status']) => {
        return get().invoices.filter(i => i.status === status);
      },
    }),
    { name: 'InvoiceStore' }
  )
);

export default useInvoiceStore;
