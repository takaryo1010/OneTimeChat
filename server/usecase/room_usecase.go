package usecase

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

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
func (uc *RoomUsecase) CreateRoom(room *model.Room) (*model.Room, error) {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	// セッションIDの生成
	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil, err
	}

	roomID := generateRoomID(uc.RoomManager) // 任意のID生成関数を使用
	room = &model.Room{
		ID:                     roomID,
		Name:                   room.Name,
		Owner:                  room.Owner,
		Expires:                room.Expires,
		RequiresAuth:           room.RequiresAuth,
		OwnerSessionID:         sessionID,
		UnauthenticatedClients: []*model.Client{},
		AuthenticatedClients:   []*model.Client{}, // 初期化
		Mu:                     sync.Mutex{},
	}

	// 部屋を作成し、マネージャーに登録
	uc.RoomManager.Rooms[roomID] = room
	appendExpireBinarySearch(uc.RoomManager, room)

	// オーナーを部屋に追加
	client := &model.Client{
		Name:      room.Owner,
		SessionID: sessionID,
		Ws:        nil,
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
func (uc *RoomUsecase) JoinRoom(roomID, clientName string) (string, error) {
	uc.RoomManager.Mu.Lock()
	room, exists := uc.RoomManager.Rooms[roomID]
	uc.RoomManager.Mu.Unlock()

	if !exists {
		return "", errors.New("room not found")
	}

	// セッションIDの生成
	generatedSessionID, err := GenerateSessionID()
	if err != nil {
		return "", err
	}

	client := &model.Client{
		Name:      clientName,
		SessionID: generatedSessionID,
		Ws:        nil,
	}

	room.Mu.Lock()
	defer room.Mu.Unlock()
	if room.RequiresAuth {
		room.UnauthenticatedClients = append(room.UnauthenticatedClients, client)
	} else {
		room.AuthenticatedClients = append(room.AuthenticatedClients, client)
	}

	fmt.Println("Authenticated Clients")
	for _, client := range room.AuthenticatedClients {
		fmt.Println("Name: ", client.Name)
		fmt.Println("SessionID: ", client.SessionID)
	}
	fmt.Println("----------------------")
	return generatedSessionID, nil

}

func (uc *RoomUsecase) Authenticate(roomID, client_session_id, owner_session_id string) error {
	uc.RoomManager.Mu.Lock()
	room, exists := uc.RoomManager.Rooms[roomID]
	uc.RoomManager.Mu.Unlock()

	if !exists {
		return errors.New("room not found")
	}

	if room.OwnerSessionID != owner_session_id {
		return errors.New("you are not the owner of this room")
	}

	room.Mu.Lock()
	defer room.Mu.Unlock()

	isClientInRoom := false

	for i, client := range room.UnauthenticatedClients {
		if client.SessionID == client_session_id {
			room.AuthenticatedClients = append(room.AuthenticatedClients, client)
			room.UnauthenticatedClients = append(room.UnauthenticatedClients[:i], room.UnauthenticatedClients[i+1:]...)
			isClientInRoom = true
			break
		}
	}

	if !isClientInRoom {
		return errors.New("client not found in the room")
	}

	return nil
}

func (uc *RoomUsecase) UpdateRoomSettings(roomID string, newRoomSettings *model.Room,owner_session_id string) (*model.Room, error) {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	room, exists := uc.RoomManager.Rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}
	if room.OwnerSessionID != owner_session_id {
		return nil, errors.New("you are not the owner of this room")
	}

	room.Name = newRoomSettings.Name
	room.RequiresAuth = newRoomSettings.RequiresAuth
	
	
	return room, nil
}

func (uc *RoomUsecase) DeleteRoom(roomID, owner_session_id string) error {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	room, exists := uc.RoomManager.Rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}
	if room.OwnerSessionID != owner_session_id {
		return errors.New("you are not the owner of this room")
	}

	delete(uc.RoomManager.Rooms, roomID)
	return nil
}