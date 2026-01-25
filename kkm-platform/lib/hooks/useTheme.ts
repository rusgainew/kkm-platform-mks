"use client";

import { useTheme } from "next-themes";
import { useEffect, useState } from "react";

type Theme = "light" | "dark" | "system";

interface UseThemeReturn {
  theme: Theme | undefined;
  systemTheme: "light" | "dark" | undefined;
  setTheme: (theme: Theme) => void;
  isDark: boolean;
  mounted: boolean;
  toggleTheme: () => void;
}

/**
 * Кастомный хук для работы с темой
 * Обертка над next-themes с дополнительными утилитами
 */
export function useThemeCustom(): UseThemeReturn {
  const { theme, setTheme, systemTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const currentTheme = mounted
    ? theme === "system"
      ? systemTheme
      : theme
    : undefined;

  const isDark = currentTheme === "dark";

  const toggleTheme = () => {
    if (mounted) {
      setTheme(isDark ? "light" : "dark");
    }
  };

  return {
    theme: theme as Theme | undefined,
    systemTheme,
    setTheme,
    isDark,
    mounted,
    toggleTheme,
  };
}
