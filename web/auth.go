package web

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/entities"
	"github.com/oleksandr/idp/usecases"
)

// NewAuthenticationHandler create a new handler that handles token-based authentication
func NewAuthenticationHandler(interactor usecases.SessionInteractor) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				token   string
				err     error
				session *entities.Session
			)

			// Extract token from header(s)
			token, err = tokenFromXHeader(r)
			if err != nil {
				token, err = tokenFromAuthorizationHeader(r)
				if err != nil {
					respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
					return
				}
			}

			// Look up the session by token and other attributes
			err = func() error {
				remoteAddr := remoteAddrFromRequest(r)
				userAgent := r.UserAgent()
				session, err = interactor.Find(token)
				if err != nil {
					return err
				}
				if session.UserAgent != userAgent || session.RemoteAddr != remoteAddr {
					return errors.New("Session not found for client")
				}
				return nil
			}()
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
				return
			}

			// Retain session
			err = interactor.Retain(*session)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Could not retain session", err)
				return
			}

			// Save session into current request context
			context.Set(r, config.CtxSessionKey, *session)

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func tokenFromXHeader(r *http.Request) (string, error) {
	token := strings.TrimSpace(r.Header.Get("X-Auth-Token"))
	if token != "" {
		return "", errors.New("Empty token")
	}
	return token, nil
}

func tokenFromAuthorizationHeader(r *http.Request) (string, error) {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return "", errors.New("Token is missing")
	}
	if s[0] != "Token" {
		return "", errors.New("Authentication schema is not supported")
	}
	p := strings.SplitN(s[1], "=", 2)
	if len(p) != 2 {
		return "", errors.New("Authentication schema is not supported")
	}
	token := strings.TrimSpace(strings.Trim(p[1], "\" "))
	if token == "" {
		return "", errors.New("Empty token")
	}
	return token, nil
}
