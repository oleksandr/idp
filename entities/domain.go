package entities

import (
	"fmt"

	"github.com/satori/go.uuid"
)

// BasicDomain contains basic domain attributes
type BasicDomain struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// Domain entities represent different tenants which contain User entities.
// A domain can be either a real domain (e.g myproject.com) or just a logical
// name of the user group you want to separate everyone into.
type Domain struct {
	BasicDomain
	UsersCount int64 `json:"-"`
}

// NewBasicDomain - a constructor for Domain entities
func NewBasicDomain(name, description string) *BasicDomain {
	d := new(BasicDomain)
	d.ID = uuid.NewV4().String()
	d.Name = name
	d.Description = description
	d.Enabled = true
	return d
}

// IsValid checks if domain is valid
func (d *BasicDomain) IsValid() (bool, error) {
	if d.Name == "" {
		return false, fmt.Errorf("Name cannot be empty!")
	}
	return true, nil
}

// BasicDomainCollection is a paginated collection of Domain entities
type BasicDomainCollection struct {
	Domains   []*BasicDomain `json:"domains"`
	Paginator *Paginator     `json:"paginator"`
}

// DomainCollection is a paginated collection of Domain entities
type DomainCollection struct {
	Domains   []*Domain  `json:"domains"`
	Paginator *Paginator `json:"paginator"`
}
