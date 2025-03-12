package controller

import (
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
	fmt.Println("WebSocket connection requested for room:", roomID, "client:", clientName)
	err := mc.RoomUsecase.HandleWebSocketConnection(c.Response(), c.Request(), roomID, clientName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return nil
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

	// フォームからルーム名を取得
	type CreateRoomRequest struct {
		Name  string `json:"name"`
		Owner string `json:"owner"`
	}
	var req CreateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	roomName := req.Name
	fmt.Println("Room name:", roomName)
	owner := req.Owner
	fmt.Println("Owner:", owner)
	// ルーム作成処理
	room, err := mc.RoomUsecase.CreateRoom(roomName, owner, sessionID)
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
	type JoinRoomRequest struct {
		ClientName string `json:"client_name"`
	}
	var req JoinRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	clientName := req.ClientName
	fmt.Println("Client name:", clientName)

	sessionID, err := GenerateSessionID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate session ID"})
	}

	err = mc.RoomUsecase.JoinRoom(roomID, clientName, sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 部屋に参加したことを確認
	fmt.Println("Client joined:", clientName)
	fmt.Println("Room ID:", roomID, "Session ID:", sessionID)

	return c.JSON(http.StatusOK, map[string]string{"roomID": roomID, "sessionID": sessionID})
}

// Other handlers (GetParticipants, UpdateRoomSettings, etc.) can follow the same pattern.
