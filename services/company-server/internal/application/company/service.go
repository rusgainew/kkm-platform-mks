package company

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/domain/events"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/domain/ports"
	"github.com/rusgainew/kkm-project-mks/company-server/internal/infrastructure/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// Service реализует бизнес-логику управления организациями
type Service struct {
	orgRepo        ports.OrganizationRepository
	empRepo        ports.EmployeeRepository
	eventPublisher ports.EventPublisher
	logger         *zap.Logger
}

// NewService создает новый экземпляр сервиса
func NewService(
	orgRepo ports.OrganizationRepository,
	empRepo ports.EmployeeRepository,
	eventPublisher ports.EventPublisher,
	logger *zap.Logger,
) *Service {
	return &Service{
		orgRepo:        orgRepo,
		empRepo:        empRepo,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// CreateOrganization создает новую организацию
func (s *Service) CreateOrganization(ctx context.Context, name, description, ownerID string) (*domain.Organization, error) {
	// Начало отслеживания времени для метрик
	startTime := time.Now()
	status := "success"
	defer func() {
		duration := time.Since(startTime)
		observability.RecordOperationDuration("CreateOrganization", duration, status)
		observability.RecordOrganizationOperation("CreateOrganization", status)
	}()

	// Создание span для распределенной трассировки
	tracer := otel.Tracer("company-service")
	_, span := tracer.Start(ctx, "CreateOrganization")
	defer span.End()

	// Добавление атрибутов в span
	span.SetAttributes(
		attribute.String("org.name", name),
		attribute.String("org.owner_id", ownerID),
		attribute.Int("org.description_length", len(description)),
	)

	s.logger.Info("Creating organization", zap.String("name", name), zap.String("owner_id", ownerID))

	// Проверка существования организации с таким именем
	exists, err := s.orgRepo.ExistsByName(ctx, name)
	if err != nil {
		status = "error"
		s.logger.Error("Failed to check organization existence", zap.Error(err))
		return nil, err
	}
	if exists {
		status = "error"
		s.logger.Warn("Organization already exists", zap.String("name", name))
		return nil, domain.ErrOrganizationExists
	}

	// Создание организации
	var descPtr *string
	if description != "" {
		descPtr = &description
	}
	org := &domain.Organization{
		ID:          uuid.New().String(),
		Name:        name,
		Description: descPtr,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := org.Validate(); err != nil {
		status = "error"
		s.logger.Error("Organization validation failed", zap.Error(err))
		return nil, err
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		status = "error"
		s.logger.Error("Failed to create organization", zap.Error(err))
		return nil, err
	}

	// Публикация события
	eventDesc := ""
	if org.Description != nil {
		eventDesc = *org.Description
	}
	event := events.OrganizationCreated{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: events.OrganizationCreatedEvent,
			Timestamp: time.Now(),
		},
		OrganizationID: org.ID,
		Name:           org.Name,
		Description:    eventDesc,
		OwnerID:        org.OwnerID,
	}

	if eventData, err := event.BaseEvent.ToJSON(); err == nil {
		observability.RecordEventPublished(events.OrganizationCreatedEvent, "success")
		_ = s.eventPublisher.Publish(ctx, events.OrganizationCreatedEvent, eventData)
	}

	s.logger.Info("Organization created successfully", zap.String("id", org.ID))
	return org, nil
}

// GetOrganization получает организацию по ID
func (s *Service) GetOrganization(ctx context.Context, id string) (*domain.Organization, error) {
	// Начало отслеживания времени для метрик
	startTime := time.Now()
	status := "success"
	defer func() {
		duration := time.Since(startTime)
		observability.RecordOperationDuration("GetOrganization", duration, status)
		observability.RecordOrganizationOperation("GetOrganization", status)
	}()

	// Создание span для распределенной трассировки
	tracer := otel.Tracer("company-service")
	_, span := tracer.Start(ctx, "GetOrganization")
	defer span.End()

	span.SetAttributes(attribute.String("org.id", id))

	s.logger.Info("Getting organization", zap.String("id", id))

	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		status = "error"
		s.logger.Error("Failed to get organization", zap.Error(err))
		return nil, err
	}

	return org, nil
}

// UpdateOrganization обновляет организацию
func (s *Service) UpdateOrganization(ctx context.Context, id, name, description string) (*domain.Organization, error) {
	// Начало отслеживания времени для метрик
	startTime := time.Now()
	status := "success"
	defer func() {
		duration := time.Since(startTime)
		observability.RecordOperationDuration("UpdateOrganization", duration, status)
		observability.RecordOrganizationOperation("UpdateOrganization", status)
	}()

	// Создание span для распределенной трассировки
	tracer := otel.Tracer("company-service")
	_, span := tracer.Start(ctx, "UpdateOrganization")
	defer span.End()

	span.SetAttributes(
		attribute.String("org.id", id),
		attribute.String("org.name", name),
		attribute.Int("org.description_length", len(description)),
	)

	s.logger.Info("Updating organization", zap.String("id", id))

	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		status = "error"
		s.logger.Error("Failed to get organization", zap.Error(err))
		return nil, err
	}

	// Обновление полей
	if name != "" {
		org.Name = name
	}
	var descPtr *string
	if description != "" {
		descPtr = &description
	}
	org.Description = descPtr
	org.UpdatedAt = time.Now()

	if err := org.Validate(); err != nil {
		status = "error"
		s.logger.Error("Organization validation failed", zap.Error(err))
		return nil, err
	}

	if err := s.orgRepo.Update(ctx, org); err != nil {
		status = "error"
		s.logger.Error("Failed to update organization", zap.Error(err))
		return nil, err
	}

	// Публикация события
	eventUpdDesc := ""
	if org.Description != nil {
		eventUpdDesc = *org.Description
	}
	event := events.OrganizationUpdated{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: events.OrganizationUpdatedEvent,
			Timestamp: time.Now(),
		},
		OrganizationID: org.ID,
		Name:           org.Name,
		Description:    eventUpdDesc,
	}

	if eventData, err := event.BaseEvent.ToJSON(); err == nil {
		observability.RecordEventPublished(events.OrganizationUpdatedEvent, "success")
		_ = s.eventPublisher.Publish(ctx, events.OrganizationUpdatedEvent, eventData)
	}

	s.logger.Info("Organization updated successfully", zap.String("id", org.ID))
	return org, nil
}

