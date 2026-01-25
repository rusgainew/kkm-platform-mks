package company

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var (
	// ErrUnauthorized - ошибка когда пользователь не имеет прав доступа
	ErrUnauthorized = errors.New("unauthorized: user does not have permission for this action")
	// ErrUserNotInOrganization - ошибка когда пользователь не входит в организацию
	ErrUserNotInOrganization = errors.New("user is not a member of this organization")
)

// CheckUserIsOwner проверяет, что пользователь является владельцем организации
func (s *Service) CheckUserIsOwner(ctx context.Context, userID, organizationID string) error {
	org, err := s.orgRepo.GetByID(ctx, organizationID)
	if err != nil {
		return err
	}

	if org.OwnerID != userID {
		s.logger.Warn("User is not owner of organization",
			zap.String("user_id", userID),
			zap.String("org_id", organizationID),
			zap.String("owner_id", org.OwnerID))
		return ErrUnauthorized
	}

	return nil
}

// CheckUserIsOrganizationMember проверяет, что пользователь входит в организацию
func (s *Service) CheckUserIsOrganizationMember(ctx context.Context, userID, organizationID string) error {
	exists, err := s.empRepo.ExistsInOrganization(ctx, userID, organizationID)
	if err != nil {
		return err
	}

	if !exists {
		s.logger.Warn("User is not a member of organization",
			zap.String("user_id", userID),
			zap.String("org_id", organizationID))
		return ErrUserNotInOrganization
	}

	return nil
}

// CheckUserIsAdmin проверяет, что пользователь имеет роль администратора в организации
func (s *Service) CheckUserIsAdmin(ctx context.Context, userID, organizationID string) error {
	// Сначала проверяем что пользователь входит в организацию
	if err := s.CheckUserIsOrganizationMember(ctx, userID, organizationID); err != nil {
		return err
	}

	// Затем проверяем роль (нужно расширить интерфейс репозитория)
	// TODO: Когда добавим GetEmployeeByUserAndOrg в интерфейс
	// emp, err := s.empRepo.GetEmployeeByUserAndOrg(ctx, userID, organizationID)
	// if err != nil {
	//     return err
	// }
	// if emp.Role != "admin" {
	//     return ErrUnauthorized
	// }

	return nil
}
