"use client";

import { useState, useCallback, useMemo } from "react";
import { Product, CartItem, CartSummary } from "@/types";
import { TAX_RATE } from "@/constants";

export function useCart() {
  const [items, setItems] = useState<CartItem[]>([]);

  const addItem = useCallback((product: Product) => {
    setItems((prev) => {
      const existingItem = prev.find((item) => item.id === product.id);
      if (existingItem) {
        return prev.map((item) =>
          item.id === product.id
            ? { ...item, quantity: item.quantity + 1 }
            : item,
        );
      }
      return [...prev, { ...product, quantity: 1 }];
    });
  }, []);

  const updateQuantity = useCallback(
    (productId: number | string, delta: number) => {
      setItems((prev) => {
        return prev
          .map((item) =>
            item.id === productId
              ? { ...item, quantity: item.quantity + delta }
              : item,
          )
          .filter((item) => item.quantity > 0);
      });
    },
    [],
  );

  const removeItem = useCallback((productId: number | string) => {
    setItems((prev) => prev.filter((item) => item.id !== productId));
  }, []);

  const clearCart = useCallback(() => {
    setItems([]);
  }, []);

  const summary: CartSummary = useMemo(() => {
    const subtotal = items.reduce(
      (sum, item) => sum + item.price * item.quantity,
      0,
    );
    const tax = subtotal * TAX_RATE;
    const total = subtotal + tax;

    return { subtotal, tax, total };
  }, [items]);

  return {
    items,
    summary,
    addItem,
    updateQuantity,
    removeItem,
    clearCart,
  };
}
