package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/gorilla/mux"
	"github.com/ksj/car-auction/internal/api"
	"github.com/ksj/car-auction/internal/config"
	"github.com/ksj/car-auction/internal/model"
	"github.com/ksj/car-auction/internal/repo"
	"github.com/ksj/car-auction/internal/service"
	"github.com/ksj/car-auction/internal/ws"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// mustOpenInMemoryDB는 순수 Go SQLite 드라이버로 메모리 DB를 열고
// AutoMigrate까지 완료한 *gorm.DB를 반환합니다. 실패 시 t.Fatal로 종료.
func mustOpenInMemoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("메모리 DB 열기 실패: %v", err)
	}
	// 모델 순서: User → Auction → Bid
	if err := db.AutoMigrate(&model.User{}, &model.Auction{}, &model.Bid{}); err != nil {
		t.Fatalf("AutoMigrate 실패: %v", err)
	}
	return db
}

func setupRouter(t *testing.T) *mux.Router {
	// 테스트용 env 세팅
	os.Setenv("DISABLE_AUTH", "true")
	// 2) config.Load() 호출 (JWT_SECRET, AUCTION_TTL 등 세팅)
	os.Setenv("DATABASE_DSN", "file::memory:?cache=shared")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("AUCTION_TTL_MINUTES", "60")
	config.Load()

	// 1) in-memory DB + AutoMigrate
	db := mustOpenInMemoryDB(t)

	hub := ws.NewHub()

	// 2) 레포 + 서비스
	auctionRepo := repo.NewAuctionRepo(db)
	bidRepo := repo.NewBidRepo(db)
	userRepo := repo.NewUserRepo(db)

	asvc := service.NewAuctionService(auctionRepo)
	bsvc := service.NewBidService(bidRepo, hub)
	usvc := service.NewUserService(userRepo)

	// 3) 라우터
	r := mux.NewRouter()
	api.RegisterUserRoutes(r, usvc)
	api.RegisterAuctionRoutes(r, asvc)
	api.RegisterBidRoutes(r, bsvc)
	return r
}

func TestBidPagination(t *testing.T) {
	router := setupRouter(t)
	server := httptest.NewServer(router)
	defer server.Close()

	// 1) 회원가입
	signupBody := map[string]string{"email": "a@b.com", "password": "pw"}
	b, _ := json.Marshal(signupBody)
	resp, _ := http.Post(server.URL+"/users/signup", "application/json", bytes.NewReader(b))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// 2) 로그인 → 토큰 획득
	resp, _ = http.Post(server.URL+"/users/login", "application/json", bytes.NewReader(b))
	var loginRes struct{ Token string }
	if err := json.NewDecoder(resp.Body).Decode(&loginRes); err != nil {
		t.Fatalf("로그인 응답 파싱 실패: %v", err)
	}
	token := loginRes.Token

	// 3) 경매 생성
	reqA := map[string]any{
		"title": "Test", "description": "D", "start_price": 100,
		"end_at": time.Now().Add(time.Hour),
	}
	bA, _ := json.Marshal(reqA)
	rA, _ := http.NewRequest("POST", server.URL+"/auctions", bytes.NewReader(bA))
	rA.Header.Set("Authorization", "Bearer "+token)
	rA.Header.Set("Content-Type", "application/json")
	resp, _ = http.DefaultClient.Do(rA)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var aRes model.Auction
	if err := json.NewDecoder(resp.Body).Decode(&aRes); err != nil {
		t.Fatalf("로그인 응답 파싱 실패: %v", err)
	}

	// 4) 여러 번 입찰
	for i := 1; i <= 12; i++ {
		bidReq := map[string]int{"amount": 100 + i}
		bb, _ := json.Marshal(bidReq)
		rb, _ := http.NewRequest("POST",
			server.URL+"/auctions/"+strconv.Itoa(int(aRes.ID))+"/bids",
			bytes.NewReader(bb))
		rb.Header.Set("Authorization", "Bearer "+token)
		rb.Header.Set("Content-Type", "application/json")
		resp, _ = http.DefaultClient.Do(rb)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	// 5) 페이지네이션 테스트: size=5, page=2 → items[5:10]
	url := server.URL + "/auctions/" + strconv.Itoa(int(aRes.ID)) +
		"/bids?page=2&size=5"
	resp, _ = http.Get(url)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var pr struct {
		Data       []model.Bid `json:"data"`
		Page, Size int
		TotalCount int64 `json:"total_count"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&pr)
	assert.Equal(t, 2, pr.Page)
	assert.Equal(t, 5, pr.Size)
	assert.Equal(t, int64(12), pr.TotalCount)
	assert.Len(t, pr.Data, 5)
	// 첫 페이지(1): bids[0..4], 두 번째 페이지(2): bids[5..9]
	assert.Equal(t, 6, pr.Data[0].Amount)
}
