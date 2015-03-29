package usecases

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/db"
	"github.com/oleksandr/idp/entities"
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
	DB    *sqlx.DB
	DBMap *gorp.DbMap
}

// Create creates a new user with a given name and description and assign it to a given domain
func (inter *UserInteractorImpl) Create(user entities.BasicUser, domainIDs []string) error {
	if ok, err := user.IsValid(); !ok {
		return fmt.Errorf("User is not valid: %v", err.Error())
	}

	var (
		err       error
		pk        int64
		domainPKs []int64
	)

	for _, id := range domainIDs {
		pk, err = inter.DBMap.SelectInt("SELECT domain_id FROM domain WHERE object_id = ?", id)
		if err != nil {
			return err
		}
		if pk == 0 {
			return entities.ErrNotFound
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
		return err
	}

	err = tx.Insert(&u)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, pk = range domainPKs {
		err = tx.Insert(&db.DomainUser{
			UserPK:   u.PK,
			DomainPK: pk,
		})
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()

	return err
}

// Update updates all attributes of a given user entity in the database
func (inter *UserInteractorImpl) Update(user entities.BasicUser, addDomainIDs []string, removeDomainIDs []string) error {
	if ok, err := user.IsValid(); !ok {
		return fmt.Errorf("User is not valid: %v", err.Error())
	}

	var (
		err       error
		userTbl   = inter.DBMap.Dialect.QuotedTableForQuery("", "user")
		u         db.User
		pk        int64
		addPKs    []int64
		removePKs []int64
	)

	err = inter.DBMap.SelectOne(&u, fmt.Sprintf("SELECT * FROM %v WHERE object_id = ?", userTbl), user.ID)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	}

	for _, id := range addDomainIDs {
		pk, err = inter.DBMap.SelectInt("SELECT domain_id FROM domain WHERE object_id = ?", id)
		if err != nil {
			return err
		}
		if pk == 0 {
			return entities.ErrNotFound
		}
		addPKs = append(addPKs, pk)
	}

	for _, id := range removeDomainIDs {
		pk, err = inter.DBMap.SelectInt("SELECT domain_id FROM domain WHERE object_id = ?", id)
		if err != nil {
			return err
		}
		if pk == 0 {
			return entities.ErrNotFound
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
		return err
	}

	_, err = tx.Update(&u)
	if err != nil {
		tx.Rollback()
		return nil
	}

	// Assign user to domains
	for _, pk = range addPKs {
		err = tx.Insert(&db.DomainUser{
			UserPK:   u.PK,
			DomainPK: pk,
		})
		if err != nil {
			tx.Rollback()
			return err
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
			return err
		}
	}

	err = tx.Commit()

	return err
}

// Delete removes user and all assigned entities from storage
func (inter *UserInteractorImpl) Delete(id string) error {
	err := db.DeleteUser(inter.DBMap, id)
	if err != nil {
		return err
	}
	return nil
}

// Find finds a user by given user ID
func (inter *UserInteractorImpl) Find(id string) (*entities.BasicUser, error) {
	var (
		u       db.User
		userTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	)
	err := inter.DBMap.SelectOne(&u, fmt.Sprintf("SELECT * FROM %v WHERE object_id = ?", userTbl), id)
	if err == sql.ErrNoRows {
		return nil, entities.ErrNotFound
	}
	return userToEntity(&u), nil
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
		return nil, entities.ErrNotFound
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
		return nil, entities.ErrNotFound
	}
	return userToEntity(&u), nil
}

// CountDomains return number of users in a domain defined by given domain ID
func (inter *UserInteractorImpl) CountDomains(userID string) (int64, error) {
	q := `SELECT count(*) FROM domain_user WHERE user_id IN (SELECT user_id FROM %v WHERE object_id = ?);`
	userTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	c, err := inter.DBMap.SelectInt(fmt.Sprintf(q, userTbl), userID)
	if err != nil {
		return -1, err
	}
	return c, nil
}

// AssignRoles assigns given set of roles to user
func (inter *UserInteractorImpl) AssignRoles(userID string, roleNames []string) error {
	var (
		err     error
		pk      int64
		u       db.User
		roles   []int64
		userTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	)

	// Find a user
	err = inter.DBMap.SelectOne(&u, fmt.Sprintf("SELECT * FROM %v WHERE object_id = ?", userTbl), userID)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	}

	// Fetch roles for assignment
	for _, name := range roleNames {
		pk, err = inter.DBMap.SelectInt("SELECT role_id FROM role WHERE name = ?", name)
		if err != nil {
			return err
		}
		if pk == 0 {
			return entities.ErrNotFound
		}
		roles = append(roles, pk)
	}

	// Assign roles in transaction
	tx, err := inter.DBMap.Begin()
	if err != nil {
		return err
	}

	for _, pk = range roles {
		err = tx.Insert(&db.UserRole{
			UserPK: u.PK,
			RolePK: pk,
		})
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()

	return err
}

// RevokeRoles revokes given set of roles from user
func (inter *UserInteractorImpl) RevokeRoles(userID string, roleNames []string) error {
	var (
		err     error
		pk      int64
		u       db.User
		roles   []int64
		userTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	)

	// Find a user
	err = inter.DBMap.SelectOne(&u, fmt.Sprintf("SELECT * FROM %v WHERE object_id = ?", userTbl), userID)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	}

	// Fetch roles for assignment
	for _, name := range roleNames {
		pk, err = inter.DBMap.SelectInt("SELECT role_id FROM role WHERE name = ?", name)
		if err != nil {
			return err
		}
		if pk == 0 {
			return entities.ErrNotFound
		}
		roles = append(roles, pk)
	}

	// Assign roles in transaction
	tx, err := inter.DBMap.Begin()
	if err != nil {
		return err
	}

	for _, pk = range roles {
		_, err = tx.Delete(&db.UserRole{
			UserPK: u.PK,
			RolePK: pk,
		})
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()

	return err
}

// List implements a paginated listing of users
func (inter *UserInteractorImpl) List(pager entities.Pager, sorter entities.Sorter) (*entities.UserCollection, error) {
	userTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	total, err := inter.DBMap.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %v", userTbl))
	if err != nil {
		return nil, err
	}
	var records []db.UserWithStats
	q := `SELECT u.*, count(du.user_id) AS domains_count FROM %v AS u
		LEFT JOIN domain_user AS du ON u.user_id = du.user_id
		GROUP BY u.user_id %v %v;`
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, userTbl, db.OrderByClause(sorter, "u"), db.LimitOffset(pager)))
	if err != nil {
		return nil, err
	}
	log.Println(records)
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
	total, err := inter.DBMap.SelectInt(fmt.Sprintf("SELECT COUNT(*) FROM %v", userTbl))
	if err != nil {
		return nil, err
	}
	var records []db.UserWithStats
	q := `SELECT u.*, count(du.user_id) AS domains_count FROM %v AS u
		LEFT JOIN domain_user AS du ON u.user_id = du.user_id
		WHERE u.user_id IN (
			SELECT DISTINCT user_id FROM domain_user WHERE domain_id
				IN (SELECT domain_id FROM domain WHERE object_id = ?)
		)
		GROUP BY u.user_id %v %v;`
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, userTbl, db.OrderByClause(sorter, "u"), db.LimitOffset(pager)), domainID)
	if err != nil {
		return nil, err
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

func userToEntity(u *db.User) *entities.BasicUser {
	e := entities.NewBasicUser(u.Name)
	u.ID = u.ID
	u.Password = u.Password
	u.Enabled = u.Enabled
	return e
}
