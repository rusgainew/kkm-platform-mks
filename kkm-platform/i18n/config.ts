/**
 * i18n Configuration
 * Поддерживаемые языки и настройки
 */

export const locales = ["ru", "en"] as const;
export type Locale = (typeof locales)[number];

export const defaultLocale: Locale = "ru";

export const localeNames: Record<Locale, string> = {
  ru: "🇷🇺 Русский",
  en: "🇬🇧 English",
};

// Маппинг язык -> флаг + название
export const localeInfo: Record<
  Locale,
  { name: string; flag: string; nativeName: string }
> = {
  ru: {
    name: "Russian",
    nativeName: "Русский",
    flag: "🇷🇺",
  },
  en: {
    name: "English",
    nativeName: "English",
    flag: "🇬🇧",
  },
};
