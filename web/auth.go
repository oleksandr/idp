package web

import (
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
			// Extract the token from header
			s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(s) != 2 {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Token is missing")
				return
			}
			if s[0] != "Token" {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Authentication schema is not supported")
				return
			}
			p := strings.SplitN(s[1], "=", 2)
			if len(p) != 2 {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Authentication schema is not supported")
				return
			}
			authToken := strings.TrimSpace(strings.Trim(p[1], "\" "))
			if authToken == "" {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Empty token")
				return
			}

			// Look up the session by token and other attributes
			remoteAddr := helpers.RemoteAddrFromRequest(r)
			userAgent := r.UserAgent()
			session, err := interactor.Find(authToken)
			if err != nil || session.UserAgent != userAgent || session.RemoteAddr != remoteAddr {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Session not found")
				return
			}

			// Validate session/user/domain
			if !session.IsValid() || session.IsExpired() {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Session is invalid/expired")
				return
			}
			if !session.Domain.Enabled {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Domain is disabled")
				return
			}
			if !session.User.Enabled {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized", "User account is disabled")
				return
			}

			// Retain session
			interactor.Retain(*session)

			// Save session into current request context
			context.Set(r, config.CtxSessionKey, *session)

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
