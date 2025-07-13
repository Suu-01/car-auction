package service

import (
	"errors"
	"time"

	"github.com/ksj/car-auction/internal/model"
	"github.com/ksj/car-auction/internal/repo"
)

// CreateAuctionRequest: API에서 받을 JSON과 1:1 매핑되는 DTO
type CreateAuctionRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartPrice  int       `json:"start_price"`
	EndAt       time.Time `json:"end_at"`
}

// UpdateAuctionRequest: 수정 가능한 필드만
type UpdateAuctionRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartPrice  *int       `json:"start_price,omitempty"`
	EndAt       *time.Time `json:"end_at,omitempty"`
}

// AuctionService 는 비즈니스 로직을 담당합니다.
type AuctionService struct {
	repo *repo.AuctionRepo
}

// NewAuctionService: repo를 주입받아 생성
func NewAuctionService(r *repo.AuctionRepo) *AuctionService {
	return &AuctionService{repo: r}
}

// ListAuctions: 전체 경매 목록 조회 (GET)
func (s *AuctionService) ListAuctions() ([]model.Auction, error) {
	return s.repo.FindAll()
}

// CreateAuction: 인증된 사용자(userID)와 요청 DTO로 새 경매 생성 (POST)
func (s *AuctionService) CreateAuction(userID uint, req CreateAuctionRequest) (*model.Auction, error) {
	a := &model.Auction{
		Title:       req.Title,
		Description: req.Description,
		StartPrice:  req.StartPrice,
		OwnerID:     userID,
		CreatedAt:   time.Now(),
		EndAt:       req.EndAt,
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

// UpdateAuction: 소유자 체크 후 업데이트
func (s *AuctionService) UpdateAuction(userID, id uint, req UpdateAuctionRequest) (*model.Auction, error) {
	// 1) 기존 경매 조회
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
	// 2) 소유자(userID)만 수정 허용
	if existing.OwnerID != userID {
		return nil, errors.New("forbidden: not owner")
	}
	// 3) 변경 요청 반영
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
	// 4) DB 업데이트
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteAuction: ID로 경매 삭제 (DELETE)
func (s *AuctionService) DeleteAuction(userID uint, auctionID uint) error {
	// 1) 경매 조회
	a, err := s.repo.FindByID(auctionID)
	if err != nil {
		return err
	}
	// 2) 소유자 검증
	if a.OwnerID != userID {
		return errors.New("forbidden")
	}
	// 3) 삭제
	return s.repo.Delete(auctionID)
}

func (s *AuctionService) GetAuction(id uint) (*model.Auction, error) {
	return s.repo.FindByID(id)
}

// PaginatedAuctions: page(1-based), size, titleFilter 로 조회
// 반환값: 경매목록, 전체개수(total), 오류
func (s *AuctionService) PaginatedAuctions(page, size int, titleFilter string) ([]model.Auction, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	// 1) 페이징된 목록 조회
	auctions, err := s.repo.FindPaginated(offset, size, titleFilter)
	if err != nil {
		return nil, 0, err
	}

	// 2) 전체 개수 조회
	total, err := s.repo.Count(titleFilter)
	if err != nil {
		return nil, 0, err
	}

	return auctions, total, nil
}
