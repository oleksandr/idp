package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/oleksandr/idp/entities"
)

func listDomains(c *cli.Context) {
	var (
		err        error
		sorter     entities.Sorter
		pager      entities.Pager
		collection *entities.DomainCollection
		paginator  entities.Paginator
	)

	sorter = entities.Sorter{"name", true}
	pager = entities.Pager{1, 100}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "ID\tNAME\tENABLED\tUSERS\tDESCRIPTION")
	fmt.Fprintln(w, "---\t\t\t\t")

	for {
		if c.String("user") != "" {
			collection, err = domainInteractor.ListByUser(c.String("user"), pager, sorter)
		} else {
			collection, err = domainInteractor.List(pager, sorter)
		}
		assertError(err)
		for _, d := range collection.Domains {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n", d.ID, d.Name, d.Enabled, d.UsersCount, d.Description)
		}
		w.Flush()
		if !paginator.HasNextPage {
			break
		}
		pager.Page++
	}
	fmt.Printf("Page %v of %v (Total records: %v)\n", collection.Paginator.Page, collection.Paginator.TotalPages(), collection.Paginator.Total)
}

func findDomain(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide an ID of the domain"))
	}
	d, err := domainInteractor.Find(c.Args().First())
	assertError(err)
	count, err := domainInteractor.CountUsers(d.ID)
	assertError(err)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "ID\tNAME\tENABLED\tUSERS\tDESCRIPTION")
	fmt.Fprintln(w, "---\t\t\t\t")
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n", d.ID, d.Name, d.Enabled, count, d.Description)
	w.Flush()
}

func addDomain(c *cli.Context) {
	if c.String("name") == "" {
		assertError(fmt.Errorf("You need to specify domain name using --name option"))
	}
	d := entities.NewBasicDomain(c.String("name"), c.String("description"))
	d.Enabled = !c.Bool("disable")
	err := domainInteractor.Create(*d)
	assertError(err)
	fmt.Printf("Domain %v created\n", d.ID)
}

func updateDomain(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide an ID of the domain"))
	}
	if c.Bool("enable") && c.Bool("disable") {
		assertError(fmt.Errorf("You can provide either --enable or --disable, but not both at the same time"))
	}

	d, err := domainInteractor.Find(c.Args().First())
	assertError(err)

	if c.String("name") != "" {
		d.Name = c.String("name")
	}
	if c.String("description") != "" {
		d.Description = c.String("description")
	}
	if c.Bool("enable") {
		d.Enabled = true
	}
	if c.Bool("disable") {
		d.Enabled = false
	}

	err = domainInteractor.Update(*d)
	assertError(err)

	fmt.Printf("Domain %v updated\n", d.ID)
}

func removeDomain(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide an ID of the domain"))
	}

	d, err := domainInteractor.Find(c.Args().First())
	assertError(err)

	err = domainInteractor.Delete(*d)
	assertError(err)

	fmt.Printf("Domain %v deleted\n", d.ID)
}
