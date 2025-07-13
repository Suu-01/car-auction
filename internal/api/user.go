package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ksj/car-auction/internal/service"
)

func RegisterUserRoutes(r *mux.Router, svc *service.UserService) {
	ur := r.PathPrefix("/users").Subrouter()
	ur.HandleFunc("/signup", svc.Signup).Methods("POST")
	ur.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Email, Password string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// svc를 사용해서 로그인 로직 호출
		tok, err := svc.Login(req.Email, req.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"token": tok})
	}).Methods("POST")
}
