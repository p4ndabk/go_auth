package user

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newTestService(t *testing.T) *Service {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return NewService(db)
}

func TestCreate(t *testing.T) {
	s := newTestService(t)

	u, err := s.Create("alice", "alice@example.com", "hashed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID == 0 {
		t.Fatal("expected user to have an ID assigned")
	}
	if u.UUID == "" {
		t.Fatal("expected user to have a UUID assigned")
	}
}

func TestCreateDuplicateEmail(t *testing.T) {
	s := newTestService(t)

	if _, err := s.Create("alice", "alice@example.com", "hashed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := s.Create("alice2", "alice@example.com", "hashed")
	if !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestCreateDuplicateUsername(t *testing.T) {
	s := newTestService(t)

	if _, err := s.Create("alice", "alice@example.com", "hashed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := s.Create("alice", "alice2@example.com", "hashed")
	if !errors.Is(err, ErrUsernameTaken) {
		t.Fatalf("expected ErrUsernameTaken, got %v", err)
	}
}

func TestGetByEmailNotFound(t *testing.T) {
	s := newTestService(t)

	_, err := s.GetByEmail("missing@example.com")
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	s := newTestService(t)

	if _, err := s.Create("alice", "alice@example.com", "hashed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := s.Create("bob", "bob@example.com", "hashed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	users, err := s.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}
