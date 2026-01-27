-- Откат миграции для analytics materialized views

-- Удаляем функцию обновления
DROP FUNCTION IF EXISTS refresh_analytics_materialized_views();

-- Удаляем индексы и материализованные представления
DROP INDEX IF EXISTS idx_analytics_operation_type_date;
DROP MATERIALIZED VIEW IF EXISTS analytics_operation_type;

DROP INDEX IF EXISTS idx_analytics_top_contractors_amount;
DROP INDEX IF EXISTS idx_analytics_top_contractors_month;
DROP MATERIALIZED VIEW IF EXISTS analytics_top_contractors;

DROP INDEX IF EXISTS idx_analytics_status_dist_date;
DROP MATERIALIZED VIEW IF EXISTS analytics_status_distribution;

DROP INDEX IF EXISTS idx_analytics_monthly_stats_month;
DROP MATERIALIZED VIEW IF EXISTS analytics_monthly_stats;

DROP INDEX IF EXISTS idx_analytics_weekly_stats_week;
DROP MATERIALIZED VIEW IF EXISTS analytics_weekly_stats;

DROP INDEX IF EXISTS idx_analytics_daily_stats_date;
DROP MATERIALIZED VIEW IF EXISTS analytics_daily_stats;
