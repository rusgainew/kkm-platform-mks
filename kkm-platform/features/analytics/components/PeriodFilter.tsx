"use client";

import { useState } from "react";
import { Calendar } from "lucide-react";
import type { TimePeriod, DateRange } from "@/types/analytics";

export interface PeriodFilterProps {
  selectedPeriod: TimePeriod;
  onPeriodChange: (period: TimePeriod, dateRange?: DateRange) => void;
  className?: string;
}

/**
 * Компонент фильтра периодов времени
 * Позволяет выбрать предустановленный период или указать custom диапазон
 */
export function PeriodFilter({
  selectedPeriod,
  onPeriodChange,
  className = "",
}: PeriodFilterProps) {
  const [showCustomRange, setShowCustomRange] = useState(false);
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");

  const periods: { value: TimePeriod; label: string }[] = [
    { value: "today", label: "Сегодня" },
    { value: "week", label: "Неделя" },
    { value: "month", label: "Месяц" },
    { value: "quarter", label: "Квартал" },
    { value: "year", label: "Год" },
    { value: "custom", label: "Свой период" },
  ];

  const handlePeriodClick = (period: TimePeriod) => {
    if (period === "custom") {
      setShowCustomRange(true);
    } else {
      setShowCustomRange(false);
      onPeriodChange(period);
    }
  };

  const handleCustomRangeApply = () => {
    if (startDate && endDate) {
      onPeriodChange("custom", { startDate, endDate });
      setShowCustomRange(false);
    }
  };

  return (
    <div className={`bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 ${className}`}>
      <div className="flex items-center gap-2 mb-3">
        <Calendar className="w-5 h-5 text-gray-500 dark:text-gray-400" />
        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
          Период
        </span>
      </div>

      {/* Quick filters */}
      <div className="flex flex-wrap gap-2 mb-4">
        {periods.map((period) => (
          <button
            key={period.value}
            onClick={() => handlePeriodClick(period.value)}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              selectedPeriod === period.value
                ? "bg-blue-600 text-white"
                : "bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600"
            }`}
          >
            {period.label}
          </button>
        ))}
      </div>

      {/* Custom date range */}
      {showCustomRange && (
        <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label htmlFor="start-date" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Начало
              </label>
              <input
                id="start-date"
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-700 dark:text-gray-100"
              />
            </div>
            <div>
              <label htmlFor="end-date" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Конец
              </label>
              <input
                id="end-date"
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-700 dark:text-gray-100"
              />
            </div>
          </div>
          <div className="mt-4 flex gap-2">
            <button
              onClick={handleCustomRangeApply}
              disabled={!startDate || !endDate}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Применить
            </button>
            <button
              onClick={() => setShowCustomRange(false)}
              className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg font-medium hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
            >
              Отмена
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
