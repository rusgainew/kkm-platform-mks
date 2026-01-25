package handlers

import (
	"context"
	"math"
	"time"

	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/invoice-query-server/internal/infrastructure/cache"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"github.com/rusgainew/kkm-project-mks/proto-lib/entities"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ListInvoicesWithFilter returns invoices with filtering and sorting
func (h *InvoiceQueryHandler) ListInvoicesWithFilter(ctx context.Context, req *pb.InvoiceFilterRequest) (*pb.APIResponse, error) {
	ctx, span := tracer.Start(ctx, "ListInvoicesWithFilter",
		trace.WithAttributes(
			attribute.String("invoice_number", req.InvoiceNumber),
			attribute.String("date_from", req.DateFrom),
			attribute.String("date_to", req.DateTo),
			attribute.String("search_text", req.SearchText),
			attribute.String("sort_field", req.SortField),
			attribute.String("sort_order", req.SortOrder),
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("ListInvoicesWithFilter called",
		zap.String("invoice_number", req.InvoiceNumber),
		zap.String("date_from", req.DateFrom),
		zap.String("date_to", req.DateTo),
		zap.Float64("min_amount", req.MinAmount),
		zap.Float64("max_amount", req.MaxAmount),
		zap.String("search_text", req.SearchText),
		zap.String("sort_field", req.SortField),
		zap.String("sort_order", req.SortOrder),
		zap.Int32("page", req.Page),
		zap.Int32("size", req.Size),
	)

	// Validate sort field
	validSortFields := map[string]bool{
		"invoice_date":   true,
		"total_amount":   true,
		"invoice_number": true,
		"created_at":     true,
		"":               true, // empty means default sort
	}
	if !validSortFields[req.SortField] {
		h.logger.Warn("Invalid sort field", zap.String("sort_field", req.SortField))
		h.metrics.ErrorsTotal.WithLabelValues("ListInvoicesWithFilter", "validation_error").Inc()
		return nil, status.Errorf(codes.InvalidArgument, "invalid sort field: %s", req.SortField)
	}

	// Validate sort order
	if req.SortOrder != "" && req.SortOrder != "ASC" && req.SortOrder != "DESC" {
		h.logger.Warn("Invalid sort order", zap.String("sort_order", req.SortOrder))
		h.metrics.ErrorsTotal.WithLabelValues("ListInvoicesWithFilter", "validation_error").Inc()
		return nil, status.Errorf(codes.InvalidArgument, "invalid sort order: %s (must be ASC or DESC)", req.SortOrder)
	}

	// Try to get from cache (with filter params in key)
	cacheKey := cache.GenerateCacheKey("invoice-query", "ListInvoicesWithFilter",
		req.InvoiceNumber, req.DateFrom, req.DateTo,
		req.MinAmount, req.MaxAmount, req.SearchText,
		req.SortField, req.SortOrder, req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("ListInvoicesWithFilter", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	// Cache miss - prepare filter and sort
	span.SetAttributes(attribute.Bool("cache_hit", false))
	filter := &ports.InvoiceFilter{
		InvoiceNumber: req.InvoiceNumber,
		DateFrom:      req.DateFrom,
		DateTo:        req.DateTo,
		MinAmount:     req.MinAmount,
		MaxAmount:     req.MaxAmount,
		SearchText:    req.SearchText,
	}

	sortField := req.SortField
	if sortField == "" {
		sortField = "created_at" // default sort
	}
	sortOrder := ports.SortOrder(req.SortOrder)
	if req.SortOrder == "" {
		sortOrder = ports.SortOrderDesc // default order
	}
	sort := &ports.InvoiceSort{
		Field: sortField,
		Order: sortOrder,
	}

	// Query database
	invoices, totalCount, err := h.repo.ListInvoicesWithFilter(ctx, filter, sort, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("ListInvoicesWithFilter").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to list invoices with filter", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("ListInvoicesWithFilter", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("ListInvoicesWithFilter", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to list invoices: %v", err)
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
	h.metrics.RequestsTotal.WithLabelValues("ListInvoicesWithFilter", "success").Inc()

	// Build response
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

// SearchInvoices performs full-text search on invoices
func (h *InvoiceQueryHandler) SearchInvoices(ctx context.Context, req *pb.SearchRequest) (*pb.APIResponse, error) {
	ctx, span := tracer.Start(ctx, "SearchInvoices",
		trace.WithAttributes(
			attribute.String("search_text", req.SearchText),
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("SearchInvoices called",
		zap.String("search_text", req.SearchText),
		zap.Int32("page", req.Page),
		zap.Int32("size", req.Size),
	)

	if req.SearchText == "" {
		h.metrics.ErrorsTotal.WithLabelValues("SearchInvoices", "validation_error").Inc()
		return nil, status.Error(codes.InvalidArgument, "search_text is required")
	}

	// Try cache
	cacheKey := cache.GenerateCacheKey("invoice-query", "SearchInvoices", req.SearchText, req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("SearchInvoices", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))
	invoices, totalCount, err := h.repo.SearchInvoices(ctx, req.SearchText, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("SearchInvoices").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to search invoices", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("SearchInvoices", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("SearchInvoices", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to search invoices: %v", err)
	}

	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	span.SetAttributes(
		attribute.Int("total_count", int(totalCount)),
		attribute.Int("total_pages", int(totalPages)),
		attribute.Int("invoices_returned", len(invoices)),
	)

	h.metrics.RequestsTotal.WithLabelValues("SearchInvoices", "success").Inc()

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

	if err := h.cache.Set(ctx, cacheKey, response); err != nil {
		h.logger.Warn("Failed to cache response", zap.String("key", cacheKey), zap.Error(err))
	}

	return response, nil
}

// GetInvoiceByNumber returns a single invoice by number
func (h *InvoiceQueryHandler) GetInvoiceByNumber(ctx context.Context, req *pb.InvoiceNumberRequest) (*pb.APIResponse, error) {
	ctx, span := tracer.Start(ctx, "GetInvoiceByNumber",
		trace.WithAttributes(
			attribute.String("invoice_number", req.InvoiceNumber),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("GetInvoiceByNumber called", zap.String("invoice_number", req.InvoiceNumber))

	if req.InvoiceNumber == "" {
		h.metrics.ErrorsTotal.WithLabelValues("GetInvoiceByNumber", "validation_error").Inc()
		return nil, status.Error(codes.InvalidArgument, "invoice_number is required")
	}

	// Try cache
	cacheKey := cache.GenerateCacheKey("invoice-query", "GetInvoiceByNumber", req.InvoiceNumber)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("GetInvoiceByNumber", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))
	invoice, err := h.repo.GetInvoiceByNumber(ctx, req.InvoiceNumber)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("GetInvoiceByNumber").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to get invoice by number", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("GetInvoiceByNumber", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("GetInvoiceByNumber", "error").Inc()
		return nil, status.Errorf(codes.NotFound, "invoice not found: %v", err)
	}

	span.SetAttributes(attribute.Bool("invoice_found", invoice != nil))
	h.metrics.RequestsTotal.WithLabelValues("GetInvoiceByNumber", "success").Inc()

	// Single invoice response - wrap in list with single item
	response := &pb.APIResponse{
		TotalElements: wrapperspb.Int32(1),
		TotalPage:     wrapperspb.Int32(1),
		Data: &pb.APIResponse_InvoiceList{
			InvoiceList: &pb.InvoiceListResponse{
				Page:          0,
				Size:          1,
				TotalElements: 1,
				TotalPage:     1,
				Invoices:      []*entities.Invoice{invoice},
			},
		},
	}

	if err := h.cache.Set(ctx, cacheKey, response); err != nil {
		h.logger.Warn("Failed to cache response", zap.String("key", cacheKey), zap.Error(err))
	}

	return response, nil
}

// GetInvoicesByDateRange returns invoices within date range
func (h *InvoiceQueryHandler) GetInvoicesByDateRange(ctx context.Context, req *pb.DateRangeRequest) (*pb.APIResponse, error) {
	ctx, span := tracer.Start(ctx, "GetInvoicesByDateRange",
		trace.WithAttributes(
			attribute.String("date_from", req.DateFrom),
			attribute.String("date_to", req.DateTo),
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("GetInvoicesByDateRange called",
		zap.String("date_from", req.DateFrom),
		zap.String("date_to", req.DateTo),
		zap.Int32("page", req.Page),
		zap.Int32("size", req.Size),
	)

	if req.DateFrom == "" || req.DateTo == "" {
		h.metrics.ErrorsTotal.WithLabelValues("GetInvoicesByDateRange", "validation_error").Inc()
		return nil, status.Error(codes.InvalidArgument, "date_from and date_to are required")
	}

	// Try cache
	cacheKey := cache.GenerateCacheKey("invoice-query", "GetInvoicesByDateRange", req.DateFrom, req.DateTo, req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("GetInvoicesByDateRange", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	span.SetAttributes(attribute.Bool("cache_hit", false))
	invoices, totalCount, err := h.repo.GetInvoicesByDateRange(ctx, req.DateFrom, req.DateTo, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("GetInvoicesByDateRange").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to get invoices by date range", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("GetInvoicesByDateRange", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("GetInvoicesByDateRange", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to get invoices by date range: %v", err)
	}

	totalPages := int32(math.Ceil(float64(totalCount) / float64(req.Size)))

	span.SetAttributes(
		attribute.Int("total_count", int(totalCount)),
		attribute.Int("total_pages", int(totalPages)),
		attribute.Int("invoices_returned", len(invoices)),
	)

	h.metrics.RequestsTotal.WithLabelValues("GetInvoicesByDateRange", "success").Inc()

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

	if err := h.cache.Set(ctx, cacheKey, response); err != nil {
		h.logger.Warn("Failed to cache response", zap.String("key", cacheKey), zap.Error(err))
	}

	return response, nil
}
