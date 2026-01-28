"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { Save, Loader2, AlertCircle } from "lucide-react";
import type { CreateInvoiceRequest, CatalogEntry, ESFInvoice } from "@/types/invoice";
import {
  ESFOperationType,
  ESFDeliveryType,
  ESFPaymentType,
} from "@/types/enums";
import { invoiceAPI } from "@/lib/api/invoice";
import { InvoiceCatalogEntriesTable } from "./InvoiceCatalogEntriesTable";
import { InvoiceTotals } from "./InvoiceTotals";
import { useInvoiceCalculations } from "../hooks/useInvoiceCalculations";
import {
  validateInvoiceForm,
  getOperationTypeLabel,
  getDeliveryTypeLabel,
  getPaymentTypeLabel,
} from "../lib/invoice-validation";

interface InvoiceFormProps {
  invoice?: ESFInvoice;
  onSuccess?: () => void;
  onCancel?: () => void;
}

interface FormData extends Omit<CreateInvoiceRequest, "catalogEntries"> {
  catalogEntries: CatalogEntry[];
}

const CURRENCY_OPTIONS = [
  { code: "KGS", name: "Киргизский сом (KGS)" },
  { code: "RUB", name: "Российский рубль (RUB)" },
  { code: "USD", name: "Доллар США (USD)" },
  { code: "EUR", name: "Евро (EUR)" },
  { code: "CNY", name: "Китайский юань (CNY)" },
];

const VAT_RATE_OPTIONS = [
  { code: "12", name: "12%" },
  { code: "0", name: "0%" },
  { code: "НО", name: "Не облагается" },
];

/**
 * Полноценная форма создания/редактирования счета-фактуры (ESF)
 * Поддерживает все поля согласно ESF API спецификации
 */
