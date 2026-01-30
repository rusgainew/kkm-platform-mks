// Файл catalog-server/internal/interfaces/grpc/catalog_handler.go содержит реализацию пакета grpc.
package grpc

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/application/catalog"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/infrastructure/middleware"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	pbDict "github.com/rusgainew/kkm-project-mks/proto-lib/dictionaries"
)

// CatalogHandler обработчик gRPC запросов для каталога
type CatalogHandler struct {
	pb.UnimplementedCatalogCommandServiceServer
	service *catalog.Service
	logger  *zap.Logger
}

// NewCatalogHandler создает новый обработчик
func NewCatalogHandler(service *catalog.Service, logger *zap.Logger) *CatalogHandler {
	return &CatalogHandler{
		service: service,
		logger:  logger,
	}
}

// CreateCatalog создает новую запись в каталоге
func (h *CatalogHandler) CreateCatalog(ctx context.Context, req *pb.CreateCatalogRequest) (*pbDict.Catalog, error) {
	h.logger.Info("creating catalog item",
		zap.String("name", req.Name),
		zap.String("number", req.Number))

	// Извлечение user_id из контекста (JWT middleware)
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		h.logger.Warn("failed to extract user_id from context", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "user authentication required")
	}

	// Для упрощения используем userID как organizationID
	organizationID := userID // TODO: получить настоящий organization_id

	// Получаем unitType безопасно
	var unitType string
	if req.UnitClassification != nil {
		unitType = req.UnitClassification.Name
	}

	// Создание через service layer
	catalog, err := h.service.CreateItem(
		ctx,
		organizationID,
		userID,
		req.Name,
		req.Number, // code
		"",         // description
		unitType,   // unitType
		0.0,        // price
		0.0,        // vatRate
	)
	if err != nil {
		h.logger.Error("failed to create catalog", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pbDict.Catalog{
		Id:                 catalog.ID.String(),
		Name:               catalog.Name,
		Number:             catalog.Code,
		TnvedCode:          req.TnvedCode,
		GkedCode:           req.GkedCode,
		UnitClassification: req.UnitClassification,
		CatalogType:        req.CatalogType,
	}, nil
}

// UpdateCatalog обновляет запись в каталоге
func (h *CatalogHandler) UpdateCatalog(ctx context.Context, req *pb.UpdateCatalogRequest) (*pbDict.Catalog, error) {
	h.logger.Info("updating catalog item", zap.String("id", req.Id))

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

	// Получаем unitType безопасно
	var unitType string
	if req.UnitClassification != nil {
		unitType = req.UnitClassification.Name
	}

	// Обновление через service layer
	catalog, err := h.service.UpdateItem(
		ctx,
		id,
		userID,
		req.Name,
		"",       // description
		unitType, // unitType
		0.0,      // price
		0.0,      // vatRate
	)
	if err != nil {
		h.logger.Error("failed to update catalog", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pbDict.Catalog{
		Id:                 catalog.ID.String(),
		Name:               catalog.Name,
		Number:             catalog.Code,
		TnvedCode:          req.TnvedCode,
		GkedCode:           req.GkedCode,
		UnitClassification: req.UnitClassification,
		CatalogType:        req.CatalogType,
	}, nil
}

// GetCatalog возвращает запись из каталога по ID
func (h *CatalogHandler) GetCatalog(ctx context.Context, req *pb.GetCatalogRequest) (*pbDict.Catalog, error) {
	h.logger.Info("getting catalog item", zap.String("id", req.Id))

	// Парсинг UUID
	id, err := uuid.Parse(req.Id)
	if err != nil {
		h.logger.Warn("invalid UUID format", zap.String("id", req.Id), zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "invalid id format")
	}

	// Получение через service layer
	catalog, err := h.service.GetItem(ctx, id)
	if err != nil {
		h.logger.Error("failed to get catalog", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pbDict.Catalog{
		Id:     catalog.ID.String(),
		Name:   catalog.Name,
		Number: catalog.Code,
	}, nil
}

// DeleteCatalog удаляет запись из каталога
func (h *CatalogHandler) DeleteCatalog(ctx context.Context, req *pb.DeleteCatalogRequest) (*pb.DeleteCatalogResponse, error) {
	h.logger.Info("deleting catalog item", zap.String("id", req.Id))

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
	err = h.service.DeleteItem(ctx, id, userID)
	if err != nil {
		h.logger.Error("failed to delete catalog", zap.Error(err))
		return nil, h.mapError(err)
	}

	return &pb.DeleteCatalogResponse{
		Success: true,
		Message: "catalog item deleted successfully",
	}, nil
}

// getUserIDFromContext извлекает user_id из контекста
func (h *CatalogHandler) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	return middleware.GetUserIDFromContext(ctx)
}

// mapError маппит domain ошибки в gRPC status codes
func (h *CatalogHandler) mapError(err error) error {
	switch err {
	case domain.ErrCatalogItemNotFound:
		return status.Error(codes.NotFound, "catalog item not found")
	case domain.ErrCatalogItemAlreadyExists:
		return status.Error(codes.AlreadyExists, "catalog item already exists")
	case domain.ErrInvalidName, domain.ErrInvalidCode:
		return status.Error(codes.InvalidArgument, "invalid catalog data")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
