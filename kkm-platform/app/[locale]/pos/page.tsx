'use client';

import Sidebar from '@/components/layout/Sidebar';
import ProductGrid from '@/components/features/ProductGrid';
import Cart from '@/components/features/Cart';
import CheckoutPanel from '@/components/features/CheckoutPanel';
import { useCart, useOnlineStatus } from '@/hooks';
import { MOCK_PRODUCTS } from '@/constants';
import { PaymentMethod } from '@/types';
import ProtectedRoute from '@/components/auth/ProtectedRoute';
import { ErrorBoundary, POSErrorFallback } from '@/components/error';

function POSContent() {
  const { isOnline, toggleOnlineStatus } = useOnlineStatus();
  const { items, summary, addItem, updateQuantity, removeItem, clearCart } = useCart();

  const handleCheckout = (paymentMethod: PaymentMethod) => {
    alert(
      `Оплата: ${paymentMethod}\nСумма: ${summary.total.toFixed(2)} ₽\n\nЧек успешно оформлен!`
    );
    clearCart();
  };

  return (
    <div className="flex h-screen bg-gray-900">
      {/* Sidebar */}
      <Sidebar isOnline={isOnline} />

      {/* Main Content */}
      <div className="flex-1 flex flex-col">
        {/* Top Bar */}
        <div className="bg-gray-800 border-b border-gray-700 px-6 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-white">Касса POS</h1>
          <div className="flex items-center gap-3">
            <div
              className={`px-3 py-1 rounded-full text-sm font-semibold ${
                isOnline
                  ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/50'
                  : 'bg-red-500/20 text-red-400 border border-red-500/50'
              }`}
            >
              {isOnline ? 'Online' : 'Offline'}
            </div>
            <button
              onClick={toggleOnlineStatus}
              className="text-sm text-blue-400 hover:text-blue-300 transition-colors"
            >
              Переключить
            </button>
          </div>
        </div>

        {/* Main Grid */}
        <div className="flex-1 overflow-hidden flex">
          {/* Left Column - Products */}
          <div className="flex-6 p-6 overflow-hidden">
            <ProductGrid products={MOCK_PRODUCTS} onAddToCart={addItem} />
          </div>

          {/* Right Column - Cart */}
          <div className="flex-4 p-6 flex flex-col">
            <Cart
              items={items}
              onUpdateQuantity={updateQuantity}
              onRemoveItem={removeItem}
            />
          </div>
        </div>

        {/* Checkout Panel */}
        <CheckoutPanel
          total={summary.total}
          onCheckout={handleCheckout}
          disabled={items.length === 0}
        />
      </div>
    </div>
  );
}

export default function POSPage() {
  return (
    <ProtectedRoute>
      <ErrorBoundary>
        <POSContent />
      </ErrorBoundary>
    </ProtectedRoute>
  );
}
