package entities

import (
	"time"

	"github.com/oleksandr/idp/config"
	"github.com/satori/go.uuid"
)

//
// Session structure represents user's time limited session
//
type Session struct {
	ID         string       `json:"sid"`
	Domain     *BasicDomain `json:"domain"`
	User       *BasicUser   `json:"user"`
	UserAgent  string       `json:"-"`
	RemoteAddr string       `json:"-"`
	CreatedOn  Time         `json:"createdOn"`
	UpdatedOn  Time         `json:"updatedOn"`
	ExpiresOn  Time         `json:"expiresOn"`
}

// NewSession create a new Session entity
func NewSession(user BasicUser, domain BasicDomain, userAgent, remoteAddr string) *Session {
	s := &Session{
		ID:         uuid.NewV4().String(),
		Domain:     &domain,
		User:       &user,
		UserAgent:  userAgent,
		RemoteAddr: remoteAddr,
	}
	now := time.Now()
	s.CreatedOn.Time = now
	s.UpdatedOn.Time = now
	s.ExpiresOn.Time = now.Add(time.Duration(config.SessionTTLMinutes) * time.Minute)
	return s
}

// IsValid checks if session has a non-empty Sid, non-null User's Id and
// a non-zero expiration time.
func (s *Session) IsValid() bool {
	return s.ID != "" && s.User.ID != "" && s.User.Enabled && !s.ExpiresOn.IsZero()
}

// IsExpired checks if the session is expired
func (s *Session) IsExpired() bool {
	return s.ExpiresOn.Sub(time.Now()) <= 0
}

//
// SessionCollection is a paginated collection of Session entities
//
type SessionCollection struct {
	Sessions  []*Session `json:"sessions"`
	Paginator *Paginator `json:"paginator"`
}