/**
 * Типы для представления участников сделки в ESF API
 * Party (Сторона) и LegalPerson (Юридическое лицо)
 */

import type { ReferenceItem } from "./reference";

/**
 * Основной тип для представления стороны сделки в ESF API
 * Используется как для поставщика, так и для покупателя
 *
 * Идентифицирует юридическое лицо (резидента КР) по его ИНН (PIN)
 */
export interface Party {
  pin: string; // Идентификационный номер налогоплательщика (ИНН)
  fullName: string; // Полное наименование организации
  mainFullName?: string; // Наименование головной организации (если филиал)
  mainPin?: string; // ИНН головной организации (если филиал)
  address?: string; // Адрес
  isResident?: boolean; // Резидент ли КР (true по умолчанию для Party)
  bankAccount?: string; // Банковский счет для платежей
  country?: ReferenceItem; // Справочник страны (обычно KG для Party)
}

/**
 * Тип LegalPerson - синоним для Party
 * Используется в ESF API для обозначения юридического лица
 * Может использоваться и в контексте резидента и нерезидента
 */
export type LegalPerson = Party;

/**
 * Расширенная информация об иностранной компании
 * Используется когда одна из сторон - иностранное лицо (не резидент КР)
 *
 * Примечание: В ESF API используется для контрактора, который не является
 * резидентом Киргизии
 */
export interface ForeignParty {
  pin?: string; // Может быть пустым для иностранцев
  fullName: string; // Полное наименование иностранной компании
  address?: string; // Адрес за границей
  country: ReferenceItem; // Страна иностранной компании
  isResident: false; // Флаг что это не резидент КР
}

/**
 * Объединённый тип для любой стороны сделки
 * Может быть как местной компанией (Party), так и иностранной (ForeignParty)
 */
export type ContractParty = Party | ForeignParty;

/**
 * Вспомогательная функция для создания Party
 * @param pin - ИНН
 * @param fullName - Полное название
 * @param address - Адрес (опционально)
 * @param mainPin - ИНН главной организации (если филиал)
 * @returns Party объект
 */
export function createParty(
  pin: string,
  fullName: string,
  address?: string,
  mainPin?: string,
): Party {
  return {
    pin,
    fullName,
    address,
    mainPin,
    isResident: true,
  };
}

/**
 * Вспомогательная функция для создания ForeignParty
 * @param fullName - Полное название иностранной компании
 * @param countryCode - Код страны (ISO 3166-1 alpha-2)
 * @param address - Адрес (опционально)
 * @returns ForeignParty объект
 */
export function createForeignParty(
  fullName: string,
  countryCode: string,
  address?: string,
): ForeignParty {
  return {
    fullName,
    address,
    country: {
      code: countryCode,
      name: countryCode,
    },
    isResident: false,
  };
}

/**
 * Функция для определения типа стороны
 * @param party - Объект стороны
 * @returns true если это местная компания, false если иностранная
 */
export function isPartyResident(party: ContractParty): party is Party {
  return party.isResident !== false;
}

/**
 * Функция для определения является ли сторона иностранной
 * @param party - Объект стороны
 * @returns true если это иностранная компания
 */
export function isForeignParty(party: ContractParty): party is ForeignParty {
  return party.isResident === false;
}

/**
 * Получить идентификатор стороны (ИНН или иностранный код)
 * @param party - Объект стороны
 * @returns Строка идентификатора
 */
export function getPartyId(party: ContractParty): string {
  if (party.pin) {
    return party.pin;
  }
  return party.fullName;
}
