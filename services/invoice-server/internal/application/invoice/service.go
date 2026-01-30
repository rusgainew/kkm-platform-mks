// Файл invoice-server/internal/application/invoice/service.go содержит реализацию пакета invoice.
package invoice

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/domain/events"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/domain/ports"
	"go.uber.org/zap"
)

// Service реализует бизнес-логику управления счетами-фактурами
type Service struct {
	invoiceRepo    ports.InvoiceRepository
	detailRepo     ports.InvoiceDetailRepository
	financialRepo  ports.FinancialDataRepository
	eventPublisher ports.EventPublisher
	logger         *zap.Logger
}

// NewService создает новый экземпляр сервиса
func NewService(
	invoiceRepo ports.InvoiceRepository,
	detailRepo ports.InvoiceDetailRepository,
	financialRepo ports.FinancialDataRepository,
	eventPublisher ports.EventPublisher,
	logger *zap.Logger,
) *Service {
	return &Service{
		invoiceRepo:    invoiceRepo,
		detailRepo:     detailRepo,
		financialRepo:  financialRepo,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// CreateInvoice создает новый счет-фактуру
func (s *Service) CreateInvoice(ctx context.Context, invoice *domain.Invoice, details []*domain.InvoiceDetail, financial *domain.FinancialData) (*domain.Invoice, error) {
	s.logger.Info("Creating invoice", zap.String("invoice_number", invoice.InvoiceNumber))

	// Проверка существования счета с таким номером
	exists, err := s.invoiceRepo.ExistsByNumber(ctx, invoice.InvoiceNumber)
	if err != nil {
		s.logger.Error("Failed to check invoice existence", zap.Error(err))
		return nil, err
	}
	if exists {
		s.logger.Warn("Invoice already exists", zap.String("invoice_number", invoice.InvoiceNumber))
		return nil, domain.ErrInvoiceExists
	}

	// Валидация
	if len(details) == 0 {
		return nil, domain.ErrEmptyDetails
	}

	// Генерация ID и UUID
	invoice.ID = uuid.New().String()
	if invoice.DocumentUUID == "" {
		invoice.DocumentUUID = uuid.New().String()
	}
	invoice.CreatedDate = time.Now()
	invoice.UpdatedAt = time.Now()
	invoice.Status = domain.StatusDraft

	if err := invoice.Validate(); err != nil {
		s.logger.Error("Invoice validation failed", zap.Error(err))
		return nil, err
	}

	// Создание счета
	if err := s.invoiceRepo.Create(ctx, invoice); err != nil {
		s.logger.Error("Failed to create invoice", zap.Error(err))
		return nil, err
	}

	// Создание позиций счета
	for _, detail := range details {
		detail.InvoiceUUID = invoice.DocumentUUID
		if err := detail.Validate(); err != nil {
			s.logger.Error("Detail validation failed", zap.Error(err))
			return nil, err
		}
	}

	if err := s.detailRepo.CreateBatch(ctx, details); err != nil {
		s.logger.Error("Failed to create invoice details", zap.Error(err))
		return nil, err
	}

	// Создание финансовых данных если есть
	if financial != nil {
		financial.InvoiceUUID = invoice.DocumentUUID
		if err := s.financialRepo.Create(ctx, financial); err != nil {
			s.logger.Error("Failed to create financial data", zap.Error(err))
			return nil, err
		}
	}

	// Публикация события
	event := events.InvoiceCreated{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: events.InvoiceCreatedEvent,
			Timestamp: time.Now(),
		},
		InvoiceID:     invoice.ID,
		DocumentUUID:  invoice.DocumentUUID,
		InvoiceNumber: invoice.InvoiceNumber,
		TotalAmount:   invoice.TotalAmount,
		CreatedBy:     invoice.CreatedBy,
	}

	if eventData, err := event.BaseEvent.ToJSON(); err == nil {
		_ = s.eventPublisher.Publish(ctx, events.InvoiceCreatedEvent, eventData)
	}

	s.logger.Info("Invoice created successfully", zap.String("id", invoice.ID))
	return invoice, nil
}

// GetInvoice получает счет-фактуру по ID
func (s *Service) GetInvoice(ctx context.Context, id string) (*domain.Invoice, error) {
	s.logger.Info("Getting invoice", zap.String("id", id))

	invoice, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get invoice", zap.Error(err))
		return nil, err
	}

	return invoice, nil
}

