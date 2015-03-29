package usecases

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/db"
	"github.com/oleksandr/idp/entities"
	"gopkg.in/gorp.v1"
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
	FindUserSpecific(userID, domainID, userAgent, remoteAddr string) (*entities.Session, error)
	List(pager entities.Pager, sorter entities.Sorter) (*entities.SessionCollection, error)
}

// SessionInteractorImpl is an actual interactor that implements SessionInteractor
type SessionInteractorImpl struct {
	DBMap *gorp.DbMap
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

	// Check if user is assigned to a domain
	userTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	var du db.DomainUser
	q := fmt.Sprintf(`SELECT du.domain_id, du.user_id FROM domain_user AS du
		   		LEFT JOIN %v AS u ON du.user_id=u.user_id
		   		LEFT JOIN domain AS d ON du.domain_id=d.domain_id
		   		WHERE u.object_id = ?
		   		AND d.object_id = ?
		   		LIMIT 1;`, userTbl)
	err := inter.DBMap.SelectOne(&du, q, session.User.ID, session.Domain.ID)
	if err != nil {
		return fmt.Errorf("Could not find user in domain: %v", err.Error())
	}

	now := time.Now().UTC()
	d := &db.Session{
		ID:         session.ID,
		UserAgent:  session.UserAgent,
		RemoteAddr: session.RemoteAddr,
		DomainPK:   du.DomainPK,
		UserPK:     du.UserPK,
		CreatedOn:  now,
		UpdatedOn:  now,
		ExpiresOn:  session.ExpiresOn.Time,
	}
	err = inter.DBMap.Insert(d)
	if err != nil {
		return err
	}
	return nil
}

// Retain prolongs session's expiration date/time till given time
func (inter *SessionInteractorImpl) Retain(session entities.Session) error {
	now := time.Now().UTC()
	expiresOn := now.Add(time.Duration(config.SessionTTLMinutes()) * time.Minute)
	r, err := inter.DBMap.Exec("UPDATE session SET expires_on = ?, updated_on = ? WHERE session_id = ?", expiresOn, now, session.ID)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	} else if err != nil {
		return err
	}
	c, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return entities.ErrNotFound
	}
	return err
}

// Delete deletes session from database
func (inter *SessionInteractorImpl) Delete(session entities.Session) error {
	var (
		s   db.Session
		err error
	)
	err = inter.DBMap.SelectOne(&s, "SELECT * FROM session WHERE session_id = ?", session.ID)
	if err == sql.ErrNoRows {
		return entities.ErrNotFound
	} else if err != nil {
		return err
	}
	_, err = inter.DBMap.Delete(&s)
	if err != nil {
		return err
	}
	return nil
}

// Purge purges all expired sessions
func (inter *SessionInteractorImpl) Purge() error {
	now := time.Now().UTC()
	_, err := inter.DBMap.Exec("DELETE FROM session WHERE expires_on <= ?", now)
	return err
}

// Find looks for a session by given session ID
func (inter *SessionInteractorImpl) Find(id string) (*entities.Session, error) {
	var (
		sv      db.SessionView
		err     error
		userTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	)
	q := fmt.Sprintf(`SELECT s.*,
			d.domain_id AS domain_id, d.object_id AS domain_object_id, d.name AS domain_name, d.is_enabled AS domain_enabled,
			u.user_id AS user_id, u.object_id AS user_object_id, u.name AS user_name, u.is_enabled AS user_enabled
		FROM session AS s
        LEFT JOIN %v AS u ON s.user_id=u.user_id
        LEFT JOIN domain AS d ON d.domain_id=s.domain_id
        WHERE s.session_id = ?
        LIMIT 1;`, userTbl)
	err = inter.DBMap.SelectOne(&sv, q, id)
	if err != nil {
		return nil, err
	}
	return sessionToEntity(&sv), nil
}

// FindUserSpecific looks for a session by given session ID, user agent and remote address
func (inter *SessionInteractorImpl) FindUserSpecific(userID, domainID, userAgent, remoteAddr string) (*entities.Session, error) {
	var (
		sv      db.SessionView
		err     error
		userTbl = inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	)
	q := fmt.Sprintf(`SELECT s.*,
			d.domain_id AS domain_id, d.object_id AS domain_object_id, d.name AS domain_name, d.is_enabled AS domain_enabled,
			u.user_id AS user_id, u.object_id AS user_object_id, u.name AS user_name, u.is_enabled AS user_enabled
		FROM session AS s
        LEFT JOIN %v AS u ON s.user_id=u.user_id
        LEFT JOIN domain AS d ON d.domain_id=s.domain_id
        WHERE u.object_id = ? AND d.object_id = ? AND s.user_agent = ? AND s.remote_addr = ?
        LIMIT 1;`, userTbl)
	err = inter.DBMap.SelectOne(&sv, q, userID, domainID, userAgent, remoteAddr)
	if err != nil {
		return nil, err
	}
	return sessionToEntity(&sv), nil
}

// List implements a paginated listing of session
func (inter *SessionInteractorImpl) List(pager entities.Pager, sorter entities.Sorter) (*entities.SessionCollection, error) {
	total, err := inter.DBMap.SelectInt("SELECT COUNT(*) FROM session")
	if err != nil {
		return nil, err
	}
	var records []db.SessionView
	userTbl := inter.DBMap.Dialect.QuotedTableForQuery("", "user")
	q := fmt.Sprintf(`SELECT s.*,
			d.domain_id AS domain_id, d.object_id AS domain_object_id, d.name AS domain_name, d.is_enabled AS domain_enabled,
			u.user_id AS user_id, u.object_id AS user_object_id, u.name AS user_name, u.is_enabled AS user_enabled
		FROM session AS s
        LEFT JOIN %v AS u ON s.user_id=u.user_id
        LEFT JOIN domain AS d ON d.domain_id=s.domain_id
        %v %v;`, userTbl, db.OrderByClause(sorter, "s"), db.LimitOffset(pager))
	_, err = inter.DBMap.Select(&records, q)
	if err != nil {
		return nil, err
	}
	c := &entities.SessionCollection{
		Sessions:  []entities.Session{},
		Paginator: *pager.CreatePaginator(len(records), total),
	}
	for _, r := range records {
		c.Sessions = append(c.Sessions, *sessionToEntity(&r))
	}
	return c, nil
}

func sessionToEntity(s *db.SessionView) *entities.Session {
	e := &entities.Session{
		ID: s.ID,
		Domain: &entities.BasicDomain{
			ID:      s.DomainID,
			Name:    s.DomainName,
			Enabled: s.DomainEnabled,
		},
		User: &entities.BasicUser{
			ID:      s.UserID,
			Name:    s.UserName,
			Enabled: s.UserEnabled,
		},
		UserAgent:  s.UserAgent,
		RemoteAddr: s.RemoteAddr,
	}
	e.CreatedOn.Time = s.CreatedOn
	e.UpdatedOn.Time = s.UpdatedOn
	e.ExpiresOn.Time = s.ExpiresOn
	return e
}
