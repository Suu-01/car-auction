package api

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/ksj/car-auction/internal/config"
)

type ctxKey string

const (
	userIDKey ctxKey = "user_id"
	roleKey   ctxKey = "user_role"
)

// AuthMiddleware は JWT を検証し、user_id と role をコンテキストに保存します
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1) Authorization ヘッダーの取得
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header is required", http.StatusUnauthorized)
			return
		}
		// 2) "Bearer トークン" 形式の検証
		log.Printf("[AUTH DBG] Authorization header: %q\n", authHeader)

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		// 3) JWT のパース
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.JwtSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			log.Printf("[AUTH DBG] token parse error: %v\n", err)
			return
		}

		// 4) Claims の取得
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		// 5) user_id と role の抽出
		uidFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "invalid user_id claim", http.StatusUnauthorized)
			return
		}
		userID := uint(uidFloat)

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "invalid role claim", http.StatusUnauthorized)
			return
		}
		log.Printf("[AUTH DBG] authenticated user=%d, role=%q\n", userID, role)

		// 6) コンテキストに保存して次のハンドラーへ
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, roleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext はコンテキストから user_id と role を取得します
func FromContext(r *http.Request) (userID uint, role string, ok bool) {
	uid, ok1 := r.Context().Value(userIDKey).(uint)
	rl, ok2 := r.Context().Value(roleKey).(string)
	return uid, rl, ok1 && ok2
}

// RequireRole は指定されたロールのみアクセスを許可するミドルウェアを返します
func RequireRole(wantRole string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, role, ok := FromContext(r)
			if !ok || role != wantRole {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
