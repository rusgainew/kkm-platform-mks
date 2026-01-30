// Файл company-server/internal/interfaces/grpc/validation_test.go содержит реализацию пакета grpc.
package grpc

import (
	"strings"
	"testing"
)

// contains проверяет, содержит ли строка подстроку (вспомогательная функция для тестов)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// TestValidateCreateOrganizationRequest тестирует валидацию создания организации
func TestValidateCreateOrganizationRequest(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"

	tests := []struct {
		name        string
		orgName     string
		description string
		ownerID     string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "Valid request",
			orgName:     "Test Organization",
			description: "A test organization",
			ownerID:     validUUID,
			wantErr:     false,
		},
		{
			name:        "Empty name",
			orgName:     "",
			description: "A test organization",
			ownerID:     validUUID,
			wantErr:     true,
			errMsg:      "name is required",
		},
		{
			name:        "Name too long",
			orgName:     string(make([]byte, 256)),
			description: "A test organization",
			ownerID:     validUUID,
			wantErr:     true,
			errMsg:      "too long",
		},
		{
			name:        "Invalid name characters",
			orgName:     "Test@#$%Organization",
			description: "A test organization",
			ownerID:     validUUID,
			wantErr:     true,
			errMsg:      "invalid characters",
		},
		{
			name:        "Description too long",
			orgName:     "Test Organization",
			description: string(make([]byte, 1001)),
			ownerID:     validUUID,
			wantErr:     true,
			errMsg:      "too long",
		},
		{
			name:        "Invalid owner UUID",
			orgName:     "Test Organization",
			description: "A test organization",
			ownerID:     "not-a-uuid",
			wantErr:     true,
			errMsg:      "valid UUID",
		},
		{
			name:        "Empty owner ID",
			orgName:     "Test Organization",
			description: "A test organization",
			ownerID:     "",
			wantErr:     true,
			errMsg:      "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateOrganizationRequest(tt.orgName, tt.description, tt.ownerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreateOrganizationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateCreateOrganizationRequest() got error %q, want to contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestValidateUpdateOrganizationRequest тестирует валидацию обновления организации
func TestValidateUpdateOrganizationRequest(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"

	tests := []struct {
		name        string
		id          string
		orgName     string
		description string
		wantErr     bool
	}{
		{
			name:        "Valid request",
			id:          validUUID,
			orgName:     "Updated Organization",
			description: "Updated description",
			wantErr:     false,
		},
		{
			name:        "Invalid ID",
			id:          "not-a-uuid",
			orgName:     "Updated Organization",
			description: "Updated description",
			wantErr:     true,
		},
		{
			name:        "Empty ID",
			id:          "",
			orgName:     "Updated Organization",
			description: "Updated description",
			wantErr:     true,
		},
		{
			name:        "Invalid name characters",
			id:          validUUID,
			orgName:     "Test@Organization",
			description: "Updated description",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateOrganizationRequest(tt.id, tt.orgName, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdateOrganizationRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidatePagination тестирует валидацию пагинации
func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name    string
		page    int32
		perPage int32
		wantErr bool
	}{
		{
			name:    "Valid pagination",
			page:    1,
			perPage: 10,
			wantErr: false,
		},
		{
			name:    "Page zero",
			page:    0,
			perPage: 10,
			wantErr: true,
		},
		{
			name:    "Page negative",
			page:    -1,
			perPage: 10,
			wantErr: true,
		},
		{
			name:    "PerPage zero",
			page:    1,
			perPage: 0,
			wantErr: true,
		},
		{
			name:    "PerPage too large",
			page:    1,
			perPage: 101,
			wantErr: true,
		},
		{
			name:    "PerPage max valid",
			page:    1,
			perPage: 100,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePagination(tt.page, tt.perPage)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePagination() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateAddMemberRequest тестирует валидацию добавления участника
func TestValidateAddMemberRequest(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"

	tests := []struct {
		name    string
		orgID   string
		userID  string
		role    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid request - admin",
			orgID:   validUUID,
			userID:  validUUID,
			role:    "admin",
			wantErr: false,
		},
		{
			name:    "Valid request - manager",
			orgID:   validUUID,
			userID:  validUUID,
			role:    "manager",
			wantErr: false,
		},
		{
			name:    "Valid request - employee",
			orgID:   validUUID,
			userID:  validUUID,
			role:    "employee",
			wantErr: false,
		},
		{
			name:    "Invalid role",
			orgID:   validUUID,
			userID:  validUUID,
			role:    "superadmin",
			wantErr: true,
			errMsg:  "invalid role",
		},
		{
			name:    "Empty role",
			orgID:   validUUID,
			userID:  validUUID,
			role:    "",
			wantErr: true,
		},
		{
			name:    "Invalid org UUID",
			orgID:   "not-uuid",
			userID:  validUUID,
			role:    "admin",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddMemberRequest(tt.orgID, tt.userID, tt.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAddMemberRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateListOrganizationsRequest тестирует валидацию списка организаций
func TestValidateListOrganizationsRequest(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"

	tests := []struct {
		name    string
		page    int32
		perPage int32
		ownerID string
		wantErr bool
	}{
		{
			name:    "Valid request without owner",
			page:    1,
			perPage: 10,
			ownerID: "",
			wantErr: false,
		},
		{
			name:    "Valid request with owner",
			page:    1,
			perPage: 10,
			ownerID: validUUID,
			wantErr: false,
		},
		{
			name:    "Invalid pagination",
			page:    0,
			perPage: 10,
			ownerID: "",
			wantErr: true,
		},
		{
			name:    "Invalid owner UUID",
			page:    1,
			perPage: 10,
			ownerID: "not-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateListOrganizationsRequest(tt.page, tt.perPage, tt.ownerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateListOrganizationsRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
