'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import CatalogForm from '@/features/catalog/components/CatalogForm';
import { createCatalogItem } from '@/lib/api/catalog';
import type { CatalogItem, CreateCatalogItemRequest, UpdateCatalogItemRequest } from '@/types/entities';

export default function CreateProductPage() {
  const router = useRouter();
  
  const handleSubmit = async (data: CreateCatalogItemRequest | UpdateCatalogItemRequest): Promise<CatalogItem> => {
    const response = await createCatalogItem(data as CreateCatalogItemRequest);
    if (response.data) {
      router.push('/catalog');
      return response.data;
    }
    throw new Error('Не удалось создать товар');
  };

  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        <CatalogForm onSubmit={handleSubmit} mode="create" />
      </div>
    </div>
  );
}
