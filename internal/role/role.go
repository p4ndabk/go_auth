package role

import (
	"database/sql"
	"time"
	"fmt"

	"github.com/google/uuid"
	"go_auth/pkg/logs"
	"go_auth/internal/application"
	"go_auth/internal/permission"
)

type Role struct {
	ID            int               `json:"id"`
	ApplicationID int               `json:"application_id"`
	UUID          string            `json:"uuid"`
	Slug          string            `json:"slug"`
	Name          string            `json:"name"`
	Description   *string           `json:"description,omitempty"`
	Active        bool              `json:"active"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`

	Application   *application.Application `json:"application,omitempty"`
	Permissions   []*permission.Permission  `json:"permissions,omitempty"`

}

func CreateRole(db *sql.DB, role *Role) error {
	role.UUID = uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	role.CreatedAt = now
	role.UpdatedAt = now

	stmt, err := db.Prepare(`
		INSERT INTO roles (application_id, uuid, slug, name, description, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		logs.Error("create role prepare statement", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		role.ApplicationID,
		role.UUID,
		role.Slug,
		role.Name,
		role.Description,
		role.Active,
		role.CreatedAt,
		role.UpdatedAt,
	)
	if err != nil {
		logs.Error("create role exec statement", err)
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		logs.Error("create role last insert id", err)
		return err
	}

	role.ID = int(lastID)
	return nil
}

func GetRoleByID(db *sql.DB, id int) (*Role, error) {
    // Busca role + application
	row := db.QueryRow(`
		SELECT 
			r.id, r.application_id, r.uuid, r.slug, r.name, r.description, r.active, r.created_at, r.updated_at,
			a.id, a.uuid, a.slug, a.name, a.description, a.active, a.created_at, a.updated_at
		FROM roles r
		JOIN applications a ON r.application_id = a.id
		WHERE r.id = ?
	`, id)

	var r Role
	var a application.Application

	err := row.Scan(
		&r.ID,
		&r.ApplicationID,
		&r.UUID,
		&r.Slug,
		&r.Name,
		&r.Description,
		&r.Active,
		&r.CreatedAt,
		&r.UpdatedAt,

		&a.ID,
		&a.UUID,
		&a.Slug,
		&a.Name,
		&a.Description,
		&a.Active,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		logs.Error("get role by id", err)
		return nil, err
	}
	r.Application = &a

	rows, err := db.Query(`
		SELECT 
			p.id, p.application_id, p.name, p.slug, p.description, p.active, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permission rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
	`, r.ID)
	if err != nil {
		logs.Error("get permissions by role id", err)
		return nil, err
	}
	defer rows.Close()

	var permissions []*permission.Permission
	for rows.Next() {
		var p permission.Permission
		err := rows.Scan(
			&p.ID,
			&p.ApplicationID,
			&p.Name,
			&p.Slug,
			&p.Description,
			&p.Active,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			logs.Error("scan permission", err)
			return nil, err
		}
		permissions = append(permissions, &p)
	}
	r.Permissions = permissions

	return &r, nil
}


func GetAllRoles(db *sql.DB) ([]*Role, error) {
	rows, err := db.Query(`
		SELECT 
			r.id, r.application_id, r.uuid, r.slug, r.name, r.description, r.active, r.created_at, r.updated_at,
			a.id, a.uuid, a.slug, a.name, a.description, a.active, a.created_at, a.updated_at
		FROM roles r
		JOIN applications a ON r.application_id = a.id
	`)
	if err != nil {
		logs.Error("get all roles", err)
		return nil, err
	}
	defer rows.Close()

	var roles []*Role

	for rows.Next() {
		var r Role
		var a application.Application

		err := rows.Scan(
			&r.ID,
			&r.ApplicationID,
			&r.UUID,
			&r.Slug,
			&r.Name,
			r.Description,
			&r.Active,
			&r.CreatedAt,
			&r.UpdatedAt,

			&a.ID,
			&a.UUID,
			&a.Slug,
			&a.Name,
			a.Description,
			a.Active,
			a.CreatedAt,
			a.UpdatedAt,
		)
		if err != nil {
			logs.Error("get all roles scan", err)
			return nil, err
		}

		r.Application = &a
		roles = append(roles, &r)
	}

	return roles, nil
}

func UpdateRole(db *sql.DB, role *Role) error {
	role.UpdatedAt = time.Now().Format(time.RFC3339)

	stmt, err := db.Prepare(`
		UPDATE roles
		SET slug = ?, name = ?, description = ?, active = ?, updated_at = ?
		WHERE id = ?
	`)
	if err != nil {
		logs.Error("update role prepare statement", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		role.Slug,
		role.Name,
		role.Description,
		role.Active,
		role.UpdatedAt,
		role.ID,
	)
	if err != nil {
		logs.Error("update role exec statement", err)
		return err
	}

	return nil
}

func DeleteRole(db *sql.DB, id int) error {
	stmt, err := db.Prepare(`
		DELETE FROM roles
		WHERE id = ?
	`)
	if err != nil {
		logs.Error("delete role prepare statement", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		logs.Error("delete role exec statement", err)
		return err
	}

	return nil
}
func GetApplicationIDByRoleID(db *sql.DB, roleID int) (int, error) {
    var appID int
    err := db.QueryRow("SELECT application_id FROM roles WHERE id = ?", roleID).Scan(&appID)
    if err != nil {
        return 0, err
    }
    return appID, nil
}

func AttachPermissions(db *sql.DB, roleID int, permissionIDs []int) error {
    tx, err := db.Begin()
    if err != nil {
        logs.Error("begin tx attach permissions", err)
        return err
    }

    _, err = tx.Exec("DELETE FROM role_permission WHERE role_id = ?", roleID)
    if err != nil {
        tx.Rollback()
        logs.Error("delete role_permission", err)
        return err
    }

    if len(permissionIDs) == 0 {
        return tx.Commit()
    }

    applicationID, err := GetApplicationIDByRoleID(db, roleID)
    if err != nil {
        tx.Rollback()
        logs.Error("get application_id by role_id", err)
        return err
    }

    stmt, err := tx.Prepare("INSERT INTO role_permission (application_id, role_id, permission_id) VALUES (?, ?, ?)")
    if err != nil {
        tx.Rollback()
        logs.Error("prepare insert role_permission", err)
        return err
    }
    defer stmt.Close()

    for _, pid := range permissionIDs {
        _, err := stmt.Exec(applicationID, roleID, pid)
        if err != nil {
            tx.Rollback()
            logs.Error(fmt.Sprintf("insert role_permission role_id=%d permission_id=%d", roleID, pid), err)
            return err
        }
    }

    if err := tx.Commit(); err != nil {
        logs.Error("commit attach permissions", err)
        return err
    }

    return nil
}

