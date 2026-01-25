'use client';

import React, { useState, useRef } from 'react';
import { Upload, Download, AlertCircle, CheckCircle, Loader2 } from 'lucide-react';
import { createUser } from '@/lib/api/users';
import type { RegisterRequest } from '@/lib/api/users';

interface CSVRow {
  email: string;
  first_name: string;
  last_name: string;
  password: string;
  role?: string;
}

interface ImportResult {
  success: number;
  failed: number;
  errors: Array<{ row: number; message: string }>;
}

export default function BulkUserImport() {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [importResult, setImportResult] = useState<ImportResult | null>(null);
  const [errorMessage, setErrorMessage] = useState('');
  const [showTemplate, setShowTemplate] = useState(false);

  // Download template CSV
  const handleDownloadTemplate = () => {
    const csvContent = `email,first_name,last_name,password,role
john.doe@example.com,John,Doe,SecurePass123!,cashier
jane.smith@example.com,Jane,Smith,SecurePass456!,manager
bob.johnson@example.com,Bob,Johnson,SecurePass789!,cashier`;

    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = 'users_template.csv';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  };

  // Parse CSV file
  const parseCSV = (content: string): CSVRow[] => {
    const lines = content.trim().split('\n');
    if (lines.length < 2) {
      throw new Error('CSV должен содержать заголовок и минимум одну строку');
    }

    const headers = lines[0].split(',').map(h => h.trim().toLowerCase());
    const rows: CSVRow[] = [];

    for (let i = 1; i < lines.length; i++) {
      const values = lines[i].split(',').map(v => v.trim());
      if (values.some(v => !v)) continue; // Skip empty rows

      const row: CSVRow = {
        email: '',
        first_name: '',
        last_name: '',
        password: '',
      };

      headers.forEach((header, idx) => {
        if (header === 'email') row.email = values[idx];
        if (header === 'first_name') row.first_name = values[idx];
        if (header === 'last_name') row.last_name = values[idx];
        if (header === 'password') row.password = values[idx];
        if (header === 'role') row.role = values[idx];
      });

      rows.push(row);
    }

    return rows;
  };

  // Validate CSV row
  const validateRow = (row: CSVRow, rowNumber: number): string | null => {
    if (!row.email) return `Строка ${rowNumber}: Email отсутствует`;
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(row.email)) 
      return `Строка ${rowNumber}: Некорректный email`;
    if (!row.first_name) return `Строка ${rowNumber}: Имя отсутствует`;
    if (!row.last_name) return `Строка ${rowNumber}: Фамилия отсутствует`;
    if (!row.password) return `Строка ${rowNumber}: Пароль отсутствует`;
    if (row.password.length < 8) 
      return `Строка ${rowNumber}: Пароль короче 8 символов`;
    if (!/[A-Z]/.test(row.password)) 
      return `Строка ${rowNumber}: Пароль без заглавной буквы`;
    if (!/[a-z]/.test(row.password)) 
      return `Строка ${rowNumber}: Пароль без строчной буквы`;
    if (!/[0-9]/.test(row.password)) 
      return `Строка ${rowNumber}: Пароль без цифры`;
    if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(row.password)) 
      return `Строка ${rowNumber}: Пароль без спецсимвола`;
    if (row.role && !['admin', 'manager', 'cashier'].includes(row.role)) 
      return `Строка ${rowNumber}: Неизвестная роль "${row.role}"`;
    return null;
  };

  // Handle file import
  const handleFileImport = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setIsLoading(true);
    setErrorMessage('');
    setImportResult(null);

    try {
      const content = await file.text();
      const rows = parseCSV(content);

      // Validate all rows first
      const errors: Array<{ row: number; message: string }> = [];
      const validRows: CSVRow[] = [];

      rows.forEach((row, idx) => {
        const error = validateRow(row, idx + 2);
        if (error) {
          errors.push({ row: idx + 2, message: error });
        } else {
          validRows.push(row);
        }
      });

      if (validRows.length === 0) {
        throw new Error('Нет валидных строк для импорта');
      }

      // Import valid rows
      let successCount = 0;
      const importErrors = [...errors];

      for (const row of validRows) {
        try {
          const payload: RegisterRequest = {
            email: row.email,
            password: row.password,
            first_name: row.first_name,
            last_name: row.last_name,
          };
          await createUser(payload);
          successCount++;
        } catch (error) {
          const errorMsg = error instanceof Error ? error.message : 'Неизвестная ошибка';
          importErrors.push({
            row: rows.indexOf(row) + 2,
            message: errorMsg,
          });
        }
      }

      setImportResult({
        success: successCount,
        failed: importErrors.length,
        errors: importErrors.slice(0, 10), // Show first 10 errors
      });
    } catch (error) {
      const msg = error instanceof Error ? error.message : 'Ошибка при импорте файла';
      setErrorMessage(msg);
    } finally {
      setIsLoading(false);
      // Reset file input
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  return (
    <div className="space-y-6">
      {/* Info Card */}
      <div className="bg-blue-900/20 border border-blue-800 rounded-lg p-4">
        <p className="text-blue-300 text-sm">
          📋 Вы можете создать нескольких пользователей одновременно, загрузив CSV файл.
          Скачайте шаблон ниже для начала.
        </p>
      </div>

      {/* Buttons */}
      <div className="flex gap-3">
        <button
          onClick={handleDownloadTemplate}
          className="flex items-center gap-2 px-4 py-2 bg-gray-800 hover:bg-gray-700 text-white rounded-lg font-semibold transition-colors"
        >
          <Download className="w-5 h-5" />
          Скачать шаблон CSV
        </button>
        <button
          onClick={() => fileInputRef.current?.click()}
          disabled={isLoading}
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isLoading ? (
            <>
              <Loader2 className="w-5 h-5 animate-spin" />
              Загрузка...
            </>
          ) : (
            <>
              <Upload className="w-5 h-5" />
              Загрузить CSV
            </>
          )}
        </button>
        <input
          ref={fileInputRef}
          type="file"
          accept=".csv"
          onChange={handleFileImport}
          disabled={isLoading}
          className="hidden"
        />
      </div>

      {/* Error Message */}
      {errorMessage && (
        <div className="p-4 bg-red-900/20 border border-red-800 rounded-lg flex items-start gap-3">
          <AlertCircle className="w-5 h-5 text-red-400 shrink-0 mt-0.5" />
          <div>
            <p className="text-red-300 font-semibold text-sm">Ошибка</p>
            <p className="text-red-400 text-sm mt-1">{errorMessage}</p>
          </div>
        </div>
      )}

      {/* Import Result */}
      {importResult && (
        <div className="p-4 border rounded-lg">
          {importResult.failed === 0 ? (
            <div className="bg-emerald-900/20 border border-emerald-800 rounded-lg p-4 flex items-start gap-3">
              <CheckCircle className="w-5 h-5 text-emerald-400 shrink-0 mt-0.5" />
              <div>
                <p className="text-emerald-300 font-semibold text-sm">Импорт успешен</p>
                <p className="text-emerald-400 text-sm mt-1">
                  ✅ Успешно создано {importResult.success} пользователей
                </p>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="bg-yellow-900/20 border border-yellow-800 rounded-lg p-4 flex items-start gap-3">
                <AlertCircle className="w-5 h-5 text-yellow-400 shrink-0 mt-0.5" />
                <div>
                  <p className="text-yellow-300 font-semibold text-sm">Импорт завершён с ошибками</p>
                  <p className="text-yellow-400 text-sm mt-1">
                    ✅ Успешно: {importResult.success} | ❌ Ошибок: {importResult.failed}
                  </p>
                </div>
              </div>

              {importResult.errors.length > 0 && (
                <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-4">
                  <p className="text-gray-300 font-semibold text-sm mb-3">Ошибки:</p>
                  <div className="space-y-2 max-h-48 overflow-y-auto">
                    {importResult.errors.map((error, idx) => (
                      <div key={idx} className="text-red-400 text-xs">
                        <span className="font-semibold">Строка {error.row}:</span> {error.message}
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {/* CSV Format Info */}
      <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-4">
        <button
          onClick={() => setShowTemplate(!showTemplate)}
          className="text-gray-300 hover:text-white font-semibold text-sm transition-colors"
        >
          {showTemplate ? '▼' : '▶'} Формат CSV файла
        </button>
        {showTemplate && (
          <div className="mt-3 space-y-2 text-gray-400 text-xs">
            <p>Файл должен содержать следующие колонки:</p>
            <div className="bg-gray-900 p-3 rounded font-mono overflow-x-auto">
              <p>email,first_name,last_name,password,role</p>
              <p className="text-gray-500 mt-2"># Примеры:</p>
              <p>john@example.com,John,Doe,SecurePass123!,cashier</p>
              <p>jane@example.com,Jane,Smith,SecurePass456!,manager</p>
            </div>
            <p className="mt-2">
              <strong>Требования к паролю:</strong> минимум 8 символов, заглавная буква, строчная буква, цифра, спецсимвол
            </p>
            <p>
              <strong>Роли:</strong> cashier (кассир), manager (менеджер), admin (администратор)
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
