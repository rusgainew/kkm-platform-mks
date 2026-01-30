// Файл catalog-query-server/internal/interfaces/grpc/handlers/catalog_query_handler.go содержит реализацию пакета handlers.
package handlers

import (
	"context"
	"math"
	"time"

	"github.com/rusgainew/kkm-project-mks/catalog-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/catalog-query-server/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/catalog-query-server/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var tracer = otel.Tracer("catalog-query-service")

type CatalogQueryHandler struct {
	pb.UnimplementedCatalogQueryServiceServer
	logger  *zap.Logger
	repo    ports.CatalogQueryRepository
	metrics *observability.MetricsCollector
	cache   *cache.RedisCache
}

func NewCatalogQueryHandler(logger *zap.Logger, repo ports.CatalogQueryRepository, metrics *observability.MetricsCollector, cache *cache.RedisCache) *CatalogQueryHandler {
	return &CatalogQueryHandler{
		logger:  logger,
		repo:    repo,
		metrics: metrics,
		cache:   cache,
	}
}

func (h *CatalogQueryHandler) ListCatalogs(ctx context.Context, req *pb.PageInfo) (*pb.APIResponse, error) {
	// Start tracing span
	ctx, span := tracer.Start(ctx, "ListCatalogs",
		trace.WithAttributes(
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("ListCatalogs called", zap.Int32("page", req.Page), zap.Int32("size", req.Size))

	// Try to get from cache first
	cacheKey := cache.GenerateCacheKey("catalog-query", "ListCatalogs", req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("ListCatalogs", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	// Cache miss - query database
	span.SetAttributes(attribute.Bool("cache_hit", false))
	catalogs, totalCount, err := h.repo.ListCatalogs(ctx, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("ListCatalogs").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to list catalogs", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("ListCatalogs", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("ListCatalogs", "error").Inc()
		return nil, err
	}

	// Calculate total pages
	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	// Add span attributes
	span.SetAttributes(
		attribute.Int("total_count", int(totalCount)),
		attribute.Int("total_pages", int(totalPages)),
		attribute.Int("catalogs_returned", len(catalogs)),
	)

	// Record success
	h.metrics.RequestsTotal.WithLabelValues("ListCatalogs", "success").Inc()

	// Build response with CatalogList data
	response := &pb.APIResponse{
		TotalElements: wrapperspb.Int32(totalCount),
		TotalPage:     wrapperspb.Int32(totalPages),
		Data: &pb.APIResponse_CatalogList{
			CatalogList: &pb.CatalogListResponse{
				Page:          req.Page,
				Size:          req.Size,
				TotalElements: totalCount,
				TotalPage:     totalPages,
				Catalogs:      catalogs,
			},
		},
	}

	// Store in cache
	if err := h.cache.Set(ctx, cacheKey, response); err != nil {
		h.logger.Warn("Failed to cache response", zap.String("key", cacheKey), zap.Error(err))
	}

	return response, nil
}
