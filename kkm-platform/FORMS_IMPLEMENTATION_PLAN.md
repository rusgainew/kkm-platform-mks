# 📋 План реализации форм в kkm-platform

**Дата:** 27 января 2026 г.  
**Цель:** Привести в соответствие все формы фронтенда с бэкенд API  
**Подход:** Backend schemas → Frontend types → Form components → Validation

---

## 🎯 Обзор сервисов и форм

### Существующие микросервисы (Backend)

1. **user-server** / **user-query-server** - Пользователи
2. **company-server** / **company-query-server** - Компании (Организации)
3. **invoice-server** / **invoice-query-server** - Счета-фактуры (ЭСФ)
4. **catalog-server** / **catalog-query-server** - Каталог товаров
5. **bank-account-server** / **bank-account-query-server** - Банковские счета
6. **document-server** / **document-query-server** - Документы
7. **foreign-company-server** / **foreign-company-query-server** - Иностранные компании

### Формы фронтенда (Status)

- ✅ **Users** - Готово (`features/users/components/UserCreateForm.tsx`)
- ⚠️ **Companies** - Частично (`features/companies/components/CompanyForm.tsx` - не все поля)
- ⚠️ **Invoices** - Частично (`features/invoices/components/InvoiceForm.tsx` - упрощенная версия)
- ❌ **Catalog** - Нет формы создания/редактирования
- ❌ **Bank Accounts** - Нет форм
- ❌ **Documents** - Есть базовая форма, но не все поля
- ❌ **Foreign Companies** - Нет форм

---

## 📊 ФАЗА 1: Аудит и структура типов (4-6 часов)

### 1.1. Создать полную типизацию для всех сущностей

**Файл:** `types/entities.ts` (новый, централизованный)

