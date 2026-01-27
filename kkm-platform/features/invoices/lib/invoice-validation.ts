"use client";

import type { CreateInvoiceRequest, CatalogEntry } from "@/types/invoice";
import {
  ESFOperationType,
  ESFDeliveryType,
  ESFPaymentType,
} from "@/types/enums";

interface FormData extends Omit<CreateInvoiceRequest, "catalogEntries"> {
  catalogEntries: CatalogEntry[];
}

interface FormErrors {
  [key: string]: string;
}

/**
 * Валидация ИНН (14 цифр для КР)
 */
export function isValidTIN(tin: string): boolean {
  return /^\d{14}$/.test(tin);
}

/**
 * Валидация даты (YYYY-MM-DD)
 */
export function isValidDate(date: string): boolean {
  const regex = /^\d{4}-\d{2}-\d{2}$/;
  if (!regex.test(date)) return false;
  const d = new Date(date);
  return d instanceof Date && !isNaN(d.getTime());
}

/**
 * Валидация номера счета
 */
export function isValidInvoiceNumber(number: string): boolean {
  return number.trim().length >= 1 && number.trim().length <= 50;
}

/**
 * Расчет суммы без НДС
 */
export function calculateAmountWithoutVAT(
  amount: number,
  vatRate: number,
): number {
  return amount / (1 + vatRate / 100);
}

/**
 * Расчет НДС
 */
export function calculateVAT(amount: number, vatRate: number): number {
  const withoutVAT = calculateAmountWithoutVAT(amount, vatRate);
  return amount - withoutVAT;
}

/**
 * Расчет общей суммы с НДС
 */
export function calculateAmountWithVAT(
  amountWithoutVAT: number,
  vatRate: number,
): number {
  return amountWithoutVAT * (1 + vatRate / 100);
}

/**
 * Валидация catalogEntry
 */
export function validateCatalogEntry(entry: CatalogEntry): string[] {
  const errors: string[] = [];

  if (!entry.id || entry.id <= 0) {
    errors.push("ID товара обязателен");
  }

  if (!entry.quantity || entry.quantity <= 0) {
    errors.push("Количество должно быть больше 0");
  }

  if (!entry.price || entry.price <= 0) {
    errors.push("Цена должна быть больше 0");
  }

  if (
    !entry.unitClassificationCode ||
    entry.unitClassificationCode.trim() === ""
  ) {
    errors.push("Код единицы измерения обязателен");
  }

  if (!entry.salesTaxCode || entry.salesTaxCode.trim() === "") {
    errors.push("Код налога обязателен");
  }

  return errors;
}

/**
 * Валидация всей формы Invoice
 */
export function validateInvoiceForm(data: FormData): FormErrors {
  const errors: FormErrors = {};

  // Тип операции
  if (!data.operationTypeCode) {
    errors.operationTypeCode = "Тип операции обязателен";
  }

  // Дата доставки
  if (!data.deliveryDate) {
    errors.deliveryDate = "Дата доставки обязательна";
  } else if (!isValidDate(data.deliveryDate)) {
    errors.deliveryDate = "Неверный формат даты (YYYY-MM-DD)";
  } else {
    const deliveryDate = new Date(data.deliveryDate);
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    if (deliveryDate > today) {
      errors.deliveryDate = "Дата доставки не может быть в будущем";
    }
  }

  // ИНН покупателя
  if (!data.contractorTin) {
    errors.contractorTin = "ИНН покупателя обязателен";
  } else if (data.isResident && !isValidTIN(data.contractorTin)) {
    errors.contractorTin = "Неверный формат ИНН (14 цифр)";
  } else if (!data.isResident && data.contractorTin.length < 5) {
    errors.contractorTin =
      "ИНН иностранного контрагента должен содержать минимум 5 символов";
  }

  // Способ доставки
  if (!data.deliveryTypeCode) {
    errors.deliveryTypeCode = "Способ доставки обязателен";
  }

  // Тип платежа
  if (!data.paymentCode) {
    errors.paymentCode = "Тип платежа обязателен";
  }

  // Валюта
  if (!data.currencyCode) {
    errors.currencyCode = "Валюта обязательна";
  }

  // Ставка НДС
  if (!data.taxRateVATCode) {
    errors.taxRateVATCode = "Ставка НДС обязательна";
  }

  // Номер счета (если указан)
  if (data.invoiceNumber && !isValidInvoiceNumber(data.invoiceNumber)) {
    errors.invoiceNumber = "Неверный формат номера счета (1-50 символов)";
  }

  // Товары/услуги
  if (!data.catalogEntries || data.catalogEntries.length === 0) {
    errors.catalogEntries = "Необходимо добавить хотя бы одну позицию";
  } else {
    const entryErrors: string[] = [];
    data.catalogEntries.forEach((entry, index) => {
      const validation = validateCatalogEntry(entry);
      if (validation.length > 0) {
        entryErrors.push(`Позиция ${index + 1}: ${validation.join(", ")}`);
      }
    });
    if (entryErrors.length > 0) {
      errors.catalogEntries = entryErrors.join("; ");
    }
  }

  // Курс валюты (если не KGS)
  if (data.currencyCode && data.currencyCode !== "KGS" && !data.currencyRate) {
    errors.currencyRate = "Курс валюты обязателен для иностранной валюты";
  }

  // Дата счета-фактуры (если указана)
  if (data.invoiceDate) {
    if (!isValidDate(data.invoiceDate)) {
      errors.invoiceDate = "Неверный формат даты счета-фактуры";
    }
  }

  return errors;
}

