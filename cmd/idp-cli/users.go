package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/oleksandr/idp/entities"
)

func listUsers(c *cli.Context) {
	var (
		err        error
		sorter     entities.Sorter
		pager      entities.Pager
		collection *entities.UserCollection
		paginator  entities.Paginator
	)

	sorter = entities.Sorter{"name", true}
	pager = entities.Pager{1, 100}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "ID\tName\tEnabled")

	for {
		collection, err = userInteractor.List(pager, sorter)
		assertError(err)
		for _, d := range collection.Users {
			fmt.Fprintf(w, "%v\t%v\t%v\n", d.ID, d.Name, d.Enabled)
		}
		w.Flush()
		if !paginator.HasNextPage {
			break
		}
		pager.Page++
	}
	fmt.Printf("Page %v of %v (Total records: %v)\n", collection.Paginator.Page, collection.Paginator.TotalPages(), collection.Paginator.Total)
}

func findUser(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide an ID of the user"))
	}
	u, err := userInteractor.Find(c.Args().First())
	assertError(err)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "ID\tName\tEnabled")
	fmt.Fprintf(w, "%v\t%v\t%v\n", u.ID, u.Name, u.Enabled)
	w.Flush()
}

func addUser(c *cli.Context) {
	if c.StringSlice("domain") == nil || len(c.StringSlice("domain")) == 0 {
		assertError(fmt.Errorf("You need to specify at least one domain ID using --domain option"))
	}
	if c.String("name") == "" {
		assertError(fmt.Errorf("You need to specify user name --name option"))
	}
	if c.String("password") == "" {
		assertError(fmt.Errorf("You need to specify user password --password option"))
	}
	u := entities.NewUser(c.String("name"))
	u.SetPassword(c.String("password"))
	u.Enabled = !c.Bool("disable")

	err := userInteractor.Create(*u, c.StringSlice("domain"))
	assertError(err)
	fmt.Printf("User %v created\n", u.ID)
}

func updateUser(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide an ID of the user"))
	}
	if c.Bool("enable") && c.Bool("disable") {
		assertError(fmt.Errorf("You can provide either --enable or --disable, but not both at the same time"))
	}

	u, err := userInteractor.Find(c.Args().First())
	assertError(err)

	if c.String("password") != "" {
		u.SetPassword(c.String("password"))
	}
	if c.Bool("enable") {
		u.Enabled = true
	}
	if c.Bool("disable") {
		u.Enabled = false
	}

	err = userInteractor.Update(*u)
	assertError(err)

	fmt.Printf("User %v updated\n", u.ID)
}

func removeUser(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide an ID of the user"))
	}

	u, err := userInteractor.Find(c.Args().First())
	assertError(err)

	err = userInteractor.Delete(*u)
	assertError(err)

	fmt.Printf("User %v deleted\n", u.ID)
}
