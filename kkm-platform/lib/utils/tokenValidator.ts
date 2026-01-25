/**
 * Валидатор JWT токенов
 * Проверяет срок действия токенов и декодирует их без обращения к серверу
 */

/**
 * Декодирует JWT токен и возвращает payload (без проверки подписи)
 * @param token - JWT токен в формате header.payload.signature
 * @returns Декодированный объект payload или null если ошибка
 */
export function decodeJwt(token: string): { exp?: number; [key: string]: unknown } | null {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) {
      console.warn('[decodeJwt] Неверный формат токена - ожидается 3 части, получено', parts.length);
      return null;
    }

    // Декодируем payload (вторая часть JWT)
    const payload = parts[1];
    // Добавляем padding если нужен для base64 декодирования
    const padded = payload + '='.repeat((4 - (payload.length % 4)) % 4);
    const decoded = JSON.parse(Buffer.from(padded, 'base64').toString());
    return decoded;
  } catch (error) {
    console.error('[decodeJwt] Ошибка декодирования токена:', error);
    return null;
  }
}

/**
 * Проверяет истек ли JWT токен
 * @param token - JWT токен для проверки
 * @param bufferSeconds - буфер в секундах перед истечением (по умолчанию 30)
 * @returns true если токен истек или истекает в ближайшие bufferSeconds секунд
 */
export function isTokenExpired(token: string, bufferSeconds: number = 30): boolean {
  const decoded = decodeJwt(token);
  if (!decoded?.exp) {
    console.warn('[isTokenExpired] Нет claim exp в токене');
    return true;
  }

  const now = Math.floor(Date.now() / 1000);
  const expiresAt = decoded.exp;
  const isExpired = now >= (expiresAt - bufferSeconds);

  if (isExpired) {
    console.warn('[isTokenExpired] Токен истек или истекает очень скоро. Истекает:', new Date(expiresAt * 1000).toISOString(), 'Сейчас:', new Date(now * 1000).toISOString());
  } else {
    const secsUntilExpiry = expiresAt - now;
    console.log('[isTokenExpired] Токен действителен еще', secsUntilExpiry, 'секунд');
  }

  return isExpired;
}

/**
 * Получить дату истечения токена
 * @param token - JWT токен
 * @returns Дата истечения или null если не найдена
 */
export function getTokenExpiryDate(token: string): Date | null {
  const decoded = decodeJwt(token);
  if (!decoded?.exp) return null;
  return new Date(decoded.exp * 1000);
}
