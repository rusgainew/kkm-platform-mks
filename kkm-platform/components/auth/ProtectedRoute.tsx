'use client';

import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/authStore';
import { Permission } from '@/types/auth';
import { ShieldAlert } from 'lucide-react';
import { ReactNode, useEffect, useState } from 'react';

interface ProtectedRouteProps {
  children: ReactNode;
  requiredPermission?: Permission;
  fallback?: ReactNode;
}

export default function ProtectedRoute({
  children,
  requiredPermission,
  fallback,
}: ProtectedRouteProps) {
  const router = useRouter();
  const hasPermission = useAuthStore((state) => state.hasPermission);
  const user = useAuthStore((state) => state.user);
  const tokens = useAuthStore((state) => state.tokens);
  const [isInitialized, setIsInitialized] = useState(false);
  const [hasRedirected, setHasRedirected] = useState(false);

  // Проверяем авторизацию после гидрации
  useEffect(() => {
    // Даём время на hydration Zustand из localStorage
    const timer = setTimeout(() => {
      if (!hasRedirected && (!user || !tokens?.accessToken)) {
        console.log('[ProtectedRoute] Not authenticated, redirecting to /auth');
        setHasRedirected(true);
        router.replace('/auth');
      }
      setIsInitialized(true);
    }, 100);

    return () => clearTimeout(timer);
  }, [hasRedirected, user, tokens, router]);

  // Показываем загрузку пока проверяем авторизацию
  if (!isInitialized) {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-950">
        <div className="text-center">
          <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-blue-500/20 mb-4">
            <div className="w-8 h-8 rounded-full border-2 border-blue-500 border-t-transparent animate-spin"></div>
          </div>
          <p className="text-gray-400">Проверка доступа...</p>
        </div>
      </div>
    );
  }

  if (!user || !tokens?.accessToken) {
    return (
      fallback || (
        <div className="flex items-center justify-center h-screen bg-gray-950">
          <div className="text-center">
            <ShieldAlert className="w-16 h-16 text-gray-600 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-300 mb-2">Требуется авторизация</h3>
            <p className="text-gray-500">Пожалуйста, войдите в систему</p>
          </div>
        </div>
      )
    );
  }

  const hasReq = requiredPermission ? hasPermission(requiredPermission) : true;

  if (requiredPermission && !hasReq){
    return (
      fallback || (
        <div className="flex items-center justify-center h-screen bg-gray-950">
          <div className="text-center max-w-md">
            <ShieldAlert className="w-16 h-16 text-red-600 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-300 mb-2">Доступ запрещен</h3>
            <p className="text-gray-500 mb-4">У вас нет прав для просмотра этого раздела</p>
            
            <div className="bg-yellow-900/30 border border-yellow-700 rounded-lg p-4 mb-4 text-left">
              <p className="text-sm text-yellow-200 mb-2">
                <span className="font-semibold">💡 Что дальше?</span>
              </p>
              <p className="text-xs text-yellow-100 mb-3">
                Обратитесь к администратору системы для получения необходимых прав доступа.
              </p>
              <p className="text-xs text-yellow-100 mb-2">
                Администратор должен:
              </p>
              <ul className="text-xs text-yellow-100 list-disc list-inside space-y-1 ml-1">
                <li>Назначить вам подходящую роль</li>
                <li>Выдать необходимые разрешения</li>
                <li>Вы сможете войти после этого</li>
              </ul>
            </div>

            {user && (
              <details className="text-xs text-left mx-auto">
                <summary className="cursor-pointer text-gray-500 hover:text-gray-300 mb-2">
                  📋 Информация для администратора
                </summary>
                <pre className="mt-2 bg-gray-900 p-3 rounded overflow-auto text-gray-400 text-xs">
{JSON.stringify({
  email: user.email,
  currentRole: user.role,
  requiredPermission,
  userPermissions: user.permissions,
  hasPermissionResult: hasReq,
}, null, 2)}
                </pre>
              </details>
            )}
          </div>
        </div>
      )
    );
  }

  return <>{children}</>;
}
