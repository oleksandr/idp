package usecases

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/db"
	"github.com/oleksandr/idp/entities"
	"gopkg.in/gorp.v1"
)

//
// RBACInteractor is an interface that defines all RBAC related use-cases
// signatures
//
type RBACInteractor interface {
	CreatePermission(p entities.BasicPermission) error
	DeletePermission(name string) error
	CreateRole(r entities.BasicRole) error
	DeleteRole(name string) error
	UpdateRoleWithPermissions(roleName string, permissions []string) error
	RemovePermissionsFromRole(permissions []string, roleName string) error
	ListPermissions(pager entities.Pager, sorter entities.Sorter) (*entities.BasicPermissionCollection, error)
	ListPermissionsByRole(roleName string, pager entities.Pager, sorter entities.Sorter) (*entities.BasicPermissionCollection, error)
	ListRoles(pager entities.Pager, sorter entities.Sorter) (*entities.BasicRoleCollection, error)
	ListRolesByUser(userID string, pager entities.Pager, sorter entities.Sorter) (*entities.BasicRoleCollection, error)
}

// RBACInteractorImpl is an actual interactor that implements RBACInteractor
type RBACInteractorImpl struct {
	DB    *sqlx.DB
	DBMap *gorp.DbMap
}

// CreatePermission creates a new permission
func (inter *RBACInteractorImpl) CreatePermission(p entities.BasicPermission) error {
	d := &db.Permission{
		Name:           p.Name,
		Description:    p.Description,
		Enabled:        p.Enabled,
		EvaluationRule: p.EvaluationRule,
	}
	err := inter.DBMap.Insert(d)
	return err
}

// DeletePermission deletes existing permission
func (inter *RBACInteractorImpl) DeletePermission(name string) error {
	var (
		p    db.Permission
		err  error
		pTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "permission")
	)
	err = inter.DBMap.SelectOne(&p, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", pTbl), name)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	} else if err != nil {
		return err
	}
	_, err = inter.DBMap.Delete(&p)
	return err
}

// FindPermission finds a role by given role name
func (inter *RBACInteractorImpl) FindPermission(name string) (*entities.BasicPermission, error) {
	var (
		p    db.Permission
		err  error
		pTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "permission")
	)
	err = inter.DBMap.SelectOne(&p, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", pTbl), name)
	if err == sql.ErrNoRows {
		return nil, entities.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return permissionToEntity(&p), nil
}

// UpdatePermission updates all attributes of a given domain entity in the database
func (inter *RBACInteractorImpl) UpdatePermission(perm entities.BasicPermission) error {
	var (
		p    db.Permission
		err  error
		pTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "permission")
	)
	err = inter.DBMap.SelectOne(&p, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", pTbl), perm.Name)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	} else if err != nil {
		return err
	}
	p.Description = perm.Description
	p.Enabled = perm.Enabled
	p.EvaluationRule = perm.EvaluationRule
	_, err = inter.DBMap.Update(&p)
	if err != nil {
		return err
	}
	return nil
}

// RenamePermission renames oldName permission to newName
func (inter *RBACInteractorImpl) RenamePermission(oldName, newName string) error {
	var (
		p    db.Permission
		err  error
		pTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "permission")
	)
	err = inter.DBMap.SelectOne(&p, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", pTbl), oldName)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	} else if err != nil {
		return err
	}
	p.Name = newName
	_, err = inter.DBMap.Update(&p)
	if err != nil {
		return err
	}
	return nil
}

// CreateRole creates a new role
func (inter *RBACInteractorImpl) CreateRole(r entities.BasicRole) error {
	d := &db.Role{
		Name:        r.Name,
		Description: r.Description,
		Enabled:     r.Enabled,
	}
	err := inter.DBMap.Insert(d)
	return err
}

