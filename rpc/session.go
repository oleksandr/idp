package rpc

import (
	"fmt"
	"log"
	"os"

	"github.com/oleksandr/idp/rpc/generated/services"
	"github.com/oleksandr/idp/usecases"
)

// AuthenticatorHandler handles implements Authenticator RPC interface
type AuthenticatorHandler struct {
	log               *log.Logger
	SessionInteractor usecases.SessionInteractor
	UserInteractor    usecases.UserInteractor
	DomainInteractor  usecases.DomainInteractor
}

// NewAuthenticatorHandler creates new AuthenticatorHandler
func NewAuthenticatorHandler() *AuthenticatorHandler {
	return &AuthenticatorHandler{
		log: log.New(os.Stdout, "[rpc] ", log.LstdFlags),
	}
}

// CreateSession is an implementation of Authentocator's CreateSession method
func (handler *AuthenticatorHandler) CreateSession(domainID string, name string, password string) (r *services.Session, err error) {
	handler.log.Println("CreateSession")
	e := services.NewBadRequest()
	e.Msg = "NOT IMPLEMENTED YET"
	return nil, e
}

// CheckSession is an implementation of Authentocator's CheckSession method
func (handler *AuthenticatorHandler) CheckSession(sessionID string) (r bool, err error) {
	handler.log.Println("CheckSession")
	e := services.NewBadRequest()
	e.Msg = "NOT IMPLEMENTED YET"
	return false, e
}

// DeleteSession is an implementation of Authentocator's DeleteSession method
func (handler *AuthenticatorHandler) DeleteSession(sessionID string) (r bool, err error) {
	handler.log.Println("DeleteSession")
	return false, fmt.Errorf("This is kind of an unknown error")
}