```typescript
/**
 * Централизованные типы для всех сущностей системы
 * Соответствуют API Gateway models
 */

// ============================================================================
// USER TYPES
// ============================================================================

export interface User {
  user_id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: UserRole;
  is_active: boolean;
  created_at: number;
  updated_at: number;
  status: string;
}

export type UserRole = "admin" | "manager" | "cashier" | "employee";

export interface CreateUserRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  role: UserRole;
}

export interface UpdateUserRequest {
  first_name?: string;
  last_name?: string;
  role?: UserRole;
  is_active?: boolean;
}

// ============================================================================
// COMPANY (ORGANIZATION) TYPES
// ============================================================================

export interface Company {
  id: string;
  name: string;
  description: string;
  owner_id: string;
  member_count: number;
  created_at: number;
  updated_at: number;
  status: CompanyStatus;
}

export type CompanyStatus = "active" | "inactive" | "suspended";

export interface CreateCompanyRequest {
  name: string;
  description: string;
}

export interface UpdateCompanyRequest {
  name?: string;
  description?: string;
  status?: CompanyStatus;
}

export interface AddMemberRequest {
  organization_id: string;
  user_id: string;
  role: "admin" | "manager" | "employee";
}

export interface Employee {
  id: string;
  user_id: string;
  organization_id: string;
  role: string;
  status: string;
  position?: string;
  department?: string;
  joined_at: number;
  last_active_at?: number;
}

// ============================================================================
// CATALOG TYPES
// ============================================================================

export interface CatalogItem {
  id: string;
  name: string;
  number: string; // Номер по каталогу
  description?: string; // Описание
  tnved_code: string; // Код ТНВЭД (обязательно!)
  category?: string; // @deprecated - использовать tnved_code
  price: number;
  currency: string;
  unit: string; // Единица измерения
  created_at: number;
  updated_at: number;
}

export interface CreateCatalogItemRequest {
  name: string;
  number: string;
  description?: string;
  tnved_code: string; // Обязательно для ЭСФ
  price: number;
  currency: string; // "KGS", "USD", "RUB"
  unit: string; // "шт", "кг", "м", "л" и т.д.
}

export interface UpdateCatalogItemRequest {
  name?: string;
  number?: string;
  description?: string;
  tnved_code?: string;
  price?: number;
  currency?: string;
  unit?: string;
}

// ============================================================================
// BANK ACCOUNT TYPES
// ============================================================================

export interface BankAccount {
  id: string;
  account_number: string;
  bank_name: string;
  bank_code: string; // БИК банка
  currency: string;
  owner_id: string;
  is_active: boolean;
  created_at: number;
  updated_at: number;
}

export interface CreateBankAccountRequest {
  account_number: string; // 20 цифр (для КР)
  bank_name: string;
  bank_code: string; // БИК
  currency: string; // "KGS", "USD", "RUB"
  owner_id: string; // ID организации
}

export interface UpdateBankAccountRequest {
  account_number?: string;
  bank_name?: string;
  bank_code?: string;
  currency?: string;
  is_active?: boolean;
}

// ============================================================================
// DOCUMENT TYPES
// ============================================================================

export interface Document {
  id: string;
  organization_id: string;
  title: string;
  content: string;
  status: DocumentStatus;
  created_by: string;
  assigned_to?: string;
  created_at: number;
  updated_at: number;
  status_changed_at: number;
  version: number;
  entries?: DocumentEntry[];
}

export type DocumentStatus =
  | "draft"
  | "pending"
  | "approved"
  | "rejected"
  | "archived";

export interface DocumentEntry {
  id: string;
  document_id: string;
  key: string;
  value: string;
  created_at: number;
  updated_at: number;
}

export interface CreateDocumentRequest {
  organization_id: string;
  title: string;
  content: string;
  assigned_to?: string;
  entries?: Array<{
    key: string;
    value: string;
  }>;
}

export interface UpdateDocumentRequest {
  title?: string;
  content?: string;
  status?: DocumentStatus;
  assigned_to?: string;
  entries?: Array<{
    key: string;
    value: string;
  }>;
}

// ============================================================================
// FOREIGN COMPANY TYPES
// ============================================================================

export interface ForeignCompany {
  id: number;
  pin: string; // Иностранный ИНН/Tax ID
  full_name: string;
  country_code: string; // ISO 3166-1 alpha-2 (RU, CN, US...)
  address?: string;
  created_at: number;
  updated_at: number;
}

export interface CreateForeignCompanyRequest {
  pin: string;
  full_name: string;
  country_code: string; // Обязательно 2 символа (ISO)
  address?: string;
}

export interface UpdateForeignCompanyRequest {
  pin?: string;
  full_name?: string;
  country_code?: string;
  address?: string;
}

// ============================================================================
// PAGINATION & RESPONSES
// ============================================================================

export interface PageInfo {
  page: number;
  size: number;
  total_count: number;
}

export interface ListResponse<T> {
  items: T[];
  page_info: PageInfo;
}

export interface APIResponse<T = any> {
  success: boolean;
  data?: T;
  error?: APIError;
  meta?: MetaData;
}

export interface APIError {
  code: string;
  message: string;
  details?: string;
}

export interface MetaData {
  page?: number;
  page_size?: number;
  total_count?: number;
  total_pages?: number;
}
```

---

## 📊 ФАЗА 2: Компоненты форм - COMPANIES (6-8 часов)

### 2.1. Создать полную форму компании

**Файл:** `features/companies/components/CompanyCreateForm.tsx` (переработать)

