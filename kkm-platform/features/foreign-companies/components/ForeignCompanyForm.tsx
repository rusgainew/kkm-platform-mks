"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import {
  createForeignCompany,
  updateForeignCompany,
  deleteForeignCompany,
  type ForeignCompany,
  type CreateForeignCompanyRequest,
  type UpdateForeignCompanyRequest,
} from "@/lib/api/foreign-companies";
import { POPULAR_COUNTRIES } from "@/types/entities";

interface ForeignCompanyFormProps {
  company?: ForeignCompany;
  onSuccess?: () => void;
  onCancel?: () => void;
}

interface FormData {
  name: string;
  tax_id: string;
  country: string;
  address: string;
  contact_email: string;
  contact_phone: string;
  currency: string;
}

interface FormErrors {
  name?: string;
  tax_id?: string;
  country?: string;
  address?: string;
  contact_email?: string;
  contact_phone?: string;
  currency?: string;
}

/**
 * Проверка валидности email
 */
function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

/**
 * Проверка валидности телефона (международный формат)
 */
function isValidPhone(phone: string): boolean {
  const phoneRegex = /^\+?[\d\s\-()]+$/;
  return phoneRegex.test(phone) && phone.replace(/\D/g, "").length >= 10;
}

/**
 * Проверка валидности кода страны (ISO 3166-1 alpha-2)
 */
function isValidCountryCode(code: string): boolean {
  return /^[A-Z]{2}$/.test(code);
}

/**
 * Проверка валидности кода валюты (ISO 4217)
 */
function isValidCurrencyCode(code: string): boolean {
  return /^[A-Z]{3}$/.test(code);
}

/**
 * Валидация PIN/Tax ID (базовая проверка длины и формата)
 */
function isValidTaxId(taxId: string): boolean {
  // Минимум 5 символов, может содержать буквы, цифры, дефисы
  return /^[A-Z0-9\-]{5,}$/.test(taxId);
}

/**
 * Форма для создания/редактирования иностранной компании
 */
