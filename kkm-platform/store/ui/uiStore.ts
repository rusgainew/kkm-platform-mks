import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';

export interface Notification {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  message: string;
  duration?: number;
  action?: {
    label: string;
    onClick: () => void;
  };
}

interface Modal {
  id: string;
  isOpen: boolean;
  title: string;
  content?: React.ReactNode;
  onConfirm?: () => void;
  onCancel?: () => void;
}

interface UIStore {
  // State
  sidebarOpen: boolean;
  sidebarCollapsed: boolean;
  theme: 'light' | 'dark';
  notifications: Notification[];
  modals: Record<string, Modal>;

  // Actions - Sidebar
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  setSidebarCollapsed: (collapsed: boolean) => void;

  // Actions - Theme
  setTheme: (theme: 'light' | 'dark') => void;
  toggleTheme: () => void;

  // Actions - Notifications
  addNotification: (notification: Omit<Notification, 'id'>) => string;
  removeNotification: (id: string) => void;
  clearNotifications: () => void;

  // Actions - Modals
  openModal: (id: string, modal: Omit<Modal, 'isOpen'>) => void;
  closeModal: (id: string) => void;
  updateModal: (id: string, updates: Partial<Modal>) => void;
  closeAllModals: () => void;
}

const useUIStore = create<UIStore>()(
  devtools(
    persist(
      (set, get) => ({
        // Initial state
        sidebarOpen: true,
        sidebarCollapsed: false,
        theme: 'dark',
        notifications: [],
        modals: {},

        // Sidebar actions
        toggleSidebar: () => {
          set({ sidebarOpen: !get().sidebarOpen });
        },

        setSidebarOpen: (open: boolean) => {
          set({ sidebarOpen: open });
        },

        setSidebarCollapsed: (collapsed: boolean) => {
          set({ sidebarCollapsed: collapsed });
        },

        // Theme actions
        setTheme: (theme: 'light' | 'dark') => {
          set({ theme });
          // Apply theme to document
          if (typeof window !== 'undefined') {
            const html = document.documentElement;
            if (theme === 'dark') {
              html.classList.add('dark');
            } else {
              html.classList.remove('dark');
            }
          }
        },

        toggleTheme: () => {
          const currentTheme = get().theme;
          const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
          get().setTheme(newTheme);
        },

        // Notification actions
        addNotification: (notification: Omit<Notification, 'id'>) => {
          const id = `notification-${Date.now()}`;
          const fullNotification: Notification = {
            ...notification,
            id,
          };

          set({
            notifications: [fullNotification, ...get().notifications],
          });

          // Auto-remove notification after duration
          if (notification.duration && notification.duration > 0) {
            setTimeout(() => {
              get().removeNotification(id);
            }, notification.duration);
          }

          return id;
        },

        removeNotification: (id: string) => {
          set({
            notifications: get().notifications.filter(n => n.id !== id),
          });
        },

        clearNotifications: () => {
          set({ notifications: [] });
        },

        // Modal actions
        openModal: (id: string, modal: Omit<Modal, 'isOpen'>) => {
          set({
            modals: {
              ...get().modals,
              [id]: {
                ...modal,
                isOpen: true,
              },
            },
          });
        },

        closeModal: (id: string) => {
          set({
            modals: {
              ...get().modals,
              [id]: {
                ...get().modals[id],
                isOpen: false,
              },
            },
          });
        },

        updateModal: (id: string, updates: Partial<Modal>) => {
          set({
            modals: {
              ...get().modals,
              [id]: {
                ...get().modals[id],
                ...updates,
              },
            },
          });
        },

        closeAllModals: () => {
          const modals = get().modals;
          const updated: Record<string, Modal> = {};
          Object.keys(modals).forEach(key => {
            updated[key] = { ...modals[key], isOpen: false };
          });
          set({ modals: updated });
        },
      }),
      {
        name: 'ui-store',
        partialize: (state) => ({
          sidebarOpen: state.sidebarOpen,
          sidebarCollapsed: state.sidebarCollapsed,
          theme: state.theme,
        }),
      }
    ),
    { name: 'UIStore' }
  )
);

export default useUIStore;
