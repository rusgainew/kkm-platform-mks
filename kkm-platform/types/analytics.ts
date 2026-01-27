/**
 * Invoice Analytics and Dashboard Types
 * Типы для аналитики счетов-фактур
 */

/**
 * Период времени для фильтрации данных
 */
export type TimePeriod =
  | "today"
  | "week"
  | "month"
  | "quarter"
  | "year"
  | "custom";

/**
 * Временной диапазон
 */
export interface DateRange {
  startDate: string; // ISO date string
  endDate: string; // ISO date string
}

/**
 * Метрика дашборда
 */
export interface DashboardMetric {
  label: string;
  value: number;
  previousValue?: number;
  change?: number; // Процент изменения
  changeType?: "increase" | "decrease" | "neutral";
  format?: "currency" | "number" | "percentage";
  icon?: string;
  currency?: string;
}

/**
 * Статистика по счетам-фактурам
 */
export interface InvoiceAnalyticsStats {
  totalRevenue: number;
  totalInvoices: number;
  averageInvoiceAmount: number;
  activeContractors: number;
  pendingInvoices: number;
  approvedInvoices: number;
  rejectedInvoices: number;
  revenueChange?: number; // Процент изменения выручки
  invoiceCountChange?: number; // Процент изменения количества
  averageAmountChange?: number; // Процент изменения среднего чека
  contractorsChange?: number; // Процент изменения количества контрагентов
  period: DateRange;
}

/**
 * Точка данных для графика
 */
export interface ChartDataPoint {
  date: string; // ISO date or formatted date
  value: number;
  label?: string;
  category?: string;
}

/**
 * Данные графика продаж
 */
export interface SalesChartData {
  data: ChartDataPoint[];
  period: DateRange;
}

/**
 * Данные для круговой диаграммы
 */
export interface PieChartData {
  name: string;
  value: number;
  percentage: number;
  color?: string;
}

/**
 * Статистика по статусам счетов
 */
export interface InvoiceStatusStats {
  pending: number;
  approved: number;
  rejected: number;
  draft: number;
  total: number;
}

/**
 * Статистика по типам операций
 */
export interface OperationTypeStats {
  operationTypeCode: string;
  operationTypeName: string;
  count: number;
  totalAmount: number;
  percentage: number;
}

/**
 * Топ контрагенты
 */
export interface TopContractor {
  tin: string;
  name: string;
  invoiceCount: number;
  totalAmount: number;
  lastInvoiceDate: string;
}

/**
 * Тренд данных
 */
export interface TrendData {
  current: number;
  previous: number;
  change: number;
  changePercentage: number;
  trend: "up" | "down" | "stable";
}

/**
 * Полные данные аналитики
 */
export interface InvoiceAnalyticsData {
  stats: InvoiceAnalyticsStats;
  salesChart: SalesChartData;
  invoiceStatusStats: InvoiceStatusStats;
  operationTypeStats: OperationTypeStats[];
  topContractors: TopContractor[];
  revenueByMonth: ChartDataPoint[];
  invoicesByStatus: PieChartData[];
  period: DateRange;
  lastUpdated: string;
}

/**
 * Фильтры аналитики
 */
export interface AnalyticsFilters {
  period: TimePeriod;
  dateRange?: DateRange;
  operationType?: string;
  status?: string;
  currency?: string;
}

/**
 * Запрос данных аналитики
 */
export interface GetAnalyticsDataRequest {
  filters: AnalyticsFilters;
}

/**
 * Ответ с данными аналитики
 */
export interface GetAnalyticsDataResponse {
  success: boolean;
  data: InvoiceAnalyticsData;
  error?: string;
}
