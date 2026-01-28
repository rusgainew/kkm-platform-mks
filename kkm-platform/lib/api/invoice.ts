/**
 * ESF API Invoice Client
 * Клиент для работы со счетами-фактурами через ESF API
 *
 * Базовый URL: http://localhost:8080/api/invoices
 * Все операции требуют JWT токен в заголовке Authorization
 */

import axios, { AxiosError, AxiosInstance } from "axios";
import type { Invoice, InvoiceDetail } from "@/types/entities";
import type { ListResponse } from "@/types/api-response";
import type {
  CreateInvoiceRequest,
  UpdateInvoiceRequest,
  AcceptOrRejectInvoiceRequest,
  SignInvoiceRequest,
  RevokeInvoiceRequest,
  InvoiceFilters,
  InvoiceOperationResponse,
} from "@/types";
import type {
  ESFAPIResponse,
  ESFCreateInvoiceResponse,
  ESFInvoiceResponse,
  ESFInvoiceListResponse,
  ESFInvoiceActionResponse,
} from "@/types/api-response";

/**
 * Invoice API Error
 */
export class InvoiceAPIError extends Error {
  constructor(
    public statusCode: number,
    public message: string,
    public details?: Record<string, unknown>,
  ) {
    super(message);
    this.name = "InvoiceAPIError";
  }
}

/**
 * Invoice API Client
 * Типизированный Axios клиент для ESF Invoice API
 */
class InvoiceAPIClient {
  private axiosInstance: AxiosInstance;
  private baseURL: string;
  private token: string | null = null;
  private retryAttempts = 3;
  private retryDelay = 1000;

  constructor(baseURL = "http://localhost:8080/api") {
    this.baseURL = baseURL;

    this.axiosInstance = axios.create({
      baseURL: `${baseURL}/invoices`,
      timeout: 30000,
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    });

    // Request interceptor - добавляет токен
    this.axiosInstance.interceptors.request.use(
      (config) => {
        if (this.token) {
          config.headers.Authorization = `Bearer ${this.token}`;
        }
        return config;
      },
      (error) => Promise.reject(error),
    );

    // Response interceptor - обработка ошибок
    this.axiosInstance.interceptors.response.use(
      (response) => response,
      (error) => this.handleError(error),
    );
  }

  /**
   * Установить JWT токен для аутентификации
   */
  setToken(token: string): void {
    this.token = token;
  }

  /**
   * Очистить токен
   */
  clearToken(): void {
    this.token = null;
  }

  /**
   * Обработка ошибок с retry логикой
   */
  private async handleError(error: AxiosError): Promise<never> {
    if (error.response?.status === 401) {
      this.clearToken();
      throw new InvoiceAPIError(401, "Unauthorized - token expired");
    }

    const message = (error.response?.data as any)?.message || error.message;
    const statusCode = error.response?.status || 500;

    throw new InvoiceAPIError(
      statusCode,
      message,
      (error.response?.data as Record<string, unknown>) || undefined,
    );
  }

  /**
   * Создать новый счет-фактуру
   *
   * @param request - Данные для создания счета
   * @returns Ответ с UUID созданного счета
   */
  async createInvoice(request: CreateInvoiceRequest): Promise<Invoice> {
    try {
      const response = await this.axiosInstance.post<ESFCreateInvoiceResponse>(
        "/",
        request,
      );

      if (!response.data.data?.documentUuid) {
        throw new InvoiceAPIError(400, "Invalid response - no documentUuid");
      }

      // Преобразовать ответ в Invoice тип
      return {
        id: response.data.data.documentUuid,
        invoice_number:
          response.data.data.invoiceNumber || request.invoiceNumber || "",
        invoice_date:
          request.invoiceDate || new Date().toISOString().split("T")[0],
        delivery_date: request.deliveryDate,
        total_amount: 0,
        is_resident: request.isResident,
        note: request.comment, // Используем comment вместо note
        status: String(response.data.data.status || "draft"), // Приводим к string
        created_at: Date.now(),
        updated_at: Date.now(),
      };
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to create invoice", { error });
    }
  }

  /**
   * Получить счет-фактуру по UUID
   *
   * @param invoiceUuid - Уникальный идентификатор счета
   * @returns Полные данные счета
   */
  async getInvoice(invoiceUuid: string): Promise<Invoice> {
    try {
      const response = await this.axiosInstance.get<ESFInvoiceResponse>(
        `/${invoiceUuid}`,
      );

      if (!response.data.data) {
        throw new InvoiceAPIError(404, "Invoice not found");
      }

      const data = response.data.data;
      return {
        documentUuid: data.documentUuid,
        invoiceNumber: data.invoiceNumber || "",
        status: {
          code: data.status,
          name: data.status,
        },
        totalAmount: data.totalAmount || 0,
        isResident: (data as any).isResident ?? true,
        legalPerson: (data as any).legalPerson || { pin: "", fullName: "" },
        contractor: (data as any).contractor || { pin: "", fullName: "" },
        ...(data as any),
      };
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to get invoice", { error });
    }
  }

