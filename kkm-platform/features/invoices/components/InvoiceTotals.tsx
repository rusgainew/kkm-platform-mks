"use client";

import { formatCurrency } from "../lib/invoice-validation";

interface InvoiceTotalsProps {
  totalAmountWithoutVAT: number;
  totalVATAmount: number;
  totalSTAmount?: number;
  totalAmount: number;
  currency?: string;
  className?: string;
}

/**
 * Компонент для отображения итоговых сумм счета-фактуры
 * Показывает сумму без НДС, НДС, НСП (если есть) и итоговую сумму
 */
export function InvoiceTotals({
  totalAmountWithoutVAT,
  totalVATAmount,
  totalSTAmount = 0,
  totalAmount,
  currency = "KGS",
  className = "",
}: InvoiceTotalsProps) {
  return (
    <div className={`bg-gray-50 dark:bg-gray-800 rounded-lg p-6 space-y-3 ${className}`}>
      <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
        Итоговые суммы
      </h3>

      {/* Сумма без НДС */}
      <div className="flex justify-between items-center pb-2 border-b border-gray-200 dark:border-gray-700">
        <span className="text-sm text-gray-600 dark:text-gray-400">
          Сумма без НДС:
        </span>
        <span className="text-base font-medium text-gray-900 dark:text-gray-100">
          {formatCurrency(totalAmountWithoutVAT, currency)}
        </span>
      </div>

      {/* НДС */}
      <div className="flex justify-between items-center pb-2 border-b border-gray-200 dark:border-gray-700">
        <span className="text-sm text-gray-600 dark:text-gray-400">
          НДС:
        </span>
        <span className="text-base font-medium text-blue-600 dark:text-blue-400">
          {formatCurrency(totalVATAmount, currency)}
        </span>
      </div>

      {/* НСП (если есть) */}
      {totalSTAmount > 0 && (
        <div className="flex justify-between items-center pb-2 border-b border-gray-200 dark:border-gray-700">
          <span className="text-sm text-gray-600 dark:text-gray-400">
            НСП:
          </span>
          <span className="text-base font-medium text-purple-600 dark:text-purple-400">
            {formatCurrency(totalSTAmount, currency)}
          </span>
        </div>
      )}

      {/* Итоговая сумма */}
      <div className="flex justify-between items-center pt-3 border-t-2 border-gray-300 dark:border-gray-600">
        <span className="text-lg font-semibold text-gray-900 dark:text-gray-100">
          Итого к оплате:
        </span>
        <span className="text-xl font-bold text-green-600 dark:text-green-400">
          {formatCurrency(totalAmount, currency)}
        </span>
      </div>

      {/* Дополнительная информация */}
      <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
        <p className="text-xs text-gray-500 dark:text-gray-400 text-center">
          Все суммы округлены до копеек
        </p>
      </div>
    </div>
  );
}
