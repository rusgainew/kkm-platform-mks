import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';

export interface User {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'manager' | 'cashier' | 'user';
  phone?: string;
  avatar?: string;
  createdAt: string;
  updatedAt: string;
}

interface AuthStore {
  // State
  user: User | null;
  token: string | null;
  refreshToken: string | null;
  loading: boolean;
  error: string | null;
  isAuthenticated: boolean;

  // Actions
  setUser: (user: User) => void;
  setToken: (token: string, refreshToken?: string) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  
  // Auth methods
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  refreshAccessToken: () => Promise<void>;
  checkAuth: () => Promise<void>;
}

const useAuthStore = create<AuthStore>()(
  devtools(
    persist(
      (set, get) => ({
        // Initial state
        user: null,
        token: null,
        refreshToken: null,
        loading: false,
        error: null,
        isAuthenticated: false,

        // Basic actions
        setUser: (user: User) => {
          set({ user, isAuthenticated: true });
        },

        setToken: (token: string, refreshToken?: string) => {
          set({ 
            token,
            refreshToken: refreshToken || get().refreshToken,
          });
          // Save to localStorage
          if (typeof window !== 'undefined') {
            localStorage.setItem('authToken', token);
            if (refreshToken) {
              localStorage.setItem('refreshToken', refreshToken);
            }
          }
        },

        setLoading: (loading: boolean) => {
          set({ loading });
        },

        setError: (error: string | null) => {
          set({ error });
        },

        clearError: () => {
          set({ error: null });
        },

        // Auth methods
        login: async (email: string, password: string) => {
          set({ loading: true, error: null });
          try {
            // API call to backend
            const response = await fetch('/api/v1/auth/login', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({ email, password }),
              credentials: 'include', // Include cookies
            });

            if (!response.ok) {
              const errorData = await response.json();
              throw new Error(errorData.message || 'Login failed');
            }

            const data = await response.json();
            const { accessToken, refreshToken, user } = data;

            get().setToken(accessToken, refreshToken);
            get().setUser(user);
            set({ loading: false, isAuthenticated: true });
          } catch (error) {
            const message = error instanceof Error ? error.message : 'Unknown error';
            set({ 
              error: message,
              loading: false,
              isAuthenticated: false,
            });
            throw error;
          }
        },

        logout: () => {
          // API call to logout endpoint (optional)
          if (typeof window !== 'undefined') {
            localStorage.removeItem('authToken');
            localStorage.removeItem('refreshToken');
          }
          
          set({
            user: null,
            token: null,
            refreshToken: null,
            isAuthenticated: false,
            error: null,
          });
        },

        refreshAccessToken: async () => {
          try {
            const refreshToken = get().refreshToken;
            if (!refreshToken) {
              throw new Error('No refresh token available');
            }

            const response = await fetch('/api/v1/auth/refresh', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({ refreshToken }),
              credentials: 'include',
            });

            if (!response.ok) {
              get().logout();
              throw new Error('Token refresh failed');
            }

            const data = await response.json();
            get().setToken(data.accessToken, data.refreshToken);
          } catch (error) {
            get().logout();
            throw error;
          }
        },

        checkAuth: async () => {
          set({ loading: true });
          try {
            const token = get().token;
            if (!token) {
              throw new Error('No token available');
            }

            const response = await fetch('/api/v1/users/me', {
              headers: {
                'Authorization': `Bearer ${token}`,
              },
              credentials: 'include',
            });

            if (!response.ok) {
              if (response.status === 401) {
                await get().refreshAccessToken();
                // Retry the request
                return get().checkAuth();
              }
              throw new Error('Failed to fetch user data');
            }

            const user = await response.json();
            get().setUser(user);
            set({ isAuthenticated: true, loading: false });
          } catch (error) {
            const message = error instanceof Error ? error.message : 'Unknown error';
            set({ 
              error: message,
              loading: false,
              isAuthenticated: false,
            });
          }
        },
      }),
      {
        name: 'auth-store', // localStorage key
        partialize: (state) => ({
          user: state.user,
          token: state.token,
          refreshToken: state.refreshToken,
          isAuthenticated: state.isAuthenticated,
        }),
      }
    ),
    { name: 'AuthStore' }
  )
);

export default useAuthStore;
