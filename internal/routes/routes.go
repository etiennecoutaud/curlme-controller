package routes

import "net/http"

// Healthz handle URL healthcheck
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
