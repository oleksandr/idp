package entities

import (
	"fmt"

	"github.com/satori/go.uuid"
)

//
// Domain entities represent different tenants which contain User entities.
// A domain can be either a real domain (e.g myproject.com) or just a logical
// name of the user group you want to separate everyone into.
//
type Domain struct {
	ID          string `json:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	UsersCount  int64  `json:"-"`
}

// NewDomain - a constructor for Domain entities
func NewDomain(name, description string) *Domain {
	d := new(Domain)
	d.ID = uuid.NewV4().String()
	d.Name = name
	d.Description = description
	d.Enabled = true
	return d
}

// IsValid checks if domain is valid
func (d *Domain) IsValid() (bool, error) {
	if d.Name == "" {
		return false, fmt.Errorf("Name cannot be empty!")
	}
	return true, nil
}

//
// DomainCollection is a paginated collection of Domain entities
//
type DomainCollection struct {
	Domains   []*Domain  `json:"domains"`
	Paginator *Paginator `json:"paginator"`
}
