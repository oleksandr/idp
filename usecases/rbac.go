package usecases

import (
	"database/sql"
	"fmt"

	"github.com/oleksandr/idp/db"
	"github.com/oleksandr/idp/entities"
	"github.com/oleksandr/idp/errs"
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
	AssertRole(userID, roleName string) (bool, error)
	AssertPermission(userID, permissionName string) (bool, error)
}

// RBACInteractorImpl is an actual interactor that implements RBACInteractor
type RBACInteractorImpl struct {
	DBMap *gorp.DbMap
}

// CreatePermission creates a new permission
func (inter *RBACInteractorImpl) CreatePermission(p entities.BasicPermission) error {
	if p.Name == "" {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Permission name cannot be empty", nil)
	}
	d := &db.Permission{
		Name:           p.Name,
		Description:    p.Description,
		Enabled:        p.Enabled,
		EvaluationRule: p.EvaluationRule,
	}
	err := inter.DBMap.Insert(d)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to create a permission", err)
	}
	return nil
}

func (inter *RBACInteractorImpl) findPermission(name string) (*db.Permission, error) {
	var (
		p    db.Permission
		err  error
		pTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "permission")
	)

	err = inter.DBMap.SelectOne(&p, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", pTbl), name)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "Permission not found by given name", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a permission", err)
	}

	return &p, nil
}

// DeletePermission deletes existing permission
func (inter *RBACInteractorImpl) DeletePermission(name string) error {
	p, err := inter.findPermission(name)
	if err != nil {
		return err
	}
	_, err = inter.DBMap.Delete(p)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to delete permission by given name", err)
	}
	return nil
}

// FindPermission finds a role by given role name
func (inter *RBACInteractorImpl) FindPermission(name string) (*entities.BasicPermission, error) {
	p, err := inter.findPermission(name)
	if err != nil {
		return nil, err
	}
	return permissionToEntity(p), nil
}

// UpdatePermission updates all attributes of a given domain entity in the database
func (inter *RBACInteractorImpl) UpdatePermission(perm entities.BasicPermission) error {
	p, err := inter.findPermission(perm.Name)
	if err != nil {
		return err
	}

	p.Description = perm.Description
	p.Enabled = perm.Enabled
	p.EvaluationRule = perm.EvaluationRule

	_, err = inter.DBMap.Update(&p)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to update permission", err)
	}
	return nil
}

// RenamePermission renames oldName permission to newName
func (inter *RBACInteractorImpl) RenamePermission(oldName, newName string) error {
	if newName == "" {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Permission name cannot be empty", nil)
	}

	p, err := inter.findPermission(oldName)
	if err != nil {
		return err
	}

	p.Name = newName

	_, err = inter.DBMap.Update(&p)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to rename permission", err)
	}
	return nil
}

// CreateRole creates a new role
func (inter *RBACInteractorImpl) CreateRole(r entities.BasicRole) error {
	if r.Name == "" {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Role name cannot be empty", nil)
	}

	d := &db.Role{
		Name:        r.Name,
		Description: r.Description,
		Enabled:     r.Enabled,
	}
	err := inter.DBMap.Insert(d)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to create a role", err)
	}
	return nil
}

func (inter *RBACInteractorImpl) findRole(name string) (*db.Role, error) {
	var (
		r    db.Role
		err  error
		rTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "role")
	)
	err = inter.DBMap.SelectOne(&r, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", rTbl), name)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "Role not found by given name", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a role", err)
	}

	return &r, nil
}

// UpdateRole updates all attributes of a given domain entity in the database
func (inter *RBACInteractorImpl) UpdateRole(role entities.BasicRole) error {
	if role.Name == "" {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Role name cannot be empty", nil)
	}

	r, err := inter.findRole(role.Name)
	if err != nil {
		return err
	}

	r.Description = role.Description
	r.Enabled = role.Enabled

	_, err = inter.DBMap.Update(r)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to update a role", err)
	}

	return nil
}

// RenameRole renames oldName role to newName
func (inter *RBACInteractorImpl) RenameRole(oldName, newName string) error {
	if newName == "" {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Role name cannot be empty", nil)
	}

	r, err := inter.findRole(oldName)
	if err != nil {
		return err
	}

	r.Name = newName

	_, err = inter.DBMap.Update(&r)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to rename a role", err)
	}
	return nil
}

// DeleteRole deletes existing role
func (inter *RBACInteractorImpl) DeleteRole(name string) error {
	err := db.DeleteRole(inter.DBMap, name)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to delete permission by given name", err)
	}
	return nil
}

// FindRole finds a role by given role name
func (inter *RBACInteractorImpl) FindRole(name string) (*entities.BasicRole, error) {
	r, err := inter.findRole(name)
	if err != nil {
		return nil, err
	}
	return roleToEntity(r), nil
}

// UpdateRoleWithPermissions adds given permissions to existing role
func (inter *RBACInteractorImpl) UpdateRoleWithPermissions(roleName string, permissions []string) error {
	var (
		err  error
		pk   int64
		pp   []int64
		pTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "permission")
	)
	r, err := inter.findRole(roleName)
	if err != nil {
		return err
	}

	for _, name := range permissions {
		pk, err = inter.DBMap.SelectInt(fmt.Sprintf("SELECT permission_id FROM %v WHERE name = ?", pTbl), name)
		if err != nil || pk == 0 {
			return errs.NewUseCaseError(errs.ErrorTypeNotFound, "Permission not found by given name", err)
		}
		pp = append(pp, pk)
	}
	if len(pp) != len(permissions) {
		return errs.NewUseCaseError(errs.ErrorTypeNotFound, fmt.Sprintf("Resolve only %v of %v provided permissions", len(pp), len(permissions)), nil)
	}

	tx, err := inter.DBMap.Begin()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to begin transaction", err)
	}
	for _, pk = range pp {
		err = tx.Insert(&db.RolePermission{
			RolePK:       r.PK,
			PermissionPK: pk,
		})
		if err != nil {
			tx.Rollback()
			return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to add permission to role", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to commit transaction", err)
	}

	return nil
}

