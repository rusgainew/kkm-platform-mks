package grpc

import (
	"context"
	"time"

	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/application/invoice"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	pbEntities "github.com/rusgainew/kkm-project-mks/proto-lib/entities"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// InvoiceCommandHandler обработчик команд для счетов
type InvoiceCommandHandler struct {
	pb.UnimplementedInvoiceCommandServiceServer
	service *invoice.Service
	metrics *observability.Metrics
	logger  *zap.Logger
}

// NewInvoiceCommandHandler создает новый обработчик команд
func NewInvoiceCommandHandler(service *invoice.Service, metrics *observability.Metrics, logger *zap.Logger) *InvoiceCommandHandler {
	return &InvoiceCommandHandler{
		service: service,
		metrics: metrics,
		logger:  logger,
	}
}

// CreateInvoice создает новый счет
func (h *InvoiceCommandHandler) CreateInvoice(ctx context.Context, req *pb.CreateInvoiceRequest) (*pbEntities.Invoice, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if req.GetInvoiceNumber() == "" {
		return nil, status.Error(codes.InvalidArgument, "invoice_number is required")
	}
	if len(req.GetDetails()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "details are required")
	}
	if req.GetLegalPerson() == nil || req.GetLegalPerson().GetPin() == "" {
		return nil, status.Error(codes.InvalidArgument, "legal_person.pin is required")
	}
	if req.GetContractor() == nil || req.GetContractor().GetPin() == "" {
		return nil, status.Error(codes.InvalidArgument, "contractor.pin is required")
	}

	parseDate := func(sv *wrapperspb.StringValue) (*time.Time, error) {
		if sv == nil || sv.GetValue() == "" {
			return nil, nil
		}
		// ожидается формат YYYY-MM-DD
		t, err := time.Parse("2006-01-02", sv.GetValue())
		if err != nil {
			return nil, err
		}
		return &t, nil
	}

	inv := &domain.Invoice{
		InvoiceNumber: req.GetInvoiceNumber(),
		TotalAmount:   req.GetTotalAmount(),
		IsResident:    req.GetIsResident(),
		Note:          req.GetNote(),
		LegalPersonID: req.GetLegalPerson().GetPin(),
		ContractorID:  req.GetContractor().GetPin(),
	}

	var err error
	inv.InvoiceDate, err = parseDate(req.GetInvoiceDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid invoice_date: %v", err)
	}
	inv.DeliveryDate, err = parseDate(req.GetDeliveryDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid delivery_date: %v", err)
	}

	// map details pb -> domain
	details := make([]*domain.InvoiceDetail, 0, len(req.GetDetails()))
	for _, d := range req.GetDetails() {
		if d == nil {
			continue
		}
		details = append(details, &domain.InvoiceDetail{
			BaseCount:        d.GetBaseCount(),
			Price:            d.GetPrice(),
			Amount:           d.GetAmount(),
			AmountWithoutVAT: d.GetAmountWithoutVat(),
			AmountVAT:        d.GetAmountVat(),
			AmountST:         d.GetAmountSt(),
			GoodsName:        d.GetGoodsName(),
			TNVEDCode:        d.GetTnvedCode(),
			GKEDCode:         d.GetGkedCode(),
			FCDNumber:        d.GetFcdNumber(),
		})
	}

	created, err := h.service.CreateInvoice(ctx, inv, details, nil)
	if err != nil {
		h.logger.Error("CreateInvoice failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	// metrics
	if h.metrics != nil {
		h.metrics.InvoicesCreated.Inc()
		h.metrics.InvoicesTotal.Inc()
		h.metrics.DetailsPerInvoice.Observe(float64(len(details)))
	}

	toStr := func(t *time.Time) *wrapperspb.StringValue {
		if t == nil {
			return nil
		}
		return wrapperspb.String(t.Format("2006-01-02"))
	}

	resp := &pbEntities.Invoice{
		DocumentUuid:                 created.DocumentUUID,
		InvoiceNumber:                created.InvoiceNumber,
		Number:                       created.Number,
		InvoiceDate:                  toStr(created.InvoiceDate),
		CreatedDate:                  wrapperspb.String(created.CreatedDate.Format("2006-01-02")),
		DeliveryDate:                 toStr(created.DeliveryDate),
		CorrectedReceiptCreationDate: toStr(created.CorrectedReceiptCreationDate),
		TotalAmount:                  created.TotalAmount,
		IsResident:                   created.IsResident,
		Note:                         created.Note,
	}
	// отдадим обратно стороны, если были даны во входе
	resp.LegalPerson = req.GetLegalPerson()
	resp.Contractor = req.GetContractor()

	return resp, nil
}

// UpdateInvoice обновляет счет
func (h *InvoiceCommandHandler) UpdateInvoice(ctx context.Context, req *pb.UpdateInvoiceRequest) (*pbEntities.Invoice, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	// Валидация ID
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	invoiceID := req.GetId()

	parseDate := func(sv *wrapperspb.StringValue) (*time.Time, error) {
		if sv == nil || sv.GetValue() == "" {
			return nil, nil
		}
		t, err := time.Parse("2006-01-02", sv.GetValue())
		if err != nil {
			return nil, err
		}
		return &t, nil
	}

	updates := &domain.Invoice{
		InvoiceNumber: req.GetInvoiceNumber(),
		TotalAmount:   req.GetTotalAmount(),
		IsResident:    req.GetIsResident(),
		Note:          req.GetNote(),
	}

	var err error
	updates.InvoiceDate, err = parseDate(req.GetInvoiceDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid invoice_date: %v", err)
	}
	updates.DeliveryDate, err = parseDate(req.GetDeliveryDate())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid delivery_date: %v", err)
	}

	// map details pb -> domain если есть
	details := make([]*domain.InvoiceDetail, 0, len(req.GetDetails()))
	for _, d := range req.GetDetails() {
		if d == nil {
			continue
		}
		details = append(details, &domain.InvoiceDetail{
			BaseCount:        d.GetBaseCount(),
			Price:            d.GetPrice(),
			Amount:           d.GetAmount(),
			AmountWithoutVAT: d.GetAmountWithoutVat(),
			AmountVAT:        d.GetAmountVat(),
			AmountST:         d.GetAmountSt(),
			GoodsName:        d.GetGoodsName(),
			TNVEDCode:        d.GetTnvedCode(),
			GKEDCode:         d.GetGkedCode(),
			FCDNumber:        d.GetFcdNumber(),
		})
	}

	// Обновление счета через service layer
	updated, err := h.service.UpdateInvoice(ctx, invoiceID, updates, details)
	if err != nil {
		h.logger.Error("UpdateInvoice failed", zap.Error(err), zap.String("id", req.GetId()))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if h.metrics != nil {
		h.metrics.OperationErrors.WithLabelValues("update", "").Desc()
	}

	toStr := func(t *time.Time) *wrapperspb.StringValue {
		if t == nil {
			return nil
		}
		return wrapperspb.String(t.Format("2006-01-02"))
	}

	resp := &pbEntities.Invoice{
		DocumentUuid:                 updated.DocumentUUID,
		InvoiceNumber:                updated.InvoiceNumber,
		Number:                       updated.Number,
		InvoiceDate:                  toStr(updated.InvoiceDate),
		CreatedDate:                  wrapperspb.String(updated.CreatedDate.Format("2006-01-02")),
		DeliveryDate:                 toStr(updated.DeliveryDate),
		CorrectedReceiptCreationDate: toStr(updated.CorrectedReceiptCreationDate),
		TotalAmount:                  updated.TotalAmount,
		IsResident:                   updated.IsResident,
		Note:                         updated.Note,
	}

	return resp, nil
}

// AcceptOrRejectInvoice принимает или отклоняет счет
func (h *InvoiceCommandHandler) AcceptOrRejectInvoice(ctx context.Context, req *pb.AcceptOrRejectRequest) (*pbEntities.Invoice, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if req.GetDocumentUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "document_uuid is required")
	}

	// Получу счет по document_uuid
	invoice, err := h.service.GetInvoiceByDocumentUUID(ctx, req.GetDocumentUuid())
	if err != nil {
		h.logger.Error("AcceptOrRejectInvoice: GetInvoiceByDocumentUUID failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if invoice == nil {
		return nil, status.Error(codes.NotFound, "invoice not found")
	}

	var result *domain.Invoice

	if req.GetAccept() {
		// Принимаем счет
		result, err = h.service.AcceptInvoice(ctx, invoice.ID)
		if err != nil {
			h.logger.Error("AcceptInvoice failed", zap.Error(err), zap.String("uuid", req.GetDocumentUuid()))
			return nil, status.Error(codes.Internal, err.Error())
		}
		if h.metrics != nil {
			h.metrics.InvoicesAccepted.Inc()
		}
	} else {
		// Отклоняем счет
		result, err = h.service.RejectInvoice(ctx, invoice.ID, req.GetRejectReason())
		if err != nil {
			h.logger.Error("RejectInvoice failed", zap.Error(err), zap.String("uuid", req.GetDocumentUuid()))
			return nil, status.Error(codes.Internal, err.Error())
		}
		if h.metrics != nil {
			h.metrics.InvoicesRejected.Inc()
		}
	}

	toStr := func(t *time.Time) *wrapperspb.StringValue {
		if t == nil {
			return nil
		}
		return wrapperspb.String(t.Format("2006-01-02"))
	}

	resp := &pbEntities.Invoice{
		DocumentUuid:                 result.DocumentUUID,
		InvoiceNumber:                result.InvoiceNumber,
		Number:                       result.Number,
		InvoiceDate:                  toStr(result.InvoiceDate),
		CreatedDate:                  wrapperspb.String(result.CreatedDate.Format("2006-01-02")),
		DeliveryDate:                 toStr(result.DeliveryDate),
		CorrectedReceiptCreationDate: toStr(result.CorrectedReceiptCreationDate),
		TotalAmount:                  result.TotalAmount,
		IsResident:                   result.IsResident,
		Note:                         result.Note,
	}

	return resp, nil
}

// RevokeInvoice отзывает счет
func (h *InvoiceCommandHandler) RevokeInvoice(ctx context.Context, req *pb.RevokeRequest) (*pbEntities.Invoice, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if req.GetDocumentUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "document_uuid is required")
	}

	// Получу счет по document_uuid
	invoice, err := h.service.GetInvoiceByDocumentUUID(ctx, req.GetDocumentUuid())
	if err != nil {
		h.logger.Error("RevokeInvoice: GetInvoiceByDocumentUUID failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if invoice == nil {
		return nil, status.Error(codes.NotFound, "invoice not found")
	}

	// Отзовем счет с указанной причиной
	reason := ""
	// На случай если в proto есть поле reason, используем его
	revoked, err := h.service.RevokeInvoice(ctx, invoice.ID, reason)
	if err != nil {
		h.logger.Error("RevokeInvoice failed", zap.Error(err), zap.String("uuid", req.GetDocumentUuid()))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if h.metrics != nil {
		h.metrics.InvoicesRevoked.Inc()
	}

	toStr := func(t *time.Time) *wrapperspb.StringValue {
		if t == nil {
			return nil
		}
		return wrapperspb.String(t.Format("2006-01-02"))
	}

	resp := &pbEntities.Invoice{
		DocumentUuid:                 revoked.DocumentUUID,
		InvoiceNumber:                revoked.InvoiceNumber,
		Number:                       revoked.Number,
		InvoiceDate:                  toStr(revoked.InvoiceDate),
		CreatedDate:                  wrapperspb.String(revoked.CreatedDate.Format("2006-01-02")),
		DeliveryDate:                 toStr(revoked.DeliveryDate),
		CorrectedReceiptCreationDate: toStr(revoked.CorrectedReceiptCreationDate),
		TotalAmount:                  revoked.TotalAmount,
		IsResident:                   revoked.IsResident,
		Note:                         revoked.Note,
	}

	return resp, nil
}

// SignInvoice подписывает счет
func (h *InvoiceCommandHandler) SignInvoice(ctx context.Context, req *pb.SignRequest) (*pbEntities.Invoice, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if req.GetDocumentUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "document_uuid is required")
	}

	// Получу счет по document_uuid, затем найду его ID в domain
	invoice, err := h.service.GetInvoiceByDocumentUUID(ctx, req.GetDocumentUuid())
	if err != nil {
		h.logger.Error("SignInvoice: GetInvoiceByDocumentUUID failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if invoice == nil {
		return nil, status.Error(codes.NotFound, "invoice not found")
	}

	// Подпишем счет (signedBy извлечем из контекста или используем по умолчанию)
	signedBy := "system"
	if req.GetSignature() != "" {
		signedBy = req.GetSignature()
	}

	signed, err := h.service.SignInvoice(ctx, invoice.ID, signedBy)
	if err != nil {
		h.logger.Error("SignInvoice failed", zap.Error(err), zap.String("uuid", req.GetDocumentUuid()))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if h.metrics != nil {
		h.metrics.InvoicesSigned.Inc()
	}

	toStr := func(t *time.Time) *wrapperspb.StringValue {
		if t == nil {
			return nil
		}
		return wrapperspb.String(t.Format("2006-01-02"))
	}

	resp := &pbEntities.Invoice{
		DocumentUuid:                 signed.DocumentUUID,
		InvoiceNumber:                signed.InvoiceNumber,
		Number:                       signed.Number,
		InvoiceDate:                  toStr(signed.InvoiceDate),
		CreatedDate:                  wrapperspb.String(signed.CreatedDate.Format("2006-01-02")),
		DeliveryDate:                 toStr(signed.DeliveryDate),
		CorrectedReceiptCreationDate: toStr(signed.CorrectedReceiptCreationDate),
		TotalAmount:                  signed.TotalAmount,
		IsResident:                   signed.IsResident,
		Note:                         signed.Note,
	}

	return resp, nil
}
