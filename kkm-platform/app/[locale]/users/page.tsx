'use client';

import UsersList from '@/components/users/UsersList';
import ProtectedRoute from '@/components/auth/ProtectedRoute';
import SPADashboardLayout from '@/components/layout/SPADashboardLayout';
import { ErrorBoundary, UsersErrorFallback } from '@/components/error';
import { UsersPageSkeleton } from '@/components/loading';
import { useState, useEffect } from 'react';

export default function UsersPage() {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Simulate loading
    const timer = setTimeout(() => setIsLoading(false), 500);
    return () => clearTimeout(timer);
  }, []);

  if (isLoading) {
    return (
      <SPADashboardLayout currentPage="users">
        <div className="p-8">
          <UsersPageSkeleton />
        </div>
      </SPADashboardLayout>
    );
  }

  return (
    <ProtectedRoute>
      <SPADashboardLayout currentPage="users">
        <ErrorBoundary>
          <UsersList />
        </ErrorBoundary>
      </SPADashboardLayout>
    </ProtectedRoute>
  );
}
