package usecases

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/dl"
	"github.com/oleksandr/idp/entities"
)

//
// SessionInteractor is an interface that defines all session related use-cases
// signatures
//
type SessionInteractor interface {
	Create(session entities.Session) error
	//Update(user entities.User) error
	//Delete(user entities.User) error
	//Find(id string) (*entities.User, error)
	List(pager entities.Pager, sorter entities.Sorter) (*entities.SessionCollection, error)
	//ListByDomain(domainID string, pager entities.Pager, sorter entities.Sorter) (*entities.UserCollection, error)
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
	if !session.Domain.Enabled {
		return fmt.Errorf("Domain is not enabled")
	}
	if !session.User.Enabled {
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
