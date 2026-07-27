package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrEmailTaken    = errors.New("email already registered")
	ErrUsernameTaken = errors.New("username already taken")
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// Create checks email/username uniqueness up front (for a precise error
// message) and also relies on the DB's unique indexes as a safety net
// against a race between the check and the insert.
func (s *Service) Create(username, email, hashedPassword string) (*User, error) {
	exists, err := s.EmailExists(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailTaken
	}

	exists, err = s.UsernameExists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameTaken
	}

	u := &User{
		UUID:      uuid.New().String(),
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	if err := s.DB.Create(u).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}

	return u, nil
}

func (s *Service) GetByEmail(email string) (*User, error) {
	var u User
	if err := s.DB.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) GetByID(id uint) (*User, error) {
	var u User
	if err := s.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) List() ([]User, error) {
	var users []User
	if err := s.DB.Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) EmailExists(email string) (bool, error) {
	var count int64
	if err := s.DB.Model(&User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Service) UsernameExists(username string) (bool, error) {
	var count int64
	if err := s.DB.Model(&User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
