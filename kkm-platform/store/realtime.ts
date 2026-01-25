import { create } from "zustand";
import { WebSocketMessage } from "@/lib/hooks/useWebSocket";

export interface RealtimeInvoice {
  id: string;
  number: string;
  status: "draft" | "issued" | "paid" | "cancelled";
  amount: number;
  updatedAt: string;
  company?: string;
}

export interface RealtimeProduct {
  id: string;
  name: string;
  sku: string;
  stock: number;
  price: number;
  updatedAt: string;
}

interface RealtimeInvoiceStore {
  invoices: Map<string, RealtimeInvoice>;
  addOrUpdateInvoice: (invoice: RealtimeInvoice) => void;
  deleteInvoice: (id: string) => void;
  getInvoice: (id: string) => RealtimeInvoice | undefined;
  getAllInvoices: () => RealtimeInvoice[];
  syncFromMessage: (message: WebSocketMessage) => void;
}

interface RealtimeCatalogStore {
  products: Map<string, RealtimeProduct>;
  addOrUpdateProduct: (product: RealtimeProduct) => void;
  deleteProduct: (id: string) => void;
  getProduct: (id: string) => RealtimeProduct | undefined;
  getAllProducts: () => RealtimeProduct[];
  syncFromMessage: (message: WebSocketMessage) => void;
}

export const useRealtimeInvoices = create<RealtimeInvoiceStore>((set, get) => ({
  invoices: new Map(),

  addOrUpdateInvoice: (invoice: RealtimeInvoice) => {
    set((state) => {
      const newInvoices = new Map(state.invoices);
      newInvoices.set(invoice.id, invoice);
      return { invoices: newInvoices };
    });
  },

  deleteInvoice: (id: string) => {
    set((state) => {
      const newInvoices = new Map(state.invoices);
      newInvoices.delete(id);
      return { invoices: newInvoices };
    });
  },

  getInvoice: (id: string) => {
    return get().invoices.get(id);
  },

  getAllInvoices: () => {
    return Array.from(get().invoices.values());
  },

  syncFromMessage: (message: WebSocketMessage) => {
    if (message.type !== "invoice") return;

    if (message.action === "update" || message.action === "create") {
      const invoice = message.data as RealtimeInvoice;
      get().addOrUpdateInvoice(invoice);
    } else if (message.action === "delete") {
      get().deleteInvoice(message.id || "");
    }
  },
}));

export const useRealtimeCatalog = create<RealtimeCatalogStore>((set, get) => ({
  products: new Map(),

  addOrUpdateProduct: (product: RealtimeProduct) => {
    set((state) => {
      const newProducts = new Map(state.products);
      newProducts.set(product.id, product);
      return { products: newProducts };
    });
  },

  deleteProduct: (id: string) => {
    set((state) => {
      const newProducts = new Map(state.products);
      newProducts.delete(id);
      return { products: newProducts };
    });
  },

  getProduct: (id: string) => {
    return get().products.get(id);
  },

  getAllProducts: () => {
    return Array.from(get().products.values());
  },

  syncFromMessage: (message: WebSocketMessage) => {
    if (message.type !== "catalog") return;

    if (message.action === "update" || message.action === "create") {
      const product = message.data as RealtimeProduct;
      get().addOrUpdateProduct(product);
    } else if (message.action === "delete") {
      get().deleteProduct(message.id || "");
    }
  },
}));
