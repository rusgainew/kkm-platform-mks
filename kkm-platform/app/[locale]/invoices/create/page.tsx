'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import InvoiceForm from '@/features/invoices/components/InvoiceForm';
import { useApiToken } from '@/lib/hooks/useApiToken';

export default function CreateInvoicePage() {
  const router = useRouter();
  const token = useApiToken();

  const handleSubmit = async (data: unknown) => {
    if (!token) {
      throw new Error('Требуется аутентификация');
    }

    const response = await fetch('/api/invoices', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Ошибка создания счета');
    }

    const invoice = await response.json();
    
    // Показываем сообщение об успехе
    setTimeout(() => {
      router.push(`/invoices/${invoice.id}`);
    }, 2000);
  };

  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        <InvoiceForm onSubmit={handleSubmit} />
      </div>
    </div>
  );
}
