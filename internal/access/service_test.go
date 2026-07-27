package access

import (
	"errors"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"go_auth/internal/application"
	"go_auth/internal/permission"
	"go_auth/internal/role"
	"go_auth/internal/user"
)

type testFixtures struct {
	access *Service
	users  *user.Service
	apps   *application.Service
	roles  *role.Service
	perms  *permission.Service
}

func newTestFixtures(t *testing.T) *testFixtures {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	err = db.AutoMigrate(
		&user.User{},
		&application.Application{},
		&role.Role{},
		&permission.Permission{},
		&role.Permission{},
		&UserApplicationRole{},
	)
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return &testFixtures{
		access: NewService(db),
		users:  user.NewService(db),
		apps:   application.NewService(db),
		roles:  role.NewService(db),
		perms:  permission.NewService(db),
	}
}

func TestAssignAndGetUserAccess(t *testing.T) {
	f := newTestFixtures(t)

	u, err := f.users.Create("alice", "alice@example.com", "hashed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app, err := f.apps.Create("main-app", "Main Application", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r, err := f.roles.Create(app.ID, "admin", "Administrator", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, err := f.perms.Create(app.ID, "Read Users", "read_users", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := f.roles.AssignPermission(r.ID, p.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := f.access.Assign(u.ID, app.ID, r.ID, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	access, err := f.access.GetUserAccess(u.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(access) != 1 {
		t.Fatalf("expected access to 1 application, got %d", len(access))
	}
	if access[0].Application.ID != app.ID {
		t.Fatalf("expected application %d, got %d", app.ID, access[0].Application.ID)
	}
	if len(access[0].Roles) != 1 || access[0].Roles[0] != "admin" {
		t.Fatalf("expected roles [admin], got %+v", access[0].Roles)
	}
	if len(access[0].Permissions) != 1 || access[0].Permissions[0] != "read_users" {
		t.Fatalf("expected permissions [read_users], got %+v", access[0].Permissions)
	}
}

func TestAssignAlreadyAssigned(t *testing.T) {
	f := newTestFixtures(t)

	u, _ := f.users.Create("alice", "alice@example.com", "hashed")
	app, _ := f.apps.Create("main-app", "Main Application", "desc")
	r, _ := f.roles.Create(app.ID, "admin", "Administrator", "desc")

	if err := f.access.Assign(u.ID, app.ID, r.ID, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := f.access.Assign(u.ID, app.ID, r.ID, nil)
	if !errors.Is(err, ErrAlreadyAssigned) {
		t.Fatalf("expected ErrAlreadyAssigned, got %v", err)
	}
}

func TestRemoveAllForUserInApplication(t *testing.T) {
	f := newTestFixtures(t)

	u, _ := f.users.Create("alice", "alice@example.com", "hashed")
	app, _ := f.apps.Create("main-app", "Main Application", "desc")
	r1, _ := f.roles.Create(app.ID, "admin", "Administrator", "desc")
	r2, _ := f.roles.Create(app.ID, "user", "User", "desc")

	if err := f.access.Assign(u.ID, app.ID, r1.ID, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := f.access.Assign(u.ID, app.ID, r2.ID, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := f.access.RemoveAllForUserInApplication(u.ID, app.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assignments, err := f.access.ListByUser(u.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(assignments) != 0 {
		t.Fatalf("expected no assignments left, got %+v", assignments)
	}
}

func TestListUsersByApplication(t *testing.T) {
	f := newTestFixtures(t)

	u, _ := f.users.Create("alice", "alice@example.com", "hashed")
	app, _ := f.apps.Create("main-app", "Main Application", "desc")
	r, _ := f.roles.Create(app.ID, "admin", "Administrator", "desc")

	if err := f.access.Assign(u.ID, app.ID, r.ID, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	users, err := f.access.ListUsersByApplication(app.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 || users[0].ID != u.ID {
		t.Fatalf("expected [%d], got %+v", u.ID, users)
	}
}