// UpdateRole updates all attributes of a given domain entity in the database
func (inter *RBACInteractorImpl) UpdateRole(role entities.BasicRole) error {
	var (
		r    db.Role
		err  error
		rTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "role")
	)
	err = inter.DBMap.SelectOne(&r, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", rTbl), role.Name)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	} else if err != nil {
		return err
	}
	r.Description = role.Description
	r.Enabled = role.Enabled
	_, err = inter.DBMap.Update(&r)
	if err != nil {
		return err
	}
	return nil
}

// RenameRole renames oldName role to newName
func (inter *RBACInteractorImpl) RenameRole(oldName, newName string) error {
	var (
		r    db.Role
		err  error
		rTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "role")
	)
	err = inter.DBMap.SelectOne(&r, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", rTbl), oldName)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	} else if err != nil {
		return err
	}
	r.Name = newName
	_, err = inter.DBMap.Update(&r)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRole deletes existing role
func (inter *RBACInteractorImpl) DeleteRole(name string) error {
	err := db.DeleteRole(inter.DBMap, name)
	if err != nil {
		return err
	}
	return nil
}

// FindRole finds a role by given role name
func (inter *RBACInteractorImpl) FindRole(name string) (*entities.BasicRole, error) {
	var (
		r    db.Role
		err  error
		rTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "role")
	)
	err = inter.DBMap.SelectOne(&r, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", rTbl), name)
	if err == sql.ErrNoRows {
		return nil, entities.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return roleToEntity(&r), nil
}

// UpdateRoleWithPermissions adds given permissions to existing role
func (inter *RBACInteractorImpl) UpdateRoleWithPermissions(roleName string, permissions []string) error {
	var (
		r    db.Role
		err  error
		pk   int64
		pp   []int64
		rTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "role")
	)
	err = inter.DBMap.SelectOne(&r, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", rTbl), roleName)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	} else if err != nil {
		return err
	}

	for _, name := range permissions {
		pk, err = inter.DBMap.SelectInt(fmt.Sprintf("SELECT permission_id FROM %v WHERE name = ?", rTbl), name)
		if err != nil {
			return err
		}
		if pk == 0 {
			return entities.ErrNotFound
		}
		pp = append(pp, pk)
	}
	if len(pp) != len(permissions) {
		return fmt.Errorf("Resolve only %v of %v provided permissions", len(pp), len(permissions))
	}

	tx, err := inter.DBMap.Begin()
	if err != nil {
		return err
	}
	for _, pk = range pp {
		err = tx.Insert(&db.RolePermission{
			RolePK:       r.PK,
			PermissionPK: pk,
		})
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()

	return err
}

// RemovePermissionsFromRole removes givens permission from existing role
func (inter *RBACInteractorImpl) RemovePermissionsFromRole(permissions []string, roleName string) error {
	var (
		r    db.Role
		err  error
		pk   int64
		pp   []int64
		rTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "role")
	)
	err = inter.DBMap.SelectOne(&r, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", rTbl), roleName)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	} else if err != nil {
		return err
	}

	for _, name := range permissions {
		pk, err = inter.DBMap.SelectInt(fmt.Sprintf("SELECT permission_id FROM %v WHERE name = ?", rTbl), name)
		if err != nil {
			return err
		}
		if pk == 0 {
			return entities.ErrNotFound
		}
		pp = append(pp, pk)
	}
	if len(pp) != len(permissions) {
		return fmt.Errorf("Resolve only %v of %v provided permissions", len(pp), len(permissions))
	}

	tx, err := inter.DBMap.Begin()
	if err != nil {
		return err
	}
	for _, pk = range pp {
		_, err = tx.Delete(&db.RolePermission{
			RolePK:       r.PK,
			PermissionPK: pk,
		})
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()

	return err
}

// ListPermissions lists existing permissions page by page
func (inter *RBACInteractorImpl) ListPermissions(pager entities.Pager, sorter entities.Sorter) (*entities.BasicPermissionCollection, error) {
	pTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "permission")
	total, err := inter.DBMap.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %v", pTbl))
	if err != nil {
		return nil, err
	}
	var records []db.Permission
	q := "SELECT * FROM %v AS p %v %v;"
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, pTbl, db.OrderByClause(sorter, "p"), db.LimitOffset(pager)))
	if err != nil {
		return nil, err
	}
	c := &entities.BasicPermissionCollection{
		Permissions: []entities.BasicPermission{},
		Paginator:   *pager.CreatePaginator(len(records), total),
	}
	for _, r := range records {
		c.Permissions = append(c.Permissions, *permissionToEntity(&r))
	}
	return c, nil
}

