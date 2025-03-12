package usecase

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

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
func (uc *RoomUsecase) CreateRoom(name, owner, generatedSessionID string) (*model.Room, error) {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	roomID := generateRoomID(uc.RoomManager) // 任意のID生成関数を使用
	room := &model.Room{
		ID:                     roomID,
		Name:                   name,
		Owner:                  owner,
		Expires:                time.Now().Add(24 * time.Hour), // 例: 24時間有効
		RequiresAuth:           true,
		UnauthenticatedClients: []*model.Client{},
		AuthenticatedClients:   []*model.Client{}, // 初期化
		Mu:                     sync.Mutex{},
	}

	// 部屋を作成し、マネージャーに登録
	uc.RoomManager.Rooms[roomID] = room
	appendExpireBinarySearch(uc.RoomManager, room)

	// オーナーを部屋に追加
	client := &model.Client{
		Name: owner,
		SessionID: generatedSessionID,
		Ws: &websocket.Conn{},
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
		Ws: &websocket.Conn{},
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
func (uc *RoomUsecase) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request, roomID, clientName string) error {
	uc.RoomManager.Mu.Lock()
	room, exists := uc.RoomManager.Rooms[roomID]
	uc.RoomManager.Mu.Unlock()

	if !exists {
		return errors.New("room not found")
	}

	conn, err := uc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}


	

	room.Mu.Lock()
	defer room.Mu.Unlock()

	// WebSocketの受信ループ
	go func() {
		defer conn.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			fmt.Println(string(msg))
			uc.broadcastToRoom(roomID, msg)
		}
	}()

	return nil
}

// broadcastToRoom broadcasts a message to all authenticated clients in a room.
func (uc *RoomUsecase) broadcastToRoom(roomID string, message []byte) {
	uc.RoomManager.Mu.Lock()
	room, exists := uc.RoomManager.Rooms[roomID]
	uc.RoomManager.Mu.Unlock()

	if !exists {
		return
	}

	room.Mu.Lock()
	defer room.Mu.Unlock()

	for _, client := range room.AuthenticatedClients {
		err := client.Ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			client.Ws.Close()
		}
	}
}