// GetInvoiceByDocumentUUID получает счет-фактуру по DocumentUUID
func (s *Service) GetInvoiceByDocumentUUID(ctx context.Context, uuid string) (*domain.Invoice, error) {
	s.logger.Info("Getting invoice by document UUID", zap.String("uuid", uuid))

	invoice, err := s.invoiceRepo.GetByDocumentUUID(ctx, uuid)
	if err != nil {
		s.logger.Error("Failed to get invoice", zap.Error(err))
		return nil, err
	}

	return invoice, nil
}

// UpdateInvoice обновляет счет-фактуру
func (s *Service) UpdateInvoice(ctx context.Context, id string, updates *domain.Invoice, details []*domain.InvoiceDetail) (*domain.Invoice, error) {
	s.logger.Info("Updating invoice", zap.String("id", id))

	invoice, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get invoice", zap.Error(err))
		return nil, err
	}

	// Обновление только в статусе Draft
	if invoice.Status != domain.StatusDraft {
		return nil, domain.ErrInvalidInvoiceStatus
	}

	// Обновление полей
	if updates.InvoiceNumber != "" {
		invoice.InvoiceNumber = updates.InvoiceNumber
	}
	if updates.Note != "" {
		invoice.Note = updates.Note
	}
	if updates.TotalAmount > 0 {
		invoice.TotalAmount = updates.TotalAmount
	}
	invoice.UpdatedAt = time.Now()

	if err := invoice.Validate(); err != nil {
		s.logger.Error("Invoice validation failed", zap.Error(err))
		return nil, err
	}

	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		s.logger.Error("Failed to update invoice", zap.Error(err))
		return nil, err
	}

	// Обновление позиций если есть
	if len(details) > 0 {
		// Удаляем старые
		if err := s.detailRepo.DeleteByInvoiceUUID(ctx, invoice.DocumentUUID); err != nil {
			s.logger.Error("Failed to delete old details", zap.Error(err))
			return nil, err
		}

		// Создаем новые
		for _, detail := range details {
			detail.InvoiceUUID = invoice.DocumentUUID
		}
		if err := s.detailRepo.CreateBatch(ctx, details); err != nil {
			s.logger.Error("Failed to create new details", zap.Error(err))
			return nil, err
		}
	}

	// Публикация события
	event := events.InvoiceUpdated{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: events.InvoiceUpdatedEvent,
			Timestamp: time.Now(),
		},
		InvoiceID:     invoice.ID,
		DocumentUUID:  invoice.DocumentUUID,
		InvoiceNumber: invoice.InvoiceNumber,
		TotalAmount:   invoice.TotalAmount,
	}

	if eventData, err := event.BaseEvent.ToJSON(); err == nil {
		_ = s.eventPublisher.Publish(ctx, events.InvoiceUpdatedEvent, eventData)
	}

	s.logger.Info("Invoice updated successfully", zap.String("id", invoice.ID))
	return invoice, nil
}

