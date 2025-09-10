package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	db *sql.DB
}

type User struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Role struct {
	ID            int    `json:"id"`
	ApplicationID int    `json:"application_id"`
	UUID          string `json:"uuid"`
	Slug          string `json:"slug"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Active        bool   `json:"active"`
}

type Permission struct {
	ID            int    `json:"id"`
	ApplicationID int    `json:"application_id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	Active        bool   `json:"active"`
}

type Application struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type UserApplicationRole struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	ApplicationID int       `json:"application_id"`
	RoleID        int       `json:"role_id"`
	ProfileID     *int      `json:"profile_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type RolePermission struct {
	ID            int       `json:"id"`
	RoleID        int       `json:"role_id"`
	PermissionID  int       `json:"permission_id"`
	ApplicationID int       `json:"application_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func New(database *sql.DB) *Service {
	return &Service{db: database}
}

// Ping checks if the database connection is alive
func (s *Service) Ping() error {
	return s.db.Ping()
}

func InitSchema(db *sql.DB, dbType string) error {
	// Check if users table exists
	var tableName string
	var checkQuery string

	switch dbType {
	case "mysql":
		checkQuery = "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'users'"
	case "sqlite":
		fallthrough
	default:
		checkQuery = "SELECT name FROM sqlite_master WHERE type='table' AND name='users'"
	}

	err := db.QueryRow(checkQuery).Scan(&tableName)
	if err == nil {
		// Tables already exist, skip schema creation
		return nil
	}

	// Load appropriate schema file
	var schemaFile string
	switch dbType {
	case "mysql":
		schemaFile = "migrations/001_initial_schema_mysql.sql"
	case "sqlite":
		fallthrough
	default:
		schemaFile = "migrations/001_initial_schema.sql"
	}

	schema, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %v", schemaFile, err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	return nil
}

