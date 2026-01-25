import { Product } from './product';

export interface CartItem extends Product {
  quantity: number;
}

export interface CartSummary {
  subtotal: number;
  tax: number;
  total: number;
}

export type PaymentMethod = 'card' | 'cash' | 'mixed';
