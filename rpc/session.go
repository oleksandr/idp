package rpc

import (
	"log"
	"os"
	"time"

	"github.com/oleksandr/idp/entities"
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
func (handler *AuthenticatorHandler) CreateSession(domain string, name string, password string, userAgent string, remoteAddr string) (r *services.Session, err error) {
	handler.log.Printf("createSession(%v, %v)", domain, name)

	d, err := handler.DomainInteractor.FindByName(domain)
	if err != nil {
		e := services.NewBadRequest()
		e.Code = "0000"
		e.Msg = e.Error()
		return nil, e
	}

	user, err := handler.UserInteractor.FindByNameInDomain(name, d.ID)
	if err != nil {
		e := services.NewBadRequest()
		e.Code = "0000"
		e.Msg = e.Error()
		return nil, e
	}

	if !user.IsPassword(password) {
		e := services.NewForbidden()
		e.Code = "0000"
		e.Msg = e.Error()
		return nil, e
	}

	session, err := handler.SessionInteractor.FindUserSpecific(user.ID, d.ID, userAgent, remoteAddr)
	if session != nil && !session.IsExpired() {
		handler.SessionInteractor.Retain(*session)
		return sessionToResponse(session), nil
	}

	// Create new session
	session = entities.NewSession(*user, *d, userAgent, remoteAddr)
	err = handler.SessionInteractor.Create(*session)
	if err != nil {
		e := services.NewForbidden()
		e.Code = "0000"
		e.Msg = e.Error()
		return nil, e
	}

	return sessionToResponse(session), nil
}

// CheckSession is an implementation of Authentocator's CheckSession method
func (handler *AuthenticatorHandler) CheckSession(sessionID string) (r bool, err error) {
	handler.log.Println("checkSession(%v)", sessionID)

	session, err := handler.SessionInteractor.Find(sessionID)
	if err != nil {
		e := services.NewBadRequest()
		e.Code = "0000"
		e.Msg = e.Error()
		return false, e
	}

	if session.IsExpired() {
		e := services.NewForbidden()
		e.Code = "0000"
		e.Msg = e.Error()
		return false, e
	}

	return true, nil
}

// DeleteSession is an implementation of Authentocator's DeleteSession method
func (handler *AuthenticatorHandler) DeleteSession(sessionID string) (r bool, err error) {
	handler.log.Println("deleteSession(%v)", sessionID)

	session, err := handler.SessionInteractor.Find(sessionID)
	if err != nil {
		e := services.NewBadRequest()
		e.Code = "0000"
		e.Msg = e.Error()
		return false, e
	}

	err = handler.SessionInteractor.Delete(*session)
	if err != nil {
		e := services.NewBadRequest()
		e.Code = "0000"
		e.Msg = e.Error()
		return false, e
	}

	return true, nil
}

func sessionToResponse(s *entities.Session) *services.Session {
	r := services.NewSession()
	r.Id = s.ID
	r.Domain = &services.Domain{
		Id:          s.Domain.ID,
		Name:        s.Domain.Name,
		Description: s.Domain.Description,
		Enabled:     s.Domain.Enabled,
	}
	r.User = &services.User{
		Id:      s.User.ID,
		Name:    s.User.Name,
		Enabled: s.User.Enabled,
	}
	r.UserAgent = s.UserAgent
	r.RemoteAddr = s.RemoteAddr
	r.CreatedOn = s.CreatedOn.Format(time.RFC3339)
	r.UpdatedOn = s.UpdatedOn.Format(time.RFC3339)
	r.ExpiresOn = s.ExpiresOn.Format(time.RFC3339)
	return r
}
