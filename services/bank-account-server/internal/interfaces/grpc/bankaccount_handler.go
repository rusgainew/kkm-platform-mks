// Файл bank-account-server/internal/interfaces/grpc/bankaccount_handler.go содержит реализацию пакета grpc.
package grpc

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rusgainew/kkm-project-mks/bank-account-server/internal/application/bankaccount"
	"github.com/rusgainew/kkm-project-mks/bank-account-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/bank-account-server/internal/infrastructure/middleware"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	pbDict "github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
)

// BankAccountHandler обработчик gRPC запросов для банковских счетов
type BankAccountHandler struct {
	pb.UnimplementedBankAccountCommandServiceServer
	service *bankaccount.Service
	logger  *zap.Logger
}

// NewBankAccountHandler создает новый обработчик
func NewBankAccountHandler(service *bankaccount.Service, logger *zap.Logger) *BankAccountHandler {
	return &BankAccountHandler{
		service: service,
		logger:  logger,
	}
}

// CreateBankAccount создает новый банковский счет
func (h *BankAccountHandler) CreateBankAccount(ctx context.Context, req *pb.CreateBankAccountRequest) (*pbDict.BankAccount, error) {
	h.logger.Info("creating bank account",
		zap.String("account_name", req.AccountName),
		zap.String("bank_account", req.BankAccount))

	// Извлечение user_id из контекста (JWT middleware)
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		h.logger.Warn("failed to extract user_id from context", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Извлекаем organization_id и bank_id из request (если есть в контексте или параметрах)
	// Для упрощения используем userID, но в реальности нужен organization_id из JWT или параметра
	organizationID := userID // TODO: получить настоящий organization_id из контекста/параметра
	bankID := uuid.New()     // TODO: получить bank_id из параметра req.Bank

	// Безопасное извлечение значений из необязательных полей
	var currencyName, bankBik, bankName string
	if req.Currency != nil {
		currencyName = req.Currency.Name
	}
	if req.Bank != nil {
		bankBik = req.Bank.Bik
		bankName = req.Bank.Name
	}

	// Создание через service layer
	bankAccount, err := h.service.CreateAccount(
		ctx,
		organizationID,
		bankID,
		userID,
		req.BankAccount,
		"",           // IBAN - если есть в запросе
		currencyName, // Используем имя из ReferenceItem
		bankBik,
		bankName,
	)
	if err != nil {
		h.logger.Error("failed to create bank account", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pbDict.BankAccount{
		Id:            bankAccount.ID.String(),
		AccountName:   req.AccountName,
		BankAccount:   bankAccount.AccountNumber,
		ContractorTin: req.ContractorTin,
		IsActive:      bankAccount.IsActive,
		Currency:      req.Currency,
		Bank:          req.Bank,
		IsResident:    req.IsResident,
	}, nil
}

// UpdateBankAccount обновляет банковский счет
func (h *BankAccountHandler) UpdateBankAccount(ctx context.Context, req *pb.UpdateBankAccountRequest) (*pbDict.BankAccount, error) {
	h.logger.Info("updating bank account", zap.String("id", req.Id))

	// Извлечение user_id из контекста
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		h.logger.Warn("failed to extract user_id from context", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Парсинг UUID
	id, err := uuid.Parse(req.Id)
	if err != nil {
		h.logger.Warn("invalid UUID format", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	// Безопасное извлечение значений из необязательных полей
	var currencyName, bankBik, bankName string
	if req.Currency != nil {
		currencyName = req.Currency.Name
	}
	if req.Bank != nil {
		bankBik = req.Bank.Bik
		bankName = req.Bank.Name
	}

	// Обновление через service layer
	bankAccount, err := h.service.UpdateAccount(
		ctx,
		id,
		userID,
		req.BankAccount,
		"", // IBAN
		currencyName,
		bankBik,
		bankName,
	)
	if err != nil {
		h.logger.Error("failed to update bank account", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pbDict.BankAccount{
		Id:            bankAccount.ID.String(),
		AccountName:   req.AccountName,
		BankAccount:   bankAccount.AccountNumber,
		ContractorTin: req.ContractorTin,
		IsActive:      bankAccount.IsActive,
		Currency:      req.Currency,
		Bank:          req.Bank,
		IsResident:    req.IsResident,
	}, nil
}

// GetBankAccount возвращает банковский счет по ID
func (h *BankAccountHandler) GetBankAccount(ctx context.Context, req *pb.GetBankAccountRequest) (*pbDict.BankAccount, error) {
	h.logger.Info("getting bank account", zap.String("id", req.Id))

	// Парсинг UUID
	id, err := uuid.Parse(req.Id)
	if err != nil {
		h.logger.Warn("invalid UUID format", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	// Получение через service layer
	bankAccount, err := h.service.GetAccount(ctx, id)
	if err != nil {
		h.logger.Error("failed to get bank account", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pbDict.BankAccount{
		Id:          bankAccount.ID.String(),
		BankAccount: bankAccount.AccountNumber,
		IsActive:    bankAccount.IsActive,
	}, nil
}

// DeleteBankAccount удаляет банковский счет
func (h *BankAccountHandler) DeleteBankAccount(ctx context.Context, req *pb.DeleteBankAccountRequest) (*pb.DeleteBankAccountResponse, error) {
	h.logger.Info("deleting bank account", zap.String("id", req.Id))

	// Извлечение user_id из контекста
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		h.logger.Warn("failed to extract user_id from context", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Парсинг UUID
	id, err := uuid.Parse(req.Id)
	if err != nil {
		h.logger.Warn("invalid UUID format", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	// Удаление через service layer
	err = h.service.DeleteAccount(ctx, id, userID)
	if err != nil {
		h.logger.Error("failed to delete bank account", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pb.DeleteBankAccountResponse{
		Success: true,
		Message: "bank account deleted successfully",
	}, nil
}

// getUserIDFromContext извлекает user_id из контекста
func (h *BankAccountHandler) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return middleware.GetUserIDFromContext(ctx)
}

// mapError маппит domain ошибки в gRPC status codes
func (h *BankAccountHandler) mapError(err error) error {
	switch err {
	case domain.ErrBankAccountNotFound:
		return status.Error(codes.NotFound, "bank account not found")
	case domain.ErrBankAccountAlreadyExists:
		return status.Error(codes.AlreadyExists, "bank account already exists")
	case domain.ErrInvalidBankAccount, domain.ErrInvalidAccountNumber:
		return status.Error(codes.InvalidArgument, "invalid bank account data")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
