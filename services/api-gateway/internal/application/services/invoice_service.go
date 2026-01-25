package services

import (
	"context"
	"fmt"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/client"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// InvoiceService сервис для работы со счетами-фактурами
type InvoiceService struct {
	connManager *client.ConnectionManager
	serviceAddr string
	metrics     *observability.Metrics
	tracer      *observability.Tracer
	logger      *zap.Logger
}

// NewInvoiceService создает новый InvoiceService
func NewInvoiceService(
	connManager *client.ConnectionManager,
	serviceAddr string,
	metrics *observability.Metrics,
	tracer *observability.Tracer,
	logger *zap.Logger,
) *InvoiceService {
	return &InvoiceService{
		connManager: connManager,
		serviceAddr: serviceAddr,
		metrics:     metrics,
		tracer:      tracer,
		logger:      logger,
	}
}

// getClient возвращает gRPC клиент для invoice service
func (s *InvoiceService) getClient(ctx context.Context) (pb.InvoiceCommandServiceClient, error) {
	conn, err := s.connManager.GetConnection(ctx, s.serviceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return pb.NewInvoiceCommandServiceClient(conn), nil
}

// CreateInvoice создает новый счет-фактуру
func (s *InvoiceService) CreateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "InvoiceService.CreateInvoice")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get invoice client", zap.Error(err))
		return nil, err
	}

	req := &pb.CreateInvoiceRequest{
		InvoiceNumber: invoice.InvoiceNumber,
		InvoiceDate:   wrapperspb.String(invoice.InvoiceDate),
		DeliveryDate:  wrapperspb.String(invoice.DeliveryDate),
		TotalAmount:   invoice.TotalAmount,
		IsResident:    invoice.IsResident,
		Note:          invoice.Note,
	}

	resp, err := client.CreateInvoice(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create invoice", zap.Error(err))
		s.metrics.IncrementErrorCount("invoice_service", "create_invoice")
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	s.logger.Info("Invoice created successfully", zap.String("uuid", resp.DocumentUuid))

	return &models.Invoice{
		ID:            resp.DocumentUuid,
		InvoiceNumber: resp.InvoiceNumber,
		InvoiceDate:   resp.InvoiceDate.GetValue(),
		DeliveryDate:  resp.DeliveryDate.GetValue(),
		TotalAmount:   resp.TotalAmount,
		IsResident:    resp.IsResident,
		Note:          resp.Note,
		Status:        resp.Status.GetCode(),
		CreatedAt:     0,
		UpdatedAt:     0,
	}, nil
}

// UpdateInvoice обновляет счет-фактуру
func (s *InvoiceService) UpdateInvoice(ctx context.Context, invoice *models.Invoice) (*models.Invoice, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "InvoiceService.UpdateInvoice")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get invoice client", zap.Error(err))
		return nil, err
	}

	req := &pb.UpdateInvoiceRequest{
		Id:            invoice.ID,
		InvoiceNumber: invoice.InvoiceNumber,
		InvoiceDate:   wrapperspb.String(invoice.InvoiceDate),
		DeliveryDate:  wrapperspb.String(invoice.DeliveryDate),
		TotalAmount:   invoice.TotalAmount,
		IsResident:    invoice.IsResident,
		Note:          invoice.Note,
	}

	resp, err := client.UpdateInvoice(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update invoice", zap.String("id", invoice.ID), zap.Error(err))
		s.metrics.IncrementErrorCount("invoice_service", "update_invoice")
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	s.logger.Info("Invoice updated successfully", zap.String("uuid", resp.DocumentUuid))

	return &models.Invoice{
		ID:            resp.DocumentUuid,
		InvoiceNumber: resp.InvoiceNumber,
		InvoiceDate:   resp.InvoiceDate.GetValue(),
		DeliveryDate:  resp.DeliveryDate.GetValue(),
		TotalAmount:   resp.TotalAmount,
		IsResident:    resp.IsResident,
		Note:          resp.Note,
		Status:        resp.Status.GetCode(),
		CreatedAt:     0,
		UpdatedAt:     0,
	}, nil
}