```typescript
'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Building2, ArrowLeft, Loader2, AlertCircle, CheckCircle, Save } from 'lucide-react';
import { CreateCompanyRequest, UpdateCompanyRequest, Company } from '@/types/entities';

interface CompanyFormProps {
  initialData?: Company;
  onSubmit: (data: CreateCompanyRequest | UpdateCompanyRequest) => Promise<void>;
  isEdit?: boolean;
}

export default function CompanyForm({ initialData, onSubmit, isEdit = false }: CompanyFormProps) {
  const router = useRouter();

  const [formData, setFormData] = useState<CreateCompanyRequest>({
    name: initialData?.name || '',
    description: initialData?.description || '',
  });

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Валидация
  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    // Название (обязательно, 3-255 символов)
    if (!formData.name.trim()) {
      newErrors.name = 'Название компании обязательно';
    } else if (formData.name.length < 3) {
      newErrors.name = 'Название должно быть не менее 3 символов';
    } else if (formData.name.length > 255) {
      newErrors.name = 'Название не должно превышать 255 символов';
    }

    // Описание (необязательно, но если есть - до 1000 символов)
    if (formData.description && formData.description.length > 1000) {
      newErrors.description = 'Описание не должно превышать 1000 символов';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    setIsLoading(true);
    setError(null);
    setSuccess(false);

    try {
      await onSubmit(formData);
      setSuccess(true);

      // Редирект через 1.5 секунды
      setTimeout(() => {
        router.push('/companies');
      }, 1500);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при сохранении компании');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6">
      <div className="max-w-3xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <button
            onClick={() => router.back()}
            className="flex items-center text-gray-400 hover:text-white mb-4 transition"
          >
            <ArrowLeft className="w-5 h-5 mr-2" />
            Назад
          </button>

          <div className="flex items-center">
            <Building2 className="w-8 h-8 text-blue-500 mr-3" />
            <h1 className="text-3xl font-bold text-white">
              {isEdit ? 'Редактировать компанию' : 'Создать компанию'}
            </h1>
          </div>
        </div>

        {/* Form Card */}
        <div className="bg-gray-800 rounded-lg shadow-xl p-8 border border-gray-700">
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Name */}
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Название компании <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className={`w-full px-4 py-3 rounded-lg bg-gray-900 border ${
                  errors.name ? 'border-red-500' : 'border-gray-700'
                } text-white placeholder-gray-500 focus:outline-none focus:border-blue-500 transition`}
                placeholder="ООО «Название компании»"
                maxLength={255}
              />
              {errors.name && (
                <p className="mt-2 text-sm text-red-500 flex items-center">
                  <AlertCircle className="w-4 h-4 mr-1" />
                  {errors.name}
                </p>
              )}
              <p className="mt-1 text-xs text-gray-500">{formData.name.length}/255 символов</p>
            </div>

            {/* Description */}
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Описание
              </label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className={`w-full px-4 py-3 rounded-lg bg-gray-900 border ${
                  errors.description ? 'border-red-500' : 'border-gray-700'
                } text-white placeholder-gray-500 focus:outline-none focus:border-blue-500 transition`}
                placeholder="Краткое описание компании..."
                rows={4}
                maxLength={1000}
              />
              {errors.description && (
                <p className="mt-2 text-sm text-red-500 flex items-center">
                  <AlertCircle className="w-4 h-4 mr-1" />
                  {errors.description}
                </p>
              )}
              <p className="mt-1 text-xs text-gray-500">
                {formData.description.length}/1000 символов
              </p>
            </div>

            {/* Error Message */}
            {error && (
              <div className="bg-red-500/10 border border-red-500 rounded-lg p-4">
                <div className="flex items-center text-red-500">
                  <AlertCircle className="w-5 h-5 mr-2" />
                  <span className="font-medium">{error}</span>
                </div>
              </div>
            )}

            {/* Success Message */}
            {success && (
              <div className="bg-green-500/10 border border-green-500 rounded-lg p-4">
                <div className="flex items-center text-green-500">
                  <CheckCircle className="w-5 h-5 mr-2" />
                  <span className="font-medium">
                    Компания успешно {isEdit ? 'обновлена' : 'создана'}!
                  </span>
                </div>
              </div>
            )}

            {/* Submit Button */}
            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700
                         text-white font-medium py-3 px-6 rounded-lg transition
                         flex items-center justify-center space-x-2"
            >
              {isLoading ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin" />
                  <span>Сохранение...</span>
                </>
              ) : (
                <>
                  <Save className="w-5 h-5" />
                  <span>{isEdit ? 'Обновить' : 'Создать'} компанию</span>
                </>
              )}
            </button>
          </form>
        </div>

        {/* Help Text */}
        <div className="mt-6 bg-blue-500/10 border border-blue-500 rounded-lg p-4">
          <p className="text-blue-400 text-sm">
            <strong>Примечание:</strong> После создания компании вы сможете добавить в неё сотрудников,
            банковские счета и другую информацию.
          </p>
        </div>
      </div>
    </div>
  );
}
```

