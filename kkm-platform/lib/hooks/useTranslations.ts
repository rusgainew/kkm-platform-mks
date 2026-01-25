"use client";

import { useCallback } from "react";
import { Locale, defaultLocale } from "@/i18n/config";

// Кэш загруженных переводов
const messageCache: Record<Locale, Record<string, any>> = {
  ru: {},
  en: {},
};

/**
 * Утилита для получения значения из вложенного объекта по точечному пути
 */
function getNestedValue(obj: any, path: string): any {
  const keys = path.split(".");
  let current = obj;

  for (const key of keys) {
    if (current && typeof current === "object" && key in current) {
      current = current[key];
    } else {
      return undefined;
    }
  }

  return current;
}

/**
 * Интерполяция переменных в строку
 * Пример: "Hello {name}" + {name: "John"} = "Hello John"
 */
function interpolate(message: string, variables?: Record<string, any>): string {
  if (!variables) return message;

  return message.replace(/\{(\w+)\}/g, (match, key) => {
    return variables[key]?.toString() || match;
  });
}

/**
 * Хук для получения переводов
 */
export function useTranslations(locale?: Locale) {
  const currentLocale = locale || defaultLocale;

  const getMessages = useCallback(async () => {
    if (
      messageCache[currentLocale] &&
      Object.keys(messageCache[currentLocale]).length > 0
    ) {
      return messageCache[currentLocale];
    }

    try {
      const messages = await import(`@/messages/${currentLocale}.json`);
      messageCache[currentLocale] = messages.default;
      return messages.default;
    } catch (error) {
      console.error(
        `Failed to load translations for locale: ${currentLocale}`,
        error
      );
      return {};
    }
  }, [currentLocale]);

  const t = useCallback(
    async (key: string, variables?: Record<string, any>): Promise<string> => {
      const messages = await getMessages();
      const value = getNestedValue(messages, key);

      if (typeof value !== "string") {
        console.warn(`Translation not found: ${key}`);
        return key;
      }

      return interpolate(value, variables);
    },
    [getMessages]
  );

  // Синхронная версия для простых случаев
  const tSync = useCallback(
    (key: string, variables?: Record<string, any>): string => {
      const messages = messageCache[currentLocale];
      const value = getNestedValue(messages, key);

      if (typeof value !== "string") {
        return key;
      }

      return interpolate(value, variables);
    },
    [currentLocale]
  );

  return { t, tSync };
}

/**
 * Синхронный хук (использует закэшированные переводы)
 */
export function useTranslationsSync(locale?: Locale) {
  const currentLocale = locale || defaultLocale;

  const tSync = useCallback(
    (key: string, variables?: Record<string, any>): string => {
      const messages = messageCache[currentLocale];

      if (!messages || Object.keys(messages).length === 0) {
        return key;
      }

      const value = getNestedValue(messages, key);

      if (typeof value !== "string") {
        return key;
      }

      return interpolate(value, variables);
    },
    [currentLocale]
  );

  return tSync;
}

/**
 * Инициализация кэша переводов при загрузке приложения
 */
export async function initializeTranslations(locale: Locale) {
  try {
    const messages = await import(`@/messages/${locale}.json`);
    messageCache[locale] = messages.default;
  } catch (error) {
    console.error(
      `Failed to initialize translations for locale: ${locale}`,
      error
    );
  }
}
