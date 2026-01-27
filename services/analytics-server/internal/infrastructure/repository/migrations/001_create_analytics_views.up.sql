-- Analytics Materialized Views для повышения производительности
-- Обновление: каждые 5 минут через cronjob или pg_cron

-- 1. Материализованное представление для дневной статистики
CREATE MATERIALIZED VIEW IF NOT EXISTS analytics_daily_stats AS
SELECT 
    DATE_TRUNC('day', created_date) AS stat_date,
    COUNT(*) AS total_invoices,
    SUM(total_amount) AS total_revenue,
    AVG(total_amount) AS average_amount,
    COUNT(DISTINCT contractor_id) AS unique_contractors,
    COUNT(CASE WHEN status IN ('draft', 'sent') THEN 1 END) AS pending_count,
    COUNT(CASE WHEN status IN ('signed', 'accepted') THEN 1 END) AS approved_count,
    COUNT(CASE WHEN status IN ('rejected', 'revoked') THEN 1 END) AS rejected_count
FROM invoices
GROUP BY DATE_TRUNC('day', created_date)
ORDER BY stat_date DESC;

-- Индекс для быстрого поиска по дате
CREATE INDEX IF NOT EXISTS idx_analytics_daily_stats_date 
ON analytics_daily_stats(stat_date DESC);

-- 2. Материализованное представление для еженедельной статистики
CREATE MATERIALIZED VIEW IF NOT EXISTS analytics_weekly_stats AS
SELECT 
    DATE_TRUNC('week', created_date) AS stat_week,
    COUNT(*) AS total_invoices,
    SUM(total_amount) AS total_revenue,
    AVG(total_amount) AS average_amount,
    COUNT(DISTINCT contractor_id) AS unique_contractors
FROM invoices
GROUP BY DATE_TRUNC('week', created_date)
ORDER BY stat_week DESC;

CREATE INDEX IF NOT EXISTS idx_analytics_weekly_stats_week 
ON analytics_weekly_stats(stat_week DESC);

-- 3. Материализованное представление для месячной статистики
CREATE MATERIALIZED VIEW IF NOT EXISTS analytics_monthly_stats AS
SELECT 
    DATE_TRUNC('month', created_date) AS stat_month,
    COUNT(*) AS total_invoices,
    SUM(total_amount) AS total_revenue,
    AVG(total_amount) AS average_amount,
    COUNT(DISTINCT contractor_id) AS unique_contractors
FROM invoices
GROUP BY DATE_TRUNC('month', created_date)
ORDER BY stat_month DESC;

CREATE INDEX IF NOT EXISTS idx_analytics_monthly_stats_month 
ON analytics_monthly_stats(stat_month DESC);

-- 4. Материализованное представление для статистики по статусам
CREATE MATERIALIZED VIEW IF NOT EXISTS analytics_status_distribution AS
SELECT 
    DATE_TRUNC('day', created_date) AS stat_date,
    status,
    COUNT(*) AS count,
    SUM(total_amount) AS amount
FROM invoices
GROUP BY DATE_TRUNC('day', created_date), status
ORDER BY stat_date DESC, count DESC;

CREATE INDEX IF NOT EXISTS idx_analytics_status_dist_date 
ON analytics_status_distribution(stat_date DESC);

-- 5. Материализованное представление для топ контрагентов
CREATE MATERIALIZED VIEW IF NOT EXISTS analytics_top_contractors AS
SELECT 
    DATE_TRUNC('month', created_date) AS stat_month,
    contractor_id,
    COUNT(*) AS invoice_count,
    SUM(total_amount) AS total_amount,
    AVG(total_amount) AS average_amount
FROM invoices
GROUP BY DATE_TRUNC('month', created_date), contractor_id
ORDER BY stat_month DESC, total_amount DESC;

CREATE INDEX IF NOT EXISTS idx_analytics_top_contractors_month 
ON analytics_top_contractors(stat_month DESC);
CREATE INDEX IF NOT EXISTS idx_analytics_top_contractors_amount 
ON analytics_top_contractors(total_amount DESC);

-- 6. Материализованное представление для типов операций (resident/import)
CREATE MATERIALIZED VIEW IF NOT EXISTS analytics_operation_type AS
SELECT 
    DATE_TRUNC('day', created_date) AS stat_date,
    CASE 
        WHEN is_resident THEN 'local'
        ELSE 'import'
    END as operation_type,
    COUNT(*) AS count,
    SUM(total_amount) AS amount
FROM invoices
GROUP BY DATE_TRUNC('day', created_date), is_resident
ORDER BY stat_date DESC;

CREATE INDEX IF NOT EXISTS idx_analytics_operation_type_date 
ON analytics_operation_type(stat_date DESC);

-- Функция для обновления всех материализованных представлений
CREATE OR REPLACE FUNCTION refresh_analytics_materialized_views()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_daily_stats;
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_weekly_stats;
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_monthly_stats;
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_status_distribution;
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_top_contractors;
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_operation_type;
END;
$$ LANGUAGE plpgsql;

-- Комментарии для документации
COMMENT ON MATERIALIZED VIEW analytics_daily_stats IS 
'Дневная агрегированная статистика по накладным. Обновляется каждые 5 минут.';

COMMENT ON MATERIALIZED VIEW analytics_weekly_stats IS 
'Еженедельная агрегированная статистика по накладным. Обновляется каждые 5 минут.';

COMMENT ON MATERIALIZED VIEW analytics_monthly_stats IS 
'Месячная агрегированная статистика по накладным. Обновляется каждые 5 минут.';

COMMENT ON MATERIALIZED VIEW analytics_status_distribution IS 
'Распределение накладных по статусам за день. Обновляется каждые 5 минут.';

COMMENT ON MATERIALIZED VIEW analytics_top_contractors IS 
'Топ контрагентов по выручке за месяц. Обновляется каждые 5 минут.';

COMMENT ON MATERIALIZED VIEW analytics_operation_type IS 
'Распределение по типам операций (местные/импортные). Обновляется каждые 5 минут.';

-- Инструкции по настройке автоматического обновления:
-- 
-- Вариант 1: Использовать pg_cron (рекомендуется):
-- SELECT cron.schedule('refresh-analytics', '*/5 * * * *', 'SELECT refresh_analytics_materialized_views()');
--
-- Вариант 2: Создать cronjob в операционной системе:
-- */5 * * * * psql -U <user> -d <database> -c "SELECT refresh_analytics_materialized_views();"
--
-- Вариант 3: Вызывать из приложения периодически
