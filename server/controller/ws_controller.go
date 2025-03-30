package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// WebSocketHandler handles WebSocket connections.
func (mc *MainController) WebSocketHandler(c echo.Context) error {
	roomID := c.QueryParam("room_id")
	clientName := c.QueryParam("client_name")
	sessionID := GetCookie(c, "session_id")
	if roomID == "" || clientName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "room_id and client_name are required"})
	}
	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "session_id is required"})
	}
	err := mc.RoomUsecase.HandleWebSocketConnection(c.Response(), c.Request(), roomID, clientName, sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	//ユーザーがWebSocketに接続したときにログを出力
	fmt.Printf("%s connected to WebSocket\n", clientName)
	return nil
}
