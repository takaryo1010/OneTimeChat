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
	Mu              sync.Mutex       // スレッドセーフにするためのミューテックス
}

// Room は個々のチャットルームを表す構造体
type Room struct {
	ID                     string     `json:"ID"`             // ルームID
	Name                   string     `json:"name"`           // ルーム名
	Owner                  string     `json:"owner"`          // ルームのオーナー
	OwnerSessionID         string     `json:"ownerSessionID"` // オーナーのセッションID
	Expires                time.Time  `json:"expires"`        // 有効期限
	RequiresAuth           bool       `json:"requiresAuth"`   // 認証が必要かどうか
	UnauthenticatedClients []*Client  // ルームへの接続許可待ちのクライアント
	AuthenticatedClients   []*Client  //ルームへの接続許可がされているクライアント
	Mu                     sync.Mutex // スレッドセーフにするためのミューテックス
}

// ResponseRoom
type ResponseRoom struct {
	ID                     string            `json:"ID"`                     // ルームID
	Name                   string            `json:"name"`                   // ルーム名
	Owner                  string            `json:"owner"`                  // ルームのオーナー
	Expires                time.Time         `json:"expires"`                // 有効期限
	RequiresAuth           bool              `json:"requiresAuth"`           // 認証が必要かどうか
	UnauthenticatedClients []*ResponseClient `json:"unauthenticatedClients"` // ルームへの接続許可待ちのクライアント
	AuthenticatedClients   []*ResponseClient `json:"authenticatedClients"`   // ルームへの接続許可がされているクライアント
}

// Client はチャットルームに参加しているユーザーを表す構造体
type Client struct {
	Name      string          // クライアント名
	ClientID  string          // クライアントID
	SessionID string          // セッションID
	Ws        *websocket.Conn // WebSocket接続
}

// ResponseClient はクライアント情報を表す構造体
type ResponseClient struct {
	Name     string `json:"name"`     // クライアント名
	ClientID string `json:"clientID"` // クライアントID
}

// Message はチャットメッセージを表す構造体
type Message struct {
	RoomID    string `json:"room_id"`   // ルームID
	Sentence  string `json:"sentence"`  // メッセージ本文
	Sender    string `json:"sender"`    // 送信者
	Timestamp int64  `json:"timestamp"` // タイムスタンプ
	Type      string `json:"type"`      // メッセージの種類
}

// Participant は参加者を表す構造体
type Participant struct {
	Name     string `json:"name"`
	ClientID string `json:"clientid"`
	IsOwner  bool   `json:"isowner"`
}
