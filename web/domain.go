package web

import (
	"net/http"

	"github.com/oleksandr/idp/usecases"
)

//
// DomainWebHandler is a collection of CRUD methods for Users API
//
type DomainWebHandler struct {
	DomainInteractor usecases.DomainInteractor
}

func (handler *DomainWebHandler) Create(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *DomainWebHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *DomainWebHandler) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *DomainWebHandler) Modify(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *DomainWebHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
