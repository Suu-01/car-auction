package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/ksj/car-auction/internal/ws"
)

// WebSocket アップグレーダー設定: 任意のオリジンを許可
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// RegisterWSRoutes は WebSocket エンドポイントを登録します。
// /ws/auctions/{id} に接続されたクライアントを指定のオークションハブに登録します。
func RegisterWSRoutes(r *mux.Router, hub *ws.Hub) {
	r.HandleFunc("/ws/auctions/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		// パスパラメータからオークション ID を取得
		vars := mux.Vars(r)
		aid, _ := strconv.Atoi(vars["id"])
		log.Printf("WS: upgrade requested for auction %d", aid)

		// WebSocket 接続へのアップグレードを試みる
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WS upgrade error: %v", err)
			return
		}
		log.Printf("WS: client connected to auction %d", aid)
		client := &ws.Client{Conn: conn, Send: make(chan []byte, 16)}
		hub.Register(uint(aid), client)

		// 読み取りゴルーチン: クライアントからのメッセージを読み捨て、切断時にクリーンアップ
		go func() {
			defer func() {
				hub.Unregister(uint(aid), client)
				conn.Close()
			}()
			for {
				if _, _, err := conn.NextReader(); err != nil {
					break
				}
			}
		}()

		// 書き込みゴルーチン: hub からのメッセージを client.Send チャネルで受け取り、クライアントに送信
		go func() {
			for msg := range client.Send {
				conn.WriteMessage(websocket.TextMessage, msg)
			}
		}()
	})
}
