package application

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
	if err := db.AutoMigrate(&Application{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return NewService(db)
}

func TestCreate(t *testing.T) {
	s := newTestService(t)

	app, err := s.Create("main-app", "Main Application", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !app.Active {
		t.Fatal("expected new application to be active by default")
	}
}

func TestCreateDuplicateSlug(t *testing.T) {
	s := newTestService(t)

	if _, err := s.Create("main-app", "Main Application", "desc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := s.Create("main-app", "Other Name", "desc")
	if !errors.Is(err, ErrSlugTaken) {
		t.Fatalf("expected ErrSlugTaken, got %v", err)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := newTestService(t)

	_, err := s.GetByID(999)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestUpdateAndDelete(t *testing.T) {
	s := newTestService(t)

	app, err := s.Create("main-app", "Main Application", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := s.Update(app.ID, "main-app", "Renamed", "new desc", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := s.GetByID(app.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Renamed" || updated.Active {
		t.Fatalf("update did not apply: %+v", updated)
	}

	if err := s.Delete(app.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := s.GetByID(app.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record to be deleted, got %v", err)
	}
}
