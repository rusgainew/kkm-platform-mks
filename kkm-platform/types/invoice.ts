/**
 * Типы для электронных счетов-фактур (ЭСФ) в ESF API
 * Включает основной тип Invoice, детали, запросы на создание/обновление
 */

import type { Party, ForeignParty, ContractParty } from "./party";
import type {
  ReferenceItem,
  VatTaxType,
  TaxInfo,
  UnitClassificationReference,
} from "./reference";

/**
 * Детальная строка в счете-фактуре
 * Представляет один товар или услугу в документе
 */
export interface ESFInvoiceDetail {
  invoiceUuid: string; // Ссылка на основной счет-фактуру
  baseCount: number; // Количество единиц товара
  price: number; // Цена за единицу
  amount: number; // Сумма: baseCount * price
  amountWithoutVAT: number; // Сумма без НДС
  amountVAT: number; // Сумма НДС
  amountST?: number; // Сумма НСП (налог специального назначения)
  fcdNumber?: string; // Номер ГТД (таможенной декларации)
  goodsName: string; // Наименование товара/услуги
  tnvedCode?: string; // ТНВЭД код товара
  gked?: string; // ГКЭД код (для услуг)
  unitClassification: UnitClassificationReference; // Единица измерения
  goodsType: ReferenceItem; // Тип товара (товар, услуга и т.д.)
  stTaxType?: TaxInfo; // Информация о специальном налоге
  goodsClassification?: ReferenceItem; // Классификация товара
}

/**
 * Запись из каталога для добавления в счет-фактуру
 * Используется при создании нового счета-фактуры
 */
export interface CatalogEntry {
  id: number; // ID записи в каталоге
  unitClassificationCode: string; // Код единицы измерения
  salesTaxCode: string; // Код налога на продажу
  customsAuthorityCode?: string; // Код таможни (опционально)
  quantity: number; // Количество
  price: number; // Цена за единицу
  vatAmount?: number; // Сумма НДС (вычисляется автоматически)
  salesTaxAmount?: number; // Сумма налога продажи
  amountWithoutTaxes?: number; // Сумма без налогов
  totalAmount?: number; // Итоговая сумма
}

/**
 * Основной тип для электронного счета-фактуры
 * Полная информация о ЭСФ документе со всеми деталями
 */
export interface ESFInvoice {
  // Идентификаторы и номера
  documentUuid: string; // Уникальный идентификатор документа
  invoiceNumber: string; // Номер счета-фактуры
  number?: string; // Альтернативный номер
  correctedReceiptUuid?: string; // UUID исправленного счета (для КСФ)

  // Даты
  invoiceDate?: string; // Дата счета-фактуры (ISO 8601)
  createdDate?: string; // Дата создания
  deliveryDate?: string; // Дата доставки
  correctedReceiptCreationDate?: string; // Дата создания исправленного счета

  // Финансовые суммы
  totalAmount: number; // Итоговая сумма
  totalAmountWithoutVAT?: number; // Сумма без НДС
  totalVATAmount?: number; // Итого НДС
  totalSTAmount?: number; // Итого НСП

  // Флаги и статусы
  isResident: boolean; // Резидент ли поставщик КР
  status?: ReferenceItem; // Статус документа (10-90)
  receiptType?: ReferenceItem; // Тип счета (РЕ, ПР, КСФ)

  // Участники сделки
  legalPerson?: Party; // Поставщик (резидент)
  contractor?: Party | ForeignParty; // Покупатель (резидент или иностранец)
  sellerBranchPin?: string; // ИНН филиала поставщика (если есть)

  // Справочники и классификации
  paymentType?: ReferenceItem; // Тип платежа (наличные, банк и т.д.)
  currency?: ReferenceItem; // Валюта (KGS, RUB, USD и т.д.)
  deliveryType?: ReferenceItem; // Тип доставки
  vatTaxType?: VatTaxType; // Ставка НДС
  country?: ReferenceItem; // Страна происхождения товаров
  deliveryCode?: string; // Код доставки
  paymentCode?: string; // Код платежа

