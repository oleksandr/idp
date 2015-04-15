package web

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/entities"
	"github.com/oleksandr/idp/errs"
	"github.com/oleksandr/idp/usecases"
)

// SessionForm used for parsing incoming data
type SessionForm struct {
	Session struct {
		User struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		} `json:"user"`
		Domain struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"domain"`
	} `json:"session"`
}

// SessionResource used for responses
type SessionResource struct {
	Session entities.Session `json:"session"`
}

//
// SessionWebHandler is a collection of CRUD methods for Sessions API
//
type SessionWebHandler struct {
	log               *log.Logger
	SessionInteractor usecases.SessionInteractor
	UserInteractor    usecases.UserInteractor
	DomainInteractor  usecases.DomainInteractor
}

// NewSessionWebHandler creates new SessionWebHandler
func NewSessionWebHandler() *SessionWebHandler {
	return &SessionWebHandler{
		log: log.New(os.Stdout, "[SessionHandler] ", log.LstdFlags),
	}
}

// Create opens a new session if none exists
func (handler *SessionWebHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Parse incoming credentials
	var form SessionForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request data", err)
		return
	}

	// Prepare arguments
	user := entities.BasicUser{}
	user.Name = form.Session.User.Name
	domain := entities.BasicDomain{}
	domain.ID = form.Session.Domain.ID
	domain.Name = form.Session.Domain.Name
	userAgent := r.UserAgent()
	remoteAddr := remoteAddrFromRequest(r)

	// Create session
	session, err := handler.SessionInteractor.CreateWithPassword(domain, user, form.Session.User.Password, userAgent, remoteAddr)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(SessionResource{Session: *session})
		return
	}

	// Handle errors
	e := err.(*errs.Error)
	handler.log.Printf("%v: %v", form, err)
	respondWithError(w, errorToHTTPStatus(e), "Failed to create session", e)
}

// Check validates if current session is valid
func (handler *SessionWebHandler) Check(w http.ResponseWriter, r *http.Request) {
	if _, ok := context.Get(r, config.CtxSessionKey).(entities.Session); ok {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// Retrieve handles a read request of a current session
func (handler *SessionWebHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	if s, ok := context.Get(r, config.CtxSessionKey).(entities.Session); ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SessionResource{Session: s})
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// Delete deletes current session
func (handler *SessionWebHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if s, ok := context.Get(r, config.CtxSessionKey).(entities.Session); ok {
		err := handler.SessionInteractor.Delete(s)
		if err == nil {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		e := err.(*errs.Error)
		handler.log.Println(e.Error())
		respondWithError(w, errorToHTTPStatus(e), "Failed to delete session", e)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}
