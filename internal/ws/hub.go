package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client は WebSocket クライアント接続情報を表します
type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

// Hub はオークションごとにクライアントを管理するハブです
type Hub struct {
	mu         sync.Mutex
	clients    map[uint]map[*Client]bool // auctionID → set of clients
	register   chan subscription
	unregister chan subscription
	broadcast  chan event
}

// event はブロードキャスト対象データを表します
type subscription struct {
	AuctionID uint
	Client    *Client
}

// subscription はハブへの登録／解除リクエストを表します
type event struct {
	AuctionID uint
	Data      []byte
}

// NewHub は新しい Hub を生成します
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		broadcast:  make(chan event),
	}
}

// Clients は指定オークションIDに接続しているクライアント数を返します
func (h *Hub) Clients(auctionID uint) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients[auctionID])
}

// Run はハブの内部イベントループを開始します
// register, unregister, broadcast チャンネルのイベントを処理します
func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.register:
			log.Printf("WS HUB: REGISTER auction=%d, total_clients=%d", sub.AuctionID, len(h.clients[sub.AuctionID])+1)
			h.mu.Lock()
			if h.clients[sub.AuctionID] == nil {
				h.clients[sub.AuctionID] = make(map[*Client]bool)
			}
			h.clients[sub.AuctionID][sub.Client] = true
			h.mu.Unlock()
		case sub := <-h.unregister:
			log.Printf("WS HUB: UNREGISTER auction=%d, remaining_clients=%d", sub.AuctionID, len(h.clients[sub.AuctionID])-1)
			h.mu.Lock()
			if conns, ok := h.clients[sub.AuctionID]; ok {
				delete(conns, sub.Client)
				if len(conns) == 0 {
					delete(h.clients, sub.AuctionID)
				}
			}
			h.mu.Unlock()
		case ev := <-h.broadcast:
			log.Printf("WS HUB: BROADCAST auction=%d, data_len=%d", ev.AuctionID, len(ev.Data))
			h.mu.Lock()
			for client := range h.clients[ev.AuctionID] {
				select {
				case client.Send <- ev.Data:
				default:
					// 채널이 가득 차면 연결 제거
					close(client.Send)
					delete(h.clients[ev.AuctionID], client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// Register は hub にクライアントを登録します
func (h *Hub) Register(auctionID uint, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns := h.clients[auctionID]
	if conns == nil {
		conns = make(map[*Client]bool)
		h.clients[auctionID] = conns
	}
	conns[c] = true
}

// Unregister は hub からクライアントを解除し、チャネルを閉じます
func (h *Hub) Unregister(auctionID uint, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns := h.clients[auctionID]; conns != nil {
		delete(conns, c)
		close(c.Send)
	}
}

// Broadcast は指定オークションIDのクライアントにメッセージを送信します
// チャネルが満杯の場合は当該クライアントを切断します
func (h *Hub) Broadcast(auctionID uint, msg []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients[auctionID] {
		select {
		case c.Send <- msg:
		default:
			// 채널이 가득 차면 끊기
			delete(h.clients[auctionID], c)
			close(c.Send)
		}
	}
}
