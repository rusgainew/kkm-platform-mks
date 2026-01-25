/**
 * Утилиты для форматирования дат
 * Работает с timestamp'ами в секундах и миллисекундах
 */

/**
 * Форматирует timestamp в дату по российскому стандарту
 * @param timestamp - Unix timestamp в секундах или миллисекундах
 * @param includeTime - Включить ли время в формат (по умолчанию false)
 * @returns Отформатированная дата строка
 */
export function formatDate(
  timestamp: number,
  includeTime: boolean = false
): string {
  // Определяем единицы: если меньше 10 млрд, то в секундах, иначе в миллисекундах
  const ms = timestamp < 10000000000 ? timestamp * 1000 : timestamp;
  const date = new Date(ms);

  if (includeTime) {
    return date.toLocaleString("ru-RU", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  }

  return date.toLocaleDateString("ru-RU", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}

/**
 * Форматирует timestamp в полный формат с временем
 * @param timestamp - Unix timestamp в секундах или миллисекундах, или строка ISO
 * @returns Отформатированная дата с временем
 */
export function formatDateTime(timestamp: number | string): string {
  if (typeof timestamp === "string") {
    // Если это строка ISO, конвертируем в миллисекунды
    const date = new Date(timestamp);
    return date.toLocaleString("ru-RU", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  }
  return formatDate(timestamp, true);
}

/**
 * Получает дату в коротком формате (пример: 24 дек 2024)
 * @param timestamp - Unix timestamp в секундах или миллисекундах
 * @returns Коротко отформатированная дата
 */
export function formatDateShort(timestamp: number): string {
  const ms = timestamp < 10000000000 ? timestamp * 1000 : timestamp;
  const date = new Date(ms);

  return date.toLocaleDateString("ru-RU", {
    year: "2-digit",
    month: "short",
    day: "numeric",
  });
}

/**
 * Получает дату объекта для функций сортировки
 * @param timestamp - Unix timestamp в секундах или миллисекундах
 * @returns Объект Date
 */
export function parseTimestamp(timestamp: number): Date {
  const ms = timestamp < 10000000000 ? timestamp * 1000 : timestamp;
  return new Date(ms);
}
