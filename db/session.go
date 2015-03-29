package db

import "time"

// Session Table
type Session struct {
	ID         string    `db:"session_id"`
	DomainPK   int64     `db:"domain_id"`
	UserPK     int64     `db:"user_id"`
	UserAgent  string    `db:"user_agent"`
	RemoteAddr string    `db:"remote_addr"`
	CreatedOn  time.Time `db:"created_on"`
	UpdatedOn  time.Time `db:"updated_on"`
	ExpiresOn  time.Time `db:"expires_on"`
}

// SessionView contains all fields for populating the entity
type SessionView struct {
	Session
	// Fields resulted as join to domain table
	DomainID      string `db:"domain_object_id"`
	DomainName    string `db:"domain_name"`
	DomainEnabled bool   `db:"domain_enabled"`
	// Fields resulted as join to user table
	UserID      string `db:"user_object_id"`
	UserName    string `db:"user_name"`
	UserEnabled bool   `db:"user_enabled"`
}
