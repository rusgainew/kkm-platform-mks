'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import BankAccountForm from '@/features/bank-accounts/components/BankAccountForm';
import { useApiToken } from '@/lib/hooks/useApiToken';

export default function CreateBankAccountPage() {
  const router = useRouter();
  const token = useApiToken();

  const handleSubmit = async (data: unknown) => {
    if (!token) {
      throw new Error('Требуется аутентификация');
    }

    const response = await fetch('/api/bank-accounts', {
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

    const account = await response.json();
    
    // Показываем сообщение об успехе
    setTimeout(() => {
      router.push(`/bank-accounts/${account.id}`);
    }, 2000);
  };

  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        <BankAccountForm onSubmit={handleSubmit} />
      </div>
    </div>
  );
}