**Файл:** `features/companies/components/CompanyUpdateForm.tsx` (обертка)

```typescript
'use client';

import React from 'react';
import CompanyForm from './CompanyForm';
import { UpdateCompanyRequest, Company } from '@/types/entities';
import { updateCompany } from '@/lib/api/companies';

interface CompanyUpdateFormProps {
  company: Company;
}

export default function CompanyUpdateForm({ company }: CompanyUpdateFormProps) {
  const handleSubmit = async (data: UpdateCompanyRequest) => {
    await updateCompany(company.id, data);
  };

  return <CompanyForm initialData={company} onSubmit={handleSubmit} isEdit />;
}
```

### 2.2. API интеграция для companies

**Файл:** `lib/api/companies.ts` (обновить)

```typescript
import apiClient from "./client";
import type {
  Company,
  CreateCompanyRequest,
  UpdateCompanyRequest,
  ListResponse,
  Employee,
  AddMemberRequest,
} from "@/types/entities";

const COMPANIES_BASE = "/companies";

export const companiesApi = {
  // Create
  async createCompany(data: CreateCompanyRequest): Promise<Company> {
    const response = await apiClient.post<{ data: Company }>(
      COMPANIES_BASE,
      data,
    );
    return response.data.data;
  },

  // Read
  async getCompany(id: string): Promise<Company> {
    const response = await apiClient.get<{ data: Company }>(
      `${COMPANIES_BASE}/${id}`,
    );
    return response.data.data;
  },

  async listCompanies(page = 1, pageSize = 20): Promise<ListResponse<Company>> {
    const response = await apiClient.get<ListResponse<Company>>(
      COMPANIES_BASE,
      {
        params: { page, page_size: pageSize },
      },
    );
    return response.data;
  },

  // Update
  async updateCompany(
    id: string,
    data: UpdateCompanyRequest,
  ): Promise<Company> {
    const response = await apiClient.put<{ data: Company }>(
      `${COMPANIES_BASE}/${id}`,
      data,
    );
    return response.data.data;
  },

  // Delete
  async deleteCompany(id: string): Promise<void> {
    await apiClient.delete(`${COMPANIES_BASE}/${id}`);
  },

  // Members management
  async getMembers(id: string): Promise<Employee[]> {
    const response = await apiClient.get<{ data: Employee[] }>(
      `${COMPANIES_BASE}/${id}/members`,
    );
    return response.data.data;
  },

  async addMember(
    id: string,
    data: Omit<AddMemberRequest, "organization_id">,
  ): Promise<Employee> {
    const response = await apiClient.post<{ data: Employee }>(
      `${COMPANIES_BASE}/${id}/members`,
      { ...data, organization_id: id },
    );
    return response.data.data;
  },

  async removeMember(id: string, memberId: string): Promise<void> {
    await apiClient.delete(`${COMPANIES_BASE}/${id}/members/${memberId}`);
  },
};
```

---

## 📊 ФАЗА 3: Компоненты форм - CATALOG (8-10 часов)

### 3.1. Создать форму каталога товаров

**Файл:** `features/catalog/components/CatalogItemForm.tsx` (новый)

