package usecases

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/db"
	"github.com/oleksandr/idp/entities"
	"github.com/oleksandr/idp/errs"
	"gopkg.in/gorp.v1"
)

//
// SessionInteractor is an interface that defines all session related use-cases
// signatures
//
type SessionInteractor interface {
	Create(domain entities.BasicDomain, user entities.BasicUser, userAgent string, remoteAddr string) (*entities.Session, error)
	CreateWithPassword(domain entities.BasicDomain, user entities.BasicUser, password, userAgent, remoteAddr string) (*entities.Session, error)
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

func (inter *SessionInteractorImpl) create(domain entities.BasicDomain, user entities.BasicUser, checkPwd bool, password string, userAgent, remoteAddr string) (*entities.Session, error) {
	var (
		err     error
		d       *db.Domain
		u       *db.User
		s       *db.Session
		session *entities.Session
	)

	// Check/find domain
	if domain.ID != "" {
		d, err = findDomainByID(inter.DBMap, domain.ID)
	} else if domain.Name != "" {
		d, err = findDomainByName(inter.DBMap, domain.Name)
	} else {
		err = errs.NewUseCaseError(errs.ErrorTypeConflict, "You need to provide domain ID or name", nil)
	}
	if err != nil {
		return nil, err
	}

	// Check/find user
	if user.ID != "" {
		u, err = findUserByID(inter.DBMap, user.ID)
	} else if user.Name != "" {
		u, err = findUserByName(inter.DBMap, user.Name)
	} else {
		err = errs.NewUseCaseError(errs.ErrorTypeConflict, "You need to provide user ID or name", nil)
	}
	if err != nil {
		return nil, err
	}

	// Check if domain/user are enabled
	if !d.Enabled {
		return nil, errs.NewUseCaseError(errs.ErrorTypeForbidden, "Domain is disabled", nil)
	}
	if !u.Enabled {
		return nil, errs.NewUseCaseError(errs.ErrorTypeForbidden, "User is disabled", nil)
	}

	// Password check
	if checkPwd {
		basicUser := userToEntity(u)
		if !basicUser.IsPassword(password) {
			return nil, errs.NewUseCaseError(errs.ErrorTypeForbidden, "Invalid passowrd", nil)
		}
	}

	// Check if user is assigned to a domain
	_, err = findUserInDomain(inter.DBMap, u.ID, d.ID)
	if err != nil {
		e := err.(*errs.Error)
		if e.Type == errs.ErrorTypeOperational {
			return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Error checking user's domain", err)
		}
		return nil, errs.NewUseCaseError(errs.ErrorTypeForbidden, "User is not in domain", err)
	}

	// Lookup existing
	session, err = inter.FindUserSpecific(u.ID, d.ID, userAgent, remoteAddr)
	if session != nil && !session.IsExpired() {
		err = inter.Retain(*session)
		if err != nil {
			return nil, err
		}
		return session, nil
	}

	// Create new session
	session = entities.NewSession(*userToEntity(u), *domainToEntity(d), userAgent, remoteAddr)
	s = &db.Session{
		ID:         session.ID,
		UserAgent:  session.UserAgent,
		RemoteAddr: session.RemoteAddr,
		DomainPK:   d.PK,
		UserPK:     u.PK,
		CreatedOn:  session.CreatedOn.Time,
		UpdatedOn:  session.UpdatedOn.Time,
		ExpiresOn:  session.ExpiresOn.Time,
	}
	err = inter.DBMap.Insert(s)
	if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to create a session", err)
	}
	return session, nil
}

// Create a user session for a given domain, user and user's agent with remote address
func (inter *SessionInteractorImpl) Create(domain entities.BasicDomain, user entities.BasicUser, userAgent string, remoteAddr string) (*entities.Session, error) {
	return inter.create(domain, user, false, "", userAgent, remoteAddr)
}

// CreateWithPassword is the same as Create() but also performs password checks
func (inter *SessionInteractorImpl) CreateWithPassword(domain entities.BasicDomain, user entities.BasicUser, password, userAgent, remoteAddr string) (*entities.Session, error) {
	return inter.create(domain, user, true, password, userAgent, remoteAddr)
}

// Retain prolongs session's expiration date/time till given time
func (inter *SessionInteractorImpl) Retain(session entities.Session) error {
	now := time.Now().UTC()
	expiresOn := now.Add(time.Duration(config.SessionTTLMinutes()) * time.Minute)

	r, err := inter.DBMap.Exec("UPDATE session SET expires_on = ?, updated_on = ? WHERE session_id = ?", expiresOn, now, session.ID)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to retain a session", err)
	}

	_, err = r.RowsAffected()
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to retain a session", err)
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
		return errs.NewUseCaseError(errs.ErrorTypeNotFound, "Session not found by given ID", err)
	} else if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a session", err)
	}

	_, err = inter.DBMap.Delete(&s)
	if err != nil {
		return errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to delete session", err)
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
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "Session not found by given ID", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a session", err)
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
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "Session not found by given ID", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of a session", err)
	}

	return sessionToEntity(&sv), nil
}

// List implements a paginated listing of session
func (inter *SessionInteractorImpl) List(pager entities.Pager, sorter entities.Sorter) (*entities.SessionCollection, error) {
	total, err := inter.DBMap.SelectInt("SELECT COUNT(*) FROM session")
	if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to count sessions", err)
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
	if err == sql.ErrNoRows {
		return nil, errs.NewUseCaseError(errs.ErrorTypeNotFound, "No sessions found", err)
	} else if err != nil {
		return nil, errs.NewUseCaseError(errs.ErrorTypeOperational, "Failed to perform a lookup of sessions", err)
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
