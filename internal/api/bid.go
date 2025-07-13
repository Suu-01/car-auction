package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ksj/car-auction/internal/service"
)

// RegisterBidRoutes 는 /auctions/{id}/bids 경로를 설정합니다.
func RegisterBidRoutes(r *mux.Router, svc *service.BidService) {
	// GET /auctions/{id}/bids?page=1&size=10
	r.HandleFunc("/auctions/{id}/bids", func(w http.ResponseWriter, r *http.Request) {
		// 1) Path param
		aid, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "invalid auction id", http.StatusBadRequest)
			return
		}
		// 2) Query param
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		size, _ := strconv.Atoi(r.URL.Query().Get("size"))

		// 3) Service 호출
		bids, total, err := svc.PaginatedBids(uint(aid), page, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 4) 응답 조립
		resp := PaginatedResponse{
			Data:       make([]any, len(bids)),
			Page:       page,
			Size:       size,
			TotalCount: total,
		}
		for i, b := range bids {
			resp.Data[i] = b
		}
		json.NewEncoder(w).Encode(resp)
	}).Methods(http.MethodGet)

	// POST /auctions/{id}/bids (인증 필요)
	post := r.Methods(http.MethodPost).Subrouter()
	post.Use(AuthMiddleware)
	post.HandleFunc("/auctions/{id}/bids", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		aid, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "invalid auction id", http.StatusBadRequest)
			return
		}
		userID, ok := FromContext(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var body struct {
			Amount int `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bid, err := svc.PlaceBid(uint(aid), userID, body.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(bid)
	}).Methods(http.MethodPost)
}
