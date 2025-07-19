package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ksj/car-auction/internal/service"
)

// signupHandler は新規ユーザー登録(signup)用のハンドラです。
func signupHandler(svc *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// サービス呼び出し: ユーザー作成とトークン生成
		token, err := svc.Signup(service.CreateUserRequest{
			Email:    req.Email,
			Password: req.Password,
			Role:     req.Role,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// レスポンス: トークンを JSON で返却
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}

// RegisterUserRoutes はユーザー関連のルートを登録します
func RegisterUserRoutes(r *mux.Router, svc *service.UserService) {
	ur := r.PathPrefix("/users").Subrouter()
	ur.HandleFunc("/signup", signupHandler(svc)).Methods(http.MethodPost)
	ur.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Email, Password string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// サービス呼び出し: ログイン -> token, role, err
		tok, role, err := svc.Login(req.Email, req.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		// レスポンス: token と role を JSON で返却
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": tok,
			"role":  role,
		})
	}).Methods("POST")
}
