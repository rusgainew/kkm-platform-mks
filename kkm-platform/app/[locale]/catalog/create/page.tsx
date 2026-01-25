'use client';

import React from 'react';
import CatalogForm from '@/features/catalog/components/CatalogForm';

export default function CreateProductPage() {
  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        <CatalogForm />
      </div>
    </div>
  );
}
