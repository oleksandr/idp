package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/usecases"
)

var (
	db                *sqlx.DB
	domainInteractor  *usecases.DomainInteractorImpl
	userInteractor    *usecases.UserInteractorImpl
	rbacInteractor    *usecases.RBACInteractorImpl
	sessionInteractor *usecases.SessionInteractorImpl
)

func main() {
	app := cli.NewApp()
	app.Name = "idp-client"
	app.Usage = "Manage Identity Provider database"
	app.Version = config.CurrentCLIVersion
	app.Author = "Oleksandr Lobunets"
	app.Email = "alexander.lobunets@gmail.com"

	// DB
	db := sqlx.MustConnect(os.Getenv(config.EnvIDPDriver), os.Getenv(config.EnvIDPDSN))
	defer db.Close()

	// Interactors
	domainInteractor = new(usecases.DomainInteractorImpl)
	domainInteractor.DB = db
	userInteractor = new(usecases.UserInteractorImpl)
	userInteractor.DB = db
	sessionInteractor = new(usecases.SessionInteractorImpl)
	sessionInteractor.DB = db
	rbacInteractor = new(usecases.RBACInteractorImpl)
	rbacInteractor.DB = db

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
							Usage: "Domain ID to assign user to",
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
						cli.StringSliceFlag{
							Name:  "add-domain",
							Usage: "Domain ID to assign user to",
							Value: &cli.StringSlice{},
						},
						cli.StringSliceFlag{
							Name:  "remove-domain",
							Usage: "Remove user from domain by given ID",
							Value: &cli.StringSlice{},
						},
						cli.StringSliceFlag{
							Name:  "assign-role",
							Usage: "Role to assign",
							Value: &cli.StringSlice{},
						},
						cli.StringSliceFlag{
							Name:  "revoke-role",
							Usage: "Role to remove",
							Value: &cli.StringSlice{},
						},
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
					Name:   "create",
					Usage:  "Create a new session",
					Action: addSession,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user",
							Usage: "User ID for a session",
						},
						cli.StringFlag{
							Name:  "domain",
							Usage: "Domain ID for a session",
						},
						cli.StringFlag{
							Name:  "agent",
							Usage: "User-Agent for a session",
						},
						cli.StringFlag{
							Name:  "remote",
							Usage: "Remote address for a session",
						},
						cli.DurationFlag{
							Name:  "ttl",
							Usage: "TTL of the session (e.g. 30s, 10m, 1h)",
						},
					},
				},
				{
					Name:   "remove",
					Usage:  "Remove an existing session by given ID from storage",
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
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "user",
							Usage: "Filter role by given user ID",
						},
					},
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
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "Role name (e.g. admin, manager)",
						},
						cli.StringFlag{
							Name:  "description",
							Usage: "Description of a role (e.g. 'Manages all users in the system')",
						},
						cli.BoolFlag{
							Name:  "disable",
							Usage: "Disable role",
						},
					},
				},
				{
					Name:   "update",
					Usage:  "Modifies an existing role",
					Action: updateRole,
					Flags: []cli.Flag{
						cli.StringSliceFlag{
							Name:  "add",
							Usage: "Array of permissions to add",
							Value: &cli.StringSlice{},
						},
						cli.StringSliceFlag{
							Name:  "remove",
							Usage: "Array of permissions to remove",
							Value: &cli.StringSlice{},
						},
						cli.StringFlag{
							Name:  "name",
							Usage: "New role's name",
						},
						cli.StringFlag{
							Name:  "description",
							Usage: "New role's description",
						},
						cli.BoolFlag{
							Name:  "enable",
							Usage: "Enable role",
						},
						cli.BoolFlag{
							Name:  "disable",
							Usage: "Disable role",
						},
					},
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
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "role",
							Usage: "Filter permissions by given role",
						},
					},
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
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "Permission name (e.g. CreatePost, DeleteCategory)",
						},
						cli.StringFlag{
							Name:  "description",
							Usage: "Description of a permission (e.g. 'Allow to create new posts')",
						},
						cli.BoolFlag{
							Name:  "disable",
							Usage: "Disable permission",
						},
					},
				},
				{
					Name:   "update",
					Usage:  "Modifies an existing permission",
					Action: updatePermission,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "New permission's name",
						},
						cli.StringFlag{
							Name:  "description",
							Usage: "New permission's description",
						},
						cli.StringFlag{
							Name:  "rule",
							Usage: "New permission's evaluation rule",
						},
						cli.BoolFlag{
							Name:  "enable",
							Usage: "Enable permission",
						},
						cli.BoolFlag{
							Name:  "disable",
							Usage: "Disable permission",
						},
					},
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
