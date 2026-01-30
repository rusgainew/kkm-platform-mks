// Файл bank-account-query-server/tests/integration/filtering_test.go содержит реализацию пакета integration_test.
package integration_test

import (
	"context"
	"testing"
	"time"

	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverAddr = "localhost:50065" // bank-account-query-server port
	timeout    = 30 * time.Second
)

// TestListBankAccountsWithFilter tests bank account filtering
func TestListBankAccountsWithFilter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err, "Failed to connect to server")
	defer conn.Close()

	client := pb.NewBankAccountQueryServiceClient(conn)

	tests := []struct {
		name        string
		request     *pb.BankAccountFilterRequest
		expectError bool
		validate    func(t *testing.T, resp *pb.APIResponse)
	}{
		{
			name: "Filter by account name",
			request: &pb.BankAccountFilterRequest{
				AccountName: "Test Account",
				Page:        0,
				Size:        10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
				bankAccountList := resp.GetBankAccountList()
				assert.NotNil(t, bankAccountList)
			},
		},
		{
			name: "Filter by bank account number",
			request: &pb.BankAccountFilterRequest{
				BankAccount: "12345678901234567890",
				Page:        0,
				Size:        10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Filter by active status",
			request: &pb.BankAccountFilterRequest{
				IsActive: true,
				Page:     0,
				Size:     10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Sort by account_name ascending",
			request: &pb.BankAccountFilterRequest{
				SortField: "account_name",
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
			name: "Sort by bank_account descending",
			request: &pb.BankAccountFilterRequest{
				SortField: "bank_account",
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
			request: &pb.BankAccountFilterRequest{
				SortField: "invalid_field",
				Page:      0,
				Size:      10,
			},
			expectError: true,
		},
		{
			name: "Invalid sort order",
			request: &pb.BankAccountFilterRequest{
				SortOrder: "RANDOM",
				Page:      0,
				Size:      10,
			},
			expectError: true,
		},
		{
			name: "Search with text",
			request: &pb.BankAccountFilterRequest{
				SearchText: "test",
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
			request: &pb.BankAccountFilterRequest{
				AccountName: "Main",
				IsActive:    true,
				SortField:   "created_at",
				SortOrder:   "DESC",
				SearchText:  "account",
				Page:        0,
				Size:        10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Pagination test",
			request: &pb.BankAccountFilterRequest{
				Page: 0,
				Size: 5,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
				bankAccountList := resp.GetBankAccountList()
				assert.NotNil(t, bankAccountList)
				assert.LessOrEqual(t, len(bankAccountList.BankAccounts), 5)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.ListBankAccountsWithFilter(ctx, tt.request)
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

// TestBankAccountCacheHit tests cache functionality
func TestBankAccountCacheHit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewBankAccountQueryServiceClient(conn)

	req := &pb.BankAccountFilterRequest{
		AccountName: "CACHE_TEST",
		Page:        0,
		Size:        10,
	}

	// First request
	start1 := time.Now()
	resp1, err := client.ListBankAccountsWithFilter(ctx, req)
	duration1 := time.Since(start1)
	require.NoError(t, err)
	assert.NotNil(t, resp1)

	// Second request (cache hit)
	start2 := time.Now()
	resp2, err := client.ListBankAccountsWithFilter(ctx, req)
	duration2 := time.Since(start2)
	require.NoError(t, err)
	assert.NotNil(t, resp2)

	t.Logf("First request: %v, Second request: %v", duration1, duration2)

	// Responses should match
	assert.Equal(t, resp1.GetBankAccountList().GetTotalElements(), resp2.GetBankAccountList().GetTotalElements())
}

// TestBankAccountPaginationEdgeCases tests pagination edge cases
func TestBankAccountPaginationEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewBankAccountQueryServiceClient(conn)

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
			req := &pb.BankAccountFilterRequest{
				Page: tt.page,
				Size: tt.size,
			}
			resp, err := client.ListBankAccountsWithFilter(ctx, req)
			assert.NoError(t, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}
