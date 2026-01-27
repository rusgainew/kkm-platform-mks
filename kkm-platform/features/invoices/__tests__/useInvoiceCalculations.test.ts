import { describe, it, expect } from "vitest";
import { renderHook } from "@testing-library/react";
import { useInvoiceCalculations } from "../hooks/useInvoiceCalculations";
import type { CatalogEntry } from "@/types/invoice";

describe("useInvoiceCalculations", () => {
  const mockEntry: CatalogEntry = {
    id: 1,
    unitClassificationCode: "796",
    salesTaxCode: "10",
    quantity: 10,
    price: 1000,
    vatAmount: 0,
    salesTaxAmount: 0,
    amountWithoutTaxes: 0,
    totalAmount: 0,
  };

  describe("totals calculation", () => {
    it("should calculate totals with 12% VAT", () => {
      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [mockEntry],
          vatRate: 12,
        }),
      );

      expect(result.current.totals.totalAmount).toBe(10000);
      expect(result.current.totals.totalAmountWithoutVAT).toBeCloseTo(
        8928.57,
        1,
      );
      expect(result.current.totals.totalVATAmount).toBeCloseTo(1071.43, 1);
    });

    it("should calculate totals with 0% VAT", () => {
      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [mockEntry],
          vatRate: 0,
        }),
      );

      expect(result.current.totals.totalAmount).toBe(10000);
      expect(result.current.totals.totalAmountWithoutVAT).toBe(10000);
      expect(result.current.totals.totalVATAmount).toBe(0);
    });

    it("should handle empty entries", () => {
      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [],
          vatRate: 12,
        }),
      );

      expect(result.current.totals.totalAmount).toBe(0);
      expect(result.current.totals.totalAmountWithoutVAT).toBe(0);
      expect(result.current.totals.totalVATAmount).toBe(0);
    });

    it("should recalculate when entries change", () => {
      const { result, rerender } = renderHook(
        ({ entries }) =>
          useInvoiceCalculations({ catalogEntries: entries, vatRate: 12 }),
        { initialProps: { entries: [mockEntry] } },
      );

      const initialTotal = result.current.totals.totalAmount;
      expect(initialTotal).toBe(10000);

      // Add another entry
      rerender({ entries: [mockEntry, { ...mockEntry, id: 2 }] });
      expect(result.current.totals.totalAmount).toBe(20000);
    });

    it("should recalculate when VAT rate changes", () => {
      const { result, rerender } = renderHook(
        ({ entries, rate }) =>
          useInvoiceCalculations({ catalogEntries: entries, vatRate: rate }),
        { initialProps: { entries: [mockEntry], rate: 12 } },
      );

      const initialVAT = result.current.totals.totalVATAmount;
      expect(initialVAT).toBeCloseTo(1071.43, 1);

      // Change VAT rate to 0
      rerender({ entries: [mockEntry], rate: 0 });
      expect(result.current.totals.totalVATAmount).toBe(0);
    });
  });

  describe("recalculateEntry", () => {
    it("should recalculate single entry with 12% VAT", () => {
      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [mockEntry],
          vatRate: 12,
        }),
      );

      const recalculated = result.current.recalculateEntry(mockEntry);

      expect(recalculated.totalAmount).toBe(10000);
      expect(recalculated.amountWithoutTaxes).toBeCloseTo(8928.57, 1);
      expect(recalculated.vatAmount).toBeCloseTo(1071.43, 1);
    });

    it("should recalculate single entry with 0% VAT", () => {
      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [mockEntry],
          vatRate: 0,
        }),
      );

      const recalculated = result.current.recalculateEntry(mockEntry);

      expect(recalculated.totalAmount).toBe(10000);
      expect(recalculated.amountWithoutTaxes).toBe(10000);
      expect(recalculated.vatAmount).toBe(0);
    });

    it("should preserve all entry fields", () => {
      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [mockEntry],
          vatRate: 12,
        }),
      );

      const recalculated = result.current.recalculateEntry(mockEntry);

      expect(recalculated.id).toBe(mockEntry.id);
      expect(recalculated.unitClassificationCode).toBe(
        mockEntry.unitClassificationCode,
      );
      expect(recalculated.salesTaxCode).toBe(mockEntry.salesTaxCode);
      expect(recalculated.quantity).toBe(mockEntry.quantity);
      expect(recalculated.price).toBe(mockEntry.price);
    });
  });

  describe("recalculateAllEntries", () => {
    it("should recalculate all entries with 12% VAT", () => {
      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [mockEntry],
          vatRate: 12,
        }),
      );

      const entries = [mockEntry, { ...mockEntry, id: 2, quantity: 5 }];

      const recalculated = result.current.recalculateAllEntries(entries);

      expect(recalculated).toHaveLength(2);
      expect(recalculated[0].totalAmount).toBe(10000);
      expect(recalculated[1].totalAmount).toBe(5000);
    });

    it("should handle empty array", () => {
      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [],
          vatRate: 12,
        }),
      );

      const recalculated = result.current.recalculateAllEntries([]);
      expect(recalculated).toEqual([]);
    });

    it("should preserve order of entries", () => {
      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [mockEntry],
          vatRate: 12,
        }),
      );

      const entries = [
        { ...mockEntry, id: 1 },
        { ...mockEntry, id: 2 },
        { ...mockEntry, id: 3 },
      ];

      const recalculated = result.current.recalculateAllEntries(entries);

      expect(recalculated[0].id).toBe(1);
      expect(recalculated[1].id).toBe(2);
      expect(recalculated[2].id).toBe(3);
    });
  });

  describe("memoization", () => {
    it("should memoize totals when inputs don't change", () => {
      const { result, rerender } = renderHook(
        ({ entries, rate }) =>
          useInvoiceCalculations({ catalogEntries: entries, vatRate: rate }),
        { initialProps: { entries: [mockEntry], rate: 12 } },
      );

      const firstTotals = result.current.totals;

      // Force re-render with same props
      rerender({ entries: [mockEntry], rate: 12 });

      // Should return same values
      expect(result.current.totals).toStrictEqual(firstTotals);
    });

    it("should recalculate when entries reference changes", () => {
      const { result, rerender } = renderHook(
        ({ entries, rate }) =>
          useInvoiceCalculations({ catalogEntries: entries, vatRate: rate }),
        { initialProps: { entries: [mockEntry], rate: 12 } },
      );

      const firstTotals = result.current.totals;

      // Create new entries array (different reference but same values)
      rerender({ entries: [{ ...mockEntry }], rate: 12 });

      // Should recalculate
      expect(result.current.totals).not.toBe(firstTotals);
    });
  });

  describe("edge cases", () => {
    it("should handle very large numbers", () => {
      const largeEntry: CatalogEntry = {
        ...mockEntry,
        quantity: 1000000,
        price: 999999,
      };

      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [largeEntry],
          vatRate: 12,
        }),
      );

      expect(result.current.totals.totalAmount).toBe(999999000000);
    });

    it("should handle decimal values", () => {
      const decimalEntry: CatalogEntry = {
        ...mockEntry,
        quantity: 1.5,
        price: 99.99,
      };

      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [decimalEntry],
          vatRate: 12,
        }),
      );

      expect(result.current.totals.totalAmount).toBeCloseTo(149.98, 1);
    });

    it("should handle zero values", () => {
      const zeroEntry: CatalogEntry = {
        ...mockEntry,
        quantity: 0,
        price: 0,
      };

      const { result } = renderHook(() =>
        useInvoiceCalculations({
          catalogEntries: [zeroEntry],
          vatRate: 12,
        }),
      );

      expect(result.current.totals.totalAmount).toBe(0);
      expect(result.current.totals.totalVATAmount).toBe(0);
    });
  });
});
