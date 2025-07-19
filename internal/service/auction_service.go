package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/ksj/car-auction/internal/model"
	"github.com/ksj/car-auction/internal/repo"
)

// CreateAuctionRequest は API から受け取る JSON と 1:1 でマッピングされる DTO です
type CreateAuctionRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartPrice  int       `json:"start_price"`
	EndAt       time.Time `json:"end_at"`

	Maker     string `json:"maker"`
	ModelName string `json:"model_name"`
	Mileage   int    `json:"mileage"`
	Year      int    `json:"year"`
	PhotoURL  string `json:"photo_url"`
}

// UpdateAuctionRequest は更新可能なフィールドのみを保持する DTO です
type UpdateAuctionRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartPrice  *int       `json:"start_price,omitempty"`
	EndAt       *time.Time `json:"end_at,omitempty"`
}

// AuctionService はオークションのビジネスロジックを担当します
type AuctionService struct {
	repo *repo.AuctionRepo
}

// NewAuctionService はリポジトリを注入して AuctionService を生成します
func NewAuctionService(r *repo.AuctionRepo) *AuctionService {
	return &AuctionService{repo: r}
}

// ListAuctions は全オークションを取得します (GET)
func (s *AuctionService) ListAuctions() ([]model.Auction, error) {
	return s.repo.FindAll()
}

// CreateAuction は認証済みユーザー (sellerID) とリクエスト DTO を使って新規オークションを作成します (POST)
func (s *AuctionService) CreateAuction(sellerID uint, req CreateAuctionRequest) (*model.Auction, error) {
	if req.Title == "" || req.StartPrice <= 0 || req.Maker == "" || req.ModelName == "" {
		return nil, errors.New("invalid request")
	}
	a := &model.Auction{
		Title:       req.Title,
		Description: req.Description,
		StartPrice:  req.StartPrice,
		Maker:       req.Maker,
		ModelName:   req.ModelName,
		Mileage:     req.Mileage,
		Year:        req.Year,
		PhotoURL:    req.PhotoURL,
		SellerID:    sellerID,
		CreatedAt:   time.Now(),
		EndAt:       req.EndAt,
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

// UpdateAuction は所有者チェック後にオークション情報を更新します
func (s *AuctionService) UpdateAuction(userID, id uint, req UpdateAuctionRequest) (*model.Auction, error) {
	// 1) 既存オークション取得
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var existing *model.Auction
	for i := range list {
		if list[i].ID == id {
			existing = &list[i]
			break
		}
	}
	if existing == nil {
		return nil, errors.New("auction not found")
	}
	// 2) 所有者チェック
	if existing.SellerID != userID {
		return nil, errors.New("forbidden: not owner")
	}
	// 3) フィールド更新
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.StartPrice != nil {
		existing.StartPrice = *req.StartPrice
	}
	if req.EndAt != nil {
		existing.EndAt = *req.EndAt
	}
	// 4) 永続化更新
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteAuction はオークションを削除します（所有者のみ実行可能）
func (s *AuctionService) DeleteAuction(auctionID, userID uint) error {
	// 1) オークション情報取得
	var auc model.Auction
	if err := s.repo.DB.First(&auc, auctionID).Error; err != nil {
		return err
	}
	// 2) 所有者チェック
	if auc.SellerID != userID {
		return errors.New("unauthorized")
	}
	// 3) 削除実行
	if err := s.repo.DeleteByID(auctionID); err != nil {
		return err
	}
	return nil
}

// GetAuction は指定された ID のオークションを取得します。
// 見つからない場合はエラーを返します。
func (s *AuctionService) GetAuction(id uint) (*model.Auction, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("auction %d not found: %w", id, err)
	}
	return a, nil
}

// PaginatedAuctions は page (1 ベース)、size、titleFilter を使って
// ページングされたオークション一覧と総件数を返します。
// 戻り値: オークション一覧 ([]model.Auction)、総件数 (int64)、エラー (error)
func (s *AuctionService) PaginatedAuctions(page, size int, titleFilter string) ([]model.Auction, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	// 1) ページングされた一覧を取得
	auctions, err := s.repo.FindPaginated(offset, size, titleFilter)
	if err != nil {
		return nil, 0, err
	}

	// 2) 総件数を取得
	total, err := s.repo.Count(titleFilter)
	if err != nil {
		return nil, 0, err
	}

	return auctions, total, nil
}