  /**
   * Получить список счетов-фактур с фильтрацией и пагинацией
   *
   * @param filters - Фильтры для поиска
   * @returns Список счетов с пагинацией
   */
  async listInvoices(filters?: InvoiceFilters): Promise<ListResponse<Invoice>> {
    try {
      const params = {
        status: filters?.status,
        invoiceNumber: filters?.invoiceNumber,
        contractorTin: filters?.contractorTin,
        dateFrom: filters?.dateFrom,
        dateTo: filters?.dateTo,
        currency: filters?.currency,
        page: filters?.page || 1,
        limit: Math.min(filters?.limit || 20, 100), // Max 100
      };

      // Удалить undefined параметры
      Object.keys(params).forEach(
        (key) =>
          (params as any)[key] === undefined && delete (params as any)[key],
      );

      const response = await this.axiosInstance.get<ESFInvoiceListResponse>(
        "/",
        { params },
      );

      const invoices: Invoice[] = (response.data.data || []).map(
        (data: any) => ({
          documentUuid: data.documentUuid,
          invoiceNumber: data.invoiceNumber || "",
          status: {
            code: data.status,
            name: data.status,
          },
          totalAmount: data.totalAmount || 0,
          isResident: data.isResident ?? true,
          legalPerson: data.legalPerson || { pin: "", fullName: "" },
          contractor: data.contractor || { pin: "", fullName: "" },
          ...data,
        }),
      );

      return {
        items: invoices,
        total: (response.data.data || []).length,
        page: params.page,
        page_size: params.limit,
      };
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to list invoices", { error });
    }
  }

  /**
   * Обновить существующий счет-фактуру
   *
   * @param request - Данные для обновления (должен содержать documentUuid)
   * @returns Обновленные данные счета
   */
  async updateInvoice(request: UpdateInvoiceRequest): Promise<Invoice> {
    try {
      if (!request.documentUuid) {
        throw new InvoiceAPIError(400, "documentUuid is required");
      }

      const response = await this.axiosInstance.put<ESFInvoiceResponse>(
        `/${request.documentUuid}`,
        request,
      );

      if (!response.data.data) {
        throw new InvoiceAPIError(400, "Invalid response");
      }

      const data = response.data.data;
      return {
        documentUuid: data.documentUuid,
        invoiceNumber: data.invoiceNumber || "",
        status: {
          code: data.status,
          name: data.status,
        },
        totalAmount: data.totalAmount || 0,
        isResident: (data as any).isResident ?? true,
        legalPerson: (data as any).legalPerson || { pin: "", fullName: "" },
        contractor: (data as any).contractor || { pin: "", fullName: "" },
        ...(data as any),
      };
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to update invoice", { error });
    }
  }

  /**
   * Принять счет-фактуру
   *
   * @param invoiceUuid - UUID счета
   * @param comment - Опциональный комментарий
   * @returns Результат операции
   */
  async acceptInvoice(
    invoiceUuid: string,
    comment?: string,
  ): Promise<InvoiceOperationResponse> {
    try {
      const request: AcceptOrRejectInvoiceRequest = {
        invoiceUuid,
        action: "accept",
        comment,
      };

      const response = await this.axiosInstance.post<ESFInvoiceActionResponse>(
        `/${invoiceUuid}/accept`,
        request,
      );

      return {
        success: true,
        invoiceUuid: response.data.data?.documentUuid,
        status: response.data.data?.status,
        message: "Invoice accepted successfully",
        responseId: response.data.responseId,
      };
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to accept invoice", { error });
    }
  }

  /**
   * Отклонить счет-фактуру
   *
   * @param invoiceUuid - UUID счета
   * @param reason - Причина отклонения (обязательно)
   * @param comment - Опциональный комментарий
   * @returns Результат операции
   */
  async rejectInvoice(
    invoiceUuid: string,
    reason: string,
    comment?: string,
  ): Promise<InvoiceOperationResponse> {
    try {
      if (!reason) {
        throw new InvoiceAPIError(400, "Reason is required for rejection");
      }

      const request: AcceptOrRejectInvoiceRequest = {
        invoiceUuid,
        action: "reject",
        reason,
        comment,
      };

      const response = await this.axiosInstance.post<ESFInvoiceActionResponse>(
        `/${invoiceUuid}/reject`,
        request,
      );

      return {
        success: true,
        invoiceUuid: response.data.data?.documentUuid,
        status: response.data.data?.status,
        message: "Invoice rejected successfully",
        responseId: response.data.responseId,
      };
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to reject invoice", { error });
    }
  }

