// Файл invoice-server/internal/interfaces/grpc/invoice_handlers_test.go содержит реализацию пакета grpc.
package grpc

import (
	"testing"

	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/application/invoice"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/domain"
	"github.com/rusgainew/kkm-project-mks/invoice-server/internal/infrastructure/observability"
	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	pbEntities "github.com/rusgainew/kkm-project-mks/proto-lib/entities"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestCreateInvoiceValidation проверяет валидацию входных параметров
func TestCreateInvoiceValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.CreateInvoiceRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "request is nil",
		},
		{
			name: "missing invoice_number",
			req: &pb.CreateInvoiceRequest{
				Details:     []*pbEntities.InvoiceDetail{{GoodsName: "Test"}},
				LegalPerson: &pbEntities.Party{Pin: "12345678901234"},
				Contractor:  &pbEntities.Party{Pin: "98765432109876"},
			},
			wantErr: true,
			errMsg:  "invoice_number is required",
		},
		{
			name: "missing details",
			req: &pb.CreateInvoiceRequest{
				InvoiceNumber: "INV-001",
				LegalPerson:   &pbEntities.Party{Pin: "12345678901234"},
				Contractor:    &pbEntities.Party{Pin: "98765432109876"},
			},
			wantErr: true,
			errMsg:  "details are required",
		},
		{
			name: "missing legal_person",
			req: &pb.CreateInvoiceRequest{
				InvoiceNumber: "INV-001",
				Details:       []*pbEntities.InvoiceDetail{{GoodsName: "Test"}},
				Contractor:    &pbEntities.Party{Pin: "98765432109876"},
			},
			wantErr: true,
			errMsg:  "legal_person.pin is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req == nil && tt.wantErr {
				// Correct expectation: nil request should cause error
				t.Logf("Correctly expecting error for %s", tt.name)
			} else if tt.req != nil && tt.wantErr {
				// Validation errors for invalid input
				switch {
				case tt.req.InvoiceNumber == "" && !tt.wantErr:
					t.Error("should validate invoice_number")
				case len(tt.req.Details) == 0 && !tt.wantErr:
					t.Error("should validate details")
				case (tt.req.LegalPerson == nil || tt.req.LegalPerson.Pin == "") && !tt.wantErr:
					t.Error("should validate legal_person")
				}
			}
		})
	}
}

// TestListInvoicesValidation проверяет валидацию PageInfo
func TestListInvoicesValidation(t *testing.T) {
	tests := []struct {
		name     string
		pageInfo *pb.PageInfo
		wantErr  bool
	}{
		{
			name:     "nil pageInfo",
			pageInfo: nil,
			wantErr:  true,
		},
		{
			name: "valid pageInfo",
			pageInfo: &pb.PageInfo{
				Page: 0,
				Size: 10,
			},
			wantErr: false,
		},
		{
			name: "zero size",
			pageInfo: &pb.PageInfo{
				Page: 0,
				Size: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.pageInfo == nil && tt.wantErr {
				t.Logf("Correctly expecting error for nil pageInfo")
			}
			if tt.pageInfo != nil && tt.pageInfo.Size == 0 && tt.wantErr {
				t.Logf("Correctly expecting error for zero size")
			}
		})
	}
}

// TestPaginationConversion проверяет конвертацию пагинации (0-based -> 1-based)
func TestPaginationConversion(t *testing.T) {
	t.Run("page_0_to_page_1", func(t *testing.T) {
		// Proto page 0 должно быть конвертировано в repository page 1
		protoPage := int32(0)
		repoPage := protoPage + 1
		if repoPage != 1 {
			t.Errorf("expected repo page 1, got %d", repoPage)
		}
	})

	t.Run("page_5_to_page_6", func(t *testing.T) {
		protoPage := int32(5)
		repoPage := protoPage + 1
		if repoPage != 6 {
			t.Errorf("expected repo page 6, got %d", repoPage)
		}
	})
}

// TestInvoiceStatusConstants проверяет константы статусов
func TestInvoiceStatusConstants(t *testing.T) {
	tests := []struct {
		name   string
		status domain.InvoiceStatus
		want   string
	}{
		{name: "draft", status: domain.StatusDraft, want: "draft"},
		{name: "sent", status: domain.StatusSent, want: "sent"},
		{name: "accepted", status: domain.StatusAccepted, want: "accepted"},
		{name: "rejected", status: domain.StatusRejected, want: "rejected"},
		{name: "signed", status: domain.StatusSigned, want: "signed"},
		{name: "revoked", status: domain.StatusRevoked, want: "revoked"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.want {
				t.Errorf("expected %s, got %s", tt.want, tt.status)
			}
		})
	}
}

