package handlers

import "net/http"

// NotImplemented is a placeholder handler for routes not yet implemented
func NotImplemented(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte(`{"error":{"code":"not_implemented","message":"This endpoint is not yet implemented"}}`))
}
