'use client';

import { Dashboard } from '@/components/dashboard';
import { MOCK_STORE_LOCATIONS, MOCK_REVENUE_TREND } from '@/constants/dashboard';
import { useDashboardData } from '@/lib/hooks/useDashboardData';
import { useAuthStore } from '@/store/authStore';
import ProtectedRoute from '@/components/auth/ProtectedRoute';
import SPADashboardLayout from '@/components/layout/SPADashboardLayout';
import { ErrorBoundary } from '@/components/error';
import { DashboardSkeleton } from '@/components/loading';

export default function DashboardPage() {
  const { stats, terminals, salesData, inventoryItems, isLoading } = useDashboardData();
  const user = useAuthStore((state) => state.user);

  if (isLoading || !stats || !terminals || !salesData) {
    return (
      <SPADashboardLayout currentPage="dashboard">
        <div className="p-8">
          <DashboardSkeleton />
        </div>
      </SPADashboardLayout>
    );
  }

  return (
    <ProtectedRoute requiredPermission="view_reports">
      <SPADashboardLayout currentPage="dashboard">
        <ErrorBoundary>
          <Dashboard
            user={user}
            stats={stats}
            salesData={salesData}
            terminals={terminals}
            inventoryItems={inventoryItems || []}
            storeLocations={MOCK_STORE_LOCATIONS}
            revenueTrend={MOCK_REVENUE_TREND}
          />
        </ErrorBoundary>
      </SPADashboardLayout>
    </ProtectedRoute>
  );
}
