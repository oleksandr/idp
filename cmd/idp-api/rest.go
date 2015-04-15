package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/usecases"
	"github.com/oleksandr/idp/web"
)

func startRESTfulServer(exitCh chan bool,
	domainInteractor usecases.DomainInteractor,
	userInteractor usecases.UserInteractor,
	sessionInteractor usecases.SessionInteractor,
	rbacInteractor usecases.RBACInteractor) {

	// Web handlers
	sessionHandler := web.NewSessionWebHandler()
	sessionHandler.SessionInteractor = sessionInteractor
	sessionHandler.UserInteractor = userInteractor
	sessionHandler.DomainInteractor = domainInteractor

	rbacHandler := web.NewRBACWebHandler()
	rbacHandler.RBACInteractor = rbacInteractor

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

	// RBAC API
	router.head(versionedRoute("/assert/role/:role"), protectedChain.ThenFunc(rbacHandler.AssertRole))
	router.head(versionedRoute("/assert/permission/:permission"), protectedChain.ThenFunc(rbacHandler.AssertPermission))

	// Utilities
	router.get("/", publicChain.ThenFunc(web.IndexHandler))

	//
	// Make a HTTP Server structure using our custom handler/router
	//
	addr := os.Getenv(config.EnvIDPRESTAddr)
	if addr == "" {
		addr = ":8000"
	}
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
		log.Println("Could not listen: %s", listenErr)
		exitCh <- true
		return
	}
	defer listener.Close()

	log.Println("RESTful API Server listening", addr)
	serveErr := s.Serve(listener)
	if serveErr != nil {
		log.Println("Error in Serve:", serveErr)
		exitCh <- true
	}
}

func versionedRoute(r string) string {
	return fmt.Sprintf("/v%v/%v", config.CurrentAPIVersion, strings.TrimLeft(r, "/"))
}
