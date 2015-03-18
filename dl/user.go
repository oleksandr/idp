package dl

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/entities"
)

// User DTO
type User struct {
	PK        int64     `db:"user_id"`
	ID        string    `db:"object_id"`
	Name      string    `db:"name"`
	Password  string    `db:"passwd"`
	Enabled   bool      `db:"is_enabled"`
	CreatedOn time.Time `db:"created_on"`
	UpdatedOn time.Time `db:"updated_on"`
}

// SaveUser update or inserts a new user
func SaveUser(db sqlx.Ext, u User) (*User, error) {
	var (
		q   string
		r   sql.Result
		err error
	)

	f, err := FindUser(db, u.ID)
	if err != nil && err != ErrNotFound {
		return nil, err
	}

	now := time.Now()
	if f != nil {
		f.UpdatedOn = now
		f.Name = u.Name
		f.Password = u.Password
		f.Enabled = u.Enabled
		if f.Password != "" {
			q = `UPDATE user SET name=?, passwd=?, is_enabled=?, updated_on=?
					WHERE user_id = ?;`
			_, err = db.Exec(q, f.Name, f.Password, f.Enabled, f.UpdatedOn, f.PK)
		} else {
			q = `UPDATE user SET name=?, is_enabled=?, updated_on=?
					WHERE user_id = ?;`
			_, err = db.Exec(q, f.Name, f.Enabled, f.UpdatedOn, f.PK)
		}
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	u.CreatedOn = now
	u.UpdatedOn = now
	q = `INSERT INTO user (object_id, name, passwd, is_enabled, created_on, updated_on)
		VALUES (?, ?, ?, ?, ?, ?);`
	r, err = db.Exec(q, u.ID, u.Name, u.Password, u.Enabled, u.CreatedOn, u.UpdatedOn)
	if err != nil {
		return nil, err
	}
	u.PK, err = r.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// DeleteUser deletes a user from database cascading
func DeleteUser(db sqlx.Ext, id string) error {
	var pk int64
	err := db.QueryRowx("SELECT user_id FROM user WHERE object_id = ?", id).Scan(&pk)
	if err == sql.ErrNoRows {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	err = ExecuteTransactionally(db, func(ext sqlx.Ext) error {
		r, err := ext.Exec("DELETE FROM domain_user WHERE user_id = ?;", pk)
		if err != nil {
			return err
		}
		r, err = ext.Exec("DELETE FROM user WHERE user_id = ?", pk)
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

// CountUsers returns a total count of user records in database
func CountUsers(db sqlx.Ext) (int64, error) {
	var count int64
	err := db.QueryRowx("SELECT count(*) FROM user;").Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// FindAllUsers returns a page of user records
func FindAllUsers(db sqlx.Ext, pager entities.Pager, sorter entities.Sorter) ([]*User, error) {
	rows, err := db.Queryx(fmt.Sprintf("SELECT * FROM user %v %v;", orderByClause(sorter), limitOffset(pager)))
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var user User
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// FindUser finds a user by given user ID
func FindUser(db sqlx.Ext, id string) (*User, error) {
	var u User
	err := db.QueryRowx("SELECT * FROM user WHERE object_id = ? LIMIT 1", id).StructScan(&u)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &u, nil
}