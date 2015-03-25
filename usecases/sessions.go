package usecases

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/dl"
	"github.com/oleksandr/idp/entities"
)

//
// SessionInteractor is an interface that defines all session related use-cases
// signatures
//
type SessionInteractor interface {
	Create(session entities.Session) error
	Retain(session entities.Session) error
	Delete(session entities.Session) error
	Purge() error
	Find(id string) (*entities.Session, error)
	FindSpecific(id, userAgent, remoteAddr string) (*entities.Session, error)
	List(pager entities.Pager, sorter entities.Sorter) (*entities.SessionCollection, error)
}

// SessionInteractorImpl is an actual interactor that implements SessionInteractor
type SessionInteractorImpl struct {
	DB *sqlx.DB
}

// Create open a new session
func (inter *SessionInteractorImpl) Create(session entities.Session) error {
	if !session.IsValid() {
		return fmt.Errorf("Session is not valid")
	}
	if session.IsExpired() {
		return fmt.Errorf("Session is expired")
	}
	if session.Domain == nil || !session.Domain.Enabled {
		return fmt.Errorf("Domain is not enabled")
	}
	if session.User == nil || !session.User.Enabled {
		return fmt.Errorf("User is not enabled")
	}
	_, err := dl.FindUserInDomain(inter.DB, session.User.ID, session.Domain.ID)
	if err != nil {
		return fmt.Errorf("Could not find user in domain: %v", err.Error())
	}
	s := dl.Session{
		ID:         session.ID,
		UserAgent:  session.UserAgent,
		RemoteAddr: session.RemoteAddr,
		ExpiresOn:  session.ExpiresOn.Time,
		DomainID:   session.Domain.ID,
		UserID:     session.User.ID,
	}
	_, err = dl.CreateSession(inter.DB, s)
	if err != nil {
		return err
	}
	return nil
}

// Retain prolongs session's expiration date/time till given time
func (inter *SessionInteractorImpl) Retain(session entities.Session) error {
	expiresOn := time.Now().UTC().Add(time.Duration(config.SessionTTLMinutes()) * time.Minute)
	err := dl.RetainSession(inter.DB, session.ID, expiresOn)
	return err
}

// Delete deletes session from database
func (inter *SessionInteractorImpl) Delete(session entities.Session) error {
	err := dl.DeleteSession(inter.DB, session.ID)
	if err != nil {
		return err
	}
	return nil
}

// Purge purges all expired sessions
func (inter *SessionInteractorImpl) Purge() error {
	err := dl.DeleteExpiredSessions(inter.DB)
	if err != nil {
		return err
	}
	return nil
}

// Find looks for a session by given session ID
func (inter *SessionInteractorImpl) Find(id string) (*entities.Session, error) {
	r, err := dl.FindSession(inter.DB, id)
	if err != nil {
		return nil, err
	}
	return sessionRecordToEntity(r), nil
}

// FindSpecific looks for a session by given session ID, user agent and remote address
func (inter *SessionInteractorImpl) FindSpecific(id, userAgent, remoteAddr string) (*entities.Session, error) {
	r, err := dl.FindSpecificSession(inter.DB, id, userAgent, remoteAddr)
	if err != nil {
		return nil, err
	}
	return sessionRecordToEntity(r), nil
}

// List implements a paginated listing of session
func (inter *SessionInteractorImpl) List(pager entities.Pager, sorter entities.Sorter) (*entities.SessionCollection, error) {
	total, err := dl.CountSessions(inter.DB)
	if err != nil {
		return nil, err
	}
	records, err := dl.FindAllSessions(inter.DB, pager, sorter)
	if err != nil {
		return nil, err
	}
	c := &entities.SessionCollection{
		Sessions:  []*entities.Session{},
		Paginator: pager.CreatePaginator(len(records), total),
	}
	for _, dto := range records {
		c.Sessions = append(c.Sessions, sessionRecordToEntity(dto))
	}
	return c, nil
}

func sessionRecordToEntity(record *dl.Session) *entities.Session {
	s := &entities.Session{
		ID: record.ID,
		Domain: &entities.BasicDomain{
			ID:      record.DomainID,
			Name:    record.DomainName,
			Enabled: record.DomainEnabled,
		},
		User: &entities.BasicUser{
			ID:      record.UserID,
			Name:    record.UserName,
			Enabled: record.UserEnabled,
		},
		UserAgent:  record.UserAgent,
		RemoteAddr: record.RemoteAddr,
	}
	s.CreatedOn.Time = record.CreatedOn
	s.UpdatedOn.Time = record.UpdatedOn
	s.ExpiresOn.Time = record.ExpiresOn
	return s
}
