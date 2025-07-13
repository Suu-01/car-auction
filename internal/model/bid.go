package model

import "time"

// Bid 는 단일 입찰 내역을 나타냅니다.
type Bid struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AuctionID uint      `json:"auction_id"`
	UserID    uint      `json:"user_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
