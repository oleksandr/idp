package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func truncateTables(c *cli.Context) {
	if !c.Bool("please") {
		assertError(fmt.Errorf("Make sure what you're doing. Say --please"))
	}
	err := dbmap.TruncateTables()
	assertError(err)
	fmt.Println("Done")
}

func dropTables(c *cli.Context) {
	if !c.Bool("please") {
		assertError(fmt.Errorf("Make sure what you're doing. Say --please"))
	}
	err := dbmap.DropTablesIfExists()
	assertError(err)
	fmt.Println("Done")
}

func createTables(c *cli.Context) {
	err := dbmap.CreateTables()
	assertError(err)
	fmt.Println("Done")
}
