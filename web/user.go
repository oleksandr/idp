package web

import (
	"net/http"

	"github.com/oleksandr/idp/usecases"
)

//
// UserWebHandler is a collection of CRUD methods for Users API
//
type UserWebHandler struct {
	UserInteractor usecases.UserInteractor
}

func (handler *UserWebHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *UserWebHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *UserWebHandler) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *UserWebHandler) Modify(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (handler *UserWebHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
