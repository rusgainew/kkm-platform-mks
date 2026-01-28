import { DashboardStats, Terminal, SalesData, InventoryItem } from '@/types';
import {
  MOCK_DASHBOARD_STATS,
  MOCK_TERMINALS,
  MOCK_SALES_DATA,
  MOCK_INVENTORY,
} from '@/constants/dashboard';

// Симуляция API endpoints для real-time обновлений

export async function fetchDashboardStats(_storeId?: string): Promise<DashboardStats> {
  // TODO: Использовать storeId для фильтрации данных конкретного магазина
  // Симуляция задержки сети
  await new Promise((resolve) => setTimeout(resolve, 300));

  // Симуляция изменений в реальном времени
  const randomVariation = Math.floor(Math.random() * 5000) - 2500;
  const randomTransactions = Math.floor(Math.random() * 10) - 5;

  return {
    ...MOCK_DASHBOARD_STATS,
    totalRevenue: Math.max(0, MOCK_DASHBOARD_STATS.totalRevenue + randomVariation),
    totalTransactions: Math.max(0, MOCK_DASHBOARD_STATS.totalTransactions + randomTransactions),
  };
}

export async function fetchTerminals(_storeId?: string): Promise<Terminal[]> {
  // TODO: Использовать storeId для фильтрации терминалов
  await new Promise((resolve) => setTimeout(resolve, 200));

  // Симуляция случайных изменений статусов
  return MOCK_TERMINALS.map((terminal) => {
    const shouldChange = Math.random() > 0.95; // 5% шанс изменения
    if (shouldChange) {
      const statuses: Terminal['status'][] = ['online', 'offline', 'paper_out'];
      const randomStatus = statuses[Math.floor(Math.random() * statuses.length)];
      return { ...terminal, status: randomStatus };
    }
    return terminal;
  });
}

export async function fetchSalesData(_storeId?: string): Promise<SalesData[]> {
  // TODO: Использовать storeId для фильтрации данных продаж
  await new Promise((resolve) => setTimeout(resolve, 250));

  // Симуляция обновления последнего часа
  return MOCK_SALES_DATA.map((item, index) => {
    if (index === MOCK_SALES_DATA.length - 1) {
      return {
        ...item,
        sales: item.sales + Math.floor(Math.random() * 1000),
        transactions: item.transactions + Math.floor(Math.random() * 5),
      };
    }
    return item;
  });
}

export async function fetchInventory(_storeId?: string): Promise<InventoryItem[]> {
  // TODO: Использовать storeId для фильтрации инвентаря
  await new Promise((resolve) => setTimeout(resolve, 180));
  return MOCK_INVENTORY;
}

// Симуляция нового заказа в real-time
export interface RealtimeTransaction {
  id: string;
  terminalId: string;
  amount: number;
  timestamp: Date;
  items: number;
}

export async function subscribeToTransactions(
  callback: (transaction: RealtimeTransaction) => void
): Promise<() => void> {
  // Симуляция WebSocket-подобного поведения
  const interval = setInterval(() => {
    const mockTransaction: RealtimeTransaction = {
      id: `TXN-${Date.now()}`,
      terminalId: `TRM-00${Math.floor(Math.random() * 8) + 1}`,
      amount: Math.floor(Math.random() * 3000) + 500,
      timestamp: new Date(),
      items: Math.floor(Math.random() * 10) + 1,
    };
    callback(mockTransaction);
  }, 5000); // Новая транзакция каждые 5 секунд

  return () => clearInterval(interval);
}