func (s *Service) CreateUser(username, email, hashedPassword string) (*User, error) {
	userUUID := uuid.New().String()

	query := `INSERT INTO users (uuid, username, email, password, created_at) 
			  VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.Exec(query, userUUID, username, email, hashedPassword, time.Now())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        int(id),
		UUID:      userUUID,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}, nil
}

func (s *Service) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, uuid, username, email, password, created_at 
			  FROM users WHERE email = ?`

	user := &User{}
	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.UUID, &user.Username,
		&user.Email, &user.Password, &user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserByID(userID int) (*User, error) {
	query := `SELECT id, uuid, username, email, password, created_at 
			  FROM users WHERE id = ?`

	user := &User{}
	err := s.db.QueryRow(query, userID).Scan(
		&user.ID, &user.UUID, &user.Username,
		&user.Email, &user.Password, &user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UserApplicationAccess represents user's access within a specific application
type UserApplicationAccess struct {
	Application Application `json:"application"`
	Roles       []string    `json:"roles"`
	Permissions []string    `json:"permissions"`
}

func (s *Service) GetUserRoles(userID int) ([]string, error) {
	query := `SELECT DISTINCT r.slug FROM roles r
			  INNER JOIN user_application_role uar ON r.id = uar.role_id
			  WHERE uar.user_id = ? AND r.active = 1`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (s *Service) GetUserPermissions(userID int) ([]string, error) {
	query := `SELECT DISTINCT p.slug FROM permissions p
			  INNER JOIN role_permissions rp ON p.id = rp.permission_id
			  INNER JOIN user_application_role uar ON rp.role_id = uar.role_id
			  WHERE uar.user_id = ? AND p.active = 1`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// GetUserAccessByApplications returns user's roles and permissions organized by application
func (s *Service) GetUserAccessByApplications(userID int) ([]UserApplicationAccess, error) {
	// First get all applications the user has access to
	applications, err := s.GetUserApplications(userID)
	if err != nil {
		return nil, err
	}

	var userAccess []UserApplicationAccess
	for _, app := range applications {
		// Get roles for this user in this application
		rolesQuery := `SELECT DISTINCT r.slug FROM roles r
					   INNER JOIN user_application_role uar ON r.id = uar.role_id
					   WHERE uar.user_id = ? AND uar.application_id = ? AND r.active = 1`
		
		roleRows, err := s.db.Query(rolesQuery, userID, app.ID)
		if err != nil {
			return nil, err
		}

		var roles []string
		for roleRows.Next() {
			var role string
			if err := roleRows.Scan(&role); err != nil {
				roleRows.Close()
				return nil, err
			}
			roles = append(roles, role)
		}
		roleRows.Close()

		// Get permissions for this user in this application
		permissionsQuery := `SELECT DISTINCT p.slug FROM permissions p
							 INNER JOIN role_permissions rp ON p.id = rp.permission_id
							 INNER JOIN user_application_role uar ON rp.role_id = uar.role_id
							 WHERE uar.user_id = ? AND uar.application_id = ? AND p.active = 1`
		
		permRows, err := s.db.Query(permissionsQuery, userID, app.ID)
		if err != nil {
			return nil, err
		}

		var permissions []string
		for permRows.Next() {
			var permission string
			if err := permRows.Scan(&permission); err != nil {
				permRows.Close()
				return nil, err
			}
			permissions = append(permissions, permission)
		}
		permRows.Close()

		userAccess = append(userAccess, UserApplicationAccess{
			Application: app,
			Roles:       roles,
			Permissions: permissions,
		})
	}

	return userAccess, nil
}

func (s *Service) EmailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	var count int
	err := s.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *Service) UsernameExists(username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = ?`

	var count int
	err := s.db.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Application CRUD operations
func (s *Service) CreateApplication(slug, name, description string) (*Application, error) {
	appUUID := uuid.New().String()

	query := `INSERT INTO applications (uuid, slug, name, description, active) 
			  VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.Exec(query, appUUID, slug, name, description, true)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Application{
		ID:          int(id),
		UUID:        appUUID,
		Slug:        slug,
		Name:        name,
		Description: description,
		Active:      true,
	}, nil
}

func (s *Service) GetApplications() ([]Application, error) {
	query := `SELECT id, uuid, slug, name, description, active FROM applications ORDER BY name`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []Application
	for rows.Next() {
		var app Application
		if err := rows.Scan(&app.ID, &app.UUID, &app.Slug, &app.Name, &app.Description, &app.Active); err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}

	return applications, nil
}

func (s *Service) GetApplicationByID(id int) (*Application, error) {
	query := `SELECT id, uuid, slug, name, description, active FROM applications WHERE id = ?`

	app := &Application{}
	err := s.db.QueryRow(query, id).Scan(
		&app.ID, &app.UUID, &app.Slug, &app.Name, &app.Description, &app.Active,
	)

	if err != nil {
		return nil, err
	}

	return app, nil
}

func (s *Service) UpdateApplication(id int, slug, name, description string, active bool) error {
	query := `UPDATE applications SET slug = ?, name = ?, description = ?, active = ? WHERE id = ?`

	_, err := s.db.Exec(query, slug, name, description, active, id)
	return err
}

func (s *Service) DeleteApplication(id int) error {
	query := `DELETE FROM applications WHERE id = ?`

	_, err := s.db.Exec(query, id)
	return err
}

func (s *Service) ApplicationSlugExists(slug string) (bool, error) {
	query := `SELECT COUNT(*) FROM applications WHERE slug = ?`

	var count int
	err := s.db.QueryRow(query, slug).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Role CRUD operations
func (s *Service) CreateRole(applicationID int, slug, name, description string) (*Role, error) {
	roleUUID := uuid.New().String()

	query := `INSERT INTO roles (application_id, uuid, slug, name, description, active) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	result, err := s.db.Exec(query, applicationID, roleUUID, slug, name, description, true)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Role{
		ID:            int(id),
		ApplicationID: applicationID,
		UUID:          roleUUID,
		Slug:          slug,
		Name:          name,
		Description:   description,
		Active:        true,
	}, nil
}

func (s *Service) GetRoles() ([]Role, error) {
	query := `SELECT id, application_id, uuid, slug, name, description, active FROM roles ORDER BY name`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.ApplicationID, &role.UUID, &role.Slug, &role.Name, &role.Description, &role.Active); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (s *Service) GetRolesByApplication(applicationID int) ([]Role, error) {
	query := `SELECT id, application_id, uuid, slug, name, description, active FROM roles WHERE application_id = ? ORDER BY name`

	rows, err := s.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.ApplicationID, &role.UUID, &role.Slug, &role.Name, &role.Description, &role.Active); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (s *Service) GetRoleByID(id int) (*Role, error) {
	query := `SELECT id, application_id, uuid, slug, name, description, active FROM roles WHERE id = ?`

	role := &Role{}
	err := s.db.QueryRow(query, id).Scan(
		&role.ID, &role.ApplicationID, &role.UUID, &role.Slug, &role.Name, &role.Description, &role.Active,
	)

	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *Service) UpdateRole(id int, slug, name, description string, active bool) error {
	query := `UPDATE roles SET slug = ?, name = ?, description = ?, active = ? WHERE id = ?`

	_, err := s.db.Exec(query, slug, name, description, active, id)
	return err
}

func (s *Service) DeleteRole(id int) error {
	query := `DELETE FROM roles WHERE id = ?`

	_, err := s.db.Exec(query, id)
	return err
}

// Permission CRUD operations
func (s *Service) CreatePermission(applicationID int, name, slug, description string) (*Permission, error) {
	query := `INSERT INTO permissions (application_id, name, slug, description, active) 
			  VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.Exec(query, applicationID, name, slug, description, true)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Permission{
		ID:            int(id),
		ApplicationID: applicationID,
		Name:          name,
		Slug:          slug,
		Description:   description,
		Active:        true,
	}, nil
}

func (s *Service) GetPermissions() ([]Permission, error) {
	query := `SELECT id, application_id, name, slug, description, active FROM permissions ORDER BY name`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var perm Permission
		if err := rows.Scan(&perm.ID, &perm.ApplicationID, &perm.Name, &perm.Slug, &perm.Description, &perm.Active); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (s *Service) GetPermissionsByApplication(applicationID int) ([]Permission, error) {
	query := `SELECT id, application_id, name, slug, description, active FROM permissions WHERE application_id = ? ORDER BY name`

	rows, err := s.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var perm Permission
		if err := rows.Scan(&perm.ID, &perm.ApplicationID, &perm.Name, &perm.Slug, &perm.Description, &perm.Active); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (s *Service) GetPermissionByID(id int) (*Permission, error) {
	query := `SELECT id, application_id, name, slug, description, active FROM permissions WHERE id = ?`

	perm := &Permission{}
	err := s.db.QueryRow(query, id).Scan(
		&perm.ID, &perm.ApplicationID, &perm.Name, &perm.Slug, &perm.Description, &perm.Active,
	)

	if err != nil {
		return nil, err
	}

	return perm, nil
}

func (s *Service) UpdatePermission(id int, name, slug, description string, active bool) error {
	query := `UPDATE permissions SET name = ?, slug = ?, description = ?, active = ? WHERE id = ?`

	_, err := s.db.Exec(query, name, slug, description, active, id)
	return err
}

func (s *Service) DeletePermission(id int) error {
	query := `DELETE FROM permissions WHERE id = ?`

	_, err := s.db.Exec(query, id)
	return err
}

// Role-Permission relationship operations
func (s *Service) AssignPermissionToRole(roleID, permissionID int) error {
	// Get application_id from role
	var roleAppID int
	roleQuery := `SELECT application_id FROM roles WHERE id = ?`
	err := s.db.QueryRow(roleQuery, roleID).Scan(&roleAppID)
	if err != nil {
		return fmt.Errorf("role not found: %v", err)
	}

	// Get application_id from permission
	var permAppID int
	permQuery := `SELECT application_id FROM permissions WHERE id = ?`
	err = s.db.QueryRow(permQuery, permissionID).Scan(&permAppID)
	if err != nil {
		return fmt.Errorf("permission not found: %v", err)
	}

	// Validate that role and permission belong to the same application
	if roleAppID != permAppID {
		return fmt.Errorf("role and permission must belong to the same application")
	}

	// Check if association already exists
	checkQuery := `SELECT COUNT(*) FROM role_permissions WHERE role_id = ? AND permission_id = ? AND application_id = ?`
	var count int
	err = s.db.QueryRow(checkQuery, roleID, permissionID, roleAppID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("permission already assigned to role")
	}

	query := `INSERT INTO role_permissions (role_id, permission_id, application_id, created_at) VALUES (?, ?, ?, ?)`
	_, err = s.db.Exec(query, roleID, permissionID, roleAppID, time.Now())
	return err
}

func (s *Service) RemovePermissionFromRole(roleID, permissionID int) error {
	query := `DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?`
	_, err := s.db.Exec(query, roleID, permissionID)
	return err
}

func (s *Service) GetRolePermissions(roleID int) ([]Permission, error) {
	query := `SELECT p.id, p.application_id, p.name, p.slug, p.description, p.active 
			  FROM permissions p
			  INNER JOIN role_permissions rp ON p.id = rp.permission_id
			  WHERE rp.role_id = ? ORDER BY p.name`

	rows, err := s.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var perm Permission
		if err := rows.Scan(&perm.ID, &perm.ApplicationID, &perm.Name, &perm.Slug, &perm.Description, &perm.Active); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// Get RolePermission records with full details
func (s *Service) GetRolePermissionDetails(roleID int) ([]RolePermission, error) {
	query := `SELECT rp.id, rp.role_id, rp.permission_id, rp.application_id, rp.created_at
			  FROM role_permissions rp
			  WHERE rp.role_id = ? ORDER BY rp.created_at`

	rows, err := s.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rolePermissions []RolePermission
	for rows.Next() {
		var rp RolePermission
		err := rows.Scan(&rp.ID, &rp.RoleID, &rp.PermissionID, &rp.ApplicationID, &rp.CreatedAt)
		if err != nil {
			return nil, err
		}
		rolePermissions = append(rolePermissions, rp)
	}

	return rolePermissions, nil
}

// Get permissions for a role within a specific application
func (s *Service) GetRolePermissionsByApplication(roleID, applicationID int) ([]Permission, error) {
	query := `SELECT p.id, p.application_id, p.name, p.slug, p.description, p.active 
			  FROM permissions p
			  INNER JOIN role_permissions rp ON p.id = rp.permission_id
			  WHERE rp.role_id = ? AND rp.application_id = ? ORDER BY p.name`

	rows, err := s.db.Query(query, roleID, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var perm Permission
		err := rows.Scan(&perm.ID, &perm.ApplicationID, &perm.Name, &perm.Slug, &perm.Description, &perm.Active)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// User-Role relationship operations
func (s *Service) AssignRoleToUser(userID, roleID int) error {
	// Check if association already exists
	checkQuery := `SELECT COUNT(*) FROM user_roles WHERE user_id = ? AND role_id = ?`
	var count int
	err := s.db.QueryRow(checkQuery, userID, roleID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("role already assigned to user")
	}

	query := `INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`
	_, err = s.db.Exec(query, userID, roleID)
	return err
}

func (s *Service) RemoveRoleFromUser(userID, roleID int) error {
	query := `DELETE FROM user_roles WHERE user_id = ? AND role_id = ?`
	_, err := s.db.Exec(query, userID, roleID)
	return err
}

// User-Application relationship operations
// DEPRECATED: Use user_application_role instead for better context
func (s *Service) AssignUserToApplication(userID, applicationID int) error {
	// This function is deprecated but kept for backward compatibility
	// Consider using AssignUserApplicationRole instead

	// Check if association already exists in new table
	checkQuery := `SELECT COUNT(*) FROM user_application_role WHERE user_id = ? AND application_id = ?`
	var count int
	err := s.db.QueryRow(checkQuery, userID, applicationID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("user already assigned to application")
	}

	// For now, we'll create a basic entry without specific role
	// In a real scenario, you should assign a specific role
	return fmt.Errorf("deprecated: use AssignUserApplicationRole instead")
}

// DEPRECATED: Use RemoveUserApplicationRole instead for better control
func (s *Service) RemoveUserFromApplication(userID, applicationID int) error {
	// This function is deprecated but kept for backward compatibility
	// Remove ALL roles for user in this application
	query := `DELETE FROM user_application_role WHERE user_id = ? AND application_id = ?`
	_, err := s.db.Exec(query, userID, applicationID)
	return err
}

func (s *Service) GetUserApplications(userID int) ([]Application, error) {
	query := `
		SELECT DISTINCT a.id, a.uuid, a.slug, a.name, a.description, a.active
		FROM applications a
		INNER JOIN user_application_role uar ON a.id = uar.application_id
		WHERE uar.user_id = ? AND a.active = 1
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []Application
	for rows.Next() {
		var app Application
		err := rows.Scan(&app.ID, &app.UUID, &app.Slug, &app.Name, &app.Description, &app.Active)
		if err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}

	return applications, nil
}

func (s *Service) GetApplicationUsers(applicationID int) ([]User, error) {
	query := `
		SELECT DISTINCT u.id, u.uuid, u.username, u.email, u.created_at
		FROM users u
		INNER JOIN user_application_role uar ON u.id = uar.user_id
		WHERE uar.application_id = ?
	`

	rows, err := s.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.UUID, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// User-Application-Role functions
func (s *Service) AssignUserApplicationRole(userID, applicationID, roleID int, profileID *int) error {
	// Check if assignment already exists
	checkQuery := `SELECT COUNT(*) FROM user_application_role WHERE user_id = ? AND application_id = ? AND role_id = ?`
	var count int
	err := s.db.QueryRow(checkQuery, userID, applicationID, roleID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("user already has this role in this application")
	}

	query := `INSERT INTO user_application_role (user_id, application_id, role_id, profile_id, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = s.db.Exec(query, userID, applicationID, roleID, profileID, time.Now())
	return err
}

func (s *Service) RemoveUserApplicationRole(userID, applicationID, roleID int) error {
	query := `DELETE FROM user_application_role WHERE user_id = ? AND application_id = ? AND role_id = ?`
	_, err := s.db.Exec(query, userID, applicationID, roleID)
	return err
}

func (s *Service) GetUserApplicationRoles(userID int) ([]UserApplicationRole, error) {
	query := `SELECT id, user_id, application_id, role_id, profile_id, created_at FROM user_application_role WHERE user_id = ?`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []UserApplicationRole
	for rows.Next() {
		var assignment UserApplicationRole
		err := rows.Scan(&assignment.ID, &assignment.UserID, &assignment.ApplicationID, &assignment.RoleID, &assignment.ProfileID, &assignment.CreatedAt)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

func (s *Service) GetApplicationRoleUsers(applicationID, roleID int) ([]User, error) {
	query := `
		SELECT DISTINCT u.id, u.uuid, u.username, u.email, u.created_at 
		FROM users u
		INNER JOIN user_application_role uar ON u.id = uar.user_id
		WHERE uar.application_id = ? AND uar.role_id = ?`

	rows, err := s.db.Query(query, applicationID, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.UUID, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *Service) GetUserRolesInApplication(userID, applicationID int) ([]Role, error) {
	query := `
		SELECT DISTINCT r.id, r.application_id, r.uuid, r.slug, r.name, r.description, r.active
		FROM roles r
		INNER JOIN user_application_role uar ON r.id = uar.role_id
		WHERE uar.user_id = ? AND uar.application_id = ?`

	rows, err := s.db.Query(query, userID, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.ID, &role.ApplicationID, &role.UUID, &role.Slug, &role.Name, &role.Description, &role.Active)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}
