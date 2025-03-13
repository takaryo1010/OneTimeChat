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
		fmt.Println("Name: ", c.Name)
		fmt.Println("SessionID: ", c.SessionID)
		fmt.Println("received SessionID: ", sessionID)
		if c.SessionID == sessionID {
			client = c
			break
		}
	}
	if client == nil {
		for _, c := range room.UnauthenticatedClients {
			fmt.Println("Name: ", c.Name)
			fmt.Println("SessionID: ", c.SessionID)
			fmt.Println("received SessionID: ", sessionID)
			if c.SessionID == sessionID {
				fmt.Println("This client is unauthenticated")
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
