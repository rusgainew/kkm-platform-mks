"use client";

import { useState } from "react";
import { Plus, Trash2, Edit2, Check, X } from "lucide-react";
import type { CatalogEntry } from "@/types/invoice";
import {
  createEmptyCatalogEntry,
  updateCatalogEntryCalculations,
  formatCurrency,
  validateCatalogEntry,
} from "../lib/invoice-validation";
import { Tooltip } from "@/components/ui/Tooltip";

interface InvoiceCatalogEntriesTableProps {
  entries: CatalogEntry[];
  vatRate: number;
  currency?: string;
  onEntriesChange: (entries: CatalogEntry[]) => void;
  disabled?: boolean;
}

interface EditingEntry {
  index: number;
  entry: CatalogEntry;
}

/**
 * Таблица позиций счета-фактуры (catalog entries)
 * Позволяет добавлять, редактировать, удалять позиции
 * Автоматически рассчитывает суммы с НДС и без НДС
 */
export function InvoiceCatalogEntriesTable({
  entries,
  vatRate,
  currency = "KGS",
  onEntriesChange,
  disabled = false,
}: InvoiceCatalogEntriesTableProps) {
  const [editingEntry, setEditingEntry] = useState<EditingEntry | null>(null);
  const [isAddingNew, setIsAddingNew] = useState(false);
  const [newEntry, setNewEntry] = useState<CatalogEntry>(
    createEmptyCatalogEntry(),
  );
  const [errors, setErrors] = useState<Record<number, string[]>>({});

  /**
   * Добавить новую позицию
   */
  const handleAddEntry = () => {
    const validationErrors = validateCatalogEntry(newEntry);
    if (validationErrors.length > 0) {
      setErrors({ [-1]: validationErrors });
      return;
    }

    const calculatedEntry = updateCatalogEntryCalculations(newEntry, vatRate);
    onEntriesChange([...entries, calculatedEntry]);
    setNewEntry(createEmptyCatalogEntry());
    setIsAddingNew(false);
    setErrors({});
  };

  /**
   * Начать редактирование позиции
   */
  const handleStartEdit = (index: number) => {
    setEditingEntry({
      index,
      entry: { ...entries[index] },
    });
  };

  /**
   * Сохранить изменения позиции
   */
  const handleSaveEdit = () => {
    if (!editingEntry) return;

    const validationErrors = validateCatalogEntry(editingEntry.entry);
    if (validationErrors.length > 0) {
      setErrors({ [editingEntry.index]: validationErrors });
      return;
    }

    const calculatedEntry = updateCatalogEntryCalculations(
      editingEntry.entry,
      vatRate,
    );
    const updatedEntries = [...entries];
    updatedEntries[editingEntry.index] = calculatedEntry;
    onEntriesChange(updatedEntries);
    setEditingEntry(null);
    setErrors({});
  };

  /**
   * Отменить редактирование
   */
  const handleCancelEdit = () => {
    setEditingEntry(null);
    setErrors({});
  };

  /**
   * Удалить позицию
   */
  const handleRemoveEntry = (index: number) => {
    const updatedEntries = entries.filter((_, i) => i !== index);
    onEntriesChange(updatedEntries);
  };

  /**
   * Обновить поле новой позиции
   */
  const handleNewEntryChange = (
    field: keyof CatalogEntry,
    value: string | number,
  ) => {
    setNewEntry((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  /**
   * Обновить поле редактируемой позиции
   */
  const handleEditEntryChange = (
    field: keyof CatalogEntry,
    value: string | number,
  ) => {
    if (!editingEntry) return;
    setEditingEntry({
      ...editingEntry,
      entry: {
        ...editingEntry.entry,
        [field]: value,
      },
    });
  };

  /**
   * Рендер строки позиции
   */
  const renderEntryRow = (entry: CatalogEntry, index: number) => {
    const isEditing = editingEntry?.index === index;
    const entryErrors = errors[index] || [];

    if (isEditing && editingEntry) {
      return (
        <tr key={index} className="bg-yellow-50 dark:bg-yellow-900/10">
          <td className="px-4 py-3">{index + 1}</td>
          <td className="px-4 py-3">
            <input
              type="number"
              value={editingEntry.entry.id}
              onChange={(e) =>
                handleEditEntryChange("id", Number(e.target.value))
              }
              className="w-full px-2 py-1 border rounded dark:bg-gray-800"
              placeholder="ID товара"
            />
          </td>
          <td className="px-4 py-3">
            <input
              type="text"
              value={editingEntry.entry.unitClassificationCode}
              onChange={(e) =>
                handleEditEntryChange("unitClassificationCode", e.target.value)
              }
              className="w-full px-2 py-1 border rounded dark:bg-gray-800"
              placeholder="Код единицы"
            />
          </td>
          <td className="px-4 py-3">
            <input
              type="number"
              step="0.01"
              value={editingEntry.entry.quantity}
              onChange={(e) =>
                handleEditEntryChange("quantity", Number(e.target.value))
              }
              className="w-full px-2 py-1 border rounded dark:bg-gray-800"
              placeholder="Количество"
            />
          </td>
          <td className="px-4 py-3">
            <input
              type="number"
              step="0.01"
              value={editingEntry.entry.price}
              onChange={(e) =>
                handleEditEntryChange("price", Number(e.target.value))
              }
              className="w-full px-2 py-1 border rounded dark:bg-gray-800"
              placeholder="Цена"
            />
          </td>
          <td className="px-4 py-3 text-right">
            {formatCurrency(
              editingEntry.entry.quantity * editingEntry.entry.price,
              currency,
            )}
          </td>
          <td className="px-4 py-3 text-right">
            <div className="flex gap-1 justify-end">
              <Tooltip content="Сохранить" position="top">
                <button
                  type="button"
                  onClick={handleSaveEdit}
                  className="p-1 text-green-600 hover:bg-green-50 dark:hover:bg-green-900/20 rounded"
                >
                  <Check size={18} />
                </button>
              </Tooltip>
              <Tooltip content="Отмена" position="top">
                <button
                  type="button"
                  onClick={handleCancelEdit}
                  className="p-1 text-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 rounded"
                >
                  <X size={18} />
                </button>
              </Tooltip>
            </div>
            {entryErrors.length > 0 && (
              <div className="text-xs text-red-500 mt-1">
                {entryErrors.join(", ")}
              </div>
            )}
          </td>
        </tr>
      );
    }

    return (
      <tr key={index} className="hover:bg-gray-50 dark:hover:bg-gray-800">
        <td className="px-4 py-3 text-center">{index + 1}</td>
        <td className="px-4 py-3">{entry.id}</td>
        <td className="px-4 py-3">{entry.unitClassificationCode}</td>
        <td className="px-4 py-3 text-right">{entry.quantity}</td>
        <td className="px-4 py-3 text-right">
          {formatCurrency(entry.price, currency)}
        </td>
        <td className="px-4 py-3 text-right font-medium">
          {formatCurrency(
            entry.totalAmount || entry.quantity * entry.price,
            currency,
          )}
        </td>
        <td className="px-4 py-3">
          <div className="flex gap-1 justify-end">
            <Tooltip content="Редактировать" position="top">
              <button
                type="button"
                onClick={() => handleStartEdit(index)}
                disabled={disabled}
                className="p-1 text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded disabled:opacity-50"
              >
                <Edit2 size={18} />
              </button>
            </Tooltip>
            <Tooltip content="Удалить" position="top">
              <button
                type="button"
                onClick={() => handleRemoveEntry(index)}
                disabled={disabled}
                className="p-1 text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded disabled:opacity-50"
              >
                <Trash2 size={18} />
              </button>
            </Tooltip>
          </div>
        </td>
      </tr>
    );
  };

  const newEntryErrors = errors[-1] || [];

  return (
    <div className="w-full">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
          Позиции счета-фактуры
        </h3>
        {!isAddingNew && (
          <button
            type="button"
            onClick={() => setIsAddingNew(true)}
            disabled={disabled}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <Plus size={18} />
            Добавить позицию
          </button>
        )}
      </div>

      {/* Таблица */}
      <div className="overflow-x-auto border border-gray-300 dark:border-gray-700 rounded-lg">
        <table className="w-full text-sm text-left">
          <thead className="bg-gray-100 dark:bg-gray-800 border-b border-gray-300 dark:border-gray-700">
            <tr>
              <th className="px-4 py-3 text-center w-12">№</th>
              <th className="px-4 py-3">ID товара</th>
              <th className="px-4 py-3">Ед. изм.</th>
              <th className="px-4 py-3 text-right">Кол-во</th>
              <th className="px-4 py-3 text-right">Цена</th>
              <th className="px-4 py-3 text-right">Сумма</th>
              <th className="px-4 py-3 text-right w-24">Действия</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
            {entries.length === 0 && !isAddingNew ? (
              <tr>
                <td
                  colSpan={7}
                  className="px-4 py-8 text-center text-gray-500 dark:text-gray-400"
                >
                  Нет добавленных позиций. Нажмите &ldquo;Добавить
                  позицию&rdquo; для начала.
                </td>
              </tr>
            ) : (
              entries.map((entry, index) => renderEntryRow(entry, index))
            )}

            {/* Строка добавления новой позиции */}
            {isAddingNew && (
              <tr className="bg-green-50 dark:bg-green-900/10">
                <td className="px-4 py-3 text-center">
                  <Plus size={18} className="inline" />
                </td>
                <td className="px-4 py-3">
                  <input
                    type="number"
                    value={newEntry.id}
                    onChange={(e) =>
                      handleNewEntryChange("id", Number(e.target.value))
                    }
                    className="w-full px-2 py-1 border rounded dark:bg-gray-800"
                    placeholder="ID товара"
                  />
                </td>
                <td className="px-4 py-3">
                  <select
                    value={newEntry.unitClassificationCode}
                    onChange={(e) =>
                      handleNewEntryChange(
                        "unitClassificationCode",
                        e.target.value,
                      )
                    }
                    className="w-full px-2 py-1 border rounded dark:bg-gray-800"
                  >
                    <option value="796">Штука (796)</option>
                    <option value="110">Килограмм (110)</option>
                    <option value="163">Литр (163)</option>
                    <option value="055">Метр (055)</option>
                    <option value="006">Метр квадратный (006)</option>
                    <option value="113">Метр кубический (113)</option>
                    <option value="255">Час (255)</option>
                  </select>
                </td>
                <td className="px-4 py-3">
                  <input
                    type="number"
                    step="0.01"
                    value={newEntry.quantity}
                    onChange={(e) =>
                      handleNewEntryChange("quantity", Number(e.target.value))
                    }
                    className="w-full px-2 py-1 border rounded dark:bg-gray-800"
                    placeholder="Количество"
                  />
                </td>
                <td className="px-4 py-3">
                  <input
                    type="number"
                    step="0.01"
                    value={newEntry.price}
                    onChange={(e) =>
                      handleNewEntryChange("price", Number(e.target.value))
                    }
                    className="w-full px-2 py-1 border rounded dark:bg-gray-800"
                    placeholder="Цена"
                  />
                </td>
                <td className="px-4 py-3 text-right font-medium">
                  {formatCurrency(newEntry.quantity * newEntry.price, currency)}
                </td>
                <td className="px-4 py-3">
                  <div className="flex gap-1 justify-end">
                    <Tooltip content="Добавить" position="top">
                      <button
                        type="button"
                        onClick={handleAddEntry}
                        className="p-1 text-green-600 hover:bg-green-50 dark:hover:bg-green-900/20 rounded"
                      >
                        <Check size={18} />
                      </button>
                    </Tooltip>
                    <Tooltip content="Отмена" position="top">
                      <button
                        type="button"
                        onClick={() => {
                          setIsAddingNew(false);
                          setNewEntry(createEmptyCatalogEntry());
                          setErrors({});
                        }}
                        className="p-1 text-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 rounded"
                      >
                        <X size={18} />
                      </button>
                    </Tooltip>
                  </div>
                  {newEntryErrors.length > 0 && (
                    <div className="text-xs text-red-500 mt-1">
                      {newEntryErrors.join(", ")}
                    </div>
                  )}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Информация о позициях */}
      {entries.length > 0 && (
        <div className="mt-3 text-sm text-gray-600 dark:text-gray-400">
          Всего позиций: <span className="font-medium">{entries.length}</span>
        </div>
      )}
    </div>
  );
}
