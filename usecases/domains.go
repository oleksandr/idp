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
// DomainInteractor is an interface that defines all domain related use-cases
// signatures
//
type DomainInteractor interface {
	Create(domain entities.BasicDomain) error
	Update(domain entities.BasicDomain) error
	Delete(id string) error
	Find(id string) (*entities.BasicDomain, error)
	FindByName(name string) (*entities.BasicDomain, error)
	CountUsers(domainID string) (int64, error)
	List(pager entities.Pager, sorter entities.Sorter) (*entities.DomainCollection, error)
	ListByUser(userID string, pager entities.Pager, sorter entities.Sorter) (*entities.DomainCollection, error)
}

// DomainInteractorImpl is an actual interactor that implements DomainInteractor
type DomainInteractorImpl struct {
	DBMap *gorp.DbMap
}

// Create creates a new domain with a given name and description
func (inter *DomainInteractorImpl) Create(domain entities.BasicDomain) error {
	if ok, err := domain.IsValid(); !ok {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Domain is invalid", err)
	}
	now := time.Now().UTC()
	d := &db.Domain{
		ID:          domain.ID,
		Name:        domain.Name,
		Description: domain.Description,
		Enabled:     domain.Enabled,
		CreatedOn:   now,
		UpdatedOn:   now,
	}
	err := inter.DBMap.Insert(d)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to create a domain", err)
	}
	return nil
}

// Update updates all attributes of a given domain entity in the database
func (inter *DomainInteractorImpl) Update(domain entities.BasicDomain) error {
	if ok, err := domain.IsValid(); !ok {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Domain is invalid", err)
	}

	d, err := findDomainByID(inter.DBMap, domain.ID)
	if err != nil {
		return err
	}

	d.ID = domain.ID
	d.Name = domain.Name
	d.Description = domain.Description
	d.Enabled = domain.Enabled
	d.UpdatedOn = time.Now().UTC()

	_, err = inter.DBMap.Update(&d)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeConflict, "Failed to updated domain", err)
	}
	return nil
}

// Delete removes domain and all assigned entities from storage
func (inter *DomainInteractorImpl) Delete(id string) error {
	err := db.DeleteDomain(inter.DBMap, id)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to delete domain by given ID", err)
	}
	return nil
}

// Find finds a domain by given domain ID
func (inter *DomainInteractorImpl) Find(id string) (*entities.BasicDomain, error) {
	d, err := findDomainByID(inter.DBMap, id)
	if err != nil {
		return nil, err
	}
	return domainToEntity(d), nil
}

// FindByName finds a domain by given domain name
func (inter *DomainInteractorImpl) FindByName(name string) (*entities.BasicDomain, error) {
	var d db.Domain
	err := inter.DBMap.SelectOne(&d, "SELECT * FROM domain WHERE name = ?", name)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "Domain not found by given name", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a domain", err)
	}
	return domainToEntity(&d), nil
}

// CountUsers return number of users in a domain defined by given domain ID
func (inter *DomainInteractorImpl) CountUsers(domainID string) (int64, error) {
	c, err := inter.DBMap.SelectInt("SELECT COUNT(*) FROM domain_user WHERE domain_id IN (SELECT domain_id from domains WHERE object_id = ?)", domainID)
	if err != nil {
		return -1, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to count users", err)
	}
	return c, nil
}

// List implements a paginated listing of domains
func (inter *DomainInteractorImpl) List(pager entities.Pager, sorter entities.Sorter) (*entities.DomainCollection, error) {
	total, err := inter.DBMap.SelectInt("SELECT COUNT(*) FROM domain")
	if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to count domains", err)
	}

	var records []db.DomainWithStats
	q := `SELECT d.*, COUNT(du.domain_id) AS users_count FROM domain AS d
        LEFT JOIN domain_user AS du ON d.domain_id = du.domain_id
        GROUP BY d.domain_id %v %v`
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, db.OrderByClause(sorter, "d"), db.LimitOffset(pager)))
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "No domains found", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of domains", err)
	}

	c := &entities.DomainCollection{
		Domains:   []entities.Domain{},
		Paginator: *pager.CreatePaginator(len(records), total),
	}
	for _, r := range records {
		c.Domains = append(c.Domains, entities.Domain{*domainToEntity(&r.Domain), r.UsersCount})
	}
	return c, nil
}

// ListByUser implements a paginated listing of domains filtered by a given user ID
func (inter *DomainInteractorImpl) ListByUser(userID string, pager entities.Pager, sorter entities.Sorter) (*entities.DomainCollection, error) {
	userTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	q := fmt.Sprintf(`SELECT count(*) FROM domain_user WHERE user_id IN (SELECT user_id FROM %v WHERE object_id = ?);`, userTbl)
	total, err := inter.DBMap.SelectInt(q, userID)
	if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to count domains for a given user", err)
	}

	var records []db.DomainWithStats
	q = `SELECT d.*, COUNT(du.domain_id) AS users_count FROM domain AS d
        LEFT JOIN domain_user AS du ON d.domain_id = du.domain_id
        WHERE d.domain_id IN (
            SELECT DISTINCT domain_id FROM domain_user WHERE user_id
                IN (SELECT user_id FROM %v WHERE object_id = ?)
        )
        GROUP BY d.domain_id %v %v;`
	_, err = inter.DBMap.Select(&records, fmt.Sprintf(q, userTbl, db.OrderByClause(sorter, "d"), db.LimitOffset(pager)), userID)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "No domains found for a given user", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of domains", err)
	}

	c := &entities.DomainCollection{
		Domains:   []entities.Domain{},
		Paginator: *pager.CreatePaginator(len(records), total),
	}
	for _, r := range records {
		c.Domains = append(c.Domains, entities.Domain{*domainToEntity(&r.Domain), r.UsersCount})
	}
	return c, nil
}

func findDomainByID(dbmap *gorp.DbMap, id string) (*db.Domain, error) {
	var (
		d   db.Domain
		err error
	)

	err = dbmap.SelectOne(&d, "SELECT * FROM domain WHERE object_id = ?", id)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "Domain not found by given ID", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a domain", err)
	}

	return &d, nil
}

func findDomainByName(dbmap *gorp.DbMap, name string) (*db.Domain, error) {
	var (
		d   db.Domain
		err error
	)

	err = dbmap.SelectOne(&d, "SELECT * FROM domain WHERE name = ?", name)
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "Domain not found by given name", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a domain", err)
	}

	return &d, nil
}

func domainToEntity(d *db.Domain) *entities.BasicDomain {
	e := entities.NewBasicDomain(d.Name, d.Description)
	e.ID = d.ID
	e.Enabled = d.Enabled
	return e
}
