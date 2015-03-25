package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oleksandr/idp/config"
)

//
// Showing info about API
//
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Version string `json:"version"`
		Message string `json:"message"`
	}{
		Version: fmt.Sprintf("v%v", config.CurrentAPIVersion),
		Message: "Welcome to Identities API",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
