-- Создание таблицы для иностранных компаний (для query-сервиса)
-- Эта таблица должна быть реплицирована из foreign-company-server
CREATE TABLE IF NOT EXISTS foreign_companies (
    id BIGSERIAL PRIMARY KEY,
    pin VARCHAR(50) NOT NULL,
    full_name VARCHAR(500) NOT NULL,
    country_code VARCHAR(2),
    address TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID NOT NULL,
    updated_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы для оптимизации запросов чтения
CREATE INDEX IF NOT EXISTS idx_foreign_companies_pin ON foreign_companies(pin);
CREATE INDEX IF NOT EXISTS idx_foreign_companies_full_name ON foreign_companies USING gin(to_tsvector('english', full_name));
CREATE INDEX IF NOT EXISTS idx_foreign_companies_country ON foreign_companies(country_code);
CREATE INDEX IF NOT EXISTS idx_foreign_companies_is_active ON foreign_companies(is_active);
CREATE INDEX IF NOT EXISTS idx_foreign_companies_created_at ON foreign_companies(created_at DESC);

-- Уникальный индекс для PIN (только среди активных компаний)
CREATE UNIQUE INDEX IF NOT EXISTS idx_foreign_companies_pin_unique ON foreign_companies(LOWER(pin)) WHERE is_active = true;

-- Комментарии к таблице
COMMENT ON TABLE foreign_companies IS 'Справочник иностранных компаний-контрагентов (query-side для CQRS)';
COMMENT ON COLUMN foreign_companies.id IS 'Уникальный идентификатор компании';
COMMENT ON COLUMN foreign_companies.pin IS 'Идентификационный номер (Tax ID, VAT, TIN, Registration Number)';
COMMENT ON COLUMN foreign_companies.full_name IS 'Полное наименование компании на языке оригинала';
COMMENT ON COLUMN foreign_companies.country_code IS 'Код страны ISO 3166-1 alpha-2 (US, DE, FR, GB и т.д.)';
COMMENT ON COLUMN foreign_companies.address IS 'Юридический адрес компании';
COMMENT ON COLUMN foreign_companies.is_active IS 'Активна ли компания в системе';
COMMENT ON COLUMN foreign_companies.created_by IS 'UUID пользователя, создавшего запись';
COMMENT ON COLUMN foreign_companies.updated_by IS 'UUID пользователя, обновившего запись';
COMMENT ON COLUMN foreign_companies.created_at IS 'Дата и время создания записи';
COMMENT ON COLUMN foreign_companies.updated_at IS 'Дата и время последнего обновления';
