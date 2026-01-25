/**
 * Companies page
 */

'use client';

import React from 'react';
import ProtectedRoute from '@/components/auth/ProtectedRoute';
import SPADashboardLayout from '@/components/layout/SPADashboardLayout';
import { CompaniesDashboard } from '@/components/companies-dashboard';
import { ErrorBoundary, CompaniesErrorFallback } from '@/components/error';

export default function CompaniesPage() {
  return (
    <ProtectedRoute requiredPermission="view_companies">
      <SPADashboardLayout currentPage="companies">
        <ErrorBoundary>
          <CompaniesDashboard />
        </ErrorBoundary>
      </SPADashboardLayout>
    </ProtectedRoute>
  );
}
