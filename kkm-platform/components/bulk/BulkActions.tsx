'use client';

import { useState, useCallback } from 'react';
import { Trash2, Check, X } from 'lucide-react';

export interface BulkActionConfig {
  id: string;
  label: string;
  icon: React.ReactNode;
  action: (ids: string[]) => Promise<void>;
  color?: 'red' | 'green' | 'blue' | 'yellow';
  confirm?: boolean;
}

interface BulkActionsProps {
  selectedIds: string[];
  actions: BulkActionConfig[];
  onClear?: () => void;
  isLoading?: boolean;
}

const colorClasses = {
  red: 'bg-red-600 hover:bg-red-700',
  green: 'bg-green-600 hover:bg-green-700',
  blue: 'bg-blue-600 hover:bg-blue-700',
  yellow: 'bg-yellow-600 hover:bg-yellow-700',
};

export function BulkActions({
  selectedIds,
  actions,
  onClear,
  isLoading,
}: BulkActionsProps) {
  const [pending, setPending] = useState<string | null>(null);
  const [confirmingAction, setConfirmingAction] = useState<string | null>(null);

  const handleAction = useCallback(
    async (actionId: string) => {
      const action = actions.find((a) => a.id === actionId);
      if (!action) return;

      if (action.confirm) {
        setConfirmingAction(actionId);
        return;
      }

      setPending(actionId);
      try {
        await action.action(selectedIds);
      } finally {
        setPending(actionId);
      }
    },
    [actions, selectedIds]
  );

  const confirmAction = useCallback(
    async (actionId: string) => {
      const action = actions.find((a) => a.id === actionId);
      if (!action) return;

      setPending(actionId);
      try {
        await action.action(selectedIds);
      } finally {
        setPending(null);
        setConfirmingAction(null);
      }
    },
    [actions, selectedIds]
  );

  if (selectedIds.length === 0) return null;

  return (
    <>
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <span className="text-sm font-medium text-blue-900">
            Выбрано: <strong>{selectedIds.length}</strong> элементов
          </span>

          <div className="flex gap-2">
            {actions.map((action) => (
              <button
                key={action.id}
                onClick={() => handleAction(action.id)}
                disabled={isLoading || pending === action.id}
                className={`flex items-center gap-2 px-3 py-2 text-white rounded-lg text-sm transition ${
                  colorClasses[action.color || 'blue']
                } disabled:opacity-50`}
              >
                {pending === action.id ? (
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                ) : (
                  action.icon
                )}
                {action.label}
              </button>
            ))}
          </div>
        </div>

        <button
          onClick={onClear}
          className="p-2 text-gray-600 hover:bg-gray-100 rounded-lg transition"
        >
          <X className="w-5 h-5" />
        </button>
      </div>

      {/* Модальное окно подтверждения */}
      {confirmingAction && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-sm">
            <h3 className="text-lg font-semibold mb-4">
              {actions.find((a) => a.id === confirmingAction)?.label}
            </h3>
            <p className="text-gray-600 mb-6">
              Вы уверены? Это действие повлияет на {selectedIds.length} элементов.
            </p>
            <div className="flex gap-4">
              <button
                onClick={() => setConfirmingAction(null)}
                className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition"
              >
                Отмена
              </button>
              <button
                onClick={() => confirmAction(confirmingAction)}
                disabled={pending === confirmingAction}
                className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 transition flex items-center justify-center gap-2"
              >
                {pending === confirmingAction ? (
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                ) : (
                  <Check className="w-4 h-4" />
                )}
                Подтвердить
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
