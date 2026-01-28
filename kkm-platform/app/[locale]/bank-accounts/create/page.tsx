'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import BankAccountForm from '@/features/bank-accounts/components/BankAccountForm';

export default function CreateBankAccountPage() {
  const router = useRouter();

  const handleSuccess = () => {
    // Перенаправляем на список счетов после успешного создания
    router.push('/bank-accounts');
  };

  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        <BankAccountForm onSuccess={handleSuccess} />
      </div>
    </div>
  );
}
