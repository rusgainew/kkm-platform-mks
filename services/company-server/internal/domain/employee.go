// Файл company-server/internal/domain/employee.go содержит реализацию пакета domain.
package domain

import (
	"time"
)

// Employee представляет участника организации
type Employee struct {
	ID             string    `json:"id" db:"id"`
	OrganizationID string    `json:"organization_id" db:"organization_id"`
	UserID         string    `json:"user_id" db:"user_id"`
	Role           string    `json:"role" db:"role"`
	JoinedAt       time.Time `json:"joined_at" db:"joined_at"`
}

// EmployeeRole определяет роли участников
type EmployeeRole string

const (
	RoleAdmin    EmployeeRole = "admin"
	RoleManager  EmployeeRole = "manager"
	RoleEmployee EmployeeRole = "employee"
)

// Validate проверяет валидность данных участника
func (e *Employee) Validate() error {
	if e.OrganizationID == "" {
		return ErrInvalidOrganizationID
	}
	if e.UserID == "" {
		return ErrInvalidUserID
	}
	if !isValidRole(e.Role) {
		return ErrInvalidRole
	}
	return nil
}

// isValidRole проверяет валидность роли
func isValidRole(role string) bool {
	validRoles := []string{
		string(RoleAdmin),
		string(RoleManager),
		string(RoleEmployee),
	}
	for _, r := range validRoles {
		if r == role {
			return true
		}
	}
	return false
}

// HasPermission проверяет, имеет ли участник определенные права
func (e *Employee) HasPermission(requiredRole EmployeeRole) bool {
	roleHierarchy := map[EmployeeRole]int{
		RoleAdmin:    3,
		RoleManager:  2,
		RoleEmployee: 1,
	}
	return roleHierarchy[EmployeeRole(e.Role)] >= roleHierarchy[requiredRole]
}
