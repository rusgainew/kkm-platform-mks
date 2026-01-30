// Файл api-gateway/internal/application/services/bank_account_query_service.go содержит реализацию пакета services.
package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
)

// BankAccountQueryService сервис для чтения банковских счетов (query side)
type BankAccountQueryService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewBankAccountQueryService создает новый BankAccountQueryService
func NewBankAccountQueryService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *BankAccountQueryService {
	return &BankAccountQueryService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для bank account query service
func (s *BankAccountQueryService) getClient(ctx context.Context) (pb.BankAccountQueryServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewBankAccountQueryServiceClient(conn), nil
}

// ListBankAccounts получает список всех банковских счетов с пагинацией
func (s *BankAccountQueryService) ListBankAccounts(ctx context.Context, pageNum, pageSize int32) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "BankAccountQueryService.ListBankAccounts")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get bank account query client", zap.Error(err))
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	req := &pb.PageInfo{
		Page: pageNum,
		Size: pageSize,
	}

	resp, err := client.ListBankAccounts(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list bank accounts", zap.Error(err))
		s.metrics.IncrementErrorCount("bank_account_query_service", "list_bank_accounts")
		return nil, fmt.Errorf("failed to list bank accounts: %w", err)
	}

	return resp, nil
}

// ListBankAccountsWithFilter получает список банковских счетов с фильтрацией
func (s *BankAccountQueryService) ListBankAccountsWithFilter(ctx context.Context, filter *pb.BankAccountFilterRequest) (*pb.APIResponse, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "BankAccountQueryService.ListBankAccountsWithFilter")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get bank account query client", zap.Error(err))
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	resp, err := client.ListBankAccountsWithFilter(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to list bank accounts with filter", zap.Error(err))
		s.metrics.IncrementErrorCount("bank_account_query_service", "list_bank_accounts_with_filter")
		return nil, fmt.Errorf("failed to list bank accounts with filter: %w", err)
	}

	return resp, nil
}
