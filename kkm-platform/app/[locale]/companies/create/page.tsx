'use client';

import React from 'react';
import CompanyForm from '@/features/companies/components/CompanyForm';

export default function CreateCompanyPage() {
  const handleSubmit = async (data: unknown) => {
    // API call will be implemented
    console.log('Create company:', data);
  };

  return <CompanyForm onSubmit={handleSubmit} />;
}