```typescript
'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Package, ArrowLeft, Loader2, AlertCircle, CheckCircle, Save } from 'lucide-react';
import type { CatalogItem, CreateCatalogItemRequest, UpdateCatalogItemRequest } from '@/types/entities';

interface CatalogItemFormProps {
  initialData?: CatalogItem;
  onSubmit: (data: CreateCatalogItemRequest | UpdateCatalogItemRequest) => Promise<void>;
  isEdit?: boolean;
}

// Справочник единиц измерения (ОКЕИ)
const UNITS = [
  { code: '796', name: 'шт' },    // штука
  { code: '006', name: 'м' },     // метр
  { code: '055', name: 'м²' },    // квадратный метр
  { code: '113', name: 'м³' },    // кубический метр
  { code: '166', name: 'кг' },    // килограмм
  { code: '112', name: 'л' },     // литр
  { code: '212', name: 'т' },     // тонна
  { code: '006', name: 'компл' }, // комплект
  { code: '715', name: 'упак' },  // упаковка
];

// Популярные валюты
const CURRENCIES = [
  { code: 'KGS', name: 'Сом (KGS)' },
  { code: 'USD', name: 'Доллар США (USD)' },
  { code: 'RUB', name: 'Российский рубль (RUB)' },
  { code: 'EUR', name: 'Евро (EUR)' },
  { code: 'CNY', name: 'Юань (CNY)' },
];

export default function CatalogItemForm({ initialData, onSubmit, isEdit = false }: CatalogItemFormProps) {
  const router = useRouter();

  const [formData, setFormData] = useState<CreateCatalogItemRequest>({
    name: initialData?.name || '',
    number: initialData?.number || '',
    description: initialData?.description || '',
    tnved_code: initialData?.tnved_code || '',
    price: initialData?.price || 0,
    currency: initialData?.currency || 'KGS',
    unit: initialData?.unit || '796',
  });

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Валидация ТНВЭД кода (10 цифр)
  const isValidTnvedCode = (code: string): boolean => {
    return /^\d{10}$/.test(code);
  };

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    // Название
    if (!formData.name.trim()) {
      newErrors.name = 'Название товара обязательно';
    } else if (formData.name.length < 3) {
      newErrors.name = 'Название должно быть не менее 3 символов';
    } else if (formData.name.length > 255) {
      newErrors.name = 'Название не должно превышать 255 символов';
    }

    // Номер по каталогу
    if (!formData.number.trim()) {
      newErrors.number = 'Номер по каталогу обязателен';
    } else if (formData.number.length > 50) {
      newErrors.number = 'Номер не должен превышать 50 символов';
    }

    // ТНВЭД код (обязательно для ЭСФ)
    if (!formData.tnved_code.trim()) {
      newErrors.tnved_code = 'ТНВЭД код обязателен для электронных счетов-фактур';
    } else if (!isValidTnvedCode(formData.tnved_code)) {
      newErrors.tnved_code = 'ТНВЭД код должен состоять из 10 цифр';
    }

    // Цена
    if (formData.price <= 0) {
      newErrors.price = 'Цена должна быть больше нуля';
    } else if (formData.price > 999999999.99) {
      newErrors.price = 'Цена слишком большая';
    }

    // Единица измерения
    if (!formData.unit) {
      newErrors.unit = 'Выберите единицу измерения';
    }

    // Валюта
    if (!formData.currency) {
      newErrors.currency = 'Выберите валюту';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    setIsLoading(true);
    setError(null);
    setSuccess(false);

    try {
      await onSubmit(formData);
      setSuccess(true);

      setTimeout(() => {
        router.push('/catalog');
      }, 1500);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при сохранении товара');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <button
            onClick={() => router.back()}
            className="flex items-center text-gray-400 hover:text-white mb-4 transition"
          >
            <ArrowLeft className="w-5 h-5 mr-2" />
            Назад
          </button>

          <div className="flex items-center">
            <Package className="w-8 h-8 text-green-500 mr-3" />
            <h1 className="text-3xl font-bold text-white">
              {isEdit ? 'Редактировать товар' : 'Добавить товар в каталог'}
            </h1>
          </div>
        </div>

        {/* Form */}
        <div className="bg-gray-800 rounded-lg shadow-xl p-8 border border-gray-700">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* Name */}
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Название товара <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className={`w-full px-4 py-3 rounded-lg bg-gray-900 border ${
                    errors.name ? 'border-red-500' : 'border-gray-700'
                  } text-white focus:border-blue-500 outline-none transition`}
                  placeholder="Например: Ноутбук ASUS ROG"
                  maxLength={255}
                />
                {errors.name && (
                  <p className="mt-2 text-sm text-red-500 flex items-center">
                    <AlertCircle className="w-4 h-4 mr-1" />
                    {errors.name}
                  </p>
                )}
              </div>

              {/* Number */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Номер по каталогу <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  value={formData.number}
                  onChange={(e) => setFormData({ ...formData, number: e.target.value })}
                  className={`w-full px-4 py-3 rounded-lg bg-gray-900 border ${
                    errors.number ? 'border-red-500' : 'border-gray-700'
                  } text-white focus:border-blue-500 outline-none transition`}
                  placeholder="CAT-001"
                  maxLength={50}
                />
                {errors.number && (
                  <p className="mt-2 text-sm text-red-500">{errors.number}</p>
                )}
              </div>

              {/* TNVED Code */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  ТНВЭД код <span className="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  value={formData.tnved_code}
                  onChange={(e) => {
                    const value = e.target.value.replace(/\D/g, '').slice(0, 10);
                    setFormData({ ...formData, tnved_code: value });
                  }}
                  className={`w-full px-4 py-3 rounded-lg bg-gray-900 border ${
                    errors.tnved_code ? 'border-red-500' : 'border-gray-700'
                  } text-white focus:border-blue-500 outline-none transition`}
                  placeholder="8471300000"
                  maxLength={10}
                />
                {errors.tnved_code && (
                  <p className="mt-2 text-sm text-red-500">{errors.tnved_code}</p>
                )}
                <p className="mt-1 text-xs text-gray-500">10 цифр (обязательно для ЭСФ)</p>
              </div>

              {/* Price */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Цена <span className="text-red-500">*</span>
                </label>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  value={formData.price}
                  onChange={(e) => setFormData({ ...formData, price: parseFloat(e.target.value) || 0 })}
                  className={`w-full px-4 py-3 rounded-lg bg-gray-900 border ${
                    errors.price ? 'border-red-500' : 'border-gray-700'
                  } text-white focus:border-blue-500 outline-none transition`}
                  placeholder="0.00"
                />
                {errors.price && (
                  <p className="mt-2 text-sm text-red-500">{errors.price}</p>
                )}
              </div>

              {/* Currency */}
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Валюта <span className="text-red-500">*</span>
                </label>
                <select
                  value={formData.currency}
                  onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
                  className="w-full px-4 py-3 rounded-lg bg-gray-900 border border-gray-700
                             text-white focus:border-blue-500 outline-none transition"
                >
                  {CURRENCIES.map((curr) => (
                    <option key={curr.code} value={curr.code}>
                      {curr.name}
                    </option>
                  ))}
                </select>
              </div>

              {/* Unit */}
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Единица измерения <span className="text-red-500">*</span>
                </label>
                <select
                  value={formData.unit}
                  onChange={(e) => setFormData({ ...formData, unit: e.target.value })}
                  className="w-full px-4 py-3 rounded-lg bg-gray-900 border border-gray-700
                             text-white focus:border-blue-500 outline-none transition"
                >
                  {UNITS.map((unit) => (
                    <option key={unit.code} value={unit.code}>
                      {unit.name} (ОКЕИ: {unit.code})
                    </option>
                  ))}
                </select>
              </div>

              {/* Description */}
              <div className="md:col-span-2">
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Описание
                </label>
                <textarea
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="w-full px-4 py-3 rounded-lg bg-gray-900 border border-gray-700
                             text-white focus:border-blue-500 outline-none transition"
                  placeholder="Дополнительное описание товара..."
                  rows={4}
                  maxLength={1000}
                />
                <p className="mt-1 text-xs text-gray-500">
                  {formData.description?.length || 0}/1000 символов
                </p>
              </div>
            </div>

            {/* Error/Success Messages */}
            {error && (
              <div className="bg-red-500/10 border border-red-500 rounded-lg p-4">
                <div className="flex items-center text-red-500">
                  <AlertCircle className="w-5 h-5 mr-2" />
                  <span>{error}</span>
                </div>
              </div>
            )}

            {success && (
              <div className="bg-green-500/10 border border-green-500 rounded-lg p-4">
                <div className="flex items-center text-green-500">
                  <CheckCircle className="w-5 h-5 mr-2" />
                  <span>Товар успешно {isEdit ? 'обновлен' : 'добавлен'}!</span>
                </div>
              </div>
            )}

            {/* Submit Button */}
            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-green-600 hover:bg-green-700 disabled:bg-gray-700
                         text-white font-medium py-3 px-6 rounded-lg transition
                         flex items-center justify-center space-x-2"
            >
              {isLoading ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin" />
                  <span>Сохранение...</span>
                </>
              ) : (
                <>
                  <Save className="w-5 h-5" />
                  <span>{isEdit ? 'Обновить' : 'Добавить'} товар</span>
                </>
              )}
            </button>
          </form>
        </div>

        {/* Help Text */}
        <div className="mt-6 bg-blue-500/10 border border-blue-500 rounded-lg p-4">
          <p className="text-blue-400 text-sm">
            <strong>Важно:</strong> ТНВЭД код обязателен для создания электронных счетов-фактур (ЭСФ).
            Вы можете найти код вашего товара в{' '}
            <a
              href="http://www.eaeunion.org/tnved"
              target="_blank"
              rel="noopener noreferrer"
              className="underline hover:text-blue-300"
            >
              классификаторе ТН ВЭД ЕАЭС
            </a>
            .
          </p>
        </div>
      </div>
    </div>
  );
}
```

