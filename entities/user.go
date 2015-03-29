package entities

import (
	"crypto/sha1"
	"fmt"

	"github.com/oleksandr/idp/config"
	"github.com/satori/go.uuid"
)

// BasicUser contains basic user attributes
type BasicUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"`
	Enabled  bool   `json:"enabled"`
}

// User entities represent users
type User struct {
	BasicUser
	DomainsCount int64 `json:"-"`
}

// NewBasicUser - a constructor for `User`
func NewBasicUser(name string) *BasicUser {
	a := new(BasicUser)
	a.ID = uuid.NewV4().String()
	a.Name = name
	a.Enabled = true
	return a
}

// SetPassword hashes a given clearTxt and assigns it to password field
func (u *BasicUser) SetPassword(clearTxt string) {
	u.Password = passwordHash(clearTxt)
}

// IsValid checks if user is valid
func (u *BasicUser) IsValid() (bool, error) {
	if u.Name == "" {
		return false, fmt.Errorf("Name cannot be empty!")
	}
	if u.Password == "" {
		return false, fmt.Errorf("Password cannot be empty!")
	}
	return true, nil
}

// IsPassword checks if a given clear text is user's password
func (u *BasicUser) IsPassword(clearTxt string) bool {
	return u.Password == passwordHash(clearTxt)
}

// passwordHash one-way hashes a string with the private HashSecret value.
func passwordHash(s string) string {
	hash := sha1.New()
	hash.Write([]byte(s))
	hash.Write([]byte(config.HashSecretSalt()))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// BasicUserCollection is a paginated collection of User entities
type BasicUserCollection struct {
	Users     []BasicUser `json:"users"`
	Paginator Paginator   `json:"paginator"`
}

// UserCollection is a paginated collection of User entities
type UserCollection struct {
	Users     []User    `json:"users"`
	Paginator Paginator `json:"paginator"`
}
