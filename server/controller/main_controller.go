package controller

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/takaryo1010/OneTimeChat/server/usecase"
)

type MainController struct {
	RoomUsecase *usecase.RoomUsecase
}

// WebSocketHandler handles WebSocket connections.
func (mc *MainController) WebSocketHandler(c echo.Context) error {
	roomID := c.QueryParam("room_id")
	clientName := c.QueryParam("client_name")

	err := mc.RoomUsecase.HandleWebSocketConnection(c.Response(), c.Request(), roomID, clientName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return nil
}
func GenerateSessionID() (string, error) {
	// 16 バイトのランダムなデータを生成
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// base64 エンコードしてセッションIDとして利用
	return base64.URLEncoding.EncodeToString(b), nil
}
func (mc *MainController) CreateRoom(c echo.Context) error {
	// セッションIDの生成
	sessionID, err := GenerateSessionID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate session ID"})
	}

	// セッションIDをクッキーに保存
	c.SetCookie(&http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // セッションの有効期限
	})

	// オーナー情報はセッションIDを基に後で取得するために設定する（仮にユーザーIDを使用）
	owner := "someUserID" // ここでは仮のユーザーIDを使用します

	// フォームからルーム名を取得
	roomName := c.FormValue("name")

	// ルーム作成処理
	room, err := mc.RoomUsecase.CreateRoom(roomName, owner)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	fmt.Println("Room created:", room.ID)
	// ルーム作成成功時に返す
	return c.JSON(http.StatusOK, room)
}

// GetRoom retrieves a specific room by ID.
func (mc *MainController) GetRoom(c echo.Context) error {
	roomID := c.Param("id")
	room, err := mc.RoomUsecase.GetRoomByID(roomID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, room)
}

// JoinRoom allows a client to join a room.
func (mc *MainController) JoinRoom(c echo.Context) error {
	roomID := c.Param("id")
	clientName := c.FormValue("client_name")
  
	err := mc.RoomUsecase.JoinRoom(roomID, clientName)
	if err != nil {
	  return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
  
	// 部屋に参加したことを確認
	fmt.Println("Client joined:", clientName)
	fmt.Println("Room ID:", roomID)
	return c.NoContent(http.StatusOK)
  }
  

// Other handlers (GetParticipants, UpdateRoomSettings, etc.) can follow the same pattern.
