import { useUIStore } from '@/store';
import { useCallback, useMemo } from 'react';

/**
 * useUI Hook - управление UI состоянием
 * 
 * Примеры использования:
 * const { sidebarOpen, theme, addNotification, openModal } = useUI();
 * 
 * // Manage sidebar
 * toggleSidebar();
 * 
 * // Manage theme
 * toggleTheme();
 * 
 * // Show notification
 * addNotification({
 *   type: 'success',
 *   message: 'User created successfully',
 *   duration: 3000,
 * });
 * 
 * // Open modal
 * openModal('confirm-delete', {
 *   title: 'Confirm Delete',
 *   content: 'Are you sure?',
 *   onConfirm: () => handleDelete(),
 * });
 */
export const useUI = () => {
  const {
    sidebarOpen,
    sidebarCollapsed,
    theme,
    notifications,
    modals,
    toggleSidebar,
    setSidebarOpen,
    setSidebarCollapsed,
    setTheme,
    toggleTheme,
    addNotification,
    removeNotification,
    clearNotifications,
    openModal,
    closeModal,
    updateModal,
    closeAllModals,
  } = useUIStore();

  const handleAddNotification = useCallback(
    (notification: any) => {
      return addNotification(notification);
    },
    [addNotification]
  );

  const handleRemoveNotification = useCallback(
    (id: string) => {
      removeNotification(id);
    },
    [removeNotification]
  );

  const handleOpenModal = useCallback(
    (id: string, modal: any) => {
      openModal(id, modal);
    },
    [openModal]
  );

  const handleCloseModal = useCallback(
    (id: string) => {
      closeModal(id);
    },
    [closeModal]
  );

  const handleUpdateModal = useCallback(
    (id: string, updates: any) => {
      updateModal(id, updates);
    },
    [updateModal]
  );

  const isDarkMode = useMemo(() => theme === 'dark', [theme]);

  return {
    // Sidebar
    sidebarOpen,
    sidebarCollapsed,
    toggleSidebar,
    setSidebarOpen,
    setSidebarCollapsed,

    // Theme
    theme,
    isDarkMode,
    setTheme,
    toggleTheme,

    // Notifications
    notifications,
    addNotification: handleAddNotification,
    removeNotification: handleRemoveNotification,
    clearNotifications,

    // Modals
    modals,
    openModal: handleOpenModal,
    closeModal: handleCloseModal,
    updateModal: handleUpdateModal,
    closeAllModals,
  };
};

export default useUI;
