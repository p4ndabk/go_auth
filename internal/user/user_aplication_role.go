package user

import (
    "time"

	"database/sql"
)

type UserApplicationRole struct {
    ID            int       `json:"id"`
    UserID        int       `json:"user_id"`
    ApplicationID int       `json:"application_id"`
    RoleID        int       `json:"role_id"`
    CreatedAt     time.Time `json:"created_at,omitempty"` // opcional, se usar timestamp
    UpdatedAt     time.Time `json:"updated_at,omitempty"` // idem
}

func CreateUserApplicationRole(db *sql.DB, uar *UserApplicationRole) error {
    _, err := db.Exec(`
        INSERT INTO user_application_role (user_id, application_id, role_id) 
        VALUES (?, ?, ?)`, uar.UserID, uar.ApplicationID, uar.RoleID)
    return err
}

func GetUserApplicationRolesByUserID(db *sql.DB, userID int) ([]*UserApplicationRole, error) {
    rows, err := db.Query(`
        SELECT id, user_id, application_id, role_id FROM user_application_role WHERE user_id = ?`, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []*UserApplicationRole
    for rows.Next() {
        var uar UserApplicationRole
        if err := rows.Scan(&uar.ID, &uar.UserID, &uar.ApplicationID, &uar.RoleID); err != nil {
            return nil, err
        }
        results = append(results, &uar)
    }
    return results, nil
}
