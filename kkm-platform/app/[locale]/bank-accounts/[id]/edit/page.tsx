'use client';

import React, { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import BankAccountForm from '@/features/bank-accounts/components/BankAccountForm';
import { useApiToken } from '@/lib/hooks/useApiToken';
import { getBankAccountById } from '@/lib/api/bank-accounts';
import { Loader2 } from 'lucide-react';

export default function EditBankAccountPage() {
  const params = useParams();
  const accountId = params.id as string;
  const [account, setAccount] = useState<unknown>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const token = useApiToken();

  useEffect(() => {
    const loadAccount = async () => {
      if (!token) {
        setError('Требуется аутентификация');
        setIsLoading(false);
        return;
      }

      try {
        setIsLoading(true);
        const data = await getBankAccountById(accountId, token);
        setAccount(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Ошибка загрузки');
      } finally {
        setIsLoading(false);
      }
    };

    loadAccount();
  }, [accountId, token]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Loader2 className="animate-spin text-green-400" size={32} />
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-red-400">{error}</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        <BankAccountForm initialData={account as any} />
      </div>
    </div>
  );
}