export function ForeignCompanyForm({
  company,
  onSuccess,
  onCancel,
}: ForeignCompanyFormProps) {
  const queryClient = useQueryClient();
  const isEditing = !!company;

  const [formData, setFormData] = useState<FormData>({
    name: company?.name || "",
    tax_id: company?.tax_id || "",
    country: company?.country || "",
    address: company?.address || "",
    contact_email: company?.contact_email || "",
    contact_phone: company?.contact_phone || "",
    currency: company?.currency || "USD",
  });

  const [errors, setErrors] = useState<FormErrors>({});
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  /**
   * Валидация формы
   */
  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};

    // Название компании (обязательное)
    if (!formData.name.trim()) {
      newErrors.name = "Название компании обязательно";
    } else if (formData.name.length < 2) {
      newErrors.name = "Название должно содержать минимум 2 символа";
    } else if (formData.name.length > 255) {
      newErrors.name = "Название не может превышать 255 символов";
    }

    // Tax ID / PIN (обязательное)
    if (!formData.tax_id.trim()) {
      newErrors.tax_id = "Налоговый номер (PIN/TIN) обязателен";
    } else if (!isValidTaxId(formData.tax_id.toUpperCase())) {
      newErrors.tax_id =
        "Неверный формат налогового номера (минимум 5 символов, только буквы, цифры, дефисы)";
    }

    // Код страны (обязательное, ISO 3166-1 alpha-2)
    if (!formData.country.trim()) {
      newErrors.country = "Код страны обязателен";
    } else if (!isValidCountryCode(formData.country.toUpperCase())) {
      newErrors.country = "Неверный формат кода страны (ISO 3166-1 alpha-2, например: US, CN, RU)";
    }

    // Адрес (обязательное)
    if (!formData.address.trim()) {
      newErrors.address = "Адрес обязателен";
    } else if (formData.address.length < 5) {
      newErrors.address = "Адрес должен содержать минимум 5 символов";
    } else if (formData.address.length > 500) {
      newErrors.address = "Адрес не может превышать 500 символов";
    }

    // Email (обязательное)
    if (!formData.contact_email.trim()) {
      newErrors.contact_email = "Email обязателен";
    } else if (!isValidEmail(formData.contact_email)) {
      newErrors.contact_email = "Неверный формат email";
    }

    // Телефон (обязательное)
    if (!formData.contact_phone.trim()) {
      newErrors.contact_phone = "Телефон обязателен";
    } else if (!isValidPhone(formData.contact_phone)) {
      newErrors.contact_phone =
        "Неверный формат телефона (международный формат, минимум 10 цифр)";
    }

    // Валюта (обязательное, ISO 4217)
    if (!formData.currency.trim()) {
      newErrors.currency = "Код валюты обязателен";
    } else if (!isValidCurrencyCode(formData.currency.toUpperCase())) {
      newErrors.currency = "Неверный формат кода валюты (ISO 4217, например: USD, EUR, CNY)";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  /**
   * Мутация для создания компании
   */
  const createMutation = useMutation({
    mutationFn: async (data: CreateForeignCompanyRequest) => {
      const response = await createForeignCompany(data);
      return response.data;
    },
    onSuccess: () => {
      toast.success("Иностранная компания успешно создана");
      queryClient.invalidateQueries({ queryKey: ["foreign-companies"] });
      onSuccess?.();
    },
    onError: (error: Error) => {
      toast.error(`Ошибка при создании: ${error.message}`);
    },
  });

  /**
   * Мутация для обновления компании
   */
  const updateMutation = useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: string;
      data: UpdateForeignCompanyRequest;
    }) => {
      const response = await updateForeignCompany(id, data);
      return response.data;
    },
    onSuccess: () => {
      toast.success("Иностранная компания успешно обновлена");
      queryClient.invalidateQueries({ queryKey: ["foreign-companies"] });
      queryClient.invalidateQueries({
        queryKey: ["foreign-company", company?.id],
      });
      onSuccess?.();
    },
    onError: (error: Error) => {
      toast.error(`Ошибка при обновлении: ${error.message}`);
    },
  });

  /**
   * Мутация для удаления компании
   */
  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await deleteForeignCompany(id);
    },
    onSuccess: () => {
      toast.success("Иностранная компания успешно удалена");
      queryClient.invalidateQueries({ queryKey: ["foreign-companies"] });
      onSuccess?.();
    },
    onError: (error: Error) => {
      toast.error(`Ошибка при удалении: ${error.message}`);
    },
  });

  /**
   * Обработчик отправки формы
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      toast.error("Пожалуйста, исправьте ошибки в форме");
      return;
    }

    // Нормализация данных
    const normalizedData = {
      name: formData.name.trim(),
      tax_id: formData.tax_id.trim().toUpperCase(),
      country: formData.country.trim().toUpperCase(),
      address: formData.address.trim(),
      contact_email: formData.contact_email.trim().toLowerCase(),
      contact_phone: formData.contact_phone.trim(),
      currency: formData.currency.trim().toUpperCase(),
    };

    if (isEditing && company) {
      updateMutation.mutate({
        id: company.id,
        data: normalizedData,
      });
    } else {
      createMutation.mutate(normalizedData);
    }
  };

  /**
   * Обработчик удаления
   */
  const handleDelete = () => {
    if (company) {
      deleteMutation.mutate(company.id);
    }
  };

  /**
   * Обработчик изменения поля
   */
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    // Очистить ошибку при изменении поля
    if (errors[name as keyof FormErrors]) {
      setErrors((prev) => ({ ...prev, [name]: undefined }));
    }
  };

  const isLoading =
    createMutation.isPending ||
    updateMutation.isPending ||
    deleteMutation.isPending;

  return (
    <div className="w-full max-w-2xl mx-auto">
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Заголовок */}
        <div className="border-b pb-4">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
            {isEditing
              ? "Редактировать иностранную компанию"
              : "Добавить иностранную компанию"}
          </h2>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            {isEditing
              ? "Обновите информацию об иностранной компании"
              : "Заполните информацию о новой иностранной компании"}
          </p>
        </div>

        {/* Основная информация */}
        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
            Основная информация
          </h3>

          {/* Название компании */}
          <div>
            <label
              htmlFor="name"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
              Название компании <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleChange}
              disabled={isLoading}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                errors.name
                  ? "border-red-500"
                  : "border-gray-300 dark:border-gray-600"
              } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              placeholder="Acme Corporation Ltd."
            />
            {errors.name && (
              <p className="mt-1 text-sm text-red-500">{errors.name}</p>
            )}
          </div>

          {/* Tax ID / PIN */}
          <div>
            <label
              htmlFor="tax_id"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
              Налоговый номер (PIN/TIN) <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              id="tax_id"
              name="tax_id"
              value={formData.tax_id}
              onChange={handleChange}
              disabled={isLoading}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 uppercase ${
                errors.tax_id
                  ? "border-red-500"
                  : "border-gray-300 dark:border-gray-600"
              } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              placeholder="12-3456789"
            />
            {errors.tax_id && (
              <p className="mt-1 text-sm text-red-500">{errors.tax_id}</p>
            )}
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Иностранный налоговый идентификационный номер (минимум 5 символов)
            </p>
          </div>

          {/* Код страны */}
          <div>
            <label
              htmlFor="country"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
              Страна <span className="text-red-500">*</span>
            </label>
            <select
              id="country"
              name="country"
              value={formData.country}
              onChange={handleChange}
              disabled={isLoading}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                errors.country
                  ? "border-red-500"
                  : "border-gray-300 dark:border-gray-600"
              } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
            >
              <option value="">Выберите страну</option>
              <optgroup label="Популярные страны">
                {POPULAR_COUNTRIES.map((country) => (
                  <option key={country.code} value={country.code}>
                    {country.name} ({country.code})
                  </option>
                ))}
              </optgroup>
            </select>
            {errors.country && (
              <p className="mt-1 text-sm text-red-500">{errors.country}</p>
            )}
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              ISO 3166-1 alpha-2 код страны (2 буквы, например: US, CN, RU)
            </p>
          </div>

          {/* Адрес */}
          <div>
            <label
              htmlFor="address"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
              Адрес <span className="text-red-500">*</span>
            </label>
            <textarea
              id="address"
              name="address"
              value={formData.address}
              onChange={handleChange}
              disabled={isLoading}
              rows={3}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                errors.address
                  ? "border-red-500"
                  : "border-gray-300 dark:border-gray-600"
              } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              placeholder="123 Main Street, Suite 100, New York, NY 10001, USA"
            />
            {errors.address && (
              <p className="mt-1 text-sm text-red-500">{errors.address}</p>
            )}
          </div>
        </div>

        {/* Контактная информация */}
        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
            Контактная информация
          </h3>

          {/* Email */}
          <div>
            <label
              htmlFor="contact_email"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
              Email <span className="text-red-500">*</span>
            </label>
            <input
              type="email"
              id="contact_email"
              name="contact_email"
              value={formData.contact_email}
              onChange={handleChange}
              disabled={isLoading}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                errors.contact_email
                  ? "border-red-500"
                  : "border-gray-300 dark:border-gray-600"
              } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              placeholder="contact@example.com"
            />
            {errors.contact_email && (
              <p className="mt-1 text-sm text-red-500">{errors.contact_email}</p>
            )}
          </div>

          {/* Телефон */}
          <div>
            <label
              htmlFor="contact_phone"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
              Телефон <span className="text-red-500">*</span>
            </label>
            <input
              type="tel"
              id="contact_phone"
              name="contact_phone"
              value={formData.contact_phone}
              onChange={handleChange}
              disabled={isLoading}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                errors.contact_phone
                  ? "border-red-500"
                  : "border-gray-300 dark:border-gray-600"
              } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              placeholder="+1 (555) 123-4567"
            />
            {errors.contact_phone && (
              <p className="mt-1 text-sm text-red-500">{errors.contact_phone}</p>
            )}
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              Международный формат телефона (минимум 10 цифр)
            </p>
          </div>

          {/* Валюта */}
          <div>
            <label
              htmlFor="currency"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
              Валюта <span className="text-red-500">*</span>
            </label>
            <select
              id="currency"
              name="currency"
              value={formData.currency}
              onChange={handleChange}
              disabled={isLoading}
              className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                errors.currency
                  ? "border-red-500"
                  : "border-gray-300 dark:border-gray-600"
              } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
            >
              <option value="USD">USD - Доллар США</option>
              <option value="EUR">EUR - Евро</option>
              <option value="CNY">CNY - Китайский юань</option>
              <option value="RUB">RUB - Российский рубль</option>
              <option value="KZT">KZT - Казахстанский тенге</option>
              <option value="UZS">UZS - Узбекский сум</option>
              <option value="TJS">TJS - Таджикский сомони</option>
              <option value="TRY">TRY - Турецкая лира</option>
              <option value="GBP">GBP - Британский фунт</option>
            </select>
            {errors.currency && (
              <p className="mt-1 text-sm text-red-500">{errors.currency}</p>
            )}
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              ISO 4217 код валюты (3 буквы)
            </p>
          </div>
        </div>

        {/* Кнопки действий */}
        <div className="flex items-center gap-3 pt-4 border-t">
          <button
            type="submit"
            disabled={isLoading}
            className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:ring-4 focus:ring-blue-300 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {isLoading
              ? "Сохранение..."
              : isEditing
              ? "Сохранить изменения"
              : "Создать компанию"}
          </button>

          {onCancel && (
            <button
              type="button"
              onClick={onCancel}
              disabled={isLoading}
              className="px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Отмена
            </button>
          )}

          {isEditing && company && (
            <button
              type="button"
              onClick={() => setShowDeleteConfirm(true)}
              disabled={isLoading}
              className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 focus:ring-4 focus:ring-red-300 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Удалить
            </button>
          )}
        </div>
      </form>

      {/* Модальное окно подтверждения удаления */}
      {showDeleteConfirm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2">
              Подтвердите удаление
            </h3>
            <p className="text-gray-600 dark:text-gray-400 mb-4">
              Вы уверены, что хотите удалить компанию &ldquo;{company?.name}&rdquo;? Это
              действие нельзя отменить.
            </p>
            <div className="flex gap-3">
              <button
                onClick={handleDelete}
                disabled={deleteMutation.isPending}
                className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 focus:ring-4 focus:ring-red-300 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {deleteMutation.isPending ? "Удаление..." : "Удалить"}
              </button>
              <button
                onClick={() => setShowDeleteConfirm(false)}
                disabled={deleteMutation.isPending}
                className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Отмена
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
