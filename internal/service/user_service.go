package service

import (
	"errors"
	"time"

	"github.com/ksj/car-auction/internal/config"
	"github.com/ksj/car-auction/internal/model"
	"github.com/ksj/car-auction/internal/repo"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService はユーザー関連のビジネスロジックを提供します
type UserService struct{ Repo *repo.UserRepo }

// NewUserService はリポジトリを注入して UserService を生成します
func NewUserService(r *repo.UserRepo) *UserService { return &UserService{Repo: r} }

// CreateUserRequest はサインアップ時に受け取るリクエスト DTO です
type CreateUserRequest struct {
	Email    string
	Password string
	Role     string // "seller" or "bidder"
}

// Signup は新規ユーザーを登録し、JWT トークンを返却します
func (s *UserService) Signup(req CreateUserRequest) (string, error) {
	// バリデーション: 必須項目とロールのチェック
	if req.Email == "" || req.Password == "" || (req.Role != "bidder" && req.Role != "seller") {
		return "", errors.New("invalid signup request")
	}

	// パスワードのハッシュ化
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// ユーザーモデルの組み立て
	u := &model.User{
		Email:     req.Email,
		Password:  string(hashBytes),
		Role:      req.Role, // ← 저장
		CreatedAt: time.Now(),
	}
	if err := s.Repo.Create(u); err != nil {
		return "", err
	}

	// JWT クレームの作成 (user_id, role, 有効期限)
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"role":    u.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tokenObj.SignedString([]byte(config.Cfg.JwtSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// Login はメールアドレスとパスワードを検証し、JWT トークンとユーザーロールを返却します
func (s *UserService) Login(email, password string) (token string, role string, err error) {
	// 1) ユーザー取得
	u, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}
	// 2) パスワード検証
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return "", "", errors.New("invalid credentials")
	}
	// 3) JWT クレームの作成 (user_id, role, 有効期限)
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"role":    u.Role,
		"exp":     time.Now().Add(config.Cfg.AuctionTTL).Unix(),
	}
	// 4) トークン生成
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(config.Cfg.JwtSecret)
	if err != nil {
		return "", "", err
	}
	// 5) token と role を返却
	return signed, u.Role, nil
}
