package repo

import (
	"github.com/ksj/car-auction/internal/model"
	"gorm.io/gorm"
)

// AuctionRepo はオークションの永続化を担当するリポジトリです
type AuctionRepo struct{ DB *gorm.DB }

// NewAuctionRepo は新しい AuctionRepo を生成します
func NewAuctionRepo(db *gorm.DB) *AuctionRepo { return &AuctionRepo{DB: db} }

// FindAll は全オークションを取得します
func (r *AuctionRepo) FindAll() ([]model.Auction, error) {
	var auctions []model.Auction
	if err := r.DB.Find(&auctions).Error; err != nil {
		return nil, err
	}
	return auctions, nil
}

// Create は新しいオークションをデータベースに保存します
func (r *AuctionRepo) Create(a *model.Auction) error {
	return r.DB.Create(a).Error
}

// FindPaginated はオフセット・リミット・タイトルフィルタを使ってオークションをページング取得します
func (r *AuctionRepo) FindPaginated(offset, limit int, titleFilter string) ([]model.Auction, error) {
	var auctions []model.Auction
	q := r.DB.Model(&model.Auction{})
	if titleFilter != "" {
		q = q.Where("title LIKE ?", "%"+titleFilter+"%")
	}
	if err := q.Offset(offset).Limit(limit).Find(&auctions).Error; err != nil {
		return nil, err
	}
	return auctions, nil
}

// Count はタイトルフィルタ適用後のオークション総件数を返します
func (r *AuctionRepo) Count(titleFilter string) (int64, error) {
	var total int64
	q := r.DB.Model(&model.Auction{})
	if titleFilter != "" {
		q = q.Where("title LIKE ?", "%"+titleFilter+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// DeleteByID は指定IDのオークションを削除します
func (r *AuctionRepo) DeleteByID(id uint) error {
	return r.DB.Delete(&model.Auction{}, id).Error
}

// Update は指定されたオークションのタイトル・説明・開始価格・終了日時を更新します
func (r *AuctionRepo) Update(a *model.Auction) error {
	return r.DB.Model(&model.Auction{}).
		Where("id = ?", a.ID).
		Updates(map[string]interface{}{
			"title":       a.Title,
			"description": a.Description,
			"start_price": a.StartPrice,
			"end_at":      a.EndAt,
		}).Error
}

// FindByID は指定IDのオークションを取得します
func (r *AuctionRepo) FindByID(id uint) (*model.Auction, error) {
	var a model.Auction
	tx := r.DB.First(&a, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &a, nil
}
