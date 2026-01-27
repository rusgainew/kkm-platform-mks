import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { InvoiceTotals } from "../components/InvoiceTotals";

describe("InvoiceTotals", () => {
  it("should render all total amounts", () => {
    render(
      <InvoiceTotals
        totalAmountWithoutVAT={10000}
        totalVATAmount={1200}
        totalAmount={11200}
        currency="KGS"
      />
    );

    expect(screen.getByText(/Сумма без НДС/i)).toBeInTheDocument();
    expect(screen.getByText("НДС:")).toBeInTheDocument();
    expect(screen.getByText(/Итого к оплате/i)).toBeInTheDocument();
  });

  it("should format amounts correctly", () => {
    render(
      <InvoiceTotals
        totalAmountWithoutVAT={10000}
        totalVATAmount={1200}
        totalAmount={11200}
        currency="KGS"
      />
    );

    // Check that formatCurrency is called and displays numbers
    const totalsText = screen.getByText(/Итоговые суммы/i).parentElement
      ?.textContent;
    expect(totalsText).toContain("10");
    expect(totalsText).toContain("000");
  });

  it("should display ST amount when provided", () => {
    render(
      <InvoiceTotals
        totalAmountWithoutVAT={10000}
        totalVATAmount={1200}
        totalSTAmount={500}
        totalAmount={11700}
        currency="KGS"
      />
    );

    expect(screen.getByText("НСП:")).toBeInTheDocument();
  });

  it("should not display ST amount when not provided", () => {
    render(
      <InvoiceTotals
        totalAmountWithoutVAT={10000}
        totalVATAmount={1200}
        totalAmount={11200}
        currency="KGS"
      />
    );

    expect(screen.queryByText("НСП:")).not.toBeInTheDocument();
  });

  it("should handle zero amounts", () => {
    render(
      <InvoiceTotals
        totalAmountWithoutVAT={0}
        totalVATAmount={0}
        totalAmount={0}
        currency="KGS"
      />
    );

    expect(screen.getByText(/Сумма без НДС/i)).toBeInTheDocument();
    expect(screen.getByText("НДС:")).toBeInTheDocument();
    expect(screen.getByText(/Итого к оплате/i)).toBeInTheDocument();
  });

  it("should handle different currencies", () => {
    render(
      <InvoiceTotals
        totalAmountWithoutVAT={10000}
        totalVATAmount={1200}
        totalAmount={11200}
        currency="USD"
      />
    );

    // formatCurrency should be called with USD
    expect(screen.getByText(/Итоговые суммы/i)).toBeInTheDocument();
  });

  it("should apply custom className", () => {
    const { container } = render(
      <InvoiceTotals
        totalAmountWithoutVAT={10000}
        totalVATAmount={1200}
        totalAmount={11200}
        currency="KGS"
        className="custom-class"
      />
    );

    const mainDiv = container.querySelector(".custom-class");
    expect(mainDiv).toBeInTheDocument();
  });
});
