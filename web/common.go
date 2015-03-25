package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/entities"
)

//
// Error structure for handler's error responses
//
type Error struct {
	Code    int                    `json:"code"`
	Title   string                 `json:"title"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func respondWithError(w http.ResponseWriter, statusCode int, title, message string) {
	e := new(Error)
	e.Code = statusCode
	e.Title = title
	e.Message = message
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]*Error{"error": e})
}

// InfoHeadersHandler is a dummy middleware to add extra headers to the response
func InfoHeadersHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// NewContentTypeHandler creates a middleware that validates the request content
// type is acompatible with the provided contentTypes list.
// It writes a HTTP 415 error if that fails.
// Only PUT, POST, and PATCH requests are considered.
func NewContentTypeHandler(contentTypes ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !(r.Method == "PUT" || r.Method == "POST" || r.Method == "PATCH") {
				next.ServeHTTP(w, r)
				return
			}
			for _, ct := range contentTypes {
				if isContentType(r.Header, ct) {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, fmt.Sprintf("Unsupported content type %q; expected one of %q", r.Header.Get("Content-Type"), contentTypes), http.StatusUnsupportedMediaType)
		}
		return http.HandlerFunc(fn)
	}
}

func isContentType(h http.Header, contentType string) bool {
	ct := h.Get("Content-Type")
	if i := strings.IndexRune(ct, ';'); i != -1 {
		ct = ct[0:i]
	}
	return ct == contentType
}

// JSONRenderingHandler sets application/json as content-type.
// This helps some typing inside of the individual handlers.
func JSONRenderingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Set headers before next call otherwise the writer might be closed already
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// LoggingHandler implements simple logging middleware to log incoming requests
// (to be extended further)
func LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		// Call the next handler
		next.ServeHTTP(w, r)
		// Guess the current domain/user
		domain := ""
		user := "anonymous"
		if s, ok := context.Get(r, config.CtxSessionKey).(entities.Session); ok {
			domain = s.Domain.Name
			user = s.User.Name
		}
		// Log to standard logger
		log.Printf("[%s@%s][%s] %q %v\n", user, domain, r.Method, r.URL.String(), time.Now().Sub(t1))
	}
	return http.HandlerFunc(fn)
}

// RecoverHandler recovers after the panic in the chain, return 500 error
// and log the error.
func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				log.Printf("PANIC: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
