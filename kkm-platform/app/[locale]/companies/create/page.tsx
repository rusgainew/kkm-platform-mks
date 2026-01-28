'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import CompanyForm from '@/features/companies/components/CompanyForm';
import { createCompany } from '@/lib/api/companies';
import type { Company, CreateCompanyRequest, UpdateCompanyRequest } from '@/types/entities';

export default function CreateCompanyPage() {
  const router = useRouter();
  
  const handleSubmit = async (data: CreateCompanyRequest | UpdateCompanyRequest): Promise<Company> => {
    const response = await createCompany(data as CreateCompanyRequest);
    if (response.data) {
      router.push('/companies');
      return response.data;
    }
    throw new Error('Не удалось создать компанию');
  };

  return <CompanyForm onSubmit={handleSubmit} mode="create" />;
}