export function InvoiceFormESF({
  invoice,
  onSuccess,
  onCancel,
}: InvoiceFormProps) {
  const queryClient = useQueryClient();
  const isEditing = !!invoice;

  // Начальные данные формы
  const [formData, setFormData] = useState<FormData>({
    operationTypeCode: invoice?.operationType || ESFOperationType.SALES,
    invoiceNumber: invoice?.invoiceNumber || "",
    deliveryDate: invoice?.deliveryDate || new Date().toISOString().split("T")[0],
    invoiceDate: invoice?.invoiceDate || "",
    contractorTin: (invoice?.contractor && "tin" in invoice.contractor) ? String(invoice.contractor.tin) : "",
    supplierBankAccount: invoice?.legalPersonBankAccount || "",
    contractorBankAccount: "",
    deliveryTypeCode: invoice?.deliveryCode || ESFDeliveryType.DIRECT,
    paymentCode: invoice?.paymentCode || ESFPaymentType.BANK_TRANSFER,
    currencyCode: invoice?.currency?.code || "KGS",
    countryCode: "KG",
    taxRateVATCode: invoice?.vatTaxType?.code || "12",
    isResident: invoice?.isResident ?? true,
    isPriceWithoutTaxes: false,
    currencyRate: 1,
    comment: invoice?.comment || "",
    catalogEntries: [],
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  // Вычисляем НДС ставку из кода
  const vatRate = formData.taxRateVATCode === "12" ? 12 : 0;

  // Используем хук для автоматических расчетов
  const { totals } = useInvoiceCalculations({
    catalogEntries: formData.catalogEntries,
    vatRate,
  });

  /**
   * Мутация для создания счета-фактуры
   */
  const createMutation = useMutation({
    mutationFn: async (data: CreateInvoiceRequest) => {
      const response = await invoiceAPI.createInvoice(data);
      return response;
    },
    onSuccess: () => {
      toast.success("Счет-фактура успешно создана");
      queryClient.invalidateQueries({ queryKey: ["invoices"] });
      onSuccess?.();
    },
    onError: (error: Error) => {
      toast.error(`Ошибка при создании: ${error.message}`);
    },
  });

  /**
   * Мутация для обновления счета-фактуры
   */
  const updateMutation = useMutation({
    mutationFn: async (data: CreateInvoiceRequest & { documentUuid: string }) => {
      const response = await invoiceAPI.updateInvoice(data);
      return response;
    },
    onSuccess: () => {
      toast.success("Счет-фактура успешно обновлена");
      queryClient.invalidateQueries({ queryKey: ["invoices"] });
      queryClient.invalidateQueries({ queryKey: ["invoice", invoice?.documentUuid] });
      onSuccess?.();
    },
    onError: (error: Error) => {
      toast.error(`Ошибка при обновлении: ${error.message}`);
    },
  });

  /**
   * Мутация для отзыва (удаления) счета-фактуры
   */
  const deleteMutation = useMutation({
    mutationFn: async (documentUuid: string) => {
      await invoiceAPI.revokeInvoice(documentUuid, "Удаление счета-фактуры");
    },
    onSuccess: () => {
      toast.success("Счет-фактура успешно удалена");
      queryClient.invalidateQueries({ queryKey: ["invoices"] });
      onSuccess?.();
    },
    onError: (error: Error) => {
      toast.error(`Ошибка при удалении: ${error.message}`);
    },
  });

  /**
   * Обработчик изменения поля
   */
  const handleFieldChange = (field: keyof FormData, value: unknown) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Очистить ошибку для этого поля
    if (errors[field]) {
      setErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  /**
   * Обработчик изменения catalogEntries
   */
  const handleCatalogEntriesChange = (entries: CatalogEntry[]) => {
    setFormData((prev) => ({ ...prev, catalogEntries: entries }));
    // Очистить ошибку catalogEntries
    if (errors.catalogEntries) {
      setErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors.catalogEntries;
        return newErrors;
      });
    }
  };

  /**
   * Валидация и отправка формы
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Валидация
    const validationErrors = validateInvoiceForm(formData);
    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      toast.error("Пожалуйста, исправьте ошибки в форме");
      return;
    }

    // Подготовка данных для отправки
    const requestData: CreateInvoiceRequest = {
      ...formData,
      catalogEntries: formData.catalogEntries,
    };

    if (isEditing && invoice) {
      updateMutation.mutate({
        ...requestData,
        documentUuid: invoice.documentUuid,
      });
    } else {
      createMutation.mutate(requestData);
    }
  };

  /**
   * Обработчик удаления
   */
  const handleDelete = () => {
    if (invoice) {
      deleteMutation.mutate(invoice.documentUuid);
    }
  };

  const isLoading =
    createMutation.isPending ||
    updateMutation.isPending ||
    deleteMutation.isPending;

  return (
    <div className="w-full max-w-6xl mx-auto">
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Заголовок */}
        <div className="border-b pb-4">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
            {isEditing
              ? `Редактировать счет-фактуру ${invoice?.invoiceNumber || ""}`
              : "Создать новый счет-фактуру (ЭСФ)"}
          </h2>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
            {isEditing
              ? "Обновите информацию о счете-фактуре"
              : "Заполните все обязательные поля для создания электронного счета-фактуры"}
          </p>
        </div>

        {/* Секция 1: Основная информация */}
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Основная информация
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {/* Тип операции */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Тип операции <span className="text-red-500">*</span>
              </label>
              <select
                value={formData.operationTypeCode}
                onChange={(e) => handleFieldChange("operationTypeCode", e.target.value)}
                disabled={isLoading}
                className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                  errors.operationTypeCode
                    ? "border-red-500"
                    : "border-gray-300 dark:border-gray-600"
                } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              >
                <option value={ESFOperationType.SALES}>
                  {getOperationTypeLabel(ESFOperationType.SALES)}
                </option>
                <option value={ESFOperationType.PURCHASE_SERVICES}>
                  {getOperationTypeLabel(ESFOperationType.PURCHASE_SERVICES)}
                </option>
                <option value={ESFOperationType.OTHER}>
                  {getOperationTypeLabel(ESFOperationType.OTHER)}
                </option>
              </select>
              {errors.operationTypeCode && (
                <p className="mt-1 text-sm text-red-500">{errors.operationTypeCode}</p>
              )}
            </div>

            {/* Номер счета */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Номер счета-фактуры
              </label>
              <input
                type="text"
                value={formData.invoiceNumber}
                onChange={(e) => handleFieldChange("invoiceNumber", e.target.value)}
                disabled={isLoading}
                className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                  errors.invoiceNumber
                    ? "border-red-500"
                    : "border-gray-300 dark:border-gray-600"
                } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
                placeholder="СФ-001 (опционально)"
              />
              {errors.invoiceNumber && (
                <p className="mt-1 text-sm text-red-500">{errors.invoiceNumber}</p>
              )}
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                Если не указан, будет сгенерирован автоматически
              </p>
            </div>

            {/* Дата доставки */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Дата доставки <span className="text-red-500">*</span>
              </label>
              <input
                type="date"
                value={formData.deliveryDate}
                onChange={(e) => handleFieldChange("deliveryDate", e.target.value)}
                disabled={isLoading}
                max={new Date().toISOString().split("T")[0]}
                className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                  errors.deliveryDate
                    ? "border-red-500"
                    : "border-gray-300 dark:border-gray-600"
                } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              />
              {errors.deliveryDate && (
                <p className="mt-1 text-sm text-red-500">{errors.deliveryDate}</p>
              )}
            </div>

            {/* Дата счета-фактуры */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Дата счета-фактуры
              </label>
              <input
                type="date"
                value={formData.invoiceDate}
                onChange={(e) => handleFieldChange("invoiceDate", e.target.value)}
                disabled={isLoading}
                className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                  errors.invoiceDate
                    ? "border-red-500"
                    : "border-gray-300 dark:border-gray-600"
                } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              />
              {errors.invoiceDate && (
                <p className="mt-1 text-sm text-red-500">{errors.invoiceDate}</p>
              )}
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                Если не указана, используется текущая дата
              </p>
            </div>

            {/* Валюта */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Валюта <span className="text-red-500">*</span>
              </label>
              <select
                value={formData.currencyCode}
                onChange={(e) => handleFieldChange("currencyCode", e.target.value)}
                disabled={isLoading}
                className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                  errors.currencyCode
                    ? "border-red-500"
                    : "border-gray-300 dark:border-gray-600"
                } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              >
                {CURRENCY_OPTIONS.map((curr) => (
                  <option key={curr.code} value={curr.code}>
                    {curr.name}
                  </option>
                ))}
              </select>
              {errors.currencyCode && (
                <p className="mt-1 text-sm text-red-500">{errors.currencyCode}</p>
              )}
            </div>

            {/* Ставка НДС */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Ставка НДС <span className="text-red-500">*</span>
              </label>
              <select
                value={formData.taxRateVATCode}
                onChange={(e) => handleFieldChange("taxRateVATCode", e.target.value)}
                disabled={isLoading}
                className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                  errors.taxRateVATCode
                    ? "border-red-500"
                    : "border-gray-300 dark:border-gray-600"
                } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              >
                {VAT_RATE_OPTIONS.map((rate) => (
                  <option key={rate.code} value={rate.code}>
                    {rate.name}
                  </option>
                ))}
              </select>
              {errors.taxRateVATCode && (
                <p className="mt-1 text-sm text-red-500">{errors.taxRateVATCode}</p>
              )}
            </div>
          </div>
        </div>

        {/* Секция 2: Контрагенты */}
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Контрагенты
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* ИНН покупателя */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                ИНН покупателя <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                value={formData.contractorTin}
                onChange={(e) => handleFieldChange("contractorTin", e.target.value)}
                disabled={isLoading}
                className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                  errors.contractorTin
                    ? "border-red-500"
                    : "border-gray-300 dark:border-gray-600"
                } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
                placeholder="01206200110100"
              />
              {errors.contractorTin && (
                <p className="mt-1 text-sm text-red-500">{errors.contractorTin}</p>
              )}
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                14 цифр для резидентов КР
              </p>
            </div>

            {/* Резидент */}
            <div className="flex items-center pt-6">
              <input
                type="checkbox"
                id="isResident"
                checked={formData.isResident}
                onChange={(e) => handleFieldChange("isResident", e.target.checked)}
                disabled={isLoading}
                className="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
              />
              <label
                htmlFor="isResident"
                className="ml-2 text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                Покупатель - резидент Кыргызстана
              </label>
            </div>

            {/* Банковский счет поставщика */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Банковский счет поставщика
              </label>
              <input
                type="text"
                value={formData.supplierBankAccount}
                onChange={(e) => handleFieldChange("supplierBankAccount", e.target.value)}
                disabled={isLoading}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
                placeholder="1234567890123456789012"
              />
            </div>

            {/* Банковский счет покупателя */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Банковский счет покупателя
              </label>
              <input
                type="text"
                value={formData.contractorBankAccount}
                onChange={(e) => handleFieldChange("contractorBankAccount", e.target.value)}
                disabled={isLoading}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
                placeholder="1234567890123456789012"
              />
            </div>
          </div>
        </div>

        {/* Секция 3: Условия доставки и оплаты */}
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Условия доставки и оплаты
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Способ доставки */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Способ доставки <span className="text-red-500">*</span>
              </label>
              <select
                value={formData.deliveryTypeCode}
                onChange={(e) => handleFieldChange("deliveryTypeCode", e.target.value)}
                disabled={isLoading}
                className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                  errors.deliveryTypeCode
                    ? "border-red-500"
                    : "border-gray-300 dark:border-gray-600"
                } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              >
                <option value={ESFDeliveryType.DIRECT}>
                  {getDeliveryTypeLabel(ESFDeliveryType.DIRECT)}
                </option>
                <option value={ESFDeliveryType.WAREHOUSE}>
                  {getDeliveryTypeLabel(ESFDeliveryType.WAREHOUSE)}
                </option>
                <option value={ESFDeliveryType.SERVICE}>
                  {getDeliveryTypeLabel(ESFDeliveryType.SERVICE)}
                </option>
                <option value={ESFDeliveryType.COURIER}>
                  {getDeliveryTypeLabel(ESFDeliveryType.COURIER)}
                </option>
                <option value={ESFDeliveryType.PICKUP}>
                  {getDeliveryTypeLabel(ESFDeliveryType.PICKUP)}
                </option>
                <option value={ESFDeliveryType.POST}>
                  {getDeliveryTypeLabel(ESFDeliveryType.POST)}
                </option>
              </select>
              {errors.deliveryTypeCode && (
                <p className="mt-1 text-sm text-red-500">{errors.deliveryTypeCode}</p>
              )}
            </div>

            {/* Тип платежа */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Тип платежа <span className="text-red-500">*</span>
              </label>
              <select
                value={formData.paymentCode}
                onChange={(e) => handleFieldChange("paymentCode", e.target.value)}
                disabled={isLoading}
                className={`w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100 ${
                  errors.paymentCode
                    ? "border-red-500"
                    : "border-gray-300 dark:border-gray-600"
                } ${isLoading ? "opacity-50 cursor-not-allowed" : ""}`}
              >
                <option value={ESFPaymentType.CASH}>
                  {getPaymentTypeLabel(ESFPaymentType.CASH)}
                </option>
                <option value={ESFPaymentType.BANK_TRANSFER}>
                  {getPaymentTypeLabel(ESFPaymentType.BANK_TRANSFER)}
                </option>
                <option value={ESFPaymentType.CHECK}>
                  {getPaymentTypeLabel(ESFPaymentType.CHECK)}
                </option>
                <option value={ESFPaymentType.CARD}>
                  {getPaymentTypeLabel(ESFPaymentType.CARD)}
                </option>
                <option value={ESFPaymentType.CREDIT}>
                  {getPaymentTypeLabel(ESFPaymentType.CREDIT)}
                </option>
                <option value={ESFPaymentType.ELECTRONIC}>
                  {getPaymentTypeLabel(ESFPaymentType.ELECTRONIC)}
                </option>
              </select>
              {errors.paymentCode && (
                <p className="mt-1 text-sm text-red-500">{errors.paymentCode}</p>
              )}
            </div>
          </div>
        </div>

        {/* Секция 4: Позиции счета-фактуры */}
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
          <InvoiceCatalogEntriesTable
            entries={formData.catalogEntries}
            vatRate={vatRate}
            currency={formData.currencyCode}
            onEntriesChange={handleCatalogEntriesChange}
            disabled={isLoading}
          />
          {errors.catalogEntries && (
            <div className="mt-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg flex items-start gap-2">
              <AlertCircle className="text-red-500 shrink-0 mt-0.5" size={18} />
              <p className="text-sm text-red-700 dark:text-red-300">{errors.catalogEntries}</p>
            </div>
          )}
        </div>

        {/* Секция 5: Итоговые суммы */}
        <InvoiceTotals
          totalAmountWithoutVAT={totals.totalAmountWithoutVAT}
          totalVATAmount={totals.totalVATAmount}
          totalAmount={totals.totalAmount}
          currency={formData.currencyCode}
        />

        {/* Секция 6: Комментарий */}
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            Дополнительная информация
          </h3>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Комментарий
            </label>
            <textarea
              value={formData.comment}
              onChange={(e) => handleFieldChange("comment", e.target.value)}
              disabled={isLoading}
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-800 dark:text-gray-100"
              placeholder="Дополнительная информация о счете-фактуре"
            />
          </div>
        </div>

        {/* Кнопки действий */}
        <div className="flex items-center gap-3 pt-4 border-t">
          <button
            type="submit"
            disabled={isLoading}
            className="flex-1 flex items-center justify-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:ring-4 focus:ring-blue-300 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
          >
            {isLoading ? (
              <>
                <Loader2 size={20} className="animate-spin" />
                {isEditing ? "Сохранение..." : "Создание..."}
              </>
            ) : (
              <>
                <Save size={20} />
                {isEditing ? "Сохранить изменения" : "Создать счет-фактуру"}
              </>
            )}
          </button>

          {onCancel && (
            <button
              type="button"
              onClick={onCancel}
              disabled={isLoading}
              className="px-6 py-3 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
            >
              Отмена
            </button>
          )}

          {isEditing && invoice && (
            <button
              type="button"
              onClick={() => setShowDeleteConfirm(true)}
              disabled={isLoading}
              className="px-6 py-3 bg-red-600 text-white rounded-lg hover:bg-red-700 focus:ring-4 focus:ring-red-300 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
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
              Вы уверены, что хотите удалить счет-фактуру{" "}
              {invoice?.invoiceNumber ? `"${invoice.invoiceNumber}"` : ""}? Это
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

// Export с обычным именем для совместимости
export default InvoiceFormESF;
