package api

import (
	"net/http"
)

// CORS ミドルウェア: 全ドメイン(*)からのアクセスを許可
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 必要に応じて "*" の代わりに "http://localhost:5173" のみを許可できます。
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// ブラウザのプリフライトリクエスト(OPTIONS)には即座に応答します。
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
