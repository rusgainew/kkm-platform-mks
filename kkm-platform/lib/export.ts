import { jsPDF } from 'jspdf';
import Papa from 'papaparse';

export interface ExportOptions {
  filename?: string;
  title?: string;
  columns?: string[];
}

/**
 * Экспортировать данные в PDF
 */
export function exportToPDF<T extends Record<string, any>>(
  data: T[],
  options: ExportOptions & { format?: 'a4' | 'letter' } = {}
) {
  const {
    filename = 'export.pdf',
    title = 'Отчет',
    columns = Object.keys(data[0] || {}),
    format = 'a4',
  } = options;

  const doc = new jsPDF({ format });

  // Добавить заголовок
  if (title) {
    doc.setFontSize(16);
    doc.text(title, 15, 15);
    doc.setFontSize(11);
    doc.text(`Дата: ${new Date().toLocaleString('ru-RU')}`, 15, 25);
  }

  // Подготовить данные для таблицы
  const tableData = data.map((row) =>
    columns.map((col) => {
      const value = row[col];
      if (value === null || value === undefined) return '';
      if (typeof value === 'object') return JSON.stringify(value);
      if (typeof value === 'boolean') return value ? 'Да' : 'Нет';
      return String(value);
    })
  );

  // @ts-ignore
  doc.autoTable({
    head: [columns.map((col) => {
      // Переводить названия на русский
      const translations: Record<string, string> = {
        id: 'ID',
        name: 'Название',
        email: 'Email',
        status: 'Статус',
        role: 'Роль',
        createdAt: 'Создано',
        updatedAt: 'Обновлено',
        amount: 'Сумма',
        totalAmount: 'Итого',
        invoiceNumber: 'Номер счета',
        dueDate: 'Срок',
        issueDate: 'Дата выставления',
        price: 'Цена',
        cost: 'Себестоимость',
        stock: 'Остаток',
        sku: 'Артикул',
      };
      return translations[col] || col;
    })],
    body: tableData,
    startY: title ? 35 : 10,
    margin: 10,
  });

  doc.save(filename);
}

/**
 * Экспортировать данные в CSV
 */
export function exportToCSV<T extends Record<string, any>>(
  data: T[],
  options: ExportOptions = {}
) {
  const { filename = 'export.csv', columns = Object.keys(data[0] || {}) } = options;

  // Подготовить данные
  const csvData = [
    columns, // Заголовки
    ...data.map((row) =>
      columns.map((col) => {
        const value = row[col];
        if (value === null || value === undefined) return '';
        if (typeof value === 'object') return JSON.stringify(value);
        // Экранировать значения если нужно
        const stringValue = String(value);
        return stringValue.includes(',') ? `"${stringValue}"` : stringValue;
      })
    ),
  ];

  const csv = Papa.unparse(csvData);
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
  const link = document.createElement('a');
  link.href = URL.createObjectURL(blob);
  link.download = filename;
  link.click();
  URL.revokeObjectURL(link.href);
}

/**
 * Экспортировать JSON
 */
export function exportToJSON<T extends Record<string, any>>(
  data: T[],
  options: ExportOptions = {}
) {
  const { filename = 'export.json' } = options;

  const json = JSON.stringify(data, null, 2);
  const blob = new Blob([json], { type: 'application/json;charset=utf-8;' });
  const link = document.createElement('a');
  link.href = URL.createObjectURL(blob);
  link.download = filename;
  link.click();
  URL.revokeObjectURL(link.href);
}
