'use client';

import { Plus, Minus, Trash2 } from 'lucide-react';
import { CartItem } from '@/types';

interface CartProps {
  items: CartItem[];
  onUpdateQuantity: (productId: number, delta: number) => void;
  onRemoveItem: (productId: number) => void;
}

export default function Cart({ items, onUpdateQuantity, onRemoveItem }: CartProps) {
  const subtotal = items.reduce((sum, item) => sum + item.price * item.quantity, 0);
  const tax = subtotal * 0.2; // 20% НДС
  const total = subtotal + tax;

  return (
    <div className="flex flex-col h-full bg-gray-800 rounded-lg border border-gray-700 shadow-xl">
      {/* Header */}
      <div className="bg-gray-800 border-b border-gray-700 p-4 rounded-t-lg">
        <h2 className="text-xl font-bold text-white">Текущий чек</h2>
      </div>

      {/* Items List */}
      <div className="flex-1 overflow-y-auto p-4 space-y-3">
        {items.length === 0 ? (
          <div className="text-center text-gray-500 py-12">
            <p>Чек пуст</p>
            <p className="text-sm mt-2">Добавьте товары</p>
          </div>
        ) : (
          items.map((item) => (
            <div key={item.id} className="bg-white rounded-lg p-3 shadow-sm border border-gray-200">
              <div className="flex justify-between items-start mb-2">
                <h3 className="font-semibold text-gray-900 flex-1">{item.name}</h3>
                <button
                  onClick={() => onRemoveItem(item.id)}
                  className="text-red-500 hover:text-red-700 ml-2"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
              <div className="flex justify-between items-center">
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => onUpdateQuantity(item.id, -1)}
                    className="w-8 h-8 rounded-md bg-gray-200 hover:bg-gray-300 flex items-center justify-center"
                  >
                    <Minus className="w-4 h-4" />
                  </button>
                  <span className="w-8 text-center font-semibold">{item.quantity}</span>
                  <button
                    onClick={() => onUpdateQuantity(item.id, 1)}
                    className="w-8 h-8 rounded-md bg-blue-600 hover:bg-blue-700 text-white flex items-center justify-center"
                  >
                    <Plus className="w-4 h-4" />
                  </button>
                </div>
                <div className="text-right">
                  <p className="text-sm text-gray-500">{item.price.toFixed(2)} ₽ × {item.quantity}</p>
                  <p className="font-bold text-gray-900">{(item.price * item.quantity).toFixed(2)} ₽</p>
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Totals */}
      <div className="bg-gray-700 border-t border-gray-600 p-4 space-y-2 rounded-b-lg">
        <div className="flex justify-between text-gray-300">
          <span>Подытог:</span>
          <span className="font-semibold">{subtotal.toFixed(2)} ₽</span>
        </div>
        <div className="flex justify-between text-gray-300">
          <span>НДС (20%):</span>
          <span className="font-semibold">{tax.toFixed(2)} ₽</span>
        </div>
        <div className="flex justify-between text-xl font-bold text-white pt-2 border-t border-gray-600">
          <span>Итого:</span>
          <span className="text-emerald-400">{total.toFixed(2)} ₽</span>
        </div>
      </div>
    </div>
  );
}