**Файл:** `lib/api/catalog.ts` (новый)

```typescript
import apiClient from "./client";
import type {
  CatalogItem,
  CreateCatalogItemRequest,
  UpdateCatalogItemRequest,
  ListResponse,
} from "@/types/entities";

const CATALOG_BASE = "/catalog";

export const catalogApi = {
  async createItem(data: CreateCatalogItemRequest): Promise<CatalogItem> {
    const response = await apiClient.post<{ data: CatalogItem }>(
      CATALOG_BASE,
      data,
    );
    return response.data.data;
  },

  async getItem(id: string): Promise<CatalogItem> {
    const response = await apiClient.get<{ data: CatalogItem }>(
      `${CATALOG_BASE}/${id}`,
    );
    return response.data.data;
  },

  async listItems(page = 1, pageSize = 50): Promise<ListResponse<CatalogItem>> {
    const response = await apiClient.get<ListResponse<CatalogItem>>(
      CATALOG_BASE,
      {
        params: { page, page_size: pageSize },
      },
    );
    return response.data;
  },

  async updateItem(
    id: string,
    data: UpdateCatalogItemRequest,
  ): Promise<CatalogItem> {
    const response = await apiClient.put<{ data: CatalogItem }>(
      `${CATALOG_BASE}/${id}`,
      data,
    );
    return response.data.data;
  },

  async deleteItem(id: string): Promise<void> {
    await apiClient.delete(`${CATALOG_BASE}/${id}`);
  },
};
```

