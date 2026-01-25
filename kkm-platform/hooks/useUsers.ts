import { useUserStore } from '@/store';
import { useCallback } from 'react';

interface UserFilter {
  role?: string;
  status?: string;
  search?: string;
}

/**
 * useUsers Hook - управление пользователями
 * 
 * Примеры использования:
 * const { users, fetchUsers, createUser, updateUser, deleteUser } = useUsers();
 * 
 * // Fetch with filters
 * await fetchUsers({ role: 'admin', status: 'active' });
 * 
 * // Create
 * await createUser({ email, name, role, phone });
 * 
 * // Update
 * await updateUser(userId, { name, role });
 * 
 * // Delete
 * await deleteUser(userId);
 */
export const useUsers = () => {
  const {
    users,
    filteredUsers,
    loading,
    error,
    filters,
    pagination,
    fetchUsers,
    fetchUserById,
    createUser,
    updateUser,
    deleteUser,
    setFilters,
    setPagination,
    clearError,
  } = useUserStore();

  const handleFetchUsers = useCallback(
    async (newFilters?: UserFilter) => {
      try {
        await fetchUsers(newFilters);
        return true;
      } catch (error) {
        return false;
      }
    },
    [fetchUsers]
  );

  const handleCreateUser = useCallback(
    async (data: any) => {
      try {
        await createUser(data);
        return true;
      } catch (error) {
        return false;
      }
    },
    [createUser]
  );

  const handleUpdateUser = useCallback(
    async (id: string, data: any) => {
      try {
        await updateUser(id, data);
        return true;
      } catch (error) {
        return false;
      }
    },
    [updateUser]
  );

  const handleDeleteUser = useCallback(
    async (id: string) => {
      try {
        await deleteUser(id);
        return true;
      } catch (error) {
        return false;
      }
    },
    [deleteUser]
  );

  return {
    users,
    filteredUsers,
    isLoading: loading,
    error,
    filters,
    pagination,
    fetchUsers: handleFetchUsers,
    fetchUserById,
    createUser: handleCreateUser,
    updateUser: handleUpdateUser,
    deleteUser: handleDeleteUser,
    setFilters,
    setPagination,
    clearError,
  };
};

export default useUsers;
