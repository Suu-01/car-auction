// internal/api/health.go
package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterHealthRoute r에 /healthz 핸들러 등록
func RegisterHealthRoute(r *mux.Router) {
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods(http.MethodGet)
}
