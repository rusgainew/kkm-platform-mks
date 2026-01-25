'use client';

import CatalogList from '@/components/catalog/CatalogList';
import ProtectedRoute from '@/components/auth/ProtectedRoute';
import SPADashboardLayout from '@/components/layout/SPADashboardLayout';
import { ErrorBoundary } from '@/components/error';
import { CatalogPageSkeleton } from '@/components/loading';
import { useState, useEffect } from 'react';

export default function CatalogPage() {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Simulate loading
    const timer = setTimeout(() => setIsLoading(false), 500);
    return () => clearTimeout(timer);
  }, []);

  if (isLoading) {
    return (
      <SPADashboardLayout currentPage="catalog">
        <div className="p-8">
          <CatalogPageSkeleton />
        </div>
      </SPADashboardLayout>
    );
  }

  return (
    <ProtectedRoute>
      <SPADashboardLayout currentPage="catalog">
        <ErrorBoundary>
          <CatalogList />
        </ErrorBoundary>
      </SPADashboardLayout>
    </ProtectedRoute>
  );
}
