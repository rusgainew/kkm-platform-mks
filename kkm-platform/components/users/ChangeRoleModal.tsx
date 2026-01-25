'use client';

import { useState } from 'react';
import { X, Loader2 } from 'lucide-react';
import type { ApiUser } from '@/lib/api/users';
import { assignUserRole } from '@/lib/api/users';

interface ChangeRoleModalProps {
  user: ApiUser | null;
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
}

const ROLE_OPTIONS = [
  { value: 'manager', label: 'Менеджер' },
  { value: 'admin', label: 'Администратор' },
  { value: 'cashier', label: 'Кассир' },
];

export default function ChangeRoleModal({ user, isOpen, onClose, onSuccess }: ChangeRoleModalProps) {
  const [selectedRoleValue, setSelectedRoleValue] = useState<string>('cashier');
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) return;

    setError(null);
    setIsLoading(true);
    try {
      await assignUserRole(user.id, {
        role: selectedRoleValue as 'user' | 'manager',
      });
      onSuccess?.();
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to change role');
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-gray-900 rounded-lg shadow-xl max-w-md w-full mx-4">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-700">
          <h2 className="text-xl font-semibold text-white">Изменить роль</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          {error && (
            <div className="p-3 bg-red-900/20 border border-red-800 rounded-lg text-sm text-red-400">
              {error}
            </div>
          )}

          <div>
            <p className="text-sm text-gray-400 mb-2">Пользователь:</p>
            <p className="text-white font-medium">
              {user?.first_name} {user?.last_name}
            </p>
            <p className="text-sm text-gray-400">{user?.email}</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Новая роль
            </label>
            <select
              value={isOpen && user ? (selectedRoleValue || user.role || 'user') : 'user'}
              onChange={(e) => setSelectedRoleValue(e.target.value)}
              className="w-full px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
            >
              {ROLE_OPTIONS.map((role) => (
                <option key={role.value} value={role.value}>
                  {role.label}
                </option>
              ))}
            </select>
          </div>

          {selectedRoleValue && user && selectedRoleValue !== user?.role && (
            <div className="p-3 bg-blue-900/20 border border-blue-800 rounded-lg text-sm text-blue-400">
              Роль будет изменена с <strong>{user?.role}</strong> на <strong>{selectedRoleValue}</strong>
            </div>
          )}
        </form>

        {/* Footer */}
        <div className="flex gap-3 p-6 border-t border-gray-700">
          <button
            onClick={onClose}
            disabled={isLoading}
            className="flex-1 px-4 py-2 bg-gray-800 hover:bg-gray-700 text-white rounded-lg transition-colors disabled:opacity-50"
          >
            Отмена
          </button>
          <button
            onClick={handleSubmit}
            disabled={isLoading || !selectedRoleValue || selectedRoleValue === user?.role}
            className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white rounded-lg transition-colors flex items-center justify-center gap-2"
          >
            {isLoading && <Loader2 className="w-4 h-4 animate-spin" />}
            Сохранить
          </button>
        </div>
      </div>
    </div>
  );
}