/**
 * Расчет итоговых сумм по catalogEntries
 */
export function calculateInvoiceTotals(
  entries: CatalogEntry[],
  vatRate: number,
): {
  totalAmountWithoutVAT: number;
  totalVATAmount: number;
  totalAmount: number;
} {
  let totalAmountWithoutVAT = 0;
  let totalVATAmount = 0;
  let totalAmount = 0;

  entries.forEach((entry) => {
    const entryTotal = entry.quantity * entry.price;
    const entryWithoutVAT =
      entry.amountWithoutTaxes ||
      calculateAmountWithoutVAT(entryTotal, vatRate);
    const entryVAT = entry.vatAmount || calculateVAT(entryTotal, vatRate);

    totalAmountWithoutVAT += entryWithoutVAT;
    totalVATAmount += entryVAT;
    totalAmount += entry.totalAmount || entryTotal;
  });

  return {
    totalAmountWithoutVAT: Math.round(totalAmountWithoutVAT * 100) / 100,
    totalVATAmount: Math.round(totalVATAmount * 100) / 100,
    totalAmount: Math.round(totalAmount * 100) / 100,
  };
}

/**
 * Форматирование суммы с валютой
 */
export function formatCurrency(
  amount: number,
  currency: string = "KGS",
): string {
  return (
    new Intl.NumberFormat("ru-KG", {
      style: "decimal",
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(amount) + ` ${currency}`
  );
}

/**
 * Получить label для типа операции
 */
export function getOperationTypeLabel(code: string): string {
  const labels: Record<string, string> = {
    [ESFOperationType.SALES]: "Реализация",
    [ESFOperationType.PURCHASE_SERVICES]: "Покупка/Услуги",
    [ESFOperationType.OTHER]: "Прочее",
  };
  return labels[code] || code;
}

/**
 * Получить label для способа доставки
 */
export function getDeliveryTypeLabel(code: string): string {
  const labels: Record<string, string> = {
    [ESFDeliveryType.DIRECT]: "Прямая доставка",
    [ESFDeliveryType.WAREHOUSE]: "Через склад",
    [ESFDeliveryType.SERVICE]: "Услуга доставки",
    [ESFDeliveryType.COURIER]: "Курьер",
    [ESFDeliveryType.PICKUP]: "Самовывоз",
    [ESFDeliveryType.POST]: "Почта",
  };
  return labels[code] || code;
}

/**
 * Получить label для типа платежа
 */
export function getPaymentTypeLabel(code: string): string {
  const labels: Record<string, string> = {
    [ESFPaymentType.CASH]: "Наличные",
    [ESFPaymentType.BANK_TRANSFER]: "Банковский перевод",
    [ESFPaymentType.CHECK]: "Чек",
    [ESFPaymentType.CARD]: "Карточка",
    [ESFPaymentType.CREDIT]: "В кредит",
    [ESFPaymentType.ELECTRONIC]: "Электронный платеж",
  };
  return labels[code] || code;
}

/**
 * Создать пустой catalogEntry
 */
export function createEmptyCatalogEntry(): CatalogEntry {
  return {
    id: 0,
    unitClassificationCode: "796", // Штука по умолчанию
    salesTaxCode: "10", // НДС по умолчанию
    quantity: 1,
    price: 0,
    vatAmount: 0,
    salesTaxAmount: 0,
    amountWithoutTaxes: 0,
    totalAmount: 0,
  };
}

/**
 * Обновить расчеты catalogEntry
 */
export function updateCatalogEntryCalculations(
  entry: CatalogEntry,
  vatRate: number,
): CatalogEntry {
  const totalAmount = entry.quantity * entry.price;
  const amountWithoutTaxes = calculateAmountWithoutVAT(totalAmount, vatRate);
  const vatAmount = calculateVAT(totalAmount, vatRate);

  return {
    ...entry,
    totalAmount: Math.round(totalAmount * 100) / 100,
    amountWithoutTaxes: Math.round(amountWithoutTaxes * 100) / 100,
    vatAmount: Math.round(vatAmount * 100) / 100,
    salesTaxAmount: 0, // Пока не реализовано
  };
}
