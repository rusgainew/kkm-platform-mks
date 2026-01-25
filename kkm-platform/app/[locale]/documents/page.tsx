'use client';

import DocumentsList from '@/components/documents/DocumentsList';
import ProtectedRoute from '@/components/auth/ProtectedRoute';
import SPADashboardLayout from '@/components/layout/SPADashboardLayout';
import { ErrorBoundary } from '@/components/error';
import { useState, useEffect } from 'react';

export default function DocumentsPage() {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Simulate loading
    const timer = setTimeout(() => setIsLoading(false), 500);
    return () => clearTimeout(timer);
  }, []);

  if (isLoading) {
    return (
      <SPADashboardLayout currentPage="documents">
        <div className="p-8">
          <div className="animate-pulse space-y-4">
            <div className="h-8 bg-gray-700 rounded w-1/4"></div>
            <div className="h-12 bg-gray-700 rounded"></div>
            <div className="h-64 bg-gray-700 rounded"></div>
          </div>
        </div>
      </SPADashboardLayout>
    );
  }

  return (
    <ProtectedRoute>
      <SPADashboardLayout currentPage="documents">
        <ErrorBoundary>
          <DocumentsList />
        </ErrorBoundary>
      </SPADashboardLayout>
    </ProtectedRoute>
  );
}
