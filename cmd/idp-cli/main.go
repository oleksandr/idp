package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oleksandr/idp/usecases"
)

var (
	db               *sqlx.DB
	domainInteractor *usecases.DomainInteractorImpl
	userInteractor   *usecases.UserInteractorImpl

//rbacInteractor   *usecases.RBACInteractorImpl
)

func main() {
	app := cli.NewApp()
	app.Name = "idp-client"
	app.Usage = "Manage Identity Provider database"
	app.Version = "0.0.1"
	app.Author = "Oleksandr Lobunets"
	app.Email = "alexander.lobunets@gmail.com"

	// DB
	db := sqlx.MustConnect("sqlite3", "/Users/alex/src/github.com/oleksandr/idp/db.sqlite3")
	defer db.Close()

	// Interactors
	domainInteractor = new(usecases.DomainInteractorImpl)
	domainInteractor.DB = db
	userInteractor = new(usecases.UserInteractorImpl)
	userInteractor.DB = db
	//rbacInteractor = new(usecases.RBACInteractorImpl)
	//sessionInteractor := new(idp.SessionInteractorImpl)

	app.Commands = []cli.Command{
		{
			Name:  "domains",
			Usage: "Manage domains",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "Print existing domains",
					Action: listDomains,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user",
							Usage: "Filter domains by given user ID",
						},
					},
				},
				{
					Name:   "find",
					Usage:  "Find an existing domains by given ID or name",
					Action: findDomain,
				},
				{
					Name:   "add",
					Usage:  "Add a new domain by providing name and description as arguments",
					Action: addDomain,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "Name to assign to domain",
						},
						cli.StringFlag{
							Name:  "description",
							Usage: "Description to assign to domain",
						},
						cli.BoolFlag{
							Name:  "disable",
							Usage: "Disable domain",
						},
					},
				},
				{
					Name:   "update",
					Usage:  "Modifies an existing domain by providing name and description as arguments",
					Action: updateDomain,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "New name to assign to domain",
						},
						cli.StringFlag{
							Name:  "description",
							Usage: "New description to assign to domain",
						},
						cli.BoolFlag{
							Name:  "enable",
							Usage: "Enable domain",
						},
						cli.BoolFlag{
							Name:  "disable",
							Usage: "Disable domain",
						},
					},
				},
				{
					Name:   "remove",
					Usage:  "Remove an existing domain",
					Action: removeDomain,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "force",
							Usage: "Delete all nested users and roles (cascading)",
						},
					},
				},
			},
		},
		{
			Name:  "users",
			Usage: "Manage users",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "List existing users",
					Action: listUsers,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "domain",
							Usage: "Filter users by given domain ID",
						},
					},
				},
				{
					Name:   "find",
					Usage:  "Find an existing user by given ID or name",
					Action: findUser,
				},
				{
					Name:   "add",
					Usage:  "Add a new user",
					Action: addUser,
					Flags: []cli.Flag{
						cli.StringSliceFlag{
							Name:  "domain",
							Usage: "Array of domain IDs to assign user to",
							Value: &cli.StringSlice{},
						},
						cli.StringFlag{
							Name:  "name",
							Usage: "Unique user name (e.g. email or login)",
						},
						cli.StringFlag{
							Name:  "password",
							Usage: "Password in clear text",
						},
						cli.BoolFlag{
							Name:  "disable",
							Usage: "Disable user",
						},
					},
				},
				{
					Name:   "update",
					Usage:  "Modifies an existing user",
					Action: updateUser,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "password",
							Usage: "New password",
						},
						cli.BoolFlag{
							Name:  "enable",
							Usage: "Enable user",
						},
						cli.BoolFlag{
							Name:  "disable",
							Usage: "Disable disable",
						},
					},
				},
				{
					Name:   "remove",
					Usage:  "Remove an existing user",
					Action: removeUser,
				},
			},
		},
		{
			Name:  "sessions",
			Usage: "Manage sessions",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "List existing sessions",
					Action: listSessions,
				},
				{
					Name:   "find",
					Usage:  "Find an existing session by given ID",
					Action: findSession,
				},
				{
					Name:   "add",
					Usage:  "Add a new session",
					Action: addSession,
				},
				{
					Name:   "remove",
					Usage:  "Remove an existing session by given ID",
					Action: removeSession,
				},
			},
		},
		{
			Name:  "roles",
			Usage: "Manage roles",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "List existing roles",
					Action: listRoles,
				},
				{
					Name:   "find",
					Usage:  "Find an existing role by given ID or name",
					Action: findRole,
				},
				{
					Name:   "add",
					Usage:  "Add a new role",
					Action: addRole,
				},
				{
					Name:   "update",
					Usage:  "Modifies an existing role",
					Action: updateRole,
				},
				{
					Name:   "remove",
					Usage:  "Remove an existing role",
					Action: removeRole,
				},
			},
		},
		{
			Name:  "permissions",
			Usage: "Manage permissions",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "List existing permissions",
					Action: listPermissions,
				},
				{
					Name:   "find",
					Usage:  "Find an existing permission by given ID or name",
					Action: findPermission,
				},
				{
					Name:   "add",
					Usage:  "Add a new permission",
					Action: addPermission,
				},
				{
					Name:   "update",
					Usage:  "Modifies an existing permission",
					Action: updatePermission,
				},
				{
					Name:   "remove",
					Usage:  "Remove an existing permission",
					Action: removePermission,
				},
			},
		},
	}

	app.Run(os.Args)
}

func assertError(err error) {
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}
