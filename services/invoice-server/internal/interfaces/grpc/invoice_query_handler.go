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

// InvoiceQueryHandler обработчик запросов для счетов
type InvoiceQueryHandler struct {
	pb.UnimplementedInvoiceQueryServiceServer
	service *invoice.Service
	metrics *observability.Metrics
	logger  *zap.Logger
}

// NewInvoiceQueryHandler создает новый обработчик запросов
func NewInvoiceQueryHandler(service *invoice.Service, metrics *observability.Metrics, logger *zap.Logger) *InvoiceQueryHandler {
	return &InvoiceQueryHandler{
		service: service,
		metrics: metrics,
		logger:  logger,
	}
}

// ListInvoices возвращает список счетов с пагинацией
func (h *InvoiceQueryHandler) ListInvoices(ctx context.Context, req *pb.PageInfo) (*pb.APIResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	// PageInfo.page — 0-базный; репозиторий ожидает 1-базный
	page := req.GetPage()
	if page < 0 {
		page = 0
	}
	repoPage := page + 1
	size := req.GetSize()
	if size <= 0 {
		size = 20
	}

	invoices, total, err := h.service.ListInvoices(ctx, repoPage, size, "")
	if err != nil {
		h.logger.Error("ListInvoices failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	// metrics: обновим общий счётчик
	if h.metrics != nil {
		h.metrics.InvoicesTotal.Set(float64(total))
	}

	mapInv := func(di *domain.Invoice) *pbEntities.Invoice {
		if di == nil {
			return nil
		}
		var invoiceDate, deliveryDate, correctedCreated *wrapperspb.StringValue
		if di.InvoiceDate != nil {
			invoiceDate = wrapperspb.String(di.InvoiceDate.Format("2006-01-02"))
		}
		if di.DeliveryDate != nil {
			deliveryDate = wrapperspb.String(di.DeliveryDate.Format("2006-01-02"))
		}
		if di.CorrectedReceiptCreationDate != nil {
			correctedCreated = wrapperspb.String(di.CorrectedReceiptCreationDate.Format("2006-01-02"))
		}
		return &pbEntities.Invoice{
			DocumentUuid:                 di.DocumentUUID,
			InvoiceNumber:                di.InvoiceNumber,
			Number:                       di.Number,
			InvoiceDate:                  invoiceDate,
			CreatedDate:                  wrapperspb.String(di.CreatedDate.Format("2006-01-02")),
			DeliveryDate:                 deliveryDate,
			CorrectedReceiptCreationDate: correctedCreated,
			TotalAmount:                  di.TotalAmount,
			IsResident:                   di.IsResident,
			Note:                         di.Note,
		}
	}

	list := make([]*pbEntities.Invoice, 0, len(invoices))
	for _, di := range invoices {
		list = append(list, mapInv(di))
	}

	// пагинация
	totalPages := int32(0)
	if size > 0 {
		totalPages = (total + size - 1) / size
	}

	resp := &pb.APIResponse{
		TotalElements: wrapperspb.Int32(total),
		TotalPage:     wrapperspb.Int32(totalPages),
		Data: &pb.APIResponse_InvoiceList{InvoiceList: &pb.InvoiceListResponse{
			Page:          page,
			Size:          size,
			TotalElements: total,
			TotalPage:     totalPages,
			Invoices:      list,
		}},
	}
	// при желании можно пробросить request_id/user_id, если они известны из метадаты

	return resp, nil
}

// ListInvoiceDetails возвращает детали счета по UUID счета с пагинацией
func (h *InvoiceQueryHandler) ListInvoiceDetails(ctx context.Context, req *pb.InvoiceDetailsRequest) (*pb.APIResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if req.GetInvoiceUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "invoice_uuid is required")
	}

	page := req.GetPage()
	if page < 0 {
		page = 0
	}
	size := req.GetSize()
	if size <= 0 || size > 100 {
		size = 20
	}

	repoPage := page + 1
	details, total, err := h.service.ListInvoiceDetails(ctx, req.GetInvoiceUuid(), repoPage, size)
	if err != nil {
		h.logger.Error("ListInvoiceDetails failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	totalPages := int32(0)
	if size > 0 {
		totalPages = (total + size - 1) / size
	}

	toPb := func(d *domain.InvoiceDetail) *pbEntities.InvoiceDetail {
		if d == nil {
			return nil
		}
		return &pbEntities.InvoiceDetail{
			Id:               d.ID,
			InvoiceUuid:      d.InvoiceUUID,
			BaseCount:        d.BaseCount,
			Price:            d.Price,
			Amount:           d.Amount,
			AmountWithoutVat: d.AmountWithoutVAT,
			AmountVat:        d.AmountVAT,
			AmountSt:         d.AmountST,
			GoodsName:        d.GoodsName,
			TnvedCode:        d.TNVEDCode,
			GkedCode:         d.GKEDCode,
			FcdNumber:        d.FCDNumber,
		}
	}

	mapped := make([]*pbEntities.InvoiceDetail, 0, len(details))
	for _, d := range details {
		mapped = append(mapped, toPb(d))
	}

	resp := &pb.APIResponse{
		TotalElements: wrapperspb.Int32(total),
		TotalPage:     wrapperspb.Int32(totalPages),
		Data: &pb.APIResponse_DetailList{DetailList: &pb.DetailListResponse{
			Page:          page,
			Size:          size,
			TotalElements: total,
			TotalPage:     totalPages,
			Details:       mapped,
		}},
	}

	return resp, nil
}

// ListInvoicesWithFilter возвращает счета с фильтрацией
func (h *InvoiceQueryHandler) ListInvoicesWithFilter(ctx context.Context, req *pb.InvoiceFilterRequest) (*pb.APIResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	// На данный момент фильтрация по статусу; можно расширить в будущем
	page := req.GetPage()
	if page < 0 {
		page = 0
	}
	repoPage := page + 1
	size := req.GetSize()
	if size <= 0 {
		size = 20
	}

	// Используем статус как строку из фильтра
	statusFilter := ""
	// NOTE: InvoiceFilterRequest должен содержать поле status; уточнить в proto.

	invoices, total, err := h.service.ListInvoices(ctx, repoPage, size, statusFilter)
	if err != nil {
		h.logger.Error("ListInvoicesWithFilter failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if h.metrics != nil {
		h.metrics.InvoicesTotal.Set(float64(total))
	}

	mapInv := func(di *domain.Invoice) *pbEntities.Invoice {
		if di == nil {
			return nil
		}
		var invoiceDate, deliveryDate, correctedCreated *wrapperspb.StringValue
		if di.InvoiceDate != nil {
			invoiceDate = wrapperspb.String(di.InvoiceDate.Format("2006-01-02"))
		}
		if di.DeliveryDate != nil {
			deliveryDate = wrapperspb.String(di.DeliveryDate.Format("2006-01-02"))
		}
		if di.CorrectedReceiptCreationDate != nil {
			correctedCreated = wrapperspb.String(di.CorrectedReceiptCreationDate.Format("2006-01-02"))
		}
		return &pbEntities.Invoice{
			DocumentUuid:                 di.DocumentUUID,
			InvoiceNumber:                di.InvoiceNumber,
			Number:                       di.Number,
			InvoiceDate:                  invoiceDate,
			CreatedDate:                  wrapperspb.String(di.CreatedDate.Format("2006-01-02")),
			DeliveryDate:                 deliveryDate,
			CorrectedReceiptCreationDate: correctedCreated,
			TotalAmount:                  di.TotalAmount,
			IsResident:                   di.IsResident,
			Note:                         di.Note,
		}
	}

	list := make([]*pbEntities.Invoice, 0, len(invoices))
	for _, di := range invoices {
		list = append(list, mapInv(di))
	}

	totalPages := int32(0)
	if size > 0 {
		totalPages = (total + size - 1) / size
	}

	resp := &pb.APIResponse{
		TotalElements: wrapperspb.Int32(total),
		TotalPage:     wrapperspb.Int32(totalPages),
		Data: &pb.APIResponse_InvoiceList{InvoiceList: &pb.InvoiceListResponse{
			Page:          page,
			Size:          size,
			TotalElements: total,
			TotalPage:     totalPages,
			Invoices:      list,
		}},
	}

	return resp, nil
}

// SearchInvoices поиск счетов по произвольному критерию
func (h *InvoiceQueryHandler) SearchInvoices(ctx context.Context, req *pb.SearchRequest) (*pb.APIResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	// SearchRequest должен содержать поле для поиска (номер, сумма и т.д.)
	// Пока используем базовый список
	page := req.GetPage()
	if page < 0 {
		page = 0
	}
	repoPage := page + 1
	size := req.GetSize()
	if size <= 0 {
		size = 20
	}

	// Вызовем ListInvoices без фильтра (можно добавить расширенную фильтрацию)
	invoices, total, err := h.service.ListInvoices(ctx, repoPage, size, "")
	if err != nil {
		h.logger.Error("SearchInvoices failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	mapInv := func(di *domain.Invoice) *pbEntities.Invoice {
		if di == nil {
			return nil
		}
		var invoiceDate, deliveryDate, correctedCreated *wrapperspb.StringValue
		if di.InvoiceDate != nil {
			invoiceDate = wrapperspb.String(di.InvoiceDate.Format("2006-01-02"))
		}
		if di.DeliveryDate != nil {
			deliveryDate = wrapperspb.String(di.DeliveryDate.Format("2006-01-02"))
		}
		if di.CorrectedReceiptCreationDate != nil {
			correctedCreated = wrapperspb.String(di.CorrectedReceiptCreationDate.Format("2006-01-02"))
		}
		return &pbEntities.Invoice{
			DocumentUuid:                 di.DocumentUUID,
			InvoiceNumber:                di.InvoiceNumber,
			Number:                       di.Number,
			InvoiceDate:                  invoiceDate,
			CreatedDate:                  wrapperspb.String(di.CreatedDate.Format("2006-01-02")),
			DeliveryDate:                 deliveryDate,
			CorrectedReceiptCreationDate: correctedCreated,
			TotalAmount:                  di.TotalAmount,
			IsResident:                   di.IsResident,
			Note:                         di.Note,
		}
	}

	list := make([]*pbEntities.Invoice, 0, len(invoices))
	for _, di := range invoices {
		list = append(list, mapInv(di))
	}

	totalPages := int32(0)
	if size > 0 {
		totalPages = (total + size - 1) / size
	}

	resp := &pb.APIResponse{
		TotalElements: wrapperspb.Int32(total),
		TotalPage:     wrapperspb.Int32(totalPages),
		Data: &pb.APIResponse_InvoiceList{InvoiceList: &pb.InvoiceListResponse{
			Page:          page,
			Size:          size,
			TotalElements: total,
			TotalPage:     totalPages,
			Invoices:      list,
		}},
	}

	return resp, nil
}

// GetInvoiceByNumber получить счет по номеру документа
func (h *InvoiceQueryHandler) GetInvoiceByNumber(ctx context.Context, req *pb.InvoiceNumberRequest) (*pb.APIResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if req.GetInvoiceNumber() == "" {
		return nil, status.Error(codes.InvalidArgument, "invoice_number is required")
	}

	// Получу счет по номеру
	invoice, err := h.service.GetInvoiceByNumber(ctx, req.GetInvoiceNumber())
	if err != nil {
		h.logger.Error("GetInvoiceByNumber failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if invoice == nil {
		return nil, status.Error(codes.NotFound, "invoice not found")
	}

	toStr := func(t *time.Time) *wrapperspb.StringValue {
		if t == nil {
			return nil
		}
		return wrapperspb.String(t.Format("2006-01-02"))
	}

	resp := &pb.APIResponse{
		Data: &pb.APIResponse_Invoice{Invoice: &pbEntities.Invoice{
			DocumentUuid:                 invoice.DocumentUUID,
			InvoiceNumber:                invoice.InvoiceNumber,
			Number:                       invoice.Number,
			InvoiceDate:                  toStr(invoice.InvoiceDate),
			CreatedDate:                  wrapperspb.String(invoice.CreatedDate.Format("2006-01-02")),
			DeliveryDate:                 toStr(invoice.DeliveryDate),
			CorrectedReceiptCreationDate: toStr(invoice.CorrectedReceiptCreationDate),
			TotalAmount:                  invoice.TotalAmount,
			IsResident:                   invoice.IsResident,
			Note:                         invoice.Note,
		}},
	}

	return resp, nil
}

// GetInvoicesByDateRange получить счета по диапазону дат
func (h *InvoiceQueryHandler) GetInvoicesByDateRange(ctx context.Context, req *pb.DateRangeRequest) (*pb.APIResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if req.GetDateFrom() == "" {
		return nil, status.Error(codes.InvalidArgument, "date_from is required")
	}
	if req.GetDateTo() == "" {
		return nil, status.Error(codes.InvalidArgument, "date_to is required")
	}

	page := req.GetPage()
	if page < 0 {
		page = 0
	}
	repoPage := page + 1
	size := req.GetSize()
	if size <= 0 {
		size = 20
	}

	invoices, total, err := h.service.ListInvoicesByDateRange(ctx, req.GetDateFrom(), req.GetDateTo(), repoPage, size)
	if err != nil {
		h.logger.Error("GetInvoicesByDateRange failed", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	mapInv := func(di *domain.Invoice) *pbEntities.Invoice {
		if di == nil {
			return nil
		}
		var invoiceDate, deliveryDate, correctedCreated *wrapperspb.StringValue
		if di.InvoiceDate != nil {
			invoiceDate = wrapperspb.String(di.InvoiceDate.Format("2006-01-02"))
		}
		if di.DeliveryDate != nil {
			deliveryDate = wrapperspb.String(di.DeliveryDate.Format("2006-01-02"))
		}
		if di.CorrectedReceiptCreationDate != nil {
			correctedCreated = wrapperspb.String(di.CorrectedReceiptCreationDate.Format("2006-01-02"))
		}
		return &pbEntities.Invoice{
			DocumentUuid:                 di.DocumentUUID,
			InvoiceNumber:                di.InvoiceNumber,
			Number:                       di.Number,
			InvoiceDate:                  invoiceDate,
			CreatedDate:                  wrapperspb.String(di.CreatedDate.Format("2006-01-02")),
			DeliveryDate:                 deliveryDate,
			CorrectedReceiptCreationDate: correctedCreated,
			TotalAmount:                  di.TotalAmount,
			IsResident:                   di.IsResident,
			Note:                         di.Note,
		}
	}

	list := make([]*pbEntities.Invoice, 0, len(invoices))
	for _, di := range invoices {
		list = append(list, mapInv(di))
	}

	totalPages := int32(0)
	if size > 0 {
		totalPages = (total + size - 1) / size
	}

	resp := &pb.APIResponse{
		TotalElements: wrapperspb.Int32(total),
		TotalPage:     wrapperspb.Int32(totalPages),
		Data: &pb.APIResponse_InvoiceList{InvoiceList: &pb.InvoiceListResponse{
			Page:          page,
			Size:          size,
			TotalElements: total,
			TotalPage:     totalPages,
			Invoices:      list,
		}},
	}

	return resp, nil
}