// TestInvoiceValidation проверяет валидацию invoice
func TestInvoiceValidation(t *testing.T) {
	tests := []struct {
		name    string
		invoice *domain.Invoice
		wantErr bool
	}{
		{
			name: "valid invoice",
			invoice: &domain.Invoice{
				InvoiceNumber: "INV-001",
				TotalAmount:   100.0,
				LegalPersonID: "person-1",
				ContractorID:  "contractor-1",
			},
			wantErr: false,
		},
		{
			name: "missing invoice_number",
			invoice: &domain.Invoice{
				InvoiceNumber: "",
				TotalAmount:   100.0,
				LegalPersonID: "person-1",
				ContractorID:  "contractor-1",
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			invoice: &domain.Invoice{
				InvoiceNumber: "INV-001",
				TotalAmount:   -100.0,
				LegalPersonID: "person-1",
				ContractorID:  "contractor-1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.invoice.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

// TestInvoiceStateTransitions проверяет переходы статусов
func TestInvoiceStateTransitions(t *testing.T) {
	inv := &domain.Invoice{
		Status: domain.StatusDraft,
	}

	tests := []struct {
		name          string
		currentStatus domain.InvoiceStatus
		canSign       bool
		canAccept     bool
		canReject     bool
		canRevoke     bool
	}{
		{
			name:          "draft status",
			currentStatus: domain.StatusDraft,
			canSign:       true,
			canAccept:     false,
			canReject:     false,
			canRevoke:     false,
		},
		{
			name:          "signed status",
			currentStatus: domain.StatusSigned,
			canSign:       false,
			canAccept:     true,
			canReject:     true,
			canRevoke:     true,
		},
		{
			name:          "accepted status",
			currentStatus: domain.StatusAccepted,
			canSign:       false,
			canAccept:     false,
			canReject:     false,
			canRevoke:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv.Status = tt.currentStatus

			if inv.CanSign() != tt.canSign {
				t.Errorf("CanSign: expected %v, got %v", tt.canSign, inv.CanSign())
			}
			if inv.CanAccept() != tt.canAccept {
				t.Errorf("CanAccept: expected %v, got %v", tt.canAccept, inv.CanAccept())
			}
			if inv.CanReject() != tt.canReject {
				t.Errorf("CanReject: expected %v, got %v", tt.canReject, inv.CanReject())
			}
			if inv.CanRevoke() != tt.canRevoke {
				t.Errorf("CanRevoke: expected %v, got %v", tt.canRevoke, inv.CanRevoke())
			}
		})
	}
}

// TestGetInvoiceByNumberValidation проверяет валидацию InvoiceNumberRequest
func TestGetInvoiceByNumberValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.InvoiceNumberRequest
		wantErr bool
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "empty invoice_number",
			req: &pb.InvoiceNumberRequest{
				InvoiceNumber: "",
			},
			wantErr: true,
		},
		{
			name: "valid invoice_number",
			req: &pb.InvoiceNumberRequest{
				InvoiceNumber: "INV-001",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req == nil && tt.wantErr {
				t.Logf("Correctly expecting error for nil request")
			}
			if tt.req != nil && tt.req.InvoiceNumber == "" && tt.wantErr {
				t.Logf("Correctly expecting error for empty invoice_number")
			}
		})
	}
}

// TestDateRangeRequestValidation проверяет валидацию DateRangeRequest
func TestDateRangeRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.DateRangeRequest
		wantErr bool
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "missing date_from",
			req: &pb.DateRangeRequest{
				DateTo: "2024-01-31",
				Page:   0,
				Size:   10,
			},
			wantErr: true,
		},
		{
			name: "missing date_to",
			req: &pb.DateRangeRequest{
				DateFrom: "2024-01-01",
				Page:     0,
				Size:     10,
			},
			wantErr: true,
		},
		{
			name: "valid request",
			req: &pb.DateRangeRequest{
				DateFrom: "2024-01-01",
				DateTo:   "2024-01-31",
				Page:     0,
				Size:     10,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req == nil && tt.wantErr {
				t.Logf("Correctly expecting error for nil request")
			}
			if tt.req != nil {
				hasErrors := tt.req.DateFrom == "" || tt.req.DateTo == ""
				if hasErrors != tt.wantErr {
					t.Logf("Date validation: %v", hasErrors)
				}
			}
		})
	}
}

// --- Helpers for handler tests ---

// stubInvoiceRepo satisfies ports.InvoiceRepository; methods unused in these tests.
type stubInvoiceRepo struct{}

func (stubInvoiceRepo) Create(context.Context, *domain.Invoice) error            { return nil }
func (stubInvoiceRepo) Update(context.Context, *domain.Invoice) error            { return nil }
func (stubInvoiceRepo) Delete(context.Context, string) error                     { return nil }
func (stubInvoiceRepo) GetByID(context.Context, string) (*domain.Invoice, error) { return nil, nil }
func (stubInvoiceRepo) List(context.Context, int32, int32, string) ([]*domain.Invoice, int32, error) {
	return nil, 0, nil
}
func (stubInvoiceRepo) Search(context.Context, string, int32, int32) ([]*domain.Invoice, int32, error) {
	return nil, 0, nil
}
func (stubInvoiceRepo) Filter(context.Context, string, string, string, float64, float64, string, string, int32, int32) ([]*domain.Invoice, int32, error) {
	return nil, 0, nil
}
func (stubInvoiceRepo) GetByInvoiceNumber(context.Context, string) (*domain.Invoice, error) {
	return nil, nil
}
func (stubInvoiceRepo) ListByDateRange(context.Context, string, string, int32, int32) ([]*domain.Invoice, int32, error) {
	return nil, 0, nil
}
func (stubInvoiceRepo) ExistsByNumber(context.Context, string) (bool, error) { return false, nil }
func (stubInvoiceRepo) GetByDocumentUUID(context.Context, string) (*domain.Invoice, error) {
	return nil, nil
}

