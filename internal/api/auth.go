package api

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ksj/car-auction/internal/config"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	if os.Getenv("DISABLE_AUTH") == "true" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 유효한 첫 번째(유일한) 테스트 유저 ID는 1이라고 가정
			ctx := context.WithValue(r.Context(), userIDKey, uint(1))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return config.Cfg.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}
		uidFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "invalid user_id claim", http.StatusUnauthorized)
			return
		}
		userID := uint(uidFloat)

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext 는 핸들러에서 context로부터 user_id를 꺼내는 헬퍼
func FromContext(r *http.Request) (uint, bool) {
	id, ok := r.Context().Value(userIDKey).(uint)
	return id, ok
}
