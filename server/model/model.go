package model

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RoomManager はチャットルーム全体の管理を行う構造体
type RoomManager struct {
	Rooms           map[string]*Room // ルームのマップ
	ExpireSortRooms []*Room          // 期限順に並べたルーム
	mu              sync.Mutex       // スレッドセーフにするためのミューテックス
}

// Room は個々のチャットルームを表す構造体
type Room struct {
	ID                     string     // ルームID
	Name                   string     // ルーム名
	Owner                  string     // ルームのオーナー
	Expires                time.Time  // 有効期限
	RequiresAuth           bool       // 認証が必要かどうか
	UnauthenticatedClients []*Client  // ルームへの接続許可待ちのクライアント
	AuthenticatedClients   []*Client  //ルームへの接続許可がされているクライアント
	Mu                     sync.Mutex // スレッドセーフにするためのミューテックス
}

// Client はチャットルームに参加しているユーザーを表す構造体
type Client struct {
	Name string          // クライアント名
	Ws   *websocket.Conn // WebSocket接続
}

// Message はチャットメッセージを表す構造体
type Message struct {
	RoomID    string `json:"room_id"`   // ルームID
	Sentence  string `json:"sentence"`  // メッセージ本文
	Sender    string `json:"sender"`    // 送信者
	Timestamp int64  `json:"timestamp"` // タイムスタンプ
}
