package role

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"go_auth/internal/permission"
)

func newTestService(t *testing.T) (*Service, *permission.Service) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&Role{}, &Permission{}, &permission.Permission{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return NewService(db), permission.NewService(db)
}

func TestCreate(t *testing.T) {
	s, _ := newTestService(t)

	r, err := s.Create(1, "admin", "Administrator", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Active {
		t.Fatal("expected new role to be active by default")
	}
}

func TestCreateDuplicateSlugSameApplication(t *testing.T) {
	s, _ := newTestService(t)

	if _, err := s.Create(1, "admin", "Administrator", "desc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := s.Create(1, "admin", "Administrator 2", "desc")
	if !errors.Is(err, ErrSlugTaken) {
		t.Fatalf("expected ErrSlugTaken, got %v", err)
	}
}

func TestAssignPermission(t *testing.T) {
	s, permSvc := newTestService(t)

	r, err := s.Create(1, "admin", "Administrator", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, err := permSvc.Create(1, "Read Users", "read_users", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := s.AssignPermission(r.ID, p.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	perms, err := s.ListPermissions(r.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(perms) != 1 || perms[0].Slug != "read_users" {
		t.Fatalf("expected [read_users], got %+v", perms)
	}
}

func TestAssignPermissionAlreadyAssigned(t *testing.T) {
	s, permSvc := newTestService(t)

	r, _ := s.Create(1, "admin", "Administrator", "desc")
	p, _ := permSvc.Create(1, "Read Users", "read_users", "desc")

	if err := s.AssignPermission(r.ID, p.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := s.AssignPermission(r.ID, p.ID)
	if !errors.Is(err, ErrAlreadyAssigned) {
		t.Fatalf("expected ErrAlreadyAssigned, got %v", err)
	}
}

func TestAssignPermissionApplicationMismatch(t *testing.T) {
	s, permSvc := newTestService(t)

	r, _ := s.Create(1, "admin", "Administrator", "desc")
	p, _ := permSvc.Create(2, "Read Users", "read_users", "desc")

	err := s.AssignPermission(r.ID, p.ID)
	if !errors.Is(err, ErrApplicationMismatch) {
		t.Fatalf("expected ErrApplicationMismatch, got %v", err)
	}
}

func TestRemovePermission(t *testing.T) {
	s, permSvc := newTestService(t)

	r, _ := s.Create(1, "admin", "Administrator", "desc")
	p, _ := permSvc.Create(1, "Read Users", "read_users", "desc")
	if err := s.AssignPermission(r.ID, p.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := s.RemovePermission(r.ID, p.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	perms, err := s.ListPermissions(r.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(perms) != 0 {
		t.Fatalf("expected no permissions after removal, got %+v", perms)
	}
}
