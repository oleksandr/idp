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
	return nil, fmt.Errorf("MAKE SURE DOMAIN + USER ARE ENABLED")
	/*
		var (
			q   string
			err error
			now = time.Now()
		)

		q = `WITH
			u(id) AS (SELECT user_id FROM user WHERE object_id = ? LIMIT 1),
			d(id) AS (SELECT domain_id FROM domain WHERE object_id = ? LIMIT 1)
			INSERT INTO user_session (
				user_session_id,
				domain_id,
				user_id,
				user_agent,
				remote_addr,
				created_on,
				updated_on,
				expires_on
			)
			SELECT ?, d.id, u.id, ?, ?, ?, ?, ? FROM u, d;`
		_, err = db.Exec(q, s.UserID, s.DomainID, s.ID, s.UserAgent, s.RemoteAddr, now, now, s.ExpiresOn)
		if err != nil {
			return nil, err
		}

		created, err := FindSession(db, s.ID)
		if err != nil {
			return nil, err
		}

		return created, nil
	*/
}

// UpdateSession updates existing session record
func UpdateSession(db sqlx.Ext, s Session) error {
	var (
		q   string
		err error
		now = time.Now()
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
	_, err = db.Exec(q, s.DomainID, s.UserID, s.UserAgent, s.RemoteAddr, s.CreatedOn, s.UpdatedOn, s.ExpiresOn, s.ID)
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
        LEFT JOIN user AS u ON s.user_id=u.user_id
        LEFT JOIN domain AS d ON d.domain_id=s.domain_id
        %v %v;`, orderByClause(sorter, "s"), limitOffset(pager))
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
	q := `SELECT s.*,
			d.domain_id AS domain_pk, d.object_id AS domain_id, d.name AS domain_name, d.is_enabled AS domain_enabled,
			u.user_id AS user_pk, u.object_id AS user_id, u.name AS user_name, u.is_enabled AS user_enabled
		FROM user_session AS s
        LEFT JOIN user AS u ON s.user_id=u.user_id
        LEFT JOIN domain AS d ON d.domain_id=s.domain_id
        WHERE s.user_session_id = ?
        LIMIT 1;`
	err := db.QueryRowx(q, id).StructScan(&s)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &s, nil
}
