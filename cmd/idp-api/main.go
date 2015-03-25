package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oleksandr/idp/usecases"
	"github.com/oleksandr/idp/web"
)

func main() {
	// TODO: move this to environment-based config (12 factors!)
	var addr = flag.String("addr", ":8080", "HTTP bind address")
	flag.Parse()

	//
	// Core setup
	//
	db := sqlx.MustConnect("sqlite3", "/Users/alex/src/github.com/oleksandr/idp/db.sqlite3")
	defer db.Close()

	// Interactors
	domainInteractor := new(usecases.DomainInteractorImpl)
	domainInteractor.DB = db
	userInteractor := new(usecases.UserInteractorImpl)
	userInteractor.DB = db
	sessionInteractor := new(usecases.SessionInteractorImpl)
	sessionInteractor.DB = db

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
	router.post("/sessions", publicChain.ThenFunc(sessionHandler.Create))
	router.head("/sessions/current", protectedChain.ThenFunc(sessionHandler.Check))
	router.get("/sessions/current", protectedChain.ThenFunc(sessionHandler.Retrieve))
	router.delete("/sessions/current", protectedChain.ThenFunc(sessionHandler.Delete))

	// Utilities
	router.get("/", publicChain.ThenFunc(web.IndexHandler))

	//
	// Make a HTTP Server structure using our custom handler/router
	//
	s := &http.Server{
		Addr:           *addr,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Create binding address listener
	listener, listenErr := net.Listen("tcp", *addr)
	if listenErr != nil {
		log.Fatalf("Could not listen: %s", listenErr)
	}
	log.Println("Listening", *addr)

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
		log.Println("Started expired sessions cleaning routine...")
		for {
			select {
			case <-time.Tick(time.Duration(30) * time.Minute):
				log.Println("TODO: implement me (clean expired sessions!")
			}
		}

	}()

	log.Fatalf("Error in Serve: %s", s.Serve(listener))
}
