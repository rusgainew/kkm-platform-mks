import { describe, it, expect } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { PeriodFilter } from "../PeriodFilter";
import type { TimePeriod, DateRange } from "@/types/analytics";

describe("PeriodFilter", () => {
  it("should render all period buttons", () => {
    const mockOnChange = () => {};
    render(
      <PeriodFilter
        selectedPeriod="month"
        onPeriodChange={mockOnChange}
      />
    );

    expect(screen.getByText("Сегодня")).toBeInTheDocument();
    expect(screen.getByText("Неделя")).toBeInTheDocument();
    expect(screen.getByText("Месяц")).toBeInTheDocument();
    expect(screen.getByText("Квартал")).toBeInTheDocument();
    expect(screen.getByText("Год")).toBeInTheDocument();
    expect(screen.getByText("Свой период")).toBeInTheDocument();
  });

  it("should highlight selected period", () => {
    const mockOnChange = () => {};
    render(
      <PeriodFilter
        selectedPeriod="week"
        onPeriodChange={mockOnChange}
      />
    );

    const weekButton = screen.getByText("Неделя");
    expect(weekButton.className).toContain("bg-blue-600");
  });

  it("should call onChange when period is clicked", () => {
    let selectedPeriod: TimePeriod = "month";
    const mockOnChange = (period: TimePeriod) => {
      selectedPeriod = period;
    };

    const { rerender } = render(
      <PeriodFilter
        selectedPeriod={selectedPeriod}
        onPeriodChange={mockOnChange}
      />
    );

    const todayButton = screen.getByText("Сегодня");
    fireEvent.click(todayButton);

    expect(selectedPeriod).toBe("today");

    rerender(
      <PeriodFilter
        selectedPeriod={selectedPeriod}
        onPeriodChange={mockOnChange}
      />
    );

    expect(todayButton.className).toContain("bg-blue-600");
  });

  it("should show custom date range inputs when 'Свой период' is clicked", () => {
    const mockOnChange = () => {};
    render(
      <PeriodFilter
        selectedPeriod="month"
        onPeriodChange={mockOnChange}
      />
    );

    const customButton = screen.getByText("Свой период");
    fireEvent.click(customButton);

    expect(screen.getByText("Начало")).toBeInTheDocument();
    expect(screen.getByText("Конец")).toBeInTheDocument();
    expect(screen.getByText("Применить")).toBeInTheDocument();
    expect(screen.getByText("Отмена")).toBeInTheDocument();
  });

  it("should call onChange with date range when custom range is applied", () => {
    let result: { period?: TimePeriod; range?: DateRange } = {};
    const mockOnChange = (period: TimePeriod, range?: DateRange) => {
      result = { period, range };
    };

    render(
      <PeriodFilter
        selectedPeriod="month"
        onPeriodChange={mockOnChange}
      />
    );

    // Open custom range
    const customButton = screen.getByText("Свой период");
    fireEvent.click(customButton);

    // Fill in dates
    const startInput = screen.getByLabelText("Начало");
    const endInput = screen.getByLabelText("Конец");
    
    fireEvent.change(startInput, { target: { value: "2024-01-01" } });
    fireEvent.change(endInput, { target: { value: "2024-01-31" } });

    // Apply
    const applyButton = screen.getByText("Применить");
    fireEvent.click(applyButton);

    expect(result.period).toBe("custom");
    expect(result.range).toEqual({
      startDate: "2024-01-01",
      endDate: "2024-01-31",
    });
  });

  it("should hide custom range when cancel is clicked", () => {
    const mockOnChange = () => {};
    render(
      <PeriodFilter
        selectedPeriod="month"
        onPeriodChange={mockOnChange}
      />
    );

    // Open custom range
    const customButton = screen.getByText("Свой период");
    fireEvent.click(customButton);

    expect(screen.getByText("Начало")).toBeInTheDocument();

    // Click cancel
    const cancelButton = screen.getByText("Отмена");
    fireEvent.click(cancelButton);

    expect(screen.queryByText("Начало")).not.toBeInTheDocument();
  });

  it("should disable apply button when dates are not filled", () => {
    const mockOnChange = () => {};
    render(
      <PeriodFilter
        selectedPeriod="month"
        onPeriodChange={mockOnChange}
      />
    );

    // Open custom range
    const customButton = screen.getByText("Свой период");
    fireEvent.click(customButton);

    const applyButton = screen.getByText("Применить") as HTMLButtonElement;
    expect(applyButton.disabled).toBe(true);

    // Fill start date only
    const startInput = screen.getByLabelText("Начало");
    fireEvent.change(startInput, { target: { value: "2024-01-01" } });

    expect(applyButton.disabled).toBe(true);

    // Fill end date
    const endInput = screen.getByLabelText("Конец");
    fireEvent.change(endInput, { target: { value: "2024-01-31" } });

    expect(applyButton.disabled).toBe(false);
  });
});
