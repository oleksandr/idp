package rpc

import (
	"log"
	"os"

	"github.com/oleksandr/idp/usecases"
)

// IdentityProviderHandler handles implements IdentityProvider RPC interface
type IdentityProviderHandler struct {
	log               *log.Logger
	SessionInteractor usecases.SessionInteractor
	UserInteractor    usecases.UserInteractor
	DomainInteractor  usecases.DomainInteractor
	RBACInteractor    usecases.RBACInteractor
}

// NewIdentityProviderHandler creates new IdentityProviderHandler
func NewIdentityProviderHandler() *IdentityProviderHandler {
	return &IdentityProviderHandler{
		log: log.New(os.Stdout, "[rpc] ", log.LstdFlags),
	}
}
