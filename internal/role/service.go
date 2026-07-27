package role

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go_auth/internal/permission"
)

var (
	ErrSlugTaken           = errors.New("role slug already exists for this application")
	ErrAlreadyAssigned     = errors.New("permission already assigned to role")
	ErrApplicationMismatch = errors.New("role and permission must belong to the same application")
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) Create(applicationID uint, slug, name, description string) (*Role, error) {
	exists, err := s.SlugExists(applicationID, slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlugTaken
	}

	r := &Role{
		ApplicationID: applicationID,
		UUID:          uuid.New().String(),
		Slug:          slug,
		Name:          name,
		Description:   description,
		Active:        true,
	}

	if err := s.DB.Create(r).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlugTaken
		}
		return nil, err
	}

	return r, nil
}

func (s *Service) List() ([]Role, error) {
	var roles []Role
	if err := s.DB.Order("name").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *Service) ListByApplication(applicationID uint) ([]Role, error) {
	var roles []Role
	if err := s.DB.Where("application_id = ?", applicationID).Order("name").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *Service) GetByID(id uint) (*Role, error) {
	var r Role
	if err := s.DB.First(&r, id).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Service) Update(id uint, slug, name, description string, active bool) error {
	return s.DB.Model(&Role{}).Where("id = ?", id).Updates(map[string]interface{}{
		"slug":        slug,
		"name":        name,
		"description": description,
		"active":      active,
	}).Error
}

func (s *Service) Delete(id uint) error {
	return s.DB.Delete(&Role{}, id).Error
}

func (s *Service) SlugExists(applicationID uint, slug string) (bool, error) {
	var count int64
	if err := s.DB.Model(&Role{}).Where("application_id = ? AND slug = ?", applicationID, slug).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// AssignPermission grants permissionID to roleID. Both must belong to the
// same application.
func (s *Service) AssignPermission(roleID, permissionID uint) error {
	var r Role
	if err := s.DB.First(&r, roleID).Error; err != nil {
		return err
	}

	var p permission.Permission
	if err := s.DB.First(&p, permissionID).Error; err != nil {
		return err
	}

	if r.ApplicationID != p.ApplicationID {
		return ErrApplicationMismatch
	}

	var count int64
	if err := s.DB.Model(&Permission{}).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrAlreadyAssigned
	}

	rp := &Permission{
		RoleID:        roleID,
		PermissionID:  permissionID,
		ApplicationID: r.ApplicationID,
		CreatedAt:     time.Now(),
	}
	if err := s.DB.Create(rp).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyAssigned
		}
		return err
	}
	return nil
}

func (s *Service) RemovePermission(roleID, permissionID uint) error {
	return s.DB.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&Permission{}).Error
}

func (s *Service) ListPermissions(roleID uint) ([]permission.Permission, error) {
	var perms []permission.Permission
	err := s.DB.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.name").
		Find(&perms).Error
	if err != nil {
		return nil, err
	}
	return perms, nil
}
