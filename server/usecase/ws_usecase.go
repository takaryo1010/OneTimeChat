package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/takaryo1010/OneTimeChat/server/model"
)

// HandleWebSocketConnection handles a WebSocket connection for a client.
func (uc *RoomUsecase) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request, roomID, clientName, sessionID string) error {
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

		if c.SessionID == sessionID {
			client = c
			break
		}
	}
	if client == nil {
		for _, c := range room.UnauthenticatedClients {

			if c.SessionID == sessionID {
				client = c
				break
			}

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
			uc.broadcastToRoom(roomID, msg, clientName, sessionID)
		}
	}()

	return nil
}

// broadcastToRoom broadcasts a message with sender information, room ID, and timestamp.
func (uc *RoomUsecase) broadcastToRoom(roomID string, sentence []byte, sender, sessionID string) {
	isClientInRoom := false
	for _, c := range uc.RoomManager.Rooms[roomID].AuthenticatedClients {
		if c.SessionID == sessionID {
			isClientInRoom = true
			break
		}
	}


	stringSentence := string(sentence)

	type MessageType struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}

	var messageType MessageType
	err := json.Unmarshal(sentence, &messageType)
	if err != nil {
		return
	}
	fmt.Println(messageType.Type)
	fmt.Println(messageType.Content)



	if !isClientInRoom && messageType.Type == "message" {
		return
	}

	uc.RoomManager.Mu.Lock()
	room, exists := uc.RoomManager.Rooms[roomID]
	uc.RoomManager.Mu.Unlock()

	if !exists {
		return
	}

	room.Mu.Lock()
	defer room.Mu.Unlock()

	// Removed the initial assignment to message here
	if messageType.Type == "message" {
		// メッセージデータを作成
		message := model.Message{
			RoomID:    roomID,
			Sentence:  messageType.Content,
			Sender:    sender,
			Timestamp: time.Now().Unix(), // 現在のUNIXタイムスタンプ
			Type:      "message",
		}
		// メッセージをJSONにエンコード
		messageJSON, err := json.Marshal(message)
		if err != nil {
			// エンコードエラー時の処理
			return
		}

		// 各クライアントにJSONメッセージを送信
		for _, client := range room.AuthenticatedClients {
			if client.SessionID == sessionID {
				continue
			}
			err := client.Ws.WriteMessage(websocket.TextMessage, messageJSON)
			if err != nil {
				client.Ws.Close()
			}
		}
	} else if messageType.Type == "participants_update" {
		// メッセージデータを作成
		message := model.Message{
			RoomID:    roomID,
			Sentence:  stringSentence,
			Sender:    sender,
			Timestamp: time.Now().Unix(), // 現在のUNIXタイムスタンプ
			Type:      "participants_update",
		}
		// メッセージをJSONにエンコード
		messageJSON, err := json.Marshal(message)
		if err != nil {
			// エンコードエラー時の処理
			return
		}

		for _, client := range room.AuthenticatedClients {
			err := client.Ws.WriteMessage(websocket.TextMessage, messageJSON)
			if err != nil {
				client.Ws.Close()
			}
		}
		for _, client := range room.UnauthenticatedClients {
			err := client.Ws.WriteMessage(websocket.TextMessage, messageJSON)
			if err != nil {
				client.Ws.Close()
			}
		}

	}
}