// RemovePermissionsFromRole removes givens permission from existing role
func (inter *RBACInteractorImpl) RemovePermissionsFromRole(permissions []string, roleName string) error {
	var (
		pk   int64
		pp   []int64
		pTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "permission")
	)
	r, err := inter.findRole(roleName)
	if err != nil {
		return err
	}

	for _, name := range permissions {
		pk, err = inter.DBMap.SelectInt(fmt.Sprintf("SELECT permission_id FROM %v WHERE name = ?", pTbl), name)
		if err != nil || pk == 0 {
			return errs.NewUseCaseError(errs.ErrorTypeNotFound, "Permission not found by given name", err)
		}
		pp = append(pp, pk)
	}
	if len(pp) != len(permissions) {
		return errs.NewUseCaseError(errs.ErrorTypeNotFound, fmt.Sprintf("Resolve only %v of %v provided permissions", len(pp), len(permissions)), nil)
	}

	tx, err := inter.DBMap.Begin()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to begin transaction", err)
	}
	for _, pk = range pp {
		_, err = tx.Delete(&db.RolePermission{
			RolePK:       r.PK,
			PermissionPK: pk,
		})
		if err != nil {
			tx.Rollback()
			return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to remove permission from role", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to commit transaction", err)
	}

	return nil
}

// ListPermissions lists existing permissions page by page
func (inter *RBACInteractorImpl) ListPermissions(pager entities.Pager, sorter entities.Sorter) (*entities.BasicPermissionCollection, error) {
	pTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "permission")
	total, err := inter.DBMap.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %v", pTbl))
	if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to count permissions", err)
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
	rTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "role")

	q := fmt.Sprintf("SELECT COUNT(*) FROM role_permission WHERE role_id IN (SELECT role_id FROM role WHERE name = ?);", rTbl)
	total, err := inter.DBMap.SelectInt(q, roleName)
	if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to count permissions", err)
	}

	var records []db.Permission
	q = `SELECT p.* FROM role_permission AS rp
		LEFT JOIN %v AS p ON p.permission_id=rp.permission_id
		LEFT JOIN %v AS r ON r.role_id=rp.role_id
		WHERE r.name=? GROUP BY p.permission_id
		%v %v;`
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, pTbl, rTbl, db.OrderByClause(sorter, "p"), db.LimitOffset(pager)), roleName)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "No permission found for given role", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of permissions for role", err)
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
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to count roles", err)
	}

	var records []db.Role
	q := "SELECT * FROM %v AS r %v %v;"
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, rTbl, db.OrderByClause(sorter, "r"), db.LimitOffset(pager)))
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "No roles found", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of roles", err)
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
	uTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	rTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "role")

	q := fmt.Sprintf("SELECT COUNT(*) FROM user_role WHERE user_id IN (SELECT user_id FROM %v WHERE object_id = ?);", uTbl)
	total, err := inter.DBMap.SelectInt(q, userID)
	if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to count roles for given user", err)
	}

	var records []db.Role
	q = `SELECT r.* FROM user_role AS ur
		LEFT JOIN %v AS r ON r.role_id=ur.role_id
		LEFT JOIN %v AS u ON u.user_id=ur.user_id
		WHERE u.object_id=? GROUP BY r.role_id
		%v %v;`
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, rTbl, uTbl, db.OrderByClause(sorter, "r"), db.LimitOffset(pager)), userID)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "No roles found for given user", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of roles for given user", err)
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

// AssertRole checks if a given user has given role assigned
func (inter *RBACInteractorImpl) AssertRole(userID, roleName string) (bool, error) {
	uTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	rTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "role")
	q := fmt.Sprintf(`SELECT COUNT(*) FROM user_role AS ur
		INNER JOIN %v AS u ON u.user_id = ur.user_id
		INNER JOIN %v AS r ON r.role_id = ur.role_id
		WHERE u.object_id=? AND r.name=?
			AND u.is_enabled=1 AND r.is_enabled=1;`, uTbl, rTbl)
	total, err := inter.DBMap.SelectInt(q, userID, roleName)
	if err != nil {
		return false, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to assert role", err)
	}
	return total > 0, nil
}

// AssertPermission checks if a given user has a given permission via any of the assigned roles
func (inter *RBACInteractorImpl) AssertPermission(userID, permissionName string) (bool, error) {
	uTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	rTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "role")
	pTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "permission")
	q := fmt.Sprintf(`SELECT COUNT(*) FROM user_role AS ur
		INNER JOIN %v AS u ON u.user_id = ur.user_id
		INNER JOIN %v AS r ON r.role_id = ur.role_id
		INNER JOIN role_permission AS rp ON rp.role_id = r.role_id
		INNER JOIN %v AS p ON p.permission_id = rp.permission_id
		WHERE u.object_id=? AND p.name=?
			AND u.is_enabled=1 AND r.is_enabled=1 AND p.is_enabled=1;`, uTbl, rTbl, pTbl)
	total, err := inter.DBMap.SelectInt(q, userID, permissionName)
	if err != nil {
		return false, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to assert permission", err)
	}
	return total > 0, nil
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
