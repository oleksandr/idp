package entities

import (
	"crypto/sha1"
	"fmt"

	"github.com/oleksandr/idp/config"
	"github.com/satori/go.uuid"
)

//
// User entities represent users.
//
type User struct {
	ID       string `json:"user_id"`
	Name     string `json:"name"`
	Password string `json:"-"`
	Enabled  bool   `json:"enabled"`
}

// NewUser - a constructor for `User`
func NewUser(name string) *User {
	a := new(User)
	a.ID = uuid.NewV4().String()
	a.Name = name
	a.Enabled = true
	return a
}

// SetPassword hashes a given clearTxt and assigns it to password field
func (u *User) SetPassword(clearTxt string) {
	u.Password = passwordHash(clearTxt)
}

// IsValid checks if user is valid
func (u *User) IsValid() (bool, error) {
	if u.Name == "" {
		return false, fmt.Errorf("Name cannot be empty!")
	}
	if u.Password == "" {
		return false, fmt.Errorf("Password cannot be empty!")
	}
	return true, nil
}

// passwordHash one-way hashes a string with the private HashSecret value.
func passwordHash(s string) string {
	hash := sha1.New()
	hash.Write([]byte(s))
	hash.Write([]byte(config.HashSecretSalt))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

//
// UserCollection is a paginated collection of User entities
//
type UserCollection struct {
	Users     []*User    `json:"users"`
	Paginator *Paginator `json:"paginator"`
}
