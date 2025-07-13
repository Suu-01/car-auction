package repo

import (
	"github.com/ksj/car-auction/internal/model"

	"gorm.io/gorm"
)

// BidRepo 는 Bid 모델의 DB 접근을 담당합니다.
type BidRepo struct {
	DB *gorm.DB
}

func NewBidRepo(db *gorm.DB) *BidRepo {
	return &BidRepo{DB: db}
}

// Create 는 새 입찰을 DB에 저장합니다.
func (r *BidRepo) Create(b *model.Bid) error {
	return r.DB.Create(b).Error
}

// FindByAuction 은 특정 경매ID의 모든 입찰을 조회합니다.
func (r *BidRepo) FindByAuction(auctionID uint) ([]model.Bid, error) {
	var bids []model.Bid
	if err := r.DB.Where("auction_id = ?", auctionID).
		Order("amount desc").
		Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

func (r *BidRepo) FindPaginated(auctionID uint, offset, limit int) ([]model.Bid, error) {
	var bids []model.Bid
	if err := r.DB.
		Where("auction_id = ?", auctionID).
		Offset(offset).
		Limit(limit).
		Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

// auctionID별 총 입찰 수 반환
func (r *BidRepo) CountByAuction(auctionID uint) (int64, error) {
	var cnt int64
	if err := r.DB.
		Model(&model.Bid{}).
		Where("auction_id = ?", auctionID).
		Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}
