package permission

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
	if err := db.AutoMigrate(&Permission{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return NewService(db)
}

func TestCreate(t *testing.T) {
	s := newTestService(t)

	p, err := s.Create(1, "Read Users", "read_users", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.Active {
		t.Fatal("expected new permission to be active by default")
	}
}

func TestCreateDuplicateSlugSameApplication(t *testing.T) {
	s := newTestService(t)

	if _, err := s.Create(1, "Read Users", "read_users", "desc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := s.Create(1, "Read Users Again", "read_users", "desc")
	if !errors.Is(err, ErrSlugTaken) {
		t.Fatalf("expected ErrSlugTaken, got %v", err)
	}
}

func TestSameSlugAllowedAcrossApplications(t *testing.T) {
	s := newTestService(t)

	if _, err := s.Create(1, "Read Users", "read_users", "desc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := s.Create(2, "Read Users", "read_users", "desc"); err != nil {
		t.Fatalf("expected slug reuse across applications to succeed, got %v", err)
	}
}

func TestListByApplication(t *testing.T) {
	s := newTestService(t)

	if _, err := s.Create(1, "Read Users", "read_users", "desc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := s.Create(2, "Read Posts", "read_posts", "desc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	perms, err := s.ListByApplication(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(perms) != 1 {
		t.Fatalf("expected 1 permission for application 1, got %d", len(perms))
	}
}
