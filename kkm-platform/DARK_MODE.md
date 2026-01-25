# 🌙 Dark Mode Implementation

**Статус:** ✅ Production Ready
**Версия:** 1.0.0
**Последнее обновление:** Session 5

## 📋 Оглавление

1. [Архитектура](#архитектура)
2. [Компоненты](#компоненты)
3. [Использование](#использование)
4. [Хуки](#хуки)
5. [Стили](#стили)
6. [Лучшие практики](#лучшие-практики)
7. [Примеры](#примеры)
8. [Конфигурация](#конфигурация)

---

## Архитектура

### Стек технологий

- **next-themes** - Управление темой с система предпочтением
- **Tailwind CSS** - Стили с поддержкой `dark:` варианта
- **React Context** - Предоставление темы компонентам
- **localStorage** - Сохранение выбора пользователя

### Поток данных

```
App
└── ThemeProvider (next-themes)
    └── Layout
        └── Компоненты
            └── useTheme() hook
                └── theme, setTheme, isDark
```

### Хранилище данных

| Хранилище    | Ключ        | Формат                        | TTL        |
| ------------ | ----------- | ----------------------------- | ---------- |
| localStorage | `kkm-theme` | 'light' \| 'dark' \| 'system' | Persistent |

---

## Компоненты

### 1. ThemeProvider

**Файл:** `components/theme/ThemeProvider.tsx`

Обертка around next-themes для всего приложения.

```tsx
import { ThemeProvider } from "@/components/theme/ThemeProvider";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html>
      <body>
        <ThemeProvider>{children}</ThemeProvider>
      </body>
    </html>
  );
}
```

**Конфигурация:**

```tsx
const ThemeProvider = ({ children }: { children: React.ReactNode }) => {
  return (
    <NextThemesProvider
      attribute="class" // Применяет class на <html>
      defaultTheme="system" // По умолчанию система
      enableSystem={true} // Следит за системой
      storageKey="kkm-theme" // Ключ localStorage
      enableColorScheme={false} // Не меняет color-scheme
    >
      {children}
    </NextThemesProvider>
  );
};
```

**Параметры:**

| Параметр       | Тип     | Описание                          |
| -------------- | ------- | --------------------------------- |
| `attribute`    | string  | 'class' - добавляет class 'dark'  |
| `defaultTheme` | string  | 'system', 'light', 'dark'         |
| `enableSystem` | boolean | Следить за системной темой        |
| `storageKey`   | string  | Ключ в localStorage               |
| `forcedTheme`  | string  | Принудительная тема (опционально) |

---

### 2. ThemeSwitcher

**Файл:** `components/theme/ThemeSwitcher.tsx`

Три варианта переключателя темы.

#### 2.1 ThemeSwitcher (Полный)

Кнопки выбора темы в ряд.

```tsx
import { ThemeSwitcher } from "@/components/theme/ThemeSwitcher";

export default function Header() {
  return <ThemeSwitcher />;
}
```

**Внешний вид:**

```
☀️ Light | 🌙 Dark | ⚙️ Auto
```

**Свойства:**

```tsx
interface ThemeSwitcherProps {
  className?: string;
}
```

**Пример:**

```tsx
<ThemeSwitcher className="gap-2" />
```

---

#### 2.2 ThemeSwitcherCompact (Компактный)

Одна кнопка для переключения между светлой и темной.

```tsx
import { ThemeSwitcherCompact } from "@/components/theme/ThemeSwitcher";

export default function Header() {
  return <ThemeSwitcherCompact />;
}
```

**Внешний вид:**

```
☀️          (Light)
🌙          (Dark, при наведении)
```

**Свойства:**

```tsx
interface ThemeSwitcherCompactProps {
  className?: string;
}
```

---

#### 2.3 ThemeSwitcherDropdown (Меню)

Выпадающее меню для выбора темы.

```tsx
import { ThemeSwitcherDropdown } from "@/components/theme/ThemeSwitcher";

export default function Header() {
  return <ThemeSwitcherDropdown />;
}
```

**Внешний вид:**

```
[⚙️ System ▼]
  ☀️ Light   ✓
  🌙 Dark
  ⚙️ Auto
```

**Свойства:**

```tsx
interface ThemeSwitcherDropdownProps {
  className?: string;
}
```

---

### 3. ThemeExamples

**Файл:** `components/theme/ThemeExamples.tsx`

Примеры компонентов с поддержкой темной темы.

#### 3.1 ThemeExampleCard

Карточка с поддержкой темной темы.

```tsx
<ThemeExampleCard />
```

**Стили:**

```tsx
<div className="bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700">
  {/* Контент */}
</div>
```

---

#### 3.2 ThemeExampleButton

Кнопки трех вариантов.

```tsx
<ThemeExampleButton />
```

**Варианты:**

1. **Primary** - Синяя (bg-blue-500 dark:bg-blue-600)
2. **Secondary** - Серая (bg-gray-200 dark:bg-gray-700)
3. **Danger** - Красная (bg-red-500 dark:bg-red-600)

---

#### 3.3 ThemeExampleForm

Форма с полями ввода.

```tsx
<ThemeExampleForm />
```

**Поля:**

- Input (текст)
- Textarea (многострочный текст)
- Select (выпадающий список)
- Checkbox (флажок)

---

#### 3.4 ThemeExampleGradient

Градиент с поддержкой темной темы.

```tsx
<ThemeExampleGradient />
```

**Стили:**

```tsx
<div className="bg-gradient-to-r from-blue-500 dark:from-blue-600 to-purple-500 dark:to-purple-600">
  {/* Контент */}
</div>
```

---

#### 3.5 ThemeExampleTable

Таблица с поддержкой темной темы.

```tsx
<ThemeExampleTable />
```

**Стили:**

```tsx
<tr className="hover:bg-gray-100 dark:hover:bg-gray-800">{/* Ячейки */}</tr>
```

---

#### 3.6 ThemeExampleAlert

Алерты четырех типов.

```tsx
<ThemeExampleAlert />
```

**Типы:**

1. **Info** - Синяя (bg-blue-50 dark:bg-blue-900/20)
2. **Success** - Зеленая (bg-green-50 dark:bg-green-900/20)
3. **Warning** - Желтая (bg-yellow-50 dark:bg-yellow-900/20)
4. **Error** - Красная (bg-red-50 dark:bg-red-900/20)

---

## Использование

### 1. Получение текущей темы

```tsx
"use client";

import { useThemeCustom } from "@/lib/hooks/useTheme";

export default function MyComponent() {
  const { theme, isDark, mounted } = useThemeCustom();

  if (!mounted) return null;

  return <div>{isDark ? "🌙 Dark" : "☀️ Light"}</div>;
}
```

---

### 2. Изменение темы

```tsx
"use client";

import { useThemeCustom } from "@/lib/hooks/useTheme";

export default function ThemeSwitcher() {
  const { setTheme, theme } = useThemeCustom();

  return <button onClick={() => setTheme("dark")}>Switch to Dark</button>;
}
```

---

### 3. Переключение между темами

```tsx
"use client";

import { useThemeCustom } from "@/lib/hooks/useTheme";

export default function ToggleTheme() {
  const { toggleTheme } = useThemeCustom();

  return <button onClick={toggleTheme}>Toggle Theme</button>;
}
```

---

### 4. Условный рендеринг

```tsx
"use client";

import { useThemeCustom } from "@/lib/hooks/useTheme";

export default function Conditional() {
  const { isDark, mounted } = useThemeCustom();

  if (!mounted) return <div>Loading...</div>;

  return isDark ? <DarkComponent /> : <LightComponent />;
}
```

---

## Хуки

### useThemeCustom()

**Файл:** `lib/hooks/useTheme.ts`

Хук для управления темой.

**Возвращает:**

```typescript
interface UseThemeReturn {
  theme: string | undefined; // 'light', 'dark', 'system'
  systemTheme: string | undefined; // 'light' или 'dark' (система)
  setTheme: (theme: string) => void; // Установить тему
  isDark: boolean; // true если темная тема
  mounted: boolean; // true когда готово
  toggleTheme: () => void; // Переключить light <-> dark
}
```

**Пример:**

```tsx
const {
  theme, // 'light' | 'dark' | 'system'
  systemTheme, // 'light' | 'dark'
  setTheme, // (theme: string) => void
  isDark, // boolean
  mounted, // boolean
  toggleTheme, // () => void
} = useThemeCustom();
```

---

## Стили

### Темные варианты в Tailwind CSS

**Синтаксис:**

```tsx
<div className="bg-white dark:bg-gray-800">
  {/* Светло: белый фон, Темно: серый */}
</div>
```

**Распространенные варианты:**

| Светлая           | Темная                 | Использование   |
| ----------------- | ---------------------- | --------------- |
| `bg-white`        | `dark:bg-gray-800`     | Фон блока       |
| `text-gray-900`   | `dark:text-white`      | Основной текст  |
| `text-gray-600`   | `dark:text-gray-400`   | Вторичный текст |
| `border-gray-200` | `dark:border-gray-700` | Границы         |
| `bg-gray-50`      | `dark:bg-gray-900`     | Фон секции      |

### Полные примеры

**Карточка:**

```tsx
<div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 shadow-sm hover:shadow-md dark:hover:shadow-lg p-6">
  <h3 className="text-lg font-bold text-gray-900 dark:text-white">Заголовок</h3>
  <p className="text-gray-600 dark:text-gray-400 mt-2">Описание</p>
</div>
```

**Кнопка:**

```tsx
<button
  className="
  bg-blue-500 hover:bg-blue-600 active:bg-blue-700
  dark:bg-blue-600 dark:hover:bg-blue-700 dark:active:bg-blue-800
  text-white font-semibold py-2 px-4 rounded-lg
  transition-colors
"
>
  Кнопка
</button>
```

**Входной элемент:**

```tsx
<input
  className="
    bg-white dark:bg-gray-800
    text-gray-900 dark:text-white
    border border-gray-300 dark:border-gray-600
    focus:border-blue-500 dark:focus:border-blue-400
    rounded-lg px-4 py-2
    placeholder-gray-500 dark:placeholder-gray-400
  "
  type="text"
  placeholder="Ввод текста"
/>
```

---

## Лучшие практики

### ✅ Рекомендуется

1. **Используйте Tailwind dark: префикс**

   ```tsx
   <div className="bg-white dark:bg-gray-800">...</div>
   ```

2. **Тестируйте оба варианта**

   - Светлую тему
   - Темную тему
   - Системную тему

3. **Проверяйте контраст текста**

   - WCAG AA (4.5:1 для основного текста)
   - WCAG AAA (7:1 для максимальной доступности)

4. **Используйте семантические цвета**

   ```tsx
   // ✓ Good
   <div className="bg-white dark:bg-gray-800">

   // ✗ Bad
   <div style={{ backgroundColor: '#ffffff' }}>
   ```

5. **Избегайте жестко закодированных цветов**

   ```tsx
   // ✓ Good - Tailwind
   className="text-gray-900 dark:text-white"

   // ✗ Bad - Hardcoded
   style={{ color: '#000000' }}
   ```

6. **Используйте useTheme для управления**
   ```tsx
   const { isDark } = useThemeCustom();
   ```

---

### ✕ Избегайте

1. **Не игнорируйте систему**

   ```tsx
   // ✗ Не делайте так
   <div className="bg-white">  // Всегда белый фон

   // ✓ Делайте так
   <div className="bg-white dark:bg-gray-800">
   ```

2. **Не используйте низкий контраст**

   ```tsx
   // ✗ Низкий контраст
   <p className="text-gray-600 dark:text-gray-500">

   // ✓ Хороший контраст
   <p className="text-gray-700 dark:text-gray-300">
   ```

3. **Не забывайте о border и shadow**

   ```tsx
   // ✓ Полная поддержка
   <div className="
     bg-white dark:bg-gray-800
     border-gray-200 dark:border-gray-700
     shadow-sm dark:shadow-lg
   ">
   ```

4. **Не забывайте о проверке mounted**

   ```tsx
   // ✓ Правильно
   if (!mounted) return null;

   // ✗ Неправильно
   return isDark ? <Dark /> : <Light />; // Hydration mismatch
   ```

---

## Примеры

### Пример 1: Простой компонент с темой

```tsx
"use client";

import { useThemeCustom } from "@/lib/hooks/useTheme";

export default function Card() {
  const { isDark, mounted } = useThemeCustom();

  if (!mounted) return null;

  return (
    <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm">
      <h2 className="text-xl font-bold text-gray-900 dark:text-white">
        {isDark ? "🌙 Dark" : "☀️ Light"}
      </h2>
      <p className="text-gray-600 dark:text-gray-400 mt-2">
        Текущая тема: {isDark ? "Темная" : "Светлая"}
      </p>
    </div>
  );
}
```

---

### Пример 2: Переключатель темы в Header

```tsx
"use client";

import { ThemeSwitcher } from "@/components/theme/ThemeSwitcher";

export default function Header() {
  return (
    <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      <div className="flex justify-between items-center p-4">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
          MyApp
        </h1>
        <ThemeSwitcher />
      </div>
    </header>
  );
}
```

---

### Пример 3: Таблица с темой

```tsx
export default function Table() {
  return (
    <table className="w-full border-collapse">
      <thead className="bg-gray-100 dark:bg-gray-800">
        <tr>
          <th className="text-left p-3 text-gray-900 dark:text-white font-bold">
            Имя
          </th>
        </tr>
      </thead>
      <tbody>
        <tr className="hover:bg-gray-100 dark:hover:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
          <td className="p-3 text-gray-600 dark:text-gray-400">Иван Петров</td>
        </tr>
      </tbody>
    </table>
  );
}
```

---

### Пример 4: Модальное окно с темой

```tsx
export default function Modal() {
  return (
    <div className="bg-black/50 dark:bg-black/70 fixed inset-0 flex items-center justify-center">
      <div className="bg-white dark:bg-gray-800 p-8 rounded-lg shadow-lg max-w-md">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">
          Заголовок
        </h2>
        <p className="text-gray-600 dark:text-gray-400 mb-6">
          Текст модального окна
        </p>
        <button className="bg-blue-500 dark:bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-600 dark:hover:bg-blue-700">
          Закрыть
        </button>
      </div>
    </div>
  );
}
```

---

## Конфигурация

### Tailwind Config

**File:** `tailwindconfig.ts`

```typescript
import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: "class", // ← Включает dark: варианты
  content: ["./app/**/*.{js,ts,jsx,tsx}", "./components/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {},
  },
  plugins: [],
};

export default config;
```

### Next.js Config

**File:** `next.config.ts`

```typescript
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Никаких специальных конфигов не требуется
  // next-themes работает с любой конфигурацией
};

export default nextConfig;
```

### Root Layout

**File:** `app/layout.tsx`

```tsx
import type { Metadata } from "next";
import { ThemeProvider } from "@/components/theme/ThemeProvider";
import { ReactNode } from "react";

export const metadata: Metadata = {
  title: "My App",
  description: "With dark mode support",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html suppressHydrationWarning>
      <body>
        <ThemeProvider>{children}</ThemeProvider>
      </body>
    </html>
  );
}
```

---

## Часто задаваемые вопросы

### Q: Как избежать мерцания при загрузке?

**A:** next-themes автоматически предотвращает мерцание. Просто добавьте `suppressHydrationWarning` в тег `<html>`:

```tsx
<html suppressHydrationWarning>
```

---

### Q: Как сохранить выбор пользователя?

**A:** next-themes автоматически сохраняет выбор в localStorage под ключом `kkm-theme`.

---

### Q: Как принудительно установить тему?

**A:** Используйте `forcedTheme` в ThemeProvider:

```tsx
<NextThemesProvider forcedTheme="dark">{children}</NextThemesProvider>
```

---

### Q: Как получить системную тему?

**A:** Используйте `systemTheme` из `useThemeCustom()`:

```tsx
const { systemTheme } = useThemeCustom();
console.log(systemTheme); // 'light' или 'dark'
```

---

### Q: Как обнаружить изменение темы?

**A:** next-themes емитует событие `themechange`:

```tsx
useEffect(() => {
  const handleThemeChange = (e: StorageEvent) => {
    if (e.key === "kkm-theme") {
      console.log("Theme changed:", e.newValue);
    }
  };

  window.addEventListener("storage", handleThemeChange);
  return () => window.removeEventListener("storage", handleThemeChange);
}, []);
```

---

## Ресурсы

- [next-themes Documentation](https://github.com/pacocoursey/next-themes)
- [Tailwind CSS Dark Mode](https://tailwindcss.com/docs/dark-mode)
- [WCAG Contrast Guidelines](https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum)
- [Material Design Dark Theme](https://material.io/design/color/dark-theme.html)

---

## Статистика

| Метрика             | Значение                                       |
| ------------------- | ---------------------------------------------- |
| Компонентов         | 3 (ThemeSwitcher варианта)                     |
| Примеров            | 6 (Card, Button, Form, Gradient, Table, Alert) |
| Страниц             | 1 (/[locale]/theme-demo)                       |
| Хуков               | 1 (useThemeCustom)                             |
| Ключей localStorage | 1 (kkm-theme)                                  |
| Поддерживаемых тем  | 3 (light, dark, system)                        |

---

## История изменений

| Версия | Дата      | Изменения     |
| ------ | --------- | ------------- |
| 1.0.0  | Session 5 | Первая версия |

---

**Документ подготовлен:** Session 5, Phase D
**Статус:** Production Ready ✅
**Последнее обновление:** [Current Date]
