package model

import "time"

type Bid struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AuctionID uint      `json:"auction_id"`
	UserID    uint      `json:"user_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
