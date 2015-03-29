package db

import (
	"time"

	"gopkg.in/gorp.v1"
)

// Domain table
type Domain struct {
	PK          int64     `db:"domain_id"`
	ID          string    `db:"object_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Enabled     bool      `db:"is_enabled"`
	CreatedOn   time.Time `db:"created_on"`
	UpdatedOn   time.Time `db:"updated_on"`
}

// DomainWithStats is a view with users count
type DomainWithStats struct {
	Domain
	UsersCount int64 `db:"users_count"`
}

// DomainUser table
type DomainUser struct {
	DomainPK int64 `db:"domain_id"`
	UserPK   int64 `db:"user_id"`
}

// DeleteDomain deletes a domain a referenced records
func DeleteDomain(dbmap *gorp.DbMap, id string) error {
	var d Domain
	err := dbmap.SelectOne(&d, "SELECT * FROM domain WHERE object_id = ?", id)
	if err != nil {
		return err
	}

	tx, err := dbmap.Begin()
	if err != nil {
		return nil
	}

	_, err = tx.Exec("DELETE FROM user_session WHERE domain_id = ?;", d.PK)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM domain_user WHERE domain_id = ?;", d.PK)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Delete(&d)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
