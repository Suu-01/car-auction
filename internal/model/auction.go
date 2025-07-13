package model

import "time"

// Auction 모델: GORM이 이 구조체를 보고 테이블을 만듭니다
type Auction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartPrice  int       `json:"start_price"`
	OwnerID     uint      `json:"owner_id"` // ← 이 줄이 반드시 있어야 합니다
	CreatedAt   time.Time `json:"created_at"`
	EndAt       time.Time `json:"end_at"`
}
