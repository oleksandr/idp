package dl

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/entities"
)

// Permission DTO
type Permission struct {
	PK             int64  `db:"permission_id"`
	Name           string `db:"name"`
	Description    string `db:"description"`
	EvaluationRule string `db:"evaluation_rule"`
	Enabled        bool   `db:"is_enabled"`
}

// Role DTO
type Role struct {
	PK          int64  `db:"role_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Enabled     bool   `db:"is_enabled"`
}

// CreatePermission create a new permission and returns updated DTO with PK
func CreatePermission(db sqlx.Ext, p Permission) (*Permission, error) {
	q := "INSERT INTO permission (name, description, evaluation_rule, is_enabled) VALUES (?, ?, ?, ?);"
	r, err := db.Exec(db.Rebind(q), p.Name, p.Description, p.EvaluationRule, p.Enabled)
	if err != nil {
		return nil, err
	}
	p.PK, err = lastInsertID(db, r, "permission", "permission_id")
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// RenamePermission changes role's name
func RenamePermission(db sqlx.Ext, pk int64, newName string) error {
	_, err := db.Exec(db.Rebind("UPDATE permission SET name = ? WHERE permission_id = ?;"), newName, pk)
	return err
}

// UpdatePermission updates all properties by PK
func UpdatePermission(db sqlx.Ext, p Permission) error {
	q := `UPDATE permission SET
		name = ?,
		description = ?,
		is_enabled = ?,
		evaluation_rule = ?
		WHERE permission_id = ?;`
	_, err := db.Exec(db.Rebind(q), p.Name, p.Description, p.Enabled, p.EvaluationRule, p.PK)
	return err
}

// DeletePermission deletes a permission from database
func DeletePermission(db sqlx.Ext, pk int64) error {
	err := ExecuteTransactionally(db.(*sqlx.DB), func(ext sqlx.Ext) error {
		_, err := ext.Exec(db.Rebind("DELETE FROM role_permission WHERE permission_id = ?"), pk)
		if err != nil {
			return err
		}
		_, err = ext.Exec(db.Rebind("DELETE FROM permission WHERE permission_id = ?"), pk)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// CountPermissions returns a total count of permission records in database
func CountPermissions(db sqlx.Ext) (int64, error) {
	var count int64
	err := db.QueryRowx("SELECT count(*) FROM permission;").Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// FindPermissionByName search for a permission by given name
func FindPermissionByName(db sqlx.Ext, name string) (*Permission, error) {
	var p Permission
	err := db.QueryRowx(db.Rebind("SELECT * FROM permission WHERE name = ? LIMIT 1"), name).StructScan(&p)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindPermissionsByNames search for a permission by given set of names
func FindPermissionsByNames(db sqlx.Ext, names []string) ([]Permission, error) {
	if len(names) == 0 {
		return nil, ErrNotFound
	}

	var (
		pp      []Permission
		args    []interface{}
		holders []string
	)

	for _, n := range names {
		holders = append(holders, "?")
		args = append(args, n)
	}

	q := fmt.Sprintf("SELECT * FROM permission WHERE name IN (%s);", strings.Join(holders, ","))
	rows, err := db.Queryx(db.Rebind(q), args...)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Permission
		err = rows.StructScan(&p)
		if err != nil {
			return nil, err
		}
		pp = append(pp, p)
	}
	return pp, nil
}

// FindAllPermissions returns a page of role records
func FindAllPermissions(db sqlx.Ext, pager entities.Pager, sorter entities.Sorter) ([]*Permission, error) {
	q := fmt.Sprintf(`SELECT * FROM permission AS p
		%v %v;`, orderByClause(sorter, "p"), limitOffset(pager))
	rows, err := db.Queryx(q)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := []*Permission{}
	for rows.Next() {
		var p Permission
		err = rows.StructScan(&p)
		if err != nil {
			return nil, err
		}
		perms = append(perms, &p)
	}
	return perms, nil
}

// FindPermissionsByRole returns a page of permission records filtered by a given role
func FindPermissionsByRole(db *sqlx.DB, roleName string, pager entities.Pager, sorter entities.Sorter) ([]*Permission, error) {
	q := fmt.Sprintf(`SELECT p.* FROM role_permission AS rp
		LEFT JOIN permission AS p ON p.permission_id=rp.permission_id
		LEFT JOIN role AS r ON r.role_id=rp.role_id
		WHERE r.name=? GROUP BY p.permission_id
		%v %v;`, orderByClause(sorter, "p"), limitOffset(pager))
	rows, err := db.Queryx(db.Rebind(q), roleName)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()
	perms := []*Permission{}
	for rows.Next() {
		var p Permission
		err = rows.StructScan(&p)
		if err != nil {
			return nil, err
		}
		perms = append(perms, &p)
	}
	return perms, nil
}

// CreateRole create a new permission and returns updated DTO with PK
func CreateRole(db sqlx.Ext, r Role) (*Role, error) {
	q := "INSERT INTO role (name, description, is_enabled) VALUES (?, ?, ?);"
	res, err := db.Exec(db.Rebind(q), r.Name, r.Description, r.Enabled)
	if err != nil {
		return nil, err
	}
	r.PK, err = lastInsertID(db, res, "role", "role_id")
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// RenameRole changes role's name
func RenameRole(db sqlx.Ext, pk int64, newName string) error {
	_, err := db.Exec(db.Rebind("UPDATE role SET name = ? WHERE role_id = ?;"), newName, pk)
	return err
}

// UpdateRole updates all properties by PK
func UpdateRole(db sqlx.Ext, role Role) error {
	q := `UPDATE role SET
		name = ?,
		description = ?,
		is_enabled = ?
		WHERE role_id = ?;`
	_, err := db.Exec(db.Rebind(q), role.Name, role.Description, role.Enabled, role.PK)
	return err
}

// DeleteRole deletes a permission from database
func DeleteRole(db sqlx.Ext, pk int64) error {
	err := ExecuteTransactionally(db.(*sqlx.DB), func(ext sqlx.Ext) error {
		_, err := ext.Exec(db.Rebind("DELETE FROM user_role WHERE role_id = ?"), pk)
		if err != nil {
			return err
		}
		_, err = ext.Exec(db.Rebind("DELETE FROM role_permission WHERE role_id = ?"), pk)
		if err != nil {
			return err
		}
		_, err = ext.Exec(db.Rebind("DELETE FROM role WHERE role_id = ?"), pk)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// CountRoles returns a total count of role records in database
func CountRoles(db sqlx.Ext) (int64, error) {
	var count int64
	err := db.QueryRowx("SELECT count(*) FROM role;").Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// CountRolesByUser returns a total count of role records by user
func CountRolesByUser(db sqlx.Ext, userID string) (int64, error) {
	var count int64
	err := db.QueryRowx(db.Rebind(fmt.Sprintf(`SELECT count(*) FROM user_role WHERE user_id IN
		(SELECT user_id FROM %v WHERE object_id = ?);`, escapeLiteral(db, "user"))), userID).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// FindRoleByName search for a role by given name
func FindRoleByName(db sqlx.Ext, name string) (*Role, error) {
	var r Role
	err := db.QueryRowx(db.Rebind("SELECT * FROM role WHERE name = ? LIMIT 1"), name).StructScan(&r)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &r, nil
}

// FindAllRoles returns a page of role records
func FindAllRoles(db sqlx.Ext, pager entities.Pager, sorter entities.Sorter) ([]*Role, error) {
	q := fmt.Sprintf(`SELECT * FROM role AS r
		%v %v;`, orderByClause(sorter, "r"), limitOffset(pager))
	rows, err := db.Queryx(q)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []*Role{}
	for rows.Next() {
		var r Role
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &r)
	}
	return roles, nil
}

// FindRolesByUser returns a page of role records filtered by a given user
func FindRolesByUser(db *sqlx.DB, userID string, pager entities.Pager, sorter entities.Sorter) ([]*Role, error) {
	q := fmt.Sprintf(`SELECT r.* FROM user_role AS ur
		LEFT JOIN role AS r ON r.role_id=ur.role_id
		LEFT JOIN %v AS u ON u.user_id=ur.user_id
		WHERE u.object_id=? GROUP BY r.role_id
		%v %v;`, escapeLiteral(db, "user"), orderByClause(sorter, "r"), limitOffset(pager))
	rows, err := db.Queryx(db.Rebind(q), userID)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []*Role{}
	for rows.Next() {
		var r Role
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &r)
	}
	return roles, nil
}

// RoleAddPermission updates a role with a given permission
func RoleAddPermission(db sqlx.Ext, rolePK, permissionPK int64) error {
	_, err := db.Exec(db.Rebind("INSERT INTO role_permission (role_id, permission_id) VALUES (?, ?);"), rolePK, permissionPK)
	return err
}

// RoleRemovePermission removes a given permission from a role
func RoleRemovePermission(db sqlx.Ext, rolePK, permissionPK int64) error {
	_, err := db.Exec(db.Rebind("DELETE FROM role_permission WHERE role_id = ? AND permission_id = ?;"), rolePK, permissionPK)
	return err
}
