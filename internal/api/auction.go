package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ksj/car-auction/internal/service"
)

// 1) GET /auctions
// オークションの一覧をページネーション付きで取得するハンドラ
func listAuctionsHandler(svc *service.AuctionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page, err := strconv.Atoi(q.Get("page"))
		if err != nil || page < 1 {
			page = 1
		}
		size, err := strconv.Atoi(q.Get("size"))
		if err != nil || size < 1 {
			size = 10
		}
		title := q.Get("title")

		items, total, err := svc.PaginatedAuctions(page, size, title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := PaginatedResponse{
			Data:       make([]any, len(items)),
			Page:       page,
			Size:       size,
			TotalCount: total,
		}
		for i, v := range items {
			resp.Data[i] = v
		}
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// 2) GET /auctions/{id}
// 指定IDのオークション詳細を取得するハンドラ
func getAuctionHandler(svc *service.AuctionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid auction id", http.StatusBadRequest)
			return
		}
		a, err := svc.GetAuction(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(a)
	}
}

// 3) POST /auctions
// 新規オークションを作成するハンドラ（販売者のみ）
func createAuctionHandler(svc *service.AuctionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, ok := FromContext(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var req service.CreateAuctionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		a, err := svc.CreateAuction(userID, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(a)
	}
}

// 4) PUT /auctions/{id}
// 既存オークションを更新するハンドラ（認証ユーザーのみ）
func updateAuctionHandler(svc *service.AuctionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _, ok := FromContext(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid auction id", http.StatusBadRequest)
			return
		}
		var req service.UpdateAuctionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		a, err := svc.UpdateAuction(userID, uint(id), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(a)
	}
}

// 5) DELETE /auctions/{id}
// オークションを削除するハンドラ（認証ユーザーのみ）
func DeleteAuctionHandler(svc *service.AuctionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// コンテキストからユーザーID取得
		uidAny := r.Context().Value(userIDKey)
		if uidAny == nil {
			http.Error(w, "no user", http.StatusUnauthorized)
			return
		}
		userID := uidAny.(uint)

		vars := mux.Vars(r)
		id64, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		auctionID := uint(id64)

		if err := svc.DeleteAuction(auctionID, userID); err != nil {
			if err.Error() == "unauthorized" {
				http.Error(w, "forbidden", http.StatusForbidden)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// RegisterAuctionRoutes: オークション関連のルートを登録
func RegisterAuctionRoutes(r *mux.Router, svc *service.AuctionService) {
	ar := r.PathPrefix("/auctions").Subrouter()

	ar.HandleFunc("", listAuctionsHandler(svc)).Methods(http.MethodGet)

	ar.HandleFunc("/{id:[0-9]+}", getAuctionHandler(svc)).Methods(http.MethodGet)

	seller := ar.Methods(http.MethodPost).Subrouter()
	seller.Use(AuthMiddleware, RequireRole("seller"))
	seller.HandleFunc("", createAuctionHandler(svc)).Methods(http.MethodPost)

	put := ar.Methods(http.MethodPut).Subrouter()
	put.Use(AuthMiddleware)
	put.HandleFunc("/{id}", updateAuctionHandler(svc)).Methods(http.MethodPut)

	del := ar.Methods(http.MethodDelete).Subrouter()
	del.Use(AuthMiddleware)
	del.HandleFunc("/{id}", DeleteAuctionHandler(svc)).Methods(http.MethodDelete)
}
