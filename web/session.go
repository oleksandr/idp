package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/entities"
	"github.com/oleksandr/idp/helpers"
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
	SessionInteractor usecases.SessionInteractor
	UserInteractor    usecases.UserInteractor
	DomainInteractor  usecases.DomainInteractor
}

// Create opens a new session if none exists
func (handler *SessionWebHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Parse incoming credentials
	var form SessionForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode request data", err.Error())
		return
	}

	// Resolve user and domain
	var (
		domain *entities.BasicDomain
		user   *entities.BasicUser
	)

	if form.Session.Domain.ID != "" {
		domain, err = handler.DomainInteractor.Find(form.Session.Domain.ID)
	} else {
		domain, err = handler.DomainInteractor.FindByName(form.Session.Domain.Name)
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to create session", err.Error())
		return
	} else if !domain.Enabled {
		respondWithError(w, http.StatusForbidden, "Failed to create session", "Domain is disabled")
		return
	}

	user, err = handler.UserInteractor.FindByNameInDomain(form.Session.User.Name, domain.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to create session", err.Error())
		return
	} else if !user.Enabled {
		respondWithError(w, http.StatusForbidden, "Failed to create session", "User is disabled")
		return
	}
	if !user.IsPassword(form.Session.User.Password) {
		respondWithError(w, http.StatusBadRequest, "Failed to create session", "Incorrect name/password")
		return
	}

	userAgent := r.UserAgent()
	remoteAddr := helpers.RemoteAddrFromRequest(r)

	// Lookup existing session
	session, err := handler.SessionInteractor.FindUserSpecific(user.ID, domain.ID, userAgent, remoteAddr)
	if session != nil && !session.IsExpired() {
		handler.SessionInteractor.Retain(*session)
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(SessionResource{Session: *session})
		return
	}

	// Create new session
	session = entities.NewSession(*user, *domain, userAgent, remoteAddr)
	err = handler.SessionInteractor.Create(*session)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to create session", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SessionResource{Session: *session})
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
		handler.SessionInteractor.Delete(s)
		w.WriteHeader(http.StatusAccepted)
		return
	}
	w.WriteHeader(http.StatusForbidden)
}
