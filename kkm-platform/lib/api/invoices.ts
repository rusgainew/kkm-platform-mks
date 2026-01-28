import axios, { AxiosError } from "axios";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001";
const client = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
});

export interface InvoiceData {
  number: string;
  date: string;
  company_id: string;
  description?: string;
  items: {
    name: string;
    quantity: number;
    price: number;
  }[];
  notes?: string;
}

export interface Invoice extends InvoiceData {
  id: string;
  status: "draft" | "sent" | "paid" | "overdue";
  total: number;
  created_at: string;
  updated_at: string;
}

export async function listInvoices(token: string): Promise<Invoice[]> {
  try {
    const { data } = await client.get("/api/invoices", {
      headers: { Authorization: `Bearer ${token}` },
    });
    return data || [];
  } catch (error) {
    throw parseApiError(error);
  }
}

export async function getInvoiceById(
  id: string,
  token: string,
): Promise<Invoice> {
  try {
    const { data } = await client.get(`/api/invoices/${id}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    return data;
  } catch (error) {
    throw parseApiError(error);
  }
}

export async function createInvoice(
  invoiceData: InvoiceData,
  token: string,
): Promise<Invoice> {
  try {
    const { data } = await client.post("/api/invoices", invoiceData, {
      headers: { Authorization: `Bearer ${token}` },
    });
    return data;
  } catch (error) {
    throw parseApiError(error);
  }
}

export async function updateInvoice(
  id: string,
  invoiceData: Partial<InvoiceData>,
  token: string,
): Promise<Invoice> {
  try {
    const { data } = await client.put(`/api/invoices/${id}`, invoiceData, {
      headers: { Authorization: `Bearer ${token}` },
    });
    return data;
  } catch (error) {
    throw parseApiError(error);
  }
}

export async function deleteInvoice(id: string, token: string): Promise<void> {
  try {
    await client.delete(`/api/invoices/${id}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
  } catch (error) {
    throw parseApiError(error);
  }
}

export async function signInvoice(id: string, token: string): Promise<Invoice> {
  try {
    const { data } = await client.post(
      `/api/invoices/${id}/sign`,
      {},
      {
        headers: { Authorization: `Bearer ${token}` },
      },
    );
    return data;
  } catch (error) {
    throw parseApiError(error);
  }
}

export async function acceptInvoice(
  id: string,
  token: string,
): Promise<Invoice> {
  try {
    const { data } = await client.post(
      `/api/invoices/${id}/accept`,
      {},
      {
        headers: { Authorization: `Bearer ${token}` },
      },
    );
    return data;
  } catch (error) {
    throw parseApiError(error);
  }
}

export async function rejectInvoice(
  id: string,
  reason?: string,
  token?: string,
): Promise<Invoice> {
  try {
    const { data } = await client.post(
      `/api/invoices/${id}/reject`,
      { reason },
      {
        headers: { Authorization: `Bearer ${token}` },
      },
    );
    return data;
  } catch (error) {
    throw parseApiError(error);
  }
}

export async function queryInvoices(
  params: {
    number?: string;
    dateFrom?: string;
    dateTo?: string;
    status?: string;
    company_id?: string;
  },
  token: string,
): Promise<Invoice[]> {
  try {
    const { data } = await client.get("/api/invoices/query", {
      params,
      headers: { Authorization: `Bearer ${token}` },
    });
    return data || [];
  } catch (error) {
    throw parseApiError(error);
  }
}

function parseApiError(error: unknown): Error {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<{
      message?: string;
      error?: string;
    }>;
    const message =
      axiosError.response?.data?.message ||
      axiosError.response?.data?.error ||
      axiosError.message ||
      "Неизвестная ошибка";
    return new Error(message);
  }
  return error instanceof Error ? error : new Error("Неизвестная ошибка");
}