// DeleteOrganization удаляет организацию
func (s *Service) DeleteOrganization(ctx context.Context, id string) error {
	// Начало отслеживания времени для метрик
	startTime := time.Now()
	status := "success"
	defer func() {
		duration := time.Since(startTime)
		observability.RecordOperationDuration("DeleteOrganization", duration, status)
		observability.RecordOrganizationOperation("DeleteOrganization", status)
	}()

	// Создание span для распределенной трассировки
	tracer := otel.Tracer("company-service")
	_, span := tracer.Start(ctx, "DeleteOrganization")
	defer span.End()

	span.SetAttributes(attribute.String("org.id", id))

	s.logger.Info("Deleting organization", zap.String("id", id))

	// Проверка существования
	_, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		status = "error"
		s.logger.Error("Failed to get organization", zap.Error(err))
		return err
	}

	if err := s.orgRepo.Delete(ctx, id); err != nil {
		status = "error"
		s.logger.Error("Failed to delete organization", zap.Error(err))
		return err
	}

	// Публикация события
	event := events.OrganizationDeleted{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: events.OrganizationDeletedEvent,
			Timestamp: time.Now(),
		},
		OrganizationID: id,
	}

	if eventData, err := event.BaseEvent.ToJSON(); err == nil {
		observability.RecordEventPublished(events.OrganizationDeletedEvent, "success")
		_ = s.eventPublisher.Publish(ctx, events.OrganizationDeletedEvent, eventData)
	}

	s.logger.Info("Organization deleted successfully", zap.String("id", id))
	return nil
}

// ListOrganizations возвращает список организаций
func (s *Service) ListOrganizations(ctx context.Context, page, perPage int32, ownerID string) ([]*domain.Organization, int32, error) {
	// Начало отслеживания времени для метрик
	startTime := time.Now()
	status := "success"
	defer func() {
		duration := time.Since(startTime)
		observability.RecordOperationDuration("ListOrganizations", duration, status)
		observability.RecordOrganizationOperation("ListOrganizations", status)
	}()

	// Создание span для распределенной трассировки
	tracer := otel.Tracer("company-service")
	_, span := tracer.Start(ctx, "ListOrganizations")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("pagination.page", int64(page)),
		attribute.Int64("pagination.per_page", int64(perPage)),
		attribute.String("org.owner_id", ownerID),
	)

	s.logger.Info("Listing organizations", zap.Int32("page", page), zap.Int32("per_page", perPage))

	orgs, total, err := s.orgRepo.List(ctx, page, perPage, ownerID)
	if err != nil {
		status = "error"
		s.logger.Error("Failed to list organizations", zap.Error(err))
		return nil, 0, err
	}

	return orgs, total, nil
}

