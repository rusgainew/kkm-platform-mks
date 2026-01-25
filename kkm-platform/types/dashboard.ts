export interface Terminal {
  id: string;
  address: string;
  status: 'online' | 'offline' | 'paper_out';
  lastSync: string;
  revenue: number;
}

export interface SalesData {
  hour: string;
  sales: number;
  transactions: number;
}

export interface InventoryItem {
  id: number;
  name: string;
  currentStock: number;
  minThreshold: number;
  status: 'critical' | 'low' | 'ok';
}

export interface DashboardStats {
  totalRevenue: number;
  activeTerminals: number;
  totalTerminals: number;
  avgCheck: number;
  totalTransactions: number;
}
