package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ksj/car-auction/internal/config"
	"github.com/ksj/car-auction/internal/model"
	"github.com/ksj/car-auction/internal/repo"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{ Repo *repo.UserRepo }

func NewUserService(r *repo.UserRepo) *UserService { return &UserService{Repo: r} }

// POST /signup
func (s *UserService) Signup(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 비밀번호 해싱
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u := &model.User{
		Email:     req.Email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}
	if err := s.Repo.Create(u); err != nil {
		http.Error(w, "email already in use", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"id": u.ID, "email": u.Email})
}

// Login: 이메일/비밀번호 검증 후 JWT 토큰을 반환합니다.
func (s *UserService) Login(email, password string) (string, error) {
	// 1) 사용자 조회
	u, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	// 2) 비밀번호 검증
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}
	// 3) 클레임 생성 (user_id + 만료시간)
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(config.Cfg.AuctionTTL).Unix(),
	}
	// 4) 토큰 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 5) 서명된 문자열로 반환
	return token.SignedString(config.Cfg.JwtSecret)
}
