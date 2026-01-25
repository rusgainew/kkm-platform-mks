'use client';

import { useUsersQuery, useCreateUserMutation, useDeleteUserMutation, type User } from '@/lib/hooks/useUsersQuery';
import { useState } from 'react';

export function UserListWithReactQuery() {
  const { data: users, isLoading, error } = useUsersQuery();
  const createMutation = useCreateUserMutation();
  const deleteMutation = useDeleteUserMutation();
  const [newUserName, setNewUserName] = useState('');
  const [newUserEmail, setNewUserEmail] = useState('');

  if (isLoading) return <div className="p-4">Загрузка...</div>;
  if (error) return <div className="p-4 text-red-600">Ошибка: {error.message}</div>;

  const handleAddUser = async () => {
    if (!newUserName || !newUserEmail) return;
    
    await createMutation.mutateAsync({
      name: newUserName,
      email: newUserEmail,
      role: 'user',
      status: 'active',
    });
    
    setNewUserName('');
    setNewUserEmail('');
  };

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Пользователи (React Query)</h1>

      {/* Добавление нового пользователя */}
      <div className="mb-6 p-4 border rounded">
        <h2 className="text-lg font-semibold mb-4">Добавить пользователя</h2>
        <div className="space-y-3">
          <input
            type="text"
            placeholder="Имя"
            value={newUserName}
            onChange={(e) => setNewUserName(e.target.value)}
            className="w-full px-3 py-2 border rounded"
          />
          <input
            type="email"
            placeholder="Email"
            value={newUserEmail}
            onChange={(e) => setNewUserEmail(e.target.value)}
            className="w-full px-3 py-2 border rounded"
          />
          <button
            onClick={handleAddUser}
            disabled={createMutation.isPending}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
          >
            {createMutation.isPending ? 'Добавление...' : 'Добавить'}
          </button>
        </div>
      </div>

      {/* Список пользователей */}
      <div className="space-y-2">
        {users?.map((user: User) => (
          <div key={user.id} className="flex items-center justify-between p-3 border rounded hover:bg-gray-50">
            <div>
              <p className="font-semibold">{user.name}</p>
              <p className="text-sm text-gray-600">{user.email}</p>
              <span className="text-xs px-2 py-1 bg-gray-200 rounded">{user.role}</span>
            </div>
            <button
              onClick={() => deleteMutation.mutate(user.id)}
              disabled={deleteMutation.isPending}
              className="px-3 py-1 text-red-600 hover:bg-red-50 rounded disabled:opacity-50"
            >
              {deleteMutation.isPending ? 'Удаление...' : 'Удалить'}
            </button>
          </div>
        ))}
      </div>

      {users?.length === 0 && <p className="text-center text-gray-500 mt-6">Нет пользователей</p>}
    </div>
  );
}
