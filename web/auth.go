package web

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/helpers"
	"github.com/oleksandr/idp/usecases"
)

// NewAuthenticationHandler create a new handler that handles token-based authentication
func NewAuthenticationHandler(interactor usecases.SessionInteractor) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			err := func() (*Session, error) {
				// Extract the token from header
				s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
				if len(s) != 2 {
					return errors.New("token is missing")
				}
				if s[0] != "Token" {
					return errors.New("authentication schema is not supported")
				}
				p := strings.SplitN(s[1], "=", 2)
				if len(p) != 2 {
					return errors.New("authentication schema is not supported")
				}
				authToken := strings.TrimSpace(strings.Trim(p[1], "\" "))
				if authToken == "" {
					return errors.New("empty token")
				}

				// Look up the session by token and other attributes
				remoteAddr := helpers.RemoteAddrFromRequest(r)
				userAgent := r.UserAgent()
				session, err := interactor.Find(authToken)
				if err != nil || session.UserAgent != userAgent || session.RemoteAddr != remoteAddr {
					return errors.New("sessio not found")
				}

				// Validate session/user/domain
				if !session.IsValid() || session.IsExpired() {
					return errors.New("invalid/expired session")
				}
				if !session.Domain.Enabled {
					return errors.New("domain is disabled")
					return
				}
				if !session.User.Enabled {
					return errors.New("user is disabled")
				}

				// Retain session
				err = interactor.Retain(*session)
				if err != nil {
					return error.New(err.Error())
				}
			}()

			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized", err.Error())
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
