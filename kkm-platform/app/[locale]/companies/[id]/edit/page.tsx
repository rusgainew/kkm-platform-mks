'use client';

import React, { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import CompanyForm from '@/features/companies/components/CompanyForm';
import { getCompany } from '@/lib/api/companies';
import { Loader2 } from 'lucide-react';

export default function EditCompanyPage() {
  const params = useParams();
  const companyId = params.id as string;
  const [company, setCompany] = useState<unknown>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadCompany = async () => {
      if (!companyId) {
        setError('ID компании не найден');
        setIsLoading(false);
        return;
      }

      try {
        setIsLoading(true);
        const data = await getCompany(companyId);
        setCompany(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Ошибка загрузки');
      } finally {
        setIsLoading(false);
      }
    };

    loadCompany();
  }, [companyId]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Loader2 className="animate-spin text-blue-400" size={32} />
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
        <CompanyForm initialData={company as any} onSubmit={async () => {}} />
      </div>
    </div>
  );
}
