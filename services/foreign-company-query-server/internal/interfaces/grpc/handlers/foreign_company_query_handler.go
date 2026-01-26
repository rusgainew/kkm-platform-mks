package handlers

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/rusgainew/kkm-project-mks/foreign-company-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/foreign-company-query-server/internal/infrastructure/cache"
	"github.com/rusgainew/kkm-project-mks/foreign-company-query-server/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var tracer = otel.Tracer("foreign-company-query-service")

type ForeignCompanyQueryHandler struct {
	pb.UnimplementedForeignCompanyQueryServiceServer
	logger  *zap.Logger
	repo    ports.ForeignCompanyQueryRepository
	metrics *observability.MetricsCollector
	cache   *cache.RedisCache
}

func NewForeignCompanyQueryHandler(logger *zap.Logger, repo ports.ForeignCompanyQueryRepository, metrics *observability.MetricsCollector, cache *cache.RedisCache) *ForeignCompanyQueryHandler {
	return &ForeignCompanyQueryHandler{
		logger:  logger,
		repo:    repo,
		metrics: metrics,
		cache:   cache,
	}
}

func (h *ForeignCompanyQueryHandler) ListForeignCompanies(ctx context.Context, req *pb.PageInfo) (*pb.APIResponse, error) {
	// Start tracing span
	ctx, span := tracer.Start(ctx, "ListForeignCompanies",
		trace.WithAttributes(
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("ListForeignCompanies called", zap.Int32("page", req.Page), zap.Int32("size", req.Size))

	// Try to get from cache first
	cacheKey := cache.GenerateCacheKey("foreign-company-query", "ListForeignCompanies", req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("ListForeignCompanies", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	// Cache miss - query database
	span.SetAttributes(attribute.Bool("cache_hit", false))
	companies, totalCount, err := h.repo.ListForeignCompanies(ctx, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("ListForeignCompanies").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to list foreign companies", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("ListForeignCompanies", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("ListForeignCompanies", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to list foreign companies: %v", err)
	}

	// Calculate total pages
	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	// Add span attributes
	span.SetAttributes(
		attribute.Int("total_count", int(totalCount)),
		attribute.Int("total_pages", int(totalPages)),
		attribute.Int("companies_returned", len(companies)),
	)

	// Record success
	h.metrics.RequestsTotal.WithLabelValues("ListForeignCompanies", "success").Inc()

	// Build response with ForeignCompanyList data
	// Note: Using CatalogList temporarily since ForeignCompanyList is not defined in proto
	// Convert ForeignCompany to Catalog format
	catalogs := make([]*dictionaries.Catalog, len(companies))
	for i, fc := range companies {
		catalogs[i] = &dictionaries.Catalog{
			Id:        strconv.FormatInt(fc.Id, 10),
			Name:      fc.FullName,
			Number:    fc.Pin,
			TnvedCode: "",
			GkedCode:  "",
		}
	}

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

// ListForeignCompaniesWithFilter returns foreign companies with filtering and sorting
func (h *ForeignCompanyQueryHandler) ListForeignCompaniesWithFilter(ctx context.Context, req *pb.ForeignCompanyFilterRequest) (*pb.APIResponse, error) {
	ctx, span := tracer.Start(ctx, "ListForeignCompaniesWithFilter",
		trace.WithAttributes(
			attribute.String("pin", req.Pin),
			attribute.String("full_name", req.FullName),
			attribute.String("country_code", req.CountryCode),
			attribute.String("search_text", req.SearchText),
			attribute.String("sort_field", req.SortField),
			attribute.String("sort_order", req.SortOrder),
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("ListForeignCompaniesWithFilter called",
		zap.String("pin", req.Pin),
		zap.String("full_name", req.FullName),
		zap.String("country_code", req.CountryCode),
		zap.String("search_text", req.SearchText),
		zap.String("sort_field", req.SortField),
		zap.String("sort_order", req.SortOrder),
		zap.Int32("page", req.Page),
		zap.Int32("size", req.Size),
	)

	// Validate sort field
	validSortFields := map[string]bool{
		"pin":          true,
		"full_name":    true,
		"country_code": true,
		"":             true, // empty means default sort
	}
	if !validSortFields[req.SortField] {
		h.logger.Warn("Invalid sort field", zap.String("sort_field", req.SortField))
		h.metrics.ErrorsTotal.WithLabelValues("ListForeignCompaniesWithFilter", "validation_error").Inc()
		return nil, status.Errorf(codes.InvalidArgument, "invalid sort field: %s", req.SortField)
	}

	// Validate sort order
	if req.SortOrder != "" && req.SortOrder != "ASC" && req.SortOrder != "DESC" {
		h.logger.Warn("Invalid sort order", zap.String("sort_order", req.SortOrder))
		h.metrics.ErrorsTotal.WithLabelValues("ListForeignCompaniesWithFilter", "validation_error").Inc()
		return nil, status.Errorf(codes.InvalidArgument, "invalid sort order: %s (must be ASC or DESC)", req.SortOrder)
	}

	// Try to get from cache
	cacheKey := cache.GenerateCacheKey("foreign-company-query", "ListForeignCompaniesWithFilter",
		req.Pin, req.FullName, req.CountryCode, req.SearchText,
		req.SortField, req.SortOrder, req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("ListForeignCompaniesWithFilter", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	// Cache miss - prepare filter and sort
	span.SetAttributes(attribute.Bool("cache_hit", false))

	filter := &ports.ForeignCompanyFilter{
		PIN:         req.Pin,
		FullName:    req.FullName,
		CountryCode: req.CountryCode,
		SearchText:  req.SearchText,
	}

	sortField := req.SortField
	if sortField == "" {
		sortField = "full_name" // default sort
	}
	sortOrder := ports.SortOrder(req.SortOrder)
	if req.SortOrder == "" {
		sortOrder = ports.SortOrderAsc // default order
	}
	sort := &ports.ForeignCompanySort{
		Field: sortField,
		Order: sortOrder,
	}

	// Query database
	companies, totalCount, err := h.repo.ListForeignCompaniesWithFilter(ctx, filter, sort, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("ListForeignCompaniesWithFilter").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to list foreign companies with filter", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("ListForeignCompaniesWithFilter", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("ListForeignCompaniesWithFilter", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to list foreign companies: %v", err)
	}

	// Calculate total pages
	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	// Add span attributes
	span.SetAttributes(
		attribute.Int("total_count", int(totalCount)),
		attribute.Int("total_pages", int(totalPages)),
		attribute.Int("companies_returned", len(companies)),
	)

	// Record success
	h.metrics.RequestsTotal.WithLabelValues("ListForeignCompaniesWithFilter", "success").Inc()

	// Build response - convert to Catalog format temporarily
	catalogs := make([]*dictionaries.Catalog, len(companies))
	for i, fc := range companies {
		catalogs[i] = &dictionaries.Catalog{
			Id:        strconv.FormatInt(fc.Id, 10),
			Name:      fc.FullName,
			Number:    fc.Pin,
			TnvedCode: "",
			GkedCode:  "",
		}
	}

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

// SearchForeignCompanies performs full-text search on foreign companies
func (h *ForeignCompanyQueryHandler) SearchForeignCompanies(ctx context.Context, req *pb.SearchRequest) (*pb.APIResponse, error) {
	ctx, span := tracer.Start(ctx, "SearchForeignCompanies",
		trace.WithAttributes(
			attribute.String("search_text", req.SearchText),
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("SearchForeignCompanies called",
		zap.String("search_text", req.SearchText),
		zap.Int32("page", req.Page),
		zap.Int32("size", req.Size))

	// Validate search text
	if req.SearchText == "" {
		h.logger.Warn("Empty search text")
		h.metrics.ErrorsTotal.WithLabelValues("SearchForeignCompanies", "validation_error").Inc()
		return nil, status.Error(codes.InvalidArgument, "search_text is required")
	}

	// Try to get from cache
	cacheKey := cache.GenerateCacheKey("foreign-company-query", "SearchForeignCompanies",
		req.SearchText, req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("SearchForeignCompanies", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	// Cache miss - search database
	span.SetAttributes(attribute.Bool("cache_hit", false))
	companies, totalCount, err := h.repo.SearchForeignCompanies(ctx, req.SearchText, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("SearchForeignCompanies").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to search foreign companies", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("SearchForeignCompanies", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("SearchForeignCompanies", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to search foreign companies: %v", err)
	}

	// Calculate total pages
	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	// Add span attributes
	span.SetAttributes(
		attribute.Int("total_count", int(totalCount)),
		attribute.Int("total_pages", int(totalPages)),
		attribute.Int("companies_returned", len(companies)),
	)

	// Record success
	h.metrics.RequestsTotal.WithLabelValues("SearchForeignCompanies", "success").Inc()

	// Build response - convert to Catalog format temporarily
	catalogs := make([]*dictionaries.Catalog, len(companies))
	for i, fc := range companies {
		catalogs[i] = &dictionaries.Catalog{
			Id:        strconv.FormatInt(fc.Id, 10),
			Name:      fc.FullName,
			Number:    fc.Pin,
			TnvedCode: "",
			GkedCode:  "",
		}
	}

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
