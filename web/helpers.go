package web

import (
	"net/http"
	"strings"
)

// RemoteAddrFromRequest returns remote address of the requesting client
func RemoteAddrFromRequest(r *http.Request) string {
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
