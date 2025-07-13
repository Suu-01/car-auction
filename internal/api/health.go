// internal/api/health.go
package api

import (
	stdlog "log"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterHealthRoute r에 /healthz 핸들러 등록
func RegisterHealthRoute(r *mux.Router) {
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			// 보통 health 체크에서는 로그만 남깁니다.
			stdlog.Println("health write failed:", err)
		}
	}).Methods(http.MethodGet)
}
