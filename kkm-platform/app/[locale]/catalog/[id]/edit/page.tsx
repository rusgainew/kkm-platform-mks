'use client';

import React, { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import CatalogForm from '@/features/catalog/components/CatalogForm';
import { getCatalogItem, updateCatalogItem } from '@/lib/api/catalog';
import { Loader2, AlertCircle } from 'lucide-react';
import type { CatalogItem, UpdateCatalogItemRequest } from '@/types/entities';

export default function EditProductPage() {
  const params = useParams();
  const router = useRouter();
  const productId = params.id as string;
  const [product, setProduct] = useState<CatalogItem | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadProduct = async () => {
      try {
        setIsLoading(true);
        const response = await getCatalogItem(productId);
        if (response.data) {
          setProduct(response.data);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Ошибка загрузки');
      } finally {
        setIsLoading(false);
      }
    };

    loadProduct();
  }, [productId]);

  const handleSubmit = async (data: UpdateCatalogItemRequest): Promise<CatalogItem> => {
    const response = await updateCatalogItem(productId, data);
    if (response.data) {
      router.push('/catalog');
      return response.data;
    }
    throw new Error('Не удалось обновить товар');
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Loader2 className="animate-spin text-purple-400" size={32} />
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="flex items-center gap-2 text-red-400">
          <AlertCircle size={24} />
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        <CatalogForm initialData={product || undefined} onSubmit={handleSubmit} mode="edit" />
      </div>
    </div>
  );
}
