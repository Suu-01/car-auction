package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ksj/car-auction/internal/service"
)

// RegisterBidRoutes はオークションの入札関連ルートを登録します
func RegisterBidRoutes(r *mux.Router, svc *service.BidService) {
	br := r.PathPrefix("/auctions/{id:[0-9]+}/bids").Subrouter()

	// GET /auctions/{id}/bids?page=&size=
	br.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		aid, _ := strconv.Atoi(vars["id"])
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		size, _ := strconv.Atoi(r.URL.Query().Get("size"))
		if page < 1 {
			page = 1
		}
		if size < 1 {
			size = 10
		}

		// ビジネスロジック呼び出し: ページネーションされた入札一覧取得
		bids, total, err := svc.PaginatedBids(uint(aid), page, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// レスポンス組立
		resp := PaginatedResponse{
			Data:       make([]any, len(bids)),
			Page:       page,
			Size:       size,
			TotalCount: total,
		}
		for i, v := range bids {
			resp.Data[i] = v
		}
		// ヘッダー設定および JSON エンコード
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}).Methods("GET")

	// POST /auctions/{id}/bids
	pr := br.Methods(http.MethodPost).Subrouter()
	pr.Use(AuthMiddleware, RequireRole("bidder"))
	pr.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		userID, _, ok := FromContext(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		aid, _ := strconv.Atoi(mux.Vars(r)["id"])
		var req struct {
			Amount int `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bid, err := svc.PlaceBid(uint(aid), userID, req.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(bid)
	}).Methods("POST")
}
