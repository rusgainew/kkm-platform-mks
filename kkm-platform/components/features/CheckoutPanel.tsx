'use client';

import { CreditCard, Banknote, Wallet } from 'lucide-react';
import { PaymentMethod } from '@/types';

interface CheckoutPanelProps {
  total: number;
  onCheckout: (paymentMethod: PaymentMethod) => void;
  disabled?: boolean;
}

export default function CheckoutPanel({ total, onCheckout, disabled = false }: CheckoutPanelProps) {
  const paymentMethods: Array<{
    id: PaymentMethod;
    label: string;
    icon: typeof CreditCard;
    color: string;
  }> = [
    { id: 'card', label: 'Оплата картой', icon: CreditCard, color: 'bg-linear-to-br from-blue-500 to-blue-600 hover:shadow-lg hover:shadow-blue-500/50' },
    { id: 'cash', label: 'Наличные', icon: Banknote, color: 'bg-linear-to-br from-emerald-500 to-emerald-600 hover:shadow-lg hover:shadow-emerald-500/50' },
    { id: 'mixed', label: 'Комбинированная', icon: Wallet, color: 'bg-linear-to-br from-purple-500 to-purple-600 hover:shadow-lg hover:shadow-purple-500/50' },
  ];

  return (
    <div className="bg-gray-800 border-t border-gray-700 p-6 shadow-xl">
      <div className="grid grid-cols-3 gap-4 max-w-4xl mx-auto">
        {paymentMethods.map((method) => {
          const Icon = method.icon;
          return (
            <button
              key={method.id}
              onClick={() => onCheckout(method.id)}
              disabled={disabled}
              className={`${method.color} text-white rounded-lg p-6 flex flex-col items-center justify-center gap-3 transition-all hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed`}
            >
              <Icon className="w-8 h-8" />
              <span className="font-bold text-lg">{method.label}</span>
              {!disabled && (
                <span className="text-sm opacity-90">{total.toFixed(2)} ₽</span>
              )}
            </button>
          );
        })}
      </div>
    </div>
  );
}