  // Комментарии и примечания
  note?: string; // Примечание к счету
  comment?: string; // Дополнительный комментарий
  correctionReasonCode?: string; // Код причины коррекции
  correctionReasonName?: string; // Название причины коррекции

  // Контактная информация для платежей
  legalPersonBankAccount?: string; // Банковский счет поставщика
  contractorBankAccount?: string; // Банковский счет покупателя
  personalAccountNumber?: string; // Лицевой счет

  // Специальные поля
  ownedCrmReceiptCode?: string; // Код ЭСФ в системе учета
  foreignName?: string; // Наименование на иностранном языке
  isPriceWithoutTaxes?: boolean; // Цены указаны без налогов
  isIndustry?: boolean; // Специальная операция промышленности
  contractNumber?: string; // Номер договора поставки
  contractDate?: string; // Дата договора

  // Валютные операции
  exchangeRate?: number; // Курс обмена (если не КГС)
  totalCurrencyValue?: number; // Сумма в иностранной валюте
  totalCurrencyValueWithoutTaxes?: number; // Сумма в иностранной валюте без налогов

  // Финансовые отчеты (для КСФ и специальных случаев)
  openingBalances?: number; // Входящее сальдо
  assessedContributionsAmount?: number; // Сумма начисленных взносов
  paidAmount?: number; // Оплаченная сумма
  penaltiesAmount?: number; // Сумма штрафов
  finesAmount?: number; // Сумма пени
  closingBalances?: number; // Исходящее сальдо
  amountToBePaid?: number; // Сумма к оплате

  // Детали и строки
  details?: ESFInvoiceDetail[]; // Список товаров/услуг в счете

  // Служебные поля
  createdAt?: number; // Timestamp создания
  updatedAt?: number; // Timestamp обновления
  operationType?: string; // Тип операции (реализация, покупка и т.д.)
}

/**
 * Запрос на создание нового счета-фактуры
 * Минимальный набор обязательных полей
 */
export interface CreateInvoiceRequest {
  // Основная информация
  operationTypeCode: string; // Тип операции (10 = Реализация, 20 = Покупка)
  invoiceNumber?: string; // Номер счета (если не авто)

  // Даты
  deliveryDate: string; // Дата доставки (YYYY-MM-DD)
  invoiceDate?: string; // Дата счета-фактуры (если не текущая)

  // Стороны сделки
  contractorTin: string; // ИНН покупателя (или иностранный код)
  supplierBankAccount?: string; // Счет поставщика
  contractorBankAccount?: string; // Счет покупателя

  // Справочники
  deliveryTypeCode: string; // Код способа доставки
  paymentCode: string; // Код типа платежа
  currencyCode: string; // Код валюты (KGS, USD и т.д.)
  countryCode?: string; // Код страны происхождения
  taxRateVATCode: string; // Код ставки НДС (например "12%")

  // Специальные флаги
  isResident: boolean; // Резидент ли поставщик
  isPriceWithoutTaxes: boolean; // Цены указаны без налогов
  isIndustry?: boolean; // Специальная операция

  // Дополнительные условия
  currencyRate?: number; // Курс обмена (если не КГС)
  affiliateTin?: string; // ИНН филиала поставщика
  isBranchDataSent?: boolean; // Данные филиала отправлены

  // Контракт и комментарии
  supplyContractNumber?: string; // Номер контракта поставки
  contractStartDate?: string; // Дата начала контракта
  comment?: string; // Примечание

  // Специальные данные
  ownedCrmReceiptCode?: string; // Код в CRM системе
  foreignName?: string; // Название на иностранном языке

  // Финансовые данные для отчетов
  openingBalances?: number;
  assessedContributionsAmount?: number;
  paidAmount?: number;
  penaltiesAmount?: number;
  finesAmount?: number;
  closingBalances?: number;
  amountToBePaid?: number;
  personalAccountNumber?: string;

