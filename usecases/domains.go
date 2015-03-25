package usecases

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/dl"
	"github.com/oleksandr/idp/entities"
)

//
// DomainInteractor is an interface that defines all domain related use-cases
// signatures
//
type DomainInteractor interface {
	Create(domain entities.BasicDomain) error
	Update(domain entities.BasicDomain) error
	Delete(domain entities.BasicDomain) error
	Find(id string) (*entities.BasicDomain, error)
	FindByName(name string) (*entities.BasicDomain, error)
	CountUsers(domainID string) (int64, error)
	List(pager entities.Pager, sorter entities.Sorter) (*entities.DomainCollection, error)
	ListByUser(userID string, pager entities.Pager, sorter entities.Sorter) (*entities.DomainCollection, error)
}

// DomainInteractorImpl is an actual interactor that implements DomainInteractor
type DomainInteractorImpl struct {
	DB *sqlx.DB
}

// Create creates a new domain with a given name and description
func (inter *DomainInteractorImpl) Create(domain entities.BasicDomain) error {
	if ok, err := domain.IsValid(); !ok {
		return fmt.Errorf("Domain is not valid: %v", err.Error())
	}
	d := dl.Domain{
		ID:          domain.ID,
		Name:        domain.Name,
		Description: domain.Description,
		Enabled:     domain.Enabled,
	}
	_, err := dl.SaveDomain(inter.DB, d)
	if err != nil {
		return err
	}
	return nil
}

// Update updates all attributes of a given domain entity in the database
func (inter *DomainInteractorImpl) Update(domain entities.BasicDomain) error {
	if ok, err := domain.IsValid(); !ok {
		return fmt.Errorf("Domain is not valid: %v", err.Error())
	}
	d := dl.Domain{
		ID:          domain.ID,
		Name:        domain.Name,
		Description: domain.Description,
		Enabled:     domain.Enabled,
	}
	_, err := dl.SaveDomain(inter.DB, d)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes domain and all assigned entities from storage
func (inter *DomainInteractorImpl) Delete(domain entities.BasicDomain) error {
	err := dl.DeleteDomain(inter.DB, domain.ID)
	if err != nil {
		return err
	}
	return nil
}

// Find finds a domain by given domain ID
func (inter *DomainInteractorImpl) Find(id string) (*entities.BasicDomain, error) {
	r, err := dl.FindDomain(inter.DB, id)
	if err != nil {
		return nil, err
	}
	return basicDomainRecordToEntity(r), nil
}

// FindByName finds a domain by given domain name
func (inter *DomainInteractorImpl) FindByName(name string) (*entities.BasicDomain, error) {
	r, err := dl.FindDomainByName(inter.DB, name)
	if err != nil {
		return nil, err
	}
	return basicDomainRecordToEntity(r), nil
}

// CountUsers return number of users in a domain defined by given domain ID
func (inter *DomainInteractorImpl) CountUsers(domainID string) (int64, error) {
	c, err := dl.CountUsersByDomain(inter.DB, domainID)
	if err != nil {
		return -1, err
	}
	return c, nil
}

// List implements a paginated listing of domains
func (inter *DomainInteractorImpl) List(pager entities.Pager, sorter entities.Sorter) (*entities.DomainCollection, error) {
	total, err := dl.CountDomains(inter.DB)
	if err != nil {
		return nil, err
	}
	records, err := dl.FindAllDomains(inter.DB, pager, sorter)
	if err != nil {
		return nil, err
	}
	c := &entities.DomainCollection{
		Domains:   []*entities.Domain{},
		Paginator: pager.CreatePaginator(len(records), total),
	}
	for _, dto := range records {
		c.Domains = append(c.Domains, domainRecordToEntity(dto))
	}
	return c, nil
}

// ListByUser implements a paginated listing of domains filtered by a given user ID
func (inter *DomainInteractorImpl) ListByUser(userID string, pager entities.Pager, sorter entities.Sorter) (*entities.DomainCollection, error) {
	total, err := dl.CountDomainsByUser(inter.DB, userID)
	if err != nil {
		return nil, err
	}
	records, err := dl.FindDomainsByUser(inter.DB, userID, pager, sorter)
	if err != nil {
		return nil, err
	}
	c := &entities.DomainCollection{
		Domains:   []*entities.Domain{},
		Paginator: pager.CreatePaginator(len(records), total),
	}
	for _, dto := range records {
		c.Domains = append(c.Domains, domainRecordToEntity(dto))
	}
	return c, nil
}

func basicDomainRecordToEntity(record *dl.Domain) *entities.BasicDomain {
	d := entities.NewBasicDomain(record.Name, record.Description)
	d.ID = record.ID
	d.Enabled = record.Enabled
	return d
}

func domainRecordToEntity(record *dl.Domain) *entities.Domain {
	d := new(entities.Domain)
	basicDomain := basicDomainRecordToEntity(record)
	d.BasicDomain = *basicDomain
	d.UsersCount = record.UsersCount
	return d
}
