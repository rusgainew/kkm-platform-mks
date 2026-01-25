package handlers

import (
	"context"
	"math"
	"time"

	"github.com/rusgainew/kkm-project-mks/bank-account-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/bank-account-query-server/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/bank-account-query-server/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var tracer = otel.Tracer("bank-account-query-service")

type BankAccountQueryHandler struct {
	pb.UnimplementedBankAccountQueryServiceServer
	logger  *zap.Logger
	repo    ports.BankAccountQueryRepository
	metrics *observability.MetricsCollector
	cache   *cache.RedisCache
}

func NewBankAccountQueryHandler(logger *zap.Logger, repo ports.BankAccountQueryRepository, metrics *observability.MetricsCollector, cache *cache.RedisCache) *BankAccountQueryHandler {
	return &BankAccountQueryHandler{
		logger:  logger,
		repo:    repo,
		metrics: metrics,
		cache:   cache,
	}
}

func (h *BankAccountQueryHandler) ListBankAccounts(ctx context.Context, req *pb.PageInfo) (*pb.APIResponse, error) {
	// Start tracing span
	ctx, span := tracer.Start(ctx, "ListBankAccounts",
		trace.WithAttributes(
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("ListBankAccounts called", zap.Int32("page", req.Page), zap.Int32("size", req.Size))

	// Try to get from cache first
	cacheKey := cache.GenerateCacheKey("bank-account-query", "ListBankAccounts", req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("ListBankAccounts", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	// Cache miss - query database
	span.SetAttributes(attribute.Bool("cache_hit", false))
	accounts, totalCount, err := h.repo.ListBankAccounts(ctx, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("ListBankAccounts").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to list bank accounts", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("ListBankAccounts", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("ListBankAccounts", "error").Inc()
		return nil, err
	}

	// Calculate total pages
	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	// Add span attributes
	span.SetAttributes(
		attribute.Int("total_count", int(totalCount)),
		attribute.Int("total_pages", int(totalPages)),
		attribute.Int("accounts_returned", len(accounts)),
	)

	// Record success
	h.metrics.RequestsTotal.WithLabelValues("ListBankAccounts", "success").Inc()

	// Build response with BankAccountList data
	response := &pb.APIResponse{
		TotalElements: wrapperspb.Int32(totalCount),
		TotalPage:     wrapperspb.Int32(totalPages),
		Data: &pb.APIResponse_BankAccountList{
			BankAccountList: &pb.BankAccountListResponse{
				Page:          req.Page,
				Size:          req.Size,
				TotalElements: totalCount,
				TotalPage:     totalPages,
				BankAccounts:  accounts,
			},
		},
	}

	// Store in cache
	if err := h.cache.Set(ctx, cacheKey, response); err != nil {
		h.logger.Warn("Failed to cache response", zap.String("key", cacheKey), zap.Error(err))
	}

	return response, nil
}
