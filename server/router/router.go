package router

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/takaryo1010/OneTimeChat/server/controller"
)

func NewRouter(mc *controller.MainController) *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // フロントエンドのオリジン
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	e.GET("/ws", mc.WebSocketHandler)
	e.POST("/room", mc.CreateRoom)
	e.GET("/room/:id", mc.GetRoom)
	e.POST("/room/:id", mc.JoinRoom)
	// e.GET("/room/:id/participants", mc.GetParticipants)
	// e.PATCH("/room/:id/settings", mc.UpdateRoomSettings)
	// e.POST("/room/:id/auth", mc.Authenticate)
	// e.DELETE("/room/:id", mc.DeleteRoom)
	// e.DELETE("/room/:id/kick/:participant_id", mc.KickParticipant)

	return e
}
