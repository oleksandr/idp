package dl

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/entities"
)

// Domain DTO
type Domain struct {
	PK          int64     `db:"domain_id"`
	ID          string    `db:"object_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Enabled     bool      `db:"is_enabled"`
	CreatedOn   time.Time `db:"created_on"`
	UpdatedOn   time.Time `db:"updated_on"`
}

// SaveDomain updates or inserts a new domain
func SaveDomain(db sqlx.Ext, d Domain) (*Domain, error) {
	var (
		q   string
		r   sql.Result
		err error
	)

	f, err := FindDomain(db, d.ID)
	if err != nil && err != ErrNotFound {
		return nil, err
	}

	now := time.Now()
	if f != nil {
		f.UpdatedOn = now
		f.Name = d.Name
		f.Description = d.Description
		f.Enabled = d.Enabled
		q = `UPDATE domain SET name=?, description=?, is_enabled=?, updated_on=?
				WHERE domain_id = ?;`
		_, err = db.Exec(q, f.Name, f.Description, f.Enabled, f.UpdatedOn, f.PK)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	d.CreatedOn = now
	d.UpdatedOn = now
	q = `INSERT INTO domain (object_id, name, description, is_enabled, created_on, updated_on)
		VALUES (?, ?, ?, ?, ?, ?);`
	r, err = db.Exec(q, d.ID, d.Name, d.Description, d.Enabled, d.CreatedOn, d.UpdatedOn)
	if err != nil {
		return nil, err
	}
	d.PK, err = r.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// DeleteDomain deletes a domain from database cascading
func DeleteDomain(db sqlx.Ext, id string) error {
	var pk int64
	err := db.QueryRowx("SELECT domain_id FROM domain WHERE object_id = ?", id).Scan(&pk)
	if err == sql.ErrNoRows {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	err = ExecuteTransactionally(db.(*sqlx.DB), func(ext sqlx.Ext) error {
		r, err := ext.Exec("DELETE FROM domain_user WHERE domain_id = ?;", pk)
		if err != nil {
			return err
		}
		r, err = ext.Exec("DELETE FROM domain WHERE domain_id = ?", pk)
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
	})
	return err
}

// CountDomains returns a total count of domain records in database
func CountDomains(db sqlx.Ext) (int64, error) {
	var count int64
	err := db.QueryRowx("SELECT count(*) FROM domain;").Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// CountDomainsByUser returns a total count of domain records belonging to a given user
func CountDomainsByUser(db sqlx.Ext, id string) (int64, error) {
	var count int64
	err := db.QueryRowx(`SELECT count(*) FROM domain_user WHERE user_id IN
		(SELECT user_id FROM user WHERE object_id = ?);`, id).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// FindAllDomains returns a page of domain records
func FindAllDomains(db sqlx.Ext, pager entities.Pager, sorter entities.Sorter) ([]*Domain, error) {
	rows, err := db.Queryx(fmt.Sprintf("SELECT * FROM domain %v %v;", orderByClause(sorter), limitOffset(pager)))
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	domains := []*Domain{}
	for rows.Next() {
		var d Domain
		err = rows.StructScan(&d)
		if err != nil {
			return nil, err
		}
		domains = append(domains, &d)
	}
	return domains, nil
}

// FindDomainsByUser returns a page of domain records filtered by a given user ID
func FindDomainsByUser(db *sqlx.DB, userID string, pager entities.Pager, sorter entities.Sorter) ([]*Domain, error) {
	q := fmt.Sprintf(`SELECT * FROM domain WHERE domain_id
		IN (
			SELECT DISTINCT domain_id FROM domain_user WHERE user_id
				IN (SELECT user_id FROM user WHERE object_id = ?)
		)
		%v %v;`, orderByClause(sorter), limitOffset(pager))
	rows, err := db.Queryx(q, userID)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	domains := []*Domain{}
	for rows.Next() {
		var domain Domain
		err = rows.StructScan(&domain)
		if err != nil {
			return nil, err
		}
		domains = append(domains, &domain)
	}
	return domains, nil
}

// FindDomain finds a domain by given domain ID
func FindDomain(db sqlx.Ext, id string) (*Domain, error) {
	var d Domain
	err := db.QueryRowx("SELECT * FROM domain WHERE object_id = ? LIMIT 1", id).StructScan(&d)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &d, nil
}

// AddUserToDomain assign a given user to
func AddUserToDomain(db sqlx.Ext, u User, d Domain) error {
	q := "INSERT INTO domain_user (domain_id, user_id) VALUES (?, ?);"
	_, err := db.Exec(q, d.PK, u.PK)
	return err
}
