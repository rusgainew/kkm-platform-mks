'use client';

import { useAuthStore } from '@/store/authStore';
import { ShoppingCart, BarChart3, Users, FileText, Smartphone, Shield, LogOut } from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function HomePage() {
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);
  const router = useRouter();
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  const handleLogout = async () => {
    setIsLoggingOut(true);
    logout();
    router.push('/');
  };

  return (
    <div className="min-h-screen bg-linear-to-br from-gray-950 via-gray-900 to-gray-950">
      {/* Navigation */}
      <nav className="bg-gray-900/50 backdrop-blur-sm border-b border-gray-800 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-10 h-10 bg-linear-to-br from-blue-500 to-emerald-500 rounded-lg flex items-center justify-center shadow-lg shadow-blue-500/30">
              <span className="text-white font-bold text-lg">К</span>
            </div>
            <span className="text-xl font-bold text-white">KKM Platform</span>
          </div>
          <div className="flex items-center gap-4">
            {user ? (
              <>
                <span className="text-gray-300">Привет, {user.name}!</span>
                {user.role === 'store_manager' ? (
                  <Link
                    href="/manager-dashboard"
                    className="px-4 py-2 bg-linear-to-br from-blue-500 to-emerald-500 text-white rounded-lg font-semibold hover:shadow-lg hover:shadow-blue-500/30 transition-all"
                  >
                    Панель менеджера
                  </Link>
                ) : user.role === 'cashier' ? (
                  <Link
                    href="/pos"
                    className="px-4 py-2 bg-linear-to-br from-blue-500 to-emerald-500 text-white rounded-lg font-semibold hover:shadow-lg hover:shadow-blue-500/30 transition-all"
                  >
                    POS Касса
                  </Link>
                ) : (
                  <Link
                    href="/dashboard"
                    className="px-4 py-2 bg-linear-to-br from-blue-500 to-emerald-500 text-white rounded-lg font-semibold hover:shadow-lg hover:shadow-blue-500/30 transition-all"
                  >
                    Панель управления
                  </Link>
                )}
                <button
                  onClick={handleLogout}
                  disabled={isLoggingOut}
                  className="px-4 py-2 bg-red-600 hover:bg-red-700 disabled:bg-red-800 text-white rounded-lg font-semibold flex items-center gap-2 transition-all"
                >
                  <LogOut className="w-4 h-4" />
                  {isLoggingOut ? 'Выход...' : 'Выход'}
                </button>
              </>
            ) : (
              <Link
                href="/auth"
                className="px-4 py-2 bg-linear-to-br from-blue-500 to-emerald-500 text-white rounded-lg font-semibold hover:shadow-lg hover:shadow-blue-500/30 transition-all"
              >
                Войти
              </Link>
            )}
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative overflow-hidden py-20">
        <div className="max-w-7xl mx-auto px-6">
          <div className="text-center">
            <h1 className="text-5xl md:text-6xl font-bold text-white mb-6">
              KKM Platform
            </h1>
            <p className="text-xl text-gray-400 mb-8 max-w-2xl mx-auto">
              Централизованная система управления кассовыми терминалами с поддержкой POS, аналитики и управления пользователями
            </p>
            {!user && (
              <Link
                href="/auth"
                className="inline-block px-8 py-3 bg-linear-to-br from-blue-500 to-emerald-500 text-white rounded-lg font-semibold hover:shadow-lg hover:shadow-blue-500/30 transition-all"
              >
                Начать работу
              </Link>
            )}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-gray-900/30">
        <div className="max-w-7xl mx-auto px-6">
          <h2 className="text-3xl font-bold text-white mb-12 text-center">
            Основные возможности
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[
              { icon: ShoppingCart, title: 'POS Система', desc: 'Полнофункциональная касса' },
              { icon: BarChart3, title: 'Аналитика', desc: 'Детальные отчеты и статистика' },
              { icon: Users, title: 'Управление', desc: 'Администрирование пользователей' },
              { icon: FileText, title: 'Документы', desc: 'Работа с счетами и счфактурами' },
              { icon: Smartphone, title: 'Мобильная', desc: 'Доступ с любых устройств' },
              { icon: Shield, title: 'Безопасность', desc: 'Защита данных и аудит' },
            ].map((feature, i) => (
              <div
                key={i}
                className="p-6 bg-gray-800/50 border border-gray-700 rounded-xl hover:border-blue-500/50 transition-all"
              >
                <feature.icon className="w-8 h-8 text-blue-500 mb-4" />
                <h3 className="text-lg font-semibold text-white mb-2">
                  {feature.title}
                </h3>
                <p className="text-gray-400">{feature.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20">
        <div className="max-w-4xl mx-auto px-6 text-center">
          <h2 className="text-3xl font-bold text-white mb-6">
            Готовы начать?
          </h2>
          <p className="text-gray-400 mb-8">
            Присоединитесь к тысячам компаний, которые используют KKM Platform для управления своим бизнесом
          </p>
          <div className="flex gap-4 justify-center flex-wrap">
            {!user ? (
              <>
                <Link
                  href="/auth"
                  className="px-8 py-3 bg-linear-to-br from-blue-500 to-emerald-500 text-white rounded-lg font-semibold hover:shadow-lg hover:shadow-blue-500/30 transition-all"
                >
                  Зарегистрироваться
                </Link>
                <Link
                  href="/auth"
                  className="px-8 py-3 border border-gray-700 text-white rounded-lg font-semibold hover:bg-gray-800 transition-all"
                >
                  Уже есть аккаунт
                </Link>
              </>
            ) : (
              <Link
                href="/dashboard"
                className="px-8 py-3 bg-linear-to-br from-blue-500 to-emerald-500 text-white rounded-lg font-semibold hover:shadow-lg hover:shadow-blue-500/30 transition-all"
              >
                Перейти в панель управления
              </Link>
            )}
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-950 border-t border-gray-800 py-8">
        <div className="max-w-7xl mx-auto px-6 text-center text-gray-400">
          <p>&copy; 2024 KKM Platform. Все права защищены.</p>
        </div>
      </footer>
    </div>
  );
}
