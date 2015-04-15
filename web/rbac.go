package web

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/entities"
	"github.com/oleksandr/idp/errs"
	"github.com/oleksandr/idp/usecases"
)

//
// RBACWebHandler is a collection of various methods for RBAC
//
type RBACWebHandler struct {
	log            *log.Logger
	RBACInteractor usecases.RBACInteractor
}

// NewRBACWebHandler creates new SessionWebHandler
func NewRBACWebHandler() *RBACWebHandler {
	return &RBACWebHandler{
		log: log.New(os.Stdout, "[RBACHandler] ", log.LstdFlags),
	}
}

// AssertRole checks if a current user has a specific role
func (handler *RBACWebHandler) AssertRole(w http.ResponseWriter, r *http.Request) {
	var (
		s   entities.Session
		ok  bool
		err error
	)

	if s, ok = context.Get(r, config.CtxSessionKey).(entities.Session); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := context.Get(r, config.CtxParamsKey).(httprouter.Params)
	role := params.ByName("role")
	handler.log.Printf("AssertRole(%v, %v)", s.User.ID, role)

	ok, err = handler.RBACInteractor.AssertRole(s.User.ID, role)
	if err != nil {
		handler.log.Println("ERROR:", err.Error())
		respondWithError(w, errorToHTTPStatus(err.(*errs.Error)), "Failed to assert role", err)
		return
	}

	if ok {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

// AssertPermission checks if a current user has a specific permission
func (handler *RBACWebHandler) AssertPermission(w http.ResponseWriter, r *http.Request) {
	var (
		s   entities.Session
		ok  bool
		err error
	)

	if s, ok = context.Get(r, config.CtxSessionKey).(entities.Session); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := context.Get(r, config.CtxParamsKey).(httprouter.Params)
	permission := params.ByName("permission")
	handler.log.Printf("AssertPermission(%v, %v)", s.User.ID, permission)

	ok, err = handler.RBACInteractor.AssertPermission(s.User.ID, permission)
	if err != nil {
		handler.log.Println("ERROR:", err.Error())
		respondWithError(w, errorToHTTPStatus(err.(*errs.Error)), "Failed to assert permission", err)
		return
	}

	if ok {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
