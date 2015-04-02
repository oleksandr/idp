package web

import (
	"net/http"
	"strings"

	"github.com/oleksandr/idp/errs"
)

// RemoteAddrFromRequest returns remote address of the requesting client
func remoteAddrFromRequest(r *http.Request) string {
	remoteAddr := r.Header.Get("X-Real-IP")
	if remoteAddr == "" {
		remoteAddr = r.Header.Get("X-Forwarded-For")
		if remoteAddr == "" {
			s := strings.SplitN(r.RemoteAddr, "]:", 2)
			if len(s) == 2 {
				return strings.TrimLeft(s[0], "[")
			}
			remoteAddr = strings.Split(r.RemoteAddr, ":")[0]
		}
	}
	return remoteAddr
}

func errorToHTTPStatus(err *errs.Error) int {
	switch err.Type {
	case errs.ErrorTypeForbidden:
		return http.StatusForbidden
	case errs.ErrorTypeConflict:
		return http.StatusBadRequest
	case errs.ErrorTypeNotFound:
		return http.StatusNotFound
	case errs.ErrorTypeOperational:
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}
