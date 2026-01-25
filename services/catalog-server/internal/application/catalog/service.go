package catalog

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/domain/events"
	"github.com/rusgainew/kkm-project-mks/catalog-server/internal/domain/ports"
	"go.uber.org/zap"
)

// Service реализует бизнес-логику управления каталогом
type Service struct {
	repo      ports.CatalogRepository
	publisher ports.EventPublisher
	logger    *zap.Logger
}

// NewService создает новый сервис каталога
func NewService(repo ports.CatalogRepository, publisher ports.EventPublisher, logger *zap.Logger) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}

// CreateItem создает новый элемент каталога
func (s *Service) CreateItem(ctx context.Context, organizationID, createdBy uuid.UUID, name, code, description, unitType string, price, vatRate float64) (*domain.CatalogItem, error) {
	// Валидация
	if name == "" {
		return nil, domain.ErrInvalidName
	}
	if code == "" {
		return nil, domain.ErrInvalidCode
	}
	if price < 0 {
		return nil, domain.ErrInvalidPrice
	}
	// unitType теперь опционален для совместимости с API
	if unitType == "" {
		unitType = "шт" // значение по умолчанию
	}

	// Проверка существования
	exists, err := s.repo.ExistsByCode(ctx, organizationID, code)
	if err != nil {
		s.logger.Error("failed to check item existence", zap.Error(err))
		return nil, fmt.Errorf("failed to check item existence: %w", err)
	}
	if exists {
		return nil, domain.ErrCatalogItemAlreadyExists
	}

	// Создание
	item := domain.NewCatalogItem(organizationID, createdBy, name, code, unitType, price, vatRate)
	item.Description = description

	if err := s.repo.Create(ctx, item); err != nil {
		s.logger.Error("failed to create catalog item", zap.Error(err))
		return nil, fmt.Errorf("failed to create catalog item: %w", err)
	}

	// Публикация события
	event := events.NewCatalogItemCreatedEvent(item.ID, organizationID, createdBy, name, code, price)
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Warn("failed to publish catalog item created event", zap.Error(err))
	}

	s.logger.Info("catalog item created",
		zap.String("item_id", item.ID.String()),
		zap.String("organization_id", organizationID.String()),
	)

	return item, nil
}

// GetItem возвращает элемент каталога по ID
func (s *Service) GetItem(ctx context.Context, id uuid.UUID) (*domain.CatalogItem, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get catalog item", zap.Error(err), zap.String("item_id", id.String()))
		return nil, err
	}
	return item, nil
}

// UpdateItem обновляет элемент каталога
func (s *Service) UpdateItem(ctx context.Context, id, updatedBy uuid.UUID, name, description, unitType string, price, vatRate float64) (*domain.CatalogItem, error) {
	// Получение элемента
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Валидация
	if name == "" {
		return nil, domain.ErrInvalidName
	}
	if price < 0 {
		return nil, domain.ErrInvalidPrice
	}
	// unitType теперь опционален - сохраняем существующий если не указан
	if unitType == "" {
		unitType = item.UnitType
	}

	// Обновление
	item.Update(name, description, unitType, price, vatRate, updatedBy)

	if err := s.repo.Update(ctx, item); err != nil {
		s.logger.Error("failed to update catalog item", zap.Error(err))
		return nil, fmt.Errorf("failed to update catalog item: %w", err)
	}

	// Публикация события
	event := events.NewCatalogItemUpdatedEvent(item.ID, item.OrganizationID, updatedBy, name, price)
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Warn("failed to publish catalog item updated event", zap.Error(err))
	}

	s.logger.Info("catalog item updated", zap.String("item_id", item.ID.String()))

	return item, nil
}

// DeleteItem удаляет элемент каталога
func (s *Service) DeleteItem(ctx context.Context, id, deletedBy uuid.UUID) error {
	// Получение элемента для проверки существования
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete catalog item", zap.Error(err))
		return fmt.Errorf("failed to delete catalog item: %w", err)
	}

	// Публикация события
	event := events.NewCatalogItemDeletedEvent(id, item.OrganizationID, deletedBy)
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Warn("failed to publish catalog item deleted event", zap.Error(err))
	}

	s.logger.Info("catalog item deleted", zap.String("item_id", id.String()))

	return nil
}

// ListItems возвращает список элементов организации
func (s *Service) ListItems(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*domain.CatalogItem, int, error) {
	items, total, err := s.repo.ListByOrganization(ctx, organizationID, limit, offset)
	if err != nil {
		s.logger.Error("failed to list catalog items", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list catalog items: %w", err)
	}
	return items, total, nil
}

// SearchItems ищет элементы по запросу
func (s *Service) SearchItems(ctx context.Context, organizationID uuid.UUID, query string, limit, offset int) ([]*domain.CatalogItem, int, error) {
	items, total, err := s.repo.Search(ctx, organizationID, query, limit, offset)
	if err != nil {
		s.logger.Error("failed to search catalog items", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to search catalog items: %w", err)
	}
	return items, total, nil
}
