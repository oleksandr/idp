package usecases

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/oleksandr/idp/db"
	"github.com/oleksandr/idp/entities"
	"github.com/oleksandr/idp/errs"
	"gopkg.in/gorp.v1"
)

//
// UserInteractor is an interface that defines all user related use-cases
// signatures
//
type UserInteractor interface {
	Create(user entities.BasicUser, domainIDs []string) error
	Update(user entities.BasicUser, addDomainIDs []string, removeDomainIDs []string) error
	Delete(id string) error
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
	DBMap *gorp.DbMap
}

// Create creates a new user with a given name and description and assign it to a given domain
func (inter *UserInteractorImpl) Create(user entities.BasicUser, domainIDs []string) error {
	if ok, err := user.IsValid(); !ok {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "user is invalid", err)
	}

	var (
		err       error
		pk        int64
		domainPKs []int64
	)

	for _, id := range domainIDs {
		pk, err = inter.DBMap.SelectInt("SELECT domain_id FROM domain WHERE object_id = ?;", id)
		if err != nil || pk == 0 {
			return errs.NewUseCaseError(errs.ErrorTypeNotFound, "Domain not found by given ID", err)
		}
		domainPKs = append(domainPKs, pk)
	}

	now := time.Now().UTC()
	u := db.User{
		ID:        user.ID,
		Name:      user.Name,
		Password:  user.Password,
		Enabled:   user.Enabled,
		CreatedOn: now,
		UpdatedOn: now,
	}

	tx, err := inter.DBMap.Begin()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to begin transaction", err)
	}

	err = tx.Insert(&u)
	if err != nil {
		tx.Rollback()
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to create user", err)
	}

	for _, pk = range domainPKs {
		err = tx.Insert(&db.DomainUser{
			UserPK:   u.PK,
			DomainPK: pk,
		})
		if err != nil {
			tx.Rollback()
			return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to assign user to a domain", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to commit transaction", err)
	}

	return nil
}

// Update updates all attributes of a given user entity in the database
func (inter *UserInteractorImpl) Update(user entities.BasicUser, addDomainIDs []string, removeDomainIDs []string) error {
	if ok, err := user.IsValid(); !ok {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "user is invalid", err)
	}

	var (
		err       error
		u         *db.User
		pk        int64
		addPKs    []int64
		removePKs []int64
	)

	u, err = findUserByID(inter.DBMap, user.ID)
	if err != nil {
		return err
	}

	for _, id := range addDomainIDs {
		pk, err = inter.DBMap.SelectInt("SELECT domain_id FROM domain WHERE object_id = ?", id)
		if err != nil || pk == 0 {
			return errs.NewUseCaseError(errs.ErrorTypeNotFound, "Domain not found by given ID", err)
		}
		addPKs = append(addPKs, pk)
	}

	for _, id := range removeDomainIDs {
		pk, err = inter.DBMap.SelectInt("SELECT domain_id FROM domain WHERE object_id = ?", id)
		if err != nil || pk == 0 {
			return errs.NewUseCaseError(errs.ErrorTypeNotFound, "Domain not found by given ID", err)
		}
		removePKs = append(removePKs, pk)
	}

	u.ID = user.ID
	u.Name = user.Name
	u.Password = user.Password
	u.Enabled = user.Enabled
	u.UpdatedOn = time.Now().UTC()

	tx, err := inter.DBMap.Begin()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to begin transaction", err)
	}

	_, err = tx.Update(u)
	if err != nil {
		tx.Rollback()
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to update user", err)
	}

	// Assign user to domains
	for _, pk = range addPKs {
		err = tx.Insert(&db.DomainUser{
			UserPK:   u.PK,
			DomainPK: pk,
		})
		if err != nil {
			tx.Rollback()
			return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to assign user to a domain", err)
		}
	}
	// Remove user from domains
	for _, pk = range removePKs {
		_, err = tx.Delete(&db.DomainUser{
			UserPK:   u.PK,
			DomainPK: pk,
		})
		if err != nil {
			tx.Rollback()
			return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to remove user from a domain", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to commit transaction", err)
	}

	return nil
}

// Delete removes user and all assigned entities from storage
func (inter *UserInteractorImpl) Delete(id string) error {
	err := db.DeleteUser(inter.DBMap, id)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to delete user by given ID", err)
	}
	return nil
}

