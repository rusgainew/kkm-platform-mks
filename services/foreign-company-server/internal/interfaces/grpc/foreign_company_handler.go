package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/infrastructure/middleware"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	pbDict "github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ServiceInterface определяет методы необходимые для работы handler
type ServiceInterface interface {
	CreateForeignCompany(ctx context.Context, pin, fullName, countryCode string, createdBy uuid.UUID) (*domain.ForeignCompany, error)
	UpdateForeignCompany(ctx context.Context, id int64, pin, fullName, countryCode, address string, updatedBy uuid.UUID) (*domain.ForeignCompany, error)
	GetForeignCompany(ctx context.Context, id int64) (*domain.ForeignCompany, error)
	DeleteForeignCompany(ctx context.Context, id int64, deletedBy uuid.UUID) error
}

// ForeignCompanyHandler обработчик gRPC запросов для иностранных компаний
type ForeignCompanyHandler struct {
	pb.UnimplementedForeignCompanyCommandServiceServer
	service ServiceInterface
	logger  *zap.Logger
}

// NewForeignCompanyHandler создает новый обработчик
func NewForeignCompanyHandler(service ServiceInterface, logger *zap.Logger) *ForeignCompanyHandler {
	return &ForeignCompanyHandler{
		service: service,
		logger:  logger,
	}
}

// CreateForeignCompany создает новую иностранную компанию
func (h *ForeignCompanyHandler) CreateForeignCompany(ctx context.Context, req *pb.CreateForeignCompanyRequest) (*pbDict.ForeignCompany, error) {
	h.logger.Info("creating foreign company",
		zap.String("pin", req.Pin),
		zap.String("full_name", req.FullName),
		zap.String("country_code", req.CountryCode))

	// Извлечение user_id из контекста (установлен JWT middleware)
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		h.logger.Warn("failed to extract user_id from context", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Создание через сервис
	company, err := h.service.CreateForeignCompany(ctx, req.Pin, req.FullName, req.CountryCode, userID)
	if err != nil {
		h.logger.Error("failed to create foreign company", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pbDict.ForeignCompany{
		Id:       company.ID,
		Pin:      company.PIN,
		FullName: company.FullName,
	}, nil
}

// UpdateForeignCompany обновляет иностранную компанию
func (h *ForeignCompanyHandler) UpdateForeignCompany(ctx context.Context, req *pb.UpdateForeignCompanyRequest) (*pbDict.ForeignCompany, error) {
	h.logger.Info("updating foreign company",
		zap.Int64("id", req.Id),
		zap.String("country_code", req.CountryCode))

	// Извлечение user_id из контекста
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		h.logger.Warn("failed to extract user_id from context", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Обновление через сервис
	company, err := h.service.UpdateForeignCompany(ctx, req.Id, req.Pin, req.FullName, req.CountryCode, req.Address, userID)
	if err != nil {
		h.logger.Error("failed to update foreign company", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pbDict.ForeignCompany{
		Id:       company.ID,
		Pin:      company.PIN,
		FullName: company.FullName,
	}, nil
}

// GetForeignCompany возвращает иностранную компанию по ID
func (h *ForeignCompanyHandler) GetForeignCompany(ctx context.Context, req *pb.GetForeignCompanyRequest) (*pbDict.ForeignCompany, error) {
	h.logger.Info("getting foreign company", zap.Int64("id", req.Id))

	// Получение через сервис
	company, err := h.service.GetForeignCompany(ctx, req.Id)
	if err != nil {
		h.logger.Error("failed to get foreign company", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pbDict.ForeignCompany{
		Id:       company.ID,
		Pin:      company.PIN,
		FullName: company.FullName,
	}, nil
}

// DeleteForeignCompany удаляет иностранную компанию
func (h *ForeignCompanyHandler) DeleteForeignCompany(ctx context.Context, req *pb.DeleteForeignCompanyRequest) (*pb.DeleteForeignCompanyResponse, error) {
	h.logger.Info("deleting foreign company", zap.Int64("id", req.Id))

	// Извлечение user_id из контекста
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		h.logger.Warn("failed to extract user_id from context", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Удаление через сервис
	err = h.service.DeleteForeignCompany(ctx, req.Id, userID)
	if err != nil {
		h.logger.Error("failed to delete foreign company", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pb.DeleteForeignCompanyResponse{
		Success: true,
		Message: "foreign company deleted successfully",
	}, nil
}

// getUserIDFromContext извлекает user_id из контекста
func (h *ForeignCompanyHandler) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return middleware.GetUserIDFromContext(ctx)
}

// mapError преобразует domain ошибки в gRPC статусы
func (h *ForeignCompanyHandler) mapError(err error) error {
	// TODO: добавить маппинг специфичных ошибок
	return status.Error(codes.Internal, err.Error())
}
