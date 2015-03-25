package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/codegangsta/cli"
	"github.com/oleksandr/idp/entities"
)

func listSessions(c *cli.Context) {
	var (
		err        error
		sorter     entities.Sorter
		pager      entities.Pager
		collection *entities.SessionCollection
		paginator  entities.Paginator
	)

	sorter = entities.Sorter{"", false}
	pager = entities.Pager{1, 100}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "ID\tDOMAIN\tUSER\tEXPIRED")
	fmt.Fprintln(w, "---\t\t\t")

	for {
		collection, err = sessionInteractor.List(pager, sorter)
		assertError(err)
		for _, s := range collection.Sessions {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", s.ID, s.Domain.Name, s.User.Name, s.IsExpired())
		}
		w.Flush()
		if !paginator.HasNextPage {
			break
		}
		pager.Page++
	}
	fmt.Printf("Page %v of %v (Total records: %v)\n", collection.Paginator.Page, collection.Paginator.TotalPages(), collection.Paginator.Total)
}

func findSession(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide an ID of the session"))
	}
	s, err := sessionInteractor.Find(c.Args().First())
	assertError(err)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "ID\tDOMAIN\tUSER\tEXPIRED")
	fmt.Fprintln(w, "---\t\t\t")
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", s.ID, s.Domain.Name, s.User.Name, s.IsExpired())
	w.Flush()
}

func addSession(c *cli.Context) {
	user, err := userInteractor.Find(c.String("user"))
	assertError(err)

	domain, err := domainInteractor.Find(c.String("domain"))
	assertError(err)

	s := entities.NewSession(*user, *domain, c.String("agent"), c.String("remote"))
	s.ExpiresOn.Time = time.Now().UTC().Add(c.Duration("ttl"))
	if s.IsExpired() {
		assertError(fmt.Errorf("Provide a valid TTL value. You cannot created an expired session"))
	}

	err = sessionInteractor.Create(*s)
	assertError(err)

	fmt.Printf("Session %v created\n", s.ID)
}

func removeSession(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide an ID of the session"))
	}
	s, err := sessionInteractor.Find(c.Args().First())
	assertError(err)
	err = sessionInteractor.Delete(*s)
	assertError(err)
	fmt.Printf("Session %v deleted\n", c.Args().First())
}
