// Файл bank-account-query-server/internal/interfaces/grpc/handlers/bank_account_query_handler_filters.go содержит реализацию пакета handlers.
package handlers

import (
	"context"
	"math"
	"time"

	"github.com/rusgainew/kkm-project-mks/bank-account-query-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/bank-account-query-server/internal/infrastructure/cache"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ListBankAccountsWithFilter returns bank accounts with filtering and sorting
func (h *BankAccountQueryHandler) ListBankAccountsWithFilter(ctx context.Context, req *pb.BankAccountFilterRequest) (*pb.APIResponse, error) {
	ctx, span := tracer.Start(ctx, "ListBankAccountsWithFilter",
		trace.WithAttributes(
			attribute.String("account_name", req.AccountName),
			attribute.String("bank_account", req.BankAccount),
			attribute.Bool("is_active", req.IsActive),
			attribute.String("search_text", req.SearchText),
			attribute.String("sort_field", req.SortField),
			attribute.String("sort_order", req.SortOrder),
			attribute.Int("page", int(req.Page)),
			attribute.Int("size", int(req.Size)),
		),
	)
	defer span.End()

	start := time.Now()
	h.logger.Info("ListBankAccountsWithFilter called",
		zap.String("account_name", req.AccountName),
		zap.String("bank_account", req.BankAccount),
		zap.Bool("is_active", req.IsActive),
		zap.String("search_text", req.SearchText),
		zap.String("sort_field", req.SortField),
		zap.String("sort_order", req.SortOrder),
		zap.Int32("page", req.Page),
		zap.Int32("size", req.Size),
	)

	// Validate sort field
	validSortFields := map[string]bool{
		"account_name": true,
		"bank_account": true,
		"created_at":   true,
		"":             true, // empty means default sort
	}
	if !validSortFields[req.SortField] {
		h.logger.Warn("Invalid sort field", zap.String("sort_field", req.SortField))
		h.metrics.ErrorsTotal.WithLabelValues("ListBankAccountsWithFilter", "validation_error").Inc()
		return nil, status.Errorf(codes.InvalidArgument, "invalid sort field: %s", req.SortField)
	}

	// Validate sort order
	if req.SortOrder != "" && req.SortOrder != "ASC" && req.SortOrder != "DESC" {
		h.logger.Warn("Invalid sort order", zap.String("sort_order", req.SortOrder))
		h.metrics.ErrorsTotal.WithLabelValues("ListBankAccountsWithFilter", "validation_error").Inc()
		return nil, status.Errorf(codes.InvalidArgument, "invalid sort order: %s (must be ASC or DESC)", req.SortOrder)
	}

	// Try to get from cache
	cacheKey := cache.GenerateCacheKey("bank-account-query", "ListBankAccountsWithFilter",
		req.AccountName, req.BankAccount, req.IsActive, req.SearchText,
		req.SortField, req.SortOrder, req.Page, req.Size)
	var cachedResponse pb.APIResponse
	if err := h.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		h.logger.Debug("Cache hit", zap.String("key", cacheKey))
		span.SetAttributes(attribute.Bool("cache_hit", true))
		h.metrics.RequestsTotal.WithLabelValues("ListBankAccountsWithFilter", "cache_hit").Inc()
		return &cachedResponse, nil
	}

	// Cache miss - prepare filter and sort
	span.SetAttributes(attribute.Bool("cache_hit", false))

	// Convert IsActive to pointer
	var isActivePtr *bool
	if req.IsActive {
		isActivePtr = &req.IsActive
	}

	filter := &ports.BankAccountFilter{
		AccountName: req.AccountName,
		BankAccount: req.BankAccount,
		IsActive:    isActivePtr,
		SearchText:  req.SearchText,
	}

	sortField := req.SortField
	if sortField == "" {
		sortField = "created_at" // default sort
	}
	sortOrder := ports.SortOrder(req.SortOrder)
	if req.SortOrder == "" {
		sortOrder = ports.SortOrderDesc // default order
	}
	sort := &ports.BankAccountSort{
		Field: sortField,
		Order: sortOrder,
	}

	// Query database
	accounts, totalCount, err := h.repo.ListBankAccountsWithFilter(ctx, filter, sort, req.Page, req.Size)

	duration := time.Since(start).Milliseconds()
	h.metrics.RequestDuration.WithLabelValues("ListBankAccountsWithFilter").Observe(float64(duration))

	if err != nil {
		span.RecordError(err)
		h.logger.Error("Failed to list bank accounts with filter", zap.Error(err))
		h.metrics.ErrorsTotal.WithLabelValues("ListBankAccountsWithFilter", "database_error").Inc()
		h.metrics.RequestsTotal.WithLabelValues("ListBankAccountsWithFilter", "error").Inc()
		return nil, status.Errorf(codes.Internal, "failed to list bank accounts: %v", err)
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
	h.metrics.RequestsTotal.WithLabelValues("ListBankAccountsWithFilter", "success").Inc()

	// Build response
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
