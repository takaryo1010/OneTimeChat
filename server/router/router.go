package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/takaryo1010/OneTimeChat/server/controller"
)

func NewRouter(mc *controller.MainController) *echo.Echo {
	e := echo.New()

	// CORS設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.0.0:3000"}, // フロントエンドのオリジン
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	}))

	// WebSocketエンドポイント
	e.GET("/ws", mc.WebSocketHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
		fmt.Println("WebSocket connection requested")
		return next
	})

	// `/room` に関するエンドポイントをグループ化
	roomGroup := e.Group("/room")
	roomGroup.POST("", mc.CreateRoom)
	roomGroup.GET("/:id", mc.GetRoom)
	roomGroup.POST("/:id", mc.JoinRoom)
	roomGroup.POST("/:id/auth", mc.Authenticate)

	// 将来的に有効化するかもしれないエンドポイント
	// roomGroup.GET("/:id/participants", mc.GetParticipants)
	// roomGroup.PATCH("/:id/settings", mc.UpdateRoomSettings)
	// roomGroup.DELETE("/:id", mc.DeleteRoom)
	// roomGroup.DELETE("/:id/kick/:participant_id", mc.KickParticipant)

	return e
}
