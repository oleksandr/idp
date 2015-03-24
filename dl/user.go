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

	// Transient/calculated attribute
	DomainsCount int64 `db:"domains_count"`
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
	err = ExecuteTransactionally(db.(*sqlx.DB), func(ext sqlx.Ext) error {
		r, err := ext.Exec("DELETE FROM user_session WHERE user_id = ?;", pk)
		if err != nil {
			return err
		}
		r, err = ext.Exec("DELETE FROM user_role WHERE user_id = ?;", pk)
		if err != nil {
			return err
		}
		r, err = ext.Exec("DELETE FROM domain_user WHERE user_id = ?;", pk)
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

// CountUsersByDomain returns a total count of user records belonging to a given domain
func CountUsersByDomain(db sqlx.Ext, id string) (int64, error) {
	var count int64
	err := db.QueryRowx(`SELECT count(*) FROM domain_user WHERE domain_id IN
		(SELECT domain_id FROM domain WHERE object_id = ?);`, id).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// FindAllUsers returns a page of user records
func FindAllUsers(db *sqlx.DB, pager entities.Pager, sorter entities.Sorter) ([]*User, error) {
	q := fmt.Sprintf(`SELECT u.*, count(du.user_id) AS domains_count FROM user AS u
		LEFT JOIN domain_user AS du ON u.user_id = du.user_id
		GROUP BY u.user_id
		%v %v;`, orderByClause(sorter, "u"), limitOffset(pager))
	rows, err := db.Queryx(q)
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

// FindUsersByDomain returns a page of user records filtered by a given domain ID
func FindUsersByDomain(db *sqlx.DB, domainID string, pager entities.Pager, sorter entities.Sorter) ([]*User, error) {
	q := fmt.Sprintf(`SELECT u.*, count(du.user_id) AS domains_count FROM user AS u
		LEFT JOIN domain_user AS du ON u.user_id = du.user_id
		WHERE u.user_id IN (
			SELECT DISTINCT user_id FROM domain_user WHERE domain_id
				IN (SELECT domain_id FROM domain WHERE object_id = ?)
		)
		GROUP BY u.user_id
		%v %v;`, orderByClause(sorter, ""), limitOffset(pager))
	rows, err := db.Queryx(q, domainID)
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

// FindUserInDomain finds a user by given user ID in a given domain
func FindUserInDomain(db sqlx.Ext, userID, domainID string) (*User, error) {
	var u User
	q := `SELECT user.* FROM domain_user
   		LEFT JOIN user ON domain_user.user_id=user.user_id
   		LEFT JOIN domain ON domain_user.domain_id=domain.domain_id
   		WHERE user.object_id = ?
   		AND domain.object_id = ?
   		LIMIT 1;`
	err := db.QueryRowx(q, userID, domainID).StructScan(&u)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &u, nil
}

// AddUserToDomain assign a given user to
func AddUserToDomain(db sqlx.Ext, userID, domainID string) error {
	q := `INSERT OR REPLACE INTO domain_user (user_id, domain_id)
		SELECT user.user_id, domain.domain_id FROM user, domain
		WHERE user.object_id = ? AND domain.object_id = ?;`
	_, err := db.Exec(q, userID, domainID)
	return err
}

// RemoveUserFromDomain removes a user from a domain
func RemoveUserFromDomain(db sqlx.Ext, userID, domainID string) error {
	q := `DELETE FROM domain_user
			WHERE domain_user.user_id IN (SELECT user_id FROM user WHERE user.object_id = ?)
			AND domain_user.domain_id IN (SELECT domain_id FROM domain WHERE domain.object_id = ?);`
	_, err := db.Exec(q, userID, domainID)
	return err
}

// AssignRoleToUser assign a given role to a user
func AssignRoleToUser(db sqlx.Ext, roleName, userID string) error {
	q := `INSERT OR REPLACE INTO user_role (user_id, role_id)
		SELECT user.user_id, role.role_id FROM user, role
		WHERE user.object_id = ? AND role.name = ?;`
	_, err := db.Exec(q, userID, roleName)
	return err
}

// RevokeRoleFromUser revoke role from a user
func RevokeRoleFromUser(db sqlx.Ext, roleName, userID string) error {
	q := `DELETE FROM user_role
		WHERE user_role.user_id IN (SELECT user_id FROM user WHERE user.object_id = ?)
		AND user_role.role_id IN (SELECT role_id FROM role WHERE role.name = ?);`
	_, err := db.Exec(q, userID, roleName)
	return err
}
