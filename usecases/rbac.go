package usecases

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/dl"
	"github.com/oleksandr/idp/entities"
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
	ListRoles(pager entities.Pager, sorter entities.Sorter) (*entities.BasicRoleCollection, error)
}

// RBACInteractorImpl is an actual interactor that implements RBACInteractor
type RBACInteractorImpl struct {
	DB *sqlx.DB
}

// CreatePermission creates a new permission
func (inter *RBACInteractorImpl) CreatePermission(p entities.BasicPermission) error {
	data := dl.Permission{
		Name:           p.Name,
		Description:    p.Description,
		Enabled:        p.Enabled,
		EvaluationRule: p.EvaluationRule,
	}
	_, err := dl.CreatePermission(inter.DB, data)
	return err
}

// DeletePermission deletes existing permission
func (inter *RBACInteractorImpl) DeletePermission(name string) error {
	perm, err := dl.FindPermissionByName(inter.DB, name)
	if err != nil {
		return err
	}
	err = dl.DeletePermission(inter.DB, perm.PK)
	return err
}

// FindPermission finds a role by given role name
func (inter *RBACInteractorImpl) FindPermission(name string) (*entities.BasicPermission, error) {
	r, err := dl.FindPermissionByName(inter.DB, name)
	if err != nil {
		return nil, err
	}
	return basicPermissionRecordToEntity(r), nil
}

// UpdatePermission updates all attributes of a given domain entity in the database
func (inter *RBACInteractorImpl) UpdatePermission(perm entities.BasicPermission) error {
	p, err := dl.FindPermissionByName(inter.DB, perm.Name)
	if err != nil {
		return err
	}
	p.Description = perm.Description
	p.Enabled = perm.Enabled
	p.EvaluationRule = perm.EvaluationRule
	err = dl.UpdatePermission(inter.DB, *p)
	if err != nil {
		return err
	}
	return nil
}

// RenamePermission renames oldName permission to newName
func (inter *RBACInteractorImpl) RenamePermission(oldName, newName string) error {
	p, err := dl.FindPermissionByName(inter.DB, oldName)
	if err != nil {
		return err
	}
	err = dl.RenamePermission(inter.DB, p.PK, newName)
	return err
}

// CreateRole creates a new role
func (inter *RBACInteractorImpl) CreateRole(r entities.BasicRole) error {
	data := dl.Role{
		Name:        r.Name,
		Description: r.Description,
		Enabled:     r.Enabled,
	}
	_, err := dl.CreateRole(inter.DB, data)
	return err
}

// UpdateRole updates all attributes of a given domain entity in the database
func (inter *RBACInteractorImpl) UpdateRole(role entities.BasicRole) error {
	r, err := dl.FindRoleByName(inter.DB, role.Name)
	if err != nil {
		return err
	}
	r.Description = role.Description
	r.Enabled = role.Enabled
	err = dl.UpdateRole(inter.DB, *r)
	if err != nil {
		return err
	}
	return nil
}

// RenameRole renames oldName role to newName
func (inter *RBACInteractorImpl) RenameRole(oldName, newName string) error {
	r, err := dl.FindRoleByName(inter.DB, oldName)
	if err != nil {
		return err
	}
	err = dl.RenameRole(inter.DB, r.PK, newName)
	return err
}

// DeleteRole deletes existing role
func (inter *RBACInteractorImpl) DeleteRole(name string) error {
	role, err := dl.FindRoleByName(inter.DB, name)
	if err != nil {
		return err
	}
	err = dl.DeleteRole(inter.DB, role.PK)
	return err
}

// FindRole finds a role by given role name
func (inter *RBACInteractorImpl) FindRole(name string) (*entities.BasicRole, error) {
	r, err := dl.FindRoleByName(inter.DB, name)
	if err != nil {
		return nil, err
	}
	return basicRoleRecordToEntity(r), nil
}

// UpdateRoleWithPermissions adds given permissions to existing role
func (inter *RBACInteractorImpl) UpdateRoleWithPermissions(roleName string, permissions []string) error {
	r, err := dl.FindRoleByName(inter.DB, roleName)
	if err != nil {
		return err
	}
	pp, err := dl.FindPermissionsByNames(inter.DB, permissions)
	if err != nil {
		return err
	}
	if len(pp) != len(permissions) {
		return fmt.Errorf("Resolve only %v of %v provided permissions", len(pp), len(permissions))
	}
	err = dl.ExecuteTransactionally(inter.DB, func(ext sqlx.Ext) error {
		for _, p := range pp {
			err = dl.RoleAddPermission(ext, r.PK, p.PK)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// RemovePermissionsFromRole removes givens permission from existing role
func (inter *RBACInteractorImpl) RemovePermissionsFromRole(permissions []string, roleName string) error {
	r, err := dl.FindRoleByName(inter.DB, roleName)
	if err != nil {
		return err
	}
	pp, err := dl.FindPermissionsByNames(inter.DB, permissions)
	if err != nil {
		return err
	}
	if len(pp) != len(permissions) {
		return fmt.Errorf("Resolve only %v of %v provided permissions", len(pp), len(permissions))
	}
	err = dl.ExecuteTransactionally(inter.DB, func(ext sqlx.Ext) error {
		for _, p := range pp {
			err = dl.RoleRemovePermission(ext, r.PK, p.PK)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// ListPermissions lists existing permissions page by page
func (inter *RBACInteractorImpl) ListPermissions(pager entities.Pager, sorter entities.Sorter) (*entities.BasicPermissionCollection, error) {
	total, err := dl.CountPermissions(inter.DB)
	if err != nil {
		return nil, err
	}
	records, err := dl.FindAllPermissions(inter.DB, pager, sorter)
	if err != nil {
		return nil, err
	}
	c := &entities.BasicPermissionCollection{
		Permissions: []*entities.BasicPermission{},
		Paginator:   pager.CreatePaginator(len(records), total),
	}
	for _, dto := range records {
		c.Permissions = append(c.Permissions, basicPermissionRecordToEntity(dto))
	}
	return c, nil
}

// ListRoles lists existing roles page by page
func (inter *RBACInteractorImpl) ListRoles(pager entities.Pager, sorter entities.Sorter) (*entities.BasicRoleCollection, error) {
	total, err := dl.CountRoles(inter.DB)
	if err != nil {
		return nil, err
	}
	records, err := dl.FindAllRoles(inter.DB, pager, sorter)
	if err != nil {
		return nil, err
	}
	c := &entities.BasicRoleCollection{
		Roles:     []*entities.BasicRole{},
		Paginator: pager.CreatePaginator(len(records), total),
	}
	for _, dto := range records {
		c.Roles = append(c.Roles, basicRoleRecordToEntity(dto))
	}
	return c, nil
}

func basicPermissionRecordToEntity(record *dl.Permission) *entities.BasicPermission {
	p := entities.NewBasicPermission(record.Name, record.Description)
	p.Enabled = record.Enabled
	p.EvaluationRule = record.EvaluationRule
	return p
}

func basicRoleRecordToEntity(record *dl.Role) *entities.BasicRole {
	r := entities.NewBasicRole(record.Name, record.Description)
	r.Enabled = record.Enabled
	return r
}
