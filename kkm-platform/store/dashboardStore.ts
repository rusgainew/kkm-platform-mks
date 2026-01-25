import { create } from 'zustand';
import { DashboardStats, Terminal, SalesData, InventoryItem } from '@/types';

interface DashboardState {
  stats: DashboardStats | null;
  terminals: Terminal[];
  salesData: SalesData[];
  inventoryItems: InventoryItem[];
  selectedStoreId: string | null;
  lastUpdate: Date | null;

  // Actions
  setStats: (stats: DashboardStats) => void;
  setTerminals: (terminals: Terminal[]) => void;
  setSalesData: (data: SalesData[]) => void;
  setInventoryItems: (items: InventoryItem[]) => void;
  setSelectedStore: (storeId: string | null) => void;
  updateTerminalStatus: (terminalId: string, status: Terminal['status']) => void;
  incrementRevenue: (amount: number) => void;
  markLastUpdate: () => void;
}

export const useDashboardStore = create<DashboardState>((set) => ({
  stats: null,
  terminals: [],
  salesData: [],
  inventoryItems: [],
  selectedStoreId: null,
  lastUpdate: null,

  setStats: (stats) => set({ stats }),
  
  setTerminals: (terminals) => set({ terminals }),
  
  setSalesData: (salesData) => set({ salesData }),
  
  setInventoryItems: (inventoryItems) => set({ inventoryItems }),
  
  setSelectedStore: (selectedStoreId) => set({ selectedStoreId }),
  
  updateTerminalStatus: (terminalId, status) =>
    set((state) => ({
      terminals: state.terminals.map((t) =>
        t.id === terminalId ? { ...t, status } : t
      ),
    })),
  
  incrementRevenue: (amount) =>
    set((state) => ({
      stats: state.stats
        ? {
            ...state.stats,
            totalRevenue: state.stats.totalRevenue + amount,
            totalTransactions: state.stats.totalTransactions + 1,
          }
        : null,
    })),
  
  markLastUpdate: () => set({ lastUpdate: new Date() }),
}));