// SignInvoice подписывает счет-фактуру
func (s *Service) SignInvoice(ctx context.Context, id string, signedBy string) (*domain.Invoice, error) {
	s.logger.Info("Signing invoice", zap.String("id", id), zap.String("signed_by", signedBy))

	invoice, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get invoice", zap.Error(err))
		return nil, err
	}

	if err := invoice.Sign(signedBy); err != nil {
		s.logger.Error("Failed to sign invoice", zap.Error(err))
		return nil, err
	}

	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		s.logger.Error("Failed to update invoice", zap.Error(err))
		return nil, err
	}

	// Публикация события
	event := events.InvoiceSigned{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: events.InvoiceSignedEvent,
			Timestamp: time.Now(),
		},
		InvoiceID:     invoice.ID,
		DocumentUUID:  invoice.DocumentUUID,
		InvoiceNumber: invoice.InvoiceNumber,
		SignedBy:      signedBy,
	}

	if eventData, err := event.BaseEvent.ToJSON(); err == nil {
		_ = s.eventPublisher.Publish(ctx, events.InvoiceSignedEvent, eventData)
	}

	s.logger.Info("Invoice signed successfully", zap.String("id", invoice.ID))
	return invoice, nil
}

// AcceptInvoice принимает счет-фактуру
func (s *Service) AcceptInvoice(ctx context.Context, id string) (*domain.Invoice, error) {
	s.logger.Info("Accepting invoice", zap.String("id", id))

	invoice, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get invoice", zap.Error(err))
		return nil, err
	}

	if err := invoice.Accept(); err != nil {
		s.logger.Error("Failed to accept invoice", zap.Error(err))
		return nil, err
	}

	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		s.logger.Error("Failed to update invoice", zap.Error(err))
		return nil, err
	}

	// Публикация события
	event := events.InvoiceAccepted{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: events.InvoiceAcceptedEvent,
			Timestamp: time.Now(),
		},
		InvoiceID:     invoice.ID,
		DocumentUUID:  invoice.DocumentUUID,
		InvoiceNumber: invoice.InvoiceNumber,
	}

	if eventData, err := event.BaseEvent.ToJSON(); err == nil {
		_ = s.eventPublisher.Publish(ctx, events.InvoiceAcceptedEvent, eventData)
	}

	s.logger.Info("Invoice accepted successfully", zap.String("id", invoice.ID))
	return invoice, nil
}

// RejectInvoice отклоняет счет-фактуру
func (s *Service) RejectInvoice(ctx context.Context, id string, reason string) (*domain.Invoice, error) {
	s.logger.Info("Rejecting invoice", zap.String("id", id), zap.String("reason", reason))

	invoice, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get invoice", zap.Error(err))
		return nil, err
	}

	if err := invoice.Reject(); err != nil {
		s.logger.Error("Failed to reject invoice", zap.Error(err))
		return nil, err
	}

	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		s.logger.Error("Failed to update invoice", zap.Error(err))
		return nil, err
	}

	// Публикация события
	event := events.InvoiceRejected{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: events.InvoiceRejectedEvent,
			Timestamp: time.Now(),
		},
		InvoiceID:     invoice.ID,
		DocumentUUID:  invoice.DocumentUUID,
		InvoiceNumber: invoice.InvoiceNumber,
		Reason:        reason,
	}

	if eventData, err := event.BaseEvent.ToJSON(); err == nil {
		_ = s.eventPublisher.Publish(ctx, events.InvoiceRejectedEvent, eventData)
	}

	s.logger.Info("Invoice rejected successfully", zap.String("id", invoice.ID))
	return invoice, nil
}

// RevokeInvoice отзывает счет-фактуру
func (s *Service) RevokeInvoice(ctx context.Context, id string, reason string) (*domain.Invoice, error) {
	s.logger.Info("Revoking invoice", zap.String("id", id), zap.String("reason", reason))

	invoice, err := s.invoiceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get invoice", zap.Error(err))
		return nil, err
	}

	if err := invoice.Revoke(); err != nil {
		s.logger.Error("Failed to revoke invoice", zap.Error(err))
		return nil, err
	}

	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		s.logger.Error("Failed to update invoice", zap.Error(err))
		return nil, err
	}

	// Публикация события
	event := events.InvoiceRevoked{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: events.InvoiceRevokedEvent,
			Timestamp: time.Now(),
		},
		InvoiceID:     invoice.ID,
		DocumentUUID:  invoice.DocumentUUID,
		InvoiceNumber: invoice.InvoiceNumber,
		Reason:        reason,
	}

	if eventData, err := event.BaseEvent.ToJSON(); err == nil {
		_ = s.eventPublisher.Publish(ctx, events.InvoiceRevokedEvent, eventData)
	}

	s.logger.Info("Invoice revoked successfully", zap.String("id", invoice.ID))
	return invoice, nil
}

