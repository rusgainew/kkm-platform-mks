'use client';

import React from 'react';
import {
  BarChart3,
  Users,
  Package,
  TrendingUp,
  AlertCircle,
  Zap,
  Clock,
  Eye,
  EyeOff,
  LogOut,
  Building2,
} from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useLogoutMutation } from '@/lib/hooks/useAuthApi';

interface DashboardLeftPanelProps {
  activeTerminals?: number;
  totalTerminals?: number;
  criticalAlerts?: number;
}

export default function DashboardLeftPanel({
  activeTerminals = 0,
  totalTerminals = 0,
  criticalAlerts = 0,
}: DashboardLeftPanelProps) {
  const [isExpanded, setIsExpanded] = React.useState(true);
  const router = useRouter();
  const logoutMutation = useLogoutMutation();

  const handleLogout = async () => {
    try {
      await logoutMutation.mutateAsync();
      router.push('/auth');
    } catch (err) {
      console.error('Logout failed:', err);
    }
  };

  const menuItems = [
    {
      id: 'overview',
      label: 'Обзор',
      icon: BarChart3,
      href: '#overview',
      isExternal: false,
    },
    {
      id: 'terminals',
      label: 'Кассы',
      icon: Zap,
      href: '#terminals',
      isExternal: false,
    },
    {
      id: 'sales',
      label: 'Продажи',
      icon: TrendingUp,
      href: '#sales',
      isExternal: false,
    },
    {
      id: 'inventory',
      label: 'Инвентарь',
      icon: Package,
      href: '#inventory',
      isExternal: false,
    },
    {
      id: 'companies',
      label: 'Компании',
      icon: Building2,
      href: '/companies',
      isExternal: true,
    },
    {
      id: 'users',
      label: 'Пользователи',
      icon: Users,
      href: '/users',
      isExternal: true,
    },
  ];

  const scrollToSection = (id: string) => {
    const element = document.querySelector(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
  };

  return (
    <div
      className={`${
        isExpanded ? 'w-64' : 'w-20'
      } bg-gray-900 border-r border-gray-800 transition-all duration-300 h-screen sticky top-0 flex flex-col shadow-2xl`}
    >
      {/* Header */}
      <div className="p-6 border-b border-gray-800 flex items-center justify-between">
        {isExpanded && (
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-linear-to-br from-blue-500 to-emerald-500 rounded-lg flex items-center justify-center">
              <BarChart3 className="w-5 h-5 text-white" />
            </div>
            <h2 className="text-white font-bold text-lg">Dashbord</h2>
          </div>
        )}
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="p-2 hover:bg-gray-800 rounded-lg text-gray-400 hover:text-white transition-colors"
          title={isExpanded ? 'Свернуть' : 'Развернуть'}
        >
          {isExpanded ? (
            <EyeOff className="w-5 h-5" />
          ) : (
            <Eye className="w-5 h-5" />
          )}
        </button>
      </div>

      {/* Status Cards */}
      {isExpanded && (
        <div className="p-4 space-y-3 border-b border-gray-800">
          {/* Terminal Status */}
          <div className="bg-gray-800/50 rounded-lg p-3">
            <div className="flex items-center gap-2 mb-2">
              <Zap className="w-4 h-4 text-emerald-400" />
              <span className="text-sm text-gray-300">Статус касс</span>
            </div>
            <div className="text-2xl font-bold text-white">
              {activeTerminals}/{totalTerminals}
            </div>
            <p className="text-xs text-gray-400 mt-1">активных из всех</p>
          </div>

          {/* Alerts */}
          {criticalAlerts > 0 && (
            <div className="bg-red-900/20 border border-red-800 rounded-lg p-3">
              <div className="flex items-center gap-2 mb-2">
                <AlertCircle className="w-4 h-4 text-red-400" />
                <span className="text-sm text-red-300">Критичные алерты</span>
              </div>
              <div className="text-2xl font-bold text-red-400">{criticalAlerts}</div>
              <p className="text-xs text-red-300 mt-1">требуют внимания</p>
            </div>
          )}

          {/* Time */}
          <div className="bg-gray-800/50 rounded-lg p-3">
            <div className="flex items-center gap-2 mb-2">
              <Clock className="w-4 h-4 text-blue-400" />
              <span className="text-sm text-gray-300">Текущее время</span>
            </div>
            <div className="text-lg font-mono text-white">
              {new Date().toLocaleTimeString('ru-RU', {
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
              })}
            </div>
          </div>
        </div>
      )}

      {/* Navigation Menu */}
      <nav className="flex-1 p-4 space-y-2 overflow-y-auto">
        {menuItems.map((item) => {
          const Icon = item.icon;
          const content = (
            <>
              <Icon className="w-5 h-5 shrink-0" />
              {isExpanded && <span className="text-sm font-medium">{item.label}</span>}
            </>
          );

          if (item.isExternal) {
            return (
              <Link
                key={item.id}
                href={item.href}
                className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all duration-200 ${
                  isExpanded ? '' : 'justify-center'
                } hover:bg-gray-800 text-gray-300 hover:text-white`}
                title={item.label}
              >
                {content}
              </Link>
            );
          }

          return (
            <button
              key={item.id}
              onClick={() => scrollToSection(item.href)}
              className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all duration-200 ${
                isExpanded ? '' : 'justify-center'
              } text-gray-300 hover:text-white hover:bg-gray-800`}
              title={item.label}
            >
              {content}
            </button>
          );
        })}
      </nav>

      {/* Footer */}
      {isExpanded && (
        <div className="p-4 border-t border-gray-800 space-y-3">
          <div className="text-xs text-gray-500">
            <p>Последнее обновление:</p>
            <p className="font-mono text-gray-400">
              {new Date().toLocaleTimeString('ru-RU')}
            </p>
          </div>
          <button
            onClick={handleLogout}
            disabled={logoutMutation.isPending}
            className="w-full flex items-center gap-2 px-4 py-2 rounded-lg bg-red-900/20 hover:bg-red-900/40 text-red-400 hover:text-red-300 transition-colors disabled:opacity-50"
            title="Выход"
          >
            <LogOut className="w-4 h-4" />
            <span className="text-sm font-medium">
              {logoutMutation.isPending ? 'Выход...' : 'Выход'}
            </span>
          </button>
        </div>
      )}
      
      {!isExpanded && (
        <div className="p-4 border-t border-gray-800">
          <button
            onClick={handleLogout}
            disabled={logoutMutation.isPending}
            className="w-full flex items-center justify-center p-2 rounded-lg text-red-400 hover:bg-red-900/40 transition-colors disabled:opacity-50"
            title="Выход"
          >
            <LogOut className="w-5 h-5" />
          </button>
        </div>
      )}
    </div>
  );
}
