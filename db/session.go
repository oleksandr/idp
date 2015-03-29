package db

import "time"

// Session Table
type Session struct {
	ID         string    `db:"session_id"`
	DomainID   string    `db:"domain_id"`
	UserID     string    `db:"user_id"`
	UserAgent  string    `db:"user_agent"`
	RemoteAddr string    `db:"remote_addr"`
	CreatedOn  time.Time `db:"created_on"`
	UpdatedOn  time.Time `db:"updated_on"`
	ExpiresOn  time.Time `db:"expires_on"`
}
