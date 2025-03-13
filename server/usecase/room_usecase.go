package usecase

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/takaryo1010/OneTimeChat/server/model"
)

type RoomUsecase struct {
	RoomManager *model.RoomManager
	upgrader    websocket.Upgrader
}



// NewRoomUsecase creates a new RoomUsecase instance.
func NewRoomUsecase() *RoomUsecase {
	return &RoomUsecase{
		RoomManager: &model.RoomManager{
			Rooms:           make(map[string]*model.Room),
			ExpireSortRooms: []*model.Room{},
			Mu:              sync.Mutex{},
		},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// CreateRoom 新しい部屋を作る
func (uc *RoomUsecase) CreateRoom(room *model.Room, generatedSessionID string) (*model.Room, error) {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	roomID := generateRoomID(uc.RoomManager) // 任意のID生成関数を使用
	room = &model.Room{
		ID:                     roomID,
		Name:                   room.Name,
		Owner:                  room.Owner,
		Expires:                room.Expires,
		RequiresAuth:           true,//TODO
		UnauthenticatedClients: []*model.Client{},
		AuthenticatedClients:   []*model.Client{}, // 初期化
		Mu:                     sync.Mutex{},
	}

	// 部屋を作成し、マネージャーに登録
	uc.RoomManager.Rooms[roomID] = room
	appendExpireBinarySearch(uc.RoomManager, room)

	// オーナーを部屋に追加
	client := &model.Client{
		Name: room.Owner,
		SessionID: generatedSessionID,
		Ws: nil,
	}

	room.AuthenticatedClients = append(room.AuthenticatedClients, client)

	return room, nil
}

// GetRoomByID retrieves a room by its ID.
func (uc *RoomUsecase) GetRoomByID(roomID string) (*model.Room, error) {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	room, exists := uc.RoomManager.Rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}
	return room, nil
}

// JoinRoom allows a client to join a room.
func (uc *RoomUsecase) JoinRoom(roomID, clientName,generatedSessionID string) error {
	uc.RoomManager.Mu.Lock()
	room, exists := uc.RoomManager.Rooms[roomID]
	uc.RoomManager.Mu.Unlock()
	
	if !exists {
		return errors.New("room not found")
	}
	
	client := &model.Client{
		Name: clientName,
		SessionID: generatedSessionID,
		Ws: nil,
	}
	
	room.Mu.Lock()
	defer room.Mu.Unlock()
	if room.RequiresAuth {
		room.AuthenticatedClients = append(room.AuthenticatedClients, client)
	} else {
		room.UnauthenticatedClients = append(room.UnauthenticatedClients, client)
	}
		
	fmt.Println("Authenticated Clients")
	for _, client := range room.AuthenticatedClients {
		fmt.Println("Name: ", client.Name)
		fmt.Println("SessionID: ", client.SessionID)
		}
	fmt.Println("----------------------")
	return nil

}

// HandleWebSocketConnection handles a WebSocket connection for a client.
func (uc *RoomUsecase) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request, roomID, clientName,sessionID string) error {
	// 部屋を取得
	uc.RoomManager.Mu.Lock()
	room, exists := uc.RoomManager.Rooms[roomID]
	uc.RoomManager.Mu.Unlock()

	if !exists {
		return errors.New("room not found")
	}

	// WebSocket 接続のアップグレード
	conn, err := uc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	// 部屋内に既に存在する仮のクライアントを検索
	var client *model.Client
	for _, c := range room.AuthenticatedClients {
		fmt.Println("Name: ", c.Name)
		fmt.Println("SessionID: ", c.SessionID)
		fmt.Println("received SessionID: ", sessionID)
		if c.SessionID == sessionID {
			client = c
			break
		}
	}

	if client == nil {
		// 仮のクライアントが見つからない場合はエラーを返す
		return errors.New("client not found in the room")
	}

	// 仮のクライアントの WebSocket 接続を更新
	client.Ws = conn

	// WebSocket 接続を確立したことをログ出力
	fmt.Printf("Client %s connected to room %s\n", clientName, roomID)

	// WebSocket のメッセージ受信ループを開始
	go func() {
		defer conn.Close() // 関数が終了したら接続を閉じる

		for {
			// クライアントからメッセージを受信
			_, msg, err := conn.ReadMessage()
			if err != nil {
				// メッセージの読み込みエラーが発生した場合、ループを終了
				break
			}
			// 受信したメッセージを他のクライアントにブロードキャスト
			uc.broadcastToRoom(roomID, msg,clientName,sessionID)
		}
	}()

	return nil
}

// broadcastToRoom broadcasts a message with sender information, room ID, and timestamp.
func (uc *RoomUsecase) broadcastToRoom(roomID string, sentence []byte, sender,sessionID string) {

	isClientInRoom := false
	for _, c := range uc.RoomManager.Rooms[roomID].AuthenticatedClients {
		if c.SessionID == sessionID {
			isClientInRoom = true
			break
		}
	}
	if !isClientInRoom {
		return
	}

	uc.RoomManager.Mu.Lock()
	room, exists := uc.RoomManager.Rooms[roomID]
	uc.RoomManager.Mu.Unlock()

	stringSentence := string(sentence)

	if !exists {
		return
	}

	room.Mu.Lock()
	defer room.Mu.Unlock()

	// メッセージデータを作成
	message := model.Message{
		RoomID:    roomID,
		Sentence:  stringSentence,
		Sender:    sender,
		Timestamp: time.Now().Unix(), // 現在のUNIXタイムスタンプ
	}

	// メッセージをJSONにエンコード
	messageJSON, err := json.Marshal(message)
	if err != nil {
		// エンコードエラー時の処理
		return
	}
	fmt.Println("Message sent:", string(messageJSON))
	fmt.Println("Sender:", sender)
	fmt.Println("SessionID:", sessionID)

	// 各クライアントにJSONメッセージを送信
	for _, client := range room.AuthenticatedClients {
		err := client.Ws.WriteMessage(websocket.TextMessage, messageJSON)
		if err != nil {
			client.Ws.Close()
		}
	}
}

