/**
 * Справочные типы для ESF API
 * Используются для представления кодирования справочников
 * (статусы, валюты, налоги, классификации и т.д.)
 */

/**
 * Основной тип справочника - код + название
 * Используется для всех справочников в ESF API:
 * - Валюты (RUB, USD, KGS и т.д.)
 * - Статусы документов (10, 20, 30 и т.д.)
 * - Типы доставки, платежей, налогов и т.д.
 */
export interface ReferenceItem {
  code: string; // Уникальный код справочника
  name: string; // Человеко-читаемое название
  id?: number; // Опционально для некоторых справочников
  shortName?: string; // Краткое название если есть
  description?: string; // Описание
}

/**
 * Информация о налоге на добавленную стоимость (НДС)
 */
export interface VatTaxType {
  code?: string; // Код в справочнике
  rate: string; // Процент: "12%", "0%", "5%" и т.д.
  name: string; // Описание: "Стандартная ставка НДС 12%"
}

/**
 * Информация о других налогах и сборах
 * (Например, НСП - налог на специальные платежи)
 */
export interface TaxInfo {
  code: string; // Код налога в справочнике
  name: string; // Название налога
  rate: string; // Ставка налога или процент
  type?: string; // Тип налога (например "excise", "sales", "special")
}

/**
 * Справочник валют - расширение ReferenceItem
 * Используется для указания валюты платежей и документов
 */
export interface CurrencyReference extends ReferenceItem {
  code: string; // ISO 4217 код (RUB, USD, KGS и т.д.)
  name: string; // Название валюты
  symbol?: string; // Символ валюты (₽, $, с и т.д.)
  rate?: number; // Курс обмена относительно базовой валюты
}

/**
 * Справочник стран
 * Используется для указания страны происхождения товаров
 */
export interface CountryReference extends ReferenceItem {
  code: string; // ISO 3166-1 код (KG, RU и т.д.)
  name: string; // Название страны
  alpha3?: string; // Трёхбуквенный код (KGZ, RUS и т.д.)
}

/**
 * Справочник типов документов
 * СФ-РЕ (Реализация), СФ-ПР (Покупка), КСФ (Коррекция)
 */
export interface DocumentTypeReference extends ReferenceItem {
  code: string; // Код типа (1 = РЕ, 2 = ПР, 3 = КСФ)
  name: string; // Название документа
}

/**
 * Справочник типов доставки
 */
export interface DeliveryTypeReference extends ReferenceItem {
  code: string; // Код типа доставки
  name: string; // Название доставки
}

/**
 * Справочник типов платежей
 */
export interface PaymentTypeReference extends ReferenceItem {
  code: string; // Код платежа
  name: string; // Название платежа
}

/**
 * Справочник классификации единиц измерения
 * Например: "796" = Штука, "110" = Килограмм и т.д.
 */
export interface UnitClassificationReference extends ReferenceItem {
  code: string; // Код единицы измерения
  name: string; // Название (Штука, КГ, Метр и т.д.)
}

/**
 * Справочник кодов классификации товаров
 * ТНВЭД - Товарная номенклатура внешнеэкономической деятельности
 */
export interface TNVEDCodeReference extends ReferenceItem {
  code: string; // ТНВЭД код (8-10 цифр)
}

/**
 * Справочник кодов ГКЭД
 * ГКЭД - Государственная классификация экономической деятельности
 */
export interface GKEDCodeReference extends ReferenceItem {
  code: string; // ГКЭД код
}

/**
 * Справочник статусов документов в ESF API
 * Используются числовые коды вместо строк
 */
export const ESF_DOCUMENT_STATUS_CODES = {
  DRAFT: "10", // Новый
  REVOKED: "20", // Отозван
  SENT: "30", // Отправлен
  ACCEPTED: "40", // Принят
  REJECTED: "50", // Отклонен
  DELETED: "70", // Удален
  ALL: "90", // Все
} as const;

/**
 * Типы способов доставки
 */
export const DELIVERY_TYPE_CODES = {
  DIRECT: "1", // Прямая доставка
  WAREHOUSE: "2", // Через склад
  SERVICE: "3", // Услуга доставки
} as const;

/**
 * Типы платежей
 */
export const PAYMENT_TYPE_CODES = {
  CASH: "1", // Наличные
  BANK_TRANSFER: "2", // Банковский перевод
  CHECK: "3", // Чек
} as const;

/**
 * Распространённые коды единиц измерения
 */
export const UNIT_CLASSIFICATION_CODES = {
  PIECE: "796", // Штука
  KILOGRAM: "110", // Килограмм
  GRAM: "163", // Грамм
  METER: "004", // Метр
  SQUARE_METER: "055", // Квадратный метр
  LITER: "112", // Литр
  HOUR: "255", // Час
  DAY: "259", // День
  KILOMETER: "037", // Километр
} as const;

/**
 * Коды валют (ISO 4217)
 */
export const ESF_CURRENCY_CODES = {
  KGS: "KGS", // Киргизский сом (основная валюта)
  RUB: "RUB", // Российский рубль
  USD: "USD", // Американский доллар
  EUR: "EUR", // Евро
  GBP: "GBP", // Британский фунт
  JPY: "JPY", // Японская йена
  CNY: "CNY", // Китайский юань
  KZT: "KZT", // Казахстанский тенге
  INR: "INR", // Индийская рупия
  TRY: "TRY", // Турецкая лира
} as const;

/**
 * Коды стран (ISO 3166-1 alpha-2)
 */
export const ESF_COUNTRY_CODES = {
  KG: "KG", // Киргизия (основная)
  RU: "RU", // Россия
  US: "US", // США
  DE: "DE", // Германия
  FR: "FR", // Франция
  GB: "GB", // Великобритания
  CN: "CN", // Китай
  JP: "JP", // Япония
  KZ: "KZ", // Казахстан
  TR: "TR", // Турция
  IN: "IN", // Индия
} as const;

/**
 * Возвращает ReferenceItem по коду из справочника
 * @param code Код справочника
 * @param dictionary Словарь кодов (например CURRENCY_CODES)
 * @param labels Словарь названий (например для отображения)
 * @returns ReferenceItem с заполненными полями
 */
export function createReferenceItem(
  code: string,
  name: string,
  id?: number,
): ReferenceItem {
  return {
    code,
    name,
    id,
  };
}

/**
 * Возвращает ReferenceItem с налоговой информацией
 */
export function createTaxInfo(
  code: string,
  name: string,
  rate: string,
  type?: string,
): TaxInfo {
  return {
    code,
    name,
    rate,
    type,
  };
}

/**
 * Возвращает VatTaxType для стандартных ставок НДС
 */
export function createVatTaxType(
  rate: string,
  description?: string,
): VatTaxType {
  const name = description || `Ставка НДС ${rate}`;
  return {
    rate,
    name,
  };
}

export type ReferenceItemType =
  | ReferenceItem
  | CurrencyReference
  | CountryReference
  | DocumentTypeReference
  | DeliveryTypeReference
  | PaymentTypeReference
  | UnitClassificationReference
  | TNVEDCodeReference
  | GKEDCodeReference;
