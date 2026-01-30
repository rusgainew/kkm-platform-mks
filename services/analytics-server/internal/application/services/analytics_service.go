// Файл analytics-server/internal/application/services/analytics_service.go содержит реализацию пакета services.
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rusgainew/kkm-project-mks/analytics-server/internal/domain/repository"
	"github.com/rusgainew/kkm-project-mks/analytics-server/internal/infrastructure/cache"
	"go.uber.org/zap"
)

// AnalyticsService сервис для аналитики с кешированием
type AnalyticsService struct {
	repo     repository.AnalyticsRepository
	cache    *cache.RedisCache
	logger   *zap.Logger
	cacheTTL time.Duration
}

// NewAnalyticsService создает новый AnalyticsService
func NewAnalyticsService(
	repo repository.AnalyticsRepository,
	cache *cache.RedisCache,
	logger *zap.Logger,
) *AnalyticsService {
	return &AnalyticsService{
		repo:     repo,
		cache:    cache,
		logger:   logger,
		cacheTTL: 5 * time.Minute, // 5 минут TTL для аналитики
	}
}

// GetStats возвращает агрегированную статистику за период с кешированием
func (s *AnalyticsService) GetStats(ctx context.Context, startDate, endDate time.Time) (*repository.AnalyticsStats, error) {
	cacheKey := s.buildCacheKey("stats", startDate, endDate)

	// Пытаемся получить из кеша
	if s.cache != nil {
		var cachedStats repository.AnalyticsStats
		err := s.cache.Get(ctx, cacheKey, &cachedStats)
		if err == nil {
			s.logger.Debug("Analytics stats retrieved from cache", zap.String("key", cacheKey))
			return &cachedStats, nil
		}
	}

	// Получаем из базы
	stats, err := s.repo.GetStats(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("get stats from repo: %w", err)
	}

	// Сохраняем в кеш
	if s.cache != nil {
		if err := s.cache.Set(ctx, cacheKey, stats); err != nil {
			s.logger.Warn("Failed to cache stats", zap.Error(err))
		}
	}

	return stats, nil
}

// GetSalesData возвращает данные продаж с кешированием
func (s *AnalyticsService) GetSalesData(ctx context.Context, startDate, endDate time.Time, granularity string) ([]repository.SalesDataPoint, error) {
	cacheKey := s.buildCacheKeyWithGranularity("sales", startDate, endDate, granularity)

	// Пытаемся получить из кеша
	if s.cache != nil {
		var cachedData []repository.SalesDataPoint
		err := s.cache.Get(ctx, cacheKey, &cachedData)
		if err == nil {
			s.logger.Debug("Sales data retrieved from cache", zap.String("key", cacheKey))
			return cachedData, nil
		}
	}

	// Получаем из базы
	data, err := s.repo.GetSalesData(ctx, startDate, endDate, granularity)
	if err != nil {
		return nil, fmt.Errorf("get sales data from repo: %w", err)
	}

	// Сохраняем в кеш
	if s.cache != nil {
		if err := s.cache.Set(ctx, cacheKey, data); err != nil {
			s.logger.Warn("Failed to cache sales data", zap.Error(err))
		}
	}

	return data, nil
}

// GetStatusDistribution возвращает распределение по статусам с кешированием
func (s *AnalyticsService) GetStatusDistribution(ctx context.Context, startDate, endDate time.Time) ([]repository.StatusDistribution, error) {
	cacheKey := s.buildCacheKey("status_dist", startDate, endDate)

	// Пытаемся получить из кеша
	if s.cache != nil {
		var cachedDist []repository.StatusDistribution
		err := s.cache.Get(ctx, cacheKey, &cachedDist)
		if err == nil {
			s.logger.Debug("Status distribution retrieved from cache", zap.String("key", cacheKey))
			return cachedDist, nil
		}
	}

	// Получаем из базы
	dist, err := s.repo.GetStatusDistribution(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("get status distribution from repo: %w", err)
	}

	// Сохраняем в кеш
	if s.cache != nil {
		if err := s.cache.Set(ctx, cacheKey, dist); err != nil {
			s.logger.Warn("Failed to cache status distribution", zap.Error(err))
		}
	}

	return dist, nil
}

// GetOperationTypeDistribution возвращает распределение по типам операций с кешированием
func (s *AnalyticsService) GetOperationTypeDistribution(ctx context.Context, startDate, endDate time.Time) ([]repository.OperationTypeDistribution, error) {
	cacheKey := s.buildCacheKey("operation_dist", startDate, endDate)

	// Пытаемся получить из кеша
	if s.cache != nil {
		var cachedDist []repository.OperationTypeDistribution
		err := s.cache.Get(ctx, cacheKey, &cachedDist)
		if err == nil {
			s.logger.Debug("Operation type distribution retrieved from cache", zap.String("key", cacheKey))
			return cachedDist, nil
		}
	}

	// Получаем из базы
	dist, err := s.repo.GetOperationTypeDistribution(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("get operation type distribution from repo: %w", err)
	}

	// Сохраняем в кеш
	if s.cache != nil {
		if err := s.cache.Set(ctx, cacheKey, dist); err != nil {
			s.logger.Warn("Failed to cache operation type distribution", zap.Error(err))
		}
	}

	return dist, nil
}