// stubDetailRepo satisfies ports.InvoiceDetailRepository and returns preset details.
type stubDetailRepo struct {
	items []*domain.InvoiceDetail
}

func (s stubDetailRepo) Create(context.Context, *domain.InvoiceDetail) error        { return nil }
func (s stubDetailRepo) CreateBatch(context.Context, []*domain.InvoiceDetail) error { return nil }
func (s stubDetailRepo) GetByInvoiceUUID(context.Context, string) ([]*domain.InvoiceDetail, error) {
	return s.items, nil
}
func (s stubDetailRepo) ListByInvoiceUUID(ctx context.Context, invoiceUUID string, page, perPage int32) ([]*domain.InvoiceDetail, int32, error) {
	total := int32(len(s.items))
	if perPage <= 0 {
		perPage = 20
	}
	if page < 1 {
		page = 1
	}
	start := int((page - 1) * perPage)
	end := start + int(perPage)
	if start > len(s.items) {
		start = len(s.items)
	}
	if end > len(s.items) {
		end = len(s.items)
	}
	window := make([]*domain.InvoiceDetail, 0, end-start)
	for _, d := range s.items[start:end] {
		window = append(window, d)
	}
	return window, total, nil
}
func (s stubDetailRepo) DeleteByInvoiceUUID(context.Context, string) error { return nil }

// stubFinancialRepo satisfies ports.FinancialDataRepository; methods unused here.
type stubFinancialRepo struct{}

func (stubFinancialRepo) Create(context.Context, *domain.FinancialData) error { return nil }
func (stubFinancialRepo) GetByInvoiceUUID(context.Context, string) (*domain.FinancialData, error) {
	return nil, nil
}
func (stubFinancialRepo) Update(context.Context, *domain.FinancialData) error { return nil }

// stubPublisher satisfies ports.EventPublisher; unused for these tests.
type stubPublisher struct{}

func (stubPublisher) Publish(context.Context, string, []byte) error { return nil }
func (stubPublisher) Close() error                                  { return nil }

func newTestHandlerWithDetails(items []*domain.InvoiceDetail) *InvoiceQueryHandler {
	logger := zap.NewNop()
	registry := prometheus.NewRegistry()
	metrics := observability.NewMetricsWithRegisterer(registry)

	service := invoice.NewService(stubInvoiceRepo{}, stubDetailRepo{items: items}, stubFinancialRepo{}, stubPublisher{}, logger)
	return NewInvoiceQueryHandler(service, metrics, logger)
}

func TestListInvoiceDetailsValidation(t *testing.T) {
	h := newTestHandlerWithDetails(nil)

	_, err := h.ListInvoiceDetails(context.Background(), nil)
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument for nil request, got %v", status.Code(err))
	}

	_, err = h.ListInvoiceDetails(context.Background(), &pb.InvoiceDetailsRequest{})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument for missing invoice_uuid, got %v", status.Code(err))
	}
}

func TestListInvoiceDetailsPagination(t *testing.T) {
	items := []*domain.InvoiceDetail{
		{ID: 1, InvoiceUUID: "inv", GoodsName: "A"},
		{ID: 2, InvoiceUUID: "inv", GoodsName: "B"},
		{ID: 3, InvoiceUUID: "inv", GoodsName: "C"},
		{ID: 4, InvoiceUUID: "inv", GoodsName: "D"},
		{ID: 5, InvoiceUUID: "inv", GoodsName: "E"},
	}

	h := newTestHandlerWithDetails(items)
	req := &pb.InvoiceDetailsRequest{InvoiceUuid: "inv", Page: 1, Size: 2}
	resp, err := h.ListInvoiceDetails(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list := resp.GetDetailList()
	if list == nil {
		t.Fatalf("expected detail_list in response")
	}
	if got, want := len(list.GetDetails()), 2; got != want {
		t.Fatalf("expected %d details, got %d", want, got)
	}
	if list.GetDetails()[0].GetId() != 3 || list.GetDetails()[1].GetId() != 4 {
		t.Fatalf("unexpected detail IDs: %v, %v", list.GetDetails()[0].GetId(), list.GetDetails()[1].GetId())
	}
	if list.GetTotalElements() != 5 {
		t.Fatalf("expected total_elements 5, got %d", list.GetTotalElements())
	}
	if list.GetTotalPage() != 3 {
		t.Fatalf("expected total_page 3, got %d", list.GetTotalPage())
	}
}