// ListInvoices возвращает список счетов-фактур
func (s *Service) ListInvoices(ctx context.Context, page, perPage int32, status string) ([]*domain.Invoice, int32, error) {
	s.logger.Info("Listing invoices", zap.Int32("page", page), zap.Int32("per_page", perPage), zap.String("status", status))

	invoices, total, err := s.invoiceRepo.List(ctx, page, perPage, status)
	if err != nil {
		s.logger.Error("Failed to list invoices", zap.Error(err))
		return nil, 0, err
	}

	return invoices, total, nil
}

// GetInvoiceByNumber получает счет-фактуру по номеру счета
func (s *Service) GetInvoiceByNumber(ctx context.Context, invoiceNumber string) (*domain.Invoice, error) {
	s.logger.Info("Getting invoice by number", zap.String("invoice_number", invoiceNumber))

	invoice, err := s.invoiceRepo.GetByInvoiceNumber(ctx, invoiceNumber)
	if err != nil {
		s.logger.Error("Failed to get invoice by number", zap.Error(err))
		return nil, err
	}

	return invoice, nil
}

// ListInvoicesByDateRange получает счета-фактуры по диапазону дат
func (s *Service) ListInvoicesByDateRange(ctx context.Context, startDate, endDate string, page, size int32) ([]*domain.Invoice, int32, error) {
	s.logger.Info("Listing invoices by date range", zap.String("start_date", startDate), zap.String("end_date", endDate), zap.Int32("page", page), zap.Int32("size", size))

	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	invoices, total, err := s.invoiceRepo.ListByDateRange(ctx, startDate, endDate, page, size)
	if err != nil {
		s.logger.Error("Failed to list invoices by date range", zap.Error(err))
		return nil, 0, err
	}

	return invoices, total, nil
}

// GetInvoiceDetails получает позиции счета-фактуры
func (s *Service) GetInvoiceDetails(ctx context.Context, invoiceUUID string) ([]*domain.InvoiceDetail, error) {
	s.logger.Info("Getting invoice details", zap.String("invoice_uuid", invoiceUUID))

	details, err := s.detailRepo.GetByInvoiceUUID(ctx, invoiceUUID)
	if err != nil {
		s.logger.Error("Failed to get invoice details", zap.Error(err))
		return nil, err
	}

	return details, nil
}

// ListInvoiceDetails получает позиции счета-фактуры с пагинацией
func (s *Service) ListInvoiceDetails(ctx context.Context, invoiceUUID string, page, size int32) ([]*domain.InvoiceDetail, int32, error) {
	s.logger.Info("Listing invoice details", zap.String("invoice_uuid", invoiceUUID), zap.Int32("page", page), zap.Int32("size", size))

	if page < 1 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}

	details, total, err := s.detailRepo.ListByInvoiceUUID(ctx, invoiceUUID, page, size)
	if err != nil {
		s.logger.Error("Failed to list invoice details", zap.Error(err))
		return nil, 0, err
	}

	return details, total, nil
}

// GetFinancialData получает финансовые данные счета
func (s *Service) GetFinancialData(ctx context.Context, invoiceUUID string) (*domain.FinancialData, error) {
	s.logger.Info("Getting financial data", zap.String("invoice_uuid", invoiceUUID))

	financial, err := s.financialRepo.GetByInvoiceUUID(ctx, invoiceUUID)
	if err != nil {
		s.logger.Error("Failed to get financial data", zap.Error(err))
		return nil, err
	}

	return financial, nil
}
