package db

import (
	"time"

	"gopkg.in/gorp.v1"
)

// User table
type User struct {
	PK        int64     `db:"user_id"`
	ID        string    `db:"object_id"`
	Name      string    `db:"name"`
	Password  string    `db:"passwd"`
	Enabled   bool      `db:"is_enabled"`
	CreatedOn time.Time `db:"created_on"`
	UpdatedOn time.Time `db:"updated_on"`
}

// UserWithStats is a view with domains count
type UserWithStats struct {
	User
	DomainsCount int64 `db:"domains_count"`
}

// UserRole table
type UserRole struct {
	UserPK int64 `db:"user_id"`
	RolePK int64 `db:"role_id"`
}

// DeleteUser deletes a domain a referenced records
func DeleteUser(dbmap *gorp.DbMap, id string) error {
	var u User
	err := dbmap.SelectOne(&u, "SELECT * FROM user WHERE object_id = ?", id)
	if err != nil {
		return err
	}

	tx, err := dbmap.Begin()
	if err != nil {
		return nil
	}

	_, err = tx.Exec("DELETE FROM user_session WHERE user_id = ?;", u.PK)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM user_role WHERE user_id = ?;", u.PK)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM domain_user WHERE user_id = ?;", u.PK)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Delete(&u)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