// AddMember добавляет участника в организацию
func (s *Service) AddMember(ctx context.Context, organizationID, userID, role string) (*domain.Employee, error) {
	// Начало отслеживания времени для метрик
	startTime := time.Now()
	status := "success"
	defer func() {
		duration := time.Since(startTime)
		observability.RecordOperationDuration("AddMember", duration, status)
		observability.RecordEmployeeOperation("AddMember", status)
	}()

	// Создание span для распределенной трассировки
	tracer := otel.Tracer("company-service")
	_, span := tracer.Start(ctx, "AddMember")
	defer span.End()

	span.SetAttributes(
		attribute.String("org.id", organizationID),
		attribute.String("user.id", userID),
		attribute.String("employee.role", role),
	)

	s.logger.Info("Adding member", zap.String("org_id", organizationID), zap.String("user_id", userID))

	// Проверка существования организации
	_, err := s.orgRepo.GetByID(ctx, organizationID)
	if err != nil {
		status = "error"
		s.logger.Error("Failed to get organization", zap.Error(err))
		return nil, err
	}

	// Проверка, не является ли пользователь уже участником
	exists, err := s.empRepo.ExistsInOrganization(ctx, userID, organizationID)
	if err != nil {
		status = "error"
		s.logger.Error("Failed to check employee existence", zap.Error(err))
		return nil, err
	}
	if exists {
		status = "error"
		s.logger.Warn("Employee already exists in organization")
		return nil, domain.ErrEmployeeExists
	}

	// Создание участника
	emp := &domain.Employee{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		UserID:         userID,
		Role:           role,
		JoinedAt:       time.Now(),
	}

	if err := emp.Validate(); err != nil {
		status = "error"
		s.logger.Error("Employee validation failed", zap.Error(err))
		return nil, err
	}

	if err := s.empRepo.Create(ctx, emp); err != nil {
		status = "error"
		s.logger.Error("Failed to create employee", zap.Error(err))
		return nil, err
	}

	// Публикация события
	event := events.EmployeeAdded{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: events.EmployeeAddedEvent,
			Timestamp: time.Now(),
		},
		EmployeeID:     emp.ID,
		OrganizationID: emp.OrganizationID,
		UserID:         emp.UserID,
		Role:           emp.Role,
	}

	if eventData, err := event.BaseEvent.ToJSON(); err == nil {
		observability.RecordEventPublished(events.EmployeeAddedEvent, "success")
		_ = s.eventPublisher.Publish(ctx, events.EmployeeAddedEvent, eventData)
	}

	s.logger.Info("Member added successfully", zap.String("employee_id", emp.ID))
	return emp, nil
}

// RemoveMember удаляет участника из организации
func (s *Service) RemoveMember(ctx context.Context, organizationID, employeeID string) error {
	// Начало отслеживания времени для метрик
	startTime := time.Now()
	status := "success"
	defer func() {
		duration := time.Since(startTime)
		observability.RecordOperationDuration("RemoveMember", duration, status)
		observability.RecordEmployeeOperation("RemoveMember", status)
	}()

	// Создание span для распределенной трассировки
	tracer := otel.Tracer("company-service")
	_, span := tracer.Start(ctx, "RemoveMember")
	defer span.End()

	span.SetAttributes(
		attribute.String("org.id", organizationID),
		attribute.String("employee.id", employeeID),
	)

	s.logger.Info("Removing member", zap.String("org_id", organizationID), zap.String("employee_id", employeeID))

	// Получение участника
	emp, err := s.empRepo.GetByID(ctx, employeeID)
	if err != nil {
		status = "error"
		s.logger.Error("Failed to get employee", zap.Error(err))
		return err
	}

	if emp.OrganizationID != organizationID {
		status = "error"
		s.logger.Error("Employee does not belong to organization")
		return domain.ErrEmployeeNotFound
	}

	if err := s.empRepo.Delete(ctx, employeeID); err != nil {
		status = "error"
		s.logger.Error("Failed to delete employee", zap.Error(err))
		return err
	}

	// Публикация события
	event := events.EmployeeRemoved{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: events.EmployeeRemovedEvent,
			Timestamp: time.Now(),
		},
		EmployeeID:     emp.ID,
		OrganizationID: emp.OrganizationID,
		UserID:         emp.UserID,
	}

	if eventData, err := event.BaseEvent.ToJSON(); err == nil {
		observability.RecordEventPublished(events.EmployeeRemovedEvent, "success")
		_ = s.eventPublisher.Publish(ctx, events.EmployeeRemovedEvent, eventData)
	}

	s.logger.Info("Member removed successfully", zap.String("employee_id", employeeID))
	return nil
}

// GetOrganizationMembers возвращает список участников организации
func (s *Service) GetOrganizationMembers(ctx context.Context, organizationID string, page, perPage int32) ([]*domain.Employee, int32, error) {
	// Начало отслеживания времени для метрик
	startTime := time.Now()
	status := "success"
	defer func() {
		duration := time.Since(startTime)
		observability.RecordOperationDuration("GetOrganizationMembers", duration, status)
		observability.RecordEmployeeOperation("GetOrganizationMembers", status)
	}()

	// Создание span для распределенной трассировки
	tracer := otel.Tracer("company-service")
	_, span := tracer.Start(ctx, "GetOrganizationMembers")
	defer span.End()

	span.SetAttributes(
		attribute.String("org.id", organizationID),
		attribute.Int64("pagination.page", int64(page)),
		attribute.Int64("pagination.per_page", int64(perPage)),
	)

	s.logger.Info("Getting organization members", zap.String("org_id", organizationID))

	// Проверка существования организации
	_, err := s.orgRepo.GetByID(ctx, organizationID)
	if err != nil {
		status = "error"
		s.logger.Error("Failed to get organization", zap.Error(err))
		return nil, 0, err
	}

	employees, total, err := s.empRepo.ListByOrganization(ctx, organizationID, page, perPage)
	if err != nil {
		status = "error"
		s.logger.Error("Failed to list employees", zap.Error(err))
		return nil, 0, err
	}

	return employees, total, nil
}