// SignInvoice подписывает счет-фактуру
func (s *InvoiceService) SignInvoice(ctx context.Context, id string, signatureData string) (*models.Invoice, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "InvoiceService.SignInvoice")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get invoice client", zap.Error(err))
		return nil, err
	}

	req := &pb.SignRequest{
		DocumentUuid: id,
		Signature:    signatureData,
	}

	resp, err := client.SignInvoice(ctx, req)
	if err != nil {
		s.logger.Error("Failed to sign invoice", zap.String("id", id), zap.Error(err))
		s.metrics.IncrementErrorCount("invoice_service", "sign_invoice")
		return nil, fmt.Errorf("failed to sign invoice: %w", err)
	}

	s.logger.Info("Invoice signed successfully", zap.String("uuid", resp.DocumentUuid))

	return &models.Invoice{
		ID:            resp.DocumentUuid,
		InvoiceNumber: resp.InvoiceNumber,
		InvoiceDate:   resp.InvoiceDate.GetValue(),
		DeliveryDate:  resp.DeliveryDate.GetValue(),
		TotalAmount:   resp.TotalAmount,
		IsResident:    resp.IsResident,
		Note:          resp.Note,
		Status:        resp.Status.GetCode(),
		CreatedAt:     0,
		UpdatedAt:     0,
	}, nil
}

// RevokeInvoice отзывает счет-фактуру
func (s *InvoiceService) RevokeInvoice(ctx context.Context, id string, reason string) (*models.Invoice, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "InvoiceService.RevokeInvoice")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get invoice client", zap.Error(err))
		return nil, err
	}

	req := &pb.RevokeRequest{
		DocumentUuid: id,
	}

	resp, err := client.RevokeInvoice(ctx, req)
	if err != nil {
		s.logger.Error("Failed to revoke invoice", zap.String("id", id), zap.Error(err))
		s.metrics.IncrementErrorCount("invoice_service", "revoke_invoice")
		return nil, fmt.Errorf("failed to revoke invoice: %w", err)
	}

	s.logger.Info("Invoice revoked successfully", zap.String("uuid", resp.DocumentUuid))

	return &models.Invoice{
		ID:            resp.DocumentUuid,
		InvoiceNumber: resp.InvoiceNumber,
		InvoiceDate:   resp.InvoiceDate.GetValue(),
		DeliveryDate:  resp.DeliveryDate.GetValue(),
		TotalAmount:   resp.TotalAmount,
		IsResident:    resp.IsResident,
		Note:          resp.Note,
		Status:        resp.Status.GetCode(),
		CreatedAt:     0,
		UpdatedAt:     0,
	}, nil
}

// AcceptOrRejectInvoice принимает или отклоняет счет-фактуру
func (s *InvoiceService) AcceptOrRejectInvoice(ctx context.Context, id string, accept bool, reason string) (*models.Invoice, error) {
	ctx, endSpan := s.tracer.StartSpan(ctx, "InvoiceService.AcceptOrRejectInvoice")
	defer endSpan()

	client, err := s.getClient(ctx)
	if err != nil {
		s.logger.Error("Failed to get invoice client", zap.Error(err))
		return nil, err
	}

	req := &pb.AcceptOrRejectRequest{
		DocumentUuid: id,
		Accept:       accept,
		RejectReason: reason,
	}

	resp, err := client.AcceptOrRejectInvoice(ctx, req)
	if err != nil {
		s.logger.Error("Failed to accept/reject invoice", zap.String("id", id), zap.Bool("accept", accept), zap.Error(err))
		s.metrics.IncrementErrorCount("invoice_service", "accept_reject_invoice")
		return nil, fmt.Errorf("failed to accept/reject invoice: %w", err)
	}

	s.logger.Info("Invoice accepted/rejected successfully", zap.String("uuid", resp.DocumentUuid))

	return &models.Invoice{
		ID:            resp.DocumentUuid,
		InvoiceNumber: resp.InvoiceNumber,
		InvoiceDate:   resp.InvoiceDate.GetValue(),
		DeliveryDate:  resp.DeliveryDate.GetValue(),
		TotalAmount:   resp.TotalAmount,
		IsResident:    resp.IsResident,
		Note:          resp.Note,
		Status:        resp.Status.GetCode(),
		CreatedAt:     0,
		UpdatedAt:     0,
	}, nil
}
