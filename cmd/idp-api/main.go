package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/justinas/alice"
	//_ "github.com/lib/pq"  # <-- need more testing
	//_ "github.com/mattn/go-sqlite3"  # <-- temporary disabled. need more testing
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/db"
	"github.com/oleksandr/idp/usecases"
	"github.com/oleksandr/idp/web"
)

func main() {
	// Essentials to startup
	addr := os.Getenv(config.EnvIDPAddr)
	if addr == "" {
		addr = ":8000"
	}
	dbmap, err := db.InitDB(os.Getenv(config.EnvIDPDriver), os.Getenv(config.EnvIDPDSN))
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer dbmap.Db.Close()
	if config.SQLTraceOn() {
		dbmap.TraceOn("[gorp]", log.New(os.Stderr, "", log.LstdFlags))
	}

	//
	// Core setup
	//

	// Interactors
	domainInteractor := new(usecases.DomainInteractorImpl)
	domainInteractor.DBMap = dbmap
	userInteractor := new(usecases.UserInteractorImpl)
	userInteractor.DBMap = dbmap
	sessionInteractor := new(usecases.SessionInteractorImpl)
	sessionInteractor.DBMap = dbmap

	// Web handlers
	//domainHandler := new(web.DomainWebHandler)
	//domainHandler.DomainInteractor = domainInteractor
	//userHandler := new(web.UserWebHandler)
	//userHandler.UserInteractor = userInteractor
	sessionHandler := new(web.SessionWebHandler)
	sessionHandler.SessionInteractor = sessionInteractor
	sessionHandler.UserInteractor = userInteractor
	sessionHandler.DomainInteractor = domainInteractor

	//
	// Middleware chain (mind the order!)
	//
	contentTypeHandler := web.NewContentTypeHandler("application/json")
	publicChain := alice.New(
		context.ClearHandler,     // cleanup ctx to avoid memory leakage
		web.LoggingHandler,       // basic requests logging
		web.RecoverHandler,       // transform panics into 500 responses
		contentTypeHandler,       // check content-type for modification requests
		web.InfoHeadersHandler,   // dummy handler to inject some info headers
		web.JSONRenderingHandler, // always set JSON content-type for this API
	)
	tokenAuthHandler := web.NewAuthenticationHandler(sessionInteractor)
	protectedChain := alice.New(
		context.ClearHandler,
		web.LoggingHandler,
		web.RecoverHandler,
		contentTypeHandler,
		web.InfoHeadersHandler,
		web.JSONRenderingHandler,
		tokenAuthHandler, // always check if request is authenticated
	)

	//
	// Routing setup
	//
	router := newRouter()

	// Domain API
	/*
		router.post("/domains", protectedChain.ThenFunc(domainHandler.Create))
		router.get("/domains/:id", protectedChain.ThenFunc(domainHandler.Retrieve))
		router.get("/domains", protectedChain.ThenFunc(domainHandler.List))
		router.patch("/domains/:id", protectedChain.ThenFunc(domainHandler.Modify))
		router.delete("/domains/:id", protectedChain.ThenFunc(domainHandler.Delete))
	*/
	// Users API
	/*
		router.post("/users", protectedChain.ThenFunc(userHandler.Create))
		router.get("/users/:id", protectedChain.ThenFunc(userHandler.Retrieve))
		router.get("/users", protectedChain.ThenFunc(userHandler.List))
		router.patch("/users/:id", protectedChain.ThenFunc(userHandler.Modify))
		router.delete("/users/:id", protectedChain.ThenFunc(userHandler.Delete))
	*/

	// Sessions API
	router.post(versionedRoute("/sessions"), publicChain.ThenFunc(sessionHandler.Create))
	router.head(versionedRoute("/sessions/current"), protectedChain.ThenFunc(sessionHandler.Check))
	router.get(versionedRoute("/sessions/current"), protectedChain.ThenFunc(sessionHandler.Retrieve))
	router.delete(versionedRoute("/sessions/current"), protectedChain.ThenFunc(sessionHandler.Delete))

	// Utilities
	router.get("/", publicChain.ThenFunc(web.IndexHandler))

	//
	// Make a HTTP Server structure using our custom handler/router
	//
	s := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Create binding address listener
	listener, listenErr := net.Listen("tcp", addr)
	if listenErr != nil {
		log.Fatalf("Could not listen: %s", listenErr)
	}
	log.Println("Listening", addr)

	// Setup signal catcher for the server's proper shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		for _ = range c {
			// Stop the HTTP server
			log.Println("Stopping the server...")
			listener.Close()
			// Tidy up and tear down
			log.Println("Tearing down...")
			log.Fatalln("Finished - bye bye. ;-)")
		}
	}()

	go func() {
		for {
			select {
			case <-time.Tick(time.Duration(30) * time.Minute):
				log.Println("Purging sessions...")
				sessionInteractor.Purge()
			}
		}

	}()

	log.Fatalf("Error in Serve: %s", s.Serve(listener))
}

func versionedRoute(r string) string {
	return fmt.Sprintf("/v%v/%v", config.CurrentAPIVersion, strings.TrimLeft(r, "/"))
}
