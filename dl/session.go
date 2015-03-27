package dl

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/entities"
)

// Session DTO
type Session struct {
	ID         string    `db:"user_session_id"`
	UserAgent  string    `db:"user_agent"`
	RemoteAddr string    `db:"remote_addr"`
	CreatedOn  time.Time `db:"created_on"`
	UpdatedOn  time.Time `db:"updated_on"`
	ExpiresOn  time.Time `db:"expires_on"`

	// Fields resulted as join to domain table
	DomainPK      int64  `db:"domain_pk"`
	DomainID      string `db:"domain_id"`
	DomainName    string `db:"domain_name"`
	DomainEnabled bool   `db:"domain_enabled"`

	// Fields resulted as join to user table
	UserPK      int64  `db:"user_pk"`
	UserID      string `db:"user_id"`
	UserName    string `db:"user_name"`
	UserEnabled bool   `db:"user_enabled"`
}

// CreateSession create a new session record based on a given dTO
func CreateSession(db sqlx.Ext, s Session) (*Session, error) {
	now := time.Now().UTC()
	q := fmt.Sprintf(`INSERT INTO user_session (
		user_session_id,
		domain_id,
		user_id,
		user_agent,
		remote_addr,
		created_on,
		updated_on,
		expires_on
	) SELECT ?, d.domain_id, u.user_id, ?, ?, ?, ?, ?
		FROM domain AS d, %v AS u
		WHERE d.object_id = ? AND u.object_id = ?;`, escapeLiteral(db, "user"))
	_, err := db.Exec(db.Rebind(q), s.ID, s.UserAgent, s.RemoteAddr, now, now, s.ExpiresOn, s.DomainID, s.UserID)
	if err != nil {
		return nil, err
	}

	created, err := FindSession(db, s.ID)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// UpdateSession updates existing session record
func UpdateSession(db sqlx.Ext, s Session) error {
	var (
		q   string
		err error
		now = time.Now().UTC()
	)
	s.UpdatedOn = now
	q = `UPDATE user_session SET
			domain_id = ?,
			user_id = ?,
			user_agent = ?,
			remote_addr = ?,
			created_on = ?,
			updated_on = ?,
			expires_on = ?
		WHERE user_session_id = ?;`
	_, err = db.Exec(db.Rebind(q), s.DomainID, s.UserID, s.UserAgent, s.RemoteAddr, s.CreatedOn, s.UpdatedOn, s.ExpiresOn, s.ID)
	return err
}

// RetainSession sets new expires_on attribute for a given session
func RetainSession(db sqlx.Ext, id string, expiresOn time.Time) error {
	r, err := db.Exec(db.Rebind("UPDATE user_session SET expires_on = ? WHERE user_session_id = ?"), expiresOn, id)
	if err != nil {
		return err
	}
	aff, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteSession deletes a session from database cascading
func DeleteSession(db sqlx.Ext, id string) error {
	r, err := db.Exec(db.Rebind("DELETE FROM user_session WHERE user_session_id = ?"), id)
	if err != nil {
		return err
	}
	aff, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteExpiredSessions purges expired sessions from database
func DeleteExpiredSessions(db sqlx.Ext) error {
	now := time.Now().UTC()
	_, err := db.Exec(db.Rebind("DELETE FROM user_session WHERE expires_on <= ?"), now)
	return err
}

// CountSessions returns a total count of session records in database
func CountSessions(db sqlx.Ext) (int64, error) {
	var count int64
	err := db.QueryRowx("SELECT count(*) FROM user_session;").Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// FindAllSessions returns a page of domain records
func FindAllSessions(db sqlx.Ext, pager entities.Pager, sorter entities.Sorter) ([]*Session, error) {
	q := fmt.Sprintf(`SELECT s.*,
			d.domain_id AS domain_pk, d.object_id AS domain_id, d.name AS domain_name, d.is_enabled AS domain_enabled,
			u.user_id AS user_pk, u.object_id AS user_id, u.name AS user_name, u.is_enabled AS user_enabled
		FROM user_session AS s
        LEFT JOIN %v AS u ON s.user_id=u.user_id
        LEFT JOIN domain AS d ON d.domain_id=s.domain_id
        %v %v;`, escapeLiteral(db, "user"), orderByClause(sorter, "s"), limitOffset(pager))
	rows, err := db.Queryx(q)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := []*Session{}
	for rows.Next() {
		var s Session
		err = rows.StructScan(&s)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &s)
	}

	return sessions, nil
}

// FindSession finds a session by given session ID
func FindSession(db sqlx.Ext, id string) (*Session, error) {
	var s Session
	q := fmt.Sprintf(`SELECT s.*,
			d.domain_id AS domain_pk, d.object_id AS domain_id, d.name AS domain_name, d.is_enabled AS domain_enabled,
			u.user_id AS user_pk, u.object_id AS user_id, u.name AS user_name, u.is_enabled AS user_enabled
		FROM user_session AS s
        LEFT JOIN %v AS u ON s.user_id=u.user_id
        LEFT JOIN domain AS d ON d.domain_id=s.domain_id
        WHERE s.user_session_id = ?
        LIMIT 1;`, escapeLiteral(db, "user"))
	err := db.QueryRowx(db.Rebind(q), id).StructScan(&s)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &s, nil
}

// FindUserSpecificSession finds a session by given session ID, user agent and remote address
func FindUserSpecificSession(db sqlx.Ext, userID, domainID string, userAgent, remoteAddr string) (*Session, error) {
	var s Session
	q := fmt.Sprintf(`SELECT s.*,
			d.domain_id AS domain_pk, d.object_id AS domain_id, d.name AS domain_name, d.is_enabled AS domain_enabled,
			u.user_id AS user_pk, u.object_id AS user_id, u.name AS user_name, u.is_enabled AS user_enabled
		FROM user_session AS s
        LEFT JOIN %v AS u ON s.user_id=u.user_id
        LEFT JOIN domain AS d ON d.domain_id=s.domain_id
        WHERE u.object_id = ? AND d.object_id = ? AND s.user_agent = ? AND s.remote_addr = ?
        AND s.expires_on > ?
        LIMIT 1;`, escapeLiteral(db, "user"))
	now := time.Now().UTC()
	err := db.QueryRowx(db.Rebind(q), userID, domainID, userAgent, remoteAddr, now).StructScan(&s)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &s, nil
}
