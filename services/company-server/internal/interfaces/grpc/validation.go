package grpc

import (
	"errors"
	"regexp"
)

// ValidationErrors - ошибки валидации
var (
	ErrEmptyName             = errors.New("organization name is required and cannot be empty")
	ErrNameTooLong           = errors.New("organization name is too long (max 255 characters)")
	ErrNameInvalidCharacters = errors.New("organization name contains invalid characters (only alphanumeric, spaces, and hyphens allowed)")
	ErrDescriptionTooLong    = errors.New("organization description is too long (max 1000 characters)")
	ErrEmptyOrganizationID   = errors.New("organization_id is required")
	ErrEmptyUserID           = errors.New("user_id is required")
	ErrEmptyEmployeeID       = errors.New("employee_id is required")
	ErrInvalidRole           = errors.New("invalid role (must be 'admin', 'manager', or 'employee')")
	ErrInvalidPageNumber     = errors.New("page number must be greater than 0")
	ErrInvalidPageSize       = errors.New("page size must be between 1 and 100")
	ErrEmptyUpdateData       = errors.New("at least one field must be provided for update")
)

// Valid roles for employees
var validRoles = map[string]bool{
	"admin":    true,
	"manager":  true,
	"employee": true,
}

// validateName проверяет корректность имени организации
func validateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}

	if len(name) > 255 {
		return ErrNameTooLong
	}

	// Позволяем только буквы, цифры, пробелы и дефисы
	nameRegex := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9\s\-\.]+$`)
	if !nameRegex.MatchString(name) {
		return ErrNameInvalidCharacters
	}

	return nil
}

// validateDescription проверяет корректность описания
func validateDescription(description string) error {
	if len(description) > 1000 {
		return ErrDescriptionTooLong
	}
	return nil
}

// validateID проверяет корректность ID (UUID format)
func validateID(id, fieldName string) error {
	if id == "" {
		return errors.New(fieldName + " is required")
	}
	// UUID regex: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	uuidRegex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
	if !uuidRegex.MatchString(id) {
		return errors.New(fieldName + " must be a valid UUID")
	}
	return nil
}

// validateRole проверяет корректность роли
func validateRole(role string) error {
	if role == "" {
		return errors.New("role is required")
	}
	if !validRoles[role] {
		return ErrInvalidRole
	}
	return nil
}

// validatePagination проверяет параметры пагинации
func validatePagination(page, perPage int32) error {
	if page <= 0 {
		return ErrInvalidPageNumber
	}

	if perPage <= 0 || perPage > 100 {
		return ErrInvalidPageSize
	}

	return nil
}

// ValidateCreateOrganizationRequest валидирует запрос создания организации
func ValidateCreateOrganizationRequest(name, description, ownerID string) error {
	if err := validateName(name); err != nil {
		return err
	}

	if err := validateDescription(description); err != nil {
		return err
	}

	if err := validateID(ownerID, "owner_id"); err != nil {
		return err
	}

	return nil
}

// ValidateUpdateOrganizationRequest валидирует запрос обновления организации
func ValidateUpdateOrganizationRequest(id, name, description string) error {
	if err := validateID(id, "organization_id"); err != nil {
		return err
	}

	// При обновлении, если поля не пусты, проверяем их
	if name != "" {
		if err := validateName(name); err != nil {
			return err
		}
	}

	if description != "" {
		if err := validateDescription(description); err != nil {
			return err
		}
	}

	return nil
}

// ValidateDeleteOrganizationRequest валидирует запрос удаления организации
func ValidateDeleteOrganizationRequest(organizationID string) error {
	return validateID(organizationID, "organization_id")
}

// ValidateGetOrganizationRequest валидирует запрос получения организации
func ValidateGetOrganizationRequest(organizationID string) error {
	return validateID(organizationID, "organization_id")
}

// ValidateListOrganizationsRequest валидирует запрос списка организаций
func ValidateListOrganizationsRequest(page, perPage int32, ownerID string) error {
	if err := validatePagination(page, perPage); err != nil {
		return err
	}

	// ownerID опционален, но если передан, должен быть валидным UUID
	if ownerID != "" {
		if err := validateID(ownerID, "owner_id"); err != nil {
			return err
		}
	}

	return nil
}

// ValidateAddMemberRequest валидирует запрос добавления участника
func ValidateAddMemberRequest(organizationID, userID, role string) error {
	if err := validateID(organizationID, "organization_id"); err != nil {
		return err
	}

	if err := validateID(userID, "user_id"); err != nil {
		return err
	}

	if err := validateRole(role); err != nil {
		return err
	}

	return nil
}

// ValidateRemoveMemberRequest валидирует запрос удаления участника
func ValidateRemoveMemberRequest(organizationID, employeeID string) error {
	if err := validateID(organizationID, "organization_id"); err != nil {
		return err
	}

	if err := validateID(employeeID, "employee_id"); err != nil {
		return err
	}

	return nil
}

// ValidateGetOrganizationMembersRequest валидирует запрос списка участников
func ValidateGetOrganizationMembersRequest(organizationID string, page, perPage int32) error {
	if err := validateID(organizationID, "organization_id"); err != nil {
		return err
	}

	if err := validatePagination(page, perPage); err != nil {
		return err
	}

	return nil
}
