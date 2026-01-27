import { describe, it, expect } from "vitest";
import type { CatalogEntry, CreateInvoiceRequest } from "@/types/invoice";
import {
  isValidTIN,
  isValidDate,
  isValidInvoiceNumber,
  calculateAmountWithoutVAT,
  calculateVAT,
  calculateAmountWithVAT,
  validateCatalogEntry,
  validateInvoiceForm,
  calculateInvoiceTotals,
  formatCurrency,
  getOperationTypeLabel,
  getDeliveryTypeLabel,
  getPaymentTypeLabel,
  createEmptyCatalogEntry,
  updateCatalogEntryCalculations,
} from "../lib/invoice-validation";

describe("invoice-validation", () => {
  describe("isValidTIN", () => {
    it("should validate correct 14-digit TIN", () => {
      expect(isValidTIN("12345678901234")).toBe(true);
      expect(isValidTIN("00000000000000")).toBe(true);
      expect(isValidTIN("99999999999999")).toBe(true);
    });

    it("should reject invalid TINs", () => {
      expect(isValidTIN("123")).toBe(false);
      expect(isValidTIN("123456789012345")).toBe(false);
      expect(isValidTIN("1234567890123a")).toBe(false);
      expect(isValidTIN("")).toBe(false);
    });
  });

  describe("isValidDate", () => {
    it("should validate correct date format", () => {
      expect(isValidDate("2024-01-15")).toBe(true);
      expect(isValidDate("2024-12-31")).toBe(true);
      expect(isValidDate("2025-06-15")).toBe(true);
    });

    it("should reject invalid dates", () => {
      expect(isValidDate("2024-13-01")).toBe(false);
      expect(isValidDate("2024-01-32")).toBe(false);
      expect(isValidDate("2024/01/15")).toBe(false);
      expect(isValidDate("")).toBe(false);
    });

    it("should validate future dates (no future check in isValidDate)", () => {
      const futureDate = new Date();
      futureDate.setFullYear(futureDate.getFullYear() + 1);
      const futureDateStr = futureDate.toISOString().split("T")[0];
      expect(isValidDate(futureDateStr)).toBe(true);
    });
  });

  describe("isValidInvoiceNumber", () => {
    it("should validate correct invoice numbers", () => {
      expect(isValidInvoiceNumber("INV-001")).toBe(true);
      expect(isValidInvoiceNumber("2024-001")).toBe(true);
      expect(isValidInvoiceNumber("A")).toBe(true);
      expect(isValidInvoiceNumber("A".repeat(50))).toBe(true);
    });

    it("should reject invalid invoice numbers", () => {
      expect(isValidInvoiceNumber("")).toBe(false);
      expect(isValidInvoiceNumber("A".repeat(51))).toBe(false);
    });
  });

  describe("calculateAmountWithoutVAT", () => {
    it("should calculate amount without VAT for 12%", () => {
      expect(calculateAmountWithoutVAT(11200, 12)).toBeCloseTo(10000, 2);
      expect(calculateAmountWithoutVAT(112, 12)).toBeCloseTo(100, 2);
      expect(calculateAmountWithoutVAT(1120, 12)).toBeCloseTo(1000, 2);
    });

    it("should handle 0% VAT", () => {
      expect(calculateAmountWithoutVAT(10000, 0)).toBe(10000);
      expect(calculateAmountWithoutVAT(100, 0)).toBe(100);
    });

    it("should handle edge cases", () => {
      expect(calculateAmountWithoutVAT(0, 12)).toBe(0);
      expect(calculateAmountWithoutVAT(1, 12)).toBeCloseTo(0.89, 2);
    });
  });

  describe("calculateVAT", () => {
    it("should calculate VAT from amount with VAT for 12%", () => {
      expect(calculateVAT(11200, 12)).toBeCloseTo(1200, 2);
      expect(calculateVAT(112, 12)).toBeCloseTo(12, 2);
      expect(calculateVAT(1120, 12)).toBeCloseTo(120, 2);
    });

    it("should handle 0% VAT", () => {
      expect(calculateVAT(10000, 0)).toBe(0);
      expect(calculateVAT(100, 0)).toBe(0);
    });

    it("should handle edge cases", () => {
      expect(calculateVAT(0, 12)).toBe(0);
      expect(calculateVAT(1.12, 12)).toBeCloseTo(0.12, 2);
    });
  });

  describe("calculateAmountWithVAT", () => {
    it("should calculate amount with VAT for 12%", () => {
      expect(calculateAmountWithVAT(10000, 12)).toBeCloseTo(11200, 2);
      expect(calculateAmountWithVAT(100, 12)).toBeCloseTo(112, 2);
      expect(calculateAmountWithVAT(1000, 12)).toBeCloseTo(1120, 2);
    });

    it("should handle 0% VAT", () => {
      expect(calculateAmountWithVAT(10000, 0)).toBe(10000);
      expect(calculateAmountWithVAT(100, 0)).toBe(100);
    });

    it("should handle edge cases", () => {
      expect(calculateAmountWithVAT(0, 12)).toBe(0);
      expect(calculateAmountWithVAT(1, 12)).toBeCloseTo(1.12, 2);
    });
  });

  describe("validateCatalogEntry", () => {
    const validEntry: CatalogEntry = {
      id: 1,
      unitClassificationCode: "796",
      salesTaxCode: "10",
      quantity: 10,
      price: 1000,
      vatAmount: 1200,
      salesTaxAmount: 0,
      amountWithoutTaxes: 10000,
      totalAmount: 11200,
    };

    it("should validate correct catalog entry", () => {
      const errors = validateCatalogEntry(validEntry);
      expect(errors).toEqual([]);
    });

    it("should reject entry with invalid id", () => {
      const entry = { ...validEntry, id: 0 };
      const errors = validateCatalogEntry(entry);
      expect(errors.length).toBeGreaterThan(0);
    });

    it("should reject entry with invalid quantity", () => {
      const entry1 = { ...validEntry, quantity: 0 };
      const errors1 = validateCatalogEntry(entry1);
      expect(errors1.length).toBeGreaterThan(0);

      const entry2 = { ...validEntry, quantity: -5 };
      const errors2 = validateCatalogEntry(entry2);
      expect(errors2.length).toBeGreaterThan(0);
    });

    it("should reject entry with invalid price", () => {
      const entry1 = { ...validEntry, price: 0 };
      const errors1 = validateCatalogEntry(entry1);
      expect(errors1.length).toBeGreaterThan(0);

      const entry2 = { ...validEntry, price: -100 };
      const errors2 = validateCatalogEntry(entry2);
      expect(errors2.length).toBeGreaterThan(0);
    });
  });

  describe("validateInvoiceForm", () => {
    const validFormData: CreateInvoiceRequest = {
      operationTypeCode: "001",
      invoiceNumber: "INV-001",
      deliveryDate: "2024-01-15",
      invoiceDate: "2024-01-15",
      currencyCode: "KGS",
      taxRateVATCode: "12",
      isResident: true,
      contractorTin: "12345678901234",
      deliveryTypeCode: "01",
      paymentCode: "1",
      catalogEntries: [
        {
          id: 1,
          unitClassificationCode: "796",
          salesTaxCode: "10",
          quantity: 10,
          price: 1000,
          vatAmount: 1200,
          salesTaxAmount: 0,
          amountWithoutTaxes: 10000,
          totalAmount: 11200,
        },
      ],
    };

    it("should validate correct form data", () => {
      const errors = validateInvoiceForm(validFormData);
      expect(errors).toEqual({});
    });

    it("should reject form with invalid TIN for resident", () => {
      const data = { ...validFormData, contractorTin: "123" };
      const errors = validateInvoiceForm(data);
      expect(errors.contractorTin).toBeDefined();
    });

    it("should reject form without catalog entries", () => {
      const data = { ...validFormData, catalogEntries: [] };
      const errors = validateInvoiceForm(data);
      expect(errors.catalogEntries).toBeDefined();
    });
  });

  describe("calculateInvoiceTotals", () => {
    const entries: CatalogEntry[] = [
      {
        id: 1,
        unitClassificationCode: "796",
        salesTaxCode: "10",
        quantity: 10,
        price: 1000,
        vatAmount: 1071.43,
        salesTaxAmount: 0,
        amountWithoutTaxes: 8928.57,
        totalAmount: 10000,
      },
      {
        id: 2,
        unitClassificationCode: "796",
        salesTaxCode: "10",
        quantity: 5,
        price: 2000,
        vatAmount: 1071.43,
        salesTaxAmount: 0,
        amountWithoutTaxes: 8928.57,
        totalAmount: 10000,
      },
    ];

    it("should calculate totals with 12% VAT from entries", () => {
      const totals = calculateInvoiceTotals(entries, 12);
      expect(totals.totalAmount).toBe(20000);
    });

    it("should calculate totals with 0% VAT", () => {
      const totals = calculateInvoiceTotals([], 0);
      expect(totals.totalAmountWithoutVAT).toBe(0);
      expect(totals.totalVATAmount).toBe(0);
      expect(totals.totalAmount).toBe(0);
    });

    it("should handle empty entries", () => {
      const totals = calculateInvoiceTotals([], 12);
      expect(totals.totalAmountWithoutVAT).toBe(0);
      expect(totals.totalVATAmount).toBe(0);
      expect(totals.totalAmount).toBe(0);
    });
  });

  describe("formatCurrency", () => {
    it("should format KGS currency", () => {
      const result = formatCurrency(10000, "KGS");
      expect(result).toContain("10");
      expect(result).toContain("000");
    });

    it("should format USD currency", () => {
      const result = formatCurrency(10000, "USD");
      expect(result).toContain("10");
      expect(result).toContain("000");
    });

    it("should handle zero amount", () => {
      expect(formatCurrency(0, "KGS")).toContain("0");
    });

    it("should handle decimal amounts", () => {
      const result = formatCurrency(1234.56, "KGS");
      expect(result).toContain("1");
      expect(result).toContain("234");
    });
  });

  describe("getOperationTypeLabel", () => {
    it("should return label for valid codes from enums", () => {
      const label = getOperationTypeLabel("1");
      expect(label).toBeTruthy();
    });

    it("should return code itself for unknown codes", () => {
      expect(getOperationTypeLabel("999")).toBe("999");
      expect(getOperationTypeLabel("")).toBe("");
    });
  });

  describe("getDeliveryTypeLabel", () => {
    it("should return label for valid codes from enums", () => {
      const label = getDeliveryTypeLabel("1");
      expect(label).toBeTruthy();
    });

    it("should return code itself for unknown codes", () => {
      expect(getDeliveryTypeLabel("99")).toBe("99");
      expect(getDeliveryTypeLabel("")).toBe("");
    });
  });

  describe("getPaymentTypeLabel", () => {
    it("should return label for valid codes from enums", () => {
      expect(getPaymentTypeLabel("1")).toBe("Наличные");
      expect(getPaymentTypeLabel("2")).toBe("Банковский перевод");
    });

    it("should return code itself for unknown codes", () => {
      expect(getPaymentTypeLabel("9")).toBe("9");
      expect(getPaymentTypeLabel("")).toBe("");
    });
  });

  describe("createEmptyCatalogEntry", () => {
    it("should create empty catalog entry with defaults", () => {
      const entry = createEmptyCatalogEntry();
      expect(entry.id).toBe(0);
      expect(entry.unitClassificationCode).toBe("796");
      expect(entry.salesTaxCode).toBe("10");
      expect(entry.quantity).toBe(1);
      expect(entry.price).toBe(0);
      expect(entry.totalAmount).toBe(0);
    });

    it("should always create same default entry (no unique IDs)", () => {
      const entry1 = createEmptyCatalogEntry();
      const entry2 = createEmptyCatalogEntry();
      expect(entry1.id).toBe(entry2.id);
    });
  });

  describe("updateCatalogEntryCalculations", () => {
    it("should recalculate totals with 12% VAT", () => {
      const entry: CatalogEntry = {
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

      const updated = updateCatalogEntryCalculations(entry, 12);
      expect(updated.totalAmount).toBe(10000);
      expect(updated.amountWithoutTaxes).toBeCloseTo(8928.57, 1);
      expect(updated.vatAmount).toBeCloseTo(1071.43, 1);
    });

    it("should recalculate totals with 0% VAT", () => {
      const entry: CatalogEntry = {
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

      const updated = updateCatalogEntryCalculations(entry, 0);
      expect(updated.totalAmount).toBe(10000);
    });

    it("should handle zero quantity", () => {
      const entry: CatalogEntry = {
        id: 1,
        unitClassificationCode: "796",
        salesTaxCode: "10",
        quantity: 0,
        price: 1000,
        vatAmount: 0,
        salesTaxAmount: 0,
        amountWithoutTaxes: 0,
        totalAmount: 0,
      };

      const updated = updateCatalogEntryCalculations(entry, 12);
      expect(updated.totalAmount).toBe(0);
    });
  });
});
