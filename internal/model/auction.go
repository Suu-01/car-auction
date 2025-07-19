package model

import "time"

// Auction モデル: GORM がこの構造体を見てテーブルを生成します
type Auction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartPrice  int       `json:"start_price"`
	CreatedAt   time.Time `json:"created_at"`
	EndAt       time.Time `json:"end_at"`

	SellerID uint  `gorm:"not null" json:"seller_id"`
	Seller   *User `gorm:"foreignKey:SellerID"`
	Bids     []Bid `gorm:"constraint:OnDelete:CASCADE;"`

	Maker     string `json:"maker"`
	ModelName string `json:"model_name"`
	Mileage   int    `json:"mileage"`
	Year      int    `json:"year"`
	PhotoURL  string `json:"photo_url"`
}
