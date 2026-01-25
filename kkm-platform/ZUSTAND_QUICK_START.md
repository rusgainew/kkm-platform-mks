# 🚀 QUICK START - ZUSTAND STORES

Полное руководство по использованию Zustand stores в компонентах.

---

## 📦 БЫСТРЫЙ СТАРТ

### 1. Основные Imports

```typescript
// Вариант 1: Импортировать custom hook (РЕКОМЕНДУЕТСЯ)
import { useAuth, useUsers, useInvoices, useCatalog, useUI } from "@/hooks";

// Вариант 2: Импортировать Zustand store напрямую (продвинутое)
import useAuthStore from "@/store/auth/authStore";
import useInvoiceStore from "@/store/invoice/invoiceStore";
```

---

## 🔐 АУТЕНТИФИКАЦИЯ (useAuth)

### Базовое использование

```typescript
import { useAuth } from "@/hooks";

export function LoginPage() {
  const { login, isLoading, error, isAuthenticated } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const success = await login(email, password);
      if (success) {
        // Автоматическое перенаправление на /dashboard
        window.location.href = "/dashboard";
      }
    } catch (err) {
      console.error("Login failed:", err);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {error && <div className="error-alert">{error}</div>}
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
      />
      <button disabled={isLoading}>
        {isLoading ? "Logging in..." : "Login"}
      </button>
    </form>
  );
}
```

### Проверка аутентификации

```typescript
export function Dashboard() {
  const { user, isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    return <Redirect to="/login" />;
  }

  return (
    <div>
      <h1>Welcome, {user?.name}!</h1>
      <p>Role: {user?.role}</p>
    </div>
  );
}
```

### Refresh Token

```typescript
export function useApiInterceptor() {
  const { refreshToken } = useAuth();

  // В lib/http.ts
  axiosInstance.interceptors.response.use(
    (response) => response,
    async (error) => {
      if (error.response?.status === 401) {
        const success = await refreshToken();
        if (success) {
          // Retry original request
          return axiosInstance(error.config);
        }
      }
      return Promise.reject(error);
    }
  );
}
```

---

## 👥 ПОЛЬЗОВАТЕЛИ (useUsers)

### Получение списка пользователей

```typescript
import { useUsers } from "@/hooks";

export function UsersList() {
  const { users, filteredUsers, isLoading, fetchUsers, setFilters } =
    useUsers();

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleSearch = (search) => {
    setFilters({ search });
  };

  const handleFilterByRole = (role) => {
    setFilters({ role });
  };

  return (
    <div>
      <input
        onChange={(e) => handleSearch(e.target.value)}
        placeholder="Search users..."
      />
      <select onChange={(e) => handleFilterByRole(e.target.value)}>
        <option value="">All roles</option>
        <option value="admin">Admin</option>
        <option value="manager">Manager</option>
        <option value="cashier">Cashier</option>
      </select>

      {isLoading ? (
        <Spinner />
      ) : (
        <table>
          <tbody>
            {filteredUsers.map((user) => (
              <tr key={user.id}>
                <td>{user.name}</td>
                <td>{user.email}</td>
                <td>{user.role}</td>
                <td>{user.status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
```

### Создание пользователя

```typescript
export function CreateUserForm({ onSuccess }) {
  const { createUser, isLoading, error } = useUsers();
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
    role: "user",
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    const success = await createUser(formData);
    if (success) {
      onSuccess?.();
      setFormData({ name: "", email: "", password: "", role: "user" });
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {error && <div className="error">{error}</div>}
      <input
        value={formData.name}
        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
        placeholder="Name"
      />
      <input
        type="email"
        value={formData.email}
        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
        placeholder="Email"
      />
      <input
        type="password"
        value={formData.password}
        onChange={(e) => setFormData({ ...formData, password: e.target.value })}
        placeholder="Password"
      />
      <select
        value={formData.role}
        onChange={(e) => setFormData({ ...formData, role: e.target.value })}
      >
        <option value="user">User</option>
        <option value="manager">Manager</option>
        <option value="cashier">Cashier</option>
        <option value="admin">Admin</option>
      </select>
      <button disabled={isLoading}>
        {isLoading ? "Creating..." : "Create User"}
      </button>
    </form>
  );
}
```

### Обновление пользователя

```typescript
export function EditUserForm({ userId }) {
  const { updateUser, isLoading } = useUsers();
  const [formData, setFormData] = useState({});

  const handleSubmit = async (e) => {
    e.preventDefault();
    const success = await updateUser(userId, formData);
    if (success) {
      // User updated successfully
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* Form fields */}
      <button disabled={isLoading}>
        {isLoading ? "Updating..." : "Update User"}
      </button>
    </form>
  );
}
```

