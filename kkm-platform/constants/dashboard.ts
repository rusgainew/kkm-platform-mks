import { DashboardStats, SalesData, Terminal, InventoryItem } from '@/types';

export const MOCK_DASHBOARD_STATS: DashboardStats = {
  totalRevenue: 125450,
  activeTerminals: 18,
  totalTerminals: 20,
  avgCheck: 856,
  totalTransactions: 312,
};

export const MOCK_SALES_DATA: SalesData[] = [
  { hour: '09:00', sales: 8500, transactions: 12 },
  { hour: '10:00', sales: 12300, transactions: 18 },
  { hour: '11:00', sales: 15600, transactions: 22 },
  { hour: '12:00', sales: 18900, transactions: 28 },
  { hour: '13:00', sales: 21200, transactions: 31 },
  { hour: '14:00', sales: 16800, transactions: 24 },
  { hour: '15:00', sales: 19400, transactions: 27 },
  { hour: '16:00', sales: 22100, transactions: 32 },
  { hour: '17:00', sales: 25300, transactions: 36 },
  { hour: '18:00', sales: 28600, transactions: 41 },
  { hour: '19:00', sales: 26700, transactions: 38 },
  { hour: '20:00', sales: 24200, transactions: 34 },
];

export const MOCK_TERMINALS: Terminal[] = [
  {
    id: 'TRM-001',
    address: 'ул. Ленина, 10, филиал №1',
    status: 'online',
    lastSync: '2 мин. назад',
    revenue: 8450,
  },
  {
    id: 'TRM-002',
    address: 'ул. Красная, 42, филиал №1',
    status: 'online',
    lastSync: '1 мин. назад',
    revenue: 9320,
  },
  {
    id: 'TRM-003',
    address: 'ул. Октябрьская, 15, филиал №2',
    status: 'online',
    lastSync: '3 мин. назад',
    revenue: 7890,
  },
  {
    id: 'TRM-004',
    address: 'пр. Мира, 88, филиал №2',
    status: 'paper_out',
    lastSync: '45 мин. назад',
    revenue: 6250,
  },
  {
    id: 'TRM-005',
    address: 'ул. Пушкина, 5, филиал №3',
    status: 'online',
    lastSync: '2 мин. назад',
    revenue: 8900,
  },
  {
    id: 'TRM-006',
    address: 'ул. Минская, 20, филиал №3',
    status: 'offline',
    lastSync: '2 часа назад',
    revenue: 0,
  },
  {
    id: 'TRM-007',
    address: 'ул. Советская, 33, филиал №1',
    status: 'online',
    lastSync: '1 мин. назад',
    revenue: 9100,
  },
  {
    id: 'TRM-008',
    address: 'ул. Комсомольская, 7, филиал №2',
    status: 'online',
    lastSync: '2 мин. назад',
    revenue: 8650,
  },
];

export const MOCK_INVENTORY: InventoryItem[] = [
  {
    id: 1,
    name: 'Переносная POS касса (SUNMI)',
    currentStock: 2,
    minThreshold: 5,
    status: 'critical',
  },
  {
    id: 2,
    name: 'Рулон чека (80x80мм)',
    currentStock: 8,
    minThreshold: 15,
    status: 'low',
  },
  {
    id: 3,
    name: 'Чернила для принтера',
    currentStock: 3,
    minThreshold: 10,
    status: 'critical',
  },
  {
    id: 4,
    name: 'USB кабель (Type-C)',
    currentStock: 12,
    minThreshold: 8,
    status: 'ok',
  },
  {
    id: 5,
    name: 'Батарейки АА (40 шт)',
    currentStock: 6,
    minThreshold: 12,
    status: 'low',
  },
  {
    id: 6,
    name: 'Сумка для транспортировки',
    currentStock: 4,
    minThreshold: 3,
    status: 'ok',
  },
];

export const MOCK_STORE_LOCATIONS = [
  { id: 'store-1', name: 'Филиал №1 (Центр)', revenue: 45000, coordinates: { x: 25, y: 30 } },
  { id: 'store-2', name: 'Филиал №2 (Северный)', revenue: 38000, coordinates: { x: 60, y: 20 } },
  { id: 'store-3', name: 'Филиал №3 (Южный)', revenue: 52000, coordinates: { x: 40, y: 70 } },
  { id: 'store-4', name: 'Филиал №4 (Восточный)', revenue: 31000, coordinates: { x: 75, y: 45 } },
  { id: 'store-5', name: 'Филиал №5 (Западный)', revenue: 42000, coordinates: { x: 15, y: 55 } },
  { id: 'store-6', name: 'Филиал №6 (Центральный рынок)', revenue: 48000, coordinates: { x: 50, y: 50 } },
];

export const MOCK_REVENUE_TREND = [
  { date: 'Пн', revenue: 98500, transactions: 245 },
  { date: 'Вт', revenue: 105300, transactions: 268 },
  { date: 'Ср', revenue: 112800, transactions: 289 },
  { date: 'Чт', revenue: 118200, transactions: 301 },
  { date: 'Пт', revenue: 125450, transactions: 312 },
  { date: 'Сб', revenue: 142300, transactions: 356 },
  { date: 'Вс', revenue: 138900, transactions: 347 },
];
