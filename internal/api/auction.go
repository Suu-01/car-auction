package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ksj/car-auction/internal/service"
)

// @Summary      페이징된 경매 목록 조회
// @Description  page, size, title(검색) 파라미터로 경매 목록을 조회합니다.
// @Tags         Auction
// @Accept       json
// @Produce      json
// @Param        page     query      int     false  "페이지 번호"      default(1)
// @Param        size     query      int     false  "페이지 크기"      default(10)
// @Param        title    query      string  false  "검색 키워드"
// @Success      200      {object}   api.PaginatedResponse
// @Failure      500      {object}   api.ErrorResponse
// @Router       /auctions [get]
func listAuctionsHandler(svc *service.AuctionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		size, _ := strconv.Atoi(r.URL.Query().Get("size"))
		title := r.URL.Query().Get("title")

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
func createAuctionHandler(svc *service.AuctionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := FromContext(r)
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
func updateAuctionHandler(svc *service.AuctionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := FromContext(r)
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
func deleteAuctionHandler(svc *service.AuctionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := FromContext(r)
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
		if err := svc.DeleteAuction(userID, uint(id)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func RegisterAuctionRoutes(r *mux.Router, svc *service.AuctionService) {
	// 1) /auctions 아래로 묶기
	ar := r.PathPrefix("/auctions").Subrouter()

	// public
	ar.HandleFunc("", listAuctionsHandler(svc)).Methods(http.MethodGet)
	ar.HandleFunc("/{id}", getAuctionHandler(svc)).Methods(http.MethodGet)

	// protected
	post := ar.Methods(http.MethodPost).Subrouter()
	post.Use(AuthMiddleware)
	post.HandleFunc("", createAuctionHandler(svc)).Methods(http.MethodPost)

	put := ar.Methods(http.MethodPut).Subrouter()
	put.Use(AuthMiddleware)
	put.HandleFunc("/{id}", updateAuctionHandler(svc)).Methods(http.MethodPut)

	del := ar.Methods(http.MethodDelete).Subrouter()
	del.Use(AuthMiddleware)
	del.HandleFunc("/{id}", deleteAuctionHandler(svc)).Methods(http.MethodDelete)
}
