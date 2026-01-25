import { create } from "zustand";
import { devtools } from "zustand/middleware";

export interface UserListState {
  id: string;
  email: string;
  name: string;
  role: "admin" | "manager" | "cashier" | "user";
  phone?: string;
  status: "active" | "inactive" | "suspended";
  createdAt: string;
  updatedAt: string;
}

interface UserFilter {
  role?: string;
  status?: string;
  search?: string;
}

interface UserStore {
  // State
  users: UserListState[];
  currentUser: UserListState | null;
  loading: boolean;
  error: string | null;
  filters: UserFilter;
  pagination: {
    page: number;
    pageSize: number;
    total: number;
  };

  // Queries
  filteredUsers: UserListState[];

  // Actions
  setUsers: (users: UserListState[]) => void;
  setCurrentUser: (user: UserListState) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  setFilters: (filters: UserFilter) => void;
  setPagination: (pagination: Partial<UserStore["pagination"]>) => void;
  clearError: () => void;

  // CRUD operations
  fetchUsers: (filters?: UserFilter) => Promise<void>;
  fetchUserById: (id: string) => Promise<void>;
  createUser: (
    data: Omit<UserListState, "id" | "createdAt" | "updatedAt">,
  ) => Promise<void>;
  updateUser: (id: string, data: Partial<UserListState>) => Promise<void>;
  deleteUser: (id: string) => Promise<void>;

  // Utility methods
  getUserById: (id: string) => UserListState | undefined;
}

const useUserStore = create<UserStore>()(
  devtools(
    (set, get) => ({
      // Initial state
      users: [],
      currentUser: null,
      loading: false,
      error: null,
      filters: {},
      pagination: {
        page: 1,
        pageSize: 10,
        total: 0,
      },

      // Computed property
      get filteredUsers() {
        const { users, filters } = get();
        let filtered = users;

        if (filters.role) {
          filtered = filtered.filter((u) => u.role === filters.role);
        }

        if (filters.status) {
          filtered = filtered.filter((u) => u.status === filters.status);
        }

        if (filters.search) {
          const search = filters.search.toLowerCase();
          filtered = filtered.filter(
            (u) =>
              u.name.toLowerCase().includes(search) ||
              u.email.toLowerCase().includes(search),
          );
        }

        return filtered;
      },

      // Basic actions
      setUsers: (users: UserListState[]) => {
        set({ users });
      },

      setCurrentUser: (user: UserListState) => {
        set({ currentUser: user });
      },

      setLoading: (loading: boolean) => {
        set({ loading });
      },

      setError: (error: string | null) => {
        set({ error });
      },

      setFilters: (filters: UserFilter) => {
        set({ filters: { ...get().filters, ...filters } });
      },

      setPagination: (pagination: Partial<UserStore["pagination"]>) => {
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
      fetchUsers: async (filters?: UserFilter) => {
        set({ loading: true, error: null });
        try {
          // Build query params
          const params = new URLSearchParams();
          if (filters?.role) params.append("role", filters.role);
          if (filters?.status) params.append("status", filters.status);
          if (filters?.search) params.append("search", filters.search);

          const response = await fetch(`/api/v1/users?${params.toString()}`);

          if (!response.ok) {
            const text = await response.text();
            let errorMessage = `Ошибка ${response.status}: ${response.statusText}`;
            try {
              const errorJson = JSON.parse(text);
              errorMessage =
                errorJson.message || errorJson.error || errorMessage;
            } catch {
              errorMessage = text || errorMessage;
            }
            console.error("[UserStore] Failed to fetch users:", errorMessage);
            throw new Error(errorMessage);
          }

          const data = await response.json();
          set({
            users: data.users,
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
          const message =
            error instanceof Error ? error.message : "Неизвестная ошибка";
          console.error("[UserStore] fetchUsers error:", message);
          set({ error: message, loading: false });
          throw error;
        }
      },

      fetchUserById: async (id: string) => {
        set({ loading: true, error: null });
        try {
          const response = await fetch(`/api/v1/users/${id}`);

          if (!response.ok) {
            throw new Error("Failed to fetch user");
          }

          const user = await response.json();
          set({ currentUser: user, loading: false });
        } catch (error) {
          const message =
            error instanceof Error ? error.message : "Unknown error";
          set({ error: message, loading: false });
          throw error;
        }
      },

      createUser: async (
        data: Omit<UserListState, "id" | "createdAt" | "updatedAt">,
      ) => {
        set({ loading: true, error: null });
        try {
          const response = await fetch("/api/v1/users", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
          });

          if (!response.ok) {
            throw new Error("Failed to create user");
          }

          const newUser = await response.json();
          set({
            users: [newUser, ...get().users],
            loading: false,
          });
        } catch (error) {
          const message =
            error instanceof Error ? error.message : "Unknown error";
          set({ error: message, loading: false });
          throw error;
        }
      },

      updateUser: async (id: string, data: Partial<UserListState>) => {
        set({ loading: true, error: null });
        try {
          const response = await fetch(`/api/v1/users/${id}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
          });

          if (!response.ok) {
            throw new Error("Failed to update user");
          }

          const updatedUser = await response.json();
          set({
            users: get().users.map((u) => (u.id === id ? updatedUser : u)),
            currentUser:
              get().currentUser?.id === id ? updatedUser : get().currentUser,
            loading: false,
          });
        } catch (error) {
          const message =
            error instanceof Error ? error.message : "Unknown error";
          set({ error: message, loading: false });
          throw error;
        }
      },

      deleteUser: async (id: string) => {
        set({ loading: true, error: null });
        try {
          const response = await fetch(`/api/v1/users/${id}`, {
            method: "DELETE",
          });

          if (!response.ok) {
            throw new Error("Failed to delete user");
          }

          set({
            users: get().users.filter((u) => u.id !== id),
            currentUser:
              get().currentUser?.id === id ? null : get().currentUser,
            loading: false,
          });
        } catch (error) {
          const message =
            error instanceof Error ? error.message : "Unknown error";
          set({ error: message, loading: false });
          throw error;
        }
      },

      // Utility
      getUserById: (id: string) => {
        return get().users.find((u) => u.id === id);
      },
    }),
    { name: "UserStore" },
  ),
);

export default useUserStore;
