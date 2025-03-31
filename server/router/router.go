package router

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/takaryo1010/OneTimeChat/server/controller"
)

func NewRouter(mc *controller.MainController) *echo.Echo {
	e := echo.New()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	
	clientURL := os.Getenv("CLIENT_URL")
	if clientURL == "" {
		log.Fatal("CLIENT_URL is not set in the environment variables")
	}

	// CORS設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.1.9:3000", clientURL}, // フロントエンドのオリジン
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	}))

	// WebSocketエンドポイント
	e.GET("/ws", mc.WebSocketHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
		return next
	})

	// `/room` に関するエンドポイントをグループ化
	roomGroup := e.Group("/room")
	roomGroup.POST("", mc.CreateRoom)
	roomGroup.GET("/:id", mc.GetRoom)
	roomGroup.POST("/:id", mc.JoinRoom)
	roomGroup.POST("/:id/auth", mc.Authenticate)

	roomGroup.GET("/:id/participants", mc.GetParticipants)
	roomGroup.PATCH("/:id/settings", mc.UpdateRoomSettings)
	roomGroup.DELETE("/:id", mc.DeleteRoom)
	roomGroup.DELETE("/:id/kick", mc.KickParticipant)
	roomGroup.DELETE("/:id/leave", mc.LeaveRoom)
	roomGroup.GET("/:id/isAuth", mc.IsAuth)

	return e
}
