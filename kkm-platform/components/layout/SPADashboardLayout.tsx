/**
 * SPA Dashboard Layout - Main container with sidebar navigation
 */

'use client';

import React, { useState } from 'react';
import {
  LayoutDashboard,
  Building2,
  Users,
  FileText,
  Settings,
  ChevronLeft,
  ChevronRight,
  Key,
  User,
  Package,
} from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/authStore';
import ChangePasswordModal from '@/components/auth/ChangePasswordModal';
import RoleSelector from '@/components/auth/RoleSelector';

interface SPADashboardLayoutProps {
  children: React.ReactNode;
  currentPage: 'dashboard' | 'manager-dashboard' | 'companies' | 'users' | 'documents' | 'catalog' | 'settings';
}

export default function SPADashboardLayout({
  children,
  currentPage,
}: SPADashboardLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [showChangePassword, setShowChangePassword] = useState(false);
  const router = useRouter();
  const user = useAuthStore((state) => state.user);

  const navigationItems = [
    {
      id: 'dashboard',
      label: 'Дашборд',
      icon: <LayoutDashboard size={20} />,
      href: '/dashboard',
    },
    {
      id: 'catalog',
      label: 'Каталог',
      icon: <Package size={20} />,
      href: '/catalog',
    },
    {
      id: 'companies',
      label: 'Компании',
      icon: <Building2 size={20} />,
      href: '/companies',
    },
    {
      id: 'users',
      label: 'Пользователи',
      icon: <Users size={20} />,
      href: '/users',
    },
    {
      id: 'documents',
      label: 'Документы',
      icon: <FileText size={20} />,
      href: '/documents',
    },
    {
      id: 'settings',
      label: 'Настройки',
      icon: <Settings size={20} />,
      href: '/settings',
    },
  ];

  return (
    <div className="min-h-screen bg-linear-to-br from-gray-950 via-gray-900 to-gray-950 flex">
      {/* Sidebar */}
      <aside
        className={`${
          sidebarOpen ? 'w-64' : 'w-20'
        } bg-gray-900/50 border-r border-gray-800 backdrop-blur-sm transition-all duration-300 flex flex-col`}
      >
        {/* Logo */}
        <div className="p-4 border-b border-gray-800 flex items-center justify-between">
          {sidebarOpen && (
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-linear-to-br from-blue-500 to-emerald-500 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold">К</span>
              </div>
              <div>
                <p className="text-white font-bold text-sm">KKM</p>
                <p className="text-gray-400 text-xs">Platform</p>
              </div>
            </div>
          )}
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="p-1 hover:bg-gray-800 rounded-lg transition"
          >
            {sidebarOpen ? (
              <ChevronLeft size={20} className="text-gray-400" />
            ) : (
              <ChevronRight size={20} className="text-gray-400" />
            )}
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 p-4 space-y-2">
          {navigationItems.map((item) => (
            <button
              key={item.id}
              onClick={() => router.push(item.href)}
              className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition ${
                currentPage === item.id
                  ? 'bg-blue-600/20 text-blue-400 border border-blue-700'
                  : 'text-gray-400 hover:bg-gray-800/50 hover:text-gray-300'
              }`}
            >
              {item.icon}
              {sidebarOpen && <span className="text-sm font-medium">{item.label}</span>}
            </button>
          ))}
        </nav>

        {/* Bottom Actions */}
        <div className="p-4 border-t border-gray-800 space-y-2">
          {/* User Info */}
          {user && sidebarOpen && (
            <div className="px-3 py-3 bg-gray-800/50 rounded-lg border border-gray-700 mb-3">
              <div className="flex items-center gap-2 mb-2">
                <div className="w-8 h-8 bg-linear-to-br from-blue-500 to-emerald-500 rounded-lg flex items-center justify-center">
                  <User size={16} className="text-white" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-semibold text-white truncate">{user.firstName} {user.lastName}</p>
                  <p className="text-xs text-gray-400 truncate">{user.email}</p>
                </div>
              </div>
              <div className="px-2 py-1 bg-gray-900 rounded text-center">
                <p className="text-xs font-medium text-blue-300">Роль: {user.role}</p>
              </div>
            </div>
          )}

          {/* Role Selector - Always visible */}
          <div className={sidebarOpen ? '' : 'flex justify-center'}>
            <RoleSelector />
          </div>

          {/* Password Button - Always visible with icon indicator */}
          <button
            onClick={() => setShowChangePassword(true)}
            className={`w-full flex items-center gap-3 px-4 py-3 text-gray-400 hover:bg-gray-800/50 hover:text-gray-300 rounded-lg transition text-sm ${
              !sidebarOpen ? 'justify-center' : ''
            }`}
            title={!sidebarOpen ? 'Изменить пароль' : ''}
          >
            <Key size={20} />
            {sidebarOpen && <span>Пароль</span>}
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 overflow-auto">
        {children}
      </main>

      {/* Modals */}
      <ChangePasswordModal
        isOpen={showChangePassword}
        onClose={() => setShowChangePassword(false)}
      />
    </div>
  );
}
