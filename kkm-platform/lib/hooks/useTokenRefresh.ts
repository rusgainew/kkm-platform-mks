'use client';

import { useEffect, useCallback, useRef } from 'react';
import { useAuthStore } from '@/store/authStore';
import { useRefreshTokensMutation } from './useAuthApi';
import { isTokenExpired } from '@/lib/utils/tokenValidator';

/**
 * Хук для автоматического обновления JWT токена перед истечением
 * - Обновляет токен за 60 сек до истечения (если есть refresh token)
 * - Проверяет состояние токена каждые 10 сек
 * - Если токен истек, обновляет его сразу же
 * - При ошибке обновления разлогинивает пользователя
 */
export function useTokenRefresh() {
  const { tokens, logout } = useAuthStore();
  const refreshMutation = useRefreshTokensMutation();
  const refreshTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const checkIntervalRef = useRef<NodeJS.Timeout | null>(null);

  /**
   * Планирует обновление токена перед истечением
   */
  const scheduleTokenRefresh = useCallback(() => {
    // Очищаем предыдущий таймер если есть
    if (refreshTimeoutRef.current) {
      clearTimeout(refreshTimeoutRef.current);
    }

    if (!tokens?.accessToken || !tokens?.refreshToken) {
      console.log('[TokenRefresh] Нет токенов для планирования обновления');
      return;
    }

    // Проверяем сейчас истек ли токен
    if (isTokenExpired(tokens.accessToken, 10)) {
      console.warn('[TokenRefresh] Токен уже истек! Обновляем немедленно...');
      setTimeout(async () => {
        try {
          console.log('[TokenRefresh] Обновляем истекший токен...');
          await refreshMutation.mutateAsync(tokens.refreshToken);
          console.log('[TokenRefresh] Истекший токен успешно обновлен');
        } catch (error) {
          console.error('[TokenRefresh] Ошибка обновления истекшего токена:', error);
          logout();
        }
      }, 0);
      return;
    }

    // Получаем время истечения токена (в секундах, по умолчанию 15 минут = 900 сек)
    const expiresIn = tokens.expiresIn || 900;
    const refreshBefore = 60; // Обновляем за 60 сек до истечения
    const delayMs = (expiresIn - refreshBefore) * 1000;

    console.log(
      '[TokenRefresh] Обновление запланировано через',
      Math.round(delayMs / 1000),
      'секунд (токен истечет через',
      expiresIn,
      'секунд)'
    );

    refreshTimeoutRef.current = setTimeout(async () => {
      try {
        console.log('[TokenRefresh] Обновляем токен перед истечением...');
        await refreshMutation.mutateAsync(tokens.refreshToken);
        console.log('[TokenRefresh] Токен успешно обновлен');
        // Планируем следующее обновление
        scheduleTokenRefresh();
      } catch (error) {
        console.error('[TokenRefresh] Ошибка обновления токена:', error);
        logout();
      }
    }, delayMs);
  }, [tokens, refreshMutation, logout]);

  /**
   * Периодическая проверка (каждые 10 сек)
   * Если токен вот-вот истечет, запускаем обновление
   */
  useEffect(() => {
    if (checkIntervalRef.current) {
      clearInterval(checkIntervalRef.current);
    }

    if (!tokens?.accessToken) return;

    checkIntervalRef.current = setInterval(() => {
      if (isTokenExpired(tokens.accessToken, 30)) {
        console.warn('[TokenRefresh] Токен истекает скоро, запускаем обновление...');
        scheduleTokenRefresh();
      }
    }, 10000); // Проверяем каждые 10 секунд

    return () => {
      if (checkIntervalRef.current) {
        clearInterval(checkIntervalRef.current);
      }
    };
  }, [tokens?.accessToken, scheduleTokenRefresh]);

  /**
   * Планируем обновление при изменении токена
   */
  useEffect(() => {
    scheduleTokenRefresh();

    return () => {
      if (refreshTimeoutRef.current) {
        clearTimeout(refreshTimeoutRef.current);
      }
    };
  }, [scheduleTokenRefresh]);
}