  // Товары в счете-фактуре
  catalogEntries: CatalogEntry[]; // Обязательно: хотя бы один товар/услуга
}

/**
 * Запрос на обновление существующего счета-фактуры
 * Может содержать частичные данные для изменения
 */
export interface UpdateInvoiceRequest extends Partial<CreateInvoiceRequest> {
  documentUuid: string; // Обязательно: ID документа
  version?: number; // Опционально: для optimistic locking
}

/**
 * Запрос на акцепт или отклонение счета-фактуры
 */
export interface AcceptOrRejectInvoiceRequest {
  invoiceUuid: string; // UUID счета
  action: "accept" | "reject"; // Действие
  reason?: string; // Причина (обязательна при отклонении)
  comment?: string; // Дополнительный комментарий
}

/**
 * Запрос на подпись счета-фактуры электронной подписью
 */
export interface SignInvoiceRequest {
  invoiceUuid: string; // UUID счета
  signature: string; // Электронная подпись (base64)
  certificateData?: string; // Данные сертификата (опционально)
  timestamp?: number; // Timestamp подписи
}

/**
 * Запрос на отзыв счета-фактуры
 */
export interface RevokeInvoiceRequest {
  invoiceUuid: string; // UUID счета
  reason: string; // Причина отзыва
  comment?: string; // Дополнительный комментарий
}

/**
 * Фильтры для поиска счетов-фактур
 */
export interface InvoiceFilters {
  status?: string; // Статус (10, 20, 30 и т.д.)
  invoiceNumber?: string; // Номер счета
  contractorTin?: string; // ИНН контрагента
  dateFrom?: string; // Дата с (YYYY-MM-DD)
  dateTo?: string; // Дата по (YYYY-MM-DD)
  createdBy?: string; // ID создателя
  isResident?: boolean; // Только резиденты
  currency?: string; // Валюта
  page?: number; // Страница (для пагинации)
  limit?: number; // Количество на странице
}

/**
 * Ответ на операцию с счетом-фактурой
 */
export interface InvoiceOperationResponse {
  success: boolean; // Успешна ли операция
  invoiceUuid?: string; // UUID счета
  invoiceNumber?: string; // Номер счета
  status?: string; // Текущий статус
  message?: string; // Сообщение о результате
  errors?: Record<string, string>; // Ошибки валидации если есть
  responseId?: string; // ID ответа от ESF API
}

/**
 * Вспомогательная функция для создания InvoiceDetail
 */
export function createInvoiceDetail(
  invoiceUuid: string,
  goodsName: string,
  quantity: number,
  price: number,
  unitCode: string,
  vatRate: string = "12",
): ESFInvoiceDetail {
  const amount = quantity * price;
  const vatAmount = (amount * Number(vatRate)) / (100 + Number(vatRate));
  const amountWithoutVAT = amount - vatAmount;

  return {
    invoiceUuid,
    baseCount: quantity,
    price,
    amount,
    amountWithoutVAT,
    amountVAT: vatAmount,
    goodsName,
    unitClassification: {
      code: unitCode,
      name: unitCode,
    },
    goodsType: {
      code: "1",
      name: "Товар",
    },
  };
}

/**
 * Вспомогательная функция для создания CatalogEntry
 */
export function createCatalogEntry(
  catalogId: number,
  quantity: number,
  price: number,
  unitCode: string = "796",
  taxCode: string = "10",
  vatRate: number = 12,
): CatalogEntry {
  const amountWithoutTaxes = quantity * price;
  const vatAmount = (amountWithoutTaxes * vatRate) / 100;
  const totalAmount = amountWithoutTaxes + vatAmount;

  return {
    id: catalogId,
    quantity,
    price,
    unitClassificationCode: unitCode,
    salesTaxCode: taxCode,
    vatAmount,
    amountWithoutTaxes,
    totalAmount,
  };
}
