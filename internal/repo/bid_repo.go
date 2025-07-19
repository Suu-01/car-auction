package repo

import (
	"github.com/ksj/car-auction/internal/model"

	"gorm.io/gorm"
)

// BidRepo は Bid モデルの DB アクセスを担当するリポジトリです
type BidRepo struct {
	DB *gorm.DB
}

// NewBidRepo は BidRepo のコンストラクタです
func NewBidRepo(db *gorm.DB) *BidRepo {
	return &BidRepo{DB: db}
}

// FindByAuction は指定されたオークションIDの入札一覧をページネーション付きで取得し、総件数を返します
func (r *BidRepo) FindByAuction(auctionID uint, page, size int) ([]model.Bid, int64, error) {
	var bids []model.Bid
	var total int64

	// 総件数を取得
	if err := r.DB.
		Model(&model.Bid{}).
		Where("auction_id = ?", auctionID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション: offset と limit, 最新順
	offset := (page - 1) * size
	if err := r.DB.
		Where("auction_id = ?", auctionID).
		Order("created_at desc").
		Offset(offset).
		Limit(size).
		Find(&bids).Error; err != nil {
		return nil, 0, err
	}

	return bids, total, nil
}

// Create は新しい入札をデータベースに保存します
func (r *BidRepo) Create(b *model.Bid) error {
	return r.DB.Create(b).Error
}

// FindPaginated は指定オークションIDの入札を offset, limit で取得します
func (r *BidRepo) FindPaginated(auctionID uint, offset, limit int) ([]model.Bid, error) {
	var bids []model.Bid
	if err := r.DB.
		Where("auction_id = ?", auctionID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

// CountByAuction は指定オークションIDの入札総件数を返します
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

// ListByAuction は指定オークションIDの入札一覧をページネーション付きで取得し、総件数を返します
func (r *BidRepo) ListByAuction(auctionID uint, page, size int) ([]model.Bid, int64, error) {
	var bids []model.Bid
	offset := (page - 1) * size
	tx := r.DB.Where("auction_id = ?", auctionID).
		Offset(offset).
		Limit(size).
		Order("created_at desc").
		Find(&bids)
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	var total int64
	r.DB.Model(&model.Bid{}).Where("auction_id = ?", auctionID).Count(&total)
	return bids, total, nil
}