### Удаление пользователя

```typescript
export function DeleteUserButton({ userId, userName }) {
  const { deleteUser, isLoading } = useUsers();
  const { addNotification } = useUI();

  const handleDelete = async () => {
    if (window.confirm(`Delete ${userName}?`)) {
      const success = await deleteUser(userId);
      if (success) {
        addNotification({
          type: "success",
          message: "User deleted successfully",
          duration: 3000,
        });
      }
    }
  };

  return (
    <button onClick={handleDelete} disabled={isLoading}>
      {isLoading ? "Deleting..." : "Delete"}
    </button>
  );
}
```

---

## 💰 СЧЕТА (useInvoices)

### Список счетов с фильтрацией

```typescript
import { useInvoices } from "@/hooks";

export function InvoicesList() {
  const {
    invoices,
    totalRevenue,
    paidAmount,
    pendingAmount,
    isLoading,
    filters,
    fetchInvoices,
    setFilters,
  } = useInvoices();

  useEffect(() => {
    fetchInvoices({ status: "pending" });
  }, []);

  return (
    <div>
      <div className="stats">
        <div>Total Revenue: ${totalRevenue}</div>
        <div>Paid: ${paidAmount}</div>
        <div>Pending: ${pendingAmount}</div>
      </div>

      <select
        value={filters.status || ""}
        onChange={(e) => setFilters({ status: e.target.value || undefined })}
      >
        <option value="">All statuses</option>
        <option value="draft">Draft</option>
        <option value="pending">Pending</option>
        <option value="paid">Paid</option>
        <option value="overdue">Overdue</option>
      </select>

      {isLoading ? (
        <Spinner />
      ) : (
        <table>
          <tbody>
            {invoices.map((invoice) => (
              <tr key={invoice.id}>
                <td>{invoice.number}</td>
                <td>{invoice.companyName}</td>
                <td>${invoice.totalAmount}</td>
                <td>{invoice.status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
```

### Создание счета

```typescript
export function CreateInvoiceForm() {
  const { createInvoice, isLoading } = useInvoices();
  const { addNotification } = useUI();
  const [items, setItems] = useState([
    { productId: "", quantity: 0, price: 0 },
  ]);
  const [formData, setFormData] = useState({
    number: "",
    companyId: "",
    items: [],
  });

  const total = useMemo(
    () => items.reduce((sum, item) => sum + item.quantity * item.price, 0),
    [items]
  );

  const handleSubmit = async (e) => {
    e.preventDefault();
    const success = await createInvoice({
      ...formData,
      items,
      totalAmount: total,
    });

    if (success) {
      addNotification({
        type: "success",
        message: "Invoice created successfully",
        duration: 3000,
      });
      // Reset form or redirect
    }
  };

  const addItem = () => {
    setItems([...items, { productId: "", quantity: 0, price: 0 }]);
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        value={formData.number}
        onChange={(e) => setFormData({ ...formData, number: e.target.value })}
        placeholder="Invoice number"
      />

      <select
        value={formData.companyId}
        onChange={(e) =>
          setFormData({ ...formData, companyId: e.target.value })
        }
      >
        <option>Select company</option>
        {/* Load companies from useCatalog or similar */}
      </select>

      <div>
        <h3>Items (Total: ${total.toFixed(2)})</h3>
        {items.map((item, idx) => (
          <div key={idx} className="item-row">
            <input
              type="number"
              value={item.quantity}
              onChange={(e) => {
                const newItems = [...items];
                newItems[idx].quantity = Number(e.target.value);
                setItems(newItems);
              }}
              placeholder="Qty"
            />
            <input
              type="number"
              value={item.price}
              onChange={(e) => {
                const newItems = [...items];
                newItems[idx].price = Number(e.target.value);
                setItems(newItems);
              }}
              placeholder="Price"
            />
            <button
              type="button"
              onClick={() => setItems(items.filter((_, i) => i !== idx))}
            >
              Remove
            </button>
          </div>
        ))}
        <button type="button" onClick={addItem}>
          Add Item
        </button>
      </div>

      <button disabled={isLoading}>
        {isLoading ? "Creating..." : "Create Invoice"}
      </button>
    </form>
  );
}
```

---

## 📦 КАТАЛОГ (useCatalog)

### Список товаров

