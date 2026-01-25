package domain

import (
	"time"
)

// Organization представляет организацию в системе
type Organization struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	OwnerID     string    `json:"owner_id" db:"owner_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Validate проверяет валидность данных организации
func (o *Organization) Validate() error {
	if o.Name == "" {
		return ErrOrganizationNameRequired
	}
	if o.OwnerID == "" {
		return ErrInvalidUserID
	}
	return nil
}

// IsOwner проверяет, является ли пользователь владельцем организации
func (o *Organization) IsOwner(userID string) bool {
	return o.OwnerID == userID
}
