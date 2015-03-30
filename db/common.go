package db

import (
	"database/sql"
	"fmt"

	"github.com/oleksandr/idp/entities"
	"gopkg.in/gorp.v1"
)

// InitDB connects to a database and creates and returns driver-specific mapper
func InitDB(driverName, DSN string) (*gorp.DbMap, error) {
	cpool, err := sql.Open(driverName, DSN)
	if err != nil {
		return nil, err
	}

	var dbmap *gorp.DbMap
	switch driverName {
	case "sqlite3":
		dbmap = &gorp.DbMap{
			Db:      cpool,
			Dialect: gorp.SqliteDialect{},
		}
	case "mysql":
		dbmap = &gorp.DbMap{
			Db:      cpool,
			Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"},
		}
	case "postgres":
		dbmap = &gorp.DbMap{
			Db:      cpool,
			Dialect: gorp.PostgresDialect{},
		}
	default:
		return nil, fmt.Errorf("Unsupported driver: %v", driverName)
	}

	tmap := dbmap.AddTableWithName(Domain{}, "domain")
	tmap.SetKeys(true, "domain_id")
	tmap.ColMap("object_id").SetUnique(true)
	tmap.ColMap("name").SetUnique(true)

	tmap = dbmap.AddTableWithName(User{}, "user")
	tmap.SetKeys(true, "user_id")
	tmap.ColMap("object_id").SetUnique(true)
	tmap.ColMap("name").SetUnique(true)

	tmap = dbmap.AddTableWithName(Session{}, "session")
	tmap.SetKeys(false, "session_id")

	tmap = dbmap.AddTableWithName(Role{}, "role")
	tmap.SetKeys(true, "role_id")
	tmap.ColMap("name").SetUnique(true)

	tmap = dbmap.AddTableWithName(Permission{}, "permission")
	tmap.SetKeys(true, "permission_id")
	tmap.ColMap("name").SetUnique(true)

	tmap = dbmap.AddTableWithName(DomainUser{}, "domain_user")
	tmap.SetKeys(false, "domain_id", "user_id")

	tmap = dbmap.AddTableWithName(RolePermission{}, "role_permission")
	tmap.SetKeys(false, "role_id", "permission_id")

	tmap = dbmap.AddTableWithName(UserRole{}, "user_role")
	tmap.SetKeys(false, "user_id", "role_id")

	return dbmap, nil
}

// OrderByClause takes a Sorter and constructs a "ORDER BY" clause if required
func OrderByClause(sorter entities.Sorter, alias string) string {
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

// LimitOffset takes a pager and returns a "LIMIT/OFFSET" clause if required
func LimitOffset(pager entities.Pager) string {
	clause := ""
	if pager.PerPage > 0 {
		clause = fmt.Sprintf("LIMIT %v OFFSET %v", pager.PerPage, (pager.Page-1)*(pager.PerPage))
	}
	return clause
}
