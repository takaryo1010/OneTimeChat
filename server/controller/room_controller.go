package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/takaryo1010/OneTimeChat/server/usecase"
	"github.com/takaryo1010/OneTimeChat/server/model"
)

type MainController struct {
	RoomUsecase *usecase.RoomUsecase
}



func (mc *MainController) CreateRoom(c echo.Context) error {
	
	
	// フォームからルーム名を取得
	var req model.Room
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	
	roomName := req.Name
	fmt.Println("Room name:", roomName)
	owner := req.Owner
	fmt.Println("Owner:", owner)
	// ルーム作成処理
	room, err := mc.RoomUsecase.CreateRoom(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// セッションIDをクッキーに保存
	c.SetCookie(&http.Cookie{
		Name:    "session_id",
		Value:   room.OwnerSessionID,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // セッションの有効期限
	})

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

	
	
	sessionID,err := mc.RoomUsecase.JoinRoom(roomID, clientName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	// セッションIDをクッキーに保存
	c.SetCookie(&http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // セッションの有効期限
	})
	// 部屋に参加したことを確認
	fmt.Println("Client joined:", clientName)
	fmt.Println("Room ID:", roomID, "Session ID:", sessionID)

	return c.JSON(http.StatusOK, map[string]string{"roomID": roomID, "sessionID": sessionID})
}


// Authenticate authenticates a client to join a room.
func (mc *MainController) Authenticate(c echo.Context) error {
	roomID := c.Param("id")
	clientSessionID := c.QueryParam("client_session_id")
	ownerSessionID := c.QueryParam("owner_session_id")
	fmt.Println("Authentication requested for room:", roomID, "client session:", clientSessionID, "owner session:", ownerSessionID)
	err := mc.RoomUsecase.Authenticate(roomID, clientSessionID, ownerSessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return nil
}