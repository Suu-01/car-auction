package service

import (
	"time"

	"github.com/ksj/car-auction/internal/model"
	"github.com/ksj/car-auction/internal/repo"
)

type BidService struct {
	Repo *repo.BidRepo
}

func NewBidService(r *repo.BidRepo) *BidService {
	return &BidService{Repo: r}
}

// PlaceBid : auctionID, userID, amount 를 받고 DB에 저장
func (s *BidService) PlaceBid(auctionID, userID uint, amount int) (*model.Bid, error) {
	// 1) 경매 시작가 조회
	var auction model.Auction
	if err := s.Repo.DB.First(&auction, auctionID).Error; err != nil {
		return nil, err
	}

	// 2) 상대 입찰가 계산
	rel := amount - int(auction.StartPrice)

	bid := &model.Bid{
		AuctionID: auctionID,
		UserID:    userID,
		Amount:    rel,
		CreatedAt: time.Now(),
	}
	// 3) 저장
	if err := s.Repo.Create(bid); err != nil {
		return nil, err
	}
	return bid, nil
}

// GET /auctions/{id}/bids
func (s *BidService) GetBids(auctionID uint) ([]model.Bid, error) {
	return s.Repo.FindByAuction(auctionID)
}

func (s *BidService) PaginatedBids(auctionID uint, page, size int) ([]model.Bid, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	// 1) 페이징된 입찰 목록
	bids, err := s.Repo.FindPaginated(auctionID, offset, size)
	if err != nil {
		return nil, 0, err
	}
	// 2) 전체 입찰 수
	total, err := s.Repo.CountByAuction(auctionID)
	if err != nil {
		return nil, 0, err
	}
	return bids, total, nil
}
