import { create } from "zustand";
import { devtools } from "zustand/middleware";
import { bearerAuth } from "@/lib/auth/bearer";

// Get API Gateway URL from environment or use correct localhost port (80 for nginx)
const getApiUrl = (endpoint: string) => {
  if (typeof window === "undefined") {
    // Server-side
    return `http://localhost/api/v1${endpoint}`;
  }
  // Client-side - use NEXT_PUBLIC_API_URL if set
  if (process.env.NEXT_PUBLIC_API_URL) {
    return `${process.env.NEXT_PUBLIC_API_URL}${endpoint}`;
  }

  const hostname = window.location.hostname;
  const isLocalhost = hostname === "localhost" || hostname === "127.0.0.1";

  // For localhost, use port 80 (nginx-proxy reverse proxy), not direct api-gateway port
  if (isLocalhost) {
    return `http://${hostname}/api/v1${endpoint}`;
  }

  // For other hosts, use https if current page is https, otherwise http
  const protocol = window.location.protocol;
  return `${protocol}//${hostname}/api/v1${endpoint}`;
};

export interface CatalogItem {
  id: string;
  name: string;
  sku: string;
  description?: string;
  category: string;
  price: number;
  costPrice: number;
  margin: number; // percentage
  profit: number; // calculated: (price - costPrice) * margin / 100
  stock: number;
  reorderLevel: number;
  image?: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

interface CatalogFilter {
  category?: string;
  search?: string;
  active?: boolean;
}

interface CatalogStore {
  // State
  items: CatalogItem[];
  currentItem: CatalogItem | null;
  categories: string[];
  loading: boolean;
  error: string | null;
  filters: CatalogFilter;
  pagination: {
    page: number;
    pageSize: number;
    total: number;
  };

  // Computed
  filteredItems: CatalogItem[];
  totalInventoryValue: number;
  lowStockItems: CatalogItem[];

  // Actions
  setItems: (items: CatalogItem[]) => void;
  setCurrentItem: (item: CatalogItem) => void;
  setCategories: (categories: string[]) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setFilters: (filters: CatalogFilter) => void;
  setPagination: (pagination: Partial<CatalogStore["pagination"]>) => void;
  clearError: () => void;

  // CRUD operations
  fetchCatalog: (filters?: CatalogFilter) => Promise<void>;
  fetchItemById: (id: string) => Promise<void>;
  createItem: (
    data: Omit<CatalogItem, "id" | "createdAt" | "updatedAt">,
  ) => Promise<void>;
  updateItem: (id: string, data: Partial<CatalogItem>) => Promise<void>;
  deleteItem: (id: string) => Promise<void>;

