package foreigncompany

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/domain/events"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/domain/ports"
)

// Service реализует бизнес-логику управления иностранными компаниями
type Service struct {
	repo      ports.ForeignCompanyRepository
	publisher ports.EventPublisher
	logger    *zap.Logger
}

// NewService создает новый экземпляр сервиса
func NewService(repo ports.ForeignCompanyRepository, publisher ports.EventPublisher, logger *zap.Logger) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}

// CreateForeignCompany создает новую иностранную компанию
func (s *Service) CreateForeignCompany(ctx context.Context, pin, fullName, countryCode string, createdBy uuid.UUID) (*domain.ForeignCompany, error) {
	// Валидация входных данных
	if pin == "" {
		return nil, domain.ErrInvalidPIN
	}
	if fullName == "" {
		return nil, domain.ErrInvalidFullName
	}
	if createdBy == uuid.Nil {
		return nil, domain.ErrInvalidCreatedBy
	}

	// Нормализация данных
	pin = strings.TrimSpace(strings.ToUpper(pin))
	fullName = strings.TrimSpace(fullName)
	countryCode = strings.TrimSpace(strings.ToUpper(countryCode))

	// Проверка существования компании с таким PIN
	exists, err := s.repo.ExistsByPIN(ctx, pin)
	if err != nil {
		s.logger.Error("failed to check foreign company existence",
			zap.String("pin", pin),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to check existence: %w", err)
	}
	if exists {
		return nil, domain.ErrForeignCompanyAlreadyExists
	}

	// Создание сущности
	company := domain.NewForeignCompany(pin, fullName, countryCode, createdBy)

	// Валидация доменной модели
	if err := company.Validate(); err != nil {
		return nil, err
	}

	// Сохранение в репозиторий
	if err := s.repo.Create(ctx, company); err != nil {
		s.logger.Error("failed to create foreign company",
			zap.String("pin", pin),
			zap.String("full_name", fullName),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create foreign company: %w", err)
	}

	// Публикация события
	event := &events.ForeignCompanyCreatedEvent{
		ID:          company.ID,
		PIN:         company.PIN,
		FullName:    company.FullName,
		CountryCode: company.CountryCode,
		CreatedBy:   company.CreatedBy,
		CreatedAt:   company.CreatedAt,
	}
	if s.publisher != nil {
		if err := s.publisher.Publish(ctx, event); err != nil {
			s.logger.Warn("failed to publish foreign company created event",
				zap.Int64("id", company.ID),
				zap.Error(err),
			)
		}
	}

	s.logger.Info("foreign company created",
		zap.Int64("id", company.ID),
		zap.String("pin", company.PIN),
		zap.String("full_name", company.FullName),
	)

	return company, nil
}

// GetForeignCompany возвращает компанию по ID
func (s *Service) GetForeignCompany(ctx context.Context, id int64) (*domain.ForeignCompany, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidID
	}

	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get foreign company",
			zap.Int64("id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return company, nil
}

// GetForeignCompanyByPIN возвращает компанию по PIN
func (s *Service) GetForeignCompanyByPIN(ctx context.Context, pin string) (*domain.ForeignCompany, error) {
	if pin == "" {
		return nil, domain.ErrInvalidPIN
	}

	pin = strings.TrimSpace(strings.ToUpper(pin))

	company, err := s.repo.GetByPIN(ctx, pin)
	if err != nil {
		s.logger.Error("failed to get foreign company by PIN",
			zap.String("pin", pin),
			zap.Error(err),
		)
		return nil, err
	}

	return company, nil
}

// UpdateForeignCompany обновляет данные компании
func (s *Service) UpdateForeignCompany(ctx context.Context, id int64, pin, fullName, countryCode, address string, updatedBy uuid.UUID) (*domain.ForeignCompany, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidID
	}
	if updatedBy == uuid.Nil {
		return nil, domain.ErrInvalidCreatedBy
	}

	// Получение существующей компании
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Нормализация данных
	pin = strings.TrimSpace(strings.ToUpper(pin))
	fullName = strings.TrimSpace(fullName)
	countryCode = strings.TrimSpace(strings.ToUpper(countryCode))
	address = strings.TrimSpace(address)

	// Если PIN изменился, проверяем уникальность
	if pin != company.PIN {
		exists, err := s.repo.ExistsByPIN(ctx, pin)
		if err != nil {
			return nil, fmt.Errorf("failed to check PIN existence: %w", err)
		}
		if exists {
			return nil, domain.ErrForeignCompanyAlreadyExists
		}
	}

	// Обновление данных
	company.Update(pin, fullName, countryCode, address, updatedBy)

	// Валидация
	if err := company.Validate(); err != nil {
		return nil, err
	}

	// Сохранение
	if err := s.repo.Update(ctx, company); err != nil {
		s.logger.Error("failed to update foreign company",
			zap.Int64("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update foreign company: %w", err)
	}

	// Публикация события
	event := &events.ForeignCompanyUpdatedEvent{
		ID:          company.ID,
		PIN:         company.PIN,
		FullName:    company.FullName,
		CountryCode: company.CountryCode,
		Address:     company.Address,
		UpdatedBy:   updatedBy,
		UpdatedAt:   company.UpdatedAt,
	}
	if s.publisher != nil {
		if err := s.publisher.Publish(ctx, event); err != nil {
			s.logger.Warn("failed to publish foreign company updated event",
				zap.Int64("id", company.ID),
				zap.Error(err),
			)
		}
	}

	s.logger.Info("foreign company updated",
		zap.Int64("id", company.ID),
		zap.String("pin", company.PIN),
	)

	return company, nil
}

// DeleteForeignCompany удаляет компанию
func (s *Service) DeleteForeignCompany(ctx context.Context, id int64, deletedBy uuid.UUID) error {
	if id <= 0 {
		return domain.ErrInvalidID
	}

	// Проверка существования
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Проверка, что компания не активна
	if company.IsActive {
		return domain.ErrCannotDeleteActiveCompany
	}

	// Удаление
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete foreign company",
			zap.Int64("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete foreign company: %w", err)
	}

	// Публикация события
	event := &events.ForeignCompanyDeletedEvent{
		ID:        id,
		DeletedBy: deletedBy,
		DeletedAt: time.Now(),
	}
	if s.publisher != nil {
		if err := s.publisher.Publish(ctx, event); err != nil {
			s.logger.Warn("failed to publish foreign company deleted event",
				zap.Int64("id", id),
				zap.Error(err),
			)
		}
	}

	s.logger.Info("foreign company deleted",
		zap.Int64("id", id),
	)

	return nil
}

// ListForeignCompanies возвращает список компаний
func (s *Service) ListForeignCompanies(ctx context.Context, limit, offset int) ([]*domain.ForeignCompany, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	companies, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("failed to list foreign companies",
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.Error(err),
		)
		return nil, 0, err
	}

	return companies, total, nil
}

