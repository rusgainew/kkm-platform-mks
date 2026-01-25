# 🎯 Шпаргалка: Управление авторизацией

## Быстрый старт

### 1️⃣ Автоматическое восстановление сессии

```typescript
// Работает автоматически! Ничего не нужно делать.
// AuthInitializer обрабатывает всё.
```

### 2️⃣ Проверка авторизации в компоненте

```typescript
import { useAuth } from "@/lib/hooks/useAuth";

const { user, hasToken, hasPermission } = useAuth();

if (!hasToken) {
  return <div>Not authenticated</div>;
}

return <div>Hello {user?.name}</div>;
```

### 3️⃣ Защита страницы

```typescript
import ProtectedRoute from "@/components/auth/ProtectedRoute";

export default function MyPage() {
  return (
    <ProtectedRoute requiredPermission="view_users">
      <MyContent />
    </ProtectedRoute>
  );
}
```

### 4️⃣ Разлогинение

```typescript
import { useAuth } from "@/lib/hooks/useAuth";

function LogoutButton() {
  const { logout } = useAuth();
  return <button onClick={logout}>Logout</button>;
}
```

## 🔧 API интеграция

### Требования к ответу сервера:

**Login response:**

```json
{
  "user": {...},
  "access_token": "jwt...",
  "refresh_token": "jwt...",
  "expires_in": 900
}
```

**Refresh response:**

```json
{
  "access_token": "jwt...",
  "refresh_token": "jwt...",
  "expires_in": 900
}
```

## 📊 Что происходит автоматически

- ✅ Токен обновляется за 1 минуту до истечения
- ✅ Сессия восстанавливается при перезагрузке
- ✅ Неавторизованные пользователи перенаправляются на /auth
- ✅ Токен всегда отправляется в Authorization заголовке

## 🐛 Отладка

```javascript
// В консоли браузера:
import { useAuthStore } from "@/store/authStore";
const { user, tokens, isAuthenticated } = useAuthStore.getState();

// Посмотреть состояние
console.log(user, tokens, isAuthenticated);

// Разлогиниться
useAuthStore.getState().logout();

// Проверить localStorage
JSON.parse(localStorage.getItem("auth-storage"));
```

## 📝 Файлы для изменения

Если нужно кастомизировать:

1. **Время обновления:** `lib/hooks/useTokenRefresh.ts` строка ~36
2. **Сообщения об ошибках:** `components/auth/ProtectedRoute.tsx`
3. **Стиль загрузки:** `components/auth/AuthInitializer.tsx`
4. **API URL:** `lib/api/users.ts` переменная `API_BASE`

## 🚀 Одна строка кода для защиты маршрута

```tsx
// Вместо:
function Page() { if (!user) redirect('/auth'); ... }

// Делай так:
<ProtectedRoute><Page /></ProtectedRoute>
```

---

**Все просто и понятно!** ✨
