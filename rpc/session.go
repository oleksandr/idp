package rpc

import (
	"log"
	"os"
	"time"

	"github.com/oleksandr/idp/entities"
	"github.com/oleksandr/idp/errs"
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

	// Prepare arguments
	u := entities.BasicUser{}
	u.Name = name
	d := entities.BasicDomain{}
	d.Name = domain

	// Create session
	session, err := handler.SessionInteractor.CreateWithPassword(d, u, password, userAgent, remoteAddr)
	if err == nil {
		return sessionToResponse(session), nil
	}

	// Handle errors
	e := err.(*errs.Error)
	return nil, errorToServiceError(e)
}

// CheckSession is an implementation of Authentocator's CheckSession method
func (handler *AuthenticatorHandler) CheckSession(sessionID string, userAgent string, remoteAddr string) (r bool, err error) {
	handler.log.Printf("checkSession(%v)", sessionID)

	session, err := handler.SessionInteractor.Find(sessionID)
	if err != nil {
		e := err.(*errs.Error)
		return false, errorToServiceError(e)
	}

	if !session.Domain.Enabled || !session.User.Enabled {
		e := services.NewForbiddenError()
		e.Msg = "Domain and/or user disabled"
		return false, e
	}

	if session.UserAgent != userAgent || session.RemoteAddr != remoteAddr {
		e := services.NewNotFoundError()
		e.Msg = "Session not found"
		return false, e
	}

	if session.IsExpired() {
		e := services.NewForbiddenError()
		e.Msg = "Session expired"
		return false, e
	}

	err = handler.SessionInteractor.Retain(*session)
	if err != nil {
		e := err.(*errs.Error)
		return false, errorToServiceError(e)
	}

	return true, nil
}

// DeleteSession is an implementation of Authentocator's DeleteSession method
func (handler *AuthenticatorHandler) DeleteSession(sessionID string, userAgent string, remoteAddr string) (r bool, err error) {
	handler.log.Printf("deleteSession(%v)", sessionID)

	session, err := handler.SessionInteractor.Find(sessionID)
	if err != nil {
		e := err.(*errs.Error)
		return false, errorToServiceError(e)
	}

	if !session.Domain.Enabled || !session.User.Enabled {
		e := services.NewForbiddenError()
		e.Msg = "Domain and/or user disabled"
		return false, e
	}

	if session.UserAgent != userAgent || session.RemoteAddr != remoteAddr {
		e := services.NewNotFoundError()
		e.Msg = "Session not found"
		return false, e
	}

	if session.IsExpired() {
		e := services.NewForbiddenError()
		e.Msg = "Session expired"
		return false, e
	}

	err = handler.SessionInteractor.Delete(*session)
	if err != nil {
		e := err.(*errs.Error)
		return false, errorToServiceError(e)
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
