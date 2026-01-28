/**
 * Analytics API Client
 * Клиент для работы с API аналитики счетов-фактур
 */

import axios, { AxiosInstance } from "axios";
import type {
  InvoiceAnalyticsData,
  AnalyticsFilters,
  GetAnalyticsDataResponse,
  InvoiceAnalyticsStats,
  SalesChartData,
} from "@/types/analytics";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

export class AnalyticsAPIClient {
  private axiosInstance: AxiosInstance;

  constructor() {
    this.axiosInstance = axios.create({
      baseURL: `${API_BASE_URL}/analytics`,
      headers: {
        "Content-Type": "application/json",
      },
    });

    // Добавляем JWT токен к каждому запросу
    this.axiosInstance.interceptors.request.use((config) => {
      const token = localStorage.getItem("token");
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });
  }

  /**
   * Получить все данные аналитики
   */
  async getAnalyticsData(
    filters: AnalyticsFilters,
  ): Promise<InvoiceAnalyticsData> {
    const params = this.buildQueryParams(filters);
    const response = await this.axiosInstance.get<GetAnalyticsDataResponse>(
      "/",
      { params },
    );

    if (!response.data.success) {
      throw new Error(response.data.error || "Failed to fetch analytics data");
    }

    return response.data.data;
  }

  /**
   * Получить статистику по счетам-фактурам
   */
  async getStats(filters: AnalyticsFilters): Promise<InvoiceAnalyticsStats> {
    const params = this.buildQueryParams(filters);
    const response = await this.axiosInstance.get<{
      success: boolean;
      data: InvoiceAnalyticsStats;
    }>("/stats", { params });

    if (!response.data.success) {
      throw new Error("Failed to fetch stats");
    }

    return response.data.data;
  }

  /**
   * Получить данные графика продаж
   */
  async getSalesChart(filters: AnalyticsFilters): Promise<SalesChartData> {
    const params = this.buildQueryParams(filters);
    const response = await this.axiosInstance.get<{
      success: boolean;
      data: SalesChartData;
    }>("/sales-chart", { params });

    if (!response.data.success) {
      throw new Error("Failed to fetch sales chart data");
    }

    return response.data.data;
  }

  /**
   * Получить статистику по статусам
   */
  async getStatusStats(filters: AnalyticsFilters) {
    const params = this.buildQueryParams(filters);
    const response = await this.axiosInstance.get("/status-stats", { params });
    return response.data;
  }

  /**
   * Получить статистику по типам операций
   */
  async getOperationTypeStats(filters: AnalyticsFilters) {
    const params = this.buildQueryParams(filters);
    const response = await this.axiosInstance.get("/operation-type-stats", {
      params,
    });
    return response.data;
  }

  /**
   * Получить топ контрагентов
   */
  async getTopContractors(filters: AnalyticsFilters, limit: number = 10) {
    const params = { ...this.buildQueryParams(filters), limit };
    const response = await this.axiosInstance.get("/top-contractors", {
      params,
    });
    return response.data;
  }

  /**
   * Получить выручку по месяцам
   */
  async getRevenueByMonth(filters: AnalyticsFilters) {
    const params = this.buildQueryParams(filters);
    const response = await this.axiosInstance.get("/revenue-by-month", {
      params,
    });
    return response.data;
  }

  /**
   * Построить query параметры из фильтров
   */
  private buildQueryParams(filters: AnalyticsFilters): Record<string, string> {
    const params: Record<string, string> = {
      period: filters.period,
    };

    if (filters.dateRange) {
      params.startDate = filters.dateRange.startDate;
      params.endDate = filters.dateRange.endDate;
    }

    if (filters.operationType) {
      params.operationType = filters.operationType;
    }

    if (filters.status) {
      params.status = filters.status;
    }

    if (filters.currency) {
      params.currency = filters.currency;
    }

    return params;
  }
}

// Экспортируем singleton instance
export const analyticsAPI = new AnalyticsAPIClient();
