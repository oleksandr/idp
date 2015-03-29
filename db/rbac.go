package db

import "gopkg.in/gorp.v1"

// Permission table
type Permission struct {
	PK             int64  `db:"permission_id"`
	Name           string `db:"name"`
	Description    string `db:"description"`
	EvaluationRule string `db:"evaluation_rule"`
	Enabled        bool   `db:"is_enabled"`
}

// Role table
type Role struct {
	PK          int64  `db:"role_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Enabled     bool   `db:"is_enabled"`
}

// RolePermission table
type RolePermission struct {
	RolePK       int64 `db:"role_id"`
	PermissionPK int64 `db:"permission_id"`
}

// DeleteRole deletes a role
func DeleteRole(dbmap *gorp.DbMap, name string) error {
	var r Role
	err := dbmap.SelectOne(&r, "SELECT * FROM role WHERE name = ?", name)
	if err != nil {
		return err
	}

	tx, err := dbmap.Begin()
	if err != nil {
		return nil
	}

	_, err = tx.Exec("DELETE FROM user_role WHERE role_id = ?;", r.PK)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM role_permission WHERE role_id = ?;", r.PK)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Delete(&r)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