// Find finds a user by given user ID
func (inter *UserInteractorImpl) Find(id string) (*entities.BasicUser, error) {
	u, err := findUserByID(inter.DBMap, id)
	if err != nil {
		return nil, err
	}
	return userToEntity(u), nil
}

// FindInDomain checks if a given user is assigned to a given domain
func (inter *UserInteractorImpl) FindInDomain(userID, domainID string) (*entities.BasicUser, error) {
	var (
		u       db.User
		userTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	)

	q := `SELECT u.* FROM domain_user
   		LEFT JOIN %v AS u ON domain_user.user_id=u.user_id
   		LEFT JOIN domain ON domain_user.domain_id=domain.domain_id
   		WHERE u.object_id = ?
   		AND domain.object_id = ?
   		LIMIT 1;`
	err := inter.DBMap.SelectOne(&u, fmt.Sprintf(q, userTbl), userID, domainID)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "User not found in a given domain", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a user in a domain", err)
	}

	return userToEntity(&u), nil
}

// FindByNameInDomain checks if a given user is assigned to a given domain
func (inter *UserInteractorImpl) FindByNameInDomain(userName, domainID string) (*entities.BasicUser, error) {
	var (
		u       db.User
		userTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	)
	q := `SELECT u.* FROM domain_user
   		LEFT JOIN %v AS u ON domain_user.user_id=u.user_id
   		LEFT JOIN domain ON domain_user.domain_id=domain.domain_id
   		WHERE u.name = ?
   		AND domain.object_id = ?
   		LIMIT 1;`
	err := inter.DBMap.SelectOne(&u, fmt.Sprintf(q, userTbl), userName, domainID)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "User not found in a given domain", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a user in a domain", err)
	}

	return userToEntity(&u), nil
}

// CountDomains return number of users in a domain defined by given domain ID
func (inter *UserInteractorImpl) CountDomains(userID string) (int64, error) {
	userTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	q := `SELECT count(*) FROM domain_user WHERE user_id IN (SELECT user_id FROM %v WHERE object_id = ?);`
	c, err := inter.DBMap.SelectInt(fmt.Sprintf(q, userTbl), userID)
	if err != nil {
		return -1, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to count domains for given user", err)
	}
	return c, nil
}

// AssignRoles assigns given set of roles to user
func (inter *UserInteractorImpl) AssignRoles(userID string, roleNames []string) error {
	var (
		err   error
		pk    int64
		u     *db.User
		roles []int64
	)

	// Find a user
	u, err = findUserByID(inter.DBMap, userID)
	if err != nil {
		return err
	}

	// Fetch roles for assignment
	for _, name := range roleNames {
		pk, err = inter.DBMap.SelectInt("SELECT role_id FROM role WHERE name = ?", name)
		if err != nil || pk == 0 {
			return errs.NewUseCaseError(errs.ErrorTypeNotFound, "Role not found by given name", err)
		}
		roles = append(roles, pk)
	}

	// Assign roles in transaction
	tx, err := inter.DBMap.Begin()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to begin transaction", err)
	}
	for _, pk = range roles {
		err = tx.Insert(&db.UserRole{
			UserPK: u.PK,
			RolePK: pk,
		})
		if err != nil {
			tx.Rollback()
			return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to assign role to user", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to commit transaction", err)
	}

	return nil
}

// RevokeRoles revokes given set of roles from user
func (inter *UserInteractorImpl) RevokeRoles(userID string, roleNames []string) error {
	var (
		err   error
		pk    int64
		u     *db.User
		roles []int64
	)

	// Find a user
	u, err = findUserByID(inter.DBMap, userID)
	if err != nil {
		return err
	}

	// Fetch roles for assignment
	for _, name := range roleNames {
		pk, err = inter.DBMap.SelectInt("SELECT role_id FROM role WHERE name = ?", name)
		if err != nil || pk == 0 {
			return errs.NewUseCaseError(errs.ErrorTypeNotFound, "Role not found by given name", err)
		}
		roles = append(roles, pk)
	}

	// Assign roles in transaction
	tx, err := inter.DBMap.Begin()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to begin transaction", err)
	}

	for _, pk = range roles {
		_, err = tx.Delete(&db.UserRole{
			UserPK: u.PK,
			RolePK: pk,
		})
		if err != nil {
			tx.Rollback()
			return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to revoke role from user", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to commit transaction", err)
	}

	return nil
}

