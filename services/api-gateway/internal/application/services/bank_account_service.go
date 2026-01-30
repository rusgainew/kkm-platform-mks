// Файл api-gateway/internal/application/services/bank_account_service.go содержит реализацию пакета services.
package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	pbDict "github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
	"go.uber.org/zap"
)

// BankAccountService сервис для работы с банковскими счетами
type BankAccountService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewBankAccountService создает новый BankAccountService
func NewBankAccountService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *BankAccountService {
	return &BankAccountService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для bank account service
func (s *BankAccountService) getClient(ctx context.Context) (pb.BankAccountCommandServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewBankAccountCommandServiceClient(conn), nil
}

// CreateBankAccount создает новый банковский счет
func (s *BankAccountService) CreateBankAccount(ctx context.Context, account *models.BankAccount) (*models.BankAccount, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "BankAccountService.CreateBankAccount")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get bank account client", zap.Error(err))
		return nil, err
	}

	req := &pb.CreateBankAccountRequest{
		AccountName:   account.AccountNumber,
		BankAccount:   account.AccountNumber,
		ContractorTin: account.OwnerID,
		IsResident:    account.IsActive,
		Currency: &pbDict.ReferenceItem{
			Name: account.Currency,
			Code: account.Currency,
		},
		Bank: &pbDict.Bank{
			Name: account.BankName,
			Bik:  account.BankCode,
		},
	}

	resp, err := client.CreateBankAccount(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create bank account", zap.Error(err))
		s.metrics.IncrementErrorCount("bank_account_service", "create_account")
		return nil, fmt.Errorf("failed to create bank account: %w", err)
	}

	s.logger.Info("Bank account created successfully", zap.String("id", resp.Id))

	return &models.BankAccount{
		ID:            resp.Id,
		AccountNumber: resp.BankAccount,
		BankName:      resp.Bank.Name,
		BankCode:      resp.Bank.Bik,
		Currency:      resp.Currency.Code,
		OwnerID:       resp.ContractorTin,
		IsActive:      resp.IsActive,
		CreatedAt:     0,
		UpdatedAt:     0,
	}, nil
}

// UpdateBankAccount обновляет банковский счет
func (s *BankAccountService) UpdateBankAccount(ctx context.Context, account *models.BankAccount) (*models.BankAccount, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "BankAccountService.UpdateBankAccount")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get bank account client", zap.Error(err))
		return nil, err
	}

	req := &pb.UpdateBankAccountRequest{
		Id:            account.ID,
		AccountName:   account.AccountNumber,
		BankAccount:   account.AccountNumber,
		ContractorTin: account.OwnerID,
		IsResident:    account.IsActive,
	}

	resp, err := client.UpdateBankAccount(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update bank account", zap.String("id", account.ID), zap.Error(err))
		s.metrics.IncrementErrorCount("bank_account_service", "update_account")
		return nil, fmt.Errorf("failed to update bank account: %w", err)
	}

	s.logger.Info("Bank account updated successfully", zap.String("id", resp.Id))

	return &models.BankAccount{
		ID:            resp.Id,
		AccountNumber: resp.BankAccount,
		BankName:      resp.Bank.Name,
		BankCode:      resp.Bank.Bik,
		Currency:      resp.Currency.Code,
		OwnerID:       resp.ContractorTin,
		IsActive:      resp.IsActive,
		CreatedAt:     0,
		UpdatedAt:     0,
	}, nil
}

// GetBankAccount возвращает банковский счет по ID
func (s *BankAccountService) GetBankAccount(ctx context.Context, id string) (*models.BankAccount, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "BankAccountService.GetBankAccount")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get bank account client", zap.Error(err))
		return nil, err
	}

	req := &pb.GetBankAccountRequest{
		Id: id,
	}

	resp, err := client.GetBankAccount(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get bank account", zap.String("id", id), zap.Error(err))
		s.metrics.IncrementErrorCount("bank_account_service", "get_account")
		return nil, fmt.Errorf("failed to get bank account: %w", err)
	}

	s.logger.Info("Bank account retrieved successfully", zap.String("id", resp.Id))

	return &models.BankAccount{
		ID:            resp.Id,
		AccountNumber: resp.BankAccount,
		BankName:      resp.Bank.GetName(),
		BankCode:      resp.Bank.GetBik(),
		Currency:      resp.Currency.GetCode(),
		OwnerID:       resp.ContractorTin,
		IsActive:      resp.IsActive,
		CreatedAt:     0,
		UpdatedAt:     0,
	}, nil
}

// DeleteBankAccount удаляет банковский счет
func (s *BankAccountService) DeleteBankAccount(ctx context.Context, id string) error {
	ctx, endSpan := s.tracer.StartSpan(ctx, "BankAccountService.DeleteBankAccount")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get bank account client", zap.Error(err))
		return err
	}

	req := &pb.DeleteBankAccountRequest{
		Id: id,
	}

	_, err = client.DeleteBankAccount(ctx, req)
	if err != nil {
		s.logger.Error("Failed to delete bank account", zap.String("id", id), zap.Error(err))
		s.metrics.IncrementErrorCount("bank_account_service", "delete_account")
		return fmt.Errorf("failed to delete bank account: %w", err)
	}

	s.logger.Info("Bank account deleted successfully", zap.String("id", id))

	return nil
}
