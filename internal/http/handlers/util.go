package handlers

import "net/http"

// NotImplemented is a placeholder handler for routes not yet implemented
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error":{"code":"not_implemented","message":"This endpoint is not yet implemented"}}`))
}
