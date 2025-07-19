package service

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/ksj/car-auction/internal/model"
	"github.com/ksj/car-auction/internal/repo"
	"github.com/ksj/car-auction/internal/ws"
	"gorm.io/gorm/clause"
)

// BidService は入札に関するビジネスロジックを提供します
type BidService struct {
	Repo *repo.BidRepo
	hub  *ws.Hub
}

// NewBidService は BidRepo と WebSocket Hub を注入して生成します
func NewBidService(r *repo.BidRepo, hub *ws.Hub) *BidService {
	return &BidService{Repo: r, hub: hub}
}

// PlaceBid はオークションID、ユーザーID、入札額を受け取り、入札処理を行います
func (s *BidService) PlaceBid(auctionID, userID uint, amount int) (*model.Bid, error) {
	// 1) トランザクション開始
	tx := s.Repo.DB.Begin()

	// 2) オークションレコードを FOR UPDATE でロック
	var auc model.Auction
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&auc, auctionID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	now := time.Now()
	// 3) 既に終了しているか確認
	if now.After(auc.EndAt) {
		tx.Rollback()
		return nil, errors.New("auction already closed")
	}

	// 4) 入札額が開始価格以下か確認
	if amount <= int(auc.StartPrice) {
		tx.Rollback()
		return nil, errors.New("bid too low")
	}

	// 5) 終了5分以内なら終了時間を5分延長
	if auc.EndAt.Sub(now) <= 5*time.Minute {
		newEnd := now.Add(5 * time.Minute)
		if err := tx.Model(&auc).
			Update("end_at", newEnd).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		auc.EndAt = newEnd
	}

	// 6) 入札を作成
	bid := &model.Bid{
		AuctionID: auctionID,
		UserID:    userID,
		Amount:    amount,
		CreatedAt: now,
	}
	if err := tx.Create(bid).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 7) トランザクションコミット
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// WebSocket で入札情報をブロードキャスト
	ev := map[string]interface{}{
		"bid":        bid,
		"new_end_at": auc.EndAt,
	}
	data, _ := json.Marshal(ev)
	log.Printf("WS: about to broadcast bid %d on auction %d; clients=%d",
		bid.ID, auctionID, s.hub.Clients(auctionID))
	s.hub.Broadcast(auctionID, data)
	log.Printf("WS: broadcast done for bid %d", bid.ID)

	return bid, nil
}

// PaginatedBids はページ番号とサイズで入札一覧と総件数を取得します
func (s *BidService) PaginatedBids(auctionID uint, page, size int) ([]model.Bid, int64, error) {
	// ページとサイズの最低値設定
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	// 1) ページングされた入札を取得
	bids, err := s.Repo.FindPaginated(auctionID, offset, size)
	if err != nil {
		return nil, 0, err
	}
	// 2) 総入札件数を取得
	total, err := s.Repo.CountByAuction(auctionID)
	if err != nil {
		return nil, 0, err
	}
	return bids, total, nil
}
