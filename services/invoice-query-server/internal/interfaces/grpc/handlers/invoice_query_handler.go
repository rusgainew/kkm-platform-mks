// Файл invoice-query-server/internal/interfaces/grpc/handlers/invoice_query_handler.go содержит реализацию пакета handlers.
package handlers

import (
	"context"
	"math"
	"time"

	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var tracer = otel.Tracer("invoice-query-service")

type InvoiceQueryHandler struct {
	pb.UnimplementedInvoiceQueryServiceServer
	logger  *zap.Logger
	repo    ports.InvoiceQueryRepository
	metrics *observability.MetricsCollector
	cache   *cache.RedisCache
}

func NewInvoiceQueryHandler(logger *zap.Logger, repo ports.InvoiceQueryRepository, metrics *observability.MetricsCollector, cache *cache.RedisCache) *InvoiceQueryHandler {
	return &InvoiceQueryHandler{
		logger:  logger,
		repo:    repo,
		metrics: metrics,
		cache:   cache,
	}
}

func (h *InvoiceQueryHandler) ListInvoices(ctx context.Context, req *pb.PageInfo) (*pb.APIResponse, error) {
	// Start tracing span
	ctx, span := tracer.Start(ctx, "ListInvoices",
		trace.WithAttributes(
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("ListInvoices called", zap.Int32("page", req.Page), zap.Int32("size", req.Size))

	// Try to get from cache first
	cacheKey := cache.GenerateCacheKey("invoice-query", "ListInvoices", req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("ListInvoices", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	// Cache miss - query database
	span.SetAttributes(attribute.Bool("cache_hit", false))
	invoices, totalCount, err := h.repo.ListInvoices(ctx, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("ListInvoices").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to list invoices", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("ListInvoices", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("ListInvoices", "error").Inc()
		return nil, err
	}

	// Calculate total pages
	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	// Add span attributes
	span.SetAttributes(
		attribute.Int("total_count", int(totalCount)),
		attribute.Int("total_pages", int(totalPages)),
		attribute.Int("invoices_returned", len(invoices)),
	)

	// Record success
	h.metrics.RequestsTotal.WithLabelValues("ListInvoices", "success").Inc()

	// Build response with InvoiceList data
	response := &pb.APIResponse{
		TotalElements: wrapperspb.Int32(totalCount),
		TotalPage:     wrapperspb.Int32(totalPages),
		Data: &pb.APIResponse_InvoiceList{
			InvoiceList: &pb.InvoiceListResponse{
				Page:          req.Page,
				Size:          req.Size,
				TotalElements: totalCount,
				TotalPage:     totalPages,
				Invoices:      invoices,
			},
		},
	}

	// Store in cache
	if err := h.cache.Set(ctx, cacheKey, response); err != nil {
		h.logger.Warn("Failed to cache response", zap.String("key", cacheKey), zap.Error(err))
	}

	return response, nil
}

func (h *InvoiceQueryHandler) ListInvoiceDetails(ctx context.Context, req *pb.InvoiceDetailsRequest) (*pb.APIResponse, error) {
	// Start tracing span
	ctx, span := tracer.Start(ctx, "ListInvoiceDetails",
		trace.WithAttributes(
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
			attribute.String("invoice_uuid", req.InvoiceUuid),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("ListInvoiceDetails called",
		zap.Int32("page", req.Page),
		zap.Int32("size", req.Size),
		zap.String("invoice_uuid", req.InvoiceUuid))

	// Try to get from cache first
	cacheKey := cache.GenerateCacheKey("invoice-query", "ListInvoiceDetails", req.InvoiceUuid, req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("ListInvoiceDetails", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	// Cache miss - query database
	span.SetAttributes(attribute.Bool("cache_hit", false))
	details, totalCount, err := h.repo.ListInvoiceDetails(ctx, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("ListInvoiceDetails").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to list invoice details", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("ListInvoiceDetails", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("ListInvoiceDetails", "error").Inc()
		return nil, err
	}

	// Calculate total pages
	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	// Add span attributes
	span.SetAttributes(
		attribute.Int("total_count", int(totalCount)),
		attribute.Int("total_pages", int(totalPages)),
		attribute.Int("details_returned", len(details)),
	)

	// Record success
	h.metrics.RequestsTotal.WithLabelValues("ListInvoiceDetails", "success").Inc()

	// Build response with DetailList data
	response := &pb.APIResponse{
		TotalElements: wrapperspb.Int32(totalCount),
		TotalPage:     wrapperspb.Int32(totalPages),
		Data: &pb.APIResponse_DetailList{
			DetailList: &pb.DetailListResponse{
				Page:          req.Page,
				Size:          req.Size,
				TotalElements: totalCount,
				TotalPage:     totalPages,
				Details:       details,
			},
		},
	}

	// Store in cache
	if err := h.cache.Set(ctx, cacheKey, response); err != nil {
		h.logger.Warn("Failed to cache response", zap.String("key", cacheKey), zap.Error(err))
	}

	return response, nil
}
