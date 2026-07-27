package access

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"go_auth/internal/application"
	"go_auth/internal/role"
	"go_auth/internal/user"
)

var ErrAlreadyAssigned = errors.New("user already has this role in this application")

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// Assign grants roleID to userID within applicationID.
func (s *Service) Assign(userID, applicationID, roleID uint, profileID *uint) error {
	var count int64
	if err := s.DB.Model(&UserApplicationRole{}).
		Where("user_id = ? AND application_id = ? AND role_id = ?", userID, applicationID, roleID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrAlreadyAssigned
	}

	uar := &UserApplicationRole{
		UserID:        userID,
		ApplicationID: applicationID,
		RoleID:        roleID,
		ProfileID:     profileID,
		CreatedAt:     time.Now(),
	}
	if err := s.DB.Create(uar).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyAssigned
		}
		return err
	}
	return nil
}

func (s *Service) Remove(userID, applicationID, roleID uint) error {
	return s.DB.Where("user_id = ? AND application_id = ? AND role_id = ?", userID, applicationID, roleID).
		Delete(&UserApplicationRole{}).Error
}

// RemoveAllForUserInApplication revokes every role a user holds in an
// application in one go.
func (s *Service) RemoveAllForUserInApplication(userID, applicationID uint) error {
	return s.DB.Where("user_id = ? AND application_id = ?", userID, applicationID).
		Delete(&UserApplicationRole{}).Error
}

func (s *Service) ListByUser(userID uint) ([]UserApplicationRole, error) {
	var assignments []UserApplicationRole
	if err := s.DB.Where("user_id = ?", userID).Order("created_at").Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (s *Service) ListApplicationsByUser(userID uint) ([]application.Application, error) {
	var apps []application.Application
	err := s.DB.Distinct("applications.*").
		Joins("JOIN user_application_role ON user_application_role.application_id = applications.id").
		Where("user_application_role.user_id = ? AND applications.active = ?", userID, true).
		Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *Service) ListUsersByApplication(applicationID uint) ([]user.User, error) {
	var users []user.User
	err := s.DB.Distinct("users.*").
		Joins("JOIN user_application_role ON user_application_role.user_id = users.id").
		Where("user_application_role.application_id = ?", applicationID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) ListUsersByApplicationRole(applicationID, roleID uint) ([]user.User, error) {
	var users []user.User
	err := s.DB.Distinct("users.*").
		Joins("JOIN user_application_role ON user_application_role.user_id = users.id").
		Where("user_application_role.application_id = ? AND user_application_role.role_id = ?", applicationID, roleID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) ListRolesForUserInApplication(userID, applicationID uint) ([]role.Role, error) {
	var roles []role.Role
	err := s.DB.Distinct("roles.*").
		Joins("JOIN user_application_role ON user_application_role.role_id = roles.id").
		Where("user_application_role.user_id = ? AND user_application_role.application_id = ?", userID, applicationID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetUserAccess returns, for every application the user has a role in, the
// distinct active role slugs and active permission slugs granted there.
// This is what gets embedded in the JWT and returned from login/me.
func (s *Service) GetUserAccess(userID uint) ([]UserApplicationAccess, error) {
	apps, err := s.ListApplicationsByUser(userID)
	if err != nil {
		return nil, err
	}

	access := make([]UserApplicationAccess, 0, len(apps))
	for _, app := range apps {
		var roleSlugs []string
		err := s.DB.Table("roles").
			Distinct("roles.slug").
			Joins("JOIN user_application_role ON user_application_role.role_id = roles.id").
			Where("user_application_role.user_id = ? AND user_application_role.application_id = ? AND roles.active = ?", userID, app.ID, true).
			Pluck("roles.slug", &roleSlugs).Error
		if err != nil {
			return nil, err
		}

		var permissionSlugs []string
		err = s.DB.Table("permissions").
			Distinct("permissions.slug").
			Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
			Joins("JOIN user_application_role ON user_application_role.role_id = role_permissions.role_id").
			Where("user_application_role.user_id = ? AND user_application_role.application_id = ? AND permissions.active = ?", userID, app.ID, true).
			Pluck("permissions.slug", &permissionSlugs).Error
		if err != nil {
			return nil, err
		}

		access = append(access, UserApplicationAccess{
			Application: app,
			Roles:       roleSlugs,
			Permissions: permissionSlugs,
		})
	}

	return access, nil
}
