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
	Create(user entities.User, domainIDs []string) error
	Update(user entities.User) error
	Delete(user entities.User) error
	Find(id string) (*entities.User, error)
	List(pager entities.Pager, sorter entities.Sorter) (*entities.UserCollection, error)
}

// UserInteractorImpl is an actual interactor that implements UserInteractor
type UserInteractorImpl struct {
	DB *sqlx.DB
}

// Create creates a new user with a given name and description and assign it to a given domain
func (inter *UserInteractorImpl) Create(user entities.User, domainIDs []string) error {
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
	// Create user and assign to provided domains
	err = dl.ExecuteTransactionally(inter.DB, func(ext sqlx.Ext) error {
		// Create a user
		u := dl.User{
			ID:       user.ID,
			Name:     user.Name,
			Password: user.Password,
			Enabled:  user.Enabled,
		}
		created, err := dl.SaveUser(ext, u)
		if err != nil {
			return err
		}
		// Assign user to domains
		for _, d := range domains {
			dl.AddUserToDomain(ext, *created, *d)
		}
		return nil
	})
	return err
}

// Update updates all attributes of a given user entity in the database
func (inter *UserInteractorImpl) Update(user entities.User) error {
	if ok, err := user.IsValid(); !ok {
		return fmt.Errorf("User is not valid: %v", err.Error())
	}
	u := dl.User{
		ID:       user.ID,
		Name:     user.Name,
		Password: user.Password,
		Enabled:  user.Enabled,
	}
	_, err := dl.SaveUser(inter.DB, u)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes user and all assigned entities from storage
func (inter *UserInteractorImpl) Delete(user entities.User) error {
	err := dl.DeleteUser(inter.DB, user.ID)
	if err != nil {
		return err
	}
	return nil
}

// Find finds a user by given user ID
func (inter *UserInteractorImpl) Find(id string) (*entities.User, error) {
	r, err := dl.FindUser(inter.DB, id)
	if err != nil {
		return nil, err
	}
	u := new(entities.User)
	u.Name = r.Name
	u.ID = r.ID
	u.Password = r.Password
	u.Enabled = r.Enabled
	return u, nil
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
	var u *entities.User
	for _, dto := range records {
		u = entities.NewUser(dto.Name)
		u.ID = dto.ID
		u.Enabled = dto.Enabled
		c.Users = append(c.Users, u)
	}
	return c, nil
}