// List implements a paginated listing of users
func (inter *UserInteractorImpl) List(pager entities.Pager, sorter entities.Sorter) (*entities.UserCollection, error) {
	userTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	total, err := inter.DBMap.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %v", userTbl))
	if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to count users", err)
	}

	var records []db.UserWithStats
	q := `SELECT u.*, count(du.user_id) AS domains_count FROM %v AS u
		LEFT JOIN domain_user AS du ON u.user_id = du.user_id
		GROUP BY u.user_id %v %v;`
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, userTbl, db.OrderByClause(sorter, "u"), db.LimitOffset(pager)))
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "No users found", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of users", err)
	}

	c := &entities.UserCollection{
		Users:     []entities.User{},
		Paginator: *pager.CreatePaginator(len(records), total),
	}
	for _, r := range records {
		c.Users = append(c.Users, entities.User{*userToEntity(&r.User), r.DomainsCount})
	}
	return c, nil
}

// ListByDomain implements a paginated listing of users filtered by given domain ID
func (inter *UserInteractorImpl) ListByDomain(domainID string, pager entities.Pager, sorter entities.Sorter) (*entities.UserCollection, error) {
	userTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	q := "SELECT count(*) FROM domain_user WHERE domain_id IN (SELECT domain_id FROM domain WHERE object_id = ?);"
	total, err := inter.DBMap.SelectInt(q, domainID)
	if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to count user for a given domain", err)
	}

	var records []db.UserWithStats
	q = `SELECT u.*, count(du.user_id) AS domains_count FROM %v AS u
		LEFT JOIN domain_user AS du ON u.user_id = du.user_id
		WHERE u.user_id IN (
			SELECT DISTINCT user_id FROM domain_user WHERE domain_id
				IN (SELECT domain_id FROM domain WHERE object_id = ?)
		)
		GROUP BY u.user_id %v %v;`
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, userTbl, db.OrderByClause(sorter, "u"), db.LimitOffset(pager)), domainID)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "No users found for a given domain", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of users for a given domain", err)
	}

	c := &entities.UserCollection{
		Users:     []entities.User{},
		Paginator: *pager.CreatePaginator(len(records), total),
	}
	for _, r := range records {
		c.Users = append(c.Users, entities.User{*userToEntity(&r.User), r.DomainsCount})
	}
	return c, nil
}

func findUserByID(dbmap *gorp.DbMap, id string) (*db.User, error) {
	var (
		u       db.User
		err     error
		userTbl = dbmap.Dialect.QuotedTableForQuery("", "user")
	)

	err = dbmap.SelectOne(&u, fmt.Sprintf("SELECT * FROM %v WHERE object_id = ?", userTbl), id)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "User not found by given ID", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a user", err)
	}

	return &u, nil
}

func findUserByName(dbmap *gorp.DbMap, name string) (*db.User, error) {
	var (
		u       db.User
		err     error
		userTbl = dbmap.Dialect.QuotedTableForQuery("", "user")
	)

	err = dbmap.SelectOne(&u, fmt.Sprintf("SELECT * FROM %v WHERE name = ?", userTbl), name)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "User not found by given name", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a user", err)
	}

	return &u, nil
}

func findUserInDomain(dbmap *gorp.DbMap, userID, domainID string) (*db.User, error) {
	var (
		u       db.User
		err     error
		userTbl = dbmap.Dialect.QuotedTableForQuery("", "user")
	)

	q := fmt.Sprintf(`SELECT u.* FROM domain_user AS du
		INNER JOIN %v AS u ON u.user_id = du.user_id
		INNER JOIN domain AS d ON d.domain_id=du.domain_id
		WHERE u.object_id = ?
		AND d.object_id = ?;`, userTbl)
	err = dbmap.SelectOne(&u, q, userID, domainID)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "User not found in domain", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a user in domain", err)
	}

	return &u, nil
}

func userToEntity(u *db.User) *entities.BasicUser {
	e := entities.NewBasicUser(u.Name)
	e.ID = u.ID
	e.Password = u.Password
	e.Enabled = u.Enabled
	return e
}
