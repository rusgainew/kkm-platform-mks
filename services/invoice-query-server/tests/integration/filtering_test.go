package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverAddr = "localhost:50064" // invoice-query-server port
	timeout    = 30 * time.Second
)

// TestListInvoicesWithFilter tests invoice filtering
func TestListInvoicesWithFilter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err, "Failed to connect to server")
	defer conn.Close()

	client := pb.NewInvoiceQueryServiceClient(conn)

	tests := []struct {
		name        string
		request     *pb.InvoiceFilterRequest
		expectError bool
		validate    func(t *testing.T, resp *pb.APIResponse)
	}{
		{
			name: "Filter by invoice number",
			request: &pb.InvoiceFilterRequest{
				InvoiceNumber: "INV-2024-001",
				Page:          0,
				Size:          10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
				invoiceList := resp.GetInvoiceList()
				assert.NotNil(t, invoiceList)
			},
		},
		{
			name: "Filter by status",
			request: &pb.InvoiceFilterRequest{
				Status: "PAID",
				Page:   0,
				Size:   10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Filter by customer name",
			request: &pb.InvoiceFilterRequest{
				CustomerName: "Test Customer",
				Page:         0,
				Size:         10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Filter by amount range",
			request: &pb.InvoiceFilterRequest{
				MinAmount: 100.0,
				MaxAmount: 1000.0,
				Page:      0,
				Size:      10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Filter by date range",
			request: &pb.InvoiceFilterRequest{
				StartDate: "2024-01-01",
				EndDate:   "2024-12-31",
				Page:      0,
				Size:      10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Sort by invoice_number ascending",
			request: &pb.InvoiceFilterRequest{
				SortField: "invoice_number",
				SortOrder: "ASC",
				Page:      0,
				Size:      10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Sort by amount descending",
			request: &pb.InvoiceFilterRequest{
				SortField: "amount",
				SortOrder: "DESC",
				Page:      0,
				Size:      10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Sort by created_at descending",
			request: &pb.InvoiceFilterRequest{
				SortField: "created_at",
				SortOrder: "DESC",
				Page:      0,
				Size:      10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Invalid sort field",
			request: &pb.InvoiceFilterRequest{
				SortField: "invalid_field",
				Page:      0,
				Size:      10,
			},
			expectError: true,
		},
		{
			name: "Invalid sort order",
			request: &pb.InvoiceFilterRequest{
				SortOrder: "RANDOM",
				Page:      0,
				Size:      10,
			},
			expectError: true,
		},
		{
			name: "Search with text",
			request: &pb.InvoiceFilterRequest{
				SearchText: "invoice",
				Page:       0,
				Size:       10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Complex filter",
			request: &pb.InvoiceFilterRequest{
				Status:       "PENDING",
				CustomerName: "Customer",
				MinAmount:    50.0,
				MaxAmount:    500.0,
				StartDate:    "2024-01-01",
				EndDate:      "2024-12-31",
				SortField:    "created_at",
				SortOrder:    "DESC",
				SearchText:   "test",
				Page:         0,
				Size:         10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Pagination test",
			request: &pb.InvoiceFilterRequest{
				Page: 0,
				Size: 5,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
				invoiceList := resp.GetInvoiceList()
				assert.NotNil(t, invoiceList)
				assert.LessOrEqual(t, len(invoiceList.Invoices), 5)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.ListInvoicesWithFilter(ctx, tt.request)
			if tt.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error: %v", err)
				if tt.validate != nil && err == nil {
					tt.validate(t, resp)
				}
			}
		})
	}
}

// TestInvoiceCacheHit tests cache functionality
func TestInvoiceCacheHit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewInvoiceQueryServiceClient(conn)

	req := &pb.InvoiceFilterRequest{
		InvoiceNumber: "CACHE_TEST",
		Page:          0,
		Size:          10,
	}

	// First request
	start1 := time.Now()
	resp1, err := client.ListInvoicesWithFilter(ctx, req)
	duration1 := time.Since(start1)
	require.NoError(t, err)
	assert.NotNil(t, resp1)

	// Second request (cache hit)
	start2 := time.Now()
	resp2, err := client.ListInvoicesWithFilter(ctx, req)
	duration2 := time.Since(start2)
	require.NoError(t, err)
	assert.NotNil(t, resp2)

	t.Logf("First request: %v, Second request: %v", duration1, duration2)

	// Responses should match
	assert.Equal(t, resp1.GetInvoiceList().GetTotalElements(), resp2.GetInvoiceList().GetTotalElements())
}

// TestInvoicePaginationEdgeCases tests pagination edge cases
func TestInvoicePaginationEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewInvoiceQueryServiceClient(conn)

	tests := []struct {
		name string
		page int32
		size int32
	}{
		{"Page 0, Size 1", 0, 1},
		{"Page 0, Size 50", 0, 50},
		{"Large page number", 500, 10},
		{"Default pagination", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.InvoiceFilterRequest{
				Page: tt.page,
				Size: tt.size,
			}
			resp, err := client.ListInvoicesWithFilter(ctx, req)
			assert.NoError(t, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

// TestInvoiceStatusFiltering tests different status filters
func TestInvoiceStatusFiltering(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewInvoiceQueryServiceClient(conn)

	statuses := []string{"DRAFT", "PENDING", "PAID", "CANCELLED", "OVERDUE"}

	for _, status := range statuses {
		t.Run(fmt.Sprintf("Status_%s", status), func(t *testing.T) {
			req := &pb.InvoiceFilterRequest{
				Status: status,
				Page:   0,
				Size:   10,
			}
			resp, err := client.ListInvoicesWithFilter(ctx, req)
			assert.NoError(t, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}





































































































































































































































































































































































































































































































}	assert.Equal(t, 0, errorCount, "Expected no errors in concurrent requests")	}		errorCount++		t.Logf("Concurrent request error: %v", err)	for err := range errors {	errorCount := 0	// Check for errors	close(errors)	}		<-done	for i := 0; i < concurrency; i++ {	// Wait for all goroutines	}		}(i)			done <- true			}				errors <- err			if err != nil {			_, err := client.ListInvoicesWithFilter(ctx, req)			}				Size:          10,				Page:          0,				InvoiceNumber: fmt.Sprintf("CONCURRENT_%03d", index),			req := &pb.InvoiceFilterRequest{		go func(index int) {	for i := 0; i < concurrency; i++ {	done := make(chan bool, concurrency)	errors := make(chan error, concurrency)	const concurrency = 10	client := pb.NewInvoiceQueryServiceClient(conn)	defer conn.Close()	require.NoError(t, err)	)		grpc.WithBlock(),		grpc.WithTransportCredentials(insecure.NewCredentials()),	conn, err := grpc.DialContext(ctx, serverAddr,	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), timeout)func TestConcurrentRequests(t *testing.T) {// TestConcurrentRequests tests handling of concurrent requests}	}		})			}				assert.NotNil(t, resp)			if err == nil {			assert.NoError(t, err)			resp, err := client.ListInvoicesWithFilter(ctx, req)			}				Size: tt.size,				Page: tt.page,			req := &pb.InvoiceFilterRequest{		t.Run(tt.name, func(t *testing.T) {	for _, tt := range tests {	}		{"Size 0 (should use default)", 0, 0},		{"Large page number", 1000, 10},		{"Page 0, Size 100", 0, 100},		{"Page 0, Size 1", 0, 1},	}{		size int32		page int32		name string	tests := []struct {	client := pb.NewInvoiceQueryServiceClient(conn)	defer conn.Close()	require.NoError(t, err)	)		grpc.WithBlock(),		grpc.WithTransportCredentials(insecure.NewCredentials()),	conn, err := grpc.DialContext(ctx, serverAddr,	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), timeout)func TestPaginationEdgeCases(t *testing.T) {// TestPaginationEdgeCases tests edge cases in pagination}	assert.Equal(t, resp1.GetInvoiceList().GetTotalElements(), resp2.GetInvoiceList().GetTotalElements())	// Responses should be identical	t.Logf("First request: %v, Second request: %v", duration1, duration2)	// Cache hit should be faster (though this is not guaranteed in all environments)	assert.NotNil(t, resp2)	require.NoError(t, err)	duration2 := time.Since(start2)	resp2, err := client.ListInvoicesWithFilter(ctx, req)	start2 := time.Now()	// Second request - should hit cache (faster)	assert.NotNil(t, resp1)	require.NoError(t, err)	duration1 := time.Since(start1)	resp1, err := client.ListInvoicesWithFilter(ctx, req)	start1 := time.Now()	// First request - should hit database	}		Size:          10,		Page:          0,		InvoiceNumber: "CACHE_TEST_001",	req := &pb.InvoiceFilterRequest{	client := pb.NewInvoiceQueryServiceClient(conn)	defer conn.Close()	require.NoError(t, err)	)		grpc.WithBlock(),		grpc.WithTransportCredentials(insecure.NewCredentials()),	conn, err := grpc.DialContext(ctx, serverAddr,	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), timeout)func TestCacheHit(t *testing.T) {// TestCacheHit tests that second identical request hits cache}	}		})			}				}					assert.NotNil(t, resp)				if err == nil {				assert.NoError(t, err)			} else {				assert.Error(t, err)			if tt.expectError {			resp, err := client.GetInvoicesByDateRange(ctx, req)			}				Size:     10,				Page:     0,				DateTo:   tt.dateTo,				DateFrom: tt.dateFrom,			req := &pb.DateRangeRequest{		t.Run(tt.name, func(t *testing.T) {	for _, tt := range tests {	}		{"Future dates", "2025-01-01", "2025-12-31", false},		{"Invalid date format", "2024/01/01", "2024/12/31", false}, // Should still work or error gracefully		{"Missing dateTo", "2024-01-01", "", true},		{"Missing dateFrom", "", "2024-12-31", true},		{"Single day", "2024-06-15", "2024-06-15", false},		{"Valid date range", "2024-01-01", "2024-12-31", false},	}{		expectError bool		dateTo      string		dateFrom    string		name        string	tests := []struct {	client := pb.NewInvoiceQueryServiceClient(conn)	defer conn.Close()	require.NoError(t, err)	)		grpc.WithBlock(),		grpc.WithTransportCredentials(insecure.NewCredentials()),	conn, err := grpc.DialContext(ctx, serverAddr,	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), timeout)func TestGetInvoicesByDateRange(t *testing.T) {// TestGetInvoicesByDateRange tests date range filtering}	}		})			}				}					}						assert.Equal(t, int32(1), invoiceList.TotalElements)						assert.NotNil(t, invoiceList)					if tt.expectFound {					invoiceList := resp.GetInvoiceList()					assert.NotNil(t, resp)				if err == nil {				assert.NoError(t, err)			} else {				assert.Error(t, err)			if tt.expectError {			resp, err := client.GetInvoiceByNumber(ctx, req)			}				InvoiceNumber: tt.invoiceNumber,			req := &pb.InvoiceNumberRequest{		t.Run(tt.name, func(t *testing.T) {	for _, tt := range tests {	}		{"Get with empty number", "", true, false},		{"Get non-existing invoice", "NONEXISTENT999", false, false},		{"Get existing invoice", "INV001", false, true},	}{		expectFound   bool		expectError   bool		invoiceNumber string		name          string	tests := []struct {	client := pb.NewInvoiceQueryServiceClient(conn)	defer conn.Close()	require.NoError(t, err)	)		grpc.WithBlock(),		grpc.WithTransportCredentials(insecure.NewCredentials()),	conn, err := grpc.DialContext(ctx, serverAddr,	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), timeout)func TestGetInvoiceByNumber(t *testing.T) {// TestGetInvoiceByNumber tests retrieval of single invoice}	}		})			}				}					assert.NotNil(t, resp)				if err == nil {				assert.NoError(t, err)			} else {				assert.Error(t, err)			if tt.expectError {			resp, err := client.SearchInvoices(ctx, req)			}				Size:       10,				Page:       0,				SearchText: tt.searchText,			req := &pb.SearchRequest{		t.Run(tt.name, func(t *testing.T) {	for _, tt := range tests {	}		{"Search with unicode", "тест", false},		{"Search with special chars", "test@123", false},		{"Search with empty string", "", true}, // Empty search should error		{"Search with number", "123", false},		{"Search with valid text", "invoice", false},	}{		expectError bool		searchText  string		name        string	tests := []struct {	client := pb.NewInvoiceQueryServiceClient(conn)	defer conn.Close()	require.NoError(t, err)	)		grpc.WithBlock(),		grpc.WithTransportCredentials(insecure.NewCredentials()),	conn, err := grpc.DialContext(ctx, serverAddr,	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), timeout)func TestSearchInvoices(t *testing.T) {// TestSearchInvoices tests full-text search functionality}	}		})			}				}					tt.validate(t, resp)				if tt.validate != nil && err == nil {				assert.NoError(t, err, "Unexpected error: %v", err)			} else {				assert.Error(t, err, "Expected error but got none")			if tt.expectError {			resp, err := client.ListInvoicesWithFilter(ctx, tt.request)		t.Run(tt.name, func(t *testing.T) {	for _, tt := range tests {	}		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size: 5,				Page: 1,			request: &pb.InvoiceFilterRequest{			name: "Pagination - page 1",		{		},			},				assert.LessOrEqual(t, len(invoiceList.Invoices), 5)				assert.NotNil(t, invoiceList)				invoiceList := resp.GetInvoiceList()				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size: 5,				Page: 0,			request: &pb.InvoiceFilterRequest{			name: "Pagination - page 0",		{		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:       10,				Page:       0,				SearchText: "invoice",				SortOrder:  "DESC",				SortField:  "created_at",				MaxAmount:  500.0,				MinAmount:  50.0,				DateTo:     "2024-12-31",				DateFrom:   "2024-01-01",			request: &pb.InvoiceFilterRequest{			name: "Complex filter - multiple conditions",		{		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:       10,				Page:       0,				SearchText: "test",			request: &pb.InvoiceFilterRequest{			name: "Search with text",		{		},			expectError: true,			},				Size:      10,				Page:      0,				SortOrder: "INVALID",			request: &pb.InvoiceFilterRequest{			name: "Invalid sort order",		{		},			expectError: true,			},				Size:      10,				Page:      0,				SortField: "invalid_field",			request: &pb.InvoiceFilterRequest{			name: "Invalid sort field",		{		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:      10,				Page:      0,				SortOrder: "DESC",				SortField: "total_amount",			request: &pb.InvoiceFilterRequest{			name: "Sort by total_amount descending",		{		},			},				assert.NotNil(t, invoiceList)				invoiceList := resp.GetInvoiceList()				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:      10,				Page:      0,				SortOrder: "ASC",				SortField: "invoice_date",			request: &pb.InvoiceFilterRequest{			name: "Sort by invoice_date ascending",		{		},			},				}					}						assert.LessOrEqual(t, amount.Value, 1000.0)						assert.GreaterOrEqual(t, amount.Value, 100.0)					if amount != nil {					amount := invoice.GetAmount()				for _, invoice := range invoiceList.Invoices {				assert.NotNil(t, invoiceList)				invoiceList := resp.GetInvoiceList()				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:      10,				Page:      0,				MaxAmount: 1000.0,				MinAmount: 100.0,			request: &pb.InvoiceFilterRequest{			name: "Filter by amount range",		{		},			},				assert.NotNil(t, invoiceList)				invoiceList := resp.GetInvoiceList()				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:     20,				Page:     0,				DateTo:   "2024-12-31",				DateFrom: "2024-01-01",			request: &pb.InvoiceFilterRequest{			name: "Filter by date range",		{		},			},				assert.True(t, invoiceList.TotalElements <= 10)				assert.NotNil(t, invoiceList)				invoiceList := resp.GetInvoiceList()				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:          10,				Page:          0,				InvoiceNumber: "INV001",			request: &pb.InvoiceFilterRequest{			name: "Filter by invoice number",		{	}{		validate    func(t *testing.T, resp *pb.APIResponse)		expectError bool		request     *pb.InvoiceFilterRequest		name        string	tests := []struct {	client := pb.NewInvoiceQueryServiceClient(conn)	defer conn.Close()	require.NoError(t, err, "Failed to connect to server")	)		grpc.WithBlock(),		grpc.WithTransportCredentials(insecure.NewCredentials()),	conn, err := grpc.DialContext(ctx, serverAddr,	// Connect to gRPC server	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), timeout)func TestListInvoicesWithFilter(t *testing.T) {// TestListInvoicesWithFilter tests various filter combinations)	timeout    = 30 * time.Second	serverAddr = "localhost:50053" // invoice-query-server portconst ()	"google.golang.org/grpc/credentials/insecure"	"google.golang.org/grpc"	"github.com/stretchr/testify/require"	"github.com/stretchr/testify/assert"	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"