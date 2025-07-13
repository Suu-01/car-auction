package repo

import (
	"github.com/ksj/car-auction/internal/model"
	"gorm.io/gorm"
)

type AuctionRepo struct{ DB *gorm.DB }

func NewAuctionRepo(db *gorm.DB) *AuctionRepo { return &AuctionRepo{DB: db} }

func (r *AuctionRepo) FindAll() ([]model.Auction, error) {
	var auctions []model.Auction
	if err := r.DB.Find(&auctions).Error; err != nil {
		return nil, err
	}
	return auctions, nil
}

func (r *AuctionRepo) Create(a *model.Auction) error {
	return r.DB.Create(a).Error
}

// FindPaginated: offset, limit, optional title 검색어로 페이징 조회
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

// Count: 전체 레코드 수(필터 적용 전/후 구분 가능)
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

func (r *AuctionRepo) Delete(id uint) error {
	return r.DB.Delete(&model.Auction{}, id).Error
}

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

func (r *AuctionRepo) FindByID(id uint) (*model.Auction, error) {
	var a model.Auction
	if err := r.DB.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}
