"use client";

import React, { useState } from "react";
import { Info } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Switch } from "@/components/ui/Switch";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createBankAccount, updateBankAccount } from "@/lib/api/bank-accounts";
import type {
  BankAccount,
  CreateBankAccountRequest,
  UpdateBankAccountRequest,
} from "@/types/entities";

interface BankAccountFormProps {
  initialData?: BankAccount | null;
  ownerId?: string;
  onSuccess?: () => void;
  onCancel?: () => void;
}

interface FormData {
  account_number: string;
  bank_name: string;
  bank_code: string;
  currency: string;
  owner_id: string;
  is_active: boolean;
}

// Валидация номера счета (20 цифр для КР)
function isValidAccountNumber(number: string): boolean {
  return /^\d{20}$/.test(number);
}

// Валидация БИК/bank_code (обычно 9 цифр для РФ, но может варьироваться)
function isValidBankCode(code: string): boolean {
  return /^\d{6,9}$/.test(code);
}

export default function BankAccountForm({
  initialData,
  ownerId,
  onSuccess,
  onCancel,
}: BankAccountFormProps) {
  const queryClient = useQueryClient();
  const isEditMode = !!initialData?.id;

  const [formData, setFormData] = useState<FormData>({
    account_number: initialData?.account_number || "",
    bank_name: initialData?.bank_name || "",
    bank_code: initialData?.bank_code || "",
    currency: initialData?.currency || "KGS",
    owner_id: initialData?.owner_id || ownerId || "",
    is_active: initialData?.is_active ?? true,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [touched, setTouched] = useState<Record<string, boolean>>({});

  // Валидация формы
  const validateForm = React.useCallback((): Record<string, string> => {
    const newErrors: Record<string, string> = {};

    if (!formData.account_number.trim()) {
      newErrors.account_number = "Номер счета обязателен";
    } else if (!isValidAccountNumber(formData.account_number)) {
      newErrors.account_number = "Номер счета должен содержать ровно 20 цифр";
    }

    if (!formData.bank_name.trim()) {
      newErrors.bank_name = "Название банка обязательно";
    } else if (formData.bank_name.length < 3) {
      newErrors.bank_name = "Название банка должно содержать минимум 3 символа";
    }

    if (!formData.bank_code.trim()) {
      newErrors.bank_code = "БИК банка обязателен";
    } else if (!isValidBankCode(formData.bank_code)) {
      newErrors.bank_code = "БИК должен содержать от 6 до 9 цифр";
    }

    if (!formData.currency) {
      newErrors.currency = "Валюта обязательна";
    }

    if (!formData.owner_id) {
      newErrors.owner_id = "Необходимо указать владельца счета";
    }

    return newErrors;
  }, [formData]);

  // Мутация создания
  const createMutation = useMutation({
    mutationFn: (data: CreateBankAccountRequest) => createBankAccount(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bank-accounts"] });
      onSuccess?.();
    },
  });

  // Мутация обновления
  const updateMutation = useMutation({
    mutationFn: (data: UpdateBankAccountRequest) =>
      updateBankAccount(initialData!.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bank-accounts"] });
      queryClient.invalidateQueries({
        queryKey: ["bank-account", initialData!.id],
      });
      onSuccess?.();
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Отметить все поля как touched
    const allFields = Object.keys(formData).reduce(
      (acc, key) => ({ ...acc, [key]: true }),
      {},
    );
    setTouched(allFields);

    const validationErrors = validateForm();
    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }

    try {
      if (isEditMode) {
        // Обновление - отправляем только измененные поля
        const updateData: UpdateBankAccountRequest = {};
        if (formData.account_number !== initialData.account_number) {
          updateData.account_number = formData.account_number;
        }
        if (formData.bank_name !== initialData.bank_name) {
          updateData.bank_name = formData.bank_name;
        }
        if (formData.bank_code !== initialData.bank_code) {
          updateData.bank_code = formData.bank_code;
        }
        if (formData.currency !== initialData.currency) {
          updateData.currency = formData.currency;
        }
        if (formData.is_active !== initialData.is_active) {
          updateData.is_active = formData.is_active;
        }

        await updateMutation.mutateAsync(updateData);
      } else {
        // Создание
        const createData: CreateBankAccountRequest = {
          account_number: formData.account_number,
          bank_name: formData.bank_name,
          bank_code: formData.bank_code,
          currency: formData.currency,
          owner_id: formData.owner_id,
        };
        await createMutation.mutateAsync(createData);
      }
    } catch (error) {
      console.error("Ошибка отправки формы:", error);
      setErrors({
        submit:
          error instanceof Error ? error.message : "Ошибка сохранения данных",
      });
    }
  };

  const handleBlur = (field: keyof FormData) => {
    setTouched((prev) => ({ ...prev, [field]: true }));
    // Validate only the touched field
    const validationErrors = validateForm();
    setErrors((prev) => ({
      ...prev,
      [field]: validationErrors[field] || "",
    }));
  };

  const isLoading = createMutation.isPending || updateMutation.isPending;
  const submitError = createMutation.error || updateMutation.error;

  return (
    <form
      onSubmit={handleSubmit}
      className="bg-gray-900 rounded-lg border border-gray-800 p-6 space-y-6"
    >
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-white">
          {isEditMode
            ? "Редактировать банковский счет"
            : "Добавить банковский счет"}
        </h2>
      </div>

      {/* Error messages */}
      {(errors.submit || submitError) && (
        <div className="p-4 bg-red-900/20 border border-red-800 text-red-300 rounded-lg">
          {errors.submit ||
            (submitError instanceof Error
              ? submitError.message
              : "Ошибка сохранения")}
        </div>
      )}

      {/* Основные реквизиты */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-gray-300">
          Банковские реквизиты
        </h3>

        {/* Номер счета */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Номер счета <span className="text-red-400">*</span>
          </label>
          <input
            type="text"
            value={formData.account_number}
            onChange={(e) => {
              const cleaned = e.target.value.replace(/\D/g, "").slice(0, 20);
              setFormData({ ...formData, account_number: cleaned });
            }}
            onBlur={() => handleBlur("account_number")}
            className={`w-full px-4 py-2 rounded-lg bg-gray-800 border ${
              touched.account_number && errors.account_number
                ? "border-red-500"
                : "border-gray-700"
            } text-white focus:border-green-500 outline-none font-mono`}
            placeholder="12345678901234567890"
            maxLength={20}
          />
          {touched.account_number && errors.account_number && (
            <p className="text-red-400 text-sm mt-1">{errors.account_number}</p>
          )}
          <p className="text-gray-400 text-xs mt-1">
            Требуется 20 цифр. Текущая длина: {formData.account_number.length}
          </p>
        </div>

        {/* Название банка */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Название банка <span className="text-red-400">*</span>
          </label>
          <input
            type="text"
            value={formData.bank_name}
            onChange={(e) =>
              setFormData({ ...formData, bank_name: e.target.value })
            }
            onBlur={() => handleBlur("bank_name")}
            className={`w-full px-4 py-2 rounded-lg bg-gray-800 border ${
              touched.bank_name && errors.bank_name
                ? "border-red-500"
                : "border-gray-700"
            } text-white focus:border-green-500 outline-none`}
            placeholder="ПАО Сбербанк"
          />
          {touched.bank_name && errors.bank_name && (
            <p className="text-red-400 text-sm mt-1">{errors.bank_name}</p>
          )}
        </div>

        {/* БИК банка */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            БИК банка <span className="text-red-400">*</span>
          </label>
          <input
            type="text"
            value={formData.bank_code}
            onChange={(e) => {
              const cleaned = e.target.value.replace(/\D/g, "").slice(0, 9);
              setFormData({ ...formData, bank_code: cleaned });
            }}
            onBlur={() => handleBlur("bank_code")}
            className={`w-full px-4 py-2 rounded-lg bg-gray-800 border ${
              touched.bank_code && errors.bank_code
                ? "border-red-500"
                : "border-gray-700"
            } text-white focus:border-green-500 outline-none font-mono`}
            placeholder="044525225"
            maxLength={9}
          />
          {touched.bank_code && errors.bank_code && (
            <p className="text-red-400 text-sm mt-1">{errors.bank_code}</p>
          )}
          <p className="text-gray-400 text-xs mt-1">
            Обычно 9 цифр для РФ, 6-9 цифр для других стран
          </p>
        </div>

        {/* Валюта */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Валюта счета <span className="text-red-400">*</span>
          </label>
          <select
            value={formData.currency}
            onChange={(e) =>
              setFormData({ ...formData, currency: e.target.value })
            }
            onBlur={() => handleBlur("currency")}
            className={`w-full px-4 py-2 rounded-lg bg-gray-800 border ${
              touched.currency && errors.currency
                ? "border-red-500"
                : "border-gray-700"
            } text-white focus:border-green-500 outline-none`}
          >
            <option value="KGS">🇰🇬 KGS - Кыргызский сом (с)</option>
            <option value="USD">🇺🇸 USD - Доллар США ($)</option>
            <option value="RUB">🇷🇺 RUB - Российский рубль (₽)</option>
            <option value="EUR">🇪🇺 EUR - Евро (€)</option>
            <option value="CNY">🇨🇳 CNY - Китайский юань (¥)</option>
          </select>
          {touched.currency && errors.currency && (
            <p className="text-red-400 text-sm mt-1">{errors.currency}</p>
          )}
        </div>
      </div>

      {/* Статус счета */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-gray-300">Статус</h3>

        <div className="p-4 bg-gray-800 rounded-lg">
          <Switch
            checked={formData.is_active}
            onChange={(checked) =>
              setFormData({ ...formData, is_active: checked })
            }
            label="Активный счет"
            labelPosition="right"
            size="md"
          />
          <p className="text-gray-400 text-sm mt-2 ml-14">
            Счет доступен для проведения операций
          </p>
        </div>
      </div>

      {/* Info box */}
      <div className="p-4 bg-blue-900/20 border border-blue-800 rounded-lg flex gap-3">
        <Info className="w-5 h-5 text-blue-300 shrink-0 mt-0.5" />
        <div className="text-blue-300 text-sm space-y-1">
          <p className="font-medium">Важная информация:</p>
          <ul className="list-disc list-inside space-y-1 ml-2">
            <li>Номер счета должен содержать ровно 20 цифр (стандарт КР)</li>
            <li>БИК банка обычно содержит 9 цифр для российских банков</li>
            <li>Все реквизиты должны быть проверены перед сохранением</li>
            <li>Неактивные счета не будут доступны для операций</li>
          </ul>
        </div>
      </div>

      {/* Buttons */}
      <div className="flex gap-3">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            disabled={isLoading}
            className="px-6 py-2 bg-gray-700 text-white rounded-lg font-medium hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Отмена
          </button>
        )}
        <Button
          type="submit"
          disabled={isLoading || Object.keys(errors).length > 0}
          loading={isLoading}
          variant="success"
          className="flex-1"
        >
          {isEditMode ? "Обновить счет" : "Создать счет"}
        </Button>
      </div>
    </form>
  );
}
