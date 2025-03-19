package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo"
	"github.com/takaryo1010/OneTimeChat/server/model"
	"github.com/takaryo1010/OneTimeChat/server/usecase"
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
	owner := req.Owner
	fmt.Println("Room name:", roomName, "by", owner)
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
	c.SetCookie(&http.Cookie{
		Name:    "room_id",
		Value:   room.ID, // ここにRoomIDを設定
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // セッションの有効期限
	})
	c.SetCookie(&http.Cookie{
		Name:    "user_name",
		Value:   url.QueryEscape(owner),
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // セッションの有効期限
	})

	c.SetCookie(&http.Cookie{
		Name:    "is_owner",
		Value:   "true",
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

	sessionID, err := mc.RoomUsecase.JoinRoom(roomID, clientName)
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
	c.SetCookie(&http.Cookie{
		Name:    "room_id",
		Value:   roomID, // ここにRoomIDを設定
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // セッションの有効期限
	})
	c.SetCookie(&http.Cookie{
		Name:    "user_name",
		Value:   url.QueryEscape(clientName),
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // セッションの有効期限
	})
	c.SetCookie(&http.Cookie{
		Name:    "is_owner",
		Value:   "false",
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // セッションの有効期限
	})
	// 部屋に参加したことを確認
	fmt.Println("Client joined:", clientName, "in room:", roomID)

	return c.JSON(http.StatusOK, map[string]string{"roomID": roomID, "sessionID": sessionID})
}

// Authenticate authenticates a client to join a room.
func (mc *MainController) Authenticate(c echo.Context) error {
	roomID := c.Param("id")
	clientSessionID := c.QueryParam("client_session_id")
	ownerSessionID := c.QueryParam("owner_session_id")
	err := mc.RoomUsecase.Authenticate(roomID, clientSessionID, ownerSessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("Client authenticated:", clientSessionID, "in room:", roomID)
	return nil
}

// 参加者を取得(notオーナー用つまり参加者がほかの参加者を確認する用)
func (mc *MainController) GetParticipants(c echo.Context) error {
	roomID := c.Param("id")
	room, err := mc.RoomUsecase.GetRoomByID(roomID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	type Participant struct {
		Name string `json:"name"`
	}
	participants := make([]Participant, 0)
	for _, client := range room.AuthenticatedClients {
		participants = append(participants, Participant{Name: client.Name})
	}

	return c.JSON(http.StatusOK, participants)
}

// ルームの設定変更(オーナー専用)
// RoomNameとrequiresAuthをjsonで必ず受け取る
func (mc *MainController) UpdateRoomSettings(c echo.Context) error {
	roomID := c.Param("id")
	ownerSessionID := c.QueryParam("owner_session_id")
	var req model.Room
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	// ルーム作成処理
	room, err := mc.RoomUsecase.UpdateRoomSettings(roomID, &req, ownerSessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("Room updated:", room.ID)
	return c.JSON(http.StatusOK, room)
}

// ルームの削除(オーナー専用)
func (mc *MainController) DeleteRoom(c echo.Context) error {
	roomID := c.Param("id")
	ownerSessionID := c.QueryParam("owner_session_id")
	err := mc.RoomUsecase.DeleteRoom(roomID, ownerSessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("Room deleted:", roomID)
	return c.JSON(http.StatusOK, map[string]string{"message": "room deleted"})
}

// 参加者をキック(オーナー専用)
func (mc *MainController) KickParticipant(c echo.Context) error {
	roomID := c.Param("id")
	ownerSessionID := c.QueryParam("owner_session_id")
	clientSessionID := c.QueryParam("client_session_id")

	if ownerSessionID == clientSessionID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "you can't kick yourself"})
	}

	if clientSessionID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "client_session_id is required"})
	}

	err := mc.RoomUsecase.KickParticipant(roomID, clientSessionID, ownerSessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("Client kicked:", clientSessionID, "in room:", roomID)
	return c.JSON(http.StatusOK, map[string]string{"message": "client kicked"})
}

// ルームから退出
func (mc *MainController) LeaveRoom(c echo.Context) error {
	roomID := c.Param("id")
	clientSessionID := c.QueryParam("client_session_id")
	err := mc.RoomUsecase.LeaveRoom(roomID, clientSessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("Client left:", clientSessionID, "in room:", roomID)
	return c.JSON(http.StatusOK, map[string]string{"message": "client left"})
}

// 認証状態を確認
func (mc *MainController) IsAuth(c echo.Context) error {
	roomID := c.Param("id")
	clientSessionID := c.QueryParam("client_session_id")
	isAuth, err := mc.RoomUsecase.IsAuth(roomID, clientSessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]bool{"isAuth": isAuth})
}
