package usecases

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/dl"
	"github.com/oleksandr/idp/entities"
)

//
// UserInteractor is an interface that defines all user related use-cases
// signatures
//
type UserInteractor interface {
	Create(user entities.BasicUser, domainIDs []string) error
	Update(user entities.BasicUser, addDomainIDs []string, removeDomainIDs []string) error
	Delete(user entities.BasicUser) error
	Find(id string) (*entities.BasicUser, error)
	FindInDomain(userID, domainID string) (*entities.BasicUser, error)
	FindByNameInDomain(userName, domainID string) (*entities.BasicUser, error)
	CountDomains(userID string) (int64, error)
	AssignRoles(userID string, roleNames []string) error
	RevokeRoles(userID string, roleNames []string) error
	List(pager entities.Pager, sorter entities.Sorter) (*entities.UserCollection, error)
	ListByDomain(domainID string, pager entities.Pager, sorter entities.Sorter) (*entities.UserCollection, error)
}

// UserInteractorImpl is an actual interactor that implements UserInteractor
type UserInteractorImpl struct {
	DB *sqlx.DB
}

// Create creates a new user with a given name and description and assign it to a given domain
func (inter *UserInteractorImpl) Create(user entities.BasicUser, domainIDs []string) error {
	if ok, err := user.IsValid(); !ok {
		return fmt.Errorf("User is not valid: %v", err.Error())
	}
	var (
		err     error
		found   *dl.Domain
		domains = []*dl.Domain{}
	)
	// Fetch domains for assignment
	for _, id := range domainIDs {
		found, err = dl.FindDomain(inter.DB, id)
		if err != nil {
			return err
		}
		domains = append(domains, found)
	}
	// Create a user
	u := dl.User{
		ID:       user.ID,
		Name:     user.Name,
		Password: user.Password,
		Enabled:  user.Enabled,
	}
	created, err := dl.SaveUser(inter.DB, u)
	if err != nil {
		return err
	}
	err = dl.ExecuteTransactionally(inter.DB, func(ext sqlx.Ext) error {
		for _, d := range domains {
			err = dl.AddUserToDomain(ext, created.ID, d.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		dl.DeleteUser(inter.DB, created.ID)
	}
	return err
}

// Update updates all attributes of a given user entity in the database
func (inter *UserInteractorImpl) Update(user entities.BasicUser, addDomainIDs []string, removeDomainIDs []string) error {
	if ok, err := user.IsValid(); !ok {
		return fmt.Errorf("User is not valid: %v", err.Error())
	}
	var (
		err    error
		found  *dl.Domain
		add    = []*dl.Domain{}
		remove = []*dl.Domain{}
	)
	// Fetch domains for assignment
	for _, id := range addDomainIDs {
		found, err = dl.FindDomain(inter.DB, id)
		if err != nil {
			return err
		}
		add = append(add, found)
	}
	// Fetch domains for removal
	for _, id := range removeDomainIDs {
		found, err = dl.FindDomain(inter.DB, id)
		if err != nil {
			return err
		}
		remove = append(remove, found)
	}
	// Update user and assign to provided domains or unassign from given domains
	err = dl.ExecuteTransactionally(inter.DB, func(ext sqlx.Ext) error {
		// Create a user
		u := dl.User{
			ID:       user.ID,
			Name:     user.Name,
			Password: user.Password,
			Enabled:  user.Enabled,
		}
		updated, err := dl.SaveUser(ext, u)
		if err != nil {
			return err
		}
		// Assign user to domains
		for _, d := range add {
			dl.AddUserToDomain(ext, updated.ID, d.ID)
		}
		// Remove user from domains
		for _, d := range remove {
			dl.RemoveUserFromDomain(ext, updated.ID, d.ID)
		}
		return nil
	})
	return err
}

// Delete removes user and all assigned entities from storage
func (inter *UserInteractorImpl) Delete(user entities.BasicUser) error {
	err := dl.DeleteUser(inter.DB, user.ID)
	if err != nil {
		return err
	}
	return nil
}

// Find finds a user by given user ID
func (inter *UserInteractorImpl) Find(id string) (*entities.BasicUser, error) {
	r, err := dl.FindUser(inter.DB, id)
	if err != nil {
		return nil, err
	}
	return basicUserRecordToEntity(r), nil
}

// FindInDomain checks if a given user is assigned to a given domain
func (inter *UserInteractorImpl) FindInDomain(userID, domainID string) (*entities.BasicUser, error) {
	r, err := dl.FindUserInDomain(inter.DB, userID, domainID)
	if err != nil {
		return nil, err
	}
	return basicUserRecordToEntity(r), nil
}

// FindByNameInDomain checks if a given user is assigned to a given domain
func (inter *UserInteractorImpl) FindByNameInDomain(userName, domainID string) (*entities.BasicUser, error) {
	r, err := dl.FindUserByNameInDomain(inter.DB, userName, domainID)
	if err != nil {
		return nil, err
	}
	return basicUserRecordToEntity(r), nil
}

// CountDomains return number of users in a domain defined by given domain ID
func (inter *UserInteractorImpl) CountDomains(userID string) (int64, error) {
	c, err := dl.CountDomainsByUser(inter.DB, userID)
	if err != nil {
		return -1, err
	}
	return c, nil
}

// AssignRoles assigns given set of roles to user
func (inter *UserInteractorImpl) AssignRoles(userID string, roleNames []string) error {
	var (
		found *dl.Role
		roles = []*dl.Role{}
	)
	// Find the user
	u, err := dl.FindUser(inter.DB, userID)
	if err != nil {
		return err
	}
	// Fetch roles for assignment
	for _, name := range roleNames {
		found, err = dl.FindRoleByName(inter.DB, name)
		if err != nil {
			return err
		}
		roles = append(roles, found)
	}
	// Assign roles in transaction
	err = dl.ExecuteTransactionally(inter.DB, func(ext sqlx.Ext) error {
		for _, r := range roles {
			err = dl.AssignRoleToUser(ext, r.Name, u.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// RevokeRoles revokes given set of roles from user
func (inter *UserInteractorImpl) RevokeRoles(userID string, roleNames []string) error {
	var (
		found *dl.Role
		roles = []*dl.Role{}
	)
	// Find the user
	u, err := dl.FindUser(inter.DB, userID)
	if err != nil {
		return err
	}
	// Fetch roles for revoking
	for _, name := range roleNames {
		found, err = dl.FindRoleByName(inter.DB, name)
		if err != nil {
			return err
		}
		roles = append(roles, found)
	}
	// Revoke roles in transaction
	err = dl.ExecuteTransactionally(inter.DB, func(ext sqlx.Ext) error {
		for _, r := range roles {
			err = dl.RevokeRoleFromUser(ext, r.Name, u.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// List implements a paginated listing of users
func (inter *UserInteractorImpl) List(pager entities.Pager, sorter entities.Sorter) (*entities.UserCollection, error) {
	total, err := dl.CountUsers(inter.DB)
	if err != nil {
		return nil, err
	}
	records, err := dl.FindAllUsers(inter.DB, pager, sorter)
	if err != nil {
		return nil, err
	}
	c := &entities.UserCollection{
		Users:     []*entities.User{},
		Paginator: pager.CreatePaginator(len(records), total),
	}
	for _, r := range records {
		c.Users = append(c.Users, userRecordToEntity(r))
	}
	return c, nil
}

// ListByDomain implements a paginated listing of users filtered by given domain ID
func (inter *UserInteractorImpl) ListByDomain(domainID string, pager entities.Pager, sorter entities.Sorter) (*entities.UserCollection, error) {
	total, err := dl.CountUsersByDomain(inter.DB, domainID)
	if err != nil {
		return nil, err
	}
	records, err := dl.FindUsersByDomain(inter.DB, domainID, pager, sorter)
	if err != nil {
		return nil, err
	}
	c := &entities.UserCollection{
		Users:     []*entities.User{},
		Paginator: pager.CreatePaginator(len(records), total),
	}
	for _, r := range records {
		c.Users = append(c.Users, userRecordToEntity(r))
	}
	return c, nil
}

func basicUserRecordToEntity(record *dl.User) *entities.BasicUser {
	u := entities.NewBasicUser(record.Name)
	u.ID = record.ID
	u.Password = record.Password
	u.Enabled = record.Enabled
	return u
}

func userRecordToEntity(record *dl.User) *entities.User {
	u := new(entities.User)
	basicUser := basicUserRecordToEntity(record)
	u.BasicUser = *basicUser
	u.DomainsCount = record.DomainsCount
	return u
}
