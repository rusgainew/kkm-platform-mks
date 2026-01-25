'use client';

import { LineChart, Line, BarChart, Bar, PieChart, Pie, Cell, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface ChartData {
  name: string;
  value: number;
  [key: string]: any;
}

interface RevenueChartProps {
  data: Array<{ date: string; revenue: number; paid: number; pending: number }>;
}

export function RevenueChart({ data }: RevenueChartProps) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h3 className="text-lg font-semibold mb-4">Доход за период</h3>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="date" />
          <YAxis />
          <Tooltip 
            formatter={(value: any) => {
              if (typeof value === 'number') {
                return `$${value.toFixed(2)}`;
              }
              return value;
            }}
          />
          <Legend />
          <Line type="monotone" dataKey="revenue" stroke="#3b82f6" name="Всего" strokeWidth={2} />
          <Line type="monotone" dataKey="paid" stroke="#10b981" name="Оплачено" strokeWidth={2} />
          <Line type="monotone" dataKey="pending" stroke="#f59e0b" name="В ожидании" strokeWidth={2} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}

interface InvoiceStatusChartProps {
  data: Array<{ status: string; count: number }>;
}

export function InvoiceStatusChart({ data }: InvoiceStatusChartProps) {
  const COLORS = {
    paid: '#10b981',
    draft: '#6b7280',
    sent: '#3b82f6',
    overdue: '#ef4444',
    cancelled: '#9ca3af',
  };

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h3 className="text-lg font-semibold mb-4">Статусы счетов</h3>
      <ResponsiveContainer width="100%" height={300}>
        <PieChart>
          <Pie
            data={data}
            cx="50%"
            cy="50%"
            labelLine={false}
            label={({ name, value }) => `${name}: ${value}`}
            outerRadius={100}
            fill="#8884d8"
            dataKey="count"
          >
            {data.map((entry, index) => (
              <Cell key={`cell-${index}`} fill={COLORS[entry.status as keyof typeof COLORS] || '#8884d8'} />
            ))}
          </Pie>
          <Tooltip />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
}

interface InventoryChartProps {
  data: Array<{ name: string; stock: number; sold: number }>;
}

export function InventoryChart({ data }: InventoryChartProps) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h3 className="text-lg font-semibold mb-4">Товары (топ 10)</h3>
      <ResponsiveContainer width="100%" height={350}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="name" angle={-45} textAnchor="end" height={100} />
          <YAxis />
          <Tooltip />
          <Legend />
          <Bar dataKey="stock" fill="#3b82f6" name="На складе" />
          <Bar dataKey="sold" fill="#10b981" name="Продано" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
