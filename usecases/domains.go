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
	Create(domain entities.Domain) error
	Update(domain entities.Domain) error
	Delete(domain entities.Domain) error
	Find(id string) (*entities.Domain, error)
	List(pager entities.Pager, sorter entities.Sorter) (*entities.DomainCollection, error)
	ListByUser(userID string, pager entities.Pager, sorter entities.Sorter) (*entities.DomainCollection, error)
}

// DomainInteractorImpl is an actual interactor that implements DomainInteractor
type DomainInteractorImpl struct {
	DB *sqlx.DB
}

// Create creates a new domain with a given name and description
func (inter *DomainInteractorImpl) Create(domain entities.Domain) error {
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
func (inter *DomainInteractorImpl) Update(domain entities.Domain) error {
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
func (inter *DomainInteractorImpl) Delete(domain *entities.Domain) error {
	err := dl.DeleteDomain(inter.DB, domain.ID)
	if err != nil {
		return err
	}
	return nil
}

// Find finds a domain by given domain ID
func (inter *DomainInteractorImpl) Find(id string) (*entities.Domain, error) {
	r, err := dl.FindDomain(inter.DB, id)
	if err != nil {
		return nil, err
	}
	d := entities.NewDomain(r.Name, r.Description)
	d.ID = r.ID
	d.Enabled = r.Enabled
	return d, nil
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

func domainRecordToEntity(record *dl.Domain) *entities.Domain {
	d := entities.NewDomain(record.Name, record.Description)
	d.ID = record.ID
	d.Enabled = record.Enabled
	return d
}
