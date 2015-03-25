package web

import (
	"encoding/json"
	"net/http"

	"github.com/oleksandr/idp/entities"
	"github.com/oleksandr/idp/usecases"
)

//
// DomainWebHandler is a collection of CRUD methods for Users API
//
type DomainWebHandler struct {
	DomainInteractor usecases.DomainInteractor
}

// Create handles a creation of a new domain
func (handler *DomainWebHandler) Create(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusNotImplemented)
}

// Retrieve handles a read request of a domain by given id
func (handler *DomainWebHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// List handles a list all entities with support of pagination
func (handler *DomainWebHandler) List(w http.ResponseWriter, r *http.Request) {
	//TODO: read from GET parameters
	pager := entities.Pager{Page: 1, PerPage: 100}
	sorter := entities.Sorter{Field: "name", Asc: true}
	collection, err := handler.DomainInteractor.List(pager, sorter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}
	json.NewEncoder(w).Encode(collection)
}

// Modify handles partial modifications of a domain
func (handler *DomainWebHandler) Modify(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Delete handles a deletion request by a given domain id
func (handler *DomainWebHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