// ListPermissionsByRole lists existing permissions by role page by page
func (inter *RBACInteractorImpl) ListPermissionsByRole(roleName string, pager entities.Pager, sorter entities.Sorter) (*entities.BasicPermissionCollection, error) {
	pTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "permission")
	total, err := inter.DBMap.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %v", pTbl))
	if err != nil {
		return nil, err
	}
	var records []db.Permission
	q := `SELECT p.* FROM role_permission AS rp
		LEFT JOIN %v AS p ON p.permission_id=rp.permission_id
		LEFT JOIN %v AS r ON r.role_id=rp.role_id
		WHERE r.name=? GROUP BY p.permission_id
		%v %v;`
	rTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "role")
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, pTbl, rTbl, db.OrderByClause(sorter, "p"), db.LimitOffset(pager)), roleName)
	if err != nil {
		return nil, err
	}
	c := &entities.BasicPermissionCollection{
		Permissions: []entities.BasicPermission{},
		Paginator:   *pager.CreatePaginator(len(records), total),
	}
	for _, r := range records {
		c.Permissions = append(c.Permissions, *permissionToEntity(&r))
	}
	return c, nil
}

// ListRoles lists existing roles page by page
func (inter *RBACInteractorImpl) ListRoles(pager entities.Pager, sorter entities.Sorter) (*entities.BasicRoleCollection, error) {
	rTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "role")
	total, err := inter.DBMap.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %v", rTbl))
	if err != nil {
		return nil, err
	}
	var records []db.Role
	q := "SELECT * FROM %v AS r %v %v;"
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, rTbl, db.OrderByClause(sorter, "r"), db.LimitOffset(pager)))
	if err != nil {
		return nil, err
	}
	c := &entities.BasicRoleCollection{
		Roles:     []entities.BasicRole{},
		Paginator: *pager.CreatePaginator(len(records), total),
	}
	for _, r := range records {
		c.Roles = append(c.Roles, *roleToEntity(&r))
	}
	return c, nil
}

// ListRolesByUser lists existing roles by user page by page
func (inter *RBACInteractorImpl) ListRolesByUser(userID string, pager entities.Pager, sorter entities.Sorter) (*entities.BasicRoleCollection, error) {
	rTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "role")
	total, err := inter.DBMap.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %v", rTbl))
	if err != nil {
		return nil, err
	}
	var records []db.Role
	q := `SELECT r.* FROM user_role AS ur
		LEFT JOIN %v AS r ON r.role_id=ur.role_id
		LEFT JOIN %v AS u ON u.user_id=ur.user_id
		WHERE u.object_id=? GROUP BY r.role_id
		%v %v;`
	uTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, rTbl, uTbl, db.OrderByClause(sorter, "r"), db.LimitOffset(pager)), userID)
	if err != nil {
		return nil, err
	}
	c := &entities.BasicRoleCollection{
		Roles:     []entities.BasicRole{},
		Paginator: *pager.CreatePaginator(len(records), total),
	}
	for _, r := range records {
		c.Roles = append(c.Roles, *roleToEntity(&r))
	}
	return c, nil
}

func permissionToEntity(p *db.Permission) *entities.BasicPermission {
	e := entities.NewBasicPermission(p.Name, p.Description)
	e.Enabled = p.Enabled
	e.EvaluationRule = p.EvaluationRule
	return e
}

func roleToEntity(r *db.Role) *entities.BasicRole {
	e := entities.NewBasicRole(r.Name, r.Description)
	e.Enabled = r.Enabled
	return e
}
