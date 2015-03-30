package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/lib/pq"  # <-- need more testing
	//_ "github.com/mattn/go-sqlite3"  # <-- temporary disabled. need more testing
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/db"
	"github.com/oleksandr/idp/usecases"
)

func main() {
	dbmap, err := db.InitDB(os.Getenv(config.EnvIDPDriver), os.Getenv(config.EnvIDPDSN))
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer dbmap.Db.Close()
	if config.SQLTraceOn() {
		dbmap.TraceOn("", log.New(os.Stderr, "[gorp] ", log.LstdFlags))
	}
	log.SetPrefix("[main] ")

	//
	// Core setup
	//
	domainInteractor := new(usecases.DomainInteractorImpl)
	domainInteractor.DBMap = dbmap
	userInteractor := new(usecases.UserInteractorImpl)
	userInteractor.DBMap = dbmap
	sessionInteractor := new(usecases.SessionInteractorImpl)
	sessionInteractor.DBMap = dbmap

	//
	// Start the servers and GC
	//
	exitCh := make(chan bool, 1)
	go startRESTfulServer(exitCh,
		domainInteractor,
		userInteractor,
		sessionInteractor)
	go startRPCServer(exitCh,
		domainInteractor,
		userInteractor,
		sessionInteractor)
	go func() {
		for {
			select {
			case <-time.Tick(time.Duration(30) * time.Minute):
				log.Println("Purging sessions...")
				sessionInteractor.Purge()
			}
		}

	}()

	// Setup signal catcher for the server's proper shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	select {
	case s := <-c:
		log.Println("Caught signal", s.String())
	case <-exitCh:
		log.Println("Caught exit from one of the servers")
	}

	log.Println("Stopping the server...")
	// Tidy up and tear down
	log.Println("Tearing down...")
	log.Fatalln("Finished - bye bye. ;-)")

}