```typescript
import { useCatalog } from "@/hooks";

export function CatalogList() {
  const {
    items,
    categories,
    lowStockItems,
    totalInventoryValue,
    isLoading,
    fetchCatalog,
  } = useCatalog();

  useEffect(() => {
    fetchCatalog({ active: true });
  }, []);

  return (
    <div>
      <div className="alerts">
        <div>Inventory Value: ${totalInventoryValue.toFixed(2)}</div>
        {lowStockItems.length > 0 && (
          <div className="warning">
            Low stock on {lowStockItems.length} items
          </div>
        )}
      </div>

      {isLoading ? (
        <Spinner />
      ) : (
        <table>
          <tbody>
            {items.map((item) => (
              <tr key={item.id}>
                <td>{item.name}</td>
                <td>{item.sku}</td>
                <td>${item.price}</td>
                <td>Stock: {item.stock}</td>
                <td>{item.stock <= item.reorderLevel ? "⚠️" : "✓"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
```

---

## 🎨 UI СОСТОЯНИЕ (useUI)

### Управление темой

```typescript
import { useUI } from "@/hooks";

export function Header() {
  const { theme, isDarkMode, toggleTheme, sidebarOpen, toggleSidebar } =
    useUI();

  return (
    <header className={isDarkMode ? "dark" : "light"}>
      <button onClick={toggleSidebar}>
        {sidebarOpen ? "✕ Close" : "☰ Menu"}
      </button>

      <button onClick={toggleTheme}>
        {isDarkMode ? "☀️ Light" : "🌙 Dark"}
      </button>
    </header>
  );
}
```

### Уведомления

```typescript
export function ActionButton() {
  const { addNotification } = useUI();

  const handleClick = () => {
    // Do something...

    addNotification({
      type: "success",
      message: "Action completed successfully!",
      duration: 3000,
    });
  };

  return <button onClick={handleClick}>Do Action</button>;
}
```

```typescript
export function ErrorHandler() {
  const { addNotification } = useUI();

  const handleError = (error) => {
    addNotification({
      type: "error",
      message: error.message,
      duration: 5000,
    });
  };

  return <>...</>;
}
```

### Модальные окна

```typescript
export function ConfirmDialog() {
  const { openModal, closeModal, modals } = useUI();
  const isOpen = modals["confirm-delete"]?.isOpen;

  const handleConfirm = () => {
    // Do something
    closeModal("confirm-delete");
  };

  return (
    <>
      <button
        onClick={() =>
          openModal("confirm-delete", {
            title: "Confirm Delete",
            content: "Are you sure?",
            onConfirm: handleConfirm,
          })
        }
      >
        Delete
      </button>

      {isOpen && (
        <Modal>
          <h2>Confirm Delete</h2>
          <p>Are you sure you want to delete this?</p>
          <button onClick={handleConfirm}>Confirm</button>
          <button onClick={() => closeModal("confirm-delete")}>Cancel</button>
        </Modal>
      )}
    </>
  );
}
```

---

## 🔄 КОМБИНИРОВАНИЕ НЕСКОЛЬКО HOOKS

```typescript
export function DashboardPage() {
  const { user } = useAuth();
  const { invoices, totalRevenue } = useInvoices();
  const { items, lowStockItems } = useCatalog();
  const { addNotification } = useUI();

  return (
    <div>
      <h1>Welcome {user?.name}</h1>

      <div className="cards">
        <Card title="Revenue" value={`$${totalRevenue}`} />
        <Card title="Stock Issues" value={lowStockItems.length} />
        <Card
          title="Pending Invoices"
          value={invoices.filter((i) => i.status === "pending").length}
        />
      </div>

      <button
        onClick={() =>
          addNotification({
            type: "info",
            message: "Dashboard updated",
            duration: 2000,
          })
        }
      >
        Refresh
      </button>
    </div>
  );
}
```

---

## ✅ BEST PRACTICES

1. **Используйте custom hooks** - `useAuth()` вместо `useAuthStore()`
2. **Обрабатывайте errors** - Всегда проверяйте `error` состояние
3. **Показывайте loading** - Используйте `isLoading` флаг
4. **Фильтруйте данные** - Используйте `setFilters()`
5. **Используйте computed** - `totalRevenue`, `lowStockItems`, и т.д.
6. **Добавляйте notifications** - Для feedback пользователю
7. **Не мутируйте состояние** - Используйте setter функции

---

## 🐛 ОТЛАДКА

### Включить DevTools

```typescript
// Zustand DevTools автоматически включены в разработке
// Откройте Redux DevTools в браузере

import { devtools } from "zustand/middleware";

// Store создается с devtools middleware
// Можно видеть все действия и состояния
```

### Логирование

```typescript
// useAuth.ts
const { user, isAuthenticated } = useAuth();
console.log("Auth state:", { user, isAuthenticated });

// Сможете видеть состояние в консоли и DevTools
```

---

**Все готово к использованию!** 🎉
