package handlers

import (
	"context"
	"math"
	"time"

	"github.com/rusgainew/kkm-project-mks/catalog-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/catalog-query-server/internal/infrastructure/cache"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ListCatalogsWithFilter returns catalogs with filtering and sorting
func (h *CatalogQueryHandler) ListCatalogsWithFilter(ctx context.Context, req *pb.CatalogFilterRequest) (*pb.APIResponse, error) {
	ctx, span := tracer.Start(ctx, "ListCatalogsWithFilter",
		trace.WithAttributes(
			attribute.String("name", req.Name),
			attribute.String("number", req.Number),
			attribute.String("tnved_code", req.TnvedCode),
			attribute.String("gked_code", req.GkedCode),
			attribute.String("search_text", req.SearchText),
			attribute.String("sort_field", req.SortField),
			attribute.String("sort_order", req.SortOrder),
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("ListCatalogsWithFilter called",
		zap.String("name", req.Name),
		zap.String("number", req.Number),
		zap.String("tnved_code", req.TnvedCode),
		zap.String("gked_code", req.GkedCode),
		zap.String("search_text", req.SearchText),
		zap.String("sort_field", req.SortField),
		zap.String("sort_order", req.SortOrder),
		zap.Int32("page", req.Page),
		zap.Int32("size", req.Size),
	)

	// Validate sort field
	validSortFields := map[string]bool{
		"name":       true,
		"number":     true,
		"created_at": true,
		"":           true, // empty means default sort
	}
	if !validSortFields[req.SortField] {
		h.logger.Warn("Invalid sort field", zap.String("sort_field", req.SortField))
		h.metrics.ErrorsTotal.WithLabelValues("ListCatalogsWithFilter", "validation_error").Inc()
		return nil, status.Errorf(codes.InvalidArgument, "invalid sort field: %s", req.SortField)
	}

	// Validate sort order
	if req.SortOrder != "" && req.SortOrder != "ASC" && req.SortOrder != "DESC" {
		h.logger.Warn("Invalid sort order", zap.String("sort_order", req.SortOrder))
		h.metrics.ErrorsTotal.WithLabelValues("ListCatalogsWithFilter", "validation_error").Inc()
		return nil, status.Errorf(codes.InvalidArgument, "invalid sort order: %s (must be ASC or DESC)", req.SortOrder)
	}

	// Try to get from cache
	cacheKey := cache.GenerateCacheKey("catalog-query", "ListCatalogsWithFilter",
		req.Name, req.Number, req.TnvedCode, req.GkedCode, req.SearchText,
		req.SortField, req.SortOrder, req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("ListCatalogsWithFilter", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	// Cache miss - prepare filter and sort
	span.SetAttributes(attribute.Bool("cache_hit", false))
	filter := &ports.CatalogFilter{
		Name:       req.Name,
		Number:     req.Number,
		TnvedCode:  req.TnvedCode,
		GkedCode:   req.GkedCode,
		SearchText: req.SearchText,
	}

	sortField := req.SortField
	if sortField == "" {
		sortField = "created_at" // default sort
	}
	sortOrder := ports.SortOrder(req.SortOrder)
	if req.SortOrder == "" {
		sortOrder = ports.SortOrderDesc // default order
	}
	sort := &ports.CatalogSort{
		Field: sortField,
		Order: sortOrder,
	}

	// Query database
	catalogs, totalCount, err := h.repo.ListCatalogsWithFilter(ctx, filter, sort, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("ListCatalogsWithFilter").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to list catalogs with filter", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("ListCatalogsWithFilter", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("ListCatalogsWithFilter", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to list catalogs: %v", err)
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
	h.metrics.RequestsTotal.WithLabelValues("ListCatalogsWithFilter", "success").Inc()

	// Build response
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