  /**
   * Подписать счет-фактуру электронной подписью
   *
   * @param invoiceUuid - UUID счета
   * @param signature - Base64-кодированная подпись
   * @param certificateData - Данные сертификата (опционально)
   * @returns Результат операции
   */
  async signInvoice(
    invoiceUuid: string,
    signature: string,
    certificateData?: string,
  ): Promise<InvoiceOperationResponse> {
    try {
      if (!signature) {
        throw new InvoiceAPIError(400, "Signature is required");
      }

      const request: SignInvoiceRequest = {
        invoiceUuid,
        signature,
        certificateData,
        timestamp: Date.now(),
      };

      const response = await this.axiosInstance.post<ESFInvoiceActionResponse>(
        `/${invoiceUuid}/sign`,
        request,
      );

      return {
        success: true,
        invoiceUuid: response.data.data?.documentUuid,
        status: response.data.data?.status,
        message: "Invoice signed successfully",
        responseId: response.data.responseId,
      };
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to sign invoice", { error });
    }
  }

  /**
   * Отозвать (отменить) счет-фактуру
   *
   * @param invoiceUuid - UUID счета
   * @param reason - Причина отзыва (обязательно)
   * @param comment - Опциональный комментарий
   * @returns Результат операции
   */
  async revokeInvoice(
    invoiceUuid: string,
    reason: string,
    comment?: string,
  ): Promise<InvoiceOperationResponse> {
    try {
      if (!reason) {
        throw new InvoiceAPIError(400, "Reason is required for revocation");
      }

      const request: RevokeInvoiceRequest = {
        invoiceUuid,
        reason,
        comment,
      };

      const response = await this.axiosInstance.post<ESFInvoiceActionResponse>(
        `/${invoiceUuid}/revoke`,
        request,
      );

      return {
        success: true,
        invoiceUuid: response.data.data?.documentUuid,
        status: response.data.data?.status,
        message: "Invoice revoked successfully",
        responseId: response.data.responseId,
      };
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to revoke invoice", { error });
    }
  }

  /**
   * Получить детали счета-фактуры
   *
   * @param invoiceUuid - UUID счета
   * @returns Массив деталей (строк) счета
   */
  async getInvoiceDetails(invoiceUuid: string): Promise<InvoiceDetail[]> {
    try {
      const response = await this.axiosInstance.get<{
        data: InvoiceDetail[];
      }>(`/${invoiceUuid}/details`);

      return response.data.data || [];
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to get invoice details", {
        error,
      });
    }
  }

  /**
   * Скачать счет-фактуру в формате PDF
   *
   * @param invoiceUuid - UUID счета
   * @returns URL для скачивания или blob PDF
   */
  async downloadInvoicePDF(invoiceUuid: string): Promise<Blob> {
    try {
      const response = await this.axiosInstance.get(
        `/${invoiceUuid}/download/pdf`,
        {
          responseType: "blob",
        },
      );

      return response.data as Blob;
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to download invoice PDF", {
        error,
      });
    }
  }

  /**
   * Скачать счет-фактуру в формате XML
   *
   * @param invoiceUuid - UUID счета
   * @returns XML string или blob
   */
  async downloadInvoiceXML(invoiceUuid: string): Promise<Blob> {
    try {
      const response = await this.axiosInstance.get(
        `/${invoiceUuid}/download/xml`,
        {
          responseType: "blob",
        },
      );

      return response.data as Blob;
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to download invoice XML", {
        error,
      });
    }
  }

  /**
   * Получить статусы счетов для карточки статуса
   *
   * @param invoiceUuid - UUID счета
   * @returns История статусов счета
   */
  async getInvoiceStatusHistory(
    invoiceUuid: string,
  ): Promise<Array<{ status: string; changedAt: string; changedBy: string }>> {
    try {
      const response = await this.axiosInstance.get<{
        data: Array<{ status: string; changedAt: string; changedBy: string }>;
      }>(`/${invoiceUuid}/status-history`);

      return response.data.data || [];
    } catch (error) {
      if (error instanceof InvoiceAPIError) throw error;
      throw new InvoiceAPIError(500, "Failed to get invoice status history", {
        error,
      });
    }
  }
}

/**
 * Экспортированный singleton instance
 * Используется во всем приложении
 */
const invoiceAPI = new InvoiceAPIClient();

export { InvoiceAPIClient, invoiceAPI };

/**
 * Экспортированные типы функций для удобства
 */
export type {
  Invoice,
  InvoiceDetail,
  CreateInvoiceRequest,
  UpdateInvoiceRequest,
  InvoiceFilters,
  InvoiceOperationResponse,
};
