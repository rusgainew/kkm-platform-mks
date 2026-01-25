// Export all stores for easy access
export { default as useAuthStore } from './auth/authStore';
export type { User } from './auth/authStore';

export { default as useUserStore } from './user/userStore';
export type { UserListState } from './user/userStore';

export { default as useInvoiceStore } from './invoice/invoiceStore';
export type { Invoice, InvoiceItem } from './invoice/invoiceStore';

export { default as useCatalogStore } from './catalog/catalogStore';
export type { CatalogItem } from './catalog/catalogStore';

export { default as useUIStore } from './ui/uiStore';
export type { Notification } from './ui/uiStore';