---

## 📊 ФАЗА 4: Остальные формы (20-25 часов)

### 4.1. Bank Accounts Form

- `features/bank-accounts/components/BankAccountForm.tsx`
- Валидация номера счета (20 цифр для КР)
- Валидация БИК банка

### 4.2. Documents Form

- `features/documents/components/DocumentForm.tsx`
- Поддержка entries (ключ-значение)
- Workflow статусов (draft → pending → approved)

### 4.3. Foreign Companies Form

- `features/foreign-companies/components/ForeignCompanyForm.tsx`
- Валидация country code (ISO 3166-1 alpha-2)
- Dropdown с популярными странами

### 4.4. Invoices Form (ESF) - самая сложная

- `features/invoices/components/InvoiceESFForm.tsx`
- Многошаговая форма (wizard)
- Интеграция с catalog для выбора товаров
- Расчет НДС и итоговых сумм
- Валидация всех обязательных полей ESF

---

## ✅ Checklist реализации

### Фаза 1: Типизация

- [ ] `types/entities.ts` - все типы сущностей
- [ ] Обновить существующие типы в `types/`
- [ ] Синхронизация с backend models

### Фаза 2: Companies

- [ ] `CompanyForm.tsx` - универсальная форма
- [ ] `CompanyCreateForm.tsx` - обертка для создания
- [ ] `CompanyUpdateForm.tsx` - обертка для редактирования
- [ ] `lib/api/companies.ts` - полное API
- [ ] Валидация на клиенте
- [ ] Unit тесты форм
- [ ] E2E тесты

