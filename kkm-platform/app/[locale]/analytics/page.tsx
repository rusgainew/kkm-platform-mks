'use client';

import { useInvoicesQuery } from '@/lib/hooks/useInvoicesQuery';
import { useCatalogQuery } from '@/lib/hooks/useCatalogQuery';
import { useUsersQuery } from '@/lib/hooks/useUsersQuery';
import { RevenueChart, InvoiceStatusChart, InventoryChart } from '@/components/dashboard/Charts';
import { KPIDashboard, createFinanceKPIs } from '@/components/dashboard/KPIMetrics';
import { ExportButtons } from '@/components/export/ExportButtons';

export default function AdvancedDashboard() {
  const { data: invoices = [], isLoading: invoicesLoading } = useInvoicesQuery();
  const { data: catalogItems = [], isLoading: catalogLoading } = useCatalogQuery();
  const { data: users = [], isLoading: usersLoading } = useUsersQuery();

  // Подготовить данные для графиков
  const invoicesByStatus = Object.entries(
    invoices.reduce((acc, inv) => {
      acc[inv.status] = (acc[inv.status] || 0) + 1;
      return acc;
    }, {} as Record<string, number>)
  ).map(([status, count]) => ({
    status,
    count,
  }));

  // Данные по доходам за последние 7 дней (mock)
  const revenueData = Array.from({ length: 7 }, (_, i) => {
    const date = new Date();
    date.setDate(date.getDate() - (6 - i));
    const dayInvoices = invoices.filter(
      (inv) =>
        new Date(inv.issueDate).toDateString() === date.toDateString()
    );

    return {
      date: date.toLocaleDateString('ru-RU', { weekday: 'short', day: '2-digit', month: '2-digit' }),
      revenue: dayInvoices.reduce((sum, inv) => sum + inv.totalAmount, 0),
      paid: dayInvoices
        .filter((inv) => inv.status === 'paid')
        .reduce((sum, inv) => sum + inv.totalAmount, 0),
      pending: dayInvoices
        .filter((inv) => inv.status !== 'paid' && inv.status !== 'cancelled')
        .reduce((sum, inv) => sum + inv.totalAmount, 0),
    };
  });

  // Топ 10 товаров по остаткам
  const topInventory = catalogItems
    .sort((a, b) => b.stock - a.stock)
    .slice(0, 10)
    .map((item) => ({
      name: item.name,
      stock: item.stock,
      sold: 0, // Это можно вычислить из реальных данных
    }));

  // Вычислить KPI
  const totalRevenue = invoices.reduce((sum, inv) => sum + inv.totalAmount, 0);
  const paidAmount = invoices
    .filter((inv) => inv.status === 'paid')
    .reduce((sum, inv) => sum + inv.totalAmount, 0);
  const pendingAmount = invoices
    .filter((inv) => inv.status === 'sent' || inv.status === 'draft')
    .reduce((sum, inv) => sum + inv.totalAmount, 0);
  const lowStockItems = catalogItems.filter((item) => item.stock < 10).length;
  const activeUsers = users.filter((u) => u.status === 'active').length;

  const kpis = createFinanceKPIs({
    totalRevenue,
    paidAmount,
    pendingAmount,
    overallUsers: users.length,
    activeUsers,
    totalOrders: invoices.length,
    lowStockItems,
  });

  const isLoading = invoicesLoading || catalogLoading || usersLoading;

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Заголовок */}
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Аналитика</h1>
        <ExportButtons
          data={invoices}
          filename="invoices_report"
          title="Отчет по счетам"
        />
      </div>

      {/* KPI Метрики */}
      <KPIDashboard metrics={kpis} />

      {/* Графики */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <RevenueChart data={revenueData} />
        <InvoiceStatusChart data={invoicesByStatus} />
      </div>

      {/* Товары */}
      {topInventory.length > 0 && (
        <InventoryChart data={topInventory} />
      )}

      {/* Информационные карточки */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-sm font-medium text-gray-600">Товаров на низком уровне</h3>
          <p className="text-2xl font-bold text-yellow-600 mt-2">{lowStockItems}</p>
          <p className="text-xs text-gray-500 mt-1">менее 10 единиц</p>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-sm font-medium text-gray-600">Просроченные счета</h3>
          <p className="text-2xl font-bold text-red-600 mt-2">
            {invoices.filter((inv) => inv.status === 'overdue').length}
          </p>
          <p className="text-xs text-gray-500 mt-1">
            $
            {invoices
              .filter((inv) => inv.status === 'overdue')
              .reduce((sum, inv) => sum + inv.totalAmount, 0)
              .toFixed(2)}
          </p>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-sm font-medium text-gray-600">Средний чек</h3>
          <p className="text-2xl font-bold text-blue-600 mt-2">
            ${invoices.length > 0 ? (totalRevenue / invoices.length).toFixed(2) : '0.00'}
          </p>
          <p className="text-xs text-gray-500 mt-1">{invoices.length} счетов</p>
        </div>
      </div>
    </div>
  );
}
