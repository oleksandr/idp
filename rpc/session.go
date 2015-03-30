package rpc

import (
	"fmt"
	"log"
	"os"

	"github.com/oleksandr/idp/rpc/generated/services"
)

type AuthenticatorHandler struct {
	log *log.Logger
}

func NewAuthenticatorHandler() *AuthenticatorHandler {
	return &AuthenticatorHandler{
		log: log.New(os.Stdout, "[rpc] ", log.LstdFlags),
	}
}

func (handler *AuthenticatorHandler) CreateSession(domainID string, name string, password string) (r *services.Session, err error) {
	handler.log.Println("CreateSession")
	e := services.NewBadRequest()
	e.Msg = "NOT IMPLEMENTED YET"
	return nil, e
}

func (handler *AuthenticatorHandler) CheckSession(sessionID string) (r bool, err error) {
	handler.log.Println("CheckSession")
	e := services.NewBadRequest()
	e.Msg = "NOT IMPLEMENTED YET"
	return false, e
}

func (handler *AuthenticatorHandler) DeleteSession(sessionID string) (r bool, err error) {
	handler.log.Println("DeleteSession")
	return false, fmt.Errorf("This is kind of an unknown error")
}
