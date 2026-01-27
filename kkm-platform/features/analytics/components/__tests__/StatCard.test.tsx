import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { StatCard } from "../StatCard";
import { TrendingUp } from "lucide-react";

describe("StatCard", () => {
  it("should render basic metrics", () => {
    render(
      <StatCard
        label="Выручка"
        value={1000000}
        format="currency"
        currency="KZT"
      />
    );

    expect(screen.getByText("Выручка")).toBeInTheDocument();
    // Check that formatted value is present (may be formatted differently)
    const values = screen.getAllByText(/1[\s,.]000[\s,.]000/);
    expect(values.length).toBeGreaterThan(0);
  });

  it("should render trend indicator for increase", () => {
    render(
      <StatCard
        label="Доход"
        value={5000}
        change={15.5}
        changeType="increase"
        format="number"
      />
    );

    expect(screen.getByText("Доход")).toBeInTheDocument();
    // Check for change percentage (may have + sign or not)
    expect(screen.getByText(/15\.5/)).toBeInTheDocument();
  });

  it("should render trend indicator for decrease", () => {
    render(
      <StatCard
        label="Расходы"
        value={3000}
        change={-10.2}
        changeType="decrease"
        format="number"
      />
    );

    expect(screen.getByText("Расходы")).toBeInTheDocument();
    expect(screen.getByText(/10\.2/)).toBeInTheDocument();
  });

  it("should render with custom icon", () => {
    const { container } = render(
      <StatCard
        label="Метрика"
        value={42}
        format="number"
        icon={<TrendingUp data-testid="custom-icon" />}
      />
    );

    expect(container.querySelector('[data-testid="custom-icon"]')).toBeInTheDocument();
  });

  it("should format percentage correctly", () => {
    render(
      <StatCard
        label="Конверсия"
        value={45.67}
        format="percentage"
      />
    );

    expect(screen.getByText("Конверсия")).toBeInTheDocument();
    expect(screen.getByText(/45\.7/)).toBeInTheDocument();
  });

  it("should handle zero values", () => {
    render(
      <StatCard
        label="Нулевой показатель"
        value={0}
        format="number"
      />
    );

    expect(screen.getByText("Нулевой показатель")).toBeInTheDocument();
    expect(screen.getByText("0")).toBeInTheDocument();
  });

  it("should handle neutral change", () => {
    render(
      <StatCard
        label="Стабильный"
        value={100}
        change={0}
        changeType="neutral"
        format="number"
      />
    );

    expect(screen.getByText("Стабильный")).toBeInTheDocument();
    expect(screen.getByText("100")).toBeInTheDocument();
    // Find all elements with "0.0%" and check at least one exists
    const changeTexts = screen.getAllByText((content, element) => {
      return element?.textContent === "0.0%";
    });
    expect(changeTexts.length).toBeGreaterThan(0);
  });
});
