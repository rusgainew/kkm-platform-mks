'use client';

import { ShoppingCart, History, Package, Settings, Wifi, WifiOff, BarChart3, Users, FileText, Home, LogOut, Building2, Loader2 } from 'lucide-react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useLogoutMutation } from '@/lib/hooks/useAuthApi';

interface SidebarProps {
  isOnline?: boolean;
}

export default function Sidebar({ isOnline = true }: SidebarProps) {
  const pathname = usePathname();
  const router = useRouter();
  const logoutMutation = useLogoutMutation();

  const handleLogout = async () => {
    try {
      // Отправляем logout запрос на сервер
      await logoutMutation.mutateAsync();
      
      // Дополнительно очищаем localStorage
      if (typeof window !== 'undefined') {
        localStorage.removeItem('auth-storage');
        localStorage.removeItem('token');
        localStorage.removeItem('refreshToken');
      }
      
      // Переходим на страницу логина
      router.push('/auth');
    } catch (err) {
      console.error('Logout failed:', err);
      
      // Даже если ошибка, очищаем данные и редиректим
      if (typeof window !== 'undefined') {
        localStorage.removeItem('auth-storage');
        localStorage.removeItem('token');
        localStorage.removeItem('refreshToken');
      }
      router.push('/auth');
    }
  };

  const menuItems = [
    { id: 'home', path: '/', icon: Home, label: 'Главная' },
    { id: 'sales', path: '/pos', icon: ShoppingCart, label: 'Касса' },
    { id: 'dashboard', path: '/dashboard', icon: BarChart3, label: 'Аналитика' },
    { id: 'catalog', path: '/catalog', icon: Package, label: 'Каталог' },
    { id: 'companies', path: '/companies', icon: Building2, label: 'Компании' },
    { id: 'employees', path: '/employees', icon: Users, label: 'Сотрудники' },
    { id: 'fiscal', path: '/fiscal', icon: FileText, label: 'Фискальные отчеты' },
    { id: 'history', path: '/history', icon: History, label: 'История' },
    { id: 'inventory', path: '/inventory', icon: Package, label: 'Склад' },
    { id: 'settings', path: '/settings', icon: Settings, label: 'Настройки' },
  ];

  const getActiveTab = () => {
    const item = menuItems.find((m) => m.path === pathname);
    return item?.id || 'sales';
  };

  return (
    <div className="w-20 bg-gray-950 border-r border-gray-800 text-white flex flex-col items-center py-6 space-y-8 shadow-2xl">
      {/* Logo */}
      <div className="mb-4">
        <div className="w-12 h-12 bg-linear-to-br from-blue-500 to-emerald-500 rounded-lg flex items-center justify-center shadow-lg shadow-blue-500/30">
          <span className="text-white font-bold text-xl">К</span>
        </div>
      </div>

      {/* Status Indicator */}
      <div className="mb-2">
        {isOnline ? (
          <Wifi className="w-6 h-6 text-emerald-400 drop-shadow-[0_0_6px_rgba(16,185,129,0.5)]" />
        ) : (
          <WifiOff className="w-6 h-6 text-red-400 drop-shadow-[0_0_6px_rgba(239,68,68,0.5)]" />
        )}
      </div>

      {/* Menu Items */}
      {menuItems.map((item) => {
        const Icon = item.icon;
        const isActive = getActiveTab() === item.id;
        return (
          <Link key={item.id} href={item.path}>
            <button
              className={`flex flex-col items-center justify-center w-14 h-14 rounded-lg transition-all duration-200 ${
                isActive
                  ? 'bg-linear-to-br from-blue-500 to-emerald-500 text-white shadow-lg shadow-blue-500/30'
                  : 'text-gray-400 hover:bg-gray-800 hover:text-white hover:shadow-lg hover:scale-105'
              }`}
              title={item.label}
            >
              <Icon className="w-6 h-6" />
            </button>
          </Link>
        );
      })}

      {/* Logout Button */}
      <div className="mt-auto pt-8 border-t border-gray-800" suppressHydrationWarning>
        <button
          onClick={handleLogout}
          disabled={logoutMutation.isPending}
          className="flex flex-col items-center justify-center w-14 h-14 rounded-lg text-gray-400 hover:bg-red-900/50 hover:text-red-400 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
          title="Выход"
        >
          {logoutMutation.isPending ? (
            <Loader2 className="w-6 h-6 animate-spin" />
          ) : (
            <LogOut className="w-6 h-6" />
          )}
        </button>
      </div>
    </div>
  );
}