  // Utility
  getItemById: (id: string) => CatalogItem | undefined;
  getItemsByCategory: (category: string) => CatalogItem[];
  calculateProfit: (costPrice: number, price: number, margin: number) => number;
}

const useCatalogStore = create<CatalogStore>()(
  devtools(
    (set, get) => ({
      // Initial state
      items: [],
      currentItem: null,
      categories: [],
      loading: false,
      error: null,
      filters: {},
      pagination: {
        page: 1,
        pageSize: 10,
        total: 0,
      },

      // Computed properties
      get filteredItems() {
        const { items, filters } = get();
        let filtered = items;

        if (filters.category) {
          filtered = filtered.filter((i) => i.category === filters.category);
        }

        if (filters.search) {
          const search = filters.search.toLowerCase();
          filtered = filtered.filter(
            (i) =>
              i.name.toLowerCase().includes(search) ||
              i.sku.toLowerCase().includes(search) ||
              i.description?.toLowerCase().includes(search),
          );
        }

        if (filters.active !== undefined) {
          filtered = filtered.filter((i) => i.active === filters.active);
        }

        return filtered;
      },

      get totalInventoryValue() {
        return get().items.reduce(
          (sum, item) => sum + item.price * item.stock,
          0,
        );
      },

      get lowStockItems() {
        return get().items.filter((item) => item.stock <= item.reorderLevel);
      },

      // Basic actions
      setItems: (items: CatalogItem[]) => {
        set({ items });
        // Extract unique categories
        const categories = [...new Set(items.map((i) => i.category))];
        set({ categories });
      },

      setCurrentItem: (item: CatalogItem) => {
        set({ currentItem: item });
      },

      setCategories: (categories: string[]) => {
        set({ categories });
      },

      setLoading: (loading: boolean) => {
        set({ loading });
      },

      setError: (error: string | null) => {
        set({ error });
      },

      setFilters: (filters: CatalogFilter) => {
        set({ filters: { ...get().filters, ...filters } });
      },

      setPagination: (pagination: Partial<CatalogStore["pagination"]>) => {
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
      fetchCatalog: async (filters?: CatalogFilter) => {
        set({ loading: true, error: null });
        try {
          const params = new URLSearchParams();
          if (filters?.category) params.append("category", filters.category);
          if (filters?.search) params.append("search", filters.search);
          if (filters?.active !== undefined)
            params.append("active", String(filters.active));

          const authHeaders = bearerAuth();

          const response = await fetch(
            getApiUrl(`/catalogs-query?${params.toString()}`),
            {
              headers: authHeaders,
              credentials: "include",
            },
          );

          if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            const errorMessage =
              errorData.error || errorData.details || `HTTP ${response.status}`;
            throw new Error(`Failed to fetch catalog: ${errorMessage}`);
          }

          const apiResponse = await response.json();

          // Извлечение товаров из различных возможных форматов
          let items = [];

          // Обработка реального формата API: { Data: { CatalogList: { catalogs: [...] } } }
          if (apiResponse.Data?.CatalogList?.catalogs) {
            items = apiResponse.Data.CatalogList.catalogs;
          }
          // Обработка формата ответа: { success, data: {...}, meta: {...} }
          else if (apiResponse.success && apiResponse.data) {
            const data = apiResponse.data || {};
            if (Array.isArray(data.catalogs)) {
              items = data.catalogs;
            } else if (Array.isArray(data.items)) {
              items = data.items;
            } else if (Array.isArray(data.catalog_items)) {
              items = data.catalog_items;
            } else if (Array.isArray(data.products)) {
              items = data.products;
            } else if (Array.isArray(data)) {
              items = data;
            }
          }
          // Обработка прямого массива
          else if (Array.isArray(apiResponse)) {
            items = apiResponse;
          }
          // Обработка формата { items: [...] }
          else if (Array.isArray(apiResponse.items)) {
            items = apiResponse.items;
          }
          // Обработка формата с ошибкой
          else if (apiResponse.error) {
            throw new Error(
              apiResponse.error?.message || "Ошибка при загрузке каталога",
            );
          }

          // Трансформируем данные из API в формат компонента
          const transformedItems = items.map((item: any, index: number) => {
            // Используем категорию из API, по умолчанию "Без категории"
            const category = item.category || "Без категории";

            return {
              id: item.id,
              name: item.name || `Товар #${index + 1}`,
              sku: item.number || item.sku || "",
              description: item.description || "",
              category,
              price: item.price || 0,
              costPrice: item.cost || item.costPrice || 0,
              margin: item.margin || 0,
              profit: item.profit || 0,
              stock: item.quantity || item.stock || 0,
              reorderLevel: item.min_quantity || item.reorderLevel || 0,
              image: item.image,
              active: item.is_active !== false && item.status !== "inactive",
              createdAt: item.created_at || new Date().toISOString(),
              updatedAt: item.updated_at || new Date().toISOString(),
            };
          });

          set({
            items: transformedItems.length > 0 ? transformedItems : [],
            pagination: {
              page:
                apiResponse.Data?.CatalogList?.page ||
                apiResponse.meta?.page ||
                0,
              pageSize:
                apiResponse.Data?.CatalogList?.size ||
                apiResponse.meta?.page_size ||
                20,
              total:
                apiResponse.Data?.CatalogList?.total_elements ||
                apiResponse.meta?.total_count ||
                transformedItems.length,
            },
            loading: false,
          });

          // Extract categories
          const categories = [
            ...new Set(transformedItems.map((i: CatalogItem) => i.category)),
          ] as string[];
          set({ categories });

          if (filters) {
            set({ filters });
          }
        } catch (error) {
          const message =
            error instanceof Error ? error.message : "Неизвестная ошибка";
          console.error("Ошибка загрузки каталога:", message);
          set({ error: message, loading: false });
          throw error;
        }
      },

      fetchItemById: async (id: string) => {
        set({ loading: true, error: null });
        try {
          const authHeaders = bearerAuth();

          const response = await fetch(getApiUrl(`/catalog/${id}`), {
            headers: authHeaders,
            credentials: "include",
          });

          if (!response.ok) {
            throw new Error("Failed to fetch item");
          }

          const item = await response.json();
          set({ currentItem: item, loading: false });
        } catch (error) {
          const message =
            error instanceof Error ? error.message : "Unknown error";
          set({ error: message, loading: false });
          throw error;
        }
      },

      createItem: async (
        data: Omit<CatalogItem, "id" | "createdAt" | "updatedAt">,
      ) => {
        set({ loading: true, error: null });
        try {
          const authHeaders = bearerAuth();

          // Трансформируем данные формы в формат API
          const apiData = {
            name: data.name,
            number: (data as any).sku || "", // sku → number
            description: data.description || "",
            tnved_code: (data as any).category || "", // category → tnved_code
            price: data.price || 0,
            currency: (data as any).currency || "RUB", // Добавляем валюту по умолчанию
            unit: (data as any).unit || "шт", // Добавляем единицу по умолчанию
          };

          const response = await fetch(getApiUrl("/catalog"), {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              ...authHeaders,
            },
            credentials: "include",
            body: JSON.stringify(apiData),
          });

          if (!response.ok) {
            throw new Error("Failed to create item");
          }

          const newItem = await response.json();
          set({
            items: [newItem, ...get().items],
            loading: false,
          });
        } catch (error) {
          const message =
            error instanceof Error ? error.message : "Unknown error";
          set({ error: message, loading: false });
          throw error;
        }
      },

      updateItem: async (id: string, data: Partial<CatalogItem>) => {
        set({ loading: true, error: null });
        try {
          const authHeaders = bearerAuth();

          const response = await fetch(getApiUrl(`/catalog/${id}`), {
            method: "PUT",
            headers: {
              "Content-Type": "application/json",
              ...authHeaders,
            },
            credentials: "include",
            body: JSON.stringify(data),
          });

          if (!response.ok) {
            throw new Error("Failed to update item");
          }

          const updatedItem = await response.json();
          set({
            items: get().items.map((i) => (i.id === id ? updatedItem : i)),
            currentItem:
              get().currentItem?.id === id ? updatedItem : get().currentItem,
            loading: false,
          });
        } catch (error) {
          const message =
            error instanceof Error ? error.message : "Unknown error";
          set({ error: message, loading: false });
          throw error;
        }
      },

      deleteItem: async (id: string) => {
        set({ loading: true, error: null });
        try {
          const authHeaders = bearerAuth();

          const response = await fetch(getApiUrl(`/catalog/${id}`), {
            method: "DELETE",
            headers: authHeaders,
            credentials: "include",
          });

          if (!response.ok) {
            throw new Error("Failed to delete item");
          }

          set({
            items: get().items.filter((i) => i.id !== id),
            currentItem:
              get().currentItem?.id === id ? null : get().currentItem,
            loading: false,
          });
        } catch (error) {
          const message =
            error instanceof Error ? error.message : "Unknown error";
          set({ error: message, loading: false });
          throw error;
        }
      },

      // Utilities
      getItemById: (id: string) => {
        return get().items.find((i) => i.id === id);
      },

      getItemsByCategory: (category: string) => {
        return get().items.filter((i) => i.category === category);
      },

      calculateProfit: (costPrice: number, price: number, margin: number) => {
        return (price - costPrice) * (margin / 100);
      },
    }),
    { name: "CatalogStore" },
  ),
);

export default useCatalogStore;