// GetTopContractors возвращает топ контрагентов с кешированием
func (s *AnalyticsService) GetTopContractors(ctx context.Context, startDate, endDate time.Time, limit int32) ([]repository.TopContractorData, error) {
	cacheKey := s.buildCacheKeyWithLimit("top_contractors", startDate, endDate, limit)

	// Пытаемся получить из кеша
	if s.cache != nil {
		var cachedData []repository.TopContractorData
		err := s.cache.Get(ctx, cacheKey, &cachedData)
		if err == nil {
			s.logger.Debug("Top contractors retrieved from cache", zap.String("key", cacheKey))
			return cachedData, nil
		}
	}

	// Получаем из базы
	data, err := s.repo.GetTopContractors(ctx, startDate, endDate, limit)
	if err != nil {
		return nil, fmt.Errorf("get top contractors from repo: %w", err)
	}

	// Сохраняем в кеш
	if s.cache != nil {
		if err := s.cache.Set(ctx, cacheKey, data); err != nil {
			s.logger.Warn("Failed to cache top contractors", zap.Error(err))
		}
	}

	return data, nil
}

// GetMonthlyRevenue возвращает месячную выручку с кешированием
func (s *AnalyticsService) GetMonthlyRevenue(ctx context.Context, startDate, endDate time.Time) ([]repository.MonthlyRevenueData, error) {
	cacheKey := s.buildCacheKey("monthly_revenue", startDate, endDate)

	// Пытаемся получить из кеша
	if s.cache != nil {
		var cachedData []repository.MonthlyRevenueData
		err := s.cache.Get(ctx, cacheKey, &cachedData)
		if err == nil {
			s.logger.Debug("Monthly revenue retrieved from cache", zap.String("key", cacheKey))
			return cachedData, nil
		}
	}

	// Получаем из базы
	data, err := s.repo.GetMonthlyRevenue(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("get monthly revenue from repo: %w", err)
	}

	// Сохраняем в кеш
	if s.cache != nil {
		if err := s.cache.Set(ctx, cacheKey, data); err != nil {
			s.logger.Warn("Failed to cache monthly revenue", zap.Error(err))
		}
	}

	return data, nil
}

// InvalidateCache инвалидирует все кеши аналитики
func (s *AnalyticsService) InvalidateCache(ctx context.Context) error {
	if s.cache == nil {
		return nil
	}

	pattern := "analytics:*"
	if err := s.cache.DeleteByPattern(ctx, pattern); err != nil {
		s.logger.Error("Failed to invalidate analytics cache", zap.Error(err))
		return fmt.Errorf("invalidate cache: %w", err)
	}

	s.logger.Info("Analytics cache invalidated", zap.String("pattern", pattern))
	return nil
}

// buildCacheKey создает ключ кеша для аналитики
func (s *AnalyticsService) buildCacheKey(prefix string, startDate, endDate time.Time) string {
	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")
	return fmt.Sprintf("analytics:%s:%s:%s", prefix, startStr, endStr)
}

// buildCacheKeyWithGranularity создает ключ кеша с учетом гранулярности
func (s *AnalyticsService) buildCacheKeyWithGranularity(prefix string, startDate, endDate time.Time, granularity string) string {
	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")
	return fmt.Sprintf("analytics:%s:%s:%s:%s", prefix, startStr, endStr, granularity)
}

// buildCacheKeyWithLimit создает ключ кеша с учетом лимита
func (s *AnalyticsService) buildCacheKeyWithLimit(prefix string, startDate, endDate time.Time, limit int32) string {
	startStr := startDate.Format("2006-01-02")
	endStr := endDate.Format("2006-01-02")
	return fmt.Sprintf("analytics:%s:%s:%s:%d", prefix, startStr, endStr, limit)
}

// CacheInfo информация о кеше для отладки
type CacheInfo struct {
	Enabled bool          `json:"enabled"`
	TTL     time.Duration `json:"ttl"`
}

// GetCacheInfo возвращает информацию о кеше
func (s *AnalyticsService) GetCacheInfo() CacheInfo {
	return CacheInfo{
		Enabled: s.cache != nil,
		TTL:     s.cacheTTL,
	}
}

// String возвращает JSON-представление CacheInfo
func (ci CacheInfo) String() string {
	data, _ := json.Marshal(ci)
	return string(data)
}
