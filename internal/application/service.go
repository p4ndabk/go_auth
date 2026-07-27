package application

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrSlugTaken = errors.New("application slug already exists")

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) Create(slug, name, description string) (*Application, error) {
	exists, err := s.SlugExists(slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlugTaken
	}

	app := &Application{
		UUID:        uuid.New().String(),
		Slug:        slug,
		Name:        name,
		Description: description,
		Active:      true,
	}

	if err := s.DB.Create(app).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlugTaken
		}
		return nil, err
	}

	return app, nil
}

func (s *Service) List() ([]Application, error) {
	var apps []Application
	if err := s.DB.Order("name").Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *Service) GetByID(id uint) (*Application, error) {
	var app Application
	if err := s.DB.First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *Service) Update(id uint, slug, name, description string, active bool) error {
	return s.DB.Model(&Application{}).Where("id = ?", id).Updates(map[string]interface{}{
		"slug":        slug,
		"name":        name,
		"description": description,
		"active":      active,
	}).Error
}

func (s *Service) Delete(id uint) error {
	return s.DB.Delete(&Application{}, id).Error
}

func (s *Service) SlugExists(slug string) (bool, error) {
	var count int64
	if err := s.DB.Model(&Application{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
