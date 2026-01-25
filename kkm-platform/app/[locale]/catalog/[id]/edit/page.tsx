'use client';

import React, { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import CatalogForm from '@/features/catalog/components/CatalogForm';
import { useCatalog } from '@/hooks';
import { Loader2, AlertCircle } from 'lucide-react';

export default function EditProductPage() {
  const params = useParams();
  const productId = params.id as string;
  const { getItemById, fetchItemById, isLoading, error } = useCatalog();
  const [product, setProduct] = useState<unknown>(null);

  useEffect(() => {
    const loadProduct = async () => {
      try {
        await fetchItemById(productId);
        const item = getItemById(productId);
        if (item) {
          setProduct(item);
        }
      } catch (err) {
        console.error('Error loading product:', err);
      }
    };

    loadProduct();
  }, [productId]);

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
        <CatalogForm initialData={product} itemId={productId} />
      </div>
    </div>
  );
}
