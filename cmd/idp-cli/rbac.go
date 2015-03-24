package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/oleksandr/idp/entities"
)

func listPermissions(c *cli.Context) {
	var (
		err        error
		sorter     entities.Sorter
		pager      entities.Pager
		collection *entities.BasicPermissionCollection
		paginator  entities.Paginator
	)

	sorter = entities.Sorter{"name", true}
	pager = entities.Pager{1, 100}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, '\t', 0)
	fmt.Fprintln(w, "NAME\tENABLED\tEVALUATION\tDESCRIPTION")
	fmt.Fprintln(w, "---\t\t\t")

	for {
		if c.String("role") != "" {
			collection, err = rbacInteractor.ListPermissionsByRole(c.String("role"), pager, sorter)
		} else {
			collection, err = rbacInteractor.ListPermissions(pager, sorter)
		}
		assertError(err)
		for _, p := range collection.Permissions {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", p.Name, p.Enabled, p.EvaluationRule, p.Description)
		}
		w.Flush()
		if !paginator.HasNextPage {
			break
		}
		pager.Page++
	}
	fmt.Printf("Page %v of %v (Total records: %v)\n", collection.Paginator.Page, collection.Paginator.TotalPages(), collection.Paginator.Total)
}

func findPermission(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide permission's name as an argument"))
	}
	p, err := rbacInteractor.FindPermission(c.Args().First())
	assertError(err)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, '\t', 0)
	fmt.Fprintln(w, "NAME\tENABLED\tEVALUATION\tDESCRIPTION")
	fmt.Fprintln(w, "---\t\t\t")
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", p.Name, p.Enabled, p.EvaluationRule, p.Description)
	w.Flush()
}

func addPermission(c *cli.Context) {
	if c.String("name") == "" {
		assertError(fmt.Errorf("You need to specify permission name using --name option"))
	}
	p := entities.NewBasicPermission(c.String("name"), c.String("description"))
	p.Enabled = !c.Bool("disable")
	err := rbacInteractor.CreatePermission(*p)
	assertError(err)
	fmt.Printf("Permission %v created\n", p.Name)
}

func updatePermission(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide a permission's name as an argument"))
	}
	if c.Bool("enable") && c.Bool("disable") {
		assertError(fmt.Errorf("You can provide either --enable or --disable, but not both at the same time"))
	}

	var (
		err error
	)

	p, err := rbacInteractor.FindPermission(c.Args().First())
	assertError(err)
	if c.String("description") != "" {
		p.Description = c.String("description")
	}
	if c.String("rule") != "" {
		p.EvaluationRule = c.String("rule")
	}
	if c.Bool("enable") {
		p.Enabled = true
	}
	if c.Bool("disable") {
		p.Enabled = false
	}
	err = rbacInteractor.UpdatePermission(*p)
	assertError(err)

	if c.String("name") != "" {
		err = rbacInteractor.RenamePermission(c.Args().First(), c.String("name"))
		assertError(err)
	}

	fmt.Printf("Permission %v updated\n", c.Args().First())
}

func removePermission(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide permission's name as an argument"))
	}
	err := rbacInteractor.DeletePermission(c.Args().First())
	assertError(err)
	fmt.Printf("Permission %v deleted\n", c.Args().First())
}

func listRoles(c *cli.Context) {
	var (
		err        error
		sorter     entities.Sorter
		pager      entities.Pager
		collection *entities.BasicRoleCollection
		paginator  entities.Paginator
	)

	sorter = entities.Sorter{"name", true}
	pager = entities.Pager{1, 100}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, '\t', 0)
	fmt.Fprintln(w, "NAME\tENABLED\tDESCRIPTION")
	fmt.Fprintln(w, "---\t\t")

	for {
		collection, err = rbacInteractor.ListRoles(pager, sorter)
		assertError(err)
		for _, r := range collection.Roles {
			fmt.Fprintf(w, "%v\t%v\t%v\n", r.Name, r.Enabled, r.Description)
		}
		w.Flush()
		if !paginator.HasNextPage {
			break
		}
		pager.Page++
	}
	fmt.Printf("Page %v of %v (Total records: %v)\n", collection.Paginator.Page, collection.Paginator.TotalPages(), collection.Paginator.Total)
}

func findRole(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide permission's name as an argument"))
	}
	r, err := rbacInteractor.FindRole(c.Args().First())
	assertError(err)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, '\t', 0)
	fmt.Fprintln(w, "NAME\tENABLEDt\tDESCRIPTION")
	fmt.Fprintln(w, "---\t\t")
	fmt.Fprintf(w, "%v\t%v\t%v\n", r.Name, r.Enabled, r.Description)
	w.Flush()
}

func addRole(c *cli.Context) {
	if c.String("name") == "" {
		assertError(fmt.Errorf("You need to specify role name using --name option"))
	}
	r := entities.NewBasicRole(c.String("name"), c.String("description"))
	r.Enabled = !c.Bool("disable")
	err := rbacInteractor.CreateRole(*r)
	assertError(err)
	fmt.Printf("Role %v created\n", r.Name)
}

func updateRole(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide a role's name as an argument"))
	}
	if c.Bool("enable") && c.Bool("disable") {
		assertError(fmt.Errorf("You can provide either --enable or --disable, but not both at the same time"))
	}

	var (
		err error
	)
	if c.StringSlice("add") != nil && len(c.StringSlice("add")) > 0 {
		err = rbacInteractor.UpdateRoleWithPermissions(c.Args().First(), c.StringSlice("add"))
		assertError(err)
	}
	if c.StringSlice("remove") != nil && len(c.StringSlice("remove")) > 0 {
		err = rbacInteractor.RemovePermissionsFromRole(c.StringSlice("remove"), c.Args().First())
		assertError(err)
	}

	r, err := rbacInteractor.FindRole(c.Args().First())
	assertError(err)
	if c.String("description") != "" {
		r.Description = c.String("description")
	}
	if c.Bool("enable") {
		r.Enabled = true
	}
	if c.Bool("disable") {
		r.Enabled = false
	}
	err = rbacInteractor.UpdateRole(*r)
	assertError(err)

	if c.String("name") != "" {
		err = rbacInteractor.RenameRole(c.Args().First(), c.String("name"))
		assertError(err)
	}

	fmt.Printf("Role %v updated\n", c.Args().First())
}

func removeRole(c *cli.Context) {
	if !c.Args().Present() {
		assertError(fmt.Errorf("You need to provide a role's name as an argument"))
	}
	err := rbacInteractor.DeleteRole(c.Args().First())
	assertError(err)
	fmt.Printf("Role %v deleted\n", c.Args().First())
}
