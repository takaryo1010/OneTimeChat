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
	sessionID := c.QueryParam("session_id")
	err := mc.RoomUsecase.HandleWebSocketConnection(c.Response(), c.Request(), roomID, clientName, sessionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	//ユーザーがWebSocketに接続したときにログを出力
	fmt.Printf("%s connected to WebSocket\n", clientName)
	return nil
}
