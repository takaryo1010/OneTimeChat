package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

// WebSocketアップグレード用
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// チャットルームのクライアント管理用マップ
var rooms = make(map[string]map[*websocket.Conn]bool)
var roomMessages = make(chan roomMessage)

// メッセージ送信時のデータ構造
type roomMessage struct {
	room    string
	message string
}
type Message struct {
	Text string `json:"text"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// WebSocket接続をアップグレード
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	// ルーム名を取得（例: クエリパラメータ "room"）
	room := r.URL.Query().Get("room")
	if room == "" {
		room = "default" // デフォルトのルーム
	}

	// ルームにクライアントを追加
	if rooms[room] == nil {
		rooms[room] = make(map[*websocket.Conn]bool)
	}
	rooms[room][ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("error reading JSON:", err)
			fmt.Println("Received data:", msg)
			delete(rooms[room], ws)
			break
		}
		// メッセージをルームにブロードキャスト
		roomMessages <- roomMessage{room: room, message: msg.Text}
	}
}

func handleMessages() {
	for {
		// ブロードキャストチャンネルからメッセージを取得
		msg := <-roomMessages
		// ルームに接続しているすべてのクライアントに送信
		for client := range rooms[msg.room] {
			err := client.WriteJSON(msg.message)
			if err != nil {
				fmt.Println("error:", err)
				client.Close()
				delete(rooms[msg.room], client)
			}
		}
	}
}

func main() {
	// ルートディレクトリ
	// fs := http.FileServer(http.Dir("./public"))
	// http.Handle("/", fs)

	// WebSocket接続
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Chat server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("error:", err)
	}
}