### Фаза 3: Catalog

- [ ] `CatalogItemForm.tsx` - форма товара
- [ ] ТНВЭД код валидация
- [ ] Справочники (единицы, валюты)
- [ ] `lib/api/catalog.ts`
- [ ] Тесты

### Фаза 4: Bank Accounts

- [ ] `BankAccountForm.tsx`
- [ ] Валидация счетов
- [ ] `lib/api/bank-accounts.ts`
- [ ] Тесты

### Фаза 5: Documents

- [ ] `DocumentForm.tsx`
- [ ] Dynamic entries
- [ ] `lib/api/documents.ts`
- [ ] Тесты

### Фаза 6: Foreign Companies

- [ ] `ForeignCompanyForm.tsx`
- [ ] Country selector
- [ ] `lib/api/foreign-companies.ts`
- [ ] Тесты

### Фаза 7: Invoices (ESF)

- [ ] `InvoiceESFForm.tsx` - многошаговая
- [ ] Шаг 1: Основная информация
- [ ] Шаг 2: Стороны сделки
- [ ] Шаг 3: Товары (из каталога)
- [ ] Шаг 4: Расчеты (НДС, итого)
- [ ] Шаг 5: Проверка и отправка
- [ ] `lib/api/invoices.ts` - обновить
- [ ] Тесты

---

## 📅 Временная оценка

| Фаза      | Компонент              | Время           |
| --------- | ---------------------- | --------------- |
| 1         | Типизация              | 4-6 часов       |
| 2         | Companies Form         | 6-8 часов       |
| 3         | Catalog Form           | 8-10 часов      |
| 4         | Bank Accounts Form     | 4-6 часов       |
| 5         | Documents Form         | 6-8 часов       |
| 6         | Foreign Companies Form | 4-6 часов       |
| 7         | Invoices ESF Form      | 20-25 часов     |
| 8         | Testing & Integration  | 10-12 часов     |
| **ИТОГО** |                        | **62-81 часов** |

---

## 🔄 Порядок реализации

1. **Неделя 1-2:** Типизация + Companies + Catalog (18-24 часа)
2. **Неделя 3:** Bank Accounts + Documents (10-14 часов)
3. **Неделя 4:** Foreign Companies + начало Invoices (24-31 час)
4. **Неделя 5:** Завершение Invoices + Testing (10-12 часов)

---

## 🛠️ Технический стек

- **Validation:** Zod schemas + custom validators
- **Forms:** React Hook Form (опционально)
- **UI:** Tailwind CSS + Lucide icons
- **State:** React Query для API calls
- **Testing:** Vitest + React Testing Library + Playwright

---

**Готово к началу реализации!**
