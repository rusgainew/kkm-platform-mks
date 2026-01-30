"use client";

import { useEffect, useState } from "react";
import { useAuthStore } from "@/store/authStore";
import { useTokenRefresh } from "@/lib/hooks/useTokenRefresh";

/**
 * Инициализатор авторизации
 * Отвечает за:
 * - Восстановление сессии пользователя при загрузке приложения
 * - Ожидание гидрации состояния из localStorage
 * - Активация автоматического обновления токенов
 * - Отображение загрузки во время инициализации
 */
export default function AuthInitializer({
  children,
}: {
  children?: React.ReactNode;
}) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const user = useAuthStore((state) => state.user);
  const tokens = useAuthStore((state) => state.tokens);
  const [isInitialized, setIsInitialized] = useState(false);
  const [isHydrated, setIsHydrated] = useState(false);

  // Активируем автоматическое обновление токена
  useTokenRefresh();

  /**
   * Ожидаем гидрации Zustand из localStorage
   */
  useEffect(() => {
    console.log("[AuthInitializer] useEffect запущен");
    setIsHydrated(useAuthStore.persist.hasHydrated());
    const unsub = useAuthStore.persist.onFinishHydration(() => {
      console.log("[AuthInitializer] Установлена гидрация = true");
      setIsHydrated(true);
    });

    return () => {
      unsub();
    };
  }, []);

  /**
   * Проверяем статус авторизации после гидрации
   */
  useEffect(() => {
    if (!isHydrated) {
      console.log("[AuthInitializer] Еще не гидрирован, пропускаем проверку");
      return;
    }

    const checkAuth = async () => {
      try {
        console.log(
          "[AuthInitializer] Проверка статуса авторизации после гидрации...",
        );
        console.log("[AuthInitializer] Аутентифицирован:", isAuthenticated);
        console.log("[AuthInitializer] Пользователь:", user?.name || "нет");
        console.log(
          "[AuthInitializer] Токены доступны:",
          !!tokens?.accessToken,
        );

        if (isAuthenticated && user && tokens?.accessToken) {
          console.log("[AuthInitializer] Сессия восстановлена из localStorage");
          console.log(
            "[AuthInitializer] Пользователь:",
            user.name,
            "Роль:",
            user.role,
          );
          console.log(
            "[AuthInitializer] Токен истечет через:",
            tokens.expiresIn,
            "секунд",
          );
          console.log(
            "[AuthInitializer] Длина токена доступа:",
            tokens.accessToken.length,
          );
        } else {
          console.log(
            "[AuthInitializer] Аутентифицированная сессия не найдена",
          );
        }
      } catch (error) {
        console.error("[AuthInitializer] Ошибка проверки авторизации:", error);
      }

      setIsInitialized(true);
    };

    checkAuth();
  }, [isHydrated, isAuthenticated, user, tokens]);

  // Если children не передан, просто инициализируем без рендеринга
  if (!children) {
    return null;
  }

  // Рендерим children только после инициализации
  if (!isInitialized) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-400 mx-auto"></div>
          <p className="text-gray-400 mt-4">Инициализация...</p>
        </div>
      </div>
    );
  }

  return <>{children}</>;
}
