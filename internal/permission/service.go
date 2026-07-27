package permission

import (
	"errors"

	"gorm.io/gorm"
)

var ErrSlugTaken = errors.New("permission slug already exists for this application")

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) Create(applicationID uint, name, slug, description string) (*Permission, error) {
	exists, err := s.SlugExists(applicationID, slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlugTaken
	}

	p := &Permission{
		ApplicationID: applicationID,
		Name:          name,
		Slug:          slug,
		Description:   description,
		Active:        true,
	}

	if err := s.DB.Create(p).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlugTaken
		}
		return nil, err
	}

	return p, nil
}

func (s *Service) List() ([]Permission, error) {
	var perms []Permission
	if err := s.DB.Order("name").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (s *Service) ListByApplication(applicationID uint) ([]Permission, error) {
	var perms []Permission
	if err := s.DB.Where("application_id = ?", applicationID).Order("name").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (s *Service) GetByID(id uint) (*Permission, error) {
	var p Permission
	if err := s.DB.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Service) Update(id uint, name, slug, description string, active bool) error {
	return s.DB.Model(&Permission{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":        name,
		"slug":        slug,
		"description": description,
		"active":      active,
	}).Error
}

func (s *Service) Delete(id uint) error {
	return s.DB.Delete(&Permission{}, id).Error
}

func (s *Service) SlugExists(applicationID uint, slug string) (bool, error) {
	var count int64
	if err := s.DB.Model(&Permission{}).Where("application_id = ? AND slug = ?", applicationID, slug).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
