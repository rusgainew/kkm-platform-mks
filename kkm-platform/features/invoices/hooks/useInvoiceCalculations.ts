"use client";

import { useMemo, useCallback } from "react";
import type { CatalogEntry } from "@/types/invoice";
import {
  calculateInvoiceTotals,
  updateCatalogEntryCalculations,
} from "../lib/invoice-validation";

interface UseInvoiceCalculationsProps {
  catalogEntries: CatalogEntry[];
  vatRate: number;
}

interface UseInvoiceCalculationsReturn {
  totals: {
    totalAmountWithoutVAT: number;
    totalVATAmount: number;
    totalAmount: number;
  };
  recalculateEntry: (entry: CatalogEntry) => CatalogEntry;
  recalculateAllEntries: (entries: CatalogEntry[]) => CatalogEntry[];
}

/**
 * Хук для автоматических расчетов в счете-фактуре
 * Обрабатывает расчет НДС, сумм без НДС и итоговых сумм
 */
export function useInvoiceCalculations({
  catalogEntries,
  vatRate,
}: UseInvoiceCalculationsProps): UseInvoiceCalculationsReturn {
  /**
   * Расчет итоговых сумм
   */
  const totals = useMemo(() => {
    return calculateInvoiceTotals(catalogEntries, vatRate);
  }, [catalogEntries, vatRate]);

  /**
   * Пересчет одной позиции
   */
  const recalculateEntry = useCallback(
    (entry: CatalogEntry): CatalogEntry => {
      return updateCatalogEntryCalculations(entry, vatRate);
    },
    [vatRate],
  );

  /**
   * Пересчет всех позиций
   */
  const recalculateAllEntries = useCallback(
    (entries: CatalogEntry[]): CatalogEntry[] => {
      return entries.map((entry) =>
        updateCatalogEntryCalculations(entry, vatRate),
      );
    },
    [vatRate],
  );

  return {
    totals,
    recalculateEntry,
    recalculateAllEntries,
  };
}
