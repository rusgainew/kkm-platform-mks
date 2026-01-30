// Файл catalog-query-server/tests/integration/filtering_test.go содержит реализацию пакета integration_test.
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
	serverAddr = "localhost:50063" // catalog-query-server port
	timeout    = 30 * time.Second
)

// TestListCatalogsWithFilter tests catalog filtering
func TestListCatalogsWithFilter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err, "Failed to connect to server")
	defer conn.Close()

	client := pb.NewCatalogQueryServiceClient(conn)

	tests := []struct {
		name        string
		request     *pb.CatalogFilterRequest
		expectError bool
		validate    func(t *testing.T, resp *pb.APIResponse)
	}{
		{
			name: "Filter by name",
			request: &pb.CatalogFilterRequest{
				Name: "Test Product",
				Page: 0,
				Size: 10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
				catalogList := resp.GetCatalogList()
				assert.NotNil(t, catalogList)
			},
		},
		{
			name: "Filter by SKU",
			request: &pb.CatalogFilterRequest{
				Sku:  "SKU12345",
				Page: 0,
				Size: 10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Filter by category",
			request: &pb.CatalogFilterRequest{
				Category: "Electronics",
				Page:     0,
				Size:     10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Filter by active status",
			request: &pb.CatalogFilterRequest{
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
			name: "Filter by price range",
			request: &pb.CatalogFilterRequest{
				MinPrice: 10.0,
				MaxPrice: 100.0,
				Page:     0,
				Size:     10,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
			},
		},
		{
			name: "Sort by name ascending",
			request: &pb.CatalogFilterRequest{
				SortField: "name",
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
			name: "Sort by price descending",
			request: &pb.CatalogFilterRequest{
				SortField: "price",
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
			request: &pb.CatalogFilterRequest{
				SortField: "invalid_field",
				Page:      0,
				Size:      10,
			},
			expectError: true,
		},
		{
			name: "Invalid sort order",
			request: &pb.CatalogFilterRequest{
				SortOrder: "RANDOM",
				Page:      0,
				Size:      10,
			},
			expectError: true,
		},
		{
			name: "Search with text",
			request: &pb.CatalogFilterRequest{
				SearchText: "laptop",
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
			request: &pb.CatalogFilterRequest{
				Name:       "Product",
				Category:   "Electronics",
				IsActive:   true,
				MinPrice:   50.0,
				MaxPrice:   500.0,
				SortField:  "created_at",
				SortOrder:  "DESC",
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
			name: "Pagination test",
			request: &pb.CatalogFilterRequest{
				Page: 0,
				Size: 5,
			},
			expectError: false,
			validate: func(t *testing.T, resp *pb.APIResponse) {
				assert.NotNil(t, resp)
				catalogList := resp.GetCatalogList()
				assert.NotNil(t, catalogList)
				assert.LessOrEqual(t, len(catalogList.Catalogs), 5)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.ListCatalogsWithFilter(ctx, tt.request)
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

// TestCatalogCacheHit tests cache functionality
func TestCatalogCacheHit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCatalogQueryServiceClient(conn)

	req := &pb.CatalogFilterRequest{
		Name: "CACHE_TEST",
		Page: 0,
		Size: 10,
	}

	// First request
	start1 := time.Now()
	resp1, err := client.ListCatalogsWithFilter(ctx, req)
	duration1 := time.Since(start1)
	require.NoError(t, err)
	assert.NotNil(t, resp1)

	// Second request (cache hit)
	start2 := time.Now()
	resp2, err := client.ListCatalogsWithFilter(ctx, req)
	duration2 := time.Since(start2)
	require.NoError(t, err)
	assert.NotNil(t, resp2)

	t.Logf("First request: %v, Second request: %v", duration1, duration2)

	// Responses should match
	assert.Equal(t, resp1.GetCatalogList().GetTotalElements(), resp2.GetCatalogList().GetTotalElements())
}

// TestCatalogPaginationEdgeCases tests pagination edge cases
func TestCatalogPaginationEdgeCases(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCatalogQueryServiceClient(conn)

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
			req := &pb.CatalogFilterRequest{
				Page: tt.page,
				Size: tt.size,
			}
			resp, err := client.ListCatalogsWithFilter(ctx, req)
			assert.NoError(t, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}



































































































































































































































































































































}	}		})			}				assert.NotNil(t, resp)			if err == nil {			assert.NoError(t, err)			resp, err := client.ListCatalogsWithFilter(ctx, req)			}				Size:       10,				Page:       0,				SearchText: tt.searchText,			req := &pb.CatalogFilterRequest{		t.Run(tt.name, func(t *testing.T) {	for _, tt := range tests {	}		{"Search unicode", "товар"},		{"Search with special chars", "test-item"},		{"Search with numbers", "123"},		{"Search simple text", "product"},	}{		searchText string		name       string	tests := []struct {	client := pb.NewCatalogQueryServiceClient(conn)	defer conn.Close()	require.NoError(t, err)	)		grpc.WithBlock(),		grpc.WithTransportCredentials(insecure.NewCredentials()),	conn, err := grpc.DialContext(ctx, serverAddr,	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), timeout)func TestCatalogSearchFunctionality(t *testing.T) {// TestCatalogSearchFunctionality tests search capabilities}	}		})			}				assert.NotNil(t, resp)			if err == nil {			assert.NoError(t, err)			resp, err := client.ListCatalogsWithFilter(ctx, req)			}				Size: tt.size,				Page: tt.page,			req := &pb.CatalogFilterRequest{		t.Run(tt.name, func(t *testing.T) {	for _, tt := range tests {	}		{"Default pagination", 0, 0},		{"Large page number", 1000, 10},		{"Page 0, Size 100", 0, 100},		{"Page 0, Size 1", 0, 1},	}{		size int32		page int32		name string	tests := []struct {	client := pb.NewCatalogQueryServiceClient(conn)	defer conn.Close()	require.NoError(t, err)	)		grpc.WithBlock(),		grpc.WithTransportCredentials(insecure.NewCredentials()),	conn, err := grpc.DialContext(ctx, serverAddr,	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), timeout)func TestCatalogPaginationEdgeCases(t *testing.T) {// TestCatalogPaginationEdgeCases tests pagination edge cases}	assert.Equal(t, resp1.GetCatalogList().GetTotalElements(), resp2.GetCatalogList().GetTotalElements())	// Responses should match	t.Logf("First request: %v, Second request: %v", duration1, duration2)	assert.NotNil(t, resp2)	require.NoError(t, err)	duration2 := time.Since(start2)	resp2, err := client.ListCatalogsWithFilter(ctx, req)	start2 := time.Now()	// Second request (cache hit)	assert.NotNil(t, resp1)	require.NoError(t, err)	duration1 := time.Since(start1)	resp1, err := client.ListCatalogsWithFilter(ctx, req)	start1 := time.Now()	// First request	}		Size: 10,		Page: 0,		Name: "CACHE_TEST",	req := &pb.CatalogFilterRequest{	client := pb.NewCatalogQueryServiceClient(conn)	defer conn.Close()	require.NoError(t, err)	)		grpc.WithBlock(),		grpc.WithTransportCredentials(insecure.NewCredentials()),	conn, err := grpc.DialContext(ctx, serverAddr,	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), timeout)func TestCatalogCacheHit(t *testing.T) {// TestCatalogCacheHit tests cache functionality}	}		})			}				}					tt.validate(t, resp)				if tt.validate != nil && err == nil {				assert.NoError(t, err, "Unexpected error: %v", err)			} else {				assert.Error(t, err, "Expected error but got none")			if tt.expectError {			resp, err := client.ListCatalogsWithFilter(ctx, tt.request)		t.Run(tt.name, func(t *testing.T) {	for _, tt := range tests {	}		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:      10,				Page:      0,				GkedCode:  "456",				TnvedCode: "123",			request: &pb.CatalogFilterRequest{			name: "Multiple TNVED and GKED codes",		{		},			},				assert.LessOrEqual(t, len(catalogList.Catalogs), 5)				assert.NotNil(t, catalogList)				catalogList := resp.GetCatalogList()				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size: 5,				Page: 0,			request: &pb.CatalogFilterRequest{			name: "Pagination test",		{		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:       10,				Page:       0,				SearchText: "item",				SortOrder:  "DESC",				SortField:  "created_at",				GkedCode:   "78",				TnvedCode:  "12",				Number:     "CAT",				Name:       "Product",			request: &pb.CatalogFilterRequest{			name: "Complex filter - all fields",		{		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:       10,				Page:       0,				SearchText: "product",			request: &pb.CatalogFilterRequest{			name: "Search with text",		{		},			expectError: true,			},				Size:      10,				Page:      0,				SortOrder: "WRONG",			request: &pb.CatalogFilterRequest{			name: "Invalid sort order",		{		},			expectError: true,			},				Size:      10,				Page:      0,				SortField: "invalid_field",			request: &pb.CatalogFilterRequest{			name: "Invalid sort field",		{		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:      10,				Page:      0,				SortOrder: "DESC",				SortField: "number",			request: &pb.CatalogFilterRequest{			name: "Sort by number descending",		{		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:      10,				Page:      0,				SortOrder: "ASC",				SortField: "name",			request: &pb.CatalogFilterRequest{			name: "Sort by name ascending",		{		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:     10,				Page:     0,				GkedCode: "789",			request: &pb.CatalogFilterRequest{			name: "Filter by GKED code",		{		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:      10,				Page:      0,				TnvedCode: "123456",			request: &pb.CatalogFilterRequest{			name: "Filter by TNVED code",		{		},			},				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size:   10,				Page:   0,				Number: "CAT001",			request: &pb.CatalogFilterRequest{			name: "Filter by number",		{		},			},				assert.NotNil(t, catalogList)				catalogList := resp.GetCatalogList()				assert.NotNil(t, resp)			validate: func(t *testing.T, resp *pb.APIResponse) {			expectError: false,			},				Size: 10,				Page: 0,				Name: "Test Product",			request: &pb.CatalogFilterRequest{			name: "Filter by name",		{	}{		validate    func(t *testing.T, resp *pb.APIResponse)		expectError bool		request     *pb.CatalogFilterRequest		name        string	tests := []struct {	client := pb.NewCatalogQueryServiceClient(conn)	defer conn.Close()	require.NoError(t, err, "Failed to connect to server")	)		grpc.WithBlock(),		grpc.WithTransportCredentials(insecure.NewCredentials()),	conn, err := grpc.DialContext(ctx, serverAddr,	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), timeout)func TestListCatalogsWithFilter(t *testing.T) {// TestListCatalogsWithFilter tests catalog filtering)	timeout    = 30 * time.Second	serverAddr = "localhost:50055" // catalog-query-server portconst ()	"google.golang.org/grpc/credentials/insecure"	"google.golang.org/grpc"	"github.com/stretchr/testify/require"	"github.com/stretchr/testify/assert"	pb "github.com/rusgainew/kkm-project-mks/proto-lib/api"
