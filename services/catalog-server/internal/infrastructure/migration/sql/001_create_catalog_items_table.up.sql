-- Создание таблицы для каталога товаров и услуг
CREATE SCHEMA IF NOT EXISTS catalogs;
SET search_path TO catalogs;

CREATE TABLE IF NOT EXISTS catalog_items (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL,
    name VARCHAR(500) NOT NULL,
    code VARCHAR(100) NOT NULL,
    description TEXT,
    unit_type VARCHAR(50) NOT NULL,
    price NUMERIC(15, 2) NOT NULL DEFAULT 0,
    vat_rate NUMERIC(5, 2) NOT NULL DEFAULT 0,
    category VARCHAR(200),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    tnved VARCHAR(20),
    gked VARCHAR(20),
    barcode VARCHAR(100),
    created_by UUID NOT NULL,
    updated_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_catalog_items_organization ON catalog_items(organization_id);
CREATE INDEX IF NOT EXISTS idx_catalog_items_code ON catalog_items(organization_id, code);
CREATE INDEX IF NOT EXISTS idx_catalog_items_name ON catalog_items(name);
CREATE INDEX IF NOT EXISTS idx_catalog_items_category ON catalog_items(category);
CREATE INDEX IF NOT EXISTS idx_catalog_items_is_active ON catalog_items(is_active);
CREATE INDEX IF NOT EXISTS idx_catalog_items_created_at ON catalog_items(created_at DESC);

-- Уникальный индекс для кода внутри организации
CREATE UNIQUE INDEX IF NOT EXISTS idx_catalog_items_org_code_unique ON catalog_items(organization_id, code);

-- Комментарии к таблице
COMMENT ON TABLE catalog_items IS 'Каталог товаров и услуг организаций';
COMMENT ON COLUMN catalog_items.id IS 'Уникальный идентификатор элемента';
COMMENT ON COLUMN catalog_items.organization_id IS 'ID организации-владельца';
COMMENT ON COLUMN catalog_items.name IS 'Наименование товара/услуги';
COMMENT ON COLUMN catalog_items.code IS 'Артикул/код товара';
COMMENT ON COLUMN catalog_items.description IS 'Описание товара/услуги';
COMMENT ON COLUMN catalog_items.unit_type IS 'Единица измерения (шт, кг, л, м)';
COMMENT ON COLUMN catalog_items.price IS 'Цена за единицу';
COMMENT ON COLUMN catalog_items.vat_rate IS 'Ставка НДС (0, 12, 20%)';
COMMENT ON COLUMN catalog_items.category IS 'Категория товара/услуги';
COMMENT ON COLUMN catalog_items.is_active IS 'Активен ли товар';
COMMENT ON COLUMN catalog_items.tnved IS 'Код ТНВЭД (для товаров)';
COMMENT ON COLUMN catalog_items.gked IS 'Код ГКЭД (для услуг)';
COMMENT ON COLUMN catalog_items.barcode IS 'Штрих-код';