// SearchForeignCompanies ищет компании по названию
func (s *Service) SearchForeignCompanies(ctx context.Context, query string, limit, offset int) ([]*domain.ForeignCompany, int, error) {
	if query == "" {
		return s.ListForeignCompanies(ctx, limit, offset)
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	query = strings.TrimSpace(query)

	companies, total, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		s.logger.Error("failed to search foreign companies",
			zap.String("query", query),
			zap.Error(err),
		)
		return nil, 0, err
	}

	return companies, total, nil
}

// DeactivateForeignCompany деактивирует компанию
func (s *Service) DeactivateForeignCompany(ctx context.Context, id int64, deactivatedBy uuid.UUID) error {
	if id <= 0 {
		return domain.ErrInvalidID
	}

	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	company.Deactivate(deactivatedBy)

	if err := s.repo.Update(ctx, company); err != nil {
		s.logger.Error("failed to deactivate foreign company",
			zap.Int64("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to deactivate: %w", err)
	}

	// Публикация события
	event := &events.ForeignCompanyDeactivatedEvent{
		ID:            id,
		DeactivatedBy: deactivatedBy,
		DeactivatedAt: time.Now(),
	}
	if s.publisher != nil {
		if err := s.publisher.Publish(ctx, event); err != nil {
			s.logger.Warn("failed to publish deactivation event",
				zap.Int64("id", id),
				zap.Error(err),
			)
		}
	}

	s.logger.Info("foreign company deactivated", zap.Int64("id", id))
	return nil
}
