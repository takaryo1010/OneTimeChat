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
func (uc *RoomUsecase) CreateRoom(room *model.Room) (*model.ResponseRoom, string, error) {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	// セッションIDの生成
	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil, "", err
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
		ClientID:  GeneratedClientID(uc.RoomManager),
		SessionID: sessionID,
		Ws:        nil,
	}

	room.AuthenticatedClients = append(room.AuthenticatedClients, client)

	res := changedForResponse(room)
	return res, client.SessionID, nil

}

// GetRoomByID retrieves a room by its ID.
func (uc *RoomUsecase) GetRoomByID(roomID string) (*model.ResponseRoom, error) {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	room, exists := uc.RoomManager.Rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}
	res := changedForResponse(room)
	return res, nil
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
		ClientID:  GeneratedClientID(uc.RoomManager),
		SessionID: generatedSessionID,
		Ws:        nil,
	}
	fmt.Println("Client joined:", clientName)
	fmt.Println("Client ID:", client.ClientID)
	room.Mu.Lock()
	defer room.Mu.Unlock()
	if room.RequiresAuth {
		room.UnauthenticatedClients = append(room.UnauthenticatedClients, client)
	} else {
		room.AuthenticatedClients = append(room.AuthenticatedClients, client)
	}

	return generatedSessionID, nil

}

func (uc *RoomUsecase) Authenticate(roomID, client_id, owner_session_id string) error {
	uc.RoomManager.Mu.Lock()
	room, exists := uc.RoomManager.Rooms[roomID]
	uc.RoomManager.Mu.Unlock()

	if !exists {
		return errors.New("room not found")
	}

	//オーナーのクライアントＩＤを取得
	room.Mu.Lock()
	defer room.Mu.Unlock()

	fmt.Println("OwnerSessionID:", room.OwnerSessionID)
	fmt.Println("OwnerSessionID:", owner_session_id)
	if owner_session_id != room.OwnerSessionID {
		return errors.New("you are not the owner of this room")
	}

	isClientInRoom := false
	fmt.Println("ClientID:", client_id)
	for i, client := range room.UnauthenticatedClients {
		if client.ClientID == client_id {
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

func (uc *RoomUsecase) UpdateRoomSettings(roomID string, newRoomSettings *model.Room, owner_session_id string) (*model.ResponseRoom, error) {
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

	res := changedForResponse(room)

	return res, nil
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

func (uc *RoomUsecase) KickParticipant(roomID, client_id, owner_session_id string) error {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	room, exists := uc.RoomManager.Rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}
	if room.OwnerSessionID != owner_session_id {
		return errors.New("you are not the owner of this room")
	}

	room.Mu.Lock()
	defer room.Mu.Unlock()

	if client_id == "" {
		return errors.New("client_id is required")
	}

	isClientInRoom := false
	for i, client := range room.AuthenticatedClients {
		if client.ClientID == client_id {
			room.AuthenticatedClients = append(room.AuthenticatedClients[:i], room.AuthenticatedClients[i+1:]...)
			isClientInRoom = true
			break
		}
	}

	if !isClientInRoom {
		return errors.New("client not found in the room")
	}

	return nil
}

func (uc *RoomUsecase) LeaveRoom(roomID, client_session_id string) error {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	room, exists := uc.RoomManager.Rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}

	room.Mu.Lock()
	defer room.Mu.Unlock()

	isClientInRoom := false
	for i, client := range room.AuthenticatedClients {
		if client.SessionID == client_session_id {
			room.AuthenticatedClients = append(room.AuthenticatedClients[:i], room.AuthenticatedClients[i+1:]...)
			isClientInRoom = true
			break
		}
	}

	if !isClientInRoom {
		return errors.New("client not found in the room")
	}

	return nil
}

func (uc *RoomUsecase) IsAuth(roomID, clientSessionID string) (bool, error) {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	room, exists := uc.RoomManager.Rooms[roomID]
	if !exists {
		return false, errors.New("room not found")
	}

	room.Mu.Lock()
	defer room.Mu.Unlock()

	for _, client := range room.AuthenticatedClients {
		if client.SessionID == clientSessionID {
			return true, nil
		}
	}

	return false, nil
}

func (uc *RoomUsecase) GetParticipants(roomID string) ([]model.Participant, []model.Participant, error) {
	uc.RoomManager.Mu.Lock()
	defer uc.RoomManager.Mu.Unlock()

	room, exists := uc.RoomManager.Rooms[roomID]
	if !exists {
		return nil, nil, errors.New("room not found")
	}

	room.Mu.Lock()
	defer room.Mu.Unlock()

	participants := make([]model.Participant, 0)
	for _, client := range room.AuthenticatedClients {
		if client.SessionID == room.OwnerSessionID {
			participants = append(participants, model.Participant{Name: client.Name, ClientID: client.ClientID, IsOwner: true})
		} else {
			participants = append(participants, model.Participant{Name: client.Name, ClientID: client.ClientID, IsOwner: false})
		}
	}

	unauthenticatedClients := make([]model.Participant, 0)
	for _, client := range room.UnauthenticatedClients {
		unauthenticatedClients = append(unauthenticatedClients, model.Participant{Name: client.Name, ClientID: client.ClientID, IsOwner: false})
	}

	return participants, unauthenticatedClients, nil
}
