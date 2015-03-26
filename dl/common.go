package dl

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"runtime"

	"github.com/jmoiron/sqlx"
	"github.com/oleksandr/idp/entities"
)

var (
	//ErrNotFound wraps sql.ErrNoRows error for data layer
	ErrNotFound = errors.New("DataLayer: not found")
)

// ExecuteTransactionally is a shortcut to run a given function in the transaction scope
func ExecuteTransactionally(db *sqlx.DB, wrappedFunc func(ext sqlx.Ext) error) error {
	defer func() error {
		var errAsError error
		const size = 4096

		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				errAsError = err
			} else if err, ok := e.(string); ok {
				errAsError = fmt.Errorf(err)
			}

			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]

			log.Printf("Got panic: %s", string(buf))

			return errAsError
		}
		return nil
	}()

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	err = wrappedFunc(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// a cross-db way of escaping table name and field names
func escapeLiteral(db sqlx.Ext, tableName string) string {
	if db.DriverName() == "postgres" {
		return fmt.Sprintf(`"%v"`, tableName)
	}
	return fmt.Sprintf("`%v`", tableName)
}

// a cross-db way of getting the last inserted id
func lastInsertID(db sqlx.Ext, result sql.Result, tableName, pkName string) (int64, error) {
	if db.DriverName() == "postgres" {
		r, err := db.QueryRowx(fmt.Sprintf("SELECT currval(pg_get_serial_sequence('%s', '%s'))", tableName, pkName)).SliceScan()
		if err != nil {
			return -1, nil
		}
		return r[0].(int64), nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

// Takes a Sorter and constructs a "ORDER BY" clause if required
func orderByClause(sorter entities.Sorter, alias string) string {
	clause := ""
	prefix := ""
	if alias != "" {
		prefix = fmt.Sprintf("%v.", alias)
	}
	if sorter.Field != "" {
		if sorter.Asc {
			clause = fmt.Sprintf("ORDER BY %v%v ASC", prefix, sorter.Field)
		} else {
			clause = fmt.Sprintf("ORDER BY %v%v ASC", prefix, sorter.Field)
		}
	}
	return clause
}

// Takes a pager and returns a "LIMIT/OFFSET" clause if required
func limitOffset(pager entities.Pager) string {
	clause := ""
	if pager.PerPage > 0 {
		clause = fmt.Sprintf("LIMIT %v OFFSET %v", pager.PerPage, (pager.Page-1)*(pager.PerPage))
	}
	return clause
}
