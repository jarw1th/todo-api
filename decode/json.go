package decode

import (
	"encoding/json"
	"net/http"
)

func JSONError(w http.ResponseWriter, err error, status int) {
	JSONResponse(w, status, map[string]string{"error": err.Error()})
}

func JSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}
